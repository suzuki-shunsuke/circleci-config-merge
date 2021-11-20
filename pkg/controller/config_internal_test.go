package controller

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_mergeConfig(t *testing.T) {
	t.Parallel()
	data := []struct {
		title string
		exp   Config
		base  Config
		child Config
	}{
		{
			title: "normal",
			exp: Config{
				Version: "2.0",
				Workflows: &Workflows{
					Workflows: map[string]Workflow{
						"build": {
							Jobs: []interface{}{
								"foo",
								"bar",
							},
						},
					},
				},
			},
			base: Config{
				Version: "2.0",
				Workflows: &Workflows{
					Workflows: map[string]Workflow{
						"build": {
							Jobs: []interface{}{
								"foo",
							},
						},
					},
				},
			},
			child: Config{
				Version: "2.0",
				Workflows: &Workflows{
					Workflows: map[string]Workflow{
						"build": {
							Jobs: []interface{}{
								"bar",
							},
						},
					},
				},
			},
		},
	}
	for _, d := range data {
		d := d
		t.Run(d.title, func(t *testing.T) {
			t.Parallel()
			m := mergeConfig(d.base, d.child)
			require.Equal(t, d.exp, m)
		})
	}
}
