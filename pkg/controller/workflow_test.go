package controller_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/suzuki-shunsuke/circleci-config-merge/pkg/controller"
)

func TestWorkflow_MasharlYAML(t *testing.T) {
	t.Parallel()
	data := []struct {
		title string
		wf    *controller.Workflow
		exp   string
	}{
		{
			title: "normal",
			wf: &controller.Workflow{
				Jobs: []interface{}{
					"foo",
				},
			},
			exp: `
jobs:
- foo
`,
		},
	}
	for _, d := range data {
		t.Run(d.title, func(t *testing.T) {
			t.Parallel()
			diff, err := testMarshalYAML(d.exp, d.wf)
			require.NoError(t, err)
			if diff != "" {
				t.Fatal(diff)
			}
		})
	}
}
