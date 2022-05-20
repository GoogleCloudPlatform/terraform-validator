package ancestrymanager

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	resources "github.com/GoogleCloudPlatform/terraform-validator/converters/google/resources"
	"github.com/GoogleCloudPlatform/terraform-validator/tfdata"

	"github.com/google/go-cmp/cmp"
	provider "github.com/hashicorp/terraform-provider-google/google"
	"go.uber.org/zap"
	cloudresourcemanager "google.golang.org/api/cloudresourcemanager/v3"
	"google.golang.org/api/option"
)

func newTestResourceManagerClient(opts []option.ClientOption) *cloudresourcemanager.Service {
	ctx := context.Background()

	rm, err := cloudresourcemanager.NewService(ctx, opts...)
	if err != nil {
		panic(err)
	}
	return rm
}

func TestGetAncestors(t *testing.T) {
	ownerProject := "foo"
	ownerAncestryPath := "organization/qux/folder/bar/project/foo"
	anotherProject := "foo2"

	// Setup a simple test server to mock the response of resource manager.
	responses := map[string]*cloudresourcemanager.Project{
		"projects/12345":     {Name: "projects/12345", Parent: "folders/bar"},
		"projects/foo":       {Name: "projects/foo", Parent: "folders/bar"},
		"folders/bar":        {Name: "folders/bar", Parent: "organizations/qux"},
		"organizations/qux":  {Name: "organizations/qux", Parent: ""},
		"projects/foo2":      {Name: "projects/foo2", Parent: "folders/bar2"},
		"folders/bar2":       {Name: "folders/bar2", Parent: "organizations/qux2"},
		"organizations/qux2": {Name: "organizations/qux2", Parent: ""},
	}
	ts := newAncestryManagerMockServer(t, responses)
	defer ts.Close()

	// option.WithEndpoint(ts.URL), option.WithoutAuthentication()
	rm := newTestResourceManagerClient([]option.ClientOption{option.WithEndpoint(ts.URL), option.WithoutAuthentication()})

	entries := map[string]string{
		ownerProject: ownerAncestryPath,
	}

	p := provider.Provider()

	// offline return errors when the cache cannot cover the request.
	// online return errors when neither cache and mock server cannot cover the request.
	cases := []struct {
		name             string
		data             resources.TerraformResourceData
		asset            *resources.Asset
		want             []string
		parent           string
		wantOnlineError  bool
		wantOfflineError bool
	}{
		{
			name: "owner project - project id",
			data: tfdata.NewFakeResourceData(
				"google_project",
				p.ResourcesMap["google_project"].Schema,
				map[string]interface{}{
					"project_id": ownerProject,
				},
			),
			asset: &resources.Asset{
				Type: "cloudresourcemanager.googleapis.com/Project",
			},
			want:   []string{"projects/foo", "folders/bar", "organizations/qux"},
			parent: "//cloudresourcemanager.googleapis.com/folders/bar",
		},
		{
			name: "owner project - project",
			data: tfdata.NewFakeResourceData(
				"google_project",
				p.ResourcesMap["google_project"].Schema,
				map[string]interface{}{
					"project": ownerProject,
				},
			),
			asset: &resources.Asset{
				Type: "cloudresourcemanager.googleapis.com/Project",
			},
			want:   []string{"projects/foo", "folders/bar", "organizations/qux"},
			parent: "//cloudresourcemanager.googleapis.com/folders/bar",
		},
		{
			name: "owner project - project number",
			data: tfdata.NewFakeResourceData(
				"google_project",
				p.ResourcesMap["google_project"].Schema,
				map[string]interface{}{
					"number": "12345",
				},
			),
			asset: &resources.Asset{
				Type: "cloudresourcemanager.googleapis.com/Project",
			},
			want:             []string{"projects/12345", "folders/bar", "organizations/qux"},
			wantOfflineError: true,
			parent:           "//cloudresourcemanager.googleapis.com/folders/bar",
		},
		{
			name: "owner project - project from config",
			data: tfdata.NewFakeResourceData(
				"google_project",
				p.ResourcesMap["google_project"].Schema,
				map[string]interface{}{},
			),
			asset: &resources.Asset{
				Type: "cloudresourcemanager.googleapis.com/Project",
			},
			want:   []string{"projects/foo", "folders/bar", "organizations/qux"},
			parent: "//cloudresourcemanager.googleapis.com/folders/bar",
		},
		{
			name: "another project",
			data: tfdata.NewFakeResourceData(
				"google_project",
				p.ResourcesMap["google_project"].Schema,
				map[string]interface{}{
					"project_id": anotherProject,
				},
			),
			asset: &resources.Asset{
				Type: "cloudresourcemanager.googleapis.com/Project",
			},
			want:             []string{"projects/foo2", "folders/bar2", "organizations/qux2"},
			wantOfflineError: true,
			parent:           "//cloudresourcemanager.googleapis.com/folders/bar2",
		},
		{
			name: "owner folder",
			data: tfdata.NewFakeResourceData(
				"google_folder_iam_policy",
				p.ResourcesMap["google_folder_iam_policy"].Schema,
				map[string]interface{}{
					"folder": "bar",
				},
			),
			asset: &resources.Asset{
				Type: "cloudresourcemanager.googleapis.com/Folder",
			},
			want:   []string{"folders/bar", "organizations/qux"},
			parent: "//cloudresourcemanager.googleapis.com/organizations/qux",
		},
		{
			name: "owner folder with prefix",
			data: tfdata.NewFakeResourceData(
				"google_folder_iam_policy",
				p.ResourcesMap["google_folder_iam_policy"].Schema,
				map[string]interface{}{
					"folder": "folders/bar",
				},
			),
			asset: &resources.Asset{
				Type: "cloudresourcemanager.googleapis.com/Folder",
			},
			want:   []string{"folders/bar", "organizations/qux"},
			parent: "//cloudresourcemanager.googleapis.com/organizations/qux",
		},
		{
			name: "another folder online",
			data: tfdata.NewFakeResourceData(
				"google_folder_iam_policy",
				p.ResourcesMap["google_folder_iam_policy"].Schema,
				map[string]interface{}{
					"folder": "bar2",
				},
			),
			asset: &resources.Asset{
				Type: "cloudresourcemanager.googleapis.com/Folder",
			},
			want:             []string{"folders/bar2", "organizations/qux2"},
			wantOfflineError: true,
			parent:           "//cloudresourcemanager.googleapis.com/organizations/qux2",
		},
		{
			// Not supporting folder create resource yet.
			name: "not exist folder online",
			data: tfdata.NewFakeResourceData(
				"google_folder_iam_policy",
				p.ResourcesMap["google_folder_iam_policy"].Schema,
				map[string]interface{}{
					"folder": "notexist",
				},
			),
			asset: &resources.Asset{
				Type: "cloudresourcemanager.googleapis.com/Folder",
			},
			wantOfflineError: true,
			wantOnlineError:  true,
		},
		{
			name: "owner org",
			data: tfdata.NewFakeResourceData(
				"google_organization_iam_policy",
				p.ResourcesMap["google_organization_iam_policy"].Schema,
				map[string]interface{}{
					"org_id": "qux",
				},
			),
			asset: &resources.Asset{
				Type: "cloudresourcemanager.googleapis.com/Organization",
			},
			want:   []string{"organizations/qux"},
			parent: "",
		},
		{
			// organization do not have ancestors except itself
			// hence offline also pass.
			name: "another org",
			data: tfdata.NewFakeResourceData(
				"google_organization_iam_policy",
				p.ResourcesMap["google_organization_iam_policy"].Schema,
				map[string]interface{}{
					"org_id": "qux2",
				},
			),
			asset: &resources.Asset{
				Type: "cloudresourcemanager.googleapis.com/Organization",
			},
			want:   []string{"organizations/qux2"},
			parent: "",
		},
		{
			name: "other resource with owner project",
			data: tfdata.NewFakeResourceData(
				"google_compute_disk",
				p.ResourcesMap["google_compute_disk"].Schema,
				map[string]interface{}{
					"project": ownerProject,
				},
			),
			asset: &resources.Asset{
				Type: "cloudresourcemanager.googleapis.com/Disk",
			},
			want:   []string{"projects/foo", "folders/bar", "organizations/qux"},
			parent: "//cloudresourcemanager.googleapis.com/projects/foo",
		},
		{
			name: "other resource online with another project",
			data: tfdata.NewFakeResourceData(
				"google_compute_disk",
				p.ResourcesMap["google_compute_disk"].Schema,
				map[string]interface{}{
					"project": anotherProject,
				},
			),
			asset: &resources.Asset{
				Type: "cloudresourcemanager.googleapis.com/Disk",
			},
			want:             []string{"projects/foo2", "folders/bar2", "organizations/qux2"},
			wantOfflineError: true,
			parent:           "//cloudresourcemanager.googleapis.com/projects/foo2",
		},
		{
			name: "custom role with org",
			data: tfdata.NewFakeResourceData(
				"google_organization_iam_custom_role",
				p.ResourcesMap["google_organization_iam_custom_role"].Schema,
				map[string]interface{}{
					"org_id": "qux",
				},
			),
			asset: &resources.Asset{
				Type: "iam.googleapis.com/Role",
			},
			want:   []string{"organizations/qux"},
			parent: "//cloudresourcemanager.googleapis.com/organizations/qux",
		},
		{
			name: "custom role with project",
			data: tfdata.NewFakeResourceData(
				"google_project_iam_custom_role",
				p.ResourcesMap["google_project_iam_custom_role"].Schema,
				map[string]interface{}{
					"project": "foo",
				},
			),
			asset: &resources.Asset{
				Type: "iam.googleapis.com/Role",
			},
			want:   []string{"projects/foo", "folders/bar", "organizations/qux"},
			parent: "//cloudresourcemanager.googleapis.com/projects/foo",
		},
		{
			name: "new project in folder",
			data: tfdata.NewFakeResourceData(
				"google_project",
				p.ResourcesMap["google_project"].Schema,
				map[string]interface{}{
					"folder_id":  "bar",
					"project_id": "new-project",
				},
			),
			asset: &resources.Asset{
				Type: "cloudresourcemanager.googleapis.com/Project",
			},
			want:   []string{"projects/new-project", "folders/bar", "organizations/qux"},
			parent: "//cloudresourcemanager.googleapis.com/folders/bar",
		},
		{
			name: "new project in organization",
			data: tfdata.NewFakeResourceData(
				"google_project",
				p.ResourcesMap["google_project"].Schema,
				map[string]interface{}{
					"org_id":     "qux",
					"project_id": "new-project",
				},
			),
			asset: &resources.Asset{
				Type: "cloudresourcemanager.googleapis.com/Project",
			},
			want:   []string{"projects/new-project", "organizations/qux"},
			parent: "//cloudresourcemanager.googleapis.com/organizations/qux",
		},
		{
			// for new projects, if it cannot find ancestors in online mode,
			// it just returns the project itself as ancestors.
			// offline will fail because no cloud resource manager.
			name: "new project without org_id or folder_id",
			data: tfdata.NewFakeResourceData(
				"google_project",
				p.ResourcesMap["google_project"].Schema,
				map[string]interface{}{
					"project_id": "new-project",
				},
			),
			asset: &resources.Asset{
				Type: "cloudresourcemanager.googleapis.com/Project",
			},
			want:             []string{"projects/new-project"},
			wantOfflineError: true,
			parent:           "//cloudresourcemanager.googleapis.com/projects/new-project",
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			cfg := &resources.Config{
				Project: ownerProject,
			}
			amOnline, err := newManager(rm, entries, zap.NewExample())
			if err != nil {
				t.Fatalf("failed to create online ancestry manager: %s", err)
			}
			got, parent, err := amOnline.Ancestors(cfg, c.data, c.asset)
			if c.wantOnlineError {
				if err == nil {
					t.Fatalf("onlineMgr.Ancestors(%v, %v, %v) = nil, want = err", cfg, c.data, c.asset)
				}
			} else {
				if err != nil {
					t.Fatalf("onlineMgr.Ancestors(%v, %v, %v) = %s, want = nil", cfg, c.data, c.asset, err)
				}
				if parent != c.parent {
					t.Errorf("onlineMgr.Ancestors(%v, %v, %v) parent = %s, want = %s", cfg, c.data, c.asset, parent, c.parent)
				}
				if diff := cmp.Diff(c.want, got); diff != "" {
					t.Errorf("onlineMgr.Ancestors(%v, %v, %v) returned unexpected diff (-want +got):\n%s", cfg, c.data, c.asset, diff)
				}
			}

			amOffline, err := newManager(nil, entries, zap.NewExample())
			if err != nil {
				t.Fatalf("failed to create offline ancestry manager: %s", err)
			}
			got, parent, err = amOffline.Ancestors(cfg, c.data, c.asset)
			if c.wantOfflineError {
				if err == nil {
					t.Fatalf("offlineMgr.Ancestors(%v, %v, %v) = nil, want = err", cfg, c.data, c.asset)
				}
			} else {
				if err != nil {
					t.Fatalf("offlineMgr.Ancestors(%v, %v, %v) = %s, want = nil", cfg, c.data, c.asset, err)
				}
				if parent != c.parent {
					t.Errorf("offlineMgr.Ancestors(%v, %v, %v) parent = %s, want = %s", cfg, c.data, c.asset, parent, c.parent)
				}
				if diff := cmp.Diff(c.want, got); diff != "" {
					t.Errorf("offlineMgr.Ancestors(%v, %v, %v) returned unexpected diff (-want +got):\n%s", cfg, c.data, c.asset, diff)
				}
			}
		})
	}
}

