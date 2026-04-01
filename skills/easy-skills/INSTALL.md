# Install Easy Skills

## For Qoder

Tell Qoder agent to install easy-skills skill:

```
/easy-skills Fetch and follow instructions from http://localhost:27842/qoder/easy-skills.md
```

Or via command line:

```bash
# Add to your Qoder skills directory
mkdir -p ~/.qoder/skills/easy-skills
cp -r skills/easy-skills/* ~/.qoder/skills/easy-skills/
```

## For Cursor

Tell Cursor to install easy-skills:

```
Fetch and follow instructions from http://localhost:27842/cursor/easy-skills.md
```

Or via command line:

```bash
# Add to your Cursor rules directory
mkdir -p ~/.cursorrules/easy-skills
cp -r skills/easy-skills/* ~/.cursorrules/easy-skills/
```

## Quick Verification

After installation, try:
```bash
easy-skills list
easy-skills status
```

## Port Configuration

The web server runs on port 27842 by default. Update this if you change the port.
