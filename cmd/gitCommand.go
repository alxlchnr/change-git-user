package cmd

import (
	"fmt"
	"os/exec"
)

type Command interface {
	Output() ([]byte, error)
	SetWorkingDir(dir string)
}

type osCommand struct {
	Cmd *exec.Cmd
}

func (cmd *osCommand) Output() ([]byte, error) {
	return cmd.Cmd.Output()
}

func (cmd *osCommand) SetWorkingDir(dir string) {
	cmd.Cmd.Dir = dir
}

var execCommand = func(name string, args ...string) Command {
	return Command(&osCommand{exec.Command(name, args...)})
}

type GitCommands struct {
}

func (*GitCommands) ChangeGitConfig(config string, value string, globally bool, unset bool, path string) {
	args := []string{"config"}
	if len(value) > 0 && !unset {
		args = append(args, config, value)
		fmt.Printf("Change %s to %s\n", config, value)
	} else {
		fmt.Printf("No %s provided. Will not change config.\n", config)
	}

	if unset {
		args = append(args, "--unset", config)
		fmt.Printf("Unset %s\n", config)
	}

	if globally {
		args = append(args, "--global")
	}

	var cmd = execCommand("git", args...)
	if !globally {
		cmd.SetWorkingDir(path)
	}

	out, _ := cmd.Output()
	fmt.Println(string(out))
}

func (*GitCommands) GetGitRemoteUrls(gitRepoPath string) string {
	fmt.Println("Found a git repository at " + gitRepoPath)
	cmd := execCommand("git", "remote", "-v")
	cmd.SetWorkingDir(gitRepoPath)
	out, _ := cmd.Output()
	remotes := string(out)
	return remotes
}

func (*GitCommands) ChangeGitRemoteUrl(newUser string, newToken string, remoteUrl *GitUrl, remote string, gitRepoPath string) {
	oldUrl := remoteUrl.ToUrl()
	remoteUrl.SetUser(newUser)
	remoteUrl.SetToken(newToken)
	newUrl := remoteUrl.ToUrl()
	fmt.Printf("Change git remote %s \nfrom %s \nto new url %s.\n", remote, oldUrl, newUrl)
	cmd := execCommand("git", "remote", "set-url", remote, newUrl)
	cmd.SetWorkingDir(gitRepoPath)
	out, _ := cmd.Output()
	fmt.Println(string(out))
}
