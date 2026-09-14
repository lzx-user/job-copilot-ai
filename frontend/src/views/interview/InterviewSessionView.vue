<script setup lang="ts">
import { InfoFilled, Promotion } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'
import { computed, nextTick, onMounted, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { generateInterviewReport, getInterviewSession, submitInterviewTurn } from '../../api/interview.api'
import type { InterviewFeedback, InterviewSessionDetail } from '../../types/interview'
import { toUserMessage } from '../../utils/error'

const route = useRoute()
const router = useRouter()
const sessionId = computed(() => String(route.params.id || ''))
const session = ref<InterviewSessionDetail | null>(null)
const answer = ref('')
const loading = ref(true)
const submitting = ref(false)
const generatingReport = ref(false)
const loadError = ref('')
const conversationRef = ref<HTMLElement | null>(null)

const latestFeedback = computed<InterviewFeedback | null>(() => {
  const candidateMessages = session.value?.messages.filter((message) => message.role === 'candidate' && message.score !== undefined) ?? []
  const message = candidateMessages.at(-1)
  if (!message || message.score === undefined || !message.feedback || !message.strengths || !message.improvements) return null
  return { score: message.score, feedback: message.feedback, strengths: message.strengths, improvements: message.improvements }
})

const currentRoundAnswered = computed(() => Boolean(session.value?.messages.some(
  (message) => message.role === 'candidate' && message.round === session.value?.currentRound,
)))

const canSubmit = computed(() => Boolean(
  session.value?.status === 'in_progress' &&
  !currentRoundAnswered.value &&
  answer.value.trim() && !submitting.value,
))

const canGenerateReport = computed(() => Boolean(
  session.value?.status === 'in_progress' &&
  session.value.currentRound === session.value.maxRounds &&
  currentRoundAnswered.value && !generatingReport.value,
))

async function loadSession() {
  loading.value = true
  loadError.value = ''
  try {
    session.value = await getInterviewSession(sessionId.value)
    await scrollToBottom()
  } catch (error) {
    loadError.value = toUserMessage(error)
  } finally {
    loading.value = false
  }
}

async function handleSubmit() {
  const content = answer.value.trim()
  if (!canSubmit.value || !content) return
  submitting.value = true
  try {
    await submitInterviewTurn(sessionId.value, content)
    answer.value = ''
    await loadSession()
    ElMessage.success('本轮回答与反馈已保存')
  } catch (error) {
    ElMessage.error(toUserMessage(error))
  } finally {
    submitting.value = false
  }
}

async function handleGenerateReport() {
  if (!canGenerateReport.value) return
  generatingReport.value = true
  try {
    await generateInterviewReport(sessionId.value)
    await loadSession()
    ElMessage.success('最终报告已生成并保存')
  } catch (error) {
    ElMessage.error(toUserMessage(error))
  } finally {
    generatingReport.value = false
  }
}

async function scrollToBottom() {
  await nextTick()
  conversationRef.value?.scrollTo({ top: conversationRef.value.scrollHeight, behavior: 'smooth' })
}

onMounted(loadSession)
</script>

<template>
  <div class="session-page">
    <div class="session-topbar">
      <el-button text @click="router.push('/app/interviews')">← 返回模拟面试</el-button>
      <strong>{{ session ? `${session.companyName} · ${session.jobTitle}` : '模拟面试' }}</strong>
      <el-button :loading="loading" @click="loadSession">刷新</el-button>
    </div>

    <div class="session-content" v-loading="loading">
      <el-alert v-if="loadError" :title="loadError" type="error" :closable="false" show-icon />
      <template v-if="session">
        <section class="job-strip panel">
          <div><small>目标公司</small><strong>{{ session.companyName }}</strong></div>
          <div><small>应聘岗位</small><strong>{{ session.jobTitle }}</strong></div>
          <div><small>面试进度</small><strong>第 {{ session.currentRound }} / {{ session.maxRounds }} 轮</strong><el-progress :percentage="session.currentRound / session.maxRounds * 100" :show-text="false" /></div>
        </section>

        <div class="conversation-layout">
          <section class="panel conversation-panel">
            <div ref="conversationRef" class="messages">
              <article v-for="message in session.messages" :key="message.id" class="message" :class="message.role">
                <span class="avatar">{{ message.role === 'interviewer' ? 'AI' : '我' }}</span>
                <div><small>第 {{ message.round }} 轮 · {{ message.role === 'interviewer' ? '面试官' : '候选人' }}</small><p>{{ message.content }}</p></div>
              </article>
            </div>
            <el-alert v-if="session.status === 'completed'" title="五轮面试已完成，最终报告已保存。" type="success" :closable="false" show-icon />
            <el-alert v-else-if="currentRoundAnswered && session.currentRound === session.maxRounds" title="5 轮问答已全部完成并保存，现在可以生成最终报告。" type="success" :closable="false" show-icon />
            <el-alert v-else-if="session.currentRound === session.maxRounds" title="已进入第 5 轮，提交本轮回答后将不再生成下一题。" type="info" :closable="false" show-icon />
            <div class="answer-box">
              <el-input v-model="answer" type="textarea" :rows="4" maxlength="10000" show-word-limit :disabled="session.status !== 'in_progress' || currentRoundAnswered" placeholder="输入本轮回答，建议用 STAR 结构说明背景、行动和结果" @keydown.ctrl.enter="handleSubmit" />
              <el-button class="gradient-button" :icon="Promotion" :loading="submitting" :disabled="!canSubmit" @click="handleSubmit">提交回答</el-button>
            </div>
          </section>

          <aside class="feedback-panel">
            <section class="panel score-card">
              <div class="panel-header"><h2 class="panel-title">最近反馈</h2><span class="muted">{{ latestFeedback ? '已保存' : '等待回答' }}</span></div>
              <div class="score"><strong>{{ latestFeedback?.score ?? '--' }}</strong><span>/100</span></div>
              <p v-if="latestFeedback">{{ latestFeedback.feedback }}</p>
            </section>
            <section class="panel feedback-card"><h3>回答优点</h3><ul v-if="latestFeedback"><li v-for="item in latestFeedback.strengths" :key="item">{{ item }}</li></ul><p v-else>提交回答后显示。</p></section>
            <section class="panel feedback-card improve"><h3>待改进点</h3><ul v-if="latestFeedback"><li v-for="item in latestFeedback.improvements" :key="item">{{ item }}</li></ul><p v-else>提交回答后显示。</p></section>
            <section class="panel next-card"><el-icon><InfoFilled /></el-icon><div><h3>提交说明</h3><p>发送期间会禁用按钮；后端同时校验轮次，重复提交不会推进两次。</p></div></section>
          </aside>
        </div>

        <section v-if="canGenerateReport" class="panel report-action">
          <div><h2>完成面试复盘</h2><p>AI 将基于岗位信息和完整五轮问答生成结构化报告；报告保存成功后，本场面试才会标记完成。</p></div>
          <el-button class="gradient-button" :loading="generatingReport" @click="handleGenerateReport">生成最终报告</el-button>
        </section>

        <section v-if="session.report" class="panel report-panel">
          <div class="panel-header"><div><span class="section-label">FINAL REPORT</span><h2 class="panel-title">最终面试报告</h2></div><el-tag type="success" effect="plain">已保存</el-tag></div>
          <div class="report-scores">
            <div><strong>{{ session.report.overallScore }}</strong><span>综合评分</span></div>
            <div><strong>{{ session.report.technicalScore }}</strong><span>技术能力</span></div>
            <div><strong>{{ session.report.expressionScore }}</strong><span>表达能力</span></div>
            <div><strong>{{ session.report.projectDepthScore }}</strong><span>项目深度</span></div>
          </div>
          <p class="report-summary">{{ session.report.summary }}</p>
          <div class="report-details">
            <div><h3>突出优势</h3><ul><li v-for="item in session.report.strengths" :key="item">{{ item }}</li></ul></div>
            <div><h3>主要不足</h3><ul><li v-for="item in session.report.weaknesses" :key="item">{{ item }}</li></ul></div>
            <div><h3>建议准备主题</h3><ul><li v-for="item in session.report.recommendedTopics" :key="item">{{ item }}</li></ul></div>
            <div><h3>回答技巧</h3><ul><li v-for="item in session.report.answerTips" :key="item">{{ item }}</li></ul></div>
          </div>
        </section>
      </template>
    </div>
  </div>
</template>

<style scoped>
.session-page { min-height: calc(100vh - var(--topbar-height)); padding-bottom: 30px; }.session-topbar{display:grid;min-height:64px;padding:8px 28px;align-items:center;grid-template-columns:1fr auto 1fr;border-bottom:1px solid var(--border-color);background:rgba(255,255,255,.9)}.session-topbar .el-button:last-child{justify-self:end}.session-content{display:grid;width:min(100%,1480px);min-height:680px;margin:0 auto;padding:18px 28px;gap:16px}.job-strip{display:grid;padding:16px 20px;grid-template-columns:repeat(3,1fr);gap:20px}.job-strip>div{display:grid;gap:5px}.job-strip small,.message small,.muted{color:var(--text-secondary)}.conversation-layout{display:grid;grid-template-columns:minmax(0,1fr) 330px;gap:16px}.conversation-panel{display:flex;min-height:590px;padding:20px;flex-direction:column;gap:14px}.messages{display:flex;max-height:470px;padding:4px;overflow:auto;flex-direction:column;gap:18px}.message{display:flex;max-width:82%;gap:10px}.message.candidate{align-self:flex-end;flex-direction:row-reverse}.avatar{display:grid;flex:0 0 40px;width:40px;height:40px;place-items:center;border-radius:12px;color:#fff;background:linear-gradient(140deg,var(--color-primary),var(--color-secondary));font-size:12px;font-weight:700}.candidate .avatar{background:#8290aa}.message>div{padding:13px 15px;border-radius:4px 15px 15px;background:#f3f6fd}.candidate>div{border-radius:15px 4px 15px;background:#edf3ff}.message p{margin:5px 0 0;white-space:pre-wrap;line-height:1.7}.answer-box{display:flex;margin-top:auto;padding:12px;align-items:flex-end;gap:10px;border:1px solid #dce3f2;border-radius:14px}.answer-box .el-input{flex:1}.answer-box :deep(.el-textarea__inner){border:0;box-shadow:none;resize:none}.feedback-panel{display:grid;align-content:start;gap:12px}.score-card,.feedback-card,.next-card{padding:18px}.score{display:flex;align-items:baseline;gap:5px}.score strong{font-size:34px;color:var(--color-primary)}.score span,.score-card p,.feedback-card p,.next-card p{color:var(--text-secondary)}.feedback-card h3,.next-card h3{margin-bottom:8px;font-size:15px}.feedback-card ul{margin:0;padding-left:20px}.feedback-card li{margin:7px 0;color:var(--text-regular);line-height:1.5}.feedback-card h3::before{margin-right:7px;color:var(--color-success);content:'✓'}.feedback-card.improve h3::before{color:var(--color-warning);content:'◇'}.next-card{display:flex;gap:10px;color:var(--color-warning)}.next-card p{margin:0}.report-action{display:flex;padding:20px;align-items:center;justify-content:space-between;gap:20px}.report-action h2{margin:0 0 6px}.report-action p,.report-summary{margin:0;color:var(--text-secondary);line-height:1.7}.report-panel{display:grid;padding:24px;gap:20px}.report-scores{display:grid;grid-template-columns:repeat(4,1fr);gap:12px}.report-scores>div{display:grid;padding:18px;place-items:center;gap:5px;border-radius:14px;background:#f5f7fc}.report-scores strong{font-size:30px;color:var(--color-primary)}.report-scores span{color:var(--text-secondary)}.report-summary{padding:16px;border-left:3px solid var(--color-primary);background:#f7f9fd}.report-details{display:grid;grid-template-columns:repeat(2,minmax(0,1fr));gap:14px}.report-details>div{padding:16px;border:1px solid var(--border-color);border-radius:14px}.report-details h3{margin:0 0 10px;font-size:15px}.report-details ul{margin:0;padding-left:20px}.report-details li{margin:7px 0;line-height:1.6}@media(max-width:1100px){.conversation-layout{grid-template-columns:1fr}.feedback-panel{grid-template-columns:repeat(2,minmax(0,1fr))}}@media(max-width:767px){.session-topbar{padding:8px 12px;grid-template-columns:1fr auto}.session-topbar strong{display:none}.session-content{padding:12px}.job-strip{grid-template-columns:1fr}.conversation-panel{min-height:520px;padding:14px}.message{max-width:94%}.answer-box,.report-action{flex-direction:column;align-items:stretch}.feedback-panel,.report-scores,.report-details{grid-template-columns:1fr}}
</style>
