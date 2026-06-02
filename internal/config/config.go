package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"
)

type Config struct {
	FolderPath string		 `json:"folder_path"`
	Subfolders bool			 `json:"subfolders"`
	Interval   time.Duration `json:"interval"`
	Random     bool			 `json:"random"`
	Enabled    bool			 `json:"enabled"`
	LastIndex  int			 `json:"last_index"`
}

func defaultConfig() Config {
	home, _ := os.UserHomeDir()
	return Config {
		FolderPath: filepath.Join(home, "Pictures"),
		Subfolders: false,
		Interval: 30 * time.Minute,
		Random: false,
		Enabled: true,
		LastIndex: 0,
	}
}

func ConfigPath() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".config", "background-changer", "config.json")
}

func LoadConfig(configPath string) (Config, error) {
	cfg := defaultConfig()
	data, err := os.ReadFile(configPath)
	if err != nil {
		if os.IsNotExist(err) {
			return cfg, nil
		}
		return cfg, err
	}
	if err = json.Unmarshal(data, &cfg); err != nil {
		return defaultConfig(), err
	}
	return cfg, nil
}

func SaveConfig(cfg Config, configPath string) error {
	if err := os.MkdirAll(filepath.Dir(configPath), os.FileMode(0o755)); err != nil {
		return err
	}

	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(configPath, data, os.FileMode(0o644))
}
