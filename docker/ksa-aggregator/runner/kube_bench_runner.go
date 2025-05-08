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

type KubeBenchRunner struct {
	JobRunner
}

func NewKubeBenchRunner(clientset *kubernetes.Clientset, namespace string, jobName string, scannerName string) *KubeBenchRunner {
	return &KubeBenchRunner{
		JobRunner: JobRunner{
			clientset:   clientset,
			namespace:   namespace,
			jobName:     jobName,
			scannerName: scannerName,
		},
	}
}

func (tr *KubeBenchRunner) Run() error {
	tr.fileName = fmt.Sprintf("kube-bench-%s.json", time.Now().Format(TimeFormat))
	job := &batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name:      tr.jobName,
			Namespace: tr.namespace,
		},
		Spec: batchv1.JobSpec{
			Template: corev1.PodTemplateSpec{
				Spec: corev1.PodSpec{
					ServiceAccountName: "kube-bench",
					SecurityContext: &corev1.PodSecurityContext{
						SeccompProfile: &corev1.SeccompProfile{
							Type: corev1.SeccompProfileTypeRuntimeDefault,
						},
					},
					RestartPolicy: corev1.RestartPolicyOnFailure,
					Containers: []corev1.Container{
						{
							Name:  "kube-bench-runner",
							Image: "docker.io/aquasec/kube-bench:v0.9.3",
							Args:  []string{"--json", "--nototals", "--outputfile", fmt.Sprintf("/var/scan/%s", tr.fileName)},
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
						{Name: "var-scan", VolumeSource: corev1.VolumeSource{PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{ClaimName: "kube-bench-scan-results"}}},
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
