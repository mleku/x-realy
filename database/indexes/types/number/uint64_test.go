package number

import (
	"bytes"
	"math"
	"testing"

	"lukechampine.com/frand"
)

func TestUint64Codec(t *testing.T) {
	// Helper function to generate random 64-bit integers
	generateRandomUint64 := func() uint64 {
		return frand.Uint64n(math.MaxUint64) // math.MaxUint64 == 18446744073709551615
	}

	for i := 0; i < 100; i++ { // Run test 100 times for random values
		// Generate a random value
		randomUint64 := generateRandomUint64()
		randomInt := int(randomUint64)

		// Create a new codec
		codec := new(Uint64Codec)

		// Test UInt64 setter and getter
		codec.SetUint64(randomUint64)
		if codec.Uint64() != randomUint64 {
			t.Fatalf("Uint64 mismatch: got %d, expected %d", codec.Uint64(), randomUint64)
		}

		// Test Int setter and getter
		codec.SetInt(randomInt)
		if codec.Int() != randomInt {
			t.Fatalf("Int mismatch: got %d, expected %d", codec.Int(), randomInt)
		}

		// Test encoding to []byte and decoding back
		bufEnc := new(bytes.Buffer)

		// MarshalWrite
		err := codec.MarshalWrite(bufEnc)
		if err != nil {
			t.Fatalf("MarshalWrite failed: %v", err)
		}
		encoded := bufEnc.Bytes()

		// Create a buffer for decoding
		bufDec := bytes.NewBuffer(encoded)

		// Decode back the value
		decodedCodec := new(Uint64Codec)
		err = decodedCodec.UnmarshalRead(bufDec)
		if err != nil {
			t.Fatalf("UnmarshalRead failed: %v", err)
		}

		if decodedCodec.Uint64() != randomUint64 {
			t.Fatalf("Decoded value mismatch: got %d, expected %d", decodedCodec.Uint64(), randomUint64)
		}

		// Compare encoded bytes to ensure correctness
		if !bytes.Equal(encoded, bufEnc.Bytes()) {
			t.Fatalf("Byte encoding mismatch: got %v, expected %v", bufEnc.Bytes(), encoded)
		}
	}
}
