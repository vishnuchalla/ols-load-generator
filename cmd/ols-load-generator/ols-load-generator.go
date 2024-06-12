package main

import (
	"context"
	"os"
	"time"

	"github.com/quay/zlog"
	"github.com/rs/zerolog"
	"github.com/urfave/cli/v2"
)

var logout zerolog.Logger
var version string

// createLogger returns a new logger with the specified log level and time format.
func createLogger(level zerolog.Level, timeFormat string) zerolog.Logger {
	return zerolog.New(&zerolog.ConsoleWriter{
		Out:        os.Stderr,
		TimeFormat: timeFormat,
	}).Level(level).
		With().
		Timestamp().
		Logger()
}

// setLogLevel sets the log level based on the value of the "-D" and "-W" flags.
// It returns an error if any during the execution.
func setLogLevel(c *cli.Context) error {
	level := zerolog.InfoLevel
	if c.Bool("W") {
		level = zerolog.WarnLevel
	}
	if c.Bool("D") {
		level = zerolog.DebugLevel
	}
	logout = createLogger(level, time.RFC3339)
	zlog.Set(&logout)
	return nil
}

// main drives the execution for ols-load-generator.
func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	app := &cli.App{
		Name:                 "ols-load-generator",
		Version:              version,
		Usage:                "A command-line tool to load test openshift lightspeed service (OLS).",
		Description:          "A command-line tool to load test openshift lightspeed service (OLS).",
		EnableBashCompletion: true,
		Before:               setLogLevel,
		Commands: []*cli.Command{
			AttackCmd,
		},
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:  "D",
				Usage: "print debugging logs",
			},
			&cli.BoolFlag{
				Name:  "W",
				Usage: "quieter log output",
			},
		},
		ExitErrHandler: func(c *cli.Context, err error) {
			if err != nil {
				if err, ok := err.(cli.ExitCoder); ok {
					err.ExitCode()
				}
				logout.Error().Err(err).Send()
			}
		},
	}
	if err := app.RunContext(ctx, os.Args); err != nil {
		logout.Error().Err(err).Send()
	}
}
