package cmd

import (
	"github.com/golang/mock/gomock"
	"reflect"
	"testing"
)

var exampleGitRemoteOutput = "origin\thttps://github.com/alxlchnr/change-git-user.git (fetch)\norigin\thttps://github.com/alxlchnr/change-git-user.git (push)"
var exampleGitRemoteOutput_2Remotes = "origin\thttps://github.com/alxlchnr/change-git-user.git (fetch)\norigin2\thttps://github.com/alxlchnr/change-git-user2.git (push)"

func TestChangeGitUser_ChangeGlobalUserConfig(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockGitCommands := NewMockGitCommands(mockCtrl)
	mockGitCommands.EXPECT().ChangeGitConfig("user.name", "newName", true, false, "").Times(1)
	mockGitCommands.EXPECT().ChangeGitConfig("user.email", "newMail", true, false, "").Times(1)

	parameters := &ChangeGitParameters{Name: "newName", Email: "newMail", Global: true, Unset: false, Path: "somePath"}

	sut := &ChangeGitUser{Parameters: parameters, GitCommands: mockGitCommands}

	sut.ChangeGlobalUserConfig()
}

func TestChangeGitUser_ChangeGlobalUserConfig__Unset(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockGitCommands := NewMockGitCommands(mockCtrl)
	mockGitCommands.EXPECT().ChangeGitConfig("user.name", "newName", true, true, "").Times(1)
	mockGitCommands.EXPECT().ChangeGitConfig("user.email", "newMail", true, true, "").Times(1)

	parameters := &ChangeGitParameters{Name: "newName", Email: "newMail", Global: false, Unset: true, Path: "somePath"}

	sut := &ChangeGitUser{Parameters: parameters, GitCommands: mockGitCommands}

	sut.ChangeGlobalUserConfig()
}

func TestChangeGitUser_changeGitUserConfig(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockGitCommands := NewMockGitCommands(mockCtrl)
	mockGitCommands.EXPECT().ChangeGitConfig("user.name", "newName", false, false, "somePath").Times(1)
	mockGitCommands.EXPECT().ChangeGitConfig("user.email", "newMail", false, false, "somePath").Times(1)

	parameters := &ChangeGitParameters{Name: "newName", Email: "newMail", Global: true, Unset: false, Path: "somePath"}

	sut := &ChangeGitUser{Parameters: parameters, GitCommands: mockGitCommands}

	sut.changeGitUserConfig("somePath", false)
}

func TestChangeGitUser_changeGitUserConfig__Unset(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockGitCommands := NewMockGitCommands(mockCtrl)
	mockGitCommands.EXPECT().ChangeGitConfig("user.name", "newName", false, true, "somePath").Times(1)
	mockGitCommands.EXPECT().ChangeGitConfig("user.email", "newMail", false, true, "somePath").Times(1)

	parameters := &ChangeGitParameters{Name: "newName", Email: "newMail", Global: true, Unset: true, Path: "somePath"}

	sut := &ChangeGitUser{Parameters: parameters, GitCommands: mockGitCommands}

	sut.changeGitUserConfig("somePath", false)
}

func TestChangeGitUser_ChangeUserConfigAndRemotes(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockGitCommands := NewMockGitCommands(mockCtrl)
	mockGitCommands.EXPECT().ChangeGitConfig("user.name", "newName", false, false, "someRepoPath").Times(1)
	mockGitCommands.EXPECT().ChangeGitConfig("user.email", "newMail", false, false, "someRepoPath").Times(1)
	mockGitCommands.EXPECT().ChangeGitRemoteUrl("newUser", "newToken", &GitUrl{Protocol: "https:", Repo: "github.com/alxlchnr/change-git-user.git"}, "origin", "someRepoPath").Times(1)
	mockGitCommands.EXPECT().GetGitRemoteUrls("someRepoPath").Return(exampleGitRemoteOutput).Times(1)

	parameters := &ChangeGitParameters{
		Name:   "newName",
		Email:  "newMail",
		User:   "newUser",
		Token:  "newToken",
		Global: false,
		Unset:  false,
		Path:   "somePath"}

	sut := &ChangeGitUser{Parameters: parameters, GitCommands: mockGitCommands}

	sut.ChangeUserConfigAndRemotes("someRepoPath")
}

func Test_gitRemotesToMap(t *testing.T) {
	oneRemote := make(map[string]*GitUrl)
	oneRemote["origin"] = &GitUrl{Protocol: "https:", Repo: "github.com/alxlchnr/change-git-user.git"}
	twoRemotes := make(map[string]*GitUrl)
	twoRemotes["origin"] = &GitUrl{Protocol: "https:", Repo: "github.com/alxlchnr/change-git-user.git"}
	twoRemotes["origin2"] = &GitUrl{Protocol: "https:", Repo: "github.com/alxlchnr/change-git-user2.git"}
	type args struct {
		remotes string
	}
	tests := []struct {
		name string
		args args
		want map[string]*GitUrl
	}{
		{
			name: "fetch & push for same remote",
			args: args{remotes: exampleGitRemoteOutput},
			want: oneRemote,
		},
		{
			name: "more than one remote",
			args: args{remotes: exampleGitRemoteOutput_2Remotes},
			want: twoRemotes,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := gitRemotesToMap(tt.args.remotes); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("gitRemotesToMap() = %v, want %v", got, tt.want)
			}
		})
	}
}
