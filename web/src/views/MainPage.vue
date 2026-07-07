<template>
  <div class="main-page" :style="bgStyle">
    <!-- 上部：顶栏（仅右侧按钮） -->
    <TopBar :edit-mode="editMode" @toggle-edit="editMode = !editMode" @open-settings="settingsVisible = true" />

    <!-- 中部：欢迎语 + 时间 + 搜索 + 常驻 + 分组 -->
    <div class="main-middle">
      <div class="hero-section" v-loading="loading" element-loading-background="transparent">
        <!-- 欢迎语 -->
        <GreetingSection />

        <!-- 搜索框 -->
        <SearchBar v-model="searchQuery" :has-match="hasFilteredResults" />

        <!-- 常驻书签（搜索时隐藏） -->
        <PinnedBar v-if="!searchQuery.trim()" :edit-mode="editMode" />

        <!-- 编辑模式工具条 -->
        <div class="edit-toolbar" v-if="editMode">
          <el-button @click="openGroupAdd()" :icon="Plus" plain>新建分组</el-button>
        </div>

        <!-- 编辑模式且无搜索：可拖动排序 -->
        <draggable
          v-if="editMode && !searchQuery.trim()"
          :list="bookmarks.groups"
          item-key="id"
          handle=".group-header"
          ghost-class="group-ghost"
          @end="onGroupDragEnd"
        >
          <template #item="{ element: group }">
            <GroupSection
              :group="group"
              :edit-mode="editMode"
              @refresh="bookmarks.fetchAll()"
              @edit-group="openGroupEdit(group)"
              @bookmarks-reordered="onBookmarksReordered"
            />
          </template>
        </draggable>
        <!-- 非编辑模式或搜索模式：直接渲染 -->
        <template v-else>
          <GroupSection
            v-for="group in filteredGroups"
            :key="group.id"
            :group="group"
            :edit-mode="editMode"
            @refresh="bookmarks.fetchAll()"
            @edit-group="openGroupEdit(group)"
          />
        </template>
      </div>
    </div>

    <!-- 下部：页脚 -->
    <footer class="main-footer" v-if="settings.footer">
      <div class="footer-content">{{ settings.footer }}</div>
    </footer>

    <!-- 弹窗们 -->
    <SettingsDialog v-model:visible="settingsVisible" />
    <GroupFormDialog
      v-model:visible="groupFormVisible"
      :edit-data="groupEditData"
      @saved="bookmarks.fetchAll()"
    />
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useAuthStore } from '@/stores/auth'
import { useBookmarksStore } from '@/stores/bookmarks'
import { useSettingsStore } from '@/stores/settings'
import { Plus } from '@element-plus/icons-vue'
import draggable from 'vuedraggable'
import { ElMessage } from 'element-plus'
import TopBar from '@/components/TopBar.vue'
import GreetingSection from '@/components/GreetingSection.vue'
import SearchBar from '@/components/SearchBar.vue'
import PinnedBar from '@/components/PinnedBar.vue'
import GroupSection from '@/components/GroupSection.vue'
import SettingsDialog from '@/components/SettingsDialog.vue'
import GroupFormDialog from '@/components/GroupFormDialog.vue'

const auth = useAuthStore()
const bookmarks = useBookmarksStore()
const settings = useSettingsStore()

const editMode = ref(false)
const loading = ref(true)
const settingsVisible = ref(false)
const groupFormVisible = ref(false)
const groupEditData = ref(null)

// 搜索关键词（通过 SearchBar 的 v-model 传入）
const searchQuery = ref('')

// 是否有书签匹配结果
const hasFilteredResults = computed(() => {
  const q = searchQuery.value.trim().toLowerCase()
  if (!q) return true
  return bookmarks.groups.some(g =>
    (g.bookmarks || []).some(b => b.name.toLowerCase().includes(q) || b.url.toLowerCase().includes(q))
  )
})

