default: docs

build: 
	GOARCH=wasm GOOS=js go build -o ../../docs/golem.wasm wasm.go

docs: 
	../../build/golem makeDocs.glm
