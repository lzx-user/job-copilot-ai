# 求职陪跑 AI 助手 V1 开发总方案（可直接开工版）

> 版本：V1 Public Beta  
> 开发周期：2026-07-14 ～ 2026-07-31  
> 产品形态：响应式 Web 应用  
> 核心目标：真实用户可以注册、填写求职档案、分析 JD、完成多轮模拟面试，并在下次登录后查看历史记录。

---

## 0. 先确定：这次到底做什么

### 0.1 一句话定位

**求职陪跑 AI 助手**是一款面向实习和校招求职者的 AI 求职准备工具，围绕“个人求职档案 → JD 匹配分析 → 模拟面试 → 复盘记录”形成完整闭环。

### 0.2 V1 成功标准

7 月 31 日前必须达到：

- 有可公开访问的网址；
- 用户可以邮箱注册、登录和退出；
- 用户可以维护个人求职档案；
- 用户粘贴真实 JD 后，可以得到结构化匹配分析；
- 用户可以基于某份 JD 完成至少 5 轮文字模拟面试；
- JD 分析和面试记录保存在数据库中；
- 刷新页面、退出后重新登录，记录仍然存在；
- 用户只能查看和删除自己的数据；
- AI Key 不暴露在浏览器和 GitHub 中；
- 对 AI 调用次数、输入长度和异常情况有基本保护；
- 至少邀请 3 名真实同学完成一次完整试用。

### 0.3 V1 不做的内容

为了月底落地，以下功能全部延后：

- 微信小程序端；
- PDF 简历自动解析；
- 语音面试；
- AI 数字人；
- RAG 知识库；
- 向量数据库；
- 多模型自由切换；
- 付费系统；
- 社区和岗位爬虫；
- 独立管理后台；
- Web IDE 的大规模 TypeScript 重构；
- 自动化测试体系。

V1 的重点不是功能最多，而是**核心流程完整、数据能保存、服务能上线、真实用户能使用**。

---

# 1. 用户与使用场景

## 1.1 目标用户

- 正在找实习的本科生、研究生；
- 正在准备秋招、春招的应届生；
- 不知道如何判断岗位匹配度的求职者；
- 缺少真实面试机会，希望提前练习的求职者；
- 投递很多岗位，但复盘记录零散的求职者。

## 1.2 核心痛点

1. JD 很长，不知道企业真正重视哪些能力；
2. 不知道自己的项目和岗位要求是否匹配；
3. 打招呼内容过于模板化；
4. 面试前不知道复习哪些知识；
5. 很少获得真实面试机会，项目讲述缺少练习；
6. 每次使用通用聊天工具都需要重新粘贴个人背景；
7. 分析和面试记录无法统一沉淀。

## 1.3 核心用户流程

```text
注册 / 登录
  ↓
完善个人求职档案
  ↓
粘贴公司、岗位和 JD
  ↓
获得结构化 JD 匹配报告
  ↓
保存报告
  ↓
基于该岗位开始模拟面试
  ↓
AI 提问 → 用户回答 → AI 反馈与追问
  ↓
结束面试并生成总结
  ↓
在历史记录中复盘
```

---

# 2. V1 功能范围

## 2.1 P0：必须完成

| 模块 | 功能 |
|---|---|
| 注册登录 | 邮箱注册、登录、退出、登录态保持 |
| 个人档案 | 目标岗位、城市、技能、项目摘要、个人优势、到岗时间 |
| JD 分析 | 输入岗位信息、AI 分析、结构化结果、保存记录 |
| 模拟面试 | 选择模式、AI 提问、多轮回答、逐轮反馈、结束总结 |
| 历史记录 | 查看 JD 分析和面试记录、详情、删除 |
| 仪表盘 | 数据统计、最近记录、快捷入口、档案完善提醒 |
| 权限安全 | 用户数据隔离、AI Key 保护、请求鉴权 |
| 稳定性 | 加载态、失败态、输入限制、调用次数限制 |
| 部署 | 前端、后端、数据库全部在线 |

## 2.2 P1：核心完成后再做

- JD 分析报告复制；
- 打招呼话术一键复制；
- 根据 JD 直接跳转模拟面试；
- 个人档案完整度；
- 面试进度条；
- 历史记录筛选；
- 移动端响应式适配；
- 用户反馈入口。

## 2.3 P2：延期到 V1.1

- 分析报告导出 PDF；
- 自定义今日计划；
- 近期面试提醒；
- PDF 简历上传；
- 语音输入；
- 语音面试；
- 面试题收藏；
- 复习知识点清单；
- 多模型切换。

---

# 3. 技术架构

## 3.1 最终技术栈

### 前端

```text
Vue 3
TypeScript
Vite
Vue Router
Pinia
Element Plus
Axios
Supabase JavaScript Client
```

### 后端

```text
Node.js
Express
TypeScript
Zod
dotenv
cors
helmet
express-rate-limit
OpenAI 兼容 SDK 或模型厂商 SDK
Supabase JavaScript Client
```

### 数据与鉴权

```text
Supabase Auth
Supabase PostgreSQL
Supabase Row Level Security（RLS）
```

### 部署

```text
前端：Vercel
后端：Render 或其他 Node.js 托管服务
数据库与登录：Supabase
```

## 3.2 系统架构

```text
┌──────────────────────────────────────┐
│              用户浏览器              │
│ Vue3 + TS + Pinia + Element Plus     │
└─────────────┬───────────────┬────────┘
              │               │
      登录与业务数据          │ AI 请求
              │               │ Authorization: Bearer <token>
              ▼               ▼
┌──────────────────┐  ┌──────────────────────────┐
│ Supabase          │  │ Express AI 服务          │
│ Auth + PostgreSQL │  │ 鉴权 / 校验 / 限流       │
│ RLS 数据隔离      │  │ Prompt / 模型调用 / 解析 │
└──────────────────┘  └─────────────┬────────────┘
                                    │
                                    ▼
                         ┌────────────────────┐
                         │ 大模型 API          │
                         │ JSON 输出 / 多轮对话│
                         └────────────────────┘
```

## 3.3 职责划分

### 前端负责

- 页面展示和交互；
- Supabase 注册、登录和退出；
- 读取与保存用户自己的业务数据；
- 获取登录 token；
- 调用 Express AI 接口；
- 展示 AI 返回结果；
- 保存分析记录和面试消息；
- 处理加载、空状态和错误提示。

### Express 后端负责

- 校验 Supabase 登录 token；
- 校验用户输入；
- 检查每日调用次数；
- 组装 Prompt；
- 调用模型；
- 校验模型返回 JSON；
- 屏蔽模型服务错误；
- 返回统一业务结构；
- 不向浏览器暴露模型 Key。

