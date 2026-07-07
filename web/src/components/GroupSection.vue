<template>
  <div class="group-section">
    <div class="group-header" :class="{ 'is-editable': editMode }">
      <div class="group-title">
        <span class="group-icon">{{ group.icon || '📁' }}</span>
        <h2 class="group-name">{{ group.name }}</h2>
        <span class="group-count">{{ group.bookmarks?.length || 0 }} 个</span>
        <div v-if="editMode" class="group-actions">
          <el-button text size="small" @click="addBookmark" title="添加书签">
            <el-icon :size="14"><Plus /></el-icon>
          </el-button>
          <el-button text size="small" @click="editGroup" title="编辑分组">
            <el-icon :size="14"><Edit /></el-icon>
          </el-button>
          <el-button text size="small" @click="deleteGroup" title="删除分组" type="danger">
            <el-icon :size="14"><Delete /></el-icon>
          </el-button>
        </div>
      </div>
    </div>
    <draggable
      v-if="editMode"
      class="bookmark-grid"
      :list="group.bookmarks"
      item-key="id"
      tag="div"
      group="bookmarks"
      handle=".bookmark-row"
      ghost-class="ghost"
      @end="onDragEnd"
    >
      <template #item="{ element }">
        <BookmarkCard
          :bookmark="element"
          :edit-mode="editMode"
          @edit="openEdit(element)"
          @pin="handlePin(element)"
          @delete="handleDelete(element)"
        />
      </template>
    </draggable>
    <div v-else class="bookmark-grid">
      <BookmarkCard
        v-for="bm in group.bookmarks"
        :key="bm.id"
        :bookmark="bm"
        :edit-mode="editMode"
        @edit="openEdit(bm)"
        @pin="handlePin(bm)"
        @delete="handleDelete(bm)"
      />
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
  group: { type: Object, required: true },
  editMode: { type: Boolean, default: false }
})
const emit = defineEmits(['refresh', 'edit-group', 'bookmarks-reordered'])

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
  } catch {}
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
  emit('edit-group', props.group)
}

function onDragEnd() {
  const ids = props.group.bookmarks.map(b => String(b.id))
  emit('bookmarks-reordered', { groupId: props.group.id, ids })
}
</script>

<style scoped>
.group-section {
  margin-bottom: 32px;
  padding: 24px;
}

.group-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 16px;
}
.group-header.is-editable {
  cursor: grab;
}
.group-header.is-editable:active {
  cursor: grabbing;
}
.group-title {
  display: flex;
  align-items: center;
  gap: 8px;
}
.group-icon {
  font-size: 22px;
  line-height: 1;
  transition: transform var(--transition-fast);
}
.group-section:hover .group-icon {
  transform: scale(1.1);
}
.group-name {
  font-size: 17px;
  font-weight: 600;
  color: var(--color-text-primary);
  text-shadow: 0 1px 2px rgba(0, 0, 0, 0.3);
}
.group-count {
  font-size: 12px;
  color: var(--color-text-secondary);
  padding: 2px 8px;
  border-radius: 10px;
}
.group-actions {
  display: flex;
  gap: 4px;
}
.group-actions .el-button {
  width: 24px;
  height: 24px;
  min-height: 24px;
  color: rgba(255, 255, 255, 0.55);
  background: transparent;
  border: none;
}
.group-actions .el-button:hover {
  color: rgba(255, 255, 255, 0.92);
  background: transparent;
}
.group-actions .el-button--danger {
  color: rgba(255, 100, 100, 0.65);
}
.group-actions .el-button--danger:hover {
  color: rgba(255, 100, 100, 0.95);
}
.group-actions .el-button .el-icon {
  margin: 0;
}
.bookmark-grid {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
}
.ghost {
  opacity: 0.4;
}

@media (max-width: 768px) {
  .group-section {
    padding: 12px 0;
    margin-bottom: 20px;
  }
  .group-header {
    margin-bottom: 10px;
  }
  .group-name {
    font-size: 15px;
  }
  .group-icon {
    font-size: 18px;
  }
  .bookmark-grid {
    gap: 6px;
  }
}
</style>