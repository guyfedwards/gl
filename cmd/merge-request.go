package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
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
	Description  string `json:"description"`
}

var (
	target      string
	source      string
	title       string
	description string
)

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
		tokenCheck()

		if title == "" {
			fmt.Printf("You must provide a title. \n$ gl merge-request -T \"This is my title\"\n")
			os.Exit(1)
		}

		url := buildURL("merge_requests")
		fmt.Println(url)

		p := &Payload{
			SourceBranch: source,
			TargetBranch: target,
			Title:        title,
			Description:  description,
		}
		j, err := json.Marshal(p)
		if err != nil {
			panic(fmt.Sprintln("Error marshalling JSON"))
		}

		client := &http.Client{}

		b := bytes.NewBuffer([]byte(j))
		req, err := http.NewRequest("POST", url, b)
		if err != nil {
			panic("Error creating request: " + err.Error())
		}

		token := viper.Get("token").(string)
		req.Header.Set("Private-Token", token)
		req.Header.Set("Content-Type", "application/json")

		resp, err := client.Do(req)
		if err != nil {
			fmt.Printf("Error making request: %v", err)
		}
		defer resp.Body.Close()

		body, _ := ioutil.ReadAll(resp.Body)
		fmt.Printf("%v", string(body))
	},
}

func init() {
	rootCmd.AddCommand(mergeRequestCmd)
	// TODO add flags target, source, title, desc, approvers etc
	mergeRequestCmd.Flags().StringVarP(&source, "source", "s", getCurrentBranch(), "Source branch")
	mergeRequestCmd.Flags().StringVarP(&target, "target", "t", "master", "Target branch")
	mergeRequestCmd.Flags().StringVarP(&title, "title", "T", "", "Title for merge request")
	mergeRequestCmd.Flags().StringVarP(&description, "description", "d", "", "Description for merge request")
}

func getCurrentBranch() string {
	var (
		cmdOut []byte
		err    error
	)

	if cmdOut, err = exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD").Output(); err != nil {
		panic("Error running `git rev-parse --abbrev-ref HEAD` \n" + err.Error())
	}

	return strings.TrimSpace(string(cmdOut))
}
