# AI Novel Backend - 基于 Eino 框架

AI 网络小说创作后端服务，充分利用 CloudWeGo Eino 框架的核心能力。

## 🚀 技术栈

- **框架**: CloudWeGo Eino v0.7.0
- **语言**: Go 1.25+
- **Web**: Gin
- **数据库**: Neon PostgreSQL
- **向量库**: Pinecone
- **LLM**: DeepSeek/OpenAI (Eino 集成)
- **部署**: 阿里云 FC 3.0

## 🎯 Eino 核心功能使用

### Components（组件）
| 组件 | 用途 | 实现 |
|------|------|------|
| **ChatModel** | LLM 对话生成 | OpenAI 兼容模型 |
| **Embedding** | 文本向量化 | OpenAI Embedding |
| **Retriever** | 语义检索 | Pinecone Retriever |
| **ToolsNode** | 工具调用 | 自定义创作工具 |
| **Lambda** | 自定义逻辑 | 数据处理节点 |
| **ChatTemplate** | 提示词模板 | 角色/情节模板 |

### Orchestration（编排）
| 编排方式 | 场景 | 说明 |
|----------|------|------|
| **Chain** | 线性流程 | 章节生成、大纲扩写 |
| **Graph** | 复杂流程 | 多 Agent 协作、状态流转 |
| **Workflow** | 工作流 | 创作审批流程 |

### Agents（智能体）
| Agent | 职责 |
|-------|------|
| **PlanningAgent** | 剧情规划、大纲设计 |
| **WritingAgent** | 章节创作、内容生成 |
| **EditingAgent** | 内容审核、润色优化 |
| **KnowledgeAgent** | 知识库管理、检索 |

## 📁 项目结构

```
ai-novel-backend/
├── cmd/
│   └── server/                  # 应用入口
├── internal/
│   ├── config/                  # 配置管理
│   ├── database/                # 数据库连接
│   ├── handler/                 # HTTP 处理器
│   ├── eino/                    # Eino 核心实现
│   │   ├── components/          # 自定义组件
│   │   ├── chains/              # Chain 编排
│   │   ├── graphs/              # Graph 编排
│   │   ├── agents/              # 多 Agent 系统
│   │   └── tools/               # 工具定义
│   ├── model/                   # 数据模型
│   ├── repository/              # 数据访问层
│   ├── service/                 # 业务逻辑层
│   └── server/                  # 服务器配置
├── scripts/                     # 构建脚本
├── config.yaml                  # 配置文件
├── s.yaml                       # FC 部署配置
└── DEVELOPMENT.md               # 开发指南
```

## 🔧 快速开始

### 环境要求
- Go 1.25+
- Neon PostgreSQL 账号
- Pinecone 账号
- DeepSeek/OpenAI API Key

### 本地开发
```bash
git clone https://github.com/cycling02/ai-novel-backend.git
cd ai-novel-backend
cp .env.example .env
# 编辑 .env 配置你的 API Key
go mod download
go run cmd/server/main.go
```

### 部署到阿里云 FC
```bash
./scripts/build.sh
s deploy
```

## 📖 API 端点

### 小说管理
- `POST /api/v1/novels` - 创建小说
- `GET /api/v1/novels` - 获取小说列表
- `GET /api/v1/novels/:id` - 获取小说详情
- `PUT /api/v1/novels/:id` - 更新小说
- `DELETE /api/v1/novels/:id` - 删除小说

### AI 创作
- `POST /api/v1/novels/:id/generate` - AI 生成章节
- `POST /api/v1/novels/:id/outline` - AI 生成大纲
- `POST /api/v1/novels/:id/suggest` - AI 情节建议
- `POST /api/v1/novels/:id/edit` - AI 润色编辑

### 知识库
- `POST /api/v1/knowledge` - 添加知识
- `GET /api/v1/knowledge` - 检索知识
- `DELETE /api/v1/knowledge/:id` - 删除知识

## 🧠 Eino 架构设计

### 章节生成 Chain
```
用户输入 → Lambda(检索知识) → ChatTemplate → ChatModel → Lambda(后处理) → 输出
```

### 多 Agent Graph
```
                    ┌──────────────┐
                    │ EntryAgent   │
                    └──────┬───────┘
                           │
         ┌─────────────────┼─────────────────┐
         │                 │                 │
    ┌────▼────┐      ┌────▼────┐      ┌────▼────┐
    │Planning │      │Writing  │      │Editing  │
    │Agent    │      │Agent    │      │Agent    │
    └────┬────┘      └────┬────┘      └────┬────┘
         │                 │                 │
         └─────────────────┼─────────────────┘
                           │
                    ┌──────▼───────┐
                    │ OutputAgent  │
                    └──────────────┘
```

## 📝 许可证

MIT
