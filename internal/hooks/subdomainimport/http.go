package subdomainimport

import (
	"fmt"
	wappalyzer "github.com/projectdiscovery/wappalyzergo"
	"github.com/sirupsen/logrus"
	"strings"
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

func (h *HTTP) handleScraperResponse(w *wappalyzer.Wappalyze, httpInfo *entities.HTTPInfo, subdomain *entities.Subdomain) colly.ResponseCallback {
	return func(r *colly.Response) {
		httpInfo.HTTPSStatusCode = r.StatusCode
		if s := r.Headers.Get("Server"); s != "" {
			httpInfo.WebServer = s
		}
		if c := r.Headers.Get("Content-Type"); c != "" {
			httpInfo.ContentType = c
		}
		if h := r.Headers; h != nil {
			httpInfo.AllHeaders = make(map[string]string)
			for k, v := range *h {
				httpInfo.AllHeaders[k] = strings.Join(v, ",")
			}
		}

		if w != nil && r.Headers != nil {
			info := w.Fingerprint(*r.Headers, r.Body)
			for k := range info {
				t := entities.Tech{Name: k}
				h.db.Where("name = ?", k).FirstOrCreate(&t)
				subdomain.Techs = append(subdomain.Techs, &t)
			}
		}
	}
}

func (h *HTTP) AfterSave(s *entities.Subdomain) error {
	httpInfo := entities.HTTPInfo{}

	httpsAddr := fmt.Sprintf("https://%s", s.Domain)
	httpAddr := fmt.Sprintf("http://%s", s.Domain)

	w, err := wappalyzer.New()
	if err != nil {
		logrus.WithError(err).Error("could not initialize wappalyzer")
	}

	httpsScraper := colly.NewCollector(collyConfig...)
	httpsScraper.SetRequestTimeout(defaultTimeout)
	httpsScraper.OnResponse(h.handleScraperResponse(w, &httpInfo, s))
	httpsScraper.OnHTML("title", func(e *colly.HTMLElement) {
		httpInfo.Title = e.Text
	})
	err = httpsScraper.Visit(httpsAddr)

	if err == nil {
		httpScraper := colly.NewCollector(collyConfig...)
		httpScraper.SetRequestTimeout(defaultTimeout)
		httpScraper.OnResponse(h.handleScraperResponse(w, &httpInfo, s))
		httpsScraper.OnHTML("title", func(e *colly.HTMLElement) {
			if httpInfo.Title == "" {
				httpInfo.Title = e.Text
			}
		})
		httpScraper.Visit(httpAddr)
	}

	s.HTTPInfo = &httpInfo
	h.db.Save(s)

	return nil
}
