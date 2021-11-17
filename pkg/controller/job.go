package controller

import (
	"fmt"
	"sort"
)

func getJobName(job interface{}) (string, error) {
	switch v := job.(type) {
	case string:
		return v, nil
	case map[interface{}]interface{}:
		for jobName, jobValue := range v {
			a, ok := jobValue.(map[interface{}]interface{})
			if !ok {
				return "", fmt.Errorf("workflow job's element must be map: %+v", jobValue)
			}
			if name, ok := a["name"]; ok {
				jobName = name
			}
			s, ok := jobName.(string)
			if !ok {
				return "", fmt.Errorf("workflow job's name must be string: %+v", jobName)
			}
			return s, nil
		}
		return "", fmt.Errorf("workflow job's element is empty")
	default:
		return "", fmt.Errorf("workflow job must be string or map")
	}
}

func sortJobs(jobs []interface{}) ([]interface{}, error) {
	type WorkflowJob struct {
		Name string
		Job  interface{}
	}
	wfJobs := make([]WorkflowJob, len(jobs))
	for i, job := range jobs {
		name, err := getJobName(job)
		if err != nil {
			return nil, fmt.Errorf("get a job name: %w", err)
		}
		wfJobs[i] = WorkflowJob{
			Name: name,
			Job:  job,
		}
	}

	sort.Slice(wfJobs, func(i, j int) bool {
		return wfJobs[i].Name < wfJobs[j].Name
	})
	arr := make([]interface{}, len(jobs))
	for i, job := range wfJobs {
		arr[i] = job.Job
	}
	return arr, nil
}
