package tfgcv

import (
	"context"
	"go.uber.org/zap"
	"path/filepath"
	"testing"
)

const (
	testDataDir      = "../test/read_planned_assets"
	testProjectName  = "gl-akopachevskyy-sql-db"
	testAncestryName = "ancestry"
)

type args struct {
	file             string
	project          string
	ancestry         string
	convertUnchanged bool
}
type testcase struct {
	name    string
	args    args
	want    int
	wantErr bool
}

func testCases(t *testing.T) []testcase {
	var tests []testcase

	tests = append(tests, []testcase{
		{
			"Test TF0_12 and JSON plan",
			args{"tf0_12plan.json", testProjectName, testAncestryName, false},
			2,
			false,
		},
		{
			"Test TF0_12 with all coverage",
			args{"tf0_12plan.allcoverage.json", "foobar", testAncestryName, false},
			9,
			false,
		},
		{
			"Test TF0_12 with no-op",
			args{"tf0_12plan.applied.json", "foobar", testAncestryName, true},
			6,
			false,
		},
		{
			"Test TF1_0 and JSON plan",
			args{"tf1_0plan.json", testProjectName, testAncestryName, false},
			2,
			false,
		},
		{
			"Test TF1_0 with all coverage",
			args{"tf1_0plan.allcoverage.json", "foobar", testAncestryName, false},
			9,
			false,
		},
		{
			"Test TF1_0 with no-op",
			args{"tf1_0plan.applied.json", "foobar", testAncestryName, true},
			6,
			false,
		},
	}...)
	return tests
}

func TestReadPlannedAssets(t *testing.T) {
	for _, tt := range testCases(t) {
		t.Run(tt.name, func(t *testing.T) {
			testFile := filepath.Join(testDataDir, tt.args.file)
			var offline bool
			offline = true
			ctx := context.Background()
			got, err := ReadPlannedAssets(ctx, testFile, tt.args.project, tt.args.ancestry, offline, tt.args.convertUnchanged, zap.NewExample())
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
