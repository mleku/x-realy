package varint

import (
	"bytes"
	"io"

	"x.realy.lol/chk"
	"x.realy.lol/varint"
)

type V struct{ val uint64 }

func New() (s *V) { return &V{} }

func (vi *V) FromUint64(ser uint64) {
	vi.val = ser
	return
}

func FromBytes(ser []byte) (s *V, err error) {
	s = &V{}
	if s.val, err = varint.Decode(bytes.NewBuffer(ser)); chk.E(err) {
		return
	}
	return
}

func (vi *V) ToUint64() (ser uint64) { return vi.val }

func (vi *V) ToUint32() (v uint32) { return uint32(vi.val) }

func (vi *V) Bytes() (b []byte) {
	buf := new(bytes.Buffer)
	varint.Encode(buf, vi.val)
	return
}

func (vi *V) MarshalWrite(w io.Writer) (err error) {
	varint.Encode(w, vi.val)
	return
}

func (vi *V) UnmarshalRead(r io.Reader) (err error) {
	vi.val, err = varint.Decode(r)
	return
}
