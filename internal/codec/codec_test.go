package codec

import (
	"bytes"
	"strings"
	"testing"
)

func TestEncryptDecryptRoundTrip(t *testing.T) {
	plaintext := []byte(`{"project":"qed","ok":true}`)
	packet, err := Encrypt(plaintext, "password")
	if err != nil {
		t.Fatalf("Encrypt() error = %v", err)
	}

	if len(packet) <= headerSize+signatureSize+resignatureSize {
		t.Fatalf("packet length = %d, want body", len(packet))
	}
	if !bytes.Equal(packet[:magicSize], magic[:]) {
		t.Fatalf("magic = %q, want %q", packet[:magicSize], magic)
	}
	if !bytes.Equal(packet[magicSize+versionSize:headerSize], supportedMethod[:]) {
		t.Fatalf("method header mismatch")
	}

	got, err := Decrypt(packet, "password")
	if err != nil {
		t.Fatalf("Decrypt() error = %v", err)
	}
	if !bytes.Equal(got, plaintext) {
		t.Fatalf("Decrypt() = %q, want %q", got, plaintext)
	}
}

func TestDecryptRejectsTamperedBodySignature(t *testing.T) {
	packet, err := Encrypt([]byte("body"), "password")
	if err != nil {
		t.Fatalf("Encrypt() error = %v", err)
	}

	tampered := append([]byte(nil), packet...)
	tampered[headerSize+saltSize] ^= 0xff

	_, err = Decrypt(tampered, "password")
	if err == nil || !strings.Contains(err.Error(), "signature") {
		t.Fatalf("Decrypt() error = %v, want signature error", err)
	}
}

func TestDecryptRejectsTamperedToolchainSignature(t *testing.T) {
	packet, err := Encrypt([]byte("body"), "password")
	if err != nil {
		t.Fatalf("Encrypt() error = %v", err)
	}

	tampered := append([]byte(nil), packet...)
	tampered[len(tampered)-1] ^= 0xff

	_, err = Decrypt(tampered, "password")
	if err == nil || !strings.Contains(err.Error(), "re-signature") {
		t.Fatalf("Decrypt() error = %v, want re-signature error", err)
	}
}

func TestDecryptRejectsWrongPassword(t *testing.T) {
	packet, err := Encrypt([]byte("secret"), "password")
	if err != nil {
		t.Fatalf("Encrypt() error = %v", err)
	}

	_, err = Decrypt(packet, "wrong-password")
	if err == nil || !strings.Contains(err.Error(), "authentication failed") {
		t.Fatalf("Decrypt() error = %v, want authentication failure", err)
	}
}
