<script setup lang="ts">
import { Briefcase, ChatDotRound, Microphone, Reading } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'
import type { Component } from 'vue'
import { onMounted, ref } from 'vue'
import { useRouter } from 'vue-router'
import { listInterviewOptions, startInterview } from '../../api/interview.api'
import AppEmpty from '../../components/common/AppEmpty.vue'
import type { InterviewOption } from '../../types/interview'
import { toUserMessage } from '../../utils/error'

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

const router = useRouter()

const options = ref<InterviewOption[]>([])
const selectedAnalysisId = ref('')
const loadingOptions = ref(true)
const starting = ref(false)
const loadError = ref('')

async function loadOptions() {
  loadingOptions.value = true
  loadError.value = ''
  try {
    const result = await listInterviewOptions()
    options.value = result.items
    selectedAnalysisId.value = result.items[0]?.analysisId ?? ''
  } catch (error) {
    loadError.value = toUserMessage(error)
  } finally {
    loadingOptions.value = false
  }
}

async function handleStart() {
  if (!selectedAnalysisId.value || starting.value) return
  starting.value = true
  try {
    const result = await startInterview(selectedAnalysisId.value)
    ElMessage.success('模拟面试已开始，第一题已保存')
    await router.push({ name: 'interview-session', params: { id: result.sessionId } })
  } catch (error) {
    ElMessage.error(toUserMessage(error))
  } finally {
    starting.value = false
  }
}

onMounted(loadOptions)
</script>

<template>
  <div class="page-shell">
    <div class="page-heading">
      <div><span class="section-label">MOCK INTERVIEW</span><h1>模拟面试</h1><p>选择已保存的 JD，由 AI 生成第一道针对性面试题。</p></div>
      <el-tag effect="plain" round>固定 5 轮 · 单轮问答已接入</el-tag>
    </div>

    <section class="panel intro-panel">
      <div class="intro-copy">
        <span class="intro-badge"><el-icon><Reading /></el-icon> 面试准备流程</span>
        <h2>选择合适的练习方式，<br /><em>每次只聚焦一个目标。</em></h2>
        <p>创建会话后进入对话页，可提交回答并获得评分、反馈与下一题。</p>
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
      <div class="panel-header">
        <div><h2 class="panel-title">选择目标岗位</h2><p>仅展示当前账号最近保存的 JD 分析。</p></div>
        <el-tag type="info">{{ options.length }} 条记录</el-tag>
      </div>

      <el-alert v-if="loadError" :title="loadError" type="error" :closable="false" show-icon />
      <div v-loading="loadingOptions">
        <el-radio-group v-if="options.length" v-model="selectedAnalysisId" class="jd-option-list">
          <el-radio v-for="option in options" :key="option.analysisId" :value="option.analysisId" border>
            <span class="jd-option-title">{{ option.companyName }} · {{ option.jobTitle }}</span>
            <span class="jd-option-score">匹配度 {{ option.matchScore }}</span>
          </el-radio>
        </el-radio-group>
        <AppEmpty v-else-if="!loadingOptions && !loadError" title="还没有可用于面试的 JD" description="先完成一份真实 JD 分析，再从这里开始面试。" icon="▤" />
      </div>

      <el-button
        v-if="options.length"
        class="gradient-button start-button"
        :loading="starting"
        :disabled="!selectedAnalysisId"
        @click="handleStart"
      >
        {{ starting ? 'AI 正在生成第一题…' : '开始模拟面试' }}
      </el-button>

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
.jd-option-list { display: grid; margin-top: 20px; gap: 12px; }
.jd-option-list :deep(.el-radio) { width: 100%; height: auto; min-height: 58px; margin: 0; padding: 14px 16px; }
.jd-option-list :deep(.el-radio__label) { display: flex; width: 100%; align-items: center; justify-content: space-between; gap: 16px; }
.jd-option-title { color: var(--text-primary); font-weight: 600; }
.jd-option-score { color: var(--color-primary); font-size: 13px; }
.start-button { width: 100%; height: 44px; margin-top: 18px; }
@media (max-width: 1000px) { .mode-grid { grid-template-columns: 1fr; } }
@media (max-width: 767px) { .intro-panel { min-height: 250px; padding: 24px; } .interview-visual { right: -70px; opacity: .35; } .mode-card { grid-template-columns: 50px 1fr; } .mode-card .el-radio { grid-column: 2; } .jd-selector { padding: 18px; } }
</style>
