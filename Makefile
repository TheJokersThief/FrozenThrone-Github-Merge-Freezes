.PHONY: build
build:
	go build -o bin/frozen_throne ./frozen_throne

.PHONY: run
run: build
	GCS_BUCKET="iamevan-test-bucket" WRITE_SECRET="secret" READ_ONLY_SECRET="read-secret" ./bin/frozen_throne