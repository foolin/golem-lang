
default: build

clean:
	rm -f golem

fmt:
	go fmt ./...

lint:
	gometalinter --disable=gocyclo --disable=goconst ./...

vet:
	go vet ./...

test:
	go test ./...

compile: 
	go build golem.go

bench_test: compile
	./golem bench_test/core_test.glm
	./golem bench_test/os_test.glm
	./golem bench_test/regexp_test.glm
	./golem bench_test/path_test.glm

build: clean fmt lint vet test compile bench_test
	go build golem.go

release: build
	rm -rf ./release
	mkdir -p ./release/linux
	mkdir -p ./release/mac
	mkdir -p ./release/windows

	GOOS=linux   GOARCH=amd64 go build -o ./release/linux/golem       golem.go
	GOOS=darwin  GOARCH=amd64 go build -o ./release/mac/golem         golem.go
	GOOS=windows GOARCH=amd64 go build -o ./release/windows/golem.exe golem.go

# To tag a release:
#     git tag -a v0.8.0 -m "version 0.8.0"
#     git push origin v0.8.0
