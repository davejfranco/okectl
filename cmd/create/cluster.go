package create

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	clusterName string
	region      string
	filePath    string
)

// clusterCmd represents the cluster command
var clusterCmd = &cobra.Command{
	Use:   "cluster",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func init() {

	//Cluster name
	clusterCmd.Flags().StringVarP(&clusterName, "name", "n", "", "Name of the cluster to create")
	if err := clusterCmd.MarkFlagRequired("name"); err != nil {
		fmt.Println(err)
	}
	//Cluster region
	clusterCmd.Flags().StringVarP(&region, "region", "r", "", "Region to create the cluster in")
	//Cluster config file
	clusterCmd.Flags().StringVarP(&filePath, "file", "f", "", "Path to the cluster config file")

	//rootCmd.AddCommand(clusterCmd)
	CreateCmd.AddCommand(clusterCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// clusterCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// clusterCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
