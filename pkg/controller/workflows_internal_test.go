package controller

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_mergeWorkflows(t *testing.T) {
	t.Parallel()
	data := []struct {
		title string
		exp   *Workflows
		base  *Workflows
		child *Workflows
	}{
		{
			title: "normal",
			exp: &Workflows{
				Version: "2.0",
				Workflows: map[string]Workflow{
					"build": {
						Jobs: []interface{}{
							"foo",
							"bar",
						},
					},
				},
			},
			base: &Workflows{
				Version: "2.0",
				Workflows: map[string]Workflow{
					"build": {
						Jobs: []interface{}{
							"foo",
						},
					},
				},
			},
			child: &Workflows{
				Workflows: map[string]Workflow{
					"build": {
						Jobs: []interface{}{
							"bar",
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
			m := mergeWorkflows(d.base, d.child)
			require.Equal(t, d.exp, m)
		})
	}
}
