package number

import (
	"bytes"
	"math"
	"testing"

	"lukechampine.com/frand"
)

func TestUint32(t *testing.T) {
	// Helper function to generate random 32-bit integers
	generateRandomUint32 := func() uint32 {
		return uint32(frand.Intn(math.MaxUint32)) // math.MaxUint32 == 4294967295
	}

	for i := 0; i < 100; i++ { // Run test 100 times for random values
		// Generate a random value
		randomUint32 := generateRandomUint32()
		randomInt := int(randomUint32)

		// Create a new codec
		codec := new(Uint32)

		// Test UInt32 setter and getter
		codec.SetUint32(randomUint32)
		if codec.Uint32() != randomUint32 {
			t.Fatalf("Uint32 mismatch: got %d, expected %d", codec.Uint32(), randomUint32)
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

		// Create a copy of encoded bytes before decoding
		bufDec := bytes.NewBuffer(encoded)

		// Decode back the value
		decoded := new(Uint32)
		err = decoded.UnmarshalRead(bufDec)
		if err != nil {
			t.Fatalf("UnmarshalRead failed: %v", err)
		}

		if decoded.Uint32() != randomUint32 {
			t.Fatalf("Decoded value mismatch: got %d, expected %d", decoded.Uint32(), randomUint32)
		}

		// Compare encoded bytes to ensure correctness
		if !bytes.Equal(encoded, bufEnc.Bytes()) {
			t.Fatalf("Byte encoding mismatch: got %v, expected %v", bufEnc.Bytes(), encoded)
		}
	}
}
