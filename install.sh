#!/bin/bash

# —————————————————————————————————————————————
#   gnome-background-changer — install script
# —————————————————————————————————————————————

APP_NAME="gnome-background-changer"
BIN_NAME="gnome-background-changer"

INSTALL_DIR="${HOME}/.local/bin"

APP_DIR="${HOME}/.local/share/applications"
APP_DESKTOP_FILE="${APP_DIR}/${APP_NAME}.desktop"

ICON_SRC="icons/tray_icon.png"
ICON_DIR="${HOME}/.local/share/icons/hicolor/256x256/apps"
ICON_DST="${ICON_DIR}/${APP_NAME}.png"

AUTOSTART_DIR="${HOME}/.config/autostart"
AUTOSTART_FILE="${AUTOSTART_DIR}/${APP_NAME}.desktop"

# ———— FLAGS ——————————————————————————————————

ENABLE_AUTOSTART=false
UNINSTALL=false
RUN_TESTS=true
 
usage() {
    cat <<EOF
Usage: $0 [OPTIONS]
 
Options:
  --autostart     Add the application to autostart
  --no-tests      Skip tests
  --uninstall     Uninstall the app
  -h, --help      Show this menu
EOF
    exit 0
}
 
for arg in "$@"; do
    case "$arg" in
        --autostart)   ENABLE_AUTOSTART=true ;;
        --no-tests)    RUN_TESTS=false ;;
        --uninstall)   UNINSTALL=true ;;
        -h|--help)     usage ;;
        *) echo "Unknown flag: $arg" >&2; usage ;;
    esac
done
 
# ———— COLORS —————————————————————————————————
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
CYAN='\033[0;36m'
NC='\033[0m'
 
info()    { echo -e "${CYAN}[INFO]${NC}  $*"; }
success() { echo -e "${GREEN}[OK]${NC}    $*"; }
warn()    { echo -e "${YELLOW}[WARN]${NC}  $*"; }
error()   { echo -e "${RED}[ERROR]${NC} $*" >&2; exit 1; }
 
# ═════════════════════════════════════════════
#  UNINSTALL
# ═════════════════════════════════════════════
uninstall() {
    info "Uninstalling the ${APP_NAME} ..."
 
    local removed=0
 
    # Remove binary
    if [[ -f "${INSTALL_DIR}/${BIN_NAME}" ]]; then
        rm -f "${INSTALL_DIR}/${BIN_NAME}"
        success "binary file removed: ${INSTALL_DIR}/${BIN_NAME}"
        removed=1
    fi
 
    # Remove icons
    if [[ -f "${HOME}/.local/share/${APP_NAME}/icons/tray_icon.png" ]]; then
        rm -f "${HOME}/.local/share/${APP_NAME}/icons/tray_icon.png"
        success "Tray icon removed."
        removed=1
    fi

    # Remove system icon
    if [[ -f "${ICON_DST}" ]]; then
        rm -f "${ICON_DST}"
        success "Icon removed: ${ICON_DST}"
        removed=1
    fi
 
    # Remove desktop entry
    if [[ -f "${APP_DESKTOP_FILE}" ]]; then
        rm -f "${APP_DESKTOP_FILE}"
        success "Desktop entry removed: ${APP_DESKTOP_FILE}"
        removed=1
    fi

    # Remove autostart entry
    if [[ -f "${AUTOSTART_FILE}" ]]; then
        rm -f "${AUTOSTART_FILE}"
        success "Autostart entry removed: ${AUTOSTART_FILE}"
        removed=1
    fi
 
    if [[ $removed -eq 0 ]]; then
        warn "Nothing found - the application was not installed."
    else
        rmdir --ignore-fail-on-non-empty \
            "${HOME}/.local/share/${APP_NAME}/icons" \
            "${HOME}/.local/share/${APP_NAME}" 2>/dev/null
        if command -v gtk-update-icon-cache &>/dev/null; then
            gtk-update-icon-cache -f -t "${HOME}/.local/share/icons/hicolor" 2>/dev/null || true
        fi
        if command -v update-desktop-database &>/dev/null; then
            update-desktop-database "${HOME}/.local/share/applications" 2>/dev/null || true
        fi
        success "The ${APP_NAME} successfully uninstalled."
    fi
    exit 0
}
 
# ═════════════════════════════════════════════
#  INSTALL
# ═════════════════════════════════════════════
check_deps() {
    info "Check dependencies ..."
    command -v go >/dev/null || \
        error "Go not found"
 
    local go_version
    go_version=$(go version | grep -oP 'go\K[0-9]+\.[0-9]+')
    local major minor
    major=$(echo "$go_version" | cut -d. -f1)
    minor=$(echo "$go_version" | cut -d. -f2)
    if [[ $major -lt 1 ]] || { [[ $major -eq 1 ]] && [[ $minor -lt 21 ]]; }; then
        error "Go 1.21+ is required. Installed version: ${go_version}"
    fi

    command -v pkg-config >/dev/null || \
        error "pkg-config not found"

    pkg-config --exists appindicator3-0.1 || \
        pkg-config --exists ayatana-appindicator3-0.1 || \
        error "Required AppIndicator library not found. On Debian/Ubuntu: sudo apt install libayatana-appindicator3-dev"

    success "Dependencies OK"
}
 
