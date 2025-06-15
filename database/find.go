package database

import (
	"bytes"

	"github.com/dgraph-io/badger/v4"

	"x.realy.lol/chk"
	"x.realy.lol/database/indexes"
	"x.realy.lol/database/indexes/prefixes"
	"x.realy.lol/database/indexes/types/idhash"
	"x.realy.lol/database/indexes/types/prefix"
	"x.realy.lol/database/indexes/types/varint"
	"x.realy.lol/errorf"
	"x.realy.lol/event"
	"x.realy.lol/filter"
	"x.realy.lol/log"
	"x.realy.lol/tags"
	"x.realy.lol/timestamp"
)

func (d *D) FindEventSerialById(evId []byte) (ser *varint.V, err error) {
	id := idhash.New()
	if err = id.FromId(evId); chk.E(err) {
		return
	}
	// find by id
	if err = d.View(func(txn *badger.Txn) (err error) {
		key := new(bytes.Buffer)
		if err = indexes.IdSearch(id).MarshalWrite(key); chk.E(err) {
			return
		}
		it := txn.NewIterator(badger.IteratorOptions{Prefix: key.Bytes()})
		defer it.Close()
		for it.Seek(key.Bytes()); it.Valid(); it.Next() {
			item := it.Item()
			k := item.KeyCopy(nil)
			buf := bytes.NewBuffer(k)
			ser = varint.New()
			if err = indexes.IdDec(id, ser).UnmarshalRead(buf); chk.E(err) {
				return
			}
		}
		return
	}); err != nil {
		return
	}
	if ser == nil {
		err = errorf.E("event %0x not found", evId)
		return
	}
	return
}

func (d *D) GetEventFromSerial(ser *varint.V) (ev *event.E, err error) {
	if err = d.View(func(txn *badger.Txn) (err error) {
		enc := indexes.EventDec(ser)
		kb := new(bytes.Buffer)
		if err = enc.MarshalWrite(kb); chk.E(err) {
			return
		}
		var item *badger.Item
		if item, err = txn.Get(kb.Bytes()); chk.E(err) {
			return
		}
		var val []byte
		if val, err = item.ValueCopy(nil); chk.E(err) {
			return
		}
		ev = event.New()
		vr := bytes.NewBuffer(val)
		if err = ev.UnmarshalRead(vr); chk.E(err) {
			return
		}
		return
	}); chk.E(err) {
		return
	}
	return
}

func (d *D) GetEventFullIndexFromSerial(ser *varint.V) (id []byte, err error) {
	if err = d.View(func(txn *badger.Txn) (err error) {
		enc := indexes.New(prefix.New(prefixes.FullIndex), ser)
		prf := new(bytes.Buffer)
		if err = enc.MarshalWrite(prf); chk.E(err) {
			return
		}
		it := txn.NewIterator(badger.IteratorOptions{Prefix: prf.Bytes()})
		defer it.Close()
		for it.Seek(prf.Bytes()); it.Valid(); it.Next() {
			item := it.Item()
			key := item.KeyCopy(nil)
			kbuf := bytes.NewBuffer(key)
			_, t, p, ki, ca := indexes.FullIndexVars()
			dec := indexes.FullIndexDec(ser, t, p, ki, ca)
			if err = dec.UnmarshalRead(kbuf); chk.E(err) {
				return
			}
			id = t.Bytes()
		}
		return
	}); chk.E(err) {
		return
	}
	return
}

func (d *D) GetEventById(evId []byte) (ev *event.E, err error) {
	var ser *varint.V
	if ser, err = d.FindEventSerialById(evId); chk.E(err) {
		return
	}
	ev, err = d.GetEventFromSerial(ser)
	return
}

// GetEventSerialsByCreatedAtRange returns the serials of events with the given since/until
// range in reverse chronological order (starting at until, going back to since).
func (d *D) GetEventSerialsByCreatedAtRange(since, until timestamp.Timestamp) (sers varint.S, err error) {
	// get the start (end) max possible index prefix
	startCreatedAt, _ := indexes.CreatedAtVars()
	startCreatedAt.FromInt(until.ToInt())
	prf := new(bytes.Buffer)
	if err = indexes.CreatedAtEnc(startCreatedAt, nil).MarshalWrite(prf); chk.E(err) {
		return
	}
	if err = d.View(func(txn *badger.Txn) (err error) {
		it := txn.NewIterator(badger.IteratorOptions{Reverse: true, Prefix: prf.Bytes()})
		defer it.Close()
		key := make([]byte, 10)
		for it.Rewind(); it.Valid(); it.Next() {
			item := it.Item()
			key = item.KeyCopy(key)
			ca, ser := indexes.CreatedAtVars()
			buf := bytes.NewBuffer(key)
			if err = indexes.CreatedAtDec(ca, ser).UnmarshalRead(buf); chk.E(err) {
				// skip it then
				continue
			}
			if ca.ToTimestamp() < since {
				break
			}
			sers = append(sers, ser)
		}
		return
	}); chk.E(err) {
		return
	}
	return
}

