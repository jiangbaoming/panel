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
