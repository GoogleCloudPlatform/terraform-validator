package tfgcv

import (
	"context"
	"path/filepath"
	"testing"

	"github.com/GoogleCloudPlatform/terraform-validator/version"
)

const (
	testDataDir      = "../test/read_planned_assets"
	testProjectName  = "gl-akopachevskyy-sql-db"
	testAncestryName = "ancestry"
)


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

func testCases(t *testing.T) []testcase {
	var tests []testcase

	if version.Supported(version.TF12) {
		tests = append(tests, []testcase{
			{
				"Test TF12 and JSON plan",
				args{"tf12plan.json", testProjectName, testAncestryName},
				2,
				false,
			},
			// TODO: Add tf11plan.tfplan to the repository.
			// See https://github.com/GoogleCloudPlatform/terraform-validator/issues/74
			// {
			// 	"Test TF12 and binary plan should error out",
			// 	args{"tf11plan.tfplan", testProjectName, testAncestryName},
			// 	0,
			// 	true,
			// },
			{
				"Test TF12 with all coverage",
				args{"tf12plan.allcoverage.json", "foobar", testAncestryName},
				9,
				false,
			},
		}...)
	}
	if version.Supported(version.TF11) {
		tests = append(tests, []testcase{
			{
				"Test TF11 and json plan should error out",
				args{"tf12plan.json", testProjectName, testAncestryName},
				0,
				true,
			},
			// TODO: Add tf11plan.tfplan to the repository.
			// See https://github.com/GoogleCloudPlatform/terraform-validator/issues/74
			// {
			// 	"Test TF11 and binary plan",
			// 	args{"tf11plan.tfplan", testProjectName, testAncestryName},
			// 	2,
			// 	false,
			// },
		}...)
	}
	return tests
}

func TestReadPlannedAssets(t *testing.T) {
	for _, tt := range testCases(t) {
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

func TestReadCurrentAssets(t *testing.T) {
	for _, tt := range testCases(t) {
		t.Run(tt.name, func(t *testing.T) {
			testFile := filepath.Join(testDataDir, tt.args.file)
			var offline bool
			if version.Supported(version.TF12) {
				offline = true
			}
			ctx := context.Background()
			got, err := ReadCurrentAssets(ctx, testFile, tt.args.project, tt.args.ancestry, offline)
			if (err != nil) != tt.wantErr {
				t.Errorf("ReadCurrentAssets() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if len(got) != tt.want {
				t.Errorf("ReadCurrentAssets() = %v, want %v", len(got), tt.want)
			}
		})
	}
}
