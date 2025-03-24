package cmd

import (
	"fmt"
	"github.com/jordyv/reconstore/internal/entities"
	"github.com/rodaine/table"
	"github.com/sirupsen/logrus"
)

type ListTechsCmd struct{}

func (c *ListTechsCmd) ShouldHandle() bool {
	return listTechsCmd.Used
}

func (c *ListTechsCmd) Handle() {
	d := db.Model(&entities.Tech{}).Order("name")

	if queryTech != "" {
		d.Where("name LIKE ?", fmt.Sprintf("%%%s%%", queryTech))
	}

	r, err := d.Rows()
	if err != nil {
		logrus.Fatal(err)
	}
	defer r.Close()

	tbl := table.New("ID", "Name")

	for r.Next() {
		var t entities.Tech
		db.ScanRows(r, &t)
		tbl.AddRow(t.ID, t.Name)
	}

	tbl.Print()
}
