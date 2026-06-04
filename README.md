# easy-easy-flip-clock

一个基于 Wails 3 的桌面翻转时钟应用。应用使用本地静态前端资源渲染 24 小时制翻转时钟，并支持 macOS 原生全屏。

## 功能特性

- 24 小时制翻转时钟
- 日期与星期显示
- 自适应窗口缩放
- macOS 原生全屏支持
- macOS `.app` / `.dmg` 打包
- 自定义应用图标

## 技术栈

- Go 1.25+
- Wails 3 alpha
- HTML / CSS / JavaScript
- jQuery
- FlipClock

本项目通过 `go.mod` 使用本地 Wails 3 代码：

```text
replace github.com/wailsapp/wails/v3 => ./third_party/wails-v3
```

因此不要删除 `third_party/wails-v3`。

## 项目结构

```text
.
├── app.go                  # 应用入口与 Wails 绑定方法
├── config.go               # 本地配置读写
├── frontend.go             # 前端资源绑定
├── frontend/               # 静态前端资源
│   ├── index.html
│   ├── css/
│   ├── js/
│   └── imgs/               # 应用图标资源
├── scripts/
│   ├── run.sh              # 本地编译并运行
│   └── build-dmg.sh        # macOS DMG 打包脚本
├── build/                  # 构建产物
└── third_party/wails-v3/   # 本地 Wails 3 依赖
```

## 环境要求

### 通用

- Go 1.22 或更高版本
- Git

### macOS 打包

macOS DMG 打包脚本依赖系统自带工具：

- `lipo`
- `hdiutil`
- `plutil`
- `iconutil`

## 本地运行

使用脚本编译并运行：

```bash
./scripts/run.sh
```

脚本会把开发二进制生成到：

```text
build/dev/easy-flip-clock
```

不会在项目根目录生成二进制文件。

## 测试

运行 Go 测试：

```bash
go test ./...
```

当前测试覆盖：

- macOS 使用 Regular activation policy 以保证原生全屏
- macOS 窗口启用 `NSWindowCollectionBehaviorFullScreenPrimary`

## macOS 打包

生成 `.dmg`：

```bash
./scripts/build-dmg.sh
```

产物路径：

```text
build/easy-flip-clock.dmg
```

打包脚本会生成 universal binary：

```text
build/macos-arm64/easy-flip-clock
build/macos-amd64/easy-flip-clock
build/macos-universal/easy-flip-clock
```

并将其封装为：

```text
easy-flip-clock.app
```

## 安装

构建 DMG 后，打开：

```bash
open build/easy-flip-clock.dmg
```

然后将 `easy-flip-clock.app` 拖拽或复制到 `/Applications`。

也可以直接用命令复制：

```bash
cp -R "/Volumes/easy-flip-clock/easy-flip-clock.app" /Applications/
```

## 全屏行为

本应用优先保证 macOS 原生全屏行为。进入全屏后，窗口应处于独立的 macOS fullscreen Space，而不是普通桌面上的放大窗口。

可用方式：

- 点击窗口左上角绿色全屏按钮
- 按 `Ctrl + Cmd + F`
- 前端中保留 `F11` 调用后端全屏方法

可以用以下命令验证是否为原生全屏：

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

结果为 `AXFullScreen=true` 时，表示当前窗口处于 macOS 原生全屏。

## Dock 图标说明

当前实现优先保证 macOS 原生全屏。运行时使用 `ActivationPolicyRegular`，这是绿色全屏按钮和 `Ctrl + Cmd + F` 稳定工作的前提。

历史验证中，运行时切换到 `ActivationPolicyAccessory` 虽然可以隐藏 Dock 图标，但会破坏 macOS 原生全屏 Space，导致全屏内容回到普通桌面层。因此该方案已回滚。

## 配置文件

应用配置存储在：

```text
~/.elect-easy-flip-clock/config.json
```

默认配置：

```json
{
  "motto": "Time flies!",
  "width": 600,
  "height": 300,
  "x": -1,
  "y": -1,
  "show_in_dock": false
}
```

注意：`show_in_dock` 字段目前不再用于切换运行时 Dock 行为，因为隐藏 Dock 会影响 macOS 原生全屏。

## 图标资源

当前应用图标资源位于：

```text
frontend/imgs/app-icon.svg
frontend/imgs/app-icon-1024.png
frontend/imgs/app.icns
```

打包脚本会将 `frontend/imgs/app.icns` 复制到 `.app` 的 `Contents/Resources/app.icns`，并在 `Info.plist` 中写入：

```xml
<key>CFBundleIconFile</key>
<string>app.icns</string>
```

## 开发注意事项

- 不要将本地开发二进制输出到项目根目录。
- `scripts/run.sh` 的输出目录是 `build/dev/`。
- `scripts/build-dmg.sh` 的输出目录是 `build/`。
- 修改 Wails 3 macOS 行为时，需要同步检查 `third_party/wails-v3/pkg/application/` 下的本地补丁。
- 如果修改全屏逻辑，请至少验证 `AXFullScreen=true`。

## 许可证

本项目使用 BSD 4-Clause License。详见 [LICENSE](LICENSE)。
