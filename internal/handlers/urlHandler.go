package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/suhas-developer07/WebscraperAndContenanalysis/internal/rabbitmq"
	"github.com/suhas-developer07/WebscraperAndContenanalysis/internal/repository"
)

type UrlHandler struct {
	Repo     *repository.PostgresRepository
	Rabbitmq *rabbitmq.RabbitmqRepo
}

func NewUrlHandler(repo *repository.PostgresRepository, rabbitmq *rabbitmq.RabbitmqRepo) *UrlHandler {
	return &UrlHandler{Repo: repo, Rabbitmq: rabbitmq}
}

func (h *UrlHandler) InsertUrlHandler(w http.ResponseWriter, r *http.Request) {
	var payload repository.InsertJobsPayload

	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		w.Header().Set("content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"message": "invalid request body formate",
		})
		return
	}
	id, urls, err := h.Repo.InsertJobs(payload)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	msg := rabbitmq.Data{
		ID:  id,
		URL: urls,
	}

	if err := h.Rabbitmq.SendURL(msg); err != nil {
		http.Error(w, "failed to publish to rabbitmq"+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"msg": "urls submitted succesfully",
	})
}
