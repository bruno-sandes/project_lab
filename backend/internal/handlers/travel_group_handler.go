package handlers

import (
	"net/http"
	"project_lab/internal/repositories"
)

type TravelGroupHandler struct {
	repo repositories.TravelGroupRepository
}

func NewTravelGroupHandler(repo repositories.TravelGroupRepository) *TravelGroupHandler {
	return &TravelGroupHandler{repo: repo}
}

func (h *TravelGroupHandler) ListGroups(w http.ResponseWriter, r *http.Request) {

}
