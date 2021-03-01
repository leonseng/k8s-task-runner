package k8sclient

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"text/template"

	log "github.com/sirupsen/logrus"
	core "k8s.io/api/core/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
)

type CreateParameters struct {
	ID        string
	Namespace string
	Image     string
	Command   []string
	Arguments []string
}

var podTemplate *template.Template

func init() {
	// define templates
	var err error
	podTemplate, err = template.New("podTemplate").Parse(testPodManifestTemplate)
	if err != nil {
		panic(err)
	}
	// do a test render here, check that it can be turned into a pod object
}

func CreatePodFromManifest(clientset *kubernetes.Clientset, params CreateParameters) error {
	// render go templates and store output into a variable - https://coderwall.com/p/ns60fq/simply-output-go-html-template-execution-to-strings
	userError := func() error {
		return fmt.Errorf("failed to create task pod. See error logs")
	}

	var podManifest bytes.Buffer
	err := podTemplate.Execute(&podManifest, params)
	if err != nil {
		log.Errorf("Failed to render Pod manifest:\n%v\n", err)
		return userError()
	}

	// create k8s objects from YAML - https://github.com/kubernetes/client-go/issues/193
	obj, groupVersionKind, err := scheme.Codecs.UniversalDeserializer().Decode([]byte(podManifest.String()), nil, nil)
	if err != nil {
		log.Errorf("Failed to decode Pod manifest into K8s Pod object:\n%v\n", err)
		return userError()
	}

	log.Debugf("%+v", obj.GetObjectKind())
	log.Debugf("%+v", groupVersionKind)
	log.Debugf("%+v", obj)

	pod := obj.(*v1.Pod)
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
