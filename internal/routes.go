package routes

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/suhas-developer07/WebscraperAndContenanalysis/internal/handlers"
	"github.com/suhas-developer07/WebscraperAndContenanalysis/internal/repository"
)

func MountRoutes(repo repository.PostgresRepository) *mux.Router {
	handler := handlers.NewUrlHandler(&repo)

	router := mux.NewRouter()

	fmt.Print("woutes working")

	router.HandleFunc("/", handler.InsertUrlHandler).Methods("POST")
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{
			"msg": "server is running in this port",
		})
	}).Methods("GET")

	return router
}
