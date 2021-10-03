package cmd

import (
	"encoding/json"
	"os"

	"github.com/integrii/flaggy"
	"github.com/jordyv/reconstore/internal/entities"
	"github.com/sirupsen/logrus"
)

var (
	jsonCmd *flaggy.Subcommand
)

func init() {
	jsonCmd = flaggy.NewSubcommand("json")
	jsonCmd.String(&queryProgramSlug, "s", "slug", "Program slug")
	jsonCmd.String(&queryDomainLike, "p", "pattern", "Domain pattern")
	jsonCmd.String(&queryTag, "t", "tag", "Query by tag")
	jsonCmd.Bool(&queryHasBounties, "b", "bounties", "Belongs to a paying program")
	jsonCmd.Bool(&queryPrivate, "z", "private", "Belongs to a private program")
	flaggy.AttachSubcommand(jsonCmd, 1)
}

type JsonCmd struct{}

func (c *JsonCmd) ShouldHandle() bool {
	return jsonCmd.Used
}

func (c *JsonCmd) Handle() {
	d := db.Model(&entities.Subdomain{}).Joins("Program").Joins("HTTPInfo").Joins("DNSInfo")

	applySubdomainQueryFilters(d)

	results := make([]entities.Subdomain, 0)
	r, err := d.Rows()
	if err != nil {
		logrus.Error(err.Error())
	}

	defer r.Close()
	for r.Next() {
		var s entities.Subdomain
		db.ScanRows(r, &s)
		results = append(results, s)
	}

	j := json.NewEncoder(os.Stdout)
	if err := j.Encode(results); err != nil {
		logrus.WithError(err).Fatal("could not encode result set to JSON")
	}
}
