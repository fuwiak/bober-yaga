package main

import (
	"fmt"
	"os"
	"strings"
	"time"
	"unicode"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type tab int

const (
	tabBricks tab = iota
	tabCreds
	tabDoctor
	tabOutput
	tabHelp
)

const tabCount = 5

type model struct {
	cfg           Config
	width, height int
	tab           tab
	cursor        int
	credCursor    int
	bricks        []Brick
	creds         []CredSpec
	flash         string
	output        string
	running       bool
	vp            viewport.Model
	ready         bool
	clock         time.Time

	// credential edit mode
	editing    bool
	editKey    string
	editBuf    string
	editMasked bool
}

type tickMsg time.Time

type runDoneMsg struct {
	title string
	out   string
	err   string
}

func newModel(cfg Config) model {
	bricks := visibleBricks(cfg)
	return model{
		cfg:        cfg,
		tab:        tabBricks,
		bricks:     bricks,
		creds:      allCredSpecs(),
		flash:      fmt.Sprintf("profile=%s · %d bricks", cfg.Profile, len(bricks)),
		vp:         viewport.New(80, 20),
		clock:      time.Now(),
		output:     "Select a brick and press Enter to run its default action.\nOutput appears on Output tab.",
		editMasked: true,
	}
}

func runTUI(cfg Config) {
	m := newModel(cfg)
	p := tea.NewProgram(m, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "yaga: %v\n", err)
		os.Exit(1)
	}
}

func (m model) Init() tea.Cmd {
	return tea.Tick(time.Second, func(t time.Time) tea.Msg { return tickMsg(t) })
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
		m.vp.Width = max(20, m.width)
		m.vp.Height = max(3, m.height-3)
		m.ready = true
		m.refreshVP()
		return m, nil

	case tickMsg:
		m.clock = time.Time(msg)
		return m, tea.Tick(time.Second, func(t time.Time) tea.Msg { return tickMsg(t) })

	case runDoneMsg:
		m.running = false
		if msg.err != "" {
			m.output = msg.title + "\n\n" + msg.out + "\n\nERR: " + msg.err
			m.flash = "failed · " + msg.title
		} else {
			m.output = msg.title + "\n\n" + msg.out
			if strings.TrimSpace(msg.out) == "" {
				m.output += "(no capture — run from shell for live output)"
			}
			m.flash = "ok · " + msg.title
		}
		m.tab = tabOutput
		m.refreshVP()
		return m, nil

	case tea.KeyMsg:
		return m.handleKey(msg)
	}
	return m, nil
}

func (m model) handleKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	if m.editing {
		return m.handleEditKey(msg)
	}

	switch msg.String() {
	case "q", "ctrl+c":
		return m, tea.Quit
	case "tab", "right":
		m.tab = (m.tab + 1) % tabCount
		m.refreshVP()
		return m, nil
	case "shift+tab", "left":
		m.tab = (m.tab + tabCount - 1) % tabCount
		m.refreshVP()
		return m, nil
	case "1":
		m.tab = tabBricks
		m.refreshVP()
		return m, nil
	case "2":
		m.tab = tabCreds
		m.refreshVP()
		return m, nil
	case "3":
		m.tab = tabDoctor
		m.refreshVP()
		return m, nil
	case "4":
		m.tab = tabOutput
		m.refreshVP()
		return m, nil
	case "5", "?":
		m.tab = tabHelp
		m.refreshVP()
		return m, nil
	case "r":
		m.bricks = visibleBricks(m.cfg)
		m.creds = allCredSpecs()
		m.flash = "refreshed"
		if m.cursor >= len(m.bricks) {
			m.cursor = max(0, len(m.bricks)-1)
		}
		m.refreshVP()
		return m, nil
	case "up", "k":
		if m.tab == tabBricks && m.cursor > 0 {
			m.cursor--
			m.refreshVP()
		} else if m.tab == tabCreds && m.credCursor > 0 {
			m.credCursor--
			m.refreshVP()
		} else {
			m.vp.LineUp(1)
		}
		return m, nil
	case "down", "j":
		if m.tab == tabBricks && m.cursor < len(m.bricks)-1 {
			m.cursor++
			m.refreshVP()
		} else if m.tab == tabCreds && m.credCursor < len(m.creds)-1 {
			m.credCursor++
			m.refreshVP()
		} else {
			m.vp.LineDown(1)
		}
		return m, nil
	case "enter", " ":
		if m.tab == tabBricks && !m.running && len(m.bricks) > 0 {
			b := m.bricks[m.cursor]
			m.running = true
			m.flash = "running " + b.ID + "…"
			return m, runBrickCmd(m.cfg, b)
		}
		if m.tab == tabCreds && len(m.creds) > 0 {
			spec := m.creds[m.credCursor]
			m.editing = true
			m.editKey = spec.Key
			m.editBuf = ""
			m.flash = "paste value · Enter save · Esc cancel · Ctrl+V ok"
			m.refreshVP()
		}
		return m, nil
	case "o":
		if m.tab == tabCreds && len(m.creds) > 0 {
			url := m.creds[m.credCursor].UIURL
			if err := openURL(url); err != nil {
				m.flash = "open failed: " + err.Error()
			} else {
				m.flash = "opened UI · " + url
			}
			m.refreshVP()
		}
		return m, nil
	case "d", "x", "backspace":
		if m.tab == tabCreds && len(m.creds) > 0 && msg.String() != "backspace" {
			key := m.creds[m.credCursor].Key
			if err := setCredential(key, ""); err != nil {
				m.flash = "clear failed"
			} else {
				m.flash = "cleared " + key
			}
			m.refreshVP()
		}
		return m, nil
	case "pgup":
		m.vp.ViewUp()
		return m, nil
	case "pgdown":
		m.vp.ViewDown()
		return m, nil
	}
	return m, nil
}

