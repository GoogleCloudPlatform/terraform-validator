package test

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/GoogleCloudPlatform/terraform-validator/converters/google"
	"github.com/GoogleCloudPlatform/terraform-validator/tfgcv"
	"github.com/GoogleCloudPlatform/terraform-validator/version"
)

func TestReadPlannedAssetsCoverage(t *testing.T) {
	if !version.Supported(version.TF12) {
		t.Skipf("TestReadPlannedAssetsCoverage runs on terraform v0.12 or above.")
	}
	cases := []struct {
		name string
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
		{name: "full_compute_firewall"},
		{name: "full_compute_instance"},
		{name: "full_container_cluster"},
		{name: "full_container_node_pool"},
		{name: "full_sql_database_instance"},
		{name: "full_storage_bucket"},
	}
	for i := range cases {
		// Allocate a variable to make sure test can run in parallel.
		c := cases[i]
		t.Run(c.name, func(t *testing.T) {
			t.Parallel()
			// Create a temporary directory for running terraform.
			dir, err := ioutil.TempDir(tmpDir, "terraform")
			if err != nil {
				log.Fatal(err)
			}
			defer os.RemoveAll(dir)

			generateTestFiles(t, "../testdata/templates", dir, c.name+".json")
			generateTestFiles(t, "../testdata/templates", dir, c.name+".tfplan.json")

			// Unmarshal payload from testfile into `want` variable.
			f := filepath.Join(dir, c.name+".json")
			payload, err := ioutil.ReadFile(f)
			if err != nil {
				t.Fatalf("cannot open %s, got: %s", f, err)
			}
			var want []google.Asset
			if err := json.Unmarshal(payload, &want); err != nil {
				t.Fatalf("cannot unmarshal JSON into assets: %s", err)
			}

			planfile := filepath.Join(dir, c.name+".tfplan.json")
			got, err := tfgcv.ReadPlannedAssets(planfile, data.Provider["project"], data.Ancestry, true)
			if err != nil {
				t.Fatalf("ReadPlannedAssets(%s, %s, %s, %t): %v", planfile, data.Provider["project"], data.Ancestry, true, err)
			}

			gotJSON := normalizeAssets(t, got, true)
			wantJSON := normalizeAssets(t, want, true)
			require.JSONEq(t, string(wantJSON), string(gotJSON))
		})
	}
}
