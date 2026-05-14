# cheatcom

A personal CLI cheatsheet to save commands you always forget and retrieve them instantly.

[![License](https://img.shields.io/github/license/Eugene7997/cheatcom)](LICENSE)
[![Build](https://img.shields.io/github/actions/workflow/status/Eugene7997/cheatcom/ci.yml)](https://github.com/Eugene7997/cheatcom/actions)
[![Version](https://img.shields.io/github/v/release/Eugene7997/cheatcom)](https://github.com/Eugene7997/cheatcom/releases)
[![Coverage](https://img.shields.io/codecov/c/github/Eugene7997/cheatcom)](https://codecov.io/gh/Eugene7997/cheatcom)

## Install

### Pre-built binaries (recommended)

Download the latest release for your platform from the [GitHub Releases page](https://github.com/Eugene7997/cheatcom/releases/latest):

| Platform | File |
|----------|------|
| Windows (x86_64) | `cheatcom_Windows_x86_64.zip` |
| Windows (arm64) | `cheatcom_Windows_arm64.zip` |
| macOS (Apple Silicon) | `cheatcom_Darwin_arm64.tar.gz` |
| macOS (Intel) | `cheatcom_Darwin_x86_64.tar.gz` |
| Linux (x86_64) | `cheatcom_Linux_x86_64.tar.gz` |
| Linux (arm64) | `cheatcom_Linux_arm64.tar.gz` |
| Linux (i386) | `cheatcom_Linux_i386.tar.gz` |

Extract the archive and place the `chc` (or `chc.exe` on Windows) binary somewhere on your `PATH`.

### Build from source

```bash
git clone https://github.com/Eugene7997/cheatcom
cd cheatcom
go build -o chc.exe ./cmd/cheatcom
```

## Usage

### Add a command

```bash
chc add "lsof -i :80" -d "find process on port" -t networking
chc add "git log --oneline --graph --decorate --all" -d "pretty git log" -t git,log
```

### Browse and pick interactively

```bash
chc
```

Opens a fuzzy-searchable list. Press:

| Key | Action |
|-----|--------|
| `enter` | Default action (copy by default — see below) |
| `r` | Run in shell |
| `p` | Print to stdout |
| `/` | Filter list |
| `q` / `esc` / `ctrl+c` | Quit |

### Edit a cheat

```bash
chc edit <id>
```

Opens an interactive form pre-filled with the current values. Navigate with:

| Key | Action |
|-----|--------|
| `tab` / `↓` | Next field |
| `shift+tab` / `↑` | Previous field |
| `enter` | Advance to next field (on last field: save and exit) |
| `esc` / `ctrl+c` | Cancel without saving |

Fields: **Command**, **Description**, **Tags** (comma-separated).

### Other commands

```bash
chc list              # print all cheats in a table
chc rm <id>           # delete by ID
```

### Set a default action

Change what `enter` does in the TUI:

```bash
chc config default copy    # enter copies to clipboard (factory default)
chc config default run     # enter executes the command
chc config default print   # enter prints to stdout
```

The setting is saved to your cheats file and persists across sessions.

### Action flags

The `--copy`, `--run`, and `--print` flags override both the configured default and the TUI key press — the list still opens for selection, but the action is fixed:

```bash
chc --copy            # select from list, then copy to clipboard
chc --run             # select from list, then execute
chc --print           # select from list, then print to stdout (useful with eval)
```

## Data file

Cheats are stored as YAML at:

- **Windows**: `%AppData%\cheatcom\cheats.yaml`
- **macOS**: `~/Library/Application Support/cheatcom/cheats.yaml`
- **Linux**: `~/.config/cheatcom/cheats.yaml`

The file is human-readable and safe to edit, back up, or sync via Dropbox/git.
