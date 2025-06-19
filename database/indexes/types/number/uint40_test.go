package number

import (
	"bytes"
	"testing"
)

func TestUint40Codec(t *testing.T) {
	// Test cases for Uint40Codec
	tests := []struct {
		name        string
		value       uint64
		expectedErr bool
	}{
		{"Minimum Value", 0, false},
		{"Maximum Value", MaxUint40, false},
		{"Value in Range", 109951162777, false}, // Example value within the range
		{"Value Exceeds Range", MaxUint40 + 1, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			codec := new(Uint40Codec)

			// Test SetUint40
			err := codec.SetUint40(tt.value)
			if tt.expectedErr {
				if err == nil {
					t.Errorf("expected error but got none")
				}
				return
			} else if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			// Test Uint40 getter
			if codec.Uint40() != tt.value {
				t.Errorf("Uint40 mismatch: got %d, expected %d", codec.Uint40(), tt.value)
			}

			// Test MarshalWrite and UnmarshalRead
			buf := new(bytes.Buffer)

			// Marshal to a buffer
			if err = codec.MarshalWrite(buf); err != nil {
				t.Fatalf("MarshalWrite failed: %v", err)
			}

			// Validate encoded size is 5 bytes
			encoded := buf.Bytes()
			if len(encoded) != 5 {
				t.Fatalf("encoded size mismatch: got %d bytes, expected 5 bytes", len(encoded))
			}

			// Decode from the buffer
			decodedCodec := new(Uint40Codec)
			if err = decodedCodec.UnmarshalRead(buf); err != nil {
				t.Fatalf("UnmarshalRead failed: %v", err)
			}

			// Validate decoded value
			if decodedCodec.Uint40() != tt.value {
				t.Errorf("Decoded value mismatch: got %d, expected %d", decodedCodec.Uint40(), tt.value)
			}
		})
	}
}
