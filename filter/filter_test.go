package filter

import (
	"encoding/json"
	"slices"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"x.realy.lol/chk"
	"x.realy.lol/event"
	"x.realy.lol/kind"
	"x.realy.lol/log"
	"x.realy.lol/timestamp"
)

func TestFilterUnmarshal(t *testing.T) {
	raw := `{"ids": ["abc"],"#e":["zzz"],"#something":["nothing","bab"],"since":1644254609,"search":"test"}`
	var f F
	err := json.Unmarshal([]byte(raw), &f)
	assert.NoError(t, err)

	assert.Condition(t, func() (success bool) {
		if f.Since == nil || f.Since.Time().UTC().Format("2006-01-02") != "2022-02-07" ||
			f.Until != nil ||
			f.Tags == nil || len(f.Tags) != 2 || !slices.Contains(f.Tags["something"], "bab") ||
			f.Search != "test" {
			return false
		}
		return true
	}, "failed to parse filter correctly")
}

func TestFilterMarshal(t *testing.T) {
	until := timestamp.Timestamp(12345678)
	filterj, err := json.Marshal(F{
		Kinds: []int{kind.TextNote, kind.RecommendServer, kind.EncryptedDirectMessage},
		Tags:  TagMap{"fruit": {"banana", "mango"}},
		Until: &until,
	})
	assert.NoError(t, err)

	expected := `{"kinds":[1,2,4],"until":12345678,"#fruit":["banana","mango"]}`
	assert.Equal(t, expected, string(filterj))
}

func TestFilterUnmarshalWithLimitZero(t *testing.T) {
	raw := `{"ids": ["abc"],"#e":["zzz"],"limit":0,"#something":["nothing","bab"],"since":1644254609,"search":"test"}`
	var f F
	err := json.Unmarshal([]byte(raw), &f)
	assert.NoError(t, err)

	assert.Condition(t, func() (success bool) {
		if f.Since == nil ||
			f.Since.Time().UTC().Format("2006-01-02") != "2022-02-07" ||
			f.Until != nil ||
			f.Tags == nil || len(f.Tags) != 2 || !slices.Contains(f.Tags["something"], "bab") ||
			f.Search != "test" {
			return false
		}
		return true
	}, "failed to parse filter correctly")
}

func TestFilterMarshalWithLimitZero(t *testing.T) {
	until := timestamp.Timestamp(12345678)
	filterj, err := json.Marshal(F{
		Kinds: []int{kind.TextNote, kind.RecommendServer, kind.EncryptedDirectMessage},
		Tags:  TagMap{"fruit": {"banana", "mango"}},
		Until: &until,
	})
	assert.NoError(t, err)

	expected := `{"kinds":[1,2,4],"until":12345678,"limit":0,"#fruit":["banana","mango"]}`
	assert.Equal(t, expected, string(filterj))
}

func TestFilterMatchingLive(t *testing.T) {
	var filter F
	var event event.E

	json.Unmarshal([]byte(`{"kinds":[1],"authors":["a8171781fd9e90ede3ea44ddca5d3abf828fe8eedeb0f3abb0dd3e563562e1fc","1d80e5588de010d137a67c42b03717595f5f510e73e42cfc48f31bae91844d59","ed4ca520e9929dfe9efdadf4011b53d30afd0678a09aa026927e60e7a45d9244"],"since":1677033299}`), &filter)
	json.Unmarshal([]byte(`{"id":"5a127c9c931f392f6afc7fdb74e8be01c34035314735a6b97d2cf360d13cfb94","pubkey":"1d80e5588de010d137a67c42b03717595f5f510e73e42cfc48f31bae91844d59","created_at":1677033299,"kind":1,"tags":[["t","japan"]],"content":"If you like my art,I'd appreciate a coin or two!!\nZap is welcome!! Thanks.\n\n\n#japan #bitcoin #art #bananaart\nhttps://void.cat/d/CgM1bzDgHUCtiNNwfX9ajY.webp","sig":"828497508487ca1e374f6b4f2bba7487bc09fccd5cc0d1baa82846a944f8c5766918abf5878a580f1e6615de91f5b57a32e34c42ee2747c983aaf47dbf2a0255"}`), &event)

	assert.True(t, filter.Matches(&event), "live filter should match")
}

