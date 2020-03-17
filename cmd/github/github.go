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

type pushhook struct {
	Refspec     string     `json:"ref"`
	Repository  Repository `json:"repository"`
	Pusher      Pusher     `json:"pusher"`
	Forced      bool       `json:"forced"`
	Deleted     bool       `json:"deleted"`
	Created     bool       `json:"created"`
	CompareLink string     `json:"compare"`
	Commits     []Commit   `json:"commits"`
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
