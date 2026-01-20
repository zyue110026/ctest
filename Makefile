# Kubernetes root (two levels up from this Makefile)
K8S_ROOT ?= $(abspath $(dir $(lastword $(MAKEFILE_LIST)))/../..)

TEST_PKG := ./test/ctest
TEST_REWRITE_PKG := ./test/ctest/test_rewrite
REPO_PATH ?=
# Defaults (can be overridden by command line)
REWRITE_TARGET ?= test/e2e
OLLAMA_MODEL ?= deepseek-coder:33b
.PHONY: help
help:
	@echo "Usage:"
	@echo "  make gen-fixtures REPO_PATH=/path/to/repo"
	@echo "  make testrewrite [REWRITE_TARGET=test/e2e] [OLLAMA_MODEL=deepseek-coder:33b]"


.PHONY: gen-fixtures
gen-fixtures:
ifndef REPO_PATH
	$(error REPO_PATH is not set. Usage: make gen-fixtures REPO_PATH=/path/to/repo)
endif
	@echo "📁 Scanning repo: $(REPO_PATH)"
	cd $(K8S_ROOT) && \
	go test $(TEST_PKG) -run TestGenerateFixtures -v -repo=$(REPO_PATH)

.PHONY: testrewrite
testrewrite:
	@echo "✏️  Rewriting tests under: $(REWRITE_TARGET)"
	@echo "🧠 Using Ollama model: $(OLLAMA_MODEL)"
	cd $(K8S_ROOT) && \
	REWRITE_TARGET=$(REWRITE_TARGET) \
	OLLAMA_MODEL=$(OLLAMA_MODEL) \
	go test $(TEST_REWRITE_PKG) -run TestRewriteWithLLM -v
