package interpretator

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
)

func RunKPHP(ctx context.Context, filename string) ([]byte, error) {
	var (
		outBuffer bytes.Buffer
		errBuffer bytes.Buffer
	)

	binaryName := "./kphp_out/" + filename
	compileCmd := exec.CommandContext(ctx, "kphp", "--mode", "cli", "-o", binaryName, filename)
	compileCmd.Stdout, compileCmd.Stderr = &outBuffer, &errBuffer

	if err := compileCmd.Run(); err != nil {
		return nil, fmt.Errorf("on compile kphp: stdOut: %s, stdErr: %s", outBuffer.String(), errBuffer.String())
	}

	outBuffer.Reset()
	errBuffer.Reset()

	runCmd := exec.CommandContext(ctx, binaryName)
	runCmd.Stdout, runCmd.Stderr = &outBuffer, &errBuffer

	if err := runCmd.Run(); err != nil {
		return nil, fmt.Errorf("on run kphp binary: stdOut: %s, stdErr: %s", outBuffer.String(), errBuffer.String())
	}

	return outBuffer.Bytes(), nil
}
