// Copyright 2019 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package test

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/GoogleCloudPlatform/terraform-validator/converters/google"
)

// TestCLI tests the "convert" and "validate" subcommand against a generated .tfplan file.
func TestCLI(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode.")
		return
	}
	// Define the reusable constraints to be use for the test cases.
	type constraint struct {
		name            string
		wantViolation   bool
		wantOutputRegex string
	}
	// Currently, we only test one rule. Moving forward, resource specific rules
	// should be added to increase the coverage.
	alwaysViolate := constraint{name: "always_violate", wantViolation: true, wantOutputRegex: "Constraint GCPAlwaysViolatesConstraintV1.always_violates_all on resource"}

	// Test cases for each type of resource is defined here.
	cases := []struct {
		name                 string
		constraints          []constraint
		compareConvertOutput compareConvertOutputFunc
	}{
		// cli-only, the following tests are not in read_test or
		// have unique paramters that separate them
		{name: "bucket"},
		{name: "disk"},
		{name: "firewall"},
		{name: "instance"},
		{name: "sql"},
		{name: "example_compute_forwarding_rule"},
		{name: "example_compute_instance"},
		{name: "example_project_create_empty_project_id"},
		{name: "example_project_iam_member_empty_project"},
		{name: "example_storage_bucket"},
		{name: "example_folder_iam_binding", compareConvertOutput: compareMergedIamBindingOutput},
		{name: "example_folder_iam_member", compareConvertOutput: compareMergedIamMemberOutput},
		{name: "example_project_create", constraints: []constraint{alwaysViolate, {name: "project_match_target", wantViolation: false, wantOutputRegex: ""}}},
		{name: "example_project_update", constraints: []constraint{alwaysViolate, {name: "project_match_target", wantViolation: true, wantOutputRegex: "Constraint GCPAlwaysViolatesConstraintV1.always_violates_project_match_target on resource"}}},
		{name: "example_project_iam_binding", compareConvertOutput: compareMergedIamBindingOutput},
		{name: "example_project_iam_member", compareConvertOutput: compareMergedIamMemberOutput},
		{name: "example_storage_bucket_iam_binding", compareConvertOutput: compareMergedIamBindingOutput},
		{name: "example_storage_bucket_iam_member", compareConvertOutput: compareMergedIamMemberOutput},
		// auto inserted tests that are not in list above or manually inserted in read_test.go
		{name: "example_access_context_manager_access_policy"},
		{name: "example_access_context_manager_service_perimeter"},
		{name: "example_bigquery_dataset"},
		{name: "example_bigquery_dataset_iam_binding"},
		{name: "example_bigquery_dataset_iam_member"},
		{name: "example_bigquery_dataset_iam_policy"},
		{name: "example_bigquery_dataset_iam_policy_empty_policy_data"},
		{name: "example_bigquery_table"},
		{name: "example_bigtable_instance"},
		{name: "example_cloud_run_mapping"},
		{name: "example_cloud_run_service"},
		{name: "example_cloud_run_service_iam_binding"},
		{name: "example_cloud_run_service_iam_member"},
		{name: "example_cloud_run_service_iam_policy"},
		{name: "example_compute_address"},
		{name: "example_compute_disk"},
		{name: "example_compute_disk_empty_image"},
		{name: "example_compute_firewall"},
		{name: "example_compute_global_address"},
		{name: "example_compute_global_forwarding_rule"},
		{name: "example_compute_network"},
		{name: "example_compute_security_policy"},
		{name: "example_compute_snapshot"},
		{name: "example_compute_ssl_policy"},
		{name: "example_compute_subnetwork"},
		{name: "example_container_cluster"},
		{name: "example_dns_managed_zone"},
		{name: "example_dns_policy"},
		{name: "example_filestore_instance"},
		{name: "example_folder_iam_policy"},
		{name: "example_folder_organization_policy"},
		{name: "example_google_cloudfunctions_function"},
		{name: "example_google_sql_database"},
		{name: "example_kms_crypto_key"},
		{name: "example_kms_crypto_key_iam_binding"},
		{name: "example_kms_crypto_key_iam_member"},
		{name: "example_kms_crypto_key_iam_policy"},
		{name: "example_kms_key_ring"},
		{name: "example_kms_key_ring_iam_binding"},
		{name: "example_kms_key_ring_iam_member"},
		{name: "example_kms_key_ring_iam_policy"},
		{name: "example_monitoring_notification_channel"},
		{name: "example_organization_iam_binding"},
		{name: "example_organization_iam_custom_role"},
		{name: "example_organization_iam_member"},
		{name: "example_organization_iam_policy"},
		{name: "example_organization_policy"},
		{name: "example_project_iam"},
		{name: "example_project_iam_custom_role"},
		{name: "example_project_iam_policy"},
		{name: "example_project_in_folder"},
		{name: "example_project_in_org"},
		{name: "example_project_organization_policy"},
		{name: "example_project_service"},
		{name: "example_pubsub_lite_reservation"},
		{name: "example_pubsub_lite_subscription"},
		{name: "example_pubsub_lite_topic"},
		{name: "example_pubsub_schema"},
		{name: "example_pubsub_subscription"},
		{name: "example_pubsub_subscription_iam_binding"},
		{name: "example_pubsub_subscription_iam_member"},
		{name: "example_pubsub_subscription_iam_policy"},
		{name: "example_pubsub_topic"},
		{name: "example_spanner_database"},
		{name: "example_spanner_database_iam_binding"},
		{name: "example_spanner_database_iam_member"},
		{name: "example_spanner_database_iam_policy"},
		{name: "example_spanner_instance_iam_binding"},
		{name: "example_spanner_instance_iam_member"},
		{name: "example_spanner_instance_iam_policy"},
		{name: "example_sql_database_instance"},
		{name: "example_storage_bucket_iam_member_random_suffix"},
		{name: "example_storage_bucket_iam_policy"},
		{name: "example_vpc_access_connector"},
		{name: "full_compute_firewall"},
		{name: "full_compute_instance"},
		{name: "full_container_cluster"},
		{name: "full_container_node_pool"},
		{name: "full_spanner_instance"},
		{name: "full_sql_database_instance"},
		{name: "full_storage_bucket"},
	}

	// Map of cases to skip to reasons for the skip
	skipCases := map[string]string{
		"TestCLI/v=0.12/tf=example_compute_forwarding_rule/offline=true/cmd=convert":                              "temperarily skip because of the predictable drift in offline mode",
		"TestCLI/v=0.12/tf=example_compute_forwarding_rule/offline=true/cmd=validate/constraint=always_violate":   "temperarily skip because of the predictable drift in offline mode",
		"TestCLI/v=0.12/tf=example_compute_instance/offline=true/cmd=convert":                                     "compute_instance doesn't work in offline mode - github.com/hashicorp/terraform-provider-google/issues/8489",
		"TestCLI/v=0.12/tf=example_compute_instance/offline=true/cmd=validate/constraint=always_violate":          "compute_instance doesn't work in offline mode - github.com/hashicorp/terraform-provider-google/issues/8489",
		"TestCLI/v=0.12/tf=example_organization_iam_binding/offline=false/cmd=convert":                            "skip because test runner doesn't have org permissions",
		"TestCLI/v=0.12/tf=example_organization_iam_binding/offline=false/cmd=validate/constraint=always_violate": "skip because test runner doesn't have org permissions",
		"TestCLI/v=0.12/tf=example_organization_iam_member/offline=false/cmd=convert":                             "skip because test runner doesn't have org permissions",
		"TestCLI/v=0.12/tf=example_organization_iam_member/offline=false/cmd=validate/constraint=always_violate":  "skip because test runner doesn't have org permissions",
		"TestCLI/v=0.12/tf=example_project_iam/offline=false/cmd=convert":                                         "example_project_iam is too complex to untangle merges with online data generically",
		"TestCLI/v=0.12/tf=example_storage_bucket_iam_member_random_suffix/offline=false":                         "test produces inconsistent results based on randomized names - github.com/GoogleCloudPlatform/terraform-validator/issues/259",
		"TestCLI/v=0.12/tf=example_storage_bucket_iam_member_random_suffix/offline=false/cmd=convert":             "test produces inconsistent results based on randomized names - github.com/GoogleCloudPlatform/terraform-validator/issues/259",
	}
	for i := range cases {
		// Allocate a variable to make sure test can run in parallel.
		c := cases[i]
		// Add default constraints if not set.
		if len(c.constraints) == 0 {
			c.constraints = []constraint{alwaysViolate}
		}

		// Add default convert comparison func if not set
		if c.compareConvertOutput == nil {
			c.compareConvertOutput = compareUnmergedConvertOutput
		}

		// Test both offline and online mode.
		for _, offline := range []bool{true, false} {
			offline := offline
			t.Run(fmt.Sprintf("v=0.12/tf=%s/offline=%t", c.name, offline), func(t *testing.T) {
				t.Parallel()
				// Create a temporary directory for running terraform.
				dir, err := ioutil.TempDir(tmpDir, "terraform")
				if err != nil {
					log.Fatal(err)
				}
				defer os.RemoveAll(dir)

				// Generate the <name>.tf and <name>_assets.json files into the temporary directory.
				generateTestFiles(t, "../testdata/templates", dir, c.name+".tf")
				generateTestFiles(t, "../testdata/templates", dir, c.name+".json")

				// Uses glob matching to match generateTestFiles internals.
				tfstateMatches, err := filepath.Glob(filepath.Join("../testdata/templates", c.name+".tfstate"))
				if err != nil {
					t.Fatalf("malformed glob: %v", err)
				}
				if tfstateMatches != nil {
					generateTestFiles(t, "../testdata/templates", dir, c.name+".tfstate")
					err = os.Rename(
						filepath.Join(dir, c.name+".tfstate"),
						filepath.Join(dir, "terraform.tfstate"),
					)
					if err != nil {
						t.Fatalf("renaming tfstate: %v", err)
					}
				}

				terraformWorkflow(t, dir, c.name)

				t.Run("cmd=convert", func(t *testing.T) {
					if reason, exists := skipCases[t.Name()]; exists {
						t.Skip(reason)
					}
					testConvertCommand(t, dir, c.name, c.name, offline, true, c.compareConvertOutput)
				})

				for _, ct := range c.constraints {
					t.Run(fmt.Sprintf("cmd=validate/constraint=%s", ct.name), func(t *testing.T) {
						if reason, exists := skipCases[t.Name()]; exists {
							t.Skip(reason)
						}
						testValidateCommand(t, ct.wantViolation, ct.wantOutputRegex, dir, c.name, offline, true, ct.name)
					})
				}
			})
		}
	}
}

