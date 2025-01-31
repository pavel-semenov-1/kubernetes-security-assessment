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
	"time"
)

type TrivyRunner struct {
	clientset *kubernetes.Clientset
	namespace string
	jobName   string
	fileName  string
}

var TimeFormat = "2006-01-02-15-04-05"

func NewTrivyRunner(clientset *kubernetes.Clientset, namespace string, jobName string) *TrivyRunner {
	return &TrivyRunner{
		clientset: clientset,
		namespace: namespace,
		jobName:   jobName,
	}
}

func (tr *TrivyRunner) Run() error {
	tr.fileName = fmt.Sprintf("trivy-%s.json", time.Now().Format(TimeFormat))
	job := &batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name:      tr.jobName,
			Namespace: tr.namespace,
		},
		Spec: batchv1.JobSpec{
			Template: corev1.PodTemplateSpec{
				Spec: corev1.PodSpec{
					ServiceAccountName: "trivy",
					SecurityContext: &corev1.PodSecurityContext{
						SeccompProfile: &corev1.SeccompProfile{
							Type: corev1.SeccompProfileTypeRuntimeDefault,
						},
					},
					InitContainers: []corev1.Container{
						{
							Name:    "cluster-preparation",
							Image:   "alpine/k8s@sha256:fea4057e9e8a0d363ac4f67e55cf6ab8a6e44a057e0e0315ca7b5068927c5fdc",
							Command: []string{"kubectl", "delete", "namespace", "trivy-temp", "--ignore-not-found", "true"},
						},
					},
					RestartPolicy: corev1.RestartPolicyOnFailure,
					Containers: []corev1.Container{
						{
							Name:  "trivy-runner",
							Image: "aquasec/trivy:latest",
							Args:  []string{"kubernetes", "--format", "json", "--output", fmt.Sprintf("/var/scan/%s", tr.fileName)},
							VolumeMounts: []corev1.VolumeMount{
								{Name: "var-scan", MountPath: "/var/scan"},
								{Name: "var-lib-cni", MountPath: "/var/lib/cni", ReadOnly: true},
								{Name: "var-lib-etcd", MountPath: "/var/lib/etcd", ReadOnly: true},
								{Name: "var-lib-kubelet", MountPath: "/var/lib/kubelet", ReadOnly: true},
								{Name: "var-lib-kube-scheduler", MountPath: "/var/lib/kube-scheduler", ReadOnly: true},
								{Name: "var-lib-kube-controller-manager", MountPath: "/var/lib/kube-controller-manager", ReadOnly: true},
								{Name: "etc-systemd", MountPath: "/etc/systemd", ReadOnly: true},
								{Name: "lib-systemd", MountPath: "/lib/systemd/", ReadOnly: true},
								{Name: "srv-kubernetes", MountPath: "/srv/kubernetes/", ReadOnly: true},
								{Name: "etc-kubernetes", MountPath: "/etc/kubernetes", ReadOnly: true},
								{Name: "usr-bin", MountPath: "/usr/local/mount-from-host/bin", ReadOnly: true},
								{Name: "etc-cni-netd", MountPath: "/etc/cni/net.d/", ReadOnly: true},
								{Name: "opt-cni-bin", MountPath: "/opt/cni/bin/", ReadOnly: true},
							},
						},
					},
					Volumes: []corev1.Volume{
						{Name: "var-scan", VolumeSource: corev1.VolumeSource{PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{ClaimName: "trivy-scan-results"}}},
						{Name: "var-lib-cni", VolumeSource: corev1.VolumeSource{HostPath: &corev1.HostPathVolumeSource{Path: "/var/lib/cni"}}},
						{Name: "var-lib-etcd", VolumeSource: corev1.VolumeSource{HostPath: &corev1.HostPathVolumeSource{Path: "/var/lib/etcd"}}},
						{Name: "var-lib-kubelet", VolumeSource: corev1.VolumeSource{HostPath: &corev1.HostPathVolumeSource{Path: "/var/lib/kubelet"}}},
						{Name: "var-lib-kube-scheduler", VolumeSource: corev1.VolumeSource{HostPath: &corev1.HostPathVolumeSource{Path: "/var/lib/kube-scheduler"}}},
						{Name: "var-lib-kube-controller-manager", VolumeSource: corev1.VolumeSource{HostPath: &corev1.HostPathVolumeSource{Path: "/var/lib/kube-controller-manager"}}},
						{Name: "etc-systemd", VolumeSource: corev1.VolumeSource{HostPath: &corev1.HostPathVolumeSource{Path: "/etc/systemd"}}},
						{Name: "lib-systemd", VolumeSource: corev1.VolumeSource{HostPath: &corev1.HostPathVolumeSource{Path: "/lib"}}},
						{Name: "srv-kubernetes", VolumeSource: corev1.VolumeSource{HostPath: &corev1.HostPathVolumeSource{Path: "/srv/kubernetes"}}},
						{Name: "etc-kubernetes", VolumeSource: corev1.VolumeSource{HostPath: &corev1.HostPathVolumeSource{Path: "/etc/kubernetes"}}},
						{Name: "usr-bin", VolumeSource: corev1.VolumeSource{HostPath: &corev1.HostPathVolumeSource{Path: "/usr/bin"}}},
						{Name: "etc-cni-netd", VolumeSource: corev1.VolumeSource{HostPath: &corev1.HostPathVolumeSource{Path: "/etc/cni/net.d/"}}},
						{Name: "opt-cni-bin", VolumeSource: corev1.VolumeSource{HostPath: &corev1.HostPathVolumeSource{Path: "/opt/cni/bin"}}},
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

func (tr *TrivyRunner) Watch(db *sql.DB) string {
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
				rows, err := db.Query("SELECT id FROM scanner where name=$1", "trivy")
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
					_, err = db.Exec("INSERT INTO report (scanner_id, filename, parsed, generated_at) VALUES ($1, $2, false, CURRENT_TIMESTAMP)", id, tr.fileName)
					if err != nil {
						panic(err)
					}
				}
				return "Job has successfully finished"
			} else if condition.Type == batchv1.JobFailed && condition.Status == "True" {
				return "Job has failed"
			}
		}
	}

	return "Unexpected error occurred"
}

func (tr *TrivyRunner) GetStatus() JobStatus {
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

func (tr *TrivyRunner) CleanUp() error {
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
