package templrendr

import (
	"github.com/mozey/errors"
)

var ErrTemplRendr = errors.NewCause("templrendr")

var ErrRouteNotFound = func(route string) error {
	return errors.NewWithCausef(ErrTemplRendr, "route not found %s", route)
}

var ErrMethodNotSupported = func(method, route string) error {
	return errors.NewWithCausef(ErrTemplRendr,
		"method %s not supported for route %s", method, route)
}
