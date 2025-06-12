package database

import (
	"bytes"
	"math"

	"github.com/dgraph-io/badger/v4"

	"x.realy.lol/chk"
	"x.realy.lol/database/indexes"
	"x.realy.lol/database/indexes/prefixes"
	"x.realy.lol/database/indexes/types/idhash"
	"x.realy.lol/database/indexes/types/prefix"
	"x.realy.lol/database/indexes/types/varint"
	"x.realy.lol/errorf"
	"x.realy.lol/event"
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
func (d *D) GetEventSerialsByCreatedAtRange(since, until timestamp.Timestamp) (sers []*varint.V, err error) {
	// get the start (end) max possible index prefix
	startCreatedAt, startSer := indexes.CreatedAtVars()
	startCreatedAt.FromInt64(until.ToInt64())
	startSer.FromUint64(math.MaxUint64)
	prf := new(bytes.Buffer)
	if err = indexes.CreatedAtEnc(startCreatedAt, startSer).MarshalWrite(prf); chk.E(err) {
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

func (d *D) GetEventSerialsByKindsCreatedAtRange(kinds []int, since, until timestamp.Timestamp) (sers []*varint.V, err error) {
	// get the start (end) max possible index prefix, one for each kind in the list
	var searchIdxs [][]byte
	for _, k := range kinds {
		kind, startCreatedAt, startSer := indexes.KindCreatedAtVars()
		kind.Set(k)
		startCreatedAt.FromInt64(until.ToInt64())
		startSer.FromUint64(math.MaxUint64)
		prf := new(bytes.Buffer)
		if err = indexes.KindCreatedAtEnc(kind, startCreatedAt, startSer).MarshalWrite(prf); chk.E(err) {
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

func (d *D) GetEventSerialsByAuthorsCreatedAtRange(pubkeys []string, since, until timestamp.Timestamp) (sers []*varint.V, err error) {
	// get the start (end) max possible index prefix, one for each kind in the list
	var searchIdxs [][]byte
	var pkDecodeErrs int
	for _, p := range pubkeys {
		pubkey, startCreatedAt, startSer := indexes.PubkeyCreatedAtVars()
		if err = pubkey.FromPubkeyHex(p); chk.E(err) {
			// gracefully ignore wrong keys
			pkDecodeErrs++
			continue
		}
		if pkDecodeErrs == len(pubkeys) {
			err = errorf.E("all pubkeys in authors field of filter failed to decode")
			return
		}
		startCreatedAt.FromInt64(until.ToInt64())
		startSer.FromUint64(math.MaxUint64)
		prf := new(bytes.Buffer)
		if err = indexes.PubkeyCreatedAtEnc(pubkey, startCreatedAt, startSer).MarshalWrite(prf); chk.E(err) {
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

func (d *D) GetEventSerialsByKindsAuthorsCreatedAtRange(kinds []int, pubkeys []string, since, until timestamp.Timestamp) (sers []*varint.V, err error) {
	// get the start (end) max possible index prefix, one for each kind in the list
	var searchIdxs [][]byte
	var pkDecodeErrs int
	for _, k := range kinds {
		for _, p := range pubkeys {
			kind, pubkey, startCreatedAt, startSer := indexes.KindPubkeyCreatedAtVars()
			if err = pubkey.FromPubkeyHex(p); chk.E(err) {
				// gracefully ignore wrong keys
				pkDecodeErrs++
				continue
			}
			if pkDecodeErrs == len(pubkeys) {
				err = errorf.E("all pubkeys in authors field of filter failed to decode")
				return
			}
			startCreatedAt.FromInt64(until.ToInt64())
			startSer.FromUint64(math.MaxUint64)
			kind.Set(k)
			prf := new(bytes.Buffer)
			if err = indexes.KindPubkeyCreatedAtEnc(kind, pubkey, startCreatedAt, startSer).MarshalWrite(prf); chk.E(err) {
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
