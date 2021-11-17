package controller_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/suzuki-shunsuke/circleci-config-merge/pkg/controller"
)

//nolint:funlen
func TestConfig_MasharlYAML(t *testing.T) {
	t.Parallel()
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
				"workflows": &controller.Workflows{
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
			t.Parallel()
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