func (m model) handleEditKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc", "ctrl+c":
		m.editing = false
		m.editBuf = ""
		m.flash = "cancelled"
		m.refreshVP()
		return m, nil
	case "enter":
		if err := setCredential(m.editKey, m.editBuf); err != nil {
			m.flash = "save failed: " + err.Error()
		} else {
			m.flash = "saved " + m.editKey + " → " + credentialsPath()
		}
		m.editing = false
		m.editBuf = ""
		m.refreshVP()
		return m, nil
	case "backspace":
		r := []rune(m.editBuf)
		if len(r) > 0 {
			m.editBuf = string(r[:len(r)-1])
		}
		m.refreshVP()
		return m, nil
	case "ctrl+u":
		m.editBuf = ""
		m.refreshVP()
		return m, nil
	default:
		if msg.Type == tea.KeyRunes {
			for _, r := range msg.Runes {
				if unicode.IsPrint(r) {
					m.editBuf += string(r)
				}
			}
			m.refreshVP()
		}
		return m, nil
	}
}

func runBrickCmd(cfg Config, b Brick) tea.Cmd {
	return func() tea.Msg {
		args := b.DefaultArgs
		title := "yaga " + b.ID
		if len(args) > 0 {
			title += " " + strings.Join(args, " ")
		}
		out, err := captureBrick(cfg, b, args)
		msg := runDoneMsg{title: title, out: out}
		if err != nil {
			msg.err = err.Error()
		}
		return msg
	}
}

