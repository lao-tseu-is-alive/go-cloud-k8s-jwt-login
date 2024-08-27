package login

import (
	"github.com/lao-tseu-is-alive/go-cloud-k8s-common/pkg/config"
	"github.com/lao-tseu-is-alive/go-cloud-k8s-jwt-login/pkg/crypto"
)

const (
	insertAdminUser = `
INSERT INTO go_user (name, username, password_hash,email, external_id, is_admin, creator, comment) 
VALUES ('Administrative Account',$1,$2,$3,$4, true, 1,  'Initial setup of Admin account')  RETURNING id;`

	updateAdminUser = `
UPDATE go_user
SET external_id   		   = $1,
    password_hash		 = $2,
    is_locked              = false,
    is_admin               = true,
    last_modification_time = CURRENT_TIMESTAMP,
    last_modification_user = 1,
    is_active              = true, 
	bad_password_count	   = 0  	-- we decide to reset counter 
WHERE username = $3;
`
)

// CreateOrUpdateAdminUserOrPanic is a function to create or update the admin user in the database
func (db *PGX) CreateOrUpdateAdminUserOrPanic() {
	db.log.Debug("trace : entering CreateOrUpdateAdminUserOrPanic()")
	adminUser := config.GetAdminUserFromFromEnvOrPanic("admin")
	adminPassword := config.GetAdminPasswordFromFromEnvOrPanic()
	adminEmail := config.GetAdminEmailFromFromEnvOrPanic("admin@example.com")
	adminExternalId := config.GetAdminIdFromFromEnvOrPanic(999999)
	passwordHash := crypto.Sha256Hash(adminPassword)
	goHash, err := crypto.HashAndSalt(passwordHash)
	if err != nil {
		db.log.Error("crypto.HashAndSalt unexpectedly failed. error : %v", err)
		panic("unable to calculate hash for the admin password ")
	}
	// check if the admin user already exists
	if db.Exist(adminUser) {
		db.log.Info("Admin user already exists, will just update password")
		_, err := db.dbi.ExecActionQuery(updateAdminUser, adminExternalId, goHash, adminUser)
		if err != nil {
			db.log.Error("CreateOrUpdateAdminUserOrPanic() could not be updated in DB. failed db.Query err: %v", err)
			panic("unable to update the admin user password")
		}
		return
	}
	// create the admin user
	newId, err := db.dbi.ExecActionQuery(insertAdminUser, adminUser, goHash, adminEmail, adminExternalId)
	if err != nil {
		db.log.Error("CreateOrUpdateAdminUserOrPanic() could not be created in DB. failed db.Query err: %v", err)
		panic("unable to create the admin user")
	}
	db.log.Info("Admin user created with id:%v", newId)
	return
}
