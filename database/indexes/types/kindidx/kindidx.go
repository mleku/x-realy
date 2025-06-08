package kindidx

import (
	"encoding/binary"
	"io"

	"x.realy.lol/errorf"
)

const Len = 2

type T struct{ val []byte }

func FromKind(kind int) (k *T) {
	k = &T{val: make([]byte, Len)}
	binary.LittleEndian.PutUint16(k.val, uint16(kind))
	return
}

func FromBytes(kindBytes []byte) (k *T, err error) {
	if len(kindBytes) != Len {
		err = errorf.E("kind must be %d bytes long, got %d", Len, len(kindBytes))
		return
	}
	k = &T{val: kindBytes}
	return
}

func (k *T) Set(ki int) {
	kk := FromKind(ki)
	k.val = kk.val
}

func (k *T) ToKind() (kind int) {
	kind = int(binary.LittleEndian.Uint16(k.val))
	return
}
func (k *T) Bytes() (b []byte) { return k.val }

func (k *T) MarshalWrite(w io.Writer) (err error) {
	_, err = w.Write(k.val)
	return
}

func (k *T) UnmarshalRead(r io.Reader) (err error) {
	if len(k.val) < Len {
		k.val = make([]byte, Len)
	} else {
		k.val = k.val[:Len]
	}
	_, err = r.Read(k.val)
	return
}
