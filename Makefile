-include .user-env
include Makefile.k3d.mk
include Makefile.inCluster.mk

PROJECT_NAME = k8s-task-runner
TEST_KUBECONFIG ?= ~/.kube/config
TEST_PORT ?= 8081

.PHONY: lint
lint:
	@ golangci-lint run

.PHONY: clean
clean: test-clean k3d-teardown

.PHONY: test-clean
test-clean:
	@ kill $$(ps -ef | grep "[ ]-port $(TEST_PORT)" | awk '{print $$2}') 2> /dev/null || true

.PHONY: test
test: test-clean
	@ go mod tidy
	@ go install
	@ k8s_task_runner -port $(TEST_PORT) -kubeconfig $(TEST_KUBECONFIG) > .log 2>&1 &
	@ K8S_TASK_RUNNER_ENDPOINT=http://localhost:$(TEST_PORT) \
			go test ./integration_tests/
