<template>
  <div v-if="bookmarksStore.pinned.length > 0 || editMode" class="dock">
    <el-button
      v-show="canScrollLeft"
      class="dock-scroll-btn dock-scroll-left"
      @click="scrollDock(-1)"
      aria-label="向左滚动"
      >‹</el-button
    >
    <div
      class="dock-viewport"
      ref="viewportRef"
      @touchstart="onTouchStart"
      @touchmove="onTouchMove"
      @touchend="onTouchEnd">
      <div class="dock-inner" :style="dockInnerStyle">
        <el-tooltip
          v-for="bm in bookmarksStore.pinned"
          :key="bm.id"
          :content="bm.name"
          placement="top"
          :show-after="100"
          :hide-after="0">
          <div
            class="dock-item"
            @click="openUrl(bm.url)"
            @mouseenter="hoverId = bm.id"
            @mouseleave="hoverId = null">
            <el-icon
              v-if="editMode"
              @click.stop="removePin(bm)"
              title="取消常驻"
              class="dock-remove"
              ><CircleClose
            /></el-icon>
            <img
              v-if="isImageUrl(bm.icon)"
              :src="bm.icon"
              class="dock-icon-img" />
            <span v-else class="dock-icon-emoji">{{ bm.icon || "🔗" }}</span>
            <span class="dock-dot" v-show="hoverId === bm.id"></span>
          </div>
        </el-tooltip>
        <div
          v-if="editMode"
          class="dock-item dock-add"
          @click="openAdd"
          title="添加常驻书签">
          <el-icon :size="20" class="dock-add-icon"><Plus /></el-icon>
        </div>
      </div>
    </div>
    <el-button
      v-show="canScrollRight"
      class="dock-scroll-btn dock-scroll-right"
      @click="scrollDock(1)"
      aria-label="向右滚动"
      >›</el-button
    >

    <el-dialog
      v-model="pickerVisible"
      title="选择常驻书签"
      width="420px"
      :close-on-click-modal="false"
      append-to-body>
      <el-input
        v-model="searchKeyword"
        placeholder="搜索书签..."
        clearable
        class="picker-search" />
      <div class="picker-list" v-if="filteredBookmarks.length > 0">
        <div
          v-for="bm in filteredBookmarks"
          :key="bm.id"
          class="picker-item"
          @click="doPin(bm)">
          <img
            v-if="isImageUrl(bm.icon)"
            :src="bm.icon"
            class="picker-icon-img" />
          <span v-else class="picker-icon-emoji">{{ bm.icon || "🔗" }}</span>
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
import { ref, computed, onMounted, onUpdated, nextTick, watch } from "vue";
import { ElMessage, ElMessageBox } from "element-plus";
import { Plus } from "@element-plus/icons-vue";
import { useBookmarksStore } from "@/stores/bookmarks";

const props = defineProps({
  editMode: { type: Boolean, default: false },
});

const bookmarksStore = useBookmarksStore();
const hoverId = ref(null);
const pickerVisible = ref(false);
const searchKeyword = ref("");
const viewportRef = ref(null);
const canScrollLeft = ref(false);
const canScrollRight = ref(false);
const scrollOffset = ref(0);
const maxScrollOffset = ref(0);

const unpinnedBookmarks = computed(() => {
  const pinnedIds = new Set(bookmarksStore.pinned.map((b) => b.id));
  const result = [];
  for (const g of bookmarksStore.groups) {
    for (const bm of g.bookmarks || []) {
      if (!pinnedIds.has(bm.id)) {
        result.push({ ...bm, groupName: g.name });
      }
    }
  }
  return result;
});

const filteredBookmarks = computed(() => {
  const kw = searchKeyword.value.trim().toLowerCase();
  if (!kw) return unpinnedBookmarks.value;
  return unpinnedBookmarks.value.filter(
    (bm) =>
      bm.name.toLowerCase().includes(kw) ||
      bm.groupName.toLowerCase().includes(kw),
  );
});

function isImageUrl(str) {
  return str && (str.startsWith("http") || str.startsWith("/uploads/"));
}
function openUrl(url) {
  if (!props.editMode && url) {
    const a = document.createElement('a')
    a.href = url
    a.target = '_blank'
    a.rel = 'noopener noreferrer'
    a.click()
  }
}

