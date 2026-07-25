package main

import (
	"bufio"
	"fmt"
	"os"
	"slices"
	"strings"
)

func main() {
	applyCredentials()
	cfg := loadConfig()
	args := os.Args[1:]

	if len(args) == 0 {
		runTUI(cfg)
		return
	}

	switch args[0] {
	case "tui", "ui", "dashboard":
		runTUI(cfg)
	case "help", "-h", "--help":
		printHelp(cfg)
	case "version", "-V", "--version":
		fmt.Println("yaga 0.2.0 (go+tui)")
	case "bricks", "plugins", "blocks":
		cmdBricks(cfg, args[1:])
	case "profile":
		cmdProfile(cfg, args[1:])
	case "credentials", "creds", "cred":
		cmdCredentials(args[1:])
	case "doctor":
		fmt.Print(renderDoctor(cfg))
	default:
		b, ok := findBrick(cfg, args[0])
		if !ok {
			fmt.Fprintf(os.Stderr, "yaga: unknown '%s' (profile=%s)\nTry: yaga bricks · yaga help\n", args[0], cfg.Profile)
			os.Exit(2)
		}
		if err := b.Run(cfg, args[1:]); err != nil {
			fmt.Fprintf(os.Stderr, "yaga %s: %v\n", b.ID, err)
			os.Exit(1)
		}
	}
}

func printHelp(cfg Config) {
	bricks := visibleBricks(cfg)
	fmt.Printf(`yaga — modular Yandex CLI (Go + Bubble Tea)

Profile: %s  ·  config: %s  ·  repo: %s

Usage:
  yaga                     TUI dashboard
  yaga <brick> [args…]     run brick
  yaga bricks              list / enable / hide
  yaga credentials         list secrets + UI URLs
  yaga profile [name]      owner | public | custom
  yaga doctor
  yaga help

Bricks:
`, cfg.Profile, cfg.Path, cfg.RepoRoot)
	for _, b := range bricks {
		fmt.Printf("  %-12s [%s]  %s\n", b.ID, padVis(b.vis()), b.Title)
		if b.Description != "" {
			fmt.Printf("               %s\n", b.Description)
		}
	}
	fmt.Println(`
TUI: Tab switch · ↑↓ bricks · Enter run default · r refresh doctor · q quit
Hide for others: yaga profile public   or   YAGA_PROFILE=public yaga`)
}

func padVis(v string) string {
	if v == VisPublic {
		return "public"
	}
	return "owner "
}

