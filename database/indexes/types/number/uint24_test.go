package number

import (
	"bytes"
	"testing"
)

func TestUint24Codec(t *testing.T) {
	tests := []struct {
		name        string
		value       uint32
		expectedErr bool
	}{
		{"Minimum Value", 0, false},
		{"Maximum Value", MaxUint24, false},
		{"Value in Range", 8374263, false},           // Example value within the range
		{"Value Exceeds Range", MaxUint24 + 1, true}, // Exceeds 24-bit limit
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			codec := new(Uint24Codec)

			// Test SetUint24
			err := codec.SetUint24(tt.value)
			if tt.expectedErr {
				if err == nil {
					t.Errorf("expected error but got none")
				}
				return
			} else if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			// Test Uint24 getter
			if codec.Uint24() != tt.value {
				t.Errorf("Uint24 mismatch: got %d, expected %d", codec.Uint24(), tt.value)
			}

			// Test MarshalWrite and UnmarshalRead
			buf := new(bytes.Buffer)

			// MarshalWrite directly to the buffer
			if err := codec.MarshalWrite(buf); err != nil {
				t.Fatalf("MarshalWrite failed: %v", err)
			}

			// Validate encoded size is 3 bytes
			encoded := buf.Bytes()
			if len(encoded) != 3 {
				t.Fatalf("encoded size mismatch: got %d bytes, expected 3 bytes", len(encoded))
			}

			// Decode from the buffer
			decodedCodec := new(Uint24Codec)
			if err := decodedCodec.UnmarshalRead(buf); err != nil {
				t.Fatalf("UnmarshalRead failed: %v", err)
			}

			// Validate decoded value
			if decodedCodec.Uint24() != tt.value {
				t.Errorf("Decoded value mismatch: got %d, expected %d", decodedCodec.Uint24(), tt.value)
			}
		})
	}
}
