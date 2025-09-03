package models

type Job struct {
	ID   int64    `json:"id"`
	URLs []string `json:"urls"`
}

type Task struct {
	JobID  int64  `json:"job_id"`
	TaskID int    `json:"task_id"`
	URL    string `json:"url"`
}

type ScrapedResult struct {
	JobID   int64  `json:"job_id"`
	TaskID  int64  `json:"task_id"`
	URL     string `json:"url"`
	RawText string `json:"raw_text"`
	Error   string `json:"error,omitempty"`
}
