package reason

import "errors"

var (
	ErrParams            = errors.New("params error")
	ErrPermission        = errors.New("permission denied")
	ErrRoleEmpty         = errors.New("role is empty")
	ErrHeaderEmpty       = errors.New("auth in the request header is empty")
	ErrTokenMode         = errors.New("token mode error")
	ErrTokenInvalid      = errors.New("token is invalid")
	ErrHeaderMalformed   = errors.New("the auth format in the request header is incorrect")
	ErrLdapGroupNotFound = errors.New("ldap group not found")
	ErrLdapUserNotFound  = errors.New("ldap user not found")
	ErrRoleNotFound      = errors.New("role does not exist")
	ErrRoleExists        = errors.New("role already exists")
	ErrUserHasSameRole   = errors.New("user has the same role")
	ErrUserNotFound      = errors.New("user does not exist")
	ErrInvalidPassword   = errors.New("password is invalid")
	ErrPolicyNotFound    = errors.New("policy does not exist")
	ErrGetIDFailed       = errors.New("id is invalid")
	ErrAdminRole         = errors.New("admin role can not modify policy")
	ErrNameInvalid       = errors.New("name must contain only letters")
)
