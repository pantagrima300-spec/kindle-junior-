package competition

import (
	"encoding/json"
	"net/http"

	"competition-backend/internal/storage"
)

type Handler struct {
	Storage *storage.Storage
}

func NewHandler(store *storage.Storage) *Handler {
	return &Handler{
		Storage: store,
	}
}

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	competitions, err := h.Storage.PocketBase.GetCompetitions(r.Context())

	if err != nil {
		http.Error(
			w,
			"failed to load competitions",
			http.StatusInternalServerError,
		)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(competitions.Items)
}

func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	if id == "" {
		http.Error(
			w,
			"competition id required",
			http.StatusBadRequest,
		)
		return
	}

	competition, err := h.Storage.PocketBase.GetCompetition(
		r.Context(),
		id,
	)

	if err != nil {
		http.Error(
			w,
			"competition not found",
			http.StatusNotFound,
		)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(competition)
}
