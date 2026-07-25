package main

import "testing"

func TestFileExists(t *testing.T) {
	if !fileExists("go.mod") {
		t.Fatal("go.mod should exist when tests run from module root")
	}
	if fileExists("this-path-should-not-exist-yaga") {
		t.Fatal("unexpected path")
	}
}

func TestDefaultConfig(t *testing.T) {
	cfg := defaultConfig()
	if cfg.Profile != "owner" {
		t.Fatalf("profile=%q want owner", cfg.Profile)
	}
}
