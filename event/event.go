package event

import (
	"bytes"
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/minio/sha256-simd"

	"x.realy.lol/chk"
	"x.realy.lol/errorf"
	"x.realy.lol/hex"
	"x.realy.lol/log"
	"x.realy.lol/p256k"
	"x.realy.lol/signer"
	"x.realy.lol/tags"
	"x.realy.lol/text"
	"x.realy.lol/timestamp"
)

type E struct {
	Id        string              `json:"id"`
	Pubkey    string              `json:"pubkey"`
	CreatedAt timestamp.Timestamp `json:"created_at"`
	Kind      int                 `json:"kind`
	Tags      tags.Tags           `json:"tags"`
	Content   string              `json:"content"`
	Sig       string              `json:"sig"`
}

func New() (ev *E) { return &E{} }

func (ev *E) Marshal() (b []byte, err error) {
	if b, err = json.Marshal(ev); chk.E(err) {
		return
	}
	return
}

func (ev *E) Unmarshal(b []byte) (err error) {
	if err = json.Unmarshal(b, ev); chk.E(err) {
		return
	}
	return
}

// Serialize does a json.Marshal and ignores errors, only logging them. Mostly a convenience for
// logging.
func (ev *E) Serialize() (b []byte) {
	var err error
	if b, err = json.Marshal(ev); chk.E(err) {
		return
	}
	return
}

// Sign an event using a provided signer initialized with a secret key, rewrite pubkey, id and
// signature as required.
func (ev *E) Sign(sign signer.I) (err error) {
	// need to change pub as this is part of the message
	ev.Pubkey = hex.Enc(sign.Pub())
	id := ev.GenIdBytes()
	ev.Id = hex.Enc(id)
	var sig []byte
	if sig, err = sign.Sign(id); chk.E(err) {
		return
	}
	ev.Sig = hex.Enc(sig)
	return
}

// Verify an event is signed by the pubkey it contains. Uses
// github.com/bitcoin-core/secp256k1 if available for faster verification.
func (ev *E) Verify() (valid bool, err error) {
	keys := p256k.Signer{}
	if err = keys.InitPub(ev.GetPubkeyBytes()); chk.E(err) {
		return
	}
	if valid, err = keys.Verify(ev.GetIdBytes(), ev.GetSigBytes()); chk.T(err) {
		// check that this isn't because of a bogus Id
		id := ev.GenIdBytes()
		if !bytes.Equal(id, ev.GetIdBytes()) {
			log.E.Ln("event Id incorrect")
			ev.Id = hex.Enc(id)
			err = nil
			if valid, err = keys.Verify(ev.GetIdBytes(), ev.GetSigBytes()); chk.E(err) {
				return
			}
			err = errorf.W("event Id incorrect but signature is valid on correct Id")
		}
		return
	}
	return
}

// ToCanonical converts the event to the canonical encoding used to derive the
// event Id.
func (ev *E) ToCanonical(dst []byte) (b []byte) {
	b = dst
	b = append(b, "[0,\""...)
	b = append(b, ev.Pubkey...)
	b = append(b, "\","...)
	b = append(b, fmt.Sprint(ev.CreatedAt)...)
	b = append(b, ',')
	b = append(b, fmt.Sprint(ev.Kind)...)
	b = append(b, ',')
	tb, _ := json.Marshal(ev.Tags)
	b = append(b, tb...)
	b = append(b, ',')
	b = text.AppendQuote(b, []byte(ev.Content), text.NostrEscape)
	b = append(b, ']')
	return
}

func (ev *E) GenIdBytes() (b []byte) {
	var can []byte
	can = ev.ToCanonical(can)
	return Hash(can)
}

func (ev *E) GetIdBytes() (i []byte) {
	var err error
	if i, err = hex.Dec(ev.Id); chk.E(err) {
		return
	}
	return
}

func (ev *E) GetSigBytes() (s []byte) {
	var err error
	if s, err = hex.Dec(ev.Sig); chk.E(err) {
		return
	}
	return
}

func (ev *E) GetPubkeyBytes() (p []byte) {
	var err error
	if p, err = hex.Dec(ev.Pubkey); chk.E(err) {
		return
	}
	return
}

func (ev *E) IdHex() (idHex string) {
	can := ev.ToCanonical(nil)
	idHex = hex.Enc(Hash(can))
	return
}

func (ev *E) CheckId() (ok bool) {
	idHex := ev.IdHex()
	return idHex == ev.Id
}

// Hash is a little helper generate a hash and return a slice instead of an
// array.
func Hash(in []byte) (out []byte) {
	h := sha256.Sum256(in)
	return h[:]
}

// this is an absolute minimum length canonical encoded event
var minimal = len(`[0,"0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef",1733739427,0,[],""]`)

// FromCanonical reverses the process of creating the canonical encoding, note
// that the signature is missing in this form. Allocate an event.T before
// calling this.
func (ev *E) FromCanonical(b []byte) (err error) {
	if len(b) < minimal {
		err = errorf.E("event is too short to be a canonical event, require at least %d got %d",
			minimal, len(b))
		return
	}
	var un []any
	if err = json.Unmarshal(b, &un); chk.E(err) {
		return
	}
	if len(un) < 5 {
		err = errorf.E("canonical event must have 5 array elements, got %d", len(un))
		return
	}
	var ok bool
	if ev.Pubkey, ok = un[1].(string); !ok {
		err = errorf.E("failed to get pubkey value, got type %v expected string", reflect.TypeOf(un[1]))
		return
	}
	var createdAt float64
	if createdAt, ok = un[2].(float64); !ok {
		err = errorf.E("failed to get created_at value, got type %v expected float64", reflect.TypeOf(un[2]))
		return
	}
	ev.CreatedAt = timestamp.New(createdAt)
	var kind float64
	if kind, ok = un[3].(float64); !ok {
		err = errorf.E("failed to get kind value, got type %v expected float64", reflect.TypeOf(un[3]))
		return
	}
	ev.Kind = int(kind)
	var tags []any
	if tags, ok = un[4].([]any); !ok {
		err = errorf.E("failed to get tags value, got type %v expected []interface", reflect.TypeOf(un[4]))
		return
	}
	if ev.Tags, err = FromSliceInterface(tags); chk.E(err) {
		return
	}
	if ev.Content, ok = un[5].(string); !ok {
		err = errorf.E("failed to get tags value, got type %v expected []interface", reflect.TypeOf(un[4]))
		return
	}
	return
}

func FromSliceInterface(in []any) (t tags.Tags, err error) {
	t = make(tags.Tags, 0)
	for _, v := range in {
		var ok bool
		var vv []any
		if vv, ok = v.([]any); !ok {
			err = errorf.E("failed to get tag value, got type %v expected []interface", reflect.TypeOf(v))
			return
		}
		var tag []string
		for _, w := range vv {
			var x string
			if x, ok = w.(string); !ok {
				err = errorf.E("failed to get tag value, got type %v expected string", reflect.TypeOf(w))
				return
			}
			tag = append(tag, x)
		}
		t = append(t, tag)
	}
	return
}
