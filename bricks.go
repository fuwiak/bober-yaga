package main

import (
	"fmt"
	"os"
	"strings"
)

// allBricks returns every klocek. Add/remove entries here (or split into bricks_*.go).
func allBricks() []Brick {
	return []Brick{
		brickCore(),
		brickWebmaster(),
		brickMetrika(),
		brickCloud(),
		brickDisk(),
		brickWordstat(),
		brickDirect(),
		brickBusiness(),
		brickMail(),
		brickSearch(),
		brickGPT(),
		brickAIStudio(),
	}
}

func brickCore() Brick {
	return Brick{
		ID: "core", Title: "Core", Visibility: VisPublic,
		Description: "Architecture / how bricks work",
		Help: `yaga core — brick model

Each brick is a Go klocek in bricks (allBricks).
visibility: public | owner
Profiles: owner (all) · public (public only) · custom (enable/disable)

  yaga bricks list
  yaga bricks hide direct
  yaga profile public
`,
		Run: func(cfg Config, args []string) error {
			fmt.Print(brickCore().Help)
			return nil
		},
	}
}

func brickWebmaster() Brick {
	return Brick{
		ID: "webmaster", Title: "Yandex Webmaster", Visibility: VisPublic,
		Aliases: []string{"wm"}, Description: "status / seo / selfcheck / microtest / boost / sitemap / feed / mirrors / recrawl / oauth",
		DefaultArgs: []string{"status"},
		Help: `yaga webmaster status|seo|selfcheck|microtest|boost|sitemap|feed|mirrors|repair|recrawl|oauth

  status     снимок Вебмастер + Метрика
  seo        чеклист позиций (ИКС, диагностика, индекс, запросы)
  selfcheck  21 «Самостоятельные проверки» + группы целевых запросов
  microtest  валидатор микроразметки: список URL + JSON-LD/OG снимок
  boost      переобход важных URL + чеклист региона/фида
  sitemap    добавить /sitemap.xml в очередь Вебмастера
  feed       загрузка performers feed
  mirrors    зеркала
  repair     перезагрузка фида
  recrawl    переобход URL (или --quota)
  oauth      ClientID/secret → access token (один раз в браузере)
`,
		Run: func(cfg Config, args []string) error {
			sub, rest := splitSub(args, "status")
			switch sub {
			case "help", "--help":
				fmt.Println(brickWebmaster().Help)
				return nil
			case "status":
				return runScript(cfg, "yandex-status.mjs", rest)
			case "seo", "positions", "rank":
				return runScript(cfg, "yandex-webmaster-seo.mjs", rest)
			case "selfcheck", "checks", "checklist":
				return runScript(cfg, "yandex-webmaster-selfcheck.mjs", rest)
			case "microtest", "microdata", "validator":
				return runScript(cfg, "yandex-webmaster-microtest.mjs", rest)
			case "boost", "index":
				return runScript(cfg, "yandex-webmaster-boost.mjs", rest)
			case "sitemap", "sitemaps":
				return runScript(cfg, "yandex-webmaster-sitemap.mjs", rest)
			case "feed":
				return runScript(cfg, "yandex-webmaster-feed.mjs", rest)
			case "mirrors":
				return runScript(cfg, "yandex-webmaster-mirrors.mjs", rest)
			case "repair":
				return runScript(cfg, "yandex-feed-repair.mjs", rest)
			case "recrawl", "crawl":
				return runScript(cfg, "yandex-webmaster-recrawl.mjs", rest)
			case "oauth", "login", "auth":
				return runScript(cfg, "yandex-webmaster-oauth.mjs", rest)
			default:
				return runScript(cfg, "yandex-status.mjs", args)
			}
		},
	}
}

