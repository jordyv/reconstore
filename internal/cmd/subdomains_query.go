package cmd

import (
	"fmt"

	"github.com/jordyv/reconstore/internal/entities"
	"github.com/sirupsen/logrus"
)

type QuerySubdomainsCmd struct{}

func (c *QuerySubdomainsCmd) ShouldHandle() bool {
	return querySubdomainsCmd.Used
}

func (c *QuerySubdomainsCmd) Handle() {
	d := db.Model(&entities.Subdomain{}).Joins("Program")

	applyQueryFilters(d)

	r, err := d.Rows()
	if err != nil {
		logrus.Error(err.Error())
	}

	defer r.Close()
	for r.Next() {
		var s entities.Subdomain
		db.ScanRows(r, &s)
		fmt.Printf("%s\n", s.Domain)
	}
}
