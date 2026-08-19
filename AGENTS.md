# 项目定位

- 项目名称：求职陪跑 AI 助手（秋招冲刺版）；
- 项目形态：面向实习和校招求职者的响应式 Web 应用；
- 核心闭环：注册/登录 → 个人档案 → JD 分析 → 模拟面试 → 历史复盘；
- 当前目标：先完成真实可运行、可部署、可用于秋招展示的 V1，再做工程化和视觉优化；
- 需求和阶段范围以 `求职陪跑_AI助手_秋招冲刺版_任务规划_精简版_2026-08-12.md` 为准。

# 当前真实进度（2026-08-12）

已经完成：

- Vue 3 + TypeScript + Vite 前端工程；
- 登录、注册、Dashboard、JD Analysis、Interview、History、Profile、Settings 路由和静态页面骨架；
- AppLayout、响应式侧栏、顶栏、移动端抽屉、设计变量和空状态；
- Supabase Client、邮箱注册/登录/退出代码、Session 初始化、刷新保持监听和路由守卫；
- Axios 请求实例；
- 旧 Express + TypeScript 后端骨架、CORS、限流、错误处理以及 `/health`、`/api/v1/health`；
- 缺少 Supabase 配置时的友好提示。

尚未完成或尚未真实联调：

- Supabase 项目配置以及真实注册、登录、刷新恢复、退出的完整联调；
- `profiles` 表、RLS 和个人档案的保存/读取/修改闭环；
- Go + Gin 后端；
- JD 分析、AI 结构化输出和 `jd_analyses` 持久化；
- `interview_sessions`、`interview_messages`、5 轮模拟面试和最终报告；
- History、Dashboard 真实统计、部署和完整端到端验收。

任何文档、页面提示和完成说明都必须反映以上真实边界。

没有真实凭据、真实数据或没有执行联调时，不得声称相关功能已经可用。

# 技术栈与迁移方向

## 前端

- Vue 3；
- TypeScript；
- Vite；
- Vue Router；
- Pinia；
- Element Plus；
- Axios；
- Supabase JavaScript Client。

## 目标后端

- Go；
- Gin；
- Supabase Auth；
- PostgreSQL / Supabase Database；
- LLM API（OpenAI 兼容接口）。

仓库中的 `backend/` 当前仍是 Node.js + Express + TypeScript 的早期骨架，只实现健康检查。

它属于待迁移的历史实现，不是后续业务后端的目标技术栈。

## 迁移规则

- 新增 JD 分析、模拟面试、报告等业务 API 时，使用 Go + Gin，不继续扩展 Express 业务能力；
- Go 后端先实现与现有健康检查等价的最小可运行服务，再逐步承接业务接口；
- 在 Go 健康检查和前端 API 地址切换完成前，不直接删除旧 `backend/`；
- 不让 Node 和 Go 同时维护同一份业务逻辑；
- 迁移步骤、目录命名、端口或环境变量发生变化时，同步更新 README、`.env.example` 和相关文档；
- 未经当前任务明确要求，不提前一次性重写全部后端；
- 不为了迁移而顺手修改与当前阶段无关的前端业务。

# 数据流转

## 1. 登录与身份

```text
登录/注册页面
  → Pinia auth store
  → Supabase JS Client
  → Supabase Auth
  → 返回 Session / User
  → store 保存认证状态
  → Vue Router 根据 Session 放行 /app/**
```

规则：

- 刷新页面时，先恢复 Supabase Session，再执行路由判断；
- 前端只能使用 Supabase URL 和 anon key；
- access token 只能在运行时用于认证请求，禁止写入日志、源码或文档；
- service role key 和 AI API Key 只能存在于后端环境变量；
- 不在前端保存或伪造可信 `user_id`；
- 用户身份以 Supabase 校验后的 Session / Token 为准。

## 2. 业务 API

```text
Vue 页面
  → API 模块 / Axios
  → Authorization: Bearer <Supabase access token>
  → Gin HTTP Handler / Auth Middleware
  → Application Service
  → Domain 业务规则
  → Port 接口
  → PostgreSQL Repository / AI Adapter / Supabase Auth Adapter
  → 统一 JSON 响应
  → 页面更新 loading / success / empty / error 状态
```

