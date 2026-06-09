<p align="center">
  <img src="./frontend/imgs/app-icon-1024.png" alt="Easy Flip Clock" width="200">
</p>

<h1 align="center">Easy Flip Clock</h1>

<p align="center">
  <img src="https://img.shields.io/github/v/release/smile-yan/easy-flip-clock?style=flat-square&color=32CD32" alt="Version">
  &nbsp;
  <img src="https://img.shields.io/badge/license-BSD--4--Clause-32CD32?style=flat-square" alt="License">
  &nbsp;
  <img src="https://img.shields.io/badge/platform-macOS%20%7C%20Windows%20%7C%20Linux-6C757D?style=flat-square&logo=apple&logoColor=white" alt="Platform">
  &nbsp;
  <img src="https://img.shields.io/badge/built%20with-Wails%203-32CD32?style=flat-square&logo=go&logoColor=white" alt="Built with Wails 3">
</p>

---

<p align="center">
  <strong>A minimalist, open-source desktop flip clock</strong><br>
  <em>24h Display · Multiple Themes · Custom Motto · Native Fullscreen · Cross-Platform · Zero Cloud</em>
</p>

<p align="center">
  <a href="https://github.com/smile-yan/easy-flip-clock/releases">
    <img src="https://img.shields.io/badge/⬇️%20Download-Latest%20Release-32CD32?style=for-the-badge&logo=github&logoColor=white" alt="Download Latest Release" height="50">
  </a>
</p>

---



<img width="1813" height="619" alt="image" src="https://github.com/user-attachments/assets/afb7e516-95a2-4c62-b672-2d42763f12f9" />


## 📋 Quick Navigation

<p align="center">

