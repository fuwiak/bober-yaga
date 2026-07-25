package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

// CredSpec describes one secret + where to get it in Yandex UI.
type CredSpec struct {
	Key         string
	Title       string
	Brick       string // related brick id
	UIURL       string // primary console URL
	How         string // short how-to
	Optional    bool
}

func credentialsPath() string {
	if p := os.Getenv("YAGA_CREDENTIALS"); p != "" {
		return p
	}
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".config", "yaga", "credentials.env")
}

func allCredSpecs() []CredSpec {
	return []CredSpec{
		{
			Key:   "YANDEX_OAUTH_TOKEN",
			Title: "OAuth access token (общий)",
			Brick: "metrika",
			UIURL: "https://oauth.yandex.ru/",
			How:   "Не ClientID/secret. Один раз: yaga webmaster oauth → получишь y0_… token",
		},
		{
			Key:   "YANDEX_METRIKA_OAUTH_TOKEN",
			Title: "OAuth Метрика (если отдельный)",
			Brick: "metrika",
			UIURL: "https://metrika.yandex.ru/",
			How:   "Кабинет Метрики → настройки / OAuth-приложение с правами metrika:read. Токен: https://oauth.yandex.ru/",
			Optional: true,
		},
		{
			Key:   "NEXT_PUBLIC_YANDEX_METRIKA_ID",
			Title: "ID счётчика Метрики",
			Brick: "metrika",
			UIURL: "https://metrika.yandex.ru/list",
			How:   "Список счётчиков → числовой ID в URL или в настройках счётчика",
		},
		{
			Key:   "YANDEX_WEBMASTER_OAUTH_TOKEN",
			Title: "OAuth access token (Вебмастер)",
			Brick: "webmaster",
			UIURL: "https://oauth.yandex.ru/",
			How:   "yaga webmaster oauth · ClientID/secret сами по себе API не открывают",
		},
		{
			Key:   "YANDEX_WEBMASTER_CLIENT_ID",
			Title: "Client ID приложения (Вебмастер)",
			Brick: "webmaster",
			UIURL: "https://oauth.yandex.ru/client/",
			How:   "oauth.yandex.ru → Мои приложения → ClientID (это не token)",
			Optional: true,
		},
		{
			Key:   "YANDEX_WEBMASTER_CLIENT_SECRET",
			Title: "Client secret приложения (Вебмастер)",
			Brick: "webmaster",
			UIURL: "https://oauth.yandex.ru/client/",
			How:   "То же приложение → Client secret · нужен для exchange code→token",
			Optional: true,
		},
		{
			Key:   "YANDEX_WEBMASTER_REFRESH_TOKEN",
			Title: "Refresh token (Вебмастер)",
			Brick: "webmaster",
			UIURL: "https://oauth.yandex.ru/",
			How:   "Пишется автоматически после yaga webmaster oauth exchange",
			Optional: true,
		},
		{
			Key:   "YANDEX_DIRECT_CLIENT_ID",
			Title: "Direct Client ID",
			Brick: "direct",
			UIURL: "https://oauth.yandex.ru/client/new",
			How:   "Новое приложение с правами Директа. Кабинет: https://direct.yandex.ru/",
		},
		{
			Key:   "YANDEX_DIRECT_CLIENT_SECRET",
			Title: "Direct Client Secret",
			Brick: "direct",
			UIURL: "https://oauth.yandex.ru/client/",
			How:   "То же приложение → Client Secret (показать)",
		},
		{
			Key:   "YANDEX_DIRECT_OAUTH_TOKEN",
			Title: "Direct access token",
			Brick: "direct",
			UIURL: "https://oauth.yandex.ru/",
			How:   "yaga direct oauth · или authorize URL с client_id. Кабинет: https://direct.yandex.ru/",
		},
		{
			Key:   "YANDEX_DIRECT_REFRESH_TOKEN",
			Title: "Direct refresh token",
			Brick: "direct",
			UIURL: "https://oauth.yandex.ru/",
			How:   "Получается вместе с access token (response_type=code). Нужен для auto-refresh",
			Optional: true,
		},
		{
			Key:   "YANDEX_AUDIENCE_OAUTH_TOKEN",
			Title: "Audience access token",
			Brick: "audience",
			UIURL: "https://audience.yandex.com/",
			How:   "npm run yandex:audience:oauth — приложение «direct wordstat» + scopes audience:read/write",
		},
		{
			Key:   "YANDEX_AUDIENCE_REFRESH_TOKEN",
			Title: "Audience refresh token",
			Brick: "audience",
			UIURL: "https://oauth.yandex.ru/",
			How:   "Пишется после yandex:audience:oauth exchange",
			Optional: true,
		},
		{
			Key:   "YANDEX_SEARCH_API_KEY",
			Title: "Cloud Search API key",
			Brick: "search",
			UIURL: "https://console.yandex.cloud/folders",
			How:   "Yandex Cloud Console → каталог → Service accounts / API keys → создать ключ. Docs: https://yandex.cloud/ru/docs/search-api/",
		},
		{
			Key:   "YANDEX_SEARCH_FOLDER_ID",
			Title: "Cloud Folder ID (Search)",
			Brick: "search",
			UIURL: "https://console.yandex.cloud/folders",
			How:   "Console → папка → скопировать folder id из URL или обзора",
		},
		{
			Key:   "YANDEX_GPT_API_KEY",
			Title: "YandexGPT API key",
			Brick: "gpt",
			UIURL: "https://console.yandex.cloud/folders",
			How:   "Cloud Console → сервисный аккаунт → API-ключ. Foundation Models: https://console.yandex.cloud/ → YandexGPT",
			Optional: true,
		},
		{
			Key:   "YANDEX_GPT_FOLDER_ID",
			Title: "Cloud Folder ID (GPT)",
			Brick: "gpt",
			UIURL: "https://console.yandex.cloud/folders",
			How:   "Тот же folder id, что для Search, или отдельный",
			Optional: true,
		},
		{
			Key:   "YANDEX_IAM_TOKEN",
			Title: "IAM token (краткоживущий)",
			Brick: "cloud",
			UIURL: "https://console.yandex.cloud/",
			How:   "yc iam create-token · или Console → CLI. Альтернатива API-ключу",
			Optional: true,
		},
		{
			Key:   "YANDEX_360_OAUTH_TOKEN",
			Title: "Yandex 360 OAuth",
			Brick: "mail",
			UIURL: "https://admin.yandex.ru/",
			How:   "Яндекс 360 для бизнеса → API. OAuth: https://yandex.ru/dev/api360/doc/ru/access",
			Optional: true,
		},
	}
}

