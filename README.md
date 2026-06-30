# ssh-wizard

Terminal SSH setup wizard built with [Bubble Tea](https://github.com/charmbracelet/bubbletea).  
Nothing Design System aesthetic — monochrome, Space Mono, instrument-panel labels.

```
  ▪▪▪▪▪▪▪▪▪▫▫▫▫▫   2 / 7

  SERVER INFO

  ✓ ALIAS                 homelab
  ✓ HOSTNAME              192.168.1.10
  ✓ USERNAME              pi
  ✓ PORT                  22
```

## Features

- Guided SSH key generation (Ed25519 / RSA 4096)
- Remote key installation via password (one-time)
- `~/.ssh/config` writer with conflict detection
- Connection test with live checklist
- VS Code Remote-SSH ready

## Requirements

- `ssh` and `ssh-keygen` on PATH (pre-installed on Linux, macOS, WSL, Git Bash)

## Install

### Linux / macOS / WSL

```bash
curl -fsSL https://raw.githubusercontent.com/duong6003/ssh-wizard/master/install.sh | sh
```

### Windows (PowerShell)

```powershell
irm https://raw.githubusercontent.com/duong6003/ssh-wizard/master/install.ps1 | iex
```

Installs to `%LOCALAPPDATA%\Programs\ssh-wizard` and adds it to PATH automatically.

### Or download manually

Grab the binary for your platform from [Releases](https://github.com/duong6003/ssh-wizard/releases).

### With Go installed

```bash
go install github.com/duong6003/ssh-wizard@latest
```

## Build from source

```bash
git clone https://github.com/duong6003/ssh-wizard.git
cd ssh-wizard
go build -o ssh-wizard .
./ssh-wizard
```

Or with Make:

```bash
make build   # build for current platform
make test    # run tests
make release # cross-compile to dist/
```

## Usage

```bash
ssh-wizard
```

Follow the on-screen prompts. At the end, connect with:

```bash
ssh <alias>
```

### Keyboard shortcuts

| Key | Action |
|-----|--------|
| `Tab` / `Enter` | Next field |
| `↑` / `↓` | Select option |
| `Ctrl+C` | Quit |
| `A` (Done screen) | Add another server |

### ASCII fallback

If your terminal doesn't support Unicode:

```bash
SSH_WIZARD_ASCII=1 ssh-wizard
```

## License

MIT
