package main

import (
	"fmt"
	//"io"
	//"log"
	//"os"
	//"text/template"
	//"time"

	//"github.com/ghodss/yaml"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	//"k8s.io/helm/pkg/chartutil"
	"k8s.io/helm/pkg/helm"
	"k8s.io/helm/pkg/proto/hapi/chart"
        "k8s.io/helm/pkg/proto/hapi/release"
        "k8s.io/helm/pkg/timeconv"
	hpf "k8s.io/helm/pkg/helm/portforwarder"
	"k8s.io/helm/pkg/kube"
)

const (
	tillerHost = "127.0.0.1:44134"
	tillerNamespace="kube-system"
	kubectx="dind"
	histMax=256
)

type releaseInfo struct {
        Revision    int32  `json:"revision"`
        Updated     string `json:"updated"`
        Status      string `json:"status"`
        Chart       string `json:"chart"`
        Description string `json:"description"`
}

type releaseHistory []releaseInfo

var printReleaseTemplate = `REVISION: {{.Release.Version}}
RELEASED: {{.ReleaseDate}}
CHART: {{.Release.Chart.Metadata.Name}}-{{.Release.Chart.Metadata.Version}}
USER-SUPPLIED VALUES:
{{.Release.Config.Raw}}
COMPUTED VALUES:
{{.ComputedValues}}
HOOKS:
{{- range .Release.Hooks }}
---
# {{.Name}}
{{.Manifest}}
{{- end }}
MANIFEST:
{{.Release.Manifest}}
`

func GetRelease() (*release.Release, error) {

	client, config, err := getKubeClient(kubectx)
        if err != nil {
               //log.Fatalf("Could not get a kube client: %s", err)
	       return nil, err
        }

	//`helmClient := helm.NewClient(helm.Host(tillerHost))


        helmClient, err := setupHelm(client, config, tillerNamespace)
	if err != nil {
               //log.Fatalf("err: %v", err)
	       return nil, err
        }
	//releases, err := helmClient.ListReleases()
	//if err != nil {
//		log.Fatalf("err: %v", err)
	//}
	//fmt.Printf("Releases: %v\n", releases)

	//r, err := helmClient.ReleaseHistory("chrt-v2", helm.WithMaxHistory(histMax))
        ////if err != nil {
        ////        log.Fatalf("Failed to get release history: %v", err)
        ////}
	////fmt.Printf("History: %v\n", relHistory)
	////if len(r.Releases) == 0 {
        ////        fmt.Println("No releases returned")
////		return
  ////      }

    ////    releaseHistory := getReleaseHistory(r.Releases)

        //var history []byte
        //var formattingError error

        //history, formattingError = yaml.Marshal(releaseHistory)

        //if formattingError != nil {
        //        log.Fatalf("Failed to format history: %v", err)
        //}

        //fmt.Println(string(history))

	//res, err := helmClient.ReleaseContent("mychart", helm.ContentReleaseVersion(1))
        //if err != nil {
        //        log.Fatalf("Failed to get content: %v", err)
        //}
        //fmt.Println(res)

	res, err := helmClient.ReleaseContent("chrt-v2", helm.ContentReleaseVersion(1))
        if err != nil {
                //log.Fatalf("Failed to get content: %v", err)
	        return nil, err
        }

	return res.Release, nil
	//printRelease(os.Stdout, res.Release)

}

func formatChartname(c *chart.Chart) string {
        if c == nil || c.Metadata == nil {
                // This is an edge case that has happened in prod, though we don't
                // know how: https://github.com/kubernetes/helm/issues/1347
                return "MISSING"
        }
        return fmt.Sprintf("%s-%s", c.Metadata.Name, c.Metadata.Version)
}

func getReleaseHistory(rls []*release.Release) (history releaseHistory) {
        for i := len(rls) - 1; i >= 0; i-- {
                r := rls[i]
                c := formatChartname(r.Chart)
                t := timeconv.String(r.Info.LastDeployed)
                s := r.Info.Status.Code.String()
                v := r.Version
                d := r.Info.Description

                rInfo := releaseInfo{
                        Revision:    v,
                        Updated:     t,
                        Status:      s,
                        Chart:       c,
                        Description: d,
                }
                history = append(history, rInfo)
        }

        return history
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

func setupHelm(kubeClient kubernetes.Interface, config *rest.Config, namespace string) (helm.Interface, error) {
        tunnel, err := setupTillerConnection(kubeClient, config, namespace)
        if err != nil {
                return nil, err
        }

        return helm.NewClient(helm.Host(fmt.Sprintf("127.0.0.1:%d", tunnel.Local))), nil
}

func setupTillerConnection(client kubernetes.Interface, config *rest.Config, namespace string) (*kube.Tunnel, error) {
        tunnel, err := hpf.New(namespace, client, config)
        if err != nil {
                return nil, fmt.Errorf("Could not get a connection to tiller: %s\nPlease ensure you have run `helm init`", err)
        }

        return tunnel, err
}

//func printRelease(out io.Writer, rel *release.Release) error {
////        if rel == nil {
////                return nil
////        }
////
////        cfg, err := chartutil.CoalesceValues(rel.Chart, rel.Config)
//        if err != nil {
//                return err
//        }
//        cfgStr, err := cfg.YAML()
//        if err != nil {
//                return err
//        }
//
//        data := map[string]interface{}{
//                "Release":        rel,
//                "ComputedValues": cfgStr,
//                "ReleaseDate":    timeconv.Format(rel.Info.LastDeployed, time.ANSIC),
//        }
//        return tpl(printReleaseTemplate, data, out)
//}
//
//func tpl(t string, vals interface{}, out io.Writer) error {
//        tt, err := template.New("_").Parse(t)
//        if err != nil {
//                return err
//        }
//        return tt.Execute(out, vals)
//}

