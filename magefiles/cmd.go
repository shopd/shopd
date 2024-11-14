package main

import (
	"fmt"
	"os/exec"

	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

func combinedOutput(cmd *exec.Cmd) (out []byte, err error) {
	log.Info().Str("cmd", cmd.String()).Msg("Exec")
	out, err = cmd.CombinedOutput()
	if err != nil {
		return out, errors.WithStack(err)
	}
	return out, nil
}

func printCombinedOutput(cmd *exec.Cmd) (err error) {
	out, err := combinedOutput(cmd)
	if err != nil {
		fmt.Println(string(out))
		return err
	}
	if len(out) > 0 {
		fmt.Println(string(out))
	}
	return nil
}
