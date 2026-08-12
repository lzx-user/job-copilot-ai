<script setup lang="ts">
import { Calendar, DocumentChecked, Microphone, Right, TrendCharts, User } from '@element-plus/icons-vue'
import type { Component } from 'vue'
import { computed } from 'vue'
import { useRouter } from 'vue-router'
import AppEmpty from '../../components/common/AppEmpty.vue'
import { useAuthStore } from '../../stores/auth.store'

interface StatItem {
  label: string
  value: number
  hint: string
  icon: Component
  tone: 'blue' | 'violet' | 'green' | 'orange'
}

const router = useRouter()
const authStore = useAuthStore()
const displayName = computed(() => authStore.user?.email?.split('@')[0] || '求职同学')

const stats: StatItem[] = [
  { label: '已分析 JD', value: 0, hint: '等待首次分析', icon: DocumentChecked, tone: 'blue' },
  { label: '模拟面试次数', value: 0, hint: '等待首次练习', icon: Microphone, tone: 'violet' },
  { label: '本周新增记录', value: 0, hint: '从今天开始积累', icon: TrendCharts, tone: 'green' },
  { label: '档案完整度', value: 0, hint: '下一阶段可保存', icon: User, tone: 'orange' },
]

const plans = [
  { title: '了解个人档案需要准备的内容', status: '建议', path: '/app/profile' },
  { title: '准备一份目标岗位 JD', status: '待准备', path: '/app/jd-analysis' },
  { title: '浏览模拟面试流程', status: '待了解', path: '/app/interviews' },
]
</script>

<template>
  <div class="page-shell dashboard-page">
    <section class="hero">
      <div class="hero-copy">
        <span class="greeting">你好，{{ displayName }} <span aria-hidden="true">👋</span></span>
        <h1>今天继续冲刺 <em>Offer</em></h1>
        <p>AI 全程陪伴你的求职之路，让每一次准备都有方向。</p>
      </div>
      <div class="hero-visual" aria-hidden="true">
        <div class="offer-card"><strong>OFFER</strong><i /><i /><span>✓</span></div>
        <div class="target">◎</div>
        <b>✦</b>
      </div>
    </section>

    <section class="stats-grid" aria-label="求职数据概览">
      <article v-for="stat in stats" :key="stat.label" class="panel stat-card">
        <span class="stat-icon" :class="stat.tone"><el-icon :size="23"><component :is="stat.icon" /></el-icon></span>
        <div><span>{{ stat.label }}</span><strong>{{ stat.value }}</strong><small>{{ stat.hint }}</small></div>
      </article>
    </section>

    <div class="dashboard-grid">
      <section class="main-column">
        <div class="shortcut-grid">
          <article class="shortcut shortcut-analysis">
            <div>
              <span class="section-label">岗位准备</span>
              <h2>开始 JD 分析</h2>
              <p>拆解岗位要求，提前整理你的能力证据。</p>
              <ul><li>岗位需求拆解</li><li>技能关键词整理</li><li>匹配度评估建议</li></ul>
              <el-button class="gradient-button" @click="router.push('/app/jd-analysis')">查看分析页 <el-icon><Right /></el-icon></el-button>
            </div>
            <div class="shortcut-icon"><el-icon><DocumentChecked /></el-icon></div>
          </article>

          <article class="shortcut shortcut-interview">
            <div>
              <span class="section-label">面试准备</span>
              <h2>开始模拟面试</h2>
              <p>熟悉练习流程，为下一阶段的真实对话做好准备。</p>
              <ul><li>三类面试模式</li><li>多轮问答流程</li><li>逐轮反馈骨架</li></ul>
              <el-button class="gradient-button" @click="router.push('/app/interviews')">查看面试页 <el-icon><Right /></el-icon></el-button>
            </div>
            <div class="shortcut-icon"><el-icon><Microphone /></el-icon></div>
          </article>
        </div>

        <section class="panel recent-panel">
          <div class="panel-header">
            <h2 class="panel-title">最近记录</h2>
            <el-button text type="primary" @click="router.push('/app/history')">查看全部</el-button>
          </div>
          <AppEmpty title="还没有求职记录" description="完成首次 JD 分析或模拟面试后，最近记录会显示在这里。" icon="⌁">
            <el-button type="primary" plain @click="router.push('/app/jd-analysis')">准备第一份 JD</el-button>
          </AppEmpty>
        </section>
      </section>

      <aside class="side-column">
        <section class="panel plan-panel">
          <div class="panel-header">
            <h2 class="panel-title"><el-icon><Calendar /></el-icon> 今日计划</h2>
            <span class="muted">0 / 3</span>
          </div>
          <button v-for="plan in plans" :key="plan.title" class="plan-row" @click="router.push(plan.path)">
            <span class="plan-check" />
            <span>{{ plan.title }}</span>
            <small>{{ plan.status }}</small>
          </button>
        </section>

        <section class="panel suggestion-panel">
          <span class="suggestion-icon">✦</span>
          <div><h3>新手建议</h3><p>先准备目标岗位与项目摘要，下一阶段完善档案后，分析会更有上下文。</p></div>
        </section>

        <section class="panel stage-panel">
          <span class="stage-mark">01</span>
          <div><strong>当前阶段</strong><p>认证基础与静态页面骨架</p></div>
        </section>
      </aside>
    </div>
  </div>
