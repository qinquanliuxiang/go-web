package service

import (
	"github.com/google/wire"
)

var ProviderService = wire.NewSet(
	NewUserSVC,
	NewRoleSVC,
	NewPolicySVC,
)
