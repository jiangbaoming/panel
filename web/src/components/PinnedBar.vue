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
