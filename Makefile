.PHONY: format install-hooks lint openapi test build run docker-image compose

# Format code with ngofumpt (https://github.com/mvdan/gofumpt).
format:
	gofumpt -w ./internal

# Add format command to local git pre-commit hook
install-hooks:
	# check gofumpt
	if ! command -v gofumpt &> /dev/null; then \
  		echo "gofumpt not found, install started..."; \
  		go install mvdan.cc/gofumpt@latest; \
    fi

	echo '#!/bin/sh\nmake format\nif [ $$? -ne 0 ]; then\n    echo "Error: make format failed. Commit aborted."\n    exit 1\nfi\nexit 0' > .git/hooks/pre-commit
	chmod +x .git/hooks/pre-commit

# Run golangci-lint (https://github.com/golangci/golangci-lint).
lint:
	golangci-lint run ./internal/...

# Generate OpenAPI v2 spec (https://github.com/swaggo/swag).
openapi:
	swag init -g api.go --parseDependency --parseInternal --dir ./internal/api --output ./openapi

# Run unit tests.
test:
	go test -v -cover ./internal/...

# Generate OpenAPI spec and build app.
build: format openapi
	CGO_ENABLED=0 cd ./cmd && go build -o ugc

# Build and run app with config file provided.
# Example:
# `make run args="-c /home/user/template_config.yaml"`
run: build
	cd ./cmd && ./ugc $(args)

docker-image:

compose:
