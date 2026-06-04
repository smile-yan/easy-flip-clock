package main

import (
	"context"
	"embed"
	"log"

	"github.com/wailsapp/wails/v3/pkg/application"
)

type App struct {
	config *Config
}

//go:embed all:frontend
var assets embed.FS

var mainWindow *application.WebviewWindow

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

func main() {
	cfg, err := Load()
	if err != nil {
		log.Printf("Failed to load config, using defaults: %v", err)
		cfg = DefaultConfig()
	}

	app := application.New(application.Options{
		Assets: application.AssetOptions{
			FS: assets,
		},
		Mac: macOptionsForConfig(cfg),
		Bind: []any{
			&App{},
		},
	})

	mainWindow = app.NewWebviewWindowWithOptions(application.WebviewWindowOptions{
		Title:  "翻转时钟",
		Width:  cfg.Width,
		Height: cfg.Height,
		X:      cfg.X,
		Y:      cfg.Y,
	})

	err = app.Run()
	if err != nil {
		log.Fatal(err)
	}
}