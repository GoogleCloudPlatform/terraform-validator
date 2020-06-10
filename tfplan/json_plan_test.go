package tfplan

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

const (
	testDataDir   = "../testdata/plans"
	testInputDir  = testDataDir + "/input/"
	testOutputDir = testDataDir + "/output/"
)

func TestReadJSONResources(t *testing.T) {
	type testcase struct {
		name     string
		filename string
	}

	var tests = []testcase{
		{
			"Test modules",
			"modules.json",
		},
		{
			"Computed references",
			"references.json",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			inputFile := filepath.Join(testInputDir, tt.filename)
			outputFile := filepath.Join(testOutputDir, tt.filename)

			input, err := ioutil.ReadFile(inputFile)
			if err != nil {
				t.Errorf("Failed to load test file %s", inputFile)
			}
			wantJSON, err := ioutil.ReadFile(outputFile)
			if err != nil {
				t.Errorf("Failed to load test file %s", outputFile)
			}

			got, err := readJSONResources(input)

			if err != nil {
				t.Fatalf("got error: %v", err)
			}

			gotJSON, err := json.Marshal(got)
			fmt.Print(string(gotJSON))
			if err != nil {
				t.Fatalf("marshaling: %v", err)
			}

			require.JSONEq(t, string(wantJSON), string(gotJSON))
		})
	}
}
