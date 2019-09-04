PLATFORMS := linux windows darwin
BUILD_DIR=./bin
NAME=terraform-validator
RELEASE_BUCKET=terraform-validator
DATE=`date +%Y-%m-%d`
LDFLAGS="-X github.com/GoogleCloudPlatform/terraform-validator/tfgcv.buildVersion=${DATE}"
BUILDTAGS=$(shell grep -q 'github.com/hashicorp/terraform v0.11' go.mod && echo 'tf_0_11')

prepare-v11:
	@echo "Prepare the environment for TF v0.11"
	cp go_tf_0_11.mod go.mod
	cp go_tf_0_11.sum go.sum
	go mod vendor

prepare-v12:
	@echo "Prepare the environment for TF v0.12"
	cp go_tf_0_12.mod go.mod
	cp go_tf_0_12.sum go.sum
	go mod vendor

test:
	# Skip integration tests in ./test/
	GO111MODULE=on go test -tags=$(BUILDTAGS) `go list -tags=$(BUILDTAGS) ./... | grep -v terraform-validator/test`

run-docker:
	docker run -it -v `pwd`:/terraform-validator -v ${GOOGLE_APPLICATION_CREDENTIALS}:/terraform-validator/credentials.json --entrypoint=/bin/bash --env TEST_PROJECT=${PROJECT_ID} --env TEST_CREDENTIALS=./credentials.json terraform-validator;

test-integration:
	go test -tags=$(BUILDTAGS) -v ./test

build-docker:
	docker build -f ./Dockerfile -t terraform-validator .

build:
	GO111MODULE=on go build -tags=$(BUILDTAGS) -ldflags ${LDFLAGS} -mod=vendor -o ${BUILD_DIR}/${NAME}

release: $(PLATFORMS)

publish:
	gsutil cp ${BUILD_DIR}/*-amd64 gs://${RELEASE_BUCKET}/releases/${DATE}

$(PLATFORMS):
	GO111MODULE=on GOOS=$@ GOARCH=amd64 CGO_ENABLED=0 go build -tags=$(BUILDTAGS) -mod=vendor -ldflags ${LDFLAGS} -o "${BUILD_DIR}/${NAME}-$@-amd64" .

clean:
	rm bin/${NAME}*

.PHONY: test test-e2e build build-docker release $(PLATFORMS) clean publish
