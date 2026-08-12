<script setup lang="ts">
import { Calendar, Refresh, Search } from '@element-plus/icons-vue'
import { ref } from 'vue'
import AppEmpty from '../../components/common/AppEmpty.vue'

const activeTab = ref('all')
const typeFilter = ref('all')
const statusFilter = ref('all')
const keyword = ref('')
</script>

<template>
  <div class="page-shell">
    <div class="page-heading">
      <div><span class="section-label">REVIEW</span><h1>历史记录</h1><p>集中查看 JD 分析与模拟面试，形成长期复盘习惯。</p></div>
      <el-tag effect="plain" round>真实数据接入前保持为空</el-tag>
    </div>

    <section class="panel history-panel">
      <div class="tabs">
        <button v-for="tab in [{key:'all',label:'全部'},{key:'jd',label:'JD 分析'},{key:'interview',label:'模拟面试'}]" :key="tab.key" :class="{active:activeTab===tab.key}" @click="activeTab=tab.key">{{ tab.label }} <span>0</span></button>
      </div>

      <div class="filters">
        <el-input v-model="keyword" :prefix-icon="Search" placeholder="搜索公司或岗位" clearable />
        <el-select v-model="typeFilter" aria-label="记录类型"><el-option label="全部类型" value="all" /><el-option label="JD 分析" value="jd" /><el-option label="模拟面试" value="interview" /></el-select>
        <el-select v-model="statusFilter" aria-label="记录状态"><el-option label="全部状态" value="all" /><el-option label="已完成" value="completed" /><el-option label="进行中" value="in-progress" /></el-select>
        <el-button :icon="Calendar">日期范围</el-button>
        <el-button :icon="Refresh" circle aria-label="刷新" />
      </div>

      <div class="table-head" aria-hidden="true"><span>公司</span><span>岗位</span><span>类型</span><span>时间</span><span>状态</span><span>操作</span></div>
      <AppEmpty title="还没有历史记录" description="本阶段不连接业务数据库，也不会用虚假企业数据填充列表。完成后续真实分析与面试后，记录会显示在这里。" icon="◷">
        <el-button type="primary" plain @click="$router.push('/app/jd-analysis')">准备第一份 JD</el-button>
      </AppEmpty>
      <div class="history-footer"><span>共 0 条记录</span><el-pagination small layout="prev, pager, next" :total="0" /></div>
    </section>
  </div>
</template>

<style scoped>
.history-panel { overflow: hidden; }
.tabs { display: flex; height: 64px; padding: 0 22px; align-items: end; gap: 28px; border-bottom: 1px solid var(--border-color); }
.tabs button { height: 52px; padding: 0; cursor: pointer; border-bottom: 2px solid transparent; color: var(--text-secondary); background: transparent; font-weight: 600; }
.tabs button.active { border-color: var(--color-primary); color: var(--color-primary); }.tabs span { display: inline-grid; min-width: 20px; height: 20px; margin-left: 5px; padding: 0 5px; place-items: center; border-radius: 99px; background: #f0f3f9; font-size: 11px; }
.filters { display: grid; padding: 20px 22px; grid-template-columns: minmax(180px,1.3fr) minmax(130px,.7fr) minmax(130px,.7fr) auto auto; gap: 10px; border-bottom: 1px solid var(--border-color); }
.table-head { display: grid; padding: 15px 24px; grid-template-columns: 1.1fr 1.2fr .8fr .8fr .7fr .5fr; gap: 12px; color: var(--text-secondary); background: #fafbfe; font-size: 12px; font-weight: 600; }
.history-panel :deep(.empty-state) { min-height: 360px; }
.history-footer { display: flex; min-height: 62px; padding: 0 22px; align-items: center; justify-content: space-between; border-top: 1px solid var(--border-color); color: var(--text-secondary); }
@media (max-width: 900px) { .filters { grid-template-columns: repeat(2,minmax(0,1fr)); }.filters .el-input { grid-column: 1/-1; }.table-head { display: none; } }
@media (max-width: 520px) { .tabs { padding-inline: 16px; gap: 16px; }.filters { padding: 16px; grid-template-columns: 1fr; }.filters .el-input { grid-column: auto; }.history-footer { padding-inline: 16px; } }
</style>

