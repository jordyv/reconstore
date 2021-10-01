package cmd

import "github.com/integrii/flaggy"

var (
	subdomainsCmd      *flaggy.Subcommand
	saveSubdomainsCmd  *flaggy.Subcommand
	querySubdomainsCmd *flaggy.Subcommand
	tagSubdomainsCmd   *flaggy.Subcommand

	saveProgramSlug string

	queryProgramSlug, queryDomainLike, queryTag string
	queryHasBounties, queryPrivate              bool

	tags []string
)

func init() {
	subdomainsCmd = flaggy.NewSubcommand("subdomains")

	saveSubdomainsCmd = flaggy.NewSubcommand("save")
	saveSubdomainsCmd.String(&saveProgramSlug, "p", "program", "Program slug")
	subdomainsCmd.AttachSubcommand(saveSubdomainsCmd, 1)

	querySubdomainsCmd = flaggy.NewSubcommand("query")
	querySubdomainsCmd.String(&queryProgramSlug, "s", "slug", "Program slug")
	querySubdomainsCmd.String(&queryDomainLike, "p", "pattern", "Domain pattern")
	querySubdomainsCmd.String(&queryTag, "t", "tag", "Query by tag")
	querySubdomainsCmd.Bool(&queryHasBounties, "b", "bounties", "Belongs to a paying program")
	querySubdomainsCmd.Bool(&queryPrivate, "z", "private", "Belongs to a private program")
	subdomainsCmd.AttachSubcommand(querySubdomainsCmd, 1)

	tagSubdomainsCmd = flaggy.NewSubcommand("tag")
	tagSubdomainsCmd.StringSlice(&tags, "t", "tags", "Tags to add to subdomain")
	subdomainsCmd.AttachSubcommand(tagSubdomainsCmd, 1)

	flaggy.AttachSubcommand(subdomainsCmd, 1)
}
