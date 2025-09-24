package main

import (
	"bytes"
	"strings"
	"testing"
)

func TestAutoPadding(t *testing.T) {
	// Test decoding with padding
	encoded := "IUAjJCVeJg=="
	r := strings.NewReader(encoded)
	var buf bytes.Buffer
	if err := DecodeStream(false, false, r, &buf); err != nil {
		t.Errorf("decode %s failed: %s", encoded, err)
	}

	// Test decoding without padding
	encoded = "IUAjJCVeJg"
	r = strings.NewReader(encoded)
	buf.Reset()
	if err := DecodeStream(false, false, r, &buf); err != nil {
		t.Errorf("decode %s failed: %s", encoded, err)
	}
}

func TestEncodeStream(t *testing.T) {
	// Test basic encoding
	input := "Hello, World!"
	r := strings.NewReader(input)
	var buf bytes.Buffer

	if err := EncodeStream(false, false, 0, r, &buf); err != nil {
		t.Errorf("encode failed: %s", err)
	}

	// Verify we can decode it back
	output := strings.TrimSpace(buf.String())
	r2 := strings.NewReader(output)
	var buf2 bytes.Buffer

	if err := DecodeStream(false, false, r2, &buf2); err != nil {
		t.Errorf("decode back failed: %s", err)
	}

	decoded := strings.TrimSpace(buf2.String())
	if decoded != input {
		t.Errorf("encode/decode roundtrip failed: expected %q, got %q", input, decoded)
	}
}

func TestURLEncoding(t *testing.T) {
	// Test URL-safe encoding
	input := "Hello>World?"
	r := strings.NewReader(input)
	var buf bytes.Buffer

	if err := EncodeStream(true, false, 0, r, &buf); err != nil {
		t.Errorf("URL encode failed: %s", err)
	}

	// Verify we can decode it back with URL mode
	output := strings.TrimSpace(buf.String())
	r2 := strings.NewReader(output)
	var buf2 bytes.Buffer

	if err := DecodeStream(true, false, r2, &buf2); err != nil {
		t.Errorf("URL decode back failed: %s", err)
	}

	decoded := strings.TrimSpace(buf2.String())
	if decoded != input {
		t.Errorf("URL encode/decode roundtrip failed: expected %q, got %q", input, decoded)
	}
}
