package cmd

import "gorm.io/gorm"

type Cmd interface {
	ShouldHandle() bool
	Handle()
}

var (
	db *gorm.DB

	Commands = []Cmd{
		&AddProgramsCmd{},
		&ListProgramsCmd{},
		&DeleteProgramsCmd{},
		&ProgramsCmd{},
		&SaveSubdomainsCmd{},
		&QuerySubdomainsCmd{},
		&TagSubdomainsCmd{},
		&SubdomainsCmd{},
		&JsonCmd{},
		&ListTechsCmd{},
		&TechsCmd{},
	}

	queryProgramSlug, queryDomainLike, queryTag, queryTech string
	queryHasBounties, queryPrivate                         bool
)

func Init(d *gorm.DB) {
	db = d
}
