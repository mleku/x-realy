package number

import (
	"encoding/binary"
	"io"
)

// Uint32Codec is a codec for encoding and decoding 32-bit unsigned integers.
type Uint32Codec struct {
	value uint32
}

// SetUint32 sets the value as a uint32.
func (c *Uint32Codec) SetUint32(value uint32) {
	c.value = value
}

// Uint32 gets the value as a uint32.
func (c *Uint32Codec) Uint32() uint32 {
	return c.value
}

// SetInt sets the value as an int, converting it to uint32.
// Values outside the range of uint32 (0–4294967295) will be truncated.
func (c *Uint32Codec) SetInt(value int) {
	c.value = uint32(value)
}

// Int gets the value as an int, converted from uint32.
func (c *Uint32Codec) Int() int {
	return int(c.value)
}

// MarshalWrite writes the uint32 value to the provided writer in BigEndian order.
func (c *Uint32Codec) MarshalWrite(w io.Writer) error {
	return binary.Write(w, binary.BigEndian, c.value)
}

// UnmarshalRead reads a uint32 value from the provided reader in BigEndian order.
func (c *Uint32Codec) UnmarshalRead(r io.Reader) error {
	return binary.Read(r, binary.BigEndian, &c.value)
}
