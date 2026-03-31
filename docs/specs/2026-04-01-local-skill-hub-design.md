# Local Skill Hub 设计文档

## 概述

Local Skill Hub (LSH) 是一个运行在用户本地的技能包管理工具，用于集中管理、分组和安装 AI IDE 技能包（skills、agents、hooks 等）到不同的 AI IDE（Qoder、Cursor 等）。

**技术栈**：Python (FastAPI) + Vue + SQLite  
**协议**：MIT

---

## 核心概念

### 概念模型

```
Source (来源)
├── Git Repo URL
├── Install Guide URL  
└── Local Files
        ↓
Package (本地副本)
├── skills/、agents/、hooks/ 等文件
├── version (git hash / timestamp)
└── install_script (解析的安装脚本)
        ↓
Group (分组)
├── 包含多个 Package
└── 目标 IDE 配置
        ↓
Installation (安装记录)
├── 目标 IDE (Qoder/Cursor)
├── 安装路径 (用户级/项目级)
├── 已安装文件列表
└── 安装时版本
```

### 关键决策

| 决策 | 选择 | 原因 |
|------|------|------|
| 安装记录位置 | LSH 侧 | 不污染目标目录 |
| 安装方式 | 默认复制，可选符号链接 | 兼容性优先 |
| 版本对比范围 | 仅本地 vs IDE | 简化设计 |
| 冲突处理 | 不检测，支持快速卸载/切换 | 用户自行管理 |

---

## 功能清单

### 1. 包管理

| 功能 | 说明 |
|------|------|
| 从 Git 导入 | `git clone` 到本地存储 |
| 从 Install Guide 导入 | 解析 INSTALL.md，归档安装流程 |
| 从本地导入 | 复制本地文件到存储 |
| 查看包列表 | 展示所有本地包 |
| 查看包详情 | 展示包内容、版本、来源 |
| 更新包 | `git pull` 或重新导入 |
| 删除包 | 删除本地副本 |

### 2. 分组管理

| 功能 | 说明 |
|------|------|
| 创建分组 | 新建空分组 |
| 添加包到分组 | 将包关联到分组 |
| 从分组移除包 | 解除关联 |
| 查看分组列表 | 展示所有分组 |
| 删除分组 | 删除分组（不删除包） |

### 3. 安装管理

| 功能 | 说明 |
|------|------|
| 安装分组 | 将分组内的包安装到指定 IDE |
| 卸载分组 | 根据安装记录精准删除已安装文件 |
| 切换分组 | 卸载当前 + 安装新分组 |
| 查看安装状态 | 展示各 IDE 的安装情况 |
| 版本对比 | 本地版本 vs IDE 安装版本 |
| 升级安装 | 重新安装以更新到新版本 |

**安装路径说明**：
- **用户级安装**：自动使用默认路径（如 `~/.qoder/`）
- **项目级安装**：需要用户输入项目目录路径（LSH 不预设项目位置）

### 4. 对话式交互

| 功能 | 说明 |
|------|------|
| 自然语言交互 | 通过聊天完成上述操作 |
| INSTALL.md 解析 | Agent 读取并理解安装指引 |
| Shell 命令执行 | Agent 执行安装命令（带审核） |

### 5. 系统设置

| 功能 | 说明 |
|------|------|
| API Key 配置 | 配置 LLM API Key |
| 存储路径配置 | 配置包存储目录 |
| IDE 路径配置 | 配置各 IDE 的安装目录 |

---

## 数据模型

```sql
-- 来源
CREATE TABLE sources (
    id INTEGER PRIMARY KEY,
    type TEXT NOT NULL,           -- 'git' | 'guide' | 'local'
    url_or_path TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 本地包
CREATE TABLE packages (
    id INTEGER PRIMARY KEY,
    source_id INTEGER REFERENCES sources(id),
    name TEXT NOT NULL,
    version TEXT,                 -- git hash 或 timestamp
    local_path TEXT NOT NULL,
    metadata_json TEXT,           -- 包元数据
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 分组
CREATE TABLE groups (
    id INTEGER PRIMARY KEY,
    name TEXT NOT NULL UNIQUE,
    description TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 分组-包关联
CREATE TABLE group_packages (
    group_id INTEGER REFERENCES groups(id) ON DELETE CASCADE,
    package_id INTEGER REFERENCES packages(id) ON DELETE CASCADE,
    PRIMARY KEY (group_id, package_id)
);

-- 安装记录
CREATE TABLE installations (
    id INTEGER PRIMARY KEY,
    group_id INTEGER REFERENCES groups(id),
    target_ide TEXT NOT NULL,     -- 'qoder' | 'cursor'
    install_scope TEXT NOT NULL,  -- 'user' | 'project'
    install_path TEXT NOT NULL,   -- 安装目标路径（项目级需用户输入）
    installed_files_json TEXT,    -- 已安装文件列表
    installed_version TEXT,       -- 安装时的版本
    installed_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 系统配置
CREATE TABLE config (
    key TEXT PRIMARY KEY,
    value TEXT,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

---

## 内置 Agent

### Agent 定义

```markdown
---
name: lsh-agent
description: Local Skill Hub 内置助手，管理技能包的导入、分组、安装和切换
---

