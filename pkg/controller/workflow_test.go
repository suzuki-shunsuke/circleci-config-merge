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
			t.Parallel()
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
