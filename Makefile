PROJECT_NAME = frozen_throne
PROJECT_ID ?= example-project
GCS_BUCKET ?= ${PROJECT_NAME}-test-bucket
WRITE_SECRET ?= secret
READ_ONLY_SECRET ?= secret-read-only

.PHONY: help
help: ## Show help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.PHONY: create_secrets
create_secrets: ## Create secret values
	echo -n "${WRITE_SECRET}" | gcloud --project ${PROJECT_ID} secrets create FT_WRITE_SECRET --replication-policy="automatic" --data-file=-
	echo -n "${READ_ONLY_SECRET}" | gcloud --project ${PROJECT_ID} secrets create FT_READ_ONLY_SECRET --replication-policy="automatic" --data-file=-

.PHONY: update_secrets
update_secrets: ## Update secret values
	echo -n "${WRITE_SECRET}" | gcloud --project ${PROJECT_ID} versions add FT_WRITE_SECRET --data-file=-
	echo -n "${READ_ONLY_SECRET}" | gcloud --project ${PROJECT_ID} versions add FT_READ_ONLY_SECRET --data-file=-

.PHONY: add_perms
add_perms: ## Add permissions for functions account to access secrets
	gcloud projects add-iam-policy-binding ${PROJECT_ID} \
		--member='serviceAccount:${PROJECT_ID}@appspot.gserviceaccount.com' \
		--role='roles/secretmanager.secretAccessor'

.PHONY: deploy_to_gfunctions
deploy_to_gfunctions: ## Deploy function to GCP Cloud Functions
	gcloud beta functions deploy ${PROJECT_NAME} \
		--region europe-west1 \
		--project ${PROJECT_ID} \
		--runtime go116 \
		--memory 128MB \
		--timeout 20s \
		--trigger-http \
		--entry-point IngestHTTP \
		--allow-unauthenticated \
		--set-env-vars GOOGLE_CLOUD_PROJECT=${PROJECT_ID},GCS_BUCKET="${GCS_BUCKET}" \
		--set-secrets 'WRITE_SECRET=FT_WRITE_SECRET:latest,READ_ONLY_SECRET=FT_READ_ONLY_SECRET:latest' \
		--max-instances 10