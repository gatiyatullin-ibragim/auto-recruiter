package profile

import (
	"archive/zip"
	"bytes"
	"errors"
	"fmt"
	"io"
	"path/filepath"
	"regexp"
	"strings"
)

const MaxResumeSize = 5 << 20

func ExtractResumeText(filename, contentType string, data []byte) (string, error) {
	if len(data) == 0 || len(data) > MaxResumeSize {
		return "", errors.New("resume must be between 1 byte and 5 MB")
	}
	switch strings.ToLower(filepath.Ext(filename)) {
	case ".txt":
		return strings.TrimSpace(string(data)), nil
	case ".docx":
		return extractDOCX(data)
	case ".pdf":
		return extractPDF(data)
	default:
		if strings.Contains(contentType, "text/plain") {
			return strings.TrimSpace(string(data)), nil
		}
		return "", fmt.Errorf("unsupported resume format %q", filepath.Ext(filename))
	}
}

func extractDOCX(data []byte) (string, error) {
	reader, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		return "", fmt.Errorf("read docx: %w", err)
	}
	for _, file := range reader.File {
		if file.Name != "word/document.xml" {
			continue
		}
		opened, err := file.Open()
		if err != nil {
			return "", fmt.Errorf("open docx document: %w", err)
		}
		xmlData, readErr := io.ReadAll(io.LimitReader(opened, MaxResumeSize))
		_ = opened.Close()
		if readErr != nil {
			return "", fmt.Errorf("read docx document: %w", readErr)
		}
		text := regexp.MustCompile(`<[^>]+>`).ReplaceAllString(string(xmlData), " ")
		return strings.Join(strings.Fields(text), " "), nil
	}
	return "", errors.New("docx document.xml not found")
}

func extractPDF(data []byte) (string, error) {
	if !bytes.HasPrefix(data, []byte("%PDF-")) {
		return "", errors.New("invalid PDF header")
	}
	text := regexp.MustCompile(`\(([^()]*)\)\s*T[Jj]`).FindAllSubmatch(data, -1)
	parts := make([]string, 0, len(text))
	for _, match := range text {
		parts = append(parts, string(match[1]))
	}
	if len(parts) == 0 {
		return "", errors.New("PDF text extraction returned no text")
	}
	return strings.Join(parts, " "), nil
}
