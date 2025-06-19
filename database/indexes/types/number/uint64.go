package number

import (
	"encoding/binary"
	"io"
)

// Uint64 is a codec for encoding and decoding 64-bit unsigned integers.
type Uint64 struct {
	value uint64
}

// SetUint64 sets the value as a uint64.
func (c *Uint64) SetUint64(value uint64) {
	c.value = value
}

// Uint64 gets the value as a uint64.
func (c *Uint64) Uint64() uint64 {
	return c.value
}

// SetInt sets the value as an int, converting it to uint64.
// Values outside the range of uint64 are truncated.
func (c *Uint64) SetInt(value int) {
	c.value = uint64(value)
}

// Int gets the value as an int, converted from uint64. May truncate if the value exceeds the range of int.
func (c *Uint64) Int() int {
	return int(c.value)
}

// MarshalWrite writes the uint64 value to the provided writer in BigEndian order.
func (c *Uint64) MarshalWrite(w io.Writer) error {
	return binary.Write(w, binary.BigEndian, c.value)
}

// UnmarshalRead reads a uint64 value from the provided reader in BigEndian order.
func (c *Uint64) UnmarshalRead(r io.Reader) error {
	return binary.Read(r, binary.BigEndian, &c.value)
}
