package subdomainimport

import (
	"fmt"
	"time"

	"github.com/gocolly/colly"
	"github.com/jordyv/reconstore/internal/entities"
	"gorm.io/gorm"
)

var (
	defaultTimeout = time.Second * 3

	collyConfig = []func(*colly.Collector){
		colly.ParseHTTPErrorResponse(),
	}
)

type HTTP struct {
	db      *gorm.DB
	scraper *colly.Collector
}

func NewHTTP(db *gorm.DB) *HTTP {
	return &HTTP{
		db:      db,
		scraper: colly.NewCollector(),
	}
}

func (h *HTTP) AfterSave(s *entities.Subdomain) error {
	httpInfo := entities.HTTPInfo{}

	httpsAddr := fmt.Sprintf("https://%s", s.Domain)
	httpAddr := fmt.Sprintf("http://%s", s.Domain)

	httpsScraper := colly.NewCollector(collyConfig...)
	httpsScraper.SetRequestTimeout(defaultTimeout)
	httpsScraper.OnResponse(func(r *colly.Response) {
		httpInfo.HTTPSStatusCode = r.StatusCode
		if s := r.Headers.Get("Server"); s != "" {
			httpInfo.WebServer = s
		}
		if c := r.Headers.Get("Content-Type"); c != "" {
			httpInfo.ContentType = c
		}
	})
	httpsScraper.OnHTML("title", func(e *colly.HTMLElement) {
		httpInfo.Title = e.Text
	})
	httpsScraper.Visit(httpsAddr)

	httpScraper := colly.NewCollector(collyConfig...)
	httpScraper.SetRequestTimeout(defaultTimeout)
	httpScraper.OnResponse(func(r *colly.Response) {
		httpInfo.HTTPSStatusCode = r.StatusCode
		if s := r.Headers.Get("Server"); httpInfo.WebServer == "" && s != "" {
			httpInfo.WebServer = s
		}
		if c := r.Headers.Get("Content-Type"); httpInfo.ContentType == "" && c != "" {
			httpInfo.ContentType = c
		}
	})
	httpsScraper.OnHTML("title", func(e *colly.HTMLElement) {
		if httpInfo.Title == "" {
			httpInfo.Title = e.Text
		}
	})
	httpScraper.Visit(httpAddr)

	s.HTTPInfo = httpInfo
	h.db.Save(s)

	return nil
}
