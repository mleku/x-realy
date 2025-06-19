package number

import (
	"errors"
	"io"
)

// MaxUint24 is the maximum value of a 24-bit unsigned integer: 2^24 - 1.
const MaxUint24 uint32 = 1<<24 - 1

// Uint24 is a codec for encoding and decoding 24-bit unsigned integers.
type Uint24 struct {
	value uint32
}

// SetUint24 sets the value as a 24-bit unsigned integer.
// If the value exceeds the maximum allowable value for 24 bits, it returns an error.
func (c *Uint24) SetUint24(value uint32) error {
	if value > MaxUint24 {
		return errors.New("value exceeds 24-bit range")
	}
	c.value = value
	return nil
}

// Uint24 gets the value as a 24-bit unsigned integer.
func (c *Uint24) Uint24() uint32 {
	return c.value
}

// SetInt sets the value as an int, converting it to a 24-bit unsigned integer.
// If the value is out of the 24-bit range, it returns an error.
func (c *Uint24) SetInt(value int) error {
	if value < 0 || uint32(value) > MaxUint24 {
		return errors.New("value exceeds 24-bit range")
	}
	c.value = uint32(value)
	return nil
}

// Int gets the value as an int, converted from the 24-bit unsigned integer.
func (c *Uint24) Int() int {
	return int(c.value)
}

// MarshalWrite encodes the 24-bit unsigned integer and writes it directly to the provided io.Writer.
// The encoding uses 3 bytes in BigEndian order.
func (c *Uint24) MarshalWrite(w io.Writer) error {
	if c.value > MaxUint24 {
		return errors.New("value exceeds 24-bit range")
	}

	// Write the 3 bytes (BigEndian order) directly to the writer
	var buf [3]byte
	buf[0] = byte((c.value >> 16) & 0xFF) // Most significant byte
	buf[1] = byte((c.value >> 8) & 0xFF)
	buf[2] = byte(c.value & 0xFF) // Least significant byte

	_, err := w.Write(buf[:]) // Write all 3 bytes to the writer
	return err
}

// UnmarshalRead reads 3 bytes directly from the provided io.Reader and decodes it into a 24-bit unsigned integer.
func (c *Uint24) UnmarshalRead(r io.Reader) error {
	// Read 3 bytes directly from the reader
	var buf [3]byte
	_, err := io.ReadFull(r, buf[:]) // Ensure exactly 3 bytes are read
	if err != nil {
		return err
	}

	// Decode the 3 bytes into a 24-bit unsigned integer
	c.value = (uint32(buf[0]) << 16) |
		(uint32(buf[1]) << 8) |
		uint32(buf[2])

	return nil
}
