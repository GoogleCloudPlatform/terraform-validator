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
              "address": "module.foo.google_compute_instance.quz1"
            }
          ],
          "child_modules": [
            {
              "address": "module.foo.bar",
              "resources": [
                {
                  "address": "module.foo.bar.google_compute_instance.quz2"
                }
              ]
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
    "address": "module.foo.google_compute_instance.quz1",
    "mode": "",
    "module": "foo",
    "name": "",
    "provider_name": "",
    "type": "",
    "values": null
  },
  {
    "address": "module.foo.bar.google_compute_instance.quz2",
    "mode": "",
    "module": "foo.bar",
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
