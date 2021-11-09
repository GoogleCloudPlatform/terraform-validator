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

release_dir="./release/${version}"
rm -rf $release_dir
mkdir -p $release_dir

echo "Go version: $(go version)"
go get github.com/google/go-licenses
echo "Downloading licenses and source code to bundle..."
# Ignore errors until https://github.com/google/go-licenses/pull/77 is merged
set +e
go-licenses save "github.com/GoogleCloudPlatform/terraform-validator" --save_path="./${release_dir}/THIRD_PARTY_NOTICES"
set -e
echo "Zipping licenses and source code..."
pushd "${release_dir}" > /dev/null
zip -rq9D "THIRD_PARTY_NOTICES.zip" "THIRD_PARTY_NOTICES/"
popd > /dev/null

architectures="amd64 arm64"
platforms="linux windows darwin"
skip_platform_arch_pairs=" windows/arm64 "

tar_gz_name=terraform-validator
ldflags="-X github.com/GoogleCloudPlatform/terraform-validator/tfgcv.buildVersion=v${version}"
release_bucket=terraform-validator

# Build release versions
for platform in ${platforms}; do
	if [[ "$platform" == "windows" ]]; then
		binary_name=terraform-validator.exe
	else
		binary_name=terraform-validator
	fi
	for arch in ${architectures}; do
		if [[ " ${skip_platform_arch_pairs[@]} " =~ " ${platform}/${arch} " ]]; then
			echo "Skipped unsupported platform/arch pair ${platform}/${arch}"
			continue
		fi

		echo "Building ${binary_name} v${version} for platform ${platform} / arch ${arch}..."
		GO111MODULE=on GOOS=${platform} GOARCH=${arch} CGO_ENABLED=0 go build -ldflags "${ldflags}" -o "${release_dir}/${binary_name}" .
		echo "Creating ${release_dir}/${tar_gz_name}_${platform}_${arch}-${version}.tar.gz"
		pushd "${release_dir}" > /dev/null
		tar -czf "${tar_gz_name}_${platform}_${arch}-${version}.tar.gz" "${binary_name}" "THIRD_PARTY_NOTICES.zip"
		popd > /dev/null
	done
done

echo "Creating Github tag v${version}"
git tag "v${version}"
git push origin "v${version}"
echo "Github tag v${version} created"

# Publish release versions
echo "Pushing releases to Google Storage"
gsutil cp ${release_dir}/*.tar.gz gs://${release_bucket}/releases/v${version}
echo "Releases pushed to Google Storage"

echo "Create a new release by visiting https://github.com/GoogleCloudPlatform/terraform-validator/releases/new?tag=${version}&title=${version}"