func (d *D) GetEventSerialsByKindsCreatedAtRange(kinds []int, since, until timestamp.Timestamp) (sers varint.S, err error) {
	// get the start (end) max possible index prefix, one for each kind in the list
	var searchIdxs [][]byte
	kind, startCreatedAt, _ := indexes.KindCreatedAtVars()
	startCreatedAt.FromInt(until.ToInt())
	for _, k := range kinds {
		kind.Set(k)
		prf := new(bytes.Buffer)
		if err = indexes.KindCreatedAtEnc(kind, startCreatedAt, nil).MarshalWrite(prf); chk.E(err) {
			return
		}
		searchIdxs = append(searchIdxs, prf.Bytes())
	}
	for _, idx := range searchIdxs {
		if err = d.View(func(txn *badger.Txn) (err error) {
			it := txn.NewIterator(badger.IteratorOptions{Reverse: true, Prefix: idx})
			defer it.Close()
			key := make([]byte, 10)
			for it.Rewind(); it.Valid(); it.Next() {
				item := it.Item()
				key = item.KeyCopy(key)
				ki, ca, ser := indexes.KindCreatedAtVars()
				buf := bytes.NewBuffer(key)
				if err = indexes.KindCreatedAtDec(ki, ca, ser).UnmarshalRead(buf); chk.E(err) {
					// skip it then
					continue
				}
				if ca.ToTimestamp() < since {
					break
				}
				sers = append(sers, ser)
			}
			return
		}); chk.E(err) {
			return
		}
	}
	return
}

func (d *D) GetEventSerialsByAuthorsCreatedAtRange(pubkeys []string, since, until timestamp.Timestamp) (sers varint.S, err error) {
	// get the start (end) max possible index prefix, one for each kind in the list
	var searchIdxs [][]byte
	var pkDecodeErrs int
	pubkey, startCreatedAt, _ := indexes.PubkeyCreatedAtVars()
	startCreatedAt.FromInt(until.ToInt())
	for _, p := range pubkeys {
		if err = pubkey.FromPubkeyHex(p); chk.E(err) {
			// gracefully ignore wrong keys
			pkDecodeErrs++
			continue
		}
		if pkDecodeErrs == len(pubkeys) {
			err = errorf.E("all pubkeys in authors field of filter failed to decode")
			return
		}
		prf := new(bytes.Buffer)
		if err = indexes.PubkeyCreatedAtEnc(pubkey, startCreatedAt, nil).MarshalWrite(prf); chk.E(err) {
			return
		}
		searchIdxs = append(searchIdxs, prf.Bytes())
	}
	for _, idx := range searchIdxs {
		if err = d.View(func(txn *badger.Txn) (err error) {
			it := txn.NewIterator(badger.IteratorOptions{Reverse: true, Prefix: idx})
			defer it.Close()
			key := make([]byte, 10)
			for it.Rewind(); it.Valid(); it.Next() {
				item := it.Item()
				key = item.KeyCopy(key)
				kind, ca, ser := indexes.KindCreatedAtVars()
				buf := bytes.NewBuffer(key)
				if err = indexes.KindCreatedAtDec(kind, ca, ser).UnmarshalRead(buf); chk.E(err) {
					// skip it then
					continue
				}
				if ca.ToTimestamp() < since {
					break
				}
				sers = append(sers, ser)
			}
			return
		}); chk.E(err) {
			return
		}
	}
	return
}

