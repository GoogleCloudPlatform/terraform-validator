#!/usr/bin/env bash

set -e

if [ $# -eq 0 ]
then
	echo "No version supplied. Run \`make release VERSION=X.X\`"
	exit 1
fi

version=$1

echo $version | grep -Eq "^v[0-9]+\.[0-9]+\.[0-9]+$" && echo "match" || echo "no match"

if ! echo $version | grep -Eq "^v[0-9]+\.[0-9]+\.[0-9]+$"
then
	echo "Invalid version: ${version}"
	echo "Please specify a semantic version prefixed with a v (e.g. v0.0.0)."
	exit 1
fi

platforms="linux windows darwin"
build_dir=./bin
name=terraform-validator
ldflags="-X github.com/GoogleCloudPlatform/terraform-validator/tfgcv.buildVersion=${version}"
release_bucket=terraform-validator

# Build release versions
for platform in ${platforms}; do
	echo "Building version ${version} for ${platform}..."
	GO111MODULE=on GOOS=${platform} GOARCH=amd64 CGO_ENABLED=0 go build -ldflags "${ldflags}" -o "${build_dir}/${name}-${platform}-amd64" .
	echo "Done building version ${version} for ${platform}"
done

echo "Creating Github tag ${version}"
git tag "${version}"
git push origin "${version}"
echo "Github tag ${version} created"

# Publish release versions
echo "Pushing releases to Google Storage"
gsutil cp ${build_dir}/*-amd64 gs://${release_bucket}/releases/${version}
echo "Releases pushed to Google Storage"

echo "Create a new release by visiting https://github.com/GoogleCloudPlatform/terraform-validator/releases/new?tag=${version}&title=${version}"
