package database

import (
	"bytes"
	"encoding/binary"

	"x.realy.lol/chk"
	"x.realy.lol/ec/schnorr"
	"x.realy.lol/errorf"
	"x.realy.lol/helpers"
	"x.realy.lol/timestamp"
)

type I string

func (i I) B() []byte { return []byte(i) }

// the following enumerations are separate from the prefix value for simpler reference.

const (
	// Event is the whole event stored in binary format
	//
	//   [ prefix ][ 8 byte serial ] [ event in binary format ]
	Event = iota

	// Config is a singular record containing a free-form configuration in JSON format
	//
	// [ prefix ] [ configuration in JSON format ]
	Config

	// Id contains a truncated 8 byte hash of an event index
	//
	// [ prefix ][ 8 bytes truncated hash of Id ][ 8 serial ]
	Id

	// FullIndex is an index designed to enable sorting and filtering of results found via
	// other indexes.
	//
	// [ prefix ][ 32 bytes full event ID ][ 8 bytes truncated hash of pubkey ][ 2 bytes kind ][ 8 bytes created_at timestamp ][ 8 serial ]
	FullIndex

	// Pubkey is an index for searching for events authored by a pubkey.
	//
	// [ prefix ][ 8 bytes truncated hash of pubkey ][ 8 serial ]
	Pubkey

	// PubkeyCreatedAt is a composite index that allows search by pubkey filtered by
	// created_at.
	//
	// [ prefix ][ 8 bytes truncated hash of pubkey ][ 8 bytes created_at ][ 8 serial ]
	PubkeyCreatedAt

	// CreatedAt is an index that allows search the timestamp on the event.
	//
	// [ prefix ][ created_at 8 bytes timestamp ][ 8 serial ]
	CreatedAt

	// FirstSeen is an index that records the timestamp of when the event was first seen.
	//
	// [ prefix ][ 8 serial ][ 8 byte timestamp ]
	FirstSeen

	// Kind is an index of event kind numbers.
	//
	// [ prefix ][ 2 bytes kind number ][ 8 serial ]
	Kind

	// TagA is an index of `a` tags, which contain kind, pubkey and hash of an arbitrary
	// text, used to create an abstract reference for a multiplicity of replaceable event with a
	// kind number. These labels also appear as `d` tags in inbound references, see
	// IdxTagLetter.
	//
	// [ prefix ][ 2 bytes kind number ][ 8 bytes hash of pubkey ][ 8 bytes hash of label ][ serial]
	TagA

	// TagEvent is a reference to an event.
	//
	// [ prefix ][ 8 bytes truncated hash of event Id ][ 8 serial ]
	TagEvent

	// TagPubkey is a reference to a user's public key identifier (author).
	//
	// [ prefix ][ 8 bytes pubkey hash ][ 8 serial ]
	TagPubkey

	// TagHashtag is a reference to a hashtag, user-created and externally labeled short
	// subject names.
	//
	// [ prefix ][ 8 bytes hash of hashtag ][ 8 serial ]
	TagHashtag

	// TagIdentifier is a `d` tag identifier that creates an arbitrary label that can be used
	// to refer to an event. This is used for parameterized replaceable events to identify them
	// with `a` tags for reference.
	//
	// [ prefix ][ 8 byte hash of identifier ][ 8 serial ]
	TagIdentifier

	// TagLetter covers all other types of single letter mandatory indexed tags, including
	// such as `d` for identifiers and things like `m` for mimetype and other kinds of
	// references, the actual letter is the second byte. The value is a truncated 8 byte hash.
	//
	// [ prefix ][ letter ][ 8 bytes hash of value field of tag ][ 8 serial ]
	TagLetter

	// TagProtected is a special tag that indicates that this event should only be accepted
	// if published by an authed user with the matching public key.
	//
	// [ prefix ][ 8 byte hash of public key ][ 8 serial ]
	TagProtected

	// TagNonstandard is an index for index keys longer than 1 character, represented as an 8
	// byte truncated hash.
	//
	// [ prefix ][ 8 byte hash of key ][ 8 byte hash of value ][ 8 serial ]
	TagNonstandard

	// FulltextWord is a fulltext word index, the index contains the whole word. This will
	// also be searchable via the use of annotations in the filter search as whole match for the
	// word and any word containing the word (contains), and ^ prefix indicates a prefix match,
	// $ indicates a suffix match, and this index also contains a sequence number for proximity
	// filtering.
	//
	// [ prefix ][ full word ][ 4 bytes word position in content field ][ 8 serial ]
	FulltextWord

	// LastAccessed is an index that stores the last time the referenced event was returned
	// in a result.
	//
	// [ prefix ][ 8 serial ] [ last accessed timestamp 8 bytes ]
	LastAccessed

	// AccessCounter is a counter that is increased when the referenced event is a result in
	// a query. This can enable a frequency of access search or sort.
	//
	// [ prefix ][ 8 serial ] [ 8 bytes access counter ]
	AccessCounter
)

