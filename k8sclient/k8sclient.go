package k8sclient

import (
	"bytes"
	"context"
	"io"

	log "github.com/sirupsen/logrus"
	core "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type CreateParameters struct {
	Id        string   `json:"id"`
	Image     string   `json:"image"`
	Command   []string `json:"command"`
	Arguments []string `json:"args"`
}

type GetStatusParameters struct {
	Id     string `json:"id"`
	Status string `json:"status"`
	Logs   string `json:"logs"`
}

func CreatePod(clientset *kubernetes.Clientset, namespace string, params CreateParameters) error {
	podName := "task-pod-" + params.Id
	pod := &core.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      podName,
			Namespace: namespace,
			Labels: map[string]string{
				"app": podName,
			},
		},
		Spec: core.PodSpec{
			Containers: []core.Container{
				{
					Name:            "task-pod",
					Image:           params.Image,
					ImagePullPolicy: core.PullAlways,
					Command:         params.Command,
					Args:            params.Arguments,
				},
			},
			RestartPolicy: core.RestartPolicyNever,
		},
	}

	pod, err := clientset.CoreV1().Pods(namespace).Create(context.TODO(), pod, metav1.CreateOptions{})
	if err != nil {
		return err
	}

	log.Infof("Test pod %s created successfully.", pod.Name)
	return nil
}

func GetPod(clientset *kubernetes.Clientset, namespace string, id string) (*core.Pod, error) {
	podName := "task-pod-" + id
	return clientset.CoreV1().Pods(namespace).Get(context.TODO(), podName, metav1.GetOptions{})
}

func GetPodLogs(clientset *kubernetes.Clientset, namespace string, id string) (string, error) {
	podName := "task-pod-" + id
	req := clientset.CoreV1().Pods(namespace).GetLogs(podName, &core.PodLogOptions{})
	podLogs, err := req.Stream(context.TODO())
	if err != nil {
		return "error in opening stream", err
	}
	defer podLogs.Close()

	buf := new(bytes.Buffer)
	_, err = io.Copy(buf, podLogs)
	if err != nil {
		return "error in copy information from podLogs to buf", nil
	}
	str := buf.String()

	return str, nil
}
