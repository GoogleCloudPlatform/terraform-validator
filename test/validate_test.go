package test

import (
	"context"
	"fmt"
	"testing"

	"github.com/GoogleCloudPlatform/terraform-validator/converters/google"
	"github.com/GoogleCloudPlatform/terraform-validator/tfgcv"
	"github.com/GoogleCloudPlatform/terraform-validator/tfplan"
	"github.com/hashicorp/terraform/terraform"
	"google.golang.org/api/cloudresourcemanager/v1"
	"google.golang.org/api/option"
)

func TestValidate(t *testing.T) {
	// Retrieving test configurations needed for building test cases.
	cfg := configure(t)

	// The test case takes a diff in the form of terraform.instanceDiff as an
	// input, and compares number of violations returned with wantViolationCount.
	cases := []struct {
		kind               string
		name               string
		diff               *terraform.InstanceDiff
		wantViolationCount int
	}{
		{
			kind: "google_organization_iam_policy",
			name: "foo_organization_iam_policy",
			diff: &terraform.InstanceDiff{
				Attributes: map[string]*terraform.ResourceAttrDiff{
					"org_id":      &terraform.ResourceAttrDiff{New: "12345678"},
					"policy_data": &terraform.ResourceAttrDiff{New: `{"bindings":[{"members":["user:jane@example.com"],"role":"roles/editor"}]}`},
				},
			},
			wantViolationCount: 2,
		},
		{
			kind: "google_organization_iam_binding",
			name: "foo_organization_iam_binding",
			diff: &terraform.InstanceDiff{
				Attributes: map[string]*terraform.ResourceAttrDiff{
					"members.#":        &terraform.ResourceAttrDiff{New: "1"},
					"members.12345678": &terraform.ResourceAttrDiff{New: "alice@gmail.com"},
					"org_id":           &terraform.ResourceAttrDiff{New: "12345678"},
					"role":             &terraform.ResourceAttrDiff{New: "role/browser"},
				},
			},
			wantViolationCount: 2,
		},
		{
			kind: "google_organization_iam_member",
			name: "foo_organization_iam_member",
			diff: &terraform.InstanceDiff{
				Attributes: map[string]*terraform.ResourceAttrDiff{
					"member": &terraform.ResourceAttrDiff{New: "alice@gmail.com"},
					"org_id": &terraform.ResourceAttrDiff{New: "12345678"},
					"role":   &terraform.ResourceAttrDiff{New: "role/browser"},
				},
			},
			wantViolationCount: 2,
		},
		{
			kind: "google_folder_iam_policy",
			name: "foo_folder_iam_policy",
			diff: &terraform.InstanceDiff{
				Attributes: map[string]*terraform.ResourceAttrDiff{
					"folder":      &terraform.ResourceAttrDiff{New: "${google_folder.department1.name}"},
					"policy_data": &terraform.ResourceAttrDiff{New: `{"bindings":[{"members":["user:jane@example.com"],"role":"roles/editor"}]}`},
				},
			},
			wantViolationCount: 2,
		},
		{
			kind: "google_folder_iam_binding",
			name: "foo_folder_iam_binding",
			diff: &terraform.InstanceDiff{
				Attributes: map[string]*terraform.ResourceAttrDiff{
					"members.#":        &terraform.ResourceAttrDiff{New: "1"},
					"members.12345678": &terraform.ResourceAttrDiff{New: "alice@gmail.com"},
					"folder":           &terraform.ResourceAttrDiff{New: "${google_folder.department1.name}"},
					"role":             &terraform.ResourceAttrDiff{New: "role/browser"},
				},
			},
			wantViolationCount: 2,
		},
		{
			kind: "google_folder_iam_member",
			name: "foo_folder_iam_member",
			diff: &terraform.InstanceDiff{
				Attributes: map[string]*terraform.ResourceAttrDiff{
					"member": &terraform.ResourceAttrDiff{New: "alice@gmail.com"},
					"folder": &terraform.ResourceAttrDiff{New: "${google_folder.department1.name}"},
					"role":   &terraform.ResourceAttrDiff{New: "role/browser"},
				},
			},
			wantViolationCount: 2,
		},
		{
			kind: "google_project_iam_policy",
			name: "foo_project_iam_policy",
			diff: &terraform.InstanceDiff{
				Attributes: map[string]*terraform.ResourceAttrDiff{
					"policy_data": &terraform.ResourceAttrDiff{New: `{"bindings":[{"members":["user:jane@example.com"],"role":"roles/editor"}]}`},
					"project":     &terraform.ResourceAttrDiff{New: cfg.project},
				},
			},
			wantViolationCount: 1,
		},
		{
			kind: "google_project_iam_binding",
			name: "foo_project_iam_binding",
			diff: &terraform.InstanceDiff{
				Attributes: map[string]*terraform.ResourceAttrDiff{
					"members.#":        &terraform.ResourceAttrDiff{New: "1"},
					"members.12345678": &terraform.ResourceAttrDiff{New: "alice@gmail.com"},
					"project":          &terraform.ResourceAttrDiff{New: cfg.project},
					"role":             &terraform.ResourceAttrDiff{New: "role/browser"},
				},
			},
			wantViolationCount: 1,
		},
		{
			kind: "google_project_iam_member",
			name: "foo_project_iam_member",
			diff: &terraform.InstanceDiff{
				Attributes: map[string]*terraform.ResourceAttrDiff{
					"member":  &terraform.ResourceAttrDiff{New: "alice@gmail.com"},
					"project": &terraform.ResourceAttrDiff{New: cfg.project},
					"role":    &terraform.ResourceAttrDiff{New: "role/browser"},
				},
			},
			wantViolationCount: 1,
		},
		{
			kind: "google_compute_firewall",
			name: "foo_compute_firewall",
			diff: &terraform.InstanceDiff{
				Attributes: map[string]*terraform.ResourceAttrDiff{
					"allow.#":                 &terraform.ResourceAttrDiff{New: "2"},
					"allow.12345678.ports.#":  &terraform.ResourceAttrDiff{New: "0"},
					"allow.12345678.protocol": &terraform.ResourceAttrDiff{New: "icmp"},
					"allow.23456789.ports.#":  &terraform.ResourceAttrDiff{New: "3"},
					"allow.23456789.ports.0":  &terraform.ResourceAttrDiff{New: "80"},
					"allow.23456789.ports.1":  &terraform.ResourceAttrDiff{New: "8080"},
					"allow.23456789.ports.2":  &terraform.ResourceAttrDiff{New: "1000-2000"},
					"allow.23456789.protocol": &terraform.ResourceAttrDiff{New: "tcp"},
					"name":                    &terraform.ResourceAttrDiff{New: "test-firewall"},
					"network":                 &terraform.ResourceAttrDiff{New: "test-network"},
					"priority":                &terraform.ResourceAttrDiff{New: "1000"},
					"source_tags.#":           &terraform.ResourceAttrDiff{New: "1"},
					"source_tags.34567890":    &terraform.ResourceAttrDiff{New: "web"},
				},
			},
			wantViolationCount: 2,
		},
		{
			kind: "google_compute_disk",
			name: "foo_compute_disk",
			diff: &terraform.InstanceDiff{
				Attributes: map[string]*terraform.ResourceAttrDiff{
					"image":                     &terraform.ResourceAttrDiff{New: "debian-8-jessie-v20170523"},
					"labels":                    &terraform.ResourceAttrDiff{New: "1"},
					"labels.environment":        &terraform.ResourceAttrDiff{New: "dev"},
					"name":                      &terraform.ResourceAttrDiff{New: "test-disk"},
					"physical_block_size_bytes": &terraform.ResourceAttrDiff{New: "4096"},
					"type":                      &terraform.ResourceAttrDiff{New: "pd-ssd"},
					"zone":                      &terraform.ResourceAttrDiff{New: "us-central1-a"},
				},
			},
			wantViolationCount: 2,
		},
		{
			kind: "google_compute_instance",
			name: "foo_compute_instance",
			diff: &terraform.InstanceDiff{
				Attributes: map[string]*terraform.ResourceAttrDiff{
					"boot_disk.0.auto_delete":               &terraform.ResourceAttrDiff{New: "true"},
					"boot_disk.0.initialize_params.#":       &terraform.ResourceAttrDiff{New: "1"},
					"boot_disk.0.initialize_params.0.image": &terraform.ResourceAttrDiff{New: "debian-cloud/debian-9"},
					"can_ip_forward":                        &terraform.ResourceAttrDiff{New: "false"},
					"deletion_protection":                   &terraform.ResourceAttrDiff{New: "false"},
					"machine_type":                          &terraform.ResourceAttrDiff{New: "n1-standard-1"},
					"metadata.%":                            &terraform.ResourceAttrDiff{New: "1"},
					"metadata.foo":                          &terraform.ResourceAttrDiff{New: "bar"},
					"metadata_startup_script":               &terraform.ResourceAttrDiff{New: "echo hi > /test.txt"},
					"name":                                  &terraform.ResourceAttrDiff{New: "test"},
					"network_interface.#":                   &terraform.ResourceAttrDiff{New: "1"},
					"network_interface.0.access_config.#":   &terraform.ResourceAttrDiff{New: "1"},
					"network_interface.0.access_config.0.assigned_nat_ip": &terraform.ResourceAttrDiff{},
					"network_interface.0.access_config.0.nat_ip":          &terraform.ResourceAttrDiff{},
					"network_interface.0.access_config.0.network_tier":    &terraform.ResourceAttrDiff{},
					"network_interface.0.network":                         &terraform.ResourceAttrDiff{New: "default"},
					"scratch_disk.#":                                      &terraform.ResourceAttrDiff{New: "1"},
					"scratch_disk.0.interface":                            &terraform.ResourceAttrDiff{New: "SCSI"},
					"service_account.#":                                   &terraform.ResourceAttrDiff{New: "1"},
					"service_account.0.scopes.#":                          &terraform.ResourceAttrDiff{New: "3"},
					"service_account.0.scopes.12345678":                   &terraform.ResourceAttrDiff{New: "https://www.googleapis.com/auth/devstorage.read_only", NewExtra: "storage-ro"},
					"service_account.0.scopes.23456789":                   &terraform.ResourceAttrDiff{New: "https://www.googleapis.com/auth/userinfo.email", NewExtra: "userinfo-email"},
					"service_account.0.scopes.34567890":                   &terraform.ResourceAttrDiff{New: "https://www.googleapis.com/auth/compute.readonly", NewExtra: "compute-ro"},
					"tags.#":                                              &terraform.ResourceAttrDiff{New: "2"},
					"tags.45678901":                                       &terraform.ResourceAttrDiff{New: "bar"},
					"tags.56789012":                                       &terraform.ResourceAttrDiff{New: "foo"},
					"zone":                                                &terraform.ResourceAttrDiff{New: "us-central1-a"},
				},
			},
			wantViolationCount: 2,
		},
		{
			kind: "google_storage_bucket",
			name: "foo_storage_bucket",
			diff: &terraform.InstanceDiff{
				Attributes: map[string]*terraform.ResourceAttrDiff{
					"location":                   &terraform.ResourceAttrDiff{New: "EU", NewExtra: "EU"},
					"name":                       &terraform.ResourceAttrDiff{New: "image-store-bucket"},
					"storage_class":              &terraform.ResourceAttrDiff{New: "STANDARD"},
					"website.#":                  &terraform.ResourceAttrDiff{New: "1"},
					"website.0.main_page_suffix": &terraform.ResourceAttrDiff{New: "index.html"},
					"website.0.not_found_page":   &terraform.ResourceAttrDiff{New: "404.html"},
				},
			},
			wantViolationCount: 2,
		},
		{
			kind: "google_sql_database_instance",
			name: "foo_sql_database_instance",
			diff: &terraform.InstanceDiff{
				Attributes: map[string]*terraform.ResourceAttrDiff{
					"database_version":            &terraform.ResourceAttrDiff{New: "POSTGRES_9_6"},
					"name":                        &terraform.ResourceAttrDiff{New: "master-instance"},
					"region":                      &terraform.ResourceAttrDiff{New: "us-central-1"},
					"settings.#":                  &terraform.ResourceAttrDiff{New: "1"},
					"settings.0.disk_autoresize":  &terraform.ResourceAttrDiff{New: "true"},
					"settings.0.pricing_plan":     &terraform.ResourceAttrDiff{New: "PER_USE"},
					"settings.0.replication_type": &terraform.ResourceAttrDiff{New: "SYNCHRONOUS"},
					"settings.0.tier":             &terraform.ResourceAttrDiff{New: "db-f1-micro"},
				},
			},
			wantViolationCount: 2,
		},
	}

	ancestry := "organization/example.com/folder/foobar"
	resourceManager, err := cloudresourcemanager.NewService(context.Background(), option.WithoutAuthentication())
	if err != nil {
		t.Fatalf("constructing resource manager client: %v", err)
	}

	for _, c := range cases {
		t.Run(fmt.Sprintf("%s/%s", c.kind, c.name), func(t *testing.T) {
			converter, err := google.NewConverter(resourceManager, cfg.project, ancestry, cfg.credentials)
			if err != nil {
				t.Fatalf("building google converter: %v", err)
			}

			// Terraform validator requires a base "google_project" resource.
			baseResourceInstanceDiff := &terraform.InstanceDiff{
				Attributes: map[string]*terraform.ResourceAttrDiff{
					"name": &terraform.ResourceAttrDiff{New: "new-foo"},
				},
			}
			baseResource := tfplan.ComposeResource(
				tfplan.Fullpath{"google_project", "foo_project", "root"},
				converter.Schemas()["google_project"].Schema,
				nil,
				baseResourceInstanceDiff,
			)
			if err := converter.AddResource(&baseResource); err != nil {
				t.Fatalf("adding base resource: %s", err)
			}

			path := tfplan.Fullpath{c.kind, c.name, "root"}
			if _, ok := converter.Schemas()[path.Kind]; !ok {
				t.Fatalf("schema %s is not supported", path.Kind)
			}
			schema := converter.Schemas()[path.Kind].Schema
			r := tfplan.ComposeResource(path, schema, nil, c.diff)
			if err := converter.AddResource(&r); err != nil {
				t.Fatalf("adding resource: %s", err)
			}
			if len(converter.Assets()) == 0 {
				t.Fatalf("got zero assets")
			}
			violations, err := tfgcv.ValidateAssets(converter.Assets(), cfg.policy)
			if err != nil {
				t.Fatalf("validating assets: %v", err)
			}
			if len(violations.Violations) != c.wantViolationCount {
				t.Fatalf("got %d violations, want %d", len(violations.Violations), c.wantViolationCount)
			}
		})
	}
}
