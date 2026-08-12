# 项目说明

- 项目名称：求职陪跑 AI 助手；
- 项目形态：响应式 Web 应用；
- 面向实习和校招求职者；
- 核心流程是个人档案、JD 分析、模拟面试和历史复盘。

# 技术栈

前端：

- Vue 3；
- TypeScript；
- Vite；
- Vue Router；
- Pinia；
- Element Plus；
- Axios；
- Supabase JavaScript Client。

后端：

- Node.js；
- Express；
- TypeScript；
- Zod；
- dotenv；
- cors；
- helmet；
- express-rate-limit；
- Supabase JavaScript Client；
- OpenAI 兼容 SDK。

# 开发规则

- Vue 页面和组件优先使用 `<script setup lang="ts">`；
- 优先使用 Composition API；
- 重要数据必须定义 `interface` 或 `type`；
- 尽量避免 `any`；
- 不使用 `@ts-ignore` 掩盖类型错误；
- 不把真实密钥写入源码；
- 不提交真实 `.env`；
- 不暴露 Supabase access token；
- 不暴露 AI API Key；
- 不提前开发当前阶段之外的功能；
- 所有修改完成后运行 typecheck 和 build；
- 没有真实运行过的命令，不得声称通过；
- 对关键逻辑添加简明中文注释；
- 注释重点解释“为什么这样写”，不要逐行翻译代码；
- 项目面向刚学习 Vue3 和 TypeScript 的开发者，避免过度封装；
- 未经用户明确允许，不执行 git push、远程部署或修改线上服务。