func TestFilterEquality(t *testing.T) {
	assert.True(t, FilterEqual(
		F{Kinds: []int{kind.EncryptedDirectMessage, kind.Deletion}},
		F{Kinds: []int{kind.EncryptedDirectMessage, kind.Deletion}},
	), "kinds filters should be equal")

	assert.True(t, FilterEqual(
		F{Kinds: []int{kind.EncryptedDirectMessage, kind.Deletion}, Tags: TagMap{"letter": {"a", "b"}}},
		F{Kinds: []int{kind.EncryptedDirectMessage, kind.Deletion}, Tags: TagMap{"letter": {"b", "a"}}},
	), "kind+tags filters should be equal")

	tm := timestamp.Now()
	assert.True(t, FilterEqual(
		F{
			Kinds: []int{kind.EncryptedDirectMessage, kind.Deletion},
			Tags:  TagMap{"letter": {"a", "b"}, "fruit": {"banana"}},
			Since: &tm,
			Ids:   []string{"aaaa", "bbbb"},
		},
		F{
			Kinds: []int{kind.Deletion, kind.EncryptedDirectMessage},
			Tags:  TagMap{"letter": {"a", "b"}, "fruit": {"banana"}},
			Since: &tm,
			Ids:   []string{"aaaa", "bbbb"},
		},
	), "kind+2tags+since+ids filters should be equal")

	assert.False(t, FilterEqual(
		F{Kinds: []int{kind.TextNote, kind.EncryptedDirectMessage, kind.Deletion}},
		F{Kinds: []int{kind.EncryptedDirectMessage, kind.Deletion, kind.Repost}},
	), "kinds filters shouldn't be equal")
}

func TestFilterClone(t *testing.T) {
	ts := timestamp.Now() - 60*60
	flt := F{
		Kinds: []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9},
		Tags:  TagMap{"letter": {"a", "b"}, "fruit": {"banana"}},
		Since: &ts,
		Ids:   []string{"9894b4b5cb5166d23ee8899a4151cf0c66aec00bde101982a13b8e8ceb972df9"},
	}
	clone := flt.Clone()
	assert.True(t, FilterEqual(flt, clone), "clone is not equal:\n %v !=\n %v", flt, clone)

	clone1 := flt.Clone()
	clone1.Ids = append(clone1.Ids, "88f0c63fcb93463407af97a5e5ee64fa883d107ef9e558472c4eb9aaaefa459d")
	assert.False(t, FilterEqual(flt, clone1), "modifying the clone ids should cause it to not be equal anymore")

	clone2 := flt.Clone()
	clone2.Tags["letter"] = append(clone2.Tags["letter"], "c")
	assert.False(t, FilterEqual(flt, clone2), "modifying the clone tag items should cause it to not be equal anymore")

	clone3 := flt.Clone()
	clone3.Tags["g"] = []string{"drt"}
	assert.False(t, FilterEqual(flt, clone3), "modifying the clone tag map should cause it to not be equal anymore")

	clone4 := flt.Clone()
	*clone4.Since++
	assert.False(t, FilterEqual(flt, clone4), "modifying the clone since should cause it to not be equal anymore")
}

func TestTheoreticalLimit(t *testing.T) {
	require.Equal(t, 6, GetTheoreticalLimit(F{Ids: []string{"a", "b", "c", "d", "e", "f"}}))
	require.Equal(t, 9, GetTheoreticalLimit(F{Authors: []string{"a", "b", "c"}, Kinds: []int{3, 0, 10002}}))
	require.Equal(t, 4, GetTheoreticalLimit(F{Authors: []string{"a", "b", "c", "d"}, Kinds: []int{10050}}))
	require.Equal(t, -1, GetTheoreticalLimit(F{Authors: []string{"a", "b", "c", "d"}}))
	require.Equal(t, -1, GetTheoreticalLimit(F{Kinds: []int{3, 0, 10002}}))
	require.Equal(t, 24, GetTheoreticalLimit(F{Authors: []string{"a", "b", "c", "d", "e", "f"}, Kinds: []int{30023, 30024}, Tags: TagMap{"d": []string{"aaa", "bbb"}}}))
	require.Equal(t, -1, GetTheoreticalLimit(F{Authors: []string{"a", "b", "c", "d", "e", "f"}, Kinds: []int{30023, 30024}}))
}

func TestFilter(t *testing.T) {
	ts := timestamp.Now() - 60*60
	now := timestamp.Now()
	flt := &F{
		Authors: []string{"1d80e5588de010d137a67c42b03717595f5f510e73e42cfc48f31bae91844d59"},
		Kinds:   []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9},
		Tags: TagMap{
			"#t": {"a", "b"},
			"#e": {"9894b4b5cb5166d23ee8899a4151cf0c66aec00bde101982a13b8e8ceb972df9"},
			"#p": {"1d80e5588de010d137a67c42b03717595f5f510e73e42cfc48f31bae91844d59"},
		},
		Until: &now,
		Since: &ts,
		Ids:   []string{"9894b4b5cb5166d23ee8899a4151cf0c66aec00bde101982a13b8e8ceb972df9"},
		// Limit: IntToPointer(10),
	}
	var err error
	var b []byte
	if b, err = json.Marshal(flt); chk.E(err) {
		t.Fatal(err)
	}
	log.I.F("%s", b)
	var f2 F
	if err = json.Unmarshal(b, &f2); chk.E(err) {
		t.Fatal(err)
	}
	log.I.S(f2)
}
