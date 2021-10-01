package cmd

import (
	"github.com/jordyv/reconstore/internal/entities"
	"github.com/rodaine/table"
	"github.com/sirupsen/logrus"
)

type ListProgramsCmd struct{}

func (c *ListProgramsCmd) ShouldHandle() bool {
	return listProgramsCmd.Used
}

func (c *ListProgramsCmd) Handle() {
	if listOrder == "" {
		listOrder = "name"
	}
	r, err := db.Model(&entities.Program{}).Order(listOrder).Rows()
	if err != nil {
		logrus.Fatal(err)
	}
	defer r.Close()

	tbl := table.New("ID", "Name", "Slug", "Platform", "Private", "Has bounties")

	for r.Next() {
		var p entities.Program
		db.ScanRows(r, &p)
		privateString := "no"
		if p.Private {
			privateString = "yes"
		}
		bountiesString := "no"
		if p.HasBounties {
			bountiesString = "yes"
		}
		tbl.AddRow(p.ID, p.Name, p.Slug, p.Platform, privateString, bountiesString)
	}

	tbl.Print()
}
