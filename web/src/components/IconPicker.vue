<template>
  <div class="icon-picker" :class="{ 'icon-picker--manage': mode === 'manage' }">

    <!-- 管理模式：分类切换 + 操作按钮 -->
    <div v-if="mode === 'manage'" class="manage-toolbar">
      <el-radio-group v-model="manageCategory" @change="onManageCategoryChange">
        <el-radio-button value="icon">图标</el-radio-button>
        <el-radio-button value="bg">背景图</el-radio-button>
      </el-radio-group>
      <div style="flex:1" />
      <el-upload
        :auto-upload="false"
        :show-file-list="false"
        accept=".zip"
        :on-change="handleZipUpload"
      >
        <el-button :icon="FolderOpened" size="small">上传 ZIP 包</el-button>
      </el-upload>
      <el-button :icon="Delete" size="small" type="danger" text @click="handleClearCategory">清空全部</el-button>
    </div>

    <!-- 选择模式：Emoji / 图库 / URL 页签切换 -->
    <div v-if="mode === 'pick'" class="icon-tabs">
      <el-button v-if="props.category !== 'bg' && !props.hideEmoji" size="small" :type="activeTab === 'emoji' ? 'primary' : 'default'" text @click="switchTab('emoji')">Emoji</el-button>
      <el-button size="small" :type="activeTab === 'gallery' ? 'primary' : 'default'" text @click="switchTab('gallery')">图库</el-button>
      <el-button v-if="mode === 'pick'" size="small" :type="activeTab === 'url' ? 'primary' : 'default'" text @click="switchTab('url')">URL</el-button>
    </div>

    <!-- Emoji 选择 -->
    <div v-if="mode === 'pick' && activeTab === 'emoji'" class="emoji-grid">
      <button
        v-for="emoji in emojiList" :key="emoji"
        type="button" class="emoji-item"
        :class="{ active: modelValue === emoji }"
        @click="emit('update:modelValue', emoji)"
      >{{ emoji }}</button>
    </div>

    <!-- 图库（选择模式 + 管理模式共用） -->
    <div v-if="mode === 'manage' || activeTab === 'gallery'" class="gallery-tab">
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
            <span v-if="mode === 'manage'" class="gal-upload-label">上传图片</span>
          </div>
        </el-upload>
        <div
          v-for="img in images" :key="img.id"
          class="gal-item"
          :class="{ selected: mode === 'pick' && modelValue === img.url }"
          @click="mode === 'pick' && emit('update:modelValue', img.url)"
        >
          <img :src="img.url" loading="lazy" />
          <template v-if="mode === 'manage'">
            <button class="gal-item-del" @click.stop="removeImg(img)"><el-icon :size="12"><Close /></el-icon></button>
            <div class="gal-item-name">{{ img.original_name }}</div>
          </template>
        </div>
      </div>
    </div>

    <!-- URL 输入 -->
    <div v-if="mode === 'pick' && activeTab === 'url'" class="url-tab">
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
import { Upload, Search, FolderOpened, Delete, Close } from '@element-plus/icons-vue'
import { getImages, uploadImage, uploadZip, deleteImage, clearImages } from '@/api'
import { ElMessage, ElMessageBox } from 'element-plus'

const props = defineProps({
  modelValue: { type: String, default: '' },
  defaultIcon: { type: String, default: '📁' },
  category: { type: String, default: 'icon' },
  hideEmoji: { type: Boolean, default: false },
  mode: { type: String, default: 'pick' }
})
const emit = defineEmits(['update:modelValue', 'refresh'])

