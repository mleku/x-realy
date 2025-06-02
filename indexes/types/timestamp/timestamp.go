package timestamp

import (
	"encoding/binary"
	"io"

	"x.realy.lol/errorf"
	timeStamp "x.realy.lol/timestamp"
)

const Len = 8

type t struct{ val []byte }

func FromInt64(timestamp int64) (ts *t) {
	ts = &t{val: make([]byte, Len)}
	binary.LittleEndian.PutUint64(ts.val, uint64(timestamp))
	return
}

func FromBytes(timestampBytes []byte) (ts *t, err error) {
	if len(timestampBytes) != Len {
		err = errorf.E("kind must be %d bytes long, got %d", Len, len(timestampBytes))
		return
	}
	ts = &t{val: timestampBytes}
	return
}

func (ts *t) ToTimestamp() (timestamp timeStamp.Timestamp) {
	return timeStamp.Timestamp(binary.LittleEndian.Uint64(ts.val))
}
func (ts *t) Bytes() (b []byte) { return ts.val }

func (ts *t) MarshalBinary(w io.Writer) { _, _ = w.Write(ts.val) }

func (ts *t) UnmarshalBinary(r io.Reader) (err error) {
	if len(ts.val) < Len {
		ts.val = make([]byte, Len)
	} else {
		ts.val = ts.val[:Len]
	}
	_, err = r.Read(ts.val)
	return
}
