package login

const (
	countUsers   = "SELECT COUNT(*) FROM go_user;"
	existUser    = "SELECT COUNT(*) FROM go_user WHERE username = $1;"
	isActiveUser = "SELECT COUNT(*) FROM go_user WHERE is_active=true AND username = $1;"
	isLockedUser = "SELECT COUNT(*) FROM go_user WHERE is_locked=true AND username = $1;"
	isAdminUser  = "SELECT COUNT(*) FROM go_user WHERE is_admin=true AND username = $1;"
	getUser      = `select id,
       name,
       email,
       username,
       password_hash,
       external_id,
       is_locked,
       is_admin,
       create_time,
       creator,
       last_modification_time,
       last_modification_user,
       is_active,
       inactivation_time,
       inactivation_reason,
       comment,
       bad_password_count
from go_user
where username = $1;
`
)
