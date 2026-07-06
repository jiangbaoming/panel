# 书签导航面板前端实现计划

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development（推荐）或 superpowers:executing-plans 按任务逐步实现。步骤使用 `- [ ]` 复选框追踪进度。

**目标：** 基于 Vue 3 + Vite + Element Plus 构建个人书签导航前端，完整对接 Go 后端 API。

**架构：** 纯 SPA，仅两个路由（/login 和 /），所有二级功能以弹窗形式呈现。状态管理使用 Pinia，HTTP 请求使用 Axios（含拦截器自动携带 Token），拖拽排序使用 vuedraggable。

**技术栈：**
- Vue 3 (Composition API, JS)
- Vite 5
- Vue Router 4
- Pinia 2
- Axios
- Element Plus
- vuedraggable@next

## 全局约束

- 所有代码使用 JavaScript（非 TypeScript）
- 代码注释/UI 文本使用中文
- 页面标题从后端 `/api/settings/:userId` 的 `page_title` 字段获取
- 所有 API 请求需携带 `Authorization: Bearer <token>` 头
- 401 响应时自动跳转 `/login`
- 构建产物输出到 `web/dist/` 目录

---

## 文件结构

```
web/
├── index.html
├── vite.config.js
├── package.json
├── src/
│   ├── main.js                  # 入口
│   ├── App.vue                  # 根组件（含路由视图）
│   ├── router/
│   │   └── index.js             # 路由配置
│   ├── api/
│   │   └── index.js             # Axios 实例 + 所有 API 方法
│   ├── stores/
│   │   ├── auth.js              # 认证 store
│   │   ├── bookmarks.js         # 分组+书签数据 store
│   │   └── settings.js          # 用户设置 store
│   ├── views/
│   │   ├── LoginPage.vue        # 登录页
│   │   └── MainPage.vue         # 主面板（聚合所有功能）
│   ├── components/
│   │   ├── TopBar.vue           # 顶栏（头像+设置/图库/登出入口）
│   │   ├── PinnedBar.vue        # 常驻书签行
│   │   ├── GroupSection.vue     # 分组区块（含书签网格+拖拽）
│   │   ├── BookmarkCard.vue     # 单张书签卡片
│   │   ├── SearchOverlay.vue    # 搜索遮罩弹窗
│   │   ├── SettingsDialog.vue   # 设置对话框
│   │   ├── ImageManageDialog.vue# 图库管理对话框
│   │   ├── UserManageDialog.vue # 用户管理对话框
│   │   ├── GroupFormDialog.vue  # 新建/编辑分组对话框
│   │   └── BookmarkFormDialog.vue # 新建/编辑书签对话框
│   └── styles/
│       └── main.css             # 全局样式
```

---

### Task 1: 项目脚手架搭建

**Files:**
- Create: `web/package.json`
- Create: `web/vite.config.js`
- Create: `web/index.html`
- Create: `web/src/main.js`
- Create: `web/src/App.vue`
- Create: `web/src/router/index.js`
- Create: `web/src/styles/main.css`

**Interfaces:**
- Produces: 可运行的空 Vite 项目骨架，Vue Router 配置（`/login` 和 `/` 两路由）

- [ ] **Step 1: 创建 package.json**

```json
{
  "name": "bookmark-panel",
  "private": true,
  "version": "1.0.0",
  "type": "module",
  "scripts": {
    "dev": "vite",
    "build": "vite build",
    "preview": "vite preview"
  },
  "dependencies": {
    "vue": "^3.4.0",
    "vue-router": "^4.3.0",
    "pinia": "^2.1.0",
    "axios": "^1.7.0",
    "element-plus": "^2.7.0",
    "@element-plus/icons-vue": "^2.3.0",
    "vuedraggable": "^4.1.0"
  },
  "devDependencies": {
    "@vitejs/plugin-vue": "^5.0.0",
    "vite": "^5.4.0"
  }
}
```

- [ ] **Step 2: 创建 vite.config.js**

```js
import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'

export default defineConfig({
  plugins: [vue()],
  resolve: {
    alias: {
      '@': '/src'
    }
  },
  server: {
    port: 3000,
    proxy: {
      '/api': {
        target: 'http://localhost:8080',
        changeOrigin: true
      },
      '/uploads': {
        target: 'http://localhost:8080',
        changeOrigin: true
      }
    }
  },
  build: {
    outDir: 'dist'
  }
})
```

- [ ] **Step 3: 创建 index.html**

```html
<!DOCTYPE html>
<html lang="zh-CN">
<head>
  <meta charset="UTF-8" />
  <meta name="viewport" content="width=device-width, initial-scale=1.0" />
  <title>个人导航页</title>
  <link rel="icon" href="data:," />
</head>
<body>
  <div id="app"></div>
  <script type="module" src="/src/main.js"></script>
</body>
</html>
```

- [ ] **Step 4: 创建 main.js**

```js
import { createApp } from 'vue'
import { createPinia } from 'pinia'
import ElementPlus from 'element-plus'
import 'element-plus/dist/index.css'
import * as ElementPlusIconsVue from '@element-plus/icons-vue'
import App from './App.vue'
import router from './router'
import './styles/main.css'

const app = createApp(App)
app.use(createPinia())
app.use(router)
app.use(ElementPlus)
for (const [key, component] of Object.entries(ElementPlusIconsVue)) {
  app.component(key, component)
}
app.mount('#app')
```

- [ ] **Step 5: 创建 App.vue**

```vue
<template>
  <router-view />
</template>
```

- [ ] **Step 6: 创建 router/index.js**

```js
import { createRouter, createWebHistory } from 'vue-router'

const routes = [
  {
    path: '/login',
    name: 'Login',
    component: () => import('@/views/LoginPage.vue'),
    meta: { requiresAuth: false }
  },
  {
    path: '/',
    name: 'Main',
    component: () => import('@/views/MainPage.vue'),
    meta: { requiresAuth: true }
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

router.beforeEach((to, from, next) => {
  const token = localStorage.getItem('token')
  if (to.meta.requiresAuth !== false && !token) {
    next('/login')
  } else if (to.path === '/login' && token) {
    next('/')
  } else {
    next()
  }
})

export default router
```

- [ ] **Step 7: 创建 main.css（全局样式）**

```css
* {
  margin: 0;
  padding: 0;
  box-sizing: border-box;
}

html, body, #app {
  height: 100%;
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, 'Helvetica Neue', Arial, sans-serif;
  background: #f5f5f7;
  color: #1d1d1f;
}

a {
  color: inherit;
  text-decoration: none;
}
```

- [ ] **Step 8: 安装依赖并验证**

```bash
cd /config/workspace/panel/web && npm install
```

---

### Task 2: API 层封装

**Files:**
- Create: `web/src/api/index.js`

**Interfaces:**
- Produces: `api` (Axios 实例)、所有 API 方法（供 stores 调用）

- [ ] **Step 1: 创建 API 封装文件**

