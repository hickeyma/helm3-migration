package main

import (
	"io"

	"github.com/spf13/cobra"

	"helm2to3/cmd/helm2to3/require"
)

func newRootCmd(out io.Writer, args []string) *cobra.Command {
	cmd := &cobra.Command{
		Use:                    "helm2to3",
		Short:                  "The Helm 2 to Helm 3 Migration Tool",
		Long:                   "The Helm 2 to Helm 3 Migration Tool",
		SilenceUsage:           true,
		Args:                   require.NoArgs,
	}
	flags := cmd.PersistentFlags()
	flags.Parse(args)

	cmd.AddCommand(
		newConvertCmd(out),
	)

	return cmd
}