func captureBrick(cfg Config, b Brick, args []string) (string, error) {
	switch b.ID {
	case "webmaster":
		sub, rest := splitSub(args, "status")
		switch sub {
		case "status":
			return runScriptCapture(cfg, "yandex-status.mjs", rest)
		case "seo", "positions", "rank":
			return runScriptCapture(cfg, "yandex-webmaster-seo.mjs", rest)
		case "selfcheck", "checks", "checklist":
			return runScriptCapture(cfg, "yandex-webmaster-selfcheck.mjs", rest)
		case "microtest", "microdata", "validator":
			return runScriptCapture(cfg, "yandex-webmaster-microtest.mjs", rest)
		case "boost", "index":
			return runScriptCapture(cfg, "yandex-webmaster-boost.mjs", rest)
		case "sitemap", "sitemaps":
			return runScriptCapture(cfg, "yandex-webmaster-sitemap.mjs", rest)
		case "feed":
			return runScriptCapture(cfg, "yandex-webmaster-feed.mjs", rest)
		case "mirrors":
			return runScriptCapture(cfg, "yandex-webmaster-mirrors.mjs", rest)
		case "recrawl", "crawl":
			return runScriptCapture(cfg, "yandex-webmaster-recrawl.mjs", rest)
		}
	case "metrika":
		sub, rest := splitSub(args, "status")
		switch sub {
		case "status":
			return runScriptCapture(cfg, "yandex-metrika-status.mjs", rest)
		case "counter":
			return runScriptCapture(cfg, "yandex-metrika-counter.mjs", rest)
		case "filters", "filter":
			return runScriptCapture(cfg, "yandex-metrika-filters.mjs", rest)
		case "organic":
			return runScriptCapture(cfg, "yandex-metrika-filters.mjs", append([]string{"--organic"}, rest...))
		case "ytm", "tag":
			return runScriptCapture(cfg, "yandex-ytm-status.mjs", rest)
		}
	case "direct":
		return runScriptCapture(cfg, "yandex-direct-campaigns.mjs", []string{"status"})
	case "business":
		sub, rest := splitSub(args, "status")
		if sub == "oauth" {
			return runScriptCapture(cfg, "yandex-business-oauth.mjs", append([]string{"status"}, rest...))
		}
		return runScriptCapture(cfg, "yandex-business.mjs", args)
	case "cloud":
		return runBinCapture("yc", []string{"version"})
	case "core":
		return brickCore().Help, nil
	case "doctor":
		return renderDoctor(cfg), nil
	}
	if b.Help != "" {
		return b.Help + "\n\n(Run from shell: yaga " + b.ID + " …)\n", nil
	}
	return "", fmt.Errorf("no capture handler; run: yaga %s", b.ID)
}

func (m *model) refreshVP() {
	m.vp.SetContent(m.body())
	m.vp.GotoTop()
}

func (m model) View() string {
	if !m.ready {
		return "…"
	}
	return lipgloss.JoinVertical(lipgloss.Left, m.header(), m.tabs(), m.vp.View(), m.footer())
}

func (m model) header() string {
	bg := lipgloss.Color("17")
	if m.cfg.Profile == "owner" {
		bg = lipgloss.Color("54")
	} else if m.cfg.Profile == "public" {
		bg = lipgloss.Color("22")
	}
	style := lipgloss.NewStyle().Background(bg).Foreground(lipgloss.Color("15")).Width(max(20, m.width))
	left := fmt.Sprintf(" YAGA · %s ", strings.ToUpper(m.cfg.Profile))
	right := trunc(m.flash, max(8, m.width-lipgloss.Width(left)-1))
	pad := max(0, m.width-lipgloss.Width(left)-lipgloss.Width(right))
	return style.Render(left + strings.Repeat(" ", pad) + right)
}

func (m model) tabs() string {
	inactive := lipgloss.NewStyle().Foreground(lipgloss.Color("245")).Padding(0, 1)
	active := lipgloss.NewStyle().Background(lipgloss.Color("62")).Foreground(lipgloss.Color("15")).Bold(true).Padding(0, 1)
	labels := []string{"1 Bricks", "2 Creds", "3 Doctor", "4 Output", "5 Help"}
	parts := make([]string, 0, len(labels))
	for i, label := range labels {
		if tab(i) == m.tab {
			parts = append(parts, active.Render(label))
		} else {
			parts = append(parts, inactive.Render(label))
		}
	}
	bar := lipgloss.JoinHorizontal(lipgloss.Top, parts...)
	return bar + strings.Repeat(" ", max(0, m.width-lipgloss.Width(bar)))
}

func (m model) footer() string {
	style := lipgloss.NewStyle().Foreground(lipgloss.Color("245")).Width(max(20, m.width))
	if m.editing {
		return style.Render(" editing credential · Enter save · Esc cancel · Ctrl+U clear ")
	}
	if m.tab == tabCreds {
		return style.Render(" ↑↓ · Enter set · o open UI · d delete · Tab switch · q quit ")
	}
	return style.Render(" Tab/←→ · ↑↓ select · Enter run · r refresh · q quit ")
}

func (m model) body() string {
	switch m.tab {
	case tabBricks:
		return m.bodyBricks()
	case tabCreds:
		return m.bodyCreds()
	case tabDoctor:
		return "\n" + renderDoctor(m.cfg)
	case tabOutput:
		if m.running {
			return "\n  running…\n"
		}
		return "\n" + m.output + "\n"
	case tabHelp:
		return helpTUI(m.cfg)
	}
	return ""
}