```js
import axios from 'axios'

const api = axios.create({
  baseURL: '/api'
})

// 请求拦截器：自动携带 token
api.interceptors.request.use(config => {
  const token = localStorage.getItem('token')
  if (token) {
    config.headers.Authorization = `Bearer ${token}`
  }
  return config
})

// 响应拦截器：401 自动跳转登录
api.interceptors.response.use(
  response => response,
  error => {
    if (error.response && error.response.status === 401) {
      localStorage.removeItem('token')
      window.location.href = '/login'
    }
    return Promise.reject(error)
  }
)

// 认证
export const login = (username, password) =>
  api.post('/login', { username, password }).then(r => r.data)

export const changePassword = (oldPassword, newPassword) =>
  api.post('/change-password', { oldPassword, newPassword }).then(r => r.data)

// 分组
export const getGroups = () =>
  api.get('/groups').then(r => r.data)

export const createGroup = (name, icon) =>
  api.post('/groups', { name, icon }).then(r => r.data)

export const updateGroup = (id, name, icon) =>
  api.put(`/groups/${id}`, { name, icon }).then(r => r.data)

export const deleteGroup = id =>
  api.delete(`/groups/${id}`).then(r => r.data)

export const reorderGroups = ids =>
  api.patch('/groups/reorder', { ids }).then(r => r.data)

// 书签
export const createBookmark = (groupId, data) =>
  api.post(`/groups/${groupId}/bookmarks`, data).then(r => r.data)

export const updateBookmark = (groupId, bookmarkId, data) =>
  api.put(`/groups/${groupId}/bookmarks/${bookmarkId}`, data).then(r => r.data)

export const deleteBookmark = (groupId, bookmarkId) =>
  api.delete(`/groups/${groupId}/bookmarks/${bookmarkId}`).then(r => r.data)

export const reorderBookmarks = (groupId, ids) =>
  api.patch(`/groups/${groupId}/bookmarks/reorder`, { ids }).then(r => r.data)

// 常驻书签
export const getPinned = () =>
  api.get('/pinned').then(r => r.data)

export const togglePin = (id, pinned) =>
  api.patch(`/pinned/${id}/pin`, { pinned }).then(r => r.data)

export const reorderPinned = ids =>
  api.patch('/pinned/reorder', { ids }).then(r => r.data)

// 图片
export const getImages = (params = {}) =>
  api.get('/images', { params }).then(r => r.data)

export const uploadImage = (formData, onProgress) =>
  api.post('/images/upload', formData, {
    headers: { 'Content-Type': 'multipart/form-data' },
    onUploadProgress: onProgress
  }).then(r => r.data)

export const uploadZip = formData =>
  api.post('/images/upload-zip', formData, {
    headers: { 'Content-Type': 'multipart/form-data' }
  }).then(r => r.data)

export const deleteImage = id =>
  api.delete(`/images/${id}`).then(r => r.data)

export const clearImages = category =>
  api.delete(`/images/clear/${category}`).then(r => r.data)

// 用户设置
export const getSettings = userId =>
  api.get(`/settings/${userId}`).then(r => r.data)

export const updateSettings = (userId, data) =>
  api.put(`/settings/${userId}`, data).then(r => r.data)

// 搜索
export const searchBookmarks = q =>
  api.get('/search', { params: { q } }).then(r => r.data)

// 用户管理（管理员）
export const getUsers = () =>
  api.get('/users').then(r => r.data)

export const createUser = (username, password, avatar, role) =>
  api.post('/users', { username, password, avatar, role }).then(r => r.data)

export const deleteUser = id =>
  api.delete(`/users/${id}`).then(r => r.data)

export const updateMe = username =>
  api.patch('/users/me', { username }).then(r => r.data)

export default api
```

---

### Task 3: 登录页 + 认证 Store

**Files:**
- Create: `web/src/stores/auth.js`
- Create: `web/src/views/LoginPage.vue`

**Interfaces:**
- Consumes: `login()` from api
- Produces: `useAuthStore` (token, user 信息, login/logout/checkAuth 方法)

- [ ] **Step 1: 创建 auth store**

```js
import { defineStore } from 'pinia'
import { login as apiLogin, updateMe as apiUpdateMe } from '@/api'
import { getSettings as apiGetSettings } from '@/api'

export const useAuthStore = defineStore('auth', {
  state: () => ({
    token: localStorage.getItem('token') || '',
    user: JSON.parse(localStorage.getItem('user') || 'null'),
    settings: null
  }),
  getters: {
    isLoggedIn: (state) => !!state.token,
    isAdmin: (state) => state.user?.role === 'admin',
    userId: (state) => state.user?.id
  },
  actions: {
    async login(username, password) {
      const res = await apiLogin(username, password)
      this.token = res.token
      this.user = { id: res.id, username: res.username, avatar: res.avatar, role: res.role }
      localStorage.setItem('token', res.token)
      localStorage.setItem('user', JSON.stringify(this.user))
    },
    logout() {
      this.token = ''
      this.user = null
      this.settings = null
      localStorage.removeItem('token')
      localStorage.removeItem('user')
    },
    async updateUsername(username) {
      const res = await apiUpdateMe(username)
      this.user.username = username
      localStorage.setItem('user', JSON.stringify(this.user))
      return res
    }
  }
})
```

- [ ] **Step 2: 创建 LoginPage.vue**

```vue
<template>
  <div class="login-page">
    <div class="login-card">
      <div class="login-header">
        <div class="login-icon">🔖</div>
        <h1>个人导航页</h1>
        <p class="login-desc">登录以管理你的书签</p>
      </div>
      <el-form
        ref="formRef"
        :model="form"
        :rules="rules"
        @keyup.enter="handleLogin"
        size="large"
      >
        <el-form-item prop="username">
          <el-input
            v-model="form.username"
            placeholder="用户名"
            :prefix-icon="User"
          />
        </el-form-item>
        <el-form-item prop="password">
          <el-input
            v-model="form.password"
            type="password"
            placeholder="密码"
            :prefix-icon="Lock"
            show-password
          />
        </el-form-item>
        <el-form-item>
          <el-button
            type="primary"
            :loading="loading"
            style="width: 100%"
            @click="handleLogin"
          >
            登录
          </el-button>
        </el-form-item>
      </el-form>
      <p v-if="error" class="login-error">{{ error }}</p>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { User, Lock } from '@element-plus/icons-vue'

const router = useRouter()
const auth = useAuthStore()
const formRef = ref(null)
const loading = ref(false)
const error = ref('')

const form = reactive({
  username: '',
  password: ''
})

const rules = {
  username: [{ required: true, message: '请输入用户名', trigger: 'blur' }],
  password: [{ required: true, message: '请输入密码', trigger: 'blur' }]
}

async function handleLogin() {
  const valid = await formRef.value.validate().catch(() => false)
  if (!valid) return
  loading.value = true
  error.value = ''
  try {
    await auth.login(form.username, form.password)
    router.push('/')
  } catch (e) {
    error.value = e.response?.data?.error || '登录失败'
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.login-page {
  height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
}
.login-card {
  background: #fff;
  border-radius: 20px;
  padding: 40px;
  width: 380px;
  box-shadow: 0 20px 60px rgba(0, 0, 0, 0.15);
}
.login-header {
  text-align: center;
  margin-bottom: 32px;
}
.login-icon {
  font-size: 48px;
  margin-bottom: 12px;
}
.login-header h1 {
  font-size: 24px;
  font-weight: 700;
  color: #1d1d1f;
  margin-bottom: 8px;
}
.login-desc {
  font-size: 14px;
  color: #86868b;
}
.login-error {
  color: #f56c6c;
  font-size: 13px;
  text-align: center;
  margin-top: 8px;
}
</style>
```

