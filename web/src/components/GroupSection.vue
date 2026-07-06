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
