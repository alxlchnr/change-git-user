package main

import (
	"flag"
	"fmt"
	"github.com/alxlchnr/change-git-user/git"
	"os"
	"path/filepath"
)

var user string
var token string
var email string
var name string
var path string
var help bool
var global bool
var unset bool

func init() {
	flag.StringVar(&user, "user", "", "the API username of the new user")
	flag.StringVar(&token, "token", "", "the API token of the new user")
	flag.StringVar(&email, "email", "", "the email of the new user")
	flag.StringVar(&name, "name", "", "the name of the new user")
	flag.StringVar(&path, "path", ".", "path where to look for git repositories")
	flag.BoolVar(&global, "global", true, "apply user name and email globally")
	flag.BoolVar(&unset, "unset", false, "unset user name and email")

	flag.BoolVar(&help, "help", false, "show help")
}

var initChangeGitUser = func(params *git.ChangeGitParameters) *git.ChangeGitUser {
	return &git.ChangeGitUser{Parameters: params, GitCommands: git.GitCommands(&git.GitCommandsImpl{})}
}
var walkDirectories = filepath.Walk

func main() {
	flag.Parse()

	if help {
		fmt.Println("change-git-user")
		flag.PrintDefaults()
	} else {
		doChangeGitUser(user, token, email, name, path, global, unset)
	}
}

func doChangeGitUser(user string, token string, email string, name string, path string, global bool, unset bool) {
	params := &git.ChangeGitParameters{
		User:   user,
		Token:  token,
		Email:  email,
		Name:   name,
		Path:   path,
		Global: global,
		Unset:  unset,
	}
	fmt.Printf("parameters:  %v\n", params)

	changeGitUser := initChangeGitUser(params)

	if global {
		go changeGitUser.ChangeGlobalUserConfig()
	}

	walkDirectories(path, func(currentPath string, info os.FileInfo, err error) error {
		if info.IsDir() && filepath.Base(currentPath) == ".git" {
			go changeGitUser.ChangeUserConfigAndRemotes(filepath.Dir(currentPath))
		}
		return nil
	})
}
