<template>
  <el-dialog
    v-model="dialogVisible"
    title="设置"
    width="720px"
    :close-on-click-modal="false"
    append-to-body
    @closed="onClosed"
  >
    <el-tabs v-model="activeTab" tab-position="left" class="settings-tabs">
      <!-- ===== 基本设置 ===== -->
      <el-tab-pane label="基本设置" name="basic">
        <div class="tab-content">
          <h3 class="tab-section-title">用户信息</h3>
          <el-form label-width="80px">
            <el-form-item label="用户名">
              <div class="user-info-row">
                <el-input v-if="editingUsername" v-model="newUsername" style="width:200px" />
                <span v-else>{{ auth.user?.username }}</span>
                <el-button v-if="!editingUsername" text @click="startEditUsername">修改</el-button>
                <template v-else>
                  <el-button type="primary" size="small" :loading="savingUsername" @click="saveUsername">保存</el-button>
                  <el-button size="small" @click="editingUsername = false">取消</el-button>
                </template>
              </div>
            </el-form-item>
            <el-form-item label="角色">
              <el-tag :type="auth.isAdmin ? 'danger' : 'info'" size="small">
                {{ auth.isAdmin ? '管理员' : '访客' }}
              </el-tag>
            </el-form-item>
            <el-form-item label="密码">
              <el-button text @click="showChangePwd = true">修改密码</el-button>
            </el-form-item>
          </el-form>

          <el-divider />

          <h3 class="tab-section-title">页面标题</h3>
          <el-input v-model="local.pageTitle" placeholder="个人导航页" style="max-width:300px" />

          <h3 class="tab-section-title" style="margin-top:20px">Favicon</h3>
          <div class="bg-icon-picker-wrap">
            <IconPicker v-model="local.pageFavicon" category="icon" default-icon="🌐" hide-emoji />
          </div>

          <el-divider />

          <h3 class="tab-section-title">自定义欢迎语</h3>
          <p style="color:#86868b;font-size:13px;margin:0 0 12px">
            显示在主页日期下方，留空则不显示。
          </p>
          <el-input
            v-model="local.welcomeMessage"
            placeholder="输入一句欢迎语，如：今天也要加油呀~"
            maxlength="80"
            show-word-limit
            style="max-width:360px"
          />

          <el-divider />

          <h3 class="tab-section-title">页脚内容</h3>
          <p style="color:#86868b;font-size:13px;margin:0 0 12px">
            显示在页面底部，支持纯文本。
          </p>
          <el-input
            v-model="local.footer"
            type="textarea"
            :rows="3"
            placeholder="输入页脚文字，留空则不显示"
          />
          <div v-if="local.footer" class="footer-preview">
            <span class="footer-preview-label">预览：</span>
            <span style="white-space:pre-line">{{ local.footer }}</span>
          </div>
        </div>
      </el-tab-pane>

      <!-- ===== 背景设置 ===== -->
      <el-tab-pane label="背景设置" name="background">
        <div class="tab-content">
          <div class="bg-preview-wrap" :class="{ 'bg-preview-empty': !local.bgImage }">
            <template v-if="local.bgImage">
              <div
                v-if="local.bgImage.startsWith('linear-gradient')"
                class="bg-preview-inner"
                :style="{ background: local.bgImage }"
              />
              <img v-else :src="local.bgImage" class="bg-preview-img" />
            </template>
            <span v-else class="bg-preview-empty-text">当前无背景</span>
          </div>

          <el-divider />

          <h3 class="tab-section-title">预设背景</h3>
          <div class="preset-grid">
            <div
              v-for="bg in presets"
              :key="bg.value"
              class="preset-item"
              :class="{ active: local.bgImage === bg.value }"
              @click="local.bgImage = bg.value"
            >
              <img v-if="bg.value" :src="bg.value" class="preset-thumb" />
              <div class="preset-label">{{ bg.label }}</div>
              <el-icon v-if="local.bgImage === bg.value" class="preset-check"><Check /></el-icon>
            </div>
          </div>

          <el-divider />

          <h3 class="tab-section-title">自定义背景</h3>
          <div class="bg-icon-picker-wrap">
            <IconPicker v-model="local.bgImage" category="bg" default-icon="🖼️" />
          </div>
        </div>
      </el-tab-pane>

      <!-- ===== 图库管理 ===== -->
      <el-tab-pane label="图库管理" name="images">
        <div class="tab-content">
          <IconPicker mode="manage" />
        </div>
      </el-tab-pane>

      <!-- ===== 用户管理（管理员） ===== -->
      <el-tab-pane v-if="auth.isAdmin" label="用户管理" name="users">
        <div class="tab-content">
          <div v-loading="userLoading">
            <el-table :data="users" size="small" style="width:100%">
              <el-table-column prop="username" label="用户名" />
              <el-table-column prop="role" label="角色" width="80">
                <template #default="{ row }">
                  <el-tag :type="row.role === 'admin' ? 'danger' : 'info'" size="small">
                    {{ row.role === 'admin' ? '管理员' : '访客' }}
                  </el-tag>
                </template>
              </el-table-column>
              <el-table-column label="操作" width="80" align="center">
                <template #default="{ row }">
                  <el-button
                    size="small"
                    type="danger"
                    text
                    :disabled="row.id === auth.user?.id"
                    @click="removeUser(row)"
                  >
                    删除
                  </el-button>
                </template>
              </el-table-column>
            </el-table>
            <div class="user-add-row">
              <el-button :icon="Plus" text @click="showAddUser = true">添加用户</el-button>
            </div>
          </div>
        </div>
      </el-tab-pane>
    </el-tabs>

    <template #footer>
      <el-button @click="dialogVisible = false">取消</el-button>
      <el-button type="primary" :loading="saving" @click="handleSave">保存设置</el-button>
    </template>

    <el-dialog title="修改密码" v-model="showChangePwd" width="380px" append-to-body>
      <el-form ref="pwdFormRef" :model="pwdForm" :rules="pwdRules" label-width="0">
        <el-form-item prop="oldPassword">
          <el-input v-model="pwdForm.oldPassword" type="password" placeholder="旧密码" show-password />
        </el-form-item>
        <el-form-item prop="newPassword">
          <el-input v-model="pwdForm.newPassword" type="password" placeholder="新密码（至少6位）" show-password />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showChangePwd = false">取消</el-button>
        <el-button type="primary" :loading="pwdLoading" @click="handleChangePwd">确认</el-button>
      </template>
    </el-dialog>

    <el-dialog title="添加用户" v-model="showAddUser" width="380px" append-to-body>
      <el-form ref="addFormRef" :model="addForm" :rules="addRules" label-width="0">
        <el-form-item prop="username">
          <el-input v-model="addForm.username" placeholder="用户名" />
        </el-form-item>
        <el-form-item prop="password">
          <el-input v-model="addForm.password" type="password" placeholder="密码" show-password />
        </el-form-item>
        <el-form-item label="角色">
          <el-radio-group v-model="addForm.role">
            <el-radio value="guest">访客</el-radio>
            <el-radio value="admin">管理员</el-radio>
          </el-radio-group>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showAddUser = false">取消</el-button>
        <el-button type="primary" :loading="addUserLoading" @click="handleAddUser">添加</el-button>
      </template>
    </el-dialog>
  </el-dialog>
