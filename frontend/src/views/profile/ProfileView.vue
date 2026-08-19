<script setup lang="ts">
import { InfoFilled } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'
import { computed, onMounted, reactive, watch } from 'vue'
import { useAuthStore } from '../../stores/auth.store'
import { useProfileStore } from '../../stores/profile.store'
import type { CandidateProfile } from '../../types/profile'

const authStore = useAuthStore()
const email = computed(() => authStore.user?.email || '尚未获取邮箱')
const avatarLetter = computed(() => email.value.slice(0, 1).toUpperCase())
const profileStore = useProfileStore()

// 页面首次挂载时读取数据库，刷新页面后也会重新执行
onMounted(() => {
  void profileStore.fetchProfile()
})

const profileForm = reactive<CandidateProfile>({
  nickname: '',
  targetRoles: [],
  expectedCities: [],
  skills: [],
  projectSummary: '',
  strengths: '',
  availability: '',
  graduationYear: '',
})

// fetchProfile 异步完成后，把 Store 的真实数据复制到页面表单。
watch(
  () => profileStore.profile,
  (profile) => {
    Object.assign(profileForm, {
      ...profile,
      // 复制数组，避免页面直接修改 Store 内部的数组
      targetRoles: [...profile.targetRoles],
      expectedCities: [...profile.expectedCities],
      skills: [...profile.skills],
    })
  },
  // 创建监听后立刻执行一次，不必等第一次变化。
  { immediate: true },
)

async function handleSave() {
  const result = await profileStore.saveProfile({
    ...profileForm,
    // 向 Store 传入数组副本，避免数据库操作持有页面的响应式数组。
    targetRoles: [...profileForm.targetRoles],
    expectedCities: [...profileForm.expectedCities],
    skills: [...profileForm.skills],
  })

  if (result === 'ignored') {
    // 重复提交或认证身份已变化时静默忽略，不向当前用户显示旧请求提示
    return
  }

  if (result === 'error') {
    ElMessage.error(profileStore.profileError || '档案保存失败，请稍后重试')
    return
  }

  ElMessage.success('个人档案已保存')
}
</script>

