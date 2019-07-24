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
