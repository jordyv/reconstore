package cmd

import (
	"fmt"
	"github.com/integrii/flaggy"
	"github.com/jordyv/reconstore/internal/config"
)

var (
	infoCmd *flaggy.Subcommand
)

func init() {
	infoCmd = flaggy.NewSubcommand("info")

	flaggy.AttachSubcommand(infoCmd, 1)
}

type InfoCmd struct{}

func (c *InfoCmd) ShouldHandle() bool {
	return infoCmd.Used
}

func (c *InfoCmd) Handle() {
	fmt.Printf(
		"info\n\ndatabase: %s (%s)\n\n",
		config.GetString(config.DBConnectionString),
		config.GetString(config.DBType),
	)
}
