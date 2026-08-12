<script setup lang="ts">
import { Check } from '@element-plus/icons-vue'
import BrandLogo from './BrandLogo.vue'

interface Props {
  eyebrow: string
  title: string
  subtitle: string
}

defineProps<Props>()
</script>

<template>
  <main class="auth-page">
    <section class="value-panel">
      <BrandLogo />
      <div class="value-copy">
        <span class="eyebrow">AI 全程陪伴你的求职准备</span>
        <h1>把每一次准备，<br />变成更清晰的方向。</h1>
        <p>从个人档案到岗位分析与模拟面试，用一个连续的工作流沉淀你的求职进步。</p>
        <ul>
          <li><span><el-icon><Check /></el-icon></span>拆解岗位要求，梳理匹配方向</li>
          <li><span><el-icon><Check /></el-icon></span>针对目标岗位，规划面试练习</li>
          <li><span><el-icon><Check /></el-icon></span>统一记录过程，方便持续复盘</li>
        </ul>
      </div>
      <div class="visual-orbit" aria-hidden="true">
        <div class="orbit-card main-card"><b>JD</b><i /><i /><i /></div>
        <div class="orbit-card small-card">✓</div>
        <span class="orbit-dot dot-a" />
        <span class="orbit-dot dot-b" />
      </div>
      <p class="stage-note">当前版本完成登录基础设施与页面骨架，AI 业务将在后续阶段接入。</p>
    </section>

    <section class="form-panel">
      <div class="mobile-brand"><BrandLogo /></div>
      <div class="auth-card">
        <span class="form-eyebrow">{{ eyebrow }}</span>
        <h2>{{ title }}</h2>
        <p class="form-subtitle">{{ subtitle }}</p>
        <slot />
      </div>
      <p class="copyright">© 2026 Job Copilot AI · 认真准备，也要保护好你的数据</p>
    </section>
  </main>
</template>

<style scoped>
.auth-page {
  display: grid;
  min-height: 100vh;
  grid-template-columns: minmax(420px, 1.05fr) minmax(480px, 0.95fr);
  background: #fff;
}

.value-panel {
  position: relative;
  display: flex;
  overflow: hidden;
  padding: 38px clamp(36px, 5vw, 76px);
  flex-direction: column;
  background:
    radial-gradient(circle at 75% 65%, rgba(127, 92, 255, 0.24), transparent 28%),
    radial-gradient(circle at 20% 10%, rgba(94, 147, 255, 0.22), transparent 34%),
    linear-gradient(150deg, #f5f8ff 0%, #eef2ff 52%, #f7f1ff 100%);
}

.value-panel::before,
.value-panel::after {
  position: absolute;
  width: 280px;
  height: 280px;
  border: 1px solid rgba(65, 105, 246, 0.13);
  border-radius: 50%;
  content: '';
}

.value-panel::before {
  right: -90px;
  bottom: 8%;
}

.value-panel::after {
  right: -15px;
  bottom: 15%;
  width: 140px;
  height: 140px;
}

.value-copy {
  position: relative;
  z-index: 2;
  max-width: 570px;
  margin-top: clamp(60px, 11vh, 130px);
}

.eyebrow,
.form-eyebrow {
  color: var(--color-primary);
  font-size: 13px;
  font-weight: 700;
  letter-spacing: 0.08em;
}

.value-copy h1 {
  margin: 18px 0 20px;
  font-size: clamp(38px, 4vw, 58px);
  line-height: 1.18;
  letter-spacing: -0.05em;
}

.value-copy > p {
  max-width: 520px;
  color: var(--text-regular);
  font-size: 17px;
}

.value-copy ul {
  display: grid;
  margin-top: 34px;
  gap: 17px;
  list-style: none;
  color: var(--text-regular);
  font-weight: 600;
}

.value-copy li {
  display: flex;
  align-items: center;
  gap: 12px;
}

.value-copy li span {
  display: grid;
  width: 26px;
  height: 26px;
  place-items: center;
  border-radius: 8px;
  color: #fff;
  background: linear-gradient(140deg, var(--color-primary), var(--color-secondary));
}

.visual-orbit {
  position: absolute;
  right: 7%;
  bottom: 7%;
  width: 250px;
  height: 220px;
  transform: rotate(-5deg);
  opacity: 0.78;
}

.orbit-card {
  position: absolute;
  display: grid;
  border: 1px solid rgba(255, 255, 255, 0.95);
  place-items: center;
  color: var(--color-primary);
  background: rgba(255, 255, 255, 0.78);
  box-shadow: 0 24px 60px rgba(65, 105, 246, 0.18);
  backdrop-filter: blur(10px);
}

.main-card {
  right: 30px;
  bottom: 25px;
  width: 150px;
  height: 180px;
  border-radius: 28px;
}

.main-card b {
  font-size: 38px;
}

.main-card i {
  display: block;
  width: 78px;
  height: 7px;
  margin-top: -20px;
  border-radius: 99px;
  background: #d9e2ff;
}

.small-card {
  z-index: 2;
  bottom: 38px;
  left: 12px;
  width: 62px;
  height: 62px;
  border-radius: 50%;
  font-size: 28px;
}

.orbit-dot {
  position: absolute;
  width: 14px;
  height: 14px;
  border-radius: 50%;
  background: var(--color-secondary);
}

.dot-a { top: 0; right: 10px; }
.dot-b { top: 56px; left: 15px; width: 8px; height: 8px; }

.stage-note {
  position: relative;
  z-index: 2;
  max-width: 520px;
  margin: auto 0 0;
  color: var(--text-secondary);
  font-size: 12px;
}

.form-panel {
  display: grid;
  padding: 36px clamp(28px, 5vw, 80px);
  place-items: center;
  align-content: center;
  background:
    radial-gradient(circle at 90% 0%, rgba(114, 91, 255, 0.08), transparent 25%),
    #fff;
}

.auth-card {
  width: min(100%, 460px);
  padding: 46px;
  border: 1px solid var(--border-color);
  border-radius: 24px;
  background: rgba(255, 255, 255, 0.94);
  box-shadow: 0 28px 80px rgba(45, 72, 130, 0.1);
}

.auth-card h2 {
  margin: 10px 0 8px;
  font-size: 30px;
  line-height: 1.3;
  letter-spacing: -0.03em;
}

.form-subtitle {
  margin-bottom: 28px;
  color: var(--text-secondary);
}

.mobile-brand {
  display: none;
}

.copyright {
  margin: 28px 0 0;
  color: var(--text-secondary);
  font-size: 12px;
  text-align: center;
}

@media (max-width: 980px) {
  .auth-page {
    grid-template-columns: 1fr;
  }

  .value-panel {
    display: none;
  }

  .form-panel {
    min-height: 100vh;
    padding: 26px;
  }

  .mobile-brand {
    display: block;
    width: min(100%, 460px);
    margin-bottom: 28px;
  }
}

@media (max-width: 560px) {
  .form-panel {
    padding: 22px 16px;
    place-items: stretch;
  }

  .auth-card {
    width: 100%;
    padding: 28px 22px;
    border-radius: 18px;
  }

  .auth-card h2 {
    font-size: 26px;
  }
}
</style>

