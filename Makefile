
default: build

clean:
	rm -f golem

test:
	go test ./...

lint:
	golint ./...

build: clean test
	go build golem.go
