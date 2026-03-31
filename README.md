# setapp

A command-line interface for browsing, searching, and managing [Setapp](https://go.setapp.com/invite/1vndseho) (referral link) applications on macOS.

Reads directly from Setapp's local SQLite database to provide fast access to app metadata.

## Installation

```sh
go install github.com/thirteen37/setapp@latest
```

## Usage

```
setapp [command] [flags]
```

### Commands

| Command | Description |
|---|---|
| `list` | List all available Setapp apps |
| `search <query>` | Search apps by name, keyword, tagline, or vendor |
| `info <app>` | Show detailed information about an app |
| `categories` | List all Setapp categories |
| `install <app>` | Install an app |
| `uninstall <app>` | Uninstall an app |
| `open <app>` | Open an installed app |
| `upgrade` | Upgrade installed apps |
| `doctor` | Check Setapp installation health |
| `home <app>` | Open an app's homepage in the browser |

### Global Flags

| Flag | Description |
|---|---|
| `--json` | Output as JSON |

## Requirements

- macOS with [Setapp](https://go.setapp.com/invite/1vndseho) (referral link) installed
- Go 1.26.1+

## License

MIT
