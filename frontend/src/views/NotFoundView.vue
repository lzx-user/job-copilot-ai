<script setup lang="ts">
import { useRouter } from 'vue-router'

const router = useRouter()
</script>

<template>
  <main class="not-found">
    <div class="code">404</div>
    <h1>这个页面暂时不存在</h1>
    <p>地址可能输入有误，或者页面已经移动。</p>
    <el-button class="gradient-button" @click="router.push('/')">返回首页</el-button>
  </main>
</template>

<style scoped>
.not-found {
  display: grid;
  min-height: 100vh;
  padding: 24px;
  place-items: center;
  align-content: center;
  text-align: center;
}

.code {
  color: var(--color-primary);
  font-size: clamp(70px, 12vw, 132px);
  font-weight: 800;
  line-height: 1;
  letter-spacing: -0.08em;
}

h1 {
  margin: 18px 0 8px;
}

p {
  margin-bottom: 24px;
  color: var(--text-secondary);
}
</style>

