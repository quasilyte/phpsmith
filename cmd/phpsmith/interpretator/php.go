package interpretator

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
)

func RunPHP(ctx context.Context, dir string) ([]byte, error) {
	var (
		outBuffer bytes.Buffer
		errBuffer bytes.Buffer
	)
	phpCmd := exec.CommandContext(ctx, "php", "-f", dir+"/main.php")
	phpCmd.Stdout, phpCmd.Stderr = &outBuffer, &errBuffer

	if err := phpCmd.Run(); err != nil {
		return nil, fmt.Errorf("on run php: stdOut: %s, stdErr: %s", outBuffer.String(), errBuffer.String())
	}

	return outBuffer.Bytes(), nil
}
