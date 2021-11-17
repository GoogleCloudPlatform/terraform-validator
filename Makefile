build_dir=./bin
name=terraform-validator

test:
	# Skip integration tests in ./test/ using -short flag
	GO111MODULE=on go test -short ./...

run-docker:
	docker run -it -v `pwd`:/terraform-validator -v ${GOOGLE_APPLICATION_CREDENTIALS}:/terraform-validator/credentials.json --entrypoint=/bin/bash --env TEST_PROJECT=${PROJECT_ID} --env TEST_CREDENTIALS=./credentials.json terraform-validator;

test-integration:
	go version
	terraform --version
	go test -v -run=CLI ./test

test-go-licenses:
	cd .. && echo "Go version: $(go version)" && go get github.com/google/go-licenses
	go-licenses check .


build-docker:
	docker build -f ./Dockerfile -t terraform-validator .

build:
	GO111MODULE=on go build -o ${build_dir}/${name}

release:
	./release.sh ${VERSION}

clean:
	rm bin/${name}*

.PHONY: test test-e2e build build-docker release clean