### Supabase 负责

- 用户注册和登录；
- PostgreSQL 数据存储；
- 使用 RLS 保证用户只能访问自己的记录。

---

# 4. 页面与路由设计

## 4.1 路由表

```text
/                       根路径，按登录状态跳转
/auth/login             登录
/auth/register          注册
/app/dashboard          仪表盘
/app/jd-analysis        JD 分析
/app/interviews         模拟面试入口
/app/interviews/:id     具体面试会话
/app/history            历史记录
/app/profile            个人档案
/app/settings           设置
```

## 4.2 路由权限

### 公开页面

- `/auth/login`
- `/auth/register`

### 登录后页面

- `/app/**`

### 守卫逻辑

```text
进入 /app/**：
  没有 Supabase Session → 跳转 /auth/login
  有 Session → 放行

进入 /auth/login：
  已登录 → 跳转 /app/dashboard
  未登录 → 放行
```

---

# 5. 页面详细方案

## 5.1 登录页

### 页面元素

- 产品 Logo 和名称；
- 邮箱输入；
- 密码输入；
- 登录按钮；
- 跳转注册；
- 错误提示；
- 加载状态；
- 简短产品价值说明。

### 验收

- 输入错误时有明确提示；
- 登录成功跳转仪表盘；
- 刷新后保持登录；
- 登录按钮请求中不可重复点击。

---

## 5.2 注册页

### 页面元素

- 邮箱；
- 密码；
- 确认密码；
- 注册按钮；
- 用户协议和隐私说明简短提示；
- 跳转登录。

### 验收

- 密码至少 8 位；
- 两次密码不一致时前端拦截；
- 注册成功后显示邮箱验证提示，或按项目配置直接进入系统；
- 注册后引导进入个人档案页面。

---

## 5.3 仪表盘

参考已确定的仪表盘样图，但 V1 对功能做减法。

### 顶部区域

- 欢迎语；
- 搜索框先保留视觉，不实现全局搜索；
- 通知图标先保留视觉；
- 用户头像与退出菜单。

### 统计卡片

- 已分析 JD 数；
- 已完成模拟面试数；
- 本周新增分析数；
- 档案完整度。

### 快捷操作

- 开始 JD 分析；
- 开始模拟面试。

### 最近记录

展示最近 3 条：

- 公司；
- 岗位；
- 类型；
- 匹配分数或面试分数；
- 时间；
- 查看详情。

### 今日计划

V1 不做用户自定义任务，采用规则生成：

```text
档案不完整 → 完善个人档案
今天无 JD 分析 → 分析一个目标岗位
今天无面试记录 → 完成一次模拟面试
有未复盘面试 → 查看面试总结
```

### 验收

- 数据来自真实数据库；
- 无数据时展示引导空状态；
- 快捷入口可正常跳转。

---

## 5.4 个人档案

参考已生成的个人档案样图。

### 字段

```ts
interface CandidateProfile {
  nickname: string
  targetRoles: string[]
  expectedCities: string[]
  skills: string[]
  projectSummary: string
  strengths: string
  availability: string
  graduationYear?: string
}
```

### 区域

1. 个人摘要卡；
2. 基础信息；
3. 技术栈标签；
4. 项目经历摘要；
5. 个人优势；
6. 档案完整度；
7. AI 档案建议（P1）。

### 档案完整度规则

```text
昵称：10%
目标岗位：15%
期望城市：10%
技能不少于 5 个：20%
项目摘要不少于 100 字：25%
个人优势不少于 50 字：10%
到岗信息：10%
```

### 保存规则

- 必填：目标岗位、技能、项目摘要；
- 项目摘要最多 5000 字；
- 技能最多 30 个；
- 保存成功显示提示；
- 更新 `updated_at`。

### 验收

- 首次用户可以创建档案；
- 再次进入可以编辑；
- 刷新后内容不丢失；
- 档案内容会自动带入 JD 分析和模拟面试。

---

## 5.5 JD 智能匹配分析

### 左侧输入

- 公司名称；
- 岗位名称；
- JD 内容；
- 简历摘要：默认使用个人档案项目摘要，可临时修改；
- 技术栈：默认使用个人档案技能；
- 开始分析按钮。

### 输入限制

```text
公司名称：1～100 字
岗位名称：1～100 字
JD：200～8000 字
简历摘要：50～5000 字
技能：最多 30 个
```

### 分析中状态

- 按钮禁用；
- 显示 Skeleton；
- 文案：“AI 正在拆解岗位要求并匹配你的经历，通常需要数十秒”；
- 禁止重复提交。

### 结果结构

```ts
interface JdAnalysisResult {
  matchScore: number
  jobSummary: string
  coreRequirements: string[]
  matchedSkills: Array<{
    skill: string
    evidence: string
  }>
  missingSkills: Array<{
    skill: string
    priority: 'high' | 'medium' | 'low'
    suggestion: string
  }>
  resumeSuggestions: string[]
  preparationTopics: string[]
  greetingMessage: string
  riskNotice: string
}
```

### 结果页面区域

1. 匹配分数；
2. 岗位核心要求；
3. 已匹配技能；
4. 缺失技能；
5. 简历优化建议；
6. 面试复习重点；
7. 打招呼话术；
8. “开始模拟面试”按钮；
9. AI 结果免责声明。

### 匹配分数提示

页面固定显示：

> 匹配分数由 AI 根据用户输入内容生成，仅用于求职准备参考，不代表企业真实筛选结果。

### 保存

分析成功后保存：

- 原始公司；
- 岗位；
- JD；
- 本次使用的简历摘要；
- 技能；
- 完整分析结果；
- 创建时间。

### 验收

- 真实 JD 可分析；
- 结果完整显示；
- AI 返回异常格式时页面不会崩溃；
- 分析结果可在历史记录查看；
- 可直接基于该分析启动模拟面试。

---

## 5.6 模拟面试入口

### 页面内容

- 选择一条历史 JD；
- 或快速输入公司、岗位和 JD；
- 选择面试模式：
  - 技术面；
  - 项目深挖；
  - 综合面；
- 设置轮数：V1 固定 5 轮；
- 开始面试。

### 建议

优先让用户从 JD 分析结果进入面试，减少重复输入。

---

## 5.7 模拟面试会话页

参考已生成的模拟面试样图。

### 顶部信息

- 公司；
- 岗位；
- 面试模式；
- 当前轮次；
- 进度条；
- 结束面试。

### 中间聊天区

