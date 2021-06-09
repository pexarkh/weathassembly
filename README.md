# Weather Service with WebAssembly and Go

## Prereqs

```
go get github.com/flosch/pongo2
go get github.com/buger/jsonparser
make bootstrap # to refresh wasm_exec.js
```

## Build and run

`make thick` or `make thin`

then

```
mk server
cd static
../server
```

open [http://localhost:3131/](http://localhost:3131/) in your browser
