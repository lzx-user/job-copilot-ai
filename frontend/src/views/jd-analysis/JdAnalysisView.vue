<script setup lang="ts">
import { InfoFilled, MagicStick, Search } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'
import type { FormInstance, FormRules } from 'element-plus'
import { computed, onMounted, reactive, ref, watch } from 'vue'
import { analyzeJD } from '../../api/jd-analysis.api'
import AppEmpty from '../../components/common/AppEmpty.vue'
import { useProfileStore } from '../../stores/profile.store'
import type { JdAnalysisResult } from '../../types/jd-analysis'
import { toUserMessage } from '../../utils/error'

interface JdDraft {
  companyName: string
  jobTitle: string
  jdContent: string
  resumeSummary: string
  skills: string[]
}

const profileStore = useProfileStore()
const formRef = ref<FormInstance>()
const submitting = ref(false)
const analysisResult = ref<JdAnalysisResult | null>(null)
const analysisError = ref('')
let hasAppliedProfile = false

const form = reactive<JdDraft>({
  companyName: '',
  jobTitle: '',
  jdContent: '',
  resumeSummary: '',
  skills: [],
})

const rules: FormRules<JdDraft> = {
  companyName: [
    { required: true, message: '请填写公司名称', trigger: 'blur' },
    { max: 100, message: '公司名称不能超过 100 字', trigger: 'blur' },
  ],
  jobTitle: [
    { required: true, message: '请填写岗位名称', trigger: 'blur' },
    { max: 100, message: '岗位名称不能超过 100 字', trigger: 'blur' },
  ],
  jdContent: [
    { required: true, message: '请粘贴岗位 JD', trigger: 'blur' },
    { min: 200, message: '岗位 JD 至少需要 200 字', trigger: 'blur' },
    { max: 8000, message: '岗位 JD 不能超过 8000 字', trigger: 'blur' },
  ],
  resumeSummary: [
    { required: true, message: '请填写或从档案带入个人经历摘要', trigger: 'blur' },
    { max: 5000, message: '个人经历摘要不能超过 5000 字', trigger: 'blur' },
  ],
  skills: [
    {
      type: 'array',
      required: true,
      min: 1,
      message: '请至少添加一个技能',
      trigger: 'change',
    },
  ],
}

const profileStatus = computed(() => {
  if (profileStore.loading) return '正在读取个人档案…'
  if (profileStore.profileError) return '个人档案读取失败，可在本页手动填写'
  if (profileStore.loaded) return '已从个人档案带入经历摘要与技能，可继续修改'
  return '尚未读取个人档案'
})

function buildResumeSummary() {
  return [profileStore.profile.projectSummary, profileStore.profile.strengths]
    .map((item) => item.trim())
    .filter(Boolean)
    .join('\n\n')
}

// 档案异步返回后只自动带入一次，避免覆盖用户随后在本页做的修改。
function applyProfileToDraft() {
  if (!profileStore.loaded || hasAppliedProfile) return

  form.resumeSummary = buildResumeSummary()
  form.skills = [...profileStore.profile.skills]
  hasAppliedProfile = true
}

watch(
  () => profileStore.loaded,
  () => applyProfileToDraft(),
  { immediate: true },
)

onMounted(async () => {
  if (!profileStore.loaded) await profileStore.fetchProfile()
  applyProfileToDraft()
})

function normalizeSkills(values: string[]) {
  const uniqueSkills = new Map<string, string>()
  values.forEach((value) => {
    const skill = value.trim()
    if (skill) uniqueSkills.set(skill.toLocaleLowerCase(), skill)
  })
  return Array.from(uniqueSkills.values())
}

async function handleAnalyze() {
  if (submitting.value) return

  form.companyName = form.companyName.trim()
  form.jobTitle = form.jobTitle.trim()
  form.jdContent = form.jdContent.trim()
  form.resumeSummary = form.resumeSummary.trim()
  form.skills = normalizeSkills(form.skills)

  try {
    await formRef.value?.validate()
  } catch {
    ElMessage.warning('请先完善标红的岗位与个人信息')
    return
  }

  submitting.value = true
  analysisResult.value = null
  analysisError.value = ''

  try {
    analysisResult.value = await analyzeJD({
      companyName: form.companyName,
      jobTitle: form.jobTitle,
      jdContent: form.jdContent,
      resumeSummary: form.resumeSummary,
      skills: form.skills,
    })
    ElMessage.success('分析完成并已保存')
  } catch (error) {
    analysisError.value = toUserMessage(error)
    ElMessage.error(analysisError.value)
  } finally {
    submitting.value = false
  }
}
</script>

