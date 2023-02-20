package config

import (
	"io/ioutil"
	"os"
	"testing"
)

func TestGetConfig(t *testing.T) {
	// Create temporary YAML file
	content := []byte(`
server:
  cache: true
  interval_time_type: "seconds"
  interval_amount: 60
`)
	tmpfile, err := ioutil.TempFile("", "config_test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())
	if _, err := tmpfile.Write(content); err != nil {
		t.Fatal(err)
	}
	if err := tmpfile.Close(); err != nil {
		t.Fatal(err)
	}

	// Test GetConfig
	config := GetConfig(tmpfile.Name())
	if !config.Server.Cache {
		t.Errorf("expected cache to be true, got %v", config.Server.Cache)
	}
	if config.Server.IntervalTimeType != "seconds" {
		t.Errorf("expected interval_time_type to be \"seconds\", got %q", config.Server.IntervalTimeType)
	}
	if config.Server.IntervalAmount != 60 {
		t.Errorf("expected interval_amount to be 60, got %d", config.Server.IntervalAmount)
	}
}