</template>

<style scoped>
.dashboard-page { padding-top: 20px; }

.hero {
  position: relative;
  display: flex;
  min-height: 190px;
  padding: 24px 30px;
  align-items: center;
  overflow: hidden;
  border: 1px solid rgba(225, 231, 248, 0.8);
  border-radius: var(--radius-lg);
  background:
    radial-gradient(circle at 72% 20%, rgba(103, 126, 249, 0.16), transparent 23%),
    linear-gradient(125deg, rgba(255,255,255,.98), rgba(247,248,255,.92));
  box-shadow: var(--shadow-card);
}

.hero-copy { position: relative; z-index: 2; }
.greeting { color: var(--text-regular); font-weight: 600; }
.hero h1 { margin: 12px 0 5px; font-size: clamp(30px, 3vw, 44px); letter-spacing: -0.04em; }
.hero h1 em { color: var(--color-primary); font-style: normal; }
.hero p { margin: 0; color: var(--text-secondary); }

.hero-visual { position: absolute; right: 8%; width: 280px; height: 165px; }
.offer-card { position: absolute; top: 10px; left: 50px; display: grid; width: 150px; height: 145px; padding: 22px; transform: rotate(8deg); border: 9px solid rgba(196,208,255,.45); border-radius: 24px; color: #8092de; background: rgba(255,255,255,.72); box-shadow: 0 20px 40px rgba(75,99,194,.18); }
.offer-card i { width: 70px; height: 7px; border-radius: 99px; background: #d8e0fa; }
.offer-card span { position: absolute; bottom: -14px; left: -20px; display: grid; width: 52px; height: 52px; place-items: center; border-radius: 50%; color: var(--color-primary); background: #fff; box-shadow: 0 12px 25px rgba(65,105,246,.22); font-size: 26px; }
.target { position: absolute; right: 0; bottom: 15px; color: var(--color-primary); font-size: 58px; }
.hero-visual b { position: absolute; top: 4px; right: 16px; color: #a7b6f4; font-size: 28px; }

.stats-grid { display: grid; margin: 18px 0; grid-template-columns: repeat(4, minmax(0,1fr)); gap: 16px; }
.stat-card { display: flex; min-height: 110px; padding: 20px; align-items: center; gap: 15px; }
.stat-card > div { display: grid; }
.stat-card span { color: var(--text-regular); }
.stat-card strong { font-size: 27px; line-height: 1.2; }
.stat-card small { color: var(--text-secondary); }
.stat-icon { display: grid; flex: 0 0 46px; width: 46px; height: 46px; place-items: center; border-radius: 14px; color: #fff !important; }
.stat-icon.blue { background: linear-gradient(140deg,#72adff,#4169f6); }
.stat-icon.violet { background: linear-gradient(140deg,#ab7dff,#7044ef); }
.stat-icon.green { background: linear-gradient(140deg,#54d7b2,#16a87a); }
.stat-icon.orange { background: linear-gradient(140deg,#ffc454,#f49c22); }

.dashboard-grid { display: grid; grid-template-columns: minmax(0, 1.75fr) minmax(290px, .75fr); gap: 18px; }
.main-column,.side-column { display: grid; align-content: start; gap: 18px; }
.shortcut-grid { display: grid; grid-template-columns: repeat(2,minmax(0,1fr)); gap: 18px; }
.shortcut { position: relative; display: flex; min-height: 280px; padding: 28px; overflow: hidden; border: 1px solid #dce5fc; border-radius: var(--radius-lg); }
.shortcut-analysis { background: linear-gradient(140deg,#f9fbff,#edf4ff); }
.shortcut-interview { border-color: #e7defc; background: linear-gradient(140deg,#fdfbff,#f4efff); }
.shortcut > div:first-child { position: relative; z-index: 2; }
.shortcut h2 { margin-bottom: 8px; font-size: 23px; }
.shortcut p { color: var(--text-regular); }
.shortcut ul { display: grid; margin: 18px 0 22px; gap: 7px; list-style: none; color: var(--text-secondary); }
.shortcut li::before { margin-right: 8px; color: var(--color-primary); content: '✓'; }
.shortcut-icon { position: absolute; right: 20px; bottom: 40px; display: grid; width: 82px; height: 98px; place-items: center; transform: rotate(8deg); border: 6px solid rgba(255,255,255,.72); border-radius: 22px; color: var(--color-primary); background: rgba(255,255,255,.75); box-shadow: 0 20px 40px rgba(65,105,246,.18); font-size: 42px; }
.shortcut-interview .shortcut-icon { color: var(--color-secondary); }
.recent-panel,.plan-panel { padding: 22px; }
.recent-panel :deep(.empty-state) { min-height: 230px; border: 1px dashed #dbe2f2; border-radius: 14px; background: #fbfcff; }
.plan-panel .panel-title { display: flex; align-items: center; gap: 8px; }
.plan-row { display: grid; width: 100%; min-height: 54px; padding: 0 2px; cursor: pointer; align-items: center; grid-template-columns: 20px 1fr auto; gap: 9px; border-bottom: 1px solid #edf0f6; text-align: left; color: var(--text-regular); background: transparent; }
.plan-row:last-child { border-bottom: 0; }
.plan-row:hover > span:nth-child(2) { color: var(--color-primary); }
.plan-check { width: 17px; height: 17px; border: 1.5px solid #cdd5e6; border-radius: 5px; }
.plan-row small { color: var(--text-secondary); }
.suggestion-panel,.stage-panel { display: flex; padding: 20px; gap: 14px; }
.suggestion-icon { display: grid; flex: 0 0 38px; width: 38px; height: 38px; place-items: center; border-radius: 12px; color: var(--color-secondary); background: #f0eaff; }
.suggestion-panel h3 { margin-bottom: 5px; }
.suggestion-panel p,.stage-panel p { margin: 0; color: var(--text-secondary); font-size: 13px; }
.stage-mark { display: grid; flex: 0 0 44px; width: 44px; height: 44px; place-items: center; border-radius: 13px; color: #fff; background: linear-gradient(140deg,var(--color-primary),var(--color-secondary)); font-weight: 800; }

@media (max-width: 1199px) {
  .stats-grid { grid-template-columns: repeat(2,minmax(0,1fr)); }
  .dashboard-grid { grid-template-columns: 1fr; }
  .side-column { grid-template-columns: repeat(2,minmax(0,1fr)); }
  .plan-panel { grid-row: span 2; }
}

@media (max-width: 767px) {
  .hero { min-height: 175px; padding: 24px 20px; }
  .hero-visual { right: -120px; opacity: .35; }
  .stats-grid { gap: 10px; }
  .stat-card { min-height: 96px; padding: 14px; }
  .stat-icon { flex-basis: 40px; width: 40px; height: 40px; }
  .shortcut-grid,.side-column { grid-template-columns: 1fr; }
  .shortcut { min-height: 265px; padding: 22px; }
  .shortcut-icon { opacity: .45; }
}

@media (max-width: 460px) {
  .stats-grid { grid-template-columns: 1fr; }
}
</style>

