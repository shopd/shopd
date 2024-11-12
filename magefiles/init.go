package main

import (
	"os"

	"github.com/mozey/logutil"
	"github.com/pkg/errors"
	"github.com/shopd/shopd/go/config"
	"github.com/shopd/shopd/go/fileutil"
)

var conf *config.Config

const paramAppDir = "APP_DIR"

func init() {
	logutil.SetupLogger(true)

	// APP_DIR must be set by the caller
	appDir := os.Getenv(paramAppDir)
	if appDir == "" {
		panic(errors.WithStack(ErrParamInvalid(paramAppDir, "")))
	}
	if !fileutil.PathExists(appDir) {
		panic(errors.WithStack(ErrParamInvalid(paramAppDir, appDir)))
	}

	conf = config.New()
}
