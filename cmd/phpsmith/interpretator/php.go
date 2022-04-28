package interpretator

import (
	"context"
	"os/exec"
)

func RunPHP(ctx context.Context, filename string) ([]byte, error) {
	return exec.CommandContext(ctx, "php", "-f", filename).CombinedOutput()
}
