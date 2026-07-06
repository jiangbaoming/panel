<template>
  <div class="main-page" :style="bgStyle">
    <!-- 上部：顶栏（仅右侧按钮） -->
    <TopBar
      @open-settings="settingsVisible = true"
      @open-images="imagesVisible = true"
      @open-users="usersVisible = true"
    />

    <!-- 中部：欢迎语 + 时间 + 搜索 + 常驻 + 分组 -->
    <div class="main-middle">
      <div class="hero-section" v-loading="loading" element-loading-background="transparent">
        <!-- 欢迎语 -->
        <div class="greeting">
          <h1 class="greeting-text">{{ greetingText }}</h1>
          <p class="time-text">{{ currentTime }}</p>
          <p class="date-text">{{ currentDate }}</p>
        </div>

        <!-- 搜索框 -->
        <div class="search-wrap" @click="searchVisible = true">
          <el-icon :size="20" class="search-icon"><Search /></el-icon>
          <span class="search-placeholder">搜索书签...</span>
          <kbd class="search-hint">Ctrl + K</kbd>
        </div>

        <!-- 常驻书签 -->
        <PinnedBar />

        <!-- 分组书签列表 -->
        <template v-for="group in bookmarks.groups" :key="group.id">
          <GroupSection
            :group="group"
            @refresh="bookmarks.fetchAll()"
            @edit-group="openGroupEdit(group)"
          />
        </template>

        <!-- 新建分组按钮 -->
        <div class="add-group-btn" v-if="bookmarks.groups.length > 0">
          <el-button @click="openGroupAdd()" :icon="Plus" text size="large">
            新建分组
          </el-button>
        </div>
      </div>
    </div>

    <!-- 下部：页脚 -->
    <footer class="main-footer" v-if="settings.footer">
      <div class="footer-content">{{ settings.footer }}</div>
    </footer>

    <!-- 弹窗们 -->
    <SearchOverlay v-model:visible="searchVisible" />
    <SettingsDialog v-model:visible="settingsVisible" />
    <ImageManageDialog v-model:visible="imagesVisible" />
    <UserManageDialog v-model:visible="usersVisible" />
    <GroupFormDialog
      v-model:visible="groupFormVisible"
      :edit-data="groupEditData"
      @saved="bookmarks.fetchAll()"
    />
  </div>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted } from 'vue'
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
import GroupFormDialog from '@/components/GroupFormDialog.vue'

const auth = useAuthStore()
const bookmarks = useBookmarksStore()
const settings = useSettingsStore()

const loading = ref(true)
const searchVisible = ref(false)
const settingsVisible = ref(false)
const imagesVisible = ref(false)
const usersVisible = ref(false)
const groupFormVisible = ref(false)
const groupEditData = ref(null)

// 实时时钟
const now = ref(new Date())
let timer = null

const greetingText = computed(() => {
  const hour = now.value.getHours()
  const name = auth.user?.username || ''
  if (hour < 6) return `夜深了，${name}`
  if (hour < 9) return `早上好，${name}`
  if (hour < 12) return `上午好，${name}`
  if (hour < 14) return `中午好，${name}`
  if (hour < 18) return `下午好，${name}`
  if (hour < 22) return `晚上好，${name}`
  return `夜深了，${name}`
})

const currentTime = computed(() => {
  return now.value.toLocaleTimeString('zh-CN', { hour: '2-digit', minute: '2-digit' })
})

const currentDate = computed(() => {
  const weekDays = ['星期日', '星期一', '星期二', '星期三', '星期四', '星期五', '星期六']
  const d = now.value
  return `${d.getFullYear()}年${d.getMonth() + 1}月${d.getDate()}日 ${weekDays[d.getDay()]}`
})

const bgStyle = computed(() => {
  if (settings.bgImage) {
    return {
      backgroundImage: `url(${settings.bgImage})`,
      backgroundSize: 'cover',
      backgroundPosition: 'center',
      backgroundAttachment: 'fixed'
    }
  }
  return {}
})

onMounted(async () => {
  // 启动时钟
  timer = setInterval(() => { now.value = new Date() }, 1000)
  // 加载数据
  try {
    await Promise.all([
      bookmarks.fetchAll(),
      settings.fetch(auth.userId)
    ])
  } catch (e) {
    // 静默处理
  } finally {
    loading.value = false
  }
})

onUnmounted(() => {
  clearInterval(timer)
  document.removeEventListener('keydown', handleKeydown)
})

// 分组操作
function openGroupAdd() {
  groupEditData.value = null
  groupFormVisible.value = true
}
function openGroupEdit(group) {
  groupEditData.value = group
  groupFormVisible.value = true
}

// Ctrl+K 搜索快捷键
function handleKeydown(e) {
  if ((e.ctrlKey || e.metaKey) && e.key === 'k') {
    e.preventDefault()
    searchVisible.value = true
  }
}
onMounted(() => document.addEventListener('keydown', handleKeydown))
</script>

<style scoped>
/* ====== 页面整体 ====== */
.main-page {
  min-height: 100vh;
  display: flex;
  flex-direction: column;
  background: #f5f5f7;
  transition: background 0.3s;
}

/* ====== 中部 ====== */
.main-middle {
  flex: 1;
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 0 24px;
}

.hero-section {
  width: 100%;
  max-width: 1100px;
  padding-bottom: 60px;
}

/* 欢迎语 + 时间 */
.greeting {
  text-align: center;
  padding: 60px 0 32px;
}

.greeting-text {
  font-size: 32px;
  font-weight: 700;
  color: #1d1d1f;
  margin: 0 0 12px;
  letter-spacing: 1px;
}

.time-text {
  font-size: 56px;
  font-weight: 300;
  color: #1d1d1f;
  margin: 0 0 4px;
  font-variant-numeric: tabular-nums;
  letter-spacing: 2px;
}

.date-text {
  font-size: 15px;
  color: #86868b;
  margin: 0;
}

/* 搜索框 */
.search-wrap {
  display: flex;
  align-items: center;
  max-width: 480px;
  margin: 0 auto 40px;
  padding: 14px 20px;
  background: #fff;
  border-radius: 16px;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.06);
  cursor: pointer;
  transition: box-shadow 0.2s, transform 0.2s;
}

.search-wrap:hover {
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.1);
  transform: translateY(-1px);
}

.search-icon {
  color: #86868b;
  margin-right: 12px;
  flex-shrink: 0;
}

.search-placeholder {
  flex: 1;
  font-size: 16px;
  color: #c7c7cc;
}

.search-hint {
  font-size: 11px;
  padding: 2px 8px;
  background: #f5f5f7;
  border-radius: 6px;
  color: #aeaeb2;
  font-family: inherit;
  flex-shrink: 0;
}

/* 新建分组 */
.add-group-btn {
  text-align: center;
  margin-top: 20px;
}

/* ====== 下部：页脚 ====== */
.main-footer {
  padding: 24px;
  text-align: center;
  border-top: 1px solid rgba(0, 0, 0, 0.06);
  max-width: 1100px;
  width: 100%;
  margin: 0 auto;
}

.footer-content {
  font-size: 13px;
  color: #aeaeb2;
  white-space: pre-line;
}
</style>
