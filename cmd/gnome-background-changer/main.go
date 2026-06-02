package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/getlantern/systray"
	"github.com/robot3human0/gnome-background-changer/internal/changer"
	"github.com/robot3human0/gnome-background-changer/internal/config"
	"github.com/robot3human0/gnome-background-changer/internal/ui"
	"github.com/robot3human0/gnome-background-changer/internal/wallpaper"
)

var (
	cfg   config.Config
	chngr *changer.Changer
)

func main() {
	var err error
	cfg, err = config.LoadConfig(config.ConfigPath())
	if err != nil {
		log.Fatal(err)
	}
	systray.Run(onReady, onExit)
}

func onReady() {
	systray.SetIcon(loadIcon())
	systray.SetTitle("wallpaper changer")
	systray.SetTooltip("wallpaper changer")

	// Menu items
	mStatus := systray.AddMenuItem("⏸  Pause", "Pause/resume wallpaper rotation")
	mNext := systray.AddMenuItem("⏭  Next Wallpaper", "Switch to next wallpaper now")
	systray.AddSeparator()
	mSettings := systray.AddMenuItem("⚙  Settings", "Open settings")
	systray.AddSeparator()
	mQuit := systray.AddMenuItem("✕  Quit", "Quit wallpaper changer")

	// Start changer
	chngr = changer.NewChanger(
		&cfg,
		func(err error) {
			ui.ShowError(fmt.Sprintf("Changer error: %v", err))
		},
		func(path string) {
			name := filepath.Base(path)
			systray.SetTooltip(fmt.Sprintf("Now: %s", name))
		},
		wallpaper.SetWallpaper,
	)

	if cfg.Enabled {
		chngr.Start()
		mStatus.SetTitle("⏸  Pause")
	} else {
		mStatus.SetTitle("▶  Resume")
	}

	// Event loop
	go func() {
		for {
			select {
			case <-mStatus.ClickedCh:
				if chngr.IsRunning() {
					chngr.Stop()
					cfg.Enabled = false
					mStatus.SetTitle("▶  Resume")
				} else {
					cfg.Enabled = true
					chngr.Start()
					mStatus.SetTitle("⏸  Pause")
				}
				config.SaveConfig(cfg, config.ConfigPath())

			case <-mNext.ClickedCh:
				if chngr.IsRunning() {
					chngr.Next()
				}

			case <-mSettings.ClickedCh:
				newCfg, ok := ui.ShowSettingsDialog(cfg)
				if ok {
					wasRunning := chngr.IsRunning()
					chngr.Stop()
					cfg = newCfg
					chngr.UpdateConfig(cfg)
					config.SaveConfig(cfg, config.ConfigPath())
					if wasRunning {
						chngr.Start()
						mStatus.SetTitle("⏸  Pause")
					}
				}

			case <-mQuit.ClickedCh:
				chngr.Stop()
				config.SaveConfig(cfg, config.ConfigPath())
				systray.Quit()
				return
			}
		}
	}()
}

func onExit() {
	os.Exit(0)
}

// loadIcon loads the tray icon from the XDG data directory
// (~/.local/share/gnome-background-changer/icons/tray_icon.png),
// falling back to a minimal embedded PNG if not found.
func loadIcon() []byte {
	dataDir := os.Getenv("XDG_DATA_HOME")
	if dataDir == "" {
		home, _ := os.UserHomeDir()
		dataDir = filepath.Join(home, ".local", "share")
	}

	iconPath := filepath.Join(dataDir, "gnome-background-changer", "icons", "tray_icon.png")
	if data, err := os.ReadFile(iconPath); err == nil {
		return data;
	}
	//Fallback: minimal half transparent yellow png
	return []byte{
		0x89, 0x50, 0x4e, 0x47, 0xd,  0xa,  0x1a, 0xa,  0x0,  0x0, 
		0x0,  0xd,  0x49, 0x48, 0x44, 0x52, 0x0,  0x0,  0x0,  0x10, 
		0x0,  0x0,  0x0,  0x10, 0x8,  0x6,  0x0,  0x0,  0x0,  0x1f, 
		0xf3, 0xff, 0x61, 0x0,  0x0,  0x0,  0x21, 0x49, 0x44, 0x41, 
		0x54, 0x78, 0x9c, 0x62, 0xf9, 0x3b, 0x9d, 0xa1, 0x81, 0x81, 
		0x2,  0xc0, 0x44, 0x89, 0xe6, 0x51, 0x3,  0x46, 0xd,  0x18, 
		0x35, 0x60, 0x30, 0x19, 0x0,  0x8,  0x0,  0x0,  0xff, 0xff, 
		0xbe, 0xcd, 0x2,  0x37, 0x2,  0x5c, 0xc9, 0x95, 0x0,  0x0, 
		0x0,  0x0,  0x49, 0x45, 0x4e, 0x44, 0xae, 0x42, 0x60, 0x82,
	}
}
