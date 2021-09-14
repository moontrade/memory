//go:build 386 || amd64 || arm || arm64 || ppc64le || mips64le || mipsle || riscv64 || wasm || tinygo.wasm
// +build 386 amd64 arm arm64 ppc64le mips64le mipsle riscv64 wasm tinygo.wasm

package mem

const (
	EmptyString = ""
)

func (p *Bytes) Reserve(length int) {
	p.ensureCapU32(uint32(length))
}
