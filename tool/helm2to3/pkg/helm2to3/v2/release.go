package v2

import (
	"k8s.io/helm/pkg/helm"
        "k8s.io/helm/pkg/proto/hapi/release"
	rls "k8s.io/helm/pkg/proto/hapi/services"
)

func GetRelease(releaseName string) (*release.Release, error) {
	helmClient, err := GetHelmClient()
	if err != nil {
	        return nil, err
        }

	res, err := helmClient.ReleaseContent(releaseName, helm.ContentReleaseVersion(1))
        if err != nil {
	        return nil, err
        }

	return res.Release, nil
}

func GetReleaseVersions(releaseName string) ([]*release.Release, error) {
	helmClient, err := GetHelmClient()
        if err != nil {
                return nil, err
        }

	r, err := helmClient.ReleaseHistory(releaseName, helm.WithMaxHistory(histMax))
        if err != nil {
                return nil, err
        }

	return r.Releases, nil
}

func DeleteRelease(releaseName string) (*rls.UninstallReleaseResponse, error) {
        opts := []helm.DeleteOption{
		 helm.DeleteDryRun(false),
                helm.DeleteDisableHooks(false),
                helm.DeletePurge(true),
                helm.DeleteTimeout(300),
                helm.DeleteDescription(""),
        }

	helmClient, err := GetHelmClient()
        if err != nil {
                return nil, err
        }

        resp, err := helmClient.DeleteRelease(releaseName, opts...)
        if err != nil {
                return nil, err
        }

        return resp, nil
}