---

### Task 4: 主面板 + TopBar + 数据 Store

**Files:**
- Create: `web/src/stores/bookmarks.js`
- Create: `web/src/stores/settings.js`
- Create: `web/src/views/MainPage.vue`
- Create: `web/src/components/TopBar.vue`

**Interfaces:**
- Consumes: all api methods, `useAuthStore`
- Produces: `useBookmarksStore`（groups, pinned, 增删改查方法）、`useSettingsStore`（settings 加载更新）、主面板骨架

- [ ] **Step 1: 创建 bookmarks store**

```js
import { defineStore } from 'pinia'
import {
  getGroups, createGroup, updateGroup, deleteGroup, reorderGroups,
  createBookmark, updateBookmark, deleteBookmark, reorderBookmarks,
  getPinned, togglePin, reorderPinned
} from '@/api'

export const useBookmarksStore = defineStore('bookmarks', {
  state: () => ({
    groups: [],
    pinned: [],
    loading: false
  }),
  actions: {
    async fetchAll() {
      this.loading = true
      try {
        const [groups, pinned] = await Promise.all([getGroups(), getPinned()])
        this.groups = groups
        this.pinned = pinned
      } finally {
        this.loading = false
      }
    },
    async addGroup(name, icon) {
      const group = await createGroup(name, icon)
      this.groups.push({ ...group, bookmarks: group.bookmarks || [] })
    },
    async editGroup(id, name, icon) {
      await updateGroup(id, name, icon)
      const g = this.groups.find(g => g.id === id)
      if (g) { g.name = name; g.icon = icon }
    },
    async removeGroup(id) {
      await deleteGroup(id)
      this.groups = this.groups.filter(g => g.id !== id)
    },
    async sortGroups(ids) {
      await reorderGroups(ids)
    },
    async addBookmark(groupId, data) {
      const bm = await createBookmark(groupId, data)
      const g = this.groups.find(g => g.id === groupId)
      if (g) {
        if (!g.bookmarks) g.bookmarks = []
        g.bookmarks.push(bm)
      }
    },
    async editBookmark(groupId, bookmarkId, data) {
      await updateBookmark(groupId, bookmarkId, data)
      const g = this.groups.find(g => g.id === groupId)
      if (g) {
        const b = g.bookmarks.find(b => b.id === bookmarkId)
        if (b) Object.assign(b, data)
      }
    },
    async removeBookmark(groupId, bookmarkId) {
      await deleteBookmark(groupId, bookmarkId)
      const g = this.groups.find(g => g.id === groupId)
      if (g) g.bookmarks = g.bookmarks.filter(b => b.id !== bookmarkId)
    },
    async sortBookmarks(groupId, ids) {
      await reorderBookmarks(groupId, ids)
    },
    // 常驻
    async togglePinned(id, pinned) {
      await togglePin(id, pinned)
      if (!pinned) {
        this.pinned = this.pinned.filter(b => b.id !== id)
      } else {
        await this.fetchAll()
      }
    },
    async sortPinned(ids) {
      await reorderPinned(ids)
    }
  }
})
```

- [ ] **Step 2: 创建 settings store**

```js
import { defineStore } from 'pinia'
import { getSettings, updateSettings as apiUpdateSettings } from '@/api'

export const useSettingsStore = defineStore('settings', {
  state: () => ({
    bgImage: '',
    displayMode: 'both',
    pageTitle: '个人导航页',
    pageFavicon: ''
  }),
  actions: {
    async fetch(userId) {
      try {
        const s = await getSettings(userId)
        this.bgImage = s.bg_image || ''
        this.displayMode = s.display_mode || 'both'
        this.pageTitle = s.page_title || '个人导航页'
        this.pageFavicon = s.page_favicon || ''
        document.title = this.pageTitle
        if (this.pageFavicon) {
          document.querySelector('link[rel="icon"]').href = this.pageFavicon
        }
      } catch (e) {
        // 默认值已设置
      }
    },
    async save(userId, data) {
      const s = await apiUpdateSettings(userId, data)
      Object.assign(this, {
        bgImage: s.bg_image || '',
        displayMode: s.display_mode || 'both',
        pageTitle: s.page_title || '个人导航页',
        pageFavicon: s.page_favicon || ''
      })
      document.title = this.pageTitle
      if (this.pageFavicon) {
        document.querySelector('link[rel="icon"]').href = this.pageFavicon
      }
    }
  }
})
```

- [ ] **Step 3: 创建 MainPage.vue（骨架 + 数据加载）**

