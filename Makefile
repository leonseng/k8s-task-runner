IMAGE_REPO ?= k3d-registry.localhost:5000
IMAGE_NAME ?= go_pytest_runner
IMAGE_VERSION ?= 0.5
IMAGE_TAG ?= $(IMAGE_REPO)/$(IMAGE_NAME):$(IMAGE_VERSION)

.PHONY: build
build:
	go mod tidy
	docker build --tag $(IMAGE_TAG) ./
	docker push $(IMAGE_TAG)

.PHONY: deploy
deploy:
	kubectl apply -f go_pytest_runner.yaml