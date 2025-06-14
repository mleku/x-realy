package indexes

import (
	"bytes"
	"math"
	"testing"
	"time"

	"github.com/minio/sha256-simd"
	"lukechampine.com/frand"

	"x.realy.lol/chk"
	"x.realy.lol/database/indexes/prefixes"
	"x.realy.lol/database/indexes/types/prefix"
	"x.realy.lol/ec/schnorr"
	"x.realy.lol/log"
)

func TestEvent(t *testing.T) {
	var err error
	for range 100 {
		ser := EventVars()
		ser.FromUint64(uint64(frand.Intn(math.MaxInt64)))
		buf := new(bytes.Buffer)
		evIdx := EventEnc(ser)
		evIdx.MarshalWrite(buf)
		bin := buf.Bytes()
		buf2 := bytes.NewBuffer(bin)
		ser2 := EventVars()
		evIdx2 := EventDec(ser2)
		if err = evIdx2.UnmarshalRead(buf2); chk.E(err) {
			t.Fatal(err)
		}
		if ser.ToUint64() != ser2.ToUint64() {
			t.Fatal("failed to recover same value as input")
		}
	}
}

func TestConfig(t *testing.T) {
	var err error
	cfg := prefix.New(prefixes.Config)
	buf := new(bytes.Buffer)
	cfg.MarshalWrite(buf)
	buf2 := bytes.NewBuffer(cfg.Bytes())
	cfg2 := prefix.New()
	if err = cfg2.UnmarshalRead(buf2); chk.E(err) {
		t.Fatal(err)
	}
	if !bytes.Equal(cfg.Bytes(), cfg2.Bytes()) {
		t.Fatal("failed to recover same value as input")
	}
}

func TestId(t *testing.T) {
	var err error
	for range 100 {
		id, ser := IdVars()
		if err = id.FromId(frand.Bytes(sha256.Size)); chk.E(err) {
			t.Fatal(err)
		}
		ser.FromUint64(uint64(frand.Intn(math.MaxInt64)))
		buf := new(bytes.Buffer)
		evIdx := IdEnc(id, ser)
		evIdx.MarshalWrite(buf)
		bin := buf.Bytes()
		buf2 := bytes.NewBuffer(bin)
		id2, ser2 := IdVars()
		evIdx2 := IdDec(id2, ser2)
		if err = evIdx2.UnmarshalRead(buf2); chk.E(err) {
			t.Fatal(err)
		}
		if !bytes.Equal(id.Bytes(), id2.Bytes()) {
			t.Fatal("failed to recover same value as input")
		}
		if ser.ToUint64() != ser2.ToUint64() {
			t.Fatal("failed to recover same value as input")
		}
	}
}

func TestFullIndex(t *testing.T) {
	var err error
	for range 100 {
		ser, id, p, ki, ca := FullIndexVars()
		if err = id.FromId(frand.Bytes(sha256.Size)); chk.E(err) {
			t.Fatal(err)
		}
		if err = p.FromPubkey(frand.Bytes(schnorr.PubKeyBytesLen)); chk.E(err) {
			t.Fatal(err)
		}
		ki.Set(frand.Intn(math.MaxUint16))
		ca.FromInt(int(time.Now().Unix()))
		ser.FromUint64(uint64(frand.Intn(math.MaxInt64)))
		buf := new(bytes.Buffer)
		fi := FullIndexEnc(ser, id, p, ki, ca)
		if err = fi.MarshalWrite(buf); chk.E(err) {
			t.Fatal(err)
		}
		// log.I.S(fi)
		bin := buf.Bytes()
		// log.I.S(bin)
		buf2 := bytes.NewBuffer(bin)
		ser2, id2, p2, ki2, ca2 := FullIndexVars()
		fi2 := FullIndexDec(ser2, id2, p2, ki2, ca2)
		if err = fi2.UnmarshalRead(buf2); chk.E(err) {
			t.Fatal(err)
		}
		if !bytes.Equal(id.Bytes(), id2.Bytes()) {
			log.I.S(id, id2)
			t.Fatal("failed to recover same value as input")
		}
		if !bytes.Equal(p.Bytes(), p2.Bytes()) {
			t.Fatal("failed to recover same value as input")
		}
		if ki.ToKind() != ki2.ToKind() {
			t.Fatal("failed to recover same value as input")
		}
		if ca.ToTimestamp() != ca2.ToTimestamp() {
			t.Fatal("failed to recover same value as input")
		}
	}
}

