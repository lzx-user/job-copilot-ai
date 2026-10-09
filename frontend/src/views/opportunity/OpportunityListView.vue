<script setup lang="ts">
import { Plus, Search } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'
import { computed, onMounted, reactive, ref } from 'vue'
import { useRouter } from 'vue-router'
import { createJob, getCareerWorkspace } from '../../api/career.api'
import AppEmpty from '../../components/common/AppEmpty.vue'
import type { CareerWorkspace } from '../../types/career'
import { toUserMessage } from '../../utils/error'

const router = useRouter()
const workspace = ref<CareerWorkspace>({ resumes: [], jobs: [], applications: [], events: [], interviews: [], retrospectives: [], offers: [] })
const loading = ref(true)
const saving = ref(false)
const dialogVisible = ref(false)
const keyword = ref('')
const errorMessage = ref('')
const form = reactive({ companyName: '', jobTitle: '', location: '', sourceUrl: '', jdContent: '' })
const statusLabels: Record<string, string> = { planned: '待投递', applied: '已投递', screening: '筛选中', interview: '面试中', offer: '已获 Offer', rejected: '未通过', withdrawn: '已放弃', accepted: '已接受' }
const rows = computed(() => workspace.value.jobs.filter((job) => `${job.companyName}${job.jobTitle}`.toLowerCase().includes(keyword.value.trim().toLowerCase())).map((job) => ({ ...job, application: workspace.value.applications.find((item) => item.jobId === job.id) })))

async function load() { loading.value = true; errorMessage.value = ''; try { workspace.value = await getCareerWorkspace() } catch (error) { errorMessage.value = toUserMessage(error) } finally { loading.value = false } }
async function submit() {
  if (!form.companyName.trim() || !form.jobTitle.trim() || form.jdContent.trim().length < 200) return ElMessage.warning('请填写公司、岗位和至少 200 字的完整 JD')
  saving.value = true
  try {
    const result = await createJob({ ...form, companyName: form.companyName.trim(), jobTitle: form.jobTitle.trim(), jdContent: form.jdContent.trim() })
    ElMessage.success('岗位与 JD 版本已保存')
    dialogVisible.value = false
    await router.push(`/app/opportunities/${result.jobId}`)
  } catch (error) { ElMessage.error(toUserMessage(error)) } finally { saving.value = false }
}
function openJob(row: { id: string }) { router.push(`/app/opportunities/${row.id}`) }
onMounted(load)
</script>

<template>
  <div class="page-shell">
    <div class="page-heading"><div><span class="section-label">OPPORTUNITY PIPELINE</span><h1>岗位与投递</h1><p>以岗位为主线连接 JD、匹配、投递、面试、复盘和 Offer。</p></div><el-button type="primary" :icon="Plus" @click="dialogVisible = true">新增岗位</el-button></div>
    <el-alert v-if="errorMessage" :title="errorMessage" type="error" :closable="false" show-icon />
    <section class="panel toolbar"><el-input v-model="keyword" :prefix-icon="Search" clearable placeholder="搜索公司或岗位" /><span>{{ rows.length }} 个岗位</span></section>
    <section v-loading="loading" class="panel table-panel">
      <el-table v-if="rows.length" :data="rows" @row-click="openJob">
        <el-table-column label="公司 / 岗位" min-width="240"><template #default="{ row }"><strong>{{ row.companyName }}</strong><small>{{ row.jobTitle }}</small></template></el-table-column>
        <el-table-column prop="location" label="地点" min-width="120"><template #default="{ row }">{{ row.location || '未填写' }}</template></el-table-column>
        <el-table-column label="当前阶段" width="130"><template #default="{ row }"><el-tag :type="row.application?.status === 'offer' || row.application?.status === 'accepted' ? 'success' : 'info'">{{ row.application ? statusLabels[row.application.status] : '未建立投递' }}</el-tag></template></el-table-column>
        <el-table-column label="JD 版本" width="100"><template #default="{ row }">v{{ row.jdVersions[0]?.versionNumber || 1 }}</template></el-table-column>
        <el-table-column label="操作" width="100"><template #default><el-button text type="primary">查看详情</el-button></template></el-table-column>
      </el-table>
      <AppEmpty v-else-if="!loading && !errorMessage" title="还没有目标岗位" description="新增岗位后，从岗位详情完成整条求职链路。" icon="⌁"><el-button type="primary" @click="dialogVisible = true">新增第一个岗位</el-button></AppEmpty>
    </section>
    <el-dialog v-model="dialogVisible" title="新增目标岗位" width="min(720px, 92vw)"><el-form label-position="top"><div class="two-column"><el-form-item label="公司"><el-input v-model="form.companyName" maxlength="100" /></el-form-item><el-form-item label="岗位"><el-input v-model="form.jobTitle" maxlength="100" /></el-form-item></div><div class="two-column"><el-form-item label="地点"><el-input v-model="form.location" maxlength="100" /></el-form-item><el-form-item label="来源链接"><el-input v-model="form.sourceUrl" placeholder="https://" /></el-form-item></div><el-form-item label="岗位 JD"><el-input v-model="form.jdContent" type="textarea" :rows="10" maxlength="8000" show-word-limit placeholder="粘贴 200～8000 字完整 JD" /></el-form-item></el-form><template #footer><el-button @click="dialogVisible = false">取消</el-button><el-button type="primary" :loading="saving" @click="submit">保存并进入详情</el-button></template></el-dialog>
  </div>
</template>

<style scoped>
.toolbar { display:flex; padding:16px; align-items:center; justify-content:space-between; gap:16px; }.toolbar .el-input{max-width:360px}.toolbar span{color:var(--text-secondary)}.table-panel{min-height:420px;padding:16px}.table-panel :deep(.el-table__row){cursor:pointer}.table-panel strong,.table-panel small{display:block}.table-panel small{margin-top:4px;color:var(--text-secondary)}.two-column{display:grid;grid-template-columns:1fr 1fr;gap:16px}@media(max-width:700px){.two-column{grid-template-columns:1fr}}
</style>
