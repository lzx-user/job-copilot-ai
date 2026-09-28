<script setup lang="ts">
import { Delete, Refresh, Search } from '@element-plus/icons-vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { computed, onMounted, ref } from 'vue'
import { useRouter } from 'vue-router'
import { deleteInterviewSession, listInterviewSessions } from '../../api/interview.api'
import { deleteJDAnalysis, getJDAnalysis, listJDAnalyses } from '../../api/jd-analysis.api'
import AppEmpty from '../../components/common/AppEmpty.vue'
import type { InterviewHistoryRecord } from '../../types/interview'
import type { JdAnalysisRecord } from '../../types/jd-analysis'
import { toUserMessage } from '../../utils/error'

type HistoryTab = 'jd' | 'interview'

const router = useRouter()
const activeTab = ref<HistoryTab>('jd')
const jdRecords = ref<JdAnalysisRecord[]>([])
const interviewRecords = ref<InterviewHistoryRecord[]>([])
const keyword = ref('')
const loading = ref(true)
const detailLoading = ref(false)
const deletingId = ref('')
const loadError = ref('')
const detail = ref<JdAnalysisRecord | null>(null)
const dialogVisible = ref(false)

const filteredJDRecords = computed(() => filterRecords(jdRecords.value))
const filteredInterviewRecords = computed(() => filterRecords(interviewRecords.value))

function filterRecords<T extends { companyName: string; jobTitle: string }>(records: T[]) {
  const value = keyword.value.trim().toLowerCase()
  if (!value) return records
  return records.filter((record) =>
    `${record.companyName} ${record.jobTitle}`.toLowerCase().includes(value),
  )
}

async function loadRecords() {
  loading.value = true
  loadError.value = ''
  try {
    const [jdResult, interviewResult] = await Promise.all([
      listJDAnalyses(),
      listInterviewSessions(),
    ])
    jdRecords.value = jdResult.items
    interviewRecords.value = interviewResult.items
  } catch (error) {
    loadError.value = toUserMessage(error)
  } finally {
    loading.value = false
  }
}

async function openJDDetail(record: JdAnalysisRecord) {
  dialogVisible.value = true
  detailLoading.value = true
  detail.value = null
  try {
    detail.value = await getJDAnalysis(record.analysisId)
  } catch (error) {
    ElMessage.error(toUserMessage(error))
    dialogVisible.value = false
  } finally {
    detailLoading.value = false
  }
}

async function removeJD(record: JdAnalysisRecord) {
  try {
    await ElMessageBox.confirm(
      `删除“${record.companyName} · ${record.jobTitle}”后，关联的面试记录和报告也会一并删除。此操作不可恢复。`,
      '确认删除 JD 分析',
      { type: 'warning', confirmButtonText: '确认删除', cancelButtonText: '取消' },
    )
    deletingId.value = record.analysisId
    await deleteJDAnalysis(record.analysisId)
    jdRecords.value = jdRecords.value.filter((item) => item.analysisId !== record.analysisId)
    interviewRecords.value = (await listInterviewSessions()).items
    ElMessage.success('JD 分析及关联面试已删除')
  } catch (error) {
    if (error !== 'cancel' && error !== 'close') ElMessage.error(toUserMessage(error))
  } finally {
    deletingId.value = ''
  }
}

async function removeInterview(record: InterviewHistoryRecord) {
  try {
    await ElMessageBox.confirm(
      `确认删除“${record.companyName} · ${record.jobTitle}”的本场面试、全部消息和报告吗？`,
      '确认删除面试',
      { type: 'warning', confirmButtonText: '确认删除', cancelButtonText: '取消' },
    )
    deletingId.value = record.sessionId
    await deleteInterviewSession(record.sessionId)
    interviewRecords.value = interviewRecords.value.filter((item) => item.sessionId !== record.sessionId)
    ElMessage.success('面试记录已删除')
  } catch (error) {
    if (error !== 'cancel' && error !== 'close') ElMessage.error(toUserMessage(error))
  } finally {
    deletingId.value = ''
  }
}

