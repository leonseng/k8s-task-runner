package main

import (
	"flag"
	"path/filepath"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"

	"github.com/leonseng/k8s_task_runner/api"

	log "github.com/sirupsen/logrus"
)

type userInputs struct {
	namespace    string
	port         int
	outOfCluster bool
	kubeconfig   string
}

func main() {
	log.SetLevel(log.DebugLevel)

	inputs := parseUserInputs()
	log.Infof("%+v\n", inputs)

	clientSet, err := getK8sClientSet(inputs.outOfCluster, inputs.kubeconfig)
	if err != nil {
		panic(err.Error())
	}

	api.HandleRequests(clientSet, inputs.namespace, inputs.port)
}

func parseUserInputs() userInputs {
	var port = flag.Int("port", 80, "Port to serve API on")
	var outOfCluster = flag.Bool("external", false, "Toggle for running k8s-task-runner out of a Kubernetes cluster")

	var kubeconfig *string
	if home := homedir.HomeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}

	flag.Parse()

	return userInputs{
		namespace:    "default",
		port:         *port,
		outOfCluster: *outOfCluster,
		kubeconfig:   *kubeconfig,
	}
}

func getK8sClientSet(outOfCluster bool, kubeconfig string) (*kubernetes.Clientset, error) {
	var config *rest.Config
	var err error
	if outOfCluster {
		config, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
	} else {
		// use the current context in kubeconfig
		config, err = rest.InClusterConfig()
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