</template>

<script setup>
import { ref, reactive, watch, computed, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Plus, Check } from '@element-plus/icons-vue'
import { useAuthStore } from '@/stores/auth'
import { useSettingsStore } from '@/stores/settings'
import { updateMe, changePassword } from '@/api'
import { getUsers, createUser, deleteUser } from '@/api'
import IconPicker from './IconPicker.vue'

const props = defineProps({ visible: Boolean })
const emit = defineEmits(['update:visible'])

const auth = useAuthStore()
const settings = useSettingsStore()

const dialogVisible = computed({
  get: () => props.visible,
  set: (v) => emit('update:visible', v)
})

const activeTab = ref('basic')
const saving = ref(false)

const local = reactive({
  pageTitle: '',
  pageFavicon: '',
  bgImage: '',
  footer: '',
  welcomeMessage: ''
})

watch(() => props.visible, (v) => {
  if (v) {
    local.pageTitle = settings.pageTitle
    local.pageFavicon = settings.pageFavicon
    local.bgImage = settings.bgImage
    local.footer = settings.footer
    local.welcomeMessage = settings.welcomeMessage
  }
})

function onClosed() {
  settings.fetch(auth.userId)
}

const presets = [
  { label: '无', value: '' },
  { label: '山脉晨曦', value: 'https://images.unsplash.com/photo-1506905925346-21bda4d32df4?w=1920&q=80' },
  { label: '海洋日落', value: 'https://images.unsplash.com/photo-1507525428034-b723cf961d3e?w=1920&q=80' },
  { label: '森林小径', value: 'https://images.unsplash.com/photo-1441974231531-c6227db76b6e?w=1920&q=80' },
  { label: '星空银河', value: 'https://images.unsplash.com/photo-1469474968028-56623f02e42e?w=1920&q=80' },
  { label: '极光夜空', value: 'https://images.unsplash.com/photo-1531366936337-7c912a4589a7?w=1920&q=80' },
  { label: '雪山倒影', value: 'https://images.unsplash.com/photo-1470071459604-3b5ec3a7fe05?w=1920&q=80' },
  { label: '沙漠星空', value: 'https://images.unsplash.com/photo-1509316785289-025f5b846b35?w=1920&q=80' },
  { label: '城市夜景', value: 'https://images.unsplash.com/photo-1519501025264-65ba15a82390?w=1920&q=80' },
  { label: '樱花大道', value: 'https://images.unsplash.com/photo-1522383225653-ed111181a951?w=1920&q=80' },
  { label: '秋叶湖畔', value: 'https://images.unsplash.com/photo-1501785888041-af3ef285b470?w=1920&q=80' },
  { label: '云雾山谷', value: 'https://images.unsplash.com/photo-1464822759023-fed622ff2c3b?w=1920&q=80' },
]