func TestGetAncestorsWithCache(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		cache     map[string][]string
		responses map[string]*cloudresourcemanager.Project
		want      []string
		wantCache map[string][]string
	}{
		{
			name:  "empty cache",
			input: "projects/abc",
			cache: make(map[string][]string),
			responses: map[string]*cloudresourcemanager.Project{
				"projects/abc": {Name: "projects/123", ProjectId: "abc", Parent: "folders/456"},
				"folders/456":  {Name: "folders/456", Parent: "folders/789"},
				"folders/789":  {Name: "folders/789", Parent: "organizations/321"},
			},
			want: []string{"projects/123", "folders/456", "folders/789", "organizations/321"},
			wantCache: map[string][]string{
				"projects/abc":      {"projects/123", "folders/456", "folders/789", "organizations/321"},
				"projects/123":      {"projects/123", "folders/456", "folders/789", "organizations/321"},
				"folders/456":       {"folders/456", "folders/789", "organizations/321"},
				"folders/789":       {"folders/789", "organizations/321"},
				"organizations/321": {"organizations/321"},
			},
		},
		{
			name:  "partial cache",
			input: "projects/abc",
			cache: map[string][]string{
				"folders/789":       {"folders/789", "organizations/321"},
				"organizations/321": {"organizations/321"},
			},
			responses: map[string]*cloudresourcemanager.Project{
				"projects/abc": {Name: "projects/123", ProjectId: "abc", Parent: "folders/456"},
				"folders/456":  {Name: "folders/456", Parent: "folders/789"},
			},
			want: []string{"projects/123", "folders/456", "folders/789", "organizations/321"},
			wantCache: map[string][]string{
				"projects/abc":      {"projects/123", "folders/456", "folders/789", "organizations/321"},
				"projects/123":      {"projects/123", "folders/456", "folders/789", "organizations/321"},
				"folders/456":       {"folders/456", "folders/789", "organizations/321"},
				"folders/789":       {"folders/789", "organizations/321"},
				"organizations/321": {"organizations/321"},
			},
		},
		{
			name:  "all response from cache",
			input: "projects/abc",
			cache: map[string][]string{
				"projects/abc": {"projects/123", "folders/456", "folders/789", "organizations/321"},
			},
			responses: map[string]*cloudresourcemanager.Project{},
			want:      []string{"projects/123", "folders/456", "folders/789", "organizations/321"},
			wantCache: map[string][]string{
				"projects/abc":      {"projects/123", "folders/456", "folders/789", "organizations/321"},
				"projects/123":      {"projects/123", "folders/456", "folders/789", "organizations/321"},
				"folders/456":       {"folders/456", "folders/789", "organizations/321"},
				"folders/789":       {"folders/789", "organizations/321"},
				"organizations/321": {"organizations/321"},
			},
		},
		{
			name:      "organization",
			input:     "organizations/321",
			cache:     map[string][]string{},
			responses: map[string]*cloudresourcemanager.Project{},
			want:      []string{"organizations/321"},
			wantCache: map[string][]string{
				"organizations/321": {"organizations/321"},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ts := newAncestryManagerMockServer(t, test.responses)
			defer ts.Close()

			m := &manager{
				errorLogger:     zap.NewExample(),
				resourceManager: newTestResourceManagerClient([]option.ClientOption{option.WithEndpoint(ts.URL), option.WithoutAuthentication()}),
				ancestorCache:   test.cache,
			}
			got, err := m.getAncestorsWithCache(test.input)
			if err != nil {
				t.Fatalf("getAncestorsWithCache(%s) = %s, want = nil", test.input, err)
			}
			if diff := cmp.Diff(test.want, got); diff != "" {
				t.Errorf("getAncestorsWithCache(%v) returned unexpected diff (-want +got):\n%s", test.input, diff)
			}
			if diff := cmp.Diff(test.wantCache, m.ancestorCache); diff != "" {
				t.Errorf("getAncestorsWithCache(%v) cache returned unexpected diff (-want +got):\n%s", test.input, diff)
			}
		})
	}
}

