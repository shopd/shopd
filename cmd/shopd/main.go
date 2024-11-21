package main

import (
	"fmt"
	"os"

	"github.com/mozey/logutil"
	shopdCmd "github.com/shopd/shopd/cmd/shopd/cmd"
)

// configBase64 is set at compile time with ldflags,
// if it's empty then read config from env at runtime
// TODO Investigate multi-tenant domain support with single shopd process.
// Supervisor process runs `shopdCmd.Execute(configBase64)`
// in separate go-routines.
// The supervisor process has a map of configBase64 for each domain
var configBase64 string

func main() {
	defer logutil.PanicHandler()
	logutil.SetupLogger(true)

	paramAppDir := "APP_DIR"
	if os.Getenv(paramAppDir) == "" {
		// Must be set on the env.
		// APP_DIR is not encoded in configBase64
		panic(fmt.Sprintf("%s is empty", paramAppDir))
	}

	// "the main.go file is very bare... one purpose: initializing Cobra"
	// https://github.com/spf13/cobra/blob/main/user_guide.md
	shopdCmd.Execute(configBase64)
}
