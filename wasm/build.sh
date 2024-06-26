GOOS=js GOARCH=wasm go build -trimpath -o ../html/dist/wasm.wasm  ../wasm/wasm.go
# cp "$(go env GOROOT)/misc/wasm/wasm_exec.js" static