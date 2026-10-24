package security

import (
	"bytes"
	"testing"
)

func TestEncryptDecryptAndTamperDetection(t *testing.T) {
	key := []byte("01234567890123456789012345678901")
	first, err := Encrypt(key, []byte("secret"))
	if err != nil {
		t.Fatal(err)
	}
	second, err := Encrypt(key, []byte("secret"))
	if err != nil {
		t.Fatal(err)
	}
	if bytes.Equal(first, second) {
		t.Fatal("encryption reused ciphertext")
	}
	got, err := Decrypt(key, first)
	if err != nil || string(got) != "secret" {
		t.Fatalf("decrypt = %q, %v", got, err)
	}
	first[len(first)-1] ^= 1
	if _, err := Decrypt(key, first); err == nil {
		t.Fatal("tampered ciphertext was accepted")
	}
}
