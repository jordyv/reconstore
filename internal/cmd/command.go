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
		&SaveSubdomainsCmd{},
		&QuerySubdomainsCmd{},
		&TagSubdomainsCmd{},
		&JsonCmd{},
	}

	queryProgramSlug, queryDomainLike, queryTag string
	queryHasBounties, queryPrivate              bool
)

func Init(d *gorm.DB) {
	db = d
}
