package controller

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Version    interface{}            `yaml:",omitempty"`
	Orbs       map[string]interface{} `yaml:",omitempty"`
	Workflows  *Workflows
	Jobs       map[string]interface{} `yaml:",omitempty"`
	Commands   map[string]interface{} `yaml:",omitempty"`
	Executors  map[string]interface{} `yaml:",omitempty"`
	Parameters map[string]interface{} `yaml:",omitempty"`
}

func readFile(filePath string, cfg *Config) error {
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("open a file "+filePath+": %w", err)
	}
	defer file.Close()
	if err := yaml.NewDecoder(file).Decode(cfg); err != nil {
		return fmt.Errorf("parse a file as YAML "+filePath+": %w", err)
	}
	return nil
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
	base.Parameters = mergeMap(base.Parameters, child.Parameters)
	return base
}
