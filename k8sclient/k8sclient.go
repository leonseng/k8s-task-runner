package k8sclient

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"io"

	log "github.com/sirupsen/logrus"
	core "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type SecretParameters struct {
	ID        string
	Namespace string
	Server    string
	Username  string
	Email     string
	Password  string
}

func CreateDockerRegistrySecret(clientset *kubernetes.Clientset, params SecretParameters) (string, error) {
	secretName := "task-secret-" + params.ID
	b64Auth := base64.StdEncoding.EncodeToString([]byte(params.Username + ":" + params.Password))

	secret := &core.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      secretName,
			Namespace: params.Namespace,
			Labels: map[string]string{
				"id": params.ID,
			},
		},
		StringData: map[string]string{
			".dockerconfigjson": fmt.Sprintf(
				`{"auths":{"%s":{"username":"%s","password":"%s","email":"%s","auth":"%s"}}}`,
				params.Server, params.Username, params.Password, params.Email, b64Auth,
			),
		},
	}

	_, err := clientset.CoreV1().Secrets(params.Namespace).Create(context.TODO(), secret, metav1.CreateOptions{})
	if err != nil {
		log.Errorf("Failed to create Secret %s:\n%s\n", secretName, err.Error())
		return "", err
	}

	log.Infof("Secret %s created successfully.", secretName)
	return secretName, nil
}

type PodParameters struct {
	ID        string
	Namespace string
	Secret    string
	Image     string
	Command   []string
	Arguments []string
	EnvVars   map[string]string
}

func CreatePodFromManifest(clientset *kubernetes.Clientset, params PodParameters) error {
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
