package probe

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

func (ph *probeHelper) createPod(ctx context.Context) error {
	var kubeconfig string
	if home := homedir.HomeDir(); home != "" {
		kubeconfig = filepath.Join(home, ".kube", "config")
	}

	// use the current context in kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		config, err = rest.InClusterConfig()
		if err != nil {
			return err
		}
	}

	// create the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return err
	}

	// create a pod
	pod := &v1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name: "host-collector-" + ph.host.Id[:6],
		},
		Spec: v1.PodSpec{
			Containers: []v1.Container{
				{
					Name:  "host-collector",
					Image: os.Getenv("HOST_COLLECTOR_IMAGE"),
					Env: []v1.EnvVar{
						{
							Name:  "HOST",
							Value: ph.sshInfo.EndPoint,
						},
						{
							Name:  "PORT",
							Value: fmt.Sprintf("%d", ph.sshInfo.Port),
						},
						{
							Name:  "USER",
							Value: ph.sshInfo.User,
						},
						{
							Name:  "PASSWORD",
							Value: ph.sshInfo.Password,
						},
						{
							Name:  "OSTYPE",
							Value: ph.sshInfo.OSType, 
						},
						{
							Name:  "HOST_ID",
							Value: ph.host.Id,
						},
						{
							Name: "INFLUXDB_TOKEN",
							Value: os.Getenv("INFLUXDB_TOKEN"),
						},
						{
							Name: "INFLUXDB_URL",
							Value: os.Getenv("INFLUXDB_URL"),
						},
						{
							Name: "INFLUXDB_ORG",
							Value: os.Getenv("INFLUXDB_ORG"),
						},
						{
							Name: "INFLUXDB_BUCKET",
							Value: os.Getenv("INFLUXDB_BUCKET"),
						},
					},
				},
			},
		},
	}

	// 在"default"命名空间创建Pod
	_, err = clientset.CoreV1().Pods("default").Create(context.TODO(), pod, metav1.CreateOptions{})
	if err != nil {
		return err
	}

	return nil
}
