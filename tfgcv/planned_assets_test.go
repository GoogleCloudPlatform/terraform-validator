package tfgcv

import (
	"testing"

	"github.com/GoogleCloudPlatform/terraform-validator/tfplan"
)

const TF11PLAN string = "../test/read_planned_assets/tf11plan.tfplan"
const TF12PLAN string = "../test/read_planned_assets/tf12plan.json"

func TestReadPlannedAssets(t *testing.T) {
	type args struct {
		path      string
		project   string
		ancestry  string
		tfVersion string
	}
	tests := []struct {
		name    string
		args    args
		want    int
		wantErr bool
	}{
		{
			"Test TF12 and JSON plan",
			args{TF12PLAN, "prj", "ancsetry", tfplan.TF12},
			2,
			false,
		},
		{
			"Test TF12 and binary plan",
			args{TF12PLAN, "prj", "ancsetry", tfplan.TF11},
			0,
			true,
		},
		{
			"Test TF11 and JSON plan",
			args{TF11PLAN, "prj", "ancsetry", tfplan.TF12},
			0,
			true,
		},
		{
			"Test TF11 and binary plan",
			args{TF12PLAN, "prj", "ancsetry", tfplan.TF12},
			2,
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ReadPlannedAssets(tt.args.path, tt.args.project, tt.args.ancestry, tt.args.tfVersion)
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
