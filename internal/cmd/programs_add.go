package cmd

import (
	"github.com/gobeam/stringy"
	"github.com/jordyv/reconstore/internal/entities"
	"github.com/sirupsen/logrus"
)

type AddProgramsCmd struct{}

func (c *AddProgramsCmd) ShouldHandle() bool {
	return addProgramsCmd.Used
}

func (c *AddProgramsCmd) Handle() {
	c.validate()

	if slug == "" {
		slug = stringy.New(programName).SnakeCase().ToLower()
	}

	var count int64
	db.Model(&entities.Program{}).Where("slug = ?", slug).Count(&count)
	if count > 0 {
		logrus.Fatalf("There is already a program with slug '%s'", slug)
	}

	db.Create(&entities.Program{
		Name:        programName,
		Slug:        slug,
		Platform:    platform,
		Private:     private,
		HasBounties: hasBounties,
	})
}

func (c *AddProgramsCmd) validate() {
	if programName == "" || platform == "" {
		logrus.Fatal("Program name and platform are required")
	}
}
