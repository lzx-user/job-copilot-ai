# 求职陪跑 AI 助手 · Job Copilot AI

面向实习和校招求职者的响应式 Web 应用。它把个人求职档案、JD 匹配分析、模拟面试和历史复盘组织成一个连续工作流。

## 当前完成阶段

本次只完成开发方案中的 2026-07-14、2026-07-15 和 2026-07-16：

- Vue 3 + TypeScript + Vite 前端工程；
- Express + TypeScript 后端工程；
- `/health` 与 `/api/v1/health` 健康检查；
- 浅色蓝紫 SaaS 设计系统、响应式侧栏、顶部栏和移动端抽屉；
- 仪表盘、JD 分析、模拟面试、历史记录、个人档案、设置页静态骨架；
- Supabase 邮箱注册、登录、刷新保持、路由守卫和退出登录代码；
- 入参校验、错误中文提示和缺少 Supabase 配置时的友好状态。

JD AI 分析、档案数据库、历史数据库、模拟面试 AI 和任何付费功能都没有提前实现。页面中的统计保持为 0，不使用虚假企业数据。

## 技术栈

前端：Vue 3、TypeScript、Vite、Vue Router、Pinia、Element Plus、Axios、Supabase JS。

后端：Node.js、Express、TypeScript、Zod、dotenv、cors、helmet、express-rate-limit。当前阶段后端不调用 AI。

## 项目目录

```text
.
├─ frontend/                 Vue 前端
│  ├─ src/api/               Axios 与 health 请求
│  ├─ src/components/        可复用布局与空状态组件
│  ├─ src/layouts/           AppLayout 公共登录后布局
│  ├─ src/lib/               Supabase 客户端
│  ├─ src/router/            懒加载路由与权限守卫
│  ├─ src/stores/            Pinia auth/ui 状态
│  ├─ src/styles/            设计变量与全局样式
│  └─ src/views/             登录、注册和业务骨架页
├─ backend/                  Express API
│  └─ src/                   环境、路由、中间件与 server
├─ docs/                     总方案、学习说明、配置步骤与参考图
├─ AGENTS.md                 长期开发规则
└─ .gitignore
```

## 环境要求

- Node.js 20+（当前验证环境为 Node.js 24）；
- npm 10+；
- 若要真实注册登录，需要一个 Supabase 项目。

## 安装与启动

先安装后端依赖：

```bash
cd backend
npm install
npm run dev
```

另开终端启动前端：

```bash
cd frontend
npm install
npm run dev
```

访问：

- 前端：<http://localhost:5173>
- 后端：<http://localhost:3000/health>
- API 健康检查：<http://localhost:3000/api/v1/health>

## 环境变量

复制示例文件：

```bash
copy frontend/.env.example frontend/.env.local
copy backend/.env.example backend/.env
```

前端 `.env.local` 至少填写：

```env
VITE_API_BASE_URL=http://localhost:3000/api/v1
VITE_SUPABASE_URL=
VITE_SUPABASE_ANON_KEY=
```

后端 `.env` 的 `FRONTEND_ORIGIN` 默认是 `http://localhost:5173`。不要把 `SUPABASE_SERVICE_ROLE_KEY`、`AI_API_KEY` 或真实 `.env` 提交到 Git。

详细的 Supabase 新手配置步骤见 [docs/Supabase配置步骤.md](docs/Supabase配置步骤.md)。

## 构建与类型检查

```bash
cd frontend
npm run typecheck
npm run build

cd ../backend
npm run typecheck
npm run build
```

## 认证行为

没有 Supabase URL 和 anon key 时，登录/注册按钮会禁用，并提示填写 `frontend/.env.local`。这不是本地假登录。配置后，页面会调用 Supabase 的 `signUp`、`signInWithPassword`、`getSession`、`onAuthStateChange` 和 `signOut`。

登录成功后只能访问 `/app/**`。刷新页面时先完成 Session 初始化，再执行路由守卫，避免短暂跳回登录页。登录页支持安全的内部 `redirect` 参数，不接受外部 URL。

## 当前已知边界

- 未提供 Supabase 凭据，因此真实注册、登录、刷新保持和退出尚未联调；
- 未接入 profiles、jd_analyses、interview_sessions、interview_messages 等业务表；
- 未接入 AI API；
- 参考图已复制到 `docs/reference/`，但没有把图片当作业务数据；
- 若不配置 Supabase，受保护业务页无法绕过认证直接访问，这是有意保留的安全行为。

## 后续开发顺序

7 月 17 日先做 `profiles` 表和 RLS，再做档案保存；随后做 JD 输入校验、后端 AI 接口、结果保存、历史记录、模拟面试和限流。每一步都应保留 loading、empty、error 和刷新恢复状态。

## 安全提醒

浏览器端 anon key 不是管理员密钥，权限由 Supabase Auth 与 RLS 决定；service role key 可以绕过 RLS，只能放后端环境变量。不要打印密码、access token、refresh token、完整 JD 或完整简历。

## 本地 Git

本次不会自动 commit 或 push。确认文件后可自行执行：

```bash
git init
git add .
git commit -m "feat: initialize job copilot web and api"
```

