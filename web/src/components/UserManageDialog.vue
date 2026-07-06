<template>
  <el-dialog title="用户管理" v-model="visible" width="500px" :close-on-click-modal="false">
    <div v-loading="loading">
      <el-table :data="users" size="small" style="width:100%">
        <el-table-column label="头像" width="60">
          <template #default="{ row }">
            <span style="font-size:24px">{{ row.avatar || '👤' }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="username" label="用户名" />
        <el-table-column prop="role" label="角色" width="80" />
        <el-table-column label="操作" width="80" align="center">
          <template #default="{ row }">
            <el-button size="small" type="danger" text @click="removeUser(row)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>
      <div class="user-mgr-add">
        <el-button :icon="Plus" @click="addVisible = true" text>添加用户</el-button>
      </div>
    </div>
    <!-- 添加用户弹窗 -->
    <el-dialog title="添加用户" v-model="addVisible" width="380px" append-to-body>
      <el-form ref="addFormRef" :model="addForm" :rules="addRules" label-width="0">
        <el-form-item prop="username">
          <el-input v-model="addForm.username" placeholder="用户名" />
        </el-form-item>
        <el-form-item prop="password">
          <el-input v-model="addForm.password" type="password" placeholder="密码" show-password />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="addVisible = false">取消</el-button>
        <el-button type="primary" :loading="addLoading" @click="handleAdd">添加</el-button>
      </template>
    </el-dialog>
  </el-dialog>
</template>

<script setup>
import { ref, watch } from 'vue'
import { getUsers, createUser, deleteUser } from '@/api'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Plus } from '@element-plus/icons-vue'

const props = defineProps({ visible: Boolean })
const emit = defineEmits(['update:visible'])

const users = ref([])
const loading = ref(false)
const addVisible = ref(false)
const addLoading = ref(false)
const addFormRef = ref(null)
const addForm = ref({ username: '', password: '' })
const addRules = {
  username: [{ required: true, message: '请输入用户名', trigger: 'blur' }],
  password: [{ required: true, min: 6, message: '密码至少6位', trigger: 'blur' }]
}

watch(() => props.visible, (v) => { if (v) fetchUsers() })

async function fetchUsers() {
  loading.value = true
  try { users.value = await getUsers() } catch (e) {} finally { loading.value = false }
}

async function removeUser(row) {
  try {
    await ElMessageBox.confirm(`确定删除用户「${row.username}」？`, '提示', { type: 'warning' })
    await deleteUser(row.id)
    ElMessage.success('已删除')
    fetchUsers()
  } catch (e) {}
}

async function handleAdd() {
  const valid = await addFormRef.value.validate().catch(() => false)
  if (!valid) return
  addLoading.value = true
  try {
    await createUser(addForm.value.username, addForm.value.password, '', 'guest')
    ElMessage.success('已添加')
    addVisible.value = false
    addForm.value = { username: '', password: '' }
    fetchUsers()
  } catch (e) {
    ElMessage.error(e.response?.data?.error || '添加失败')
  } finally {
    addLoading.value = false
  }
}
</script>

<style scoped>
.user-mgr-add {
  margin-top: 12px;
  text-align: center;
}
</style>
