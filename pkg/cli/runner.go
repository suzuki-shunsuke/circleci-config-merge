package cli

import (
	"context"
	"io"

	"github.com/suzuki-shunsuke/circleci-config-merge/pkg/constant"
	"github.com/urfave/cli/v2"
)

type Runner struct {
	Stdin  io.Reader
	Stdout io.Writer
	Stderr io.Writer
}

func (runner Runner) Run(ctx context.Context, args ...string) error {
	app := cli.App{
		Name:    "circleci-config-merge",
		Usage:   "generate CircleCI configuration file by merging multiple files. https://github.com/suzuki-shunsuke/circleci-config-merge",
		Version: constant.Version,
		Commands: []*cli.Command{
			{
				Name:   "merge",
				Usage:  "generate CircleCI configuration file by merging multiple files",
				Action: runner.action,
				Flags:  []cli.Flag{},
			},
		},
	}

	return app.RunContext(ctx, args)
}
