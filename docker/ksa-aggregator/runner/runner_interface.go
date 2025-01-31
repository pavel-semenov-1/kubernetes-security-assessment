package runner

import "database/sql"

// Runner defines a common interface for all security scanner runners
type Runner interface {
	Run() error
	GetStatus() JobStatus
	CleanUp() error
	Watch(*sql.DB) string
}

type JobStatus struct {
	ActivePods    int32 `json:"active_pods"`
	SucceededPods int32 `json:"succeeded_pods"`
	FailedPods    int32 `json:"failed_pods"`
}

func (js *JobStatus) Active() bool {
	return js.ActivePods > 0
}

func (js *JobStatus) Succeeded() bool {
	return js.SucceededPods > 0
}

func (js *JobStatus) Failed() bool {
	return js.FailedPods > 0
}
