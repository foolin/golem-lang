
default: build

clean:
	rm -f golem

fmt:
	go fmt ./...

lint:
	golint ./...
	gometalinter --disable=gocyclo --disable=goconst ./...

vet:
	go vet ./...

test:
	go test ./...

compile: 
	go build golem.go

bench_test: compile
	./golem bench_test/os_test.glm
	./golem bench_test/regexp_test.glm
	./golem bench_test/path_test.glm

build: clean fmt lint vet test compile bench_test
	go build golem.go
