package main

import (
	"os"
	"strings"
	"testing"

	"github.com/wailsapp/wails/v3/pkg/application"
)

func TestMacOptionsUseRegularPolicyForNativeFullscreenByDefault(t *testing.T) {
	options := macOptionsForConfig(DefaultConfig())

	if options.ActivationPolicy != application.ActivationPolicyRegular {
		t.Fatalf("expected default config to use regular activation for native fullscreen, got %d", options.ActivationPolicy)
	}
}

func TestMacOptionsUseRegularPolicyEvenWhenDockIsConfiguredHidden(t *testing.T) {
	cfg := DefaultConfig()
	cfg.ShowInDock = false

	options := macOptionsForConfig(cfg)

	if options.ActivationPolicy != application.ActivationPolicyRegular {
		t.Fatalf("expected native fullscreen to keep regular activation policy, got %d", options.ActivationPolicy)
	}
}

func TestDarwinWindowEnablesNativeFullscreen(t *testing.T) {
	source, err := os.ReadFile("third_party/wails-v3/pkg/application/webview_window_darwin.go")
	if err != nil {
		t.Fatal(err)
	}

	if !strings.Contains(string(source), "NSWindowCollectionBehaviorFullScreenPrimary") {
		t.Fatal("darwin windows must opt in to native fullscreen so macOS shows fullscreen instead of zoom")
	}
}
