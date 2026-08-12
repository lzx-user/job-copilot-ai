<script setup lang="ts">
import { Edit, InfoFilled, Plus, User } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'
import { computed } from 'vue'
import { useAuthStore } from '../../stores/auth.store'

const authStore = useAuthStore()
const email = computed(() => authStore.user?.email || '尚未获取邮箱')
const avatarLetter = computed(() => email.value.slice(0,1).toUpperCase())

function explainAvailability() {
  ElMessage.info('个人档案保存将在 7 月 17 日阶段完成，当前不会写入数据库。')
}
</script>

<template>
  <div class="page-shell">
    <div class="page-heading">
      <div><span class="section-label">CANDIDATE PROFILE</span><h1>个人档案</h1><p>完善个人信息后，AI 才能提供更贴近你的 JD 分析与面试练习。</p></div>
      <el-alert class="stage-alert" title="个人档案保存将在下一阶段完成" type="warning" :closable="false" show-icon />
    </div>

    <div class="profile-layout">
      <main class="profile-main">
        <section class="panel identity-card">
          <div class="avatar">{{ avatarLetter }}</div>
          <div class="identity-copy"><span class="identity-tag">个人档案</span><h2>待完善昵称</h2><p>目标岗位尚未填写</p><small>{{ email }}</small></div>
          <div class="mini-ring"><strong>0%</strong><span>已完善</span></div>
        </section>

        <section class="panel section-card">
          <div class="panel-header"><h2 class="panel-title">基础信息</h2><el-button :icon="Edit" disabled>编辑</el-button></div>
          <div class="info-grid">
            <div><span>昵称</span><strong>待完善</strong></div><div><span>目标岗位</span><strong>待完善</strong></div><div><span>期望城市</span><strong>待完善</strong></div><div><span>可实习时间</span><strong>待完善</strong></div><div><span>登录邮箱</span><strong>{{ email }}</strong></div><div><span>毕业年份</span><strong>待完善</strong></div>
          </div>
        </section>

        <section class="panel section-card">
          <div class="panel-header"><h2 class="panel-title">技术栈与关键词</h2><el-button :icon="Edit" disabled>编辑</el-button></div>
          <div class="empty-inline"><el-icon :size="24"><Plus /></el-icon><span>还没有添加技能，下一阶段可填写最多 30 个关键词。</span></div>
        </section>

        <section class="panel section-card text-section">
          <div class="panel-header"><h2 class="panel-title">项目经历摘要</h2><el-button :icon="Edit" disabled>编辑</el-button></div>
          <p>还没有项目摘要。建议提前准备项目背景、你的职责、关键难点、解决方案和可验证结果。</p>
        </section>

        <section class="panel section-card text-section">
          <div class="panel-header"><h2 class="panel-title">个人优势 / 自我介绍</h2><el-button :icon="Edit" disabled>编辑</el-button></div>
          <p>还没有个人优势说明。下一阶段填写时，请优先使用真实经历和具体证据。</p>
        </section>

        <div class="profile-actions"><el-button class="gradient-button" @click="explainAvailability">保存档案（下一阶段开放）</el-button><el-button disabled>预览简历摘要</el-button></div>
      </main>

      <aside class="profile-aside">
        <section class="panel completeness-card">
          <h2 class="panel-title">档案完整度</h2>
          <div class="large-ring"><div><strong>0%</strong><span>待完善</span></div></div>
          <ul><li><i />基础信息待填写</li><li><i />技术栈待补充</li><li><i />项目经历待填写</li><li><i />个人优势待填写</li></ul>
        </section>
        <section class="panel advice-card"><span class="advice-icon">✦</span><div><h3>档案填写建议</h3><p>只填写真实经历。清晰、具体、可验证的信息，比堆叠关键词更有帮助。</p></div></section>
        <section class="panel privacy-card"><el-icon><InfoFilled /></el-icon><div><h3>数据说明</h3><p>当前页面不保存信息。后续接入数据库时，将使用 RLS 隔离每位用户的数据。</p></div></section>
      </aside>
    </div>
  </div>
</template>

