package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

const (
	wantModule = "github.com/fuwiak/bober-yaga"
	wantRepo   = "https://github.com/fuwiak/bober-yaga"
	oldRepo    = "github.com/fuwiak/yaga"
)

func TestGoModModulePath(t *testing.T) {
	raw, err := os.ReadFile("go.mod")
	if err != nil {
		t.Fatal(err)
	}
	first := strings.SplitN(string(raw), "\n", 2)[0]
	want := "module " + wantModule
	if first != want {
		t.Fatalf("go.mod module line = %q, want %q", first, want)
	}
	if strings.Contains(string(raw), oldRepo) {
		t.Fatalf("go.mod still references old path %q", oldRepo)
	}
}

func TestDocsUseNewRepoName(t *testing.T) {
	files := []string{"README.md", "CONTRIBUTING.md"}
	for _, name := range files {
		raw, err := os.ReadFile(name)
		if err != nil {
			t.Fatalf("%s: %v", name, err)
		}
		text := string(raw)

		if strings.Contains(text, "https://"+oldRepo) || strings.Contains(text, oldRepo+".git") {
			t.Errorf("%s: still points at old GitHub repo %q", name, oldRepo)
		}
		if !strings.Contains(text, wantRepo) {
			t.Errorf("%s: missing new repo URL %q", name, wantRepo)
		}
		if !strings.Contains(text, "cd bober-yaga") {
			t.Errorf("%s: clone instructions must use `cd bober-yaga`", name)
		}
		if strings.Contains(text, "\ncd yaga\n") || strings.HasSuffix(text, "\ncd yaga") {
			t.Errorf("%s: stale `cd yaga` after clone", name)
		}
	}

	readme, err := os.ReadFile("README.md")
	if err != nil {
		t.Fatal(err)
	}
	first := strings.SplitN(string(readme), "\n", 2)[0]
	if first != "# bober-yaga" {
		t.Fatalf("README title = %q, want %q", first, "# bober-yaga")
	}
	if strings.HasPrefix(string(readme), "# yaga\n") {
		t.Fatal("README still titled # yaga")
	}
}

func TestNoStaleGitHubNameInTrackedDocs(t *testing.T) {
	// Guard rails: any markdown/go.mod under module root must not advertise the old repo URL.
	entries, err := os.ReadDir(".")
	if err != nil {
		t.Fatal(err)
	}
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		name := e.Name()
		ext := filepath.Ext(name)
		if name != "go.mod" && ext != ".md" {
			continue
		}
		raw, err := os.ReadFile(name)
		if err != nil {
			t.Fatal(err)
		}
		if strings.Contains(string(raw), "https://"+oldRepo) {
			t.Errorf("%s: contains old URL https://%s", name, oldRepo)
		}
	}
}
