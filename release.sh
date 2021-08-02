#!/usr/bin/env bash

set -e

if [ $# -eq 0 ]
then
	echo "No version supplied. Run \`make release VERSION=X.X\`"
	exit 1
fi

version=$1

if ! echo $version | grep -Eq "^v[0-9]+\.[0-9]+\.[0-9]+$"
then
	echo "Invalid version: ${version}"
	echo "Please specify a semantic version prefixed with a v (e.g. v0.0.0)."
	exit 1
fi

architectures="amd64 arm64"
platforms="linux windows darwin"
skip_platform_arch_pairs=" windows/arm64 "

build_dir=./bin
name=terraform-validator
ldflags="-X github.com/GoogleCloudPlatform/terraform-validator/tfgcv.buildVersion=${version}"
release_bucket=terraform-validator

# Build release versions
for platform in ${platforms}; do
	for arch in ${architectures}; do
		if [[ " ${skip_platform_arch_pairs[@]} " =~ " ${platform}/${arch} " ]]; then
			echo "Skipped unsupported platform/arch pair ${platform}/${arch}"
			continue
		fi


		echo "Building version ${version} for platform ${platform} / arch ${arch}..."
		GO111MODULE=on GOOS=${platform} GOARCH=${arch} CGO_ENABLED=0 go build -ldflags "${ldflags}" -o "${build_dir}/${name}-${platform}-${arch}" .
		echo "Done building version ${version} for platform ${platform} / arch ${arch}"
	done
done

echo "Creating Github tag ${version}"
git tag "${version}"
git push origin "${version}"
echo "Github tag ${version} created"

# Publish release versions
echo "Pushing releases to Google Storage"
for arch in ${architectures}; do
	gsutil cp ${build_dir}/*-${arch} gs://${release_bucket}/releases/${version}
done
echo "Releases pushed to Google Storage"

echo "Create a new release by visiting https://github.com/GoogleCloudPlatform/terraform-validator/releases/new?tag=${version}&title=${version}"
