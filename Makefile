
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
	mkdir -p ./release/golem/linux
	mkdir -p ./release/golem/mac
	mkdir -p ./release/golem/windows
	GOOS=linux   GOARCH=amd64 go build -o ./release/golem/linux/golem       golem.go
	GOOS=darwin  GOARCH=amd64 go build -o ./release/golem/mac/golem         golem.go
	GOOS=windows GOARCH=amd64 go build -o ./release/golem/windows/golem.exe golem.go
	cd release && zip -r ./golem.zip golem
	cd release && tar czf ./golem.tar.gz golem
