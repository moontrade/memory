package main

import (
	"errors"
	"io"
	"unsafe"
)

type Bytes8112 [8112]byte

func NewBytes8112(s string) *Bytes8112 {
	v := Bytes8112{}
	v.set(s)
	return &v
}
func (s *Bytes8112) set(v string) {
	copy(s[0:], v)
}
func (s *Bytes8112) Len() int {
	return 8112
}
func (s *Bytes8112) Cap() int {
	return 8112
}
func (s *Bytes8112) Unsafe() string {
	return *(*string)(unsafe.Pointer(s))
}
func (s *Bytes8112) String() string {
	return string(s[0:s.Len()])
}
func (s *Bytes8112) Bytes() []byte {
	return s[0:s.Len()]
}
func (s *Bytes8112) Clone() *Bytes8112 {
	v := Bytes8112{}
	copy(s[0:], v[0:])
	return &v
}
func (s *Bytes8112) Mut() *Bytes8112Mut {
	return *(**Bytes8112Mut)(unsafe.Pointer(&s))
}
func (s *Bytes8112) MarshalBinaryTo(b []byte) []byte {
	return append(b, (*(*[8112]byte)(unsafe.Pointer(&s)))[0:]...)
}
func (s *Bytes8112) MarshalBinary() ([]byte, error) {
	var v []byte
	return append(v, (*(*[8112]byte)(unsafe.Pointer(&s)))[0:]...), nil
}
func (s *Bytes8112) UnmarshalBinary(b []byte) error {
	if len(b) < 8112 {
		return io.ErrShortBuffer
	}
	v := (*Bytes8112)(unsafe.Pointer(&b[0]))
	*s = *v
	return nil
}

type Bytes8112Mut struct {
	Bytes8112
}

func (s *Bytes8112Mut) Set(v string) {
	s.set(v)
}

// Block8
type Block8 struct {
	head BlockHeader
	body Bytes8112
}

func (s *Block8) String() string {
	return ""
	//return fmt.Sprintf("%v", s.MarshalMap(nil))
}

func (s *Block8) MarshalMap(m map[string]interface{}) map[string]interface{} {
	if m == nil {
		m = make(map[string]interface{})
	}
	m["head"] = s.Head().MarshalMap(nil)
	m["body"] = s.Body()
	return m
}

func (s *Block8) MarshalBinaryTo(b []byte) []byte {
	return append(b, (*(*[8208]byte)(unsafe.Pointer(s)))[0:]...)
}
func (s *Block8) MarshalBinary() ([]byte, error) {
	var v []byte
	return append(v, (*(*[8208]byte)(unsafe.Pointer(s)))[0:]...), nil
}
func (s *Block8) Read(b []byte) (n int, err error) {
	if len(b) < 8208 {
		return -1, errors.New("short buffer")
	}
	v := (*Block8)(unsafe.Pointer(&b[0]))
	*v = *s
	return 8208, nil
}
func (s *Block8) UnmarshalBinary(b []byte) error {
	if len(b) < 8208 {
		return errors.New("short buffer")
	}
	v := (*Block8)(unsafe.Pointer(&b[0]))
	*s = *v
	return nil
}
func (s *Block8) Clone() *Block8 {
	v := &Block8{}
	*v = *s
	return v
}
func (s *Block8) Bytes() []byte {
	return (*(*[8208]byte)(unsafe.Pointer(s)))[0:]
}
func (s *Block8) Mut() *Block8Mut {
	return (*Block8Mut)(unsafe.Pointer(s))
}
func (s *Block8) Head() *BlockHeader {
	return &s.head
}
func (s *Block8) Body() *Bytes8112 {
	return &s.body
}

// Block8
type Block8Mut struct {
	Block8
}

func (s *Block8Mut) Clone() *Block8Mut {
	v := &Block8Mut{}
	*v = *s
	return v
}
func (s *Block8Mut) Freeze() *Block8 {
	return (*Block8)(unsafe.Pointer(s))
}
func (s *Block8Mut) Head() *BlockHeaderMut {
	return s.head.Mut()
}
func (s *Block8Mut) SetHead(v *BlockHeader) *Block8Mut {
	s.head = *v
	return s
}
func (s *Block8Mut) SetBody(v Bytes8112) *Block8Mut {
	s.body = v
	return s
}

// BlockHeader
type BlockHeader struct {
	streamID    int64
	id          int64
	created     int64
	completed   int64
	min         int64
	max         int64
	start       int64
	end         int64
	savepoint   int64
	savepointR  int64
	count       uint16
	seq         uint16
	size        uint16
	sizeU       uint16
	sizeX       uint16
	compression Compression
	_           [5]byte // Padding
}

func (s *BlockHeader) String() string {
	return ""
	//return fmt.Sprintf("%v", s.MarshalMap(nil))
}

func (s *BlockHeader) MarshalMap(m map[string]interface{}) map[string]interface{} {
	if m == nil {
		m = make(map[string]interface{})
	}
	m["streamID"] = s.StreamID()
	m["id"] = s.Id()
	m["created"] = s.Created()
	m["completed"] = s.Completed()
	m["min"] = s.Min()
	m["max"] = s.Max()
	m["start"] = s.Start()
	m["end"] = s.End()
	m["savepoint"] = s.Savepoint()
	m["savepointR"] = s.SavepointR()
	m["count"] = s.Count()
	m["seq"] = s.Seq()
	m["size"] = s.Size()
	m["sizeU"] = s.SizeU()
	m["sizeX"] = s.SizeX()
	m["compression"] = s.Compression()
	return m
}

