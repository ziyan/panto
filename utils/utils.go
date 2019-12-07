package utils

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"

	"github.com/op/go-logging"
)

var log = logging.MustGetLogger("utils")

var Hostname = func() string {
	hostname, err := os.Hostname()
	if err != nil {
		panic(err)
	}
	return hostname
}()

var Executable = func() string {
	executable, err := os.Executable()
	if err != nil {
		panic(err)
	}
	return executable
}()

func DeployExecutable(ctx context.Context, remote, path string) error {
	return exec.CommandContext(ctx, "rsync", "--perms", "--times", "--compress", "--checksum", Executable, fmt.Sprintf("%s:%s", remote, path)).Run()
}

func RunRemoteExecutable(ctx context.Context, remote, path string, args ...string) (*exec.Cmd, io.ReadWriteCloser, error) {
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
