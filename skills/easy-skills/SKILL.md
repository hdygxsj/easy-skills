# Easy Skills

Easy Skills Hub CLI tool for managing AI IDE Skills.

## Overview

Easy Skills is a local skill hub that helps you manage AI IDE skills across different IDEs (Qoder, Cursor). It provides a unified way to:
- Register skill packages from various sources
- Track which skills are installed where
- Switch between skill sets easily

## Available Commands

### List Packages
```bash
easy-skills list --target qoder
```
Shows all registered packages for Qoder.

### Check Status
```bash
easy-skills status --ide qoder
```
Shows which packages are installed in Qoder (user/project scope).

### Register a Package
```bash
easy-skills register --name superpowers --target qoder --source https://github.com/obra/superpowers
```
Register a new package to the local Hub.

### Show Package Info
```bash
easy-skills info --name superpowers --target qoder
```
Show detailed information including all components.

## Workflow Examples

### Install a Skill Package
1. Check current status: `easy-skills status --ide qoder`
2. List available packages: `easy-skills list --target qoder`
3. Register if needed: `easy-skills register --name <pkg> --target qoder --source <url>`
4. (Installation would be done via IDE's native mechanism)

### Switch Skill Sets
1. Check current: `easy-skills status --ide qoder`
2. Note installed packages
3. Uninstall old packages via IDE
4. Register new packages: `easy-skills register --name <new-pkg> --target qoder --source <url>`

## Output Format

All commands output JSON:
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
