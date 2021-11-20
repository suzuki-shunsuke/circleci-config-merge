package controller

type Workflow struct {
	Triggers interface{}   `yaml:",omitempty"`
	Jobs     []interface{} `yaml:",omitempty"`
	When     string        `yaml:",omitempty"`
}