// Prefix is a map of the constant names above to the two byte prefix string of an index
// prefix.
var Prefix = map[int]I{
	Event:           "ev",
	Config:          "cf",
	Id:              "id",
	FullIndex:       "fi",
	Pubkey:          "pk",
	PubkeyCreatedAt: "pc",
	CreatedAt:       "ca",
	FirstSeen:       "fs",
	Kind:            "ki",
	TagA:            "ta",
	TagEvent:        "te",
	TagPubkey:       "tp",
	TagHashtag:      "tt",
	TagIdentifier:   "td",
	TagLetter:       "t*",
	TagProtected:    "t-",
	TagNonstandard:  "t?",
	FulltextWord:    "fw",
	LastAccessed:    "la",
	AccessCounter:   "ac",
}

// SplitLengthsFromPosition cuts a slice into segments of a given length starting after the 2
// byte prefix.
func SplitLengthsFromPosition(b []byte, positions ...int) (segments [][]byte, err error) {
	if len(positions) == 0 {
		err = errorf.E("must specify segment lengths")
		return
	}
	var total int
	for _, v := range positions {
		total += v
	}
	prev := 2
	if total > len(b)-prev {
		err = errorf.E("index is not long enough to split for this type %s %d - require %d", b[:2], len(b)-prev, total)
		return
	}
	for _, i := range positions {
		segments = append(segments, b[prev:prev+i])
		prev = i
	}
	return
}

const (
	serial    = 8
	idHash    = 8
	pubHash   = 8
	fullId    = 32
	kind      = 2
	createdAt = 8
	timeStamp = 8
	hash      = 8
	letter    = 1
	wordPos   = 4
)

type IdxEvent struct {
	Serial []byte
}

func NewIdxEvent(ser []byte) (idx []byte, err error) {
	if len(ser) != serial {
		err = errorf.E("invalid length for serial, got %d require %d", len(ser), serial)
		return
	}
	idx = append([]byte(Prefix[Event]), ser...)
	return
}

func IdxToEvent(idx []byte) (ie *IdxEvent, err error) {
	var segments [][]byte
	if segments, err = SplitIndex(idx); chk.E(err) {
		return
	}
	ie = &IdxEvent{Serial: segments[0]}
	return
}

type IdxFullIndex struct {
	FullId    []byte
	PubHash   []byte
	Kind      int
	CreatedAt timestamp.Timestamp
	Serial    []byte
}

