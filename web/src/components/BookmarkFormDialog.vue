<template>
  <el-dialog
    :title="editData ? '编辑书签' : '添加书签'"
    v-model="dialogVisible"
    width="520px"
    :close-on-click-modal="false"
    append-to-body
    @closed="resetForm"
  >
    <el-form ref="formRef" :model="form" :rules="rules" label-width="80px">
      <el-form-item label="名称" prop="name">
        <el-input v-model="form.name" placeholder="书签名称" />
      </el-form-item>
      <el-form-item label="地址" prop="url">
        <el-input v-model="form.url" placeholder="https://" />
      </el-form-item>
      <el-form-item label="图标">
        <IconPicker v-model="form.icon" default-icon="🔗" />
      </el-form-item>
      <el-form-item label="背景色">
        <ColorPicker v-model="form.bg_color" />
      </el-form-item>
      <el-form-item label="图标背景">
        <ColorPicker v-model="form.icon_bg" />
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
import IconPicker from './IconPicker.vue'
import ColorPicker from './ColorPicker.vue'

const props = defineProps({
  visible: Boolean,
  groupId: { type: Number, required: true },
  editData: { type: Object, default: null },
  pinned: { type: Boolean, default: false }
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
  bg_color: 'rgba(255, 255, 255, 0.06)',
  icon_bg: 'rgba(255, 255, 255, 0.06)'
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
      bg_color: props.editData.bg_color || 'rgba(255, 255, 255, 0.06)',
      icon_bg: props.editData.icon_bg || 'rgba(255, 255, 255, 0.06)'
    })
  }
})

function resetForm() {
  form.name = ''; form.url = ''; form.icon = ''
  form.bg_color = 'rgba(255, 255, 255, 0.06)'; form.icon_bg = 'rgba(255, 255, 255, 0.06)'
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
      const bm = await bookmarksStore.addBookmark(props.groupId, { ...form })
      if (props.pinned && bm) {
        await bookmarksStore.togglePinned(bm.id, true)
      }
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