package tfplan

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReadJSONResources(t *testing.T) {
	data := []byte(`
{
  "planned_values": {
    "root_module": {
      "child_modules": [
        {
          "address": "module.foo",
          "resources": [
            {
              "address": "module.foo.google_compute_instance.bar"
            }
          ]
        }
      ]
    }
  }
}
`)
	wantJSON := []byte(`
[
  {
    "address": "module.foo.google_compute_instance.bar",
    "mode": "",
    "module": "foo",
    "name": "",
    "provider_name": "",
    "type": "",
    "values": null
  }
]
`)
	got, err := readJSONResources(data)
	if err != nil {
		t.Fatalf("got error: %v", err)
	}
	gotJSON, err := json.Marshal(got)
	if err != nil {
		t.Fatalf("marshaling: %v", err)
	}
	require.JSONEq(t, string(wantJSON), string(gotJSON))
}
