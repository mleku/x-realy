package kindidx

import (
	"encoding/binary"
	"io"

	"x.realy.lol/errorf"
)

const Len = 2

type t struct{ val []byte }

func FromKind(kind int) (k *t) {
	k = &t{val: make([]byte, Len)}
	binary.LittleEndian.PutUint16(k.val, uint16(kind))
	return
}

func FromBytes(kindBytes []byte) (k *t, err error) {
	if len(kindBytes) != Len {
		err = errorf.E("kind must be %d bytes long, got %d", Len, len(kindBytes))
		return
	}
	k = &t{val: kindBytes}
	return
}

func (k *t) ToKind() (kind int) {
	kind = int(binary.LittleEndian.Uint16(k.val))
	return
}
func (k *t) Bytes() (b []byte) { return k.val }

func (k *t) MarshalBinary(w io.Writer) { _, _ = w.Write(k.val) }

func (k *t) UnmarshalBinary(r io.Reader) (err error) {
	if len(k.val) < Len {
		k.val = make([]byte, Len)
	} else {
		k.val = k.val[:Len]
	}
	_, err = r.Read(k.val)
	return
}
