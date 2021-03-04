-include .user-env
PROJECT_NAME = k8s-task-runner

TEST_CLUSTER = $(PROJECT_NAME)
TEST_REGISTRY = $(PROJECT_NAME).registry.localhost
TEST_REGISTRY_PORT = 5000
ifdef HTTP_PROXY
K3D_PROXY_VARS := -e "http_proxy=$(HTTP_PROXY)@server[0]" \
								-e "https_proxy=$(HTTP_PROXY)@server[0]" \
								-e "no_proxy=k3d-$(TEST_REGISTRY)@server[0]"
endif

IMAGE_REPO ?= k3d-$(TEST_REGISTRY):$(TEST_REGISTRY_PORT)
IMAGE_NAME ?= k8s_task_runner
IMAGE_VERSION ?= 0.5
IMAGE_TAG ?= $(IMAGE_REPO)/$(IMAGE_NAME):$(IMAGE_VERSION)

.PHONY: lint
lint:
	@ golangci-lint run

.PHONY: clean
clean: k3d-teardown

################
# Manage image #
################
.PHONY: image-build
image-build:
	go mod tidy
	docker build --tag $(IMAGE_TAG) ./
	docker push $(IMAGE_TAG)

#######################
# Manage test objects #
#######################
.PHONY: test-out-of-cluster-setup
test-out-of-cluster-setup: test-out-of-cluster-clean
	@ go run main.go --external --port 8081 > .log 2>&1 &

.PHONY: test-out-of-cluster-clean
test-out-of-cluster-clean:
	@ kill $$(ps -ef | grep "[ ]--external --port 8081" | awk '{print $$2}') 2> /dev/null || true

.PHONY: test-out-of-cluster
test-out-of-cluster:
	@ go test -run TestOutOfCluster ./integration_tests/

.PHONY: test-in-cluster-setup
test-in-cluster-setup: test-in-cluster-clean
	@ kubectl apply -f integration_tests/k8s_task_runner.yaml
	@ kubectl wait --for=condition=available --timeout=60s deployments/k8s-task-runner
	@ sleep 5

.PHONY: test-in-cluster-clean
test-in-cluster-clean:
	@ kubectl delete -f integration_tests/k8s_task_runner.yaml || true

.PHONY: test-in-cluster
test-in-cluster:
	@ go test -run TestInCluster ./integration_tests/

###########################
# Manage test environment #
###########################
.PHONY: k3d-teardown
k3d-teardown:
	@ if k3d cluster list $(TEST_CLUSTER); then \
			k3d cluster delete $(TEST_CLUSTER); \
		fi
	@ if k3d registry list k3d-$(TEST_REGISTRY); then \
			k3d registry delete k3d-$(TEST_REGISTRY); \
		fi

.PHONY: k3d-setup
k3d-setup:
	@ if ! k3d registry list k3d-$(TEST_REGISTRY); then \
			k3d registry create $(TEST_REGISTRY) --port $(TEST_REGISTRY_PORT); \
		fi
ifdef HTTP_PROXY
	@ echo "Using proxy: $(HTTP_PROXY)"
endif
	@ if ! k3d cluster list $(TEST_CLUSTER); then \
			k3d cluster create $(TEST_CLUSTER) \
				$(K3D_PROXY_VARS) \
				--registry-use k3d-$(TEST_REGISTRY):$(TEST_REGISTRY_PORT) \
				-p "8080:80@loadbalancer" \
				; \
		fi
	@ kubectl config use-context k3d-$(TEST_CLUSTER)
