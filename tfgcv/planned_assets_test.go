package tfgcv

import (
	"encoding/json"
	"io/ioutil"
	"path/filepath"
	"testing"

	"github.com/GoogleCloudPlatform/terraform-validator/converters/google"
	"github.com/GoogleCloudPlatform/terraform-validator/version"
	"github.com/golang/protobuf/proto"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

const (
	testDataDir = "../test/read_planned_assets"
	testProjectName = "gl-akopachevskyy-sql-db"
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

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testFile := filepath.Join(testDataDir, tt.args.file)
			var offline bool
			if version.Supported(version.TF12) {
				offline = true
			}
			got, err := ReadPlannedAssets(testFile, tt.args.project, tt.args.ancestry, offline)
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

func TestReadPlannedAssetsCoverage(t *testing.T) {
	if !version.Supported(version.TF12) {
		t.Skipf("TestReadPlannedAssetsCoverage runs on terraform v0.12 or above.")
	}

	testdata := "testdata"
	defaultProject := "foobar"
	defaultAncestry := "organization/12345/folder/67890"
	cases := []struct {
		name     string
	}{
		{name: "example_compute_disk"},
		{name: "example_compute_firewall"},
		{name: "example_compute_instance"},
		{name: "example_container_cluster"},
		{name: "example_organization_iam_binding"},
		{name: "example_organization_iam_member"},
		{name: "example_organization_iam_policy"},
		{name: "example_project"},
		{name: "example_project_iam_binding"},
		{name: "example_project_iam_member"},
		{name: "example_project_iam_policy"},
		{name: "example_sql_database_instance"},
		{name: "example_storage_bucket"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			ancestry := defaultAncestry
			project := defaultProject
			testPlanFile := filepath.Join(testdata, c.name+"_tfplan.json")
			wantAssetPath := filepath.Join(testdata, c.name+"_assets.json")
			wantJSON, err := ioutil.ReadFile(wantAssetPath)
			if err != nil {
				t.Fatalf("cannot open %s, got: %s", wantAssetPath, err)
			}
			var want []google.Asset
			if err := json.Unmarshal(wantJSON, &want); err != nil {
				t.Fatalf("cannot unmarshal JSON into assets: %s", err)
			}
			assets, err := ReadPlannedAssets(testPlanFile, project, ancestry, true)
			if err != nil {
				t.Fatalf("ReadPlannedAssets(%s, %s, %s, %t): %v", testPlanFile, project, ancestry, true, err)
			}
			got := remarshal(t, assets)
			if diff := cmp.Diff(want, got, cmp.Comparer(proto.Equal), cmpopts.IgnoreUnexported(google.Asset{})); diff != "" {
				t.Errorf("ReadPlannedAssets(%s, %s, %s, %t) returned diff (-want +got):\n%s", testPlanFile, project, ancestry, true, diff)
			}
		})
	}
}

// remarshal runs Marshal and unmarshal to set type correctly.
//
// This is necessary as some of the types are interface and
// the result returned from ShowAssets of a different primitive (float64 vs int).
// For example, there are cases like float64(1000) is returned but the
// test fixture in JSON would decode that as int(1000).
func remarshal(t *testing.T, assets []google.Asset) []google.Asset {
	t.Helper()
	var got []google.Asset
	payload, err := json.Marshal(assets)
	if err != nil {
		t.Fatalf("cannot marshal returned assets: %v", err)
	}
	if err := json.Unmarshal(payload, &got); err != nil {
		t.Fatalf("cannot unmarshal returned assets: %v", err)
	}
	return got
}

