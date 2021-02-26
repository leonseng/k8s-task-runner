# Kubernetes Task Runner

Exposes a REST API to create a single-run pod on a Kubernetes cluster, and another API to retrieve the task status and logs.

## Requirement

The following binaries are required to run the tests:
- [kubectl](https://kubernetes.io/docs/tasks/tools/install-kubectl/)
- [k3d](https://k3d.io/#installation)
- [golangci-lint](https://golangci-lint.run/usage/install/#local-installation)

## Instructions

To setup for testing, run `make init`, which performs the following:
- Spin up a test Kubernetes cluster using k3d (Or run `make k3d-setup`)
- Build and push `k8s-task-runner` image into the k3d docker repository (Or run `make image-build`)

To test, run `make test`.

