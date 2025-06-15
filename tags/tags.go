package tags

import (
	"encoding/json"
	"errors"
	"iter"
	"slices"
	"strings"

	"x.realy.lol/chk"
	"x.realy.lol/ec/schnorr"
	"x.realy.lol/helpers"
	"x.realy.lol/hex"
	"x.realy.lol/ints"
	"x.realy.lol/normalize"
)

type Tag []string

// StartsWith checks if a tag contains a prefix.
// for example,
//
//	["p", "abcdef...", "wss://relay.com"]
//
// would match against
//
//	["p", "abcdef..."]
//
// or even
//
//	["p", "abcdef...", "wss://"]
func (tag Tag) StartsWith(prefix []string) bool {
	prefixLen := len(prefix)

	if prefixLen > len(tag) {
		return false
	}
	// check initial elements for equality
	for i := 0; i < prefixLen-1; i++ {
		if prefix[i] != tag[i] {
			return false
		}
	}
	// check last element just for a prefix
	return strings.HasPrefix(tag[prefixLen-1], prefix[prefixLen-1])
}

func (tag Tag) Key() string {
	if len(tag) > 0 {
		return tag[0]
	}
	return ""
}

func (tag Tag) Value() string {
	if len(tag) > 1 {
		return tag[1]
	}
	return ""
}

func (tag Tag) Relay() string {
	if len(tag) > 2 && (tag[0] == "e" || tag[0] == "p") {
		return normalize.Url(tag[2])
	}
	return ""
}

type Tags []Tag

// GetD gets the first "d" tag (for parameterized replaceable events) value or ""
func (tags Tags) GetD() string {
	for _, v := range tags {
		if v.StartsWith([]string{"d", ""}) {
			return v[1]
		}
	}
	return ""
}

// GetFirst gets the first tag in tags that matches the prefix, see [Tag.StartsWith]
func (tags Tags) GetFirst(tagPrefix []string) *Tag {
	for _, v := range tags {
		if v.StartsWith(tagPrefix) {
			return &v
		}
	}
	return nil
}

// GetLast gets the last tag in tags that matches the prefix, see [Tag.StartsWith]
func (tags Tags) GetLast(tagPrefix []string) *Tag {
	for i := len(tags) - 1; i >= 0; i-- {
		v := tags[i]
		if v.StartsWith(tagPrefix) {
			return &v
		}
	}
	return nil
}

// GetAll gets all the tags that match the prefix, see [Tag.StartsWith]
func (tags Tags) GetAll(tagPrefix []string) Tags {
	result := make(Tags, 0, len(tags))
	for _, v := range tags {
		if v.StartsWith(tagPrefix) {
			result = append(result, v)
		}
	}
	return result
}

func (tags Tags) GetAllExactKeys(key string) Tags {
	result := make(Tags, 0, len(tags))
	for _, v := range tags {
		if v.StartsWith([]string{key}) {
			if v.Key() == key {
				result = append(result, v)
			}
		}
	}
	return result
}

// All returns an iterator for all the tags that match the prefix, see [Tag.StartsWith]
func (tags Tags) All(tagPrefix []string) iter.Seq2[int, Tag] {
	return func(yield func(int, Tag) bool) {
		for i, v := range tags {
			if v.StartsWith(tagPrefix) {
				if !yield(i, v) {
					break
				}
			}
		}
	}
}

// FilterOut returns a new slice with only the elements that match the prefix, see [Tag.StartsWith]
func (tags Tags) FilterOut(tagPrefix []string) Tags {
	filtered := make(Tags, 0, len(tags))
	for _, v := range tags {
		if !v.StartsWith(tagPrefix) {
			filtered = append(filtered, v)
		}
	}
	return filtered
}

// FilterOutInPlace removes all tags that match the prefix, but potentially reorders the tags in unpredictable ways, see [Tag.StartsWith]
func (tags *Tags) FilterOutInPlace(tagPrefix []string) {
	for i := 0; i < len(*tags); i++ {
		tag := (*tags)[i]
		if tag.StartsWith(tagPrefix) {
			// remove this by swapping the last tag into this place
			last := len(*tags) - 1
			(*tags)[i] = (*tags)[last]
			*tags = (*tags)[0:last]
			i-- // this is so we can match this just swapped item in the next iteration
		}
	}
}

// AppendUnique appends a tag if it doesn't exist yet, otherwise does nothing.
// the uniqueness comparison is done based only on the first 2 elements of the tag.
func (tags Tags) AppendUnique(tag Tag) Tags {
	n := len(tag)
	if n > 2 {
		n = 2
	}

	if tags.GetFirst(tag[:n]) == nil {
		return append(tags, tag)
	}
	return tags
}

func (t *Tags) Scan(src any) error {
	var jtags []byte

	switch v := src.(type) {
	case []byte:
		jtags = v
	case string:
		jtags = []byte(v)
	default:
		return errors.New("couldn't scan tags, it's not a json string")
	}

	json.Unmarshal(jtags, &t)
	return nil
}

func (tags Tags) ContainsAny(tagName string, values []string) bool {
	for _, tag := range tags {
		if len(tag) < 2 {
			continue
		}

		if tag[0] != tagName {
			continue
		}

		if slices.Contains(values, tag[1]) {
			return true
		}
	}

	return false
}

// Marshal Tag. Used for Serialization so string escaping should be as in RFC8259.
func (tag Tag) marshalTo(dst []byte) []byte {
	dst = append(dst, '[')
	for i, s := range tag {
		if i > 0 {
			dst = append(dst, ',')
		}
		dst = helpers.EscapeString(dst, s)
	}
	dst = append(dst, ']')
	return dst
}

// MarshalTo appends the JSON encoded byte of Tags as [][]string to dst.
// String escaping is as described in RFC8259.
func (tags Tags) marshalTo(dst []byte) []byte {
	dst = append(dst, '[')
	for i, tag := range tags {
		if i > 0 {
			dst = append(dst, ',')
		}
		dst = tag.marshalTo(dst)
	}
	dst = append(dst, ']')
	return dst
}

type Tag_a struct {
	Kind   int
	Pubkey []byte
	Ident  string
}

func (tags Tags) Get_a_Tags() (atags []Tag_a) {
	a := tags.GetAll([]string{"a"})
	var err error
	if len(a) > 0 {
		for _, v := range a {
			if v[0] == "a" && len(v) > 1 {
				var atag Tag_a
				if atag, err = Decode_a_Tag(v[1]); chk.E(err) {
					continue
				}
				atags = append(atags, atag)
			}
		}
	}
	return
}

func Decode_a_Tag(a string) (ta Tag_a, err error) {
	// try to split it
	parts := strings.Split(a, ":")
	// there must be a kind first
	ki := ints.New(0)
	if _, err = ki.Unmarshal([]byte(parts[0])); chk.E(err) {
		return
	}
	ta = Tag_a{
		Kind: int(ki.Uint16()),
	}
	if len(parts) < 2 {
		return
	}
	// next must be a pubkey
	if len(parts[1]) != 2*schnorr.PubKeyBytesLen {
		return
	}
	var pk []byte
	if pk, err = hex.Dec(parts[1]); err != nil {
		return
	}
	ta.Pubkey = pk
	// there possibly can be nothing after this
	if len(parts) >= 3 {
		// third part is the identifier (d tag)
		ta.Ident = parts[2]
	}
	return
}
