
default: build

clean:
	rm -f golem

fmt:
	go fmt ./...

lint:
	golint ./...
	gometalinter ./...

vet:
	go vet ./...

test:
	go test ./...

build: clean fmt lint vet test
	go build golem.go
