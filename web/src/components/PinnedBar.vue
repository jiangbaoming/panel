<template>
  <div v-if="bookmarksStore.pinned.length > 0 || editMode" class="dock">
    <div class="dock-inner">
      <div
        v-for="(bm, idx) in bookmarksStore.pinned"
        :key="bm.id"
        class="dock-item"
        :class="{ 'dock-item--hover': hoverId === bm.id }"
        @click="openUrl(bm.url)"
        @mouseenter="onHover(bm.id, idx)"
        @mouseleave="hoverId = null"
      >
        <span class="dock-tooltip">{{ bm.name }}</span>
        <el-button
          v-if="editMode"
          class="dock-remove"
          @click.stop="removePin(bm)"
          title="取消常驻"
        >×</el-button>
        <img v-if="isImageUrl(bm.icon)" :src="bm.icon" class="dock-icon-img" />
        <span v-else class="dock-icon-emoji">{{ bm.icon || '🔗' }}</span>
        <span class="dock-dot" v-show="hoverId === bm.id"></span>
      </div>
      <div
        v-if="editMode"
        class="dock-item dock-add"
        @click="openAdd"
        title="添加常驻书签"
      >
        <el-icon :size="24" class="dock-add-icon"><Plus /></el-icon>
      </div>
    </div>
    <!-- 选择常驻书签弹窗 -->
    <el-dialog
      v-model="pickerVisible"
      title="选择常驻书签"
      width="420px"
      :close-on-click-modal="false"
      append-to-body
    >
      <el-input
          v-model="searchKeyword"
          placeholder="搜索书签..."
          clearable
          class="picker-search"
        />
        <div class="picker-list" v-if="filteredBookmarks.length > 0">
        <div
          v-for="bm in filteredBookmarks"
          :key="bm.id"
          class="picker-item"
          @click="doPin(bm)"
        >
          <img v-if="isImageUrl(bm.icon)" :src="bm.icon" class="picker-icon-img" />
          <span v-else class="picker-icon-emoji">{{ bm.icon || '🔗' }}</span>
          <div class="picker-info">
            <span class="picker-name">{{ bm.name }}</span>
            <span class="picker-group">{{ bm.groupName }}</span>
          </div>
        </div>
      </div>
      <div v-else class="picker-empty">没有可常驻的书签</div>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Plus } from '@element-plus/icons-vue'
import { useBookmarksStore } from '@/stores/bookmarks'

const props = defineProps({
  editMode: { type: Boolean, default: false }
})

const bookmarksStore = useBookmarksStore()
const hoverId = ref(null)
const hoverIdx = ref(-1)
const pickerVisible = ref(false)
const searchKeyword = ref('')

const unpinnedBookmarks = computed(() => {
  const pinnedIds = new Set(bookmarksStore.pinned.map(b => b.id))
  const result = []
  for (const g of bookmarksStore.groups) {
    for (const bm of (g.bookmarks || [])) {
      if (!pinnedIds.has(bm.id)) {
        result.push({ ...bm, groupName: g.name })
      }
    }
  }
  return result
})

const filteredBookmarks = computed(() => {
  const kw = searchKeyword.value.trim().toLowerCase()
  if (!kw) return unpinnedBookmarks.value
  return unpinnedBookmarks.value.filter(
    bm => bm.name.toLowerCase().includes(kw) || bm.groupName.toLowerCase().includes(kw)
  )
})

function isImageUrl(str) {
  return str && (str.startsWith('http') || str.startsWith('/uploads/'))
}
function openUrl(url) {
  if (url) window.open(url, '_blank')
}
function onHover(id, idx) {
  hoverId.value = id
  hoverIdx.value = idx
}

function openAdd() {
  searchKeyword.value = ''
  pickerVisible.value = true
}

async function doPin(bm) {
  try {
    await bookmarksStore.togglePinned(bm.id, true)
    ElMessage.success('已常驻')
    pickerVisible.value = false
  } catch (e) {
    ElMessage.error('操作失败')
  }
}

async function removePin(bm) {
  try {
    await ElMessageBox.confirm(`确定取消常驻「${bm.name}」？`, '提示', { type: 'warning' })
    await bookmarksStore.togglePinned(bm.id, false)
    ElMessage.success('已取消常驻')
  } catch (e) {}
}
</script>

<style scoped>
.dock {
  display: flex;
  justify-content: center;
  margin: 32px 0 24px;
  animation: dockSlideUp 0.6s cubic-bezier(0.16, 1, 0.3, 1) 0.2s both;
  perspective: 1200px;
}

