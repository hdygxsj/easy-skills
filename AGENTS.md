# Local Skill Hub - Agent 知识库

> 本文档记录 local-skill-hub 项目开发所需的知识、规范和决策。

---

## 项目定位

**local-skill-hub** 是一个运行在用户本地电脑的 skill 管理工具，用于：
1. 集中管理从各种来源获取的 skill/agent/hooks 等技能包
2. 支持 skill 分组，可以灵活组合和打包
3. 将技能包安装到不同的 AI IDE（Qoder、Cursor 等）
4. 提供冲突检测和快速切换能力
5. 支持通过对话方式安装技能包

**技术栈**：
- 前后端分离项目
- 本地运行，使用偏口（如 3847 等）
- 支持配置大模型 API Key 进行对话

**开源协议**：MIT

---

## AI IDE 安装规则参考

### Qoder 安装规则

#### 目录结构

| 类型 | 用户级路径 | 项目级路径 | 说明 |
|------|-----------|-----------|------|
| Skills | `~/.qoder/skills/{skill-name}/SKILL.md` | `.qoder/skills/{skill-name}/SKILL.md` | 技能包 |
| Agents | `~/.qoder/agents/{agentName}.md` | `.qoder/agents/{agentName}.md` | 自定义智能体 |
| Hooks | `~/.qoder/hooks/*.sh` | `.qoder/hooks/*.sh` | 钩子脚本 |
| Settings | `~/.qoder/settings.json` | `.qoder/settings.json` | 配置文件 |

#### Agent 文件格式

```markdown
---
name: agent-name
description: 简短描述功能和专长
tools: Read, Grep, Glob, Bash
---

你是一位资深代码审查员，负责确保代码质量。
...
```

支持的工具列表：
- `Bash` - 执行 shell 命令
- `create_file` - 创建/覆盖文件
- `search_replace` - 编辑文件
- `search_file` - 检索文件
- `grep_code` - 搜索文件内容
- `read_file` - 读取文件内容
- `fetch_content` - 获取网页内容
- `search_web` - 网页搜索
- `todo_write` - 写入 TODO
- `ask_user_question` - 向用户提问
- `search_memory` - 搜索记忆
- `update_memory` - 更新记忆
- `switch_mode` - 切换模式
- `create_plan` - 创建计划
- `run_preview` - 预览 Web 应用
- `mcp__<server>__<tool>` - MCP 工具

#### Hooks 配置格式

**settings.json 结构**：
```json
{
  "hooks": {
    "事件名": [
      {
        "matcher": "匹配条件",
        "hooks": [
          {
            "type": "command",
            "command": "脚本路径"
          }
        ]
      }
    ]
  }
}
```

**支持的事件**：
| 事件名称 | 触发时机 | 可阻断 |
|----------|----------|--------|
| UserPromptSubmit | 用户提交 Prompt 后 | 是 |
| PreToolUse | 工具调用执行前 | 是 |
| PostToolUse | 工具调用成功后 | 否 |
| PostToolUseFailure | 工具调用失败后 | 否 |
| Stop | Agent 完成响应时 | 否 |

**matcher 匹配规则**：
- 不填或 `"*"` - 匹配所有
- 精确值 - 精确匹配，如 `"Bash"`
- `|` 分隔 - 匹配多个值，如 `"Write | Edit"`
- 正则表达式 - 如 `"mcp__.*"`

#### 安装指引 (INSTALL.md) 格式

典型的安装指引包含以下步骤：

1. **确认安装位置** - 用户级 (`~/.qoder/`) 或项目级 (`.qoder/`)
2. **Clone 仓库** - `git clone` 到指定位置
3. **创建符号链接** - 将 skills/agents 链接到对应目录
4. **配置 Hooks** - 创建 hook 脚本并配置 settings.json
5. **重启 IDE** - 使配置生效

