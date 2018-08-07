package ctl

import (
	"context"
)

//Base struct to all okectl operations
type Base struct {
	CompartmentID string
	Context       context.Context
}
