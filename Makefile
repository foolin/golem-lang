
default: build

test:
	go test ./...

lint:
	golint ./...

clean:
	rm -f golem

build: test clean
	go build golem.go
