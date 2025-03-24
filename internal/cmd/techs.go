package cmd

import (
	"fmt"
	"github.com/integrii/flaggy"
)

var (
	techsCmd     *flaggy.Subcommand
	listTechsCmd *flaggy.Subcommand
)

func init() {
	techsCmd = flaggy.NewSubcommand("techs")

	listTechsCmd = flaggy.NewSubcommand("list")
	listTechsCmd.Description = "List all technologies in the database"
	listTechsCmd.String(&queryTech, "e", "tech", "Query by name")
	techsCmd.AttachSubcommand(listTechsCmd, 1)

	flaggy.AttachSubcommand(techsCmd, 1)
}

type TechsCmd struct{}

func (c *TechsCmd) ShouldHandle() bool {
	return techsCmd.Used
}

func (c *TechsCmd) Handle() {
	fmt.Print("techs\n\n  Usage:\n    techs list\n\n")
}