[Features](#-features) ·
[Screenshots](#-screenshots) ·
[Quick Start](#-quick-start) ·
[Tech Stack](#-tech-stack) ·
[Contributing](#-contributing)

</p>

---

## ✨ Features

**Easy Flip Clock** is not just a clock — it is a calm, distraction-free timepiece for your desktop. Built with Wails 3 and native web technologies, it blends a tactile flip animation with a clean, modern UI.

|                                   | Common Online Clocks                  | **Easy Flip Clock**                                                          |
| :-------------------------------- | :------------------------------------ | :---------------------------------------------------------------------------- |
| Works fully offline               | No — needs network                    | **Yes — runs as a local native app**                                          |
| Native macOS fullscreen (Space)   | No                                    | **Yes — `NSWindowCollectionBehaviorFullScreenPrimary` & `AXFullScreen=true`** |
| 24h / 12h toggle                  | Limited                               | **Yes — switchable with one click**                                           |
| Multiple themes                   | Few                                   | **Dark / Light / Sepia / Blue / Forest / Sunset / Midnight / Ocean**         |
| Custom motto / daily quote        | No                                    | **Yes — edit in `~/.easy-flip-clock/config.json`**                           |
| Lunar calendar (农历)             | No                                    | **Yes — displayed next to the date bar**                                     |
| In-app auto-update check          | No                                    | **Yes — check GitHub Releases from settings**                                |
| Open source & free                | Often paid / freemium                 | **BSD 4-Clause**                                                              |

---

## 🕰️ Screenshots

<p align="center">
  <img src="./frontend/imgs/app-icon-1024.png" alt="Easy Flip Clock — App Icon" width="240">
</p>

### Three display styles

Choose what fits your desk best — a minimal `HH:MM` with `AM/PM` indicator, or a full `HH:MM:SS` countdown.

<table>
  <tr>
    <td align="center" width="50%">
      <strong>Minimal · without seconds</strong><br>
      <em>A calm reading with AM/PM badge.</em>
    </td>
    <td align="center" width="50%">
      <strong>Full · with seconds</strong><br>
      <em>The classic split-flap look, ticking every second.</em>
    </td>
  </tr>
</table>

### A theme for every mood

Pick from a hand-tuned palette — each theme is a complete re-skin of the card, divider, and background, not just a single accent color change.

> **Dark · Light · Sepia · Blue · Forest · Sunset · Midnight · Ocean**

### Custom motto

Below the clock, a single line displays a daily motto (default: `君子三思而后行`). Edit `~/.easy-flip-clock/config.json` to make it your own.

---

## 🎯 Highlights

### Native macOS fullscreen

This is the core engineering investment. Most "fullscreen" implementations in cross-platform frameworks just enlarge the window. Easy Flip Clock uses a real macOS fullscreen Space — so the green traffic-light button and `Ctrl + Cmd + F` actually work the way you'd expect.

> Verified live via AppleScript:
> ```bash
> osascript -e 'tell application "System Events" to tell process "easy-flip-clock" to tell window 1 to return (value of attribute "AXFullScreen")'
> # → true
> ```

### Multiple themes, multiple personalities

Every theme is hand-tuned — backgrounds, card gradients, divider shadows, and the flip animation easing all change together. No theme is just a recolor of dark.

### Local-first, zero telemetry

The app runs entirely from local assets. There is no analytics, no account, no network call except the explicit "check for update" action in settings.

### Tunable config

`~/.easy-flip-clock/config.json` controls motto, window size, position, and Dock behavior — fully scriptable, no GUI required.

---

## 🚀 Quick Start

### System Requirements

- **macOS**: 11.0 (Big Sur) or higher
- **Windows**: Windows 10 or higher
- **Linux**: Ubuntu 20.04+ / equivalent
- **Memory**: 100MB or less (it's a clock)
- **Storage**: ~20MB

### Install

<p>
  <a href="https://github.com/smile-yan/easy-flip-clock/releases">
    <img src="https://img.shields.io/badge/Download-Latest%20Release-32CD32?style=for-the-badge&logo=github&logoColor=white" alt="Download Latest Release" height="50">
  </a>
</p>

Click the button above and pick the installer for your platform.

### Get started in 3 steps

1. **Install** the app from the Releases page
2. **Launch** — the clock starts ticking immediately, no login, no setup
3. **Customize** — open Settings (`⚙`) to pick a theme, toggle 12/24h, or set your motto

---

## 🛠️ Tech Stack

- **Go 1.25+** — application runtime
- **Wails 3 (alpha)** — native window + webview shell
- **HTML / CSS / JavaScript** — the entire UI
- **jQuery** + **FlipClock** — the split-flap animation
- **macOS Cocoa** — `ActivationPolicyRegular` + fullscreen Space

This project uses a local copy of Wails 3 via `go.mod`:

```text
replace github.com/wailsapp/wails/v3 => ./third_party/wails-v3
```

> ⚠️ Do not delete `third_party/wails-v3` — it carries local macOS fullscreen patches.

---

## 💻 Build from Source

### Prerequisites

- Go 1.22 or higher
- Git
- (macOS only) `lipo`, `hdiutil`, `plutil`, `iconutil` — all pre-installed on macOS

### Local run

```bash
./scripts/run.sh
```

The dev binary is generated to `build/dev/easy-flip-clock` — never to the project root.

### Run tests

```bash
go test ./...
```

Coverage:

- macOS uses `ActivationPolicyRegular` to guarantee native fullscreen
- macOS window enables `NSWindowCollectionBehaviorFullScreenPrimary`

### Build a macOS DMG

```bash
./scripts/build-dmg.sh
```

Outputs:

```text
build/macos-arm64/easy-flip-clock
build/macos-amd64/easy-flip-clock
build/macos-universal/easy-flip-clock
easy-flip-clock.app
build/easy-flip-clock.dmg
```

### Install the built `.app`

```bash
open build/easy-flip-clock.dmg
cp -R "/Volumes/easy-flip-clock/easy-flip-clock.app" /Applications/
```

---

## ⚙️ Configuration

Stored at:

```text
~/.easy-flip-clock/config.json
```

Default:

```json
{
  "motto": "君子三思而后行",
  "width": 600,
  "height": 300,
  "x": -1,
  "y": -1,
  "show_in_dock": false
}
```

> Note: `show_in_dock` is currently a no-op at runtime — hiding the Dock icon breaks the macOS fullscreen Space, so the app always uses `ActivationPolicyRegular`.

---

## 🧱 Project Structure

```text
.
├── app.go                  # App entry + Wails bindings (fullscreen, config)
├── config.go               # Local config read/write
├── frontend.go             # Frontend asset binding
├── update.go               # GitHub release auto-update
├── app_test.go             # macOS fullscreen behavior tests
├── frontend/               # Static frontend assets
│   ├── index.html
│   ├── css/                # flipclock.css + styles.css
│   ├── js/                 # jQuery + flipclock.min.js
│   └── imgs/               # App icon (svg / png / icns)
├── scripts/
│   ├── run.sh              # Local compile + run
│   └── build-dmg.sh        # macOS DMG packager
├── build/                  # Build artifacts
└── third_party/wails-v3/   # Local Wails 3 dependency
```

---

## 🪟 Fullscreen Behavior

The app guarantees **real** macOS fullscreen — not a maximized window on the desktop layer.

How to enter:

- Click the green fullscreen button in the title bar
- Press `Ctrl + Cmd + F`
- Press `F11` (front-end triggers the back-end method)

Verify with:

```bash
osascript <<'APPLESCRIPT'
tell application "System Events"
    tell process "easy-flip-clock"
        tell window 1
            return "AXFullScreen=" & (value of attribute "AXFullScreen" as text)
        end tell
    end tell
end tell
APPLESCRIPT
```

If you see `AXFullScreen=true` — the window is in a real macOS fullscreen Space. 🎉

---

## ❓ FAQ

<details>
<summary><strong>Q: Why does the app show a Dock icon?</strong></summary>
A: Hiding the Dock icon (via `ActivationPolicyAccessory`) breaks macOS native fullscreen — the window falls back to a regular desktop-level maximize. The app uses `ActivationPolicyRegular` to keep the green fullscreen button and `Ctrl + Cmd + F` working correctly.
</details>

<details>
<summary><strong>Q: Can I run it on Windows / Linux?</strong></summary>
A: The Wails 3 project supports all three platforms in principle. The packaged release currently ships a macOS DMG first — Windows / Linux builds are welcome as community PRs.
</details>

<details>
<summary><strong>Q: Where is my config stored?</strong></summary>
A: At `~/.easy-flip-clock/config.json` on macOS / Linux, and `%USERPROFILE%\.easy-flip-clock\config.json` on Windows.
</details>

<details>
<summary><strong>Q: How do I add a new theme?</strong></summary>
A: Add a CSS class under `frontend/css/styles.css`, then register the option value in `frontend/index.html`'s theme `<select>`. No rebuild needed for theme changes — just reload the app.
</details>

<details>
<summary><strong>Q: Does it phone home?</strong></summary>
A: No. The only network call is the explicit "Check for update" action in Settings, which talks to the public GitHub Releases API.
</details>

---

## 🤝 Contributing

PRs are welcome — themes, bug fixes, Windows / Linux packaging, and new flip animations are all good first issues.

1. Fork this project
2. Create a feature branch (`git checkout -b feature/AmazingTheme`)
3. Commit your changes (`git commit -m 'Add: sunset-glow theme'`)
4. Push to the branch (`git push origin feature/AmazingTheme`)
5. Open a Pull Request

### Development notes

- Never output dev binaries to the project root — they go to `build/dev/`.
- DMG output goes to `build/`.
- Modifying Wails 3 macOS behavior? Synchronize the local patch under `third_party/wails-v3/pkg/application/`.
- If you change the fullscreen path, verify `AXFullScreen=true` end-to-end.

---

## 📄 License

This project is licensed under the [BSD 4-Clause License](LICENSE).

---

## Contributors

<p align="center">
  <a href="https://github.com/smile-yan/easy-flip-clock/graphs/contributors">
    <img src="https://contrib.rocks/image?repo=smile-yan/easy-flip-clock&max=100" alt="Contributors" />
  </a>
</p>

## ⭐ Star History

<p align="center">
  <a href="https://www.star-history.com/#smile-yan/easy-flip-clock&Date" target="_blank">
    <img src="https://api.star-history.com/svg?repos=smile-yan/easy-flip-clock&type=Date" alt="Star History" width="600">
  </a>
</p>

<div align="center">

**If you like it, give it a star ⭐**

[Report Bug](https://github.com/smile-yan/easy-flip-clock/issues) · [Request Feature](https://github.com/smile-yan/easy-flip-clock/issues)

</div>
