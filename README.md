# Kubernetes Task Runner

Exposes a REST API to create a single-run pod on a Kubernetes cluster, and another API to retrieve the task status and logs.

## Requirement

The following binaries are required to run the tests:
- [kubectl](https://kubernetes.io/docs/tasks/tools/install-kubectl/)
- [k3d](https://k3d.io/#installation)
- [golangci-lint](https://golangci-lint.run/usage/install/#local-installation)

## Usage

```
$ k8s_task_runner --help
Usage of k8s_task_runner:
  -inCluster
        Toggle for running k8s-task-runner in a Kubernetes cluster
  -kubeconfig string
        absolute path to the kubeconfig file (default "/etc/k8s-task-runner/.kube/config")
  -port int
        Port to serve API on (default 80)
```

## Testing

> If you are behind a proxy, create a file called `.user-env` at the root directory of the project with following contents:
> ```
> HTTP_PROXY=<PROXY_URL>:<PROXY_PORT>  # e.g. HTTP_PROXY=http://10.0.2.15:3128
> ```

`k8s-task-runner` requires a Kubernetes cluster to interact with. You can spin up one by running
```
make k3d-setup
```
Once the test cluster has been set up, pick one of the following test scenarios:

- To run `k8s-task-runner` as a `go` binary external to the Kubernetes cluster, run
  ```
  make test
  ```
  which starts the API server on `localhost:8081` and executes the `go` tests against it
- To run `k8s-task-runner` as a Pod within the Kubernetes cluster, run
  ```
  make in-cluster-test
  ```
  which
  1. builds a new `k8s-task-runner` image and push it into the `k3d` Docker registry
  1. deploys the necessary Kubernetes objects to start serving `k8s-task-runner` on `localhost:8080`
  1. executes the `go` tests on `localhost:8080`

## API documentation

The OpenAPI 3.0 specification for `k8s-task-runner` can be found [here](./swaggerui/openapi.yaml). The Swagger UI can be accessed by `<host>/swaggerui/` endpoint of the API.

## Todo

- [x] Replace `k8s.io/api/core/v1.*` in [k8sclient](./k8sclient/k8sclient.go) with Kubernetes manifest YAML files
- [x] Add ability for users to pass in Docker credentials to pull images from private repos
- [X] Add basic health check API to detect if server is ready
- [X] Add ability to specify environment variables in task pods
- [ ] Some mechanism to clean up old task pods and secrets
- [ ] Allow app configuration via environment variables
- [ ] Improve API documentation. OpenAPI?
- [ ] Image versioning needs work. It's currently statically defined under `IMAGE_VERSION` in the [Makefile](./Makefile)
- [ ] Automatically pushing `k8s-task-runner` image to a remote Docker registry (Dockerhub?).
- [ ] Create Helm chart
