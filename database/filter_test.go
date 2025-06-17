package database

import (
	"bufio"
	"bytes"
	"testing"

	"x.realy.lol/apputil"
	"x.realy.lol/chk"
	"x.realy.lol/database/indexes/types/varint"
	"x.realy.lol/event"
	"x.realy.lol/filter"
	"x.realy.lol/interrupt"
	"x.realy.lol/log"
)

func TestD_Filter(t *testing.T) {
	var err error
	d := New()
	tmpDir := "testrealy"
	dbExists := !apputil.FileExists(tmpDir)
	if err = d.Init(tmpDir); chk.E(err) {
		t.Fatal(err)
	}
	interrupt.AddHandler(func() {
		d.Close()
	})
	if dbExists {
		buf := bytes.NewBuffer(ExampleEvents)
		scan := bufio.NewScanner(buf)
		scan.Buffer(make([]byte, 5120000), 5120000)
		var count, errs int
		for scan.Scan() {
			b := scan.Bytes()
			ev := event.New()
			if err = ev.Unmarshal(b); chk.E(err) {
				t.Fatalf("%s:\n%s", err, b)
			}
			// verify the signature on the event
			var ok bool
			if ok, err = ev.Verify(); chk.E(err) {
				errs++
				continue
			}
			if !ok {
				errs++
				log.E.F("event signature is invalid\n%s", b)
				continue
			}
			count++
			if count%1000 == 0 {
				log.I.F("unmarshaled %d events", count)
			}
			if err = d.StoreEvent(ev); chk.E(err) {
				continue
			}
		}
		log.I.F("stored %d events", count)
	}
	// fetch some kind 0
	var sers []*varint.V
	if sers, err = d.Filter(filter.F{
		Kinds: []int{0},
		Limit: filter.IntToPointer(50),
	}, nil); chk.E(err) {
		t.Fatal(err)
	}
	// log.I.S(sers)
	var fids [][]byte
	for _, ser := range sers {
		var evIds []byte
		if evIds, err = d.GetEventIdFromSerial(ser); chk.E(err) {
			// continue
			log.I.S(ser)
			t.Fatal(err)
		}
		fids = append(fids, evIds)
	}
	log.I.S(fids)
}
