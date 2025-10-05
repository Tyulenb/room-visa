package storage

import (
	"bytes"
	"image"
	"image/color"
	"image/jpeg"
	"image/png"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/google/uuid"
)

func CreateTemp(t *testing.T, buf *bytes.Buffer) multipart.File {
    t.Helper()
    tmp, err := os.CreateTemp("", "photo-*.jpeg")
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

func makePNG(t *testing.T) multipart.File {
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

    return CreateTemp(t, buf)
}

func makeJPEG(t *testing.T) multipart.File {
    t.Helper()
    buf := &bytes.Buffer{}
    img := image.NewRGBA(image.Rect(0, 0, 10, 10))
    for y := 0; y < 10; y++ {
        for x := 0; x < 10; x++ {
            img.Set(x,y, color.RGBA{R:255, B:255})
        }
    }
    if err := jpeg.Encode(buf, img, nil); err != nil {
        t.Fatalf("Error during encoding img: %v", err)
    }

    return CreateTemp(t, buf)
}

func makeBadFile(t *testing.T) multipart.File {
    t.Helper()
    buf := bytes.NewBufferString("NOT AN IMAGE")
    return CreateTemp(t, buf)
}

func TestSaveGet(t *testing.T) {
    tmpDir := t.TempDir()
	ps := &PhotoStorage{path: tmpDir}
    name := uuid.NewString()

    tests := []struct{
        name string
        img multipart.File
        format string
        expErr bool
        err error
    }{
        {
            name: "PNG",
            img: makePNG(t),
            format: "png",
            expErr: false,
            err: nil,
        },
        {
            name: "JPEG",
            img: makeJPEG(t),
            format: "jpeg",
            expErr: false,
            err: nil,
        },
        {
            name: "BAD",
            img: makeBadFile(t),
            format: "jpeg",
            expErr: true,
            err: image.ErrFormat,
        },
    }
    
    for _, tt := range tests {
        err := ps.Save(name, tt.img)
        if err != nil {
            if tt.expErr && tt.err == err {
                return
            }else {
                t.Fatalf("Test: %s\nSave(name, img) failed: %v", tt.name, err)
            }
        }
        
        fullPath := filepath.Join(tmpDir, name)
        photoStats, err := os.Stat(fullPath)
        if err != nil {
            t.Fatalf("Test: %s\nStat failed: %v", tt.name, err)
        }

        if photoStats.Size() == 0 {
            t.Fatalf("Test: %s\nFile size is 0, empty image", tt.name)
        }

        file, err := ps.GetByName(name)
        if err != nil {
            t.Fatalf("Test: %s\nFailed GetByName(name): %v", err, tt.name)
        }

        _, format, err := image.Decode(file)
        if err != nil {
            t.Fatalf("Test: %s\nFailed during file decoding: %v", err, tt.name)
        }

        frmt := strings.ToLower(format)
        if (frmt != tt.format) {
            t.Fatalf("Test: %s\nExpected %s, but got: %s", tt.name, tt.format, frmt)
        }
    }
}
