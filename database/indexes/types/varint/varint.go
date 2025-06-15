package varint

import (
	"bytes"
	"io"

	"x.realy.lol/chk"
	"x.realy.lol/varint"
)

type V struct{ val uint64 }

type S []*V

func New() (s *V) { return &V{} }

func (vi *V) FromUint64(ser uint64) {
	vi.val = ser
	return
}

func FromBytes(ser []byte) (s *V, err error) {
	s = &V{}
	if s.val, err = varint.Decode(bytes.NewBuffer(ser)); chk.E(err) {
		return
	}
	return
}

func (vi *V) ToUint64() (ser uint64) { return vi.val }

func (vi *V) ToInt() (ser int) { return int(vi.val) }

func (vi *V) ToUint32() (v uint32) { return uint32(vi.val) }

func (vi *V) Bytes() (b []byte) {
	buf := new(bytes.Buffer)
	varint.Encode(buf, vi.val)
	return
}

func (vi *V) MarshalWrite(w io.Writer) (err error) {
	varint.Encode(w, vi.val)
	return
}

func (vi *V) UnmarshalRead(r io.Reader) (err error) {
	vi.val, err = varint.Decode(r)
	return
}

// DeduplicateInOrder removes duplicates from a slice of V.
func DeduplicateInOrder(s S) (v S) {
	// for larger slices, this uses a lot less memory, at the cost of slower execution.
	if len(s) > 10000 {
	skip:
		for i, sa := range s {
			for j, sb := range s {
				if i != j && sa.val == sb.val {
					continue skip
				}
			}
			v = append(v, sa)
		}
	} else {
		// for small slices, this is faster but uses more memory.
		seen := map[uint64]*V{}
		for _, val := range s {
			if _, ok := seen[val.val]; !ok {
				v = append(v, val)
				seen[val.val] = val
			}
		}
	}
	return
}

// Intersect deduplicates and performs a set intersection on two slices.
func Intersect(a, b []*V) (sers []*V) {
	// first deduplicate to eliminate unnecessary iterations
	a = DeduplicateInOrder(a)
	b = DeduplicateInOrder(b)
	for _, as := range a {
		for _, bs := range b {
			if as.val == bs.val {
				// if the match is found, add to the result and move to the next candidate from
				// the "a" serial list.
				sers = append(sers, as)
				break
			}
		}
	}
	return
}
