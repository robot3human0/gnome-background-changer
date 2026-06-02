# gnome-background-changer

Automatic wallpaper rotation daemon for GNOME with a system tray icon.

## Features
* Rotates wallpapers from a chosen folder on a configurable interval
* Sequential or random rotation
* System tray menu: pause, next, settings, quit
* Supports light and dark GNOME themes simultaneously

## Requirements
* GNOME desktop with `gsettings`
* Go 1.21+ (only for building from source)
* pkg-config package
* Appindicator library (ayatana-appindicator3)


## Installation

### Quick install (recommended)
```bash
git clone https://github.com/robot3human0/gnome-background-changer
cd gnome-background-changer
bash install.sh
```

### Options
| Flag          | Description                 |
|---------------|-----------------------------|
| `--autostart` | Add to GNOME autostart      |
| `--no-tests`  | Skip tests before building  |
| `--uninstall` | Remove the app              |

## Usage

After installation run:
```bash
gnome-background-changer
```

The app appears in the system tray. Click to access the menu.
On first run, open **Settings** to choose a wallpaper source folder and interval.

## License
GPL v3 — see [LICENSE](LICENSE).