// 搜索过滤：匹配书签名称或 URL
const filteredGroups = computed(() => {
  const q = searchQuery.value.trim().toLowerCase()
  if (!q) return bookmarks.groups
  return bookmarks.groups.map(g => ({
    ...g,
    bookmarks: (g.bookmarks || []).filter(b =>
      b.name.toLowerCase().includes(q) || b.url.toLowerCase().includes(q)
    )
  })).filter(g => g.bookmarks.length > 0)
})

const bgStyle = computed(() => {
  if (settings.bgImage) {
    const isUrl = settings.bgImage.startsWith('http') || settings.bgImage.startsWith('/')
    return {
      backgroundImage: isUrl ? `url(${settings.bgImage})` : settings.bgImage,
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

// 分组操作
function openGroupAdd() {
  groupEditData.value = null
  groupFormVisible.value = true
}
function openGroupEdit(group) {
  groupEditData.value = group
  groupFormVisible.value = true
}

async function onGroupDragEnd() {
  const ids = bookmarks.groups.map(g => String(g.id))
  try {
    await bookmarks.sortGroups(ids)
  } catch {}
}

async function onBookmarksReordered(payload) {
  try {
    await bookmarks.sortBookmarks(payload.groupId, payload.ids)
  } catch {}
}

</script>

<style scoped>
/* ====== 页面整体 ====== */
.main-page {
  min-height: 100vh;
  display: flex;
  flex-direction: column;
  background: var(--color-bg);
  background-image:
    radial-gradient(ellipse at 50% 0%, rgba(64, 158, 255, 0.04) 0%, transparent 60%),
    radial-gradient(ellipse at 80% 100%, rgba(64, 158, 255, 0.03) 0%, transparent 40%);
  transition: background 0.3s;
  position: relative;
  isolation: isolate;
}
.main-page::before {
  content: '';
  position: fixed;
  inset: 0;
  background: rgba(0, 0, 0, 0.45);
  z-index: -1;
  pointer-events: none;
}

/* ====== 中部 ====== */
.main-middle {
  flex: 1;
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 0 24px;
  --color-text-primary: rgba(255, 255, 255, 0.92);
  --color-text-secondary: rgba(255, 255, 255, 0.65);
  --color-text-tertiary: rgba(255, 255, 255, 0.45);
}

.hero-section {
  width: 100%;
  max-width: 1100px;
  padding-bottom: 60px;
}

@media (max-width: 768px) {
  .main-middle {
    padding: 0 12px;
  }
  .hero-section {
    padding-bottom: 30px;
  }
  .edit-toolbar {
    padding: 0 12px;
  }
  .main-footer {
    padding: 16px 12px;
  }
}

/* 新建分组（顶部） */
.edit-toolbar {
  display: flex;
  gap: 8px;
  padding: 0px 24px;
}
.edit-toolbar .el-button {
  background: transparent;
  color: rgba(255, 255, 255, 0.7);
  border-color: rgba(255, 255, 255, 0.25);
}
.edit-toolbar .el-button:hover {
  color: rgba(255, 255, 255, 0.92);
  border-color: rgba(255, 255, 255, 0.5);
  background: transparent;
}

.group-ghost {
  opacity: 0.3;
  filter: blur(2px);
}

/* 新建分组（底部，已废弃但保留） */
.add-group-btn {
  text-align: center;
  margin-top: 20px;
}

/* ====== 下部：页脚 ====== */
.main-footer {
  padding: 24px;
  text-align: center;
  border-top: 1px solid var(--color-border);
  max-width: 1100px;
  width: 100%;
  margin: 0 auto;
  --color-text-primary: rgba(255, 255, 255, 0.92);
  --color-text-secondary: rgba(255, 255, 255, 0.65);
  --color-text-tertiary: rgba(255, 255, 255, 0.45);
}

.footer-content {
  font-size: 13px;
  color: var(--color-text-tertiary);
  white-space: pre-line;
  text-shadow: 0 1px 2px rgba(0, 0, 0, 0.25);
}
</style>