func NewIdxFullIndex(fi, pk []byte, ki int, ca timestamp.Timestamp, ser []byte) (idx []byte, err error) {
	if len(fi) != fullId {
		err = errorf.E("invalid length for fullID, got %d require %d", len(fi), fullId)
		return
	}
	if len(pk) != schnorr.PubKeyBytesLen {
		err = errorf.E("invalid length for pubkey, got %d require %d", len(pk), schnorr.PubKeyBytesLen)
		return
	}
	if len(ser) != serial {
		err = errorf.E("invalid length for serial, got %d require %d", len(ser), serial)
		return
	}
	kiB := make([]byte, kind)
	binary.LittleEndian.PutUint16(kiB, uint16(ki))
	caB := make([]byte, createdAt)
	binary.LittleEndian.PutUint64(caB, uint64(ca))
	buf := new(bytes.Buffer)
	buf.Write([]byte(Prefix[FullIndex]))
	buf.Write(fi)
	buf.Write(helpers.Hash(pk)[:pubHash])
	buf.Write(kiB)
	buf.Write(caB)
	buf.Write(ser)
	idx = buf.Bytes()
	return
}

func IdxToFullIndex(idx []byte) (ie *IdxFullIndex, err error) {
	var segments [][]byte
	if segments, err = SplitIndex(idx); chk.E(err) {
		return
	}
	ie = &IdxFullIndex{
		FullId:    segments[0],
		PubHash:   segments[1],
		Kind:      int(binary.LittleEndian.Uint16(segments[2])),
		CreatedAt: timestamp.Timestamp(binary.LittleEndian.Uint64(segments[3])),
		Serial:    segments[4],
	}
	return
}

type IdxPubkey struct {
	PubHash []byte
	Serial  []byte
}

func NewIdxPubkey(pk, ser []byte) (idx []byte, err error) {
	if len(pk) != schnorr.PubKeyBytesLen {
		err = errorf.E("invalid length for pubkey, got %d require %d", len(pk), schnorr.PubKeyBytesLen)
		return
	}
	if len(ser) != serial {
		err = errorf.E("invalid length for serial, got %d require %d", len(ser), serial)
		return
	}
	buf := new(bytes.Buffer)
	buf.Write([]byte(Prefix[Pubkey]))
	buf.Write(helpers.Hash(pk)[:pubHash])
	buf.Write(ser)
	idx = buf.Bytes()
	return
}

func IdxToPubkey(idx []byte) (ie *IdxPubkey, err error) {
	var segments [][]byte
	if segments, err = SplitIndex(idx); chk.E(err) {
		return
	}
	ie = &IdxPubkey{
		PubHash: segments[0],
		Serial:  segments[1],
	}
	return
}

type IdxPubkeyCreatedAt struct {
	PubHash   []byte
	CreatedAt timestamp.Timestamp
	Serial    []byte
}

func NewIdxPubkeyCreatedAt(pk []byte, ca timestamp.Timestamp, ser []byte) (idx []byte, err error) {
	if len(pk) != schnorr.PubKeyBytesLen {
		err = errorf.E("invalid length for pubkey, got %d require %d", len(pk), schnorr.PubKeyBytesLen)
		return
	}
	if len(ser) != serial {
		err = errorf.E("invalid length for serial, got %d require %d", len(ser), serial)
		return
	}
	caB := make([]byte, createdAt)
	binary.LittleEndian.PutUint64(caB, uint64(ca))
	buf := new(bytes.Buffer)
	buf.Write([]byte(Prefix[PubkeyCreatedAt]))
	buf.Write(helpers.Hash(pk)[:pubHash])
	buf.Write(caB)
	buf.Write(ser)
	idx = buf.Bytes()
	return
}

func IdxToPubkeyCreatedAt(idx []byte) (ie *IdxPubkeyCreatedAt, err error) {
	var segments [][]byte
	if segments, err = SplitIndex(idx); chk.E(err) {
		return
	}
	ie = &IdxPubkeyCreatedAt{
		PubHash:   segments[0],
		CreatedAt: timestamp.Timestamp(binary.LittleEndian.Uint64(segments[1])),
		Serial:    segments[2],
	}
	return
}

type IdxCreatedAt struct {
	CreatedAt timestamp.Timestamp
	Serial    []byte
}

