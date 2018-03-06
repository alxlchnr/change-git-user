package git

import (
	"reflect"
	"testing"
)

type mockCommand struct {
	Dir string
}

var outputMock = func() ([]byte, error) {
	return nil, nil
}

func (mCmd *mockCommand) Output() ([]byte, error) {
	return outputMock()
}

func (mCmd *mockCommand) SetWorkingDir(dir string) {
	mCmd.Dir = dir
}

var mockExecCommand = func(t *testing.T, mockedOutput string, expectedPath string, expectedName string, expectedArgs ...string) {
	execCommand = func(name string, arg ...string) Command {
		if name != expectedName {
			t.Errorf("Expected command with name %v, got %v", expectedName, name)
		}
		if !reflect.DeepEqual(expectedArgs, arg) {
			t.Errorf("Expected arguments %v, got %v", expectedArgs, arg)
		}
		cmd := &mockCommand{}
		outputMock = func() ([]byte, error) {
			if cmd.Dir != expectedPath {
				t.Errorf("Expected working directory %v, got %v", expectedPath, cmd.Dir)
			}
			return []byte(mockedOutput), nil
		}
		return Command(cmd)
	}
}

func TestGitCommands_ChangeGitConfig(t *testing.T) {
	old := execCommand
	defer func() { execCommand = old }()

	type args struct {
		config   string
		value    string
		globally bool
		unset    bool
		path     string
	}
	tests := []struct {
		name                string
		args                args
		expectedCommandArgs []string
	}{
		{
			name: "should set config to new value",
			args: args{
				config:   "user.email",
				value:    "new@gmail.com",
				globally: false,
				unset:    false,
				path:     "aLocalPath",
			},
			expectedCommandArgs: []string{"config", "user.email", "new@gmail.com"},
		}, {
			name: "should set config to new value globally",
			args: args{
				config:   "user.email",
				value:    "new@gmail.com",
				globally: true,
				unset:    false,
				path:     "",
			},
			expectedCommandArgs: []string{"config", "user.email", "new@gmail.com", "--global"},
		}, {
			name: "should unset config to new value globally",
			args: args{
				config:   "user.email",
				value:    "new@gmail.com",
				globally: true,
				unset:    true,
				path:     "",
			},
			expectedCommandArgs: []string{"config", "--unset", "user.email", "--global"},
		}, {
			name: "should unset config to new value",
			args: args{
				config:   "user.email",
				value:    "new@gmail.com",
				globally: false,
				unset:    true,
				path:     "aLocalPath",
			},
			expectedCommandArgs: []string{"config", "--unset", "user.email"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &GitCommandsImpl{}
			mockExecCommand(t, "it worked", tt.args.path, "git", tt.expectedCommandArgs...)
			g.ChangeGitConfig(tt.args.config, tt.args.value, tt.args.globally, tt.args.unset, tt.args.path)
		})
	}
}

func TestGitCommands_GetGitRemoteUrls(t *testing.T) {
	type args struct {
		gitRepoPath string
	}
	tests := []struct {
		name                  string
		args                  args
		want                  string
		expectedCommandArgs   []string
		expectedCommandOutput string
	}{
		{
			name:                  "get remote URLs",
			args:                  args{gitRepoPath: "aLocalGitRepo"},
			want:                  "expectedOutput",
			expectedCommandArgs:   []string{"remote", "-v"},
			expectedCommandOutput: "expectedOutput",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &GitCommandsImpl{}
			mockExecCommand(t, tt.expectedCommandOutput, tt.args.gitRepoPath, "git", tt.expectedCommandArgs...)
			if got := g.GetGitRemoteUrls(tt.args.gitRepoPath); got != tt.want {
				t.Errorf("GitCommands.GetGitRemoteUrls() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGitCommands_ChangeGitRemoteUrl(t *testing.T) {
	type args struct {
		newUser     string
		newToken    string
		remoteUrl   *GitUrl
		remote      string
		gitRepoPath string
	}
	tests := []struct {
		name                string
		args                args
		expectedCommandArgs []string
	}{
		{
			name: "get remote URLs",
			args: args{
				gitRepoPath: "aLocalGitRepo",
				newUser:     "aUserName",
				newToken:    "aGithubToken",
				remote:      "origin",
				remoteUrl:   NewGitUrl("https://old.github.url/abcd"),
			},
			expectedCommandArgs: []string{"remote", "set-url", "origin", "https://aUserName:aGithubToken@old.github.url/abcd"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &GitCommandsImpl{}
			mockExecCommand(t, "it worked...", tt.args.gitRepoPath, "git", tt.expectedCommandArgs...)
			g.ChangeGitRemoteUrl(tt.args.newUser, tt.args.newToken, tt.args.remoteUrl, tt.args.remote, tt.args.gitRepoPath)
		})
	}
}
