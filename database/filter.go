package database

import (
	"math"

	"x.realy.lol/chk"
	"x.realy.lol/database/indexes/types/varint"
	"x.realy.lol/filter"
	"x.realy.lol/hex"
	"x.realy.lol/timestamp"
)

type Bitfield byte

const (
	hasIds     Bitfield = 1
	hasKinds   Bitfield = 2
	hasAuthors Bitfield = 4
	hasTags    Bitfield = 8
	hasSince   Bitfield = 16
	hasUntil   Bitfield = 32
	hasLimit   Bitfield = 64
	hasSearch  Bitfield = 128
)

func ToBitfield(f *filter.F) (b Bitfield) {
	if len(f.Ids) != 0 {
		b += hasIds
	}
	if len(f.Kinds) != 0 {
		b += hasKinds
	}
	if len(f.Authors) != 0 {
		b += hasAuthors
	}
	if len(f.Kinds) != 0 {
		b += hasTags
	}
	if f.Since != nil {
		b += hasSince
	}
	if f.Until != nil {
		b += hasUntil
	}
	if f.Limit != nil {
		b += hasLimit
	}
	if f.Search != "" {
		b += hasSearch
	}
	return
}

func (d *D) Filter(f filter.F) (evSerials []*varint.V, err error) {
	bf := ToBitfield(&f)
	// first, if there is Ids these override everything else
	if bf&hasIds != 0 {
		for _, v := range f.Ids {
			var id []byte
			if id, err = hex.Dec(v); chk.E(err) {
				// just going to ignore it i guess
				continue
			}
			var ev *varint.V
			if ev, err = d.FindEventSerialById(id); chk.E(err) {
				// just going to ignore it i guess
				continue
			}
			evSerials = append(evSerials, ev)
		}
		return
	}
	var since, until timestamp.Timestamp
	if bf&hasSince != 0 {
		since = *f.Since
	}
	if bf&hasUntil != 0 {
		until = *f.Until
	} else {
		until = math.MaxInt64
	}
	// next, check for filters that only have since and/or until
	if bf&hasSince != 0 || bf&hasUntil != 0 {
		if evSerials, err = d.GetEventSerialsByCreatedAtRange(since, until); chk.E(err) {
			return
		}
		return
	}
	// next, kinds
	if bf&hasKinds == hasKinds && ^hasKinds&bf == 0 {
		if evSerials, err = d.GetEventSerialsByKindsCreatedAtRange(f.Kinds, since, until); chk.E(err) {
			return
		}
		return
	}
	// next authors
	if bf&hasAuthors == hasAuthors && ^hasAuthors&bf == 0 {
		if evSerials, err = d.GetEventSerialsByAuthorsCreatedAtRange(f.Authors, since, until); chk.E(err) {
			return
		}
		return
	}
	// next authors/kinds
	ak := hasAuthors + hasKinds
	if bf&(ak) == ak && ^ak&bf == 0 {
		if evSerials, err = d.GetEventSerialsByKindsAuthorsCreatedAtRange(f.Kinds, f.Authors, since, until); chk.E(err) {
			return
		}
		return
	}
	// next authors/tags
	at := hasAuthors + hasTags
	if bf&(at) == at && ^at&bf == 0 {
		if evSerials, err = d.GetEventSerialsByKindsAuthorsCreatedAtRange(f.Kinds, f.Authors, since, until); chk.E(err) {
			return
		}
		return
	}

	// next kinds/tags

	// next
	return
}
