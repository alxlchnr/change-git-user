# change-git-user
[![Go Report Card](https://goreportcard.com/badge/github.com/alxlchnr/change-git-user)](https://goreportcard.com/report/github.com/alxlchnr/change-git-user)
[![Build Status](https://travis-ci.org/alxlchnr/change-git-user.svg?branch=master)](https://travis-ci.org/alxlchnr/change-git-user)
[![codecov](https://codecov.io/gh/alxlchnr/change-git-user/branch/master/graph/badge.svg)](https://codecov.io/gh/alxlchnr/change-git-user)
[![Maintainability](https://api.codeclimate.com/v1/badges/383ef543a78b8907ba6c/maintainability)](https://codeclimate.com/github/alxlchnr/change-git-user/maintainability)
[![MIT Licence](https://badges.frapsoft.com/os/mit/mit.svg?v=103)](https://opensource.org/licenses/mit-license.php)

This small command line tool helps you to change the remote urls of your GIT repositories. (Given an GIT-URL with following schema: https://{user}:{token}@{repo-url})
Starting with a provided folder it will search for GIT repositories recursively and will apply the changes.

Additionally it can set your GIT user data globally for your computer.

## Prerequisites
You need to have GIT installed on your computer.

## Install
You can download releases of the command line tool from this github repository or if you have Go installed 
on your computer you can also get the executable by executing
        
        go get github.com/alxlchnr/change-git-user


## Possible parameters

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


If some of the parameters are not set, they will not be changed.