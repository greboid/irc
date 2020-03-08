package bot

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha1"
	"crypto/subtle"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

type WebhookReceiver struct {
}

type UserAccount struct {
	Login             string `json:"login"`
	ID                int    `json:"id"`
	NodeID            string `json:"node_id"`
	AvatarURL         string `json:"avatar_url"`
	GravatarID        string `json:"gravatar_id"`
	URL               string `json:"url"`
	HTMLURL           string `json:"html_url"`
	FollowersURL      string `json:"followers_url"`
	FollowingURL      string `json:"following_url"`
	GistsURL          string `json:"gists_url"`
	StarredURL        string `json:"starred_url"`
	SubscriptionsURL  string `json:"subscriptions_url"`
	OrganizationsURL  string `json:"organizations_url"`
	ReposURL          string `json:"repos_url"`
	EventsURL         string `json:"events_url"`
	ReceivedEventsURL string `json:"received_events_url"`
	Type              string `json:"type"`
	SiteAdmin         bool   `json:"site_admin"`
}

type Author struct {
	Name  string    `json:"name"`
	Email string    `json:"email"`
	Date  time.Time `json:"date"`
}

type GitCommit struct {
	Sha     string `json:"sha"`
	URL     string `json:"url"`
	HTMLURL string `json:"html_url"`
}

type Verification struct {
	Verified  bool   `json:"verified"`
	Reason    string `json:"reason"`
	Signature string `json:"signature"`
	Payload   string `json:"payload"`
}

type GithubCommit struct {
	Author       Author       `json:"author"`
	Committer    Author       `json:"committer"`
	Message      string       `json:"message"`
	Tree         GitCommit    `json:"tree"`
	URL          string       `json:"url"`
	CommentCount int          `json:"comment_count"`
	Verification Verification `json:"verification"`
}

type Commit struct {
	Sha         string       `json:"sha"`
	NodeID      string       `json:"node_id"`
	Commit      GithubCommit `json:"commit"`
	URL         string       `json:"url"`
	HTMLURL     string       `json:"html_url"`
	CommentsURL string       `json:"comments_url"`
	Author      UserAccount  `json:"author"`
	Committer   UserAccount  `json:"committer"`
	Parents     []GitCommit  `json:"parents"`
}

type Branch struct {
	Name      string    `json:"name"`
	Commit    GitCommit `json:"commit"`
	Protected bool      `json:"protected"`
}

type License struct {
	Key    string `json:"key"`
	Name   string `json:"name"`
	SpdxID string `json:"spdx_id"`
	URL    string `json:"url"`
	NodeID string `json:"node_id"`
}

