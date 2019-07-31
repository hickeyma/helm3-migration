# 1. Overview

Helm v3 introduces quite a lot of change in the underlying architecture and plumbing from the previous release, 
Helm v2. One key change is around Release state storage. The changes includes the Kubernetes resource for storage and the 
release object metadata contained in the resource. Releases will also be on a per user namespace instead of the Tiller 
namespace used, like the default namespace `kube-system`.

# 2. Requirement

When Helm v3 is installed in a cluster that is currently managed by a Helm v2 installation, the Helm v3 installation 
should be able to manage the existing v2 Releases.

Note: This proposal covers the migration use case of Helm v3 managing existing Helm v2 releases (i.e. converting v2 releases to v3 releases). Other migration use cases are covered by documentation (https://github.com/helm/helm/pull/5582).

# 3. Proposal

A standalone migration tool that migrates from Helm v2 to Helm v3. (@prydonius 
https://github.com/helm/community/issues/67#issuecomment-448033387) 

The primary function of the tool is to:

- Automatically back up Helm v2 Release and convert them to Helm v3 Release

The suggestion is for a simple, Helm-org supported plugin named `helm 2to3`. The plugin should concentrate at the 
start on its primary function of converting releases from v2 to v3 through the `convert` subcommand. It should be able 
to be extended if need be. (@jdolitsky https://github.com/helm/community/issues/67#issuecomment-448045222)

```console
$ helm 2to3 convert myrelease --dry-run

NOTE: This is in dry-run mode, the following actions will not be executed.
Run without --dry-run to take the actions described below:

[Helm 3] Release "myrelease" will be created.
[Helm 3] ReleaseVersion "myrelease.v1" will be created.
[Helm 3] ReleaseVersion "myrelease.v2" will be created.
[Helm 3] ReleaseVersion "myrelease.v3" will be created.
[Helm 2] ReleaseVersion "myrelease.v1" will be deleted.
[Helm 2] ReleaseVersion "myrelease.v2" will be deleted.
[Helm 2] ReleaseVersion "myrelease.v3" will be deleted.

$ helm 2to3 convert myrelease

[Helm 3] Release "myrelease" will be created.
[Helm 3] ReleaseVersion "myrelease.v1" will be created.
[Helm 3] ReleaseVersion "myrelease.v1" created.
[Helm 3] ReleaseVersion "myrelease.v2" will be created.
[Helm 3] ReleaseVersion "myrelease.v2" created.
[Helm 3] ReleaseVersion "myrelease.v3" will be created.
[Helm 3] ReleaseVersion "myrelease.v3" created.
[Helm 3] Release "myrelease" created.
[Helm 2] Release "myrelease" will be deleted.
[Helm 2] ReleaseVersion "myrelease.v1" will be deleted.
[Helm 2] ReleaseVersion "myrelease.v1" deleted.
[Helm 2] ReleaseVersion "myrelease.v2" will be deleted.
[Helm 2] ReleaseVersion "myrelease.v2" deleted.
[Helm 2] ReleaseVersion "myrelease.v3" will be deleted.
[Helm 2] ReleaseVersion "myrelease.v3" deleted.
[Helm 2] Release "myrelease" deleted.
Release "myrelease" was converted successfully from Helm 2 to Helm 3. 
```

## 3.1 Flow

Steps when converting Helm v2 release object to Helm 3 release object:

- Get v2 Release info
- Retrieve versions
- For each version:
  - Get v2 Release object from Helm v2 state storage
  - Map v2 Release object to v3 Release object
  - Add the v3 Release object to Helm v3 state storage (v2 version deployed namespace)

## 3.2 Assumptions

The following will be assumed:

- Underlying Kubernetes resources will be unchanged.
- The namespace(s) which the release versions were deployed to in Helm v2 system is/are still available in the cluster. Otherwise the namespace(s) need to be created manually.

## 3.3 Extensions required to the Helm v3 client

Some possible extensions to the Helm v3 client:

- For Release version object:
  - Be able to set `modifiedAt`
  - Be able to set `status`
  - Be able to set `creationTimestamp`
  - Be able to set `version`
  - Enable release objects to be added only (i.e. without adding the undelying kubernetes resources)

## 3.4 Helm convert release object code example <July 2019>

**Note: The code here is accurate as of July 2019. There maybe changes to client code after this date and before Helm v3 is released.**

Steps when converting Helm v2 release object to Helm 3 release object:

```
fmt.Printf("[Helm 3] ReleaseVersion \"%s\" will be created.", releaseVersionName)

log.Printf("[Helm 2] Get v2 Release object ....")
v2Rel, err := v2.GetRelease(releaseName, releaseVersion)
if err != nil {
        return err
}

log.Printf("[Helm 2/3] Map v2 Release object to v3 Release object....")
v3Chrt, err := v3.Mapv2ChartTov3Chart(v2Rel.Chart)
if err != nil {
        return err
}

log.Printf("[Helm 3] Add the v3 Release object ....")
cfg := v3.SetupConfig(namespace)
client := v3.GetInstallClient(cfg)
client.Namespace = namespace
client.ReleaseName = releaseVersionName
_, err = client.Run(v3Chrt)
if err != nil {
        return err
}

fmt.Printf("[Helm 3] ReleaseVersion \"%s\" created.", releaseVersionName)
```

Retrieving Helm v2 release object:

```
package v2

import (
	"k8s.io/helm/pkg/helm"
        "k8s.io/helm/pkg/proto/hapi/release"
)

func GetRelease(releaseName string, releaseVersion int32) (*release.Release, error) {
        helmClient, err := GetHelmClient()
	if err != nil {
	       return nil, err
        }

	res, err := helmClient.ReleaseContent(releaseName, helm.ContentReleaseVersion(releaseVersion))
        if err != nil {
	        return nil, err
        }

	return res.Release, nil
}
```

Connecting to Helm v2 client:

```
package v2

import (
	"fmt"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"

	"k8s.io/helm/pkg/helm"
	hpf "k8s.io/helm/pkg/helm/portforwarder"
	"k8s.io/helm/pkg/kube"
)

//hardcoded for test purposes
const (
	tillerNamespace="kube-system"
	kubectx="dind"
	histMax=256
)

func GetHelmClient() (helm.Interface, error) {
	client, config, err := getKubeClient(kubectx)
        if err != nil {
               return nil, err
        }

        tunnel, err := setupTillerConnection(client, config, tillerNamespace)
        if err != nil {
                return nil, err
        }

        return helm.NewClient(helm.Host(fmt.Sprintf("127.0.0.1:%d", tunnel.Local))), nil
}

// getKubeClient creates a Kubernetes config and client for a given kubeconfig context.
func getKubeClient(context string) (kubernetes.Interface, *rest.Config, error) {
        _, config, err := configForContext(context)
        if err != nil {
                return nil, nil, err
        }
        client, err := kubernetes.NewForConfig(config)
        if err != nil {
                return nil, nil, fmt.Errorf("could not get Kubernetes client: %s", err)
        }
        return client, config, nil
}

// configForContext creates a Kubernetes REST client configuration for a given kubeconfig context.
func configForContext(context string) (clientcmd.ClientConfig, *rest.Config, error) {
        clientConfig := kube.GetConfig(context, "")
        config, err := clientConfig.ClientConfig()
        if err != nil {
                return nil, nil, fmt.Errorf("could not get Kubernetes config for context %q: %s", context, err)
        }
        return clientConfig, config, nil
}

func setupTillerConnection(client kubernetes.Interface, config *rest.Config, namespace string) (*kube.Tunnel, error) {
        tunnel, err := hpf.New(namespace, client, config)
        if err != nil {
                return nil, fmt.Errorf("Could not get a connection to tiller: %s\nPlease ensure you have run `helm init`", err)
        }

        return tunnel, err
}
```

Converting Helm v2 release object to Helm v3 release object:

```
package v3
  
import (
	"github.com/golang/protobuf/ptypes/any"

        "helm.sh/helm/pkg/chart"

	v2ChartUtil "k8s.io/helm/pkg/chartutil"
	v2Chart "k8s.io/helm/pkg/proto/hapi/chart"
)

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
```

Connecting to Helm v3 client:

```
package v3
  
import (
	"context"
        "fmt"
        "log"
        "os"
        "path/filepath"
	"sync"
        
	auth "github.com/deislabs/oras/pkg/auth/docker"
        "k8s.io/cli-runtime/pkg/genericclioptions"

	"helm.sh/helm/pkg/action"
	"helm.sh/helm/pkg/cli"
        "helm.sh/helm/pkg/kube"
        "helm.sh/helm/pkg/registry"
        "helm.sh/helm/pkg/storage"
        "helm.sh/helm/pkg/storage/driver"
)

var (
        settings   cli.EnvSettings
        config     genericclioptions.RESTClientGetter
        configOnce sync.Once
)

func GetInstallClient(cfg *action.Configuration) *action.Install{
	return action.NewInstall(cfg)
}

func SetupConfig(namespace string) (*action.Configuration) {

        actionConfig, err := newActionConfig(namespace)
	if err != nil {
                panic(err)
        }

        // set defaults from environment
        //settings.Init()

        // Add the registry client based on settings
        // TODO: Move this elsewhere (first, settings.Init() must move)
        // TODO: handle errors, dont panic
        credentialsFile := filepath.Join(settings.Home.Registry(), registry.CredentialsFileBasename)
        client, err := auth.NewClient(credentialsFile)
        if err != nil {
                panic(err)
        }
        resolver, err := client.Resolver(context.Background())
        if err != nil {
                panic(err)
        }
        actionConfig.RegistryClient = registry.NewClient(&registry.ClientOptions{
                Debug: settings.Debug,
                //Out:   out,
                Authorizer: registry.Authorizer{
                        Client: client,
                },
                Resolver: registry.Resolver{
                        Resolver: resolver,
                },
                CacheRootDir: settings.Home.Registry(),
        })

        return actionConfig
}

func newActionConfig(namespace string) (*action.Configuration, error) {
        kc := kube.New(kubeConfig())
        kc.Log = logf

        clientset, err := kc.KubernetesClientSet()
        if err != nil {
		return nil, err
        }

        var store *storage.Storage
        switch os.Getenv("HELM_DRIVER") {
        case "secret", "secrets", "":
                d := driver.NewSecrets(clientset.CoreV1().Secrets(namespace))
                d.Log = logf
                store = storage.Init(d)
        case "configmap", "configmaps":
                d := driver.NewConfigMaps(clientset.CoreV1().ConfigMaps(namespace))
                d.Log = logf
                store = storage.Init(d)
        case "memory":
                d := driver.NewMemory()
                store = storage.Init(d)
        default:
                // Not sure what to do here.
                panic("Unknown driver in HELM_DRIVER: " + os.Getenv("HELM_DRIVER"))
        }

        return &action.Configuration{
                RESTClientGetter: kubeConfig(),
                KubeClient:       kc,
                Releases:         store,
                Log:              logf,
        }, nil
}

func kubeConfig() genericclioptions.RESTClientGetter {
        configOnce.Do(func() {
                config = kube.GetConfig(settings.KubeConfig, settings.KubeContext, settings.Namespace)
        })
        return config
}

func logf(format string, v ...interface{}) {
        if settings.Debug {
                format = fmt.Sprintf("[debug] %s\n", format)
                log.Output(2, fmt.Sprintf(format, v...))
        }
}

```

# 4. Reference

- Migration was raised at *KubeCon/CloudNativeCon Seattle 2018* at the *Helm Deep Dive session*. 
Ref: https://github.com/helm/community/issues/67.
