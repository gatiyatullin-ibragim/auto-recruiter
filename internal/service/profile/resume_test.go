package profile

import (
	"strings"
	"testing"
)

func TestExtractResumeText(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name, filename, contentType, data, want string
	}{
		{"text", "resume.txt", "text/plain", " Ada Lovelace ", "Ada Lovelace"},
		{"text by content type", "resume", "text/plain", " Ada ", "Ada"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ExtractResumeText(tt.filename, tt.contentType, []byte(tt.data))
			if err != nil || got != tt.want {
				t.Fatalf("ExtractResumeText() = %q, %v; want %q, nil", got, err, tt.want)
			}
		})
	}
}

func TestExtractResumeTextRejectsInvalidInputs(t *testing.T) {
	t.Parallel()
	if _, err := ExtractResumeText("resume.exe", "application/octet-stream", []byte("data")); err == nil {
		t.Fatal("unsupported extension was accepted")
	}
	if _, err := ExtractResumeText("resume.txt", "text/plain", nil); err == nil {
		t.Fatal("empty resume was accepted")
	}
	if _, err := ExtractResumeText("resume.pdf", "application/pdf", []byte("not pdf")); err == nil {
		t.Fatal("invalid PDF was accepted")
	}
	if _, err := ExtractResumeText("resume.txt", "text/plain", []byte(strings.Repeat("x", MaxResumeSize+1))); err == nil {
		t.Fatal("oversized resume was accepted")
	}
}
