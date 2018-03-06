package git

type ChangeGitParameters struct {
	User   string
	Token  string
	Email  string
	Name   string
	Path   string
	Help   bool
	Global bool
	Unset  bool
}