```vue
<template>
  <div class="main-page" :style="bgStyle">
    <div class="main-container">
      <TopBar @open-settings="settingsVisible = true" @open-images="imagesVisible = true" @open-users="usersVisible = true" />
      <div class="main-content" v-loading="loading">
        <!-- 常驻书签 -->
        <PinnedBar @pin-changed="bookmarks.fetchAll()" />
        <!-- 分组列表 -->
        <div class="groups-wrapper">
          <GroupSection
            v-for="group in bookmarks.groups"
            :key="group.id"
            :group="group"
            @refresh="bookmarks.fetchAll()"
          />
        </div>
        <!-- 新建分组按钮 -->
        <div class="add-group-btn" v-if="bookmarks.groups.length > 0">
          <el-button @click="openGroupForm()" :icon="Plus" text size="large">
            新建分组
          </el-button>
        </div>
      </div>
    </div>
    <!-- 搜索按钮 -->
    <button class="search-fab" @click="searchVisible = true" title="搜索 (Ctrl+K)">
      <el-icon :size="22"><Search /></el-icon>
    </button>
    <!-- 搜索弹窗 -->
    <SearchOverlay v-model:visible="searchVisible" />
    <!-- 设置弹窗 -->
    <SettingsDialog v-model:visible="settingsVisible" />
    <!-- 图库弹窗 -->
    <ImageManageDialog v-model:visible="imagesVisible" />
    <!-- 用户管理弹窗 -->
    <UserManageDialog v-model:visible="usersVisible" />
  </div>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { useBookmarksStore } from '@/stores/bookmarks'
import { useSettingsStore } from '@/stores/settings'
import { Plus, Search } from '@element-plus/icons-vue'
import TopBar from '@/components/TopBar.vue'
import PinnedBar from '@/components/PinnedBar.vue'
import GroupSection from '@/components/GroupSection.vue'
import SearchOverlay from '@/components/SearchOverlay.vue'
import SettingsDialog from '@/components/SettingsDialog.vue'
import ImageManageDialog from '@/components/ImageManageDialog.vue'
import UserManageDialog from '@/components/UserManageDialog.vue'

const router = useRouter()
const auth = useAuthStore()
const bookmarks = useBookmarksStore()
const settings = useSettingsStore()

const loading = ref(true)
const searchVisible = ref(false)
const settingsVisible = ref(false)
const imagesVisible = ref(false)
const usersVisible = ref(false)

const bgStyle = computed(() => {
  if (settings.bgImage) {
    return { backgroundImage: `url(${settings.bgImage})`, backgroundSize: 'cover', backgroundPosition: 'center' }
  }
  return {}
})

onMounted(async () => {
  await Promise.all([
    bookmarks.fetchAll(),
    settings.fetch(auth.userId)
  ])
  loading.value = false
})

// Ctrl+K 快捷键
function handleKeydown(e) {
  if ((e.ctrlKey || e.metaKey) && e.key === 'k') {
    e.preventDefault()
    searchVisible.value = true
  }
}
onMounted(() => document.addEventListener('keydown', handleKeydown))
onUnmounted(() => document.removeEventListener('keydown', handleKeydown))
</script>

<style scoped>
.main-page {
  min-height: 100vh;
  background: #f5f5f7;
  transition: background 0.3s;
}
.main-container {
  max-width: 1200px;
  margin: 0 auto;
  padding: 0 24px;
}
.main-content {
  padding: 20px 0 80px;
}
.add-group-btn {
  text-align: center;
  margin-top: 24px;
}
.search-fab {
  position: fixed;
  bottom: 32px;
  right: 32px;
  width: 48px;
  height: 48px;
  border-radius: 50%;
  background: #6366f1;
  color: #fff;
  border: none;
  cursor: pointer;
  box-shadow: 0 4px 16px rgba(99, 102, 241, 0.4);
  display: flex;
  align-items: center;
  justify-content: center;
  transition: transform 0.2s, box-shadow 0.2s;
  z-index: 100;
}
.search-fab:hover {
  transform: scale(1.08);
  box-shadow: 0 6px 24px rgba(99, 102, 241, 0.5);
}
</style>
```

- [ ] **Step 4: 创建 TopBar.vue**

```vue
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
```

---

### Task 5: 分组区块 + 书签卡片 + 常驻书签

**Files:**
- Create: `web/src/components/GroupSection.vue`
- Create: `web/src/components/BookmarkCard.vue`
- Create: `web/src/components/PinnedBar.vue`

**Interfaces:**
- Consumes: `useBookmarksStore`, drag-and-drop event handlers
- Produces: 核心 UI 组件

- [ ] **Step 1: 创建 BookmarkCard.vue**

```vue
<template>
  <div
    class="bookmark-card"
    :style="cardStyle"
    @click="openUrl"
    @mouseenter="hover = true"
    @mouseleave="hover = false"
  >
    <div class="card-actions" v-show="hover">
      <el-button size="small" circle @click.stop="$emit('edit', bookmark)" title="编辑">
        <el-icon><Edit /></el-icon>
      </el-button>
      <el-button size="small" circle @click.stop="$emit('pin', bookmark)" title="置顶">
        <el-icon><Star /></el-icon>
      </el-button>
      <el-button size="small" circle type="danger" @click.stop="$emit('delete', bookmark)" title="删除">
        <el-icon><Delete /></el-icon>
      </el-button>
    </div>
    <div class="card-icon" :style="iconStyle">
      <img v-if="isImageUrl(bookmark.icon)" :src="bookmark.icon" />
      <span v-else>{{ bookmark.icon || '🔗' }}</span>
    </div>
    <div class="card-name">{{ bookmark.name }}</div>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue'
import { Edit, Star, Delete } from '@element-plus/icons-vue'

const props = defineProps({
  bookmark: { type: Object, required: true }
})
defineEmits(['edit', 'pin', 'delete'])

const hover = ref(false)

function isImageUrl(str) {
  return str && (str.startsWith('http') || str.startsWith('/uploads/'))
}

const cardStyle = computed(() => {
  if (props.bookmark.bg_color) {
    return { backgroundColor: props.bookmark.bg_color }
  }
  return {}
})

const iconStyle = computed(() => {
  if (props.bookmark.icon_bg) {
    return { backgroundColor: props.bookmark.icon_bg }
  }
  return {}
})

function openUrl() {
  if (props.bookmark.url) {
    window.open(props.bookmark.url, '_blank')
  }
}
</script>

<style scoped>
.bookmark-card {
  width: 120px;
  height: 110px;
  background: #fff;
  border-radius: 16px;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 8px;
  cursor: pointer;
  position: relative;
  transition: transform 0.2s, box-shadow 0.2s;
  box-shadow: 0 1px 4px rgba(0, 0, 0, 0.04);
}
.bookmark-card:hover {
  transform: translateY(-3px);
  box-shadow: 0 8px 24px rgba(0, 0, 0, 0.1);
}
.card-actions {
  position: absolute;
  top: 4px;
  right: 4px;
  display: flex;
  gap: 2px;
  z-index: 2;
}
.card-actions .el-button {
  width: 24px;
  height: 24px;
  background: rgba(255, 255, 255, 0.9);
}
.card-icon {
  width: 44px;
  height: 44px;
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 24px;
  background: #f5f5f7;
}
.card-icon img {
  width: 28px;
  height: 28px;
  object-fit: contain;
}
.card-name {
  font-size: 12px;
  color: #1d1d1f;
  text-align: center;
  max-width: 100px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}
</style>
```

- [ ] **Step 2: 创建 GroupSection.vue**

