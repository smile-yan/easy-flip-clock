package main

import (
	"context"
	"embed"
	"fmt"
	"log"
	"os"

	"github.com/wailsapp/wails/v3/pkg/application"
)

type App struct {
	config *Config
}

//go:embed all:frontend
var assets embed.FS

var mainWindow *application.WebviewWindow
var globalApp *application.App

func (a *App) ToggleFullscreen() {
	log.Println("ToggleFullscreen called, mainWindow:", mainWindow)
	if mainWindow != nil {
		if mainWindow.IsFullscreen() {
			log.Println("Unfullscreen")
			mainWindow.UnFullscreen()
		} else {
			log.Println("Fullscreen")
			mainWindow.Fullscreen()
		}
	}
}

func (a *App) startup(runtime any) {
	cfg, err := Load()
	if err != nil {
		log.Printf("Failed to load config, using defaults: %v", err)
		cfg = DefaultConfig()
	}
	a.config = cfg
}

func (a *App) shutdown() {
}

func (a *App) GetConfig() map[string]any {
	if a.config == nil {
		log.Println("[App] GetConfig called with nil config, using defaults")
		a.config = DefaultConfig()
	}
	applyConfigDefaults(a.config)
	log.Printf("[App] GetConfig: motto=%q theme=%s style=%s time_format=%s show_date=%v show_seconds=%v show_lunar=%v show_motto=%v",
		a.config.Motto, a.config.Theme, a.config.Style, a.config.TimeFormat,
		a.config.ShowDate, a.config.ShowSeconds, a.config.ShowLunar, a.config.ShowMotto)
	return map[string]any{
		"motto":        a.config.Motto,
		"width":        a.config.Width,
		"height":       a.config.Height,
		"x":            a.config.X,
		"y":            a.config.Y,
		"show_in_dock": a.config.ShowInDock,
		"theme":        a.config.Theme,
		"style":        a.config.Style,
		"time_format":  a.config.TimeFormat,
		"show_date":    a.config.ShowDate,
		"show_seconds": a.config.ShowSeconds,
		"show_lunar":   a.config.ShowLunar,
		"show_motto":   a.config.ShowMotto,
		"color":        a.config.Color,
	}
}

func (a *App) SaveConfig(motto string, showInDock bool) error {
	if a.config == nil {
		a.config = DefaultConfig()
	}
	a.config.Motto = motto
	a.config.ShowInDock = showInDock
	return Save(a.config)
}

// SaveTheme 仅更新主题设置并持久化。
func (a *App) SaveTheme(theme string) error {
	if a.config == nil {
		a.config = DefaultConfig()
	}
	if !isValidTheme(theme) {
		return fmt.Errorf("unsupported theme: %s", theme)
	}
	a.config.Theme = theme
	return Save(a.config)
}

// SaveSettings 一次性保存设置面板中可调的所有字段。
// 参数使用 map[string]any，因为当前 Wails v3 绑定层不会把 JS 对象自动转换成 Go 结构体。
// Pro 字段（color）即使传入也会被保存，仅作为占位。
func (a *App) SaveSettings(payload map[string]any) error {
	log.Printf("[App] SaveSettings called with payload=%+v", payload)
	if a.config == nil {
		log.Println("[App] SaveSettings: config is nil, initializing defaults")
		a.config = DefaultConfig()
	}
	applyConfigDefaults(a.config)

	// motto 允许空字符串（用户清空座右铭），不要用非空守卫
	if v, ok := payload["motto"].(string); ok {
		a.config.Motto = v
	}
	if v, ok := payload["show_in_dock"].(bool); ok {
		a.config.ShowInDock = v
	}
	if v, ok := payload["theme"].(string); ok && isValidTheme(v) {
		a.config.Theme = v
	}
	if v, ok := payload["style"].(string); ok && isValidStyle(v) {
		a.config.Style = v
	}
	if v, ok := payload["time_format"].(string); ok && isValidTimeFormat(v) {
		a.config.TimeFormat = v
	}
	if v, ok := payload["show_date"].(bool); ok {
		a.config.ShowDate = v
	}
	if v, ok := payload["show_seconds"].(bool); ok {
		a.config.ShowSeconds = v
	}
	if v, ok := payload["show_lunar"].(bool); ok {
		a.config.ShowLunar = v
	}
	if v, ok := payload["show_motto"].(bool); ok {
		a.config.ShowMotto = v
	}
	// color 字段属于 Pro 功能，目前只存不生效。
	if v, ok := payload["color"].(string); ok {
		a.config.Color = v
	}
	log.Printf("[App] SaveSettings: updated config motto=%q theme=%s style=%s time_format=%s", a.config.Motto, a.config.Theme, a.config.Style, a.config.TimeFormat)
	if err := Save(a.config); err != nil {
		log.Printf("[App] SaveSettings failed: %v", err)
		return err
	}
	log.Println("[App] SaveSettings succeeded")
	return nil
}

func isValidTheme(theme string) bool {
	for _, t := range AvailableThemes {
		if t == theme {
			return true
		}
	}
	return false
}

func isValidStyle(style string) bool {
	for _, s := range AvailableStyles {
		if s == style {
			return true
		}
	}
	return false
}

func isValidTimeFormat(format string) bool {
	for _, f := range AvailableTimeFormats {
		if f == format {
			return true
		}
	}
	return false
}

// applyConfigDefaults 对从老版本配置文件读出的 Config 做字段补全。
func applyConfigDefaults(cfg *Config) {
	if cfg.Theme == "" {
		cfg.Theme = DefaultTheme
	}
	if cfg.Style == "" {
		cfg.Style = DefaultStyle
	}
	if cfg.TimeFormat == "" {
		cfg.TimeFormat = DefaultTimeFormat
	}
}

