package test

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/GoogleCloudPlatform/terraform-validator/converters/google"
	"github.com/GoogleCloudPlatform/terraform-validator/tfgcv"
	"go.uber.org/zap"
)

func TestReadPlannedAssetsCoverage(t *testing.T) {
	cases := []struct {
		name string
	}{
		// read-only, the following tests are not in cli_test or
		// have unique paramters that seperate them
		{name: "example_project_create"},
		{name: "example_project_update"},
		{name: "example_project_iam_binding"},
		{name: "example_project_iam_member"},
		{name: "example_storage_bucket_iam_binding"},
		{name: "example_storage_bucket_iam_member"},
		// auto inserted tests that are not in list above or manually inserted in cli_test.go
		{name: "example_access_context_manager_service_perimeter"},
		{name: "example_bigquery_dataset"},
		{name: "example_bigquery_dataset_iam_binding"},
		{name: "example_bigquery_dataset_iam_member"},
		{name: "example_bigquery_dataset_iam_policy"},
		{name: "example_bigtable_instance"},
		{name: "example_compute_disk"},
		{name: "example_compute_disk_empty_image"},
		{name: "example_compute_firewall"},
		{name: "example_compute_global_forwarding_rule"},
		{name: "example_compute_network"},
		{name: "example_compute_snapshot"},
		{name: "example_compute_subnetwork"},
		{name: "example_container_cluster"},
		{name: "example_dns_managed_zone"},
		{name: "example_filestore_instance"},
		{name: "example_kms_crypto_key"},
		{name: "example_kms_crypto_key_iam_binding"},
		{name: "example_kms_crypto_key_iam_member"},
		{name: "example_kms_crypto_key_iam_policy"},
		{name: "example_kms_key_ring"},
		{name: "example_kms_key_ring_iam_binding"},
		{name: "example_kms_key_ring_iam_member"},
		{name: "example_kms_key_ring_iam_policy"},
		{name: "example_organization_iam_binding"},
		{name: "example_organization_iam_member"},
		{name: "example_organization_iam_policy"},
		{name: "example_project_iam"},
		{name: "example_project_iam_policy"},
		{name: "example_project_in_folder"},
		{name: "example_project_in_org"},
		{name: "example_project_organization_policy"},
		{name: "example_project_service"},
		{name: "example_pubsub_subscription"},
		{name: "example_pubsub_topic"},
		{name: "example_sql_database_instance"},
		{name: "example_storage_bucket"},
		{name: "example_storage_bucket_iam_member_random_suffix"},
		{name: "example_storage_bucket_iam_policy"},
		{name: "full_compute_firewall"},
		{name: "full_compute_instance"},
		{name: "full_container_cluster"},
		{name: "full_container_node_pool"},
		{name: "full_spanner_instance"},
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
			ctx := context.Background()
			got, err := tfgcv.ReadPlannedAssets(ctx, planfile, data.Provider["project"], data.Ancestry, true, false, zap.NewExample())
			if err != nil {
				t.Fatalf("ReadPlannedAssets(%s, %s, %s, %t): %v", planfile, data.Provider["project"], data.Ancestry, true, err)
			}

			expectedAssets := normalizeAssets(t, want, true)
			actualAssets := normalizeAssets(t, got, true)
			require.ElementsMatch(t, actualAssets, expectedAssets)
		})
	}
}
