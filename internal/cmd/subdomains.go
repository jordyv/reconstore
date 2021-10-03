package cmd

import (
	"fmt"

	"github.com/integrii/flaggy"
	"github.com/jordyv/reconstore/internal/entities"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

var (
	subdomainsCmd      *flaggy.Subcommand
	saveSubdomainsCmd  *flaggy.Subcommand
	querySubdomainsCmd *flaggy.Subcommand
	tagSubdomainsCmd   *flaggy.Subcommand
	jsonSubdomainsCmd  *flaggy.Subcommand

	saveProgramSlug string

	queryProgramSlug, queryDomainLike, queryTag string
	queryHasBounties, queryPrivate              bool

	tags []string
)

func init() {
	subdomainsCmd = flaggy.NewSubcommand("subdomains")

	saveSubdomainsCmd = flaggy.NewSubcommand("save")
	saveSubdomainsCmd.String(&saveProgramSlug, "p", "program", "Program slug")
	subdomainsCmd.AttachSubcommand(saveSubdomainsCmd, 1)

	querySubdomainsCmd = flaggy.NewSubcommand("query")
	querySubdomainsCmd.String(&queryProgramSlug, "s", "slug", "Program slug")
	querySubdomainsCmd.String(&queryDomainLike, "p", "pattern", "Domain pattern")
	querySubdomainsCmd.String(&queryTag, "t", "tag", "Query by tag")
	querySubdomainsCmd.Bool(&queryHasBounties, "b", "bounties", "Belongs to a paying program")
	querySubdomainsCmd.Bool(&queryPrivate, "z", "private", "Belongs to a private program")
	subdomainsCmd.AttachSubcommand(querySubdomainsCmd, 1)

	tagSubdomainsCmd = flaggy.NewSubcommand("tag")
	tagSubdomainsCmd.StringSlice(&tags, "t", "tags", "Tags to add to subdomain")
	subdomainsCmd.AttachSubcommand(tagSubdomainsCmd, 1)

	jsonSubdomainsCmd = flaggy.NewSubcommand("json")
	jsonSubdomainsCmd.String(&queryProgramSlug, "s", "slug", "Program slug")
	jsonSubdomainsCmd.String(&queryDomainLike, "p", "pattern", "Domain pattern")
	jsonSubdomainsCmd.String(&queryTag, "t", "tag", "Query by tag")
	jsonSubdomainsCmd.Bool(&queryHasBounties, "b", "bounties", "Belongs to a paying program")
	jsonSubdomainsCmd.Bool(&queryPrivate, "z", "private", "Belongs to a private program")
	subdomainsCmd.AttachSubcommand(jsonSubdomainsCmd, 1)

	flaggy.AttachSubcommand(subdomainsCmd, 1)
}

func applyQueryFilters(d *gorm.DB) {
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
}
