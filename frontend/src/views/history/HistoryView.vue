<script setup lang="ts">
import { Refresh, Search } from '@element-plus/icons-vue'
import { computed, onMounted, ref } from 'vue'
import { getJDAnalysis, listJDAnalyses } from '../../api/jd-analysis.api'
import AppEmpty from '../../components/common/AppEmpty.vue'
import type { JdAnalysisRecord } from '../../types/jd-analysis'
import { toUserMessage } from '../../utils/error'

const records = ref<JdAnalysisRecord[]>([])
const keyword = ref('')
const loading = ref(true)
const detailLoading = ref(false)
const loadError = ref('')
const detail = ref<JdAnalysisRecord | null>(null)
const dialogVisible = ref(false)

const filteredRecords = computed(() => {
  const value = keyword.value.trim().toLowerCase()
  if (!value) return records.value
  return records.value.filter((record) =>
    `${record.companyName} ${record.jobTitle}`.toLowerCase().includes(value),
  )
})

async function loadRecords() {
  loading.value = true
  loadError.value = ''
  try {
    records.value = (await listJDAnalyses()).items
  } catch (error) {
    loadError.value = toUserMessage(error)
  } finally {
    loading.value = false
  }
}

async function openDetail(record: JdAnalysisRecord) {
  dialogVisible.value = true
  detailLoading.value = true
  detail.value = null
  try {
    detail.value = await getJDAnalysis(record.analysisId)
  } catch (error) {
    loadError.value = toUserMessage(error)
    dialogVisible.value = false
  } finally {
    detailLoading.value = false
  }
}

function formatDate(value: string) {
  return new Intl.DateTimeFormat('zh-CN', { dateStyle: 'medium', timeStyle: 'short' }).format(new Date(value))
}

onMounted(loadRecords)
</script>

<template>
  <div class="page-shell">
    <div class="page-heading">
      <div><span class="section-label">REVIEW</span><h1>历史记录</h1><p>查看当前账号保存的 JD 分析，重新登录后仍可恢复。</p></div>
      <el-tag effect="plain" round>{{ records.length }} 条 JD 分析</el-tag>
    </div>

    <section class="panel history-panel">
      <div class="tabs"><button class="active">JD 分析 <span>{{ records.length }}</span></button><button disabled>模拟面试 <span>0</span></button></div>
      <div class="filters">
        <el-input v-model="keyword" :prefix-icon="Search" placeholder="搜索公司或岗位" clearable />
        <el-button :icon="Refresh" :loading="loading" @click="loadRecords">刷新</el-button>
      </div>
      <el-alert v-if="loadError" :title="loadError" type="error" :closable="false" show-icon />
      <div v-loading="loading" class="records-area">
        <div v-if="filteredRecords.length" class="record-list">
          <button v-for="record in filteredRecords" :key="record.analysisId" class="record-row" @click="openDetail(record)">
            <span><strong>{{ record.companyName }}</strong><small>{{ record.jobTitle }}</small></span>
            <span class="score">{{ record.matchScore }} 分</span>
            <time>{{ formatDate(record.createdAt) }}</time>
            <em>查看详情 →</em>
          </button>
        </div>
        <AppEmpty v-else-if="!loading && !loadError" title="还没有 JD 分析记录" description="完成一份真实 JD 分析后，记录会保存并显示在这里。" icon="◷">
          <el-button type="primary" plain @click="$router.push('/app/jd-analysis')">准备第一份 JD</el-button>
        </AppEmpty>
      </div>
      <div class="history-footer">共 {{ filteredRecords.length }} 条记录</div>
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
.tabs button { height: 52px; border: 0; border-bottom: 2px solid transparent; color: var(--text-secondary); background: transparent; font-weight: 600; }
.tabs button.active { border-color: var(--color-primary); color: var(--color-primary); }
.tabs span { display: inline-grid; min-width: 20px; height: 20px; margin-left: 5px; place-items: center; border-radius: 99px; background: #f0f3f9; font-size: 11px; }
.filters { display: flex; padding: 18px 22px; gap: 10px; border-bottom: 1px solid var(--border-color); }
.filters .el-input { max-width: 420px; }
.records-area { min-height: 300px; }
.record-list { display: grid; }
.record-row { display: grid; padding: 18px 24px; align-items: center; grid-template-columns: minmax(180px,1fr) 90px 190px 100px; gap: 16px; border: 0; border-bottom: 1px solid var(--border-color); text-align: left; background: #fff; cursor: pointer; }
.record-row:hover { background: #f8faff; }
.record-row span:first-child { display: grid; gap: 4px; }
.record-row small, .record-row time { color: var(--text-secondary); }
.record-row .score { color: var(--color-primary); font-weight: 700; }
.record-row em { color: var(--color-primary); font-style: normal; }
.history-footer { padding: 18px 22px; color: var(--text-secondary); }
.detail-body { min-height: 180px; }
.detail-title { display: flex; justify-content: space-between; gap: 20px; }
.detail-title h2 { margin: 0; }.detail-title p, .summary { color: var(--text-secondary); }
.detail-title > strong { color: var(--color-primary); font-size: 34px; }.detail-title small { font-size: 14px; }
.detail-body section { margin-top: 22px; }.detail-body h3 { font-size: 15px; }
.detail-body li { margin: 7px 0; color: var(--text-regular); line-height: 1.6; }
@media (max-width: 700px) { .record-row { grid-template-columns: 1fr auto; }.record-row time, .record-row em { display: none; } }
</style>