// TestCLIWithoutProject tests the "convert" and "validate" subcommand
// against a generated .tfplan file, and it does not have --project option.
func TestCLIWithoutProject(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode.")
		return
	}
	// Define the reusable constraints to be use for the test cases.
	type constraint struct {
		name            string
		wantViolation   bool
		wantOutputRegex string
	}
	// Currently, we only test one rule. Moving forward, resource specific rules
	// should be added to increase the coverage.
	alwaysViolate := constraint{name: "always_violate", wantViolation: true, wantOutputRegex: "Constraint GCPAlwaysViolatesConstraintV1.always_violates_all on resource"}

	// Test cases for each type of resource is defined here.
	cases := []struct {
		name                 string
		constraints          []constraint
		compareConvertOutput compareConvertOutputFunc
	}{
		{name: "example_project_create_empty_project_id"},
		{name: "example_storage_bucket"},
		{name: "example_project_iam_member_empty_project"},
	}

	for i := range cases {
		// Allocate a variable to make sure test can run in parallel.
		c := cases[i]
		// Add default constraints if not set.
		if len(c.constraints) == 0 {
			c.constraints = []constraint{alwaysViolate}
		}

		// Add default convert comparison func if not set
		if c.compareConvertOutput == nil {
			c.compareConvertOutput = compareUnmergedConvertOutput
		}

		// Test both offline and online mode.
		for _, offline := range []bool{true, false} {
			offline := offline
			t.Run(fmt.Sprintf("v=0.12/tf=%s/offline=%t", c.name, offline), func(t *testing.T) {
				t.Parallel()
				// Create a temporary directory for running terraform.
				dir, err := ioutil.TempDir(tmpDir, "terraform")
				if err != nil {
					log.Fatal(err)
				}
				defer os.RemoveAll(dir)

				// Generate the <name>.tf and <name>_assets.json files into the temporary directory.
				generateTestFiles(t, "../testdata/templates", dir, c.name+".tf")
				generateTestFiles(t, "../testdata/templates", dir, c.name+"_without_default_project.json")

				// Uses glob matching to match generateTestFiles internals.
				tfstateMatches, err := filepath.Glob(filepath.Join("../testdata/templates", c.name+".tfstate"))
				if err != nil {
					t.Fatalf("malformed glob: %v", err)
				}
				if tfstateMatches != nil {
					generateTestFiles(t, "../testdata/templates", dir, c.name+".tfstate")
					err = os.Rename(
						filepath.Join(dir, c.name+".tfstate"),
						filepath.Join(dir, "terraform.tfstate"),
					)
					if err != nil {
						t.Fatalf("renaming tfstate: %v", err)
					}
				}

				terraformWorkflow(t, dir, c.name)

				t.Run("cmd=convert", func(t *testing.T) {
					testConvertCommand(t, dir, c.name, c.name+"_without_default_project", offline, false, c.compareConvertOutput)
				})

				for _, ct := range c.constraints {
					t.Run(fmt.Sprintf("cmd=validate/constraint=%s", ct.name), func(t *testing.T) {
						testValidateCommand(t, ct.wantViolation, ct.wantOutputRegex, dir, c.name, offline, false, ct.name)
					})
				}
			})
		}
	}
}

