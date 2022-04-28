package interpretator

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"os/exec"
	"sync"
)

var mu sync.Mutex

func RunKPHP(ctx context.Context, dir string, seed int64) ([]byte, error) {
	var (
		outBuffer bytes.Buffer
		errBuffer bytes.Buffer
	)

	binaryName := dir + "/" + dir
	mu.Lock()
	defer mu.Unlock()
	compileCmd := exec.CommandContext(ctx, "kphp", "--mode", "cli", "-o", binaryName, dir+"/main.php")
	compileCmd.Stdout, compileCmd.Stderr = &outBuffer, &errBuffer

	if err := compileCmd.Run(); err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok && exitErr.ExitCode() > 1 {
			log.Printf("non one or zero exit code found: %d, seed: %d \n", exitErr.ExitCode(), seed)
		}
		return nil, fmt.Errorf("on compile kphp: stdOut: %s, stdErr: %s", outBuffer.String(), errBuffer.String())
	}

	outBuffer.Reset()
	errBuffer.Reset()

	runCmd := exec.CommandContext(ctx, binaryName)
	runCmd.Stdout, runCmd.Stderr = &outBuffer, &errBuffer

	if err := runCmd.Run(); err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok && exitErr.ExitCode() > 1 {
			log.Printf("non one or zero exit code found: %d, seed: %d \n", exitErr.ExitCode(), seed)
		}
		return nil, fmt.Errorf("on run kphp binary: stdOut: %s, stdErr: %s", outBuffer.String(), errBuffer.String())
	}

	return outBuffer.Bytes(), nil
}
