# Kubernetes Task Runner

Exposes a REST API to create a single-run pod on a Kubernetes cluster, and another API to retrieve the task status and logs.

## Requirement

The following binaries are required to run the tests:
- [kubectl](https://kubernetes.io/docs/tasks/tools/install-kubectl/)
- [k3d](https://k3d.io/#installation)
- [golangci-lint](https://golangci-lint.run/usage/install/#local-installation)

## Testing

> If you are behind a proxy, create a file called `.user-env` at the root directory of the project with following contents:
> ```
> HTTP_PROXY=<PROXY_URL>:<PROXY_PORT>  # e.g. HTTP_PROXY=http://10.0.2.15:3128
> ```

`k8s-task-runner` requires a Kubernetes cluster to interact with. You can spin up one by running `make k3d-setup`. Once the test cluster has been set up, pick one of the following test scenarios:

- Out of cluster

  To run `k8s-task-runner` as a `go` binary external to the Kubernetes cluster, run
  1. `make test-out-of-cluster-setup` to start the program, which should be listening on `localhost:8081`.
  1. `go test -run TestOutOfCluster ./integration_tests/` to test `k8s-task-runner` on `localhost:8081`
- In cluster

  To run `k8s-task-runner` as a Pod within the Kubernetes cluster, run the following in sequence:
  1. `make image-build` to build a new `k8s-task-runner` image and push it into the test Docker registry
  1. `make test-in-cluster-setup` to deploy the necessary Kubernetes objects to start serving `k8s-task-runner` on `localhost:8080`
  1. `go test -run TestInCluster ./integration_tests/` to test `k8s-task-runner` on `localhost:8080`

## Todo

- [ ] Replace `k8s.io/api/core/v1.*` in [k8sclient](./k8sclient/k8sclient.go) with Kubernetes manifest YAML files
- [ ] Image versioning needs work. It's currently statically defined under `IMAGE_VERSION` in the [Makefile](./Makefile)
- [ ] Automatically pushing `k8s-task-runner` image to a remote Docker registry (Dockerhub?).
- [ ] Create Helm chart
