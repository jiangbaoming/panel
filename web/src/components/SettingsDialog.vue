<template>
  <el-dialog title="页面设置" v-model="visible" width="480px" @closed="fetchSettings">
    <el-form label-width="90px">
      <el-form-item label="页面标题">
        <el-input v-model="local.title" placeholder="个人导航页" />
      </el-form-item>
      <el-form-item label="背景图片">
        <div class="bg-row">
          <el-input v-model="local.bgImage" placeholder="图片 URL">
            <template #append>
              <el-button @click="pickImage('bg')">选择</el-button>
            </template>
          </el-input>
          <el-button v-if="local.bgImage" type="danger" text @click="local.bgImage = ''">清除</el-button>
        </div>
        <img v-if="local.bgImage" :src="local.bgImage" class="bg-preview" />
      </el-form-item>
      <el-form-item label="Favicon">
        <el-input v-model="local.favicon" placeholder="图标 URL">
          <template #append>
            <el-button @click="pickImage('favicon')">选择</el-button>
          </template>
        </el-input>
      </el-form-item>
    </el-form>
    <template #footer>
      <el-button @click="visible = false">取消</el-button>
      <el-button type="primary" :loading="loading" @click="handleSave">保存</el-button>
    </template>
    <!-- 图片选择弹窗 -->
    <ImagePicker v-model:visible="pickerVisible" @select="onImageSelected" />
  </el-dialog>
</template>

<script setup>
import { ref, reactive, watch } from 'vue'
import { ElMessage } from 'element-plus'
import { useAuthStore } from '@/stores/auth'
import { useSettingsStore } from '@/stores/settings'
import ImagePicker from './ImagePicker.vue'

const props = defineProps({ visible: Boolean })
const emit = defineEmits(['update:visible'])

const auth = useAuthStore()
const settings = useSettingsStore()
const loading = ref(false)
const pickerVisible = ref(false)
const pickTarget = ref('bg')

const local = reactive({ title: '', bgImage: '', favicon: '' })

watch(() => props.visible, (v) => {
  if (v) {
    local.title = settings.pageTitle
    local.bgImage = settings.bgImage
    local.favicon = settings.pageFavicon
  }
})

function pickImage(target) {
  pickTarget.value = target
  pickerVisible.value = true
}

function onImageSelected(url) {
  if (pickTarget.value === 'bg') local.bgImage = url
  else local.favicon = url
}

async function fetchSettings() {
  await settings.fetch(auth.userId)
}

async function handleSave() {
  loading.value = true
  try {
    await settings.save(auth.userId, {
      page_title: local.title,
      bg_image: local.bgImage,
      page_favicon: local.favicon
    })
    ElMessage.success('已保存')
    emit('update:visible', false)
  } catch (e) {
    ElMessage.error('保存失败')
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.bg-row {
  display: flex;
  gap: 8px;
  margin-bottom: 8px;
}
.bg-preview {
  max-width: 100%;
  max-height: 120px;
  border-radius: 8px;
  object-fit: cover;
}
</style>
