package wallpaper

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

var supportExts = map[string]bool {
	".jpg":  true,
	".jpeg": true,
	".png":  true,
	".svg":  true,
	".webp": true,
	".bmp":  true,
	".tiff": true,
	".tif":  true,
}

// SetWallpaper sets the desktop background using gsettings.
// It sets both the light and dark URI so it works regardless of the theme.
func SetWallpaper(path string) error {
	uri := fmt.Sprintf("file://%s", path)

	cmds := [][]string {
		{"gsettings", "set", "org.gnome.desktop.background", "picture-uri", uri},
		{"gsettings", "set", "org.gnome.desktop.background", "picture-uri-dark", uri},
		// also set picture-options to "zoom" for best fit
		{"gsettings", "set", "org.gnome.desktop.background", "picture-options", "zoom"},
	}

	for _, args := range cmds {
		cmd := exec.Command(args[0],args[1:]...)
		// gsettings needs DBUS_SESSION_BUS_ADDRESS from the user session
		cmd.Env = os.Environ()
		if out, err := cmd.CombinedOutput(); err != nil {
			return fmt.Errorf("gsettings error: %w - %s", err, string(out))
		}
	}

	return nil
}

// ListWallpapers returns all supported image files in the given folder.
func ListWallpapers(folder string) ([]string, error) {
	entries, err := os.ReadDir(folder)
	if err != nil {
		return nil, err
	}

	var files []string
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		ext := strings.ToLower(filepath.Ext(e.Name()))
		if supportExts[ext] {
			files = append(files, filepath.Join(folder, e.Name()))
		}
	}

	if len(files) == 0 {
		return nil, fmt.Errorf("no supported images in %s", folder)
	}

	return files, nil
}
