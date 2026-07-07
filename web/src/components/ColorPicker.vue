<template>
  <div class="color-picker-wrap">
    <el-color-picker
      v-model="color"
      show-alpha
      :predefine="presetColors"
      @change="onChange"
    >
      <template #default="{ color: currentColor }">
        <div class="color-btn" :style="btnStyle">
          <span class="color-btn-text">{{ displayHex(currentColor) }}</span>
        </div>
      </template>
    </el-color-picker>
  </div>
</template>

<script setup>
import { computed } from 'vue'

const props = defineProps({
  modelValue: { type: String, default: '' }
})
const emit = defineEmits(['update:modelValue'])

const color = computed({
  get: () => props.modelValue || null,
  set: (val) => emit('update:modelValue', val || '')
})

function onChange(val) {
  emit('update:modelValue', val || '')
}

function displayHex(c) {
  if (!c) return '#FFFFFF'
  if (c === 'rgba(0,0,0,0)') return '#00000000'
  return c.toUpperCase()
}

function getContrastColor(hex) {
  if (!hex || hex === 'rgba(0,0,0,0)') return 'rgba(0,0,0,0.6)'
  const h = hex.replace('#', '')
  if (h.length < 6) return 'rgba(0,0,0,0.6)'
  const r = parseInt(h.slice(0, 2), 16)
  const g = parseInt(h.slice(2, 4), 16)
  const b = parseInt(h.slice(4, 6), 16)
  const luminance = (0.299 * r + 0.587 * g + 0.114 * b) / 255
  return luminance > 0.5 ? 'rgba(0,0,0,0.6)' : 'rgba(255,255,255,0.85)'
}

const btnStyle = computed(() => {
  if (!props.modelValue) return { backgroundColor: '#fff' }
  if (props.modelValue === 'rgba(0,0,0,0)') {
    return {
      background: 'conic-gradient(#e8e8ec 25%, #fff 25%, #fff 50%, #e8e8ec 50%, #e8e8ec 75%, #fff 75%)'
    }
  }
  return {
    backgroundColor: props.modelValue,
    color: getContrastColor(props.modelValue)
  }
})

const presetColors = [
  'rgba(0,0,0,0)',
  '#ffffff',
  '#f5f5f7',
  '#e8e8ec',
  '#aeaeb2',
  '#1d1d1f',
  '#ef4444',
  '#f472b6',
  '#f9a8d4',
  '#a855f7',
  '#6366f1',
  '#3b82f6',
  '#38bdf8',
  '#06b6d4',
  '#22c55e',
  '#34d399',
  '#84cc16',
  '#eab308',
  '#f59e0b',
  '#f97316',
  '#a16207',
]
</script>

<style scoped>
.color-picker-wrap {
  display: inline-flex;
  display: flex;
  width: 100%;
}
.color-btn {
  width: 60px;
  height: 36px;
  border-radius: var(--radius-sm);
  border: none;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: opacity var(--transition-fast);
  padding: 0;
  background: #fff;
}
.color-btn:hover {
  opacity: 0.85;
}
.color-btn-text {
  font-size: 11px;
  font-weight: 500;
  font-family: monospace;
  color: inherit;
}

:deep(.el-color-picker) {
  width: 100% !important;
}
:deep(.el-color-picker__trigger) {
  width: 100% !important;
  padding: 0 !important;
  border: none !important;
  border-radius: var(--radius-sm) !important;
  box-shadow: none !important;
}
</style>