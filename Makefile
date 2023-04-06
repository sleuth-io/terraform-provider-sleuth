.PHONY: docs release help format

# Help system from https://marmelab.com/blog/2016/02/29/auto-documented-makefile.html
.DEFAULT_GOAL := help
TEST?=$$(go list ./... | grep -v 'vendor')
HOSTNAME=sleuth.io
NAMESPACE=core
NAME=sleuth
BINARY=terraform-provider-${NAME}
VERSION=0.3.0-dev
OS_ARCH=linux_amd64

# go source files, ignore vendor directory
SRC = $(shell find . -type f -name '*.go' -not -path "./vendor/*")

help:
	@grep -E '^[a-zA-Z0-9_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

build: ## Build the binary
	go build -o ${BINARY}

format: ## Format the source code with gofmt
	@gofmt -l -w $(SRC)

release: ## Releases the current version as a snapshot
	goreleaser release --rm-dist --snapshot --skip-publish  --skip-sign

install: build ## Builds and installs locally
	mkdir -p ~/.terraform.d/plugins/${HOSTNAME}/${NAMESPACE}/${NAME}/${VERSION}/${OS_ARCH}
	mv ${BINARY} ~/.terraform.d/plugins/${HOSTNAME}/${NAMESPACE}/${NAME}/${VERSION}/${OS_ARCH}

test: ## Runs the tests
	go test -i $(TEST) || exit 1                                                   
	echo $(TEST) | xargs -t -n4 go test $(TESTARGS) -timeout=30s -parallel=4                    

docs: ## Generates docs
	go run github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs

testacc: ## Runs acceptance tests
	TF_ACC=1 go test $(TEST) -v $(TESTARGS) -timeout 120m

dev: ## Runs terraform against your local dev env
	test -s main.tf || (echo "**** Set up main.tf first from main.tf.example *** "; exit 1)
	rm -f terraform.tfstate && terraform plan && terraform apply

