package k8sclient

import (
	"bytes"
	"context"
	"io"
	"strconv"

	log "github.com/sirupsen/logrus"
	core "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type CreateParameters struct {
	Id        string   `json:"id"`
	Image     string   `json:"image"`
	Arguments []string `json:"args"`
	Timeout   int      `json:"timeout"`
	OutputDir string   `json:"outputDir"`
}

type GetStatusParameters struct {
	Id     string `json:"id"`
	Status string `json:"status"`
	Logs   string `json:"logs"`
}

func CreatePod(clientset *kubernetes.Clientset, namespace string, params CreateParameters) error {
	podName := "pytest-pod-" + params.Id
	pod := &core.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      podName,
			Namespace: namespace,
			Labels: map[string]string{
				"app":       podName,
				"lifespan":  strconv.Itoa(params.Timeout), // to be used by a cleanup pod
				"outputDir": params.OutputDir,
			},
		},
		Spec: core.PodSpec{
			Containers: []core.Container{
				{
					Name:            "pytest-pod",
					Image:           params.Image,
					ImagePullPolicy: core.PullAlways,
					Args:            params.Arguments,
				},
			},
			RestartPolicy: core.RestartPolicyNever,
		},
		// todo: mount outputdir to NFS {id}/
	}

	pod, err := clientset.CoreV1().Pods(namespace).Create(context.TODO(), pod, metav1.CreateOptions{})
	if err != nil {
		return err
	}

	log.Infof("Test pod %s created successfully.", podName)
	return nil
}

func GetPod(clientset *kubernetes.Clientset, namespace string, id string) (*core.Pod, error) {
	podName := "pytest-pod-" + id
	return clientset.CoreV1().Pods(namespace).Get(context.TODO(), podName, metav1.GetOptions{})
}

func GetPodLogs(clientset *kubernetes.Clientset, namespace string, id string) (string, error) {
	podName := "pytest-pod-" + id
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
