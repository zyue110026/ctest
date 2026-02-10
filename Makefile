# ---------------------------------------
# Kubernetes root (two levels up from this Makefile)
# ---------------------------------------
# This sets K8S_ROOT to the absolute path two levels above this Makefile.
# It is used as the working directory for all go test commands.
K8S_ROOT ?= $(abspath $(dir $(lastword $(MAKEFILE_LIST)))/../..)

# ---------------------------------------
# Packages
# ---------------------------------------
TEST_PKG := ./test/ctest                 # Package containing test fixture generation
TEST_REWRITE_PKG := ./test/ctest/test_rewrite  # Package containing rewrite test

ETCD_BIN := $(K8S_ROOT)/third_party/etcd/etcd
ETCD_DIR := $(K8S_ROOT)/third_party/etcd





# ---------------------------------------
# Optional inputs (can override from command line)
# ---------------------------------------
REPO_PATH ?=                             # Path to the repository for generating fixtures
REWRITE_TARGET ?= test/e2e               # Target directory or file to rewrite
OLLAMA_MODEL ?= gpt-oss:120b-cloud       # Ollama model to use for rewriting
OVERWRITE_REWRITTEN ?= false             # Whether to overwrite already rewritten files (true/false)

# ---------------------------------------
# Help
# ---------------------------------------
.PHONY: help
help:
	@echo "Usage:"
	@echo "  make gen-fixtures REPO_PATH=/path/to/repo"
	@echo "    Generate test fixtures for the specified repository."
	@echo ""
	@echo "  make testrewrite [REWRITE_TARGET=test/e2e] [OLLAMA_MODEL=deepseek-coder:33b] [OVERWRITE_REWRITTEN=false]"
	@echo "    Rewrite Go test files using Ollama. Optional environment variables:"
	@echo "      REWRITE_TARGET       Directory or file to rewrite (default: test/e2e)"
	@echo "      OLLAMA_MODEL         Ollama model to use (default: deepseek-coder:33b)"
	@echo "      OVERWRITE_REWRITTEN  Whether to overwrite already rewritten files (default: false)"
	@echo ""
	@echo "  make test-integration"
	@echo "    Run Kubernetes integration tests with etcd setup."
	@echo "    Logs output to test/ctest/logs/ctest_integration_logs_YYYYMMDDTHHMMSS.html."
	@echo ""
	@echo "  make ctest-e2e"
	@echo "    Run end-to-end CTest tests using kubetest2 + kind."
	@echo "    Automatically checks and installs Go, kubetest2, and kind if missing."
	@echo "    Uses ginkgo focus file: ctest and built binaries."
	@echo ""
	@echo "  make ctest-unit"
	@echo "    Run unit tests with names prefixed by TestCtest, excluding test/ folder."
	@echo "    Logs output to test/ctest/logs/ctest_unit_logs_YYYYMMDDTHHMMSS.html."

# ---------------------------------------
# Generate Fixtures
# ---------------------------------------
.PHONY: gen-fixtures
gen-fixtures:
ifndef REPO_PATH
	# Fail if REPO_PATH is not specified
	$(error REPO_PATH is not set. Usage: make gen-fixtures REPO_PATH=/path/to/repo)
endif
	@echo "📁 Scanning repo: $(REPO_PATH)"
	# Run the Go test that generates fixtures
	cd $(K8S_ROOT) && \
	go test $(TEST_PKG) -run TestGenerateFixtures -v -repo=$(REPO_PATH)

# ---------------------------------------
# Rewrite Tests Using Ollama
# ---------------------------------------
.PHONY: testrewrite
testrewrite:
	@echo "✏️  Rewriting tests under: $(REWRITE_TARGET)"
	@echo "🧠 Using Ollama model: $(OLLAMA_MODEL)"
	@echo "⚡ Overwrite rewritten files: $(OVERWRITE_REWRITTEN)"
	# Set environment variables and run the rewrite test
	cd $(K8S_ROOT) && \
	REWRITE_TARGET=$(REWRITE_TARGET) \
	OLLAMA_MODEL=$(OLLAMA_MODEL) \
	OVERWRITE_REWRITTEN=$(OVERWRITE_REWRITTEN) \
	go test -timeout 24h $(TEST_REWRITE_PKG) -run TestRewriteWithLLM -v