- AI 面试官消息；
- 用户消息；
- 当前问题；
- 输入框；
- 发送按钮。

### 右侧反馈区

每轮 AI 返回：

```ts
interface InterviewTurnResult {
  score: number
  feedback: string
  strengths: string[]
  improvements: string[]
  nextQuestion: string | null
}
```

### 交互规则

1. 创建会话；
2. AI 生成第一题；
3. 用户提交回答；
4. AI 评价并生成下一题；
5. 保存用户回答、反馈和下一题；
6. 达到第 5 轮，提示结束；
7. 用户可提前结束；
8. 结束后生成最终报告。

### 输入限制

```text
单次回答：1～4000 字
请求中不可再次发送
最近上下文：最多传最近 8～10 条消息
```

### 刷新恢复

会话和消息全部保存在数据库。

刷新后：

- 重新读取会话；
- 加载历史消息；
- 恢复当前轮次；
- 如果状态为 `completed`，展示最终报告。

### 验收

- 至少完成 5 轮；
- 每轮有分数、优点、改进点和追问；
- 刷新不丢消息；
- 可提前结束；
- 结束后有总结报告。

---

## 5.8 面试最终报告

```ts
interface InterviewFinalReport {
  overallScore: number
  technicalScore: number
  expressionScore: number
  projectDepthScore: number
  strengths: string[]
  weaknesses: string[]
  recommendedTopics: string[]
  answerTips: string[]
  summary: string
}
```

页面展示：

- 综合评分；
- 分项评分；
- 主要优点；
- 主要薄弱点；
- 下一步复习建议；
- 项目讲述建议；
- 返回历史记录。

---

## 5.9 历史记录

参考已生成的历史记录样图。

### 分类

- 全部；
- JD 分析；
- 模拟面试。

### 列表字段

- 公司；
- 岗位；
- 类型；
- 分数；
- 状态；
- 创建时间；
- 操作。

### 操作

- 查看详情；
- 基于分析开始面试；
- 删除。

### 删除规则

- 删除前二次确认；
- 删除面试会话时，同时删除对应消息；
- 不允许删除他人记录。

### V1 筛选

- 类型；
- 公司或岗位关键词；
- 时间倒序。

---

## 5.10 设置页

V1 只实现：

- 退出登录；
- 数据与隐私说明；
- AI 使用说明；
- 删除账号入口暂时展示“联系删除”或延后；
- 模型参数不开放给普通用户。

---

# 6. 视觉与样式规范

## 6.1 设计方向

严格沿用已确定参考图：

- 明亮浅色 SaaS 风格；
- 白色卡片；
- 蓝紫渐变；
- 轻边框和柔和阴影；
- 大圆角；
- 清晰信息层级；
- 不堆砌高饱和颜色；
- 页面功能多，但保持留白。

## 6.2 CSS 变量

```css
:root {
  --color-primary: #4169f6;
  --color-primary-hover: #3158df;
  --color-secondary: #7c4dff;
  --color-success: #20b878;
  --color-warning: #f5a623;
  --color-danger: #ef5350;

  --text-primary: #17213c;
  --text-regular: #4f5d78;
  --text-secondary: #8791a8;
  --border-color: #e8ecf5;
  --background-page: #f7f9ff;
  --background-card: #ffffff;

  --radius-sm: 8px;
  --radius-md: 12px;
  --radius-lg: 18px;

  --shadow-card: 0 8px 30px rgba(45, 72, 130, 0.07);
  --shadow-hover: 0 12px 36px rgba(45, 72, 130, 0.12);

  --sidebar-width: 240px;
  --topbar-height: 68px;
}
```

## 6.3 页面尺寸

```text
桌面主设计宽度：1440px
侧边栏：240px
主内容最大宽度：1440px 内自适应
卡片间距：16～20px
页面左右边距：28～32px
卡片圆角：14～18px
按钮高度：40～44px
```

## 6.4 响应式规则

### ≥ 1200px

- 完整侧边栏；
- JD 页面左右两栏；
- 面试页聊天 + 反馈双栏。

### 768～1199px

- 侧边栏收缩；
- JD 结果区改为上下布局；
- 面试反馈区放在聊天区下方。

### < 768px

- 侧边栏改抽屉；
- 统计卡片两列；
- 所有表单单列；
- 历史记录使用卡片而非宽表格；
- 输入框固定在底部时避免遮挡。

## 6.5 Element Plus 使用边界

直接使用：

- Form；
- Input；
- Select；
- Button；
- Tag；
- Dialog；
- Drawer；
- Table；
- Progress；
- Skeleton；
- Empty；
- Message；
- Dropdown。

需要自定义样式：

- 仪表盘统计卡片；
- JD 结果卡片；
- 聊天气泡；
- 面试反馈面板；
- 页面 Hero；
- 侧边栏选中态；
- 渐变主按钮。

---

# 7. 前端工程设计

## 7.1 推荐目录

```text
job-copilot-web/
├─ public/
├─ src/
│  ├─ api/
│  │  ├─ request.ts
│  │  ├─ ai.api.ts
│  │  ├─ profile.api.ts
│  │  ├─ analysis.api.ts
│  │  └─ interview.api.ts
│  ├─ assets/
│  │  ├─ icons/
│  │  └─ images/
│  ├─ components/
│  │  ├─ common/
│  │  │  ├─ AppCard.vue
│  │  │  ├─ AppEmpty.vue
│  │  │  ├─ AppLoading.vue
│  │  │  ├─ ScoreRing.vue
│  │  │  └─ SkillTags.vue
│  │  ├─ dashboard/
│  │  ├─ analysis/
│  │  ├─ interview/
│  │  ├─ history/
│  │  └─ profile/
│  ├─ layouts/
│  │  └─ AppLayout.vue
│  ├─ router/
│  │  └─ index.ts
│  ├─ stores/
│  │  ├─ auth.store.ts
│  │  ├─ profile.store.ts
│  │  ├─ analysis.store.ts
│  │  ├─ interview.store.ts
│  │  └─ ui.store.ts
│  ├─ styles/
│  │  ├─ variables.css
│  │  ├─ reset.css
│  │  ├─ element-overrides.css
│  │  └─ global.css
│  ├─ types/
│  │  ├─ api.ts
│  │  ├─ auth.ts
│  │  ├─ profile.ts
│  │  ├─ analysis.ts
│  │  ├─ interview.ts
│  │  └─ database.ts
│  ├─ utils/
│  │  ├─ format.ts
│  │  ├─ validators.ts
│  │  ├─ storage.ts
│  │  └─ error.ts
│  ├─ views/
│  │  ├─ auth/
│  │  ├─ dashboard/
│  │  ├─ jd-analysis/
│  │  ├─ interview/
│  │  ├─ history/
│  │  ├─ profile/
│  │  └─ settings/
│  ├─ lib/
│  │  └─ supabase.ts
│  ├─ App.vue
│  └─ main.ts
├─ .env.example
├─ package.json
└─ vite.config.ts
```

