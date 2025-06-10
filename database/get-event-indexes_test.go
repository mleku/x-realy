package database

import (
	"bufio"
	"bytes"
	_ "embed"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"testing"
	"time"

	"x.realy.lol/apputil"
	"x.realy.lol/chk"
	"x.realy.lol/event"
	"x.realy.lol/log"
	"x.realy.lol/units"
)

var ExampleEvents []byte

func init() {
	var err error
	if !apputil.FileExists("examples.jsonl") {
		var req *http.Request
		req, err = http.NewRequest("GET", "https://files.mleku.dev/examples.jsonl", nil)
		if err != nil {
			panic("wtf")
		}
		var res *http.Response
		if res, err = http.DefaultClient.Do(req); chk.E(err) {
			panic("wtf")
		}
		var fh *os.File
		if fh, err = os.OpenFile("examples.jsonl", os.O_CREATE|os.O_RDWR, 0600); chk.E(err) {
			panic("wtf")
		}
		if _, err = io.Copy(fh, res.Body); chk.E(err) {
			panic("wtf")
		}
		res.Body.Close()
	}
	log.I.F("loading file...")
	var oh *os.File
	if oh, err = os.Open("examples.jsonl"); chk.E(err) {
		panic("wtf")
	}
	buf := new(bytes.Buffer)
	if _, err = io.Copy(buf, oh); chk.E(err) {
		panic("wtf")
	}
	ExampleEvents = buf.Bytes()
	oh.Close()
}

func TestGetEventIndexes(t *testing.T) {
	var err error
	d := New()
	tmpDir := filepath.Join(os.TempDir(), "testrealy")
	if err = d.Init(tmpDir); chk.E(err) {
		t.Fatal(err)
	}
	defer d.Close()
	defer os.RemoveAll(tmpDir)
	buf := bytes.NewBuffer(ExampleEvents)
	scan := bufio.NewScanner(buf)
	scan.Buffer(make([]byte, 5120000), 5120000)
	var count, errs, encErrs, datasize, size, binsize int
	start := time.Now()
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
		// check the event encodes to binary, decodes, and produces the identical canonical form
		binE := new(bytes.Buffer)
		if err = ev.MarshalWrite(binE); chk.E(err) {
			log.I.F("bogus tags probably: %s", b)
			encErrs++
			continue
		}
		ev2 := event.New()
		bin2 := bytes.NewBuffer(binE.Bytes())
		if err = ev2.UnmarshalRead(bin2); chk.E(err) {
			encErrs++
			continue
		}
		var can1, can2 []byte
		ev.ToCanonical(can1)
		ev2.ToCanonical(can2)
		if !bytes.Equal(can1, can2) {
			encErrs++
			log.I.S(can1, can2)
			continue
		}
		binsize += len(binE.Bytes())
		var valid bool
		if valid, err = ev.Verify(); chk.E(err) {
			log.I.F("%s", b)
			encErrs++
			continue
		}
		if !valid {
			t.Fatalf("event failed to verify\n%s", b)
		}
		var indices [][]byte
		if indices, _, err = d.GetEventIndexes(ev); chk.E(err) {
			t.Fatal(err)
		}
		datasize += len(b)
		for _, v := range indices {
			size += len(v)
		}
		_ = indices
		count++
	}
	log.I.F("unmarshaled, verified and indexed %d events in %s, %d Mb of indexes from %d Mb of events, %d Mb as binary, failed verify %d, failed encode %d", count, time.Now().Sub(start), size/units.Mb, datasize/units.Mb, binsize/units.Mb, errs, encErrs)
	d.Close()
	os.RemoveAll(tmpDir)
}

var _ = `wdawdad\nhttps://cdn.discordapp.com/attachments/1277777226397388800/1278018649860472874/grain.png?ex=66cf471e&is=66cdf59e&hm=790aced618bb517ebd560e1fd3def537351ef130e239c0ee86d43ff63c44a146&","sig":"2abc1b3bb119071209daba6bf2b6c76cdad036249aad624938a5a2736739d6c139adb7aa94d24550bc53a972e75c40549513a74d9ace8c4435a5d262c172300b`
var _ = ``
