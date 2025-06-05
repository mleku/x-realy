package fulltext

import (
	"io"

	"x.realy.lol/chk"
	"x.realy.lol/varint"
)

type T struct {
	val []byte
}

func New() (ft *T) { return &T{} }

func (ft *T) FromWord(word []byte) { ft.val = word }

func (ft *T) Bytes() (b []byte) { return ft.val }

func (ft *T) MarshalWrite(w io.Writer) {
	varint.Encode(w, uint64(len(ft.val)))
	_, _ = w.Write(ft.val)
}

func (ft *T) UnmarshalRead(r io.Reader) (err error) {
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