func (s *BlockHeader) MarshalBinaryTo(b []byte) []byte {
	return append(b, (*(*[96]byte)(unsafe.Pointer(s)))[0:]...)
}
func (s *BlockHeader) MarshalBinary() ([]byte, error) {
	var v []byte
	return append(v, (*(*[96]byte)(unsafe.Pointer(s)))[0:]...), nil
}
func (s *BlockHeader) Read(b []byte) (n int, err error) {
	if len(b) < 96 {
		return -1, errors.New("short buffer")
	}
	v := (*BlockHeader)(unsafe.Pointer(&b[0]))
	*v = *s
	return 96, nil
}
func (s *BlockHeader) UnmarshalBinary(b []byte) error {
	if len(b) < 96 {
		return errors.New("short buffer")
	}
	v := (*BlockHeader)(unsafe.Pointer(&b[0]))
	*s = *v
	return nil
}
func (s *BlockHeader) Clone() *BlockHeader {
	v := &BlockHeader{}
	*v = *s
	return v
}
func (s *BlockHeader) Bytes() []byte {
	return (*(*[96]byte)(unsafe.Pointer(s)))[0:]
}
func (s *BlockHeader) Mut() *BlockHeaderMut {
	return (*BlockHeaderMut)(unsafe.Pointer(s))
}
func (s *BlockHeader) StreamID() int64 {
	return s.streamID
}
func (s *BlockHeader) Id() int64 {
	return s.id
}
func (s *BlockHeader) Created() int64 {
	return s.created
}
func (s *BlockHeader) Completed() int64 {
	return s.completed
}
func (s *BlockHeader) Min() int64 {
	return s.min
}
func (s *BlockHeader) Max() int64 {
	return s.max
}
func (s *BlockHeader) Start() int64 {
	return s.start
}
func (s *BlockHeader) End() int64 {
	return s.end
}
func (s *BlockHeader) Savepoint() int64 {
	return s.savepoint
}
func (s *BlockHeader) SavepointR() int64 {
	return s.savepointR
}
func (s *BlockHeader) Count() uint16 {
	return s.count
}
func (s *BlockHeader) Seq() uint16 {
	return s.seq
}
func (s *BlockHeader) Size() uint16 {
	return s.size
}
func (s *BlockHeader) SizeU() uint16 {
	return s.sizeU
}
func (s *BlockHeader) SizeX() uint16 {
	return s.sizeX
}
func (s *BlockHeader) Compression() Compression {
	return s.compression
}

// BlockHeader
type BlockHeaderMut struct {
	BlockHeader
}

func (s *BlockHeaderMut) Clone() *BlockHeaderMut {
	v := &BlockHeaderMut{}
	*v = *s
	return v
}
func (s *BlockHeaderMut) Freeze() *BlockHeader {
	return (*BlockHeader)(unsafe.Pointer(s))
}
func (s *BlockHeaderMut) SetStreamID(v int64) *BlockHeaderMut {
	s.streamID = v
	return s
}
func (s *BlockHeaderMut) SetId(v int64) *BlockHeaderMut {
	s.id = v
	return s
}
func (s *BlockHeaderMut) SetCreated(v int64) *BlockHeaderMut {
	s.created = v
	return s
}
func (s *BlockHeaderMut) SetCompleted(v int64) *BlockHeaderMut {
	s.completed = v
	return s
}
func (s *BlockHeaderMut) SetMin(v int64) *BlockHeaderMut {
	s.min = v
	return s
}
func (s *BlockHeaderMut) SetMax(v int64) *BlockHeaderMut {
	s.max = v
	return s
}
func (s *BlockHeaderMut) SetStart(v int64) *BlockHeaderMut {
	s.start = v
	return s
}
func (s *BlockHeaderMut) SetEnd(v int64) *BlockHeaderMut {
	s.end = v
	return s
}
func (s *BlockHeaderMut) SetSavepoint(v int64) *BlockHeaderMut {
	s.savepoint = v
	return s
}
func (s *BlockHeaderMut) SetSavepointR(v int64) *BlockHeaderMut {
	s.savepointR = v
	return s
}
func (s *BlockHeaderMut) SetCount(v uint16) *BlockHeaderMut {
	s.count = v
	return s
}
func (s *BlockHeaderMut) SetSeq(v uint16) *BlockHeaderMut {
	s.seq = v
	return s
}
func (s *BlockHeaderMut) SetSize(v uint16) *BlockHeaderMut {
	s.size = v
	return s
}
func (s *BlockHeaderMut) SetSizeU(v uint16) *BlockHeaderMut {
	s.sizeU = v
	return s
}
func (s *BlockHeaderMut) SetSizeX(v uint16) *BlockHeaderMut {
	s.sizeX = v
	return s
}
func (s *BlockHeaderMut) SetCompression(v Compression) *BlockHeaderMut {
	s.compression = v
	return s
}

type Compression byte
