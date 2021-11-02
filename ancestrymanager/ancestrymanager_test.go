package ancestrymanager

import (
	"context"
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"

	"go.uber.org/zap"
	cloudresourcemanager "google.golang.org/api/cloudresourcemanager/v1"
	"google.golang.org/api/option"
)

type testRetriever struct {
	online bool
	opts   []option.ClientOption
}

func (tr *testRetriever) NewResourceManagerClient(userAgent string) *cloudresourcemanager.Service {
	ctx := context.Background()

	if tr.online {
		rm, err := cloudresourcemanager.NewService(ctx, tr.opts...)
		if err != nil {
			panic(err)
		}
		return rm
	}
	return nil
}

func TestAncestryPath(t *testing.T) {
	cases := []struct {
		name           string
		input          []*cloudresourcemanager.Ancestor
		expectedOutput string
	}{
		{
			name:           "Empty",
			input:          []*cloudresourcemanager.Ancestor{},
			expectedOutput: "",
		},
		{
			name: "ProjectOrganization",
			input: []*cloudresourcemanager.Ancestor{
				{
					ResourceId: &cloudresourcemanager.ResourceId{
						Id:   "my-prj",
						Type: "project",
					},
				},
				{
					ResourceId: &cloudresourcemanager.ResourceId{
						Id:   "my-org",
						Type: "organization",
					},
				},
			},
			expectedOutput: "organization/my-org/project/my-prj",
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			output := ancestryPath(c.input)
			if output != c.expectedOutput {
				t.Errorf("expected output %q, got %q", c.expectedOutput, output)
			}
		})
	}
}

func TestGetAncestry(t *testing.T) {
	ownerProject := "foo"
	ownerAncestry := "organization/qux/folder/bar"
	ownerAncestryPath := "organization/qux/folder/bar/project/foo"
	anotherProject := "foo2"

	// Setup a simple test server to mock the response of resource manager.
	cache := map[string][]*cloudresourcemanager.Ancestor{
		ownerProject: []*cloudresourcemanager.Ancestor{
			{ResourceId: &cloudresourcemanager.ResourceId{Id: "foo", Type: "project"}},
			{ResourceId: &cloudresourcemanager.ResourceId{Id: "bar", Type: "folder"}},
			{ResourceId: &cloudresourcemanager.ResourceId{Id: "qux", Type: "organization"}},
		},
		anotherProject: []*cloudresourcemanager.Ancestor{
			{ResourceId: &cloudresourcemanager.ResourceId{Id: "foo2", Type: "project"}},
			{ResourceId: &cloudresourcemanager.ResourceId{Id: "bar2", Type: "folder"}},
			{ResourceId: &cloudresourcemanager.ResourceId{Id: "qux2", Type: "organization"}},
		},
	}
	ts := newAncestryManagerMockServer(t, cache)
	defer ts.Close()

	// option.WithEndpoint(ts.URL), option.WithoutAuthentication()
	trOnline := &testRetriever{online: true, opts: []option.ClientOption{option.WithEndpoint(ts.URL), option.WithoutAuthentication()}}

	entries := map[string]string{
		ownerProject: ownerAncestry,
	}
	amOnline, err := New(trOnline, entries, "", false, zap.NewExample())
	if err != nil {
		t.Fatalf("failed to create online ancestry manager: %s", err)
	}

	trOffline := &testRetriever{online: false}
	amOffline, err := New(trOffline, entries, "", true, zap.NewExample())
	if err != nil {
		t.Fatalf("failed to create offline ancestry manager: %s", err)
	}

	cases := []struct {
		name      string
		target    AncestryManager
		query     string
		wantError bool
		want      string
	}{
		{name: "owner_project_online", target: amOnline, query: ownerProject, want: ownerAncestryPath},
		{name: "owner_project_offline", target: amOffline, query: ownerProject, want: ownerAncestryPath},
		{name: "another_project_online", target: amOnline, query: anotherProject, want: "organization/qux2/folder/bar2/project/foo2"},
		{name: "another_project_offline", target: amOffline, query: anotherProject, wantError: true},
		{name: "missed_project_online", target: amOnline, query: "notexist", wantError: true},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got, err := c.target.GetAncestry(c.query)
			if !c.wantError && err != nil {
				t.Fatalf("GetAncestry(%s) returns error: %s", c.query, err)
			}
			if c.wantError && err == nil {
				t.Fatalf("GetAncestry(%s) returns no error, want error", c.query)
			}
			if got != c.want {
				t.Errorf("GetAncestry(%s): got=%s, want=%s", c.query, got, c.want)
			}
		})
	}
}

func newAncestryManagerMockServer(t *testing.T, cache map[string][]*cloudresourcemanager.Ancestor) *httptest.Server {
	t.Helper()
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		re := regexp.MustCompile(`([^/]*):getAncestry`)
		path := re.FindStringSubmatch(r.URL.Path)
		if path == nil || cache[path[1]] == nil {
			w.WriteHeader(http.StatusForbidden)
			return
		}
		payload, err := (&cloudresourcemanager.GetAncestryResponse{Ancestor: cache[path[1]]}).MarshalJSON()
		if err != nil {
			t.Errorf("failed to MarshalJSON: %s", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Write(payload)
	}))
	return ts
}
