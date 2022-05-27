package tfgcv

import (
	"context"
	"path/filepath"
	"testing"

	"go.uber.org/zap"
)

const (
	testDataDir      = "../test/read_planned_assets"
	testProjectName  = "gl-akopachevskyy-sql-db"
	testAncestryName = "folders/1234"
)

type args struct {
	file             string
	project          string
	ancestry         string
	convertUnchanged bool
}
type testcase struct {
	name          string
	args          args
	want          int
	wantErr       bool
	ancestryCache map[string]string
}

func testCases(t *testing.T) []testcase {
	var tests []testcase

	tests = append(tests, []testcase{
		{
			"Test TF0_12 and JSON plan",
			args{"tf0_12plan.json", testProjectName, testAncestryName, false},
			2,
			false,
			map[string]string{
				"projects/" + testProjectName: testAncestryName + "/projects/" + testProjectName,
			},
		},
		{
			"Test TF0_12 with all coverage",
			args{"tf0_12plan.allcoverage.json", "foobar", testAncestryName, false},
			9,
			false,
			map[string]string{
				"projects/foobar":         testAncestryName,
				"folders/my-folder":       "organization/test-org",
				"organizations/123456789": "",
			},
		},
		{
			"Test TF0_12 with no-op",
			args{"tf0_12plan.applied.json", "foobar", testAncestryName, true},
			6,
			false,
			map[string]string{
				"projects/345":    testAncestryName,
				"projects/foobar": testAncestryName,
				"folders/567":     "organization/123",
			},
		},
		{
			"Test TF1_0 and JSON plan",
			args{"tf1_0plan.json", testProjectName, testAncestryName, false},
			2,
			false,
			map[string]string{
				"projects/" + testProjectName: testAncestryName,
				"folders/my-folder":           "organization/test-org",
			},
		},
		{
			"Test TF1_0 with all coverage",
			args{"tf1_0plan.allcoverage.json", "foobar", testAncestryName, false},
			9,
			false,
			map[string]string{
				"projects/foobar":   testAncestryName,
				"folders/my-folder": "organization/test-org",
			},
		},
		{
			"Test TF1_0 with no-op",
			args{"tf1_0plan.applied.json", "foobar", testAncestryName, true},
			6,
			false,
			map[string]string{
				"projects/345":    testAncestryName,
				"projects/foobar": testAncestryName,
				"folders/567":     "organization/123",
			},
		},
	}...)
	return tests
}

func TestReadPlannedAssets(t *testing.T) {
	for _, tt := range testCases(t) {
		t.Run(tt.name, func(t *testing.T) {
			testFile := filepath.Join(testDataDir, tt.args.file)
			offline := true
			ctx := context.Background()
			got, err := ReadPlannedAssets(ctx, testFile, tt.args.project, tt.ancestryCache, offline, tt.args.convertUnchanged, zap.NewExample(), "")
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
