package main

import (
	"reflect"
	"testing"
)

func TestGitUrl_SetUser(t *testing.T) {
	type fields struct {
		Protocol string
		User     string
		Token    string
		Repo     string
	}
	type args struct {
		newUser string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *GitUrl
	}{
		{
			name:   "no new value",
			fields: fields{User: "oldToken"},
			args:   args{newUser: ""},
			want:   &GitUrl{User: "oldToken"},
		}, {
			name:   "new value",
			fields: fields{User: "oldToken"},
			args:   args{newUser: "newToken"},
			want:   &GitUrl{User: "newToken"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &GitUrl{
				Protocol: tt.fields.Protocol,
				User:     tt.fields.User,
				Token:    tt.fields.Token,
				Repo:     tt.fields.Repo,
			}
			g.SetUser(tt.args.newUser)
			if !reflect.DeepEqual(g, tt.want) {
				t.Errorf("SetUser(%v) = %v, want %v", tt.args.newUser, g, tt.want)
			}
		})
	}
}

func TestGitUrl_SetToken(t *testing.T) {
	type fields struct {
		Protocol string
		User     string
		Token    string
		Repo     string
	}
	type args struct {
		newToken string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *GitUrl
	}{
		{
			name:   "no new value",
			fields: fields{Token: "oldToken"},
			args:   args{newToken: ""},
			want:   &GitUrl{Token: "oldToken"},
		}, {
			name:   "new value",
			fields: fields{Token: "oldToken"},
			args:   args{newToken: "newToken"},
			want:   &GitUrl{Token: "newToken"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &GitUrl{
				Protocol: tt.fields.Protocol,
				User:     tt.fields.User,
				Token:    tt.fields.Token,
				Repo:     tt.fields.Repo,
			}
			g.SetToken(tt.args.newToken)
			if !reflect.DeepEqual(g, tt.want) {
				t.Errorf("SetToken(%v) = %v, want %v", tt.args.newToken, g, tt.want)
			}
		})
	}
}

func TestGitUrl_ToUrl(t *testing.T) {
	type fields struct {
		Protocol string
		User     string
		Token    string
		Repo     string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name:   "no credentials",
			fields: fields{Protocol: "http:", User: "", Token: "", Repo: "github.com/xyz"},
			want:   "http://github.com/xyz",
		}, {
			name:   "no user",
			fields: fields{Protocol: "http:", User: "", Token: "aToken", Repo: "github.com/xyz"},
			want:   "http://github.com/xyz",
		}, {
			name:   "only user",
			fields: fields{Protocol: "http:", User: "aUserName", Token: "", Repo: "github.com/xyz"},
			want:   "http://aUserName@github.com/xyz",
		}, {
			name:   "with credentials",
			fields: fields{Protocol: "http:", User: "aUserName", Token: "aToken", Repo: "github.com/xyz"},
			want:   "http://aUserName:aToken@github.com/xyz",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &GitUrl{
				Protocol: tt.fields.Protocol,
				User:     tt.fields.User,
				Token:    tt.fields.Token,
				Repo:     tt.fields.Repo,
			}
			if got := g.ToUrl(); got != tt.want {
				t.Errorf("GitUrl.ToUrl() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewGitUrl(t *testing.T) {
	type args struct {
		url string
	}
	tests := []struct {
		name string
		args args
		want *GitUrl
	}{
		{
			name: "No user and token in url",
			args: args{url: "https://github.com/xyz"},
			want: &GitUrl{Protocol: "https:", Repo: "github.com/xyz"},
		}, {
			name: "No user and token in url,but @ in url",
			args: args{url: "https://@github.com/xyz"},
			want: &GitUrl{Protocol: "https:", Repo: "github.com/xyz"},
		}, {
			name: "No token in url",
			args: args{url: "https://aUserName@github.com/xyz"},
			want: &GitUrl{Protocol: "https:", User: "aUserName", Repo: "github.com/xyz"},
		}, {
			name: "No user in url",
			args: args{url: "https://:aToken@github.com/xyz"},
			want: &GitUrl{Protocol: "https:", User: "", Token: "aToken", Repo: "github.com/xyz"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewGitUrl(tt.args.url); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewGitUrl() = %v, want %v", got, tt.want)
			}
		})
	}
}
