```shell
# Dependencies
apt install git llvm clang ninja-build cmake

# Build
git clone https://github.com/moontrade/tinygo
cd tinygo
make llvm-source
export CC=clang
export CXX=clang++
make llvm-build
make
./build/tinygo help
ldd ./build/tinygo

# Release
make release
git submodule update --init
tar -xvf path/to/release.tar.gz
./tinygo/bin/tinygo
```