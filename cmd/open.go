package cmd

import (
	"fmt"
	"log"
	"os/exec"
	"runtime"
	"strings"

	"github.com/spf13/cobra"
)

// openCmd represents the open command
var openCmd = &cobra.Command{
	Aliases: []string{"o"},
	Args:    cobra.NoArgs,
	Use:     "open",
	Short:   "Open in gitlab",
	Long:    `Open current repo in gitlab. To go to a specific page, pass options`,
	Run: func(cmd *cobra.Command, args []string) {
		openBrowser(getRemoteURL())
	},
}

func init() {
	rootCmd.AddCommand(openCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// openCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// openCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func openBrowser(url string) {
	var err error

	switch runtime.GOOS {
	case "darwin":
		err = exec.Command("open", url).Start()
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("start", url).Start()
	default:
		err = fmt.Errorf("Unsupported platform")
	}
	if err != nil {
		log.Fatal(err)
	}
}

func getRemoteURL() string {
	var (
		cmdOut []byte
		err    error
	)
	if cmdOut, err = exec.Command("git", "remote", "-v").Output(); err != nil {
		log.Fatal("Error executing git command: ", err)
	}

	rem := strings.Split(string(cmdOut), "\n")[0]
	return replaceString(rem)
}

func replaceString(s string) string {
	r := strings.NewReplacer("git@", "http://", ":", "/", ".git", "")
	return r.Replace(strings.Fields(s)[1])
}
