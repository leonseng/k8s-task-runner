package k8sclient

import (
	"bytes"
	"fmt"
	"text/template"

	log "github.com/sirupsen/logrus"
	v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes/scheme"
)

var podTemplate *template.Template

func init() {
	// define templates
	var err error
	podTemplate, err = template.New("podTemplate").Parse(testPodManifestTemplate)
	if err != nil {
		panic(err)
	}
	// do a test render here, check that it can be turned into a pod object
	_, err = manifestToPodObject(
		CreateParameters{
			ID:        "123",
			Namespace: "default",
			Image:     "busybox:1.28",
			Command:   []string{"date"},
			Arguments: []string{"--rfc-2822"},
		},
	)
	if err != nil {
		panic(err)
	}
}

func manifestToPodObject(params CreateParameters) (*v1.Pod, error) {
	var podManifest bytes.Buffer
	err := podTemplate.Execute(&podManifest, params)
	if err != nil {
		log.Errorf("Failed to render Pod manifest\n")
		return nil, err
	}

	// create k8s objects from YAML - https://github.com/kubernetes/client-go/issues/193
	obj, groupVersionKind, err := scheme.Codecs.UniversalDeserializer().Decode(podManifest.Bytes(), nil, nil)
	if err != nil {
		log.Errorf("Failed to decode Pod manifest into K8s Pod object\n")
		return nil, err
	}

	log.Debugf("%+v", obj.GetObjectKind())
	log.Debugf("%+v", groupVersionKind)
	log.Debugf("%+v", obj)

	if groupVersionKind.Kind != "Pod" {
		log.Errorf("Rendered manifest is not of type Pod.\n")
		return nil, fmt.Errorf("rendered manifest is not of type Pod")
	}

	return obj.(*v1.Pod), nil
}
