package main

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type Config struct {
	Profile  string   `json:"profile"`
	Enabled  []string `json:"enabled"`
	Disabled []string `json:"disabled"`
	Hidden   []string `json:"hidden"`
	Path     string   `json:"-"`
	RepoRoot string   `json:"-"`
}

func defaultConfig() Config {
	return Config{Profile: "owner"}
}

func configPath() string {
	if p := os.Getenv("YAGA_CONFIG"); p != "" {
		return p
	}
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".config", "yaga", "config.json")
}

func loadConfig() Config {
	cfg := defaultConfig()
	cfg.Path = configPath()
	cfg.RepoRoot = findRepoRoot()
	if b, err := os.ReadFile(cfg.Path); err == nil {
		_ = json.Unmarshal(b, &cfg)
		cfg.Path = configPath()
		cfg.RepoRoot = findRepoRoot()
	}
	if p := os.Getenv("YAGA_PROFILE"); p != "" {
		cfg.Profile = p
	}
	if cfg.Profile == "" {
		cfg.Profile = "owner"
	}
	return cfg
}

func saveConfig(cfg Config) error {
	dir := filepath.Dir(cfg.Path)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}
	out := struct {
		Profile  string   `json:"profile"`
		Enabled  []string `json:"enabled"`
		Disabled []string `json:"disabled"`
		Hidden   []string `json:"hidden"`
	}{cfg.Profile, cfg.Enabled, cfg.Disabled, cfg.Hidden}
	b, err := json.MarshalIndent(out, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(cfg.Path, append(b, '\n'), 0o644)
}

func findRepoRoot() string {
	if r := os.Getenv("YAGA_REPO"); r != "" {
		return r
	}
	// cli/yaga → repo root
	exe, err := os.Executable()
	if err == nil {
		dir := filepath.Dir(exe)
		for i := 0; i < 6; i++ {
			if fileExists(filepath.Join(dir, "package.json")) && fileExists(filepath.Join(dir, "scripts")) {
				return dir
			}
			dir = filepath.Dir(dir)
		}
	}
	wd, _ := os.Getwd()
	dir := wd
	for i := 0; i < 8; i++ {
		if fileExists(filepath.Join(dir, "package.json")) && fileExists(filepath.Join(dir, "scripts")) {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}
	// fallback: relative to source when go run
	return filepath.Clean(filepath.Join(mustSourceDir(), "..", ".."))
}

func mustSourceDir() string {
	wd, _ := os.Getwd()
	return wd
}

func fileExists(p string) bool {
	_, err := os.Stat(p)
	return err == nil
}

func scriptsDir(cfg Config) string {
	return filepath.Join(cfg.RepoRoot, "scripts")
}
