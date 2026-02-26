package api

import (
	"encoding/json"
	"net/http"

	"github.com/unnecessary-special-projects/ghist/internal/models"
	"github.com/unnecessary-special-projects/ghist/internal/store"
)

func (s *Server) handleListTasks(w http.ResponseWriter, r *http.Request) {
	status := r.URL.Query().Get("status")
	milestone := r.URL.Query().Get("milestone")
	priority := r.URL.Query().Get("priority")
	taskType := r.URL.Query().Get("type")

	tasks, err := s.store.ListTasks(status, milestone, priority, taskType)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if tasks == nil {
		tasks = []models.Task{}
	}

	writeJSON(w, http.StatusOK, tasks)
}

type createTaskRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Status      string `json:"status"`
	Milestone   string `json:"milestone"`
	Priority    string `json:"priority"`
	Type        string `json:"type"`
	LegacyID    string `json:"legacy_id"`
}

func (s *Server) handleCreateTask(w http.ResponseWriter, r *http.Request) {
	var req createTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON")
		return
	}

	if req.Title == "" {
		writeError(w, http.StatusBadRequest, "title is required")
		return
	}

	task, err := s.store.CreateTask(store.CreateTaskInput{
		Title:       req.Title,
		Description: req.Description,
		Status:      req.Status,
		Milestone:   req.Milestone,
		Priority:    req.Priority,
		Type:        req.Type,
		LegacyID:    req.LegacyID,
	})
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusCreated, task)
}

func (s *Server) handleGetTask(w http.ResponseWriter, r *http.Request) {
	id, err := models.ParseTaskID(r.PathValue("id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid task id")
		return
	}

	task, err := s.store.GetTask(id)
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, task)
}

type updateTaskRequest struct {
	Title       *string `json:"title"`
	Description *string `json:"description"`
	Plan        *string `json:"plan"`
	Status      *string `json:"status"`
	Milestone   *string `json:"milestone"`
	CommitHash  *string `json:"commit_hash"`
	Priority    *string `json:"priority"`
	Type        *string `json:"type"`
	LegacyID    *string `json:"legacy_id"`
}

func (s *Server) handleUpdateTask(w http.ResponseWriter, r *http.Request) {
	id, err := models.ParseTaskID(r.PathValue("id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid task id")
		return
	}

	var req updateTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON")
		return
	}

	task, err := s.store.UpdateTask(id, store.TaskUpdate{
		Title:       req.Title,
		Description: req.Description,
		Plan:        req.Plan,
		Status:      req.Status,
		Milestone:   req.Milestone,
		CommitHash:  req.CommitHash,
		Priority:    req.Priority,
		Type:        req.Type,
		LegacyID:    req.LegacyID,
	})
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, task)
}

func (s *Server) handleDeleteTask(w http.ResponseWriter, r *http.Request) {
	id, err := models.ParseTaskID(r.PathValue("id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid task id")
		return
	}

	if err := s.store.DeleteTask(id); err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"status": "deleted"})
}