## 7.2 Pinia Store 职责

### authStore

```ts
session
user
initializing
initialize()
login()
register()
logout()
```

### profileStore

```ts
profile
loading
loaded
fetchProfile()
saveProfile()
calculateCompleteness()
```

### analysisStore

```ts
currentForm
currentResult
analyzing
history
analyze()
saveAnalysis()
fetchHistory()
reset()
```

### interviewStore

```ts
currentSession
messages
currentFeedback
sending
startInterview()
loadSession()
sendAnswer()
endInterview()
```

### uiStore

```ts
sidebarCollapsed
mobileDrawerVisible
globalLoading
```

## 7.3 Axios 实例

```ts
interface ApiResponse<T> {
  code: string
  message: string
  data: T
  requestId?: string
}
```

请求拦截器：

- 获取 Supabase access token；
- 添加 `Authorization: Bearer <token>`；
- 添加 `Content-Type: application/json`；
- 设置 45 秒超时。

响应拦截器：

- 统一读取 `code`；
- 401 时退出或刷新登录状态；
- 超时显示友好文案；
- 不直接展示服务器堆栈。

---

# 8. 后端工程设计

## 8.1 推荐目录

```text
job-copilot-server/
├─ src/
│  ├─ config/
│  │  ├─ env.ts
│  │  ├─ ai.ts
│  │  └─ supabase.ts
│  ├─ controllers/
│  │  ├─ analysis.controller.ts
│  │  └─ interview.controller.ts
│  ├─ middleware/
│  │  ├─ auth.middleware.ts
│  │  ├─ validate.middleware.ts
│  │  ├─ error.middleware.ts
│  │  ├─ request-id.middleware.ts
│  │  └─ rate-limit.middleware.ts
│  ├─ prompts/
│  │  ├─ analyze-jd.prompt.ts
│  │  ├─ interview-start.prompt.ts
│  │  ├─ interview-turn.prompt.ts
│  │  └─ interview-report.prompt.ts
│  ├─ routes/
│  │  ├─ health.routes.ts
│  │  └─ ai.routes.ts
│  ├─ schemas/
│  │  ├─ common.schema.ts
│  │  ├─ analysis.schema.ts
│  │  └─ interview.schema.ts
│  ├─ services/
│  │  ├─ ai.service.ts
│  │  ├─ usage.service.ts
│  │  ├─ analysis.service.ts
│  │  └─ interview.service.ts
│  ├─ types/
│  │  ├─ express.d.ts
│  │  └─ api.ts
│  ├─ utils/
│  │  ├─ json.ts
│  │  ├─ retry.ts
│  │  └─ logger.ts
│  ├─ app.ts
│  └─ server.ts
├─ .env.example
├─ package.json
└─ tsconfig.json
```

## 8.2 中间件顺序

```text
requestId
→ helmet
→ cors
→ express.json
→ 基础 IP 限流
→ 路由
→ 404
→ 全局错误处理
```

AI 路由内部：

```text
用户鉴权
→ Zod 参数校验
→ 用户每日配额校验
→ Controller
→ Service
```

---

# 9. 数据库设计

## 9.1 profiles

```sql
create table public.profiles (
  user_id uuid primary key references auth.users(id) on delete cascade,
  nickname text not null default '',
  target_roles text[] not null default '{}',
  expected_cities text[] not null default '{}',
  skills text[] not null default '{}',
  project_summary text not null default '',
  strengths text not null default '',
  availability text not null default '',
  graduation_year text,
  created_at timestamptz not null default now(),
  updated_at timestamptz not null default now()
);
```

## 9.2 jd_analyses

```sql
create table public.jd_analyses (
  id uuid primary key default gen_random_uuid(),
  user_id uuid not null references auth.users(id) on delete cascade,
  company_name text not null,
  job_title text not null,
  jd_content text not null,
  resume_summary text not null,
  skills text[] not null default '{}',
  match_score int not null check (match_score between 0 and 100),
  analysis_result jsonb not null,
  created_at timestamptz not null default now(),
  updated_at timestamptz not null default now()
);

create index jd_analyses_user_created_idx
on public.jd_analyses(user_id, created_at desc);
```

## 9.3 interview_sessions

```sql
create type public.interview_mode as enum (
  'technical',
  'project',
  'comprehensive'
);

create type public.interview_status as enum (
  'in_progress',
  'completed',
  'abandoned'
);

create table public.interview_sessions (
  id uuid primary key default gen_random_uuid(),
  user_id uuid not null references auth.users(id) on delete cascade,
  analysis_id uuid references public.jd_analyses(id) on delete set null,
  company_name text not null,
  job_title text not null,
  jd_content text not null,
  mode public.interview_mode not null,
  status public.interview_status not null default 'in_progress',
  current_round int not null default 1,
  max_rounds int not null default 5,
  overall_score int check (overall_score between 0 and 100),
  final_report jsonb,
  created_at timestamptz not null default now(),
  updated_at timestamptz not null default now(),
  completed_at timestamptz
);

create index interview_sessions_user_created_idx
on public.interview_sessions(user_id, created_at desc);
```

## 9.4 interview_messages

```sql
create type public.message_role as enum (
  'assistant',
  'user'
);

create table public.interview_messages (
  id uuid primary key default gen_random_uuid(),
  user_id uuid not null references auth.users(id) on delete cascade,
  session_id uuid not null references public.interview_sessions(id) on delete cascade,
  role public.message_role not null,
  content text not null,
  round_number int not null,
  feedback jsonb,
  created_at timestamptz not null default now()
);

create index interview_messages_session_created_idx
on public.interview_messages(session_id, created_at);
```

## 9.5 usage_records

V1 使用事件记录表，简单可靠。

```sql
create type public.usage_feature as enum (
  'jd_analysis',
  'interview_turn',
  'interview_report'
);

create table public.usage_records (
  id bigint generated always as identity primary key,
  user_id uuid not null references auth.users(id) on delete cascade,
  feature public.usage_feature not null,
  created_at timestamptz not null default now()
);

create index usage_records_user_feature_created_idx
on public.usage_records(user_id, feature, created_at desc);
```

---

# 10. RLS 权限策略

所有业务表开启 RLS：

