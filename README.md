[![CircleCI](https://circleci.com/gh/TheJokersThief/frozen-throne/tree/main.svg?style=svg)](https://circleci.com/gh/TheJokersThief/frozen-throne/tree/main)

# Frozen Throne (Merge Freezes)
Frozen Throne is an API deployed to GCP Cloud Run for gating PR merges on Github.

![image](https://user-images.githubusercontent.com/1175876/145129803-ce719ddc-f8ba-4c90-a5eb-90dd1d116965.png)

<!-- TOC -->

- [Frozen Throne Merge Freezes](#frozen-throne-merge-freezes)
- [API](#api)
    - [Authentication](#authentication)
    - [Example requests](#example-requests)
- [Deployment](#deployment)
    - [Pre-requisites](#pre-requisites)
    - [Create Secrets](#create-secrets)
    - [Deploy to Cloud Run](#deploy-to-cloud-run)

<!-- /TOC -->

# API

| Endpoint         | Description              | POST data            |
|------------------|--------------------------|----------------------|
| /freeze/{repo}   | Freeze the github {repo} | `user`               |
| /unfreeze/{repo} | Unfreezes a repo         | `user`               |
| /github-webhook  |                          | github webhook event |

## Authentication

The freeze and unfreeze endpoints both require authentication in the form of a header in the request.

```
X-Access-Token: WRITE_SECRET
```

## Example requests

```bash
$ curl -X POST -H "X-Access-Token: SECRET" localhost:8080/freeze/frozen-throne -d "user=thejokersthief"
{"frozen":true}

$ curl -X POST -H "X-Access-Token: SECRET" localhost:8080/unfreeze/frozen-throne -d "user=thejokersthief"
{"frozen":false}
```

# Deployment

## Pre-requisites

* Go 1.16
* [Have created a Github App](https://docs.github.com/en/developers/apps/building-github-apps/creating-a-github-app) and have noted the Github App ID, and [have generated a private key](https://docs.github.com/en/developers/apps/building-github-apps/authenticating-with-github-apps).
* Have generated a secret for both your webhook verification and a write-access API key (`openssl rand -base64 48`)

## Create Secrets
The deployed cloud function uses the GCP Secret Manager to store secret values for:

1. The Write Secret token
1. The Read-Only Secret token
1. The secret used to sign webhooks from Github
1. The Github App ID
1. The Github App's private key

The first 4 of these can be created with the following command:

```bash
PROJECT_ID=<ID> \
WRITE_SECRET=<secret> \
WEBHOOK_SECRET=<secret> \
GITHUB_APP_ID=<secret int> \
    make create_secrets
```

And you can _update_ the secrets by using the same command, but replacing `create_secrets` with `update_secrets`.

The final secret is a private key associated with the Github app. This is a `.pem` file and can be added with the following command:

```bash
export PROJECT_ID="example"
export PATH_TO_PEM_FILE="some/file/path"
gcloud --project ${PROJECT_ID} secrets create FT_GITHUB_PRIVATE_KEY --replication-policy="automatic" --data-file=${PATH_TO_PEM_FILE}
```

## Deploy to Cloud Run

Now that you've got all your secrets set up, you are good to deploy to Cloud Run. This involves two stages:

1. Build a cloud image
2. Deploy the image to Cloud Run

```bash
make build
PROJECT_ID=<ID> make deploy
```
