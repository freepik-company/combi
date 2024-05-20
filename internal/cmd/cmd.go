package cmd

import (
	"combi/internal/cmd/daemon"
	"combi/internal/cmd/version"

	"github.com/spf13/cobra"
)

const (
	descriptionShort = `TODO`
	descriptionLong  = `
	TODO`
)

func NewRootCommand(name string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   name,
		Short: descriptionShort,
		Long:  descriptionLong,
	}

	cmd.AddCommand(
		version.NewCommand(),
		daemon.NewCommand(),
	)

	return cmd
}
