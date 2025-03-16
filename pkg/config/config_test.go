package config

import (
	"os"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	// Create a temporary config file
	configContent := `
rules:
  - source: "quay.io"
    target: "my-ecr-quay"
  - source: "docker.io"
    target: "my-ecr-docker"
`
	tmpFile, err := os.CreateTemp("", "config.yaml")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	if _, err := tmpFile.Write([]byte(configContent)); err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}
	if err := tmpFile.Close(); err != nil {
		t.Fatalf("Failed to close temp file: %v", err)
	}

	// Load the config
	config, err := LoadConfig(tmpFile.Name())
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	// Verify the config content
	if len(config.Rules) != 2 {
		t.Fatalf("Expected 2 rules, got %d", len(config.Rules))
	}
	if config.Rules[0].Source != "quay.io" || config.Rules[0].Target != "my-ecr-quay" {
		t.Errorf("Unexpected rule: %+v", config.Rules[0])
	}
	if config.Rules[1].Source != "docker.io" || config.Rules[1].Target != "my-ecr-docker" {
		t.Errorf("Unexpected rule: %+v", config.Rules[1])
	}
}