func cmdBricks(cfg Config, args []string) {
	sub := "list"
	if len(args) > 0 {
		sub = args[0]
	}
	switch sub {
	case "list", "ls":
		fmt.Printf("profile=%s\n\n", cfg.Profile)
		for _, m := range allBrickMeta(cfg) {
			mark := "○"
			if m.Visible {
				mark = "●"
			}
			fmt.Printf("  %s %-12s vis=%-6s  %s\n", mark, m.ID, m.vis(), m.Title)
		}
		fmt.Println("\n● visible  ○ hidden")
		fmt.Println("yaga bricks enable|disable|hide|unhide <id>")
	case "enable", "disable", "hide", "unhide":
		if len(args) < 2 {
			fmt.Fprintf(os.Stderr, "usage: yaga bricks %s <id>\n", sub)
			os.Exit(2)
		}
		id := args[1]
		switch sub {
		case "enable":
			cfg.Profile = "custom"
			cfg.Disabled = slices.DeleteFunc(cfg.Disabled, func(s string) bool { return s == id })
			if !slices.Contains(cfg.Enabled, id) {
				cfg.Enabled = append(cfg.Enabled, id)
			}
		case "disable":
			cfg.Profile = "custom"
			cfg.Enabled = slices.DeleteFunc(cfg.Enabled, func(s string) bool { return s == id })
			if !slices.Contains(cfg.Disabled, id) {
				cfg.Disabled = append(cfg.Disabled, id)
			}
		case "hide":
			if !slices.Contains(cfg.Hidden, id) {
				cfg.Hidden = append(cfg.Hidden, id)
			}
		case "unhide":
			cfg.Hidden = slices.DeleteFunc(cfg.Hidden, func(s string) bool { return s == id })
		}
		if err := saveConfig(cfg); err != nil {
			fmt.Fprintf(os.Stderr, "save: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("saved %s %s → %s\n", sub, id, cfg.Path)
	default:
		fmt.Fprintf(os.Stderr, "usage: yaga bricks [list|enable|disable|hide|unhide]\n")
		os.Exit(2)
	}
}

func cmdProfile(cfg Config, args []string) {
	if len(args) == 0 {
		fmt.Printf("current: %s\n", cfg.Profile)
		fmt.Println("  owner   all bricks")
		fmt.Println("  public  only visibility=public")
		fmt.Println("  custom  enable/disable lists")
		return
	}
	name := args[0]
	if name != "owner" && name != "public" && name != "custom" {
		fmt.Fprintf(os.Stderr, "unknown profile %q\n", name)
		os.Exit(2)
	}
	cfg.Profile = name
	if err := saveConfig(cfg); err != nil {
		fmt.Fprintf(os.Stderr, "save: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("profile=%s → %s\n", name, cfg.Path)
}

func cmdCredentials(args []string) {
	sub := "list"
	if len(args) > 0 {
		sub = args[0]
	}
	switch sub {
	case "list", "ls":
		fmt.Print(renderCredentialsHelp())
	case "set":
		if len(args) < 2 {
			fmt.Fprintf(os.Stderr, "usage: yaga credentials set KEY [value]\n")
			os.Exit(2)
		}
		key := args[1]
		var value string
		if len(args) >= 3 {
			value = strings.Join(args[2:], " ")
		} else {
			fmt.Fprintf(os.Stderr, "value for %s: ", key)
			sc := bufio.NewScanner(os.Stdin)
			if sc.Scan() {
				value = sc.Text()
			}
		}
		if err := setCredential(key, value); err != nil {
			fmt.Fprintf(os.Stderr, "save: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("saved %s → %s\n", key, credentialsPath())
	case "unset", "rm", "delete":
		if len(args) < 2 {
			fmt.Fprintf(os.Stderr, "usage: yaga credentials unset KEY\n")
			os.Exit(2)
		}
		if err := setCredential(args[1], ""); err != nil {
			fmt.Fprintf(os.Stderr, "save: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("cleared %s\n", args[1])
	case "path":
		fmt.Println(credentialsPath())
	case "open":
		key := ""
		if len(args) >= 2 {
			key = args[1]
		}
		for _, s := range allCredSpecs() {
			if key == "" || s.Key == key {
				fmt.Printf("%s\n  %s\n", s.Key, s.UIURL)
				if key != "" {
					_ = openURL(s.UIURL)
					return
				}
			}
		}
	default:
		fmt.Fprintf(os.Stderr, "usage: yaga credentials [list|set|unset|path|open]\n")
		os.Exit(2)
	}
}

func renderDoctor(cfg Config) string {
	var b strings.Builder
	fmt.Fprintf(&b, "yaga doctor · profile=%s · bricks=%d\n\n", cfg.Profile, len(visibleBricks(cfg)))
	checks := []struct {
		label string
		ok    bool
	}{
		{"YANDEX_OAUTH_TOKEN / METRIKA", envSet("YANDEX_OAUTH_TOKEN", "YANDEX_METRIKA_OAUTH_TOKEN")},
		{"YANDEX_WEBMASTER_OAUTH_TOKEN", envSet("YANDEX_WEBMASTER_OAUTH_TOKEN")},
		{"YANDEX_DIRECT tokens", envSet("YANDEX_DIRECT_OAUTH_TOKEN", "YANDEX_DIRECT_REFRESH_TOKEN")},
		{"YANDEX_DIRECT_CLIENT_ID", envSet("YANDEX_DIRECT_CLIENT_ID")},
		{"YANDEX_SEARCH_API_KEY", envSet("YANDEX_SEARCH_API_KEY")},
		{"yc on PATH", which("yc")},
		{"yandex-disk on PATH", which("yandex-disk")},
		{"node on PATH", which("node")},
		{"repo scripts/", fileExists(scriptsDir(cfg))},
	}
	for _, c := range checks {
		mark := "○"
		if c.ok {
			mark = "✓"
		}
		fmt.Fprintf(&b, "  %s %s\n", mark, c.label)
	}
	fmt.Fprintf(&b, "\nVisible: ")
	ids := visibleBricks(cfg)
	for i, br := range ids {
		if i > 0 {
			b.WriteString(", ")
		}
		b.WriteString(br.ID)
	}
	b.WriteString("\n")
	return b.String()
}
