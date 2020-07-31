package policy

import (
	"path/filepath"
	"sort"
	"testing"

	"github.com/GoogleCloudPlatform/terraform-validator/converters/google"
	"github.com/google/go-cmp/cmp"
)

const (
	testDataDir        = "./testdata"
	pathPrefix         = "//cloudresourcemanager.googleapis.com"
	testOrgID          = "123456789"
	testOrgPath        = pathPrefix + "/organizations/" + testOrgID
	testFolderName     = "example-tf-folder"
	testFolderID       = "123456789"
	testFolderNamePath = pathPrefix + "/" + testFolderName
	testFolderIDPath   = pathPrefix + "/folders/" + testFolderID
	testProjectName    = "example-tf-project"
	testProjectPath    = pathPrefix + "/projects/" + testProjectName
	testAncestryName   = "ancestry"
)

// TestOverlayWithCreates verifies that the overlay handles iam resource creation correctly.
func TestOverlayWithCreates(t *testing.T) {
	policy := func() *google.IAMPolicy {
		return &google.IAMPolicy{
			Bindings: []google.IAMBinding{
				{
					Role:    "roles/bigquery.admin",
					Members: []string{"user:user2@example.com"},
				},
				{
					Role:    "roles/editor",
					Members: []string{"user:user1@example.com"},
				},
			},
		}
	}
	wantOverlay := Overlay{
		testOrgPath:        policy(),
		testFolderNamePath: policy(),
		testProjectPath:    policy(),
	}

	tests := []struct {
		name string
		file string
	}{
		{
			name: "binding_member",
			file: "tf12plan_create_binding_member.json",
		},
		{
			name: "iam_policy",
			file: "tf12plan_create_iam_policy.json",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			file := filepath.Join(testDataDir, tt.file)
			got, err := BuildOverlay(file, testProjectName, testAncestryName, true)
			sortOverlay(got)
			if err != nil {
				t.Fatalf("BuildOverly(%s) error: %v", file, err)
			}
			if diff := cmp.Diff(wantOverlay, got); diff != "" {
				t.Fatalf("BuildOverlay(%s)) diff (-want +got): %s", file, diff)
			}
		})
	}
}

// TestOverlayWithUpdates verifies that the overlay handles iam resource creation correctly.
func TestOverlayWithUpdates(t *testing.T) {
	policy := func() *google.IAMPolicy {
		return &google.IAMPolicy{
			Bindings: []google.IAMBinding{
				{
					Role:    "roles/bigquery.admin",
					Members: []string{"user:user4@example.com"},
				},
				{
					Role:    "roles/editor",
					Members: []string{"user:user3@example.com"},
				},
			},
		}
	}
	wantOverlay := Overlay{
		testOrgPath:      policy(),
		testFolderIDPath: policy(),
		testProjectPath:  policy(),
	}

	tests := []struct {
		name string
		file string
	}{
		{
			name: "binding_member",
			file: "tf12plan_update_binding_member.json",
		},
		{
			name: "iam_policy",
			file: "tf12plan_update_iam_policy.json",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			file := filepath.Join(testDataDir, tt.file)
			got, err := BuildOverlay(file, testProjectName, testAncestryName, true)
			sortOverlay(got)
			if err != nil {
				t.Fatalf("BuildOverly(%s) error: %v", file, err)
			}
			if diff := cmp.Diff(wantOverlay, got); diff != "" {
				t.Fatalf("BuildOverlay(%s)) diff (-want +got): %s", file, diff)
			}
		})
	}
}

// TestOverlayWithDeletes verifies that the overlay handles iam resource creation correctly.
func TestOverlayWithDeletes(t *testing.T) {
	policy := func() *google.IAMPolicy {
		return &google.IAMPolicy{
			Bindings: []google.IAMBinding{},
		}
	}
	wantOverlay := Overlay{
		testOrgPath:      policy(),
		testFolderIDPath: policy(),
		testProjectPath:  policy(),
	}

	tests := []struct {
		name string
		file string
	}{
		{
			name: "binding_member",
			file: "tf12plan_delete_binding_member.json",
		},
		{
			name: "iam_policy",
			file: "tf12plan_delete_iam_policy.json",
		},
		{
			name: "delete_resources",
			file: "tf12plan_delete_resources.json",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			file := filepath.Join(testDataDir, tt.file)
			got, err := BuildOverlay(file, testProjectName, testAncestryName, true)
			sortOverlay(got)
			if err != nil {
				t.Fatalf("BuildOverly(%s) error: %v", file, err)
			}
			if diff := cmp.Diff(wantOverlay, got); diff != "" {
				t.Fatalf("BuildOverlay(%s)) diff (-want +got): %s", file, diff)
			}
		})
	}
}

// sortBindings provides a deterministic sorting of all fields in the Overlay.
// Sorting is done in place.
func sortOverlay(overlay Overlay) {
	for _, policy := range overlay {
		if policy == nil {
			continue
		}
		// Sort the bindings themselves by role.
		sort.Slice(policy.Bindings, func(i, j int) bool {
			return policy.Bindings[i].Role < policy.Bindings[j].Role
		})
		// For each binding, ensure members are ordered.
		for _, b := range policy.Bindings {
			sort.Strings(b.Members)
		}
	}
}
