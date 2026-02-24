package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/unnecessary-special-projects/ghist/internal/models"
)

type createEventRequest struct {
	Type    string `json:"type"`
	Message string `json:"message"`
	TaskID  *int64 `json:"task_id"`
}

func (s *Server) handleCreateEvent(w http.ResponseWriter, r *http.Request) {
	var req createEventRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON")
		return
	}
	if req.Message == "" {
		writeError(w, http.StatusBadRequest, "message is required")
		return
	}
	event, err := s.store.CreateEvent(req.Type, req.Message, "{}", req.TaskID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusCreated, event)
}

func (s *Server) handleListTaskEvents(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid task id")
		return
	}
	events, err := s.store.ListEventsByTask(id)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if events == nil {
		events = []models.Event{}
	}
	writeJSON(w, http.StatusOK, events)
}

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
