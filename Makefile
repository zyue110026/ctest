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
