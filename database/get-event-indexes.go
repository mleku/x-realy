package database

import (
	"bytes"
	"time"

	"x.realy.lol/chk"
	"x.realy.lol/database/indexes"
	"x.realy.lol/database/indexes/types/fullid"
	identhash "x.realy.lol/database/indexes/types/identHash"
	"x.realy.lol/database/indexes/types/idhash"
	"x.realy.lol/database/indexes/types/kindidx"
	"x.realy.lol/database/indexes/types/letter"
	"x.realy.lol/database/indexes/types/pubhash"
	"x.realy.lol/database/indexes/types/serial"
	"x.realy.lol/database/indexes/types/timestamp"
	"x.realy.lol/event"
	"x.realy.lol/hex"
	"x.realy.lol/tags"
)

// GetEventIndexes generates a set of indexes for a new event record. The first record is the
// key that should have the binary encoded event as its value.
func (d *D) GetEventIndexes(ev *event.E) (indices [][]byte, ser *serial.S, err error) {
	// log.I.F("getting event indices for\n%s", ev.Serialize())
	// get a new serial
	ser = serial.New()
	var s uint64
	if s, err = d.Serial(); chk.E(err) {
		return
	}
	ser.FromSerial(s)
	// create the event id key
	id := idhash.New()
	var idb []byte
	if idb, err = ev.IdBytes(); chk.E(err) {
		return
	}
	if err = id.FromId(idb); chk.E(err) {
		return
	}
	evIDB := new(bytes.Buffer)
	if err = indexes.IdEnc(id, ser).MarshalWrite(evIDB); chk.E(err) {
		return
	}
	indices = append(indices, evIDB.Bytes())
	// create the full index key
	fid := fullid.New()
	if err = fid.FromId(idb); chk.E(err) {
		return
	}
	p := pubhash.New()
	var pk []byte
	if pk, err = ev.PubBytes(); chk.E(err) {
		return
	}
	if err = p.FromPubkey(pk); chk.E(err) {
		return
	}
	ki := kindidx.FromKind(ev.Kind)
	ca := &timestamp.T{}
	ca.FromInt64(int64(ev.CreatedAt))
	evIFiB := new(bytes.Buffer)
	if err = indexes.FullIndexEnc(fid, p, ki, ca, ser).MarshalWrite(evIFiB); chk.E(err) {
		return
	}
	indices = append(indices, evIFiB.Bytes())
	// pubkey index
	evIPkB := new(bytes.Buffer)
	if err = indexes.PubkeyEnc(p, ser).MarshalWrite(evIPkB); chk.E(err) {
		return
	}
	indices = append(indices, evIPkB.Bytes())
	// pubkey-created_at index
	evIPkCaB := new(bytes.Buffer)
	if err = indexes.PubkeyCreatedAtEnc(p, ca, ser).MarshalWrite(evIPkCaB); chk.E(err) {
		return
	}
	indices = append(indices, evIPkCaB.Bytes())
	// created_at index
	evICaB := new(bytes.Buffer)
	if err = indexes.CreatedAtEnc(ca, ser).MarshalWrite(evICaB); chk.E(err) {
		return
	}
	indices = append(indices, evICaB.Bytes())
	// FirstSeen index
	evIFsB := new(bytes.Buffer)
	fs := &timestamp.T{}
	fs.FromInt64(time.Now().Unix())
	if err = indexes.FirstSeenEnc(ser, fs).MarshalWrite(evIFsB); chk.E(err) {
		return
	}
	indices = append(indices, evIFsB.Bytes())
	// Kind index
	evIKiB := new(bytes.Buffer)
	if err = indexes.KindEnc(ki, ser).MarshalWrite(evIKiB); chk.E(err) {
		return
	}
	indices = append(indices, evIKiB.Bytes())
	// tags
	// TagA index
	var atags []tags.Tag_a
	var tagAs []indexes.TagA
	atags = ev.Tags.Get_a_Tags()
	for _, v := range atags {
		aki, apk, aid, _ := indexes.TagAVars()
		aki.Set(v.Kind)
		if err = apk.FromPubkey(v.Pubkey); chk.E(err) {
			continue
		}
		if err = aid.FromIdent([]byte(v.Ident)); chk.E(err) {
			continue
		}
		tagAs = append(tagAs, indexes.TagA{
			Ki: aki, P: apk, Id: aid, Ser: ser,
		})
	}
	for _, v := range tagAs {
		evITaB := new(bytes.Buffer)
		if err = indexes.TagAEnc(v.Ki, v.P, v.Id, ser).MarshalWrite(evITaB); chk.E(err) {
			return
		}
		indices = append(indices, evITaB.Bytes())
	}
	// TagEvent index
	eTags := ev.Tags.GetAllExactKeys("e")
	for _, v := range eTags {
		eid := v.Value()
		var eh []byte
		if eh, err = hex.Dec(eid); chk.E(err) {
			err = nil
			continue
		}
		ih := idhash.New()
		if err = ih.FromId(eh); chk.E(err) {
			err = nil
			continue
		}
		evIeB := new(bytes.Buffer)
		if err = indexes.TagEventEnc(ih, ser).MarshalWrite(evIeB); chk.E(err) {
			return
		}
		indices = append(indices, evIeB.Bytes())
	}
	// TagPubkey index
	pTags := ev.Tags.GetAllExactKeys("p")
	for _, v := range pTags {
		pt := v.Value()
		var pkb []byte
		if pkb, err = hex.Dec(pt); err != nil {
			err = nil
			continue
		}
		ph := pubhash.New()
		if err = ph.FromPubkey(pkb); chk.E(err) {
			err = nil
			continue
		}
		evIpB := new(bytes.Buffer)
		if err = indexes.TagPubkeyEnc(ph, ser).MarshalWrite(evIpB); chk.E(err) {
			return
		}
		indices = append(indices, evIpB.Bytes())
	}
	// TagHashtag index
	ttags := ev.Tags.GetAllExactKeys("t")
	for _, v := range ttags {
		ht := v.Value()
		hh := identhash.New()
		if err = hh.FromIdent([]byte(ht)); chk.E(err) {
			err = nil
			continue
		}
		evIhB := new(bytes.Buffer)
		if err = indexes.TagHashtagEnc(hh, ser).MarshalWrite(evIhB); chk.E(err) {
			return
		}
		indices = append(indices, evIhB.Bytes())
	}
	// TagIdentifier index
	dtags := ev.Tags.GetAllExactKeys("d")
	for _, v := range dtags {
		dt := v.Value()
		dh := identhash.New()
		if err = dh.FromIdent([]byte(dt)); chk.E(err) {
			err = nil
			continue
		}
		evIidB := new(bytes.Buffer)
		if err = indexes.TagIdentifierEnc(dh, ser).MarshalWrite(evIidB); chk.E(err) {
			return
		}
		indices = append(indices, evIidB.Bytes())
	}
	// TagLetter index, TagProtected, TagNonstandard
	for _, v := range ev.Tags {
		key := v.Key()
		if len(key) == 1 {
			switch key {
			case "t", "p", "e":
				// we already made indexes for these letters
				continue
			case "-":
				// TagProtected
				evIprotB := new(bytes.Buffer)
				if err = indexes.TagProtectedEnc(p, ser).MarshalWrite(evIprotB); chk.E(err) {
					return
				}
				indices = append(indices, evIprotB.Bytes())
			default:
				if !((key[0] >= 'a' && key[0] <= 'z') || (key[0] >= 'A' && key[0] <= 'Z')) {
					// this is not a single letter tag or protected. nonstandard
					nk, nv := identhash.New(), identhash.New()
					_ = nk.FromIdent([]byte(key))
					if len(v) > 1 {
						_ = nv.FromIdent([]byte(v.Value()))
					} else {
						_ = nv.FromIdent([]byte{})
					}
					evInsB := new(bytes.Buffer)
					if err = indexes.TagNonstandardEnc(nk, nv, ser).MarshalWrite(evInsB); chk.E(err) {
						return
					}
					indices = append(indices, evInsB.Bytes())
					continue
				}
			}
			// we have a single letter that is not e, p or t
			l := letter.New(key[0])
			val := identhash.New()
			// this can be empty, but the hash would still be distinct
			if err = val.FromIdent([]byte(v.Value())); chk.E(err) {
				continue
			}
			evIlB := new(bytes.Buffer)
			if err = indexes.TagLetterEnc(l, val, ser).MarshalWrite(evIlB); chk.E(err) {
				return
			}
			indices = append(indices, evIlB.Bytes())
		} else {
			// TagNonstandard
			nk, nv := identhash.New(), identhash.New()
			_ = nk.FromIdent([]byte(key))
			if len(v) > 1 {
				_ = nv.FromIdent([]byte(v.Value()))
			} else {
				_ = nv.FromIdent([]byte{})
			}
			evInsB := new(bytes.Buffer)
			if err = indexes.TagNonstandardEnc(nk, nv, ser).MarshalWrite(evInsB); chk.E(err) {
				return
			}
			indices = append(indices, evInsB.Bytes())
		}
	}
	// FullTextWord index
	var ftk [][]byte
	if ftk, err = d.GetFulltextKeys(ev, ser); chk.E(err) {
		return
	}
	indices = append(indices, ftk...)
	return
}
