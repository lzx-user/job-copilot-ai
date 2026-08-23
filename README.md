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
- 前端健康检查已切换到 Go 的统一响应 Contract。

尚未真实完成：

- Supabase 真实凭据下的完整认证联调；
- `profiles` 表、RLS 和个人档案持久化；
- JD AI 分析与 `jd_analyses` 持久化；
- 五轮模拟面试、最终报告和历史恢复；
- Dashboard、History 真实数据及部署验收。

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
│  ├─ internal/config/               环境配置
│  ├─ internal/adapter/http/         Router、Handler、Middleware
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
FRONTEND_ORIGIN=http://localhost:5173
```

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
