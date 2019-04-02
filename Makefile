PLATFORMS := linux windows darwin
BUILD_DIR=./bin
NAME=terraform-validator
RELEASE_BUCKET=terraform-validator
DATE=`date +%Y-%m-%d`

test:
	# Skip integration tests in ./test/
	GO111MODULE=on go test `go list ./... | grep -v terraform-validator/test`

test-e2e: build-docker
	set -e ;\
	CONTAINER=$$(docker create --env TEST_PROJECT=${PROJECT} --env TEST_CREDENTIALS=./credentials.json terraform-validator go test -v ./test) ;\
	echo $$CONTAINER ;\
	docker cp ${CREDENTIALS} $$CONTAINER:/terraform-validator ;\
	docker start --attach $$CONTAINER ;\

build-docker:
	docker build -f ./Dockerfile -t terraform-validator .

build:
	GO111MODULE=on go build -mod=vendor -o ${BUILD_DIR}/${NAME}

release: $(PLATFORMS)

publish:
	gsutil cp ${BUILD_DIR}/*-amd64 gs://${RELEASE_BUCKET}/releases/${DATE}

$(PLATFORMS):
	GO111MODULE=on GOOS=$@ GOARCH=amd64 CGO_ENABLED=0 go build -mod=vendor -o "${BUILD_DIR}/${NAME}-$@-amd64" .

clean:
	rm bin/${NAME}*

.PHONY: test test-e2e build build-docker release $(PLATFORMS) clean publish
