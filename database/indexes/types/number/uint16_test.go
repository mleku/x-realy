package number

import (
	"bytes"
	"math"
	"testing"

	"lukechampine.com/frand"
)

func TestUint16(t *testing.T) {
	// Helper function to generate random 16-bit integers
	generateRandomUint16 := func() uint16 {
		return uint16(frand.Intn(math.MaxUint16)) // math.MaxUint16 == 65535
	}

	for i := 0; i < 100; i++ { // Run test 100 times for random values
		// Generate a random value
		randomUint16 := generateRandomUint16()
		randomInt := int(randomUint16)

		// Create a new encodedUint16
		encodedUint16 := new(Uint16)

		// Test UInt16 setter and getter
		encodedUint16.Set(randomUint16)
		if encodedUint16.Get() != randomUint16 {
			t.Fatalf("Get mismatch: got %d, expected %d", encodedUint16.Get(), randomUint16)
		}

		// Test GetInt setter and getter
		encodedUint16.SetInt(randomInt)
		if encodedUint16.GetInt() != randomInt {
			t.Fatalf("GetInt mismatch: got %d, expected %d", encodedUint16.GetInt(), randomInt)
		}

		// Test encoding to []byte and decoding back
		bufEnc := new(bytes.Buffer)

		// MarshalWrite
		err := encodedUint16.MarshalWrite(bufEnc)
		if err != nil {
			t.Fatalf("MarshalWrite failed: %v", err)
		}
		encoded := bufEnc.Bytes()

		// Create a copy of encoded bytes before decoding
		bufDec := bytes.NewBuffer(encoded)

		// Decode back the value
		decodedUint16 := new(Uint16)
		err = decodedUint16.UnmarshalRead(bufDec)
		if err != nil {
			t.Fatalf("UnmarshalRead failed: %v", err)
		}

		if decodedUint16.Get() != randomUint16 {
			t.Fatalf("Decoded value mismatch: got %d, expected %d", decodedUint16.Get(), randomUint16)
		}

		// Compare encoded bytes to ensure correctness
		if !bytes.Equal(encoded, bufEnc.Bytes()) {
			t.Fatalf("Byte encoding mismatch: got %v, expected %v", bufEnc.Bytes(), encoded)
		}
	}
}