```sql
alter table public.profiles enable row level security;
alter table public.jd_analyses enable row level security;
alter table public.interview_sessions enable row level security;
alter table public.interview_messages enable row level security;
alter table public.usage_records enable row level security;
```

## 10.1 profiles

```sql
create policy "users can read own profile"
on public.profiles for select
using (auth.uid() = user_id);

create policy "users can insert own profile"
on public.profiles for insert
with check (auth.uid() = user_id);

create policy "users can update own profile"
on public.profiles for update
using (auth.uid() = user_id)
with check (auth.uid() = user_id);
```

## 10.2 通用用户数据策略

对 `jd_analyses`、`interview_sessions`、`interview_messages` 分别创建：

```sql
create policy "users can select own rows"
on public.jd_analyses for select
using (auth.uid() = user_id);

create policy "users can insert own rows"
on public.jd_analyses for insert
with check (auth.uid() = user_id);

create policy "users can update own rows"
on public.jd_analyses for update
using (auth.uid() = user_id)
with check (auth.uid() = user_id);

create policy "users can delete own rows"
on public.jd_analyses for delete
using (auth.uid() = user_id);
```

将表名替换后，为两个面试表创建相同策略。

`usage_records` 不允许前端直接插入，建议仅由后端使用 service role 写入；前端也不需要读取。

---

# 11. API 设计

## 11.1 健康检查

```http
GET /health
```

返回：

```json
{
  "status": "ok",
  "timestamp": "2026-07-14T10:00:00.000Z"
}
```

## 11.2 JD 分析

```http
POST /api/v1/ai/analyze-jd
Authorization: Bearer <supabase_access_token>
Content-Type: application/json
```

请求：

```json
{
  "companyName": "某公司",
  "jobTitle": "前端开发实习生",
  "jd": "岗位职责与要求……",
  "profile": {
    "targetRoles": ["前端开发实习生"],
    "skills": ["React", "Vue3", "TypeScript"],
    "projectSummary": "项目经历摘要……",
    "strengths": "个人优势……"
  }
}
```

响应：

```json
{
  "code": "OK",
  "message": "分析完成",
  "data": {
    "matchScore": 78,
    "jobSummary": "该岗位主要负责……",
    "coreRequirements": ["Vue3 业务开发", "TypeScript"],
    "matchedSkills": [
      {
        "skill": "TypeScript",
        "evidence": "个人技能与项目描述中包含 TypeScript 实践"
      }
    ],
    "missingSkills": [
      {
        "skill": "uni-app",
        "priority": "medium",
        "suggestion": "补充一个移动端项目或说明相关学习进度"
      }
    ],
    "resumeSuggestions": ["……"],
    "preparationTopics": ["……"],
    "greetingMessage": "您好，我是……",
    "riskNotice": "结果仅供求职准备参考"
  }
}
```

## 11.3 开始面试

```http
POST /api/v1/ai/interview/start
```

请求：

```json
{
  "companyName": "某公司",
  "jobTitle": "前端开发实习生",
  "jd": "……",
  "mode": "technical",
  "profile": {
    "skills": ["React", "Vue3", "TypeScript"],
    "projectSummary": "……"
  }
}
```

响应：

```json
{
  "code": "OK",
  "message": "面试已开始",
  "data": {
    "firstQuestion": "请介绍一下你最有代表性的前端项目。",
    "opening": "你好，我将担任本次 AI 面试官。"
  }
}
```

## 11.4 提交一轮回答

```http
POST /api/v1/ai/interview/turn
```

请求：

```json
{
  "companyName": "某公司",
  "jobTitle": "前端开发实习生",
  "jd": "……",
  "mode": "technical",
  "round": 2,
  "profileSummary": "……",
  "history": [
    {
      "role": "assistant",
      "content": "请介绍一下你的 Web IDE 项目。"
    },
    {
      "role": "user",
      "content": "该项目主要用于……"
    }
  ],
  "answer": "我的回答……"
}
```

响应：

```json
{
  "code": "OK",
  "message": "本轮反馈已生成",
  "data": {
    "score": 82,
    "feedback": "回答覆盖了项目目标和主要技术……",
    "strengths": ["项目背景清晰", "技术选型有解释"],
    "improvements": ["缺少具体难点", "缺少量化结果"],
    "nextQuestion": "多人同时编辑时，你如何处理冲突？"
  }
}
```

## 11.5 结束面试

```http
POST /api/v1/ai/interview/report
```

请求：

```json
{
  "companyName": "某公司",
  "jobTitle": "前端开发实习生",
  "mode": "technical",
  "transcript": [
    {
      "question": "……",
      "answer": "……",
      "score": 82,
      "feedback": "……"
    }
  ]
}
```

响应：

```json
{
  "code": "OK",
  "message": "面试报告已生成",
  "data": {
    "overallScore": 80,
    "technicalScore": 82,
    "expressionScore": 78,
    "projectDepthScore": 79,
    "strengths": ["……"],
    "weaknesses": ["……"],
    "recommendedTopics": ["……"],
    "answerTips": ["……"],
    "summary": "……"
  }
}
```

---

# 12. 统一错误格式

```json
{
  "code": "MODEL_TIMEOUT",
  "message": "AI 服务响应超时，请稍后重试",
  "data": null,
  "requestId": "req_xxx"
}
```

## 12.1 错误码

| code | 页面提示 |
|---|---|
| AUTH_REQUIRED | 登录状态已失效，请重新登录 |
| VALIDATION_ERROR | 输入内容不完整或格式不正确 |
| PROFILE_REQUIRED | 请先完善个人档案 |
| LIMIT_EXCEEDED | 今日使用次数已达上限 |
| MODEL_TIMEOUT | AI 服务响应超时，请稍后重试 |
| MODEL_OUTPUT_INVALID | 本次分析结果生成异常，请重新分析 |
| MODEL_UNAVAILABLE | AI 服务暂时繁忙，请稍后再试 |
| DATABASE_ERROR | 数据保存失败，请稍后重试 |
| INTERNAL_ERROR | 服务出现异常，请稍后重试 |

---

# 13. AI Prompt 设计

## 13.1 JD 分析系统 Prompt

```text
你是一名谨慎、客观的校招求职分析助手。

你的任务是根据用户提供的岗位 JD 和个人求职档案生成求职准备建议。

要求：
1. 不得声称可以预测企业真实录用结果；
2. 匹配分数只表示文本层面的准备程度；
3. 所有判断必须尽量引用用户档案中的具体依据；
4. 不得捏造用户没有提供的经历；
5. 建议要适合实习生和应届生，不要套用社招高级岗位标准；
6. 返回严格 JSON，不输出 Markdown 代码块；
7. 数组没有内容时返回空数组，不得省略字段；
8. matchScore 必须是 0～100 的整数。
```

