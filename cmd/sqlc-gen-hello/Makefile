sqlc-gen-hello.wasm: main.go
	tinygo build -o sqlc-gen-hello.wasm -gc=leaking -scheduler=none -wasm-abi=generic -target=wasi main.go
	openssl sha256 sqlc-gen-hello.wasm
