PROJECT_NAME = frozen_throne
PROJECT_ID ?= example-project
GCS_BUCKET ?= ${PROJECT_NAME}-test-bucket

WRITE_SECRET ?= secret
READ_ONLY_SECRET ?= secret-read-only
WEBHOOK_SECRET ?= secretysecret
GITHUB_APP_ID ?= 1

.PHONY: help
help: ## Show help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'


.PHONY: build-linux
build-linux: ## Build go binary for linux
	GOOS=linux GOARCH=amd64 go build -o bin/linux/${PROJECT_NAME} ./cmd/run_frozen_throne.go

.PHONY: build-darwin
build-darwin: ## Build go binary for mac OS
	GOOS=darwin GOARCH=amd64 go build -o bin/darwin/${PROJECT_NAME} ./cmd/run_frozen_throne.go

.PHONY: run
run:
	WRITE_SECRET=${WRITE_SECRET} \
	READ_ONLY_SECRET=${READ_ONLY_SECRET} \
	WEBHOOK_SECRET=${WEBHOOK_SECRET} \
	GITHUB_APP_ID=${GITHUB_APP_ID} \
	./bin/darwin/frozen_throne

.PHONY: build
build: build-linux build-darwin ## Build all binaries

.PHONY: create_secrets
create_secrets: ## Create secret values
	echo -n "${WRITE_SECRET}" | gcloud --project ${PROJECT_ID} secrets create FT_WRITE_SECRET --replication-policy="automatic" --data-file=-
	echo -n "${READ_ONLY_SECRET}" | gcloud --project ${PROJECT_ID} secrets create FT_READ_ONLY_SECRET --replication-policy="automatic" --data-file=-
	echo -n "${WEBHOOK_SECRET}" | gcloud --project ${PROJECT_ID} secrets create FT_WEBHOOK_SECRET --replication-policy="automatic" --data-file=-
	echo -n "${GITHUB_APP_ID}" | gcloud --project ${PROJECT_ID} secrets create FT_GITHUB_APP_ID --replication-policy="automatic" --data-file=-
# echo -n "${GITHUB_PRIVATE_KEY}" | gcloud --project ${PROJECT_ID} secrets create FT_GITHUB_PRIVATE_KEY --replication-policy="automatic" --data-file=-


.PHONY: update_secrets
update_secrets: ## Update secret values
	echo "${WRITE_SECRET}" | tr -d \\n | gcloud --project ${PROJECT_ID} secrets versions add FT_WRITE_SECRET --data-file=-
	echo "${READ_ONLY_SECRET}" | tr -d \\n | gcloud --project ${PROJECT_ID} secrets versions add FT_READ_ONLY_SECRET --data-file=-
	echo "${WEBHOOK_SECRET}" | tr -d \\n | gcloud --project ${PROJECT_ID} secrets versions add FT_WEBHOOK_SECRET --data-file=-
	echo "${GITHUB_APP_ID}" | tr -d \\n | gcloud --project ${PROJECT_ID} secrets versions add FT_GITHUB_APP_ID --data-file=-
# echo "${GITHUB_PRIVATE_KEY}" | tr -d \\n | gcloud --project ${PROJECT_ID} secrets versions add FT_GITHUB_PRIVATE_KEY --data-file=-

.PHONY: cloud_build
cloud_build: ## Build image and push it to GCR
	gcloud builds submit \
		--project ${PROJECT_ID} \
		--config cloudbuild.yaml \
		--substitutions _PROJECT_NAME=${PROJECT_NAME},_TAG=latest

.PHONY: deploy
deploy: cloud_build ## Deploy function to GCP Cloud Functions
	gcloud run deploy frozen-throne \
		--region europe-west1 \
		--project ${PROJECT_ID} \
		--memory 128Mi \
		--timeout 20s \
		--set-env-vars GOOGLE_CLOUD_PROJECT=${PROJECT_ID},GCS_BUCKET="${GCS_BUCKET}" \
		--set-secrets 'WRITE_SECRET=FT_WRITE_SECRET:latest,READ_ONLY_SECRET=FT_READ_ONLY_SECRET:latest,WEBHOOK_SECRET=FT_WEBHOOK_SECRET:latest,GITHUB_APP_ID=FT_GITHUB_APP_ID:latest,GITHUB_APP_PRIVATE_KEY=FT_GITHUB_PRIVATE_KEY:latest' \
		--max-instances 10 \
		--port 8080 \
		--allow-unauthenticated \
		--image gcr.io/${PROJECT_ID}/${PROJECT_NAME}:latest