func TestGetAncestorsWithCache_Fail(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		cache     map[string][]string
		responses map[string]*cloudresourcemanager.Project
		wantErr   string
	}{
		{
			name:  "no parent response",
			input: "projects/abc",
			cache: make(map[string][]string),
			responses: map[string]*cloudresourcemanager.Project{
				"projects/abc": {Name: "projects/123", ProjectId: "projects/abc", Parent: "folders/not-exist"},
			},
			wantErr: "no response",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ts := newAncestryManagerMockServer(t, test.responses)
			defer ts.Close()

			crm := newTestResourceManagerClient(
				[]option.ClientOption{
					option.WithEndpoint(ts.URL),
					option.WithoutAuthentication(),
				},
			)
			m := &manager{
				errorLogger:     zap.NewExample(),
				resourceManager: crm,
				ancestorCache:   test.cache,
			}
			_, err := m.getAncestorsWithCache(test.input)
			if err == nil {
				t.Fatalf("getAncestorsWithCache(%s) = nil, want = %s", test.input, test.wantErr)
			}
		})
	}
}

func newAncestryManagerMockServer(t *testing.T, responses map[string]*cloudresourcemanager.Project) *httptest.Server {
	t.Helper()
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		name := strings.TrimPrefix(r.URL.Path, "/v3/")
		resp, ok := responses[name]
		if !ok {
			w.WriteHeader(http.StatusForbidden)
			w.Write([]byte(fmt.Sprintf("no response for request path %s", "/v3/"+name)))
			return
		}
		payload, err := resp.MarshalJSON()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("failed to MarshalJSON: %s", err)))
			return
		}
		w.Write(payload)
	}))
	return ts
}

