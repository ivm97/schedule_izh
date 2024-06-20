package main

import (
	"encoding/json"
	"testing"
)

func TestReadConf(t *testing.T) {
	data := readConf("settings/config.json")
	if js, err := json.Marshal(data); err != nil {
		t.Errorf("Expected correct json, but got... %v", js)
	}
}
