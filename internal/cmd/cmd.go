package cmd

import (
	"gcu/internal/cmd/run"
	"gcu/internal/cmd/version"

	"github.com/spf13/cobra"
)

const (
	descriptionShort = `TODO`
	descriptionLong  = `
	TODO`
)

func NewRootCommand(name string) *cobra.Command {
	c := &cobra.Command{
		Use:   name,
		Short: descriptionShort,
		Long:  descriptionLong,
	}

	c.AddCommand(
		version.NewCommand(),
		run.NewCommand(),
	)

	return c
}
