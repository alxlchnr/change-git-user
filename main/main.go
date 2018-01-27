package main

import (
	"flag"
	"path/filepath"
	"os"
	"os/exec"
	"fmt"
	"strings"
)

var user string
var token string
var email string
var name string
var path string
var help bool

func init() {
	flag.StringVar(&user, "user", "", "the API username of the new user")
	flag.StringVar(&token, "token", "", "the API token of the new user")
	flag.StringVar(&email, "email", "", "the email of the new user")
	flag.StringVar(&name, "name", "", "the name of the new user")
	flag.StringVar(&path, "path", ".", "path where to look for git repositories")

	flag.BoolVar(&help, "help", false, "show help")
}

func main() {
	flag.Parse()

	if help {
		fmt.Println("change-git-user")
		flag.PrintDefaults()
	} else {
		fmt.Printf("user:  %s\n", user)
		fmt.Printf("token: %s\n", token)
		fmt.Printf("email: %s\n", email)
		fmt.Printf("name:  %s\n", name)
		fmt.Printf("path:  %s\n", path)

		go setGitUserNameGlobally()

		go setGitUserEmailGlobally()

		filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
			if info.IsDir() && filepath.Base(path) == ".git" {
				go processGitRepo(path)
			}
			return nil;
		})
	}
}
func processGitRepo(path string) {
	gitRepoPath := filepath.Dir(path)
	remotes := getGitRemoteUrls(gitRepoPath)
	remotesToUrls := gitRemotesToMap(remotes)
	for key, value := range remotesToUrls {
		go changeGitRemoteUrl(value, key, gitRepoPath)
	}
}

func changeGitRemoteUrl(remoteUrl *GitUrl, remote string, gitRepoPath string) {
	oldUrl := remoteUrl.ToUrl()
	remoteUrl.SetUser(user)
	remoteUrl.SetToken(token)
	newUrl := remoteUrl.ToUrl()
	fmt.Printf("Change git remote %s \nfrom %s \nto new url %s.\n", remote, oldUrl, newUrl)
	cmd := exec.Command("git", "remote", "set-url", remote, newUrl)
	cmd.Dir = gitRepoPath
	out, _ := cmd.Output()
	fmt.Println(string(out))
}

func setGitUserNameGlobally() {
	if len(name) > 0 {
		fmt.Printf("Set new global git user %s\n", name)
		cmd := exec.Command("git", "config", "--global", "user.name", name)
		out, _ := cmd.Output()
		fmt.Println(string(out))
	} else {
		fmt.Println("No user name provided. Will not change global config.")
	}
}

func setGitUserEmailGlobally() {
	if len(email) > 0 {
		fmt.Printf("Set new global git user email %s\n", name)
		cmd := exec.Command("git", "config", "--global", "user.email", email)
		out, _ := cmd.Output()
		fmt.Println(string(out))
	} else {
		fmt.Println("No email address provided. Will not change global config.")
	}
}

func getGitRemoteUrls(gitRepoPath string) string {
	fmt.Println("Found a git repository at " + gitRepoPath)
	cmd := exec.Command("git", "remote", "-v")
	cmd.Dir = gitRepoPath
	out, _ := cmd.Output()
	remotes := string(out)
	return remotes
}

func gitRemotesToMap(remotes string) map[string]*GitUrl {
	remoteUrls := strings.Split(remotes, "\n")
	remotesToUrls := make(map[string]*GitUrl)
	for _, line := range remoteUrls {
		if len(line) > 0 {
			split := strings.Split(strings.Split(line, " ")[0], "\t")
			remotesToUrls[split[0]] = NewGitUrl(split[1])
		}
	}
	return remotesToUrls
}
