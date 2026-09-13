# 求职陪跑 AI 助手 · Job Copilot AI

面向实习和校招求职者的响应式 Web 应用。项目把个人求职档案、JD 匹配分析、模拟面试和历史复盘组织成连续工作流。

## 当前真实进度

已经完成：

- Vue 3 + TypeScript + Vite 前端工程及主要页面骨架；
- Supabase 邮箱注册、登录、Session 恢复、路由守卫和退出代码；
- Axios API Client；
- Go + Gin 后端基础服务；
- `GET /health` 与 `GET /api/v1/health`；
- Go Config、统一 JSON、404 Error Response 和 CORS；
- 前端健康检查已切换到 Go 的统一响应 Contract；
- `POST /api/v1/ai/analyze-jd` 的 Gin Handler、Supabase Bearer Token 校验、Application Service 和 OpenAI 兼容 AI Adapter；
- 8.22 阶段的八字段 Prompt、严格 JSON 解码、结构化结果校验和 AI 异常响应；
- 使用已配置的真实 LLM 完成过一次无敏感数据的结构化响应验证；
- `jd_analyses` migration、RLS、最小角色权限、Supabase Repository 和前端真实分析结果页；
- 使用演示数据完成登录、Go API、真实 LLM、结构校验、数据库写入和结果展示的端到端验收。
- `interview_sessions`、`interview_messages` migration 与 RLS，以及五轮面试的 Session、Message、Status 领域规则。
- `POST /api/v1/ai/interview/start` 的鉴权、JD 上下文读取、AI 第一题生成和会话启动持久化代码链路。
- `profiles` 建表、RLS、最小权限 migration，以及现有档案 Store 的保存、读取、修改和刷新恢复代码链路；
- JD 分析历史列表、单条详情 API 与真实数据页面；
- `POST /api/v1/ai/interview/turn` 的五轮回答校验、AI 评分反馈、前四轮下一题生成、重复提交保护和原子持久化代码链路；
- 面试会话读取与对话页刷新恢复代码链路，第 5 轮只保存回答与反馈，不生成第 6 题；
- 面试 AI 上下文只读取最近 8 条消息，并对 JD、档案、历史消息和当前回答分别限长。
- `POST /api/v1/ai/interview/report` 的九字段最终报告生成、严格结构校验、报告持久化和 Session 原子完成代码链路；
- 面试页的报告生成状态、四项评分、总结、优势、不足、推荐复习主题和回答建议展示，以及刷新恢复代码链路。

尚未真实完成：

- Supabase 真实凭据下的完整认证联调；
- 新增 migration 在目标 Supabase 项目的远端执行，以及真实账号下的 Profile、JD 历史、面试单轮完整联调；
- 面试历史列表；
- Dashboard 真实统计及部署验收。

## 技术栈

前端：Vue 3、TypeScript、Vite、Vue Router、Pinia、Element Plus、Axios、Supabase JS。

后端：Go、Gin。后续按实际业务接入 Supabase Auth、PostgreSQL 和 OpenAI 兼容 LLM API。

## 项目目录

```text
.
├─ frontend/                         Vue 前端
│  ├─ src/api/                       Axios 与 API 请求
│  ├─ src/components/                公共组件
│  ├─ src/router/                    路由与权限守卫
│  ├─ src/stores/                    Pinia 状态
│  └─ src/views/                     页面
├─ backend/                          Go + Gin 后端
│  ├─ cmd/server/main.go             启动入口
│  ├─ internal/domain/analysis/      JD 分析实体与业务规则
│  ├─ internal/domain/interview/     模拟面试实体、状态与五轮规则
│  ├─ internal/application/analysis/ JD 分析用例编排
│  ├─ internal/application/interview/ 模拟面试启动用例编排
│  ├─ internal/port/                 AI、认证与仓储接口
│  ├─ internal/adapter/              HTTP、AI 与认证适配器
│  ├─ internal/config/               环境配置
│  ├─ pkg/response/                  统一 JSON
│  ├─ go.mod
│  └─ .env.example
├─ docs/
├─ AGENTS.md
└─ .gitignore
```

## 环境要求

- Go 1.26.4 或与 `backend/go.mod` 兼容的版本；
- Node.js 20+；
- npm 10+；
- 真实注册登录需要 Supabase 项目。

## 本地启动

终端 1，启动 Go 后端：

```bash
cd backend
cp .env.example .env
go run ./cmd/server
```

终端 2，启动前端：

