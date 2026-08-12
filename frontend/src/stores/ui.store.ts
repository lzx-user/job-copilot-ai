import { defineStore } from 'pinia'
import { ref } from 'vue'

export const useUiStore = defineStore('ui', () => {
  const sidebarCollapsed = ref(false)
  const mobileDrawerVisible = ref(false)

  function setSidebarCollapsed(value: boolean) {
    sidebarCollapsed.value = value
  }

  function toggleMobileDrawer() {
    mobileDrawerVisible.value = !mobileDrawerVisible.value
  }

  function closeMobileDrawer() {
    mobileDrawerVisible.value = false
  }

  return {
    sidebarCollapsed,
    mobileDrawerVisible,
    setSidebarCollapsed,
    toggleMobileDrawer,
    closeMobileDrawer,
  }
})

