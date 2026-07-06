<template>
  <el-dialog
    :title="editData ? '编辑分组' : '新建分组'"
    v-model="dialogVisible"
    width="480px"
    :close-on-click-modal="false"
    append-to-body
    @closed="resetForm"
  >
    <el-form ref="formRef" :model="form" :rules="rules" label-width="60px">
      <el-form-item label="名称" prop="name">
        <el-input v-model="form.name" placeholder="分组名称" />
      </el-form-item>
      <el-form-item label="图标">
        <IconPicker v-model="form.icon" default-icon="📁" />
      </el-form-item>
    </el-form>
    <template #footer>
      <el-button @click="dialogVisible = false">取消</el-button>
      <el-button type="primary" :loading="loading" @click="handleSave">保存</el-button>
    </template>
  </el-dialog>
</template>

<script setup>
import { ref, reactive, computed } from 'vue'
import { ElMessage } from 'element-plus'
import { useBookmarksStore } from '@/stores/bookmarks'
import IconPicker from './IconPicker.vue'

const props = defineProps({
  visible: Boolean,
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
const form = reactive({ name: '', icon: '' })
const rules = {
  name: [{ required: true, message: '请输入分组名称', trigger: 'blur' }]
}

async function handleSave() {
  const valid = await formRef.value.validate().catch(() => false)
  if (!valid) return
  loading.value = true
  try {
    if (props.editData) {
      await bookmarksStore.editGroup(props.editData.id, form.name, form.icon)
    } else {
      await bookmarksStore.addGroup(form.name, form.icon)
    }
    ElMessage.success(props.editData ? '已更新' : '已创建')
    dialogVisible.value = false
    emit('saved')
  } catch (e) {
    ElMessage.error(e.response?.data?.error || '保存失败')
  } finally {
    loading.value = false
  }
}

function resetForm() {
  form.name = ''; form.icon = ''
  formRef.value?.resetFields()
}
</script>