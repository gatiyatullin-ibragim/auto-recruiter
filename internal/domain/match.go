package domain

type Match struct {
	// Match — результат сравнения профиля кандидата с вакансией.
	UserID    string   `json:"user_id"`
	VacancyID string   `json:"vacancy_id"`
	Score     float64  `json:"score"`
	Reasons   []string `json:"reasons"`
}

func (m Match) Accepted(threshold float64) bool {
	return m.Score >= threshold
}
