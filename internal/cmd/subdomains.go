package cmd

import (
	"fmt"
	"github.com/integrii/flaggy"
)

var (
	subdomainsCmd      *flaggy.Subcommand
	saveSubdomainsCmd  *flaggy.Subcommand
	querySubdomainsCmd *flaggy.Subcommand
	tagSubdomainsCmd   *flaggy.Subcommand

	saveThreads     int = 12
	saveProgramSlug string
	update          bool

	tags []string
)

func init() {
	subdomainsCmd = flaggy.NewSubcommand("subdomains")

	saveSubdomainsCmd = flaggy.NewSubcommand("save")
	saveSubdomainsCmd.String(&saveProgramSlug, "p", "program", "Program slug")
	saveSubdomainsCmd.Bool(&update, "u", "update", "Update if already stored")
	saveSubdomainsCmd.Int(&saveThreads, "t", "threads", "Number of threads")
	subdomainsCmd.AttachSubcommand(saveSubdomainsCmd, 1)

	querySubdomainsCmd = flaggy.NewSubcommand("query")
	querySubdomainsCmd.String(&queryProgramSlug, "s", "slug", "Program slug")
	querySubdomainsCmd.String(&queryDomainLike, "p", "pattern", "Domain pattern")
	querySubdomainsCmd.String(&queryTag, "t", "tag", "Query by tag")
	querySubdomainsCmd.String(&queryTech, "e", "tech", "Query by tech")
	querySubdomainsCmd.Bool(&queryHasBounties, "b", "bounties", "Belongs to a paying program")
	querySubdomainsCmd.Bool(&queryPrivate, "z", "private", "Belongs to a private program")
	subdomainsCmd.AttachSubcommand(querySubdomainsCmd, 1)

	tagSubdomainsCmd = flaggy.NewSubcommand("tag")
	tagSubdomainsCmd.StringSlice(&tags, "t", "tags", "Tags to add to subdomain")
	subdomainsCmd.AttachSubcommand(tagSubdomainsCmd, 1)

	flaggy.AttachSubcommand(subdomainsCmd, 1)
}

type SubdomainsCmd struct{}

func (c *SubdomainsCmd) ShouldHandle() bool {
	return subdomainsCmd.Used
}

func (c *SubdomainsCmd) Handle() {
	fmt.Print("subdomains\n\n  Usage:\n    subdomains [query|save|tag]\n\n")
}
