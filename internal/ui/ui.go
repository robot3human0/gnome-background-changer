package ui

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/ncruces/zenity"
	"github.com/robot3human0/gnome-background-changer/internal/config"
)

var (
	minMinutesInterval = 1
	maxMinutesInterval = 24 * 60
)

// show an error dialog
func ShowError(msg string) {
	zenity.Error(msg, zenity.Title("Background Changer Error"))
}

// opens a simple settings dialog using zenity.
// returns updated config and true if user confirmed, or false if canceled.
func ShowSettingsDialog(current config.Config) (config.Config, bool) {
	cfg := current

	// Step 1: folder picker
	folder, err := zenity.SelectFile(
		zenity.Title("Select Background Folder"),
		zenity.Directory(),
		zenity.Filename(cfg.FolderPath),
	)
	if err != nil {
		// user canceled
		return cfg, false
	}

	cfg.FolderPath = folder

	// Step 2: Interval
	minutes := int(cfg.Interval.Minutes())
	if minutes < minMinutesInterval || minutes > maxMinutesInterval {
		minutes = 30
	}
	intervalStr, err := zenity.Entry(
		fmt.Sprintf("Change interval (current: %d minutes):", minutes),
		zenity.Title("Background Interval"),
		zenity.EntryText(strconv.Itoa(minutes)),
	)
	if err != nil {
		return cfg, false
	}
	intervalStr = strings.TrimSpace(intervalStr)
	if mins, err := strconv.Atoi(intervalStr); err == nil && mins > 0 {
		cfg.Interval = time.Duration(mins) * time.Minute
	}

	// Step 3: order mode
	err = zenity.Question(
		"Rotate backgrounds in random order?\n\nYes = Random  No = Sequential",
		zenity.Title("Background Rotation Order"),
		zenity.OKLabel("Random"),
		zenity.CancelLabel("Sequential"),
	)
	cfg.Random = err == nil // nil means user choose OK (Random)

	return cfg, true
}
