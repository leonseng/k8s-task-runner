package main

import (
	"flag"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"

	"github.com/leonseng/k8s_task_runner/api"

	log "github.com/sirupsen/logrus"
)

type userInputs struct {
	namespace  string
	port       int
	inCluster  bool
	kubeconfig string
}

func main() {
	log.SetLevel(log.DebugLevel)

	inputs := parseUserInputs()
	log.Infof("%+v\n", inputs)

	clientSet, err := getK8sClientSet(inputs)
	if err != nil {
		panic(err.Error())
	}

	api.HandleRequests(
		api.ApplicationConfiguration{
			Port:          inputs.port,
			K8sClientSet:  clientSet,
			TaskNamespace: inputs.namespace,
		},
	)
}

func parseUserInputs() userInputs {
	var port = flag.Int("port", 80, "Port to serve API on")
	var inCluster = flag.Bool("inCluster", false, "Toggle for running k8s-task-runner in a Kubernetes cluster")
	var kubeconfig = flag.String("kubeconfig", "/etc/k8s-task-runner/.kube/config", "absolute path to the kubeconfig file")

	flag.Parse()

	return userInputs{
		namespace:  "default",
		port:       *port,
		inCluster:  *inCluster,
		kubeconfig: *kubeconfig,
	}
}

func getK8sClientSet(input userInputs) (*kubernetes.Clientset, error) {
	var config *rest.Config
	var err error
	if input.inCluster {
		// use the current context in kubeconfig
		config, err = rest.InClusterConfig()
	} else {
		config, err = clientcmd.BuildConfigFromFlags("", input.kubeconfig)
	}
	if err != nil {
		return nil, err
	}

	// creates the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	return clientset, nil
}
