package utils

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
)

func RunRemoteExecutable(ctx context.Context, remote, path string, args ...string) (*exec.Cmd, io.ReadWriteCloser, error) {

	if err := exec.CommandContext(ctx, "rsync", "--perms", "--times", "--compress", "--checksum", Executable, fmt.Sprintf("%s:%s", remote, path)).Run(); err != nil {
		return nil, nil, err
	}

	cmd := exec.CommandContext(ctx, "ssh", append([]string{remote, path}, args...)...)
	cmd.Stderr = os.Stderr

	wc, err := cmd.StdinPipe()
	if err != nil {
		return nil, nil, err
	}
	defer func() {
		if wc != nil {
			wc.Close()
		}
	}()

	rc, err := cmd.StdoutPipe()
	if err != nil {
		return nil, nil, err
	}
	defer func() {
		if rc != nil {
			rc.Close()
		}
	}()

	if err := cmd.Start(); err != nil {
		return nil, nil, err
	}

	rwc := NewReadWriteCloser(rc, wc)
	rc = nil
	wc = nil
	return cmd, rwc, nil
}
