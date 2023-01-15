package controller_test

import (
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/require"
	"github.com/suzuki-shunsuke/circleci-config-merge/pkg/controller"
	"gopkg.in/yaml.v2"
)

func testMarshalYAML(exp, data interface{}) (string, error) {
	b, err := yaml.Marshal(data)
	if err != nil {
		return "", fmt.Errorf("marshal data as YAML: %w", err)
	}
	var cfgMap interface{}
	if err := yaml.Unmarshal(b, &cfgMap); err != nil {
		return "", fmt.Errorf("unmarshal data as YAML: %w", err)
	}
	var expMap interface{}
	if s, ok := exp.(string); ok {
		if err := yaml.Unmarshal([]byte(s), &expMap); err != nil {
			return "", fmt.Errorf("unmarshal exp as YAML: %w", err)
		}
	} else {
		b, err := yaml.Marshal(exp)
		if err != nil {
			return "", fmt.Errorf("marshal exp as YAML: %w", err)
		}
		if err := yaml.Unmarshal(b, &expMap); err != nil {
			return "", fmt.Errorf("unmarshal exp as YAML: %w", err)
		}
	}
	return cmp.Diff(expMap, cfgMap), nil
}

//nolint:funlen
func TestConfig_MasharlYAML(t *testing.T) {
	t.Parallel()
	data := []struct {
		title string
		cfg   *controller.Config
		exp   string
	}{
		{
			title: "normal",
			cfg: &controller.Config{
				Version: "2.1",
				Orbs: map[string]interface{}{
					"foo": "circleci/hello-build@0.0.5",
				},
				Workflows: &controller.Workflows{
					Version: "2",
					Workflows: map[string]*controller.Workflow{
						"build": {
							When: "<< pipeline.parameters.run_integration_tests >>",
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
			exp: `
version: "2.1"
orbs:
  foo: circleci/hello-build@0.0.5
workflows:
  version: "2"
  build:
    when: << pipeline.parameters.run_integration_tests >>
    jobs:
    - foo
    - bar
jobs:
  foo: {}
  bar: {}
commands:
  test: {}
executors:
  golang: {}
`,
		},
		{
			title: "logical statement",
			// https://circleci.com/docs/configuration-reference/#logic-statement-examples
			cfg: &controller.Config{
				Version: "2.1",
				Orbs: map[string]interface{}{
					"foo": "circleci/hello-build@0.0.5",
				},
				Workflows: &controller.Workflows{
					Version: "2",
					Workflows: map[string]*controller.Workflow{
						"build": {
							When: map[interface{}]interface{}{
								"or": []map[interface{}]interface{}{
									{
										"equal": []interface{}{
											"main",
											"<< pipeline.git.branch >>",
										},
									},
									{
										"equal": []interface{}{
											"staging",
											"<< pipeline.git.branch >>",
										},
									},
								},
							},
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
			exp: `
version: "2.1"
orbs:
  foo: circleci/hello-build@0.0.5
workflows:
  version: "2"
  build:
    when:
      or:
        - equal: [ main, << pipeline.git.branch >> ]
        - equal: [ staging, << pipeline.git.branch >> ]
    jobs:
    - foo
    - bar
jobs:
  foo: {}
  bar: {}
commands:
  test: {}
executors:
  golang: {}
`,
		},
	}
	for _, d := range data {
		d := d
		t.Run(d.title, func(t *testing.T) {
			t.Parallel()
			diff, err := testMarshalYAML(d.exp, d.cfg)
			require.Nil(t, err)
			if diff != "" {
				t.Fatal(diff)
			}
		})
	}
}
