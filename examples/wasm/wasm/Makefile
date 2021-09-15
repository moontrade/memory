default:
	export GOROOT=/usr/local/Cellar/go/1.16.6/libexec; /Users/cmo/repos/moontrade/tinygo/build/tinygo build -o main.wasm -target wasi -scheduler coroutines -gc provided --print-allocs=. --no-debug .
	#export GOROOT=/usr/local/Cellar/go/1.16.6/libexec; /Users/cmo/repos/moontrade/tinygo/build/tinygo build -o conservative.wasm -target wasi -scheduler coroutines -gc conservative --print-allocs=. --no-debug .
	#go build -o maingo
	wasmer compile --llvm --dylib -o main.dylib main.wasm