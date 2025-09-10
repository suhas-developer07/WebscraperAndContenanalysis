package main

import (
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

	searchHandler := search.NewHandler(esClient)

	http.HandleFunc("/search", searchHandler.HandleSearch)
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("✅ Healthy"))
	})

	log.Println("🚀 Query Service running on :8081")
	log.Fatal(http.ListenAndServe(":8081", nil))
}