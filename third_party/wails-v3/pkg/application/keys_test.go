package application

import "testing"

func TestParseAcceleratorDistinguishesCtrlFromCmd(t *testing.T) {
	accelerator, err := parseAccelerator("Ctrl+Command+F")
	if err != nil {
		t.Fatalf("parseAccelerator returned error: %v", err)
	}

	modifiers := map[modifier]bool{}
	for _, value := range accelerator.Modifiers {
		modifiers[value] = true
	}

	if !modifiers[ControlKey] {
		t.Fatalf("expected Ctrl to parse as ControlKey, got %v", accelerator.Modifiers)
	}
	if !modifiers[CmdOrCtrlKey] {
		t.Fatalf("expected Command to parse as CmdOrCtrlKey, got %v", accelerator.Modifiers)
	}
	if accelerator.String() != "cmd+ctrl+f" {
		t.Fatalf("expected accelerator string cmd+ctrl+f, got %q", accelerator.String())
	}
}
