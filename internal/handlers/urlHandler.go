package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/suhas-developer07/WebscraperAndContenanalysis/internal/repository"
)

type UrlHandler struct {
	Repo *repository.PostgresRepository
}

func NewUrlHandler(repo *repository.PostgresRepository) *UrlHandler {
	return &UrlHandler{Repo: repo}
}

func (h *UrlHandler) InsertUrlHandler(w http.ResponseWriter, r *http.Request) {
	var payload repository.InsertJobsPayload

	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"msg": "invalid body formate",
		})
		return
	}

	if err := h.Repo.InsertJobs(payload); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"msg": "urls submitted successfully",
	})
}
