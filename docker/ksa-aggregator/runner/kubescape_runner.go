package runner

import (
	"context"
	"fmt"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"time"
)

type KubescapeRunner struct {
	JobRunner
}

func NewKubescapeRunner(clientset *kubernetes.Clientset, namespace string, jobName string, scannerName string) *KubescapeRunner {
	return &KubescapeRunner{
		JobRunner: JobRunner{
			clientset:   clientset,
			namespace:   namespace,
			jobName:     jobName,
			scannerName: scannerName,
		},
	}
}

func (tr *KubescapeRunner) Run() error {
	tr.fileName = fmt.Sprintf("kubescape-%s.json", time.Now().Format(TimeFormat))

	job := &batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name:      tr.jobName,
			Namespace: tr.namespace,
		},
		Spec: batchv1.JobSpec{
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Name: "kubescape",
				},
				Spec: corev1.PodSpec{
					ServiceAccountName: "kubescape-discovery",
					RestartPolicy:      corev1.RestartPolicyOnFailure,
					SecurityContext: &corev1.PodSecurityContext{
						SeccompProfile: &corev1.SeccompProfile{
							Type: corev1.SeccompProfileTypeRuntimeDefault,
						},
					},
					Containers: []corev1.Container{
						{
							Name:            "kubescape",
							Image:           "ksa/kubescape:latest",
							ImagePullPolicy: corev1.PullIfNotPresent,
							Command: []string{
								"sh",
								"-c",
								"kubescape scan --format json --output " + fmt.Sprintf("/var/scan/%s", tr.fileName),
							},
							VolumeMounts: []corev1.VolumeMount{
								{Name: "var-scan", MountPath: "/var/scan"},
							},
						},
					},
					Volumes: []corev1.Volume{
						{Name: "var-scan", VolumeSource: corev1.VolumeSource{PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{ClaimName: "kubescape-scan-results"}}},
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
