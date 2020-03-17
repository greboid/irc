package main

import (
	"context"
	"crypto/hmac"
	"crypto/sha1"
	"crypto/subtle"
	"crypto/tls"
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/greboid/irc/rpc"
	"github.com/kouhin/envflag"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"time"
)

var (
	RPCHost      = flag.String("rpc-host", "localhost", "gRPC server to connect to")
	RPCPort      = flag.Int("rpc-port", 8001, "gRPC server port")
	RPCToken     = flag.String("rpc-token", "", "gRPC authentication token")
	Channel      = flag.String("channel", "", "Channel to send messages to")
	GithubSecret = flag.String("github-secret", "", "Github secret for validating webhooks")
	WebPort      = flag.Int("web-port", 8000, "Web port for receiving github webhooks")
)

type github struct {
	client rpc.IRCPluginClient
}

func main() {
	if err := envflag.Parse(); err != nil {
		log.Fatalf("Unable to load config: %s", err.Error())
	}
	github := github{}
	log.Printf("Creating Github RPC Client")
	client, err := github.doRPC()
	if err != nil {
		log.Fatalf("Unable to create RPC Client: %s", err.Error())
	}
	github.client = client
	log.Printf("Starting github web server")
	github.doWeb()

}

func (g *github) doRPC() (rpc.IRCPluginClient, error) {
	creds := credentials.NewTLS(&tls.Config{InsecureSkipVerify: true})
	conn, err := grpc.Dial(fmt.Sprintf("%s:%d", *RPCHost, *RPCPort), grpc.WithTransportCredentials(creds))
	client := rpc.NewIRCPluginClient(conn)
	_, err = client.Ping(rpc.CtxWithToken(context.Background(), "bearer", *RPCToken), &rpc.Empty{})
	if err != nil {
		return nil, err
	}
	return client, nil
}

func (g *github) doWeb() {
	mux := http.NewServeMux()
	mux.HandleFunc("/github", g.handleGithub)
	server := http.Server{
		Addr:    fmt.Sprintf(":%d", *WebPort),
		Handler: mux,
	}
	go func() {
		log.Print(server.ListenAndServe())
	}()
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, os.Kill)
	log.Printf("Waiting for stop")
	<-stop
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Unable to shutdown: %s", err.Error())
	}
}

func (g *github) handleGithub(writer http.ResponseWriter, request *http.Request) {
	bodyBytes, err := ioutil.ReadAll(request.Body)
	defer func() { _ = request.Body.Close() }()
	if err != nil {
		log.Printf("Error reading body: %s", err.Error())
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	header := strings.SplitN(request.Header.Get("X-Hub-Signature"), "=", 2)
	if header[0] != "sha1" {
		log.Printf("Error: %s", "Bad header")
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	if !CheckGithubSecret(bodyBytes, header[1], *GithubSecret) {
		log.Printf("Error: %s", "Bad hash")
		writer.WriteHeader(http.StatusBadRequest)
	}
	data := GitHubCommitHook{}
	err = json.Unmarshal(bodyBytes, &data)
	if err != nil {
		log.Printf("Error unmarshalling: %s", err.Error())
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	log.Printf("Sending %s: %s", *Channel, fmt.Sprintf("%s %s", data.Repository.Name, data.Commit.Author.Login))
	_, err = g.client.SendChannelMessage(rpc.CtxWithToken(context.Background(), "bearer", *RPCToken), &rpc.ChannelMessage{
		Channel: *Channel,
		Message: fmt.Sprintf("%s %s", data.Repository.Name, data.Commit.Author.Login),
	})
	if err != nil {
		log.Printf("Error sending to channel: %s", err.Error())
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	_, _ = writer.Write([]byte("Delivered."))
	pretty, err := json.MarshalIndent(data, "", "\t")
	log.Printf("%s", string(pretty))
}

func CheckGithubSecret(bodyBytes []byte, headerSecret string, githubSecret string) bool {
	h := hmac.New(sha1.New, []byte(githubSecret))
	h.Write(bodyBytes)
	expected := fmt.Sprintf("%s", hex.EncodeToString(h.Sum(nil)))
	return len(expected) == len(headerSecret) && subtle.ConstantTimeCompare([]byte(expected), []byte(headerSecret)) == 1
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