规则：

- Go 后端验证用户身份后，把可信 `user_id` 传入应用层；
- 所有用户业务数据都按可信 `user_id` 查询和写入；
- 用户业务表通过 RLS 做数据隔离；
- 不信任前端提交的 `user_id`；
- 不允许客户端选择、覆盖或伪造数据所有者；
- 前端不直接调用 LLM；
- Prompt、AI Key、重试、输出解析和字段校验全部在 Go 后端；
- API 响应保持统一结构；
- 业务错误、鉴权错误、参数校验错误、数据库错误和上游 AI 错误需要可区分。

## 3. JD 分析闭环

```text
Profile + 公司 + 岗位 + JD
  → POST /api/v1/ai/analyze-jd
  → Application 组织上下文
  → AIClient Port
  → LLM 返回 JSON
  → Go struct 解码
  → 业务字段校验
  → 保存 jd_analyses
  → 返回结构化结果
  → Result / History 展示
```

结构化结果至少包括：

- `matchScore`；
- `jobSummary`；
- `coreRequirements`；
- `matchedSkills`；
- `missingSkills`；
- `resumeSuggestions`；
- `preparationTopics`；
- `greetingMessage`。

规则：

- `matchScore` 必须为 0～100 的整数；
- AI 返回缺字段、格式错误或数据越界时，不得作为正常结果保存；
- 模型调用失败时，不生成伪造成功数据；
- 页面不能因为模型返回异常内容而直接崩溃。

## 4. 模拟面试闭环

```text
选择 JD
  → 创建 interview_session
  → AI 生成第一题
  → 保存 message
  → 用户提交回答
  → 读取 Profile / JD / Session / 最近 N 条 Messages
  → AI 评分、反馈并生成下一题
  → 保存消息
  → 更新 currentRound
  → 第 5 轮结束
  → 生成并保存最终报告
  → History 可恢复和复盘
```

规则：

- 一场面试默认最多 5 轮；
- `completed` 的 Session 不能继续提交答案；
- 必须防止重复提交；
- 必须保证刷新后可以从数据库恢复当前 Session 和 Messages；
- 传给 LLM 的上下文必须限制长度；
- 不无限拼接全部历史消息；
- 最终报告生成成功后才能将 Session 标记为完成；
- 消息保存、轮次更新和最终状态必须保持一致。

# Go 后端架构原则

采用轻量 DDD + 六边形架构，但只按当前功能创建必要文件。

```text
HTTP 请求
  → Adapter（Gin Handler / Middleware）
  → Application（用例编排）
  → Domain（实体和业务规则）
  → Port（依赖接口）
  → Adapter（PostgreSQL / Supabase / LLM）
```

## Domain

负责表达稳定业务规则。

例如：

- 一场面试最多 5 轮；
- `completed` Session 不能继续回答；
- 分数必须在 0～100；
- Session 状态如何变化。

Domain 不依赖：

- Gin；
- Supabase SDK；
- PostgreSQL Driver；
- 具体 LLM SDK。

## Application

负责：

- 用例执行顺序；
- 调用 Domain；
- 调用 Repository；
- 调用 AI Port；
- 组织事务边界；
- 传递上下文。

Application 不处理 HTTP 展示细节。

## Port

为外部能力定义最小接口，例如：

- Auth；
- Repository；
- AI Client。

Port 只定义当前业务真正需要的方法，不提前设计大量未来接口。

## Adapter

负责：

- Gin HTTP；
- Auth Middleware；
- Supabase；
- PostgreSQL；
- LLM；
- 外部系统错误转换。

## 架构约束

- 先确定当前接口并完成最小闭环，业务稳定后再抽象；
- 不为了“像 DDD”创建大量空目录；
- 不创建没有实际价值的空接口；
- 不创建只有一层机械转发的无意义 Service；
- 不提前实现当前 V1 不需要的通用平台能力；
- Go 的 `context.Context` 沿调用链传递；
- 错误必须显式处理；
- 不使用 `panic` 代替正常业务错误；
- AI JSON 必须使用明确 struct 解码并执行字段校验；
- 不把未验证的任意 JSON 直接返回前端。

# V1 开发顺序

严格按以下顺序推进。

每次只完成当前阶段的可验收结果：

