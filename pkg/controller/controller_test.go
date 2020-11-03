package controller_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/suzuki-shunsuke/circleci-config-merge/pkg/controller"
	"gopkg.in/yaml.v2"
)

func TestWorkflows_UnmasharlYAML(t *testing.T) {
	data := []struct {
		title string
		exp   controller.Workflows
		yaml  []byte
		isErr bool
	}{
		{
			title: "normal",
			exp: controller.Workflows{
				Version: "2.1",
				Workflows: map[string]controller.Workflow{
					"build": {
						Jobs: []interface{}{
							"foo",
						},
					},
				},
			},
			yaml: []byte(`
version: "2.1"
build:
  jobs:
  - foo
`),
		},
	}
	for _, d := range data {
		d := d
		t.Run(d.title, func(t *testing.T) {
			wfs := controller.Workflows{}
			err := yaml.Unmarshal(d.yaml, &wfs)
			if d.isErr {
				require.NotNil(t, err)
				return
			}
			require.Nil(t, err)
			require.Equal(t, d.exp, wfs)
		})
	}
}

func TestWorkflows_MasharlYAML(t *testing.T) {
	data := []struct {
		title string
		wfs   controller.Workflows
		isErr bool
		exp   interface{}
	}{
		{
			title: "normal",
			wfs: controller.Workflows{
				Version: "2.1",
				Workflows: map[string]controller.Workflow{
					"build": {
						Jobs: []interface{}{
							"foo",
						},
					},
				},
			},
			exp: map[string]interface{}{
				"version": "2.1",
				"build": controller.Workflow{
					Jobs: []interface{}{
						"foo",
					},
				},
			},
		},
	}
	for _, d := range data {
		d := d
		t.Run(d.title, func(t *testing.T) {
			b, err := d.wfs.MarshalYAML()
			if d.isErr {
				require.NotNil(t, err)
				return
			}
			require.Nil(t, err)
			require.Equal(t, d.exp, b)
		})
	}
}

func TestWorkflow_MasharlYAML(t *testing.T) {
	data := []struct {
		title string
		wf    controller.Workflow
		isErr bool
		exp   interface{}
	}{
		{
			title: "normal",
			wf: controller.Workflow{
				Jobs: []interface{}{
					"foo",
				},
			},
			exp: map[string]interface{}{
				"jobs": []interface{}{
					"foo",
				},
			},
		},
	}
	for _, d := range data {
		d := d
		t.Run(d.title, func(t *testing.T) {
			b, err := d.wf.MarshalYAML()
			if d.isErr {
				require.NotNil(t, err)
				return
			}
			require.Nil(t, err)
			require.Equal(t, d.exp, b)
		})
	}
}

//nolint:funlen
func TestConfig_MasharlYAML(t *testing.T) {
	data := []struct {
		title string
		cfg   controller.Config
		isErr bool
		exp   interface{}
	}{
		{
			title: "normal",
			cfg: controller.Config{
				Version: "2.1",
				Orbs: map[string]interface{}{
					"foo": "circleci/hello-build@0.0.5",
				},
				Workflows: controller.Workflows{
					Version: "2",
					Workflows: map[string]controller.Workflow{
						"build": {
							Jobs: []interface{}{
								"foo", "bar",
							},
						},
					},
				},
				Jobs: map[string]interface{}{
					"foo": map[string]interface{}{},
					"bar": map[string]interface{}{},
				},
				Commands: map[string]interface{}{
					"test": map[string]interface{}{},
				},
				Executors: map[string]interface{}{
					"golang": map[string]interface{}{},
				},
			},
			exp: map[string]interface{}{
				"version": "2.1",
				"orbs": map[string]interface{}{
					"foo": "circleci/hello-build@0.0.5",
				},
				"workflows": controller.Workflows{
					Version: "2",
					Workflows: map[string]controller.Workflow{
						"build": {
							Jobs: []interface{}{
								"foo", "bar",
							},
						},
					},
				},
				"jobs": map[string]interface{}{
					"foo": map[string]interface{}{},
					"bar": map[string]interface{}{},
				},
				"commands": map[string]interface{}{
					"test": map[string]interface{}{},
				},
				"executors": map[string]interface{}{
					"golang": map[string]interface{}{},
				},
			},
		},
	}
	for _, d := range data {
		d := d
		t.Run(d.title, func(t *testing.T) {
			b, err := d.cfg.MarshalYAML()
			if d.isErr {
				require.NotNil(t, err)
				return
			}
			require.Nil(t, err)
			require.Equal(t, d.exp, b)
		})
	}
}
