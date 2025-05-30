package event

import (
	"bufio"
	"bytes"
	"encoding/json"
	"reflect"
	"testing"

	"x.realy.lol/chk"
	"x.realy.lol/errorf"
	"x.realy.lol/event/examples"
	"x.realy.lol/log"
	"x.realy.lol/p256k"
)

func TestE(t *testing.T) {
	scanner := bufio.NewScanner(bytes.NewBuffer(examples.Cache))
	var err error
	for scanner.Scan() {
		b := scanner.Bytes()
		ev1 := &E{}
		if err = json.Unmarshal(b, ev1); chk.E(err) {
			t.Fatal(err)
		}
		if !ev1.CheckId() {
			t.Fatalf("failed to verify event:\ngot: %s\nexpect: %s", ev1.IdHex(), ev1.Id)
		}
		can := ev1.ToCanonical(nil)
		ev2 := &E{}
		if err = ev2.FromCanonical(can); chk.E(err) {
			return
		}
		ev2.Sig = ev1.Sig
		ev2.Id = ev1.IdHex()
		if !bytes.Equal(ev1.Serialize(), ev2.Serialize()) {
			log.I.F("%v\n%s\n%s", ev1.Serialize(), ev2.Serialize())
			t.Fatal("failed to unmarshal via ToCanonical/FromCanonical")
		}
		var b2 []byte
		if b2, err = json.Marshal(ev1); chk.E(err) {
			t.Fatal(err)
		}
		ev3 := &E{}
		if err = json.Unmarshal(b2, ev3); chk.E(err) {
			t.Fatal(err)
		}
		if !reflect.DeepEqual(ev1, ev3) {
			t.Fatal("failed to unmarshal via json.Marshal/json.Unmarshal")
		}
	}
}

func TestSignVerify(t *testing.T) {
	var err error
	signer := new(p256k.Signer)
	if err = signer.Generate(); chk.E(err) {
		t.Fatal(err)
	}
	scanner := bufio.NewScanner(bytes.NewBuffer(examples.Cache))
	for scanner.Scan() {
		b := scanner.Bytes()
		ev1 := &E{}
		if err = json.Unmarshal(b, ev1); chk.E(err) {
			t.Fatal(err)
		}
		var ok bool
		if ok, err = ev1.Verify(); !ok {
			err = errorf.E("failed to verify original signature on example event\n%s", b)
			t.Fatal(err)
		}
		if err = ev1.Sign(signer); chk.E(err) {
			t.Fatal(err)
		}
		if ok, err = ev1.Verify(); chk.E(err) {
			t.Fatal(err)
		}
		if !ok {
			t.Fatalf("failed to sign event\n%s", ev1.Serialize())
		}
	}
}
