package timestamp

import (
	"bytes"
	"io"

	"x.realy.lol/chk"
	"x.realy.lol/database/indexes/types/varint"
	timeStamp "x.realy.lol/timestamp"
)

const Len = 8

type T struct{ val int }

func (ts *T) FromInt(t int)     { ts.val = t }
func (ts *T) FromInt64(t int64) { ts.val = int(t) }

func FromBytes(timestampBytes []byte) (ts *T, err error) {
	v := varint.New()
	if err = v.UnmarshalRead(bytes.NewBuffer(timestampBytes)); chk.E(err) {
		return
	}
	ts = &T{val: v.ToInt()}
	return
}

func (ts *T) ToTimestamp() (timestamp timeStamp.Timestamp) {
	return
}
func (ts *T) Bytes() (b []byte, err error) {
	v := varint.New()
	buf := new(bytes.Buffer)
	if err = v.MarshalWrite(buf); chk.E(err) {
		return
	}
	b = buf.Bytes()
	return
}

func (ts *T) MarshalWrite(w io.Writer) (err error) {
	v := varint.New()
	if err = v.MarshalWrite(w); chk.E(err) {
		return
	}
	return
}

func (ts *T) UnmarshalRead(r io.Reader) (err error) {
	v := varint.New()
	if err = v.UnmarshalRead(r); chk.E(err) {
		return
	}
	ts.val = v.ToInt()
	return
}