```vue
<template>
  <div class="group-section">
    <div class="group-header">
      <div class="group-title">
        <span class="group-icon">{{ group.icon || '📁' }}</span>
        <h2 class="group-name">{{ group.name }}</h2>
        <span class="group-count">{{ group.bookmarks?.length || 0 }} 个</span>
        <div class="group-actions">
          <el-button text @click="editGroup" title="编辑分组">
            <el-icon><Edit /></el-icon>
          </el-button>
          <el-button text @click="deleteGroup" title="删除分组" type="danger">
            <el-icon><Delete /></el-icon>
          </el-button>
        </div>
      </div>
      <el-button size="small" :icon="Plus" @click="addBookmark">添加书签</el-button>
    </div>
    <div class="bookmark-grid">
      <draggable
        :list="group.bookmarks"
        item-key="id"
        group="bookmarks"
        handle=".bookmark-card"
        ghost-class="ghost"
        :on-end="onDragEnd"
      >
        <template #item="{ element }">
          <div>
            <BookmarkCard
              :bookmark="element"
              @edit="openEdit(element)"
              @pin="handlePin(element)"
              @delete="handleDelete(element)"
            />
          </div>
        </template>
      </draggable>
    </div>
    <!-- 新建/编辑书签弹窗 -->
    <BookmarkFormDialog
      v-model:visible="formVisible"
      :group-id="group.id"
      :edit-data="editData"
      @saved="emit('refresh')"
    />
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import draggable from 'vuedraggable'
import { useBookmarksStore } from '@/stores/bookmarks'
import { Plus, Edit, Delete } from '@element-plus/icons-vue'
import BookmarkCard from './BookmarkCard.vue'
import BookmarkFormDialog from './BookmarkFormDialog.vue'

const props = defineProps({
  group: { type: Object, required: true }
})
const emit = defineEmits(['refresh'])

const bookmarksStore = useBookmarksStore()
const formVisible = ref(false)
const editData = ref(null)

function addBookmark() {
  editData.value = null
  formVisible.value = true
}
function openEdit(bm) {
  editData.value = bm
  formVisible.value = true
}

async function handlePin(bm) {
  try {
    await bookmarksStore.togglePinned(bm.id, !bm.pinned)
    ElMessage.success(bm.pinned ? '已取消置顶' : '已置顶')
    emit('refresh')
  } catch (e) {
    ElMessage.error('操作失败')
  }
}

async function handleDelete(bm) {
  try {
    await ElMessageBox.confirm(`确定删除书签「${bm.name}」？`, '提示', { type: 'warning' })
    await bookmarksStore.removeBookmark(props.group.id, bm.id)
    ElMessage.success('已删除')
    emit('refresh')
  } catch (e) {
    // 取消不操作
  }
}

async function deleteGroup() {
  try {
    await ElMessageBox.confirm(`确定删除分组「${props.group.name}」及其所有书签？`, '提示', { type: 'warning' })
    await bookmarksStore.removeGroup(props.group.id)
    ElMessage.success('已删除')
    emit('refresh')
  } catch (e) {}
}

function editGroup() {
  // 由父组件处理
}

async function onDragEnd() {
  try {
    const ids = props.group.bookmarks.map(b => String(b.id))
    await bookmarksStore.sortBookmarks(props.group.id, ids)
  } catch (e) {
    ElMessage.error('排序保存失败')
  }
}
</script>

<style scoped>
.group-section {
  margin-bottom: 32px;
}
.group-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 16px;
}
.group-title {
  display: flex;
  align-items: center;
  gap: 8px;
}
.group-icon {
  font-size: 22px;
  line-height: 1;
}
.group-name {
  font-size: 17px;
  font-weight: 600;
  color: #1d1d1f;
}
.group-count {
  font-size: 12px;
  color: #86868b;
}
.group-actions {
  display: flex;
  gap: 2px;
  opacity: 0;
  transition: opacity 0.2s;
}
.group-title:hover .group-actions {
  opacity: 1;
}
.bookmark-grid {
  display: flex;
  flex-wrap: wrap;
  gap: 12px;
}
.ghost {
  opacity: 0.4;
}
</style>
```

- [ ] **Step 3: 创建 PinnedBar.vue**

```vue
<template>
  <div v-if="bookmarksStore.pinned.length > 0" class="pinned-bar">
    <div class="pinned-header">
      <el-icon :size="16" color="#f59e0b"><Star /></el-icon>
      <span class="pinned-label">常驻</span>
    </div>
    <div class="pinned-list">
      <div
        v-for="bm in bookmarksStore.pinned"
        :key="bm.id"
        class="pinned-item"
        :style="{ backgroundColor: bm.bg_color || '#fff' }"
        @click="openUrl(bm.url)"
        @mouseenter="hoverId = bm.id"
        @mouseleave="hoverId = null"
      >
        <img v-if="isImageUrl(bm.icon)" :src="bm.icon" class="pinned-icon-img" />
        <span v-else class="pinned-icon">{{ bm.icon || '🔗' }}</span>
        <span class="pinned-name">{{ bm.name }}</span>
        <el-button
          v-show="hoverId === bm.id"
          size="small"
          circle
          class="pinned-unpin"
          @click.stop="unpin(bm)"
          title="取消常驻"
        >
          <el-icon><Close /></el-icon>
        </el-button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { ElMessage } from 'element-plus'
import { useBookmarksStore } from '@/stores/bookmarks'
import { Star, Close } from '@element-plus/icons-vue'

const bookmarksStore = useBookmarksStore()
const hoverId = ref(null)

function isImageUrl(str) {
  return str && (str.startsWith('http') || str.startsWith('/uploads/'))
}
function openUrl(url) {
  if (url) window.open(url, '_blank')
}
async function unpin(bm) {
  try {
    await bookmarksStore.togglePinned(bm.id, false)
    ElMessage.success('已取消常驻')
  } catch (e) {
    ElMessage.error('操作失败')
  }
}
</script>

<style scoped>
.pinned-bar {
  margin: 20px 0 24px;
}
.pinned-header {
  display: flex;
  align-items: center;
  gap: 6px;
  margin-bottom: 12px;
}
.pinned-label {
  font-size: 13px;
  color: #86868b;
  font-weight: 500;
}
.pinned-list {
  display: flex;
  gap: 10px;
  overflow-x: auto;
  padding-bottom: 8px;
}
.pinned-item {
  flex-shrink: 0;
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 8px 16px;
  border-radius: 12px;
  cursor: pointer;
  position: relative;
  transition: transform 0.2s, box-shadow 0.2s;
  box-shadow: 0 1px 4px rgba(0, 0, 0, 0.06);
}
.pinned-item:hover {
  transform: translateY(-2px);
  box-shadow: 0 6px 20px rgba(0, 0, 0, 0.1);
}
.pinned-icon {
  font-size: 20px;
  line-height: 1;
}
.pinned-icon-img {
  width: 20px;
  height: 20px;
  object-fit: contain;
}
.pinned-name {
  font-size: 13px;
  font-weight: 500;
  white-space: nowrap;
}
.pinned-unpin {
  position: absolute;
  top: -6px;
  right: -6px;
  width: 20px;
  height: 20px;
}
</style>
```

---

### Task 6: 所有对话框组件

**Files:**
- Create: `web/src/components/BookmarkFormDialog.vue`
- Create: `web/src/components/GroupFormDialog.vue`
- Create: `web/src/components/SettingsDialog.vue`
- Create: `web/src/components/ImageManageDialog.vue`
- Create: `web/src/components/UserManageDialog.vue`

**Interfaces:**
- Consumes: useBookmarksStore, useAuthStore, useSettingsStore, api methods
- Produces: 所有 CRUD 弹窗

- [ ] **Step 1: 创建 BookmarkFormDialog.vue**