type compareConvertOutputFunc func(t *testing.T, expected []google.Asset, actual []google.Asset, offline bool)

func compareUnmergedConvertOutput(t *testing.T, expected []google.Asset, actual []google.Asset, offline bool) {
	expectedAssets := normalizeAssets(t, expected, offline)
	actualAssets := normalizeAssets(t, actual, offline)
	require.ElementsMatch(t, expectedAssets, actualAssets)
}

// For merged IAM members, only consider whether the expected members are present.
func compareMergedIamMemberOutput(t *testing.T, expected []google.Asset, actual []google.Asset, offline bool) {
	var normalizedActual []google.Asset
	for i := range expected {
		expectedAsset := expected[i]
		actualAsset := actual[i]

		// Copy actualAsset
		normalizedActualAsset := actualAsset

		expectedBindings := map[string]map[string]struct{}{}
		for _, binding := range expectedAsset.IAMPolicy.Bindings {
			expectedBindings[binding.Role] = map[string]struct{}{}
			for _, member := range binding.Members {
				expectedBindings[binding.Role][member] = struct{}{}
			}
		}

		iamPolicy := google.IAMPolicy{}
		for _, binding := range actualAsset.IAMPolicy.Bindings {
			if expectedMembers, exists := expectedBindings[binding.Role]; exists {
				iamBinding := google.IAMBinding{
					Role: binding.Role,
				}
				for _, member := range binding.Members {
					if _, exists := expectedMembers[member]; exists {
						iamBinding.Members = append(iamBinding.Members, member)
					}
				}
				iamPolicy.Bindings = append(iamPolicy.Bindings, iamBinding)
			}
		}
		normalizedActualAsset.IAMPolicy = &iamPolicy
		normalizedActual = append(normalizedActual, normalizedActualAsset)
	}

	expectedAssets := normalizeAssets(t, expected, offline)
	actualAssets := normalizeAssets(t, normalizedActual, offline)
	require.ElementsMatch(t, expectedAssets, actualAssets)
}

// For merged IAM bindings, only consider whether the expected bindings are as expected.
func compareMergedIamBindingOutput(t *testing.T, expected []google.Asset, actual []google.Asset, offline bool) {
	var normalizedActual []google.Asset
	for i := range expected {
		expectedAsset := expected[i]
		actualAsset := actual[i]

		// Copy actualAsset
		normalizedActualAsset := actualAsset

		expectedBindings := map[string]struct{}{}
		for _, binding := range expectedAsset.IAMPolicy.Bindings {
			expectedBindings[binding.Role] = struct{}{}
		}

		iamPolicy := google.IAMPolicy{}
		for _, binding := range actualAsset.IAMPolicy.Bindings {
			if _, exists := expectedBindings[binding.Role]; exists {
				iamPolicy.Bindings = append(iamPolicy.Bindings, binding)
			}
		}
		normalizedActualAsset.IAMPolicy = &iamPolicy
		normalizedActual = append(normalizedActual, normalizedActualAsset)
	}

	expectedAssets := normalizeAssets(t, expected, offline)
	actualAssets := normalizeAssets(t, normalizedActual, offline)
	require.ElementsMatch(t, expectedAssets, actualAssets)
}
