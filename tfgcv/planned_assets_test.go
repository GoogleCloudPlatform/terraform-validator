package tfgcv

import (
	"context"
	"path/filepath"
	"testing"

	"github.com/GoogleCloudPlatform/terraform-validator/version"
)

const (
	testDataDir      = "../testdata/plans"
	testProjectName  = "gcp-foundation-shared-devops"
	testAncestryName = "ancestry"
)

func TestReadPlannedAssets(t *testing.T) {
	type args struct {
		file     string
		project  string
		ancestry string
	}
	type testcase struct {
		name    string
		args    args
		want    int
		wantErr bool
	}

	var tests []testcase

	if version.Supported(version.TF12) {
		tests = append(tests, []testcase{
			// {
			// 	"Test TF12 and JSON plan",
			// 	args{"tf12plan.json", testProjectName, testAncestryName},
			// 	2,
			// 	false,
			// },
			// {
			// 	"Test TF12 with all coverage",
			// 	args{"tf12plan.allcoverage.json", "foobar", testAncestryName},
			// 	9,
			// 	false,
			// },
			{
				"Test TF12 with computed resources",
				args{"tf12plan.computed.json", "foobar", testAncestryName},
				9,
				false,
			},
		}...)
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testFile := filepath.Join(testDataDir, tt.args.file)
			var offline bool
			if version.Supported(version.TF12) {
				offline = true
			}
			ctx := context.Background()
			got, err := ReadPlannedAssets(ctx, testFile, tt.args.project, tt.args.ancestry, offline)
			if (err != nil) != tt.wantErr {
				t.Errorf("ReadPlannedAssets() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if len(got) != tt.want {
				t.Errorf("ReadPlannedAssets() = %v, want %v", len(got), tt.want)
			}
		})
	}
}