type Repository struct {
	ID               int         `json:"id"`
	NodeID           string      `json:"node_id"`
	Name             string      `json:"name"`
	FullName         string      `json:"full_name"`
	Private          bool        `json:"private"`
	Owner            UserAccount `json:"owner"`
	HTMLURL          string      `json:"html_url"`
	Description      string      `json:"description"`
	Fork             bool        `json:"fork"`
	URL              string      `json:"url"`
	ForksURL         string      `json:"forks_url"`
	KeysURL          string      `json:"keys_url"`
	CollaboratorsURL string      `json:"collaborators_url"`
	TeamsURL         string      `json:"teams_url"`
	HooksURL         string      `json:"hooks_url"`
	IssueEventsURL   string      `json:"issue_events_url"`
	EventsURL        string      `json:"events_url"`
	AssigneesURL     string      `json:"assignees_url"`
	BranchesURL      string      `json:"branches_url"`
	TagsURL          string      `json:"tags_url"`
	BlobsURL         string      `json:"blobs_url"`
	GitTagsURL       string      `json:"git_tags_url"`
	GitRefsURL       string      `json:"git_refs_url"`
	TreesURL         string      `json:"trees_url"`
	StatusesURL      string      `json:"statuses_url"`
	LanguagesURL     string      `json:"languages_url"`
	StargazersURL    string      `json:"stargazers_url"`
	ContributorsURL  string      `json:"contributors_url"`
	SubscribersURL   string      `json:"subscribers_url"`
	SubscriptionURL  string      `json:"subscription_url"`
	CommitsURL       string      `json:"commits_url"`
	GitCommitsURL    string      `json:"git_commits_url"`
	CommentsURL      string      `json:"comments_url"`
	IssueCommentURL  string      `json:"issue_comment_url"`
	ContentsURL      string      `json:"contents_url"`
	CompareURL       string      `json:"compare_url"`
	MergesURL        string      `json:"merges_url"`
	ArchiveURL       string      `json:"archive_url"`
	DownloadsURL     string      `json:"downloads_url"`
	IssuesURL        string      `json:"issues_url"`
	PullsURL         string      `json:"pulls_url"`
	MilestonesURL    string      `json:"milestones_url"`
	NotificationsURL string      `json:"notifications_url"`
	LabelsURL        string      `json:"labels_url"`
	ReleasesURL      string      `json:"releases_url"`
	DeploymentsURL   string      `json:"deployments_url"`
	UpdatedAt        time.Time   `json:"updated_at"`
	GitURL           string      `json:"git_url"`
	SSHURL           string      `json:"ssh_url"`
	CloneURL         string      `json:"clone_url"`
	SvnURL           string      `json:"svn_url"`
	Homepage         *string     `json:"homepage"`
	Size             int         `json:"size"`
	StargazersCount  int         `json:"stargazers_count"`
	WatchersCount    int         `json:"watchers_count"`
	Language         string      `json:"language"`
	HasIssues        bool        `json:"has_issues"`
	HasProjects      bool        `json:"has_projects"`
	HasDownloads     bool        `json:"has_downloads"`
	HasWiki          bool        `json:"has_wiki"`
	HasPages         bool        `json:"has_pages"`
	ForksCount       int         `json:"forks_count"`
	MirrorURL        *string     `json:"mirror_url"`
	Archived         bool        `json:"archived"`
	Disabled         bool        `json:"disabled"`
	OpenIssuesCount  int         `json:"open_issues_count"`
	License          License     `json:"license"`
	Forks            int         `json:"forks"`
	OpenIssues       int         `json:"open_issues"`
	Watchers         int         `json:"watchers"`
	DefaultBranch    string      `json:"default_branch"`
}

type GitHubCommitHook struct {
	ID          int64       `json:"id"`
	Sha         string      `json:"sha"`
	Name        string      `json:"name"`
	TargetURL   string      `json:"target_url"`
	AvatarURL   string      `json:"avatar_url"`
	Context     string      `json:"context"`
	Description string      `json:"description"`
	State       string      `json:"state"`
	Commit      Commit      `json:"commit"`
	Branches    []Branch    `json:"branches"`
	CreatedAt   time.Time   `json:"created_at"`
	UpdatedAt   time.Time   `json:"updated_at"`
	Repository  Repository  `json:"repository"`
	Sender      UserAccount `json:"sender"`
}

// payload is the raw body received, signature is the X-GitHub-Delivery header, *githubSecret is the shared secret as text
func CheckGithubSecret(r *http.Request, githubSecret string) bool {
	//Get body, put it back
	bodyBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return false
	}
	r.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
	//Get signature
	gotHash := strings.SplitN(r.Header.Get("X-Hub-Signature"), "=", 2)
	if gotHash[0] != "sha1" {
		return false
	}
	h := hmac.New(sha1.New, []byte(githubSecret))
	h.Write(bodyBytes)
	expected := fmt.Sprintf("sha1=%s", hex.EncodeToString(h.Sum(nil)))
	return len(expected) == len(gotHash[1]) && subtle.ConstantTimeCompare([]byte(expected), []byte(gotHash[1])) == 1
}