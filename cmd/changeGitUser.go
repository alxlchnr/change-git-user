package cmd

import (
	"strings"
)

type ChangeGitUser struct {
	Parameters  *ChangeGitParameters
	GitCommands GitCommands
}

func (cgu *ChangeGitUser) ChangeGlobalUserConfig() {
	cgu.changeGitUserConfig("", true)
}
func (cgu *ChangeGitUser) changeGitUserConfig(path string, globally bool) {
	cgu.GitCommands.ChangeGitConfig("user.name", cgu.Parameters.Name, globally, cgu.Parameters.Unset, path)
	cgu.GitCommands.ChangeGitConfig("user.email", cgu.Parameters.Email, globally, cgu.Parameters.Unset, path)
}

func (cgu *ChangeGitUser) ChangeUserConfigAndRemotes(gitRepoPath string) {
	remotes := cgu.GitCommands.GetGitRemoteUrls(gitRepoPath)
	remotesToUrls := gitRemotesToMap(remotes)

	for key, value := range remotesToUrls {
		cgu.GitCommands.ChangeGitRemoteUrl(cgu.Parameters.User, cgu.Parameters.Token, value, key, gitRepoPath)
		if !cgu.Parameters.Global {
			cgu.changeGitUserConfig(gitRepoPath, false)
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
