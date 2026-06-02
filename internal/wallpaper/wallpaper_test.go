package wallpaper

import (
	"os"
	"path/filepath"
	"testing"
)

/* 
	The SetWallpaper() function just a shall out, so it don't need any test.
*/

// Test: func ListWallpapers(folder string) ([]string, error)
// returns correct files from folder with pictures
func TestListWallpapers_ReturnImages(t *testing.T) {
	dir := t.TempDir()
	for _, name := range []string{"a.jpg", "b.jpg", "c.jpg"} {
		os.WriteFile(filepath.Join(dir, name), []byte{}, os.FileMode(0o644))
	}
	
	files, err := ListWallpapers(dir)
	if err != nil {
		t.Fatal(err)
	}
	if len(files) != 3 {
		t.Errorf("got %d files, want 3", len(files))
	}
}

// filter NO image files
func TestListWallpapers_FiltersNoImages(t *testing.T) {
	dir := t.TempDir()
	os.WriteFile(filepath.Join(dir, "image.jpg"), []byte{}, os.FileMode(0o644))
	os.WriteFile(filepath.Join(dir, "document.pdf"), []byte{}, os.FileMode(0o644))
	os.WriteFile(filepath.Join(dir, "readme.txt"), []byte{}, os.FileMode(0o644))

	files, err := ListWallpapers(dir)
	if err != nil {
		t.Fatal(err)
	}

	if len(files) != 1 {
		t.Errorf("got %d files, want 1", len(files))
	}
}

// Empty folder must return the error
func TestListWallpaper_EmptyDir(t *testing.T) {
	dir := t.TempDir()

	_, err := ListWallpapers(dir)
	if err == nil {
		t.Error("expected error for empty dir. got nil")
	}
}
