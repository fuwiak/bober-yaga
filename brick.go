package main

import "slices"

// Visibility controls who sees a brick.
const (
	VisPublic = "public" // visible on profile=public
	VisOwner  = "owner"  // only profile=owner (default)
)

// Brick is one klocek — add a file / register in allBricks().
type Brick struct {
	ID          string
	Title       string
	Description string
	Visibility  string // public | owner
	Aliases     []string
	// DefaultArgs used when TUI runs the brick with Enter
	DefaultArgs []string
	// Help text for `yaga <id>` / `yaga <id> help`
	Help string
	// Run executes the brick. args is everything after the brick name.
	Run func(cfg Config, args []string) error
}

func (b Brick) vis() string {
	if b.Visibility == "" {
		return VisOwner
	}
	return b.Visibility
}

func isVisible(b Brick, cfg Config) bool {
	id := b.ID
	if slices.Contains(cfg.Hidden, id) {
		return false
	}
	if slices.Contains(cfg.Disabled, id) {
		return false
	}
	switch cfg.Profile {
	case "public":
		return b.vis() == VisPublic
	case "custom":
		if len(cfg.Enabled) > 0 {
			return slices.Contains(cfg.Enabled, id) || id == "core"
		}
		return !slices.Contains(cfg.Disabled, id)
	default: // owner
		return true
	}
}

func visibleBricks(cfg Config) []Brick {
	var out []Brick
	for _, b := range allBricks() {
		if isVisible(b, cfg) {
			out = append(out, b)
		}
	}
	return out
}

func findBrick(cfg Config, name string) (Brick, bool) {
	for _, b := range visibleBricks(cfg) {
		if b.ID == name {
			return b, true
		}
		if slices.Contains(b.Aliases, name) {
			return b, true
		}
	}
	return Brick{}, false
}

func allBrickMeta(cfg Config) []BrickMeta {
	var out []BrickMeta
	for _, b := range allBricks() {
		out = append(out, BrickMeta{
			Brick:   b,
			Visible: isVisible(b, cfg),
		})
	}
	return out
}

type BrickMeta struct {
	Brick
	Visible bool
}
