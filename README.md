# Elite Ball Knowledge

A terminal app for the FIFA World Cup 2026 — live scores, fixtures, standings, and top scorers, refreshed every 45 seconds. Built in Go.

## Install

Requires Go 1.23+.

```bash
git clone https://github.com/charliehustlr1792/fifawc26.git
cd fifawc26
go build -o fifawc26 ./cmd/fifawc26
```

## Setup

Get a free API key from [football-data.org](https://www.football-data.org/client/register), then:

```bash
# PowerShell
$env:FIFAWC26_API_KEY = "your_token_here"

# bash / zsh
export FIFAWC26_API_KEY="your_token_here"
```

## Usage

```bash
fifawc26                          # launch TUI
fifawc26 standings                # group tables
fifawc26 matches --team Brazil    # filter fixtures
fifawc26 scorers -n 20            # top scorers
fifawc26 --help                   # all commands
```

### TUI keys

| Key | Action |
|---|---|
| `1` `2` `3` / `Tab` | Switch tabs |
| `A`–`L` / `0` | Filter standings by group / clear |
| `↑` `↓` / `j` `k` | Move cursor on matches |
| `Enter` / `Esc` | Open detail / go back |
| `r` / `q` | Refresh / quit |

## Stack

football-data.org v4 API · BoltDB cache (`~/.fifawc26/cache.db`) · Cobra CLI · Bubble Tea + Lip Gloss · go-pretty tables.

## License

MIT