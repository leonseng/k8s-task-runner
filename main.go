package main

import (
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"

	"github.com/leonseng/go_pytest_runner/api"

	log "github.com/sirupsen/logrus"
)

var PytestRunnerNamespace = "default"
var PytestRunnerPort = 80

func main() {
	log.SetLevel(log.DebugLevel)

	config, err := rest.InClusterConfig()
	if err != nil {
		panic(err.Error())
	}

	// creates the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	api.HandleRequests(clientset, PytestRunnerNamespace, PytestRunnerPort)
}