@keyframes dockSlideUp {
  from {
    opacity: 0;
    transform: translateY(20px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

.dock-inner {
  display: flex;
  align-items: flex-end;
  gap: 4px;
  padding: 10px 16px;
}

.dock-item {
  position: relative;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: flex-end;
  width: 56px;
  height: 56px;
  cursor: pointer;
  transition: transform 0.3s cubic-bezier(0.25, 0.8, 0.25, 1.2);
  transform-origin: bottom center;
  border-radius: 12px;
  flex-shrink: 0;
}

.dock-item:hover {
  transform: scale(1.35);
}

.dock-tooltip {
  position: absolute;
  top: -36px;
  left: 50%;
  transform: translateX(-50%);
  padding: 4px 10px;
  background: rgba(0, 0, 0, 0.75);
  color: #fff;
  font-size: 12px;
  border-radius: 6px;
  white-space: nowrap;
  pointer-events: none;
  opacity: 0;
  transition: opacity 0.2s ease;
}
.dock-tooltip::after {
  content: '';
  position: absolute;
  bottom: -4px;
  left: 50%;
  transform: translateX(-50%);
  width: 0;
  height: 0;
  border-left: 5px solid transparent;
  border-right: 5px solid transparent;
  border-top: 5px solid rgba(0, 0, 0, 0.75);
}
.dock-item:hover .dock-tooltip {
  opacity: 1;
}

.dock-icon-img {
  width: 44px;
  height: 44px;
  object-fit: contain;
  border-radius: 8px;
  flex-shrink: 0;
}

.dock-icon-emoji {
  font-size: 36px;
  line-height: 1;
  flex-shrink: 0;
}

.dock-dot {
  width: 4px;
  height: 4px;
  border-radius: 50%;
  background: rgba(255, 255, 255, 0.6);
  margin-top: 6px;
  flex-shrink: 0;
}

.dock-remove {
  position: absolute;
  top: -4px;
  right: -4px;
  width: 18px;
  height: 18px;
  border-radius: 50%;
  border: none;
  background: rgba(255, 80, 80, 0.85);
  color: #fff;
  font-size: 12px;
  line-height: 18px;
  text-align: center;
  cursor: pointer;
  z-index: 2;
  opacity: 0;
  transition: opacity 0.2s;
}
.dock-item:hover .dock-remove {
  opacity: 1;
}

.dock-add {
  background: none;
  border: 1px dashed rgba(255, 255, 255, 0.3);
  cursor: pointer;
  width: 40px;
  height: 40px;
  line-height: 40px;
  text-align: center;
  align-items: center;
  justify-content: center;
}
.dock-add:hover {
  border-color: rgba(255, 255, 255, 0.6);
}
.dock-add-icon {
  color: rgba(255, 255, 255, 0.5);
}

.picker-search {
  margin-bottom: 12px;
}

.picker-list {
  max-height: 360px;
  overflow-y: auto;
}
.picker-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 10px 12px;
  border-radius: 8px;
  cursor: pointer;
  transition: background 0.15s;
}
.picker-item:hover {
  background: rgba(0, 0, 0, 0.05);
}
.picker-icon-img {
  width: 32px;
  height: 32px;
  object-fit: contain;
  border-radius: 6px;
  flex-shrink: 0;
}
.picker-icon-emoji {
  font-size: 24px;
  line-height: 1;
  flex-shrink: 0;
}
.picker-info {
  display: flex;
  flex-direction: column;
  gap: 2px;
  min-width: 0;
}
.picker-name {
  font-size: 14px;
  color: #333;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}
.picker-group {
  font-size: 12px;
  color: #999;
}
.picker-empty {
  text-align: center;
  color: #999;
  padding: 32px 0;
}

@media (max-width: 768px) {
  .dock-inner {
    gap: 2px;
    padding: 8px 8px;
  }
  .dock-item {
    width: 42px;
    height: 42px;
    border-radius: 10px;
  }
  .dock-icon-img {
    width: 32px;
    height: 32px;
    border-radius: 6px;
  }
  .dock-icon-emoji {
    font-size: 26px;
  }
  .dock-dot {
    width: 3px;
    height: 3px;
    margin-top: 4px;
  }
  .dock-tooltip {
    display: none;
  }
  .dock-add {
    width: 32px;
    height: 32px;
  }
  .picker-list {
    max-height: 240px;
  }
}
</style>