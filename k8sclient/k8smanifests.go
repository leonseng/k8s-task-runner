package k8sclient

var testPodManifestTemplate = `
apiVersion: v1
kind: Pod
metadata:
  labels:
    run: task-pod
    id: "{{.ID}}"
  name: task-pod-{{.ID}}
  namespace: {{.Namespace}}
spec:
  {{- if .Secret}}
  imagePullSecrets:
  - name: {{.Secret}}
  {{- end}}
  containers:
  - image: {{.Image}}
    name: task-pod
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
    {{- if .EnvVars}}
    env:
      {{- range $key, $value := .EnvVars}}
    - name: "{{$key}}"
      value: "{{$value}}"
      {{- end}}
    {{- end}}
    imagePullPolicy: Always
  restartPolicy: Never
`