func (d *D) GetEventSerialsByKindsAuthorsCreatedAtRange(kinds []int, pubkeys []string, since, until timestamp.Timestamp) (sers varint.S, err error) {
	// get the start (end) max possible index prefix, one for each kind in the list
	var searchIdxs [][]byte
	var pkDecodeErrs int
	kind, pubkey, startCreatedAt, _ := indexes.KindPubkeyCreatedAtVars()
	startCreatedAt.FromInt(until.ToInt())
	for _, k := range kinds {
		for _, p := range pubkeys {
			if err = pubkey.FromPubkeyHex(p); chk.E(err) {
				// gracefully ignore wrong keys
				pkDecodeErrs++
				continue
			}
			if pkDecodeErrs == len(pubkeys) {
				err = errorf.E("all pubkeys in authors field of filter failed to decode")
				return
			}
			kind.Set(k)
			prf := new(bytes.Buffer)
			if err = indexes.KindPubkeyCreatedAtEnc(kind, pubkey, startCreatedAt, nil).MarshalWrite(prf); chk.E(err) {
				return
			}
			searchIdxs = append(searchIdxs, prf.Bytes())
		}
	}
	for _, idx := range searchIdxs {
		if err = d.View(func(txn *badger.Txn) (err error) {
			it := txn.NewIterator(badger.IteratorOptions{Reverse: true, Prefix: idx})
			defer it.Close()
			key := make([]byte, 10)
			for it.Rewind(); it.Valid(); it.Next() {
				item := it.Item()
				key = item.KeyCopy(key)
				ki, ca, ser := indexes.KindCreatedAtVars()
				buf := bytes.NewBuffer(key)
				if err = indexes.KindCreatedAtDec(ki, ca, ser).UnmarshalRead(buf); chk.E(err) {
					// skip it then
					continue
				}
				if ca.ToTimestamp() < since {
					break
				}
				sers = append(sers, ser)
			}
			return
		}); chk.E(err) {
			return
		}
	}
	return
}

// GetEventSerialsByTagsCreatedAtRange searches for events that match the tags in a filter and
// returns the list of serials that were found.
func (d *D) GetEventSerialsByTagsCreatedAtRange(t filter.TagMap) (sers varint.S, err error) {
	if len(t) < 1 {
		err = errorf.E("no tags provided")
		return
	}
	var searchIdxs [][]byte
	for tk, tv := range t {
		// the key of each element of the map must be `#X` where X is a-zA-Z
		if len(tk) != 2 {
			continue
		}
		if tk[0] != '#' {
			log.E.F("invalid tag map key '%s'", tk)
		}
		switch tk[1] {
		case 'a':
			// not sure if this is a thing. maybe a prefix search?
			for _, ta := range tv {
				var atag tags.Tag_a
				if atag, err = tags.Decode_a_Tag(ta); chk.E(err) {
					err = nil
					continue
				}
				if atag.Kind == 0 {
					err = nil
					continue
				}
				ki, pk, ident, _ := indexes.TagAVars()
				ki.Set(atag.Kind)
				if atag.Pubkey == nil {
					err = nil
					continue
				}
				if err = pk.FromPubkey(atag.Pubkey); chk.E(err) {
					err = nil
					continue
				}
				if len(atag.Ident) < 1 {
				}
				if err = ident.FromIdent([]byte(atag.Ident)); chk.E(err) {
					err = nil
				}
				buf := new(bytes.Buffer)
				if err = indexes.TagAEnc(ki, pk, ident, nil).MarshalWrite(buf); chk.E(err) {
					err = nil
					continue
				}
				searchIdxs = append(searchIdxs, buf.Bytes())
			}
		case 'd':
			// d tags are identifiers used to mark replaceable events to create a namespace,
			// that the references can be used to replace them, or referred to using 'a' tags.
			for _, td := range tv {
				ident, _ := indexes.TagIdentifierVars()
				if err = ident.FromIdent([]byte(td)); chk.E(err) {
					err = nil
					continue
				}
				buf := new(bytes.Buffer)
				if err = indexes.TagIdentifierEnc(ident, nil).MarshalWrite(buf); chk.E(err) {
					err = nil
					continue
				}
				searchIdxs = append(searchIdxs, buf.Bytes())
			}
		case 'e':
			// e tags refer to events. they can have a third field such as 'root' and 'reply'
			// but this third field isn't indexed.
			for _, te := range tv {
				evt, _ := indexes.TagEventVars()
				if err = evt.FromIdHex(te); chk.E(err) {
					err = nil
					continue
				}
				buf := new(bytes.Buffer)
				if err = indexes.TagEventEnc(evt, nil).MarshalWrite(buf); chk.E(err) {
					err = nil
					continue
				}
				searchIdxs = append(searchIdxs, buf.Bytes())
			}
		case 'p':
			// p tags are references to author pubkeys of events. usually a 64 character hex
			// string but sometimes is a hashtag in follow events.
			for _, te := range tv {
				pk, _ := indexes.TagPubkeyVars()
				if err = pk.FromPubkeyHex(te); chk.E(err) {
					err = nil
					continue
				}
				buf := new(bytes.Buffer)
				if err = indexes.TagPubkeyEnc(pk, nil).MarshalWrite(buf); chk.E(err) {
					err = nil
					continue
				}
				searchIdxs = append(searchIdxs, buf.Bytes())
			}
		case 't':
			// t tags are hashtags, arbitrary strings that can be used to assist search for
			// topics.
			for _, tt := range tv {
				ht, _ := indexes.TagHashtagVars()
				if err = ht.FromIdent([]byte(tt)); chk.E(err) {
					err = nil
					continue
				}
				buf := new(bytes.Buffer)
				if err = indexes.TagHashtagEnc(ht, nil).MarshalWrite(buf); chk.E(err) {
					err = nil
					continue
				}
				searchIdxs = append(searchIdxs, buf.Bytes())
			}
		default:
			// everything else is arbitrary strings, that may have application specific
			// semantics.
			for _, tl := range tv {
				l, val, _ := indexes.TagLetterVars()
				l.Set(tk[1])
				if err = val.FromIdent([]byte(tl)); chk.E(err) {
					err = nil
					continue
				}
				buf := new(bytes.Buffer)
				if err = indexes.TagLetterEnc(l, val, nil).MarshalWrite(buf); chk.E(err) {
					err = nil
					continue
				}
				searchIdxs = append(searchIdxs, buf.Bytes())
			}
		}
	}
	return
}

