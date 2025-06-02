package fullid

import (
	"io"

	"github.com/minio/sha256-simd"

	"x.realy.lol/errorf"
)

const Len = sha256.Size

type t struct {
	val []byte
}

func New() (fi *t) { return &t{make([]byte, Len)} }

func (fi *t) FromId(id []byte) (err error) {
	if len(id) != Len {
		err = errorf.E("invalid Id length, got %d require %d", len(id), Len)
		return
	}
	fi.val = id
	return
}
func (fi *t) Bytes() (b []byte) { return fi.val }

func (fi *t) MarshalBinary(w io.Writer) { _, _ = w.Write(fi.val) }

func (fi *t) UnmarshalBinary(r io.Reader) (err error) {
	if len(fi.val) < Len {
		fi.val = make([]byte, Len)
	} else {
		fi.val = fi.val[:Len]
	}
	_, err = r.Read(fi.val)
	return
}