<template>
  <div class="page-shell">
    <div class="page-heading">
      <div>
        <span class="section-label">CANDIDATE PROFILE</span>
        <h1>个人档案</h1>
        <p>完善个人信息后，AI 才能提供更贴近你的 JD 分析与面试练习。</p>
      </div>
      <el-alert
        class="stage-alert"
        title="档案数据将按当前登录用户进行隔离"
        type="info"
        :closable="false"
        show-icon
      />
    </div>

    <el-alert
      v-if="profileStore.profileError"
      class="profile-error"
      :title="profileStore.profileError"
      type="error"
      :closable="false"
      show-icon
    />

    <div class="profile-layout">
      <main class="profile-main">
        <section class="panel identity-card">
          <div class="avatar">{{ avatarLetter }}</div>
          <div class="identity-copy">
            <span class="identity-tag">个人档案</span>
            <h2>{{ profileForm.nickname || '待完善昵称' }}</h2>
            <p>
              {{
                profileForm.targetRoles.length
                  ? profileForm.targetRoles.join(' / ')
                  : '目标岗位尚未填写'
              }}
            </p>
            <small>{{ email }}</small>
          </div>
          <div class="mini-ring"><strong>0%</strong><span>已完善</span></div>
        </section>

        <section v-loading="profileStore.loading" class="panel section-card">
          <div class="panel-header">
            <h2 class="panel-title">基础信息</h2>
          </div>

          <el-form
            :model="profileForm"
            label-position="top"
            :disabled="profileStore.loading || profileStore.saving"
          >
            <div class="profile-form-grid">
              <el-form-item label="昵称">
                <el-input
                  v-model="profileForm.nickname"
                  placeholder="例如：小林"
                  maxlength="30"
                  show-word-limit
                />
              </el-form-item>

              <el-form-item label="毕业年份">
                <el-input
                  v-model="profileForm.graduationYear"
                  placeholder="例如：2027"
                  maxlength="4"
                />
              </el-form-item>

              <el-form-item label="目标岗位">
                <el-select
                  v-model="profileForm.targetRoles"
                  multiple
                  filterable
                  allow-create
                  default-first-option
                  placeholder="输入岗位后按回车添加"
                />
              </el-form-item>

              <el-form-item label="期望城市">
                <el-select
                  v-model="profileForm.expectedCities"
                  multiple
                  filterable
                  allow-create
                  default-first-option
                  placeholder="输入城市后按回车添加"
                />
              </el-form-item>

              <el-form-item label="可实习 / 到岗时间">
                <el-input
                  v-model="profileForm.availability"
                  placeholder="例如：2026 年 9 月起，可连续实习 6 个月"
                />
              </el-form-item>

              <el-form-item label="登录邮箱">
                <el-input :model-value="email" disabled />
              </el-form-item>
            </div>
          </el-form>
        </section>

        <section class="panel section-card">
          <div class="panel-header">
            <h2 class="panel-title">技术栈与关键词</h2>
          </div>

          <el-select
            v-model="profileForm.skills"
            class="profile-tags-select"
            multiple
            filterable
            allow-create
            default-first-option
            :multiple-limit="30"
            :disabled="profileStore.loading || profileStore.saving"
            placeholder="输入技能后按回车添加，例如 Vue 3、TypeScript、Go"
          />
        </section>

        <section class="panel section-card">
          <div class="panel-header">
            <h2 class="panel-title">项目经历摘要</h2>
          </div>

          <el-input
            v-model="profileForm.projectSummary"
            type="textarea"
            :rows="6"
            maxlength="1500"
            show-word-limit
            :disabled="profileStore.loading || profileStore.saving"
            placeholder="介绍项目背景、你的职责、关键难点、解决方案和可验证结果"
          />
        </section>

        <section class="panel section-card">
          <div class="panel-header">
            <h2 class="panel-title">个人优势 / 自我介绍</h2>
          </div>

          <el-input
            v-model="profileForm.strengths"
            type="textarea"
            :rows="5"
            maxlength="1000"
            show-word-limit
            :disabled="profileStore.loading || profileStore.saving"
            placeholder="结合真实经历，说明你的技术特点、学习能力和协作优势"
          />
        </section>

        <div class="profile-actions">
          <el-button
            class="gradient-button"
            :loading="profileStore.saving"
            :disabled="profileStore.loading"
            @click="handleSave"
          >
            保存档案
          </el-button>
          <el-button disabled>预览简历摘要</el-button>
        </div>
      </main>

      <aside class="profile-aside">
        <section class="panel completeness-card">
          <h2 class="panel-title">档案完整度</h2>
          <div class="large-ring">
            <div><strong>0%</strong><span>待完善</span></div>
          </div>
          <ul>
            <li><i />基础信息待填写</li>
            <li><i />技术栈待补充</li>
            <li><i />项目经历待填写</li>
            <li><i />个人优势待填写</li>
          </ul>
        </section>
        <section class="panel advice-card">
          <span class="advice-icon">✦</span>
          <div>
            <h3>档案填写建议</h3>
            <p>
              只填写真实经历。清晰、具体、可验证的信息，比堆叠关键词更有帮助。
            </p>
          </div>
        </section>
        <section class="panel privacy-card">
          <el-icon><InfoFilled /></el-icon>
          <div>
            <h3>数据说明</h3>
            <p>
              个人档案通过 Supabase RLS
              按登录用户隔离，页面不会让用户自行指定数据所有者。
            </p>
          </div>
        </section>
      </aside>
    </div>
  </div>
</template>