func TestParseAncestryPath(t *testing.T) {
	tests := []struct {
		name string
		path string
		want []string
	}{
		{
			name: "all kinds of resource",
			path: "organizations/123/folders/456/projects/789",
			want: []string{"projects/789", "folders/456", "organizations/123"},
		},
		{
			name: "multiple folders",
			path: "organizations/123/folders/456/folders/789",
			want: []string{"folders/789", "folders/456", "organizations/123"},
		},
		{
			name: "normalize resource name",
			path: "organization/123/folder/456/project/789",
			want: []string{"projects/789", "folders/456", "organizations/123"},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := parseAncestryPath(test.path)
			if err != nil {
				t.Fatalf("parseAncestryPath(%s) = %s, want = nil", test.path, err)
			}
			if diff := cmp.Diff(test.want, got); diff != "" {
				t.Errorf("parseAncestryPath(%v) returned unexpected diff (-want +got):\n%s", test.path, diff)
			}
		})
	}
}

func TestParseAncestryPath_Fail(t *testing.T) {
	tests := []struct {
		name    string
		path    string
		wantErr string
	}{
		{
			name:    "malform with single word",
			path:    "organizations",
			wantErr: "unexpected format",
		},
		{
			name:    "malform",
			path:    "organizations/123/folders",
			wantErr: "unexpected format",
		},
		{
			name:    "invalid keyword",
			path:    "org/123/folders/123",
			wantErr: "invalid ancestry path",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, err := parseAncestryPath(test.path)
			if err == nil {
				t.Fatalf("parseAncestryPath(%s) = nil, want = %s", test.path, test.wantErr)
			}
		})
	}
}

