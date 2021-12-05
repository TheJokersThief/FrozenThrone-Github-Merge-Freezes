[![CircleCI](https://circleci.com/gh/TheJokersThief/frozen-throne/tree/main.svg?style=svg)](https://circleci.com/gh/TheJokersThief/frozen-throne/tree/main)

# Frozen Throne (Merge Freezes)
FrozenThrone is an API deployed to GCP Cloud functions for gating PR merges on Github.

# API

| Method | Description                             | Params            |
|--------|-----------------------------------------|-------------------|
| GET    | Retrieves the current status for a repo | token, repo       |
| POST   | Freezes a repo                          | token, repo, user |
| PATCH  | Unfreezes a repo                        | token, repo, user |

`GET` params are passed via the URL.

`POST`/`PATCH` params are passed via HTTP form data.

# Deployment
## Create Secrets
The deployed cloud function uses the GCP Secret Manager to store secret values for:

* The Write Secret token
* The Read-Only Secret token

```
PROJECT_ID=<ID> make create_secrets
```

Or update existing secrets with

```
PROJECT_ID=<ID> \
WRITE_SECRET=<secret> \
READ_ONLY_SECRET=<secret>\
    make create_secrets
```

## Deploy function

```
PROJECT_ID=<ID> make deploy_to_gfunctions
```