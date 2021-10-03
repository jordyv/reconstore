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
		&JsonSubdomainsCmd{},
	}
)

func Init(d *gorm.DB) {
	db = d
}
