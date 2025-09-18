package storage

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"testing"
)

func TestPhotoStorage_Save(t *testing.T) {
	tmpDir := t.TempDir()

	ps := NewPhotoStorage(tmpDir)
	name := "test.jpg"

	content := []byte("TEST FILE")
	srcPath := filepath.Join(tmpDir, "src.tmp")

	if err := os.WriteFile(srcPath, content, 0644); err != nil {
		t.Fatal(err)
	}

	srcFile, err := os.Open(srcPath)
	if err != nil {
		t.Fatal(err)
	}

	if err := ps.Save(name, srcFile); err != nil {
		t.Fatalf("Save returned error: %v", err)
	}

	gotFile, err := ps.GetByName(name)
	if err != nil {
		t.Fatalf("GetByName returned err: %v", err)
	}
	defer gotFile.Close()

	gotBytes, err := io.ReadAll(gotFile)
	if err != nil {
		t.Fatalf("read saved file: %v", err)
	}
	if !bytes.Equal(gotBytes, content) {
		t.Fatalf("file content mismatch: got %q want %q", string(gotBytes), string(content))
	}
}
