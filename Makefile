
default: build

clean:
	rm -f golem
	rm -rf ./release

fmt:
	go fmt ./...

lint:
	gometalinter.v2 --disable=gocyclo --disable=goconst ./...

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

#build: clean fmt lint vet test compile bench_test
build: clean fmt test compile bench_test

release: build
	mkdir -p ./release/golem/linux
	mkdir -p ./release/golem/mac
	mkdir -p ./release/golem/windows
	GOOS=linux   GOARCH=amd64 go build -o ./release/golem/linux/golem       golem.go
	GOOS=darwin  GOARCH=amd64 go build -o ./release/golem/mac/golem         golem.go
	GOOS=windows GOARCH=amd64 go build -o ./release/golem/windows/golem.exe golem.go
	cd release && zip -r ./golem-0.8.2.zip golem
	cd release && tar czf ./golem-0.8.2.tar.gz golem
