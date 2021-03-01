package k8sclient

type podManifestData struct {
	namespace string
	podName   string
	image     string
	command   []string
	arguments []string
}

var testPodManifestTemplate = `
apiVersion: v1
kind: Pod
metadata:
  labels:
    run: task-pod-{{ .ID }}
  name: task-pod-{{ .ID }}
  namespace: {{ .Namespace }}
spec:
  containers:
  - image: {{ .Image }}
    name: test
    resources: {}
    command: [
      {{- range .Command}}
      "{{.}}",
      {{- end }}
    ]
    args: [
      {{- range .Arguments}}
      "{{.}}",
      {{- end}}
    ]
    imagePullPolicy: Always
  restartPolicy: Never
`
