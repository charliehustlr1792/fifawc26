<p align="center">
  <img src="./ascii-art-text.png" alt="Elite Ball Knowledge" />
</p>

<p align="center">
  A terminal app for the FIFA World Cup 2026. Live scores, fixtures, standings, top scorers.
</p>

<p align="center">
  <a href="https://github.com/charliehustlr1792/fifawc26/releases/latest">
    <img src="https://img.shields.io/github/v/release/charliehustlr1792/fifawc26?style=flat-square" alt="Release" />
  </a>
  <a href="https://github.com/charliehustlr1792/fifawc26/blob/main/LICENSE">
    <img src="https://img.shields.io/github/license/charliehustlr1792/fifawc26?style=flat-square" alt="License" />
  </a>
</p>

## Install

**Homebrew** (macOS, Linux)
```bash
brew install charliehustlr1792/tap/fifawc26
```

**Scoop** (Windows)
```powershell
scoop bucket add fifawc26 https://github.com/charliehustlr1792/scoop-bucket
scoop install fifawc26
```

**Go**
```bash
go install github.com/charliehustlr1792/fifawc26/cmd/fifawc26@latest
```

**Manual:** grab the binary for your OS from [Releases](https://github.com/charliehustlr1792/fifawc26/releases/latest) and put it on your PATH.

## First Run

```bash
fifawc26
```

On first launch you pick a tier:

1. **With API key** (recommended). Full data via [football-data.org](https://www.football-data.org/client/register). Free, takes 30 seconds to register, the app guides you through it.
2. **Keyless**. Works out of the box. Live scores limited; some data thinner.

Switch tiers anytime:
```bash
fifawc26 setup
```

## Features

- Group standings with W/D/L, GF/GA/GD, points, form
- Full fixture list with filters by team, matchday, status
- Match detail view with kickoff time, score, half-time, winner
- Team detail with squad and that team's fixtures
- Top scorers leaderboard
- Auto-refresh every 45 seconds
- Local cache so you don't burn API quota

## TUI Keys

| Key | Action |
|---|---|
| `1` `2` `3` / `Tab` | Switch tabs |
| `A` to `L` | Filter standings by group |
| `0` | Show all groups |
| `↑` `↓` | Move cursor |
| `Enter` | Open detail (team or match) |
| `t` | Open home team from match detail |
| `Esc` | Back |
| `r` | Refresh |
| `q` | Quit |

## CLI

```bash
fifawc26 standings
fifawc26 matches --team Brazil
fifawc26 matches --status SCHEDULED
fifawc26 scorers -n 20
fifawc26 version
```

## License

MIT