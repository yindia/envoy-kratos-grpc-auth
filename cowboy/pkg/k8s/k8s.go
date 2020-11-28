package k8s

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	config2 "github.com/evalsocket/envoy-kratos-grpc-auth/cowboy/pkg/config"
	"github.com/evalsocket/envoy-kratos-grpc-auth/cowboy/pkg/model"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

type K8sClient struct {
	region map[string]*kubernetes.Clientset
}

func NewK8s() *K8sClient {
	k := K8sClient{}
	for _, r := range config2.GetRegions() {
		config, err := clientcmd.BuildConfigFromFlags("", config2.GetKubeconfig(r))
		if err != nil {
			fmt.Println(err)
			panic(err.Error())
		}

		// create the clientset
		clientset, err := kubernetes.NewForConfig(config)
		if err != nil {
			fmt.Println(err)
			panic(err.Error())
		}
		if k.region[r] == nil {
			k.region = make(map[string]*kubernetes.Clientset)
		}
		k.region[r] = clientset
	}
	return &k
}

//CreateApplication
func (k *K8sClient) CreateApplication(appConfig model.AppConfig) (*appsv1.Deployment, *corev1.Service, error) {
	deployment, err := k.createDeployment(appConfig)
	if err != nil {
		return nil, nil, err
	}
	service, err := k.createService(appConfig)
	if err != nil {
		return nil, nil, err
	}
	deployment, err = k.region[appConfig.Region].AppsV1().Deployments(appConfig.Org).Create(context.Background(), deployment, metav1.CreateOptions{})
	if err != nil {
		return nil, nil, err
	}
	service, err = k.region[appConfig.Region].CoreV1().Services(appConfig.Org).Create(context.Background(), service, metav1.CreateOptions{})
	if err != nil {
		return nil, nil, err
	}
	return deployment, service, nil
}

func (k *K8sClient) createDeployment(appConfig model.AppConfig) (*appsv1.Deployment, error) {
	var volumes []corev1.Volume
	for _, v := range appConfig.Volume {
		volumes = append(volumes, v)
	}

	dep := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      fmt.Sprintf("%v-deployment", appConfig.Name),
			Namespace: appConfig.Org,
			Labels:    k.createLabels(appConfig),
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: appConfig.Replica,
			Strategy: appConfig.Strategy,
			Selector: &metav1.LabelSelector{
				MatchLabels: k.createLabels(appConfig),
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: k.createLabels(appConfig),
				},
				Spec: corev1.PodSpec{
					Containers: k.getContainer(appConfig),
				},
			},
		},
	}
	if len(volumes) > 0 {
		dep.Spec.Template.Spec.Volumes = volumes
	}
	return dep, nil
}

func (k *K8sClient) getContainer(appConfig model.AppConfig) []corev1.Container {
	po := strings.Split(appConfig.Ports, ",")

	var ports []corev1.ContainerPort

	for _, p := range po {
		tp, err := strconv.Atoi(p)
		if err != nil {
			fmt.Println("error in converting port")
		}
		ports = append(ports, corev1.ContainerPort{
			ContainerPort: int32(tp),
			Name:          "http",
		})
	}
	var command []string
	for _, c := range appConfig.Command {
		command = append(command, c)
	}

	var volumeMounts []corev1.VolumeMount
	for _, v := range appConfig.VolumeMounts {
		volumeMounts = append(volumeMounts, v)
	}

	// Attach envFrom if exist
	var envFrom []corev1.EnvFromSource
	for _, v := range appConfig.EnvFrom {
		envFrom = append(envFrom, v)
	}

	// Attach volumeMounts if exist
	var env []corev1.EnvVar
	for _, v := range appConfig.Env {
		env = append(env, v)
	}

	var imagePullPolicy corev1.PullPolicy = "Always"
	if appConfig.ImagePullPolicy != "" {
		imagePullPolicy = appConfig.ImagePullPolicy
	}

	dep := []corev1.Container{{
		Image:           appConfig.Image,
		Name:            fmt.Sprintf("%v-container", appConfig.Name),
		Command:         command,
		Ports:           ports,
		ImagePullPolicy: imagePullPolicy,
	}}
	if appConfig.WorkingDir != "" {
		dep[0].WorkingDir = appConfig.WorkingDir
	}
	if len(volumeMounts) > 0 {
		dep[0].VolumeMounts = volumeMounts
	}
	if len(envFrom) > 0 {
		dep[0].EnvFrom = envFrom
	}
	if len(env) > 0 {
		dep[0].Env = env
	}
	return dep
}

func (k *K8sClient) createService(appConfig model.AppConfig) (*corev1.Service, error) {

	po := strings.Split(appConfig.Ports, ",")

	var ports []corev1.ServicePort
	targetPort := intstr.IntOrString{
		Type:   1,
		IntVal: 8000,
	}
	for _, p := range po {
		tp, err := strconv.Atoi(p)
		if err != nil {
			fmt.Println("error in converting port")
		}
		ports = append(ports, corev1.ServicePort{
			Port:       int32(tp),
			Name:       "http",
			TargetPort: targetPort,
		})
	}

	service := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      fmt.Sprintf("%v-service", appConfig.Name),
			Namespace: appConfig.Org,
			Labels:    k.createLabels(appConfig),
		},
		Spec: corev1.ServiceSpec{
			Ports:    ports,
			Selector: k.createLabels(appConfig),
			Type:     "LoadBalancer",
		},
	}
	return service, nil
}

func (k *K8sClient) createLabels(appConfig model.AppConfig) map[string]string {
	return map[string]string{
		"name": appConfig.Name,
		"org":  appConfig.Org,
	}
}
