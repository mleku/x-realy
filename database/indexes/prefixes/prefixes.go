package prefixes

import (
	"io"
)

const Len = 2

type I string

const (
	// Event is the whole event stored in binary format
	//
	//   [ prefix ][ 8 byte serial ] [ event in binary format ]
	Event = iota

	// Config is a singular record containing a free-form configuration in JSON format
	//
	// [ prefix ] [ configuration in JSON format ]
	Config

	// Id contains a truncated 8 byte hash of an event index. This is the secondary key of an
	// event, the primary key is the serial found in the Event.
	//
	// [ prefix ][ 8 bytes truncated hash of Id ][ 8 serial ]
	Id

	// FullIndex is an index designed to enable sorting and filtering of results found via
	// other indexes, without having to decode the event.
	//
	// [ prefix ][ 8 serial ][ 32 bytes full event ID ][ 8 bytes truncated hash of pubkey ][ 2 bytes kind ][ 8 bytes created_at timestamp ]
	FullIndex

	// ------------------------------------------------------------------------
	//
	// The following are search indexes. This first category are primarily for kind, pubkey and
	// created_at timestamps. These compose a set of 3 primary indexes alone, two that combine
	// with the timestamp, and a third that combines all three, covering every combination of
	// these.
	_

	// Pubkey is an index for searching for events authored by a pubkey.
	//
	// [ prefix ][ 8 bytes truncated hash of pubkey ][ 8 serial ]
	Pubkey

	// Kind is an index of event kind numbers.
	//
	// [ prefix ][ 2 bytes kind number ][ 8 serial ]
	Kind

	// CreatedAt is an index that allows search the timestamp on the event.
	//
	// [ prefix ][ created_at 8 bytes timestamp ][ 8 serial ]
	CreatedAt

	// PubkeyCreatedAt is a composite index that allows search by pubkey filtered by
	// created_at.
	//
	// [ prefix ][ 8 bytes truncated hash of pubkey ][ 8 bytes created_at ][ 8 serial ]
	PubkeyCreatedAt

	// KindCreatedAt is an index of kind and created_at timestamp.
	//
	// [ prefix ][ 2 bytes kind number ][ created_at 8 bytes timestamp ][ 8 bytes serial ]
	KindCreatedAt

	// KindPubkeyCreatedAt is an index of kind and created_at timestamp.
	//
	// [ prefix ][ 2 bytes kind number ][ 8 bytes hash of pubkey ][ created_at 8 bytes timestamp ][ 8 bytes serial ]
	KindPubkeyCreatedAt

	// ------------------------------------------------------------------------
	//
	// The following are search indexes for tags, which are references to other categories,
	// including events, replaceable event identities (d tags), public keys, hashtags, and
	// arbitrary other kinds of keys including standard single letter and nonstandard word keys.
	//
	// Combining them with the previous set of 6 indexes involves using one query from the
	// previous section according to the filter, and one or more of these tag indexes, to
	// acquire a list of event serials from each query, and then intersecting the result sets
	// from each one to yield the matches.
	_

	// TagA is an index of `a` tags, which contain kind, pubkey and hash of an arbitrary
	// text, used to create an abstract reference for a multiplicity of replaceable event with a
	// kind number. These labels also appear as `d` tags in inbound references, see
	// IdxTagIdentifier.
	//
	// [ prefix ][ 2 bytes kind number ][ 8 bytes hash of pubkey ][ 8 bytes hash of label ][ serial]
	TagA

	// TagIdentifier is a `d` tag identifier that creates an arbitrary label that can be used
	// to refer to an event. This is used for parameterized replaceable events to identify them
	// with `a` tags for reference.
	//
	// [ prefix ][ 8 byte hash of identifier ][ 8 serial ]
	TagIdentifier

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

	// ------------------------------------------------------------------------
	_

	// FulltextWord is a fulltext word index, the index contains the whole word. This will
	// also be searchable via the use of annotations in the filter search as whole match for the
	// word and any word containing the word (contains), and ^ prefix indicates a prefix match,
	// $ indicates a suffix match, and this index also contains a sequence number for proximity
	// filtering.
	//
	// [ prefix ][ varint word len ][ full word ][ 4 bytes word position in content field ][ 8 serial ]
	FulltextWord

	// ------------------------------------------------------------------------
	//
	// The following keys are event metadata that are needed to enable other types of
	// functionality such as garbage collection and metadata queries.
	_

	// FirstSeen is an index that records the timestamp of when the event was first seen.
	//
	// [ prefix ][ 8 serial ][ 8 byte timestamp ]
	FirstSeen

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

func (i I) Write(w io.Writer) (n int, err error) { return w.Write([]byte(i)) }

func Prefix(prf int) (i I) {
	switch prf {
	case Event:
		return "ev"
	case Config:
		return "cf"
	case Id:
		return "id"
	case FullIndex:
		return "fi"
	case Pubkey:
		return "pk"
	case PubkeyCreatedAt:
		return "pc"
	case CreatedAt:
		return "ca"
	case FirstSeen:
		return "fs"
	case Kind:
		return "ki"
	case KindCreatedAt:
		return "kc"
	case KindPubkeyCreatedAt:
		return "kp"
	case TagA:
		return "ta"
	case TagEvent:
		return "te"
	case TagPubkey:
		return "tp"
	case TagHashtag:
		return "tt"
	case TagIdentifier:
		return "td"
	case TagLetter:
		return "t*"
	case TagProtected:
		return "t-"
	case TagNonstandard:
		return "t?"
	case FulltextWord:
		return "fw"
	case LastAccessed:
		return "la"
	case AccessCounter:
		return "ac"
	}
	return
}
