package frozen_throne_server

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/bradleyfalzon/ghinstallation/v2"
	"github.com/google/go-github/github"
)

func WebhookHandler(w http.ResponseWriter, r *http.Request) {
	webhookSecret := []byte(serverConfig.WebhookSecret)
	log.Printf("%v %v", serverConfig.WebhookSecret, webhookSecret)
	payload, err := github.ValidatePayload(r, webhookSecret)
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

	var resp StatusResponse
	switch e := event.(type) {
	case *github.CheckSuiteEvent:
		resp = processStatusCheck(*e.Installation.ID, *e.Org.Name, *e.Repo.Name, *e.CheckSuite.HeadSHA)
	case *github.PullRequestEvent:
		resp = processStatusCheck(*e.Installation.ID, *e.Repo.Owner.Login, *e.Repo.Name, *e.PullRequest.Head.SHA)
	case *github.PushEvent:
		resp = processStatusCheck(*e.Installation.ID, *e.Repo.Organization, *e.Repo.Name, *e.HeadCommit.SHA)
	default:
		//
	}

	json.NewEncoder(w).Encode(resp)
}

func processStatusCheck(installationID int64, org string, repo string, headSHA string) StatusResponse {
	ctx := context.Background()
	itr, err := ghinstallation.New(
		http.DefaultTransport, serverConfig.GithubAppID, installationID,
		[]byte(serverConfig.GithubAppPrivateKey))
	if err != nil {
		log.Printf("error initialising integration transport err=%s\n", err)
	}

	// Use installation transport with client.
	client := github.NewClient(&http.Client{Transport: itr})

	_, statusErr := ft.Check(repo)

	var returnResp StatusResponse
	var status, title, text, conclusion string
	conclusion = "success"
	if statusErr == nil {
		// If the status error is nil, that means it is frozen
		status = "in_progress"
		title = fmt.Sprintf("Repo \"%s\" is frozen", repo)
		text = "All merges have been blocked."
		returnResp = StatusResponse{}
	} else {
		status = "completed"
		title = fmt.Sprintf("Repo \"%s\" is not frozen", repo)
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
		Conclusion: &conclusion,
	}
	_, resp, err := client.Checks.CreateCheckRun(ctx, org, repo, checkOptions)
	log.Printf("%v | %v", err, resp)
	return returnResp
}
