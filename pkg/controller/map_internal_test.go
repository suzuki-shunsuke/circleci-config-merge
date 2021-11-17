package controller

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_mergeMap(t *testing.T) {
	t.Parallel()
	data := []struct {
		title string
		exp   map[string]interface{}
		base  map[string]interface{}
		child map[string]interface{}
	}{
		{
			title: "normal",
			exp: map[string]interface{}{
				"foo": struct{}{},
				"bar": struct{}{},
			},
			base: map[string]interface{}{
				"foo": struct{}{},
			},
			child: map[string]interface{}{
				"bar": struct{}{},
			},
		},
	}
	for _, d := range data {
		d := d
		t.Run(d.title, func(t *testing.T) {
			t.Parallel()
			m := mergeMap(d.base, d.child)
			require.Equal(t, d.exp, m)
		})
	}
}
