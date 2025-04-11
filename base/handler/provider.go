package handler

import (
	"github.com/google/wire"
)

var ProviderHandler = wire.NewSet(
	wire.Bind(new(BindResponseInterface), new(*BindRequest)),
	NewResponse,
)
