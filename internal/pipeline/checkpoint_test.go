package pipeline

import (
	"os"
	"path/filepath"
	"testing"
)

func tempCheckpointPath(t *testing.T) string {
	t.Helper()
	return filepath.Join(t.TempDir(), "checkpoint.json")
}

func TestCheckpoint_GetReturnsZeroForUnknownSource(t *testing.T) {
	cp, err := NewCheckpoint(tempCheckpointPath(t))
	if err != nil {
		t.Fatalf("NewCheckpoint: %v", err)
	}
	if got := cp.Get("unknown"); got != 0 {
		t.Errorf("expected 0, got %d", got)
	}
}

func TestCheckpoint_SetAndGet(t *testing.T) {
	cp, err := NewCheckpoint(tempCheckpointPath(t))
	if err != nil {
		t.Fatalf("NewCheckpoint: %v", err)
	}
	if err := cp.Set("src-a", 42); err != nil {
		t.Fatalf("Set: %v", err)
	}
	if got := cp.Get("src-a"); got != 42 {
		t.Errorf("expected 42, got %d", got)
	}
}

func TestCheckpoint_PersistsAcrossReload(t *testing.T) {
	path := tempCheckpointPath(t)
	cp, _ := NewCheckpoint(path)
	_ = cp.Set("src-b", 99)

	cp2, err := NewCheckpoint(path)
	if err != nil {
		t.Fatalf("reload: %v", err)
	}
	if got := cp2.Get("src-b"); got != 99 {
		t.Errorf("expected 99 after reload, got %d", got)
	}
}

func TestCheckpoint_Delete(t *testing.T) {
	path := tempCheckpointPath(t)
	cp, _ := NewCheckpoint(path)
	_ = cp.Set("src-c", 7)
	if err := cp.Delete("src-c"); err != nil {
		t.Fatalf("Delete: %v", err)
	}
	if got := cp.Get("src-c"); got != 0 {
		t.Errorf("expected 0 after delete, got %d", got)
	}
}

func TestCheckpoint_DeletePersists(t *testing.T) {
	path := tempCheckpointPath(t)
	cp, _ := NewCheckpoint(path)
	_ = cp.Set("src-d", 55)
	_ = cp.Delete("src-d")

	cp2, _ := NewCheckpoint(path)
	if got := cp2.Get("src-d"); got != 0 {
		t.Errorf("expected 0 in reloaded checkpoint, got %d", got)
	}
}

func TestCheckpoint_MissingFileIsNotError(t *testing.T) {
	path := filepath.Join(t.TempDir(), "does_not_exist.json")
	_, err := NewCheckpoint(path)
	if err != nil {
		t.Errorf("expected no error for missing file, got %v", err)
	}
}

func TestCheckpoint_InvalidJSONReturnsError(t *testing.T) {
	path := tempCheckpointPath(t)
	_ = os.WriteFile(path, []byte("not-json"), 0o644)
	_, err := NewCheckpoint(path)
	if err == nil {
		t.Error("expected error for invalid JSON, got nil")
	}
}