func NewIdxCreatedAt(ca timestamp.Timestamp, ser []byte) (idx []byte, err error) {
	if len(ser) != serial {
		err = errorf.E("invalid length for serial, got %d require %d", len(ser), serial)
		return
	}
	caB := make([]byte, createdAt)
	binary.LittleEndian.PutUint64(caB, uint64(ca))
	buf := new(bytes.Buffer)
	buf.Write([]byte(Prefix[CreatedAt]))
	buf.Write(caB)
	buf.Write(ser)
	idx = buf.Bytes()
	return
}

func IdxToCreatedAt(idx []byte) (ie *IdxCreatedAt, err error) {
	var segments [][]byte
	if segments, err = SplitIndex(idx); chk.E(err) {
		return
	}
	ie = &IdxCreatedAt{
		CreatedAt: timestamp.Timestamp(binary.LittleEndian.Uint64(segments[0])),
		Serial:    segments[1],
	}
	return
}

type IdxFirstSeen struct {
	Serial    []byte
	Timestamp timestamp.Timestamp
}

func NewIdxFirstSeen(ser []byte, fs timestamp.Timestamp) (idx []byte, err error) {
	if len(ser) != serial {
		err = errorf.E("invalid length for serial, got %d require %d", len(ser), serial)
		return
	}
	fsB := make([]byte, timeStamp)
	binary.LittleEndian.PutUint64(fsB, uint64(fs))
	buf := new(bytes.Buffer)
	buf.Write([]byte(Prefix[FirstSeen]))
	buf.Write(fsB)
	buf.Write(ser)
	idx = buf.Bytes()
	return
}

func IdxToFirstSeen(idx []byte) (ie *IdxFirstSeen, err error) {
	var segments [][]byte
	if segments, err = SplitIndex(idx); chk.E(err) {
		return
	}
	ie = &IdxFirstSeen{
		Serial:    segments[0],
		Timestamp: timestamp.Timestamp(binary.LittleEndian.Uint64(segments[1])),
	}
	return
}

type IdxKind struct {
	Kind   int
	Serial []byte
}

func NewIdxKind(ki int, ser []byte) (idx []byte, err error) {
	if len(ser) != serial {
		err = errorf.E("invalid length for serial, got %d require %d", len(ser), serial)
		return
	}
	kiB := make([]byte, kind)
	binary.LittleEndian.PutUint16(kiB, uint16(ki))
	buf := new(bytes.Buffer)
	buf.Write([]byte(Prefix[Kind]))
	buf.Write(kiB)
	buf.Write(ser)
	idx = buf.Bytes()
	return
}

func IdxToKind(idx []byte) (ie *IdxKind, err error) {
	var segments [][]byte
	if segments, err = SplitIndex(idx); chk.E(err) {
		return
	}
	ie = &IdxKind{
		Kind:   int(binary.LittleEndian.Uint16(segments[0])),
		Serial: segments[1],
	}
	return
}

type IdxTagA struct {
	Kind           int
	PubHash        []byte
	IdentifierHash []byte
	Serial         []byte
}

func NewIdxTagA(ki int, pk, identifier, ser []byte) (idx []byte, err error) {
	if len(pk) != schnorr.PubKeyBytesLen {
		err = errorf.E("invalid length for pubkey, got %d require %d", len(pk), schnorr.PubKeyBytesLen)
		return
	}
	if len(ser) != serial {
		err = errorf.E("invalid length for serial, got %d require %d", len(ser), serial)
		return
	}
	kiB := make([]byte, kind)
	binary.LittleEndian.PutUint16(kiB, uint16(ki))
	buf := new(bytes.Buffer)
	buf.Write([]byte(Prefix[TagA]))
	buf.Write(kiB)
	buf.Write(helpers.Hash(pk)[:pubHash])
	buf.Write(helpers.Hash(identifier)[:hash])
	buf.Write(ser)
	idx = buf.Bytes()
	return
}

