package cli

import (
	"fmt"

	"github.com/suzuki-shunsuke/circleci-config-merge/pkg/controller"
	"github.com/urfave/cli/v2"
)

func (runner *Runner) setCLIArg(c *cli.Context, params controller.Params) controller.Params {
	arr := c.Args().Slice()
	args := make(map[string]struct{}, len(arr))
	for _, a := range arr {
		args[a] = struct{}{}
	}
	params.Files = args
	return params
}

func (runner *Runner) action(c *cli.Context) error {
	// files...
	params := controller.Params{}
	params = runner.setCLIArg(c, params)

	ctrl, params, err := controller.New(c.Context, params)
	if err != nil {
		return fmt.Errorf("initialize a controller: %w", err)
	}

	return ctrl.Run(c.Context, &params) //nolint:wrapcheck
}
