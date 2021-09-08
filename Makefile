TEST?=$$(go list ./... | grep -v 'vendor')
HOSTNAME=sleuth.io
NAMESPACE=core
NAME=sleuth
BINARY=terraform-provider-${NAME}
VERSION=0.1-dev
OS_ARCH=linux_amd64

# go source files, ignore vendor directory
SRC = $(shell find . -type f -name '*.go' -not -path "./vendor/*")

default: install

build:
	go build -o ${BINARY}

format:
	@gofmt -l -w $(SRC)

release:
	goreleaser release --rm-dist --snapshot --skip-publish  --skip-sign

install: build
	mkdir -p ~/.terraform.d/plugins/${HOSTNAME}/${NAMESPACE}/${NAME}/${VERSION}/${OS_ARCH}
	mv ${BINARY} ~/.terraform.d/plugins/${HOSTNAME}/${NAMESPACE}/${NAME}/${VERSION}/${OS_ARCH}

test: 
	go test -i $(TEST) || exit 1                                                   
	echo $(TEST) | xargs -t -n4 go test $(TESTARGS) -timeout=30s -parallel=4                    


testacc: 
	TF_ACC=1 go test $(TEST) -v $(TESTARGS) -timeout 120m   
