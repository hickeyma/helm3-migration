package main
  
import (
        "fmt"
        //"log"
        //"os"

	"v2v3test/v2"
	"v2v3test/v3"
)

func main() {
        //var chartPath = "/home/usr1/test/helm-charts/chrt-v2"
	var relName = "chrt-v2"

        fmt.Printf("Get v2 release info ....\n")
	v2Rel, err := v2.GetRelease()
	if err != nil {
		fmt.Printf("ERROR: Failed to get content: %v", err)
                return
        }
	//v2.PrintRelease(os.Stdout, v2Rel)

        fmt.Printf("map v2 release info to equivalent v3 info....\n")
	v3Chrt := v3.Mapv2ChrtTov3Chrt(v2Rel.Chart)
	//fmt.Printf("In main, chart values: %q\n", v3Chrt.Values)
        //fmt.Printf("v3 Chart: %q\n", v3Chrt)

        fmt.Printf("Add v2 release info to v3 state ... \n")
        cfg := v3.SetupConfig()
        client := v3.GetInstallClient(cfg)

        //loadedChrt, err := v3.LoadChart(chartPath)
        //if err != nil {
        //        fmt.Printf("Error loading chart: %q\n", err)
//		return
 //       } 
//	fmt.Printf("In main, loaded chart values: %q\n", loadedChrt.Values)

	client.Namespace = getNamespace()
	client.ReleaseName = relName
	//client.DryRun = true

	//rel, err := client.Run(loadedChrt)
	rel, err := client.Run(v3Chrt)
	if err != nil {
		fmt.Printf("ERROR: migrating v2 release to v3: %q\n", err)
		return
        }

	fmt.Printf("Chart details .........\n")
	fmt.Printf("%q\n\n", rel)
	fmt.Printf("Succeeded: Migrated v2 info to v3\n")
}

func getNamespace() string {
        //if ns, _, err := kubeConfig().ToRawKubeConfigLoader().Namespace(); err == nil {
        //        return ns
        //}
        return "default"
}
