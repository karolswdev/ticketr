# Global Installation Guide

This guide covers the global installation of Ticketr v3, directory structure, migration from v2.x, and platform-specific considerations.

## Table of Contents
- [Installation Methods](#installation-methods)
- [Directory Structure](#directory-structure)
- [Environment Variables](#environment-variables)
- [Migration from v2.x](#migration-from-v2x)
- [Platform-Specific Notes](#platform-specific-notes)
- [Troubleshooting](#troubleshooting)

## Installation Methods

### Method 1: Go Install (Recommended)

```bash
go install github.com/karolswdev/ticktr@latest
```

The binary will be installed to `$GOPATH/bin` or `$HOME/go/bin`.

### Method 2: Build from Source

```bash
git clone https://github.com/karolswdev/ticktr.git
cd ticketr
go build -o ticketr cmd/ticketr/main.go
sudo mv ticketr /usr/local/bin/
```

## Directory Structure

Ticketr v3 follows platform conventions for storing configuration, data, and cache files.

### Linux/Unix (XDG Base Directory Specification)

```
~/.config/ticketr/         # Configuration files
├── config.yaml           # Main configuration
└── workspaces.yaml      # Workspace definitions

~/.local/share/ticketr/    # Persistent data
├── ticketr.db           # SQLite database
└── templates/           # Ticket templates

~/.cache/ticketr/          # Cache files
├── jira_schema.json     # Cached field mappings
└── logs/                # Application logs
```

### macOS

```
~/Library/Application Support/ticketr/  # Config and data
├── config.yaml
├── workspaces.yaml
├── ticketr.db
└── templates/

~/Library/Caches/ticketr/              # Cache files
├── jira_schema.json
└── logs/
```

### Windows

```
%APPDATA%\ticketr\              # Configuration
├── config.yaml
└── workspaces.yaml

%LOCALAPPDATA%\ticketr\         # Data and cache
├── ticketr.db
├── templates/
└── cache/
```

## Environment Variables

Override default paths using environment variables:

```bash
# Override configuration directory
export TICKETR_CONFIG_HOME=/custom/config

# Override data directory
export TICKETR_DATA_HOME=/custom/data

# Override cache directory
export TICKETR_CACHE_HOME=/custom/cache
```

### XDG Variables (Linux/Unix)

Ticketr respects standard XDG environment variables:

```bash
export XDG_CONFIG_HOME=$HOME/.config
export XDG_DATA_HOME=$HOME/.local/share
export XDG_CACHE_HOME=$HOME/.cache
```

## Migration from v2.x

### Understanding the Changes

**v2.x Structure (Old):**
```
~/.ticketr/
├── config.yaml
├── ticketr.db
└── cache/
```

**v3.x Structure (New):**
- Configuration moved to XDG-compliant directories
- Better separation of config, data, and cache
- Platform-specific paths

### Manual Migration Steps

1. **Backup existing data:**
```bash
cp -r ~/.ticketr ~/.ticketr.backup
```

2. **Create new workspace:**
```bash
ticketr workspace create default \
  --url "https://your.atlassian.net" \
  --project "PROJ" \
  --set-default
```

3. **Copy templates (optional):**
```bash
cp -r ~/.ticketr/templates/* ~/.local/share/ticketr/templates/
```

4. **Verify migration:**
```bash
ticketr workspace list
```

## Platform-Specific Notes

### Linux

Install dependencies for keychain support:
```bash
# Debian/Ubuntu
sudo apt-get install libsecret-1-dev

# Fedora
sudo dnf install libsecret-devel

# Arch
sudo pacman -S libsecret
```

### macOS

On first run, you'll be prompted for keychain access. Click "Always Allow" for seamless operation.

### Windows

Credentials are stored in Windows Credential Manager. You can view them at:
Control Panel → Credential Manager → Windows Credentials

## Troubleshooting

### Common Issues

**"Permission denied" when creating directories**
- Ensure you have write permissions to your home directory
- Check parent directory permissions

**"Failed to access keychain"**
- Linux: Install libsecret
- macOS: Unlock Keychain Access
- Windows: Check Credential Manager service

**Directory not created automatically**
- Run `ticketr workspace create` to initialize directories
- Manually create with: `mkdir -p ~/.config/ticketr`

### Getting Help

- GitHub Issues: https://github.com/karolswdev/ticktr/issues
- Documentation: https://github.com/karolswdev/ticktr/docs

## See Also

- [README.md](../README.md)
- [Workspace Management Guide](workspace-management-guide.md)
- [Architecture Documentation](ARCHITECTURE.md)
- [V3 Implementation Roadmap](v3-implementation-roadmap.md)