package testutil

import (
	"testing"

	"github.com/matryer/is"
	"github.com/mozey/logutil"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"github.com/shopd/shopd/go/config"
)

// I wraps is.I
type I struct {
	*is.I
	t *testing.T
}

// NoErr logs the error with stack trace if err is not nil and exits
func (is *I) NoErr(err error) {
	if err != nil {
		type stackTracer interface {
			StackTrace() errors.StackTrace
		}
		_, ok := err.(stackTracer)
		if !ok {
			// Add stack trace to err if it doesn't have one
			// See comments in github.com/pkg/errors
			// "Although the stackTracer interface is not exported by this package,
			// it is considered a part of its stable public interface"
			err = errors.WithStack(err)
		}
		log.Error().Stack().Err(err).Msg("")
	}
	// TODO Behaviour here differs for is.NewRelaxed,
	// currently this wrapper doesn't support relaxed mode
	is.I.NoErr(err)
}

func New(t *testing.T, is *is.I) *I {
	return &I{t: t, I: is}
}

func Setup(t *testing.T) *I {
	logutil.SetupLogger(true)
	return New(t, is.New(t))
}

func SetupConf(t *testing.T) (wrapper *I, conf *config.Config) {
	logutil.SetupLogger(true)
	wrapper = New(t, is.New(t))
	// TODO LoadFile support for .env? Tests must always use dev config
	// conf, err := config.LoadFile("dev")
	conf = config.New()
	return wrapper, conf
}
