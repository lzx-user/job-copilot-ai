# Supabase 配置步骤（新手版）

本项目没有提供真实 Supabase 密钥。请把密钥填写到本地文件，不要写进源码或提交 Git。

## 1. 创建项目

1. 打开 [supabase.com](https://supabase.com/) 并注册账号；
2. 点击 New project；
3. 选择组织，项目名可以填写 `job-copilot-ai`；
4. 生成一个新的数据库密码，保存到密码管理器；不要把它当成前端 anon key；
5. 等待项目初始化完成。

## 2. 开启邮箱认证

进入左侧 Authentication → Providers → Email：

- 确认 Email provider 已开启；
- 可以先保留 Confirm email 开启，这样注册后需要点击验证邮件；
- 如果关闭 Confirm email，注册成功后可能直接得到 Session。代码会根据 Supabase 实际返回判断，不假设两种配置相同。

不同版本控制台的菜单名称可能略有差异，寻找 Authentication、Providers、Email 这几个关键词即可。

## 3. 找到 URL 和 anon key

进入 Project Settings → API：

- Project URL 就是 `VITE_SUPABASE_URL`；
- `anon` / `public` key 就是 `VITE_SUPABASE_ANON_KEY`；
- 某些控制台称它为 Publishable key，确认它对应浏览器公开 key；
- `service_role` / Secret key 不是 anon key，永远不要复制到 frontend。

## 4. 配置前端

在 `frontend` 目录创建 `.env.local`（不要改 `.env.example`）：

```env
VITE_APP_NAME=求职陪跑 AI 助手
VITE_API_BASE_URL=http://localhost:3000/api/v1
VITE_SUPABASE_URL=https://你的项目-ref.supabase.co
VITE_SUPABASE_ANON_KEY=你的-anon-public-key
```

Vite 在启动时读取环境变量，修改后必须停止并重新运行 `npm run dev`。只刷新浏览器通常不会加载新变量。

## 5. 为什么 anon key 可以放前端

anon key 只标识项目，不代表管理员权限。真正的数据权限由登录身份和 PostgreSQL RLS 策略决定。后续创建 `profiles`、`jd_analyses` 等表时必须开启 RLS，并用 `auth.uid()` 限制只能访问当前用户。

service role key 能绕过 RLS，只能放在后端托管平台的环境变量中。本次前端代码不会读取它。

## 6. 启动并测试

终端 1：

```bash
cd backend
npm install
npm run dev
```

终端 2：

```bash
cd frontend
npm install
npm run dev
```

打开前端：

1. 进入注册页，填写合法邮箱和至少 8 位密码；
2. 如果开启邮箱验证，打开邮件并点击验证链接；
3. 回到登录页输入同一账号；
4. 登录成功后刷新页面，应仍在受保护页面；
5. 从右上角菜单或设置页退出；
6. 再次打开 `/app/dashboard`，应被送回登录页。

## 7. 常见问题

### 显示“Supabase 尚未配置”

确认文件名是 `frontend/.env.local`，变量名有 `VITE_` 前缀，没有多余引号，并重启 Vite。

### invalid login credentials

检查邮箱和密码，确认用户已注册；Supabase 的错误提示不会显示密码或 token。

### Email not confirmed

检查邮箱的验证链接，也可以在 Authentication → Users 查看用户状态。开发阶段邮件可能进入垃圾箱。

### 邮件收不到

检查垃圾邮件、邮箱地址是否正确、项目 Email provider 是否启用。Supabase 免费项目的邮件发送可能有频率限制，频繁测试时应稍等。

### CORS 或后端无法连接

确认后端正在 `3000` 端口运行，`backend/.env` 的 `FRONTEND_ORIGIN` 与浏览器地址完全一致，并检查登录页底部的后端服务状态。

### 认证状态没有刷新

停止 Vite 后重新启动，清除浏览器该站点的旧存储，再重新登录。不要用 localStorage 手工伪造 Session。

