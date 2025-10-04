package storage

import (
	"bytes"
	"image"
	"image/color"
	"image/png"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"testing"
)

func makePNG(t *testing.T) multipart.File{
    t.Helper()
    buf := &bytes.Buffer{}
    img := image.NewRGBA(image.Rect(0, 0, 10, 10))
    for y := 0; y < 10; y++ {
        for x := 0; x < 10; x++ {
            img.Set(x,y, color.RGBA{R:255, G:255})
        }
    }
    if err := png.Encode(buf, img); err != nil {
        t.Fatalf("Error during encoding img: %v", err)
    }

    tmp, err := os.CreateTemp("", "photo-*.png")
    if err != nil {
        t.Fatalf("Error during file creation: %v", err)
    }
    if _, err := tmp.Write(buf.Bytes()); err != nil {
        t.Fatalf("Error during file writing: %v", err)
    }
    if _, err := tmp.Seek(0, io.SeekStart); err != nil {
        t.Fatalf("Error during seek: %v", err)
    }
    return tmp

}

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
