<script setup lang="ts">
interface Props {
  title: string
  description: string
  icon?: string
}

withDefaults(defineProps<Props>(), {
  icon: '✦',
})
</script>

<template>
  <div class="empty-state">
    <div class="empty-icon" aria-hidden="true">{{ icon }}</div>
    <h3>{{ title }}</h3>
    <p>{{ description }}</p>
    <slot />
  </div>
</template>

<style scoped>
.empty-state {
  display: grid;
  min-height: 210px;
  padding: 36px 20px;
  place-items: center;
  align-content: center;
  text-align: center;
}

.empty-icon {
  display: grid;
  width: 58px;
  height: 58px;
  margin-bottom: 16px;
  place-items: center;
  border-radius: 18px;
  color: var(--color-primary);
  background: linear-gradient(145deg, #edf3ff, #f3ecff);
  font-size: 26px;
}

h3 {
  margin-bottom: 6px;
  font-size: 17px;
}

p {
  max-width: 420px;
  margin-bottom: 18px;
  color: var(--text-secondary);
}
</style>

