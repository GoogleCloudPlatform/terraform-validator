#!/usr/bin/env bash

set -e

if [ $# -eq 0 ]
then
	echo "No version supplied. Run \`make release VERSION=X.X\`"
	exit 1
fi

version=$1

if ! echo $version | grep -Eq "^[0-9]+\.[0-9]+\.[0-9]+$"
then
	echo "Invalid version: ${version}"
	echo "Please specify a semantic version with no prefix (e.g. X.X.X)."
	exit 1
fi

echo "Listing supported resources to ${version}-resources.txt"
make build
bin/terraform-validator list-supported-resources > ${version}-resources.txt


echo "Creating Github tag v${version}"
git tag "v${version}"
git push origin "v${version}"
echo "Github tag v${version} created"

echo "Create a new release by visiting https://github.com/GoogleCloudPlatform/terraform-validator/releases/new?tag=v${version}&title=v${version}"