function formatDate(value: string) {
  return new Intl.DateTimeFormat('zh-CN', { dateStyle: 'medium', timeStyle: 'short' }).format(new Date(value))
}

function statusLabel(status: InterviewHistoryRecord['status']) {
  if (status === 'completed') return '已完成'
  if (status === 'in_progress') return '进行中'
  return '准备中'
}

onMounted(loadRecords)
</script>

<template>
  <div class="page-shell">
    <div class="page-heading">
      <div><span class="section-label">REVIEW</span><h1>历史记录</h1><p>查看当前账号保存的 JD 分析和模拟面试，重新登录后仍可恢复。</p></div>
      <el-tag effect="plain" round>{{ jdRecords.length + interviewRecords.length }} 条记录</el-tag>
    </div>

    <section class="panel history-panel">
      <div class="tabs">
        <button :class="{ active: activeTab === 'jd' }" @click="activeTab = 'jd'">JD 分析 <span>{{ jdRecords.length }}</span></button>
        <button :class="{ active: activeTab === 'interview' }" @click="activeTab = 'interview'">模拟面试 <span>{{ interviewRecords.length }}</span></button>
      </div>
      <div class="filters">
        <el-input v-model="keyword" :prefix-icon="Search" placeholder="搜索公司或岗位" clearable />
        <el-button :icon="Refresh" :loading="loading" @click="loadRecords">刷新</el-button>
      </div>
      <el-alert v-if="loadError" :title="loadError" type="error" :closable="false" show-icon />
      <div v-loading="loading" class="records-area">
        <template v-if="activeTab === 'jd'">
          <div v-if="filteredJDRecords.length" class="record-list">
            <div v-for="record in filteredJDRecords" :key="record.analysisId" class="record-row" role="button" tabindex="0" @click="openJDDetail(record)" @keydown.enter.self="openJDDetail(record)">
              <span><strong>{{ record.companyName }}</strong><small>{{ record.jobTitle }}</small></span>
              <span class="score">{{ record.matchScore }} 分</span>
              <time>{{ formatDate(record.createdAt) }}</time>
              <span class="record-actions"><em>查看详情</em><el-button text type="danger" :icon="Delete" :loading="deletingId === record.analysisId" @click.stop="removeJD(record)">删除</el-button></span>
            </div>
          </div>
          <AppEmpty v-else-if="!loading && !loadError" title="还没有 JD 分析记录" description="完成一份真实 JD 分析后，记录会保存并显示在这里。" icon="◷">
            <el-button type="primary" plain @click="router.push('/app/jd-analysis')">准备第一份 JD</el-button>
          </AppEmpty>
        </template>

        <template v-else>
          <div v-if="filteredInterviewRecords.length" class="record-list">
            <div v-for="record in filteredInterviewRecords" :key="record.sessionId" class="record-row interview-row" role="button" tabindex="0" @click="router.push(`/app/interviews/${record.sessionId}`)" @keydown.enter.self="router.push(`/app/interviews/${record.sessionId}`)">
              <span><strong>{{ record.companyName }}</strong><small>{{ record.jobTitle }}</small></span>
              <span><el-tag :type="record.status === 'completed' ? 'success' : 'warning'" effect="plain">{{ statusLabel(record.status) }}</el-tag><small>第 {{ record.currentRound }}/{{ record.maxRounds }} 轮</small></span>
              <span class="score">{{ record.overallScore === null ? '--' : `${record.overallScore} 分` }}</span>
              <time>{{ formatDate(record.updatedAt) }}</time>
              <span class="record-actions"><em>{{ record.status === 'completed' ? '查看报告' : '继续面试' }}</em><el-button text type="danger" :icon="Delete" :loading="deletingId === record.sessionId" @click.stop="removeInterview(record)">删除</el-button></span>
            </div>
          </div>
          <AppEmpty v-else-if="!loading && !loadError" title="还没有模拟面试记录" description="选择已保存的 JD 开始五轮模拟面试。" icon="◷">
            <el-button type="primary" plain @click="router.push('/app/interviews')">开始模拟面试</el-button>
          </AppEmpty>
        </template>
      </div>
      <div class="history-footer">当前显示 {{ activeTab === 'jd' ? filteredJDRecords.length : filteredInterviewRecords.length }} 条记录</div>
    </section>

    <el-dialog v-model="dialogVisible" title="JD 分析详情" width="min(760px, 92vw)">
      <div v-loading="detailLoading" class="detail-body">
        <template v-if="detail">
          <div class="detail-title"><div><h2>{{ detail.companyName }}</h2><p>{{ detail.jobTitle }} · {{ formatDate(detail.createdAt) }}</p></div><strong>{{ detail.matchScore }}<small>/100</small></strong></div>
          <p class="summary">{{ detail.jobSummary }}</p>
          <section><h3>核心要求</h3><ul><li v-for="item in detail.coreRequirements" :key="item">{{ item }}</li></ul></section>
          <section><h3>已匹配能力</h3><ul><li v-for="item in detail.matchedSkills" :key="item">{{ item }}</li></ul></section>
          <section><h3>待补能力</h3><ul><li v-for="item in detail.missingSkills" :key="item">{{ item }}</li></ul></section>
          <section><h3>简历建议</h3><ul><li v-for="item in detail.resumeSuggestions" :key="item">{{ item }}</li></ul></section>
          <section><h3>准备主题</h3><ul><li v-for="item in detail.preparationTopics" :key="item">{{ item }}</li></ul></section>
        </template>
      </div>
    </el-dialog>
  </div>
