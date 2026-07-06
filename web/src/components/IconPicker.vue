<template>
  <div class="icon-picker">

    <!-- 模式切换（背景图模式下隐藏 Emoji） -->
    <div class="icon-tabs">
      <el-button v-if="props.category !== 'bg'" size="small" :type="activeTab === 'emoji' ? 'primary' : 'default'" text @click="switchTab('emoji')">Emoji</el-button>
      <el-button size="small" :type="activeTab === 'gallery' ? 'primary' : 'default'" text @click="switchTab('gallery')">图库</el-button>
      <el-button size="small" :type="activeTab === 'url' ? 'primary' : 'default'" text @click="switchTab('url')">URL</el-button>
    </div>

    <!-- Emoji 选择 -->
    <div v-if="activeTab === 'emoji'" class="emoji-grid">
      <button
        v-for="emoji in emojiList" :key="emoji"
        type="button" class="emoji-item"
        :class="{ active: modelValue === emoji }"
        @click="emit('update:modelValue', emoji)"
      >{{ emoji }}</button>
    </div>

    <!-- 图库选择（内嵌） -->
    <div v-if="activeTab === 'gallery'" class="gallery-tab">
      <div class="gal-toolbar">
        <el-input
          v-model="searchQuery" placeholder="搜索..."
          clearable :prefix-icon="Search" size="small" style="flex:1"
          @input="onSearchInput"
        />
      </div>
      <div v-loading="loading" class="gal-grid" @scroll="onGalScroll">
        <el-upload :auto-upload="false" :show-file-list="false" accept="image/*" :on-change="handleUpload">
          <div class="gal-upload-trigger">
            <el-icon :size="22"><Upload /></el-icon>
          </div>
        </el-upload>
        <div
          v-for="img in images" :key="img.id"
          class="gal-item" :class="{ selected: modelValue === img.url }"
          @click="emit('update:modelValue', img.url)"
        >
          <img :src="img.url" loading="lazy" />
        </div>
      </div>
    </div>

    <!-- URL 输入 -->
    <div v-if="activeTab === 'url'" class="url-tab">
      <el-input v-model="urlInput" placeholder="输入图标图片 URL" size="small" clearable @input="urlPreview = ''">
        <template #append>
          <el-button @click="confirmUrl">确认</el-button>
        </template>
      </el-input>
      <div v-if="urlPreview" class="url-preview"><img :src="urlPreview" @error="urlPreview = ''" /></div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { Upload, Search } from '@element-plus/icons-vue'
import { getImages, uploadImage } from '@/api'
import { ElMessage } from 'element-plus'

const props = defineProps({
  modelValue: { type: String, default: '' },
  defaultIcon: { type: String, default: '📁' },
  category: { type: String, default: 'icon' }
})
const emit = defineEmits(['update:modelValue'])

const activeTab = ref('emoji')
const urlInput = ref('')
const urlPreview = ref('')
const images = ref([])
const loading = ref(false)
const searchQuery = ref('')
const page = ref(1)
const pageSize = ref(24)
const total = ref(0)
const hasMore = ref(true)
let searchTimer = null

function isImageUrl(str) {
  return str && (str.startsWith('http') || str.startsWith('/uploads/'))
}

const emojiList = [
  '📁', '🔗', '🛠️', '💬', '🎬', '⚙️', '🎨', '⏱️',
  '📝', '💚', '🟢', '🐳', '⛏️', '📘', '🧑‍💻', '📰',
  '🚀', '🏀', '✨', '📌', '🎯', '📷', '🔤', '📋',
  '📄', '📗', '📑', '✏️', '📊', '🖼️', '🌍', '☁️',
  '🐙', '🔍', '🔎', '🌐', '💡', '📺', '📱', '📕',
  '🎵', '🎶', '🎥', '▶️', '🤖', '🧠', '📖', '💻',
  '📦', '👤', '⭐', '🔥', '💎', '🎮', '✉️', '🔔',
  '⚡', '🛡️', '🎈', '📌', '🧩', '🎪', '🏆', '🎭'
]

function switchTab(tab) {
  activeTab.value = tab
  if (tab === 'gallery') fetchImages()
  if (tab === 'url' && isImageUrl(props.modelValue)) {
    urlInput.value = props.modelValue
    urlPreview.value = props.modelValue
  }
}

onMounted(() => {
  if (props.category === 'bg') {
    activeTab.value = 'gallery'
    fetchImages()
  }
})

function confirmUrl() {
  const v = urlInput.value.trim()
  if (v) { emit('update:modelValue', v); urlPreview.value = v }
}

