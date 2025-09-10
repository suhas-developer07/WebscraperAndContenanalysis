package search

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

type ElasticsearchClient interface {
	Search(index string, query map[string]interface{}) (map[string]interface{}, error)
}

type Handler struct {
	esClient ElasticsearchClient
}

func NewHandler(esClient ElasticsearchClient) *Handler {
	return &Handler{esClient: esClient}
}

func (h *Handler) HandleSearch(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	category := r.URL.Query().Get("category")
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	pageSize, _ := strconv.Atoi(r.URL.Query().Get("page_size"))
	
	if page < 1 { page = 1 }
	if pageSize < 1 || pageSize > 100 { pageSize = 10 }

	searchReq := SearchRequest{
		Query:    query,
		Category: category,
		Page:     page,
		PageSize: pageSize,
	}

	esQuery := BuildQuery(searchReq)
	result, err := h.esClient.Search("analyzed_content", esQuery)
	if err != nil {
		http.Error(w, fmt.Sprintf("Search failed: %v", err), http.StatusInternalServerError)
		return
	}

	// Parse the response (this is the standard Elasticsearch response format)
	hits, ok := result["hits"].(map[string]interface{})
	if !ok {
		http.Error(w, "Invalid search response format", http.StatusInternalServerError)
		return
	}

	// Extract total hits
	total := int64(0)
	if totalObj, ok := hits["total"].(map[string]interface{}); ok {
		if totalVal, ok := totalObj["value"].(float64); ok {
			total = int64(totalVal)
		}
	}

	// Extract search results
	results := make([]map[string]interface{}, 0)
	if hitList, ok := hits["hits"].([]interface{}); ok {
		for _, hit := range hitList {
			if hitMap, ok := hit.(map[string]interface{}); ok {
				results = append(results, hitMap)
			}
		}
	}

	// Prepare response
	response := map[string]interface{}{
		"results":  results,
		"total":    total,
		"page":     page,
		"page_size": pageSize,
		"has_more": (int64(page)*int64(pageSize)) < total,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}