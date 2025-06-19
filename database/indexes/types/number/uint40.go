package number

import (
	"errors"
	"io"
)

// MaxUint40 is the maximum value of a 40-bit unsigned integer: 2^40 - 1.
const MaxUint40 uint64 = 1<<40 - 1

// Uint40 is a codec for encoding and decoding 40-bit unsigned integers.
type Uint40 struct{ value uint64 }

// SetUint40 sets the value as a 40-bit unsigned integer.
// If the value exceeds the maximum allowable value for 40 bits, it returns an error.
func (c *Uint40) SetUint40(value uint64) error {
	if value > MaxUint40 {
		return errors.New("value exceeds 40-bit range")
	}
	c.value = value
	return nil
}

// Uint40 gets the value as a 40-bit unsigned integer.
func (c *Uint40) Uint40() uint64 { return c.value }

// SetInt sets the value as an int, converting it to a 40-bit unsigned integer.
// If the value is out of the 40-bit range, it returns an error.
func (c *Uint40) SetInt(value int) error {
	if value < 0 || uint64(value) > MaxUint40 {
		return errors.New("value exceeds 40-bit range")
	}
	c.value = uint64(value)
	return nil
}

// Int gets the value as an int, converted from the 40-bit unsigned integer.
// Note: If the value exceeds the int range, it will be truncated.
func (c *Uint40) Int() int { return int(c.value) }

// MarshalWrite encodes the 40-bit unsigned integer and writes it to the provided writer.
// The encoding uses 5 bytes in BigEndian order.
func (c *Uint40) MarshalWrite(w io.Writer) (err error) {
	if c.value > MaxUint40 {
		return errors.New("value exceeds 40-bit range")
	}
	// Buffer for the 5 bytes
	buf := make([]byte, 5)
	// Write the upper 5 bytes (ignoring the most significant 3 bytes of uint64)
	buf[0] = byte((c.value >> 32) & 0xFF) // Most significant byte
	buf[1] = byte((c.value >> 24) & 0xFF)
	buf[2] = byte((c.value >> 16) & 0xFF)
	buf[3] = byte((c.value >> 8) & 0xFF)
	buf[4] = byte(c.value & 0xFF) // Least significant byte
	_, err = w.Write(buf)
	return err
}

// UnmarshalRead reads 5 bytes from the provided reader and decodes it into a 40-bit unsigned integer.
func (c *Uint40) UnmarshalRead(r io.Reader) (err error) {
	// Buffer for the 5 bytes
	buf := make([]byte, 5)
	_, err = r.Read(buf)
	if err != nil {
		return err
	}
	// Decode the 5 bytes into a 40-bit unsigned integer
	c.value = (uint64(buf[0]) << 32) |
		(uint64(buf[1]) << 24) |
		(uint64(buf[2]) << 16) |
		(uint64(buf[3]) << 8) |
		uint64(buf[4])

	return nil
}
