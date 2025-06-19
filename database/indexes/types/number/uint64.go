package number

import (
	"encoding/binary"
	"io"
)

// Uint64Codec is a codec for encoding and decoding 64-bit unsigned integers.
type Uint64Codec struct {
	value uint64
}

// SetUint64 sets the value as a uint64.
func (c *Uint64Codec) SetUint64(value uint64) {
	c.value = value
}

// Uint64 gets the value as a uint64.
func (c *Uint64Codec) Uint64() uint64 {
	return c.value
}

// SetInt sets the value as an int, converting it to uint64.
// Values outside the range of uint64 are truncated.
func (c *Uint64Codec) SetInt(value int) {
	c.value = uint64(value)
}

// Int gets the value as an int, converted from uint64. May truncate if the value exceeds the range of int.
func (c *Uint64Codec) Int() int {
	return int(c.value)
}

// MarshalWrite writes the uint64 value to the provided writer in BigEndian order.
func (c *Uint64Codec) MarshalWrite(w io.Writer) error {
	return binary.Write(w, binary.BigEndian, c.value)
}

// UnmarshalRead reads a uint64 value from the provided reader in BigEndian order.
func (c *Uint64Codec) UnmarshalRead(r io.Reader) error {
	return binary.Read(r, binary.BigEndian, &c.value)
}
