package controller_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/suzuki-shunsuke/circleci-config-merge/pkg/controller"
	"gopkg.in/yaml.v2"
)

func TestWorkflows_UnmasharlYAML(t *testing.T) {
	t.Parallel()
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
				Workflows: map[string]*controller.Workflow{
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
			t.Parallel()
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
	t.Parallel()
	data := []struct {
		title string
		wfs   *controller.Workflows
		exp   interface{}
	}{
		{
			title: "normal",
			wfs: &controller.Workflows{
				Version: "2.1",
				Workflows: map[string]*controller.Workflow{
					"build": {
						When: "<< pipeline.parameters.run_integration_tests >>",
						Jobs: []interface{}{
							"foo",
						},
					},
				},
			},
			exp: `
version: "2.1"
build:
  when: << pipeline.parameters.run_integration_tests >>
  jobs:
  - foo
`,
		},
	}
	for _, d := range data {
		d := d
		t.Run(d.title, func(t *testing.T) {
			t.Parallel()
			diff, err := testMarshalYAML(d.exp, d.wfs)
			require.Nil(t, err)
			if diff != "" {
				t.Fatal(diff)
			}
		})
	}
}
