package controller

import (
	"fmt"
)

type Workflows struct {
	Version   interface{}         `yaml:",omitempty"`
	Workflows map[string]Workflow `yaml:",omitempty"`
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
