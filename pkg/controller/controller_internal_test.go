package controller

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestController_mergeMap(t *testing.T) {
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
	ctrl := Controller{}
	for _, d := range data {
		d := d
		t.Run(d.title, func(t *testing.T) {
			t.Parallel()
			m := ctrl.mergeMap(d.base, d.child)
			require.Equal(t, d.exp, m)
		})
	}
}

func TestController_mergeWorkflows(t *testing.T) {
	t.Parallel()
	data := []struct {
		title string
		exp   Workflows
		base  Workflows
		child Workflows
	}{
		{
			title: "normal",
			exp: Workflows{
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
			base: Workflows{
				Version: "2.0",
				Workflows: map[string]Workflow{
					"build": {
						Jobs: []interface{}{
							"foo",
						},
					},
				},
			},
			child: Workflows{
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
	ctrl := Controller{}
	for _, d := range data {
		d := d
		t.Run(d.title, func(t *testing.T) {
			t.Parallel()
			m := ctrl.mergeWorkflows(d.base, d.child)
			require.Equal(t, d.exp, m)
		})
	}
}

func TestController_mergeConfig(t *testing.T) {
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
				Workflows: Workflows{
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
				Workflows: Workflows{
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
				Workflows: Workflows{
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
	ctrl := Controller{}
	for _, d := range data {
		d := d
		t.Run(d.title, func(t *testing.T) {
			t.Parallel()
			m := ctrl.mergeConfig(d.base, d.child)
			require.Equal(t, d.exp, m)
		})
	}
}

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
