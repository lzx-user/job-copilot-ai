<script setup lang="ts">
import { InfoFilled, MagicStick, Plus, Search } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'
import { reactive } from 'vue'
import AppEmpty from '../../components/common/AppEmpty.vue'

interface JdDraft {
  companyName: string
  jobTitle: string
  jdContent: string
  resumeSummary: string
}

// reactive 适合把同一个表单的多个字段组织在一起；当前阶段只保存在页面内，不写数据库。
const form = reactive<JdDraft>({
  companyName: '',
  jobTitle: '',
  jdContent: '',
  resumeSummary: '',
})

function explainAvailability() {
  ElMessage.info('JD AI 分析将在 7 月 19 日阶段接入，本页当前只提供输入骨架。')
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
      <el-tag effect="plain" round>静态骨架 · 后续阶段开放</el-tag>
    </div>

    <div class="analysis-layout">
      <section class="panel form-panel">
        <div class="panel-header"><h2 class="panel-title">岗位信息</h2><span class="draft-badge">草稿未保存</span></div>
        <el-form :model="form" label-position="top">
          <div class="two-column-form">
            <el-form-item label="公司名称">
              <el-input v-model="form.companyName" maxlength="100" placeholder="例如：某科技公司" />
            </el-form-item>
            <el-form-item label="岗位名称">
              <el-input v-model="form.jobTitle" maxlength="100" placeholder="例如：前端开发实习生" />
            </el-form-item>
          </div>
          <el-form-item label="岗位 JD">
            <el-input v-model="form.jdContent" type="textarea" :rows="8" maxlength="8000" show-word-limit placeholder="粘贴完整岗位职责与任职要求（后续分析要求 200～8000 字）" />
          </el-form-item>
          <el-form-item label="我的简历摘要">
            <el-input v-model="form.resumeSummary" type="textarea" :rows="6" maxlength="5000" show-word-limit placeholder="下一阶段完善个人档案后，可自动带入项目与能力摘要" />
          </el-form-item>

          <div class="skill-area">
            <div><strong>技术栈</strong><small>最多 30 个</small></div>
            <div class="skill-tags">
              <el-tag type="info" effect="plain">尚未添加技能</el-tag>
              <el-button plain :icon="Plus" @click="explainAvailability">添加技术栈</el-button>
            </div>
          </div>

          <el-button class="gradient-button analyze-button" :icon="MagicStick" @click="explainAvailability">开始分析</el-button>
          <p class="form-hint"><el-icon><InfoFilled /></el-icon> 当前不会调用 AI，也不会保存你输入的内容。</p>
        </el-form>
      </section>

      <section class="panel result-panel">
        <div class="result-header">
          <div><h2 class="panel-title">AI 分析结果</h2><p>完成分析后，将按参考图结构展示匹配建议。</p></div>
          <span class="result-chip"><el-icon><Search /></el-icon> 等待分析</span>
        </div>
        <AppEmpty title="还没有分析结果" description="先准备公司、岗位、JD、简历摘要与技术栈。AI 分析会在后续阶段安全地通过后端接入。" icon="⌕">
          <div class="result-preview">
            <span>匹配分数</span><span>核心要求</span><span>能力差距</span><span>优化建议</span>
          </div>
        </AppEmpty>
        <div class="disclaimer">AI 分析结果仅供求职准备参考，请结合自身情况判断和优化。</div>
      </section>
    </div>
  </div>
</template>

<style scoped>
.analysis-layout { display: grid; grid-template-columns: minmax(360px,.82fr) minmax(480px,1.18fr); gap: 18px; }
.form-panel,.result-panel { min-width: 0; padding: 24px; }
.draft-badge { padding: 5px 9px; border-radius: 8px; color: var(--text-secondary); background: #f4f6fa; font-size: 12px; }
.two-column-form { display: grid; grid-template-columns: repeat(2,minmax(0,1fr)); gap: 14px; }
.skill-area { display: grid; margin: 4px 0 22px; gap: 12px; }
.skill-area > div:first-child { display: flex; justify-content: space-between; }
.skill-area small { color: var(--text-secondary); }
.skill-tags { display: flex; flex-wrap: wrap; gap: 9px; }
.analyze-button { width: 100%; height: 46px; }
.form-hint { display: flex; margin: 12px 0 0; align-items: center; justify-content: center; gap: 5px; color: var(--text-secondary); font-size: 12px; }
.result-panel { display: flex; flex-direction: column; }
.result-header { display: flex; padding-bottom: 20px; justify-content: space-between; gap: 16px; border-bottom: 1px solid var(--border-color); }
.result-header p { margin: 5px 0 0; color: var(--text-secondary); }
.result-chip { display: inline-flex; height: 34px; padding: 0 11px; align-items: center; gap: 6px; border-radius: 9px; color: var(--color-primary); background: #edf3ff; white-space: nowrap; }
.result-panel :deep(.empty-state) { min-height: 470px; }
.result-preview { display: grid; width: min(100%,500px); margin-top: 10px; grid-template-columns: repeat(4,minmax(0,1fr)); gap: 8px; }
.result-preview span { padding: 9px; border: 1px solid var(--border-color); border-radius: 9px; color: var(--text-secondary); background: #fbfcff; font-size: 12px; }
.disclaimer { margin-top: auto; padding: 14px; border-radius: 10px; color: var(--text-secondary); background: #f8f9fd; font-size: 12px; text-align: center; }

@media (max-width: 1199px) { .analysis-layout { grid-template-columns: 1fr; } .result-panel :deep(.empty-state) { min-height: 300px; } }
@media (max-width: 600px) { .two-column-form { grid-template-columns: 1fr; gap: 0; } .form-panel,.result-panel { padding: 18px; } .result-preview { grid-template-columns: repeat(2,minmax(0,1fr)); } }
</style>

