package reason

import "errors"

var (
	ErrParams            = errors.New("params error")
	ErrPermission        = errors.New("permission denied")
	ErrHeaderEmpty       = errors.New("auth in the request header is empty")
	ErrTokenMode         = errors.New("token mode error")
	ErrTokenInvalid      = errors.New("token is invalid")
	ErrHeaderMalformed   = errors.New("the auth format in the request header is incorrect")
	ErrLdapGroupNotFound = errors.New("ldap group not found")
	ErrLdapUserNotFound  = errors.New("ldap user not found")
	ErrRoleNotFound      = errors.New("role does not exist")
	ErrRoleHasUser       = errors.New("role has user")
	ErrRoleIsEmpty       = errors.New("role is empty")
	ErrRoleExists        = errors.New("role already exists")
	ErrUserNotFound      = errors.New("user does not exist")
	ErrUserIsDisable     = errors.New("user has been disabled")
	ErrUserExists        = errors.New("user already exists")
	ErrUserIsEmpty       = errors.New("user is empty")
	ErrEncryptPassword   = errors.New("failed to encrypt password")
	ErrAdminUserNotAllow = errors.New("admin user cannot operate")
	ErrInvalidPassword   = errors.New("password is invalid")
	ErrPolicyNotFound    = errors.New("policy does not exist")
	ErrPolicyUsedByRole  = errors.New("policy has been used by role")
	ErrNameInvalid       = errors.New("name must contain only letters")
)
