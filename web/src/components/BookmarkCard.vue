<template>
  <div
    class="bookmark-row"
    :class="{ 'is-editable': editMode }"
    :style="rowStyle"
    @click="openUrl"
  >
    <div class="row-icon" :style="iconStyle">
      <img v-if="isImageUrl(bookmark.icon)" :src="bookmark.icon" />
      <span v-else>{{ bookmark.icon || '🔗' }}</span>
    </div>
    <div class="row-body">
      <span class="row-name">{{ bookmark.name }}</span>
      <span class="row-url">{{ displayUrl }}</span>
    </div>
    <div class="row-actions" v-show="editMode" @click.stop>
      <el-button text size="small" @click="$emit('edit', bookmark)" title="编辑">
        <el-icon :size="13"><Edit /></el-icon>
      </el-button>
      <el-button text size="small" @click="$emit('pin', bookmark)" :title="isPinned ? '取消常驻' : '常驻'">
        <el-icon :size="13"><StarFilled v-if="isPinned" /><Star v-else /></el-icon>
      </el-button>
      <el-button text size="small" @click="$emit('delete', bookmark)" title="删除" type="danger">
        <el-icon :size="13"><Delete /></el-icon>
      </el-button>
    </div>
  </div>
</template>

<script setup>
import { computed } from 'vue'
import { Edit, Star, StarFilled, Delete } from '@element-plus/icons-vue'
import { useBookmarksStore } from '@/stores/bookmarks'

const props = defineProps({
  bookmark: { type: Object, required: true },
  editMode: { type: Boolean, default: false }
})
defineEmits(['edit', 'pin', 'delete'])

const bookmarksStore = useBookmarksStore()
const isPinned = computed(() => bookmarksStore.pinned.some(b => b.id === props.bookmark.id))


function isImageUrl(str) {
  return str && (str.startsWith('http') || str.startsWith('/uploads/'))
}

const displayUrl = computed(() => {
  try {
    const u = new URL(props.bookmark.url)
    return u.hostname + u.pathname.replace(/\/$/, '')
  } catch {
    return props.bookmark.url || ''
  }
})

const rowStyle = computed(() => {
  if (props.bookmark.bg_color) {
    return { '--card-bg': props.bookmark.bg_color }
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
.bookmark-row {
  --card-bg: rgba(255, 255, 255, 0.06);
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 12px 16px;
  border-radius: var(--radius-md);
  position: relative;
  width: 240px;
  height: 68px;
  cursor: pointer;
  transition: box-shadow 0.25s ease, transform 0.25s ease, border-color 0.25s ease;
  border: 1px solid transparent;
  isolation: isolate;
}
.bookmark-row.is-editable {
  cursor: grab;
}
.bookmark-row.is-editable:active {
  cursor: grabbing;
}
.bookmark-row::before {
  content: '';
  position: absolute;
  inset: 0;
  border-radius: inherit;
  background: var(--card-bg);
  opacity: 0.4;
  transition: opacity 0.3s ease;
  z-index: -1;
}
.bookmark-row:hover::before {
  opacity: 1;
}
.bookmark-row:hover {
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.06);
  transform: translateY(-1px);
  border-color: rgba(255, 255, 255, 0.08);
}
.bookmark-row:active {
  transform: translateY(0);
  box-shadow: 0 1px 4px rgba(0, 0, 0, 0.04);
}

.row-icon {
  width: 42px;
  height: 42px;
  border-radius: var(--radius-md);
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 22px;
  background: rgba(255, 255, 255, 0.06);
  flex-shrink: 0;
  transition: transform var(--transition-fast);
}
.bookmark-row:hover .row-icon {
  transform: scale(1.05);
}
.row-icon img {
  width: 26px;
  height: 26px;
  object-fit: contain;
  border-radius: 4px;
}

.row-body {
  flex: 1;
  min-width: 0;
  display: flex;
  flex-direction: column;
  gap: 2px;
}
.row-name {
  font-size: 14px;
  font-weight: 500;
  color: var(--color-text-primary);
  overflow: hidden;
  text-overflow: ellipsis;
  text-shadow: 0 1px 2px rgba(0, 0, 0, 0.3);
  white-space: nowrap;
  flex: 1;
  min-width: 0;
}
.row-url {
  font-size: 11px;
  color: var(--color-text-tertiary);
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.row-actions {
  display: flex;
  align-items: center;
  gap: 1px;
  flex-shrink: 0;
}
.row-actions .el-button {
  width: 22px;
  height: 22px;
  min-height: 22px;
  color: rgba(255, 255, 255, 0.55);
  background: transparent;
  border: none;
}
.row-actions .el-button:hover {
  color: rgba(255, 255, 255, 0.92);
  background: transparent;
}
.row-actions .el-button--danger {
  color: rgba(255, 100, 100, 0.65);
}
.row-actions .el-button--danger:hover {
  color: rgba(255, 100, 100, 0.95);
}
.row-actions .el-button .el-icon {
  vertical-align: middle;
}

@media (max-width: 768px) {
  .bookmark-row {
    width: calc(50% - 3px);
    height: auto;
    min-height: 56px;
    padding: 10px 12px;
    gap: 8px;
  }
  .row-icon {
    width: 34px;
    height: 34px;
    font-size: 18px;
  }
  .row-icon img {
    width: 22px;
    height: 22px;
  }
  .row-name {
    font-size: 13px;
  }
  .row-url {
    font-size: 10px;
  }
}
</style>