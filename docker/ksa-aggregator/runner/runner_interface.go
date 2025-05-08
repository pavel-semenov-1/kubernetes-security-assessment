package runner

import (
	"context"
	"database/sql"
	"fmt"
	batchv1 "k8s.io/api/batch/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"
)

// Runner defines a common interface for all security scanner runners
type Runner interface {
	Run() error
	GetStatus() JobStatus
	CleanUp() error
	Watch(*sql.DB) (int, string)
}

// JobRunner is base struct for all runners
type JobRunner struct {
	clientset   kubernetes.Interface
	namespace   string
	jobName     string
	scannerName string
	fileName    string
}

// JobStatus is used to track the status of the Kubernetes jobs
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

var TimeFormat = "2006-01-02-15-04-05"

// GetStatus returns the JobStatus object, which represents the status of the Kubernetes job for this runner
func (jr *JobRunner) GetStatus() JobStatus {
	job, err := jr.clientset.BatchV1().Jobs(jr.namespace).Get(context.TODO(), jr.jobName, metav1.GetOptions{})
	if err != nil {
		return JobStatus{
			ActivePods:    0,
			SucceededPods: 0,
			FailedPods:    0,
		}
	}

	return JobStatus{
		ActivePods:    job.Status.Active,
		SucceededPods: job.Status.Succeeded,
		FailedPods:    job.Status.Failed,
	}
}

// CleanUp removes Kubernetes job and its pods
func (jr *JobRunner) CleanUp() error {
	propagationPolicy := metav1.DeletePropagationForeground
	err := jr.clientset.BatchV1().Jobs(jr.namespace).Delete(context.TODO(), jr.jobName, metav1.DeleteOptions{
		PropagationPolicy: &propagationPolicy,
	})
	if err != nil {
		return fmt.Errorf("failed to delete job %s: %w", jr.jobName, err)
	}

	// Wait for the job to be deleted
	watchInterface, err := jr.clientset.BatchV1().Jobs(jr.namespace).Watch(context.TODO(), metav1.ListOptions{
		FieldSelector: fmt.Sprintf("metadata.name=%s", jr.jobName),
	})
	if err != nil {
		return fmt.Errorf("failed to watch job %s: %w", jr.jobName, err)
	}
	defer watchInterface.Stop()

	for event := range watchInterface.ResultChan() {
		switch event.Type {
		case watch.Deleted:
			fmt.Printf("Job %s and its pods have been deleted\n", jr.jobName)
			return nil
		case watch.Error:
			return fmt.Errorf("error watching job %s: %v", jr.jobName, event.Object)
		}
	}

	return nil
}

// Watch method waits for the job to finish and then returns the generated reportId
func (jr *JobRunner) Watch(db *sql.DB) (int, string) {
	fieldSelector := fmt.Sprintf("metadata.name=%s", jr.jobName)
	listOptions := metav1.ListOptions{FieldSelector: fieldSelector}
	watcher, err := jr.clientset.BatchV1().Jobs(jr.namespace).Watch(context.TODO(), listOptions)
	if err != nil {
		panic(err)
	}
	defer watcher.Stop()

	for event := range watcher.ResultChan() {
		job, ok := event.Object.(*batchv1.Job)
		if !ok {
			continue
		}

		for _, condition := range job.Status.Conditions {
			if condition.Type == batchv1.JobComplete && condition.Status == "True" {
				rows, err := db.Query("SELECT id FROM scanner where name=$1", jr.scannerName)
				if err != nil {
					panic(err)
				}
				defer rows.Close()

				if rows.Next() {
					var id int
					err = rows.Scan(&id)
					if err != nil {
						panic(err)
					}
					var reportId int
					err := db.QueryRow(
						"INSERT INTO report (scanner_id, filename, parsed, generated_at) VALUES ($1, $2, false, CURRENT_TIMESTAMP) RETURNING id",
						id, jr.fileName,
					).Scan(&reportId)
					if err != nil {
						panic(err)
					}
					return reportId, "Job has successfully finished"
				}
				return 0, "No such scanner"
			} else if condition.Type == batchv1.JobFailed && condition.Status == "True" {
				return 0, "Job has failed"
			}
		}
	}

	return 0, "Unexpected error occurred"
}
