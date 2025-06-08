package serial

import (
	"encoding/binary"
	"io"

	"x.realy.lol/errorf"
)

const Len = 8

type S struct{ val []byte }

func New() (s *S) { return &S{make([]byte, Len)} }

func (s *S) FromSerial(ser uint64) {
	binary.LittleEndian.PutUint64(s.val, ser)
	return
}

func FromBytes(ser []byte) (s *S, err error) {
	if len(ser) != Len {
		err = errorf.E("serial must be %d bytes long, got %d", Len, len(ser))
		return
	}
	s = &S{val: ser}
	return
}

func (s *S) ToSerial() (ser uint64) {
	ser = binary.LittleEndian.Uint64(s.val)
	return
}

func (s *S) Bytes() (b []byte) { return s.val }

func (s *S) MarshalWrite(w io.Writer) (err error) {
	_, err = w.Write(s.val)
	return
}

func (s *S) UnmarshalRead(r io.Reader) (err error) {
	if len(s.val) < Len {
		s.val = make([]byte, Len)
	} else {
		s.val = s.val[:Len]
	}
	_, err = r.Read(s.val)
	return
}
