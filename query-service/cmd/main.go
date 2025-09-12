package main

import (
	"github.com/rs/cors"
	"log"
	"net/http"
	"query-service/internal/elasticsearch"
	"query-service/internal/search"
)

func main() {
	esClient, err := elasticsearch.NewClient("http://localhost:9200")
	if err != nil {
		log.Fatalf("Failed to connect to Elasticsearch: %v", err)
	}

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:5173"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
	})

	searchHandler := search.NewHandler(esClient)

	http.HandleFunc("/search", searchHandler.HandleSearch)
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(" Healthy"))
	})

	log.Println(" Query Service running on :17029")
	log.Fatal(http.ListenAndServe(":17029", c.Handler(http.DefaultServeMux)))
}
