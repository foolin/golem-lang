default: build

clean:
	rm -rf build

fmt:
	go fmt ./...

vet:
	go vet -shadow ./...

test:
	go test ./...

compile: 
	go build -o build/golem cli/golem.go

build: clean fmt vet test compile
	cd bench_test && ../build/golem benchTest.glm

lint:
	gometalinter.v2  --disable=gocyclo  --disable=goconst ./...
