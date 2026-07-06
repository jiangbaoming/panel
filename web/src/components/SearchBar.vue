<template>
  <div class="search-bar">
    <div class="search-wrap" :class="{ 'search-has-suggest': searchQuery.trim() && !hasMatch }">
      <!-- 搜索引擎选择 -->
      <el-dropdown trigger="click" @command="selectEngine">
        <div class="search-engine-btn">
          <span class="engine-icon">{{ currentEngine.icon }}</span>
        </div>
        <template #dropdown>
          <el-dropdown-menu>
            <el-dropdown-item
              v-for="eng in engines"
              :key="eng.key"
              :command="eng.key"
              :class="{ 'is-selected': selectedEngine === eng.key }"
            >
              <span style="margin-right:6px">{{ eng.icon }}</span>
              {{ eng.name }}
            </el-dropdown-item>
          </el-dropdown-menu>
        </template>
      </el-dropdown>

      <input
        ref="inputRef"
        v-model="searchQuery"
        class="search-input"
        placeholder="搜索书签..."
        type="text"
        @keydown.enter="doWebSearch"
      />
      <kbd v-if="!searchQuery" class="search-hint">Ctrl + K</kbd>
      <el-button v-else text size="small" class="search-clear" @click="searchQuery = ''">
        <el-icon><Close /></el-icon>
      </el-button>
    </div>

    <!-- 联网搜索建议 -->
    <div v-if="searchQuery.trim() && !hasMatch" class="web-search-suggest">
      <div
        v-for="eng in engines"
        :key="eng.key"
        class="suggest-item"
        @click="openSearch(eng)"
      >
        <span class="suggest-icon">{{ eng.icon }}</span>
        <span>在 <b>{{ eng.name }}</b> 中搜索「{{ searchQuery.trim() }}」</span>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { Close } from '@element-plus/icons-vue'

const props = defineProps({
  modelValue: { type: String, default: '' },
  hasMatch: { type: Boolean, default: true }
})

const emit = defineEmits(['update:modelValue'])

const engines = [
  { key: 'google', name: 'Google', icon: '🔍', url: 'https://www.google.com/search?q=' },
  { key: 'bing', name: 'Bing', icon: '🔎', url: 'https://www.bing.com/search?q=' },
  { key: 'baidu', name: '百度', icon: '🌐', url: 'https://www.baidu.com/s?wd=' },
  { key: 'github', name: 'GitHub', icon: '🐙', url: 'https://github.com/search?q=' },
  { key: 'zhihu', name: '知乎', icon: '💡', url: 'https://www.zhihu.com/search?type=content&q=' },
]

const selectedEngine = ref('google')
const currentEngine = computed(() => engines.find(e => e.key === selectedEngine.value) || engines[0])

const inputRef = ref(null)
const searchQuery = computed({
  get: () => props.modelValue,
  set: (v) => emit('update:modelValue', v)
})

function selectEngine(key) {
  selectedEngine.value = key
}

function doWebSearch() {
  const q = searchQuery.value.trim()
  if (!q) return
  openSearch(currentEngine.value)
}

function openSearch(eng) {
  window.open(eng.url + encodeURIComponent(searchQuery.value.trim()), '_blank')
}

// Ctrl+K 聚焦
function handleKeydown(e) {
  if ((e.ctrlKey || e.metaKey) && e.key === 'k') {
    e.preventDefault()
    inputRef.value?.focus()
  }
}
onMounted(() => document.addEventListener('keydown', handleKeydown))
onUnmounted(() => document.removeEventListener('keydown', handleKeydown))
</script>

<style scoped>
.search-bar {
  max-width: 600px;
  margin: 0 auto 40px;
  animation: searchSlideIn 0.5s ease 0.1s both;
}

@keyframes searchSlideIn {
  from {
    opacity: 0;
    transform: translateY(8px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

.search-wrap {
  display: flex;
  align-items: center;
  padding: 12px 18px;
  background: var(--color-surface);
  border-radius: var(--radius-lg);
  box-shadow: var(--shadow-md);
  transition: box-shadow var(--transition-normal), transform var(--transition-normal), border-radius var(--transition-normal);
  border: 1px solid transparent;
}

.search-wrap:focus-within {
  box-shadow: 0 4px 24px rgba(64, 158, 255, 0.15);
  transform: translateY(-1px);
  border-color: rgba(64, 158, 255, 0.2);
}

.search-wrap.search-has-suggest {
  border-radius: var(--radius-lg) var(--radius-lg) 0 0;
  border-bottom-color: transparent;
}

/* 搜索引擎按钮 */
.search-engine-btn {
  display: flex;
  align-items: center;
  padding: 4px 10px 4px 6px;
  cursor: pointer;
  border-radius: var(--radius-sm);
  transition: background var(--transition-fast);
  margin-right: 10px;
  flex-shrink: 0;
}
.search-engine-btn:hover {
  background: var(--color-surface-hover);
}
.engine-icon {
  font-size: 18px;
  line-height: 1;
}

.search-input {
  flex: 1;
  border: none;
  outline: none;
  font-size: 16px;
  color: #1d1d1f;
  font-family: inherit;
  background: transparent;
}
.search-input::placeholder {
  color: #aaa;
}

.search-hint {
  font-size: 11px;
  padding: 3px 8px;
  background: #f5f5f7;
  border-radius: 6px;
  color: #999;
  font-family: inherit;
  flex-shrink: 0;
  border: 1px solid rgba(0, 0, 0, 0.06);
}

.search-clear {
  flex-shrink: 0;
  color: #999;
  transition: color var(--transition-fast);
}
.search-clear:hover {
  color: #666;
}

/* 联网搜索建议 */
.web-search-suggest {
  background: var(--color-surface);
  border-radius: 0 0 var(--radius-lg) var(--radius-lg);
  box-shadow: 0 8px 24px rgba(0, 0, 0, 0.08);
  overflow: hidden;
  padding: 4px 0;
  border: 1px solid rgba(0, 0, 0, 0.04);
  border-top: none;
}
.suggest-item {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 10px 20px;
  cursor: pointer;
  font-size: 14px;
  color: #1d1d1f;
  transition: background var(--transition-fast);
}
.suggest-item:hover {
  background: var(--color-surface-hover);
}
.suggest-icon {
  font-size: 16px;
}

@media (max-width: 768px) {
  .search-bar {
    max-width: 100%;
    margin: 0 auto 24px;
  }
  .search-wrap {
    padding: 10px 14px;
  }
  .search-input {
    font-size: 14px;
  }
  .search-engine-btn {
    padding: 2px 6px 2px 4px;
    margin-right: 6px;
  }
  .engine-icon {
    font-size: 16px;
  }
  .search-hint {
    display: none;
  }
}
</style>