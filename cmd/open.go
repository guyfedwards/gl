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
		openBrowser(strings.Join(getRemoteParts(), "/"))
	},
}

func init() {
	rootCmd.AddCommand(openCmd)
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
