package frozen_throne_server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/bradleyfalzon/ghinstallation/v2"
	"github.com/google/go-github/github"
)

func WebhookHandler(w http.ResponseWriter, r *http.Request) {
	payload, err := github.ValidatePayload(r, []byte(serverConfig.WebhookSecret))
	if err != nil {
		log.Printf("error validating request body: err=%s\n", err)
		return
	}
	defer r.Body.Close()

	event, err := github.ParseWebHook(github.WebHookType(r), payload)
	if err != nil {
		log.Printf("could not parse webhook: err=%s\n", err)
		return
	}

	installationID := r.Header.Get("X-GitHub-Hook-Installation-Target-ID")
	installationIDInt, parseErr := strconv.ParseInt(installationID, 10, 64)
	if parseErr != nil {
		log.Printf("could not parse webhook: err=%s\n", err)
		return
	}

	// Wrap the shared transport for use with the integration ID 1 authenticating with installation ID 99.
	itr, err := ghinstallation.NewKeyFromFile(
		http.DefaultTransport, serverConfig.GithubAppID, installationIDInt, "2016-10-19.private-key.pem")
	if err != nil {
		log.Printf("could not parse webhook: err=%s\n", err)
		return
	}

	// Use installation transport with client.
	client := github.NewClient(&http.Client{Transport: itr})

	switch e := event.(type) {
	case *github.CheckSuiteEvent:
		processStatusCheck(client, *e.Org.Name, *e.Repo.Name, *e.CheckSuite.HeadSHA)
	case *github.PullRequestEvent:
		processStatusCheck(client, *e.Repo.Organization.Name, *e.Repo.Name, *e.PullRequest.Head.SHA)
	case *github.PushEvent:
		processStatusCheck(client, *e.Repo.Organization, *e.Repo.Name, *e.HeadCommit.SHA)
	default:
		//
	}
}

func processStatusCheck(client *github.Client, org string, repo string, headSHA string) StatusResponse {
	ctx := context.Background()

	_, statusErr := ft.Check(repo)

	var returnResp StatusResponse
	var status, title, text string
	if statusErr == nil {
		// If the status error is nil, that means it is frozen
		status = "in_progress"
		title = fmt.Sprintf("%s is frozen", repo)
		text = "All merges have been blocked."
		returnResp = StatusResponse{}
	} else {
		status = "completed"
		title = fmt.Sprintf("%s is not frozen", repo)
		text = "All merges are okay."
	}

	checkOptions := github.CreateCheckRunOptions{
		Name:      "frozen-throne",
		HeadSHA:   headSHA,
		Status:    &status,
		StartedAt: &github.Timestamp{Time: time.Now()},
		Output: &github.CheckRunOutput{
			Title:   &title,
			Summary: &title,
			Text:    &text,
		},
	}
	client.Checks.CreateCheckRun(ctx, org, repo, checkOptions)
	return returnResp
}
