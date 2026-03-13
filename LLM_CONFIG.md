# LLM 配置指南

本文档说明如何配置不同的 LLM 提供商。

## 🎯 支持的 LLM 提供商

| 提供商 | 模型 |  Base URL | 特点 |
|--------|------|-----------|------|
| **DeepSeek** | deepseek-chat, deepseek-coder | `https://api.deepseek.com` | 性价比高，中文优秀 |
| **OpenAI** | gpt-4o, gpt-4o-mini, gpt-4-turbo | `https://api.openai.com/v1` | 全球领先，多语言 |
| **MiniMax** | MiniMax-M2.5, MiniMax-Text-01 | `https://api.minimaxi.com/v1` | 国产大模型，长文本 |
| **智谱 AI** | glm-5, glm-4-air, glm-4-flash | `https://open.bigmodel.cn/api/paas/v4` | 国产大模型，多模态 |
| **月之暗面** | kimi-k2.5, kimi-latest, kimi-plus | `https://api.moonshot.cn/v1` | 超长上下文，Kimi |

---

## ⚙️ 配置方法

### 方法 1: 环境变量（推荐）

编辑 `.env` 文件：

```bash
# MiniMax M2.5
LLM_PROVIDER=minimax
LLM_BASE_URL=https://api.minimaxi.com/v1
LLM_MODEL=MiniMax-M2.5
MINIMAX_API_KEY=your-minimax-api-key

# 智谱 GLM-5
LLM_PROVIDER=zhipu
LLM_BASE_URL=https://open.bigmodel.cn/api/paas/v4
LLM_MODEL=glm-5
ZHIPU_API_KEY=your-zhipu-api-key

# 月之暗面 Kimi K2.5
LLM_PROVIDER=moonshot
LLM_BASE_URL=https://api.moonshot.cn/v1
LLM_MODEL=kimi-k2.5
MOONSHOT_API_KEY=your-moonshot-api-key
```

### 方法 2: 配置文件

编辑 `config.yaml`：

```yaml
llm:
  provider: "minimax"  # 或 zhipu, moonshot, deepseek, openai
  api_key: "your-api-key"
  base_url: "https://api.minimaxi.com/v1"
  model: "MiniMax-M2.5"
```

---

## 📋 各平台详细配置

### 1. MiniMax (M2.5)

**获取 API Key:**
1. 访问 https://platform.minimaxi.com/
2. 登录/注册账号
3. 进入 API 管理页面
4. 创建 API Key

**配置:**
```bash
LLM_PROVIDER=minimax
LLM_BASE_URL=https://api.minimaxi.com/v1
LLM_MODEL=MiniMax-M2.5
MINIMAX_API_KEY=your-minimax-api-key
```

**可用模型:**
- `MiniMax-M2.5` - 最新最强模型
- `MiniMax-Text-01` - 文本专用模型

**特点:**
- ✅ 支持超长文本（256K tokens）
- ✅ 中文理解能力强
- ✅ 适合小说创作

---

### 2. 智谱 AI (GLM-5 / GLM-4.7)

**获取 API Key:**
1. 访问 https://open.bigmodel.cn/
2. 登录/注册账号
3. 进入 API 控制台
4. 创建 API Key

**配置:**
```bash
LLM_PROVIDER=zhipu
LLM_BASE_URL=https://open.bigmodel.cn/api/paas/v4
LLM_MODEL=glm-5
ZHIPU_API_KEY=your-zhipu-api-key
```

**可用模型:**
- `glm-5` - 最新版本
- `glm-4-air` - 平衡版
- `glm-4-flash` - 快速版
- `glm-4-long` - 长文本版
- `glm-4v` - 视觉版

**特点:**
- ✅ 多模态支持
- ✅ 中文优化
- ✅ 性价比高

---

### 3. 月之暗面 (Kimi K2.5)

**获取 API Key:**
1. 访问 https://platform.moonshot.cn/
2. 登录/注册账号
3. 进入 API 管理
4. 创建 API Key

