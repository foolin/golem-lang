default: build

clean:
	rm -rf build

fmt:
	go fmt ./...

vet:
	go vet ./...

test:
	go test ./...

compile: 
	go build -o build/golem golem.go
	mkdir -p build/lib/encoding
	mkdir -p build/lib/os
	mkdir -p build/lib/path
	mkdir -p build/lib/regexp
	cp lib/encoding/encoding.glm build/lib/encoding/encoding.glm 
	cp lib/os/os.glm             build/lib/os/os.glm 
	cp lib/path/path.glm         build/lib/path/path.glm 
	cp lib/regexp/regexp.glm     build/lib/regexp/regexp.glm 

build: clean fmt vet test compile 
	cd bench_test && ../build/golem core_test.glm
	cd bench_test && ../build/golem encoding_test.glm
	cd bench_test && ../build/golem os_test.glm
	cd bench_test && ../build/golem path_test.glm
	cd bench_test && ../build/golem regexp_test.glm

lint:
	gometalinter.v2  --disable=gocyclo  --disable=goconst ./...

