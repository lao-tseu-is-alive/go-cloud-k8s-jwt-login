package login

import (
	"context"
	"errors"
	"fmt"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/lao-tseu-is-alive/go-cloud-k8s-common/pkg/database"
	"github.com/lao-tseu-is-alive/go-cloud-k8s-common/pkg/golog"
)

type PGX struct {
	Conn *pgxpool.Pool
	dbi  database.DB
	log  golog.MyLogger
}

func (db *PGX) Get(login string) (*User, error) {
	db.log.Debug("trace : entering Get(%v)", login)
	if !db.Exist(login) {
		msg := fmt.Sprintf(UserDoesNotExist, login)
		db.log.Warn(msg)
		return nil, errors.New(msg)
	}
	res := &User{}
	err := pgxscan.Get(context.Background(), db.Conn, res, getUser, login)
	if err != nil {
		db.log.Error(SelectFailedInNWithErrorE, "Get", err)
		return nil, err
	}
	if res == nil {
		db.log.Info(FunctionNReturnedNoResults, "Get")
		return nil, pgx.ErrNoRows
	}
	return res, nil
}

func (db *PGX) Exist(login string) bool {
	db.log.Debug("trace : entering Exist(%v)", login)
	count, err := db.dbi.GetQueryInt(existUser, login)
	if err != nil {
		db.log.Error("Exist(%v) could not be retrieved from DB. failed db.Query err: %v", login, err)
		return false
	}
	if count > 0 {
		db.log.Info(" Exist(%v) id does exist  count:%v", login, count)
		return true
	} else {
		db.log.Info(" Exist(%v) id does not exist count:%v", login, count)
		return false
	}
}

func (db *PGX) IsUserActive(login string) bool {
	db.log.Debug("trace : entering IsUserActive(%s)", login)
	count, err := db.dbi.GetQueryInt(isActiveUser, login)
	if err != nil {
		db.log.Error("IsUserActive(%s) could not be retrieved from DB. failed db.Query err: %v", login, err)
		return false
	}
	if count > 0 {
		db.log.Info(" IsUserActive(%s) is true  count:%v", login, count)
		return true
	} else {
		db.log.Info(" IsUserActive(%s) is false count:%v", login, count)
		return false
	}
}

func (db *PGX) IsAdmin(login string) bool {
	db.log.Debug("trace : entering IsAdmin(%s)", login)
	count, err := db.dbi.GetQueryInt(isAdminUser, login)
	if err != nil {
		db.log.Error("IsAdmin(%s) could not be retrieved from DB. failed db.Query err: %v", login, err)
		return false
	}
	if count > 0 {
		db.log.Info(" IsAdmin(%s) is true  count:%v", login, count)
		return true
	} else {
		db.log.Info(" IsAdmin(%s) is false count:%v", login, count)
		return false
	}
}

func (db *PGX) IsLocked(login string) bool {
	db.log.Debug("trace : entering IsLocked(%s)", login)
	count, err := db.dbi.GetQueryInt(isLockedUser, login)
	if err != nil {
		db.log.Error("IsLocked(%s) could not be retrieved from DB. failed db.Query err: %v", login, err)
		return false
	}
	if count > 0 {
		db.log.Info(" IsLocked(%s) is true  count:%v", login, count)
		return true
	} else {
		db.log.Info(" IsLocked(%s) is false count:%v", login, count)
		return false
	}
}

// NewPgxDB will instantiate a new storage of type postgres and ensure schema exist
func NewPgxDB(db database.DB, log golog.MyLogger) (Storage, error) {
	var psql PGX
	pgConn, err := db.GetPGConn()
	if err != nil {
		return nil, err
	}
	psql.Conn = pgConn
	psql.dbi = db
	psql.log = log
	var numberOfRows int
	errTypeThingTable := pgConn.QueryRow(context.Background(), countUsers).Scan(&numberOfRows)
	if errTypeThingTable != nil {
		log.Error("Unable to retrieve the number of users error: %v", err)
		return nil, err
	}

	if numberOfRows > 0 {
		log.Info("'database contains %d records'", numberOfRows)
	} else {
		log.Warn("«go_thing.type_thing» is empty ! it should contain at least one row")
		return nil, errors.New("problem with initial content of database «go_thing.type_thing» should not be empty ")
	}

	return &psql, err
}
