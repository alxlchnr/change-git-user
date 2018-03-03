# change-git-user
[![Go Report Card](https://goreportcard.com/badge/github.com/alxlchnr/change-git-user?style=flat-square)](https://goreportcard.com/report/github.com/alxlchnr/change-git-user)
[![Build Status](https://travis-ci.org/alxlchnr/change-git-user.svg?branch=master)](https://travis-ci.org/alxlchnr/change-git-user)

This small command line tool helps you to change the remote urls of your GIT repositories. (Given an GIT-URL with following schema: https://{user}:{token}@{repo-url})
Additionally it can set your GIT user data globally for your computer.

## Possible parameters
<pre><code>
  -email string
      	the email of the new user
    -global
      	apply user name and email globally (default true)
    -help
      	show help
    -name string
      	the name of the new user
    -path string
      	path where to look for git repositories (default ".")
    -token string
      	the API token of the new user
    -unset
      	unset user name and email
    -user string
      	the API username of the new user
</code></pre>

If some of the parameters are not set, they will not be changed.

**You may have to mark the tool as executable on your computer before first usage**
