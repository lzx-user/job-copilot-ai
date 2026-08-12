<script setup lang="ts">
interface Props {
  compact?: boolean
}

withDefaults(defineProps<Props>(), {
  compact: false,
})
</script>

<template>
  <div class="brand" :class="{ 'is-compact': compact }" aria-label="求职陪跑 AI 助手">
    <span class="brand-mark" aria-hidden="true">✦</span>
    <span v-if="!compact" class="brand-name">求职陪跑 <strong>AI</strong> 助手</span>
  </div>
</template>

<style scoped>
.brand {
  display: flex;
  align-items: center;
  gap: 12px;
  min-width: 0;
}

.brand-mark {
  display: grid;
  flex: 0 0 38px;
  width: 38px;
  height: 38px;
  place-items: center;
  border-radius: 12px;
  color: #fff;
  background: linear-gradient(140deg, #2f5df4 5%, #8d62ff 92%);
  box-shadow: 0 8px 18px rgba(65, 105, 246, 0.28);
  font-size: 22px;
}

.brand-name {
  overflow: hidden;
  color: var(--text-primary);
  font-size: 18px;
  font-weight: 700;
  letter-spacing: -0.02em;
  white-space: nowrap;
}

.brand-name strong {
  color: var(--color-primary);
}

.is-compact {
  justify-content: center;
}
</style>

