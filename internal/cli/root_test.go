package cli

import (
	"bytes"
	"strings"
	"testing"
)

func TestRootCommandEncryptDecrypt(t *testing.T) {
	var encrypted bytes.Buffer
	var stderr bytes.Buffer

	encrypt := NewRootCommand(strings.NewReader("hello qed"), &encrypted, &stderr)
	encrypt.SetArgs([]string{"-e", "password"})
	if err := encrypt.Execute(); err != nil {
		t.Fatalf("encrypt Execute() error = %v; stderr = %s", err, stderr.String())
	}
	if encrypted.Len() == 0 {
		t.Fatalf("encrypt output is empty")
	}

	var decrypted bytes.Buffer
	decrypt := NewRootCommand(bytes.NewReader(encrypted.Bytes()), &decrypted, &stderr)
	decrypt.SetArgs([]string{"-d", "password"})
	if err := decrypt.Execute(); err != nil {
		t.Fatalf("decrypt Execute() error = %v; stderr = %s", err, stderr.String())
	}
	if got := decrypted.String(); got != "hello qed" {
		t.Fatalf("decrypt output = %q, want %q", got, "hello qed")
	}
}

func TestRootCommandRequiresSingleMode(t *testing.T) {
	cmd := NewRootCommand(strings.NewReader("input"), &bytes.Buffer{}, &bytes.Buffer{})
	cmd.SetArgs([]string{"password"})

	err := cmd.Execute()
	if err == nil || !strings.Contains(err.Error(), "exactly one") {
		t.Fatalf("Execute() error = %v, want mode validation", err)
	}
}
