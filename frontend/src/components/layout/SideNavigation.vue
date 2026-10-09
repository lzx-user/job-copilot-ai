<script setup lang="ts">
import {
  Clock,
  Briefcase,
  DocumentChecked,
  House,
  Microphone,
  Setting,
  Tickets,
  User,
} from '@element-plus/icons-vue'
import type { Component } from 'vue'
import { computed, ref } from 'vue'
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
const hovered = ref(false)
const focusWithin = ref(false)
const items: NavigationItem[] = [
  { label: '仪表盘', path: '/app/dashboard', icon: House },
  { label: '岗位与投递', path: '/app/opportunities', icon: Briefcase },
  { label: '简历资产', path: '/app/resumes', icon: Tickets },
  { label: 'JD 分析', path: '/app/jd-analysis', icon: DocumentChecked },
  { label: '模拟面试', path: '/app/interviews', icon: Microphone },
  { label: '历史记录', path: '/app/history', icon: Clock },
  { label: '个人档案', path: '/app/profile', icon: User },
  { label: '设置', path: '/app/settings', icon: Setting },
]

const shouldCollapse = computed(() => props.collapsed && !props.drawer)
const previewExpanded = computed(() => shouldCollapse.value && (hovered.value || focusWithin.value))
const isVisuallyCollapsed = computed(() => shouldCollapse.value && !previewExpanded.value)

function handleFocusOut(event: FocusEvent) {
  const navigation = event.currentTarget as HTMLElement
  focusWithin.value = navigation.contains(event.relatedTarget as Node | null)
}

function isActive(path: string) {
  return path === '/app/interviews' || path === '/app/opportunities'
    ? route.path.startsWith(path)
    : route.path === path
}
</script>

<template>
  <div
    class="side-navigation"
    :class="{
      collapsed: isVisuallyCollapsed,
      'preview-expanded': previewExpanded,
    }"
    @mouseenter="hovered = true"
    @mouseleave="hovered = false"
    @focusin="focusWithin = true"
    @focusout="handleFocusOut"
  >
    <div class="logo-area">
      <BrandLogo :compact="isVisuallyCollapsed" />
    </div>

    <nav class="nav-list" aria-label="主要导航">
      <RouterLink
        v-for="item in items"
        :key="item.path"
        :to="item.path"
        class="nav-item"
        :class="{ active: isActive(item.path) }"
        :title="isVisuallyCollapsed ? item.label : undefined"
        @click="emit('navigate')"
      >
        <el-icon :size="20"><component :is="item.icon" /></el-icon>
        <span :aria-hidden="isVisuallyCollapsed">{{ item.label }}</span>
      </RouterLink>
    </nav>

    <div class="sidebar-footer">
      <div class="planning-card" :aria-hidden="isVisuallyCollapsed">
        <div class="planning-badge">✦ 求职主线</div>
        <strong>以岗位为中心</strong>
        <span>匹配 → 投递 → 面试 → Offer</span>
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
  overflow: hidden;
  border-right: 1px solid var(--border-color);
  background: rgba(255, 255, 255, 0.96);
  transition: width 0.25s ease, box-shadow 0.2s ease;
}

.side-navigation.collapsed {
  width: 84px;
  padding-inline: 12px;
}

.side-navigation.preview-expanded {
  box-shadow: 16px 0 32px rgba(45, 72, 130, 0.14);
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
  position: relative;
  display: flex;
  height: 50px;
  padding: 0 15px;
  align-items: center;
  gap: 14px;
  border-radius: 13px;
  color: var(--text-regular);
  font-weight: 600;
  transition: color 0.2s ease, background 0.2s ease;
}

.nav-item > span {
  position: absolute;
  left: 49px;
  width: calc(var(--sidebar-width) - 80px);
  opacity: 1;
  transform: translateX(0);
  white-space: nowrap;
  transition: opacity 0.12s ease 0.14s, transform 0.16s ease 0.12s;
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

.collapsed .nav-item > span {
  visibility: hidden;
  opacity: 0;
  transform: translateX(-6px);
  transition: opacity 0.06s ease, transform 0.06s ease, visibility 0s linear 0.06s;
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
  width: calc(var(--sidebar-width) - 32px);
  min-height: 140px;
  padding: 16px;
  gap: 4px;
  background: linear-gradient(145deg, #eaf3ff, #e7d9ff);
  opacity: 1;
  transform: translateX(0);
  transition: opacity 0.14s ease 0.12s, transform 0.18s ease 0.1s;
}

.collapsed .planning-card {
  visibility: hidden;
  opacity: 0;
  transform: translateX(-8px);
  transition: opacity 0.08s ease, transform 0.08s ease, visibility 0s linear 0.08s;
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
