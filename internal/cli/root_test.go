package cli

import (
	"bytes"
	"io/fs"
	"strings"
	"testing"
	"time"
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

func TestRootCommandRejectsEncryptToTerminal(t *testing.T) {
	stdout := &terminalWriter{}
	cmd := NewRootCommand(strings.NewReader("secret"), stdout, &bytes.Buffer{})
	cmd.SetArgs([]string{"-e", "password"})

	err := cmd.Execute()
	if err == nil || !strings.Contains(err.Error(), "binary output to terminal") {
		t.Fatalf("Execute() error = %v, want terminal binary output error", err)
	}
	if stdout.Len() != 0 {
		t.Fatalf("stdout length = %d, want 0", stdout.Len())
	}
}

func TestRootCommandVersion(t *testing.T) {
	var stdout bytes.Buffer
	cmd := NewRootCommand(strings.NewReader(""), &stdout, &bytes.Buffer{})
	cmd.SetArgs([]string{"version"})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}

	got := stdout.String()
	for _, want := range []string{"qed dev", "commit: none", "built: unknown"} {
		if !strings.Contains(got, want) {
			t.Fatalf("version output = %q, want to contain %q", got, want)
		}
	}
}

type terminalWriter struct {
	bytes.Buffer
}

func (w *terminalWriter) Stat() (fs.FileInfo, error) {
	return terminalFileInfo{}, nil
}

type terminalFileInfo struct{}

func (terminalFileInfo) Name() string       { return "stdout" }
func (terminalFileInfo) Size() int64        { return 0 }
func (terminalFileInfo) Mode() fs.FileMode  { return fs.ModeCharDevice }
func (terminalFileInfo) ModTime() time.Time { return time.Time{} }
func (terminalFileInfo) IsDir() bool        { return false }
func (terminalFileInfo) Sys() any           { return nil }