# ---------------------------------------
# Integration Test with etcd check and logs
# ---------------------------------------
.PHONY: ctest-integration
ctest-integration:
	@mkdir -p $(K8S_ROOT)/test/ctest/logs
	@LOG_FILE=$(K8S_ROOT)/test/ctest/logs/ctest_integration_logs_$(date +%Y%m%dT%H%M%S).html; \
	echo "📂 Entering Kubernetes root: $(K8S_ROOT)"; \
	cd $(K8S_ROOT) && \
	echo "🔍 Checking for etcd at $(ETCD_BIN)..." && \
	if [ ! -x "$(ETCD_BIN)" ]; then \
		echo "⚠️  etcd not found. Installing..."; \
		./hack/install-etcd.sh; \
		echo "✅ etcd installed."; \
		echo 'export PATH="$$PATH:$(K8S_ROOT)/third_party/etcd"' >> ~/.profile; \
		export PATH="$$PATH:$(K8S_ROOT)/third_party/etcd"; \
	else \
		echo "✅ etcd is already installed."; \
	fi && \
	echo "🏃 Running integration tests (prefix TestCtest)..." && \
	make test-integration \
		GOFLAGS=-v \
		KUBE_COVER=y \
		KUBE_TEST_ARGS="-run ^TestCtest" \
		2>&1 | tee $$LOG_FILE


# ---------------------------------------
# CTest E2E using kubetest2 + kind
# ---------------------------------------
.PHONY: ctest-e2e
ctest-e2e:
	@echo "🔍 Checking Go..."
	@if ! command -v go >/dev/null 2>&1; then \
		echo "❌ Go is not installed or not in PATH."; \
		echo "   Please install Go first: https://go.dev/dl/"; \
		exit 1; \
	else \
		echo "✅ Go found."; \
	fi

	@echo "🔍 Checking kubetest2..."
	@if ! command -v kubetest2 >/dev/null 2>&1; then \
		echo "⚠️  kubetest2 not found. Installing..."; \
		go install sigs.k8s.io/kubetest2/...@latest; \
		echo "✅ kubetest2 installed."; \
	else \
		echo "✅ kubetest2 already installed."; \
	fi

	@echo "🔍 Checking kind..."
	@if ! command -v kind >/dev/null 2>&1; then \
		echo "⚠️  kind not found. Installing..."; \
		go install sigs.k8s.io/kind@latest; \
		echo "✅ kind installed."; \
	else \
		echo "✅ kind already installed."; \
	fi

	@echo "🏃 Running CTest E2E (ginkgo focus: ctest)..."
	kubetest2 kind --build --up --down --test ginkgo -v 4 -- \
		--test-args="--ginkgo.focus-file=ctest" \
		--use-built-binaries

# ---------------------------------------
# Run Unit Tests (prefix TestCtest, exclude test/ folder)
# ---------------------------------------
.PHONY: ctest-unit
ctest-unit:
	@mkdir -p $(K8S_ROOT)/test/ctest/logs
	@LOG_FILE=$(K8S_ROOT)/test/ctest/logs/ctest_unit_logs_$(shell date +%Y%m%dT%H%M%S).html; \
	echo "📂 Entering Kubernetes root: $(K8S_ROOT)"; \
	cd $(K8S_ROOT) && \
	echo "🏃 Running unit tests (prefix TestCtest, excluding test/)..."; \
	PKGS=$$(go list ./... | grep -v '^k8s.io/kubernetes/test/') && \
	go test -timeout 24h $$PKGS -run '^TestCtest' -v 2>&1 | tee $$LOG_FILE
