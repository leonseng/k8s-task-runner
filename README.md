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

To setup for testing, run `make init`, which performs the following:
- Spin up a test Kubernetes cluster using k3d (Or run `make k3d-setup`)
- Build and push `k8s-task-runner` image into the k3d docker repository (Or run `make image-build`)

To test, run `make test`.


## Todo

- [ ] Replace `k8s.io/api/core/v1.*` in [k8sclient](./k8sclient/k8sclient.go) with Kubernetes manifest YAML files
- [ ] Image versioning needs work. It's currently statically defined under `IMAGE_VERSION` in the [Makefile](./Makefile)
- [ ] Automatically pushing `k8s-task-runner` image to a remote Docker registry (Dockerhub?).
- [ ] Create Helm chart
