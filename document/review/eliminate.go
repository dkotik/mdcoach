package review

import (
	"log/slog"
	"math"
)

func (r *Review) KeepAtMostQuestions(limit int) {
	if limit < 1 {
		limit = 1
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	if limit > len(r.questions) {
		return
	}
	for _, eliminate := range r.questions[limit:] {
		delete(r.known, eliminate)
	}
	r.logger.Info(
		"reduced question stack size",
		slog.Int("original", len(r.questions)),
		slog.Int("eliminated", len(r.questions)-limit),
		slog.Int("remaining", limit),
	)
	r.questions = r.questions[:limit]
}

func (r *Review) KeepAtMostPercent(limit float64) {
	if limit >= 100 {
		return
	}
	if limit < 1 {
		limit = 1
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	r.KeepAtMostQuestions(int(
		math.Ceil(limit * float64(len(r.questions)) / 100),
	))
}
