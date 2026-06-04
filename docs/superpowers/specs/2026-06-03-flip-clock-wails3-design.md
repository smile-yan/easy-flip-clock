# Flip Clock (Wails3) 设计文档

## 概述

使用 Wails3 + Go 重写翻转时钟桌面应用，前端保持原有 HTML/CSS/JS，Go 负责窗口管理、配置持久化托盘和快捷键。

## 功能范围

| 功能 | 描述 | 优先级 |
|------|------|--------|
| 翻转时钟显示 | 24小时制，时分秒翻转动画 | P0 |
| 日期+星期显示 | 顶部显示 `YYYY-MM-DD 星期X` | P0 |
| 格言显示 | 底部显示自定义格言 | P0 |
| 窗口自适应缩放 | 根据窗口大小自动缩放 | P0 |
| 配置持久化 | JSON 文件存储格言、窗口位置、大小、程序坞图标显示状态 | P1 |
| 程序坞图标显示控制 | 用户可配置是否在程序坞显示图标，默认不显示 | P1 |

## 技术选型

- **框架**: Wails3
- **前端**: HTML/CSS/JS（保持原结构）
- **后端**: Go
- **配置存储**: JSON 文件 `~/.elect-flip-clock/config.json`
- **持久化库**: `github.com/mitchellh/go-homedir` + 标准 `os`/`json`

## 目录结构

```
elect-flip-clock/
├── frontend/               # 前端资源（保持原结构）
│   ├── index.html
│   ├── main.js
│   ├── renderer.js
│   ├── preload.js
│   ├── css/
│   ├── js/
│   └── imgs/
├── internal/              # Go 内部包
│   └── config/
│       └── config.go      # 配置管理（加载/保存）
├── src/                   # Wails Go 入口
│   └── main.go           # 窗口创建、托盘、快捷键绑定
├── wails.json            # Wails 项目配置
├── app.go                # Wails App struct
├── frontend.go           # Wails 前后端绑定
└── docs/
    └── specs/
        └── 2026-06-03-flip-clock-wails3-design.md
```

## 核心模块

### 1. config (internal/config/config.go)

**职责**: 管理应用配置（JSON 序列化）

**数据结构**:
```go
type Config struct {
    Motto         string `json:"motto"`      // 格言
    Width         int    `json:"width"`      // 窗口宽度
    Height        int    `json:"height"`     // 窗口高度
    X             int    `json:"x"`          // 窗口 X 位置
    Y             int    `json:"y"`          // 窗口 Y 位置
    ShowInDock    bool   `json:"show_in_dock"` // 是否在程序坞显示图标，默认 false
}
```

**接口**:
- `Load() (*Config, error)` — 从 `~/.elect-flip-clock/config.json` 加载配置，不存在则返回默认配置（ShowInDock 默认为 false）
- `Save(*Config) error` — 保存配置到 JSON 文件

### 2. main (src/main.go)

**职责**: Wails3 应用入口，窗口管理，托盘，快捷键

**Wails 配置**:
- 窗口标题: `翻转时钟`
- 尺寸: 600x300（默认）
- macOS 下根据 `ShowInDock` 配置决定是否隐藏 dock 图标（默认隐藏）
- 启用 `devtools: true`（开发调试用）

**托盘菜单** (可选扩展):
- 显示窗口
- 退出

**快捷键** (可选扩展):
- `Cmd+Shift+F`: 切换窗口显示/隐藏

### 3. frontend (frontend.go)

**职责**: Go 与前端的桥接层，通过 `wails:invoke` 暴露方法给前端调用

**暴露给前端的方法**:
- `GetConfig() map[string]interface{}` — 获取当前配置（包含 ShowInDock）
- `SaveConfig(motto string, showInDock bool) error` — 保存格言和程序坞显示状态

## 数据流

```
[前端 JS]  --wails:invoke-->  [GoBridge]  --调用-->  [Config模块]  --读写-->  [JSON文件]
```

## 窗口生命周期

1. 应用启动 → `Load()` 读取配置 → 创建窗口（使用配置的尺寸/位置）
2. 窗口关闭 → `Save()` 保存当前配置
3. 应用退出 → 配置已持久化

## 构建配置

| 平台 | 输出名 | 备注 |
|------|--------|------|
| macOS | `elect-flip-clock.app` | 根据 ShowInDock 配置决定是否隐藏 dock 图标 |
| Windows | `elect-flip-clock.exe` | |
| Linux | `elect-flip-clock` | |