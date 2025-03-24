package cmd

import (
	"fmt"
	"github.com/jordyv/reconstore/internal/entities"
	"github.com/sirupsen/logrus"
)

type QueryHostsCmd struct{}

func (c *QueryHostsCmd) ShouldHandle() bool {
	return queryHostsCmd.Used
}

func (c *QueryHostsCmd) Handle() {
	d := db.Model(&entities.Subdomain{}).
		Select("subdomains.*, http_infos.*").
		Joins("Program").
		Joins("JOIN http_infos on http_infos.id = subdomains.http_info_id").
		Where("subdomains.http_info_id IS NOT NULL")

	applySubdomainQueryFilters(d)

	r, err := d.Rows()
	if err != nil {
		logrus.Error(err.Error())
	}

	defer r.Close()
	for r.Next() {
		var s struct {
			Domain          string
			HTTPStatusCode  int
			HTTPSStatusCode int
		}
		db.ScanRows(r, &s)
		if s.HTTPStatusCode > 0 || s.HTTPSStatusCode > 0 {
			protocol := "http"
			statusCode := s.HTTPStatusCode
			if s.HTTPSStatusCode > 0 {
				protocol = "https"
				statusCode = s.HTTPSStatusCode
			}
			fmt.Printf("%d %s://%s\n", statusCode, protocol, s.Domain)
		}
	}
}