1. 配置 Supabase，真实联调注册、登录、Session、路由守卫和退出；
2. 创建 `profiles` 表、RLS，完成档案保存、读取、修改和刷新恢复；
3. 初始化 Go + Gin，实现 config、router、handler、统一响应、CORS 和 `GET /health`；
4. 按实际需要建立 Domain、Application、Port、Adapter 最小骨架；
5. 完成 JD 输入、Go 分析 API、AI 结构化校验、`jd_analyses` 保存和结果页；
6. 完成面试 Session、消息持久化、5 轮问答和最终报告；
7. 用真实数据完成 History、Dashboard 和数据隔离；
8. 完成部署、完整流程回归、README、架构图、截图、Demo 和秋招材料。

除非用户明确调整规划，否则不要跨阶段提前开发后续功能。

如果当前任务只属于某一阶段，只处理当前阶段需要的内容。

# V1 范围限制

8 月 31 日前不增加以下内容：

- RAG；
- 向量数据库；
- PDF 解析；
- 语音面试；
- 多模型切换；
- 支付；
- 社区；
- 岗位爬虫；
- 微服务；
- Kafka；
- CQRS；
- Kubernetes；
- Redis 集群；
- 复杂 DDD；
- 过度抽象；
- 与当前闭环无关的通用平台能力；
- 通知；
- 复杂搜索；
- 自定义计划；
- 近期提醒；
- Pro 功能。

未经用户明确要求，不主动建议把上述内容加入当前 V1。

# Agent 执行方式

默认采用直接执行模式。

普通开发任务不单独生成长 Plan，也不因为任务简单而停留在规划阶段。

默认执行流程：

```text
读取必要文件
  → 确认当前实现
  → 判断影响范围
  → 直接修改
  → 运行必要检查
  → 汇报修改结果
```

以下任务默认直接执行：

- 单页面开发或调整；
- 样式修改；
- 表单字段修改；
- 明确 Bug 修复；
- 单个 API；
- 单个 Use Case；
- 接口联调；
- 类型错误修复；
- 当前阶段内的小范围重构；
- 新增当前需求明确要求的文件。

只有以下情况需要先给出简短方案：

- 跨多个核心领域的大规模重构；
- 数据库存在破坏性迁移；
- 当前需求与现有架构明显冲突；
- 修改范围无法从现有代码判断；
- 存在多个会明显影响后续架构的实现方向。

即使需要方案，也只给必要的短计划，随后继续实现。

除非用户明确要求只分析，否则不要停留在 Plan 阶段。

# 修改前检查

开始修改前必须先读取与当前任务直接相关的现有代码。

至少确认：

- 当前目录和文件是否已经存在；
- 当前项目实际使用的命名和组织方式；
- 是否已有可复用的组件；
- 是否已有可复用的类型；
- 是否已有 API 封装；
- 是否已有 Store；
- 是否已有 Middleware；
- 是否已有工具函数；
- 当前接口实际请求和响应结构；
- 当前用户未提交的修改是否会受到影响。

`AGENTS.md` 描述的是目标规范，不代表所有目标结构已经存在。

不得因为文档中存在某个目录、接口或模块名称，就假设代码已经实现。

优先修改现有实现。

不要重复创建功能相同的：

- API Client；
- Store；
- Type；
- Component；
- Middleware；
- Repository；
- Utility。

# 最小修改原则

每次任务只修改完成当前需求所必须的代码。

禁止：

- 顺手重构无关模块；
- 批量修改无关格式；
- 修改与需求无关的文件命名；
- 为未来可能使用的功能提前抽象；
- 因个人偏好替换已经可以工作的库或实现；
- 将一个局部需求扩展为全项目重构；
- 因为发现其他 Bug 就顺手一起修改；
- 覆盖用户当前未提交但与任务无关的修改。

如果发现其他问题但不影响当前任务：

- 不主动修改；
- 可以在最终说明中简要指出。

# API Contract 规则

前后端接口以已经确定的 API Contract 为边界。

未经当前任务明确要求，不随意修改：

- API Path；
- HTTP Method；
- Request 字段；
- Response 字段；
- 字段命名；
- 错误结构；
- 状态码语义。

统一响应优先保持：

