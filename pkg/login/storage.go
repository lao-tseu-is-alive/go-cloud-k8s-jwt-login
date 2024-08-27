package login

import (
	"fmt"
	"github.com/lao-tseu-is-alive/go-cloud-k8s-common/pkg/database"
	"github.com/lao-tseu-is-alive/go-cloud-k8s-common/pkg/golog"
)

type Storage interface {
	//CreateOrUpdateAdminUserOrPanic creates or updates the admin user in the store
	CreateOrUpdateAdminUserOrPanic()
	// Get returns the user with the specified user login.
	Get(login string) (*User, error)
	// Exist returns true only if a user with the specified login exists in store.
	Exist(login string) bool
	// IsUserActive returns true if the user with the specified login has the is_active attribute set to true
	IsUserActive(login string) bool
	// IsAdmin returns true if the user with the specified login has the is_admin attribute set to true
	IsAdmin(login string) bool
	// IsLocked returns true if the user with the specified login has the is_locked attribute set to true
	IsLocked(login string) bool
}

func GetStorageInstanceOrPanic(dbDriver string, db database.DB, l golog.MyLogger) Storage {
	var store Storage
	var err error
	switch dbDriver {
	case "pgx":
		store, err = NewPgxDB(db, l)
		if err != nil {
			panic(fmt.Sprintf("error doing NewPgxDB(pgConn : %v", err))
		}

	default:
		panic("unsupported DB driver type")
	}
	return store
}
