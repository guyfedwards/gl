package main

import (
	"flag"
	"fmt"
	"log"
	"os/exec"
	"runtime"
	"strings"
)

func main() {
	ciPtr := flag.Bool("c", false, "Open CI/CD page")
	regPtr := flag.Bool("r", false, "Open registry page")
	merPtr := flag.Bool("m", false, "Open merge requests page")
	issPtr := flag.Bool("i", false, "Open issues page")
	wikPtr := flag.Bool("w", false, "Open wiki page")
	setPtr := flag.Bool("s", false, "Open settings page")

	flag.Parse()

	var page string

	if *ciPtr {
		page = "pipelines"
	} else if *regPtr {
		page = "container_registry"
	} else if *merPtr {
		page = "merge_requests"
	} else if *issPtr {
		page = "issues"
	} else if *wikPtr {
		page = "wikis/home"
	} else if *setPtr {
		page = "edit"
	}

	openBrowser(getRemoteURL() + "/" + page)
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
