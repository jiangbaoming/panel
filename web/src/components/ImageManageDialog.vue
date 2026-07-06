<template>
  <el-dialog title="图库管理" v-model="visible" width="700px" :close-on-click-modal="false">
    <div class="img-mgr">
      <div class="img-mgr-toolbar">
        <el-upload
          :auto-upload="false"
          :show-file-list="false"
          accept="image/*"
          multiple
          :on-change="handleUpload"
        >
          <el-button :icon="Upload" type="primary">上传图片</el-button>
        </el-upload>
        <el-upload
          :auto-upload="false"
          :show-file-list="false"
          accept=".zip"
          :on-change="handleZipUpload"
        >
          <el-button :icon="FolderOpened">上传 ZIP 包</el-button>
        </el-upload>
      </div>
      <div v-loading="loading" class="img-grid">
        <div v-for="img in images" :key="img.id" class="img-card">
          <img :src="img.url" />
          <div class="img-card-actions">
            <el-button size="small" @click="copyUrl(img.url)">复制链接</el-button>
            <el-button size="small" type="danger" @click="remove(img)">删除</el-button>
          </div>
          <div class="img-card-name">{{ img.original_name }}</div>
        </div>
        <el-empty v-if="!loading && images.length === 0" description="暂无图片" />
      </div>
    </div>
  </el-dialog>
</template>

<script setup>
import { ref, watch } from 'vue'
import { getImages, uploadImage, uploadZip, deleteImage } from '@/api'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Upload, FolderOpened } from '@element-plus/icons-vue'

const props = defineProps({ visible: Boolean })
const emit = defineEmits(['update:visible'])

const images = ref([])
const loading = ref(false)

watch(() => props.visible, (v) => { if (v) fetchImages() })

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
    await uploadImage(formData)
    ElMessage.success('上传成功')
    fetchImages()
  } catch (e) {
    ElMessage.error('上传失败')
  }
}

async function handleZipUpload(file) {
  const formData = new FormData()
  formData.append('file', file.raw || file)
  try {
    const res = await uploadZip(formData)
    ElMessage.success(`导入 ${res.imported || 0} 个图标`)
    fetchImages()
  } catch (e) {
    ElMessage.error('导入失败')
  }
}

function copyUrl(url) {
  navigator.clipboard.writeText(url).then(() => ElMessage.success('已复制'))
}

async function remove(img) {
  try {
    await ElMessageBox.confirm(`确定删除「${img.original_name}」？`, '提示', { type: 'warning' })
    await deleteImage(img.id)
    ElMessage.success('已删除')
    fetchImages()
  } catch (e) {}
}
</script>

<style scoped>
.img-mgr-toolbar {
  display: flex;
  gap: 8px;
  margin-bottom: 16px;
}
.img-grid {
  display: flex;
  flex-wrap: wrap;
  gap: 12px;
  max-height: 480px;
  overflow-y: auto;
}
.img-card {
  width: 120px;
  border-radius: 8px;
  overflow: hidden;
  background: #f5f5f7;
  position: relative;
}
.img-card img {
  width: 100%;
  height: 90px;
  object-fit: cover;
  display: block;
}
.img-card-actions {
  display: flex;
  gap: 4px;
  padding: 4px;
}
.img-card-name {
  font-size: 11px;
  color: #86868b;
  padding: 0 4px 6px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}
</style>
