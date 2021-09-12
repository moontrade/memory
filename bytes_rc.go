package mem

type BytesRC struct {
	Bytes
	count int32
}

func NewRC(p Bytes) BytesRC {
	return BytesRC{Bytes: p, count: 1}
}

func (p *BytesRC) Clone() BytesRC {
	p.count++
	return *p
}

func (p *BytesRC) Drop() {
	if p == nil || p.addr == 0 {
		return
	}
	p.count--

	if p.count == 0 {
		p.Bytes.Drop()
		*p = BytesRC{}
	}
}