**配置:**
```bash
LLM_PROVIDER=moonshot
LLM_BASE_URL=https://api.moonshot.cn/v1
LLM_MODEL=kimi-k2.5
MOONSHOT_API_KEY=your-moonshot-api-key
```

**可用模型:**
- `kimi-k2.5` - 最新版本
- `kimi-latest` - 最新稳定版
- `kimi-plus` - 增强版
- `kimi-vision` - 视觉版

**特点:**
- ✅ 超长上下文（128K-256K tokens）
- ✅ 优秀的中文理解
- ✅ 适合长篇小说创作

---

### 4. DeepSeek (默认)

**获取 API Key:**
1. 访问 https://platform.deepseek.com/
2. 登录/注册账号
3. 进入 API 管理
4. 创建 API Key

**配置:**
```bash
LLM_PROVIDER=deepseek
LLM_BASE_URL=https://api.deepseek.com
LLM_MODEL=deepseek-chat
DEEPSEEK_API_KEY=your-deepseek-api-key
```

**可用模型:**
- `deepseek-chat` - 对话模型
- `deepseek-coder` - 代码模型

**特点:**
- ✅ 性价比极高
- ✅ 中文优秀
- ✅ 代码能力强

---

### 5. OpenAI

**获取 API Key:**
1. 访问 https://platform.openai.com/
2. 登录/注册账号
3. 进入 API Keys 页面
4. 创建 API Key

**配置:**
```bash
LLM_PROVIDER=openai
LLM_BASE_URL=https://api.openai.com/v1
LLM_MODEL=gpt-4o
OPENAI_API_KEY=your-openai-api-key
```

**可用模型:**
- `gpt-4o` - 最强模型
- `gpt-4o-mini` - 轻量版
- `gpt-4-turbo` - 快速版
- `o1-preview` - 推理模型

**特点:**
- ✅ 全球领先
- ✅ 多语言支持
- ✅ 生态完善

---

## 🔄 切换 LLM 提供商

### 使用环境变量
```bash
# 从 DeepSeek 切换到 MiniMax
export LLM_PROVIDER=minimax
export MINIMAX_API_KEY=your-key
export LLM_MODEL=MiniMax-M2.5

# 重启服务
go run cmd/server/main.go
```

### 使用配置文件
编辑 `config.yaml`，修改 `llm.provider` 和相关配置，然后重启服务。

---

## 💰 价格参考（2026 年）

| 提供商 | 模型 | 输入价格 | 输出价格 | 上下文 |
|--------|------|----------|----------|--------|
| DeepSeek | deepseek-chat | ¥0.002/1K | ¥0.008/1K | 64K |
| MiniMax | M2.5 | ¥0.001/1K | ¥0.004/1K | 256K |
| 智谱 AI | glm-5 | ¥0.005/1K | ¥0.015/1K | 128K |
| 月之暗面 | kimi-k2.5 | ¥0.002/1K | ¥0.006/1K | 256K |
| OpenAI | gpt-4o | $0.005/1K | $0.015/1K | 128K |

*价格为参考价，请以官方最新定价为准*

---

## 🎯 推荐场景

| 场景 | 推荐模型 | 理由 |
|------|----------|------|
| 网络小说创作 | Kimi K2.5 / MiniMax M2.5 | 超长上下文，适合长篇 |
| 代码生成 | DeepSeek Coder / GPT-4o | 代码能力强 |
| 日常对话 | DeepSeek / GLM-4 | 性价比高 |
| 多语言内容 | GPT-4o | 多语言支持好 |
| 预算有限 | DeepSeek / MiniMax | 价格低 |

---

## 🔧 验证配置

```bash
# 健康检查
curl http://localhost:8080/health

# 测试 AI 生成
curl -X POST http://localhost:8080/api/v1/novels/{id}/generate \
  -H "Content-Type: application/json" \
  -d '{"chapter_title":"第一章","outline":"主角登场..."}'
```

---

**最后更新**: 2026-03-14
