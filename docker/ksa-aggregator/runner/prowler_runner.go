package runner

import (
	"context"
	"database/sql"
	"fmt"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"
	"strings"
	"time"
)

type ProwlerRunner struct {
	clientset *kubernetes.Clientset
	namespace string
	jobName   string
	fileName  string
}

func NewProwlerRunner(clientset *kubernetes.Clientset, namespace string, jobName string) *ProwlerRunner {
	return &ProwlerRunner{
		clientset: clientset,
		namespace: namespace,
		jobName:   jobName,
	}
}

func (tr *ProwlerRunner) Run() error {
	tr.fileName = fmt.Sprintf("prowler-%s.ocsf.json", time.Now().Format(TimeFormat))
	index := strings.Index(tr.fileName, ".")
	var fileNamePrefix = tr.fileName[:index]
	job := &batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name:      tr.jobName,
			Namespace: tr.namespace,
		},
		Spec: batchv1.JobSpec{
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Name: "prowler",
				},
				Spec: corev1.PodSpec{
					ServiceAccountName: "prowler",
					RestartPolicy:      corev1.RestartPolicyOnFailure,
					SecurityContext: &corev1.PodSecurityContext{
						SeccompProfile: &corev1.SeccompProfile{
							Type: corev1.SeccompProfileTypeRuntimeDefault,
						},
					},
					Containers: []corev1.Container{
						{
							Name:            "prowler",
							Image:           "ksa/prowler",
							ImagePullPolicy: corev1.PullIfNotPresent,
							Command: []string{
								"prowler", "kubernetes", "--output-formats", "json-ocsf", "--status", "FAIL", "MANUAL", "-F", fileNamePrefix, "-o", "/var/scan", "-z",
							},
							VolumeMounts: []corev1.VolumeMount{
								{
									Name:      "var-scan",
									MountPath: "/var/scan",
								},
							},
						},
					},
					Volumes: []corev1.Volume{
						{
							Name: "var-scan",
							VolumeSource: corev1.VolumeSource{
								PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
									ClaimName: "prowler-scan-results",
								},
							},
						},
					},
				},
			},
		},
	}

	_, err := tr.clientset.BatchV1().Jobs(tr.namespace).Create(context.TODO(), job, metav1.CreateOptions{})
	if err != nil {
		return fmt.Errorf("failed to create job %s: %w", tr.jobName, err)
	}

	return nil
}

func (tr *ProwlerRunner) Watch(db *sql.DB) (int, string) {
	fieldSelector := fmt.Sprintf("metadata.name=%s", tr.jobName)
	listOptions := metav1.ListOptions{FieldSelector: fieldSelector}
	watcher, err := tr.clientset.BatchV1().Jobs(tr.namespace).Watch(context.TODO(), listOptions)
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
				rows, err := db.Query("SELECT id FROM scanner where name=$1", "prowler")
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
						id, tr.fileName,
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

func (tr *ProwlerRunner) GetStatus() JobStatus {
	job, err := tr.clientset.BatchV1().Jobs(tr.namespace).Get(context.TODO(), tr.jobName, metav1.GetOptions{})
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

func (tr *ProwlerRunner) CleanUp() error {
	propagationPolicy := metav1.DeletePropagationForeground
	err := tr.clientset.BatchV1().Jobs(tr.namespace).Delete(context.TODO(), tr.jobName, metav1.DeleteOptions{
		PropagationPolicy: &propagationPolicy,
	})
	if err != nil {
		return fmt.Errorf("failed to delete job %s: %w", tr.jobName, err)
	}

	// Wait for the job to be deleted
	watchInterface, err := tr.clientset.BatchV1().Jobs(tr.namespace).Watch(context.TODO(), metav1.ListOptions{
		FieldSelector: fmt.Sprintf("metadata.name=%s", tr.jobName),
	})
	if err != nil {
		return fmt.Errorf("failed to watch job %s: %w", tr.jobName, err)
	}
	defer watchInterface.Stop()

	for event := range watchInterface.ResultChan() {
		switch event.Type {
		case watch.Deleted:
			fmt.Printf("Job %s and its pods have been deleted\n", tr.jobName)
			return nil
		case watch.Error:
			return fmt.Errorf("error watching job %s: %v", tr.jobName, event.Object)
		}
	}

	return nil
}
