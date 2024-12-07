package router

import (
	"github.com/mozey/errors"
)

var ErrRouter = errors.NewCause("router")

var ErrRouteNotFound = func(route string) error {
	return errors.NewWithCausef(ErrRouter, "route not found %s", route)
}