const editingUsername = ref(false)
const newUsername = ref('')
const savingUsername = ref(false)

function startEditUsername() {
  newUsername.value = auth.user?.username || ''
  editingUsername.value = true
}

async function saveUsername() {
  if (!newUsername.value.trim()) {
    ElMessage.warning('用户名不能为空')
    return
  }
  savingUsername.value = true
  try {
    const res = await updateMe(newUsername.value.trim())
    if (res.token) {
      auth.token = res.token
      localStorage.setItem('token', res.token)
    }
    auth.user.username = newUsername.value.trim()
    localStorage.setItem('user', JSON.stringify(auth.user))
    ElMessage.success('用户名已更新')
    editingUsername.value = false
  } finally {
    savingUsername.value = false
  }
}

const showChangePwd = ref(false)
const pwdLoading = ref(false)
const pwdFormRef = ref(null)
const pwdForm = reactive({ oldPassword: '', newPassword: '' })
const pwdRules = {
  oldPassword: [{ required: true, message: '请输入旧密码', trigger: 'blur' }],
  newPassword: [{ required: true, min: 6, message: '新密码至少6位', trigger: 'blur' }]
}

async function handleChangePwd() {
  const valid = await pwdFormRef.value.validate().catch(() => false)
  if (!valid) return
  pwdLoading.value = true
  try {
    await changePassword(pwdForm.oldPassword, pwdForm.newPassword)
    ElMessage.success('密码修改成功')
    showChangePwd.value = false
    pwdForm.oldPassword = ''
    pwdForm.newPassword = ''
  } finally {
    pwdLoading.value = false
  }
}

const users = ref([])
const userLoading = ref(false)
const showAddUser = ref(false)
const addUserLoading = ref(false)
const addFormRef = ref(null)
const addForm = reactive({ username: '', password: '', role: 'guest' })
const addRules = {
  username: [{ required: true, message: '请输入用户名', trigger: 'blur' }],
  password: [{ required: true, min: 6, message: '密码至少6位', trigger: 'blur' }]
}

watch(activeTab, (tab) => {
  if (tab === 'users' && auth.isAdmin) fetchUsers()
})

async function fetchUsers() {
  userLoading.value = true
  try { users.value = await getUsers() } catch {} finally { userLoading.value = false }
}

