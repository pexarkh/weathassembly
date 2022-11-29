.PHONY: info
info:
	@echo "TARGETS:"
	@echo "  bootstrap"
	@echo "  thick"
	@echo "  thin"
	@echo "  console"
	@echo "  server"
	@echo "  clean"

.PHONY: bootstrap
bootstrap:
	go mod tidy
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
