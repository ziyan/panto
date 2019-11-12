package cli

import (
	"os"

	"github.com/op/go-logging"
	"github.com/urfave/cli"
)

var log = logging.MustGetLogger("cli")

func configureLogging(level, format string) {
	logging.SetBackend(logging.NewBackendFormatter(
		logging.NewLogBackend(os.Stderr, "", 0),
		logging.MustStringFormatter(format),
	))
	if level, err := logging.LogLevel(level); err == nil {
		logging.SetLevel(level, "")
	}
	log.Debugf("log level set to %s", logging.GetLevel(""))
}

func Run(args []string) {

	app := cli.NewApp()
	app.EnableBashCompletion = true
	app.Name = "panto"
	app.Version = "0.1.0"
	app.Usage = "A tool to collate system log and events."

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "log-level",
			Value: "INFO",
			Usage: "log level",
		},
		cli.StringFlag{
			Name:  "log-format",
			Value: "%{color}%{time:2006-01-02T15:04:05.000Z07:00} [%{level:.4s}] [%{shortfile} %{shortfunc}] %{message}%{color:reset}",
			Usage: "log format",
		},
	}

	app.Action = func(c *cli.Context) error {
		configureLogging(c.String("log-level"), c.String("log-format"))

		return nil
	}

	app.Run(args)
}