func loadCredentialsFile() map[string]string {
	out := map[string]string{}
	path := credentialsPath()
	f, err := os.Open(path)
	if err != nil {
		return out
	}
	defer f.Close()
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		line := strings.TrimSpace(sc.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		k, v, ok := strings.Cut(line, "=")
		if !ok {
			continue
		}
		k = strings.TrimSpace(k)
		v = strings.TrimSpace(v)
		v = strings.Trim(v, `"'`)
		if k != "" {
			out[k] = v
		}
	}
	return out
}

func saveCredentialsFile(vals map[string]string) error {
	path := credentialsPath()
	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		return err
	}
	var b strings.Builder
	b.WriteString("# yaga credentials — do not commit\n")
	b.WriteString("# UI links: see yaga TUI → Credentials tab\n\n")
	for _, spec := range allCredSpecs() {
		if v, ok := vals[spec.Key]; ok && v != "" {
			b.WriteString(spec.Key)
			b.WriteByte('=')
			b.WriteString(shellQuote(v))
			b.WriteByte('\n')
		}
	}
	// keep unknown keys
	known := map[string]bool{}
	for _, s := range allCredSpecs() {
		known[s.Key] = true
	}
	for k, v := range vals {
		if known[k] || v == "" {
			continue
		}
		b.WriteString(k)
		b.WriteByte('=')
		b.WriteString(shellQuote(v))
		b.WriteByte('\n')
	}
	return os.WriteFile(path, []byte(b.String()), 0o600)
}

func shellQuote(s string) string {
	if strings.ContainsAny(s, " \t#\"'\\$") {
		return strconvQuote(s)
	}
	return s
}

func strconvQuote(s string) string {
	return `"` + strings.ReplaceAll(strings.ReplaceAll(s, `\`, `\\`), `"`, `\"`) + `"`
}

// applyCredentials loads file into process env (file wins only if env empty).
func applyCredentials() {
	vals := loadCredentialsFile()
	for k, v := range vals {
		if strings.TrimSpace(os.Getenv(k)) == "" && v != "" {
			_ = os.Setenv(k, v)
		}
	}
}

func credValue(key string) string {
	if v := strings.TrimSpace(os.Getenv(key)); v != "" {
		return v
	}
	return loadCredentialsFile()[key]
}

func credIsSet(key string) bool {
	return credValue(key) != ""
}

func maskSecret(v string) string {
	if v == "" {
		return "(empty)"
	}
	r := []rune(v)
	if len(r) <= 8 {
		return strings.Repeat("•", len(r))
	}
	return string(r[:4]) + strings.Repeat("•", 6) + string(r[len(r)-4:])
}

func openURL(url string) error {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "darwin":
		cmd = exec.Command("open", url)
	case "linux":
		cmd = exec.Command("xdg-open", url)
	default:
		cmd = exec.Command("open", url)
	}
	return cmd.Start()
}

func setCredential(key, value string) error {
	vals := loadCredentialsFile()
	value = strings.TrimSpace(value)
	if value == "" {
		delete(vals, key)
	} else {
		vals[key] = value
	}
	if err := saveCredentialsFile(vals); err != nil {
		return err
	}
	if value == "" {
		_ = os.Unsetenv(key)
	} else {
		_ = os.Setenv(key, value)
	}
	return nil
}

func renderCredentialsHelp() string {
	var b strings.Builder
	b.WriteString("Credentials — где взять в UI Яндекса\n\n")
	for _, s := range allCredSpecs() {
		set := "○"
		if credIsSet(s.Key) {
			set = "✓"
		}
		fmt.Fprintf(&b, "%s %s\n", set, s.Key)
		fmt.Fprintf(&b, "   %s\n", s.Title)
		fmt.Fprintf(&b, "   UI: %s\n", s.UIURL)
		fmt.Fprintf(&b, "   %s\n\n", s.How)
	}
	fmt.Fprintf(&b, "File: %s\n", credentialsPath())
	return b.String()
}