用户 Prompt 包含：

```text
公司名称
岗位名称
完整 JD
目标岗位
技能
项目摘要
个人优势
```

## 13.2 面试系统 Prompt

```text
你是一名严格但友好的前端/全栈实习面试官。

目标：
- 根据岗位 JD 和用户项目经历进行针对性提问；
- 每次只提出一个清晰问题；
- 用户回答后给出具体反馈，不要只说“很好”；
- 必须指出回答中已覆盖的点和缺少的点；
- 不得捏造用户项目细节；
- 追问优先围绕岗位要求和用户真实项目；
- 返回严格 JSON。
```

## 13.3 输出校验

后端必须使用 Zod 校验模型结果：

```ts
const analysisResultSchema = z.object({
  matchScore: z.number().int().min(0).max(100),
  jobSummary: z.string().min(1),
  coreRequirements: z.array(z.string()).max(10),
  matchedSkills: z.array(
    z.object({
      skill: z.string(),
      evidence: z.string()
    })
  ).max(15),
  missingSkills: z.array(
    z.object({
      skill: z.string(),
      priority: z.enum(['high', 'medium', 'low']),
      suggestion: z.string()
    })
  ).max(15),
  resumeSuggestions: z.array(z.string()).max(10),
  preparationTopics: z.array(z.string()).max(12),
  greetingMessage: z.string(),
  riskNotice: z.string()
})
```

解析失败处理：

```text
第一次输出校验失败
→ 使用“修复为指定 JSON 结构”的 Prompt 重试一次
→ 第二次仍失败
→ 返回 MODEL_OUTPUT_INVALID
```

---

# 14. 鉴权与安全

## 14.1 前端可公开的内容

```env
VITE_SUPABASE_URL=
VITE_SUPABASE_ANON_KEY=
VITE_API_BASE_URL=
```

Supabase anon key 可以用于前端，但前提是所有表正确启用 RLS。

## 14.2 绝不能放前端的内容

```env
AI_API_KEY=
SUPABASE_SERVICE_ROLE_KEY=
```

它们只能存在后端环境变量中。

## 14.3 后端鉴权

前端调用 AI API 时：

```http
Authorization: Bearer <supabase_access_token>
```

Express 中间件：

1. 读取 Bearer token；
2. 调用 Supabase Auth 校验用户；
3. 得到 `user.id`；
4. 写入 `req.user`；
5. 后续限流和日志均使用真实 user id。

## 14.4 日志规范

允许记录：

- requestId；
- userId 的部分脱敏值；
- 接口名；
- 耗时；
- 状态；
- 模型错误码。

不要记录：

- 完整简历；
- 完整 JD；
- 用户密码；
- access token；
- AI API Key；
- 完整聊天内容。

---

# 15. AI 使用次数限制

Public Beta 建议：

```text
每个用户每天：
JD 分析：5 次
面试对话：20 轮
面试总结：3 次
```

后端流程：

```text
鉴权
→ 查询今天 usage_records 中该功能数量
→ 达到上限：返回 LIMIT_EXCEEDED
→ 未达到：调用模型
→ 模型成功后插入 usage_records
```

模型失败或超时，不计入成功次数。

同时增加 IP 级别基础限流，防止恶意请求。

---

# 16. 加载、空状态和失败状态

## 16.1 必须有的加载态

- 登录；
- 保存档案；
- JD 分析；
- 开始面试；
- 提交回答；
- 生成最终报告；
- 历史记录加载。

## 16.2 必须有的空状态

- 没有档案；
- 没有 JD 分析；
- 没有模拟面试；
- 搜索无结果；
- 分析记录被删除。

## 16.3 必须有的失败态

- 网络断开；
- AI 超时；
- 登录失效；
- 数据保存失败；
- 模型返回格式错误；
- 使用次数到达上限。

---

# 17. 环境变量

## 17.1 前端 `.env.example`

```env
VITE_APP_NAME=求职陪跑 AI 助手
VITE_API_BASE_URL=http://localhost:3000/api/v1
VITE_SUPABASE_URL=
VITE_SUPABASE_ANON_KEY=
```

## 17.2 后端 `.env.example`

```env
NODE_ENV=development
PORT=3000
FRONTEND_ORIGIN=http://localhost:5173

SUPABASE_URL=
SUPABASE_ANON_KEY=
SUPABASE_SERVICE_ROLE_KEY=

AI_API_KEY=
AI_BASE_URL=
AI_MODEL=

AI_TIMEOUT_MS=45000
JD_DAILY_LIMIT=5
INTERVIEW_TURN_DAILY_LIMIT=20
INTERVIEW_REPORT_DAILY_LIMIT=3
```

---

# 18. 初始化命令

## 18.1 前端

```bash
npm create vite@latest job-copilot-web -- --template vue-ts
cd job-copilot-web

npm install
npm install vue-router pinia axios element-plus
npm install @element-plus/icons-vue @supabase/supabase-js
```

## 18.2 后端

```bash
mkdir job-copilot-server
cd job-copilot-server

npm init -y
npm install express cors dotenv zod helmet express-rate-limit
npm install @supabase/supabase-js openai
npm install -D typescript tsx @types/node @types/express @types/cors
npx tsc --init
```

## 18.3 后端脚本

```json
{
  "scripts": {
    "dev": "tsx watch src/server.ts",
    "build": "tsc",
    "start": "node dist/server.js",
    "typecheck": "tsc --noEmit"
  }
}
```

---

# 19. 开发顺序与依赖关系

必须按下面顺序，不要看到哪页好看就先做哪页。

```text
项目骨架
→ AppLayout 和设计变量
→ Supabase 注册登录
→ 个人档案
→ JD 输入页静态 UI
→ JD AI 接口
→ JD 结果展示
→ 保存与历史记录
→ 模拟面试数据结构
→ 面试 AI 接口
→ 面试聊天 UI
→ 面试最终报告
→ 仪表盘真实统计
→ 限流和异常处理
→ 部署
→ 真实用户试用
```

原因：

- 个人档案是 JD 和面试的上下文；
- JD 记录是模拟面试的入口；
- 历史记录依赖前两个核心功能；
- 仪表盘依赖已有真实数据；
- 先做仪表盘会产生大量假数据和返工。

---

# 20. 7 月 14 日～7 月 31 日逐日计划

## 7 月 14 日：项目初始化

完成：

