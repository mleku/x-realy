package identhash

import (
	"io"

	"x.realy.lol/helpers"
)

const Len = 8

type T struct{ val []byte }

func New() (i *T) { return &T{make([]byte, Len)} }

func (i *T) FromIdent(id []byte) (err error) {
	i.val = helpers.Hash(id)[:Len]
	return
}

func (i *T) Bytes() (b []byte) { return i.val }

func (i *T) MarshalWrite(w io.Writer) { _, _ = w.Write(i.val) }

func (i *T) UnmarshalRead(r io.Reader) (err error) {
	if len(i.val) < Len {
		i.val = make([]byte, Len)
	} else {
		i.val = i.val[:Len]
	}
	_, err = r.Read(i.val)
	return
}
