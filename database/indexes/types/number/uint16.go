package number

import (
	"encoding/binary"
	"io"
)

// Uint16 is a codec for encoding and decoding 16-bit unsigned integers.
type Uint16 struct {
	value uint16
}

// Set sets the value as a uint16.
func (c *Uint16) Set(value uint16) {
	c.value = value
}

// Get gets the value as a uint16.
func (c *Uint16) Get() uint16 {
	return c.value
}

// SetInt sets the value as an int, converting it to uint16. Truncates values outside uint16 range (0-65535).
func (c *Uint16) SetInt(value int) {
	c.value = uint16(value)
}

// GetInt gets the value as an int, converted from uint16.
func (c *Uint16) GetInt() int {
	return int(c.value)
}

// MarshalWrite writes the uint16 value to the provided writer in BigEndian order.
func (c *Uint16) MarshalWrite(w io.Writer) error {
	return binary.Write(w, binary.BigEndian, c.value)
}

// UnmarshalRead reads a uint16 value from the provided reader in BigEndian order.
func (c *Uint16) UnmarshalRead(r io.Reader) error {
	return binary.Read(r, binary.BigEndian, &c.value)
}
