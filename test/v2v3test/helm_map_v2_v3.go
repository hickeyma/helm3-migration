package main
  
import (
	"context"
        "fmt"
	"io"
        "log"
        "os"
        "path/filepath"
	"sync"
	"text/template"
	"time"
        
	auth "github.com/deislabs/oras/pkg/auth/docker"
        "k8s.io/cli-runtime/pkg/genericclioptions"

	"helm.sh/helm/pkg/action"
        "helm.sh/helm/pkg/chart"
	"helm.sh/helm/pkg/cli"
        "helm.sh/helm/pkg/chart/loader"
        "helm.sh/helm/pkg/kube"
        "helm.sh/helm/pkg/registry"
        "helm.sh/helm/pkg/storage"
        "helm.sh/helm/pkg/storage/driver"

	v2Chart "k8s.io/helm/pkg/proto/hapi/chart"
	v2ChartUtil "k8s.io/helm/pkg/chartutil"
	v2Release "k8s.io/helm/pkg/proto/hapi/release"
	v2TimeConv "k8s.io/helm/pkg/timeconv"
)

func main() {
        var chartPath = "/home/usr1/test/helm-charts/chrt-v2"
	var relName = "chrt-v2"

        fmt.Printf("Get v2 release info ....")
	v2Rel, err := GetRelease()
	if err != nil {
                log.Fatalf("Failed to get content: %v", err)
                return
        }
	printRelease(os.Stdout, v2Rel)

	v3Chrt := mapv2ChrtTov3Chrt(v2Rel.Chart)
        fmt.Printf("v3 Chart: %q\n", v3Chrt)
        fmt.Printf("Add v2 release info to v3 state ... ")
        cfg := setupConfig()
        client := action.NewInstall(cfg)

        loadedChrt, err := loadChart(chartPath)
        if err != nil {
                fmt.Printf("Error loading chart: %q\n", err)
		return
        } 

	client.Namespace = getNamespace()
	client.ReleaseName = relName
	//client.DryRun = true

	rel, err := client.Run(loadedChrt)
	if err != nil {
                fmt.Printf("Error loading chart: %q\n", err)
        } else {
		fmt.Printf("Chart details .........\n")
		fmt.Printf("%q\n", rel)
	}

}

var (
        settings   cli.EnvSettings
        config     genericclioptions.RESTClientGetter
        configOnce sync.Once
)

func loadChart(chartPath string) (*chart.Chart, error) {
        // Check chart requirements to make sure all dependencies are present in /charts
        chartRequested, err := loader.Load(chartPath)
        if err != nil {
                return nil, err
        }

	if req := chartRequested.Metadata.Dependencies; req != nil {
                // If checkDependencies returns an error, we have unfulfilled dependencies.
                // As of Helm 2.4.0, this is treated as a stopping condition:
                // https://github.com/kubernetes/helm/issues/2209
                if err := action.CheckDependencies(chartRequested, req); err != nil {
                        return nil, err
                }
        }  
        return chartRequested, nil
}

func setupConfig() (*action.Configuration) {

        actionConfig := newActionConfig(false)

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

func newActionConfig(allNamespaces bool) *action.Configuration {
        kc := kube.New(kubeConfig())
        kc.Log = logf

        clientset, err := kc.KubernetesClientSet()
        if err != nil {
                // TODO return error
                log.Fatal(err)
        }
        var namespace string
        if !allNamespaces {
                namespace = getNamespace()
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
        }
}

func kubeConfig() genericclioptions.RESTClientGetter {
        configOnce.Do(func() {
                config = kube.GetConfig(settings.KubeConfig, settings.KubeContext, settings.Namespace)
        })
        return config
}

func getNamespace() string {
        //if ns, _, err := kubeConfig().ToRawKubeConfigLoader().Namespace(); err == nil {
        //        return ns
        //}
        return "default"
}

func logf(format string, v ...interface{}) {
        if settings.Debug {
                format = fmt.Sprintf("[debug] %s\n", format)
                log.Output(2, fmt.Sprintf(format, v...))
        }
}

func mapv2ChrtTov3Chrt(v2Chrt *v2Chart.Chart) (*chart.Chart) {
	 v3Chrt := new(chart.Chart)

	 v3Chrt.Metadata = mapMetadata(v2Chrt)

	 v3Chrt.Lock = new(chart.Lock)

	 return v3Chrt
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
	//TODO: metadata.Type = 

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

func printRelease(out io.Writer, rel *v2Release.Release) error {
        if rel == nil {
                return nil
        }

        cfg, err := v2ChartUtil.CoalesceValues(rel.Chart, rel.Config)
        if err != nil {
                return err
        }
        cfgStr, err := cfg.YAML()
        if err != nil {
                return err
        }

        data := map[string]interface{}{
                "Release":        rel,
                "ComputedValues": cfgStr,
                "ReleaseDate":    v2TimeConv.Format(rel.Info.LastDeployed, time.ANSIC),
        }
        return tpl(printReleaseTemplate, data, out)
}

func tpl(t string, vals interface{}, out io.Writer) error {
        tt, err := template.New("_").Parse(t)
        if err != nil {
                return err
        }
        return tt.Execute(out, vals)
}
