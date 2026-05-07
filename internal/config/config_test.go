package config_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/yourorg/logpipe/internal/config"
)

func writeTemp(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "logpipe-*.yaml")
	if err != nil {
		t.Fatalf("create temp file: %v", err)
	}
	if _, err := f.WriteString(content); err != nil {
		t.Fatalf("write temp file: %v", err)
	}
	f.Close()
	return f.Name()
}

func TestLoad_ValidConfig(t *testing.T) {
	yaml := `
server:
  addr: "0.0.0.0"
  port: 8080
sources:
  - name: app-logs
    type: docker
    labels:
      env: prod
sinks:
  - name: stdout-sink
    type: stdout
    target: ""
`
	path := writeTemp(t, yaml)
	cfg, err := config.Load(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Server.Port != 8080 {
		t.Errorf("expected port 8080, got %d", cfg.Server.Port)
	}
	if len(cfg.Sources) != 1 || cfg.Sources[0].Name != "app-logs" {
		t.Errorf("unexpected sources: %+v", cfg.Sources)
	}
}

func TestLoad_MissingFile(t *testing.T) {
	_, err := config.Load(filepath.Join(t.TempDir(), "nonexistent.yaml"))
	if err == nil {
		t.Fatal("expected error for missing file, got nil")
	}
}

func TestLoad_InvalidPort(t *testing.T) {
	yaml := `
server:
  addr: "0.0.0.0"
  port: 0
sources: []
sinks: []
`
	path := writeTemp(t, yaml)
	_, err := config.Load(path)
	if err == nil {
		t.Fatal("expected validation error for port 0, got nil")
	}
}

func TestLoad_DuplicateSourceName(t *testing.T) {
	yaml := `
server:
  addr: "0.0.0.0"
  port: 9090
sources:
  - name: dup
    type: file
  - name: dup
    type: stdin
sinks: []
`
	path := writeTemp(t, yaml)
	_, err := config.Load(path)
	if err == nil {
		t.Fatal("expected duplicate name error, got nil")
	}
}

func TestLoad_EmptySources(t *testing.T) {
	yaml := `
server:
  addr: "127.0.0.1"
  port: 4040
sources: []
sinks:
  - name: file-sink
    type: file
    target: "/tmp/out.log"
`
	path := writeTemp(t, yaml)
	cfg, err := config.Load(path)
	if err != nil {
		t.Fatalf("unexpected error for empty sources: %v", err)
	}
	if len(cfg.Sources) != 0 {
		t.Errorf("expected 0 sources, got %d", len(cfg.Sources))
	}
}