func TestInitAncestryCache(t *testing.T) {
	tests := []struct {
		name    string
		entries map[string]string
		want    map[string][]string
	}{
		{
			name: "empty ancestry",
			entries: map[string]string{
				"test-proj": "",
			},
			want: map[string][]string{},
		},
		{
			name: "empty key",
			entries: map[string]string{
				"": "organizations/123/folders/345",
			},
			want: map[string][]string{},
		},
		{
			name: "default key to project",
			entries: map[string]string{
				"test-proj": "organizations/123/folders/345",
			},
			want: map[string][]string{
				"projects/test-proj": {"projects/test-proj", "folders/345", "organizations/123"},
				"folders/345":        {"folders/345", "organizations/123"},
				"organizations/123":  {"organizations/123"},
			},
		},
		{
			name: "key has prefix folders/",
			entries: map[string]string{
				"folders/345": "organizations/123",
			},
			want: map[string][]string{
				"folders/345":       {"folders/345", "organizations/123"},
				"organizations/123": {"organizations/123"},
			},
		},
		{
			name: "key has prefix projects/",
			entries: map[string]string{
				"projects/test-proj": "organizations/123",
			},
			want: map[string][]string{
				"projects/test-proj": {"projects/test-proj", "organizations/123"},
				"organizations/123":  {"organizations/123"},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			m := &manager{
				ancestorCache: make(map[string][]string),
			}
			err := m.initAncestryCache(test.entries)
			if err != nil {
				t.Fatalf("initAncestryCache(%v) = %s, want = nil", test.entries, err)
			}
			if diff := cmp.Diff(test.want, m.ancestorCache); diff != "" {
				t.Errorf("initAncestryCache(%v) returned unexpected diff (-want +got):\n%s", test.entries, diff)
			}
		})
	}
}