func TestPubkey(t *testing.T) {
	var err error
	for range 100 {
		p, ser := PubkeyVars()
		if err = p.FromPubkey(frand.Bytes(schnorr.PubKeyBytesLen)); chk.E(err) {
			t.Fatal(err)
		}
		ser.FromUint64(uint64(frand.Intn(math.MaxInt64)))
		buf := new(bytes.Buffer)
		fi := PubkeyEnc(p, ser)
		fi.MarshalWrite(buf)
		// log.I.S(fi)
		bin := buf.Bytes()
		// log.I.S(bin)
		buf2 := bytes.NewBuffer(bin)
		p2, ser2 := PubkeyVars()
		fi2 := PubkeyDec(p2, ser2)
		if err = fi2.UnmarshalRead(buf2); chk.E(err) {
			t.Fatal(err)
		}

		if ser.ToUint64() != ser2.ToUint64() {
			t.Fatal("failed to recover same value as input")
		}
	}
}

func TestPubkeyCreatedAt(t *testing.T) {
	var err error
	for range 100 {
		p, ca, ser := PubkeyCreatedAtVars()
		if err = p.FromPubkey(frand.Bytes(schnorr.PubKeyBytesLen)); chk.E(err) {
			t.Fatal(err)
		}
		ca.FromInt(int(time.Now().Unix()))
		ser.FromUint64(uint64(frand.Intn(math.MaxInt64)))
		buf := new(bytes.Buffer)
		fi := PubkeyCreatedAtEnc(p, ca, ser)
		fi.MarshalWrite(buf)
		bin := buf.Bytes()
		buf2 := bytes.NewBuffer(bin)
		p2, ca2, ser2 := PubkeyCreatedAtVars()
		fi2 := PubkeyCreatedAtDec(p2, ca2, ser2)
		if err = fi2.UnmarshalRead(buf2); chk.E(err) {
			t.Fatal(err)
		}
		if ca.ToTimestamp() != ca2.ToTimestamp() {
			t.Fatal("failed to recover same value as input")
		}
		if ser.ToUint64() != ser2.ToUint64() {
			t.Fatal("failed to recover same value as input")
		}
	}
}

func TestCreatedAt(t *testing.T) {
	var err error
	for range 100 {
		ca, ser := CreatedAtVars()
		ca.FromInt(int(time.Now().Unix()))
		ser.FromUint64(uint64(frand.Intn(math.MaxInt64)))
		buf := new(bytes.Buffer)
		fi := CreatedAtEnc(ca, ser)
		fi.MarshalWrite(buf)
		bin := buf.Bytes()
		buf2 := bytes.NewBuffer(bin)
		ca2, ser2 := CreatedAtVars()
		fi2 := CreatedAtDec(ca2, ser2)
		if err = fi2.UnmarshalRead(buf2); chk.E(err) {
			t.Fatal(err)
		}
		if ca.ToTimestamp() != ca2.ToTimestamp() {
			t.Fatal("failed to recover same value as input")
		}
		if ser.ToUint64() != ser2.ToUint64() {
			t.Fatal("failed to recover same value as input")
		}
	}
}

func TestFirstSeen(t *testing.T) {
	var err error
	for range 100 {
		ser, ts := FirstSeenVars()
		ts.FromInt(int(time.Now().Unix()))
		ser.FromUint64(uint64(frand.Intn(math.MaxInt64)))
		buf := new(bytes.Buffer)
		fs := FirstSeenEnc(ser, ts)
		fs.MarshalWrite(buf)
		bin := buf.Bytes()
		buf2 := bytes.NewBuffer(bin)
		ser2, ca2 := FirstSeenVars()
		fs2 := FirstSeenDec(ser2, ca2)
		if err = fs2.UnmarshalRead(buf2); chk.E(err) {
			t.Fatal(err)
		}
		if ser.ToUint64() != ser2.ToUint64() {
			t.Fatal("failed to recover same value as input")
		}
		if ts.ToTimestamp() != ca2.ToTimestamp() {
			t.Fatal("failed to recover same value as input")
		}
	}
}

