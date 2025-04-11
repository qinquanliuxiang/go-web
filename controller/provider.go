package controller

import "github.com/google/wire"

var ProviderContr = wire.NewSet(
	NewUserCtrl,
	NewRoleCtrl,
	NewPolicyCtrl,
)
