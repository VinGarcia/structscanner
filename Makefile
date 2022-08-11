
GOBIN=$(shell go env GOPATH)/bin

test: setup
	$(GOBIN)/richgo test $(args) ./...

lint: setup
	go vet -structtag=false ./...
	$(GOBIN)/staticcheck ./...

setup: $(GOBIN)/richgo $(GOBIN)/staticcheck
$(GOBIN)/staticcheck:
	go install honnef.co/go/tools/cmd/staticcheck@latest
$(GOBIN)/richgo:
	go install github.com/kyoh86/richgo@latest

