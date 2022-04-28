package interpretator

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"sync"
)

var mu sync.Mutex

func RunKPHP(ctx context.Context, dir string) ([]byte, error) {
	var (
		outBuffer bytes.Buffer
		errBuffer bytes.Buffer
	)

	mu.Lock()
	defer mu.Unlock()
	binaryName := "./kphp_out/" + dir
	compileCmd := exec.CommandContext(ctx, "kphp", "--mode", "cli", "-o", binaryName, dir+"/main.php")
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
