package api

import (
	"encoding/json"
	"net/http"
)

func (s *Server) handleGetMilestoneOrder(w http.ResponseWriter, r *http.Request) {
	order, err := s.store.GetMilestoneOrder()
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, order)
}

func (s *Server) handleSetMilestoneOrder(w http.ResponseWriter, r *http.Request) {
	var order []string
	if err := json.NewDecoder(r.Body).Decode(&order); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON array")
		return
	}
	if err := s.store.SetMilestoneOrder(order); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	s.hub.broadcast()
	writeJSON(w, http.StatusOK, order)
}
