# Flip Clock (Wails3) Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** 使用 Wails3 重写翻转时钟，Go 负责窗口管理和配置持久化，前端保持原 HTML/CSS/JS 结构

**Architecture:** Wails3 应用，前端资源在 `frontend/` 目录，Go 代码在 `src/` 和 `internal/` 目录，通过 `wailsjs` bridge 打通前后端

**Tech Stack:** Wails3, Go, HTML/CSS/JS, jQuery FlipClock

---

## 文件结构

```
elect-flip-clock/
├── frontend/               # 前端资源（从根目录迁移）
│   ├── index.html
│   ├── main.js
│   ├── renderer.js
│   ├── preload.js
│   ├── css/
│   ├── js/
│   └── imgs/
├── internal/
│   └── config/
│       └── config.go      # 配置管理模块
├── src/
│   └── main.go           # Wails main.go（窗口创建）
├── app.go                 # Wails App struct
├── frontend.go            # Wails 前后端绑定（暴露方法给前端）
├── wails.json            # Wails 项目配置
└── docs/
    └── specs/
        └── 2026-06-03-flip-clock-wails3-design.md
```

---

## Task 1: 初始化 Wails3 项目

**Files:**
- Create: `wails.json`
- Create: `app.go`
- Create: `frontend.go`
- Create: `src/main.go`
- Create: `internal/config/config.go`
- Delete: `main.js`, `preload.js`, `renderer.js` (Electron 入口)

- [ ] **Step 1: 创建 wails.json 配置文件**

```json
{
  "$schema": "https://wails.io/schemas/config.v2.json",
  "name": "elect-flip-clock",
  "outputfilename": "elect-flip-clock",
  "frontend:install": "",
  "frontend:build": "",
  "frontend:dev": "",
  "author": "",
  "version": "",
  "outputType": "desktop",
  "platform": ["darwin", "windows", "linux"],
  "strict": {}
}
```

- [ ] **Step 2: 创建 app.go (Wails App struct)**

```go
package main

import (
    "embed"
    "log"
)

//go:embed all:frontend
var assets embed.FS

type App struct {
    runtime any
}

func (a *App) startup(runtime any) {
    a.runtime = runtime
}

func (a *App) shutdown() {
}
```

- [ ] **Step 3: 创建 frontend.go (前后端绑定)**

```go
package main

import (
    "elect-flip-clock/internal/config"
)

// GetConfig 获取当前配置
func GetConfig() map[string]interface{} {
    cfg, err := config.Load()
    if err != nil {
        cfg = &config.Config{}
    }
    return map[string]interface{}{
        "motto":       cfg.Motto,
        "width":       cfg.Width,
        "height":      cfg.Height,
        "x":           cfg.X,
        "y":           cfg.Y,
        "showInDock":  cfg.ShowInDock,
    }
}

// SaveConfig 保存配置
func SaveConfig(motto string, showInDock bool) error {
    cfg, err := config.Load()
    if err != nil {
        cfg = &config.Config{}
    }
    cfg.Motto = motto
    cfg.ShowInDock = showInDock
    return config.Save(cfg)
}
```

- [ ] **Step 4: 创建 src/main.go (Wails 入口)**

```go
package main

import (
    "embed"
    "fmt"
)

//go:embed all:frontend
var assets embed.FS

func main() {
    // Wails3 入口，不再需要 main 函数
    // 由 wails 框架自动生成
}
```

- [ ] **Step 5: 删除 Electron 相关文件**

```bash
rm -f main.js preload.js renderer.js
```

---

## Task 2: 实现 Config 模块

**Files:**
- Create: `internal/config/config.go`

- [ ] **Step 1: 创建 internal/config/config.go**

```go
package config

import (
    "encoding/json"
    ""os"
    "path/filepath"

    "github.com/mitchellh/go-homedir"
)

type Config struct {
    Motto      string `json:"motto"`
    Width      int    `json:"width"`
    Height     int    `json:"height"`
    X          int    `json:"x"`
    Y          int    `json:"y"`
    ShowInDock bool   `json:"show_in_dock"`
}

func configPath() (string, error) {
    home, err := go-homedir.Dir()
    if err != nil {
        return "", err
    }
    dir := filepath.Join(home, ".elect-flip-clock")
    if err := os.MkdirAll(dir, 0755); err != nil {
        return "", err
    }
    return filepath.Join(dir, "config.json"), nil
}

func Load() (*Config, error) {
    path, err := configPath()
    if err != nil {
        return nil, err
    }

    data, err := os.ReadFile(path)
    if err != nil {
        if os.IsNotExist(err) {
            return &Config{
                Motto:      "君子三思而后行",
                Width:      600,
                Height:     300,
                ShowInDock: false,
            }, nil
        }
        return nil, err
    }

    var cfg Config
    if err := json.Unmarshal(data, &cfg); err != nil {
        return nil, err
    }
    return &cfg, nil
}

func Save(cfg *Config) error {
    path, err := configPath()
    if err != nil {
        return err
    }

    data, err := json.MarshalIndent(cfg, "", "  ")
    if err != nil {
        return err
    }

    return os.WriteFile(path, data, 0644)
}
```

- [ ] **Step 2: 提交 config 模块**

```bash
git add internal/config/config.go wails.json app.go frontend.go
git commit -m "feat: add wails3 project scaffold and config module"
```

---

