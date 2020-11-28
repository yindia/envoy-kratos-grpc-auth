package model

import (
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
)

type AppConfig struct {
	Name string `json:"name"`
	Org string `json:"org"`
	Image string `json:"image"`
	Volume []corev1.Volume `json:"volume"`
	Replica *int32 `json:"replica"`
	Region string `json:"region"`
	Strategy appsv1.DeploymentStrategy `json:"strategy"`
	Ports string `json:"replica"`
	//Ports []corev1.ContainerPort `json:"ports"`
	ServicePort []corev1.ServicePort `json:"servicePort"`
	Command []string `json:"command"`
	VolumeMounts  []corev1.VolumeMount `json:"volumeMount"`
	EnvFrom  []corev1.EnvFromSource `json:"envFrom"`
	Env  []corev1.EnvVar `json:"env"`
	ImagePullPolicy  corev1.PullPolicy `json:"imagePullPolicy"`
	LivenessProbe *corev1.Probe `json:"livenessProbe"`
	ReadinessProbe *corev1.Probe `json:"readinessProbe"`
	WorkingDir string `json:"workingDir"`
}

type DBModel struct {
	Deployment appsv1.Deployment `json:"deployment"`
	Service corev1.Service`json:"service"`
	Org string `json:"org"`
	ID string `json:"id"`
}