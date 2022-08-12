package php

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"os/exec"
)

type Runner struct{}

func (Runner) Run(ctx context.Context, dir string, seed int64) ([]byte, error) {
	var (
		outBuffer bytes.Buffer
		errBuffer bytes.Buffer
	)
	phpCmd := exec.CommandContext(ctx, "php", "-f", dir+"/main.php")
	phpCmd.Stdout, phpCmd.Stderr = &outBuffer, &errBuffer

	if err := phpCmd.Run(); err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok && exitErr.ExitCode() > 1 {
			log.Printf("non one or zero exit code found: %d, seed: %d \n", exitErr.ExitCode(), seed)
		}
		return nil, fmt.Errorf("on run php: stdOut: %s, stdErr: %s", outBuffer.String(), errBuffer.String())
	}

	return outBuffer.Bytes(), nil
}

func (Runner) Name() string {
	return "php_runner"
}
