package database

import (
	"bytes"
	"time"

	"x.realy.lol/chk"
	"x.realy.lol/database/indexes"
	"x.realy.lol/database/indexes/types/timestamp"
	"x.realy.lol/database/indexes/types/varint"
	"x.realy.lol/errorf"
	"x.realy.lol/event"
)

func (d *D) StoreEvent(ev *event.E) (err error) {
	var ev2 *event.E
	if ev2, err = d.FindEventById(ev.GetIdBytes()); err != nil {
		// so we didn't find it?
	}
	if ev2 != nil {
		// we did found it
		if ev.Id == ev2.Id {
			err = errorf.E("duplicate event")
			return
		}
	}
	var ser *varint.V
	var idxs [][]byte
	if idxs, ser, err = d.GetEventIndexes(ev); chk.E(err) {
		return
	}
	_ = idxs
	evK := new(bytes.Buffer)
	if err = indexes.EventEnc(ser).MarshalWrite(evK); chk.E(err) {
		return
	}
	ts := &timestamp.T{}
	ts.FromInt64(time.Now().Unix())
	// FirstSeen
	fsI := new(bytes.Buffer)
	if err = indexes.FirstSeenEnc(ser, ts).MarshalWrite(fsI); chk.E(err) {
		return
	}
	idxs = append(idxs, fsI.Bytes())
	// write indexes; none of the above have values.
	for _, v := range idxs {
		if err = d.Set(v, nil); chk.E(err) {
			return
		}
	}
	// LastAccessed
	laI := new(bytes.Buffer)
	if err = indexes.LastAccessedEnc(ser).MarshalWrite(laI); chk.E(err) {
		return
	}
	if err = d.Set(laI.Bytes(), ts.Bytes()); chk.E(err) {
		return
	}
	// AccessCounter
	acI := new(bytes.Buffer)
	if err = indexes.AccessCounterEnc(ser).MarshalWrite(acI); chk.E(err) {
		return
	}
	ac := varint.New()
	if err = d.Set(acI.Bytes(), ac.Bytes()); chk.E(err) {
		return
	}
	// lastly, the event
	evk := new(bytes.Buffer)
	if err = indexes.EventEnc(ser).MarshalWrite(evk); chk.E(err) {
		return
	}
	evV := new(bytes.Buffer)
	if err = ev.MarshalWrite(evV); chk.E(err) {
		return
	}
	if err = d.Set(evk.Bytes(), evV.Bytes()); chk.E(err) {
		return
	}
	return
}
