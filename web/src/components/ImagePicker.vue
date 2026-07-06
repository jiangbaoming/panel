<template>
  <el-dialog title="选择图片" v-model="dialogVisible" width="600px">
    <div class="img-picker">
      <div class="img-picker-toolbar">
        <el-upload
          :auto-upload="false"
          :show-file-list="false"
          accept="image/*"
          :on-change="handleUpload"
        >
          <el-button :icon="Upload" size="small">上传图片</el-button>
        </el-upload>
      </div>
      <div v-loading="loading" class="img-grid">
        <div
          v-for="img in images"
          :key="img.id"
          class="img-item"
          :class="{ selected: selected === img.url }"
          @click="select(img.url)"
        >
          <img :src="img.url" :title="img.original_name" />
        </div>
        <el-empty v-if="!loading && images.length === 0" description="暂无图片" />
      </div>
    </div>
    <template #footer>
      <el-button @click="dialogVisible = false">取消</el-button>
      <el-button type="primary" :disabled="!selected" @click="confirm">选择</el-button>
    </template>
  </el-dialog>
</template>

<script setup>
import { ref, watch, computed } from 'vue'
import { getImages, uploadImage } from '@/api'
import { ElMessage } from 'element-plus'
import { Upload } from '@element-plus/icons-vue'

const props = defineProps({ visible: Boolean })
const emit = defineEmits(['update:visible', 'select'])

const dialogVisible = computed({
  get: () => props.visible,
  set: (v) => emit('update:visible', v)
})

const images = ref([])
const loading = ref(false)
const selected = ref('')

watch(() => props.visible, (v) => {
  if (v) { selected.value = ''; fetchImages() }
})

async function fetchImages() {
  loading.value = true
  try {
    const res = await getImages({ pageSize: 200 })
    images.value = res.data || []
  } finally {
    loading.value = false
  }
}

async function handleUpload(file) {
  const formData = new FormData()
  formData.append('file', file.raw || file)
  try {
    const res = await uploadImage(formData)
    ElMessage.success('上传成功')
    fetchImages()
  } catch (e) {
    ElMessage.error('上传失败')
  }
}

function select(url) {
  selected.value = url
}

function confirm() {
  if (selected.value) {
    emit('select', selected.value)
  }
  dialogVisible.value = false
}
</script>

<style scoped>
.img-grid {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
  max-height: 360px;
  overflow-y: auto;
  margin-top: 12px;
}
.img-item {
  width: 80px;
  height: 80px;
  border-radius: 8px;
  overflow: hidden;
  cursor: pointer;
  border: 2px solid transparent;
  transition: border-color 0.2s;
}
.img-item:hover { border-color: #6366f1; }
.img-item.selected { border-color: #6366f1; }
.img-item img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}
</style>