// GetEventSerialsByAuthorsTagsCreatedAtRange first performs
func (d *D) GetEventSerialsByAuthorsTagsCreatedAtRange(t filter.TagMap, pubkeys []string, since, until timestamp.Timestamp) (sers varint.S, err error) {
	var acSers, tagSers varint.S
	if acSers, err = d.GetEventSerialsByAuthorsCreatedAtRange(pubkeys, since, until); chk.E(err) {
		return
	}
	// now we have the most limited set of serials that are included by the pubkeys, we can then
	// construct the tags searches for all of these serials to filter out the events that don't
	// have both author AND one of the tags.
	if tagSers, err = d.GetEventSerialsByTagsCreatedAtRange(t); chk.E(err) {
		return
	}
	// remove the serials that are not present in both lists.
	sers = varint.Intersect(acSers, tagSers)
	return
}

// GetEventSerialsByKindsTagsCreatedAtRange first performs
func (d *D) GetEventSerialsByKindsTagsCreatedAtRange(t filter.TagMap, kinds []int, since, until timestamp.Timestamp) (sers varint.S, err error) {
	var acSers, tagSers varint.S
	if acSers, err = d.GetEventSerialsByKindsCreatedAtRange(kinds, since, until); chk.E(err) {
		return
	}
	// now we have the most limited set of serials that are included by the pubkeys, we can then
	// construct the tags searches for all of these serials to filter out the events that don't
	// have both author AND one of the tags.
	if tagSers, err = d.GetEventSerialsByTagsCreatedAtRange(t); chk.E(err) {
		return
	}
	// remove the serials that are not present in both lists.
	sers = varint.Intersect(acSers, tagSers)
	return
}

// GetEventSerialsByKindsAuthorsTagsCreatedAtRange first performs
func (d *D) GetEventSerialsByKindsAuthorsTagsCreatedAtRange(t filter.TagMap, kinds []int, pubkeys []string, since, until timestamp.Timestamp) (sers varint.S, err error) {
	var acSers, tagSers varint.S
	if acSers, err = d.GetEventSerialsByKindsAuthorsCreatedAtRange(kinds, pubkeys, since, until); chk.E(err) {
		return
	}
	// now we have the most limited set of serials that are included by the pubkeys, we can then
	// construct the tags searches for all of these serials to filter out the events that don't
	// have both author AND one of the tags.
	if tagSers, err = d.GetEventSerialsByTagsCreatedAtRange(t); chk.E(err) {
		return
	}
	// remove the serials that are not present in both lists.
	sers = varint.Intersect(acSers, tagSers)
	return
}

func (d *D) GetFullIndexesFromSerials(sers varint.S) (index []indexes.FullIndex, err error) {
	for _, ser := range sers {
		if err = d.View(func(txn *badger.Txn) (err error) {
			buf := new(bytes.Buffer)
			if err = indexes.FullIndexEnc(ser, nil, nil, nil, nil).MarshalWrite(buf); chk.E(err) {
				return
			}
			prf := buf.Bytes()
			it := txn.NewIterator(badger.IteratorOptions{Prefix: prf})
			defer it.Close()
			for it.Seek(prf); it.Valid(); {
				item := it.Item()
				key := item.KeyCopy(nil)
				kBuf := bytes.NewBuffer(key)
				s, t, p, k, c := indexes.FullIndexVars()
				if err = indexes.FullIndexDec(s, t, p, k, c).UnmarshalRead(kBuf); chk.E(err) {
					return
				}
				index = append(index, indexes.FullIndex{
					Ser:       s,
					Id:        t,
					Pubkey:    p,
					Kind:      k,
					CreatedAt: c,
				})
				return
			}
			return
		}); chk.E(err) {
			// just skip then.
		}
	}
	return
}
