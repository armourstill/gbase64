package gbase64

import "testing"

func TestAutoPadding(t *testing.T) {
	encoded := "IUAjJCVeJg=="
	if _, err := Decode(STD, []byte(encoded)); err != nil {
		t.Errorf("decode %s failed: %s", encoded, err)
	}
	encoded = "IUAjJCVeJg"
	if _, err := Decode(STD, []byte(encoded)); err != nil {
		t.Errorf("decode %s failed: %s", encoded, err)
	}
}
