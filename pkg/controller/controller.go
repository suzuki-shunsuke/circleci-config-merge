package controller

import (
	"context"
	"fmt"
	"os"
	"sort"

	"gopkg.in/yaml.v2"
)

type Params struct {
	Files map[string]struct{}
}

type Workflows struct {
	Version   interface{}         `yaml:",omitempty"`
	Workflows map[string]Workflow `yaml:",omitempty"`
}

func (wfs *Workflows) UnmarshalYAML(unmarshal func(interface{}) error) error { //nolint:cyclop
	m := map[string]interface{}{}
	if err := unmarshal(&m); err != nil {
		return err
	}
	wfMap := map[string]Workflow{}
	for workflowName, v := range m {
		if workflowName == "version" {
			wfs.Version = v
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
					if wf.others == nil {
						wf.others = map[string]interface{}{}
					}
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

func (wfs *Workflows) MarshalYAML() (interface{}, error) {
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

func getJobName(job interface{}) (string, error) {
	switch v := job.(type) {
	case string:
		return v, nil
	case map[interface{}]interface{}:
		for jobName, jobValue := range v {
			a, ok := jobValue.(map[interface{}]interface{})
			if !ok {
				return "", fmt.Errorf("workflow job's element must be map: %+v", jobValue)
			}
			if name, ok := a["name"]; ok {
				jobName = name
			}
			s, ok := jobName.(string)
			if !ok {
				return "", fmt.Errorf("workflow job's name must be string: %+v", jobName)
			}
			return s, nil
		}
		return "", fmt.Errorf("workflow job's element is empty")
	default:
		return "", fmt.Errorf("workflow job must be string or map")
	}
}

func sortJobs(jobs []interface{}) ([]interface{}, error) {
	type WorkflowJob struct {
		Name string
		Job  interface{}
	}
	wfJobs := make([]WorkflowJob, len(jobs))
	for i, job := range jobs {
		name, err := getJobName(job)
		if err != nil {
			return nil, fmt.Errorf("get a job name: %w", err)
		}
		wfJobs[i] = WorkflowJob{
			Name: name,
			Job:  job,
		}
	}

	sort.Slice(wfJobs, func(i, j int) bool {
		return wfJobs[i].Name < wfJobs[j].Name
	})
	arr := make([]interface{}, len(jobs))
	for i, job := range wfJobs {
		arr[i] = job.Job
	}
	return arr, nil
}

func (wf *Workflow) MarshalYAML() (interface{}, error) {
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

func (cfg *Config) MarshalYAML() (interface{}, error) {
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

func readFile(filePath string) (Config, error) {
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

func mergeMap(base, child map[string]interface{}) map[string]interface{} {
	if base == nil {
		return child
	}
	for k, v := range child {
		base[k] = v
	}
	return base
}

func mergeWorkflows(base, child Workflows) Workflows {
	if child.Version != nil {
		base.Version = child.Version
	}
	for k, childWorkflow := range child.Workflows {
		if baseWorkflow, ok := base.Workflows[k]; ok {
			baseWorkflow.Jobs = append(baseWorkflow.Jobs, childWorkflow.Jobs...)
			if childWorkflow.Triggers != nil {
				baseWorkflow.Triggers = childWorkflow.Triggers
			}
			baseWorkflow.others = mergeMap(baseWorkflow.others, childWorkflow.others)
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

func mergeConfig(base, child Config) Config {
	if child.Version != nil {
		base.Version = child.Version
	}
	base.Orbs = mergeMap(base.Orbs, child.Orbs)
	base.Workflows = mergeWorkflows(base.Workflows, child.Workflows)
	base.Jobs = mergeMap(base.Jobs, child.Jobs)
	base.Commands = mergeMap(base.Commands, child.Commands)
	base.Executors = mergeMap(base.Executors, child.Executors)
	return base
}

func (ctrl *Controller) Run(ctx context.Context, params Params) error {
	cfg := Config{}
	for filePath := range params.Files {
		m, err := readFile(filePath)
		if err != nil {
			return fmt.Errorf("read a file "+filePath+": %w", err)
		}
		cfg = mergeConfig(cfg, m)
	}
	for k, workflow := range cfg.Workflows.Workflows {
		jobs, err := sortJobs(workflow.Jobs)
		if err == nil {
			workflow.Jobs = jobs
		}
		cfg.Workflows.Workflows[k] = workflow
	}
	if err := yaml.NewEncoder(ctrl.Stdout).Encode(cfg); err != nil {
		return fmt.Errorf("encode a merged config as YAML: %w", err)
	}
	return nil
}
