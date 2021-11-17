module github.com/GoogleCloudPlatform/terraform-validator

replace github.com/hashicorp/terraform-provider-google/v3 => /Users/slevenick/go-workspace/src/github.com/terraform-providers/terraform-provider-google

require (
	cloud.google.com/go/bigtable v1.10.1
	github.com/apparentlymart/go-cidr v1.1.0
	github.com/davecgh/go-spew v1.1.1
	github.com/forseti-security/config-validator v0.0.0-20210621194145-08e4202b50d8
	github.com/go-lintpack/lintpack v0.5.2 // indirect
	github.com/golang/protobuf v1.5.2
	github.com/golangci/errcheck v0.0.0-20181223084120-ef45e06d44b6 // indirect
	github.com/golangci/go-tools v0.0.0-20190318055746-e32c54105b7c // indirect
	github.com/golangci/goconst v0.0.0-20180610141641-041c5f2b40f3 // indirect
	github.com/golangci/gocyclo v0.0.0-20180528134321-2becd97e67ee // indirect
	github.com/golangci/gosec v0.0.0-20190211064107-66fb7fc33547 // indirect
	github.com/golangci/ineffassign v0.0.0-20190609212857-42439a7714cc // indirect
	github.com/golangci/prealloc v0.0.0-20180630174525-215b22d4de21 // indirect
	github.com/hashicorp/errwrap v1.0.0
	github.com/hashicorp/go-cleanhttp v0.5.2
	github.com/hashicorp/terraform-json v0.12.0
	github.com/hashicorp/terraform-plugin-sdk/v2 v2.5.0
	github.com/hashicorp/terraform-provider-google v1.20.0 // indirect
	github.com/hashicorp/terraform-provider-google/v3 v3.0.0-20211117190748-d65dc11ef418 // indirect
	github.com/klauspost/cpuid v1.2.0 // indirect
	github.com/mitchellh/go-homedir v1.1.0
	github.com/pkg/errors v0.9.1
	github.com/shirou/gopsutil v0.0.0-20180427012116-c95755e4bcd7 // indirect
	github.com/shirou/w32 v0.0.0-20160930032740-bb4de0191aa4 // indirect
	github.com/spf13/cobra v1.1.3
	github.com/stretchr/testify v1.7.0
	go.uber.org/atomic v1.9.0 // indirect
	go.uber.org/multierr v1.7.0 // indirect
	go.uber.org/zap v1.19.1
	golang.org/x/oauth2 v0.0.0-20211005180243-6b3c2da341f1
	google.golang.org/api v0.60.0
	google.golang.org/genproto v0.0.0-20211021150943-2b146023228c
	google.golang.org/grpc v1.40.0
	google.golang.org/protobuf v1.27.1 // indirect
	gopkg.in/airbrake/gobrake.v2 v2.0.9 // indirect
	gopkg.in/gemnasium/logrus-airbrake-hook.v2 v2.1.2 // indirect
	sourcegraph.com/sqs/pbtypes v0.0.0-20180604144634-d3ebe8f20ae4 // indirect
)

go 1.14
