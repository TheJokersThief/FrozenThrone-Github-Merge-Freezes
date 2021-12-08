[![CircleCI](https://circleci.com/gh/TheJokersThief/frozen-throne/tree/main.svg?style=svg)](https://circleci.com/gh/TheJokersThief/frozen-throne/tree/main)

# Frozen Throne (Merge Freezes)
FrozenThrone is an API deployed to GCP Cloud functions for gating PR merges on Github.

![image](https://user-images.githubusercontent.com/1175876/145129803-ce719ddc-f8ba-4c90-a5eb-90dd1d116965.png)

# API


# Deployment
## Create Secrets
The deployed cloud function uses the GCP Secret Manager to store secret values for:

* The Write Secret token
* The Read-Only Secret token

```bash
PROJECT_ID=<ID> make create_secrets
```

Or update existing secrets with

```bash
PROJECT_ID=<ID> \
WRITE_SECRET=<secret> \
READ_ONLY_SECRET=<secret> \
    make create_secrets
```

## Deploy to Cloud Run

```bash
PROJECT_ID=<ID> make deploy
```