- 建立前端和后端仓库；
- 初始化 Vue3 + TS + Vite；
- 初始化 Express + TS；
- 安装依赖；
- 建立目录；
- 配置 ESLint 可选，优先不折腾复杂规则；
- 写 `.env.example`；
- 首次 Git 提交。

验收：

```text
前端 localhost:5173 正常打开
后端 localhost:3000/health 返回 ok
前端可调用 health 接口
```

## 7 月 15 日：整体布局与设计系统

完成：

- `variables.css`；
- `AppLayout`；
- 左侧导航；
- 顶部栏；
- 路由；
- 仪表盘、JD、面试、历史、档案空页面；
- 响应式基础。

验收：

```text
所有路由可以跳转
布局与参考图风格统一
浏览器缩放时不明显错位
```

## 7 月 16 日：Supabase 与登录

完成：

- 创建 Supabase 项目；
- 配置邮箱 Auth；
- 创建 `supabase.ts`；
- 登录页；
- 注册页；
- authStore；
- 路由守卫；
- 退出登录。

验收：

```text
可以真实注册
可以真实登录
刷新保持登录
未登录不能访问 /app/**
```

## 7 月 17 日：个人档案数据库与页面

完成：

- 建 `profiles` 表；
- 开启 RLS；
- 写权限策略；
- 按参考图实现档案页；
- 首次保存和读取；
- 计算档案完整度。

验收：

```text
用户 A 看不到用户 B 档案
保存后刷新内容仍在
```

## 7 月 18 日：JD 分析输入页

完成：

- 公司和岗位表单；
- JD 输入框；
- 自动带入档案摘要；
- 技能标签；
- 表单校验；
- 分析结果页面静态结构；
- Skeleton。

验收：

```text
不合法输入无法提交
页面结构与参考图接近
```

## 7 月 19 日：JD AI 后端

完成：

- `/ai/analyze-jd`；
- 鉴权中间件；
- Zod 请求校验；
- Prompt；
- 模型调用；
- JSON 校验；
- 错误格式；
- 前后端联调。

验收：

```text
真实 JD 可返回完整 JSON
AI Key 不出现在前端
错误 JSON 不会导致服务器崩溃
```

## 7 月 20 日：JD 结果与保存

完成：

- 结果卡片；
- 分数环；
- 已匹配/缺失技能；
- 建议；
- 打招呼复制；
- 建 `jd_analyses` 表；
- RLS；
- 成功后保存记录。

验收：

```text
一条真实分析完整显示
数据库中存在对应记录
重新登录后仍能查看
```

## 7 月 21 日：历史记录第一版

完成：

- 分析记录列表；
- 详情 Drawer；
- 删除；
- 搜索公司/岗位；
- 空状态。

验收：

```text
只能看到自己的记录
删除后列表同步更新
```

## 7 月 22 日：模拟面试会话结构

完成：

- 建 `interview_sessions`；
- 建 `interview_messages`；
- RLS；
- 面试入口页；
- 模式选择；
- 基于 JD 记录创建会话；
- `/interview/start`。

验收：

```text
可以选择一条 JD 开始面试
数据库生成 session
AI 返回第一题
```

## 7 月 23 日：聊天页面和单轮反馈

完成：

- 面试页面布局；
- 聊天气泡；
- 输入框；
- 提交回答；
- `/interview/turn`；
- 右侧评分和反馈；
- 保存消息。

验收：

```text
至少完成 3 轮连续对话
每轮消息和反馈写入数据库
```

## 7 月 24 日：会话恢复与 5 轮流程

完成：

- 刷新恢复；
- 当前轮次；
- 进度条；
- 发送防重复；
- 最近消息上下文；
- 达到 5 轮引导结束。

验收：

```text
刷新后继续面试
不会重复发送
第 5 轮后可结束
```

## 7 月 25 日：最终面试报告

完成：

- `/interview/report`；
- 最终报告结构；
- 分项评分；
- 优点；
- 薄弱点；
- 复习建议；
- 会话状态更新为 completed。

验收：

```text
完整面试生成报告
历史记录可打开报告
```

## 7 月 26 日：仪表盘真实数据

完成：

- 统计卡片；
- 最近记录；
- 快捷入口；
- 规则生成今日计划；
- 无数据引导。

验收：

```text
统计数字来自数据库
新用户和已有数据用户都能正常显示
```

## 7 月 27 日：交互与异常加固

完成：

- 全部 loading；
- 全部 empty；
- 超时；
- 登录失效；
- AI JSON 失败重试一次；
- 输入长度限制；
- 统一错误提示；
- usage_records；
- 每日调用限制。

验收：

```text
断网、超时、达到限额时页面不崩
错误提示用户看得懂
```

## 7 月 28 日：部署

完成：

- Supabase 正式环境；
- 后端部署；
- 前端部署；
- CORS；
- 环境变量；
- 线上注册；
- 线上 AI 调用；
- 线上数据保存。

验收：

```text
使用公网链接完成注册→档案→JD分析→面试
```

## 7 月 29 日：真实用户试用

邀请至少 3 人，每人执行：

1. 注册；
2. 填档案；
3. 分析一份 JD；
4. 完成至少 3 轮面试；
5. 查看历史记录；
6. 删除一条测试记录；
7. 提交反馈。

记录：

- 哪一步不理解；
- 哪个按钮找不到；
- AI 内容是否过长；
- 手机是否溢出；
- 是否数据丢失；
- 是否出现接口失败。

## 7 月 30 日：修复问题

优先级：

1. 无法使用；
2. 数据丢失；
3. 权限问题；
4. AI 返回异常；
5. 移动端错位；
6. 文案和样式细节。

不新增大功能。

## 7 月 31 日：发布与求职材料

完成：

- README；
- 项目截图；
- 2～3 分钟演示录屏；
- 架构图；
- 测试账号说明或开放注册；
- 已知问题；
- GitHub 仓库整理；
- 简历项目经历；
- Boss 在线简历更新；
- 项目讲述稿。

---

# 21. 每日固定执行节奏

```text
09:30～11:00  投递和打招呼
11:00～11:30  跟进回复、记录投递
13:30～16:30  V1 核心开发
16:45～18:30  联调、修复和提交代码
19:30～20:30  八股复习
20:30～21:00  项目讲述和当日复盘
```

每日开发底线：

- 至少完成一个可验收的小功能；
- 至少一次 Git 提交；
- 不允许连续两天只看教程；
- 当天学到的内容当天落到项目；
- 出现面试邀约时，面试准备优先。

---

# 22. Git 分支与提交规范

个人项目不需要复杂 Git Flow。

