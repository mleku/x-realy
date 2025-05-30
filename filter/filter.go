package filter

import (
	"encoding/json"
	"slices"

	"x.realy.lol/event"
	"x.realy.lol/helpers"
	"x.realy.lol/kind"
	"x.realy.lol/timestamp"
)

type Filters []Filter

type Filter struct {
	IDs     []string
	Kinds   []int
	Authors []string
	Tags    TagMap
	Since   *timestamp.Timestamp
	Until   *timestamp.Timestamp
	Limit   int
	Search  string

	// LimitZero is or must be set when there is a "limit":0 in the filter, and not when "limit" is just omitted
	LimitZero bool `json:"-"`
}

type TagMap map[string][]string

func (eff Filters) String() string {
	j, _ := json.Marshal(eff)
	return string(j)
}

func (eff Filters) Match(event *event.E) bool {
	for _, filter := range eff {
		if filter.Matches(event) {
			return true
		}
	}
	return false
}

func (eff Filters) MatchIgnoringTimestampConstraints(event *event.E) bool {
	for _, filter := range eff {
		if filter.MatchesIgnoringTimestampConstraints(event) {
			return true
		}
	}
	return false
}

func (ef Filter) String() string {
	j, _ := json.Marshal(ef)
	return string(j)
}

func (ef Filter) Matches(event *event.E) bool {
	if !ef.MatchesIgnoringTimestampConstraints(event) {
		return false
	}

	if ef.Since != nil && event.CreatedAt < *ef.Since {
		return false
	}

	if ef.Until != nil && event.CreatedAt > *ef.Until {
		return false
	}

	return true
}

func (ef Filter) MatchesIgnoringTimestampConstraints(event *event.E) bool {
	if event == nil {
		return false
	}

	if ef.IDs != nil && !slices.Contains(ef.IDs, event.Id) {
		return false
	}

	if ef.Kinds != nil && !slices.Contains(ef.Kinds, event.Kind) {
		return false
	}

	if ef.Authors != nil && !slices.Contains(ef.Authors, event.Pubkey) {
		return false
	}

	for f, v := range ef.Tags {
		if v != nil && !event.Tags.ContainsAny(f, v) {
			return false
		}
	}

	return true
}

func FilterEqual(a Filter, b Filter) bool {
	if !helpers.Similar(a.Kinds, b.Kinds) {
		return false
	}

	if !helpers.Similar(a.IDs, b.IDs) {
		return false
	}

	if !helpers.Similar(a.Authors, b.Authors) {
		return false
	}

	if len(a.Tags) != len(b.Tags) {
		return false
	}

	for f, av := range a.Tags {
		if bv, ok := b.Tags[f]; !ok {
			return false
		} else {
			if !helpers.Similar(av, bv) {
				return false
			}
		}
	}

	if !helpers.ArePointerValuesEqual(a.Since, b.Since) {
		return false
	}

	if !helpers.ArePointerValuesEqual(a.Until, b.Until) {
		return false
	}

	if a.Search != b.Search {
		return false
	}

	if a.LimitZero != b.LimitZero {
		return false
	}

	return true
}

func (ef Filter) Clone() Filter {
	clone := Filter{
		IDs:       slices.Clone(ef.IDs),
		Authors:   slices.Clone(ef.Authors),
		Kinds:     slices.Clone(ef.Kinds),
		Limit:     ef.Limit,
		Search:    ef.Search,
		LimitZero: ef.LimitZero,
	}

	if ef.Tags != nil {
		clone.Tags = make(TagMap, len(ef.Tags))
		for k, v := range ef.Tags {
			clone.Tags[k] = slices.Clone(v)
		}
	}

	if ef.Since != nil {
		since := *ef.Since
		clone.Since = &since
	}

	if ef.Until != nil {
		until := *ef.Until
		clone.Until = &until
	}

	return clone
}

// GetTheoreticalLimit gets the maximum number of events that a normal filter would ever return, for example, if
// there is a number of "ids" in the filter, the theoretical limit will be that number of ids.
//
// It returns -1 if there are no theoretical limits.
//
// The given .Limit present in the filter is ignored.
func GetTheoreticalLimit(filter Filter) int {
	if len(filter.IDs) > 0 {
		return len(filter.IDs)
	}

	if len(filter.Kinds) == 0 {
		return -1
	}

	if len(filter.Authors) > 0 {
		allAreReplaceable := true
		for _, k := range filter.Kinds {
			if !kind.IsReplaceableKind(k) {
				allAreReplaceable = false
				break
			}
		}
		if allAreReplaceable {
			return len(filter.Authors) * len(filter.Kinds)
		}

		if len(filter.Tags["d"]) > 0 {
			allAreAddressable := true
			for _, k := range filter.Kinds {
				if !kind.IsAddressableKind(k) {
					allAreAddressable = false
					break
				}
			}
			if allAreAddressable {
				return len(filter.Authors) * len(filter.Kinds) * len(filter.Tags["d"])
			}
		}
	}

	return -1
}
