package controller

import (
	"fmt"

	"github.com/mitchellh/mapstructure"
)

type Workflows struct {
	Version   interface{}          `yaml:",omitempty"`
	Workflows map[string]*Workflow `yaml:",omitempty"`
}

func mergeWorkflows(base, child *Workflows) *Workflows {
	if child.Version != nil {
		base.Version = child.Version
	}
	for k, childWorkflow := range child.Workflows {
		if baseWorkflow, ok := base.Workflows[k]; ok {
			baseWorkflow.Jobs = append(baseWorkflow.Jobs, childWorkflow.Jobs...)
			if childWorkflow.Triggers != nil {
				baseWorkflow.Triggers = childWorkflow.Triggers
			}
			base.Workflows[k] = baseWorkflow
		} else {
			if base.Workflows == nil {
				base.Workflows = map[string]*Workflow{
					k: childWorkflow,
				}
			} else {
				base.Workflows[k] = childWorkflow
			}
		}
	}
	return base
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

func (wfs *Workflows) UnmarshalYAML(unmarshal func(interface{}) error) error {
	m := map[string]interface{}{}
	if err := unmarshal(&m); err != nil {
		return err
	}
	if version, ok := m["version"]; ok {
		wfs.Version = version
		delete(m, "version")
	}
	wfMap := map[string]*Workflow{}
	if err := mapstructure.Decode(m, &wfMap); err != nil {
		return fmt.Errorf("decode map to structure: %w", err)
	}
	wfs.Workflows = wfMap
	return nil
}
