package cmd

import (
	"github.com/integrii/flaggy"
)

var (
	programsCmd       *flaggy.Subcommand
	addProgramsCmd    *flaggy.Subcommand
	listProgramsCmd   *flaggy.Subcommand
	deleteProgramsCmd *flaggy.Subcommand

	programName, slug, platform string
	private, hasBounties        bool

	listOrder string
)

func init() {
	programsCmd = flaggy.NewSubcommand("programs")

	addProgramsCmd = flaggy.NewSubcommand("add")
	addProgramsCmd.Description = "Add a new program"
	addProgramsCmd.String(&programName, "n", "name", "Name")
	addProgramsCmd.String(&slug, "s", "slug", "Slug")
	addProgramsCmd.String(&platform, "p", "platform", "Bug bounty platform")
	addProgramsCmd.Bool(&private, "z", "private", "Private program")
	addProgramsCmd.Bool(&hasBounties, "b", "bounties", "Has bounty payments")
	programsCmd.AttachSubcommand(addProgramsCmd, 1)

	listProgramsCmd = flaggy.NewSubcommand("list")
	listProgramsCmd.Description = "List all programs in the database"
	listProgramsCmd.String(&listOrder, "s", "sort", "Sort by field (id, name, slug, private, has_bounties), default by name")
	programsCmd.AttachSubcommand(listProgramsCmd, 1)

	deleteProgramsCmd = flaggy.NewSubcommand("delete")
	deleteProgramsCmd.Description = "Delete a program"
	deleteProgramsCmd.String(&slug, "s", "slug", "Slug of the program")
	programsCmd.AttachSubcommand(deleteProgramsCmd, 1)

	flaggy.AttachSubcommand(programsCmd, 1)
}
