package timestamp

import (
	"time"

	"golang.org/x/exp/constraints"
)

type Timestamp int64

func Now() Timestamp { return Timestamp(time.Now().Unix()) }

func New[T constraints.Integer | constraints.Float](t T) Timestamp {
	return Timestamp(t)
}

func (t Timestamp) Time() time.Time { return time.Unix(int64(t), 0) }
func (t Timestamp) ToInt64() int64  { return int64(t) }
func (t Timestamp) ToInt() int      { return int(t) }