<style scoped>
.stage-alert {
  width: min(100%, 390px);
}
.profile-layout {
  display: grid;
  grid-template-columns: minmax(0, 1fr) 350px;
  gap: 18px;
}
.profile-main,
.profile-aside {
  display: grid;
  align-content: start;
  gap: 16px;
}
.profile-form-grid {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 4px 16px;
}
.profile-form-grid :deep(.el-form-item) {
  margin-bottom: 16px;
}
.profile-tags-select {
  width: 100%;
}
.profile-error {
  margin-bottom: 16px;
}
.identity-card {
  display: flex;
  min-height: 150px;
  padding: 24px 28px;
  align-items: center;
  gap: 20px;
  background: linear-gradient(125deg, #fff, #f8f9ff);
}
.avatar {
  display: grid;
  flex: 0 0 78px;
  width: 78px;
  height: 78px;
  place-items: center;
  border: 7px solid #eef2ff;
  border-radius: 50%;
  color: #fff;
  background: linear-gradient(140deg, #7487a8, #3e4d69);
  font-size: 27px;
  font-weight: 800;
}
.identity-copy {
  min-width: 0;
}
.identity-copy h2 {
  margin: 5px 0 2px;
}
.identity-copy p {
  margin: 0;
  color: var(--text-regular);
}
.identity-copy small {
  display: block;
  overflow: hidden;
  max-width: 360px;
  margin-top: 7px;
  color: var(--text-secondary);
  text-overflow: ellipsis;
  white-space: nowrap;
}
.identity-tag {
  padding: 4px 8px;
  border-radius: 7px;
  color: var(--color-primary);
  background: #edf3ff;
  font-size: 12px;
}
.mini-ring {
  display: grid;
  width: 82px;
  height: 82px;
  margin-left: auto;
  place-items: center;
  align-content: center;
  border: 7px solid #edf0f6;
  border-radius: 50%;
}
.mini-ring strong {
  font-size: 20px;
}
.mini-ring span {
  color: var(--text-secondary);
  font-size: 11px;
}
.section-card {
  padding: 22px;
}
.info-grid {
  display: grid;
  grid-template-columns: repeat(3, minmax(0, 1fr));
  gap: 12px;
}
.info-grid > div {
  display: grid;
  min-width: 0;
  padding: 12px 14px;
  gap: 4px;
  border-radius: 10px;
  background: #f8f9fd;
}
.info-grid span {
  color: var(--text-secondary);
  font-size: 12px;
}
.info-grid strong {
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}
.empty-inline {
  display: flex;
  min-height: 70px;
  padding: 18px;
  align-items: center;
  gap: 12px;
  border: 1px dashed #d8dfef;
  border-radius: 12px;
  color: var(--text-secondary);
  background: #fbfcff;
}
.text-section p {
  margin: 0;
  color: var(--text-secondary);
}
.profile-actions {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 14px;
}
.completeness-card {
  display: grid;
  padding: 22px;
  place-items: center;
}
.completeness-card h2 {
  justify-self: start;
}
.large-ring {
  display: grid;
  width: 150px;
  height: 150px;
  margin: 20px 0;
  place-items: center;
  border: 11px solid #edf0f6;
  border-radius: 50%;
}
.large-ring > div {
  display: grid;
  place-items: center;
}
.large-ring strong {
  font-size: 34px;
}
.large-ring span {
  color: var(--text-secondary);
}
.completeness-card ul {
  display: grid;
  width: 100%;
  gap: 11px;
  list-style: none;
  color: var(--text-regular);
}
.completeness-card li {
  display: flex;
  align-items: center;
  gap: 9px;
}
.completeness-card li i {
  width: 17px;
  height: 17px;
  border: 2px solid #d2d9e8;
  border-radius: 50%;
}
.advice-card,
.privacy-card {
  display: flex;
  padding: 20px;
  gap: 12px;
}
.advice-icon {
  display: grid;
  flex: 0 0 38px;
  width: 38px;
  height: 38px;
  place-items: center;
  border-radius: 12px;
  color: var(--color-secondary);
  background: #f1eaff;
}
.advice-card h3,
.privacy-card h3 {
  margin-bottom: 5px;
}
.advice-card p,
.privacy-card p {
  margin: 0;
  color: var(--text-secondary);
  font-size: 13px;
}
.privacy-card > .el-icon {
  flex: 0 0 auto;
  margin-top: 3px;
  color: var(--color-primary);
}
@media (max-width: 1100px) {
  .profile-layout {
    grid-template-columns: 1fr;
  }
  .profile-aside {
    grid-template-columns: repeat(2, minmax(0, 1fr));
  }
  .completeness-card {
    grid-row: span 2;
  }
}
@media (max-width: 700px) {
  .stage-alert {
    width: 100%;
  }
  .identity-card {
    padding: 20px;
    flex-wrap: wrap;
  }
  .mini-ring {
    margin-left: 0;
  }
  .info-grid {
    grid-template-columns: 1fr;
  }
  .profile-aside {
    grid-template-columns: 1fr;
  }
  .profile-actions {
    grid-template-columns: 1fr;
  }
  .profile-form-grid {
    grid-template-columns: 1fr;
  }
}
</style>
