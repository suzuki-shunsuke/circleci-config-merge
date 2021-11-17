package controller

import (
	"context"
	"fmt"
	"os"

	"gopkg.in/yaml.v2"
)

type Params struct {
	Files map[string]struct{}
}

func New(ctx context.Context, params Params) (Controller, Params, error) {
	return Controller{
		Stdout: os.Stdout,
		Stderr: os.Stderr,
	}, params, nil
}

func (ctrl *Controller) Run(ctx context.Context, params *Params) error {
	cfg := Config{}
	for filePath := range params.Files {
		child := Config{}
		if err := readFile(filePath, &child); err != nil {
			return fmt.Errorf("read a file "+filePath+": %w", err)
		}
		cfg = mergeConfig(cfg, child)
	}
	for k, workflow := range cfg.Workflows.Workflows {
		jobs, err := sortJobs(workflow.Jobs)
		if err == nil {
			workflow.Jobs = jobs
		}
		cfg.Workflows.Workflows[k] = workflow
	}
	if err := yaml.NewEncoder(ctrl.Stdout).Encode(&cfg); err != nil {
		return fmt.Errorf("encode a merged config as YAML: %w", err)
	}
	return nil
}
