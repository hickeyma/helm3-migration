package main

import (
	"fmt"
	"io"

	"github.com/spf13/cobra"

	"helm2to3/cmd/helm2to3/require"
	"helm2to3/pkg/helm2to3/v2"
	"helm2to3/pkg/helm2to3/v3"
	v2Rel "k8s.io/helm/pkg/proto/hapi/release"

)

const convertDesc = `
This command converts a Helm v2 release to v2 release format.
`
type convertOptions struct {
	releaseName    string
}

func newConvertCmd(out io.Writer) *cobra.Command {
	o := &convertOptions{}

	cmd := &cobra.Command{
		Use:   "convert RELEASE",
		Short: "Converts a Helm v2 release to v3 release",
		Long:  convertDesc,
		Args:  require.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			o.releaseName = args[0]
			return o.run(out)
		},
	}

	return cmd
}

func (o *convertOptions) run(out io.Writer) error {
	fmt.Printf("[Helm 3] Release \"%s\" will be created.\n", o.releaseName)

	//fmt.Printf("Get v2 release versions ....\n")
        v2Releases, err := v2.GetReleaseVersions(o.releaseName)
        if err != nil {
                return err
        }

	var isUpgrade = make(map[string]bool)
	for i := len(v2Releases) - 1; i >= 0; i-- {
		release := v2Releases[i]
		if _, ok := isUpgrade[release.Namespace]; !ok {
			isUpgrade[release.Namespace] = false
		}
		fmt.Printf("[Helm 3] ReleaseVersion \"%s\" will be created.\n", fmt.Sprintf("%s.v%d", release.Name, release.Version))
		//err := createReleaseVersion(release, o.releaseName, isUpgrade[release.Namespace])
		err := createReleaseVersion(release, o.releaseName, isUpgrade[release.Namespace])
                if err != nil {
                        return err
	        }
		isUpgrade[release.Namespace] = true
		fmt.Printf("[Helm 3] ReleaseVersion \"%s\" created.\n", fmt.Sprintf("%s.v%d", release.Name, release.Version))
        }

	fmt.Printf("[Helm 3] Release \"%s\"  created.\n", o.releaseName)

	fmt.Printf("[Helm 2] Release \"%s\" will be deleted.\n", o.releaseName)

	_, err = v2.DeleteRelease(o.releaseName)
        if err != nil {
                return err
        }

	fmt.Printf("[Helm 2] Release \"%s\"  deleted.\n", o.releaseName)

	return nil
}

func createReleaseVersion(release *v2Rel.Release, releaseName string, isUpgrade bool) error {
        //fmt.Printf("Map v2 chart to equivalent v3 chart....\n")
        v3Chrt, err := v3.Mapv2ChartTov3Chart(release.Chart)
        if err != nil {
                return err
        }

        //fmt.Printf("Add v2 release version to v3 ... \n")
	if isUpgrade {
		return v3.UpgradeRelease(v3Chrt, releaseName, release.Namespace)
	} else {
		return v3.InstallRelease(v3Chrt, releaseName, release.Namespace)
	}
}

