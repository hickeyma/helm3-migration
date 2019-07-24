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
