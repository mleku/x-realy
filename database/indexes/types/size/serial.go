package size

import (
	"encoding/binary"
	"io"

	"x.realy.lol/errorf"
)

const Len = 4

type T struct{ val []byte }

func New() (s *T) { return &T{make([]byte, Len)} }

func (s *T) FromUint32(n uint32) {
	s.val = make([]byte, Len)
	binary.LittleEndian.PutUint32(s.val, n)
	return
}

func FromBytes(val []byte) (s *T, err error) {
	if len(val) != Len {
		err = errorf.E("size must be %d bytes long, got %d", Len, len(val))
		return
	}
	s = &T{val: val}
	return
}

func (s *T) ToUint32() (ser uint32) {
	ser = binary.LittleEndian.Uint32(s.val)
	return
}

func (s *T) Bytes() (b []byte) { return s.val }

func (s *T) MarshalWrite(w io.Writer) (err error) {
	_, err = w.Write(s.val)
	return
}

func (s *T) UnmarshalRead(r io.Reader) (err error) {
	if len(s.val) < Len {
		s.val = make([]byte, Len)
	} else {
		s.val = s.val[:Len]
	}
	_, err = r.Read(s.val)
	return
}
