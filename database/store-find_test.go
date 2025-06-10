package database

import (
	"bufio"
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"x.realy.lol/chk"
	"x.realy.lol/event"
	"x.realy.lol/interrupt"
	"x.realy.lol/log"
)

func TestD_StoreEvent(t *testing.T) {
	var err error
	d := New()
	tmpDir := filepath.Join(os.TempDir(), "testrealy")
	os.RemoveAll(tmpDir)
	if err = d.Init(tmpDir); chk.E(err) {
		t.Fatal(err)
	}
	buf := bytes.NewBuffer(ExampleEvents)
	scan := bufio.NewScanner(buf)
	scan.Buffer(make([]byte, 5120000), 5120000)
	var count, errs int
	var evIds [][]byte
	interrupt.AddHandler(func() {
		d.Close()
		os.RemoveAll(tmpDir)
	})
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
		if count%10000 == 0 {
			log.I.F("unmarshaled %d events", count)
			// break
		}
		if err = d.StoreEvent(ev); chk.E(err) {
			continue
		}
		evIds = append(evIds, ev.GetIdBytes())
	}
	log.I.F("completed unmarshalling %d events", count)
	for _, v := range evIds {
		var ev *event.E
		if ev, err = d.FindEvent(v); chk.E(err) {
			t.Fatal(err)
		}
		_ = ev
		// log.I.S(ev)
	}
	log.I.F("stored and retrieved %d events", len(evIds))
	return
}
