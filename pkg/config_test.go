package pkg

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

func TestGetDefaultConfigPath(t *testing.T) {
	t.Run("XDG_CONFIG_HOME is set", func(t *testing.T) {
		t.Setenv("XDG_CONFIG_HOME", "/home/testuser")
		t.Setenv("HOME", "/home/fallback") // even if both are set, XDG_CONFIG_HOME should take precedence

		expected := filepath.Join("/home/testuser", "hyprgo.conf")
		got := GetDefaultConfigPath()
		if got != expected {
			t.Errorf("expected %s, got %s", expected, got)
		}
	})

	t.Run("XDG_CONFIG_HOME is not set", func(t *testing.T) {
		t.Setenv("HOME", "/home/testuser")
		_ = os.Unsetenv("XDG_CONFIG_HOME")

		expected := filepath.Join("/home/testuser", ".config", "hyprgo.conf")
		got := GetDefaultConfigPath()
		if got != expected {
			t.Errorf("expected %s, got %s", expected, got)
		}
	})
}

func TestGetWallpapers(t *testing.T) {
	t.Run("Print wallpapers", func(t *testing.T) {
		fmt.Println(GetWallpapers("/home/personal"))
	})

}