func TestKind(t *testing.T) {
	var err error
	for range 100 {
		ki, ser := KindVars()
		ki.Set(frand.Intn(math.MaxUint16))
		ser.FromUint64(uint64(frand.Intn(math.MaxInt64)))
		buf := new(bytes.Buffer)
		kIdx := KindEnc(ki, ser)
		kIdx.MarshalWrite(buf)
		bin := buf.Bytes()
		buf2 := bytes.NewBuffer(bin)
		ki2, ser2 := KindVars()
		fi2 := KindDec(ki2, ser2)
		if err = fi2.UnmarshalRead(buf2); chk.E(err) {
			t.Fatal(err)
		}
		if ki.ToKind() != ki2.ToKind() {
			t.Fatal("failed to recover same value as input")
		}
		if ser.ToUint64() != ser2.ToUint64() {
			t.Fatal("failed to recover same value as input")
		}
	}
}

func TestTagA(t *testing.T) {
	var err error
	for range 100 {
		ki, p, id, ser := TagAVars()
		if err = id.FromIdent(frand.Bytes(frand.Intn(16) + 8)); chk.E(err) {
			t.Fatal(err)
		}
		if err = p.FromPubkey(frand.Bytes(schnorr.PubKeyBytesLen)); chk.E(err) {
			t.Fatal(err)
		}
		ki.Set(frand.Intn(math.MaxUint16))
		ser.FromUint64(uint64(frand.Intn(math.MaxInt64)))
		buf := new(bytes.Buffer)
		fi := TagAEnc(ki, p, id, ser)
		fi.MarshalWrite(buf)
		bin := buf.Bytes()
		buf2 := bytes.NewBuffer(bin)
		ki2, p2, id2, ser2 := TagAVars()
		fi2 := TagADec(ki2, p2, id2, ser2)
		if err = fi2.UnmarshalRead(buf2); chk.E(err) {
			t.Fatal(err)
		}
		if !bytes.Equal(id.Bytes(), id2.Bytes()) {
			t.Fatal("failed to recover same value as input")
		}
		if !bytes.Equal(p.Bytes(), p2.Bytes()) {
			t.Fatal("failed to recover same value as input")
		}
		if ki.ToKind() != ki2.ToKind() {
			t.Fatal("failed to recover same value as input")
		}
		if ser.ToUint64() != ser2.ToUint64() {
			t.Fatal("failed to recover same value as input")
		}
	}
}

func TestTagEvent(t *testing.T) {
	var err error
	for range 100 {
		id, ser := TagEventVars()
		if err = id.FromId(frand.Bytes(sha256.Size)); chk.E(err) {
			t.Fatal(err)
		}
		ser.FromUint64(uint64(frand.Intn(math.MaxInt64)))
		buf := new(bytes.Buffer)
		evIdx := TagEventEnc(id, ser)
		evIdx.MarshalWrite(buf)
		bin := buf.Bytes()
		buf2 := bytes.NewBuffer(bin)
		id2, ser2 := TagEventVars()
		evIdx2 := TagEventDec(id2, ser2)
		if err = evIdx2.UnmarshalRead(buf2); chk.E(err) {
			t.Fatal(err)
		}
		if !bytes.Equal(id.Bytes(), id2.Bytes()) {
			t.Fatal("failed to recover same value as input")
		}
		if ser.ToUint64() != ser2.ToUint64() {
			t.Fatal("failed to recover same value as input")
		}
	}
}

func TestTagPubkey(t *testing.T) {
	var err error
	for range 100 {
		p, ser := TagPubkeyVars()
		if err = p.FromPubkey(frand.Bytes(schnorr.PubKeyBytesLen)); chk.E(err) {
			t.Fatal(err)
		}
		ser.FromUint64(uint64(frand.Intn(math.MaxInt64)))
		buf := new(bytes.Buffer)
		fi := TagPubkeyEnc(p, ser)
		fi.MarshalWrite(buf)
		// log.I.S(fi)
		bin := buf.Bytes()
		// log.I.S(bin)
		buf2 := bytes.NewBuffer(bin)
		p2, ser2 := TagPubkeyVars()
		fi2 := TagPubkeyDec(p2, ser2)
		if err = fi2.UnmarshalRead(buf2); chk.E(err) {
			t.Fatal(err)
		}
		if ser.ToUint64() != ser2.ToUint64() {
			t.Fatal("failed to recover same value as input")
		}
	}
}