```json
{
  "code": "OK",
  "message": "success",
  "data": {},
  "requestId": "optional"
}
```

前端 JSON 字段统一使用 camelCase，例如：

```text
matchScore
currentRound
companyName
jobTitle
```

Go struct 使用 Go 命名风格，但 JSON Tag 必须保持 API Contract：

```go
type AnalysisResult struct {
	MatchScore int `json:"matchScore"`
}
```

如果必须修改 API Contract：

1. 先确认所有调用方；
2. 同步修改前端 TypeScript 类型；
3. 同步修改 API 调用；
4. 同步修改 Go DTO / Struct；
5. 同步修改 README 或接口文档；
6. 完成前后端联调后才能声称修改完成。

# 数据库规则

数据库 Schema 是跨前后端共享契约，不随意修改。

新增或修改数据库字段前必须确认：

- 当前表结构；
- TypeScript 类型；
- Go Domain / DTO；
- Repository；
- RLS Policy；
- 页面实际使用字段。

禁止为了修复普通页面问题直接：

- 删除字段；
- 删除表；
- 删除 RLS Policy；
- 重建已有业务表；
- 清空已有数据。

涉及数据库结构变化时：

- 使用明确 SQL migration；
- 不直接覆盖已有生产数据；
- 新增字段优先考虑兼容现有数据；
- 同步更新相关 Go struct；
- 同步更新 TypeScript interface；
- 同步检查 Repository；
- 同步检查 RLS Policy；
- 同步检查页面读写逻辑。

任何用户业务表默认包含可信的 `user_id` 所有权关系。

不得使用前端传入的 `user_id` 作为最终数据所有者。

最终数据所有者必须来自后端验证后的真实登录用户身份。

# AI 功能开发规则

LLM 输出不是可信数据源。

所有 AI 业务接口必须遵循：

```text
输入校验
  → 构建 Prompt
  → 调用 LLM
  → 解析结构化输出
  → Go Struct 解码
  → 业务字段校验
  → 必要时保存数据库
  → 返回前端
```

禁止：

- 将 LLM 原始任意 JSON 直接透传前端；
- 使用字符串拼接方式手工构造业务 JSON；
- 因模型偶发缺字段让页面直接崩溃；
- 将完整 Prompt 写入日志；
- 将完整 JD 写入日志；
- 将完整简历写入日志；
- 将完整面试答案写入日志；
- 在浏览器中保存 AI API Key；
- 模型失败后保存伪造的成功结果。

AI 调用必须考虑：

- timeout；
- `context.Context` cancellation；
- 上游 4xx / 5xx；
- 非法 JSON；
- 缺少字段；
- 字段类型错误；
- 数值越界；
- 重复提交；
- 用户取消或页面离开。

所有 AI 评分字段必须校验：

```text
0 <= score <= 100
```

模型失败时：

- 返回明确业务错误；
- 不保存“成功”的分析结果；
- 不更新错误的面试轮次；
- 不将失败 Session 标记为 completed。

# 安全与隐私规则

禁止在源码、日志、README、截图、测试数据中写入：

- 密码；
- access token；
- refresh token；
- service role key；
- AI API Key；
- 用户完整简历；
- 用户完整 JD；
- 用户完整面试内容。

前端允许公开的环境变量只包括真正设计为公开使用的值，例如：

```text
VITE_SUPABASE_URL
VITE_SUPABASE_ANON_KEY
VITE_API_BASE_URL
```

真实 `.env` 不提交 Git。

只提交：

```text
.env.example
```

后端日志如需记录用户，只记录必要的脱敏标识。

# 修改后的检查

根据修改范围执行最小必要检查。

## 前端

普通 TS / Vue 修改：

```bash
npm run typecheck
```

以下情况再执行：

```bash
npm run build
```

适用场景：

- 修改构建配置；
- 修改依赖；
- 修改路由；
- 修改入口文件；
- 完成一个完整阶段功能；
- 准备提交一个较完整功能。

## Go

修改过的 `.go` 文件必须执行：

```bash
gofmt
```

普通业务修改：

- 优先运行相关 package test。

阶段功能完成时：

```bash
go test ./...
```

涉及以下内容时：

- 启动入口；
- 配置；
- 依赖；
- Router；
- Middleware；
- Build 相关修改；

