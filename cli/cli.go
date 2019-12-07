package cli

import (
	"context"
	"os"
	"os/exec"
	"syscall"
	"time"

	"github.com/op/go-logging"
	"github.com/urfave/cli"

	"github.com/ziyan/panto/utils"
)

var log = logging.MustGetLogger("cli")

func Run(args []string) {

	app := cli.NewApp()
	app.EnableBashCompletion = true
	app.Name = "panto"
	app.Version = "0.1.0"
	app.Usage = "A tool to collate system log and events."

	app.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:  "log-level",
			Value: "INFO",
			Usage: "log level",
		},
		&cli.StringFlag{
			Name:  "log-format",
			Value: "%{color}%{time:2006-01-02T15:04:05.000Z07:00} [%{level:.4s}] [%{shortfile} %{shortfunc}] %{message}%{color:reset}",
			Usage: "log format",
		},
		&cli.StringFlag{
			Name:  "remote-executable-path",
			Value: "/tmp/panto",
		},
		&cli.StringSliceFlag{
			Name:  "remote",
			Usage: "remote hostname",
		},
	}

	app.Before = func(c *cli.Context) error {
		logging.SetBackend(logging.NewBackendFormatter(
			logging.NewLogBackend(os.Stderr, "", 0),
			logging.MustStringFormatter(c.String("log-format")),
		))
		if level, err := logging.LogLevel(c.String("log-level")); err == nil {
			logging.SetLevel(level, "")
		}
		log.Debugf("log level set to %s", logging.GetLevel(""))
		log.Debugf("hostname is %s", utils.Hostname)
		return nil
	}

	app.Commands = []*cli.Command{
		&cli.Command{
			Name: "server",
			Action: func(c *cli.Context) error {
				time.Sleep(60 * time.Second)
				return nil
			},
		},
	}

	app.Action = func(c *cli.Context) error {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		for _, remote := range c.StringSlice("remote") {
			if err := utils.DeployExecutable(ctx, remote, c.String("remote-executable-path")); err != nil {
				return err
			}
		}

		var cmds []*exec.Cmd
		defer func() {
			cancel()
			for _, cmd := range cmds {
				if err := cmd.Wait(); err != nil {
					log.Errorf("command %s exited with: %s", cmd, err)
				}
			}
		}()

		for _, remote := range c.StringSlice("remote") {
			cmd, _, err := utils.RunRemoteExecutable(ctx, remote, c.String("remote-executable-path"), "--log-level", c.String("log-level"), "server")
			if err != nil {
				return err
			}
			cmds = append(cmds, cmd)
		}

		time.Sleep(2 * time.Second)

		for _, cmd := range cmds {
			if cmd.Process != nil {
				if err := cmd.Process.Signal(syscall.SIGTERM); err != nil {
					log.Errorf("failed to signal command %s: %s", cmd, err)
				}
			}
		}

		time.Sleep(1 * time.Second)
		return nil
	}

	app.Run(args)
}