func IdxToTagA(idx []byte) (ie *IdxTagA, err error) {
	var segments [][]byte
	if segments, err = SplitIndex(idx); chk.E(err) {
		return
	}
	ie = &IdxTagA{
		Kind:           int(binary.LittleEndian.Uint16(segments[0])),
		PubHash:        segments[1],
		IdentifierHash: segments[2],
		Serial:         segments[3],
	}
	return
}

type IdxTagEvent struct {
	IdHash []byte
	Serial []byte
}

func NewIdxTagEvent(id, ser []byte) (idx []byte, err error) {
	if len(id) != 32 {
		err = errorf.E("invalid length for event id, got %d require %d", len(id), 32)
		return
	}
	if len(ser) != serial {
		err = errorf.E("invalid length for serial, got %d require %d", len(ser), serial)
		return
	}
	buf := new(bytes.Buffer)
	buf.Write([]byte(Prefix[TagEvent]))
	buf.Write(helpers.Hash(id)[:hash])
	buf.Write(ser)
	idx = buf.Bytes()
	return
}

func IdxToTagEvent(idx []byte) (ie *IdxTagEvent, err error) {
	var segments [][]byte
	if segments, err = SplitIndex(idx); chk.E(err) {
		return
	}
	ie = &IdxTagEvent{
		IdHash: segments[0],
		Serial: segments[1],
	}
	return
}

type IdxTagPubkey struct {
	PubHash []byte
	Serial  []byte
}

func NewIdxTagPubkey(pk []byte, ser []byte) (idx []byte, err error) {
	if len(pk) != schnorr.PubKeyBytesLen {
		err = errorf.E("invalid length for pubkey, got %d require %d", len(pk), schnorr.PubKeyBytesLen)
		return
	}
	if len(ser) != serial {
		err = errorf.E("invalid length for serial, got %d require %d", len(ser), serial)
		return
	}
	buf := new(bytes.Buffer)
	buf.Write([]byte(Prefix[TagPubkey]))
	buf.Write(helpers.Hash(pk)[:pubHash])
	buf.Write(ser)
	idx = buf.Bytes()
	return
}

func IdxToTagPubkey(idx []byte) (ie *IdxTagPubkey, err error) {
	var segments [][]byte
	if segments, err = SplitIndex(idx); chk.E(err) {
		return
	}
	ie = &IdxTagPubkey{
		PubHash: segments[0],
		Serial:  segments[1],
	}
	return
}

type IdxTagHashtag struct {
	Hashtag []byte
	Serial  []byte
}

func NewIdxTagHashtag(ht []byte, ser []byte) (idx []byte, err error) {
	if len(ser) != serial {
		err = errorf.E("invalid length for serial, got %d require %d", len(ser), serial)
		return
	}
	buf := new(bytes.Buffer)
	buf.Write([]byte(Prefix[TagHashtag]))
	buf.Write(helpers.Hash(ht)[:pubHash])
	buf.Write(ser)
	idx = buf.Bytes()
	return
}

func IdxToTagHashtag(idx []byte) (ie *IdxTagHashtag, err error) {
	var segments [][]byte
	if segments, err = SplitIndex(idx); chk.E(err) {
		return
	}
	ie = &IdxTagHashtag{
		Hashtag: segments[0],
		Serial:  segments[1],
	}
	return
}

type IdxTagIdentifier struct {
	IdentifierHash []byte
	Serial         []byte
}

func NewIdxTagIdentifier(ident, ser []byte) (idx []byte, err error) {
	if len(ser) != serial {
		err = errorf.E("invalid length for serial, got %d require %d", len(ser), serial)
		return
	}
	buf := new(bytes.Buffer)
	buf.Write([]byte(Prefix[TagIdentifier]))
	buf.Write(helpers.Hash(ident)[:hash])
	buf.Write(ser)
	idx = buf.Bytes()
	return
}

