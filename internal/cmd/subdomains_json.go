package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/jordyv/reconstore/internal/entities"
	"github.com/sirupsen/logrus"
)

type JsonSubdomainsCmd struct{}

func (c *JsonSubdomainsCmd) ShouldHandle() bool {
	return jsonSubdomainsCmd.Used
}

func (c *JsonSubdomainsCmd) Handle() {
	d := db.Model(&entities.Subdomain{}).Joins("Program").Joins("HTTPInfo").Joins("DNSInfo")

	applyQueryFilters(d)

	results := make([]entities.Subdomain, 0)
	r, err := d.Rows()
	if err != nil {
		logrus.Error(err.Error())
	}

	defer r.Close()
	for r.Next() {
		var s entities.Subdomain
		db.ScanRows(r, &s)
		results = append(results, s)
	}

	b, err := json.Marshal(results)
	if err != nil {
		logrus.WithError(err).Fatal("could not encode result set to JSON")
	}

	fmt.Println(string(b))
}
