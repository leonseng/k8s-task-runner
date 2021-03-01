package k8sclient

import (
	"bytes"
	"context"
	"fmt"
	"io"

	log "github.com/sirupsen/logrus"
	core "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type CreateParameters struct {
	ID        string
	Namespace string
	Image     string
	Command   []string
	Arguments []string
}

func CreatePodFromManifest(clientset *kubernetes.Clientset, params CreateParameters) error {
	// render go templates and store output into a variable - https://coderwall.com/p/ns60fq/simply-output-go-html-template-execution-to-strings
	userError := func() error {
		return fmt.Errorf("failed to create task pod. See error logs")
	}

	pod, err := manifestToPodObject(params)
	if err != nil {
		log.Errorf(err.Error())
		return userError()
	}

	pod, err = clientset.CoreV1().Pods(params.Namespace).Create(context.TODO(), pod, metav1.CreateOptions{})
	if err != nil {
		log.Errorf("Failed to create pod:\n%v\n", err)
		return userError()
	}

	log.Infof("Test pod %s created successfully.", pod.Name)
	return nil
}

func CreatePod(clientset *kubernetes.Clientset, params CreateParameters) error {
	podName := "task-pod-" + params.ID
	pod := &core.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      podName,
			Namespace: params.Namespace,
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

	pod, err := clientset.CoreV1().Pods(params.Namespace).Create(context.TODO(), pod, metav1.CreateOptions{})
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
