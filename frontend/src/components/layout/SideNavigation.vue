<script setup lang="ts">
import {
  Clock,
  DocumentChecked,
  House,
  Microphone,
  Setting,
  User,
} from '@element-plus/icons-vue'
import type { Component } from 'vue'
import { computed } from 'vue'
import { useRoute } from 'vue-router'
import BrandLogo from '../common/BrandLogo.vue'

interface Props {
  collapsed?: boolean
  drawer?: boolean
}

const props = withDefaults(defineProps<Props>(), {
  collapsed: false,
  drawer: false,
})

const emit = defineEmits<{
  navigate: []
}>()

interface NavigationItem {
  label: string
  path: string
  icon: Component
}

const route = useRoute()
const items: NavigationItem[] = [
  { label: '仪表盘', path: '/app/dashboard', icon: House },
  { label: 'JD 分析', path: '/app/jd-analysis', icon: DocumentChecked },
  { label: '模拟面试', path: '/app/interviews', icon: Microphone },
  { label: '历史记录', path: '/app/history', icon: Clock },
  { label: '个人档案', path: '/app/profile', icon: User },
  { label: '设置', path: '/app/settings', icon: Setting },
]

const shouldCollapse = computed(() => props.collapsed && !props.drawer)

function isActive(path: string) {
  return path === '/app/interviews'
    ? route.path.startsWith('/app/interviews')
    : route.path === path
}
</script>

<template>
  <div class="side-navigation" :class="{ collapsed: shouldCollapse }">
    <div class="logo-area">
      <BrandLogo :compact="shouldCollapse" />
    </div>

    <nav class="nav-list" aria-label="主要导航">
      <RouterLink
        v-for="item in items"
        :key="item.path"
        :to="item.path"
        class="nav-item"
        :class="{ active: isActive(item.path) }"
        :title="shouldCollapse ? item.label : undefined"
        @click="emit('navigate')"
      >
        <el-icon :size="20"><component :is="item.icon" /></el-icon>
        <span v-if="!shouldCollapse">{{ item.label }}</span>
      </RouterLink>
    </nav>

    <div class="sidebar-footer">
      <div v-if="!shouldCollapse" class="planning-card">
        <div class="planning-badge">✦ 产品路线</div>
        <strong>更多 AI 能力</strong>
        <span>功能规划中</span>
      </div>
      <div v-if="!shouldCollapse" class="progress-card">
        <div><span>本周求职进度</span><strong>0%</strong></div>
        <span class="progress-track"><i /></span>
        <small>新用户暂无任务记录</small>
      </div>
    </div>
  </div>
</template>

<style scoped>
.side-navigation {
  display: flex;
  width: var(--sidebar-width);
  height: 100%;
  min-height: 100vh;
  padding: 0 16px 18px;
  flex-direction: column;
  border-right: 1px solid var(--border-color);
  background: rgba(255, 255, 255, 0.96);
  transition: width 0.25s ease;
}

.side-navigation.collapsed {
  width: 84px;
  padding-inline: 12px;
}

.logo-area {
  display: flex;
  height: var(--topbar-height);
  align-items: center;
  padding-inline: 6px;
}

.nav-list {
  display: grid;
  gap: 8px;
  padding-top: 22px;
}

.nav-item {
  display: flex;
  height: 50px;
  padding: 0 15px;
  align-items: center;
  gap: 14px;
  border-radius: 13px;
  color: var(--text-regular);
  font-weight: 600;
  transition: 0.2s ease;
}

.nav-item:hover {
  color: var(--color-primary);
  background: #f4f7ff;
}

.nav-item.active {
  color: var(--color-primary);
  background: linear-gradient(100deg, #edf3ff, #f0efff);
}

.collapsed .nav-item {
  justify-content: center;
  padding: 0;
}

.sidebar-footer {
  display: grid;
  margin-top: auto;
  gap: 12px;
}

.planning-card,
.progress-card {
  overflow: hidden;
  border: 1px solid var(--border-color);
  border-radius: 14px;
}

.planning-card {
  display: grid;
  padding: 16px;
  gap: 4px;
  background: linear-gradient(145deg, #eaf3ff, #e7d9ff);
}

.planning-badge {
  margin-bottom: 4px;
  color: var(--color-primary);
  font-weight: 700;
}

.planning-card span {
  color: var(--text-regular);
  font-size: 12px;
}

.progress-card {
  display: grid;
  padding: 14px;
  gap: 10px;
  background: #fff;
}

.progress-card > div {
  display: flex;
  justify-content: space-between;
}

.progress-track {
  height: 7px;
  overflow: hidden;
  border-radius: 999px;
  background: #edf0f7;
}

.progress-track i {
  display: block;
  width: 0;
  height: 100%;
}

.progress-card small {
  color: var(--text-secondary);
}
</style>

