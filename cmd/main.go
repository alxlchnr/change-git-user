package cmd

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

var user string
var token string
var email string
var name string
var path string
var help bool
var global bool
var unset bool
var gitCommands *GitCommands
var walkDirectories = filepath.Walk

func init() {
	flag.StringVar(&user, "user", "", "the API username of the new user")
	flag.StringVar(&token, "token", "", "the API token of the new user")
	flag.StringVar(&email, "email", "", "the email of the new user")
	flag.StringVar(&name, "name", "", "the name of the new user")
	flag.StringVar(&path, "path", ".", "path where to look for git repositories")
	flag.BoolVar(&global, "global", true, "apply user name and email globally")
	flag.BoolVar(&unset, "unset", false, "unset user name and email")

	flag.BoolVar(&help, "help", false, "show help")
	gitCommands = &GitCommands{}
}

func main() {
	flag.Parse()

	if help {
		fmt.Println("change-git-user")
		flag.PrintDefaults()
	} else {
		params := &ChangeGitParameters{
			User:   user,
			Token:  token,
			Email:  email,
			Name:   name,
			Path:   path,
			Global: global,
			Unset:  unset,
		}
		fmt.Printf("parameters:  %v\n", params)

		if global {
			go gitCommands.ChangeGitConfig("user.name", params.Name, params.Global, params.Unset, "")
			go gitCommands.ChangeGitConfig("user.email", params.Email, params.Global, params.Unset, "")
		}

		walkDirectories(path, func(path string, info os.FileInfo, err error) error {
			if info.IsDir() && filepath.Base(path) == ".git" {
				go processGitRepo(params)
			}
			return nil
		})
	}
}
func processGitRepo(params *ChangeGitParameters) {
	gitRepoPath := filepath.Dir(path)
	remotes := gitCommands.GetGitRemoteUrls(gitRepoPath)
	remotesToUrls := gitRemotesToMap(remotes)
	for key, value := range remotesToUrls {
		go gitCommands.ChangeGitRemoteUrl(params.User, params.Token, value, key, gitRepoPath)
		if !params.Global {
			go gitCommands.ChangeGitConfig("user.email", params.Email, false, params.Unset, gitRepoPath)
			go gitCommands.ChangeGitConfig("user.name", params.Name, false, params.Unset, gitRepoPath)
		}
	}
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
