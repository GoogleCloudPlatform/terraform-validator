package test

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	crmv1 "google.golang.org/api/cloudresourcemanager/v1"

	resources "github.com/GoogleCloudPlatform/terraform-validator/converters/google/resources"
	"github.com/GoogleCloudPlatform/terraform-validator/tfdata"
	"github.com/GoogleCloudPlatform/terraform-validator/tfplan"
	provider "github.com/hashicorp/terraform-provider-google-beta/google-beta"
)

func TestIAMFetchFullResource(t *testing.T) {
	cases := []struct {
		name string
	}{
		{name: "example_bigquery_dataset_iam_binding"},
		{name: "example_bigquery_dataset_iam_member"},
		{name: "example_cloud_run_service_iam_binding"},
		{name: "example_cloud_run_service_iam_member"},
		{name: "example_compute_instance_iam_binding"},
		{name: "example_compute_instance_iam_member"},
		{name: "example_folder_iam_binding"},
		{name: "example_folder_iam_member"},
		{name: "example_kms_crypto_key_iam_binding"},
		{name: "example_kms_crypto_key_iam_member"},
		{name: "example_kms_key_ring_iam_binding"},
		{name: "example_kms_key_ring_iam_member"},
		{name: "example_organization_iam_binding"},
		{name: "example_organization_iam_member"},
		{name: "example_project_iam_binding"},
		{name: "example_project_iam_member"},
		{name: "example_pubsub_subscription_iam_binding"},
		{name: "example_pubsub_subscription_iam_member"},
		{name: "example_secret_manager_secret_iam_binding"},
		{name: "example_secret_manager_secret_iam_member"},
		{name: "example_spanner_database_iam_binding"},
		{name: "example_spanner_database_iam_member"},
		{name: "example_spanner_instance_iam_binding"},
		{name: "example_spanner_instance_iam_member"},
		{name: "example_storage_bucket_iam_binding"},
		{name: "example_storage_bucket_iam_member"},
	}

	converters := resources.ResourceConverters()
	schema := provider.Provider()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var payload []byte
		var err error
		if strings.Contains(r.URL.String(), "dataset") {
			obj := map[string][]interface{}{
				"access": {&crmv1.Policy{}},
			}
			payload, err = json.Marshal(obj)
		} else {
			payload, err = (&crmv1.Policy{}).MarshalJSON()
		}
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("failed to MarshalJSON: %s", err)))
			return
		}
		w.Write(payload)
	}))

	// Using Cleanup instead of defer because t.Parallel() does not block t.Run.
	t.Cleanup(func() {
		server.Close()
	})

	cfg := resources.NewTestConfig(server)

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			t.Parallel()

			// Create a temporary directory for generating tfplan.json from template.
			dir, err := ioutil.TempDir(tmpDir, "terraform")
			if err != nil {
				t.Fatal(err)
			}
			defer os.RemoveAll(dir)

			generateTestFiles(t, "../testdata/templates", dir, c.name+".tfplan.json")
			path := filepath.Join(dir, c.name+".tfplan.json")

			data, err := ioutil.ReadFile(path)
			if err != nil {
				t.Fatalf("opening JSON plan file: %s", err)
			}
			changes, err := tfplan.ReadResourceChanges(data)
			if err != nil {
				t.Fatalf("ReadResourceChanges failed: %s", err)
			}

			for _, rc := range changes {
				resource := schema.ResourcesMap[rc.Type]
				rd := tfdata.NewFakeResourceData(
					rc.Type,
					resource.Schema,
					rc.Change.After.(map[string]interface{}),
				)
				for _, converter := range converters[rd.Kind()] {
					if converter.FetchFullResource == nil {
						continue
					}
					_, err := converter.FetchFullResource(rd, cfg)
					if err != nil {
						t.Errorf("FetchFullResource() = %s, want = nil", err)
					}
				}
			}

		})
	}
}
