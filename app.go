package main

import (
	"context"
	"embed"
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
		a.config = DefaultConfig()
	}
	return map[string]any{
		"motto":        a.config.Motto,
		"width":        a.config.Width,
		"height":       a.config.Height,
		"x":            a.config.X,
		"y":            a.config.Y,
		"show_in_dock": a.config.ShowInDock,
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
		ActivationPolicy:                              application.ActivationPolicyRegular,
		ApplicationShouldTerminateAfterLastWindowClosed: true,
	}
}

// createCustomMenuBar 创建自定义菜单栏，移除 File 和 Edit 菜单
func createCustomMenuBar() *application.Menu {
	menu := application.NewMenu()

	// 添加 App 菜单（应用菜单）
	appMenu := menu.AddSubmenu("easy-flip-clock")
	appMenu.Add("关于").OnClick(func(ctx *application.Context) {
		globalApp.ShowAboutDialog()
	})
	appMenu.AddSeparator()
	appMenu.Add("设置").OnClick(func(ctx *application.Context) {
		// TODO: 打开设置界面
		log.Println("打开设置")
	})
	appMenu.Add("检查更新").OnClick(func(ctx *application.Context) {
		// TODO: 检查更新逻辑
		log.Println("检查更新")
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

	// 读取应用图标
	iconPath := "frontend/imgs/app-icon-1024.png"
	iconData, err := os.ReadFile(iconPath)
	if err != nil {
		log.Printf("Failed to load icon: %v", err)
	}

	// 关于对话框描述：版本号 + 开源地址
	description := "v0.0.1\nhttps://github.com/smile-yan/easy-flip-clock"

	globalApp = application.New(application.Options{
		Name:        "easy-flip-clock",
		Description: description,
		Icon:        iconData,
		Assets: application.AssetOptions{
			FS: assets,
		},
		Mac: macOptionsForConfig(cfg),
		Bind: []any{
			&App{},
		},
	})

	// 设置自定义菜单栏
	globalApp.SetMenu(createCustomMenuBar())

	mainWindow = globalApp.NewWebviewWindowWithOptions(application.WebviewWindowOptions{
		Title:  "翻转时钟",
		Width:  cfg.Width,
		Height: cfg.Height,
		X:      cfg.X,
		Y:      cfg.Y,
	})

	err = globalApp.Run()
	if err != nil {
		log.Fatal(err)
	}
}