run_tests() {
    info "Running tests..."
    if go test ./internal/...; then
        success "All tests passed."
    else
        error "Tests failed. Installation aborted"
    fi
}
 
build() {
    info "Building ${APP_NAME}..."
    go build -v -o "${BIN_NAME}" "./cmd/${APP_NAME}" || \
        error "Build failed"
    success "Binary built: ./${BIN_NAME}"
}
 
install_bin() {
    [[ -f "${BIN_NAME}" ]] || error "Binary not found: ${BIN_NAME}"

    info "Installing binary to ${INSTALL_DIR}..."
    
    mkdir -p "${INSTALL_DIR}"
    mv "${BIN_NAME}" "${INSTALL_DIR}/${BIN_NAME}"
    chmod +x "${INSTALL_DIR}/${BIN_NAME}"

    success "Installed: ${INSTALL_DIR}/${BIN_NAME}"
 
    # Warn if ~/.local/bin is not in PATH
    if [[ ":${PATH}:" != *":${INSTALL_DIR}:"* ]]; then
        warn "${INSTALL_DIR} is not in PATH."
        warn "Add to ~/.bashrc or ~/.zshrc:"
        warn "  export PATH=\"\${HOME}/.local/bin:\${PATH}\""
    fi
}
 
install_icon() {
    if [[ ! -f "${ICON_SRC}" ]]; then
        warn "Icon not found (${ICON_SRC}), skipping."
        return
    fi
    info "Installing icon..."
    local xdg_icon_dir="${HOME}/.local/share/${APP_NAME}/icons"
    mkdir -p "${xdg_icon_dir}"
    cp "${ICON_SRC}" "${xdg_icon_dir}/tray_icon.png"
    # for system theme
    mkdir -p "${ICON_DIR}"
    cp "${ICON_SRC}" "${ICON_DST}" 
    if command -v gtk-update-icon-cache &>/dev/null; then
        gtk-update-icon-cache -f -t "${HOME}/.local/share/icons/hicolor" 2>/dev/null || true
    fi
    success "Icon installed: ${xdg_icon_dir}"
}

install_desktop_entry() {
    info "Creating desktop entry..."
    mkdir -p "${APP_DIR}"
    cat > "${APP_DESKTOP_FILE}" <<EOF
[Desktop Entry]
Type=Application
Name=GNOME Background Changer
Categories=Utility;DesktopSettings;
Comment=Automatic wallpaper rotation for GNOME
Exec=${INSTALL_DIR}/${BIN_NAME}
Icon=${APP_NAME}
Terminal=false
Hidden=false
StartupNotify=false
Version=1.0
EOF
    success "Desktop entry configured"
}

install_autostart() {
    info "Creating autostart entry..."
    mkdir -p "${AUTOSTART_DIR}"
    cat > "${AUTOSTART_FILE}" <<EOF
[Desktop Entry]
Type=Application
Name=GNOME Background Changer
Categories=Utility;DesktopSettings;
Comment=Automatic wallpaper rotation for GNOME
Exec=${INSTALL_DIR}/${BIN_NAME}
Icon=${APP_NAME}
Terminal=false
Hidden=false
StartupNotify=false
X-GNOME-Autostart-enabled=true
EOF
    success "Autostart configured"
}
 
# ═════════════════════════════════════════════
#  MAIN
# ═════════════════════════════════════════════
main() {
    echo -e "${CYAN}"
    echo "╔══════════════════════════════════════════╗"
    echo "║     gnome-background-changer installer   ║"
    echo "╚══════════════════════════════════════════╝"
    echo -e "${NC}"
 
    [[ "$UNINSTALL" == true ]] && uninstall
 
    # Ensure we run from project root
    if [[ ! -f "go.mod" ]]; then
        error "go.mod not found. Run this script from the project root."
    fi
 
    check_deps
    [[ "$RUN_TESTS" == true ]] && run_tests
    build
    install_bin
    install_icon
    install_desktop_entry
    [[ "$ENABLE_AUTOSTART" == true ]] && install_autostart

    if command -v update-desktop-database &>/dev/null; then
        update-desktop-database "${HOME}/.local/share/applications" 2>/dev/null || true
    fi
 
    echo ""
    echo -e "${GREEN}══════════════════════════════════════════${NC}"
    success "Installation complete!"
    if [[ "$ENABLE_AUTOSTART" == true ]]; then
        info  "The app will start automatically on next GNOME login."
    else
        info  "Run manually: ${BIN_NAME}"
        info  "To enable autostart: $0 --autostart"
    fi
    echo -e "${GREEN}══════════════════════════════════════════${NC}"
}
 
main
