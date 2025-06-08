package event

import (
	"io"

	"x.realy.lol/chk"
	"x.realy.lol/ec/schnorr"
	"x.realy.lol/errorf"
	"x.realy.lol/hex"
	"x.realy.lol/timestamp"
	"x.realy.lol/varint"
)

// todo: maybe we should make e and p tag values binary to reduce space usage

// MarshalWrite writes a binary encoding of an event.
//
// NOTE: Event must not be nil or this will panic. Use event.New.
//
// [ 32 bytes Id ]
// [ 32 bytes Pubkey ]
// [ varint CreatedAt ]
// [ 2 bytes Kind ]
// [ varint Tags length ]
//
//	[ varint tag length ]
//	  [ varint tag element length ]
//	  [ tag element data ]
//	...
//
// [ varint Content length ]
// [ 64 bytes Sig ]
func (ev *E) MarshalWrite(w io.Writer) (err error) {
	if ev == nil {
		panic("cannot marshal a nil event")
	}
	_, _ = w.Write(ev.GetIdBytes())
	_, _ = w.Write(ev.GetPubkeyBytes())
	varint.Encode(w, uint64(ev.CreatedAt))
	varint.Encode(w, uint64(ev.Kind))
	varint.Encode(w, uint64(len(ev.Tags)))
	for _, x := range ev.Tags {
		varint.Encode(w, uint64(len(x)))
		// e and p tag values should be hex
		var isBin bool
		if len(x) > 1 && (x[0] == "e" || x[0] == "p") {
			isBin = true
		}
		for i, y := range x {
			if i == 1 && isBin {
				var b []byte
				b, err = hex.Dec(y)
				if err != nil {
					err = errorf.E("e or p tag value not hex: %s", err.Error())
					return
				}
				if len(b) != 32 {
					err = errorf.E("e or p tag value with invalid decoded byte length %d", len(b))
					return
				}
				varint.Encode(w, uint64(len(b)))
				_, _ = w.Write(b)
			} else {
				varint.Encode(w, uint64(len(y)))
				_, _ = w.Write([]byte(y))
			}
		}
	}
	varint.Encode(w, uint64(len(ev.Content)))
	_, _ = w.Write([]byte(ev.Content))
	_, _ = w.Write(ev.GetSigBytes())
	return err
}

// UnmarshalRead decodes an event in binary form into an allocated event struct.
//
// NOTE: Event must not be nil or this will panic. Use event.New.
func (ev *E) UnmarshalRead(r io.Reader) (err error) {
	if ev == nil {
		panic("cannot unmarshal into nil event struct")
	}
	id := make([]byte, 32)
	if _, err = r.Read(id); chk.E(err) {
		return
	}
	ev.Id = hex.Enc(id)
	pubkey := make([]byte, 32)
	if _, err = r.Read(pubkey); chk.E(err) {
		return
	}
	ev.Pubkey = hex.Enc(pubkey)
	var ca uint64
	if ca, err = varint.Decode(r); chk.E(err) {
		return
	}
	ev.CreatedAt = timestamp.New(ca)
	var k uint64
	if k, err = varint.Decode(r); chk.E(err) {
		return
	}
	ev.Kind = int(k)
	var nTags uint64
	if nTags, err = varint.Decode(r); chk.E(err) {
		return
	}
	for range nTags {
		var nField uint64
		if nField, err = varint.Decode(r); chk.E(err) {
			return
		}
		var t []string
		var isBin bool
		for i := range nField {
			var lenField uint64
			if lenField, err = varint.Decode(r); chk.E(err) {
				return
			}
			field := make([]byte, lenField)
			if _, err = r.Read(field); chk.E(err) {
				return
			}
			// if it is first field, length 1 and is e or p, the value field must be binary
			if i == 0 && len(field) == 1 && (field[0] == 'e' || field[0] == 'p') {
				isBin = true
			}
			if i == 1 && isBin {
				// this is a binary value, was an e or p tag key, 32 bytes long, encode value
				// field to hex
				f := make([]byte, 64)
				_ = hex.EncBytes(f, field)
				field = f
			}
			t = append(t, string(field))
		}
		ev.Tags = append(ev.Tags, t)
	}
	var cLen uint64
	if cLen, err = varint.Decode(r); chk.E(err) {
		return
	}
	content := make([]byte, cLen)
	if _, err = r.Read(content); chk.E(err) {
		return
	}
	ev.Content = string(content)
	sig := make([]byte, schnorr.SignatureSize)
	if _, err = r.Read(sig); chk.E(err) {
		return
	}
	ev.Sig = hex.Enc(sig)
	return
}
