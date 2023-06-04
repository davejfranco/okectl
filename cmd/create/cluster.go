package create

import (
	"fmt"
	"okectl/pkg/ctl"
	"okectl/pkg/oci"

	"github.com/spf13/cobra"
)

var (
	clusterName   string
	compartmentID string
	k8sVersion    string
	filePath      string
)

// clusterCmd represents the cluster command
var clusterCmd = &cobra.Command{
	Use:   "cluster",
	Short: "this command creates a cluster",
	Long: `Create a cluster using the following command:

okectl create cluster --name <cluster-name> --compartment-id <compartment-id>
or
okectl create cluster -f <cluster-config-file>.yaml`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		if filePath == "" {
			cmd.MarkFlagRequired("compartment-id")
			cmd.MarkFlagRequired("name")
			cmd.SetUsageTemplate(fmt.Sprintf("%s\nRequired flags:\n  -c, --compartment-id\n  -n, --name\n", cmd.UsageString()))
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		if filePath != "" {
			ctl.CreateClusterFromFile(filePath)
			return
		}
		if k8sVersion == "" {
			availableVersions := oci.GetKubernetesVersion()
			k8sVersion = availableVersions[len(availableVersions)-1]
		}
		ctl.CreateCluster(clusterName, compartmentID, k8sVersion)
	},
}

func init() {

	//Cluster name
	clusterCmd.Flags().StringVarP(&clusterName, "name", "n", "", "Name of the cluster to create")
	//Compartment ID
	clusterCmd.Flags().StringVarP(&compartmentID, "compartment-id", "c", "", "The OCID of the compartment where the cluster will be created")
	//Kubernetes Cluster version
	clusterCmd.Flags().StringVarP(&k8sVersion, "version", "v", "", "Version of kubernetes to use. If is not provided, the latest version will be used")

	//Cluster config file
	clusterCmd.Flags().StringVarP(&filePath, "file", "f", "", "Path to the cluster config file")
	//Add commands to create pallet
	CreateCmd.AddCommand(clusterCmd)

}
