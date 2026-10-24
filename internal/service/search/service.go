package search

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"sync/atomic"

	"github.com/gatiyatullin-ibragim/auto-recruiter/internal/domain"
	"github.com/gatiyatullin-ibragim/auto-recruiter/internal/scraper"
	"github.com/gatiyatullin-ibragim/auto-recruiter/internal/service/aggregator"
)

type Event struct {
	JobID   string `json:"job_id"`
	Source  string `json:"source,omitempty"`
	Status  string `json:"status"`
	Found   int    `json:"found,omitempty"`
	Message string `json:"message,omitempty"`
}

type Service struct {
	aggregator *aggregator.Service
	vacancies  domain.VacancyRepository
	mu         sync.RWMutex
	subs       map[string]map[chan Event]struct{}
	history    map[string][]Event
	sequence   atomic.Uint64
}

func NewService(a *aggregator.Service, vacancies domain.VacancyRepository) *Service {
	return &Service{aggregator: a, vacancies: vacancies, subs: make(map[string]map[chan Event]struct{}), history: make(map[string][]Event)}
}

func (s *Service) Start(ctx context.Context, criteria domain.SearchCriteria) string {
	jobID := fmt.Sprintf("search-%d", s.sequence.Add(1))
	go s.run(context.WithoutCancel(ctx), jobID, criteria)
	return jobID
}

func (s *Service) run(ctx context.Context, jobID string, criteria domain.SearchCriteria) {
	s.publish(Event{JobID: jobID, Status: "started", Message: "Поиск запущен"})
	query := scraper.Query{Keywords: append([]string{criteria.Position}, criteria.Roles...), Remote: criteria.Remote}
	items, err := s.aggregator.Collect(ctx, query)
	for _, item := range items {
		if saveErr := s.vacancies.Save(ctx, item); saveErr != nil {
			s.publish(Event{JobID: jobID, Status: "failed", Message: saveErr.Error()})
			return
		}
	}
	if err != nil {
		s.publish(Event{JobID: jobID, Status: "partial", Found: len(items), Message: err.Error()})
	} else {
		s.publish(Event{JobID: jobID, Status: "completed", Found: len(items)})
	}
	s.publish(Event{JobID: jobID, Status: "done"})
}

func (s *Service) Subscribe(ctx context.Context, jobID string) <-chan Event {
	ch := make(chan Event, 16)
	s.mu.Lock()
	for _, event := range s.history[jobID] {
		ch <- event
	}
	if s.subs[jobID] == nil {
		s.subs[jobID] = make(map[chan Event]struct{})
	}
	s.subs[jobID][ch] = struct{}{}
	s.mu.Unlock()
	go func() {
		<-ctx.Done()
		s.mu.Lock()
		delete(s.subs[jobID], ch)
		close(ch)
		s.mu.Unlock()
	}()
	return ch
}

func (s *Service) publish(event Event) {
	s.mu.Lock()
	s.history[event.JobID] = append(s.history[event.JobID], event)
	defer s.mu.Unlock()
	for ch := range s.subs[event.JobID] {
		select {
		case ch <- event:
		default:
		}
	}
}

func EncodeEvent(event Event) ([]byte, error) {
	return json.Marshal(event)
}