<style scoped>
.stage-alert { width: min(100%,390px); }
.profile-layout { display: grid; grid-template-columns: minmax(0,1fr) 350px; gap: 18px; }.profile-main,.profile-aside { display: grid; align-content: start; gap: 16px; }
.identity-card { display: flex; min-height: 150px; padding: 24px 28px; align-items: center; gap: 20px; background: linear-gradient(125deg,#fff,#f8f9ff); }
.avatar { display: grid; flex: 0 0 78px; width: 78px; height: 78px; place-items: center; border: 7px solid #eef2ff; border-radius: 50%; color: #fff; background: linear-gradient(140deg,#7487a8,#3e4d69); font-size: 27px; font-weight: 800; }
.identity-copy { min-width: 0; }.identity-copy h2 { margin: 5px 0 2px; }.identity-copy p { margin: 0; color: var(--text-regular); }.identity-copy small { display: block; overflow: hidden; max-width: 360px; margin-top: 7px; color: var(--text-secondary); text-overflow: ellipsis; white-space: nowrap; }
.identity-tag { padding: 4px 8px; border-radius: 7px; color: var(--color-primary); background: #edf3ff; font-size: 12px; }.mini-ring { display: grid; width: 82px; height: 82px; margin-left: auto; place-items: center; align-content: center; border: 7px solid #edf0f6; border-radius: 50%; }.mini-ring strong { font-size: 20px; }.mini-ring span { color: var(--text-secondary); font-size: 11px; }
.section-card { padding: 22px; }.info-grid { display: grid; grid-template-columns: repeat(3,minmax(0,1fr)); gap: 12px; }.info-grid > div { display: grid; min-width: 0; padding: 12px 14px; gap: 4px; border-radius: 10px; background: #f8f9fd; }.info-grid span { color: var(--text-secondary); font-size: 12px; }.info-grid strong { overflow: hidden; text-overflow: ellipsis; white-space: nowrap; }
.empty-inline { display: flex; min-height: 70px; padding: 18px; align-items: center; gap: 12px; border: 1px dashed #d8dfef; border-radius: 12px; color: var(--text-secondary); background: #fbfcff; }.text-section p { margin: 0; color: var(--text-secondary); }
.profile-actions { display: grid; grid-template-columns: 1fr 1fr; gap: 14px; }.completeness-card { display: grid; padding: 22px; place-items: center; }.completeness-card h2 { justify-self: start; }.large-ring { display: grid; width: 150px; height: 150px; margin: 20px 0; place-items: center; border: 11px solid #edf0f6; border-radius: 50%; }.large-ring > div { display: grid; place-items: center; }.large-ring strong { font-size: 34px; }.large-ring span { color: var(--text-secondary); }.completeness-card ul { display: grid; width: 100%; gap: 11px; list-style: none; color: var(--text-regular); }.completeness-card li { display: flex; align-items: center; gap: 9px; }.completeness-card li i { width: 17px; height: 17px; border: 2px solid #d2d9e8; border-radius: 50%; }
.advice-card,.privacy-card { display: flex; padding: 20px; gap: 12px; }.advice-icon { display: grid; flex: 0 0 38px; width: 38px; height: 38px; place-items: center; border-radius: 12px; color: var(--color-secondary); background: #f1eaff; }.advice-card h3,.privacy-card h3 { margin-bottom: 5px; }.advice-card p,.privacy-card p { margin: 0; color: var(--text-secondary); font-size: 13px; }.privacy-card > .el-icon { flex: 0 0 auto; margin-top: 3px; color: var(--color-primary); }
@media (max-width: 1100px) { .profile-layout { grid-template-columns: 1fr; }.profile-aside { grid-template-columns: repeat(2,minmax(0,1fr)); }.completeness-card { grid-row: span 2; } }
@media (max-width: 700px) { .stage-alert { width: 100%; }.identity-card { padding: 20px; flex-wrap: wrap; }.mini-ring { margin-left: 0; }.info-grid { grid-template-columns: 1fr; }.profile-aside { grid-template-columns: 1fr; }.profile-actions { grid-template-columns: 1fr; } }
</style>

