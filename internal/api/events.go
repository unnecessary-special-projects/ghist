package api

import (
	"net/http"
	"strconv"

	"github.com/coderstone/ghist/internal/models"
)

func (s *Server) handleListEvents(w http.ResponseWriter, r *http.Request) {
	limit := 20
	if l := r.URL.Query().Get("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil && parsed > 0 {
			limit = parsed
		}
	}

	events, err := s.store.ListEvents(limit)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if events == nil {
		events = []models.Event{}
	}

	writeJSON(w, http.StatusOK, events)
}

func (s *Server) handleStatus(w http.ResponseWriter, r *http.Request) {
	counts, err := s.store.TaskCountsByStatus()
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	milestones, err := s.store.MilestoneInfo()
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	events, err := s.store.ListEvents(5)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	total := 0
	for _, c := range counts {
		total += c
	}

	if milestones == nil {
		milestones = []models.MilestoneInfo{}
	}
	if events == nil {
		events = []models.Event{}
	}

	summary := models.StatusSummary{
		TotalTasks:    total,
		TasksByStatus: counts,
		Milestones:    milestones,
		RecentEvents:  events,
	}

	writeJSON(w, http.StatusOK, summary)
}
