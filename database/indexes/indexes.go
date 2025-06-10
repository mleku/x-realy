package indexes

import (
	"io"

	"x.realy.lol/chk"
	"x.realy.lol/codec"
	"x.realy.lol/database/indexes/prefixes"
	"x.realy.lol/database/indexes/types/fullid"
	"x.realy.lol/database/indexes/types/fulltext"
	"x.realy.lol/database/indexes/types/identHash"
	"x.realy.lol/database/indexes/types/idhash"
	"x.realy.lol/database/indexes/types/kindidx"
	"x.realy.lol/database/indexes/types/letter"
	"x.realy.lol/database/indexes/types/prefix"
	"x.realy.lol/database/indexes/types/pubhash"
	"x.realy.lol/database/indexes/types/timestamp"
	"x.realy.lol/database/indexes/types/varint"
)

type Encs []codec.I

// T is a wrapper around an array of codec.I. The caller provides the Encs so they can then call
// the accessor function of the codec.I implementation.
type T struct {
	Encs
}

// New creates a new indexes. The helper functions below have an encode and decode variant, the
// decode variant does not add the prefix encoder because it has been read by prefixes.Identify.
func New(encoders ...codec.I) (i *T) { return &T{encoders} }

func (t *T) MarshalWrite(w io.Writer) (err error) {
	for _, e := range t.Encs {
		if err = e.MarshalWrite(w); chk.E(err) {
			return
		}
	}
	return
}

func (t *T) UnmarshalRead(r io.Reader) (err error) {
	for _, e := range t.Encs {
		if err = e.UnmarshalRead(r); chk.E(err) {
			return
		}
	}
	return
}

func EventVars() (ser *varint.V) {
	ser = varint.New()
	return
}
func EventEnc(ser *varint.V) (enc *T) {
	return New(prefix.New(prefixes.Event), ser)
}
func EventDec(ser *varint.V) (enc *T) {
	return New(prefix.New(), ser)
}

func IdVars() (id *idhash.T, ser *varint.V) {
	id = idhash.New()
	ser = varint.New()
	return
}
func IdEnc(id *idhash.T, ser *varint.V) (enc *T) {
	return New(prefix.New(prefixes.Id), id, ser)
}
func IdSearch(id *idhash.T) (enc *T) {
	return New(prefix.New(prefixes.Id), id)
}
func IdDec(id *idhash.T, ser *varint.V) (enc *T) {
	return New(prefix.New(), id, ser)
}

func FullIndexVars() (ser *varint.V, t *fullid.T, p *pubhash.T, ki *kindidx.T,
	ca *timestamp.T) {
	ser = varint.New()
	t = fullid.New()
	p = pubhash.New()
	ki = kindidx.FromKind(0)
	ca = &timestamp.T{}
	return
}
func FullIndexEnc(ser *varint.V, t *fullid.T, p *pubhash.T, ki *kindidx.T,
	ca *timestamp.T) (enc *T) {
	return New(prefix.New(prefixes.FullIndex), ser, t, p, ki, ca)
}
func FullIndexDec(ser *varint.V, t *fullid.T, p *pubhash.T, ki *kindidx.T,
	ca *timestamp.T) (enc *T) {
	return New(prefix.New(), ser, t, p, ki, ca)
}

func PubkeyVars() (p *pubhash.T, ser *varint.V) {
	p = pubhash.New()
	ser = varint.New()
	return
}
func PubkeyEnc(p *pubhash.T, ser *varint.V) (enc *T) {
	return New(prefix.New(prefixes.Pubkey), p, ser)
}
func PubkeyDec(p *pubhash.T, ser *varint.V) (enc *T) {
	return New(prefix.New(), p, ser)
}

func PubkeyCreatedAtVars() (p *pubhash.T, ca *timestamp.T, ser *varint.V) {
	p = pubhash.New()
	ca = &timestamp.T{}
	ser = varint.New()
	return
}
func PubkeyCreatedAtEnc(p *pubhash.T, ca *timestamp.T, ser *varint.V) (enc *T) {
	return New(prefix.New(prefixes.PubkeyCreatedAt), p, ca, ser)
}
func PubkeyCreatedAtDec(p *pubhash.T, ca *timestamp.T, ser *varint.V) (enc *T) {
	return New(prefix.New(), p, ca, ser)
}

