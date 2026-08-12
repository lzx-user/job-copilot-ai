<script setup lang="ts">
import { ArrowDown, Bell, Menu, Search } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'
import { computed, onBeforeUnmount, onMounted, ref, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import SideNavigation from '../components/layout/SideNavigation.vue'
import { useAuthStore } from '../stores/auth.store'
import { useUiStore } from '../stores/ui.store'

const route = useRoute()
const router = useRouter()
const authStore = useAuthStore()
const uiStore = useUiStore()
const viewportWidth = ref(typeof window === 'undefined' ? 1440 : window.innerWidth)

const isMobile = computed(() => viewportWidth.value < 768)
const userEmail = computed(() => authStore.user?.email ?? '求职同学')
const avatarLetter = computed(() => userEmail.value.slice(0, 1).toUpperCase())

function syncResponsiveLayout() {
  viewportWidth.value = window.innerWidth
  // 中等屏幕收起侧栏只保留图标；手机端改为抽屉，避免挤压主内容。
  uiStore.setSidebarCollapsed(viewportWidth.value >= 768 && viewportWidth.value < 1200)
  if (viewportWidth.value >= 768) uiStore.closeMobileDrawer()
}

async function handleUserCommand(command: string) {
  if (command !== 'logout') return
  const result = await authStore.logout()
  if (result.success) {
    await router.replace('/auth/login')
    ElMessage.success('已安全退出登录')
  } else {
    ElMessage.error(authStore.authError)
  }
}

watch(
  () => route.fullPath,
  () => uiStore.closeMobileDrawer(),
)

onMounted(() => {
  syncResponsiveLayout()
  window.addEventListener('resize', syncResponsiveLayout, { passive: true })
})

onBeforeUnmount(() => window.removeEventListener('resize', syncResponsiveLayout))
</script>

<template>
  <div class="app-layout">
    <aside v-if="!isMobile" class="desktop-sidebar">
      <SideNavigation :collapsed="uiStore.sidebarCollapsed" />
    </aside>

    <el-drawer
      v-model="uiStore.mobileDrawerVisible"
      class="mobile-nav-drawer"
      direction="ltr"
      :size="270"
      :with-header="false"
    >
      <SideNavigation drawer @navigate="uiStore.closeMobileDrawer" />
    </el-drawer>

    <div class="app-main">
      <header class="topbar">
        <button v-if="isMobile" class="icon-button menu-button" aria-label="打开导航" @click="uiStore.toggleMobileDrawer">
          <el-icon :size="22"><Menu /></el-icon>
        </button>

        <div class="topbar-spacer" />
        <div class="search-placeholder">
          <el-icon><Search /></el-icon>
          <span>搜索职位、公司或功能…</span>
          <kbd>⌘ K</kbd>
        </div>
        <button class="icon-button notification-button" aria-label="通知（功能规划中）">
          <el-icon :size="20"><Bell /></el-icon>
          <i />
        </button>

        <el-dropdown trigger="click" @command="handleUserCommand">
          <button class="user-menu">
            <span class="avatar">{{ avatarLetter }}</span>
            <span class="user-email">{{ userEmail }}</span>
            <el-icon><ArrowDown /></el-icon>
          </button>
          <template #dropdown>
            <el-dropdown-menu>
              <el-dropdown-item command="logout" :disabled="authStore.submitting">退出登录</el-dropdown-item>
            </el-dropdown-menu>
          </template>
        </el-dropdown>
      </header>

      <main class="route-content">
        <RouterView />
      </main>
    </div>
  </div>
</template>

<style scoped>
.app-layout {
  display: flex;
  min-height: 100vh;
}

.desktop-sidebar {
  position: sticky;
  z-index: 20;
  top: 0;
  flex: 0 0 auto;
  height: 100vh;
}

.app-main {
  min-width: 0;
  flex: 1;
}

.topbar {
  position: sticky;
  z-index: 15;
  top: 0;
  display: flex;
  height: var(--topbar-height);
  padding: 0 28px;
  align-items: center;
  gap: 14px;
  border-bottom: 1px solid rgba(232, 236, 245, 0.82);
  background: rgba(255, 255, 255, 0.84);
  backdrop-filter: blur(18px);
}

.topbar-spacer {
  flex: 1;
}

.search-placeholder {
  display: flex;
  width: min(320px, 28vw);
  height: 40px;
  padding: 0 12px;
  align-items: center;
  gap: 9px;
  border: 1px solid var(--border-color);
  border-radius: 11px;
  color: var(--text-secondary);
  background: #fff;
}

.search-placeholder span {
  overflow: hidden;
  flex: 1;
  text-overflow: ellipsis;
  white-space: nowrap;
}

kbd {
  padding: 2px 6px;
  border-radius: 6px;
  color: var(--text-secondary);
  background: #f2f4f8;
  font: 11px inherit;
}

.icon-button,
.user-menu {
  display: inline-flex;
  cursor: pointer;
  align-items: center;
  justify-content: center;
  color: var(--text-regular);
  background: transparent;
}

.icon-button {
  position: relative;
  width: 40px;
  height: 40px;
  border-radius: 10px;
}

.icon-button:hover {
  background: #f3f6fc;
}

.notification-button i {
  position: absolute;
  top: 6px;
  right: 7px;
  width: 7px;
  height: 7px;
  border: 2px solid #fff;
  border-radius: 50%;
  background: var(--color-danger);
}

.user-menu {
  min-width: 0;
  max-width: 230px;
  gap: 9px;
}

.avatar {
  display: grid;
  flex: 0 0 38px;
  width: 38px;
  height: 38px;
  place-items: center;
  border-radius: 50%;
  color: #fff;
  background: linear-gradient(140deg, #7d8fab, #4b5d7a);
  font-weight: 700;
}

.user-email {
  overflow: hidden;
  max-width: 150px;
  color: var(--text-primary);
  text-overflow: ellipsis;
  white-space: nowrap;
}

.route-content {
  min-height: calc(100vh - var(--topbar-height));
}

@media (max-width: 767px) {
  .topbar {
    padding: 0 14px;
    gap: 6px;
  }

  .search-placeholder {
    display: none;
  }

  .user-email {
    display: none;
  }

  .user-menu {
    gap: 4px;
  }
}
</style>

