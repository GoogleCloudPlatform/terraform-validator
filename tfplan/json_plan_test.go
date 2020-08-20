package tfplan

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

func newPlan(t *testing.T) jsonPlan {
	t.Helper()
	data := []byte(`
{
  "planned_values": {
    "root_module": {
      "child_modules": [
        {
          "address": "module.foo",
          "resources": [
            {
              "address": "module.foo.google_compute_instance.quz1",
              "values": {
								"key1": "value1",
								"nestedKey1": { "insideKey1": "insideValue1"}
              }
            }
          ],
          "child_modules": [
            {
              "address": "module.foo.bar",
              "resources": [
                {
                  "address": "module.foo.bar.google_compute_instance.quz2",
									"values": {"key2": "value2"}
                },
								{
									"address": "module.foo.bar.google_compute_instance.quz4",
									"values": {"key4": "value4"}
								}  
              ]
            }
          ]
        }
      ]
    }
  },
	"resource_changes": [
		{
			"address": "module.foo.google_compute_instance.quz1",
			"change": {
				"actions": ["delete", "create"],
				"before": {"key1": "value1"},
				"after": {
					"key1": "value1",
					"nestedKey1": { "insideKey1": "insideValue1"}
				}
			}
		},
		{
			"address": "module.foo.bar.google_compute_instance.quz2",
			"change": {
				"actions": ["noop"],
				"before": {"key2": "value2"},
				"after": {"key2": "value2"}
			}
		},
		{
			"address": "module.foo.bar.google_compute_instance.quz3",
			"change": {
				"actions": ["delete"],
				"before": {"key3": "value3"},
				"after": {}
			}
		},  
		{
			"address": "module.foo.bar.google_compute_instance.quz4",
			"change": {
				"actions": ["create"],
				"before": {},
				"after": {"key4": "value4"}
			}
		}  
	]
}
`)
	plan := jsonPlan{}
	err := json.Unmarshal(data, &plan)
	if err != nil {
		t.Fatalf("parsing %s: %v", string(data), err)
	}
	return plan
}

func TestReadPlannedJSONResources(t *testing.T) {
	wantJSON := []byte(`
[
  {
    "address": "module.foo.google_compute_instance.quz1",
    "mode": "",
    "module": "foo",
    "name": "",
    "provider_name": "",
    "type": "",
		"values": {
			"key1": "value1",
			"nestedKey1": { "insideKey1": "insideValue1"}
		}
  },
  {
    "address": "module.foo.bar.google_compute_instance.quz2",
    "mode": "",
    "module": "foo.bar",
    "name": "",
    "provider_name": "",
    "type": "",
		"values": {"key2": "value2"}
  },
  {
    "address": "module.foo.bar.google_compute_instance.quz4",
    "mode": "",
    "module": "foo.bar",
    "name": "",
    "provider_name": "",
    "type": "",
		"values": {"key4": "value4"}
  }
]
`)
	got := readPlannedJSONResources(newPlan(t))
	gotJSON, err := json.Marshal(got)
	if err != nil {
		t.Fatalf("marshaling: %v", err)
	}
	require.JSONEq(t, string(wantJSON), string(gotJSON))
}

func TestReadCurrentJSONResource(t *testing.T) {
	wantJSON := []byte(`
[
  {
    "address": "module.foo.google_compute_instance.quz1",
    "mode": "",
    "module": "",
    "name": "",
    "provider_name": "",
    "type": "",
    "values": {"key1": "value1"}
  },
  {
    "address": "module.foo.bar.google_compute_instance.quz2",
    "mode": "",
    "module": "",
    "name": "",
    "provider_name": "",
    "type": "",
    "values": {"key2": "value2"}
  },
  {
    "address": "module.foo.bar.google_compute_instance.quz3",
    "mode": "",
    "module": "",
    "name": "",
    "provider_name": "",
    "type": "",
		"values": {"key3": "value3"}
  }
]
`)
	got := readCurrentJSONResources(newPlan(t))
	gotJSON, err := json.Marshal(got)
	if err != nil {
		t.Fatalf("marshaling: %v", err)
	}
	require.JSONEq(t, string(wantJSON), string(gotJSON))
}
