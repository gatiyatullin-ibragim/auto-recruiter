package contacts

import (
	"net/url"
	"regexp"
	"strings"
	"time"
)

type Contact struct {
	Email      string    `json:"email"`
	Type       string    `json:"type"`
	SourceURL  string    `json:"source_url"`
	Confidence string    `json:"confidence"`
	FoundAt    time.Time `json:"found_at"`
}

var emailPattern = regexp.MustCompile(`(?i)\b[A-Z0-9._%+\-]+@[A-Z0-9.\-]+\.[A-Z]{2,}\b`)

func Extract(text, sourceURL string) []Contact {
	text = strings.NewReplacer("[at]", "@", "(at)", "@", " at ", "@", "[dot]", ".", "(dot)", ".", " dot ", ".").Replace(text)
	text = regexp.MustCompile(`\s*@\s*`).ReplaceAllString(text, "@")
	text = regexp.MustCompile(`\s*\.\s*`).ReplaceAllString(text, ".")
	seen := make(map[string]struct{})
	result := make([]Contact, 0)
	for _, email := range emailPattern.FindAllString(strings.ToLower(text), -1) {
		if !allowed(email) {
			continue
		}
		if _, ok := seen[email]; ok {
			continue
		}
		seen[email] = struct{}{}
		result = append(result, Contact{Email: email, Type: classify(email), SourceURL: sourceURL, Confidence: "high", FoundAt: time.Now().UTC()})
	}
	return result
}

func allowed(email string) bool {
	parts := strings.Split(email, "@")
	if len(parts) != 2 || strings.Contains(parts[1], ".css") || strings.Contains(parts[1], ".js") {
		return false
	}
	return !strings.HasPrefix(parts[0], "noreply") && !strings.HasPrefix(parts[0], "no-reply") &&
		!strings.HasPrefix(parts[0], "sentry") && parts[1] != "example.com"
}

func classify(email string) string {
	local := strings.Split(email, "@")[0]
	if strings.Contains(local, "hr") || strings.Contains(local, "recruit") || strings.Contains(local, "career") || strings.Contains(local, "job") {
		return "hr"
	}
	return "general"
}

func SourceHost(sourceURL string) string {
	parsed, err := url.Parse(sourceURL)
	if err != nil {
		return ""
	}
	return parsed.Hostname()
}
