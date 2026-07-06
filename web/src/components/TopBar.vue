<template>
  <div class="topbar">
    <div class="topbar-left">
      <span class="topbar-title">{{ settings.pageTitle || '个人导航页' }}</span>
    </div>
    <div class="topbar-right">
      <el-button text circle @click="$emit('openImages')" title="图库">
        <el-icon :size="20"><Picture /></el-icon>
      </el-button>
      <el-button v-if="auth.isAdmin" text circle @click="$emit('openUsers')" title="用户管理">
        <el-icon :size="20"><UserFilled /></el-icon>
      </el-button>
      <el-dropdown trigger="click" @command="handleCommand">
        <span class="topbar-user">
          <span class="topbar-avatar">{{ auth.user?.avatar || '👤' }}</span>
          <span class="topbar-name">{{ auth.user?.username }}</span>
          <el-icon><ArrowDown /></el-icon>
        </span>
        <template #dropdown>
          <el-dropdown-menu>
            <el-dropdown-item command="settings">
              <el-icon><Setting /></el-icon>设置
            </el-dropdown-item>
            <el-dropdown-item command="password">
              <el-icon><EditPen /></el-icon>修改密码
            </el-dropdown-item>
            <el-dropdown-item command="logout" divided>
              <el-icon><SwitchButton /></el-icon>退出登录
            </el-dropdown-item>
          </el-dropdown-menu>
        </template>
      </el-dropdown>
    </div>
    <!-- 修改密码弹窗 -->
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
import { useSettingsStore } from '@/stores/settings'
import { changePassword } from '@/api'
import { ElMessage } from 'element-plus'
import { Picture, Setting, EditPen, SwitchButton, ArrowDown, UserFilled } from '@element-plus/icons-vue'

const emit = defineEmits(['openSettings', 'openImages', 'openUsers'])
const router = useRouter()
const auth = useAuthStore()
const settings = useSettingsStore()

const pwdVisible = ref(false)
const pwdLoading = ref(false)
const pwdFormRef = ref(null)
const pwdForm = reactive({ oldPassword: '', newPassword: '' })
const pwdRules = {
  oldPassword: [{ required: true, message: '请输入旧密码', trigger: 'blur' }],
  newPassword: [{ required: true, min: 6, message: '新密码至少6位', trigger: 'blur' }]
}

function handleCommand(cmd) {
  if (cmd === 'settings') emit('openSettings')
  else if (cmd === 'password') pwdVisible.value = true
  else if (cmd === 'logout') { auth.logout(); router.push('/login') }
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
  justify-content: space-between;
  padding: 16px 0;
  border-bottom: 1px solid rgba(0, 0, 0, 0.06);
}
.topbar-title {
  font-size: 18px;
  font-weight: 700;
  color: #1d1d1f;
}
.topbar-right {
  display: flex;
  align-items: center;
  gap: 4px;
}
.topbar-user {
  display: flex;
  align-items: center;
  gap: 8px;
  cursor: pointer;
  padding: 4px 12px;
  border-radius: 20px;
  transition: background 0.2s;
}
.topbar-user:hover {
  background: rgba(0, 0, 0, 0.04);
}
.topbar-avatar {
  font-size: 24px;
  line-height: 1;
}
.topbar-name {
  font-size: 14px;
  color: #1d1d1f;
  font-weight: 500;
}
</style>
