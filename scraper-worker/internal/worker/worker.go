package worker

import (
	"log"

	"strings"
	"time"

	"github.com/gocolly/colly/v2"
	"github.com/suhas-developer07/webScraperContentAnalysis/scraper-worker/internal/models.go"
)

type Scraper struct {
}

func NewScraper() *Scraper {
	return &Scraper{}
}

func (s *Scraper) ProcessTask(task models.Task) models.ScrapedResult {
	result := models.ScrapedResult{
		JobID:  task.JobID,
		TaskID: task.JobID,
		URL:    task.URL,
	}

	var textParts []string

	c := colly.NewCollector(
		colly.UserAgent("ScrapperBOT/1.0 (+http://localhost:8080)"),
		colly.AllowURLRevisit(),
	)

	c.Limit(&colly.LimitRule{
		DomainGlob:  "*",
		Parallelism: 1,
		RandomDelay: 2 * time.Second,
	})

	c.OnHTML("body", func(h *colly.HTMLElement) {

		text := h.Text

		if text != "" {
			textParts = append(textParts, strings.TrimSpace(text))
		}
	})

	c.OnError(func(r *colly.Response, err error) {
		result.Error = err.Error()
		log.Printf("Request URL :%s failed with response :%v error :%v", r.Request.URL, r, err)
	})

	err := c.Visit(task.URL)
	if err != nil {
		result.Error = err.Error()
		return result
	}

	c.Wait()

	fullText := strings.Join(textParts, "\n")

	cleanedtext := cleanText(fullText)

	result.RawText = cleanedtext

	return result
}

func cleanText(text string) string {
	text = strings.Join(strings.Fields(text), "")
	return text
}
