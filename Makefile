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
	cd bench_test && ../build/golem core_test.glm
	cd bench_test && ../build/golem os_test.glm
	cd bench_test && ../build/golem path_test.glm
	cd bench_test && ../build/golem regexp_test.glm

build: clean fmt lint vet test compile bench_test