## Task 3: 迁移前端文件到 frontend/

**Files:**
- Move: `index.html` → `frontend/index.html`
- Move: `css/` → `frontend/css/`
- Move: `js/` → `frontend/js/`
- Move: `imgs/` → `frontend/imgs/`
- Move: `demos/` → `frontend/demos/`

- [ ] **Step 1: 创建 frontend 目录并迁移文件**

```bash
mkdir -p frontend
mv index.html frontend/
mv css frontend/
mv js frontend/
mv imgs frontend/
mv demos frontend/
```

- [ ] **Step 2: 更新 frontend/index.html 中的路径引用**

将 `<script src="js/flipclock.min.js">` 和 `<link rel="stylesheet" href="css/flipclock.css">` 等路径更新为相对路径（保持不变因为结构一致）

- [ ] **Step 3: 提交迁移**

```bash
git add frontend/
git commit -m "feat: move frontend assets to frontend/ directory"
```

---

## Task 4: 实现窗口管理和 Dock 隐藏

**Files:**
- Modify: `app.go`
- Modify: `frontend.go`

- [ ] **Step 1: 更新 app.go 添加窗口配置**

```go
package main

import (
    "embed"
    "log"

    "elect-flip-clock/internal/config"
    "github.com/wailsapp/wails/v3/pkg/application"
)

//go:embed all:frontend
var assets embed.FS

type App struct {
    runtime   any
    app       *application.Application
    config    *config.Config
}

func (a *App) startup(runtime any) {
    a.runtime = runtime

    cfg, err := config.Load()
    if err != nil {
        log.Printf("failed to load config: %v", err)
        cfg = &config.Config{}
    }
    a.config = cfg

    app := runtime.(*application.Application)
    a.app = app

    // 设置窗口尺寸
    if cfg.Width > 0 && cfg.Height > 0 {
        app.DefaultWindow().SetSize(cfg.Width, cfg.Height)
    }

    // 设置窗口位置
    if cfg.X > 0 || cfg.Y > 0 {
        app.DefaultWindow().SetPosition(cfg.X, cfg.Y)
    }

    // macOS: 根据 ShowInDock 配置决定是否显示 dock 图标
    #ifdef darwin
    if !cfg.ShowInDock {
        app.SetDockVisible(false)
    }
    #endif
}

func (a *App) shutdown() {
    // 保存窗口状态
    if a.app != nil && a.config != nil {
        win := a.app.DefaultWindow()
        pos := win.Position()
        size := win.Size()
        a.config.X = pos.X
        a.config.Y = pos.Y
        a.config.Width = size.Width
        a.config.Height = size.Height
        config.Save(a.config)
    }
}
```

- [ ] **Step 2: 更新 frontend.go 添加窗口状态保存**

```go
package main

import (
    "elect-flip-clock/internal/config"
)

// GetConfig 获取当前配置
func GetConfig() map[string]interface{} {
    cfg, err := config.Load()
    if err != nil {
        cfg = &config.Config{}
    }
    return map[string]interface{}{
        "motto":       cfg.Motto,
        "width":       cfg.Width,
        "height":      cfg.Height,
        "x":           cfg.X,
        "y":           cfg.Y,
        "showInDock":  cfg.ShowInDock,
    }
}

// SaveConfig 保存配置
func SaveConfig(motto string, showInDock bool) error {
    cfg, err := config.Load()
    if err != nil {
        cfg = &config.Config{}
    }
    cfg.Motto = motto
    cfg.ShowInDock = showInDock
    return config.Save(cfg)
}
```

- [ ] **Step 3: 提交**

```bash
git add app.go frontend.go
git commit -m "feat: add window management and dock visibility control"
```

---

## Task 5: 删除 Electron 相关文件

**Files:**
- Delete: `electron.js` (如有)
- Delete: `package.json`, `package-lock.json`
- Delete: `pnpm-lock.yaml`
- Delete: `node_modules/`
- Delete: `demos/` (如果没迁移的话)

- [ ] **Step 1: 删除 Electron 和 Node 相关文件**

```bash
rm -rf node_modules
rm -f package.json package-lock.json pnpm-lock.yaml
rm -f main.js preload.js renderer.js
```

- [ ] **Step 2: 清理 .gitignore 中无关项，提交清理**

```bash
git add -A
git commit -m "chore: remove electron and node dependencies"
```

---

## Task 6: 验证构建

- [ ] **Step 1: 初始化 Go 模块**

```bash
go mod init elect-flip-clock
wails doctor
```

- [ ] **Step 2: 本地开发验证**

```bash
wails dev
```

- [ ] **Step 3: 生产构建验证**

```bash
wails build
```

---

## 依赖

| 库 | 用途 |
|----|------|
| `github.com/wailsapp/wails/v3` | Wails3 框架 |
| `github.com/mitchellh/go-homedir` | 获取用户目录 |
| `github.com/jsummers/lofig` | （如有需要）绘图 |

---

## 验证清单

- [ ] `wails dev` 能启动应用
- [ ] 翻转时钟正常显示
- [ ] 日期+星期正确显示
- [ ] 格言正确显示
- [ ] 窗口可调整大小
- [ ] 关闭应用后重新打开，窗口位置和尺寸保持
- [ ] macOS 下 `ShowInDock=false` 时 dock 图标隐藏
- [ ] `wails build` 能成功构建