package main

import (
	"embed"
	_ "github.com/golang-migrate/migrate/v4/database/pgx/v5"
	"github.com/lao-tseu-is-alive/go-cloud-k8s-common/pkg/config"
	"github.com/lao-tseu-is-alive/go-cloud-k8s-common/pkg/database"
	"github.com/lao-tseu-is-alive/go-cloud-k8s-common/pkg/gohttp"
	"github.com/lao-tseu-is-alive/go-cloud-k8s-common/pkg/golog"
	"github.com/lao-tseu-is-alive/go-cloud-k8s-common/pkg/tools"
	"github.com/lao-tseu-is-alive/go-cloud-k8s-jwt-login/pkg/version"
	"io/fs"
	"log"
	"net/http"
	"runtime"
)

const (
	APP               = "goCloudK8sJwtLoginServer"
	defaultPort       = 8888
	defaultServerIp   = "0.0.0.0"
	defaultServerPath = "/"
	defaultWebRootDir = "front/dist/"
	defaultAdminId    = 98765
	defaultAdminUser  = "goadmin"
	defaultAdminEmail = "goadmin@lausanne.ch"
)

// content holds our static web server content.
//
//go:embed all:front/dist
var content embed.FS

// sqlMigrations holds our db migrations sql files using https://github.com/golang-migrate/migrate
// in the line above you SHOULD have the same path  as const defaultSqlDbMigrationsPath
//
//go:embed db/migrations/*.sql
var sqlMigrations embed.FS

func GetMyDefaultHandler(s *gohttp.Server, webRootDir string, content embed.FS) http.HandlerFunc {
	handlerName := "GetMyDefaultHandler"
	logger := s.GetLog()
	logger.Debug("Initial call to %s with webRootDir:%s", handlerName, webRootDir)
	RootPathGetCounter := s.RootPathGetCounter
	// Create a subfolder filesystem to serve only the content of front/dist
	subFS, err := fs.Sub(content, "front/dist")
	if err != nil {
		logger.Fatal("Error creating sub-filesystem: %v", err)
	}

	// Create a file server handler for the embed filesystem
	handler := http.FileServer(http.FS(subFS))

	return func(w http.ResponseWriter, r *http.Request) {
		gohttp.TraceRequest(handlerName, r, logger)
		RootPathGetCounter.Inc()
		handler.ServeHTTP(w, r)
	}
}

func GetProtectedHandler(server *gohttp.Server, l golog.MyLogger) http.HandlerFunc {
	handlerName := "GetProtectedHandler"
	return func(w http.ResponseWriter, r *http.Request) {
		gohttp.TraceRequest(handlerName, r, l)
		// get the user from the context
		claims := gohttp.GetJwtCustomClaimsFromContext(r)

		currentUserId := claims.User.UserId
		// check if user is admin
		if !claims.User.IsAdmin {
			l.Error("User %d is not admin: %+v", currentUserId, claims)
			http.Error(w, "User is not admin", http.StatusForbidden)
			return
		}
		// respond with protected data
		err := server.JsonResponse(w, claims)
		if err != nil {
			http.Error(w, "Error responding with protected data", http.StatusInternalServerError)
			return
		}
	}
}

func main() {

	l, err := golog.NewLogger("zap", golog.DebugLevel, APP)
	if err != nil {
		log.Fatalf("💥💥 error golog.NewLogger error: %v'\n", err)
	}
	l.Info("🚀🚀 Starting App:'%s', ver:%s, build:%s, from: %s", APP, version.VERSION, version.Build, version.REPOSITORY)

	dbDsn := config.GetPgDbDsnUrlFromEnvOrPanic("127.0.0.1", 5432,
		tools.ToSnakeCase(version.APP), version.AppSnake, "prefer")

	db, err := database.GetInstance("pgx", dbDsn, runtime.NumCPU(), l)
	if err != nil {
		l.Fatal("💥💥 error doing database.GetInstance(postgres, dbDsn  : %v\n", err)
	}
	defer db.Close()

	// checking database connection
	dbVersion, err := db.GetVersion()
	if err != nil {
		l.Fatal("💥💥 error doing dbConn.GetVersion() error: %v", err)
	}
	l.Info("connected to db version : %s", dbVersion)
	// checking metadata information
	//metadataService := metadata.Service{		Log: l,		Db:  db,	}

	myVersionReader := gohttp.NewSimpleVersionReader(APP, version.VERSION, version.REPOSITORY, version.Build)
	server := gohttp.CreateNewServerFromEnvOrFail(
		defaultPort,
		defaultServerIp,
		defaultAdminUser,
		defaultAdminEmail,
		defaultAdminId,
		myVersionReader,
		l)

	mux := server.GetRouter()
	myJwt := server.JwtCheck
	mux.Handle("POST /login", gohttp.GetLoginPostHandler(server))
	// Protected endpoint (using jwtMiddleware)
	mux.Handle("GET /protected", myJwt.JwtMiddleware(GetProtectedHandler(server, l)))
	mux.Handle("GET /*", gohttp.NewPrometheusMiddleware(
		server.GetPrometheusRegistry(), nil).
		WrapHandler("GET /*", GetMyDefaultHandler(server, defaultWebRootDir, content)),
	)

	server.StartServer()
}
