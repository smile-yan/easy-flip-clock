#!/bin/bash
set -e

cd "$(dirname "$0")/.."

OUTPUT_NAME="easy-flip-clock"
APP_NAME="${OUTPUT_NAME}.app"
DMG_NAME="${OUTPUT_NAME}.dmg"
DMG_TMP="output/${OUTPUT_NAME}-tmp.dmg"
DMG_OUT="output/${DMG_NAME}"
VOL_NAME="${OUTPUT_NAME}"
ICON_FILE="frontend/imgs/app.icns"
BUILD_DIR="output"
TMP_DIR="/tmp/build-dmg-$$"

echo "=== Building macOS DMG ==="

# Clean previous builds
rm -rf "${BUILD_DIR}"/macos-* 2>/dev/null || true
rm -f "${DMG_TMP}" "${DMG_OUT}" 2>/dev/null || true
mkdir -p "${BUILD_DIR}"

BUNDLE_ID="com.easyflipclock.${OUTPUT_NAME}"

# Build universal binary (arm64 + amd64)
echo "[1/4] Building wails app for darwin/arm64..."
mkdir -p "${BUILD_DIR}/macos-arm64"
CGO_ENABLED=1 GOOS=darwin GOARCH=arm64 go build -o "${BUILD_DIR}/macos-arm64/${OUTPUT_NAME}" .

echo "[2/4] Building wails app for darwin/amd64..."
mkdir -p "${BUILD_DIR}/macos-amd64"
CGO_ENABLED=1 GOOS=darwin GOARCH=amd64 go build -o "${BUILD_DIR}/macos-amd64/${OUTPUT_NAME}" .

echo "[3/4] Creating universal binary..."
mkdir -p "${BUILD_DIR}/macos-universal"
lipo -create -output "${BUILD_DIR}/macos-universal/${OUTPUT_NAME}" \
    "${BUILD_DIR}/macos-arm64/${OUTPUT_NAME}" \
    "${BUILD_DIR}/macos-amd64/${OUTPUT_NAME}"

# Create .app bundle structure
echo "[4/4] Creating .app bundle and DMG..."

rm -rf "${TMP_DIR}"
mkdir -p "${TMP_DIR}/${APP_NAME}/Contents/MacOS"
mkdir -p "${TMP_DIR}/${APP_NAME}/Contents/Resources"

# Copy binary
cp "${BUILD_DIR}/macos-universal/${OUTPUT_NAME}" "${TMP_DIR}/${APP_NAME}/Contents/MacOS/${OUTPUT_NAME}"
chmod +x "${TMP_DIR}/${APP_NAME}/Contents/MacOS/${OUTPUT_NAME}"

# Copy icon if exists
if [ -f "${ICON_FILE}" ]; then
    cp "${ICON_FILE}" "${TMP_DIR}/${APP_NAME}/Contents/Resources/app.icns"
fi

# Create Info.plist
cat > "${TMP_DIR}/${APP_NAME}/Contents/Info.plist" << EOF
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
    <key>CFBundleExecutable</key>
    <string>${OUTPUT_NAME}</string>
    <key>CFBundleIdentifier</key>
    <string>${BUNDLE_ID}</string>
    <key>CFBundleName</key>
    <string>${OUTPUT_NAME}</string>
    <key>CFBundleDisplayName</key>
    <string>${OUTPUT_NAME}</string>
    <key>CFBundleIconFile</key>
    <string>app.icns</string>
    <key>CFBundlePackageType</key>
    <string>APPL</string>
    <key>CFBundleShortVersionString</key>
    <string>1.0.0</string>
    <key>CFBundleVersion</key>
    <string>1</string>
    <key>LSMinimumSystemVersion</key>
    <string>10.15</string>
    <key>NSHumanReadableCopyright</key>
    <string>Copyright © 2024. All rights reserved.</string>
    <key>NSPrincipalClass</key>
    <string>NSApplication</string>
    <key>LSUIElement</key>
    <true/>
</dict>
</plist>
EOF

# Copy frontend assets
cp -r frontend "${TMP_DIR}/${APP_NAME}/Contents/Resources/frontend"

# Create .app bundle
rm -rf "output/${APP_NAME}"
cp -r "${TMP_DIR}/${APP_NAME}" "output/"

# Create uncompressed read-write DMG
MOUNT_DIR="/Volumes/${VOL_NAME}"
hdiutil create \
    -volname "${VOL_NAME}" \
    -srcfolder "output/${APP_NAME}" \
    -ov -format UDRW \
    "${DMG_TMP}"

# Mount the DMG
ATTACH_OUT="$(hdiutil attach "${DMG_TMP}" -mountpoint "${MOUNT_DIR}" -nobrowse 2>&1)"
DISK_DEV="$(printf '%s\n' "${ATTACH_OUT}" | awk '/GUID_partition_scheme/ {print $1; exit}')"

# Add Applications symlink
ln -sf /Applications "${MOUNT_DIR}/Applications"

# Add volume icon
if [ -f "${ICON_FILE}" ]; then
    cp "${ICON_FILE}" "${MOUNT_DIR}/.VolumeIcon.icns"
    SetFile -a C "${MOUNT_DIR}"
fi

# Force detach
hdiutil detach "${DISK_DEV:-$MOUNT_DIR}" -force

# Convert to compressed UDZO
hdiutil convert "${DMG_TMP}" -format UDZO -o "${DMG_OUT}"
rm -f "${DMG_TMP}"

# Embed icon resource into the DMG file itself using Rez
ICON_TMP="/tmp/${OUTPUT_NAME}-dmg-icon.icns"
ICON_RSRC="/tmp/${OUTPUT_NAME}-dmg-icon.rsrc"
cp "${ICON_FILE}" "${ICON_TMP}"
sips -i "${ICON_TMP}" >/dev/null
DeRez -only icns "${ICON_TMP}" > "${ICON_RSRC}"
Rez -append "${ICON_RSRC}" -o "${DMG_OUT}"
rm -f "${ICON_TMP}" "${ICON_RSRC}"
SetFile -a C "${DMG_OUT}"

echo ""
echo "=== Done! ==="
echo "DMG created at: ${DMG_OUT}"
ls -lh "${DMG_OUT}"

# Cleanup
rm -rf "${TMP_DIR}"