async function removeUser(row) {
  try {
    await ElMessageBox.confirm(`确定删除用户「${row.username}」？`, '提示', { type: 'warning' })
    await deleteUser(row.id)
    ElMessage.success('已删除')
    fetchUsers()
  } catch {}
}

async function handleAddUser() {
  const valid = await addFormRef.value.validate().catch(() => false)
  if (!valid) return
  addUserLoading.value = true
  try {
    await createUser(addForm.username, addForm.password, '', addForm.role)
    ElMessage.success('已添加')
    showAddUser.value = false
    addForm.username = ''
    addForm.password = ''
    addForm.role = 'guest'
    fetchUsers()
  } finally {
    addUserLoading.value = false
  }
}

async function handleSave() {
  saving.value = true
  try {
    await settings.save(auth.userId, {
      page_title: local.pageTitle,
      page_favicon: local.pageFavicon,
      bg_image: local.bgImage,
      footer: local.footer,
      welcome_message: local.welcomeMessage
    })
    ElMessage.success('已保存')
    dialogVisible.value = false
  } finally {
    saving.value = false
  }
}

onMounted(() => {
  if (auth.isAdmin) fetchUsers()
})
</script>

<style scoped>
.settings-tabs {
  min-height: 420px;
}
.tab-content {
  padding: 4px;
  height: 400px;
  overflow-y: auto;
  padding: 4px;
  border: 1px solid #e8e8ec;
  border-radius: var(--radius-md);
}
.tab-section-title {
  font-size: 15px;
  font-weight: 600;
  color: #1d1d1f;
  margin: 0 0 12px;
}

.user-info-row {
  display: flex;
  align-items: center;
  gap: 8px;
}

.preset-grid {
  display: grid;
  grid-template-columns: repeat(6, 1fr);
  gap: 8px;
}
.preset-item {
  height: 52px;
  border-radius: 8px;
  cursor: pointer;
  position: relative;
  border: 1px solid #e8e8ec;
  transition: border-color 0.2s, transform 0.2s;
  overflow: hidden;
  background: #f5f5f7;
}
.preset-thumb {
  width: 100%;
  height: 100%;
  object-fit: cover;
  display: block;
}
.preset-item:not(:has(.preset-thumb)) {
  display: flex;
  align-items: center;
  justify-content: center;
}
.preset-item:hover {
  transform: scale(1.03);
  border-color: var(--color-primary);
}
.preset-item.active {
  border-color: var(--color-primary);
  border-width: 2px;
}
.preset-label {
  position: absolute;
  bottom: 0;
  left: 0;
  right: 0;
  font-size: 10px;
  color: #fff;
  text-shadow: 0 1px 3px rgba(0,0,0,0.5);
  background: linear-gradient(transparent, rgba(0,0,0,0.5));
  padding: 16px 6px 4px;
  text-align: center;
  line-height: 1.3;
}
.preset-item:not(:has(.preset-thumb)) .preset-label {
  position: static;
  background: none;
  text-shadow: none;
  color: #86868b;
  padding: 0;
}
.preset-check {
  position: absolute;
  top: 4px;
  right: 4px;
  color: #fff;
  background: #6366f1;
  border-radius: 50%;
  font-size: 12px;
  padding: 3px;
}

.bg-icon-picker-wrap {
  padding: 0px;
  max-height: 400px;
  overflow: auto;
}

.bg-preview-wrap {
  border-radius: 12px;
  border: 1px solid #e8e8ec;
  overflow: hidden;
  height: 80px;
  display: flex;
  align-items: center;
  justify-content: center;
}
.bg-preview-empty {
  background: #f5f5f7;
  border: 2px dashed #d4d4d8;
}
.bg-preview-empty-text {
  color: #aeaeb2;
  font-size: 14px;
}
.bg-preview-inner {
  width: 100%;
  height: 100%;
}
.bg-preview-img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.footer-preview {
  margin-top: 12px;
  padding: 12px;
  background: #f5f5f7;
  border: 1px solid #e8e8ec;
  border-radius: 8px;
  font-size: 13px;
  color: #86868b;
}
.footer-preview-label {
  color: #aeaeb2;
}

.user-add-row {
  margin-top: 12px;
  text-align: center;
}
</style>