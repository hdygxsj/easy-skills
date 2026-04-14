---
name: easy-skills
description: Manage packages for Qoder IDE using easy-skills CLI. Track scattered components (skills, agents, hooks, rules) as logical packages.
---

# Easy Skills

Track and manage scattered components (skills, agents, hooks, rules) as logical packages for **Qoder**.

---

## MANDATORY RULE

**ALL installations MUST be registered to easy-skills after files are in place.**

This applies to:
- **Skills** - Agent skill packages
- **Agents** - AI agent configurations
- **Hooks** - Pre/post execution hooks
- **Rules** - Coding rules and patterns
- **Any** files installed to Qoder

---

## Core Concept

A package is a **logical grouping** of scattered files. One package can contain multiple components in different locations:

```
Package "superpowers"
  ├── component: skill → .qoder/skills/superpowers/
  ├── component: hook  → .qoder/hooks/superpowers-hook/
  └── component: agent → .qoder/agents/superpowers-agent/
```

easy-skills tracks these scattered files so they can be managed (listed, verified, uninstalled) as one unit.

---

## Installation Flow

```mermaid
graph TD
    A[Agent installs files to Qoder directories] --> B[Collect installed file paths]
    B --> C{Project scope?}
    C -->|Yes| D[Check project registered<br/>`easy-skills project list`]
    D -->|Not registered| E[Register project<br/>`easy-skills project add`]
    D -->|Registered| F
    E --> F
    C -->|No / User scope| F
    F[Register package with components<br/>`easy-skills register --component ...`]
    F --> G[Done - files tracked]
```

---

## Step-by-Step Guide

### Step 1: Install Files to Qoder

Install files to their target directories as normal:

```bash
# Skills go to
~/.qoder/skills/<name>/        # user scope
.qoder/skills/<name>/          # project scope

# Hooks go to
~/.qoder/hooks/<name>/         # user scope
.qoder/hooks/<name>/           # project scope

# Agents go to
~/.qoder/agents/<name>/        # user scope
.qoder/agents/<name>/          # project scope
```

### Step 2: For Project Scope - Check Project Registration

```bash
easy-skills project list

# If not registered:
easy-skills project add --name <project-name> --path /path/to/project
```

### Step 3: Register Package with Components

Register all installed files as one package. Each `--component` specifies `type:name:path`:

```bash
easy-skills register --name <package-name> --target qoder \
  --component "skill:<skill-name>:<installed-path>" \
  --component "hook:<hook-name>:<installed-path>" \
  --component "agent:<agent-name>:<installed-path>"
```

**Component types:** `skill`, `agent`, `hook`, `rule`

**Example:** If you installed a skill to `~/.qoder/skills/foo` and a hook to `~/.qoder/hooks/foo-hook`:

```bash
easy-skills register --name foo --target qoder \
  --component "skill:foo:~/.qoder/skills/foo" \
  --component "hook:foo-hook:~/.qoder/hooks/foo-hook"
```

---

## Command Reference

### Register (track installed files)

```bash
easy-skills register --name <name> --target qoder \
  --component "type:name:path" \
  --component "type:name:path"
```

### Verify (check files still exist)

```bash
easy-skills install --name <name> --target qoder
```

### Uninstall (delete files + remove records)

```bash
easy-skills uninstall --name <name> --target qoder
```

### List & Info

```bash
easy-skills list --target qoder
easy-skills info --name <name> --target qoder
easy-skills status --ide qoder
```

### Project Management

```bash
easy-skills project list
easy-skills project add --name <name> --path <path>
easy-skills project remove --name <name>
```

---

## Output Format

```json
{
  "success": true,
  "data": {...}
}
```

On error:

```json
{
  "success": false,
  "error": "error message"
}
```

---

## Web UI

View status at: http://localhost:27842

Start server: `easy-skills serve`