function onSearchInput() {
  clearTimeout(searchTimer)
  searchTimer = setTimeout(() => { page.value = 1; fetchImages() }, 300)
}

async function fetchImages(append = false) {
  loading.value = true
  try {
    const res = await getImages({ category: props.category, q: searchQuery.value || undefined, page: page.value, pageSize: pageSize.value })
    if (append) {
      images.value = [...images.value, ...(res.data || [])]
    } else {
      images.value = res.data || []
    }
    total.value = res.total || 0
    hasMore.value = images.value.length < total.value
  } finally { loading.value = false }
}

function onGalScroll(e) {
  const el = e.target
  if (el.scrollTop + el.clientHeight >= el.scrollHeight - 30 && hasMore.value && !loading.value) {
    page.value++
    fetchImages(true)
  }
}

async function handleUpload(file) {
  const formData = new FormData()
  formData.append('file', file.raw || file)
  formData.append('category', props.category)
  try {
    await uploadImage(formData)
    ElMessage.success('上传成功')
    fetchImages()
  } catch (e) { ElMessage.error('上传失败') }
}
</script>

<style scoped>
.icon-picker {
  display: flex;
  flex-direction: column;
  gap: 6px;
  width: 100%;
  height: 320px;
  border: 1px solid #e8e8ec;
  border-radius: var(--radius-md);
  box-sizing: border-box;
}

/* 模式切换 */
.icon-tabs {
  display: flex;
  gap: 2px;
  background: #f5f5f7;
  border-radius: var(--radius-sm);
  padding: 3px;
  margin: 6px 8px 0;
  flex-shrink: 0;
}
.icon-tabs .el-button {
  flex: 1;
  border-radius: 6px;
  font-size: 13px;
  padding: 4px 0;
  height: 30px;
  justify-content: center;
}

/* emoji 网格 */
.emoji-grid {
  display: flex;
  flex-wrap: wrap;
  gap: 5px;
  flex: 1;
  overflow-y: auto;
  padding: 0 8px 6px;
  align-content: flex-start;
}
.emoji-item {
  width: 38px;
  height: 38px;
  border-radius: var(--radius-sm);
  border: 1px solid transparent;
  background: #f5f5f7;
  font-size: 18px;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all var(--transition-fast);
  padding: 0;
}
.emoji-item:hover {
  background: #e8e8ec;
  transform: scale(1.12);
}
.emoji-item.active {
  background: var(--color-primary-light);
  border-color: var(--color-primary);
}

/* 图库 tab */
.gallery-tab {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 8px;
  min-height: 0;
  padding: 0 8px 6px;
}
.gal-toolbar {
  display: flex;
  align-items: center;
  flex-shrink: 0;
}
.gal-grid {
  display: flex;
  flex-wrap: wrap;
  gap: 6px;
  overflow-y: auto;
  flex: 1;
  align-content: flex-start;
}
.gal-upload-trigger {
  width: 72px;
  height: 72px;
  border-radius: var(--radius-sm);
  border: 1px dashed #d4d4d8;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  background: #f5f5f7;
  color: var(--color-text-tertiary);
  transition: border-color var(--transition-fast), background var(--transition-fast), color var(--transition-fast);
  flex-shrink: 0;
}
.gal-upload-trigger:hover {
  border-color: var(--color-primary);
  background: var(--color-primary-light);
  color: var(--color-primary);
}
.gal-item {
  width: 72px;
  height: 72px;
  border-radius: var(--radius-sm);
  overflow: hidden;
  cursor: pointer;
  border: 2px solid transparent;
  background: #f5f5f7;
  flex-shrink: 0;
  transition: border-color var(--transition-fast), transform var(--transition-fast);
}
.gal-item:hover {
  border-color: var(--color-primary);
  transform: translateY(-1px);
}
.gal-item.selected {
  border-color: var(--color-primary);
}
.gal-item img {
  width: 100%;
  height: 100%;
  object-fit: cover;
  display: block;
}
.gal-pagination {
  display: flex;
  justify-content: center;
  flex-shrink: 0;
}

/* URL 输入 */
.url-tab {
  display: flex;
  flex-direction: column;
  gap: 8px;
  padding: 0 8px 6px;
}
.url-preview {
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 8px;
}
.url-preview img {
  width: 56px;
  height: 56px;
  object-fit: contain;
  border-radius: var(--radius-sm);
  border: 1px solid #e8e8ec;
  background: #f5f5f7;
}
</style>