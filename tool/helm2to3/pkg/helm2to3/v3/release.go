package v3
  
import (
	"github.com/golang/protobuf/ptypes/any"

        "helm.sh/helm/pkg/chart"

	v2ChartUtil "k8s.io/helm/pkg/chartutil"
	v2Chart "k8s.io/helm/pkg/proto/hapi/chart"
)


func InstallRelease(v3Chrt *chart.Chart, releaseName string, namespace string) error {
	cfg := SetupConfig(namespace)
        client := GetInstallClient(cfg)
        client.Namespace = namespace
        client.ReleaseName = releaseName
	_, err := client.Run(v3Chrt)
	return err
}

func UpgradeRelease(v3Chrt *chart.Chart, releaseName string, namespace string) error {
	cfg := SetupConfig(namespace)
        client := GetUpgradeClient(cfg)
        client.MaxHistory = 256
	_, err := client.Run(releaseName, v3Chrt)
	return err
}

func Mapv2ChartTov3Chart(v2Chrt *v2Chart.Chart) (*chart.Chart, error) {
	 v3Chrt := new(chart.Chart)

	 v3Chrt.Metadata = mapMetadata(v2Chrt)
	 v3Chrt.Templates = mapTemplates(v2Chrt.Templates)
	 err := mapDependencies(v2Chrt.Dependencies, v3Chrt)
	 if err != nil {
		 return nil, err
         }
	 v3Chrt.Values, err = mapValues(v2Chrt.Values)
	 if err != nil {
		 return nil, err
         }
	 v3Chrt.Files = mapFiles(v2Chrt.Files)
	 //TODO
	 //v3Chrt.Schema
	 //TODO
	 //v3Chrt.Lock = new(chart.Lock)

	 return v3Chrt, nil
}

func mapMetadata(v2Chrt *v2Chart.Chart) (*chart.Metadata) {
	metadata := new(chart.Metadata)

	metadata.Name = v2Chrt.Metadata.Name
        metadata.Home = v2Chrt.Metadata.Home
        metadata.Sources = v2Chrt.Metadata.Sources
        metadata.Version = v2Chrt.Metadata.Version
        metadata.Description = v2Chrt.Metadata.Description
        metadata.Keywords = v2Chrt.Metadata.Keywords
	metadata.Maintainers = mapMaintainers(v2Chrt.Metadata.Maintainers)
        metadata.Icon = v2Chrt.Metadata.Icon
        metadata.APIVersion = v2Chrt.Metadata.ApiVersion
        metadata.Condition = v2Chrt.Metadata.Condition
        metadata.Tags = v2Chrt.Metadata.Tags
        metadata.AppVersion = v2Chrt.Metadata.AppVersion
        metadata.Deprecated = v2Chrt.Metadata.Deprecated
        metadata.Annotations = v2Chrt.Metadata.Annotations
        metadata.KubeVersion  = v2Chrt.Metadata.KubeVersion
	//TODO: metadata.Dependencies = 
	metadata.Type =  "application"

	return metadata
}

func mapMaintainers(v2Maintainers []*v2Chart.Maintainer) ([]*chart.Maintainer) {
	maintainers := []*chart.Maintainer{}

	for _, val := range v2Maintainers {
		maintainer := new(chart.Maintainer)
		maintainer.Name = val.Name
		maintainer.Email = val.Email
		maintainer.URL = val.Url
		maintainers = append(maintainers, maintainer)
	}
	return maintainers
}

func mapTemplates(v2Templates []*v2Chart.Template) ([]*chart.File) {
	files := []*chart.File{}

	 for _, val := range v2Templates {
		 file := new(chart.File)
		 file.Name = val.Name
		 file.Data = val.Data
		 files = append(files, file)
	 }
	 return files
}

func mapDependencies(v2Dependencies []*v2Chart.Chart, chart *chart.Chart) error {
	for _, val := range v2Dependencies {
		dependency, err := Mapv2ChartTov3Chart(val)
		if err != nil {
			return err
		}
		chart.AddDependency(dependency)
	}
	return nil
}

func mapValues(v2Config *v2Chart.Config) (map[string]interface{}, error) {
	values, err := v2ChartUtil.ReadValues([]byte(v2Config.Raw))
        if err != nil {
		return nil, err
        }

	return values, nil
}

func mapFiles(v2Files []*any.Any) ([]*chart.File) {
	files := []*chart.File{}
	for _, f := range v2Files {
		file := new(chart.File)
		file.Name = f.TypeUrl
                file.Data = f.Value
                files = append(files, file)
	 }
	 return files

}
