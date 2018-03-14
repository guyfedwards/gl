package cmd

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"os/exec"
	"strings"

	"github.com/spf13/viper"
)

// TODO make work with both ssh/https remotes
func getRemoteParts() []string {
	var (
		cmdOut []byte
		err    error
	)
	if cmdOut, err = exec.Command("git", "remote", "get-url", "origin").Output(); err != nil {
		log.Fatal("Error executing git command: ", err)
	}

	r := strings.TrimSpace(string(cmdOut))
	replaced := replaceString(r)

	// naive implementation assuming no path to gitlab host
	s := strings.Split(replaced, "/")
	base := s[:3]
	path := s[3:]
	return []string{strings.Join(base, "/"), strings.Join(path, "/")}
}

func replaceString(s string) string {
	r := strings.NewReplacer("git@", "https://", ":", "/", ".git", "")
	return r.Replace(s)
}

func tokenCheck() {
	if !viper.IsSet("token") {
		fmt.Println(`You must set Gitlab token before proceeding.
$ gl config set token <token>
		`)
		os.Exit(1)
	}
}

func buildURL(endpoint string) string {
	p := url.PathEscape(projectFlag)
	base := getRemoteParts()[0]

	return strings.Join([]string{base, "/api/v4/", p, "/", endpoint}, "")
}

func getCurrentProject() string {
	p := getRemoteParts()[1]

	return p
}

// TODO get user projects and match with remote string
// func getUserProjects() {
// 	resp, err := http.Get()
// }
