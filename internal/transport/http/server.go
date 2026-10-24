package httptransport

import (
	"log/slog"
	nethttp "net/http"

	"github.com/go-chi/chi/v5"

	"github.com/gatiyatullin-ibragim/auto-recruiter/internal/handler"
	"github.com/gatiyatullin-ibragim/auto-recruiter/internal/repository/postgres"
	"github.com/gatiyatullin-ibragim/auto-recruiter/internal/scraper/demo"
	"github.com/gatiyatullin-ibragim/auto-recruiter/internal/service"
	"github.com/gatiyatullin-ibragim/auto-recruiter/internal/service/aggregator"
	"github.com/gatiyatullin-ibragim/auto-recruiter/internal/service/search"
	"github.com/gatiyatullin-ibragim/auto-recruiter/internal/transport/http/middleware"
)

type Server struct {
	// Server — composition root HTTP-транспорта и его зависимостей.
	logger  *slog.Logger
	handler *handler.HTTPHandler
}

func NewServer(logger *slog.Logger) *Server {
	// NewServer — создаёт in-memory dependencies для учебного запуска.
	vacancyStore := postgres.NewVacancyStore()
	candidateStore := postgres.NewCandidateStore()
	aggregatorService := aggregator.NewService(demo.NewClient())
	searchService := search.NewService(aggregatorService, vacancyStore)
	return &Server{
		logger: logger,
		handler: handler.NewHTTPHandler(
			service.NewVacancyService(vacancyStore),
			service.NewCandidateService(candidateStore), searchService,
		),
	}
}
func (s *Server) Handler() nethttp.Handler {
	// Handler — регистрирует маршруты и возвращает middleware chain.
	mux := chi.NewRouter()
	mux.Get("/health", s.handler.Health)
	mux.Get("/vacancies", s.handler.ListVacancies)
	mux.Post("/vacancies", s.handler.CreateVacancy)
	mux.Get("/vacancies/{id}", s.handler.GetVacancy)
	mux.Get("/api/profile", s.handler.GetProfile)
	mux.Put("/api/profile", s.handler.UpdateProfile)
	mux.Get("/api/criteria", s.handler.GetCriteria)
	mux.Put("/api/criteria", s.handler.UpdateCriteria)
	mux.Post("/api/resume", s.handler.UploadResume)
	mux.Post("/api/search", s.handler.StartSearch)
	mux.Get("/api/search/{jobID}/events", s.handler.SearchEvents)
	return middleware.SecurityHeaders(middleware.Logging(s.logger, mux))
}
