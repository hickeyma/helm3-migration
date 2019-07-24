package v2

import (
	"k8s.io/helm/pkg/helm"
        "k8s.io/helm/pkg/proto/hapi/release"
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
