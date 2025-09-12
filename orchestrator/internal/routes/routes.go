package routes

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/suhas-developer07/WebscraperAndContenanalysis/internal/handlers"
	"github.com/suhas-developer07/WebscraperAndContenanalysis/internal/rabbitmq"
	"github.com/suhas-developer07/WebscraperAndContenanalysis/internal/repository"
)

func MountRoutes(repo repository.PostgresRepository, rabbitmq rabbitmq.RabbitmqRepo,) *mux.Router {
	handler := handlers.NewUrlHandler(&repo, &rabbitmq)

	router := mux.NewRouter()

	router.HandleFunc("/jobs", handler.InsertUrlHandler).Methods("POST")
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{
			"msg": "server is running in this port",
		})
	}).Methods("GET")

	return router
}
