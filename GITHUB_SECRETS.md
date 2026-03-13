# GitHub Secrets & Variables 配置指南

本文档列出在 GitHub 上需要配置的所有环境变量和密钥。

## 📍 配置位置

1. 打开项目：https://github.com/cycling02/ai-novel-backend
2. 进入 **Settings** → **Secrets and variables** → **Actions**
3. 分别配置 **Secrets** 和 **Variables**

---

## 🔐 Secrets (加密密钥)

这些是敏感信息，必须加密存储。

| 名称 | 说明 | 示例值 | 必需 | 用途 |
|------|------|--------|------|------|
| `ALIYUN_ACCOUNT_ID` | 阿里云账号 ID | `1234567890123456` | ✅ | FC 部署 |
| `ALIYUN_ACCESS_KEY_ID` | 阿里云 AccessKey ID | `LTAI5t...` | ✅ | FC 部署 |
| `ALIYUN_ACCESS_KEY_SECRET` | 阿里云 AccessKey Secret | `xxxxxxxxxxxxx` | ✅ | FC 部署 |
| `DATABASE_URL` | Neon PostgreSQL 连接字符串 | `postgres://user:pass@host:5432/db?sslmode=require` | ✅ | 数据库连接 |
| `PINECONE_API_KEY` | Pinecone API Key | `pcsk_xxxxxxxxxxxxx` | ✅ | 向量数据库 |
| `DEEPSEEK_API_KEY` | DeepSeek API Key | `sk-xxxxxxxxxxxxx` | ✅ | LLM 调用 |
| `OPENAI_API_KEY` | OpenAI API Key (可选) | `sk-xxxxxxxxxxxxx` | ❌ | 备用 LLM |
| `JWT_SECRET` | JWT 签名密钥 | `your-super-secret-key` | ✅ | 用户认证 |
| `CODECOV_TOKEN` | Codecov Token (可选) | `xxxxxxxxxxxxx` | ❌ | 测试覆盖率 |

### 如何获取

#### 阿里云密钥
1. 访问 https://ram.console.aliyun.com/
2. 创建 AccessKey 或使用现有的
3. 账号 ID 在右上角头像处查看

#### Neon PostgreSQL
1. 访问 https://console.neon.tech/
2. 选择你的项目
3. 复制 Connection String

#### Pinecone
1. 访问 https://app.pinecone.io/
2. 进入 API Keys 页面
3. 复制 API Key

#### DeepSeek
1. 访问 https://platform.deepseek.com/
2. 进入 API Keys 页面
3. 创建并复制 API Key

#### Codecov (可选)
1. 访问 https://about.codecov.io/
2. 注册并添加 GitHub 项目
3. 获取 Upload Token

---

## ⚙️ Variables (普通变量)

这些是非敏感配置，可以在工作流中直接使用。

| 名称 | 说明 | 示例值 | 必需 | 用途 |
|------|------|--------|------|------|
| `LLM_PROVIDER` | LLM 提供商 | `deepseek` 或 `openai` | ✅ | 选择 LLM |
| `LLM_BASE_URL` | LLM API 地址 | `https://api.deepseek.com` | ✅ | API 端点 |
| `LLM_MODEL` | 模型名称 | `deepseek-chat` | ✅ | 使用的模型 |
| `ALIYUN_REGION` | 阿里云区域 | `ap-southeast-1` | ✅ | FC 区域 |
| `FC_SERVICE_NAME` | FC 服务名称 | `ai-novel-service` | ❌ | 服务名 |
| `FC_FUNCTION_NAME` | FC 函数名称 | `ai-novel-backend` | ❌ | 函数名 |

---

## 🌍 Environments (环境)

建议配置两个环境：`production` 和 `staging`

### 配置位置
**Settings** → **Environments** → 点击 **New environment**

### production 环境
| 类型 | 名称 | 值 |
|------|------|-----|
| Secret | `DATABASE_URL` | 生产数据库连接 |
| Secret | `PINECONE_API_KEY` | 生产 Pinecone Key |
| Secret | `DEEPSEEK_API_KEY` | 生产 DeepSeek Key |
| Variable | `ALIYUN_REGION` | `ap-southeast-1` |

### staging 环境 (可选)
| 类型 | 名称 | 值 |
|------|------|-----|
| Secret | `DATABASE_URL` | 测试数据库连接 |
| Secret | `PINECONE_API_KEY` | 测试 Pinecone Key |
| Secret | `DEEPSEEK_API_KEY` | 测试 DeepSeek Key |

---

## 📋 快速配置清单

### 必需配置 (最小化部署)
- [ ] `ALIYUN_ACCOUNT_ID`
- [ ] `ALIYUN_ACCESS_KEY_ID`
- [ ] `ALIYUN_ACCESS_KEY_SECRET`
- [ ] `DATABASE_URL`
- [ ] `PINECONE_API_KEY`
- [ ] `DEEPSEEK_API_KEY`
- [ ] `JWT_SECRET`
- [ ] `LLM_PROVIDER` (Variable)
- [ ] `LLM_BASE_URL` (Variable)
- [ ] `LLM_MODEL` (Variable)

### 可选配置
- [ ] `OPENAI_API_KEY` (备用 LLM)
- [ ] `CODECOV_TOKEN` (测试覆盖率)

---

## 🔧 验证配置

配置完成后，可以通过以下方式验证：

### 1. 手动触发工作流
1. 进入 **Actions** 标签
2. 选择 **Deploy to Aliyun FC**
3. 点击 **Run workflow**
4. 选择分支和环境
5. 点击 **Run workflow**

### 2. 查看日志
- 检查 **Actions** 中的运行日志
- 确认所有步骤都成功
- 查看部署输出的域名

### 3. 测试 API
```bash
# 健康检查
curl https://your-fc-domain.com/health

# 创建测试小说
curl -X POST https://your-fc-domain.com/api/v1/novels \
  -H "Content-Type: application/json" \
  -d '{"title":"测试小说","description":"测试","genre":"玄幻"}'
```

---

## 🛡️ 安全建议

1. **定期轮换密钥** - 每 90 天更新一次 AccessKey 和 API Key
2. **使用最小权限** - 阿里云 RAM 用户只授予 FC 相关权限
3. **启用分支保护** - 保护 `main` 分支，要求 PR 审查
4. **限制环境变量访问** - 敏感 Secret 只在特定环境可用
5. **监控异常访问** - 开启 GitHub 安全告警

---

## 📝 本地测试

在本地测试时，复制 `.env.example` 为 `.env`：

```bash
cp .env.example .env
```

然后编辑 `.env` 填入你的配置：

```bash
# .env
DATABASE_URL=postgres://...
PINECONE_API_KEY=pcsk_...
DEEPSEEK_API_KEY=sk-...
JWT_SECRET=your-secret-key
```

---

**最后更新**: 2026-03-14  
**版本**: v1.0.0
