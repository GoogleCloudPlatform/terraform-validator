PLATFORMS := linux windows darwin
BUILD_DIR=./bin
NAME=terraform-validator
RELEASE_BUCKET=terraform-validator
DATE=`date +%Y-%m-%d`
LDFLAGS="-X github.com/GoogleCloudPlatform/terraform-validator/tfgcv.buildVersion=${DATE}"

test:
	# Skip integration tests in ./test/ using -short flag
	GO111MODULE=on go test -mod=vendor -short ./...

run-docker:
	docker run -it -v `pwd`:/terraform-validator -v ${GOOGLE_APPLICATION_CREDENTIALS}:/terraform-validator/credentials.json --entrypoint=/bin/bash --env TEST_PROJECT=${PROJECT_ID} --env TEST_CREDENTIALS=./credentials.json terraform-validator;

test-integration:
	go test -mod=vendor -v -run=CLI ./test

build-docker:
	docker build -f ./Dockerfile -t terraform-validator .

build:
	GO111MODULE=on go build -mod=vendor -ldflags ${LDFLAGS} -o ${BUILD_DIR}/${NAME}

release: $(PLATFORMS)

publish:
	gsutil cp ${BUILD_DIR}/*-amd64 gs://${RELEASE_BUCKET}/releases/${DATE}

$(PLATFORMS):
	GO111MODULE=on GOOS=$@ GOARCH=amd64 CGO_ENABLED=0 go build -mod=vendor -ldflags ${LDFLAGS} -o "${BUILD_DIR}/${NAME}-$@-amd64" .

clean:
	rm bin/${NAME}*

.PHONY: test test-e2e build build-docker release $(PLATFORMS) clean publish
