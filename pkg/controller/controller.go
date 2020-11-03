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

type Workflows struct {
	Version   interface{}         `yaml:",omitempty"`
	Workflows map[string]Workflow `yaml:",omitempty"`
}

func (wfs *Workflows) UnmarshalYAML(unmarshal func(interface{}) error) error {
	m := map[string]interface{}{}
	if err := unmarshal(&m); err != nil {
		return err
	}
	wfMap := map[string]Workflow{}
	for workflowName, v := range m {
		if workflowName == "version" {
			ver, ok := v.(string)
			if !ok {
				return fmt.Errorf("workflow version must be string: %+v", ver)
			}
			wfs.Version = ver
			continue
		}
		wf := Workflow{}
		a, ok := v.(map[interface{}]interface{})
		if !ok {
			return fmt.Errorf("workflow must be map: %+v", v)
		}
		for k, v := range a {
			switch k {
			case "triggers":
				wf.Triggers = v
			case "jobs":
				arr, ok := v.([]interface{})
				if !ok {
					return fmt.Errorf("workflow jobs must be an array: workflow: %s: %+v", workflowName, v)
				}
				wf.Jobs = arr
			default:
				if key, ok := k.(string); ok {
					wf.others[key] = v
				} else {
					return fmt.Errorf("workflow's key must be a string: workflow: %s: %+v", workflowName, k)
				}
			}
		}
		wfMap[workflowName] = wf
	}
	wfs.Workflows = wfMap
	return nil
}

func (wfs Workflows) MarshalYAML() (interface{}, error) {
	m := make(map[string]interface{}, len(wfs.Workflows))
	for k, v := range wfs.Workflows {
		m[k] = v
	}
	if wfs.Version != nil {
		m["version"] = wfs.Version
	}
	return m, nil
}

type Workflow struct {
	Triggers interface{}            `yaml:",omitempty"`
	Jobs     []interface{}          `yaml:",omitempty"`
	others   map[string]interface{} `yaml:",omitempty"`
}

func (wf Workflow) MarshalYAML() (interface{}, error) {
	m := map[string]interface{}{}
	for k, v := range wf.others {
		m[k] = v
	}
	if wf.Triggers != nil {
		m["triggers"] = wf.Triggers
	}
	if wf.Jobs != nil {
		m["jobs"] = wf.Jobs
	}
	return m, nil
}

type Config struct {
	Version   interface{}            `yaml:",omitempty"`
	Orbs      map[string]interface{} `yaml:",omitempty"`
	Workflows Workflows
	Jobs      map[string]interface{} `yaml:",omitempty"`
	Commands  map[string]interface{} `yaml:",omitempty"`
	Executors map[string]interface{} `yaml:",omitempty"`
	others    map[string]interface{}
}

func (cfg Config) MarshalYAML() (interface{}, error) {
	m := map[string]interface{}{}
	for k, v := range cfg.others {
		m[k] = v
	}
	if cfg.Version != nil {
		m["version"] = cfg.Version
	}
	if cfg.Orbs != nil {
		m["orbs"] = cfg.Orbs
	}
	m["workflows"] = cfg.Workflows
	if cfg.Jobs != nil {
		m["jobs"] = cfg.Jobs
	}
	if cfg.Commands != nil {
		m["commands"] = cfg.Commands
	}
	if cfg.Executors != nil {
		m["executors"] = cfg.Executors
	}
	return m, nil
}

func New(ctx context.Context, params Params) (Controller, Params, error) {
	return Controller{
		Stdout: os.Stdout,
		Stderr: os.Stderr,
	}, params, nil
}

func (ctrl Controller) readFile(filePath string) (Config, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return Config{}, fmt.Errorf("open a file "+filePath+": %w", err)
	}
	defer file.Close()
	m := Config{}
	if err := yaml.NewDecoder(file).Decode(&m); err != nil {
		return m, fmt.Errorf("parse a file as YAML "+filePath+": %w", err)
	}
	return m, nil
}

func (ctrl Controller) mergeMap(base, child map[string]interface{}) map[string]interface{} {
	if base == nil {
		return child
	}
	for k, v := range child {
		base[k] = v
	}
	return base
}

func (ctrl Controller) mergeWorkflows(base, child Workflows) Workflows {
	if child.Version != nil {
		base.Version = child.Version
	}
	for k, childWorkflow := range child.Workflows {
		if baseWorkflow, ok := base.Workflows[k]; ok {
			baseWorkflow.Jobs = append(baseWorkflow.Jobs, childWorkflow.Jobs...)
			if childWorkflow.Triggers != nil {
				baseWorkflow.Triggers = childWorkflow.Triggers
			}
			baseWorkflow.others = ctrl.mergeMap(baseWorkflow.others, childWorkflow.others)
			base.Workflows[k] = baseWorkflow
		} else {
			if base.Workflows == nil {
				base.Workflows = map[string]Workflow{
					k: childWorkflow,
				}
			} else {
				base.Workflows[k] = childWorkflow
			}
		}
	}
	return base
}

func (ctrl Controller) mergeConfig(base, child Config) Config {
	if child.Version != nil {
		base.Version = child.Version
	}
	base.Orbs = ctrl.mergeMap(base.Orbs, child.Orbs)
	base.Workflows = ctrl.mergeWorkflows(base.Workflows, child.Workflows)
	base.Jobs = ctrl.mergeMap(base.Jobs, child.Jobs)
	base.Commands = ctrl.mergeMap(base.Commands, child.Commands)
	base.Executors = ctrl.mergeMap(base.Executors, child.Executors)
	return base
}

func (ctrl Controller) Run(ctx context.Context, params Params) error {
	cfg := Config{}
	for filePath := range params.Files {
		m, err := ctrl.readFile(filePath)
		if err != nil {
			return fmt.Errorf("read a file"+filePath+": %w", err)
		}
		cfg = ctrl.mergeConfig(cfg, m)
	}
	if err := yaml.NewEncoder(ctrl.Stdout).Encode(cfg); err != nil {
		return fmt.Errorf("encode a merged config as YAML: %w", err)
	}
	return nil
}
