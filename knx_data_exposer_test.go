package main

import (
	"testing"

	"go.uber.org/zap/zaptest"
)

func TestDatapointStructureGeneration(t *testing.T) {
	logger := zaptest.NewLogger(t)
	m := map[string][]string{
		"DPT_1000": {"0/0/0", "0/0/1"},
		"DPT_2000": {"1/0/0"},
	}

	r := create_datapoint_structure(logger, m)
	if val, ok := r["0/0/0"]; ok {
		if val != "1.000" {
			t.Fatalf("Expected \"1.000\" as value for key \"0/0/0\" but got: %s", val)
		}
	} else {
		t.Fatalf("group address 0/0/0 was not found as a key in the datapoint structure")
	}
	if val, ok := r["0/0/1"]; ok {
		if val != "1.000" {
			t.Fatalf("Expected \"1.000\" as value for key \"0/0/1\" but got: %s", val)
		}
	} else {
		t.Fatalf("group address 0/0/1 was not found as a key in the datapoint structure")
	}
	if val, ok := r["1/0/0"]; ok {
		if val != "2.000" {
			t.Fatalf("Expected \"2.000\" as value for key \"1/0/0\" but got: %s", val)
		}
	} else {
		t.Fatalf("group address 1/0/0 was not found as a key in the datapoint structure")
	}
}
