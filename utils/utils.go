package utils

import (
	"os"

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