function openAdd() {
  searchKeyword.value = "";
  pickerVisible.value = true;
}

async function doPin(bm) {
  try {
    await bookmarksStore.togglePinned(bm.id, true);
    ElMessage.success("已常驻");
    pickerVisible.value = false;
  } catch {}
}

async function removePin(bm) {
  try {
    await ElMessageBox.confirm(`确定取消常驻「${bm.name}」？`, "提示", {
      type: "warning",
    });
    await bookmarksStore.togglePinned(bm.id, false);
    ElMessage.success("已取消常驻");
  } catch {}
}

function updateScrollState() {
  const el = viewportRef.value;
  if (!el || !el.firstElementChild) return;
  const innerWidth = el.firstElementChild.scrollWidth;
  const overflow = innerWidth - el.clientWidth;
  maxScrollOffset.value = Math.max(0, overflow);
  canScrollLeft.value = scrollOffset.value < -2;
  canScrollRight.value = scrollOffset.value > -(maxScrollOffset.value - 2);
}

function scrollDock(dir) {
  const step = 100;
  const next = scrollOffset.value - dir * step;
  scrollOffset.value = Math.max(-maxScrollOffset.value, Math.min(0, next));
}

const touchStartX = ref(0);
const touchStartOffset = ref(0);
const isTouching = ref(false);

const dockInnerStyle = computed(() => ({
  transform: `translateX(${scrollOffset.value}px)`,
  transition: isTouching.value ? "none" : "transform 0.3s ease",
}));

function onTouchStart(e) {
  touchStartX.value = e.touches[0].clientX;
  touchStartOffset.value = scrollOffset.value;
  isTouching.value = true;
}

function onTouchMove(e) {
  const dx = e.touches[0].clientX - touchStartX.value;
  const next = touchStartOffset.value + dx;
  scrollOffset.value = Math.max(-maxScrollOffset.value, Math.min(0, next));
}

function onTouchEnd() {
  isTouching.value = false;
}

onMounted(() => nextTick(updateScrollState));
onUpdated(() => nextTick(updateScrollState));
watch(
  () => bookmarksStore.pinned.length,
  () => {
    scrollOffset.value = 0;
    nextTick(updateScrollState);
  },
);
</script>

<style scoped>
.dock {
  display: flex;
  align-items: flex-end;
  justify-content: center;
  margin: 32px 0 24px;
  animation: dockSlideUp 0.6s cubic-bezier(0.16, 1, 0.3, 1) 0.2s both;
  perspective: 1200px;
  max-width: 100%;
  overflow: hidden;
  position: relative;
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

.dock-viewport {
  overflow: hidden;
  flex: 1 1 auto;
  min-width: 0;
}

.dock-inner {
  display: flex;
  align-items: flex-end;
  justify-content: center;
  gap: 4px;
  padding: 10px 16px;
  flex-wrap: nowrap;
}

.dock-scroll-btn {
  flex-shrink: 0;
  width: 28px;
  height: 28px;
  border: none;
  border-radius: 50%;
  background: rgba(255, 255, 255, 0.12);
  color: rgba(255, 255, 255, 0.6);
  font-size: 18px;
  line-height: 1;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  transition:
    background 0.15s,
    color 0.15s;
  margin-bottom: 10px;
  -webkit-tap-highlight-color: transparent;
}
.dock-scroll-btn:hover {
  background: rgba(255, 255, 255, 0.22);
  color: #fff;
}

/* el-tooltip 包裹层不做弹性压缩 */
.el-tooltip__trigger {
  flex-shrink: 0;
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

.dock-icon-img {
  width: 44px;
  height: 44px;
  object-fit: contain;
  border-radius: 8px;
  flex-shrink: 0;
}
.dock-icon-emoji {
  width: 44px;
  height: 44px;
  font-size: 32px;
  display: flex;
  align-items: center;
  justify-content: center;
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
  top: 8px;
  right: 4px;
  width: 15px;
  height: 15px;
  border-radius: 50%;
  color: #e91515;
  font-size: 12px;
  line-height: 18px;
  text-align: center;
  align-items: center;
  justify-content: center;
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
  .dock-add {
    width: 32px;
    height: 32px;
  }
  .picker-list {
    max-height: 240px;
  }
}
</style>