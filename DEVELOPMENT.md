# 开发指南

## 📋 环境要求

- Go 1.25+
- Neon PostgreSQL 账号
- Pinecone 账号
- DeepSeek/OpenAI API Key
- 阿里云 CLI（部署用）

## 🚀 快速开始

### 1. 克隆项目

```bash
git clone https://github.com/cycling02/ai-novel-backend.git
cd ai-novel-backend
```

### 2. 配置环境变量

```bash
cp .env.example .env
# 编辑 .env 填入你的配置
```

**必需配置:**
- `DATABASE_URL` - Neon PostgreSQL 连接字符串
- `PINECONE_API_KEY` - Pinecone API Key
- `DEEPSEEK_API_KEY` - DeepSeek API Key
- `JWT_SECRET` - JWT 密钥

### 3. 安装依赖

```bash
go mod download
```

### 4. 运行服务

```bash
go run cmd/server/main.go
```

服务启动在 `http://localhost:8080`

### 5. 测试 API

```bash
# 健康检查
curl http://localhost:8080/health

# 创建小说
curl -X POST http://localhost:8080/api/v1/novels \
  -H "Content-Type: application/json" \
  -d '{"title":"我的小说","description":"简介","genre":"玄幻"}'

# AI 生成章节
curl -X POST http://localhost:8080/api/v1/novels/{id}/generate \
  -H "Content-Type: application/json" \
  -d '{"chapter_title":"第一章","outline":"主角登场..."}'
```

## 🏗️ 项目架构

### 目录结构

```
ai-novel-backend/
├── cmd/server/              # 应用入口
├── internal/
│   ├── config/              # 配置管理
│   ├── database/            # 数据库连接和迁移
│   ├── eino/                # Eino 框架集成
│   │   ├── agents/          # 多 Agent 系统
│   │   ├── chains/          # Chain 编排
│   │   ├── components/      # Eino 组件
│   │   └── graphs/          # Graph 编排
│   ├── handler/             # HTTP 处理器
│   ├── model/               # 数据模型
│   ├── repository/          # 数据访问层
│   ├── server/              # 服务器配置
│   └── service/             # 业务逻辑层
├── scripts/                 # 构建脚本
├── config.yaml              # 配置文件
└── s.yaml                   # FC 部署配置
```

### Eino 架构

#### Components（组件）
- **ChatModel** - LLM 对话生成（支持 DeepSeek/OpenAI）
- **Embedding** - 文本向量化
- **Retriever** - Pinecone 语义检索
- **ChatTemplate** - 提示词模板
- **Tools** - 工具调用（角色查询、世界观查询等）

#### Chains（链式编排）
- **ChapterGenerateChain** - 章节生成
- **OutlineExpandChain** - 大纲扩写
- **PlotSuggestChain** - 情节建议
- **ContentEditChain** - 内容编辑

#### Graphs（图编排）
- **NovelCreationGraph** - 小说创作完整流程

#### Agents（智能体）
- **PlanningAgent** - 剧情规划
- **WritingAgent** - 章节创作
- **EditingAgent** - 内容审核
- **ResearchAgent** - 资料检索

## 📦 构建

### 本地构建

```bash
go build -o ai-novel-backend ./cmd/server
```

### 使用脚本

```bash
./scripts/build.sh
```

### 交叉编译（FC 部署）

```bash
GOOS=linux GOARCH=amd64 go build -o ai-novel-backend ./cmd/server
```

## 🚀 部署

### 阿里云 FC 3.0

```bash
# 安装 Serverless Devs
npm install -g @serverless-devs/s

# 配置阿里云账号
s config add

# 部署
s deploy
```

## 🧪 测试

```bash
# 运行测试
go test ./...

# 覆盖率
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

## 📖 API 文档

### 小说管理

| 方法 | 端点 | 说明 |
|------|------|------|
| POST | `/api/v1/novels` | 创建小说 |
| GET | `/api/v1/novels` | 获取小说列表 |
| GET | `/api/v1/novels/:id` | 获取小说详情 |
| PUT | `/api/v1/novels/:id` | 更新小说 |
| DELETE | `/api/v1/novels/:id` | 删除小说 |

### AI 创作

| 方法 | 端点 | 说明 |
|------|------|------|
| POST | `/api/v1/novels/:id/generate` | AI 生成章节 |
| POST | `/api/v1/novels/:id/outline` | AI 生成大纲 |
| POST | `/api/v1/novels/:id/suggest` | AI 情节建议 |
| POST | `/api/v1/content/edit` | AI 润色编辑 |

### 健康检查

| 方法 | 端点 | 说明 |
|------|------|------|
| GET | `/health` | 整体健康检查 |
| GET | `/health/ready` | 就绪检查 |
| GET | `/health/live` | 存活检查 |

## 🔧 故障排查

### 数据库连接失败

检查 `DATABASE_URL` 是否正确，确保 Neon 项目已创建数据库。

### LLM 调用失败

检查 API Key 是否正确，确保有足够的配额。

### Pinecone 连接失败

检查 `PINECONE_API_KEY` 和 `PINECONE_INDEX_NAME` 是否正确。

## 📝 许可证

MIT
