package cmd

import (
	"fmt"
	"github.com/jordyv/reconstore/internal/entities"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

func applySubdomainQueryFilters(d *gorm.DB) {
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

	if queryTech != "" {
		d.Joins("LEFT JOIN subdomain_techs ON subdomains.id = subdomain_techs.subdomain_id LEFT JOIN techs ON techs.id = subdomain_techs.tech_id").Where("techs.name LIKE ?", fmt.Sprintf("%%%s%%", queryTech))
	}
}
