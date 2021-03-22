-include .user-env
include Makefile.k3d.mk
include Makefile.inCluster.mk

PROJECT_NAME = k8s-task-runner
TEST_KUBECONFIG ?= ~/.kube/config

.PHONY: lint
lint:
	@ golangci-lint run

.PHONY: clean
clean: test-clean k3d-teardown

.PHONY: test-clean
test-clean:
	@ kill $$(ps -ef | grep "[ ]-port 8081" | awk '{print $$2}') 2> /dev/null || true

.PHONY: test
test: test-clean
	@ go mod tidy
	@ go install
	@ k8s_task_runner -port 8081 -kubeconfig $(TEST_KUBECONFIG) > .log 2>&1 &
	@ go test -run TestOutOfCluster ./integration_tests/
