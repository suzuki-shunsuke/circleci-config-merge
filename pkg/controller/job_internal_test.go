package controller

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_getJobName(t *testing.T) {
	t.Parallel()
	data := []struct {
		title string
		exp   string
		isErr bool
		job   interface{}
	}{
		{
			title: "normal",
			exp:   "foo",
			job:   "foo",
		},
		{
			title: "map",
			exp:   "zoo",
			job: map[interface{}]interface{}{
				"zoo": map[interface{}]interface{}{
					"requires": []interface{}{
						"bar",
					},
				},
			},
		},
	}
	for _, d := range data {
		d := d
		t.Run(d.title, func(t *testing.T) {
			t.Parallel()
			name, err := getJobName(d.job)
			if d.isErr {
				require.NotNil(t, err)
				return
			}
			require.Nil(t, err)
			require.Equal(t, d.exp, name)
		})
	}
}

func Test_sortJobs(t *testing.T) {
	t.Parallel()
	data := []struct {
		title string
		exp   []interface{}
		isErr bool
		jobs  []interface{}
	}{
		{
			title: "normal",
			exp: []interface{}{
				"bar",
				"foo",
			},
			jobs: []interface{}{
				"foo",
				"bar",
			},
		},
		{
			title: "map",
			exp: []interface{}{
				"foo",
				map[interface{}]interface{}{
					"zoo": map[interface{}]interface{}{
						"requires": []interface{}{
							"bar",
						},
					},
				},
			},
			jobs: []interface{}{
				map[interface{}]interface{}{
					"zoo": map[interface{}]interface{}{
						"requires": []interface{}{
							"bar",
						},
					},
				},
				"foo",
			},
		},
	}
	for _, d := range data {
		d := d
		t.Run(d.title, func(t *testing.T) {
			t.Parallel()
			jobs, err := sortJobs(d.jobs)
			if d.isErr {
				require.NotNil(t, err)
				return
			}
			require.Nil(t, err)
			require.Equal(t, d.exp, jobs)
		})
	}
}
