package cmd

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/user"
	"path"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var projectFlag string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "gl",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	//	Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.Flags().StringVarP(&projectFlag, "project", "p", getCurrentProject(), "Project to run command against. Defaults to current repo.")
}

func getHomeDir() string {
	u, err := user.Current()
	if err != nil {
		fmt.Println("Can't get your home directory.")
		os.Exit(1)
	}
	return u.HomeDir
}

func getConfigPath() string {
	return path.Join(getHomeDir(), ".gl.yaml")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	// Find home directory.
	home, err := homedir.Dir()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Search config in home directory with name ".gl" (without extension).
	viper.AddConfigPath(home)
	viper.SetConfigName(".gl")

	viper.AutomaticEnv() // read in environment variables that match

	// create config file if not exists
	_, err = os.Stat(getConfigPath())
	if err != nil {
		err := ioutil.WriteFile(getConfigPath(), []byte{}, 0644)
		if err != nil {
			fmt.Printf("Could not create config file: %v", err)
			os.Exit(1)
		}
	}

	_ = viper.ReadInConfig()
}
