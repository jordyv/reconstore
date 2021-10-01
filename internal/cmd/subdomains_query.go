package cmd

import (
	"fmt"

	"github.com/jordyv/reconstore/internal/entities"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type QuerySubdomainsCmd struct{}

func (c *QuerySubdomainsCmd) ShouldHandle() bool {
	return querySubdomainsCmd.Used
}

func (c *QuerySubdomainsCmd) Handle() {
	d := db.Model(&entities.Subdomain{}).Joins("Program")
	if queryProgramSlug != "" {
		var program entities.Program
		p := db.Where("slug = ?", queryProgramSlug).First(&program)
		if p.Error == gorm.ErrRecordNotFound {
			logrus.Fatal("No program found with this slug")
		}
		if program.ID > 0 {
			d.Where("program_id = ?", program.ID)
		}
	}

	if queryDomainLike != "" {
		d.Where("domain LIKE ?", fmt.Sprintf("%%%s%%", queryDomainLike))
	}

	if queryHasBounties {
		d.Where("program.has_bounties = true")
	}

	if queryPrivate {
		d.Where("program.private = true")
	}

	if queryTag != "" {
		d.Joins("LEFT JOIN subdomain_tags ON subdomains.id = subdomain_tags.subdomain_id LEFT JOIN tags ON tags.id = subdomain_tags.tag_id").Where("tags.name = ?", queryTag)
	}

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