```text
main：始终保持可部署
feat/xxx：较大功能临时分支
fix/xxx：线上问题
```

提交示例：

```text
feat(auth): 完成 Supabase 邮箱登录与路由守卫
feat(profile): 实现个人档案保存和完整度计算
feat(ai): 接入 JD 匹配分析接口
feat(interview): 实现多轮面试消息持久化
fix(history): 修复删除记录后列表未刷新
docs(readme): 补充部署说明和产品截图
```

---

# 23. 完成定义（Definition of Done）

一个功能只有同时满足以下条件，才算完成：

- 正常流程可用；
- 输入有校验；
- 请求有 loading；
- 失败有提示；
- 数据能保存；
- 刷新后状态合理；
- 不访问他人数据；
- 页面在 1366px 宽度不明显错位；
- 代码没有明显 `any` 滥用；
- 已完成 Git 提交。

“页面画出来但接口没通”不算完成。  
“接口能调但刷新丢数据”不算完成。  
“自己能用但别人注册不了”不算落地。

---

# 24. 上线前检查清单

## 功能

- [ ] 注册
- [ ] 登录
- [ ] 退出
- [ ] 路由保护
- [ ] 新建档案
- [ ] 编辑档案
- [ ] JD 分析
- [ ] 保存分析
- [ ] 查看历史分析
- [ ] 删除分析
- [ ] 开始面试
- [ ] 多轮回答
- [ ] 刷新恢复
- [ ] 结束面试
- [ ] 最终报告
- [ ] 历史面试
- [ ] 仪表盘统计

## 安全

- [ ] AI Key 不在前端
- [ ] service role key 不在前端
- [ ] `.env` 已加入 `.gitignore`
- [ ] 所有业务表启用 RLS
- [ ] 用户只能读写自己的数据
- [ ] AI 接口要求 Bearer token
- [ ] 输入长度有限制
- [ ] 调用次数有限制
- [ ] 日志不输出完整简历

## 体验

- [ ] 所有按钮有 loading
- [ ] 所有列表有空状态
- [ ] AI 超时有提示
- [ ] 手机端可完成核心流程
- [ ] 长文本不会撑破布局
- [ ] 删除有二次确认
- [ ] AI 分数有免责声明

## 发布

- [ ] 公网网址
- [ ] GitHub 仓库
- [ ] README
- [ ] 产品截图
- [ ] 演示录屏
- [ ] 3 名真实用户完成试用
- [ ] 已知问题列表
- [ ] 简历项目描述

---

# 25. README 结构

```text
1. 产品介绍
2. 在线体验
3. 产品截图
4. 核心功能
5. 技术栈
6. 系统架构
7. 数据流
8. 项目目录
9. 本地启动
10. 环境变量
11. 数据库初始化
12. 部署说明
13. 隐私与 AI 说明
14. 已知问题
15. 后续规划
```

---

# 26. 简历项目描述草稿

## 求职陪跑 AI 助手｜全栈独立开发

**技术栈：Vue3 / TypeScript / Vite / Pinia / Element Plus / Node.js / Express / Supabase / PostgreSQL / 大模型 API**

项目介绍：

> 面向实习与校招求职场景，独立设计并上线 AI 求职陪跑 Web 应用，围绕个人求职档案、JD 匹配分析、多轮模拟面试和历史复盘形成完整使用闭环，支持真实用户注册、跨设备数据持久化与公网访问。

技术亮点：

- 基于 Vue3、TypeScript、Pinia 和 Element Plus 构建响应式前端，拆分档案、JD 分析、模拟面试和历史记录等业务模块，并统一处理表单校验、加载状态和异常反馈；
- 使用 Supabase Auth、PostgreSQL 与 RLS 实现用户鉴权和数据隔离，保存分析报告、面试会话及消息记录，支持刷新恢复和跨设备访问；
- 在 Express 服务端封装大模型调用，根据岗位 JD 与用户档案生成结构化匹配报告，并使用 Zod 校验模型输出，避免异常内容直接进入前端；
- 实现多轮模拟面试流程，保存上下文消息并生成逐轮评分、优点、改进建议和最终复盘报告；
- 增加 AI 调用配额、输入长度限制、超时处理和统一错误码，保护模型密钥并控制公开测试版调用成本。

注意：没有真实数据时，不要在简历中编造用户数、调用次数和提升百分比。上线后可以记录真实数据再补充。

---

# 27. 面试讲述框架

## 27.1 为什么做

> 我在投实习过程中发现，岗位 JD 很分散，自己经常不知道该突出哪些项目，也缺少真实面试机会。通用聊天工具每次都要重新输入个人背景，所以我把个人求职档案、JD 分析、模拟面试和历史复盘整合成一个可长期使用的产品。

## 27.2 最大难点

优先讲三项：

1. AI 返回内容不稳定，如何用结构化 Prompt 和 Zod 校验；
2. 多轮面试上下文如何控制长度并持久化；
3. Supabase Auth 与 RLS 如何保证用户数据隔离。

## 27.3 为什么不做 RAG

> V1 的核心数据是用户档案和单份 JD，上下文规模有限，直接使用结构化 Prompt 已经可以完成核心需求。为了在月底前完成可落地版本，我没有过早引入 RAG 和向量数据库，后续在接入面试题库、企业资料库时再考虑。

## 27.4 为什么选择 Supabase

> 项目需要真实注册、跨设备持久化和用户数据隔离。Supabase 能同时提供 Auth 和 PostgreSQL，并可使用 RLS 控制用户只能访问自己的记录，减少自建完整用户系统的时间，把主要精力放在核心 AI 业务上。

---

# 28. 今天立即执行的任务

按顺序完成，不继续扩展需求：

```text
1. 创建 job-copilot-web
2. 创建 job-copilot-server
3. 安装依赖
4. 建立前后端目录
5. 前端创建 7 个空路由页面
6. 完成 AppLayout、侧边栏和顶部栏
7. 后端完成 GET /health
8. 前端请求 /health 并在控制台打印
9. 创建 Git 仓库并提交
```

今天的完成标准：

> 打开前端能看到与参考图同风格的页面骨架；点击左侧菜单可以切换页面；后端健康检查正常；前端与后端已经连通。

完成这些后，再进入 Supabase 登录，不要提前写 AI 功能。

---

# 29. 最终原则

1. **先闭环，后美化。**
2. **先真实数据，后仪表盘。**
3. **先一个模型，后多模型。**
4. **先文字面试，后语音。**
5. **先真实上线，后增加功能。**
6. **任何未完成的功能不写成已完成。**
7. **每天必须留下可运行成果，而不是只留下学习笔记。**
