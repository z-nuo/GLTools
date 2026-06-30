.PHONY: fmt vet test check

GO_FILES := $(shell find . -name '*.go' -not -path './vendor/*')

fmt:
	@unformatted="$$(gofmt -l $(GO_FILES))"; \
	if [ -n "$$unformatted" ]; then \
		echo "$$unformatted"; \
		exit 1; \
	fi

vet:
	go vet ./...

test:
	go test ./...

check: fmt vet test