const activeTab = ref(props.hideEmoji ? 'gallery' : 'emoji')
const urlInput = ref('')
const urlPreview = ref('')
const images = ref([])
const loading = ref(false)
const searchQuery = ref('')
const page = ref(1)
const pageSize = ref(24)
const total = ref(0)
const hasMore = ref(true)
const manageCategory = ref('icon')
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
  '⚡', '🛡️', '🎈', '📌', '🧩', '🎪', '🏆', '🎭',
  '🏠', '🏢', '🏖️', '🏕️', '🏔️', '🏗️', '🏡', '🏜️',
  '🎄', '🎃', '🎁', '🎀', '🎉', '🎊', '🎋', '🎍',
  '🐶', '🐱', '🦊', '🐼', '🐨', '🐯', '🦁', '🐸',
  '🐧', '🐤', '🦄', '🐝', '🐞', '🦋', '🐠', '🐬',
  '🌈', '🌟', '💫', '⭐', '🌙', '☀️', '⛅', '🌧️',
  '🌊', '🔥', '💧', '🌿', '🍀', '🌻', '🌹', '🌸',
  '🌺', '🍎', '🍊', '🍋', '🍉', '🍇', '🍒', '🍕',
  '🎧', '🎤', '🎸', '🎹', '🎺', '🎻', '🥁', '🎷',
  '⚽', '🏀', '🏈', '⚾', '🎾', '🏐', '🎱', '🏓',
  '🚗', '🚲', '✈️', '🚢', '🚁', '🛸', '🚂', '🚌',
  '🔑', '🔒', '🔓', '🔨', '🔧', '🔩', '🔪', '💣',
  '💊', '💉', '🩺', '🩹', '🧬', '🔬', '🔭', '📡',
  '❤️', '💙', '💜', '💛', '🧡', '🖤', '🤍', '🤎',
  '💯', '✅', '❌', '➕', '➖', '❓', '❗', '💤',
  '🕐', '🕑', '🕒', '🕓', '🕔', '🕕', '🕖', '🕗',
  '♻️', '🔄', '🔃', '↗️', '⬆️', '⬇️', '⏩', '⏪',
  '🎓', '🏫', '📚', '📐', '📏', '📎', '🖇️', '📌',
  '🗂️', '🗃️', '🗄️', '🗑️', '🖊️', '🖋️', '🖌️', '🖍️',
  '💼', '💰', '💳', '💵', '💴', '💶', '💷', '🏦',
  '📈', '📉', '📊', '📋', '📌', '📎', '🖇️', '🖊️',
  '🍔', '🍟', '🍕', '🌭', '🍿', '🧃', '🍩', '🍪',
  '🎮', '🕹️', '🎲', '🎯', '🧩', '♟️', '🎰', '🎳',
  '🌍', '🌎', '🌏', '🗺️', '🏴', '🏳️', '🏳️‍🌈', '🇨🇳',
  '🔇', '🔈', '🔉', '🔊', '📢', '📣', '📯', '🔔',
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
  if (props.mode === 'manage') {
    fetchImages()
    return
  }
  if (props.category === 'bg') {
    activeTab.value = 'gallery'
  }
  if (activeTab.value === 'gallery') {
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
    const cat = props.mode === 'manage' ? manageCategory.value : props.category
    const res = await getImages({ category: cat, q: searchQuery.value || undefined, page: page.value, pageSize: pageSize.value })
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
  const cat = props.mode === 'manage' ? manageCategory.value : props.category
  formData.append('category', cat)
  try {
    await uploadImage(formData)
    ElMessage.success('上传成功')
    page.value = 1
    fetchImages()
  } catch {}
}

async function handleZipUpload(file) {
  const formData = new FormData()
  formData.append('file', file.raw || file)
  formData.append('category', manageCategory.value)
  try {
    const res = await uploadZip(formData)
    ElMessage.success(`导入 ${res.imported || 0} 个文件`)
    page.value = 1
    fetchImages()
  } catch {}
}

async function removeImg(img) {
  try {
    await ElMessageBox.confirm(`确定删除「${img.original_name}」？`, '提示', { type: 'warning' })
    await deleteImage(img.id)
    ElMessage.success('已删除')
    fetchImages()
  } catch {}
}

async function handleClearCategory() {
  const label = manageCategory.value === 'icon' ? '图标' : '背景图'
  try {
    await ElMessageBox.confirm(
      `确定清空所有「${label}」图片？此操作不可恢复。`,
      '清空确认',
      { type: 'warning', confirmButtonText: '确认清空', cancelButtonText: '取消' }
    )
    await clearImages(manageCategory.value)
    ElMessage.success(`已清空所有${label}`)
    fetchImages()
  } catch {}
}

function onManageCategoryChange() {
  page.value = 1
  fetchImages()
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

/* 管理模式 */
.icon-picker--manage {
  height: 390px;
}
.manage-toolbar {
  display: flex;
  align-items: center;
  gap: 8px;
  flex-shrink: 0;
  padding: 8px 8px 0;
}
.gal-upload-label {
  display: none;
}
.icon-picker--manage .gal-upload-trigger {
  width: 80px;
  height: 80px;
  flex-direction: column;
  gap: 4px;
}
.icon-picker--manage .gal-upload-label {
  display: block;
  font-size: 11px;
  color: var(--color-text-tertiary);
}
.icon-picker--manage .gal-item {
  position: relative;
  width: 80px;
  height: auto;
  min-height: 80px;
  border: 1px solid #e8e8ec;
  cursor: default;
  flex-direction: column;
  overflow: visible;
}
.icon-picker--manage .gal-item img {
  width: 100%;
  height: 56px;
  object-fit: cover;
  display: block;
  flex-shrink: 0;
  border-radius: var(--radius-sm) var(--radius-sm) 0 0;
}
.icon-picker--manage .gal-item:hover {
  border-color: var(--color-primary);
  transform: none;
  box-shadow: 0 2px 8px rgba(0,0,0,0.08);
}
.gal-item-del {
  position: absolute;
  top: 2px;
  right: 2px;
  width: 18px;
  height: 18px;
  border-radius: 50%;
  border: none;
  background: rgba(255, 80, 80, 0.9);
  color: #fff;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  opacity: 0;
  transition: opacity 0.15s;
  z-index: 2;
  padding: 0;
  line-height: 1;
}
.gal-item:hover .gal-item-del {
  opacity: 1;
}
.gal-item-name {
  font-size: 11px;
  color: #999;
  text-align: center;
  padding: 2px 4px 4px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  max-width: 100%;
  box-sizing: border-box;
}
</style>