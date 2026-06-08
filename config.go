package main

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"

	homedir "github.com/mitchellh/go-homedir"
)

type Config struct {
	Motto       string `json:"motto"`
	Width       int    `json:"width"`
	Height      int    `json:"height"`
	X           int    `json:"x"`
	Y           int    `json:"y"`
	ShowInDock  bool   `json:"show_in_dock"`
	Theme       string `json:"theme"`
	Style       string `json:"style"`
	TimeFormat  string `json:"time_format"`
	ShowDate    bool   `json:"show_date"`
	ShowSeconds bool   `json:"show_seconds"`
	ShowLunar   bool   `json:"show_lunar"`
	ShowMotto   bool   `json:"show_motto"`
	Color       string `json:"color"`
}

// DefaultTheme is the theme used when no theme is configured.
const DefaultTheme = "dark"

// AvailableThemes 列出所有可用的主题名称。
var AvailableThemes = []string{"dark", "light", "sepia", "blue"}

// DefaultStyle is the clock style used when none is configured.
const DefaultStyle = "with-seconds"

// AvailableStyles 列出所有可用的钟面样式。
var AvailableStyles = []string{"with-seconds", "without-seconds"}

// DefaultTimeFormat is the time format used when none is configured.
const DefaultTimeFormat = "24h"

// AvailableTimeFormats 列出所有可用的时间格式。
var AvailableTimeFormats = []string{"24h", "12h"}

func DefaultConfig() *Config {
	return &Config{
		Motto:       "Time flies!",
		Width:       600,
		Height:      300,
		X:           -1,
		Y:           -1,
		ShowInDock:  false,
		Theme:       DefaultTheme,
		Style:       DefaultStyle,
		TimeFormat:  DefaultTimeFormat,
		ShowDate:    true,
		ShowSeconds: true,
		ShowLunar:   false,
		ShowMotto:   true,
		Color:       "",
	}
}

func getConfigPath() (string, error) {
	homeDir, err := homedir.Dir()
	if err != nil {
		return "", err
	}
	configDir := filepath.Join(homeDir, ".easy-flip-clock")
	return filepath.Join(configDir, "config.json"), nil
}

func Load() (*Config, error) {
	configPath, err := getConfigPath()
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return DefaultConfig(), nil
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
	configPath, err := getConfigPath()
	if err != nil {
		return err
	}

	dir := filepath.Dir(configPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(configPath, data, 0644)
}