<template>
  <div class="page-shell">
    <div class="page-heading">
      <div>
        <span class="section-label">JOB MATCH</span>
        <h1>JD 智能匹配分析</h1>
        <p>整理岗位信息与个人经历，为后续 AI 匹配分析做好准备。</p>
      </div>
      <el-tag effect="plain" round>8.23 · 真实分析闭环</el-tag>
    </div>

    <el-alert
      class="stage-alert"
      title="分析由 Go 后端调用 AI；结构化结果校验通过后保存到数据库"
      type="success"
      :closable="false"
      show-icon
    />

    <div class="analysis-layout">
      <section class="panel form-panel">
        <div class="panel-header">
          <h2 class="panel-title">岗位信息</h2>
          <span class="draft-badge">草稿未保存</span>
        </div>

        <el-form
          ref="formRef"
          :model="form"
          :rules="rules"
          label-position="top"
          :disabled="submitting"
        >
          <div class="two-column-form">
            <el-form-item label="公司名称" prop="companyName">
              <el-input v-model="form.companyName" maxlength="100" placeholder="例如：某科技公司" />
            </el-form-item>
            <el-form-item label="岗位名称" prop="jobTitle">
              <el-input v-model="form.jobTitle" maxlength="100" placeholder="例如：全栈开发实习生" />
            </el-form-item>
          </div>

          <el-form-item label="岗位 JD" prop="jdContent">
            <el-input
              v-model="form.jdContent"
              type="textarea"
              :rows="8"
              maxlength="8000"
              show-word-limit
              placeholder="粘贴完整岗位职责与任职要求（200～8000 字）"
            />
          </el-form-item>

          <el-form-item label="我的经历摘要" prop="resumeSummary">
            <el-input
              v-model="form.resumeSummary"
              type="textarea"
              :rows="6"
              maxlength="5000"
              show-word-limit
              placeholder="将从个人档案自动带入，也可以针对当前岗位调整"
            />
            <p class="profile-hint" :class="{ 'is-error': profileStore.profileError }">
              {{ profileStatus }}
            </p>
          </el-form-item>

          <el-form-item label="技术栈" prop="skills">
            <el-select
              v-model="form.skills"
              class="skills-select"
              multiple
              filterable
              allow-create
              default-first-option
              :multiple-limit="30"
              placeholder="输入技能后按回车添加，例如 Vue 3、TypeScript、Go"
            />
          </el-form-item>

          <el-button
            class="gradient-button analyze-button"
            :icon="MagicStick"
            :loading="submitting"
            @click="handleAnalyze"
          >
            {{ submitting ? 'AI 正在分析…' : '开始匹配分析' }}
          </el-button>
          <p class="form-hint">
            <el-icon><InfoFilled /></el-icon>
            提交后将消耗一次 AI 调用，只有合法结果才会保存。
          </p>
        </el-form>
      </section>

      <section v-loading="submitting" class="panel result-panel">
        <div class="result-header">
          <div>
            <h2 class="panel-title">匹配分析结果</h2>
            <p>从岗位要求、能力匹配与准备方向三个维度辅助决策。</p>
          </div>
          <span class="result-chip">
            <el-icon><Search /></el-icon>
            {{ analysisResult ? '分析已保存' : '等待分析' }}
          </span>
        </div>

        <el-alert
          v-if="analysisError"
          class="result-error"
          :title="analysisError"
          type="error"
          :closable="false"
          show-icon
        />

        <div v-if="analysisResult" class="analysis-result">
          <div class="score-summary">
            <el-progress
              type="dashboard"
              :percentage="analysisResult.matchScore"
              :width="132"
              :stroke-width="10"
            />
            <div>
              <h3>{{ analysisResult.greetingMessage }}</h3>
              <p>{{ analysisResult.jobSummary }}</p>
            </div>
          </div>

          <div class="result-section">
            <h3>岗位核心要求</h3>
            <ul><li v-for="item in analysisResult.coreRequirements" :key="item">{{ item }}</li></ul>
          </div>

          <div class="skill-groups">
            <div class="result-section">
              <h3>已匹配能力</h3>
              <div v-if="analysisResult.matchedSkills.length" class="tag-list">
                <el-tag v-for="item in analysisResult.matchedSkills" :key="item" type="success" effect="light">{{ item }}</el-tag>
              </div>
              <p v-else class="empty-copy">暂未识别到明确匹配项</p>
            </div>
            <div class="result-section">
              <h3>待补充能力</h3>
              <div v-if="analysisResult.missingSkills.length" class="tag-list">
                <el-tag v-for="item in analysisResult.missingSkills" :key="item" type="warning" effect="light">{{ item }}</el-tag>
              </div>
              <p v-else class="empty-copy">暂无明显能力缺口</p>
            </div>
          </div>

          <div class="result-section">
            <h3>简历优化建议</h3>
            <ol><li v-for="item in analysisResult.resumeSuggestions" :key="item">{{ item }}</li></ol>
          </div>

          <div class="result-section">
            <h3>面试准备主题</h3>
            <div class="tag-list">
              <el-tag v-for="item in analysisResult.preparationTopics" :key="item" effect="plain">{{ item }}</el-tag>
            </div>
          </div>
        </div>

        <AppEmpty
          v-else
          :title="analysisError ? '本次分析未完成' : '还没有分析结果'"
          :description="analysisError ? '请根据上方提示检查服务状态后重试。' : '填写公司、岗位、JD、个人经历摘要与技术栈，然后开始真实分析。'"
          icon="⌕"
        />

        <div class="disclaimer">AI 分析结果仅供求职准备参考，请结合自身情况判断和优化。</div>
      </section>
    </div>
  </div>
