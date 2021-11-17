package controller

type Workflow struct {
	Triggers interface{}            `yaml:",omitempty"`
	Jobs     []interface{}          `yaml:",omitempty"`
	others   map[string]interface{} `yaml:",omitempty"`
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