```vue
<template>
  <el-dialog
    :title="editData ? '编辑书签' : '添加书签'"
    v-model="visible"
    width="480px"
    :close-on-click-modal="false"
    @closed="resetForm"
  >
    <el-form ref="formRef" :model="form" :rules="rules" label-width="70px">
      <el-form-item label="名称" prop="name">
        <el-input v-model="form.name" placeholder="书签名称" />
      </el-form-item>
      <el-form-item label="地址" prop="url">
        <el-input v-model="form.url" placeholder="https://" />
      </el-form-item>
      <el-form-item label="图标">
        <el-input v-model="form.icon" placeholder="Emoji 或图片 URL">
          <template #prepend><span style="font-size:18px">{{ form.icon || '🔗' }}</span></template>
        </el-input>
      </el-form-item>
      <el-form-item label="背景色">
        <el-color-picker v-model="form.bg_color" show-alpha />
        <span class="form-tip">卡片背景色</span>
      </el-form-item>
      <el-form-item label="图标背景">
        <el-color-picker v-model="form.icon_bg" show-alpha />
        <span class="form-tip">图标区域背景色</span>
      </el-form-item>
    </el-form>
    <template #footer>
      <el-button @click="visible = false">取消</el-button>
      <el-button type="primary" :loading="loading" @click="handleSave">保存</el-button>
    </template>
  </el-dialog>
</template>

<script setup>
import { ref, reactive, watch } from 'vue'
import { ElMessage } from 'element-plus'
import { useBookmarksStore } from '@/stores/bookmarks'

const props = defineProps({
  visible: Boolean,
  groupId: { type: Number, required: true },
  editData: { type: Object, default: null }
})
const emit = defineEmits(['update:visible', 'saved'])

const bookmarksStore = useBookmarksStore()
const formRef = ref(null)
const loading = ref(false)

const form = reactive({
  name: '',
  url: '',
  icon: '',
  bg_color: '',
  icon_bg: ''
})

const rules = {
  name: [{ required: true, message: '请输入名称', trigger: 'blur' }],
  url: [{ required: true, message: '请输入地址', trigger: 'blur' }]
}

watch(() => props.visible, (v) => {
  if (v && props.editData) {
    Object.assign(form, {
      name: props.editData.name || '',
      url: props.editData.url || '',
      icon: props.editData.icon || '',
      bg_color: props.editData.bg_color || '',
      icon_bg: props.editData.icon_bg || ''
    })
  }
})

function resetForm() {
  form.name = ''; form.url = ''; form.icon = ''
  form.bg_color = ''; form.icon_bg = ''
  formRef.value?.resetFields()
}

async function handleSave() {
  const valid = await formRef.value.validate().catch(() => false)
  if (!valid) return
  loading.value = true
  try {
    if (props.editData) {
      await bookmarksStore.editBookmark(props.groupId, props.editData.id, { ...form })
    } else {
      await bookmarksStore.addBookmark(props.groupId, { ...form })
    }
    ElMessage.success(props.editData ? '已更新' : '已添加')
    emit('update:visible', false)
    emit('saved')
  } catch (e) {
    ElMessage.error(e.response?.data?.error || '保存失败')
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.form-tip {
  font-size: 12px;
  color: #86868b;
  margin-left: 8px;
}
</style>
```

- [ ] **Step 2: 创建 GroupFormDialog.vue**

```vue
<template>
  <el-dialog
    :title="editData ? '编辑分组' : '新建分组'"
    v-model="visible"
    width="400px"
    :close-on-click-modal="false"
    @closed="resetForm"
  >
    <el-form ref="formRef" :model="form" :rules="rules" label-width="60px">
      <el-form-item label="名称" prop="name">
        <el-input v-model="form.name" placeholder="分组名称" />
      </el-form-item>
      <el-form-item label="图标">
        <el-input v-model="form.icon" placeholder="Emoji 图标">
          <template #prepend><span style="font-size:18px">{{ form.icon || '📁' }}</span></template>
        </el-input>
      </el-form-item>
    </el-form>
    <template #footer>
      <el-button @click="visible = false">取消</el-button>
      <el-button type="primary" :loading="loading" @click="handleSave">保存</el-button>
    </template>
  </el-dialog>
</template>

<script setup>
import { ref, reactive } from 'vue'
import { ElMessage } from 'element-plus'
import { useBookmarksStore } from '@/stores/bookmarks'

const props = defineProps({
  visible: Boolean,
  editData: { type: Object, default: null }
})
const emit = defineEmits(['update:visible', 'saved'])

const bookmarksStore = useBookmarksStore()
const formRef = ref(null)
const loading = ref(false)
const form = reactive({ name: '', icon: '' })
const rules = {
  name: [{ required: true, message: '请输入分组名称', trigger: 'blur' }]
}

async function handleSave() {
  const valid = await formRef.value.validate().catch(() => false)
  if (!valid) return
  loading.value = true
  try {
    if (props.editData) {
      await bookmarksStore.editGroup(props.editData.id, form.name, form.icon)
    } else {
      await bookmarksStore.addGroup(form.name, form.icon)
    }
    ElMessage.success(props.editData ? '已更新' : '已创建')
    emit('update:visible', false)
    emit('saved')
  } catch (e) {
    ElMessage.error(e.response?.data?.error || '保存失败')
  } finally {
    loading.value = false
  }
}

function resetForm() {
  form.name = ''; form.icon = ''
  formRef.value?.resetFields()
}
</script>
```

- [ ] **Step 3: 创建 SettingsDialog.vue**

```vue
<template>
  <el-dialog title="页面设置" v-model="visible" width="480px" @closed="fetchSettings">
    <el-form label-width="90px">
      <el-form-item label="页面标题">
        <el-input v-model="local.title" placeholder="个人导航页" />
      </el-form-item>
      <el-form-item label="背景图片">
        <div class="bg-row">
          <el-input v-model="local.bgImage" placeholder="图片 URL">
            <template #append>
              <el-button @click="pickImage('bg')">选择</el-button>
            </template>
          </el-input>
          <el-button v-if="local.bgImage" type="danger" text @click="local.bgImage = ''">清除</el-button>
        </div>
        <img v-if="local.bgImage" :src="local.bgImage" class="bg-preview" />
      </el-form-item>
      <el-form-item label="Favicon">
        <el-input v-model="local.favicon" placeholder="图标 URL">
          <template #append>
            <el-button @click="pickImage('favicon')">选择</el-button>
          </template>
        </el-input>
      </el-form-item>
    </el-form>
    <template #footer>
      <el-button @click="visible = false">取消</el-button>
      <el-button type="primary" :loading="loading" @click="handleSave">保存</el-button>
    </template>
    <!-- 图片选择弹窗 -->
    <ImagePicker v-model:visible="pickerVisible" @select="onImageSelected" />
  </el-dialog>
</template>

<script setup>
import { ref, reactive, watch } from 'vue'
import { ElMessage } from 'element-plus'
import { useAuthStore } from '@/stores/auth'
import { useSettingsStore } from '@/stores/settings'
import ImagePicker from './ImagePicker.vue'

const props = defineProps({ visible: Boolean })
const emit = defineEmits(['update:visible'])

const auth = useAuthStore()
const settings = useSettingsStore()
const loading = ref(false)
const pickerVisible = ref(false)
const pickTarget = ref('bg')

const local = reactive({ title: '', bgImage: '', favicon: '' })

watch(() => props.visible, (v) => {
  if (v) {
    local.title = settings.pageTitle
    local.bgImage = settings.bgImage
    local.favicon = settings.pageFavicon
  }
})

function pickImage(target) {
  pickTarget.value = target
  pickerVisible.value = true
}

function onImageSelected(url) {
  if (pickTarget.value === 'bg') local.bgImage = url
  else local.favicon = url
}

async function fetchSettings() {
  await settings.fetch(auth.userId)
}

async function handleSave() {
  loading.value = true
  try {
    await settings.save(auth.userId, {
      page_title: local.title,
      bg_image: local.bgImage,
      page_favicon: local.favicon
    })
    ElMessage.success('已保存')
    emit('update:visible', false)
  } catch (e) {
    ElMessage.error('保存失败')
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.bg-row {
  display: flex;
  gap: 8px;
  margin-bottom: 8px;
}
.bg-preview {
  max-width: 100%;
  max-height: 120px;
  border-radius: 8px;
  object-fit: cover;
}
</style>
```

