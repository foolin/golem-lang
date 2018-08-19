default: build

clean:
	rm -rf build
	rm -rf release

fmt:
	go fmt ./...

lint:
	gometalinter.v2 \
		--disable=gocyclo \
		--disable=goconst \
		--exclude="lib/(.*)is unused \(U1000\) \(megacheck\)" \
		./...

vet:
	go vet ./...

test:
	go test ./...

compile: 
	go build -o build/golem golem.go
	mkdir -p build/lib
	go build -buildmode=plugin -o build/lib/os/os.so         lib/os/os.go
	go build -buildmode=plugin -o build/lib/path/path.so     lib/path/path.go
	go build -buildmode=plugin -o build/lib/regexp/regexp.so lib/regexp/regexp.go
	cp lib/os/os.glm         build/lib/os/os.glm 
	cp lib/path/path.glm     build/lib/path/path.glm 
	cp lib/regexp/regexp.glm build/lib/regexp/regexp.glm 

bench_test: test compile
	build/golem bench_test/core_test.glm
	build/golem bench_test/os_test.glm
	build/golem bench_test/path_test.glm
	build/golem bench_test/regexp_test.glm

build: clean fmt lint vet test compile bench_test

# Cross-compiling plugins doesn't work. So, creating the release tarball for 
# a given platform must be done on a machine that runs that platform.
release:
	rm -rf release
	mkdir -p release/golem
	cp -r build/ release/golem
	cd release && tar czf golem-${PLATFORM}-${VERSION}.tar.gz golem

