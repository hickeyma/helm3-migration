package main

import (
	"fmt"
	"io"

	"github.com/spf13/cobra"

	"helm2to3/cmd/helm2to3/require"
	"helm2to3/pkg/helm2to3/v2"
	"helm2to3/pkg/helm2to3/v3"

)

const convertDesc = `
This command converts a Helm v2 release to v2 release format.
`
type convertOptions struct {
	namespace string
	releaseName    string
}

func newConvertCmd(out io.Writer) *cobra.Command {
	o := &convertOptions{}

	cmd := &cobra.Command{
		Use:   "convert RELEASE NAMESPACE",
		Short: "Converts a Helm v2 release to v2 release format",
		Long:  convertDesc,
		Args:  require.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			o.releaseName = args[0]
			o.namespace = args[1]
			return o.run(out)
		},
	}

	return cmd
}

func (o *convertOptions) run(out io.Writer) error {
	fmt.Fprintf(out, "Convert %s\n", o.releaseName)

	fmt.Printf("Get v2 release info ....\n")
        v2Rel, err := v2.GetRelease(o.releaseName)
        if err != nil {
                return err
        }

        fmt.Printf("Map v2 release info to equivalent v3 info....\n")
        v3Chrt, err := v3.Mapv2ChartTov3Chart(v2Rel.Chart)
        if err != nil {
                return err
        }

        fmt.Printf("Add v2 release info to v3 state ... \n")
        cfg := v3.SetupConfig(o.namespace)
        client := v3.GetInstallClient(cfg)

        client.Namespace = o.namespace
        client.ReleaseName = o.releaseName

        _, err = client.Run(v3Chrt)
        if err != nil {
                return err
        }
        fmt.Printf("Migrated v2 info to v3\n")

	return nil
}
