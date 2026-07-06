<template>
  <el-dialog
    :title="editData ? '编辑书签' : '添加书签'"
    v-model="dialogVisible"
    width="480px"
    :close-on-click-modal="false"
    @closed="resetForm"
  >
    <el-form ref="formRef" :model="form" :rules="rules" label-width="70px">
      <el-form-item label="名称" prop="name">
        <el-input v-model="form.name" placeholder="书签名称" />
      </el-form-item>
      <el-form-item label="地址" prop="url">
        <el-input v-model="form.url" placeholder="https://" />
      </el-form-item>
      <el-form-item label="图标">
        <el-input v-model="form.icon" placeholder="Emoji 或图片 URL">
          <template #prepend><span style="font-size:18px">{{ form.icon || '🔗' }}</span></template>
        </el-input>
      </el-form-item>
      <el-form-item label="背景色">
        <el-color-picker v-model="form.bg_color" show-alpha />
        <span class="form-tip">卡片背景色</span>
      </el-form-item>
      <el-form-item label="图标背景">
        <el-color-picker v-model="form.icon_bg" show-alpha />
        <span class="form-tip">图标区域背景色</span>
      </el-form-item>
    </el-form>
    <template #footer>
      <el-button @click="dialogVisible = false">取消</el-button>
      <el-button type="primary" :loading="loading" @click="handleSave">保存</el-button>
    </template>
  </el-dialog>
</template>

<script setup>
import { ref, reactive, watch, computed } from 'vue'
import { ElMessage } from 'element-plus'
import { useBookmarksStore } from '@/stores/bookmarks'

const props = defineProps({
  visible: Boolean,
  groupId: { type: Number, required: true },
  editData: { type: Object, default: null }
})
const emit = defineEmits(['update:visible', 'saved'])

const dialogVisible = computed({
  get: () => props.visible,
  set: (v) => emit('update:visible', v)
})

const bookmarksStore = useBookmarksStore()
const formRef = ref(null)
const loading = ref(false)

const form = reactive({
  name: '',
  url: '',
  icon: '',
  bg_color: '',
  icon_bg: ''
})

const rules = {
  name: [{ required: true, message: '请输入名称', trigger: 'blur' }],
  url: [{ required: true, message: '请输入地址', trigger: 'blur' }]
}

watch(() => props.visible, (v) => {
  if (v && props.editData) {
    Object.assign(form, {
      name: props.editData.name || '',
      url: props.editData.url || '',
      icon: props.editData.icon || '',
      bg_color: props.editData.bg_color || '',
      icon_bg: props.editData.icon_bg || ''
    })
  }
})

function resetForm() {
  form.name = ''; form.url = ''; form.icon = ''
  form.bg_color = ''; form.icon_bg = ''
  formRef.value?.resetFields()
}

async function handleSave() {
  const valid = await formRef.value.validate().catch(() => false)
  if (!valid) return
  loading.value = true
  try {
    if (props.editData) {
      await bookmarksStore.editBookmark(props.groupId, props.editData.id, { ...form })
    } else {
      await bookmarksStore.addBookmark(props.groupId, { ...form })
    }
    ElMessage.success(props.editData ? '已更新' : '已添加')
    dialogVisible.value = false
    emit('saved')
  } catch (e) {
    ElMessage.error(e.response?.data?.error || '保存失败')
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.form-tip {
  font-size: 12px;
  color: #86868b;
  margin-left: 8px;
}
</style>
