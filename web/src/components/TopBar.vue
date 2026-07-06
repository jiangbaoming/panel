<template>
  <div class="topbar">
    <!-- 左侧留空 -->
    <div class="topbar-left"></div>
    <!-- 右侧操作按钮 -->
    <div class="topbar-right">
      <el-button text @click="$emit('openSettings')" title="设置">
        <el-icon :size="20"><Setting /></el-icon>
        <span class="btn-label">设置</span>
      </el-button>
      <el-button text @click="$emit('openImages')" title="图库">
        <el-icon :size="20"><Picture /></el-icon>
        <span class="btn-label">图库</span>
      </el-button>
      <el-button v-if="auth.isAdmin" text @click="$emit('openUsers')" title="用户管理">
        <el-icon :size="20"><UserFilled /></el-icon>
        <span class="btn-label">用户</span>
      </el-button>
      <el-button text @click="handleLogout" title="退出登录">
        <el-icon :size="20"><SwitchButton /></el-icon>
        <span class="btn-label">退出</span>
      </el-button>
    </div>
    <!-- 修改密码弹窗（独立） -->
    <el-dialog v-model="pwdVisible" title="修改密码" width="400px">
      <el-form ref="pwdFormRef" :model="pwdForm" :rules="pwdRules" label-width="0">
        <el-form-item prop="oldPassword">
          <el-input v-model="pwdForm.oldPassword" type="password" placeholder="旧密码" show-password />
        </el-form-item>
        <el-form-item prop="newPassword">
          <el-input v-model="pwdForm.newPassword" type="password" placeholder="新密码（至少6位）" show-password />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="pwdVisible = false">取消</el-button>
        <el-button type="primary" :loading="pwdLoading" @click="handleChangePwd">确认</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { changePassword } from '@/api'
import { ElMessage } from 'element-plus'
import { Setting, Picture, UserFilled, SwitchButton } from '@element-plus/icons-vue'

const emit = defineEmits(['openSettings', 'openImages', 'openUsers'])
const router = useRouter()
const auth = useAuthStore()

const pwdVisible = ref(false)
const pwdLoading = ref(false)
const pwdFormRef = ref(null)
const pwdForm = reactive({ oldPassword: '', newPassword: '' })
const pwdRules = {
  oldPassword: [{ required: true, message: '请输入旧密码', trigger: 'blur' }],
  newPassword: [{ required: true, min: 6, message: '新密码至少6位', trigger: 'blur' }]
}

defineExpose({ openPassword: () => { pwdVisible.value = true } })

async function handleLogout() {
  auth.logout()
  router.push('/login')
}

async function handleChangePwd() {
  const valid = await pwdFormRef.value.validate().catch(() => false)
  if (!valid) return
  pwdLoading.value = true
  try {
    await changePassword(pwdForm.oldPassword, pwdForm.newPassword)
    ElMessage.success('密码修改成功')
    pwdVisible.value = false
    pwdForm.oldPassword = ''
    pwdForm.newPassword = ''
  } catch (e) {
    ElMessage.error(e.response?.data?.error || '修改失败')
  } finally {
    pwdLoading.value = false
  }
}
</script>

<style scoped>
.topbar {
  display: flex;
  align-items: center;
  justify-content: flex-end;
  padding: 16px 24px;
  max-width: 1100px;
  margin: 0 auto;
  width: 100%;
}
.topbar-left {
  flex: 1;
}
.topbar-right {
  display: flex;
  align-items: center;
  gap: 4px;
}
.btn-label {
  font-size: 13px;
  margin-left: 4px;
}
</style>
