package lib

import (
	"os/exec"
	"strings"
)

type Git struct {
	RepositoryPath string
}

func (g Git) LogParser(rawGitLog string) {
	log := strings.Split(rawGitLog, "~")

	if len(log) != 3 {
		return
	}

	commit := Commit{Hash: log[0], Author: log[1], Subject: log[2]}
	commit.SendNotification()
	commit.Save()
}

func (g Git) Fetch() {
	cmd := exec.Command("git", "fetch")
	cmd.Dir = g.RepositoryPath
	cmd.Run()
}

func (g Git) Log(ch chan string, filename string) {
	cmd := exec.Command("git", "log", "--oneline", "--all", "--pretty=%H~%an~%s", `--since="5 minutes ago"`, `*/`+filename+"*")
	cmd.Dir = g.RepositoryPath
	output, _ := cmd.Output()
	ch <- string(output)
}
