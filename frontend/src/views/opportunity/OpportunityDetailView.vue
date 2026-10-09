<script setup lang="ts">
import { Back, MagicStick, Microphone, Plus } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'
import { computed, onMounted, reactive, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { addApplicationEvent, createApplication, createJDVersion, createRealInterview, getCareerWorkspace, saveOffer, saveRetrospective, updateRealInterview } from '../../api/career.api'
import type { CareerWorkspace } from '../../types/career'
import { toUserMessage } from '../../utils/error'

const route = useRoute()
const router = useRouter()
const workspace = ref<CareerWorkspace>({ resumes: [], jobs: [], applications: [], events: [], interviews: [], retrospectives: [], offers: [] })
const loading = ref(true)
const saving = ref(false)
const errorMessage = ref('')
const activeTab = ref('overview')
const jdDialogVisible = ref(false)
const newJDContent = ref('')
const job = computed(() => workspace.value.jobs.find((item) => item.id === route.params.id))
const application = computed(() => workspace.value.applications.find((item) => item.jobId === job.value?.id))
const events = computed(() => workspace.value.events.filter((item) => item.applicationId === application.value?.id).sort((a, b) => +new Date(b.occurredAt) - +new Date(a.occurredAt)))
const interviews = computed(() => workspace.value.interviews.filter((item) => item.applicationId === application.value?.id))
const offer = computed(() => workspace.value.offers.find((item) => item.applicationId === application.value?.id))
const latestJD = computed(() => job.value?.jdVersions[0])
const allVersions = computed(() => workspace.value.resumes.flatMap((resume) => resume.versions.map((version) => ({ ...version, resumeTitle: resume.title }))))
const statusLabels: Record<string, string> = { planned: '待投递', applied: '已投递', screening: '筛选中', interview: '面试中', offer: '已获 Offer', rejected: '未通过', withdrawn: '已放弃', accepted: '已接受' }
const eventLabels: Record<string, string> = { planned: '加入计划', applied: '已投递', screening: '进入筛选', written_test: '笔试', interview_scheduled: '安排面试', interview_completed: '完成面试', offer_received: '收到 Offer', rejected: '未通过', withdrawn: '主动放弃', offer_accepted: '接受 Offer', note: '备注' }

const applicationForm = reactive({ resumeVersionId: '', status: 'planned', source: '', appliedAt: '', deadline: '' })
const eventForm = reactive({ eventType: 'applied', occurredAt: '', outcome: '', notes: '' })
const interviewForm = reactive({ roundName: '', scheduledAt: '', durationMinutes: 60 as number | null, format: '线上', locationOrLink: '', result: 'scheduled', notes: '' })
const retrospectiveForm = reactive({ interviewId: '', questionsText: '', selfAssessment: '', strengthsText: '', weaknessesText: '', actionsText: '' })
const offerForm = reactive({ status: 'pending', receivedAt: '', deadline: '', salarySummary: '', notes: '', decidedAt: '' })
const interviewResults = reactive<Record<string, string>>({})

function optionalDate(value: string) { return value ? new Date(value).toISOString() : null }
function lines(value: string) { return value.split('\n').map((item) => item.trim()).filter(Boolean) }
async function load() {
  loading.value = true; errorMessage.value = ''
  try {
    workspace.value = await getCareerWorkspace()
    applicationForm.resumeVersionId ||= allVersions.value[0]?.id ?? ''
    retrospectiveForm.interviewId ||= interviews.value[0]?.id ?? ''
    interviews.value.forEach((item) => { interviewResults[item.id] = item.result })
    if (offer.value) Object.assign(offerForm, { status: offer.value.status, receivedAt: offer.value.receivedAt.slice(0, 16), deadline: offer.value.deadline?.slice(0, 16) ?? '', salarySummary: offer.value.salarySummary, notes: offer.value.notes, decidedAt: offer.value.decidedAt?.slice(0, 16) ?? '' })
  } catch (error) { errorMessage.value = toUserMessage(error) } finally { loading.value = false }
}
async function run(action: () => Promise<unknown>, message: string) { saving.value = true; try { await action(); ElMessage.success(message); await load() } catch (error) { ElMessage.error(toUserMessage(error)) } finally { saving.value = false } }
async function submitApplication() { if (!job.value) return; await run(() => createApplication({ jobId: job.value!.id, resumeVersionId: applicationForm.resumeVersionId || null, status: applicationForm.status, source: applicationForm.source, appliedAt: optionalDate(applicationForm.appliedAt), deadline: optionalDate(applicationForm.deadline) }), '投递记录已建立') }
async function submitEvent() { if (!application.value) return; await run(() => addApplicationEvent({ applicationId: application.value!.id, eventType: eventForm.eventType, occurredAt: optionalDate(eventForm.occurredAt), outcome: eventForm.outcome, notes: eventForm.notes }), '进展已记录') }
async function submitInterview() { if (!application.value || !interviewForm.roundName.trim()) return ElMessage.warning('请填写面试轮次'); await run(() => createRealInterview({ applicationId: application.value!.id, ...interviewForm, scheduledAt: optionalDate(interviewForm.scheduledAt) }), '真实面试已记录') }
async function updateInterview(item: { id: string; notes: string }) { await run(() => updateRealInterview(item.id, { result: interviewResults[item.id] || 'pending', notes: item.notes }), '面试结果已更新') }
async function submitRetrospective() { if (!retrospectiveForm.interviewId) return; await run(() => saveRetrospective({ interviewId: retrospectiveForm.interviewId, questions: lines(retrospectiveForm.questionsText), selfAssessment: retrospectiveForm.selfAssessment, strengths: lines(retrospectiveForm.strengthsText), weaknesses: lines(retrospectiveForm.weaknessesText), followUpActions: lines(retrospectiveForm.actionsText), aiAnalysis: '' }), '面试复盘已保存') }
async function submitOffer() { if (!application.value) return; await run(() => saveOffer({ applicationId: application.value!.id, status: offerForm.status, receivedAt: optionalDate(offerForm.receivedAt), deadline: optionalDate(offerForm.deadline), salarySummary: offerForm.salarySummary, notes: offerForm.notes, decidedAt: optionalDate(offerForm.decidedAt) }), 'Offer 信息已保存') }
async function submitJDVersion() { if (!job.value || newJDContent.value.trim().length < 200) return ElMessage.warning('新版 JD 至少需要 200 字'); await run(() => createJDVersion(job.value!.id, newJDContent.value.trim()), '新版 JD 已保存'); jdDialogVisible.value = false; newJDContent.value = '' }
function analyze() { if (!job.value || !latestJD.value) return; router.push({ path: '/app/jd-analysis', query: { jobId: job.value.id, jdVersionId: latestJD.value.id, resumeVersionId: application.value?.resumeVersionId || allVersions.value[0]?.id || '' } }) }
function mockInterview() { router.push({ path: '/app/interviews', query: { applicationId: application.value?.id || '', resumeVersionId: application.value?.resumeVersionId || '' } }) }
function formatDate(value: string) { return new Intl.DateTimeFormat('zh-CN', { dateStyle: 'medium', timeStyle: 'short' }).format(new Date(value)) }
onMounted(load)
</script>

<template>
  <div v-loading="loading" class="page-shell">
    <el-button text :icon="Back" @click="router.push('/app/opportunities')">返回岗位列表</el-button>
    <el-alert v-if="errorMessage" :title="errorMessage" type="error" :closable="false" show-icon />
    <template v-if="job">
      <section class="panel job-hero"><div><span class="section-label">{{ job.companyName }}</span><h1>{{ job.jobTitle }}</h1><p>{{ job.location || '地点未填写' }} · JD v{{ latestJD?.versionNumber }}</p></div><div><el-tag size="large" :type="application?.status === 'offer' || application?.status === 'accepted' ? 'success' : 'info'">{{ application ? statusLabels[application.status] : '未建立投递' }}</el-tag><el-button :icon="MagicStick" @click="analyze">AI 匹配</el-button><el-button v-if="application" :icon="Microphone" @click="mockInterview">针对性模拟</el-button></div></section>
      <el-tabs v-model="activeTab" class="panel detail-tabs">
        <el-tab-pane label="岗位概览" name="overview"><div class="tab-content"><div class="section-row"><h3>岗位 JD · v{{ latestJD?.versionNumber }}</h3><el-button text type="primary" @click="newJDContent = latestJD?.content || ''; jdDialogVisible = true">保存新版 JD</el-button></div><p class="long-text">{{ latestJD?.content }}</p><a v-if="job.sourceUrl" :href="job.sourceUrl" target="_blank" rel="noreferrer">打开岗位来源</a></div></el-tab-pane>
        <el-tab-pane label="投递时间线" name="application"><div class="tab-content">
          <template v-if="!application"><h3>建立投递记录</h3><p class="muted">岗位与投递分开保存：只有真正准备或完成投递时才创建 Application。</p><el-form label-position="top"><el-form-item label="使用的简历版本"><el-select v-model="applicationForm.resumeVersionId" clearable><el-option v-for="version in allVersions" :key="version.id" :label="`${version.resumeTitle} · v${version.versionNumber}`" :value="version.id" /></el-select></el-form-item><div class="form-grid"><el-form-item label="初始阶段"><el-select v-model="applicationForm.status"><el-option label="待投递" value="planned" /><el-option label="已投递" value="applied" /></el-select></el-form-item><el-form-item label="来源"><el-input v-model="applicationForm.source" placeholder="官网 / 内推 / 招聘平台" /></el-form-item><el-form-item label="投递时间"><el-date-picker v-model="applicationForm.appliedAt" type="datetime" value-format="YYYY-MM-DDTHH:mm" /></el-form-item><el-form-item label="截止时间"><el-date-picker v-model="applicationForm.deadline" type="datetime" value-format="YYYY-MM-DDTHH:mm" /></el-form-item></div><el-button type="primary" :loading="saving" @click="submitApplication">建立投递</el-button></el-form></template>
          <template v-else><div class="section-row"><h3>进展记录</h3><el-tag>{{ statusLabels[application.status] }}</el-tag></div><el-timeline v-if="events.length"><el-timeline-item v-for="event in events" :key="event.id" :timestamp="formatDate(event.occurredAt)" placement="top"><strong>{{ eventLabels[event.eventType] || event.eventType }}</strong><p v-if="event.outcome || event.notes">{{ event.outcome }} {{ event.notes }}</p></el-timeline-item></el-timeline><el-divider /><h3>新增进展</h3><div class="form-grid"><el-select v-model="eventForm.eventType"><el-option v-for="(label, value) in eventLabels" :key="value" :label="label" :value="value" /></el-select><el-date-picker v-model="eventForm.occurredAt" type="datetime" value-format="YYYY-MM-DDTHH:mm" placeholder="发生时间（默认现在）" /><el-input v-model="eventForm.outcome" placeholder="结果摘要" /><el-input v-model="eventForm.notes" placeholder="备注" /></div><el-button class="action-button" :icon="Plus" :loading="saving" @click="submitEvent">保存进展</el-button></template>
        </div></el-tab-pane>
        <el-tab-pane label="真实面试与复盘" name="interviews"><div class="tab-content"><template v-if="application"><h3>真实面试轮次</h3><div v-if="interviews.length" class="interview-list"><article v-for="item in interviews" :key="item.id"><strong>{{ item.roundName }}</strong><span>{{ item.scheduledAt ? formatDate(item.scheduledAt) : '时间未定' }}</span><p>{{ item.notes }}</p><div class="result-editor"><el-select v-model="interviewResults[item.id]" size="small"><el-option label="已安排" value="scheduled" /><el-option label="已完成" value="completed" /><el-option label="通过" value="passed" /><el-option label="未通过" value="failed" /><el-option label="已取消" value="cancelled" /></el-select><el-button size="small" :loading="saving" @click="updateInterview(item)">更新结果</el-button></div></article></div><div class="form-grid"><el-input v-model="interviewForm.roundName" placeholder="轮次，如：一面 / HR 面" /><el-date-picker v-model="interviewForm.scheduledAt" type="datetime" value-format="YYYY-MM-DDTHH:mm" placeholder="面试时间" /><el-input v-model="interviewForm.format" placeholder="线上 / 线下" /><el-input v-model="interviewForm.locationOrLink" placeholder="地点或会议链接" /></div><el-input v-model="interviewForm.notes" class="block-field" type="textarea" :rows="2" placeholder="面试前备注" /><el-button :loading="saving" @click="submitInterview">记录面试轮次</el-button><el-divider /><h3>面试复盘</h3><el-select v-model="retrospectiveForm.interviewId" placeholder="选择要复盘的轮次"><el-option v-for="item in interviews" :key="item.id" :label="item.roundName" :value="item.id" /></el-select><div class="retro-grid"><el-input v-model="retrospectiveForm.questionsText" type="textarea" :rows="4" placeholder="面试问题，每行一条" /><el-input v-model="retrospectiveForm.selfAssessment" type="textarea" :rows="4" placeholder="自我复盘" /><el-input v-model="retrospectiveForm.strengthsText" type="textarea" :rows="3" placeholder="表现好的地方，每行一条" /><el-input v-model="retrospectiveForm.weaknessesText" type="textarea" :rows="3" placeholder="待改进，每行一条" /><el-input v-model="retrospectiveForm.actionsText" type="textarea" :rows="3" placeholder="后续行动，每行一条" /></div><el-button type="primary" :loading="saving" :disabled="!retrospectiveForm.interviewId" @click="submitRetrospective">保存复盘</el-button></template><el-empty v-else description="请先在投递时间线建立投递记录" /></div></el-tab-pane>
        <el-tab-pane label="Offer" name="offer"><div class="tab-content"><template v-if="application"><h3>{{ offer ? '更新 Offer' : '记录 Offer' }}</h3><div class="form-grid"><el-select v-model="offerForm.status"><el-option label="待决定" value="pending" /><el-option label="已接受" value="accepted" /><el-option label="已拒绝" value="declined" /><el-option label="已过期" value="expired" /></el-select><el-date-picker v-model="offerForm.receivedAt" type="datetime" value-format="YYYY-MM-DDTHH:mm" placeholder="收到时间" /><el-date-picker v-model="offerForm.deadline" type="datetime" value-format="YYYY-MM-DDTHH:mm" placeholder="决策截止" /><el-input v-model="offerForm.salarySummary" placeholder="薪资摘要（避免敏感明细）" /></div><el-input v-model="offerForm.notes" class="block-field" type="textarea" :rows="4" placeholder="Offer 备注与比较因素" /><el-button type="primary" :loading="saving" @click="submitOffer">保存 Offer</el-button></template><el-empty v-else description="请先建立投递记录" /></div></el-tab-pane>
      </el-tabs>
    </template>
    <el-empty v-else-if="!loading && !errorMessage" description="岗位不存在或已删除"><el-button @click="router.push('/app/opportunities')">返回列表</el-button></el-empty>
    <el-dialog v-model="jdDialogVisible" title="保存新版 JD" width="min(720px, 92vw)"><el-input v-model="newJDContent" type="textarea" :rows="12" maxlength="8000" show-word-limit /><template #footer><el-button @click="jdDialogVisible = false">取消</el-button><el-button type="primary" :loading="saving" @click="submitJDVersion">保存新版本</el-button></template></el-dialog>
  </div>
</template>

<style scoped>
.job-hero{display:flex;margin-top:12px;padding:26px;align-items:center;justify-content:space-between;gap:20px}.job-hero h1{margin:5px 0}.job-hero p{margin:0;color:var(--text-secondary)}.job-hero>div:last-child{display:flex;align-items:center;gap:10px;flex-wrap:wrap}.detail-tabs{margin-top:18px;padding:0 24px 24px}.tab-content{min-height:420px;padding:18px 4px}.long-text{white-space:pre-wrap;line-height:1.8;color:var(--text-regular)}.muted{color:var(--text-secondary)}.form-grid{display:grid;margin:14px 0;grid-template-columns:repeat(2,minmax(0,1fr));gap:14px}.form-grid>*{width:100%}.section-row{display:flex;align-items:center;justify-content:space-between}.action-button,.block-field{margin-bottom:14px}.interview-list{display:grid;margin-bottom:18px;gap:10px}.interview-list article{padding:14px;border:1px solid var(--border-color);border-radius:12px}.interview-list span{margin-left:12px;color:var(--text-secondary)}.interview-list p{margin-bottom:0}.result-editor{display:flex;margin-top:10px;gap:8px}.result-editor .el-select{max-width:140px}.retro-grid{display:grid;margin:14px 0;grid-template-columns:1fr 1fr;gap:12px}@media(max-width:760px){.job-hero{align-items:flex-start;flex-direction:column}.form-grid,.retro-grid{grid-template-columns:1fr}.detail-tabs{padding-inline:14px}}
</style>
