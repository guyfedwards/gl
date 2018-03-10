package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// Payload for creating merge request
type Payload struct {
	SourceBranch string `json:"source_branch"`
	TargetBranch string `json:"target_branch"`
	Title        string `json:"title"`
}

// mergeRequestCmd represents the merge command
var mergeRequestCmd = &cobra.Command{
	Args:    cobra.OnlyValidArgs,
	Aliases: []string{"mr"},
	Use:     "merge-request",
	Short:   "Create merge request",
	Long: `Create merge request for a branch. Defaults to current branch and
origin/master. Example usage:

$ gl merge`,
	Run: func(cmd *cobra.Command, args []string) {
		if !viper.IsSet("token") {
			fmt.Println(`You must set Gitlab token before proceeding.
$ gl config set token <token>
			`)
			os.Exit(1)
		}
		//TODO get current branch as default source
		var (
			source string = getCurrentBranch()
			target string = "master"
			title  string = "WIP: Default Title"
		)

		fmt.Println(source, target, title)

		// TODO url encode project
		e := url.PathEscape("frontend/common-ui")
		fmt.Println(e)

		//store api version in config
		// TODO create from remote
		url := "https://gitlab.algomi.net/api/v4/projects/frontend%2Fcommon-ui/merge_requests"

		p := &Payload{
			SourceBranch: "ge-proptypes",
			TargetBranch: "master",
			Title:        "WIP: test",
		}
		j, _ := json.Marshal(p)

		// client := &http.Client{}

		b := bytes.NewBuffer([]byte(j))
		req, err := http.NewRequest("POST", url, b)
		if err != nil {
			panic("Error creating request: " + err.Error())
		}

		token := viper.Get("token").(string)
		req.Header.Set("Private-Token", token)
		req.Header.Set("Content-Type", "application/json")

		// resp, err := client.Do(req)
		// if err != nil {
		// 	fmt.Printf("Error making request: %v", err)
		// }
		// defer resp.Body.Close()

		// body, _ := ioutil.ReadAll(resp.Body)
		// fmt.Printf("%v", string(body))
	},
}

func init() {
	rootCmd.AddCommand(mergeRequestCmd)

}

func buildURL(base string, group string, project string) string {

	return strings.Join([]string{"https://gitlab.algomi.net", "/api/v4/", ""}, "")
}

func getCurrentBranch() string {
	var (
		cmdOut []byte
		err    error
	)

	if cmdOut, err = exec.Command("git", "rev-parse", "--abbrev-ref HEAD").Output(); err != nil {
		panic("Error running `git rev-parse --abbrev-ref HEAD` \n" + err.Error())
	}

	return string(cmdOut)
}

func encodeProjectName(url string) string {
	return ""
}
