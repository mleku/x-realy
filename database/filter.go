package database

import (
	"bytes"
	"math"
	"sort"

	"x.realy.lol/chk"
	"x.realy.lol/database/indexes"
	"x.realy.lol/database/indexes/types/pubhash"
	"x.realy.lol/database/indexes/types/varint"
	"x.realy.lol/filter"
	"x.realy.lol/hex"
	"x.realy.lol/log"
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
	if len(f.Tags) != 0 {
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

// Filter runs a nip-01 type query on a provided filter and returns the database serial keys of
// the matching events, excluding a list of authors also provided from the result.
func (d *D) Filter(f filter.F, exclude []*pubhash.T) (evSerials varint.S, err error) {
	var evs varint.S
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
			evs = append(evs, ev)
		}
		return
	}
	var since, until *timestamp.Timestamp
	if bf&hasSince == 0 || bf&hasUntil == 0 {
		if bf&hasSince != 0 {
			since = f.Since
		}
		if bf&hasUntil != 0 {
			until = f.Until
		} else {
			m := timestamp.Timestamp(math.MaxInt64)
			until = &m
		}
	}
	limit := f.Limit
	var postLimit bool
	if limit != nil {
		// put a reasonable cap on unlimited. the actual results may be a lot less for composite
		// searches that intersect with tags.
		limit = filter.IntToPointer(10000)
	} else {
		// this means trim the result at the end before returning.
		postLimit = true
	}
	log.I.F("%b %b", bf, bf&(hasSince+hasUntil+hasLimit))
	bf = bf &^ hasLimit
	// next, check for filters that only have since and/or until
	if bf&(hasSince+hasUntil) != 0 && ^(hasUntil+hasSince)&bf == 0 {
		if evs, err = d.GetEventSerialsByCreatedAtRange(since, until, limit, postLimit); chk.E(err) {
			return
		}
		goto done
	}
	// next, kinds
	if bf&hasKinds == hasKinds && ^hasKinds&bf == 0 {
		log.I.F("kinds")
		if evs, err = d.GetEventSerialsByKinds(f.Kinds, f.Limit); chk.E(err) {
			return
		}
		goto done
	}
	// next, kinds/created_at
	if (bf&hasKinds+hasSince == hasKinds+hasSince ||
		bf&hasKinds+hasUntil == hasKinds+hasUntil ||
		bf&hasKinds+hasUntil+hasSince == hasKinds+hasUntil+hasSince) &&
		^(hasKinds+hasUntil+hasSince)&bf == 0 {
		if evs, err = d.GetEventSerialsByKindsCreatedAtRange(f.Kinds, since, until, limit); chk.E(err) {
			return
		}
		goto done
	}
	// next authors
	if bf&hasAuthors == hasAuthors && ^hasAuthors&bf == 0 {
		if evs, err = d.GetEventSerialsByAuthorsCreatedAtRange(f.Authors, since, until, limit); chk.E(err) {
			return
		}
		goto done
	}
	// next authors/kinds

	if ak := hasAuthors + hasKinds; bf&(ak) == ak && ^ak&bf == 0 {
		if evs, err = d.GetEventSerialsByKindsAuthorsCreatedAtRange(f.Kinds, f.Authors, since, until, limit); chk.E(err) {
			return
		}
		goto done
	}
	// if there is tags, assemble them into an array of tags with the
	if bf&hasTags != 0 && bf&^hasTags == 0 {
		if evs, err = d.GetEventSerialsByTagsCreatedAtRange(f.Tags, limit); chk.E(err) {

		}
	}
	// next authors/tags
	if at := hasAuthors + hasTags; bf&(at) == at && ^at&bf == 0 {
		if evs, err = d.GetEventSerialsByAuthorsTagsCreatedAtRange(f.Tags, f.Authors, since, until, limit); chk.E(err) {
			return
		}
		goto done
	}
	// next kinds/tags
	if kt := hasKinds + hasTags; bf&(kt) == kt && ^kt&bf == 0 {
		if evs, err = d.GetEventSerialsByKindsTagsCreatedAtRange(f.Tags, f.Kinds, since,
			until, limit); chk.E(err) {
			return
		}
		goto done
	}
	// next kinds/authors/tags
	if kat := hasAuthors + hasTags; bf&(kat) == kat && ^kat&bf == 0 {
		if evs, err = d.GetEventSerialsByKindsAuthorsTagsCreatedAtRange(f.Tags, f.Kinds, f.Authors, since, until, limit); chk.E(err) {
			return
		}
		goto done
	}
done:
	// scan the FullIndex for these serials, and sort them by descending created_at
	var index []indexes.FullIndex
	if index, err = d.GetFullIndexesFromSerials(evs); chk.E(err) {
		return
	}
	// sort by reverse chronological order
	sort.Slice(index, func(i, j int) bool {
		return index[i].CreatedAt.ToTimestamp() > index[j].CreatedAt.ToTimestamp()
	})
	for _, item := range index {
		for _, x := range exclude {
			if bytes.Equal(item.Pubkey.Bytes(), x.Bytes()) {
				continue
			}
		}
		evSerials = append(evSerials, item.Ser)
	}
	return
}
