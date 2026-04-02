# Easy Skills

Use `easy-skills` CLI to manage AI IDE Skills for this project.

## Overview

Easy Skills is a local skill hub that manages skill packages across different IDEs (Qoder, Cursor).

## Commands

### Check Current Status

```bash
# Check installed packages in Cursor
easy-skills status --ide cursor

# Check installed packages in Qoder
easy-skills status --ide qoder
```

### List Available Packages

```bash
# List all packages for Cursor
easy-skills list --target cursor

# List all packages for Qoder
easy-skills list --target qoder
```

### Register a Package

```bash
# Register a new package to Hub
easy-skills register --name <package-name> --target cursor --source <git-url>
```

### View Package Details

```bash
# Show package info including components
easy-skills info --name <package-name> --target cursor
```

## Workflow Examples

### Before Starting a Task

1. Check what's installed: `easy-skills status --ide cursor`
2. List available packages: `easy-skills list --target cursor`
3. Note which skills are available for this project type

### When User Asks to Use a Skill

1. Check if skill is in installed list
2. If not installed, register it first: `easy-skills register --name <pkg> --target <ide> --source <url>`
3. Use the skill per its documentation

## Output Format

All commands return JSON:

```json
{
  "success": true,
  "data": {
    "packages": [...]
  }
}
```

On error:

```json
{
  "success": false,
  "error": "error message"
}
```

## Web UI

View skill status visually at: http://localhost:27842
