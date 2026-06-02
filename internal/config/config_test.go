package config

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"
)

// Test: func defaultConfig() Config
func TestDefaultConfig(t *testing.T) {
	home, _ := os.UserHomeDir()
	cfg := defaultConfig()

	tests := []struct {
		name string
		got  any
		want any
	}{
		{"FolderPath", cfg.FolderPath, filepath.Join(home, "Pictures")},
		{"Subfolders", cfg.Subfolders, false},
		{"Interval", cfg.Interval, 30 * time.Minute},
		{"Random", cfg.Random, false},
		{"Enabled", cfg.Enabled, true},
		{"LastIndex", cfg.LastIndex, 0},
	}

	for _, tt := range tests {
		if tt.got != tt.want {
			t.Errorf("%s: got %v, want %v", tt.name, tt.got, tt.want)
		}
	}
}

// Test: func configPath() string
func TestConfigPath(t *testing.T) {
	home, _ := os.UserHomeDir()
	hardcodedPath := home + "/.config/background-changer/config.json"
	if path := ConfigPath(); hardcodedPath != path {
		t.Errorf("Paths are not equal: got %v, want %v", path, hardcodedPath)
	}
}

// Test: func LoadConfig() Config
func TestLoadConfig(t *testing.T) {
	const (
		folderPath = "/path/to/pictures"
		subfolders = true
		interval   = time.Duration(15 * time.Minute)
		random     = true
		enabled    = false
		lastIndex  = 69
	)
	tempDirPath := t.TempDir()
	fakeConfigPath := filepath.Join(tempDirPath, "config.json")
	configJson := fmt.Sprintf(
		`{"folder_path":"%s","subfolders":%v,"interval":%d,"random":%v,"enabled":%v,"last_index":%v}`,
		folderPath,
		subfolders,
		interval,
		random,
		enabled,
		lastIndex,
	)

	file, err := os.Create(fakeConfigPath)
	if err != nil {
		t.Error(err)
		return
	}
	file.WriteString(configJson)
	file.Close()

	cfg, err := LoadConfig(fakeConfigPath)
	if err != nil {
		t.Fatal(err)
	}
	tests := []struct {
		name string
		got  any
		want any
	}{
		{"FolderPath", cfg.FolderPath, folderPath},
		{"Subfolders", cfg.Subfolders, subfolders},
		{"Interval", cfg.Interval, interval},
		{"Random", cfg.Random, random},
		{"Enabled", cfg.Enabled, enabled},
		{"LastIndex", cfg.LastIndex, lastIndex},
	}

	for _, tt := range tests {
		if tt.got != tt.want {
			t.Errorf("%s: got %v, want %v", tt.name, tt.got, tt.want)
		}
	}
}

func TestLoadConfig_NoFile(t *testing.T) {
	const randomInvalidPath = "/random/invalid/path"

	defaultCfg := defaultConfig()
	cfg, err := LoadConfig(randomInvalidPath)
	if err != nil {
		t.Fatal("Expected nil error for missing file, fot:", err)
	}

	tests := []struct {
		name string
		got  any
		want any
	}{
		{"FolderPath", cfg.FolderPath, defaultCfg.FolderPath},
		{"Subfolders", cfg.Subfolders, defaultCfg.Subfolders},
		{"Interval", cfg.Interval, defaultCfg.Interval},
		{"Random", cfg.Random, defaultCfg.Random},
		{"Enabled", cfg.Enabled, defaultCfg.Enabled},
		{"LastIndex", cfg.LastIndex, defaultCfg.LastIndex},
	}

	for _, tt := range tests {
		if tt.got != tt.want {
			t.Errorf("%s: got %v, want %v", tt.name, tt.got, tt.want)
		}
	}
}

func TestLoadConfig_InvalidJSON(t *testing.T) {
	f, err := os.CreateTemp(t.TempDir(), "config.json")
	if err != nil {
		t.Fatal(err)
	}
	f.WriteString(`{not valid json}`)
	f.Close()

	_, err = LoadConfig(f.Name())
	if err == nil {
		t.Error("expected error for invalid JSON, got nil")
	}
}

// Test: func SaveConfig(cfg Config) error
func TestSaveConfig(t *testing.T) {
	const (
		folderPath = "/path/to/pictures"
		subfolders = true
		interval   = time.Duration(15 * time.Minute)
		random     = true
		enabled    = false
		lastIndex  = 69
	)

	cfg := Config{
		FolderPath: folderPath,
		Subfolders: subfolders,
		Interval: interval,
		Random: random,
		Enabled: enabled,
		LastIndex: lastIndex,
	}

	tempFolder := t.TempDir()
	configPath := filepath.Join(tempFolder, "temp_config.json")

	err := SaveConfig(cfg, configPath)
	if err != nil {
		t.Error(err)
	}

	loadedSavedCfg, err := LoadConfig(configPath)
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name string
		got  any
		want any
	}{
		{"FolderPath", loadedSavedCfg.FolderPath, cfg.FolderPath},
		{"Subfolders", loadedSavedCfg.Subfolders, cfg.Subfolders},
		{"Interval", loadedSavedCfg.Interval, cfg.Interval},
		{"Random", loadedSavedCfg.Random, cfg.Random},
		{"Enabled", loadedSavedCfg.Enabled, cfg.Enabled},
		{"LastIndex", loadedSavedCfg.LastIndex, cfg.LastIndex},
	}

	for _, tt := range tests {
		if tt.got != tt.want {
			t.Errorf("%s: got %v, want %v", tt.name, tt.got, tt.want)
		}
	}
}
