package warden

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestPrimitiveQuery(t *testing.T) {
	type skenas struct {
		payload  map[string]interface{}
		matcher  []string
		expected bool
		title    string
	}

	for i, skena := range []skenas{
		{
			payload: jsonToGeneric(` {"user_id": "1"} `),
			matcher: []string{
				"a/user_id/si:2,3",
			},
			expected: false,
		},
		{
			payload: jsonToGeneric(` {"user_id": "3"} `),
			matcher: []string{
				"a/user_id/si:2,3",
			},
			expected: true,
		},
		{
			payload: jsonToGeneric(` {"user_id": "1"} `),
			matcher: []string{
				"a/user_id/sq:1",
			},
			expected: true,
		},
		{
			payload: jsonToGeneric(` {"user_id": "1", "active" : true} `),
			matcher: []string{
				"a/active/bt:t",
			},
			expected: true,
		},
		{
			payload: jsonToGeneric(` {"user_id": "1", "active" : false, "block" : false} `),
			matcher: []string{
				"a/active/bf:f",
				"a/block/bf:t",
			},
			expected: true,
			title:    "multiple matcher allow",
		},
		{
			payload: jsonToGeneric(` {"user_id": "1", "active" : false, "block" : false} `),
			matcher: []string{
				"a/active/bf:f",
				"d/block/bf:t",
			},
			expected: false,
			title:    "multiple matcher combine allow & deny",
		},
		{
			payload: jsonToGeneric(` {"user_id": "1", "active" : false, "block" : true} `),
			matcher: []string{
				"a/active/bf:f",
				"d/block/bf:t",
			},
			expected: true,
			title:    "multiple matcher combine allow & deny, deny unmatch",
		},
		{
			payload: jsonToGeneric(` {"user_id": "1", "active" : false, "block" : true} `),
			matcher: []string{
				"a/active/bf/user_id/sq:1/block/bt",
			},
			expected: true,
			title:    "one selector multiple attribute all match",
		},
		{
			payload: jsonToGeneric(` {"user_id": "1", "active" : false, "block" : true} `),
			matcher: []string{
				"a/active/bf/user_id/sq:1/block/bf",
			},
			expected: false,
			title:    "one selector multiple attribute one unmatch",
		},
		{
			payload: jsonToGeneric(` {"user_id": "1", "active" : false, "block" : true} `),
			matcher: []string{
				"a/active/bf/user_id/sq:1/block/bf",
				"d/block/bf",
			},
			expected: false,
			title:    "mutli selector multiple attribute & one match in deny",
		},
		{
			payload: jsonToGeneric(` {"address" : {"city" : "seatle"} } `),
			matcher: []string{
				"a/address.city/sq:seatle",
			},
			expected: true,
			title:    "nested attribute matcher",
		},
		{
			payload: jsonToGeneric(` {"address" : {"city" : "seatle"} } `),
			matcher: []string{
				"d/address.city/sq:seatle",
			},
			expected: false,
			title:    "nested attribute matcher",
		},
		{
			payload: jsonToGeneric(` {"address" : {"city" : "newyork"} } `),
			matcher: []string{
				"d/address.city/sq:seatle",
				"a/address.city/sq:newyork",
			},
			expected: true,
			title:    "nested attribute matcher",
		},
		{
			payload: jsonToGeneric(` {"address" : {"city" : "seattle"} } `),
			matcher: []string{
				"d/eyeColor/sq:blue",
				"a/eyeColor/sq:red",
				"d/address.city/sq:seattle",
			},
			expected: false,
			title:    "nested attribute matcher",
		},
	} {
		title := skena.title
		if title == "" {
			title = fmt.Sprintf("Skena %d", i)
		}
		t.Run(title, func(t *testing.T) {
			q := NewPrimitiveQuery(skena.matcher)
			if q.Match(skena.payload) != skena.expected {
				t.Errorf("unmatched")
			}
		})
	}
}

func jsonToGeneric(s string) map[string]interface{} {
	var p map[string]interface{}
	json.Unmarshal([]byte(s), &p)
	return p
}
