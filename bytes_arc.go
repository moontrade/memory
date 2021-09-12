package mem

import "sync/atomic"

type BytesARC struct {
	Bytes
	count int32
}

func (p *BytesARC) Clone() BytesARC {
	if p == nil {
		return BytesARC{}
	}
	for {
		count := atomic.LoadInt32(&p.count)
		if count <= 0 {
			return BytesARC{}
		}
		if atomic.CompareAndSwapInt32(&p.count, count, count+1) {
			return *p
		}
	}
}

func (p *BytesARC) Drop() {
	if p == nil || p.addr == 0 {
		return
	}
	if atomic.AddInt32(&p.count, -1) == 0 {
		p.Bytes.Drop()
		p.addr = 0
	}
}