func TestInitAncestryCache_Fail(t *testing.T) {
	tests := []struct {
		name    string
		entries map[string]string
	}{
		{
			name: "typo",
			entries: map[string]string{
				"foldres/def": "organizations/123",
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			m := &manager{
				ancestorCache: make(map[string][]string),
			}
			err := m.initAncestryCache(test.entries)
			if err == nil {
				t.Fatalf("initAncestryCache(%v) = nil, want = err", test.entries)
			}
		})
	}
}

func TestParseAncestryKey(t *testing.T) {
	tests := []struct {
		name string
		key  string
		want string
	}{
		{
			name: "not contain /",
			key:  "proj",
			want: "projects/proj",
		},
		{
			name: "contain projects/",
			key:  "projects/1",
			want: "projects/1",
		},
		{
			name: "contain folders/",
			key:  "folders/1",
			want: "folders/1",
		},
		{
			name: "contain organizations/",
			key:  "organizations/1",
			want: "organizations/1",
		},
		{
			name: "contain project/",
			key:  "project/1",
			want: "projects/1",
		},
		{
			name: "contain folder/",
			key:  "folder/1",
			want: "folders/1",
		},
		{
			name: "contain organization/",
			key:  "organization/1",
			want: "organizations/1",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := parseAncestryKey(test.key)
			if err != nil {
				t.Fatalf("parseAncestryKey(%v) = %v, want = nil", test.key, err)
			}
			if got != test.want {
				t.Errorf("parseAncestryKey(%v) = %v, want = %v", test.key, got, test.want)
			}
		})
	}
}

func TestParseAncestryKey_Fail(t *testing.T) {
	tests := []struct {
		name string
		key  string
	}{
		{
			name: "invalid spell",
			key:  "org/1",
		},
		{
			name: "multiple /",
			key:  "folders/123/folders/456",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := parseAncestryKey(test.key)
			if err == nil {
				t.Fatalf("parseAncestryKey(%v) = %v, want error", test.key, got)
			}
		})
	}
}