func brickMetrika() Brick {
	return Brick{
		ID: "metrika", Title: "Yandex Metrika", Visibility: VisPublic,
		Aliases: []string{"metrica", "ym"}, Description: "status / counter / filters / ytm",
		DefaultArgs: []string{"status"},
		Help: `yaga metrika status|counter|filters|organic|ytm

  status    снимок счётчика + цели
  counter   найти/создать счётчик
  filters   фильтры + «не учитывать мои» + органика
  organic   только органика за период (--from/--to)
  ytm       статус Яндекс Тег Менеджера
`,
		Run: func(cfg Config, args []string) error {
			sub, rest := splitSub(args, "status")
			switch sub {
			case "help", "--help":
				fmt.Println(brickMetrika().Help)
				return nil
			case "status":
				return runScript(cfg, "yandex-metrika-status.mjs", rest)
			case "counter":
				return runScript(cfg, "yandex-metrika-counter.mjs", rest)
			case "filters", "filter":
				return runScript(cfg, "yandex-metrika-filters.mjs", rest)
			case "organic":
				return runScript(cfg, "yandex-metrika-filters.mjs", append([]string{"--organic"}, rest...))
			case "ytm", "tag":
				return runScript(cfg, "yandex-ytm-status.mjs", rest)
			default:
				return runScript(cfg, "yandex-metrika-status.mjs", args)
			}
		},
	}
}

func brickCloud() Brick {
	return Brick{
		ID: "cloud", Title: "Yandex Cloud (yc)", Visibility: VisPublic,
		Aliases: []string{"yc"}, Description: "passthrough → yc",
		DefaultArgs: []string{"--help"},
		Help: `yaga cloud <yc-args…>   (requires yc on PATH)`,
		Run: func(cfg Config, args []string) error {
			if len(args) == 0 || args[0] == "help" || args[0] == "--help" {
				fmt.Println(brickCloud().Help)
				if len(args) == 0 {
					return nil
				}
			}
			if len(args) > 0 && args[0] == "--" {
				args = args[1:]
			}
			return runBin("yc", args)
		},
	}
}

func brickDisk() Brick {
	return Brick{
		ID: "disk", Title: "Yandex Disk", Visibility: VisPublic,
		Aliases: []string{"ydisk"}, Description: "passthrough → yandex-disk",
		DefaultArgs: []string{"status"},
		Help: `yaga disk status|start|stop|sync`,
		Run: func(cfg Config, args []string) error {
			if len(args) == 0 || args[0] == "help" || args[0] == "--help" {
				fmt.Println(brickDisk().Help)
				if len(args) == 0 {
					return nil
				}
			}
			if len(args) > 0 && args[0] == "--" {
				args = args[1:]
			}
			return runBin("yandex-disk", args)
		},
	}
}

func brickWordstat() Brick {
	return Brick{
		ID: "wordstat", Title: "Wordstat", Visibility: VisOwner,
		Aliases: []string{"ws"}, Description: "keyword research pipeline",
		DefaultArgs: []string{"run"},
		Help: `yaga wordstat run|search "<phrase>"`,
		Run: func(cfg Config, args []string) error {
			sub, rest := splitSub(args, "run")
			switch sub {
			case "help", "--help":
				fmt.Println(brickWordstat().Help)
				return nil
			case "run", "pipeline":
				return runScript(cfg, "yandex-wordstat.mjs", rest)
			case "search":
				phrase := strings.TrimSpace(strings.Join(rest, " "))
				if phrase == "" {
					return fmt.Errorf(`usage: yaga wordstat search "…"`)
				}
				os.Setenv("WORDSTAT_PHRASE", phrase)
				return runScript(cfg, "yandex-wordstat.mjs", nil)
			default:
				return runScript(cfg, "yandex-wordstat.mjs", args)
			}
		},
	}
}

func brickDirect() Brick {
	return Brick{
		ID: "direct", Title: "Yandex Direct", Visibility: VisOwner,
		Aliases: []string{"dir"}, Description: "campaigns + oauth (private)",
		DefaultArgs: []string{"campaigns", "status"},
		Help: `yaga direct campaigns […]
yaga direct oauth […]`,
		Run: func(cfg Config, args []string) error {
			sub, rest := splitSub(args, "campaigns")
			switch sub {
			case "help", "--help":
				fmt.Println(brickDirect().Help)
				return nil
			case "campaigns", "campaign":
				if len(rest) == 0 {
					rest = []string{"status"}
				}
				return runScript(cfg, "yandex-direct-campaigns.mjs", rest)
			case "oauth":
				return runScript(cfg, "yandex-direct-oauth.mjs", rest)
			default:
				return runScript(cfg, "yandex-direct-campaigns.mjs", args)
			}
		},
	}
}