你是 Local Skill Hub 的内置助手。你可以帮助用户：
1. 导入技能包（从 Git、Install Guide 或本地文件）
2. 管理分组（创建、添加、移除）
3. 安装/卸载/切换技能包到 AI IDE
4. 查看版本状态和升级

你有以下工具可以调用...
```

### 工具列表

| 工具 | 说明 |
|------|------|
| `import_from_git` | 从 Git URL 克隆/拉取包 |
| `import_from_guide` | 从 INSTALL.md URL 解析并归档 |
| `import_from_local` | 从本地路径导入 |
| `list_packages` | 列出所有本地包 |
| `get_package_info` | 获取包详情 |
| `create_group` | 创建分组 |
| `list_groups` | 列出所有分组 |
| `add_to_group` | 将包添加到分组 |
| `remove_from_group` | 从分组移除包 |
| `install_group` | 安装分组到指定 IDE |
| `uninstall_group` | 卸载分组 |
| `switch_group` | 切换分组 |
| `list_installations` | 列出安装状态 |
| `compare_versions` | 对比版本 |
| `upgrade_installation` | 升级安装 |
| `execute_shell` | 执行 shell 命令 |
| `read_file` | 读取文件内容 |
| `fetch_url` | 获取 URL 内容 |

---

## 项目结构

```
local-skill-hub/
├── backend/                     # Python 后端
│   ├── app/
│   │   ├── __init__.py
│   │   ├── main.py             # FastAPI 入口
│   │   ├── api/                # API 路由
│   │   │   ├── __init__.py
│   │   │   ├── packages.py
│   │   │   ├── groups.py
│   │   │   ├── installations.py
│   │   │   ├── chat.py
│   │   │   └── config.py
│   │   ├── services/           # 业务逻辑
│   │   │   ├── __init__.py
│   │   │   ├── package_service.py
│   │   │   ├── group_service.py
│   │   │   ├── install_service.py
│   │   │   └── git_service.py
│   │   ├── models/             # SQLAlchemy 模型
│   │   │   ├── __init__.py
│   │   │   └── models.py
│   │   ├── installers/         # IDE 安装器
│   │   │   ├── __init__.py
│   │   │   ├── base.py
│   │   │   ├── qoder.py
│   │   │   └── cursor.py
│   │   ├── agent/              # Agent 核心
│   │   │   ├── __init__.py
│   │   │   ├── core.py
│   │   │   ├── tools.py
│   │   │   └── prompts.py
│   │   └── database.py         # 数据库连接
│   ├── tests/                  # 测试用例
│   │   ├── __init__.py
│   │   ├── test_packages.py
│   │   ├── test_groups.py
│   │   ├── test_installations.py
│   │   └── test_agent.py
│   ├── requirements.txt
│   └── pytest.ini
├── frontend/                   # Vue 前端
│   ├── src/
│   │   ├── main.js
│   │   ├── App.vue
│   │   ├── views/
│   │   │   ├── ChatView.vue
│   │   │   ├── PackagesView.vue
│   │   │   ├── GroupsView.vue
│   │   │   ├── InstallationsView.vue
│   │   │   └── SettingsView.vue
│   │   ├── components/
│   │   │   ├── ChatMessage.vue
│   │   │   ├── PackageCard.vue
│   │   │   ├── GroupCard.vue
│   │   │   └── ...
│   │   ├── api/
│   │   │   └── index.js
│   │   └── router/
│   │       └── index.js
│   ├── package.json
│   └── vite.config.js
├── agents/                     # 内置 Agent 配置
│   └── lsh-agent.md
├── data/                       # 数据目录 (gitignore)
│   ├── packages/               # 本地包存储
│   └── lsh.db                  # SQLite 数据库
├── docs/
│   └── specs/
├── AGENTS.md                   # 项目知识库
├── LICENSE                     # MIT
├── .gitignore
└── ...
```

---

## IDE 安装规则

### Qoder

| 类型 | 用户级路径 | 项目级路径 |
|------|-----------|-----------|
| Skills | `~/.qoder/skills/{name}/` | `.qoder/skills/{name}/` |
| Agents | `~/.qoder/agents/{name}.md` | `.qoder/agents/{name}.md` |
| Hooks | `~/.qoder/hooks/{name}.sh` | `.qoder/hooks/{name}.sh` |
| Settings | `~/.qoder/settings.json` | `.qoder/settings.json` |

### Cursor

> 待补充：需进一步确认 Cursor 的目录结构

---

## settings.json 合并策略

安装 hooks 时需要合并 settings.json：

```python
def merge_settings(existing: dict, new_hooks: dict) -> dict:
    result = existing.copy()
    if 'hooks' not in result:
        result['hooks'] = {}
    
    for event, handlers in new_hooks.get('hooks', {}).items():
        if event not in result['hooks']:
            result['hooks'][event] = []
        result['hooks'][event].extend(handlers)
    
    return result
```

---

## 安全考虑

1. **API Key 存储**：存储在 `data/config.json`，不提交到版本控制
2. **Shell 命令执行**：Agent 执行命令前需要用户确认（或配置白名单）
3. **文件操作**：仅在允许的目录内操作

---

## 后续扩展

- [ ] 支持更多 AI IDE（Codex、OpenCode、Gemini CLI）
- [ ] 远程版本检测和更新提醒
- [ ] 技能包依赖管理
- [ ] 技能包市场/分享
