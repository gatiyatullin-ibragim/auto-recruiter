package handler

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"

	"github.com/gatiyatullin-ibragim/auto-recruiter/internal/domain"
	"github.com/gatiyatullin-ibragim/auto-recruiter/internal/service"
	"github.com/gatiyatullin-ibragim/auto-recruiter/internal/service/profile"
	"github.com/gatiyatullin-ibragim/auto-recruiter/internal/service/search"
)

type HTTPHandler struct {
	vacancies  *service.VacancyService
	candidates *service.CandidateService
	search     *search.Service
}

const defaultCandidateID = "default"

func NewHTTPHandler(vacancies *service.VacancyService, candidates *service.CandidateService, searchService *search.Service) *HTTPHandler {
	return &HTTPHandler{vacancies: vacancies, candidates: candidates, search: searchService}
}

func (h *HTTPHandler) StartSearch(w http.ResponseWriter, r *http.Request) {
	candidate, err := h.profile(r)
	if err != nil {
		writeError(w, err)
		return
	}
	var criteria domain.SearchCriteria
	if err := decodeJSON(r, &criteria); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	if criteria.Position == "" {
		criteria = candidate.Criteria
	}
	writeJSON(w, http.StatusAccepted, map[string]string{"job_id": h.search.Start(r.Context(), criteria)})
}

func (h *HTTPHandler) SearchEvents(w http.ResponseWriter, r *http.Request) {
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "streaming unsupported", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	for event := range h.search.Subscribe(r.Context(), chi.URLParam(r, "jobID")) {
		data, err := search.EncodeEvent(event)
		if err != nil {
			continue
		}
		_, _ = w.Write([]byte("event: progress\ndata: "))
		_, _ = w.Write(data)
		_, _ = w.Write([]byte("\n\n"))
		flusher.Flush()
		if event.Status == "done" {
			return
		}
	}
}

func (h *HTTPHandler) Health(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (h *HTTPHandler) GetVacancy(w http.ResponseWriter, r *http.Request) {
	item, err := h.vacancies.Get(r.Context(), chi.URLParam(r, "id"))
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, item)
}

func (h *HTTPHandler) ListVacancies(w http.ResponseWriter, r *http.Request) {
	criteria := domain.SearchCriteria{
		Remote:         r.URL.Query().Get("remote") == "true",
		EmploymentType: r.URL.Query().Get("employment_type"),
	}
	items, err := h.vacancies.List(r.Context(), criteria)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, items)
}

func (h *HTTPHandler) CreateVacancy(w http.ResponseWriter, r *http.Request) {
	var vacancy domain.Vacancy
	if err := decodeJSON(r, &vacancy); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	if strings.TrimSpace(vacancy.ID) == "" {
		vacancy.ID = vacancy.ExternalID
	}
	if err := h.vacancies.Create(r.Context(), vacancy); err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, vacancy.Normalize())
}

func (h *HTTPHandler) GetProfile(w http.ResponseWriter, r *http.Request) {
	candidate, err := h.profile(r)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, candidate)
}

func (h *HTTPHandler) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	var candidate domain.Candidate
	if err := decodeJSON(r, &candidate); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	candidate.ID = defaultCandidateID
	if err := candidate.Validate(); err != nil {
		writeError(w, err)
		return
	}
	if err := h.candidates.Save(r.Context(), candidate); err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, candidate)
}

func (h *HTTPHandler) GetCriteria(w http.ResponseWriter, r *http.Request) {
	candidate, err := h.profile(r)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, candidate.Criteria)
}

func (h *HTTPHandler) UpdateCriteria(w http.ResponseWriter, r *http.Request) {
	var criteria domain.SearchCriteria
	if err := decodeJSON(r, &criteria); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	candidate, err := h.profile(r)
	if err != nil {
		writeError(w, err)
		return
	}
	candidate.Criteria = criteria
	if err := candidate.Validate(); err != nil {
		writeError(w, err)
		return
	}
	if err := h.candidates.Save(r.Context(), candidate); err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, criteria)
}

func (h *HTTPHandler) UploadResume(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseMultipartForm(profile.MaxResumeSize); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid multipart form"})
		return
	}
	file, header, err := r.FormFile("file")
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "file field is required"})
		return
	}
	defer file.Close()
	data, err := io.ReadAll(io.LimitReader(file, profile.MaxResumeSize+1))
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "read resume: " + err.Error()})
		return
	}
	text, err := profile.ExtractResumeText(header.Filename, header.Header.Get("Content-Type"), data)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	candidate, err := h.profile(r)
	if err != nil {
		writeError(w, err)
		return
	}
	candidate.ResumeText = text
	if err := h.candidates.Save(r.Context(), candidate); err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"filename": header.Filename, "characters": len(text)})
}

func (h *HTTPHandler) profile(r *http.Request) (domain.Candidate, error) {
	candidate, err := h.candidates.Get(r.Context(), defaultCandidateID)
	if errors.Is(err, domain.ErrNotFound) {
		candidate = domain.Candidate{
			ID: defaultCandidateID, Name: "Candidate", Email: "candidate@example.com",
		}
		if saveErr := h.candidates.Save(r.Context(), candidate); saveErr != nil {
			return domain.Candidate{}, saveErr
		}
		return candidate, nil
	}
	return candidate, err
}

func decodeJSON(r *http.Request, target any) error {
	r.Body = io.NopCloser(io.LimitReader(r.Body, 1<<20))
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(target); err != nil {
		return err
	}
	var extra any
	if err := decoder.Decode(&extra); err != io.EOF {
		return errors.New("request body must contain one JSON object")
	}
	return nil
}

func writeError(w http.ResponseWriter, err error) {
	status := http.StatusInternalServerError
	switch {
	case errors.Is(err, domain.ErrNotFound):
		status = http.StatusNotFound
	case errors.Is(err, domain.ErrInvalidInput):
		status = http.StatusBadRequest
	case errors.Is(err, domain.ErrConflict):
		status = http.StatusConflict
	}
	writeJSON(w, status, map[string]string{"error": err.Error()})
}

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}
