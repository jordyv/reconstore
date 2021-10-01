package cmd

import (
	"github.com/jordyv/reconstore/internal/entities"
	"github.com/sirupsen/logrus"
)

type DeleteProgramsCmd struct{}

func (c *DeleteProgramsCmd) ShouldHandle() bool {
	return deleteProgramsCmd.Used
}

func (c *DeleteProgramsCmd) Handle() {
	c.validate()

	var program entities.Program
	db.Where("slug = ?", slug).Find(&program)

	if program.ID > 0 {
		if r := db.Where("program_id = ?", program.ID).Delete(&entities.Subdomain{}); r.Error != nil {
			logrus.WithError(r.Error).Warn("couldn't delete all subdomains for program")
			return
		}
		if r := db.Where("slug = ?", slug).Delete(&entities.Program{}); r.Error != nil {
			logrus.WithError(r.Error).Fatalf("couldn't delete program with slug '%s'", slug)
		}
	} else {
		logrus.Fatal("No program with this slug found")
	}
}

func (c *DeleteProgramsCmd) validate() {
	if slug == "" {
		logrus.Fatal("Slug is required")
	}
}
