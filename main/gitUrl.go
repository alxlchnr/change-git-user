package main

import "strings"

type GitUrl struct {
	Protocol string
	User     string
	Token    string
	Repo     string
}

func (g *GitUrl) SetUser(newUser string) {
	if len(newUser) > 0 {
		g.User = newUser
	}
}

func (g *GitUrl) SetToken(newToken string) {
	if len(newToken) > 0 {
		g.Token = newToken
	}
}

func (g *GitUrl) ToUrl() string {
	if len(g.Token) > 0 && len(g.User) > 0 {
		return g.Protocol + "//" + g.User + ":" + g.Token + "@" + g.Repo
	}
	if len(g.User) > 0 {
		return g.Protocol + "//" + g.User + "@" + g.Repo
	}
	return g.Protocol + "//" + g.Repo
}

func NewGitUrl(url string) *GitUrl {
	splittetAfterProtocol := strings.Split(url, "//")
	splittedAfterCredentials := strings.Split(splittetAfterProtocol[1], "@")
	credentialsSplitted := strings.Split(splittedAfterCredentials[0], ":")
	if len(credentialsSplitted) == 2 {
		return &GitUrl{
			splittetAfterProtocol[0],
			credentialsSplitted[0],
			credentialsSplitted[1],
			splittedAfterCredentials[1],
		}
	} else {
		if len(splittedAfterCredentials) == 2 {

			return &GitUrl{
				splittetAfterProtocol[0],
				splittedAfterCredentials[0],
				"",
				splittedAfterCredentials[1],
			}
		} else {
			return &GitUrl{
				splittetAfterProtocol[0],
				"",
				"",
				splittedAfterCredentials[0],
			}
		}
	}
}
