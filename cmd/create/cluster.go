package create

import (
	"fmt"
	"okectl/pkg/ctl"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	clusterName   string
	compartmentID string
	k8sVersion    string
	region        string
	cfgfile       string
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
		if cfgfile == "" {
			cmd.MarkFlagRequired("compartment-id")
			cmd.MarkFlagRequired("name")
			cmd.SetUsageTemplate(fmt.Sprintf("%s\nRequired flags:\n  -c, --compartment-id\n  -n, --name\n", cmd.UsageString()))
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		if cfgfile != "" {
			clusterDetails := viper.GetViper().GetStringMapString("cluster")
			if len(clusterDetails) == 0 {
				fmt.Println("Cluster details are required")
				os.Exit(1)
			}

			clusterName = clusterDetails["name"]
			compartmentID = clusterDetails["compartmentid"] //maps remove the uppercase from the key
			k8sVersion = clusterDetails["version"]
			region = clusterDetails["region"]

			if clusterName == "" || compartmentID == "" {
				fmt.Println("cluster name and compartment id are required")
				os.Exit(1)
			}
		}
		ctl.NewCluster(clusterName, compartmentID, k8sVersion, region)
	},
}

func init() {
	cobra.OnInitialize(initConfig)
	//Cluster name
	clusterCmd.Flags().StringVarP(&clusterName, "name", "n", "", "Name of the cluster to create")
	//Compartment ID
	clusterCmd.Flags().StringVarP(&compartmentID, "compartment-id", "c", "", "The OCID of the compartment where the cluster will be created")
	//Kubernetes Cluster version
	clusterCmd.Flags().StringVarP(&k8sVersion, "version", "v", "", "Version of kubernetes to use. If is not provided, the latest version will be used")
	//Cluster region
	clusterCmd.Flags().StringVarP(&region, "region", "r", "", "Region where the cluster will be created")
	//Cluster config file
	clusterCmd.Flags().StringVarP(&cfgfile, "config-file", "f", "", "Path to the cluster config file")
	//Add commands to create pallet
	CreateCmd.AddCommand(clusterCmd)

}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgfile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgfile)
	}

	//viper.AutomaticEnv() // read in environment variables that match
	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}
