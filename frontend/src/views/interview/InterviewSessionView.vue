<script setup lang="ts">
import { ChatDotRound, InfoFilled, Microphone, Promotion } from '@element-plus/icons-vue'
import { computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'

const route = useRoute()
const router = useRouter()
const sessionLabel = computed(() => String(route.params.id || 'preview'))
</script>

<template>
  <div class="session-page">
    <div class="session-topbar">
      <el-button text @click="router.push('/app/interviews')">← 返回模拟面试</el-button>
      <strong>模拟面试会话骨架</strong>
      <el-button disabled>结束面试</el-button>
    </div>

    <div class="session-content">
      <section class="mode-strip panel">
        <div class="mode selected"><span><el-icon><Microphone /></el-icon></span><div><strong>技术面</strong><small>考察专业技术能力</small></div></div>
        <div class="mode"><span><el-icon><ChatDotRound /></el-icon></span><div><strong>项目深挖</strong><small>深入项目细节和思路</small></div></div>
        <div class="progress-area"><small>面试进度</small><strong>第 0 / 5 轮</strong><el-progress :percentage="0" :show-text="false" /></div>
      </section>

      <section class="job-strip panel">
        <div><small>会话标识</small><strong>{{ sessionLabel }}</strong></div>
        <div><small>应聘岗位</small><strong>尚未选择 JD</strong></div>
        <div><small>面试状态</small><strong class="waiting">等待后续阶段接入</strong></div>
      </section>

      <div class="conversation-layout">
        <section class="panel conversation-panel">
          <div class="conversation-empty">
            <span class="bot-avatar">✦</span>
            <div>
              <h2>面试对话区</h2>
              <p>真实问题、用户回答和逐轮追问会在后续阶段显示在这里。</p>
              <el-alert title="当前不会向 AI 发送任何请求，也不会保存聊天内容。" type="info" :closable="false" show-icon />
            </div>
          </div>
          <div class="message-skeleton" aria-hidden="true"><i /><i /><i /></div>
          <div class="answer-box">
            <el-input type="textarea" :rows="3" disabled placeholder="面试功能将在后续阶段接入" />
            <el-button class="gradient-button" :icon="Promotion" disabled>发送</el-button>
          </div>
        </section>

        <aside class="feedback-panel">
          <section class="panel score-card">
            <div class="panel-header"><h2 class="panel-title">面试反馈</h2><span class="muted">等待回答</span></div>
            <div class="score"><strong>--</strong><span>/100</span><i /></div>
          </section>
          <section class="panel feedback-card"><h3>回答优点</h3><p>提交第一轮回答后显示具体优点。</p></section>
          <section class="panel feedback-card improve"><h3>待改进点</h3><p>提交第一轮回答后显示改进建议。</p></section>
          <section class="panel next-card"><el-icon><InfoFilled /></el-icon><div><h3>下一轮建议</h3><p>面试功能将在后续阶段接入。</p></div></section>
        </aside>
      </div>
    </div>
  </div>
</template>

<style scoped>
.session-page { min-height: calc(100vh - var(--topbar-height)); padding-bottom: 30px; }
.session-topbar { display: grid; height: 64px; padding: 0 28px; align-items: center; grid-template-columns: 1fr auto 1fr; border-bottom: 1px solid var(--border-color); background: rgba(255,255,255,.82); }
.session-topbar .el-button:last-child { justify-self: end; }
.session-content { display: grid; width: min(100%,1480px); margin: 0 auto; padding: 18px 28px; gap: 16px; }
.mode-strip { display: grid; padding: 14px 18px; align-items: center; grid-template-columns: 220px 220px 1fr; gap: 14px; }
.mode { display: flex; padding: 10px; align-items: center; gap: 10px; border: 1px solid transparent; border-radius: 12px; color: var(--text-secondary); }
.mode.selected { border-color: #adc0ff; color: var(--color-primary); background: #f8faff; }
.mode > span { display: grid; width: 38px; height: 38px; place-items: center; border-radius: 10px; background: #edf3ff; }
.mode div,.progress-area { display: grid; }
.mode small,.progress-area small,.job-strip small { color: var(--text-secondary); }
.progress-area { width: min(100%,400px); justify-self: end; }
.job-strip { display: grid; padding: 16px 20px; grid-template-columns: repeat(3,1fr); gap: 20px; }
.job-strip > div { display: grid; gap: 4px; }
.waiting { color: var(--color-warning); }
.conversation-layout { display: grid; grid-template-columns: minmax(0,1fr) 330px; gap: 16px; }
.conversation-panel { display: flex; min-height: 600px; padding: 22px; flex-direction: column; }
.conversation-empty { display: flex; max-width: 680px; padding: 18px; gap: 14px; }
.conversation-empty h2 { margin-bottom: 5px; font-size: 18px; }
.conversation-empty p { color: var(--text-secondary); }
.bot-avatar { display: grid; flex: 0 0 44px; width: 44px; height: 44px; place-items: center; border-radius: 14px; color: #fff; background: linear-gradient(140deg,var(--color-primary),var(--color-secondary)); font-size: 22px; }
.message-skeleton { display: grid; width: 56%; margin: 30px 20px; padding: 18px; gap: 12px; border-radius: 14px; background: #f6f8fd; }
.message-skeleton i { width: 88%; height: 10px; border-radius: 99px; background: #e6eaf3; }
.message-skeleton i:nth-child(2) { width: 100%; }.message-skeleton i:nth-child(3) { width: 62%; }
.answer-box { display: flex; margin-top: auto; padding: 12px; align-items: flex-end; gap: 10px; border: 1px solid #dce3f2; border-radius: 14px; }
.answer-box .el-input { flex: 1; }.answer-box :deep(.el-textarea__inner) { border: 0; box-shadow: none; resize: none; }
.feedback-panel { display: grid; align-content: start; gap: 12px; }
.score-card,.feedback-card,.next-card { padding: 18px; }
.score { display: flex; align-items: baseline; gap: 5px; }
.score strong { font-size: 34px; }.score span { color: var(--text-secondary); }.score i { width: 58px; height: 58px; margin-left: auto; border: 7px solid #edf0f6; border-radius: 50%; }
.feedback-card h3,.next-card h3 { margin-bottom: 6px; font-size: 15px; }.feedback-card p,.next-card p { margin: 0; color: var(--text-secondary); font-size: 13px; }
.feedback-card h3::before { margin-right: 7px; color: var(--color-success); content: '✓'; }.feedback-card.improve h3::before { color: var(--color-warning); content: '◇'; }
.next-card { display: flex; gap: 10px; color: var(--color-warning); }.next-card p { color: var(--text-secondary); }
@media (max-width: 1199px) { .conversation-layout { grid-template-columns: 1fr; } .feedback-panel { grid-template-columns: repeat(2,minmax(0,1fr)); } }
@media (max-width: 767px) { .session-topbar { height: auto; min-height: 62px; padding: 8px 12px; grid-template-columns: 1fr auto; }.session-topbar strong { display: none; }.session-content { padding: 12px; }.mode-strip { grid-template-columns: 1fr 1fr; }.progress-area { grid-column: 1/-1; width: 100%; justify-self: stretch; }.job-strip { grid-template-columns: 1fr; }.conversation-panel { min-height: 520px; padding: 14px; }.conversation-empty { padding: 8px 2px; }.message-skeleton { width: 86%; margin-inline: 4px; }.answer-box { flex-direction: column; align-items: stretch; }.feedback-panel { grid-template-columns: 1fr; } }
</style>

