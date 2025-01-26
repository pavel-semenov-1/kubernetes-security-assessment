package runner

import (
	"context"
	"fmt"
	"time"

	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type TrivyRunner struct {
	clientset *kubernetes.Clientset
	namespace string
	jobName   string
}

var TimeFormat = "2006-01-02-15-04-05"

func NewTrivyRunner(clientset *kubernetes.Clientset, namespace string) *TrivyRunner {
	return &TrivyRunner{
		clientset: clientset,
		namespace: namespace,
	}
}

func (tr *TrivyRunner) Run() string {
	job := &batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "trivy-runner",
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
							Args:  []string{"kubernetes", "--format", "json", "--output", fmt.Sprintf("/var/scan/trivy-%s.json", time.Now().Format(TimeFormat))},
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

	createdJob, err := tr.clientset.BatchV1().Jobs(tr.namespace).Create(context.TODO(), job, metav1.CreateOptions{})
	if err != nil {
		panic("Failed to create job: " + err.Error())
	}

	tr.jobName = createdJob.Name

	return tr.jobName
}

func (tr *TrivyRunner) GetStatus() JobStatus {
	if tr.jobName == "" {
		return JobStatus{}
	}

	job, err := tr.clientset.BatchV1().Jobs(tr.namespace).Get(context.TODO(), tr.jobName, metav1.GetOptions{})
	if err != nil {
		panic(fmt.Sprintf("Failed to get job %s: %v", tr.jobName, err))
	}

	return JobStatus{
		active_pods:    job.Status.Active,
		succeeded_pods: job.Status.Succeeded,
		failed_pods:    job.Status.Failed,
	}
}
