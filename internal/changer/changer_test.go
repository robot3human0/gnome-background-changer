package changer

import (
	"os"
	"path/filepath"
	"sync"
	"testing"
	"time"

	"github.com/robot3human0/gnome-background-changer/internal/config"
)

// Util function: create a temp folder with fake pictures inside
func makeFakeWallpaperDir(t *testing.T, names ...string) string {
	t.Helper()
	dir := t.TempDir()
	for _, name := range names {
		os.WriteFile(filepath.Join(dir, name), []byte{}, os.FileMode(0o644))
	}
	return dir
}

// Start() function runs the rotation and call the setter
func TestChanger_StartCallSetter(t *testing.T) {
	dir := makeFakeWallpaperDir(t, "a.jpg", "b.jpg")

	var mu sync.Mutex
	var called []string

	cfg := config.Config {
		FolderPath: dir,
		Interval: 100 * time.Millisecond,
	}

	ch := NewChanger(&cfg, nil, nil, func(path string) error {
		mu.Lock()
		called = append(called, path)
		mu.Unlock()
		return nil
	})

	ch.Start()
	time.Sleep(350 * time.Millisecond)
	ch.Stop()

	mu.Lock()
	n := len(called)
	mu.Unlock()

	if n < 2 {
		t.Errorf("expected at least 2 setter calls, got %d", n)
	}
}

// Stop() function stops the rotation
func TestChanger_Stop(t *testing.T) {
	dir := makeFakeWallpaperDir(t, "aboba.png", "biobab.webp")

	cfg := config.Config {
		FolderPath: dir,
		Interval: 50 * time.Millisecond,
	}

	ch := NewChanger(&cfg, nil, nil, func(path string) error { return nil })

	ch.Start()
	ch.Stop()
	
	if ch.IsRunning() {
		t.Error("expected changer to be stopped")
	}
}

// Stop() function is idempotent - the second call doesn't panic
func TestChanger_StopIdempotent(t *testing.T) {
	dir := makeFakeWallpaperDir(t, "giga.bmp")
	cfg := config.Config {
		FolderPath: dir,
		Interval: time.Minute,
	}
	ch := NewChanger(&cfg, nil, nil, func(path string) error { return nil })

	ch.Start()
	ch.Stop()
	ch.Stop() // must NOT panic
}

// Next() function immediately change the background
func TestChanger_Next(t *testing.T) {
	dir := makeFakeWallpaperDir(t, "o.jpg", "o.png", "o.bmp")

	var mu sync.Mutex
	var called []string

	cfg := config.Config {
		FolderPath: dir,
		Interval: time.Hour,
	}
	
	ch := NewChanger(&cfg, nil, nil, func(path string) error {
		mu.Lock()
		called = append(called, path)
		mu.Unlock()
		return nil
	})

	ch.Start()
	ch.Next()
	ch.Stop()

	mu.Lock()
	n := len(called)
	mu.Unlock()

	if n < 2 {
		t.Errorf("expected at least 2 calls (Start + Next), got %d", n)
	}
}