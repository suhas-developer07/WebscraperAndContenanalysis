package search

import "strings"

type SearchRequest struct {
	Query    string `json:"q"`
	Category string `json:"category"`
	Page     int    `json:"page"`
	PageSize int    `json:"page_size"`
}

func BuildQuery(req SearchRequest) map[string]interface{} {
	from := (req.Page - 1) * req.PageSize
	
	query := map[string]interface{}{
		"from": from,
		"size": req.PageSize,
		"sort": []map[string]interface{}{
			{"_score": map[string]interface{}{"order": "desc"}},
		},
	}

	boolQuery := map[string]interface{}{}

	if req.Query != "" {
		boolQuery["must"] = []map[string]interface{}{
			{
				"multi_match": map[string]interface{}{
					"query":    req.Query,
					"fields":   []string{"summary", "key_entities"},
					"fuzziness": "AUTO",
				},
			},
		}
	}

	if req.Category != "" {
		boolQuery["filter"] = []map[string]interface{}{
			{
				"term": map[string]interface{}{
					"domain_category": strings.ToLower(req.Category),
				},
			},
		}
	}

	if len(boolQuery) > 0 {
		query["query"] = map[string]interface{}{
			"bool": boolQuery,
		}
	} else {
		query["query"] = map[string]interface{}{
			"match_all": map[string]interface{}{},
		}
	}

	return query
}