func IdxToTagIdentifier(idx []byte) (ie *IdxTagIdentifier, err error) {
	var segments [][]byte
	if segments, err = SplitIndex(idx); chk.E(err) {
		return
	}
	ie = &IdxTagIdentifier{
		IdentifierHash: segments[0],
		Serial:         segments[1],
	}
	return
}

type IdxTagLetter struct {
	Letter    byte
	ValueHash []byte
	Serial    []byte
}

func NewIdxTagLetter(let byte, val []byte, ser []byte) (idx []byte, err error) {
	if len(ser) != serial {
		err = errorf.E("invalid length for serial, got %d require %d", len(ser), serial)
		return
	}
	buf := new(bytes.Buffer)
	buf.Write([]byte(Prefix[TagLetter]))
	buf.Write([]byte{let})
	buf.Write(helpers.Hash(val)[:hash])
	buf.Write(ser)
	idx = buf.Bytes()
	return
}

func IdxToTagLetter(idx []byte) (ie *IdxTagLetter, err error) {
	var segments [][]byte
	if segments, err = SplitIndex(idx); chk.E(err) {
		return
	}
	ie = &IdxTagLetter{
		Letter:    segments[0][0],
		ValueHash: segments[1],
		Serial:    segments[2],
	}
	return
}

type IdxTagProtected struct {
	Serial []byte
}

func NewIdxTagProtected(ser []byte) (idx []byte, err error) {
	if len(ser) != serial {
		err = errorf.E("invalid length for serial, got %d require %d", len(ser), serial)
		return
	}
	buf := new(bytes.Buffer)
	buf.Write([]byte(Prefix[TagProtected]))
	buf.Write(ser)
	idx = buf.Bytes()
	return
}

func IdxToTagProtected(idx []byte) (ie *IdxTagProtected, err error) {
	var segments [][]byte
	if segments, err = SplitIndex(idx); chk.E(err) {
		return
	}
	ie = &IdxTagProtected{
		Serial: segments[0],
	}
	return
}

type IdxTagNonstandard struct {
	KeyHash   []byte
	ValueHash []byte
	Serial    []byte
}

func NewIdxTagNonstandard(key, val []byte, ser []byte) (idx []byte, err error) {
	if len(ser) != serial {
		err = errorf.E("invalid length for serial, got %d require %d", len(ser), serial)
		return
	}
	buf := new(bytes.Buffer)
	buf.Write([]byte(Prefix[TagNonstandard]))
	buf.Write(helpers.Hash(key)[:hash])
	buf.Write(helpers.Hash(val)[:hash])
	buf.Write(ser)
	idx = buf.Bytes()
	return
}

func IdxToTagNonstandard(idx []byte) (ie *IdxTagNonstandard, err error) {
	var segments [][]byte
	if segments, err = SplitIndex(idx); chk.E(err) {
		return
	}
	ie = &IdxTagNonstandard{
		KeyHash:   segments[0],
		ValueHash: segments[1],
		Serial:    segments[2],
	}
	return
}

type IdxFullTextWord struct {
	Word    []byte
	WordPos uint32
	Serial  []byte
}

func NewIdxFulltextWord(word []byte, wordPos uint32, ser []byte) (idx []byte, err error) {
	if len(ser) != serial {
		err = errorf.E("invalid length for serial, got %d require %d", len(ser), serial)
		return
	}
	wpB := make([]byte, kind)
	binary.LittleEndian.PutUint32(wpB, wordPos)
	buf := new(bytes.Buffer)
	buf.Write([]byte(Prefix[FulltextWord]))
	buf.Write(word)
	buf.Write(wpB)
	buf.Write(ser)
	idx = buf.Bytes()
	return
}

func IdxToFullTextWord(idx []byte) (ie *IdxFullTextWord, err error) {
	var segments [][]byte
	if segments, err = SplitIndex(idx); chk.E(err) {
		return
	}
	ie = &IdxFullTextWord{
		Word:    segments[0],
		WordPos: binary.LittleEndian.Uint32(segments[1]),
		Serial:  segments[2],
	}
	return
}