需要确认：

```bash
go build ./...
```

## 旧 Node 后端

只有实际修改旧 `backend/` 时，才运行它对应的：

```bash
npm run typecheck
npm run build
```

不要因为只修改 Go 或前端而无意义运行旧 Node 后端检查。

## 无法检查时

如果因为：

- 网络；
- 私有依赖；
- 环境变量；
- 缺少真实凭据；
- 本地环境缺少工具链；

导致无法执行检查：

- 明确说明哪个命令没有执行；
- 明确说明具体原因；
- 不得写成“测试通过”；
- 不得编造运行结果。

# 通用开发规则

- Vue 页面和组件优先使用 `<script setup lang="ts">` 和 Composition API；
- 重要数据必须定义 `interface` 或 `type`；
- 尽量避免 `any`；
- 不使用 `@ts-ignore` 掩盖类型错误；
- 页面根据实际功能覆盖 loading、empty、error 等状态；
- 涉及持久化数据的页面需要考虑刷新恢复；
- Dashboard 和 History 只展示真实数据；
- 不使用虚假企业、虚假统计或虚假面试记录作为正式功能数据；
- 表单前后端都要校验；
- 后端校验是最终可信边界；
- 对关键逻辑添加简明中文注释；
- 注释重点解释“为什么这样写”；
- 不逐行翻译代码；
- 项目面向正在学习 Vue 3、TypeScript 和 Go 的开发者；
- 代码保持直接、清晰；
- 避免过度封装；
- 不把真实密钥写入源码；
- 不提交真实 `.env`；
- 没有真实运行过的命令、接口或流程，不得声称通过；
- 发现与当前任务无关的用户改动时保留；
- 不覆盖；
- 不回退；
- 不删除；
- 未经用户明确允许，不执行 `git push`；
- 未经用户明确允许，不执行远程部署；
- 未经用户明确允许，不修改线上服务；
- 未经用户明确允许，不提交真实外部数据；
- 未经用户明确允许，不删除数据库表、生产数据或远程资源。

# Git 修改规则

Agent 可以：

- 查看 `git status`；
- 查看 `git diff`；
- 查看必要的 Git 历史；
- 创建当前任务所需的本地代码修改。

未经用户明确要求，不执行：

```bash
git push
git reset --hard
git clean -fd
git rebase
git merge
git checkout -- .
git restore .
```

不擅自回退用户已有修改。

如果发现工作区已有用户修改：

- 先识别哪些文件属于当前任务；
- 当前任务无关文件保持原样；
- 不用自己的版本覆盖用户正在修改的内容。

# 阶段完成标准

“完成”必须意味着对应阶段的真实链路已经达到当前可验证条件。

不能仅因为以下内容存在就声称功能完成：

- UI 已画出；
- Mock 已写；
- 函数已创建；
- Interface 已定义；
- Route 已注册；
- SQL 已写但没有执行；
- API 已写但没有真实调用；
- 测试代码存在但没有运行。

根据功能适用范围，阶段验收需要检查：

- 成功路径；
- 失败路径；
- loading 状态；
- 错误状态；
- 刷新恢复；
- 重复提交；
- 权限边界；
- 数据持久化。

不是所有功能都需要验证全部项目。

例如：

- `/health` 不需要验证刷新恢复；
- 纯展示组件不需要验证数据库权限。

使用 Mock 时必须明显标注。

完成真实接口后，及时移除不再需要的 Mock。

最终优先级始终是：

```text
先闭环
  → 再工程化
  → 再美化
  → 再扩展
```

# Agent 最终汇报格式

完成代码修改后，最终说明保持简洁。

至少包含以下内容：

## 修改内容

说明实际修改了什么。

## 检查结果

说明实际运行了哪些检查以及结果。

例如：

```text
- npm run typecheck：通过
- go test ./...：通过
```

如果没有运行：

```text
- 未执行真实 Supabase 联调：当前环境未配置 Supabase 凭据
```

## 未完成或风险

只有真实存在时才写。

不要把以下状态：

- 未运行；
- 未验证；
- 推测可用；

描述成：

- 已完成；
- 已通过；
- 已上线。

如果当前任务已经完整完成且没有其他问题，不需要额外扩展无关建议。