</template>

<style scoped>
.stage-alert { margin-bottom: 18px; }
.analysis-layout { display: grid; grid-template-columns: minmax(360px,.82fr) minmax(480px,1.18fr); gap: 18px; }
.form-panel,.result-panel { min-width: 0; padding: 24px; }
.draft-badge { padding: 5px 9px; border-radius: 8px; color: var(--text-secondary); background: #f4f6fa; font-size: 12px; }
.two-column-form { display: grid; grid-template-columns: repeat(2,minmax(0,1fr)); gap: 14px; }
.skills-select { width: 100%; }
.profile-hint { margin: 6px 0 0; color: var(--text-secondary); font-size: 12px; }
.profile-hint.is-error { color: var(--color-danger, #f56c6c); }
.analyze-button { width: 100%; height: 46px; }
.form-hint { display: flex; margin: 12px 0 0; align-items: center; justify-content: center; gap: 5px; color: var(--text-secondary); font-size: 12px; }
.result-panel { display: flex; flex-direction: column; }
.result-header { display: flex; padding-bottom: 20px; justify-content: space-between; gap: 16px; border-bottom: 1px solid var(--border-color); }
.result-header p { margin: 5px 0 0; color: var(--text-secondary); }
.result-chip { display: inline-flex; height: 34px; padding: 0 11px; align-items: center; gap: 6px; border-radius: 9px; color: var(--color-primary); background: #edf3ff; white-space: nowrap; }
.result-panel :deep(.empty-state) { min-height: 470px; }
.result-error { margin-top: 18px; }
.analysis-result { display: grid; padding: 24px 0; gap: 16px; }
.score-summary { display: flex; padding: 20px; align-items: center; gap: 24px; border-radius: 14px; background: linear-gradient(135deg,#f2f7ff,#f8f5ff); }
.score-summary h3 { margin: 0 0 8px; font-size: 20px; }
.score-summary p { margin: 0; color: var(--text-secondary); line-height: 1.7; }
.result-section { padding: 18px; border: 1px solid var(--border-color); border-radius: 12px; background: #fff; }
.result-section h3 { margin: 0 0 12px; font-size: 15px; }
.result-section ul,.result-section ol { display: grid; margin: 0; padding-left: 20px; gap: 8px; color: var(--text-secondary); line-height: 1.6; }
.skill-groups { display: grid; grid-template-columns: repeat(2,minmax(0,1fr)); gap: 12px; }
.tag-list { display: flex; flex-wrap: wrap; gap: 8px; }
.empty-copy { margin: 0; color: var(--text-secondary); font-size: 13px; }
.disclaimer { margin-top: auto; padding: 14px; border-radius: 10px; color: var(--text-secondary); background: #f8f9fd; font-size: 12px; text-align: center; }

@media (max-width: 1199px) { .analysis-layout { grid-template-columns: 1fr; } .result-panel :deep(.empty-state) { min-height: 300px; } }
@media (max-width: 600px) { .two-column-form,.skill-groups { grid-template-columns: 1fr; gap: 0; } .skill-groups { gap: 12px; } .form-panel,.result-panel { padding: 18px; } .score-summary { flex-direction: column; text-align: center; } }
</style>
