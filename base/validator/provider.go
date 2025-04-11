package validator

import "github.com/google/wire"

var ProviderValidator = wire.NewSet(
	wire.Bind(new(CheckReqInterface), new(*Validator)),
	NewValidator,
)