</template>

<style scoped>
.history-panel { overflow: hidden; }
.tabs { display: flex; height: 64px; padding: 0 22px; align-items: end; gap: 28px; border-bottom: 1px solid var(--border-color); }
.tabs button { height: 52px; border: 0; border-bottom: 2px solid transparent; color: var(--text-secondary); background: transparent; font-weight: 600; cursor: pointer; }
.tabs button.active { border-color: var(--color-primary); color: var(--color-primary); }
.tabs span { display: inline-grid; min-width: 20px; height: 20px; margin-left: 5px; place-items: center; border-radius: 99px; background: #f0f3f9; font-size: 11px; }
.filters { display: flex; padding: 18px 22px; gap: 10px; border-bottom: 1px solid var(--border-color); }
.filters .el-input { max-width: 420px; }
.records-area { min-height: 300px; }
.record-list { display: grid; }
.record-row { display: grid; padding: 18px 24px; align-items: center; grid-template-columns: minmax(180px,1fr) 90px 190px 190px; gap: 16px; border-bottom: 1px solid var(--border-color); background: #fff; cursor: pointer; }
.record-row.interview-row { grid-template-columns: minmax(180px,1fr) 120px 80px 170px 190px; }
.record-row:hover, .record-row:focus-visible { background: #f8faff; outline: none; }
.record-row > span:first-child, .interview-row > span:nth-child(2) { display: grid; gap: 4px; }
.record-row small, .record-row time { color: var(--text-secondary); }
.record-row .score { color: var(--color-primary); font-weight: 700; }
.record-actions { display: flex; align-items: center; justify-content: flex-end; gap: 6px; }
.record-actions em { color: var(--color-primary); font-style: normal; }
.history-footer { padding: 18px 22px; color: var(--text-secondary); }
.detail-body { min-height: 180px; }
.detail-title { display: flex; justify-content: space-between; gap: 20px; }
.detail-title h2 { margin: 0; }.detail-title p, .summary { color: var(--text-secondary); }
.detail-title > strong { color: var(--color-primary); font-size: 34px; }.detail-title small { font-size: 14px; }
.detail-body section { margin-top: 22px; }.detail-body h3 { font-size: 15px; }
.detail-body li { margin: 7px 0; color: var(--text-regular); line-height: 1.6; }
@media (max-width: 900px) { .record-row, .record-row.interview-row { grid-template-columns: 1fr auto; }.record-row time, .record-row .score { display: none; }.record-actions { grid-column: 1 / -1; justify-content: flex-start; } }
</style>
