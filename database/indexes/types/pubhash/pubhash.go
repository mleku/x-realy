package pubhash

import (
	"io"

	"x.realy.lol/chk"
	"x.realy.lol/ec/schnorr"
	"x.realy.lol/errorf"
	"x.realy.lol/helpers"
	"x.realy.lol/hex"
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

func (ph *T) FromPubkeyHex(pk string) (err error) {
	if len(pk) != schnorr.PubKeyBytesLen*2 {
		err = errorf.E("invalid Pubkey length, got %d require %d", len(pk), schnorr.PubKeyBytesLen*2)
		return
	}
	var pkb []byte
	if pkb, err = hex.Dec(pk); chk.E(err) {
		return
	}
	ph.val = helpers.Hash(pkb)[:Len]
	return
}

func (ph *T) Bytes() (b []byte) { return ph.val }

func (ph *T) MarshalWrite(w io.Writer) (err error) {
	_, err = w.Write(ph.val)
	return
}

func (ph *T) UnmarshalRead(r io.Reader) (err error) {
	if len(ph.val) < Len {
		ph.val = make([]byte, Len)
	} else {
		ph.val = ph.val[:Len]
	}
	_, err = r.Read(ph.val)
	return
}