func TestTagHashtag(t *testing.T) {
	var err error
	for range 100 {
		id, ser := TagHashtagVars()
		if err = id.FromIdent(frand.Bytes(frand.Intn(16) + 8)); chk.E(err) {
			t.Fatal(err)
		}
		ser.FromUint64(uint64(frand.Intn(math.MaxInt64)))
		buf := new(bytes.Buffer)
		fi := TagHashtagEnc(id, ser)
		fi.MarshalWrite(buf)
		bin := buf.Bytes()
		buf2 := bytes.NewBuffer(bin)
		id2, ser2 := TagHashtagVars()
		fi2 := TagHashtagDec(id2, ser2)
		if err = fi2.UnmarshalRead(buf2); chk.E(err) {
			t.Fatal(err)
		}
		if !bytes.Equal(id.Bytes(), id2.Bytes()) {
			t.Fatal("failed to recover same value as input")
		}
		if ser.ToUint64() != ser2.ToUint64() {
			t.Fatal("failed to recover same value as input")
		}
	}
}

func TestTagIdentifier(t *testing.T) {
	var err error
	for range 100 {
		id, ser := TagIdentifierVars()
		if err = id.FromIdent(frand.Bytes(frand.Intn(16) + 8)); chk.E(err) {
			t.Fatal(err)
		}
		ser.FromUint64(uint64(frand.Intn(math.MaxInt64)))
		buf := new(bytes.Buffer)
		fi := TagIdentifierEnc(id, ser)
		fi.MarshalWrite(buf)
		bin := buf.Bytes()
		buf2 := bytes.NewBuffer(bin)
		id2, ser2 := TagIdentifierVars()
		fi2 := TagIdentifierDec(id2, ser2)
		if err = fi2.UnmarshalRead(buf2); chk.E(err) {
			t.Fatal(err)
		}
		if !bytes.Equal(id.Bytes(), id2.Bytes()) {
			t.Fatal("failed to recover same value as input")
		}
		if ser.ToUint64() != ser2.ToUint64() {
			t.Fatal("failed to recover same value as input")
		}
	}
}

func TestTagLetter(t *testing.T) {
	var err error
	for range 100 {
		l, id, ser := TagLetterVars()
		if err = id.FromIdent(frand.Bytes(frand.Intn(16) + 8)); chk.E(err) {
			t.Fatal(err)
		}
		lb := frand.Bytes(1)
		l.Set(lb[0])
		ser.FromUint64(uint64(frand.Intn(math.MaxInt64)))
		buf := new(bytes.Buffer)
		fi := TagLetterEnc(l, id, ser)
		fi.MarshalWrite(buf)
		bin := buf.Bytes()
		buf2 := bytes.NewBuffer(bin)
		l2, id2, ser2 := TagLetterVars()
		fi2 := TagLetterDec(l2, id2, ser2)
		if err = fi2.UnmarshalRead(buf2); chk.E(err) {
			t.Fatal(err)
		}
		if l.Letter() != l2.Letter() {
			t.Fatal("failed to recover same value as input")
		}
		if !bytes.Equal(id.Bytes(), id2.Bytes()) {
			t.Fatal("failed to recover same value as input")
		}
		if ser.ToUint64() != ser2.ToUint64() {
			t.Fatal("failed to recover same value as input")
		}
	}
}

func TestTagProtected(t *testing.T) {
	var err error
	for range 100 {
		p, ser := TagProtectedVars()
		if err = p.FromPubkey(frand.Bytes(schnorr.PubKeyBytesLen)); chk.E(err) {
			t.Fatal(err)
		}
		ser.FromUint64(uint64(frand.Intn(math.MaxInt64)))
		buf := new(bytes.Buffer)
		fi := TagProtectedEnc(p, ser)
		fi.MarshalWrite(buf)
		// log.I.S(fi)
		bin := buf.Bytes()
		// log.I.S(bin)
		buf2 := bytes.NewBuffer(bin)
		p2, ser2 := TagProtectedVars()
		fi2 := TagProtectedDec(p2, ser2)
		if err = fi2.UnmarshalRead(buf2); chk.E(err) {
			t.Fatal(err)
		}

		if ser.ToUint64() != ser2.ToUint64() {
			t.Fatal("failed to recover same value as input")
		}
	}
}

