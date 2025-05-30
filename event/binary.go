package event

import (
	"io"

	"x.realy.lol/chk"
	"x.realy.lol/ec/schnorr"
	"x.realy.lol/hex"
	"x.realy.lol/timestamp"
	"x.realy.lol/varint"
)

// MarshalBinary writes a binary encoding of an event.
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
func (ev *E) MarshalBinary(w io.Writer) {
	_, _ = w.Write(ev.GetIdBytes())
	_, _ = w.Write(ev.GetPubkeyBytes())
	varint.Encode(w, uint64(ev.CreatedAt))
	varint.Encode(w, uint64(ev.Kind))
	varint.Encode(w, uint64(len(ev.Tags)))
	for _, x := range ev.Tags {
		varint.Encode(w, uint64(len(x)))
		for _, y := range x {
			varint.Encode(w, uint64(len(y)))
			_, _ = w.Write([]byte(y))
		}
	}
	varint.Encode(w, uint64(len(ev.Content)))
	_, _ = w.Write([]byte(ev.Content))
	_, _ = w.Write(ev.GetSigBytes())
	return
}

func (ev *E) UnmarshalBinary(r io.Reader) (err error) {
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
		for range nField {
			var lenField uint64
			if lenField, err = varint.Decode(r); chk.E(err) {
				return
			}
			field := make([]byte, lenField)
			if _, err = r.Read(field); chk.E(err) {
				return
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
