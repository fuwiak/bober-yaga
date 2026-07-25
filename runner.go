package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func runScript(cfg Config, name string, args []string) error {
	script := scriptsDir(cfg) + "/" + name
	if !fileExists(script) {
		return fmt.Errorf("script not found: %s (repo=%s)", script, cfg.RepoRoot)
	}
	cmd := exec.Command("node", append([]string{script}, args...)...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	cmd.Dir = cfg.RepoRoot
	cmd.Env = os.Environ()
	return cmd.Run()
}

func runScriptCapture(cfg Config, name string, args []string) (string, error) {
	script := scriptsDir(cfg) + "/" + name
	if !fileExists(script) {
		return "", fmt.Errorf("script not found: %s", script)
	}
	cmd := exec.Command("node", append([]string{script}, args...)...)
	cmd.Dir = cfg.RepoRoot
	cmd.Env = os.Environ()
	out, err := cmd.CombinedOutput()
	return string(out), err
}

func runBin(name string, args []string) error {
	cmd := exec.Command(name, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	cmd.Env = os.Environ()
	err := cmd.Run()
	if err != nil {
		if ee, ok := err.(*exec.Error); ok && ee.Err == exec.ErrNotFound {
			return fmt.Errorf("%s not found on PATH", name)
		}
	}
	return err
}

func runNodeArgs(cfg Config, args []string) error {
	cmd := exec.Command("node", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	cmd.Dir = cfg.RepoRoot
	cmd.Env = os.Environ()
	return cmd.Run()
}

func runBinCapture(name string, args []string) (string, error) {
	cmd := exec.Command(name, args...)
	cmd.Env = os.Environ()
	out, err := cmd.CombinedOutput()
	if err != nil {
		if ee, ok := err.(*exec.Error); ok && ee.Err == exec.ErrNotFound {
			return "", fmt.Errorf("%s not found on PATH", name)
		}
	}
	return string(out), err
}

func which(name string) bool {
	_, err := exec.LookPath(name)
	return err == nil
}

func envSet(keys ...string) bool {
	for _, k := range keys {
		if strings.TrimSpace(os.Getenv(k)) != "" {
			return true
		}
	}
	return false
}
