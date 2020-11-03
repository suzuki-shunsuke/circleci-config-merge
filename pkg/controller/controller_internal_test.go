package controller

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestController_mergeMap(t *testing.T) {
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
			m := ctrl.mergeMap(d.base, d.child)
			require.Equal(t, d.exp, m)
		})
	}
}

func TestController_mergeWorkflows(t *testing.T) {
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
			m := ctrl.mergeWorkflows(d.base, d.child)
			require.Equal(t, d.exp, m)
		})
	}
}

func TestController_mergeConfig(t *testing.T) {
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
			m := ctrl.mergeConfig(d.base, d.child)
			require.Equal(t, d.exp, m)
		})
	}
}
