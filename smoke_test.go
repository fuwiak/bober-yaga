package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// Smoke: build binary once, then check version + bricks list still work after rename.
func TestCLISmoke(t *testing.T) {
	dir := t.TempDir()
	bin := filepath.Join(dir, "yaga")

	build := exec.Command("go", "build", "-o", bin, ".")
	build.Env = os.Environ()
	if out, err := build.CombinedOutput(); err != nil {
		t.Fatalf("go build: %v\n%s", err, out)
	}

	t.Run("version", func(t *testing.T) {
		out, err := exec.Command(bin, "version").CombinedOutput()
		if err != nil {
			t.Fatalf("version: %v\n%s", err, out)
		}
		got := string(out)
		if !strings.Contains(got, "yaga") {
			t.Fatalf("version output missing CLI name: %q", got)
		}
	})

	t.Run("bricks_list", func(t *testing.T) {
		out, err := exec.Command(bin, "bricks", "list").CombinedOutput()
		if err != nil {
			t.Fatalf("bricks list: %v\n%s", err, out)
		}
		got := string(out)
		// At least one known public brick should appear.
		if !strings.Contains(got, "webmaster") && !strings.Contains(got, "metrika") {
			t.Fatalf("bricks list missing expected bricks:\n%s", got)
		}
	})

	t.Run("help", func(t *testing.T) {
		out, err := exec.Command(bin, "help").CombinedOutput()
		if err != nil {
			t.Fatalf("help: %v\n%s", err, out)
		}
		if !strings.Contains(string(out), "yaga") {
			t.Fatalf("help output unexpected:\n%s", out)
		}
	})
}
