package middleware

import (
	"github.com/google/wire"
	"qqlx/store"
	"qqlx/store/rbac"
)

var ProviderMiddleware = wire.NewSet(
	wire.Bind(new(store.Authorizer), new(*rbac.Authentication)),
	rbac.NewAuthentication,
	NewAuthorization,
)