func (m model) bodyCreds() string {
	var b strings.Builder
	b.WriteString("\n  Credentials — set secrets + UI where to find them\n")
	b.WriteString("  " + strings.Repeat("─", 52) + "\n")
	fmt.Fprintf(&b, "  file: %s\n\n", credentialsPath())

	if m.editing {
		display := m.editBuf
		if m.editMasked && len(display) > 0 {
			display = maskSecret(display)
		}
		fmt.Fprintf(&b, "  SET %s\n", m.editKey)
		fmt.Fprintf(&b, "  > %s▌\n\n", display)
		fmt.Fprintf(&b, "  Paste token/value and press Enter.\n")
		return b.String()
	}

	for i, spec := range m.creds {
		cursor := "  "
		if i == m.credCursor {
			cursor = "▸ "
		}
		mark := "○"
		val := credValue(spec.Key)
		if val != "" {
			mark = "✓"
		}
		opt := ""
		if spec.Optional {
			opt = " (opt)"
		}
		fmt.Fprintf(&b, "%s%s %-32s %s%s\n", cursor, mark, spec.Key, spec.Title, opt)
		if i == m.credCursor {
			fmt.Fprintf(&b, "               brick: %s\n", spec.Brick)
			fmt.Fprintf(&b, "               value: %s\n", maskSecret(val))
			fmt.Fprintf(&b, "               UI:    %s\n", spec.UIURL)
			fmt.Fprintf(&b, "               how:   %s\n", wrapHint(spec.How, 60))
		}
	}
	b.WriteString("\n  Keys: Enter=set · o=open UI in browser · d=delete saved value\n")
	return b.String()
}

func wrapHint(s string, width int) string {
	if len(s) <= width {
		return s
	}
	var out strings.Builder
	words := strings.Fields(s)
	line := ""
	first := true
	for _, w := range words {
		if line == "" {
			line = w
			continue
		}
		if len(line)+1+len(w) > width {
			if !first {
				out.WriteString("\n               ")
			}
			out.WriteString(line)
			first = false
			line = w
		} else {
			line += " " + w
		}
	}
	if line != "" {
		if !first {
			out.WriteString("\n               ")
		}
		out.WriteString(line)
	}
	return out.String()
}

func (m model) bodyBricks() string {
	var b strings.Builder
	b.WriteString(fmt.Sprintf("\n  Bricks for profile=%s (Enter = default action)\n", m.cfg.Profile))
	b.WriteString("  " + strings.Repeat("─", 48) + "\n\n")
	if len(m.bricks) == 0 {
		b.WriteString("  (none visible — try yaga profile owner)\n")
		return b.String()
	}
	for i, br := range m.bricks {
		cursor := "  "
		if i == m.cursor {
			cursor = "▸ "
		}
		b.WriteString(fmt.Sprintf("%s%-12s [%s]  %s\n", cursor, br.ID, br.vis(), br.Title))
		if br.Description != "" {
			b.WriteString(fmt.Sprintf("               %s\n", br.Description))
		}
	}
	if m.cursor >= 0 && m.cursor < len(m.bricks) {
		br := m.bricks[m.cursor]
		b.WriteString("\n  " + strings.Repeat("─", 48) + "\n")
		b.WriteString(fmt.Sprintf("  default: yaga %s %s\n", br.ID, strings.Join(br.DefaultArgs, " ")))
	}
	return b.String()
}

func helpTUI(cfg Config) string {
	return fmt.Sprintf(`
  yaga — Go + Bubble Tea (same stack as bober)

  Tabs
    1 Bricks   pick klocek, Enter runs default
    2 Creds    add/edit credentials + UI URLs
    3 Doctor   tokens / binaries status
    4 Output   last captured run
    5 Help

  Credentials
    Enter  set value (saved to ~/.config/yaga/credentials.env)
    o      open Yandex UI page in browser
    d      delete saved value

  Shell
    yaga
    yaga credentials          list + UI links
    yaga credentials set KEY
    yaga webmaster status
    yaga webmaster seo
    yaga webmaster recrawl <url>
    yaga profile public

  Config: %s
  Creds:  %s
  Repo:   %s
`, cfg.Path, credentialsPath(), cfg.RepoRoot)
}

func trunc(s string, n int) string {
	r := []rune(s)
	if n <= 0 {
		return ""
	}
	if len(r) <= n {
		return s
	}
	if n <= 1 {
		return string(r[:n])
	}
	return string(r[:n-1]) + "…"
}
