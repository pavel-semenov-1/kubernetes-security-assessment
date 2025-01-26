package runner

// Runner defines a common interface for all security scanner runners
type Runner interface {
	Run() string
	GetStatus() JobStatus
}

type JobStatus struct {
	active_pods    int32
	succeeded_pods int32
	failed_pods    int32
}

type JobState struct {
	active    bool
	succeeded bool
	failed    bool
}

func (js *JobStatus) Active() bool {
	return js.active_pods > 0
}

func (js *JobStatus) Succeeded() bool {
	return js.succeeded_pods > 0
}

func (js *JobStatus) Failed() bool {
	return js.failed_pods > 0
}

func GetJobState(status JobStatus) JobState {
	return JobState{
		active:    status.Active(),
		succeeded: status.Succeeded(),
		failed:    status.Failed(),
	}
}