**示例 (Superpowers Qoder 安装)**：
```bash
# Step 1: Clone the repository
git clone https://github.com/hdygxsj/superpowers.git ~/.qoder/superpowers

# Step 2: Create symlinks
mkdir -p ~/.qoder/skills ~/.qoder/agents
for skill in ~/.qoder/superpowers/skills/*/; do
  ln -sf "$skill" ~/.qoder/skills/$(basename "$skill")
done
for agent in ~/.qoder/superpowers/agents/*.md; do
  ln -sf "$agent" ~/.qoder/agents/
done

# Step 3: Configure Hooks
# ... 创建 hook 脚本和 settings.json

# Step 4: Restart Qoder
```

### Cursor 安装规则

> 待补充：需要进一步了解 Cursor 的 skill/agent/hooks 安装机制

#### 参考文档
- https://cursor.com/cn/docs/skills
- https://cursor.com/cn/docs/hooks
- https://cursor.com/cn/docs/rules
- https://cursor.com/cn/docs/subagents

---

## 技能包格式

### 标准技能包结构

```
{package-name}/
├── .qoder/                    # Qoder 配置
│   ├── skills/
│   │   └── {skill-name}/
│   │       └── SKILL.md
│   ├── agents/
│   │   └── {agent-name}.md
│   ├── hooks/
│   │   └── {hook-name}.sh
│   └── settings.json
├── .cursor/                   # Cursor 配置 (待确认)
│   └── ...
├── INSTALL.md                 # 安装指引
├── README.md                  # 包说明
└── metadata.json              # 元数据（版本、依赖等）
```

### metadata.json 格式

```json
{
  "name": "package-name",
  "version": "1.0.0",
  "description": "技能包描述",
  "targets": ["qoder", "cursor"],
  "skills": ["skill-a", "skill-b"],
  "agents": ["agent-x"],
  "hooks": ["hook-y"],
  "dependencies": {
    "other-package": ">=1.0.0"
  }
}
```

### INSTALL.md 格式要求

INSTALL.md 应包含：
1. **前置条件** - 需要的软件和环境
2. **安装步骤** - 具体的命令和操作
3. **验证方法** - 如何确认安装成功
4. **更新方法** - 如何更新到新版本

---

## 核心功能设计

### 1. 技能包管理

- **导入**：从 URL、Git 仓库、本地文件导入技能包
- **存储**：在本地仓库统一管理所有技能包
- **分组**：将技能包按用途分组（如 "开发流程"、"代码审查"）
- **元数据**：记录技能包的名称、版本、来源、依赖等

### 安装路径说明

- **用户级安装**：自动使用默认路径（如 `~/.qoder/`）
- **项目级安装**：LSH 不预设项目位置，需要用户输入项目目录路径

### 2. 冲突检测

- **同名检测**：检测是否有同名的 skill/agent/hook
- **依赖冲突**：检测依赖版本是否兼容
- **覆盖风险**：检测安装是否会导致已有配置被覆盖

### 3. 快速切换

- **分组切换**：一键切换整个技能包组
- **临时禁用**：临时禁用某个 skill 或 agent
- **快照管理**：保存当前配置快照，快速回滚

### 4. 安装到 IDE

- **目标选择**：选择安装到哪个 IDE
- **安装模式**：用户级或项目级安装
- **符号链接**：使用符号链接避免文件复制
- **配置合并**：智能合并 settings.json

### 5. 对话式安装

- **自然语言安装**：用户通过对话描述需求
- **自动解析**：从 INSTALL.md 提取安装流程
- **API 转发**：未来支持通过 HTTP 接口接收安装指令

---

## 设计决策记录

### 2026-04-01: 项目启动

- **决策**：采用前后端分离架构
- **原因**：便于后续扩展和支持多种 AI IDE
- **状态**：已确认

### 2026-04-01: 安装位置策略

- **决策**：支持用户级 (`~/.qoder/`) 和项目级 (`.qoder/`) 两种安装方式
- **原因**：用户级全局生效，项目级可版本控制
- **状态**：已确认

---

## 待探索

- [ ] Cursor 的 skill/agent/hooks 具体安装机制
- [ ] 更多 AI IDE 的安装规则（如 Codex、OpenCode、Gemini CLI）
- [ ] 技能包依赖管理的具体实现方式
- [ ] 冲突解决策略（优先覆盖、提示用户、拒绝安装等）