func TestTagNonstandard(t *testing.T) {
	var err error
	for range 100 {
		k, v, ser := TagNonstandardVars()
		if err = k.FromIdent(frand.Bytes(frand.Intn(16) + 8)); chk.E(err) {
			t.Fatal(err)
		}
		if err = v.FromIdent(frand.Bytes(frand.Intn(16) + 8)); chk.E(err) {
			t.Fatal(err)
		}
		ser.FromUint64(uint64(frand.Intn(math.MaxInt64)))
		buf := new(bytes.Buffer)
		fi := TagNonstandardEnc(k, v, ser)
		fi.MarshalWrite(buf)
		bin := buf.Bytes()
		buf2 := bytes.NewBuffer(bin)
		k2, v2, ser2 := TagNonstandardVars()
		fi2 := TagNonstandardDec(k2, v2, ser2)
		if err = fi2.UnmarshalRead(buf2); chk.E(err) {
			t.Fatal(err)
		}
		if !bytes.Equal(k.Bytes(), k2.Bytes()) {
			t.Fatal("failed to recover same value as input")
		}
		if !bytes.Equal(v.Bytes(), v2.Bytes()) {
			t.Fatal("failed to recover same value as input")
		}
		if ser.ToUint64() != ser2.ToUint64() {
			t.Fatal("failed to recover same value as input")
		}
	}
}

func TestFulltextWord(t *testing.T) {
	var err error
	for range 100 {
		fw, pos, ser := FullTextWordVars()
		fw.FromWord(frand.Bytes(frand.Intn(10) + 5))
		pos.FromUint64(uint64(frand.Intn(math.MaxUint32)))
		ser.FromUint64(uint64(frand.Intn(math.MaxInt64)))
		buf := new(bytes.Buffer)
		fi := FullTextWordEnc(fw, pos, ser)
		if err = fi.MarshalWrite(buf); chk.E(err) {
			t.Fatal(err)
		}
		bin := buf.Bytes()
		buf2 := bytes.NewBuffer(bin)
		fw2, pos2, ser2 := FullTextWordVars()
		fi2 := FullTextWordDec(fw2, pos2, ser2)
		if err = fi2.UnmarshalRead(buf2); chk.E(err) {
			t.Fatal(err)
		}
		if !bytes.Equal(fw.Bytes(), fw2.Bytes()) {
			t.Fatal("failed to recover same value as input")
		}
		if pos.ToUint32() != pos2.ToUint32() {
			t.Fatal("failed to recover same value as input")
		}
		if ser.ToUint64() != ser2.ToUint64() {
			t.Fatal("failed to recover same value as input")
		}
	}
}

func TestLastAccessed(t *testing.T) {
	var err error
	for range 100 {
		ser := LastAccessedVars()
		ser.FromUint64(uint64(frand.Intn(math.MaxInt64)))
		buf := new(bytes.Buffer)
		fi := LastAccessedEnc(ser)
		fi.MarshalWrite(buf)
		bin := buf.Bytes()
		buf2 := bytes.NewBuffer(bin)
		ser2 := LastAccessedVars()
		fi2 := LastAccessedDec(ser2)
		if err = fi2.UnmarshalRead(buf2); chk.E(err) {
			t.Fatal(err)
		}
		if ser.ToUint64() != ser2.ToUint64() {
			t.Fatal("failed to recover same value as input")
		}
	}
}

func TestAccessCounter(t *testing.T) {
	var err error
	for range 100 {
		ser := AccessCounterVars()
		ser.FromUint64(uint64(frand.Intn(math.MaxInt64)))
		buf := new(bytes.Buffer)
		fi := AccessCounterEnc(ser)
		fi.MarshalWrite(buf)
		bin := buf.Bytes()
		buf2 := bytes.NewBuffer(bin)
		ser2 := AccessCounterVars()
		fi2 := AccessCounterDec(ser2)
		if err = fi2.UnmarshalRead(buf2); chk.E(err) {
			t.Fatal(err)
		}
		if ser.ToUint64() != ser2.ToUint64() {
			t.Fatal("failed to recover same value as input")
		}
	}
}
