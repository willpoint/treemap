package main

import (
	"encoding/json"
	"testing"
)

func TestTNode_Weight(t *testing.T) {

	d := []byte(`
    {
       "name": "A",
       "children": [
           {"name": "B", "size": 2813},
           {"name": "C", "size": 813},
           {"name": "D", "children": [
               {"name": "E", "size": 382},
               {"name": "F", "size": 1032}
           ]}
       ]
    }
    `)
	tm := new(TNode)
	err := json.Unmarshal(d, tm)
	if err != nil {
		t.Errorf("failed unmarshalling: %v", err)
	}
	var expected float64 = 2813 + 813 + 382 + 1032
	if got := tm.Weight(); got != expected {
		t.Errorf("expected %f got %f", expected, got)
	}
}