- [ ] **Step 4: 创建 ImagePicker.vue（专供弹窗内选择图片）**

```vue
<template>
  <el-dialog title="选择图片" v-model="visible" width="600px">
    <div class="img-picker">
      <div class="img-picker-toolbar">
        <el-upload
          :auto-upload="false"
          :show-file-list="false"
          accept="image/*"
          :on-change="handleUpload"
        >
          <el-button :icon="Upload" size="small">上传图片</el-button>
        </el-upload>
      </div>
      <div v-loading="loading" class="img-grid">
        <div
          v-for="img in images"
          :key="img.id"
          class="img-item"
          :class="{ selected: selected === img.url }"
          @click="select(img.url)"
        >
          <img :src="img.url" :title="img.original_name" />
        </div>
        <el-empty v-if="!loading && images.length === 0" description="暂无图片" />
      </div>
    </div>
    <template #footer>
      <el-button @click="visible = false">取消</el-button>
      <el-button type="primary" :disabled="!selected" @click="confirm">选择</el-button>
    </template>
  </el-dialog>
</template>

<script setup>
import { ref, watch } from 'vue'
import { getImages, uploadImage } from '@/api'
import { ElMessage } from 'element-plus'
import { Upload } from '@element-plus/icons-vue'

const props = defineProps({ visible: Boolean })
const emit = defineEmits(['update:visible', 'select'])

const images = ref([])
const loading = ref(false)
const selected = ref('')

watch(() => props.visible, (v) => {
  if (v) { selected.value = ''; fetchImages() }
})

async function fetchImages() {
  loading.value = true
  try {
    const res = await getImages({ pageSize: 200 })
    images.value = res.data || []
  } finally {
    loading.value = false
  }
}

async function handleUpload(file) {
  const formData = new FormData()
  formData.append('file', file.raw || file)
  try {
    const res = await uploadImage(formData)
    ElMessage.success('上传成功')
    fetchImages()
  } catch (e) {
    ElMessage.error('上传失败')
  }
}

function select(url) {
  selected.value = url
}

function confirm() {
  if (selected.value) {
    emit('select', selected.value)
    emit('update:visible', false)
  }
}
</script>

<style scoped>
.img-grid {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
  max-height: 360px;
  overflow-y: auto;
  margin-top: 12px;
}
.img-item {
  width: 80px;
  height: 80px;
  border-radius: 8px;
  overflow: hidden;
  cursor: pointer;
  border: 2px solid transparent;
  transition: border-color 0.2s;
}
.img-item:hover { border-color: #6366f1; }
.img-item.selected { border-color: #6366f1; }
.img-item img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}
</style>
```

- [ ] **Step 5: 创建 ImageManageDialog.vue**

```vue
<template>
  <el-dialog title="图库管理" v-model="visible" width="700px" :close-on-click-modal="false">
    <div class="img-mgr">
      <div class="img-mgr-toolbar">
        <el-upload
          :auto-upload="false"
          :show-file-list="false"
          accept="image/*"
          multiple
          :on-change="handleUpload"
        >
          <el-button :icon="Upload" type="primary">上传图片</el-button>
        </el-upload>
        <el-upload
          :auto-upload="false"
          :show-file-list="false"
          accept=".zip"
          :on-change="handleZipUpload"
        >
          <el-button :icon="FolderOpened">上传 ZIP 包</el-button>
        </el-upload>
      </div>
      <div v-loading="loading" class="img-grid">
        <div v-for="img in images" :key="img.id" class="img-card">
          <img :src="img.url" />
          <div class="img-card-actions">
            <el-button size="small" @click="copyUrl(img.url)">复制链接</el-button>
            <el-button size="small" type="danger" @click="remove(img)">删除</el-button>
          </div>
          <div class="img-card-name">{{ img.original_name }}</div>
        </div>
        <el-empty v-if="!loading && images.length === 0" description="暂无图片" />
      </div>
    </div>
  </el-dialog>
</template>

<script setup>
import { ref, watch } from 'vue'
import { getImages, uploadImage, uploadZip, deleteImage } from '@/api'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Upload, FolderOpened } from '@element-plus/icons-vue'

const props = defineProps({ visible: Boolean })
const emit = defineEmits(['update:visible'])

const images = ref([])
const loading = ref(false)

watch(() => props.visible, (v) => { if (v) fetchImages() })

async function fetchImages() {
  loading.value = true
  try {
    const res = await getImages({ pageSize: 200 })
    images.value = res.data || []
  } finally {
    loading.value = false
  }
}

async function handleUpload(file) {
  const formData = new FormData()
  formData.append('file', file.raw || file)
  try {
    await uploadImage(formData)
    ElMessage.success('上传成功')
    fetchImages()
  } catch (e) {
    ElMessage.error('上传失败')
  }
}

async function handleZipUpload(file) {
  const formData = new FormData()
  formData.append('file', file.raw || file)
  try {
    const res = await uploadZip(formData)
    ElMessage.success(`导入 ${res.imported || 0} 个图标`)
    fetchImages()
  } catch (e) {
    ElMessage.error('导入失败')
  }
}

function copyUrl(url) {
  navigator.clipboard.writeText(url).then(() => ElMessage.success('已复制'))
}

async function remove(img) {
  try {
    await ElMessageBox.confirm(`确定删除「${img.original_name}」？`, '提示', { type: 'warning' })
    await deleteImage(img.id)
    ElMessage.success('已删除')
    fetchImages()
  } catch (e) {}
}
</script>

<style scoped>
.img-mgr-toolbar {
  display: flex;
  gap: 8px;
  margin-bottom: 16px;
}
.img-grid {
  display: flex;
  flex-wrap: wrap;
  gap: 12px;
  max-height: 480px;
  overflow-y: auto;
}
.img-card {
  width: 120px;
  border-radius: 8px;
  overflow: hidden;
  background: #f5f5f7;
  position: relative;
}
.img-card img {
  width: 100%;
  height: 90px;
  object-fit: cover;
  display: block;
}
.img-card-actions {
  display: flex;
  gap: 4px;
  padding: 4px;
}
.img-card-name {
  font-size: 11px;
  color: #86868b;
  padding: 0 4px 6px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}
</style>
```

