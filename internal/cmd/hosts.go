package cmd

import (
	"fmt"
	"github.com/integrii/flaggy"
)

var (
	hostsCmd      *flaggy.Subcommand
	queryHostsCmd *flaggy.Subcommand
)

func init() {
	hostsCmd = flaggy.NewSubcommand("hosts")

	queryHostsCmd = flaggy.NewSubcommand("query")
	queryHostsCmd.String(&queryProgramSlug, "s", "slug", "Program slug")
	queryHostsCmd.String(&queryDomainLike, "p", "pattern", "Domain pattern")
	queryHostsCmd.String(&queryTag, "t", "tag", "Query by tag")
	queryHostsCmd.String(&queryTech, "e", "tech", "Query by tech")
	queryHostsCmd.Bool(&queryHasBounties, "b", "bounties", "Belongs to a paying program")
	queryHostsCmd.Bool(&queryPrivate, "z", "private", "Belongs to a private program")
	hostsCmd.AttachSubcommand(queryHostsCmd, 1)

	flaggy.AttachSubcommand(hostsCmd, 1)
}

type HostsCmd struct{}

func (c *HostsCmd) ShouldHandle() bool {
	return hostsCmd.Used
}

func (c *HostsCmd) Handle() {
	fmt.Print("hosts\n\n  Usage:\n    hosts query\n\n")
}
