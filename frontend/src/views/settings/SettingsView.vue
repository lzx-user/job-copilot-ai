<script setup lang="ts">
import { Delete, InfoFilled, Lock, SwitchButton, Warning } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'
import { useRouter } from 'vue-router'
import { useAuthStore } from '../../stores/auth.store'

const router = useRouter()
const authStore = useAuthStore()

async function logout() {
  const result = await authStore.logout()
  if (result.success) {
    await router.replace('/auth/login')
    ElMessage.success('已安全退出登录')
  } else ElMessage.error(authStore.authError)
}
</script>

<template>
  <div class="page-shell settings-page">
    <div class="page-heading"><div><span class="section-label">SETTINGS</span><h1>设置</h1><p>了解产品边界、AI 使用方式与数据隐私原则。</p></div></div>
    <div class="settings-layout">
      <section class="panel setting-card"><span class="setting-icon blue"><el-icon><InfoFilled /></el-icon></span><div><h2>产品说明</h2><p>求职陪跑 AI 助手面向实习与校招准备，当前完成工程基础、静态页面骨架和 Supabase 认证代码。</p><el-tag effect="plain">V1 · 7.14–7.16 阶段</el-tag></div></section>
      <section class="panel setting-card"><span class="setting-icon violet"><el-icon><Warning /></el-icon></span><div><h2>AI 使用说明</h2><p>当前阶段不会调用 AI。后续生成的分析与评分只用于求职准备参考，不代表企业真实筛选或录用结果。</p></div></section>
      <section class="panel setting-card"><span class="setting-icon green"><el-icon><Lock /></el-icon></span><div><h2>数据隐私说明</h2><p>密码由 Supabase Auth 处理，源码不会保存密码或 token。业务表接入时必须启用 RLS，保证用户只能访问自己的数据。</p></div></section>
      <section class="panel account-card"><div><h2>账号操作</h2><p>退出后本机 Session 会被清理，再次访问受保护页面需要重新登录。</p></div><el-button type="danger" plain :icon="SwitchButton" :loading="authStore.submitting" @click="logout">退出登录</el-button></section>
      <section class="panel danger-card"><div><h2>删除账号</h2><p>V1 当前阶段尚未开放自助删除。请勿通过删除本地文件代替服务端账号删除。</p></div><el-button :icon="Delete" disabled>后续阶段评估</el-button></section>
    </div>
  </div>
</template>

<style scoped>
.settings-layout { display: grid; grid-template-columns: repeat(2,minmax(0,1fr)); gap: 16px; }.setting-card { display: flex; min-height: 170px; padding: 24px; gap: 16px; }.setting-card:first-child { grid-column: 1/-1; }.setting-card h2,.account-card h2,.danger-card h2 { margin-bottom: 7px; font-size: 18px; }.setting-card p,.account-card p,.danger-card p { color: var(--text-secondary); }.setting-icon { display: grid; flex: 0 0 46px; width: 46px; height: 46px; place-items: center; border-radius: 14px; font-size: 22px; }.setting-icon.blue { color: var(--color-primary); background: #edf3ff; }.setting-icon.violet { color: var(--color-secondary); background: #f1ebff; }.setting-icon.green { color: var(--color-success); background: #e9faf4; }.account-card,.danger-card { display: flex; padding: 22px 24px; align-items: center; justify-content: space-between; gap: 20px; }.danger-card { grid-column: 1/-1; border-color: #f4e0e0; }.account-card p,.danger-card p { margin: 0; }
@media (max-width: 800px) { .settings-layout { grid-template-columns: 1fr; }.setting-card:first-child,.danger-card { grid-column: auto; }.account-card,.danger-card { align-items: stretch; flex-direction: column; }.account-card .el-button,.danger-card .el-button { width: 100%; } }
</style>