func CreatedAtVars() (ca *timestamp.T, ser *varint.V) {
	ca = &timestamp.T{}
	ser = varint.New()
	return
}
func CreatedAtEnc(ca *timestamp.T, ser *varint.V) (enc *T) {
	return New(prefix.New(prefixes.CreatedAt), ca, ser)
}
func CreatedAtDec(ca *timestamp.T, ser *varint.V) (enc *T) {
	return New(prefix.New(), ca, ser)
}

func FirstSeenVars() (ser *varint.V, ts *timestamp.T) {
	ts = &timestamp.T{}
	ser = varint.New()
	return
}
func FirstSeenEnc(ser *varint.V, ts *timestamp.T) (enc *T) {
	return New(prefix.New(prefixes.FirstSeen), ser, ts)
}
func FirstSeenDec(ser *varint.V, ts *timestamp.T) (enc *T) {
	return New(prefix.New(), ser, ts)
}

func KindVars() (ki *kindidx.T, ser *varint.V) {
	ki = kindidx.FromKind(0)
	ser = varint.New()
	return
}
func KindEnc(ki *kindidx.T, ser *varint.V) (enc *T) {
	return New(prefix.New(prefixes.Kind), ki, ser)
}
func KindDec(ki *kindidx.T, ser *varint.V) (enc *T) {
	return New(prefix.New(), ki, ser)
}

func KindCreatedAtVars() (ki *kindidx.T, ca *timestamp.T, ser *varint.V) {
	ki = kindidx.FromKind(0)
	ca = &timestamp.T{}
	ser = varint.New()
	return
}
func KindCreatedAtEnc(ki *kindidx.T, ca *timestamp.T, ser *varint.V) (enc *T) {
	return New(prefix.New(prefixes.Kind), ki, ca, ser)
}
func KindCreatedAtDec(ki *kindidx.T, ca *timestamp.T, ser *varint.V) (enc *T) {
	return New(prefix.New(), ki, ca, ser)
}

func KindPubkeyCreatedAtVars() (ki *kindidx.T, p *pubhash.T, ca *timestamp.T, ser *varint.V) {
	ki = kindidx.FromKind(0)
	ser = varint.New()
	return
}
func KindPubkeyCreatedAtEnc(ki *kindidx.T, p *pubhash.T, ca *timestamp.T, ser *varint.V) (enc *T) {
	return New(prefix.New(prefixes.Kind), ki, p, ca, ser)
}
func KindPubkeyCreatedAtDec(ki *kindidx.T, p *pubhash.T, ca *timestamp.T, ser *varint.V) (enc *T) {
	return New(prefix.New(), ki, p, ca, ser)
}

type TagA struct {
	Ki  *kindidx.T
	P   *pubhash.T
	Id  *identhash.T
	Ser *varint.V
}

func TagAVars() (ki *kindidx.T, p *pubhash.T, id *identhash.T, ser *varint.V) {
	ki = kindidx.FromKind(0)
	p = pubhash.New()
	id = identhash.New()
	ser = varint.New()
	return
}
func TagAEnc(ki *kindidx.T, p *pubhash.T, id *identhash.T, ser *varint.V) (enc *T) {
	return New(prefix.New(prefixes.TagA), ki, p, id, ser)
}
func TagADec(ki *kindidx.T, p *pubhash.T, id *identhash.T, ser *varint.V) (enc *T) {
	return New(prefix.New(), ki, p, id, ser)
}

func TagEventVars() (id *idhash.T, ser *varint.V) {
	id = idhash.New()
	ser = varint.New()
	return
}
func TagEventEnc(id *idhash.T, ser *varint.V) (enc *T) {
	return New(prefix.New(prefixes.TagEvent), id, ser)
}
func TagEventDec(id *idhash.T, ser *varint.V) (enc *T) {
	return New(prefix.New(), id, ser)
}

