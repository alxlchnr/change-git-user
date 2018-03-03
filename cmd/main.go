package cmd

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
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
		fmt.Printf("global:  %v\n", global)
		fmt.Printf("unset:  %v\n", unset)

		if global {
			go setGitUserName(name, global, unset, "")

			go setGitUserEmail(email, global, unset, "")
		}

		filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
			if info.IsDir() && filepath.Base(path) == ".git" {
				go processGitRepo(user, token, path)
			}
			return nil
		})
	}
}
func processGitRepo(newUser string, newToken string, path string) {
	gitRepoPath := filepath.Dir(path)
	remotes := getGitRemoteUrls(gitRepoPath)
	remotesToUrls := gitRemotesToMap(remotes)
	for key, value := range remotesToUrls {
		go changeGitRemoteUrl(newUser, newToken, value, key, gitRepoPath)
		if !global {
			go setGitUserEmail(email, false, unset, gitRepoPath)
			go setGitUserName(name, false, unset, gitRepoPath)
		}
	}
}

func changeGitRemoteUrl(newUser string, newToken string, remoteUrl *GitUrl, remote string, gitRepoPath string) {
	oldUrl := remoteUrl.ToUrl()
	remoteUrl.SetUser(newUser)
	remoteUrl.SetToken(newToken)
	newUrl := remoteUrl.ToUrl()
	fmt.Printf("Change git remote %s \nfrom %s \nto new url %s.\n", remote, oldUrl, newUrl)
	cmd := exec.Command("git", "remote", "set-url", remote, newUrl)
	cmd.Dir = gitRepoPath
	out, _ := cmd.Output()
	fmt.Println(string(out))
}

func setGitUserName(newValue string, globally bool, unset bool, path string) {
	setGitConfig("user.name", newValue, globally, unset, path)
}

func setGitUserEmail(newValue string, globally bool, unset bool, path string) {
	setGitConfig("user.email", newValue, globally, unset, path)
}

func setGitConfig(config string, value string, globally bool, unset bool, path string) {
	args := []string{"config"}
	if len(value) > 0 {
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

	cmd := exec.Command("git", args...)
	if !globally {
		cmd.Dir = path
	}

	out, _ := cmd.Output()
	fmt.Println(string(out))
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
