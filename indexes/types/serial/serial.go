package serial

import (
	"encoding/binary"
	"io"

	"x.realy.lol/errorf"
)

const Len = 8

type t struct{ val []byte }

func New() (s *t) { return &t{make([]byte, Len)} }

func FromSerial(ser uint64) (s *t) {
	s = &t{val: make([]byte, Len)}
	binary.LittleEndian.PutUint64(s.val, ser)
	return
}

func FromBytes(ser []byte) (s *t, err error) {
	if len(ser) != Len {
		err = errorf.E("serial must be %d bytes long, got %d", Len, len(ser))
		return
	}
	s = &t{val: ser}
	return
}

func (s *t) ToSerial() (ser uint64) {
	ser = binary.LittleEndian.Uint64(s.val)
	return
}

func (s *t) Bytes() (b []byte) { return s.val }

func (s *t) MarshalBinary(w io.Writer) { _, _ = w.Write(s.val) }

func (s *t) UnmarshalBinary(r io.Reader) (err error) {
	if len(s.val) < Len {
		s.val = make([]byte, Len)
	} else {
		s.val = s.val[:Len]
	}
	_, err = r.Read(s.val)
	return
}
