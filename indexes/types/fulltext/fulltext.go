package fulltext

import (
	"io"

	"x.realy.lol/chk"
	"x.realy.lol/varint"
)

type t struct {
	val []byte
}

func New() (ft *t) { return &t{} }

func FromWord(word []byte) (ft *t, err error) {
	ft = &t{val: word}
	return
}

func (ft *t) Bytes() (b []byte) { return ft.val }

func (ft *t) MarshalBinary(w io.Writer) {
	varint.Encode(w, uint64(len(ft.val)))
	_, _ = w.Write(ft.val)
}

func (ft *t) UnmarshalBinary(r io.Reader) (err error) {
	var l uint64
	if l, err = varint.Decode(r); chk.E(err) {
		return
	}
	wl := int(l)
	if len(ft.val) < wl {
		ft.val = make([]byte, wl)
	} else {
		ft.val = ft.val[:wl]
	}
	_, err = r.Read(ft.val)
	return
}