type IdxLastAccessed struct {
	Serial []byte
}

func NewIdxLastAccessed(ser []byte) (idx []byte, err error) {
	if len(ser) != serial {
		err = errorf.E("invalid length for serial, got %d require %d", len(ser), serial)
		return
	}
	idx = append([]byte(Prefix[LastAccessed]), ser...)
	return
}

func IdxToLastAccessed(idx []byte) (ie *IdxLastAccessed, err error) {
	var segments [][]byte
	if segments, err = SplitIndex(idx); chk.E(err) {
		return
	}
	ie = &IdxLastAccessed{
		Serial: segments[0],
	}
	return
}

type IdxAccessCounter struct {
	Serial []byte
}

func NewIdxAccessCounter(ser []byte) (idx []byte, err error) {
	if len(ser) != serial {
		err = errorf.E("invalid length for serial, got %d require %d", len(ser), serial)
		return
	}
	idx = append([]byte(Prefix[AccessCounter]), ser...)
	return
}

func IdxToAccessCounter(idx []byte) (ie *IdxAccessCounter, err error) {
	var segments [][]byte
	if segments, err = SplitIndex(idx); chk.E(err) {
		return
	}
	ie = &IdxAccessCounter{
		Serial: segments[0],
	}
	return
}

func SplitIndex(idx []byte) (segments [][]byte, err error) {
	switch I(idx[:2]) {
	case Prefix[Event]:
		segments, err = SplitLengthsFromPosition(idx, serial)
	case Prefix[Config]:
		return
	case Prefix[Id]:
		segments, err = SplitLengthsFromPosition(idx, idHash, serial)
	case Prefix[FullIndex]:
		segments, err = SplitLengthsFromPosition(idx, fullId, pubHash, kind, createdAt, serial)
	case Prefix[Pubkey]:
		segments, err = SplitLengthsFromPosition(idx, pubHash, serial)
	case Prefix[PubkeyCreatedAt]:
		segments, err = SplitLengthsFromPosition(idx, pubHash, createdAt, serial)
	case Prefix[CreatedAt]:
		segments, err = SplitLengthsFromPosition(idx, createdAt, serial)
	case Prefix[FirstSeen]:
		segments, err = SplitLengthsFromPosition(idx, serial, timeStamp)
	case Prefix[Kind]:
		segments, err = SplitLengthsFromPosition(idx, kind, serial)
	case Prefix[TagA]:
		segments, err = SplitLengthsFromPosition(idx, kind, pubHash, hash, serial)
	case Prefix[TagEvent]:
		segments, err = SplitLengthsFromPosition(idx, idHash, serial)
	case Prefix[TagPubkey]:
		segments, err = SplitLengthsFromPosition(idx, pubHash, serial)
	case Prefix[TagHashtag]:
		segments, err = SplitLengthsFromPosition(idx, hash, serial)
	case Prefix[TagIdentifier]:
		segments, err = SplitLengthsFromPosition(idx, hash, serial)
	case Prefix[TagLetter]:
		segments, err = SplitLengthsFromPosition(idx, letter, hash, serial)
	case Prefix[TagProtected]:
		segments, err = SplitLengthsFromPosition(idx, serial)
	case Prefix[TagNonstandard]:
		segments, err = SplitLengthsFromPosition(idx, hash, hash, serial)
	case Prefix[FulltextWord]:
		wordLen := len(idx) - 2 - wordPos - serial
		segments, err = SplitLengthsFromPosition(idx, wordLen, wordPos, serial)
	case Prefix[LastAccessed]:
		segments, err = SplitLengthsFromPosition(idx, serial)
	case Prefix[AccessCounter]:
		segments, err = SplitLengthsFromPosition(idx, serial)
	}
	return
}
