# Weather Service with WebAssembly and Go

## Prereqs

```
go mod tidy
make bootstrap # to refresh wasm_exec.js
```

## Build and run

### build wasm

`make thick` or `make thin`

### then run

```
make server
cd static
../server
```

open [http://localhost:3131/](http://localhost:3131/) in your browser
