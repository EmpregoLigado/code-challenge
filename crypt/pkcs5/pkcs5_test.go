package pkcs5

import (
	"encoding/hex"
	"testing"
)

var key = []byte{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x10, 0x11, 0x12, 0x13, 0x14, 0x15, 0x16}

func TestEncrypt(t *testing.T) {
	result, err := encrypt([]byte("password"), key)
	if err != nil {
		t.Fatalf(err.Error())
	}
	resultstring := hex.EncodeToString(result)
	if resultstring != "a5b1346a2f19bc886c60b408e0060d51" {
		t.Fatalf("expected a5b1346a2f19bc886c60b408e0060d51, got %s", resultstring)
	}
}

func TestErrors(t *testing.T) {
	_, err := encrypt([]byte("string"), []byte{0x01, 0x02, 0x03, 0x04, 0x05, 0x06})
	if err == nil {
		t.Fatalf("should have failed with invalid key size")
	}

	_, err = encrypt([]byte(""), key)
	if err == nil {
		t.Fatalf("should have failed with plain content empty")
	}

	_, err = decrypt([]byte("string"), []byte{0x01, 0x02, 0x03, 0x04, 0x05, 0x06})
	if err == nil {
		t.Fatalf("should have failed with invalid key size")
	}

	_, err = decrypt([]byte(""), key)
	if err == nil {
		t.Fatalf("should have failed with plain content empty")
	}

}

func TestDecrypt(t *testing.T) {
	crypt, err := hex.DecodeString("a5b1346a2f19bc886c60b408e0060d51")
	if err != nil {
		t.Fatalf(err.Error())
	}
	result, err := decrypt(crypt, key)
	if err != nil {
		t.Fatalf(err.Error())
	}
	if string(result) != "password" {
		t.Fatalf("expected password, got %s", string(result))
	}
}

func TestCipher(t *testing.T) {
	cipher := Cipher{Key: key}

	crypt, err := cipher.Encrypt("password")
	if err != nil {
		t.Fatalf(err.Error())
	}
	if crypt != "a5b1346a2f19bc886c60b408e0060d51" {
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

func TestCipherErrors(t *testing.T) {
	cipher := Cipher{Key: []byte{0x01, 0x02, 0x03, 0x04, 0x05, 0x06}}

	_, err := cipher.Encrypt("string")
	if err == nil {
		t.Fatalf("should have failed with invalid key size")
	}

	_, err = cipher.Encrypt("")
	if err == nil {
		t.Fatalf("should have failed with plain content empty")
	}

	_, err = cipher.Decrypt("string")
	if err == nil {
		t.Fatalf("should have failed with invalid key size")
	}

	_, err = cipher.Decrypt("")
	if err == nil {
		t.Fatalf("should have failed with plain content empty")
	}
}
