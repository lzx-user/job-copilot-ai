<script setup lang="ts">
import { ElMessage, type FormInstance, type FormRules } from 'element-plus'
import { reactive, ref } from 'vue'
import { useRouter } from 'vue-router'
import AuthShell from '../../components/common/AuthShell.vue'
import { isSupabaseConfigured } from '../../lib/supabase'
import { useAuthStore } from '../../stores/auth.store'
import type { RegisterFormValues } from '../../types/auth'
import { isValidEmail } from '../../utils/validators'

const router = useRouter()
const authStore = useAuthStore()
const formRef = ref<FormInstance>()

const form = reactive<RegisterFormValues>({
  email: '',
  password: '',
  confirmPassword: '',
  acceptedTerms: false,
})

const rules: FormRules<RegisterFormValues> = {
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
  confirmPassword: [
    { required: true, message: '请再次输入密码', trigger: 'blur' },
    {
      validator: (_rule, value: string, callback) => {
        value === form.password ? callback() : callback(new Error('两次输入的密码不一致'))
      },
      trigger: ['blur', 'change'],
    },
  ],
  acceptedTerms: [
    {
      validator: (_rule, value: boolean, callback) => {
        value ? callback() : callback(new Error('请先阅读并同意使用说明与隐私提示'))
      },
      trigger: 'change',
    },
  ],
}

async function submitRegister() {
  authStore.clearError()
  try {
    await formRef.value?.validate()
  } catch {
    return
  }

  const result = await authStore.register(form.email.trim(), form.password)
  if (!result.success) return

  if (result.requiresEmailConfirmation) {
    ElMessage.success('注册成功，请前往邮箱完成验证后登录')
    await router.replace({ path: '/auth/login', query: { email: form.email.trim() } })
    return
  }

  ElMessage.success('注册成功，欢迎加入')
  await router.replace('/app/dashboard')
}
</script>

<template>
  <AuthShell eyebrow="Create account" title="创建你的账号" subtitle="用邮箱注册，开始建立连续的求职准备记录。">
    <el-alert
      v-if="!isSupabaseConfigured"
      class="config-alert"
      title="Supabase 尚未配置"
      description="配置完成后，这里会调用真实的邮箱注册，不使用本地假账号。"
      type="warning"
      :closable="false"
      show-icon
    />
    <el-alert v-if="authStore.authError" class="error-alert" :title="authStore.authError" type="error" :closable="false" show-icon />

    <el-form ref="formRef" :model="form" :rules="rules" label-position="top" @submit.prevent="submitRegister">
      <el-form-item label="邮箱" prop="email">
        <el-input v-model="form.email" type="email" autocomplete="email" placeholder="name@example.com" />
      </el-form-item>
      <el-form-item label="密码" prop="password">
        <el-input v-model="form.password" type="password" autocomplete="new-password" placeholder="至少 8 位" show-password />
      </el-form-item>
      <el-form-item label="确认密码" prop="confirmPassword">
        <el-input v-model="form.confirmPassword" type="password" autocomplete="new-password" placeholder="再次输入密码" show-password @keyup.enter="submitRegister" />
      </el-form-item>
      <el-form-item prop="acceptedTerms">
        <el-checkbox v-model="form.acceptedTerms">
          我已阅读并同意产品使用说明与数据隐私提示
        </el-checkbox>
      </el-form-item>
      <el-button
        class="gradient-button submit-button"
        native-type="submit"
        :loading="authStore.submitting"
        :disabled="authStore.submitting || !isSupabaseConfigured"
      >
        注册账号
      </el-button>
    </el-form>

    <p class="privacy-note">我们不会在源码或日志中记录你的密码与登录 token。</p>
    <p class="switch-link">已经有账号？<RouterLink to="/auth/login">直接登录</RouterLink></p>
  </AuthShell>
</template>

<style scoped>
.config-alert,
.error-alert {
  margin-bottom: 18px;
}

.submit-button {
  width: 100%;
  margin-top: 4px;
}

.privacy-note {
  margin: 16px 0 0;
  color: var(--text-secondary);
  font-size: 12px;
  text-align: center;
}

.switch-link {
  margin: 14px 0 0;
  color: var(--text-secondary);
  text-align: center;
}

.switch-link a {
  color: var(--color-primary);
  font-weight: 700;
}

:deep(.el-checkbox) {
  height: auto;
  align-items: flex-start;
  white-space: normal;
}

:deep(.el-checkbox__label) {
  line-height: 1.5;
}
</style>
