# yaga

**yaga** (Яга) — модульный CLI для сервисов Яндекса (Go + Charm Bubble Tea).

> **Неофициальный порт / клиент.** Это сторонний инструмент, **не** продукт Яндекса и **не** связан с ООО «Яндекс» официально. Используйте на свой страх и риск; API и условия сервисов Яндекса могут меняться.

Каждый сервис — **brick** (кирпичик) в `bricks.go`. Профили скрывают private-функции.

Репозиторий: [`github.com/fuwiak/yaga`](https://github.com/fuwiak/yaga).  
Brick’и, которые вызывают `scripts/*.mjs`, ожидают checkout **bober-ai** через `YAGA_REPO`.

## Установка

```bash
git clone https://github.com/fuwiak/yaga.git
cd yaga
go build -o yaga .
./install.sh            # ~/bin/yaga → ./run
./.githooks/install.sh  # хуки стиля коммитов (для контрибьюторов)
```

## Использование

```bash
yaga                         # TUI
yaga webmaster status
yaga webmaster oauth
yaga webmaster seo
yaga webmaster sitemap       # добавить /sitemap.xml в очередь Вебмастера
yaga webmaster microtest
yaga webmaster boost
yaga metrika status
yaga business status         # Яндекс Бизнес Partner API
yaga direct campaigns status
yaga bricks
yaga profile public
yaga doctor
yaga credentials
```

Если нужны скрипты из bober-ai:

```bash
export YAGA_REPO=/path/to/bober-ai
```

### TUI

| Вкладка | |
|---------|--|
| 1 Bricks | список brick’ов, Enter = default |
| 2 Creds | секреты + ссылки на UI Яндекса |
| 3 Doctor | токены / бинарники |
| 4 Output | результат |
| 5 Help | |

Секреты: `~/.config/yaga/credentials.env` (не коммитить).

## Профили

| Профиль | |
|---------|--|
| `owner` | всё (для владельца) |
| `public` | только `visibility: public` |
| `custom` | enable / disable / hide |

Конфиг: `~/.config/yaga/config.json`

## Стиль коммитов

Префиксы в духе pandas: `ENH:`, `BUG:`, `CI:`, …  
В сообщениях коммитов **не** упоминать Cursor / Codex / Copilot.

Подробнее: [CONTRIBUTING.md](CONTRIBUTING.md).

## CI / CD

- **CI** — `go vet`, сборка, тесты, smoke на push/PR; в PR проверяется префикс `TYPE:` у коммитов
- **Release** — тег `v*` → мультиархивные бинарники в GitHub Releases

```bash
git tag v0.2.1
git push origin v0.2.1
```

## Стек

- Go 1.22+
- bubbletea + lipgloss + bubbles

## Отказ от ответственности

Яндекс® и названия сервисов (Вебмастер, Метрика, Директ и др.) — товарные знаки соответствующих правообладателей. **yaga** — независимый неофициальный клиент для удобной работы с публичными API и кабинетами, без гарантий совместимости и без официальной поддержки со стороны Яндекса.