func brickBusiness() Brick {
	return Brick{
		ID: "business", Title: "Yandex Business", Visibility: VisOwner,
		Aliases: []string{"biz", "sprav", "geoadv"}, Description: "Partner API: card search/create + ads",
		DefaultArgs: []string{"status"},
		Help: `yaga business status|search|card|campaigns|rubrics|regions|owner|limits|oauth

  status              ping API + companyId
  search <text>       поиск организации
  card                локальный профиль + API search
  card create         create-company из data/yandex-business-company.json
  campaigns           список/create/price/launch РК (агентство)
  rubrics|regions     подсказки
  limits              что API не умеет (edit card / stories)
  oauth               bootstrap токена

Docs: https://yandex.ru/dev/business-api/doc/ref/index.html
UI:   https://yandex.ru/sprav/companies
`,
		Run: func(cfg Config, args []string) error {
			sub, rest := splitSub(args, "status")
			switch sub {
			case "help", "--help":
				fmt.Println(brickBusiness().Help)
				return nil
			case "oauth":
				return runScript(cfg, "yandex-business-oauth.mjs", rest)
			default:
				return runScript(cfg, "yandex-business.mjs", args)
			}
		},
	}
}

func brickMail() Brick {
	return Brick{
		ID: "mail", Title: "Yandex 360 Mail", Visibility: VisOwner,
		Aliases: []string{"mail360", "y360"}, Description: "corporate mailbox",
		DefaultArgs: []string{},
		Help: `yaga mail [mailbox args…]`,
		Run: func(cfg Config, args []string) error {
			if len(args) > 0 && (args[0] == "help" || args[0] == "--help") {
				fmt.Println(brickMail().Help)
				return nil
			}
			if len(args) > 0 && (args[0] == "mailbox" || args[0] == "box") {
				return runScript(cfg, "yandex360-mailbox.mjs", args[1:])
			}
			return runScript(cfg, "yandex360-mailbox.mjs", args)
		},
	}
}

func brickSearch() Brick {
	return Brick{
		ID: "search", Title: "Yandex Search API", Visibility: VisOwner,
		Description: "Cloud Search via legacy node helper",
		DefaultArgs: []string{"help"},
		Help:        `yaga search web "<query>"  (needs YANDEX_SEARCH_API_KEY + folder)`,
		Run: func(cfg Config, args []string) error {
			return runLegacyNode(cfg, "search", args)
		},
	}
}

func brickGPT() Brick {
	return Brick{
		ID: "gpt", Title: "YandexGPT", Visibility: VisOwner,
		Aliases: []string{"yandexgpt", "llm"}, Description: "foundation models chat",
		DefaultArgs: []string{"help"},
		Help:        `yaga gpt chat <prompt…>`,
		Run: func(cfg Config, args []string) error {
			return runLegacyNode(cfg, "gpt", args)
		},
	}
}

func brickAIStudio() Brick {
	return Brick{
		ID: "aistudio", Title: "Yandex AI Studio", Visibility: VisOwner,
		Aliases: []string{"studio", "rag"}, Description: "passthrough → yandex-ai-studio",
		DefaultArgs: []string{"--help"},
		Help:        `yaga aistudio -- <yandex-ai-studio args>`,
		Run: func(cfg Config, args []string) error {
			if len(args) == 0 || args[0] == "help" || args[0] == "--help" {
				fmt.Println(brickAIStudio().Help)
				if len(args) == 0 {
					return nil
				}
			}
			if len(args) > 0 && args[0] == "--" {
				args = args[1:]
			}
			return runBin("yandex-ai-studio", args)
		},
	}
}

func splitSub(args []string, def string) (string, []string) {
	if len(args) == 0 {
		return def, nil
	}
	return args[0], args[1:]
}

func runLegacyNode(cfg Config, id string, args []string) error {
	legacy := cfg.RepoRoot + "/cli/yaga/node/bin/yaga.js"
	if !fileExists(legacy) {
		return fmt.Errorf("legacy node runner missing at %s", legacy)
	}
	return runNodeArgs(cfg, append([]string{legacy, id}, args...))
}