func TagPubkeyVars() (p *pubhash.T, ser *varint.V) {
	p = pubhash.New()
	ser = varint.New()
	return
}
func TagPubkeyEnc(p *pubhash.T, ser *varint.V) (enc *T) {
	return New(prefix.New(prefixes.TagPubkey), p, ser)
}
func TagPubkeyDec(p *pubhash.T, ser *varint.V) (enc *T) {
	return New(prefix.New(), p, ser)
}

func TagHashtagVars() (hashtag *identhash.T, ser *varint.V) {
	hashtag = identhash.New()
	ser = varint.New()
	return
}
func TagHashtagEnc(hashtag *identhash.T, ser *varint.V) (enc *T) {
	return New(prefix.New(prefixes.TagHashtag), hashtag, ser)
}
func TagHashtagDec(hashtag *identhash.T, ser *varint.V) (enc *T) {
	return New(prefix.New(), hashtag, ser)
}

func TagIdentifierVars() (ident *identhash.T, ser *varint.V) {
	ident = identhash.New()
	ser = varint.New()
	return
}
func TagIdentifierEnc(ident *identhash.T, ser *varint.V) (enc *T) {
	return New(prefix.New(prefixes.TagIdentifier), ident, ser)
}
func TagIdentifierDec(ident *identhash.T, ser *varint.V) (enc *T) {
	return New(prefix.New(), ident, ser)
}

func TagLetterVars() (l *letter.T, val *identhash.T, ser *varint.V) {
	l = letter.New(0)
	val = identhash.New()
	ser = varint.New()
	return
}
func TagLetterEnc(l *letter.T, val *identhash.T, ser *varint.V) (enc *T) {
	return New(prefix.New(prefixes.TagLetter), l, val, ser)
}
func TagLetterDec(l *letter.T, val *identhash.T, ser *varint.V) (enc *T) {
	return New(prefix.New(), l, val, ser)
}

func TagProtectedVars() (p *pubhash.T, ser *varint.V) {
	p = pubhash.New()
	ser = varint.New()
	return
}
func TagProtectedEnc(p *pubhash.T, ser *varint.V) (enc *T) {
	return New(prefix.New(prefixes.TagProtected), p, ser)
}
func TagProtectedDec(p *pubhash.T, ser *varint.V) (enc *T) {
	return New(prefix.New(), p, ser)
}

func TagNonstandardVars() (key, value *identhash.T, ser *varint.V) {
	key = identhash.New()
	value = identhash.New()
	ser = varint.New()
	return
}
func TagNonstandardEnc(key, value *identhash.T, ser *varint.V) (enc *T) {
	return New(prefix.New(prefixes.TagNonstandard), key, value, ser)
}
func TagNonstandardDec(key, value *identhash.T, ser *varint.V) (enc *T) {
	return New(prefix.New(), key, value, ser)
}

func FullTextWordVars() (fw *fulltext.T, pos *varint.V, ser *varint.V) {
	fw = fulltext.New()
	pos = varint.New()
	ser = varint.New()
	return
}
func FullTextWordEnc(fw *fulltext.T, pos *varint.V, ser *varint.V) (enc *T) {
	return New(prefix.New(prefixes.FulltextWord), fw, pos, ser)
}
func FullTextWordDec(fw *fulltext.T, pos *varint.V, ser *varint.V) (enc *T) {
	return New(prefix.New(), fw, pos, ser)
}

func LastAccessedVars() (ser *varint.V) {
	ser = varint.New()
	return
}
func LastAccessedEnc(ser *varint.V) (enc *T) {
	return New(prefix.New(prefixes.LastAccessed), ser)
}
func LastAccessedDec(ser *varint.V) (enc *T) {
	return New(prefix.New(), ser)
}

func AccessCounterVars() (ser *varint.V) {
	ser = varint.New()
	return
}
func AccessCounterEnc(ser *varint.V) (enc *T) {
	return New(prefix.New(prefixes.AccessCounter), ser)
}
func AccessCounterDec(ser *varint.V) (enc *T) {
	return New(prefix.New(), ser)
}