```bash
cd frontend
npm install
cp .env.example .env.local
npm run dev
```

访问：

- 前端：<http://localhost:5173>
- 基础健康检查：<http://localhost:8080/health>
- V1 API 健康检查：<http://localhost:8080/api/v1/health>

JD 分析接口：`POST http://localhost:8080/api/v1/ai/analyze-jd`（需要 Supabase access token）。AI 结果通过八字段校验后才会写入 `jd_analyses`。

启动模拟面试：`POST http://localhost:8080/api/v1/ai/interview/start`（需要 Supabase access token），请求体为：

```json
{
  "analysisId": "当前用户已有的 JD 分析 UUID"
}
```

成功响应的 `data` 包含 `sessionId`、`status`、`currentRound`、`maxRounds` 和 `question`。第一题和 Session 第 1 轮状态通过数据库函数原子保存；AI 失败时不会返回或保存伪造问题。

提交面试回答：`POST http://localhost:8080/api/v1/ai/interview/turn`（需要 Supabase access token），请求体为：

```json
{
  "sessionId": "当前用户进行中的面试会话 UUID",
  "answer": "本轮回答"
}
```

成功响应包含 `score`、`feedback`、`strengths`、`improvements`、`nextQuestion` 和推进后的轮次。第 1～4 轮会原子保存回答、反馈、下一题与新轮次；第 5 轮保存回答和反馈后 `nextQuestion` 为 `null`。相同轮次重复提交会返回冲突。`GET /api/v1/ai/interview/sessions/:id` 用于刷新恢复。

生成最终报告：`POST http://localhost:8080/api/v1/ai/interview/report`（需要 Supabase access token），请求体为：

```json
{
  "sessionId": "已完成五轮问答的面试会话 UUID"
}
```

成功响应包含 `overallScore`、`technicalScore`、`expressionScore`、`projectDepthScore`、`strengths`、`weaknesses`、`recommendedTopics`、`answerTips` 和 `summary`。只有五轮问题、回答和反馈均已持久化后才能生成；报告写入与 Session 标记为 `completed` 通过数据库函数原子完成。重复请求已完成的 Session 会返回已保存报告，不会再次调用模型。

JD 历史使用 `GET /api/v1/ai/jd-analyses` 和 `GET /api/v1/ai/jd-analyses/:id`。新环境需要按文件名顺序执行 `supabase/migrations/` 下的 SQL migration。当前尚未完成面试历史列表。

健康检查响应：

```json
{
  "code": "OK",
  "data": {
    "status": "ok"
  }
}
```

## 环境变量

前端：

```env
VITE_APP_NAME=求职陪跑 AI 助手
VITE_API_BASE_URL=http://localhost:8080/api/v1
VITE_SUPABASE_URL=
VITE_SUPABASE_ANON_KEY=
```

后端：

```env
APP_PORT=8080
APP_ENV=development
# 多个来源可使用英文逗号分隔；本地地址会自动兼容 localhost、127.0.0.1 和 ::1
FRONTEND_ORIGIN=http://localhost:5173
SUPABASE_URL=
SUPABASE_ANON_KEY=
AI_API_BASE_URL=
AI_API_KEY=
AI_MODEL=
AI_TIMEOUT_SECONDS=90
```

`AI_API_BASE_URL` 填 OpenAI 兼容 API 的 `/v1` 基础地址，后端会请求 `/chat/completions`。当认证或 AI 配置缺失时，健康检查仍可启动，业务接口会返回明确的未配置错误，不会生成假结果。

真实 `.env` 和 `.env.local` 不提交 Git。Supabase service role key、数据库密码和 AI API Key 只能放后端运行环境，不能放入浏览器代码。

## 检查命令

```bash
cd backend
gofmt -w ./cmd ./internal ./pkg
go test ./...
go build ./...

cd ../frontend
npm run typecheck
npm run build
```

## 认证边界

没有 Supabase URL 和 anon key 时，登录和注册页会提示配置缺失，不会创建本地假登录。配置后，前端会调用 Supabase Auth；当前仓库尚未使用真实凭据完成注册、登录、刷新恢复和退出的端到端联调。

后续用户业务数据必须使用经过验证的真实用户身份，并通过 RLS 隔离。不得信任前端提交的 `user_id`。

## 历史实现说明

仓库早期使用过 Express + TypeScript 健康检查骨架。2026-08-23 已完成基础后端向 Go + Gin 的迁移，当前 `backend/` 是唯一业务后端，不再维护 Express 业务逻辑。
