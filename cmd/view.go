package cmd

import (
	"fmt"
	"io/ioutil"

	"github.com/spf13/cobra"
)

// viewCmd represents the view command
var viewCmd = &cobra.Command{
	Args:  cobra.NoArgs,
	Use:   "view",
	Short: "View contents of config file",
	Long:  `View contents of config file`,
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := ioutil.ReadFile(getConfigPath())
		if err != nil {
			panic(fmt.Errorf("Could not read config file: %v", err))
		}

		fmt.Printf("Config file: %s \n----------- \n%v", getConfigPath(), string(cfg))
	},
}

func init() {
	configCmd.AddCommand(viewCmd)
}