func (a *App) BeforeClose(ctx context.Context) bool {
	if a.config != nil {
		if err := Save(a.config); err != nil {
			log.Printf("Failed to save config on close: %v", err)
		}
	}
	// 返回 false，close button 可以关闭应用（BeforeClose 返回值只阻止 JS 触发，Cmd+Q 和 close button 不受影响）
	return false
}

func macOptionsForConfig(cfg *Config) application.MacOptions {
	return application.MacOptions{
		ActivationPolicy: application.ActivationPolicyRegular,
		ApplicationShouldTerminateAfterLastWindowClosed: true,
	}
}

// createCustomMenuBar 创建自定义菜单栏，移除 File 和 Edit 菜单
func createCustomMenuBar(result *UpdateResult) *application.Menu {
	menu := application.NewMenu()

	// 添加 App 菜单（应用菜单）
	appMenu := menu.AddSubmenu("easy-flip-clock")
	appMenu.Add("关于").OnClick(func(ctx *application.Context) {
		globalApp.ShowAboutDialog()
	})
	appMenu.AddSeparator()
	appMenu.Add("设置").OnClick(func(ctx *application.Context) {
		// 通过事件通知前端打开设置弹窗
		globalApp.Events.Emit(&application.WailsEvent{
			Name: "open-settings",
		})
	})
	appMenu.Add("检查更新").OnClick(func(ctx *application.Context) {
		log.Printf("检查更新结果: %+v", result)

		if result.HasUpdate {
			// 发现新版本，显示询问对话框
			// 注意：不使用 AttachToWindow，因为 sheet 模式下自定义按钮回调可能不触发
			dialog := application.QuestionDialog()
			dialog.SetTitle("发现新版本")
			dialog.SetMessage(fmt.Sprintf("当前版本: %s\n最新版本: %s\n\n更新说明:\n%s\n\n是否前往下载？", result.CurrentVer, result.LatestVer, result.ReleaseNote))

			// 添加"前往下载"按钮，点击后打开浏览器
			downloadBtn := dialog.AddButton("前往下载")
			downloadBtn.SetAsDefault()
			downloadBtn.OnClick(func() {
				log.Printf("用户点击下载，打开URL: %s", result.DownloadURL)
				globalApp.BrowserOpenURL(result.DownloadURL)
			})

			// 添加"稍后再说"按钮
			laterBtn := dialog.AddButton("稍后再说")
			laterBtn.SetAsCancel()

			log.Println("显示更新对话框...")
			dialog.Show()
			log.Println("对话框已显示")
		} else {
			// 已是最新版本，显示提示对话框
			dialog := application.InfoDialog()
			dialog.SetTitle("检查更新")
			dialog.SetMessage(fmt.Sprintf("当前已是最新版本 %s", result.CurrentVer))
			dialog.Show()
		}
	})
	appMenu.AddSeparator()
	quitItem := appMenu.Add("退出")
	quitItem.SetAccelerator("cmd+q")
	quitItem.OnClick(func(ctx *application.Context) {
		globalApp.Quit()
	})

	// 添加 Window 菜单（窗口菜单）
	windowMenu := menu.AddSubmenu("窗口")
	windowMenu.AddRole(application.Minimize)
	windowMenu.AddRole(application.Zoom)
	windowMenu.AddSeparator()
	fullscreenItem := windowMenu.Add("进入全屏")
	fullscreenItem.SetAccelerator("ctrl+cmd+f")
	fullscreenItem.OnClick(func(ctx *application.Context) {
		if mainWindow != nil {
			if mainWindow.IsFullscreen() {
				mainWindow.UnFullscreen()
			} else {
				mainWindow.Fullscreen()
			}
		}
	})

	return menu
}

func main() {
	cfg, err := Load()
	if err != nil {
		log.Printf("Failed to load config, using defaults: %v", err)
		cfg = DefaultConfig()
	}
	log.Printf("[main] Loaded config motto=%q theme=%s style=%s time_format=%s", cfg.Motto, cfg.Theme, cfg.Style, cfg.TimeFormat)

	// 读取应用图标
	iconPath := "frontend/imgs/app-icon-1024.png"
	iconData, err := os.ReadFile(iconPath)
	if err != nil {
		log.Printf("Failed to load icon: %v", err)
	}

	// 关于对话框描述：版本号 + 开源地址
	versionResult := CheckForUpdate()
	description := versionResult.CurrentVer + "\nhttps://github.com/smile-yan/easy-flip-clock"

	app := &App{config: cfg}

	globalApp = application.New(application.Options{
		Name:        "easy-flip-clock",
		Description: description,
		Icon:        iconData,
		Assets: application.AssetOptions{
			FS: assets,
		},
		Mac: macOptionsForConfig(cfg),
		Bind: []any{
			app,
		},
	})

	// 设置自定义菜单栏
	globalApp.SetMenu(createCustomMenuBar(versionResult))

	mainWindow = globalApp.NewWebviewWindowWithOptions(application.WebviewWindowOptions{
		Title:  "翻转时钟",
		Width:  cfg.Width,
		Height: cfg.Height,
		X:      cfg.X,
		Y:      cfg.Y,
		ShouldClose: func(window *application.WebviewWindow) bool {
			log.Println("[main] ShouldClose triggered, saving config")
			if err := Save(app.config); err != nil {
				log.Printf("[main] Failed to save config on close: %v", err)
			} else {
				log.Println("[main] Config saved on close")
			}
			return true
		},
	})

	err = globalApp.Run()
	if err != nil {
		log.Fatal(err)
	}
}
