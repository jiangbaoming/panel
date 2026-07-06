<template>
  <div class="main-page" :style="bgStyle">
    <div class="main-container">
      <TopBar
        @open-settings="settingsVisible = true"
        @open-images="imagesVisible = true"
        @open-users="usersVisible = true"
      />
      <div class="main-content" v-loading="loading">
        <!-- 常驻书签 -->
        <PinnedBar @pin-changed="bookmarks.fetchAll()" />

        <!-- 分组列表 -->
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

    <!-- 分组表单弹窗 -->
    <GroupFormDialog
      v-model:visible="groupFormVisible"
      :edit-data="groupEditData"
      @saved="bookmarks.fetchAll()"
    />
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
import GroupFormDialog from '@/components/GroupFormDialog.vue'

const router = useRouter()
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
  min-height: 60vh;
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
