package middleware

import (
	"qqlx/base/interfaces"
	"qqlx/store/rbac"

	"github.com/google/wire"
)

var ProviderMiddleware = wire.NewSet(
	wire.Bind(new(interfaces.Authorizer), new(*rbac.Authentication)),
	rbac.NewAuthentication,
	NewAuthorization,
)
