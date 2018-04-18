package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os/exec"
	"strconv"
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
	AssigneeID   int    `json:"assignee_id"`
}

type Res struct {
	WebURL string `json:"web_url"`
}

var (
	target      string
	source      string
	title       string
	description string
	assigneeID  string
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
			log.Fatal("You must provide a title. \n$ gl merge-request -t \"This is my title\"\n")
		}

		url := buildURL("merge_requests")

		aID, err := strconv.Atoi(assigneeID)
		if err != nil {
			log.Fatal("Error converting assigneeID to int: ", err)
		}

		p := &Payload{
			SourceBranch: source,
			TargetBranch: target,
			Title:        title,
			Description:  description,
			AssigneeID:   aID,
		}
		j, err := json.Marshal(p)
		if err != nil {
			log.Fatal("Error marshalling JSON", err)
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

		var res Res
		err = json.Unmarshal(body, &res)
		if err != nil {
			fmt.Printf("Error unmarshalling JSON response: %v", err)
		}
		fmt.Printf("%v\n", res.WebURL)
	},
}

func init() {
	rootCmd.AddCommand(mergeRequestCmd)
	mergeRequestCmd.Flags().StringVarP(&source, "source", "S", getCurrentBranch(), "Source branch")
	mergeRequestCmd.Flags().StringVarP(&target, "target", "T", "master", "Target branch")
	mergeRequestCmd.Flags().StringVarP(&title, "title", "t", "", "Title for merge request")
	mergeRequestCmd.Flags().StringVarP(&description, "description", "d", "", "Description for merge request")
	mergeRequestCmd.Flags().StringVarP(&assigneeID, "assigneeID", "a", "", "Assignee id ")
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
