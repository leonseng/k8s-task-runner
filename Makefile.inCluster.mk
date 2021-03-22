IN_CLUSTER_IMAGE_REPO ?= k3d-$(TEST_REGISTRY):$(TEST_REGISTRY_PORT)
IN_CLUSTER_IMAGE_NAME ?= k8s_task_runner
IN_CLUSTER_IMAGE_VERSION ?= 0.5
IN_CLUSTER_IMAGE_TAG ?= $(IN_CLUSTER_IMAGE_REPO)/$(IN_CLUSTER_IMAGE_NAME):$(IN_CLUSTER_IMAGE_VERSION)

.PHONY: in-cluster-clean
in-cluster-clean:
	@ kubectl delete -f integration_tests/k3d_ingress.yaml || true
	@ kubectl delete -f integration_tests/k8s_task_runner.yaml || true

.PHONY: in-cluster-test
in-cluster-test: in-cluster-clean
	# build and load image into cluster
	@ go mod tidy
	@ docker build --tag $(IN_CLUSTER_IMAGE_TAG) ./
	@ docker push $(IN_CLUSTER_IMAGE_TAG)

	# setup
	@ kubectl apply -f integration_tests/k8s_task_runner.yaml
	@ kubectl wait --for=condition=available --timeout=60s deployments/k8s-task-runner
	@ kubectl apply -f integration_tests/k3d_ingress.yaml

	# run test
	@ K8S_TASK_RUNNER_ENDPOINT=http://localhost:8080 \
			go test ./integration_tests/