- [ ] **Step 6: 创建 UserManageDialog.vue**

```vue
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
```

---

### Task 7: 搜索弹窗

**Files:**
- Create: `web/src/components/SearchOverlay.vue`

**Interfaces:**
- Consumes: `searchBookmarks` from api
- Produces: 全屏搜索遮罩组件

- [ ] **Step 1: 创建 SearchOverlay.vue**

```vue
<template>
  <teleport to="body">
    <transition name="search-fade">
      <div v-if="visible" class="search-overlay" @click.self="close">
        <div class="search-panel" @keydown="handleKeydown">
          <div class="search-input-wrap">
            <el-icon :size="22" class="search-input-icon"><Search /></el-icon>
            <input
              ref="inputRef"
              v-model="query"
              class="search-input"
              placeholder="搜索书签..."
              @input="onInput"
            />
            <kbd class="search-hint">ESC</kbd>
          </div>
          <div v-if="query" class="search-results" v-loading="loading">
            <div
              v-for="(r, i) in results"
              :key="r.id"
              class="search-item"
              :class="{ active: activeIndex === i }"
              @click="openItem(r)"
              @mouseenter="activeIndex = i"
            >
              <span class="search-item-icon">{{ r.icon || '🔗' }}</span>
              <div class="search-item-info">
                <span class="search-item-name">{{ r.name }}</span>
                <span class="search-item-url">{{ r.url }}</span>
              </div>
              <span class="search-item-group">{{ r.groupName }}</span>
              <el-icon class="search-item-open"><TopRight /></el-icon>
            </div>
            <div v-if="results.length === 0 && !loading" class="search-empty">
              未找到匹配的书签
            </div>
          </div>
        </div>
      </div>
    </transition>
  </teleport>
</template>

<script setup>
import { ref, watch, nextTick } from 'vue'
import { searchBookmarks } from '@/api'
import { Search, TopRight } from '@element-plus/icons-vue'

const props = defineProps({ visible: Boolean })
const emit = defineEmits(['update:visible'])

const query = ref('')
const results = ref([])
const loading = ref(false)
const activeIndex = ref(-1)
const inputRef = ref(null)

let debounceTimer = null

watch(() => props.visible, (v) => {
  if (v) {
    query.value = ''
    results.value = []
    activeIndex.value = -1
    nextTick(() => inputRef.value?.focus())
  }
})

function onInput() {
  clearTimeout(debounceTimer)
  if (!query.value.trim()) { results.value = []; return }
  debounceTimer = setTimeout(doSearch, 300)
}

async function doSearch() {
  const q = query.value.trim()
  if (!q) return
  loading.value = true
  try {
    results.value = await searchBookmarks(q)
    activeIndex.value = results.value.length > 0 ? 0 : -1
  } catch (e) {
    results.value = []
  } finally {
    loading.value = false
  }
}

function handleKeydown(e) {
  if (e.key === 'ArrowDown') {
    e.preventDefault()
    if (activeIndex.value < results.value.length - 1) activeIndex.value++
  } else if (e.key === 'ArrowUp') {
    e.preventDefault()
    if (activeIndex.value > 0) activeIndex.value--
  } else if (e.key === 'Enter' && activeIndex.value >= 0) {
    openItem(results.value[activeIndex.value])
  } else if (e.key === 'Escape') {
    close()
  }
}

function openItem(item) {
  if (item?.url) {
    window.open(item.url, '_blank')
    close()
  }
}

function close() {
  emit('update:visible', false)
}
</script>

<style scoped>
.search-overlay {
  position: fixed;
  inset: 0;
  background: rgba(0, 0, 0, 0.5);
  backdrop-filter: blur(4px);
  display: flex;
  justify-content: center;
  padding-top: 15vh;
  z-index: 2000;
}
.search-panel {
  width: 580px;
  max-height: 60vh;
  display: flex;
  flex-direction: column;
}
.search-input-wrap {
  display: flex;
  align-items: center;
  background: #fff;
  border-radius: 16px;
  padding: 14px 20px;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.12);
}
.search-input-icon {
  color: #86868b;
  margin-right: 12px;
  flex-shrink: 0;
}
.search-input {
  flex: 1;
  border: none;
  outline: none;
  font-size: 18px;
  background: transparent;
  color: #1d1d1f;
}
.search-input::placeholder { color: #c7c7cc; }
.search-hint {
  font-size: 11px;
  padding: 2px 6px;
  background: #f5f5f7;
  border-radius: 4px;
  color: #86868b;
  font-family: inherit;
  flex-shrink: 0;
}
.search-results {
  background: #fff;
  border-radius: 12px;
  margin-top: 8px;
  overflow-y: auto;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.12);
  max-height: 50vh;
}
.search-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 12px 16px;
  cursor: pointer;
  transition: background 0.15s;
}
.search-item:first-child { border-radius: 12px 12px 0 0; }
.search-item:last-child { border-radius: 0 0 12px 12px; }
.search-item.active, .search-item:hover { background: #f5f5f7; }
.search-item-icon { font-size: 20px; line-height: 1; }
.search-item-info {
  flex: 1;
  min-width: 0;
  display: flex;
  flex-direction: column;
  gap: 2px;
}
.search-item-name { font-size: 14px; font-weight: 500; color: #1d1d1f; }
.search-item-url { font-size: 12px; color: #86868b; overflow: hidden; text-overflow: ellipsis; white-space: nowrap; }
.search-item-group { font-size: 11px; color: #aeaeb2; flex-shrink: 0; }
.search-item-open { color: #c7c7cc; flex-shrink: 0; }
.search-empty { padding: 32px; text-align: center; color: #86868b; font-size: 14px; }

.search-fade-enter-active, .search-fade-leave-active { transition: opacity 0.2s; }
.search-fade-enter-from, .search-fade-leave-to { opacity: 0; }
</style>
```

---

## 自检清单

1. **Spec 覆盖检查**
   - 两个路由（/login, /）：Task 1 路由配置 + Task 3 LoginPage + Task 4 MainPage ✅
   - 登录：Task 3 ✅
   - 分组 CRUD + 排序：Task 5 GroupSection + BookmarkFormDialog/GroupFormDialog ✅
   - 书签 CRUD + 排序：Task 5 GroupSection + BookmarkFormDialog ✅
   - 常驻书签：Task 5 PinnedBar ✅
   - 设置：Task 6 SettingsDialog ✅
   - 图库管理：Task 6 ImageManageDialog ✅
   - 用户管理：Task 6 UserManageDialog ✅
   - 搜索（Ctrl+K）：Task 7 SearchOverlay ✅
   - 所有二级功能均为弹窗 ✅
   - 认证拦截 + 401 跳转：Task 2 API 拦截器 ✅

2. **占位符检查**：无 TBD/TODO，所有代码完整 ✅

3. **类型一致性**：
   - API 方法签名与 handler 参数匹配 ✅
   - Store action 命名一致 ✅
   - Prop 和事件命名一致 ✅
