# Windows 闪退 Bug 修复清单

> 排查日期：2026-06-16
> 范围：`app.go` / `update.go` / `config.go` / `scripts/build-all.sh` / `.github/workflows/release.yml`
> 症状：Windows 下双击 exe 出现"闪退"（窗口不出现、进程立即或很快退出）

---

## Bug 1（🔴 P0 · 最可能根因）：`os.ReadFile` 使用相对路径读取资源

### 位置
- `app.go:290-294` —— 读取 `frontend/imgs/app-icon-1024.png`
- `update.go:29-46` —— 读取 `wails.json`

### 问题
Windows 下双击 exe 启动时，**当前工作目录不可控**（可能是 `C:\Windows\System32`、用户目录或别的位置），绝不是 exe 所在目录。两处 `os.ReadFile` 都会失败：
- `wails.json` 读不到 → 版本回退为 `0.0.0`（本身不闪退）
- 图标读不到 → `iconData = nil`，某些 Wails 内部路径上可能触发 nil 指针

### 修复
- 方案 A（推荐）：把两个文件都改用 `//go:embed` 编译进二进制
- 方案 B：用 `os.Executable()` 取得 exe 自身目录后拼接路径

```go
//go:embed all:frontend
var assets embed.FS

//go:embed wails.json
var wailsJSON []byte
```

---

## Bug 2（🔴 P0）：`http.Client{}` 无超时 + 启动时同步阻塞

### 位置
- `app.go:297` —— `versionResult := CheckForUpdate()` 在 `globalApp.Run()` 之前同步调用
- `update.go:138` —— `client := &http.Client{}`（**无 Timeout**）

### 问题
- 同步 HTTP 请求 → 整个 app 启动被阻塞
- 无超时 → 网络受限 / DNS 慢 / 防火墙拦截时会挂死几十秒到几分钟
- 用户体感："双击 exe → 什么都没出现 → 关掉"——常被误判为闪退

### 修复
- 给 `http.Client` 设置超时（建议 `Timeout: 5 * time.Second`）
- 把 `CheckForUpdate()` 挪到 goroutine 异步执行（或延后到 `startup` 钩子里）

```go
client := &http.Client{Timeout: 5 * time.Second}
```

---

## Bug 3（🟠 P1）：编译参数缺少 `-H windowsgui`，双击时闪黑窗

### 位置
- `scripts/build-all.sh:24`
- `.github/workflows/release.yml:87`

### 问题
Go 默认把 Windows 可执行文件编译为**控制台子系统**。双击时会先弹一个 cmd 黑窗，再显示 GUI。如果 GUI 启动稍慢，用户看到"黑窗闪一下就关掉"，误判为闪退。

### 修复
在两处 `ldflags` 里追加 `-H windowsgui`：

```bash
go build -ldflags="-s -w -H windowsgui" -o "${BUILD_DIR}/windows-amd64/${OUTPUT_NAME}.exe" .
```

---

## Bug 4（🟠 P1）：快捷键写死了 `cmd+q` / `ctrl+cmd+f`，Windows 语义错误

### 位置
- `app.go:256` —— `quitItem.SetAccelerator("cmd+q")`
- `app.go:267` —— `fullscreenItem.SetAccelerator("ctrl+cmd+f")`

### 问题
- `cmd+q`：Wails 把 `cmd` 解析为 `CmdOrCtrl`，Windows 上等价 `Ctrl+Q`——**能用**
- `ctrl+cmd+f`：会变成 `Ctrl+Cmd+F`，Windows 上**根本没反应**（按 Windows 习惯应该是 `F11` 或 `Ctrl+Shift+F`）

### 修复
```go
// 退出：用 CmdOrCtrl+q 跨平台一致
quitItem.SetAccelerator("CmdOrCtrl+q")

// 全屏：改成 F11 或 Ctrl+Shift+F
fullscreenItem.SetAccelerator("F11")
```

---

## Bug 5（🟡 P2）：`BeforeClose` 和 `ShouldClose` 双写配置

### 位置
- `app.go:188-196` —— `App.BeforeClose`（JS 拦截钩子）
- `app.go:324-332` —— `WebviewWindowOptions.ShouldClose`（真正的关闭拦截）

### 问题
Wails v3 中两者都会在窗口关闭时触发 → 同一份 JSON 配置会被 `Save` 写两次。不会闪退，但是：
- 冗余 IO
- 极端情况下 `app.config` 还未初始化（启动失败时）会 panic

### 修复
只保留 `ShouldClose`（真正生效的那个），删掉 `App.BeforeClose` 里重复的保存逻辑。

---

## 推荐修复顺序

1. Bug 3（一行 `ldflags` 改动，影响最大）
2. Bug 2（启动速度、避免被误判闪退）
3. Bug 1（健壮性，避免 nil 风险）
4. Bug 4（跨平台 UX）
5. Bug 5（代码整洁）

---

## 验证步骤

修复后在 cmd 里直接运行 exe（不要双击）来定位：
```cmd
cd C:\path\to\easy-flip-clock.exe
easy-flip-clock.exe
```

- 无闪退 → Bug 3 修复成功
- 报 `Failed to load icon` / `读取 wails.json 失败` → 还有 Bug 1
- 卡住几秒再显示 → 之前是 Bug 2
- 打印 `A FATAL ERROR HAS OCCURRED` → 属于环境问题，需装 [WebView2 Runtime](https://developer.microsoft.com/microsoft-edge/webview2/)
