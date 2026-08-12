<script setup lang="ts">
import { ElMessage, type FormInstance, type FormRules } from 'element-plus'
import { computed, onMounted, reactive, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { getHealthStatus } from '../../api/health.api'
import AuthShell from '../../components/common/AuthShell.vue'
import { isSupabaseConfigured } from '../../lib/supabase'
import { useAuthStore } from '../../stores/auth.store'
import type { ServiceStatus } from '../../types/api'
import type { LoginFormValues } from '../../types/auth'
import { isSafeInternalRedirect, isValidEmail } from '../../utils/validators'

const router = useRouter()
const route = useRoute()
const authStore = useAuthStore()
const formRef = ref<FormInstance>()
const serviceStatus = ref<ServiceStatus>('checking')

const form = reactive<LoginFormValues>({
  email: typeof route.query.email === 'string' ? route.query.email : '',
  password: '',
})

const rules: FormRules<LoginFormValues> = {
  email: [
    { required: true, message: '请输入邮箱', trigger: 'blur' },
    {
      validator: (_rule, value: string, callback) => {
        isValidEmail(value) ? callback() : callback(new Error('请输入正确的邮箱格式'))
      },
      trigger: ['blur', 'change'],
    },
  ],
  password: [
    { required: true, message: '请输入密码', trigger: 'blur' },
    { min: 8, message: '密码至少需要 8 位', trigger: ['blur', 'change'] },
  ],
}

const serviceLabel = computed(() => ({
  checking: '检查中',
  online: '服务正常',
  offline: '暂时无法连接',
})[serviceStatus.value])

async function submitLogin() {
  authStore.clearError()
  try {
    // Element Plus 在校验失败时会 reject；此处捕获后直接停止，避免无效请求打到 Supabase。
    await formRef.value?.validate()
  } catch {
    return
  }

  const result = await authStore.login(form.email.trim(), form.password)
  if (!result.success) return

  const redirect = route.query.redirect
  await router.replace(isSafeInternalRedirect(redirect) ? redirect : '/app/dashboard')
  ElMessage.success('登录成功，欢迎回来')
}

onMounted(async () => {
  try {
    const health = await getHealthStatus()
    serviceStatus.value = health.status === 'ok' ? 'online' : 'offline'
    console.info('后端健康检查：', health.status)
  } catch {
    serviceStatus.value = 'offline'
  }
})
</script>

<template>
  <AuthShell eyebrow="Welcome back" title="欢迎回来" subtitle="登录后继续你的求职准备与复盘。">
    <el-alert
      v-if="!isSupabaseConfigured"
      class="config-alert"
      title="Supabase 尚未配置"
      description="请先按照 docs/Supabase配置步骤.md 填写 frontend/.env.local。"
      type="warning"
      :closable="false"
      show-icon
    />
    <el-alert v-if="authStore.authError" class="error-alert" :title="authStore.authError" type="error" :closable="false" show-icon />

    <el-form ref="formRef" :model="form" :rules="rules" label-position="top" @submit.prevent="submitLogin">
      <el-form-item label="邮箱" prop="email">
        <el-input v-model="form.email" type="email" autocomplete="email" placeholder="name@example.com" />
      </el-form-item>
      <el-form-item label="密码" prop="password">
        <el-input v-model="form.password" type="password" autocomplete="current-password" placeholder="至少 8 位" show-password @keyup.enter="submitLogin" />
      </el-form-item>
      <el-button
        class="gradient-button submit-button"
        native-type="submit"
        :loading="authStore.submitting"
        :disabled="authStore.submitting || !isSupabaseConfigured"
      >
        登录
      </el-button>
    </el-form>

    <p class="switch-link">还没有账号？<RouterLink to="/auth/register">免费注册</RouterLink></p>
    <div class="service-status" :class="serviceStatus">
      <span class="status-dot" />
      后端服务状态：{{ serviceLabel }}
    </div>
  </AuthShell>
</template>

<style scoped>
.config-alert,
.error-alert {
  margin-bottom: 18px;
}

.submit-button {
  width: 100%;
  margin-top: 6px;
}

.switch-link {
  margin: 20px 0 0;
  color: var(--text-secondary);
  text-align: center;
}

.switch-link a {
  color: var(--color-primary);
  font-weight: 700;
}

.service-status {
  display: flex;
  margin-top: 22px;
  align-items: center;
  justify-content: center;
  gap: 9px;
  color: var(--text-secondary);
  font-size: 12px;
}

.service-status.online { color: var(--color-success); }
.service-status.offline { color: var(--color-danger); }
</style>
