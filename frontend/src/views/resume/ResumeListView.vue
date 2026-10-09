<script setup lang="ts">
import { Plus } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'
import { computed, onMounted, reactive, ref } from 'vue'
import { createResume, createResumeVersion, getCareerWorkspace } from '../../api/career.api'
import AppEmpty from '../../components/common/AppEmpty.vue'
import type { CareerWorkspace } from '../../types/career'
import { toUserMessage } from '../../utils/error'

const emptyWorkspace: CareerWorkspace = { resumes: [], jobs: [], applications: [], events: [], interviews: [], retrospectives: [], offers: [] }
const workspace = ref<CareerWorkspace>(emptyWorkspace)
const loading = ref(true)
const saving = ref(false)
const dialogVisible = ref(false)
const editingResumeId = ref('')
const errorMessage = ref('')
const form = reactive({ title: '', content: '', skills: [] as string[], projectSummary: '', note: '', isDefault: false })
const resumes = computed(() => workspace.value.resumes)

async function load() {
  loading.value = true
  errorMessage.value = ''
  try { workspace.value = await getCareerWorkspace() } catch (error) { errorMessage.value = toUserMessage(error) } finally { loading.value = false }
}

async function submit() {
  if (!form.title.trim() || !form.content.trim()) return ElMessage.warning('请填写简历名称和文本内容')
  saving.value = true
  try {
		const payload = { content: form.content.trim(), skills: form.skills.map((item) => item.trim()).filter(Boolean), projectSummary: form.projectSummary, note: form.note }
		if (editingResumeId.value) await createResumeVersion(editingResumeId.value, payload)
		else await createResume({ ...payload, title: form.title.trim(), isDefault: form.isDefault })
    ElMessage.success('简历版本已保存')
    dialogVisible.value = false
    Object.assign(form, { title: '', content: '', skills: [], projectSummary: '', note: '', isDefault: false })
		editingResumeId.value = ''
    await load()
  } catch (error) { ElMessage.error(toUserMessage(error)) } finally { saving.value = false }
}

function openNewVersion(resumeId: string) {
	const resume = resumes.value.find((item) => item.id === resumeId)
	const latest = resume?.versions[0]
	if (!resume || !latest) return
	editingResumeId.value = resumeId
	Object.assign(form, { title: resume.title, content: latest.content, skills: [...latest.skills], projectSummary: latest.projectSummary, note: '', isDefault: resume.isDefault })
	dialogVisible.value = true
}
function openNewResume() {
	editingResumeId.value = ''
	Object.assign(form, { title: '', content: '', skills: [], projectSummary: '', note: '', isDefault: false })
	dialogVisible.value = true
}

function formatDate(value: string) { return new Intl.DateTimeFormat('zh-CN', { dateStyle: 'medium' }).format(new Date(value)) }
onMounted(load)
</script>

<template>
  <div class="page-shell">
    <div class="page-heading">
      <div><span class="section-label">RESUME VERSIONS</span><h1>简历资产</h1><p>保存投递时真正使用的简历版本，让匹配分析与结果可追溯。</p></div>
      <el-button type="primary" :icon="Plus" @click="openNewResume">新增简历</el-button>
    </div>
    <el-alert v-if="errorMessage" :title="errorMessage" type="error" :closable="false" show-icon />
    <section v-loading="loading" class="panel content-panel">
      <div v-if="resumes.length" class="resume-grid">
        <article v-for="resume in resumes" :key="resume.id" class="resume-card">
          <div><el-tag v-if="resume.isDefault" size="small">默认</el-tag><small>共 {{ resume.versions.length }} 个版本</small></div>
          <h3>{{ resume.title }}</h3>
          <p>{{ resume.versions[0]?.note || '暂无版本备注' }}</p>
          <span>最新保存于 {{ formatDate(resume.updatedAt) }}</span>
          <el-button text type="primary" @click="openNewVersion(resume.id)">基于当前内容新增版本</el-button>
          <el-collapse><el-collapse-item title="查看最新版本内容"><p class="resume-content">{{ resume.versions[0]?.content }}</p></el-collapse-item></el-collapse>
        </article>
      </div>
      <AppEmpty v-else-if="!loading && !errorMessage" title="还没有简历版本" description="先保存一份文本简历，后续投递、匹配和面试将引用同一个版本。" icon="▤">
        <el-button type="primary" @click="openNewResume">新增第一份简历</el-button>
      </AppEmpty>
    </section>

    <el-dialog v-model="dialogVisible" :title="editingResumeId ? '新增简历版本' : '新增简历'" width="min(680px, 92vw)" @closed="editingResumeId = ''">
      <el-form label-position="top">
        <el-form-item v-if="!editingResumeId" label="简历名称"><el-input v-model="form.title" maxlength="100" placeholder="例如：后端开发校招简历" /></el-form-item>
        <el-form-item label="简历文本"><el-input v-model="form.content" type="textarea" :rows="9" maxlength="20000" show-word-limit placeholder="粘贴脱敏后的简历正文" /></el-form-item>
        <el-form-item label="技能"><el-select v-model="form.skills" multiple filterable allow-create default-first-option placeholder="输入后回车添加" /></el-form-item>
        <el-form-item label="项目摘要"><el-input v-model="form.projectSummary" type="textarea" :rows="3" maxlength="5000" /></el-form-item>
        <el-form-item label="版本备注"><el-input v-model="form.note" maxlength="200" placeholder="例如：突出 Go 与工程化" /></el-form-item>
        <el-checkbox v-if="!editingResumeId" v-model="form.isDefault">设为默认简历</el-checkbox>
      </el-form>
      <template #footer><el-button @click="dialogVisible = false">取消</el-button><el-button type="primary" :loading="saving" @click="submit">保存</el-button></template>
    </el-dialog>
  </div>
</template>

<style scoped>
.content-panel { min-height: 420px; padding: 24px; }
.resume-grid { display: grid; grid-template-columns: repeat(auto-fill,minmax(290px,1fr)); gap: 16px; }
.resume-card { padding: 20px; border: 1px solid var(--border-color); border-radius: 16px; background: #fff; }
.resume-card > div:first-child { display: flex; justify-content: space-between; color: var(--text-secondary); }
.resume-card h3 { margin: 15px 0 7px; }
.resume-card > p,.resume-card > span { color: var(--text-secondary); font-size: 13px; }
.resume-content { white-space: pre-wrap; color: var(--text-regular); line-height: 1.7; }
</style>
