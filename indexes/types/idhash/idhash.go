package idhash

import (
	"io"

	"github.com/minio/sha256-simd"

	"x.realy.lol/errorf"
	"x.realy.lol/helpers"
)

const Len = 8

type t struct{ val []byte }

func New() (i *t) { return &t{make([]byte, Len)} }

func (i *t) FromId(id []byte) (err error) {
	if len(id) != sha256.Size {
		err = errorf.E("invalid Id length, got %d require %d", len(id), sha256.Size)
		return
	}
	i.val = helpers.Hash(id)
	return
}

func (i *t) Bytes() (b []byte) { return i.val }

func (i *t) MarshalBinary(w io.Writer) { _, _ = w.Write(i.val) }

func (i *t) UnmarshalBinary(r io.Reader) (err error) {
	if len(i.val) < Len {
		i.val = make([]byte, Len)
	} else {
		i.val = i.val[:Len]
	}
	_, err = r.Read(i.val)
	return
}
