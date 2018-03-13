package null

import (
	"testing"
)

func TestCipher(t *testing.T) {
	cipher := Cipher{}

	crypt, err := cipher.Encrypt("password")
	if err != nil {
		t.Fatalf(err.Error())
	}
	if crypt != "password" {
		t.Fatalf("wrong encryption")
	}
	result, err := cipher.Decrypt(crypt)
	if err != nil {
		t.Fatalf(err.Error())
	}
	if result != "password" {
		t.Fatalf("expected password, got %s", result)
	}
}
