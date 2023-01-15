package controller

type Workflow struct {
	Triggers interface{}   `yaml:",omitempty"`
	Jobs     []interface{} `yaml:",omitempty"`
	When     interface{}   `yaml:",omitempty"`
}
