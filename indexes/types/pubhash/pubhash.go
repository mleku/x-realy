package pubhash

import (
	"io"

	"x.realy.lol/ec/schnorr"
	"x.realy.lol/errorf"
	"x.realy.lol/helpers"
)

const Len = 8

type T struct{ val []byte }

func New() (ph *T) { return &T{make([]byte, Len)} }

func (ph *T) FromPubkey(pk []byte) (err error) {
	if len(pk) != schnorr.PubKeyBytesLen {
		err = errorf.E("invalid Pubkey length, got %d require %d", len(pk), schnorr.PubKeyBytesLen)
		return
	}
	ph.val = helpers.Hash(pk)[:Len]
	return
}

func (ph *T) Bytes() (b []byte) { return ph.val }

func (ph *T) MarshalWrite(w io.Writer) { _, _ = w.Write(ph.val) }

func (ph *T) UnmarshalRead(r io.Reader) (err error) {
	if len(ph.val) < Len {
		ph.val = make([]byte, Len)
	} else {
		ph.val = ph.val[:Len]
	}
	_, err = r.Read(ph.val)
	return
}
