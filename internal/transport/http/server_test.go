package httptransport

import (
	"bytes"
	"encoding/json"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	"log/slog"

	"github.com/gatiyatullin-ibragim/auto-recruiter/internal/domain"
)

func TestServerProfileLifecycle(t *testing.T) {
	t.Parallel()
	server := NewServer(slog.New(slog.NewTextHandler(bytes.NewBuffer(nil), nil)))
	ts := httptest.NewServer(server.Handler())
	defer ts.Close()

	response, err := http.Get(ts.URL + "/api/profile")
	if err != nil {
		t.Fatal(err)
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		t.Fatalf("GET /api/profile status = %d", response.StatusCode)
	}
	var candidate domain.Candidate
	if err := json.NewDecoder(response.Body).Decode(&candidate); err != nil {
		t.Fatal(err)
	}
	if candidate.ID != "default" {
		t.Fatalf("default profile ID = %q", candidate.ID)
	}
}

func TestServerRejectsInvalidVacancyURL(t *testing.T) {
	t.Parallel()
	server := NewServer(slog.Default())
	ts := httptest.NewServer(server.Handler())
	defer ts.Close()

	body := bytes.NewBufferString(`{"id":"v1","source":"demo","title":"Go","company":"Acme","url":"not-a-url"}`)
	response, err := http.Post(ts.URL+"/vacancies", "application/json", body)
	if err != nil {
		t.Fatal(err)
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusBadRequest {
		t.Fatalf("POST /vacancies status = %d, want 400", response.StatusCode)
	}
}

func TestServerResumeUpload(t *testing.T) {
	t.Parallel()
	server := NewServer(slog.Default())
	ts := httptest.NewServer(server.Handler())
	defer ts.Close()

	var body bytes.Buffer
	writer := multipart.NewWriter(&body)
	part, err := writer.CreateFormFile("file", "resume.txt")
	if err != nil {
		t.Fatal(err)
	}
	_, _ = part.Write([]byte("Go developer"))
	if err := writer.Close(); err != nil {
		t.Fatal(err)
	}
	request, err := http.NewRequest(http.MethodPost, ts.URL+"/api/resume", &body)
	if err != nil {
		t.Fatal(err)
	}
	request.Header.Set("Content-Type", writer.FormDataContentType())
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		t.Fatal(err)
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		t.Fatalf("POST /api/resume status = %d", response.StatusCode)
	}
}
