TEST_CLUSTER = $(PROJECT_NAME)
TEST_REGISTRY = $(PROJECT_NAME).registry.localhost
TEST_REGISTRY_PORT = 5000
ifdef HTTP_PROXY
K3D_PROXY_VARS := -e "http_proxy=$(HTTP_PROXY)@server[0]" \
								-e "https_proxy=$(HTTP_PROXY)@server[0]" \
								-e "no_proxy=k3d-$(TEST_REGISTRY)@server[0]"
endif

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
