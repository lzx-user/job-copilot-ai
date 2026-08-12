<script setup lang="ts">
import { Briefcase, ChatDotRound, Microphone, Reading } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'
import type { Component } from 'vue'
import AppEmpty from '../../components/common/AppEmpty.vue'

interface InterviewMode {
  title: string
  description: string
  icon: Component
  tone: string
}

const modes: InterviewMode[] = [
  { title: '技术面', description: '围绕岗位技能与基础知识进行针对性练习', icon: Microphone, tone: 'blue' },
  { title: '项目深挖', description: '练习项目背景、难点、取舍和个人贡献', icon: Briefcase, tone: 'violet' },
  { title: '综合面', description: '兼顾经历表达、岗位动机与协作思考', icon: ChatDotRound, tone: 'green' },
]

function explainAvailability() {
  ElMessage.info('模拟面试业务将在 7 月 22 日之后接入，当前只展示流程骨架。')
}
</script>

<template>
  <div class="page-shell">
    <div class="page-heading">
      <div><span class="section-label">MOCK INTERVIEW</span><h1>模拟面试</h1><p>先熟悉练习模式，后续可基于真实 JD 开启多轮文字面试。</p></div>
      <el-tag effect="plain" round>固定 5 轮 · 后续阶段开放</el-tag>
    </div>

    <section class="panel intro-panel">
      <div class="intro-copy">
        <span class="intro-badge"><el-icon><Reading /></el-icon> 面试准备流程</span>
        <h2>选择合适的练习方式，<br /><em>每次只聚焦一个目标。</em></h2>
        <p>真实面试功能会在后续阶段结合个人档案与岗位 JD 接入。当前页面不生成问题，也不保存会话。</p>
      </div>
      <div class="interview-visual" aria-hidden="true"><span>AI</span><i /><b>•••</b></div>
    </section>

    <section class="mode-grid">
      <article v-for="mode in modes" :key="mode.title" class="panel mode-card" :class="mode.tone">
        <span class="mode-icon"><el-icon :size="25"><component :is="mode.icon" /></el-icon></span>
        <div><h3>{{ mode.title }}</h3><p>{{ mode.description }}</p></div>
        <el-radio :model-value="''" :value="mode.title" disabled>选择</el-radio>
      </article>
    </section>

    <section class="panel jd-selector">
      <div class="panel-header"><div><h2 class="panel-title">选择目标岗位</h2><p>后续将从已保存的 JD 分析记录中选择。</p></div><el-tag type="info">暂无记录</el-tag></div>
      <AppEmpty title="还没有可用于面试的 JD" description="当前阶段不会创建虚假岗位。后续完成一份真实 JD 分析后，可以从这里开始面试。" icon="▤">
        <el-button class="gradient-button" @click="explainAvailability">开始面试（后续开放）</el-button>
      </AppEmpty>
    </section>
  </div>
</template>

<style scoped>
.intro-panel { position: relative; min-height: 230px; padding: 34px; overflow: hidden; background: linear-gradient(125deg,#fff,#f1f5ff 60%,#f5efff); }
.intro-copy { position: relative; z-index: 2; max-width: 670px; }
.intro-badge { display: inline-flex; padding: 6px 10px; align-items: center; gap: 6px; border-radius: 9px; color: var(--color-primary); background: #edf3ff; font-weight: 600; }
.intro-panel h2 { margin: 18px 0 10px; font-size: clamp(26px,3vw,37px); line-height: 1.3; }
.intro-panel h2 em { color: var(--color-primary); font-style: normal; }
.intro-panel p { max-width: 600px; margin: 0; color: var(--text-regular); }
.interview-visual { position: absolute; right: 8%; bottom: -30px; display: grid; width: 180px; height: 180px; place-items: center; border: 12px solid rgba(255,255,255,.65); border-radius: 44px; color: #fff; background: linear-gradient(140deg,#6685ff,#8a55f3); box-shadow: 0 28px 55px rgba(100,75,219,.28); font-size: 44px; font-weight: 800; transform: rotate(7deg); }
.interview-visual i { position: absolute; top: -15px; right: -15px; width: 45px; height: 45px; border: 8px solid #dfe6ff; border-radius: 50%; }
.interview-visual b { position: absolute; bottom: 28px; font-size: 20px; letter-spacing: 4px; }
.mode-grid { display: grid; margin: 18px 0; grid-template-columns: repeat(3,minmax(0,1fr)); gap: 16px; }
.mode-card { display: grid; min-height: 140px; padding: 20px; align-items: center; grid-template-columns: 50px 1fr auto; gap: 14px; }
.mode-card h3 { margin-bottom: 4px; }
.mode-card p { margin: 0; color: var(--text-secondary); font-size: 13px; }
.mode-icon { display: grid; width: 48px; height: 48px; place-items: center; border-radius: 14px; color: var(--color-primary); background: #edf3ff; }
.mode-card.violet .mode-icon { color: var(--color-secondary); background: #f1ebff; }
.mode-card.green .mode-icon { color: var(--color-success); background: #e9faf4; }
.jd-selector { padding: 24px; }
.panel-header p { margin: 5px 0 0; color: var(--text-secondary); }
.jd-selector :deep(.empty-state) { min-height: 270px; border: 1px dashed #dbe2f2; border-radius: 14px; background: #fbfcff; }
@media (max-width: 1000px) { .mode-grid { grid-template-columns: 1fr; } }
@media (max-width: 767px) { .intro-panel { min-height: 250px; padding: 24px; } .interview-visual { right: -70px; opacity: .35; } .mode-card { grid-template-columns: 50px 1fr; } .mode-card .el-radio { grid-column: 2; } .jd-selector { padding: 18px; } }
</style>

