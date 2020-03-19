package main

import (
	"crypto/hmac"
	"crypto/sha1"
	"crypto/subtle"
	"encoding/hex"
	"fmt"
)

func CheckGithubSecret(bodyBytes []byte, headerSecret string, githubSecret string) bool {
	h := hmac.New(sha1.New, []byte(githubSecret))
	h.Write(bodyBytes)
	expected := fmt.Sprintf("%s", hex.EncodeToString(h.Sum(nil)))
	return len(expected) == len(headerSecret) && subtle.ConstantTimeCompare([]byte(expected), []byte(headerSecret)) == 1
}

//Push events
type pushhook struct {
	Refspec     string     `json:"ref"`
	Repository  Repository `json:"repository"`
	Pusher      Pusher     `json:"pusher"`
	Forced      bool       `json:"forced"`
	Deleted     bool       `json:"deleted"`
	Created     bool       `json:"created"`
	CompareLink string     `json:"compare"`
	Commits     []Commit   `json:"commits"`
	Baserefspec string     `json:"base_ref"`
}

type Repository struct {
	FullName string `json:"full_name"`
}

type Pusher struct {
	Name string `json:"name"`
}

type Commit struct {
	ID        string `json:"id"`
	Message   string `json:"message"`
	URL       string `json:"url"`
	Author    Author `json:"author"`
	Committer Author `json:"committer"`
}

type Author struct {
	User string `json:"username"`
}

//Pull Request events
type prhook struct {
	Action      string      `json:"action"`
	PullRequest PullRequest `json:"pull_request"`
	Repository Repository `json:"repository"`
}

type PullRequest struct {
	Url        string     `json:"url"`
	State      string     `json:"state"`
	Title      string     `json:"title"`
	Body       string     `json:"body"`
	User       User       `json:"user"`
}

type User struct {
	Login string `json:"login"`
}
