.PHONY: bootstrap
bootstrap:
	cp -v $(shell go env GOROOT)"/misc/wasm/wasm_exec.js" static/

.PHONY: thick
thick:
	GOOS=js GOARCH=wasm go build -o static/weather.wasm cmd/thick/main.go

.PHONY: thin
thin:
	GOOS=js GOARCH=wasm go build -o static/weather.wasm cmd/thin/main.go

.PHONY: console
console:
	go build -o console cmd/console/console.go

.PHONY: server
server:
	go build -o server cmd/server/server.go

.PHONY: clean
clean:weather
	rm -vf console server console.log static/weather.wasm

.PHONY: get
get:
	go get github.com/flosch/pongo2
	go get github.com/buger/jsonparser