package login

import "time"

type User struct {
	Id                   int32      `json:"id"`
	Name                 string     `json:"name"`
	Email                string     `json:"email"`
	Username             string     `json:"username"`
	PasswordHash         string     `json:"password_hash"`
	ExternalId           *int32     `json:"external_id,omitempty"`
	IsLocked             bool       `json:"is_locked"`
	IsAdmin              bool       `json:"is_admin"`
	CreateTime           time.Time  `json:"create_time"`
	Creator              int32      `json:"creator"`
	LastModificationTime *time.Time `json:"last_modification_time,omitempty"`
	LastModificationUser *int32     `json:"last_modification_user,omitempty"`
	IsActive             bool       `json:"is_active"`
	InactivationTime     *time.Time `json:"inactivation_time,omitempty"`
	InactivationReason   *string    `json:"inactivation_reason,omitempty"`
	Comment              *string    `json:"comment,omitempty"`
	BadPasswordCount     int32      `json:"bad_password_count"`
}

type JwtInfo struct {
	Authenticated bool   `json:"authenticated"`
	Method        string `json:"method"`
	Message       string `json:"message"`
	Token         string `json:"token"`
	Login         string `json:"login"`
	Email         string `json:"email"`
	UserId        int    `json:"user_id"`
}
