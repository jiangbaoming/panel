<template>
  <teleport to="body">
    <transition name="search-fade">
      <div v-if="visible" class="search-overlay" @click.self="close">
        <div class="search-panel" @keydown="handleKeydown">
          <div class="search-input-wrap">
            <el-icon :size="22" class="search-input-icon"><Search /></el-icon>
            <input
              ref="inputRef"
              v-model="query"
              class="search-input"
              placeholder="搜索书签..."
              @input="onInput"
            />
            <kbd class="search-hint">ESC</kbd>
          </div>
          <div v-if="query" class="search-results" v-loading="loading">
            <div
              v-for="(r, i) in results"
              :key="r.id"
              class="search-item"
              :class="{ active: activeIndex === i }"
              @click="openItem(r)"
              @mouseenter="activeIndex = i"
            >
              <span class="search-item-icon">{{ r.icon || '🔗' }}</span>
              <div class="search-item-info">
                <span class="search-item-name">{{ r.name }}</span>
                <span class="search-item-url">{{ r.url }}</span>
              </div>
              <span class="search-item-group">{{ r.groupName }}</span>
              <el-icon class="search-item-open"><TopRight /></el-icon>
            </div>
            <div v-if="results.length === 0 && !loading" class="search-empty">
              未找到匹配的书签
            </div>
          </div>
        </div>
      </div>
    </transition>
  </teleport>
</template>

<script setup>
import { ref, watch, nextTick } from 'vue'
import { searchBookmarks } from '@/api'
import { Search, TopRight } from '@element-plus/icons-vue'

const props = defineProps({ visible: Boolean })
const emit = defineEmits(['update:visible'])

const query = ref('')
const results = ref([])
const loading = ref(false)
const activeIndex = ref(-1)
const inputRef = ref(null)

let debounceTimer = null

watch(() => props.visible, (v) => {
  if (v) {
    query.value = ''
    results.value = []
    activeIndex.value = -1
    nextTick(() => inputRef.value?.focus())
  }
})

function onInput() {
  clearTimeout(debounceTimer)
  if (!query.value.trim()) { results.value = []; return }
  debounceTimer = setTimeout(doSearch, 300)
}

async function doSearch() {
  const q = query.value.trim()
  if (!q) return
  loading.value = true
  try {
    results.value = await searchBookmarks(q)
    activeIndex.value = results.value.length > 0 ? 0 : -1
  } catch (e) {
    results.value = []
  } finally {
    loading.value = false
  }
}

function handleKeydown(e) {
  if (e.key === 'ArrowDown') {
    e.preventDefault()
    if (activeIndex.value < results.value.length - 1) activeIndex.value++
  } else if (e.key === 'ArrowUp') {
    e.preventDefault()
    if (activeIndex.value > 0) activeIndex.value--
  } else if (e.key === 'Enter' && activeIndex.value >= 0) {
    openItem(results.value[activeIndex.value])
  } else if (e.key === 'Escape') {
    close()
  }
}

function openItem(item) {
  if (item?.url) {
    window.open(item.url, '_blank')
    close()
  }
}

function close() {
  emit('update:visible', false)
}
</script>

<style scoped>
.search-overlay {
  position: fixed;
  inset: 0;
  background: rgba(0, 0, 0, 0.5);
  backdrop-filter: blur(4px);
  display: flex;
  justify-content: center;
  padding-top: 15vh;
  z-index: 2000;
}
.search-panel {
  width: 580px;
  max-height: 60vh;
  display: flex;
  flex-direction: column;
}
.search-input-wrap {
  display: flex;
  align-items: center;
  background: #fff;
  border-radius: 16px;
  padding: 14px 20px;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.12);
}
.search-input-icon {
  color: #86868b;
  margin-right: 12px;
  flex-shrink: 0;
}
.search-input {
  flex: 1;
  border: none;
  outline: none;
  font-size: 18px;
  background: transparent;
  color: #1d1d1f;
}
.search-input::placeholder { color: #c7c7cc; }
.search-hint {
  font-size: 11px;
  padding: 2px 6px;
  background: #f5f5f7;
  border-radius: 4px;
  color: #86868b;
  font-family: inherit;
  flex-shrink: 0;
}
.search-results {
  background: #fff;
  border-radius: 12px;
  margin-top: 8px;
  overflow-y: auto;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.12);
  max-height: 50vh;
}
.search-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 12px 16px;
  cursor: pointer;
  transition: background 0.15s;
}
.search-item:first-child { border-radius: 12px 12px 0 0; }
.search-item:last-child { border-radius: 0 0 12px 12px; }
.search-item.active, .search-item:hover { background: #f5f5f7; }
.search-item-icon { font-size: 20px; line-height: 1; }
.search-item-info {
  flex: 1;
  min-width: 0;
  display: flex;
  flex-direction: column;
  gap: 2px;
}
.search-item-name { font-size: 14px; font-weight: 500; color: #1d1d1f; }
.search-item-url { font-size: 12px; color: #86868b; overflow: hidden; text-overflow: ellipsis; white-space: nowrap; }
.search-item-group { font-size: 11px; color: #aeaeb2; flex-shrink: 0; }
.search-item-open { color: #c7c7cc; flex-shrink: 0; }
.search-empty { padding: 32px; text-align: center; color: #86868b; font-size: 14px; }

.search-fade-enter-active, .search-fade-leave-active { transition: opacity 0.2s; }
.search-fade-enter-from, .search-fade-leave-to { opacity: 0; }
</style>
