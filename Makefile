TARGET=trlock
GOBIN=$(GOPATH)/bin
GOOS=$(shell uname | tr A-Z a-z)

tools:
	@echo $(GOBIN)
	go install golang.org/x/tools/cmd/goimports@v0.8.0
	go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.4.0
	go get github.com/hekmon/transmissionrpc/v3@v3.0.0
	go get github.com/sirupsen/logrus@v1.9.3

.PHONY: fmt
fmt:
	go fmt ./...
.PHONY: lint
lint:
	$(GOBIN)/golangci-lint run --fix
.PHONY: import
import:
	$(GOBIN)/goimports -l -w .

clean:
	rm $(TARGET)

.PHONY: build
build:
	@echo $(GOFILES) $(GOOS)
	GOOS=$(GOOS) GOARCH=amd64 go build -o $(TARGET) .

.PHONY: run
run:
	go run main.go

.PHONY: all
all: fmt lint import build
