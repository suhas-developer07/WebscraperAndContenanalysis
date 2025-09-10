package models

type SearchResult struct {
	JobID          int      `json:"job_id"`
	TaskID         int      `json:"task_id"`
	URL            string   `json:"url"`
	ContentType    string   `json:"content_type"`
	DomainCategory string   `json:"domain_category"`
	Summary        string   `json:"summary"`
	KeyEntities    []string `json:"key_entities"`
	SentimentTone  string   `json:"sentiment_tone"`
	Score          float64  `json:"score,omitempty"`
}

type SearchResponse struct {
	Results  []SearchResult `json:"results"`
	Total    int64          `json:"total"`
	Page     int            `json:"page"`
	PageSize int            `json:"page_size"`
	HasMore  bool           `json:"has_more"`
}

type SearchRequest struct {
	Query    string `json:"q"`
	Category string `json:"category"`
	Page     int    `json:"page"`
	PageSize int    `json:"page_size"`
}
