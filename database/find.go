package database

import (
	"bytes"

	"github.com/dgraph-io/badger/v4"

	"x.realy.lol/chk"
	"x.realy.lol/database/indexes"
	"x.realy.lol/event"
)

func (d *D) FindEvent(evId []byte) (ev *event.E, err error) {
	id, ser := indexes.IdVars()
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
			if err = indexes.IdDec(id, ser).UnmarshalRead(buf); chk.E(err) {
				return
			}
		}
		return
	}); err != nil {
		return
	}
	if err = d.View(func(txn *badger.Txn) (err error) {
		evk := new(bytes.Buffer)
		if err = indexes.EventEnc(ser).MarshalWrite(evk); chk.E(err) {
			return
		}
		it := txn.NewIterator(badger.IteratorOptions{Prefix: evk.Bytes()})
		defer it.Close()
		for it.Seek(evk.Bytes()); it.Valid(); {
			item := it.Item()
			var val []byte
			if val, err = item.ValueCopy(nil); chk.E(err) {
				return
			}
			ev = event.New()
			if err = ev.UnmarshalRead(bytes.NewBuffer(val)); chk.E(err) {
				return
			}
			return
		}
		return
	}); chk.E(err) {
		return
	}
	return
}
