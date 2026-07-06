<template>
  <div class="topbar">
    <div class="topbar-left"></div>
    <div class="topbar-right">
      <el-button text class="edit-btn" :class="{ active: editMode }" @click="$emit('toggleEdit')" :title="editMode ? '退出编辑' : '编辑模式'">
        <el-icon :size="18">
          <Edit v-if="!editMode" />
          <CloseBold v-else />
        </el-icon>
      </el-button>
      <el-dropdown trigger="click" @command="handleCommand">
        <el-button text class="topbar-trigger">
          <el-icon :size="18"><MoreFilled /></el-icon>
        </el-button>
        <template #dropdown>
          <el-dropdown-menu>
            <el-dropdown-item command="openSettings">
              <el-icon :size="16"><Setting /></el-icon>
              <span>设置</span>
            </el-dropdown-item>
            <el-dropdown-item command="logout" divided>
              <el-icon :size="16"><SwitchButton /></el-icon>
              <span>退出</span>
            </el-dropdown-item>
          </el-dropdown-menu>
        </template>
      </el-dropdown>
    </div>
  </div>
</template>

<script setup>
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { MoreFilled, Edit, CloseBold, Setting, SwitchButton } from '@element-plus/icons-vue'

defineProps({ editMode: Boolean })
const emit = defineEmits(['openSettings', 'toggleEdit'])
const router = useRouter()
const auth = useAuthStore()

function handleCommand(cmd) {
  if (cmd === 'toggleEdit') emit('toggleEdit')
  else if (cmd === 'openSettings') emit('openSettings')
  else if (cmd === 'logout') {
    auth.logout()
    router.push('/login')
  }
}
</script>

<style scoped>
.topbar {
  display: flex;
  align-items: center;
  justify-content: flex-end;
  padding: 12px 20px;
  margin: 0 auto;
  width: 100%;
}
.topbar-left {
  flex: 1;
}
.topbar-right {
  display: flex;
  align-items: center;
  gap: 2px;
}
.edit-btn {
  padding: 8px;
  border-radius: var(--radius-sm);
  transition: background var(--transition-fast), color var(--transition-fast);
  color: rgba(255, 255, 255, 0.65);
}
.edit-btn:hover {
  color: rgba(255, 255, 255, 0.92);
}
.edit-btn.active {
  background: rgba(64, 158, 255, 0.2);
  color: rgba(255, 255, 255, 0.92);
}
.topbar-trigger {
  padding: 8px;
  border-radius: var(--radius-sm);
  transition: background var(--transition-fast);
  color: rgba(255, 255, 255, 0.65);
}
.topbar-trigger:hover {
  color: rgba(255, 255, 255, 0.92);
}

@media (max-width: 768px) {
  .topbar {
    padding: 8px 12px;
  }
}
</style>