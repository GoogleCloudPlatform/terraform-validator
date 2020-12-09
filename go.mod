module github.com/GoogleCloudPlatform/terraform-validator

require (
	cloud.google.com/go v0.73.0 // indirect
	cloud.google.com/go/bigtable v1.6.0 // indirect
	cloud.google.com/go/storage v1.12.0 // indirect
	github.com/GoogleCloudPlatform/terraform-google-conversion v0.0.0-20201209001225-610fadaf66b5
	github.com/Microsoft/go-winio v0.4.16 // indirect
	github.com/agext/levenshtein v1.2.3 // indirect
	github.com/apparentlymart/go-cidr v1.1.0 // indirect
	github.com/aws/aws-sdk-go v1.36.5 // indirect
	github.com/fatih/color v1.10.0 // indirect
	github.com/forseti-security/config-validator v0.0.0-20201204181659-b3da694a3d79
	github.com/gammazero/deque v0.0.0-20201010052221-3932da5530cc // indirect
	github.com/gammazero/workerpool v1.1.1 // indirect
	github.com/go-git/go-git/v5 v5.2.0 // indirect
	github.com/go-logr/logr v0.3.0 // indirect
	github.com/go-openapi/validate v0.20.0 // indirect
	github.com/golang/glog v0.0.0-20160126235308-23def4e6c14b
	github.com/golang/protobuf v1.4.3
	github.com/google/go-cmp v0.5.4
	github.com/google/gofuzz v1.2.0 // indirect
	github.com/hashicorp/errwrap v1.1.0 // indirect
	github.com/hashicorp/go-getter v1.5.1 // indirect
	github.com/hashicorp/go-hclog v0.15.0 // indirect
	github.com/hashicorp/go-multierror v1.1.0 // indirect
	github.com/hashicorp/go-plugin v1.4.0 // indirect
	github.com/hashicorp/go-uuid v1.0.2 // indirect
	github.com/hashicorp/hcl/v2 v2.8.0 // indirect
	github.com/hashicorp/terraform-config-inspect v0.0.0-20201102131242-0c45ba392e51 // indirect
	github.com/hashicorp/terraform-exec v0.11.0 // indirect
	github.com/hashicorp/terraform-json v0.7.0 // indirect
	github.com/hashicorp/terraform-plugin-sdk v1.16.0
	github.com/hashicorp/terraform-svchost v0.0.0-20200729002733-f050f53b9734 // indirect
	github.com/hashicorp/yamux v0.0.0-20200609203250-aecfd211c9ce // indirect
	github.com/kevinburke/ssh_config v0.0.0-20201106050909-4977a11b4351 // indirect
	github.com/mitchellh/cli v1.1.2 // indirect
	github.com/mitchellh/go-testing-interface v1.14.1 // indirect
	github.com/mitchellh/go-wordwrap v1.0.1 // indirect
	github.com/mitchellh/hashstructure v1.1.0 // indirect
	github.com/oklog/run v1.1.0 // indirect
	github.com/open-policy-agent/frameworks/constraint v0.0.0-20201118071520-0d37681951a4 // indirect
	github.com/open-policy-agent/opa v0.25.2 // indirect
	github.com/pkg/errors v0.9.1
	github.com/posener/complete v1.2.3 // indirect
	github.com/prometheus/client_golang v1.8.0 // indirect
	github.com/prometheus/common v0.15.0 // indirect
	github.com/spf13/afero v1.5.1 // indirect
	github.com/spf13/cobra v1.1.1
	github.com/stoewer/go-strcase v1.2.0 // indirect
	github.com/stretchr/testify v1.6.1
	github.com/terraform-providers/terraform-provider-google v1.20.1-0.20200228174759-ed183996d331
	github.com/vmihailenco/msgpack v4.0.4+incompatible // indirect
	github.com/vmihailenco/tagparser v0.1.2 // indirect
	github.com/xanzy/ssh-agent v0.3.0 // indirect
	github.com/zclconf/go-cty v1.7.0 // indirect
	github.com/zclconf/go-cty-yaml v1.0.2 // indirect
	golang.org/x/crypto v0.0.0-20201208171446-5f87f3452ae9 // indirect
	golang.org/x/exp v0.0.0-20201203231725-fa01524bc59d // indirect
	golang.org/x/lint v0.0.0-20201208152925-83fdc39ff7b5 // indirect
	golang.org/x/net v0.0.0-20201209123823-ac852fbbde11 // indirect
	golang.org/x/oauth2 v0.0.0-20201208152858-08078c50e5b5 // indirect
	golang.org/x/sys v0.0.0-20201207223542-d4d67f95c62d // indirect
	golang.org/x/term v0.0.0-20201207232118-ee85cb95a76b // indirect
	golang.org/x/time v0.0.0-20201208040808-7e3f01d25324 // indirect
	golang.org/x/tools v0.0.0-20201208233053-a543418bbed2 // indirect
	google.golang.org/api v0.36.0
	google.golang.org/genproto v0.0.0-20201209185603-f92720507ed4
	google.golang.org/grpc v1.34.0 // indirect
	honnef.co/go/tools v0.0.1-2020.1.6 // indirect
	k8s.io/apiextensions-apiserver v0.20.0 // indirect
	k8s.io/client-go v11.0.0+incompatible // indirect
	k8s.io/kubectl v0.20.0 // indirect
	sigs.k8s.io/controller-runtime v0.6.4 // indirect
)

replace github.com/GoogleCloudPlatform/terraform-google-conversion => ./vendor/github.com/GoogleCloudPlatform/terraform-google-conversion

go 1.13
