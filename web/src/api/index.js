import axios from 'axios'
import { ElMessage } from 'element-plus'

const api = axios.create({
  baseURL: '/api'
})

// 请求拦截器：自动携带 token
api.interceptors.request.use(config => {
  const token = localStorage.getItem('token')
  if (token) {
    config.headers.Authorization = `Bearer ${token}`
  }
  return config
})

// 响应拦截器：统一解包 + 全局错误提示
api.interceptors.response.use(
  response => {
    const res = response.data
    if (res && typeof res.code === 'number') {
      if (res.code === 0) {
        return res.data
      }
      if (!response.config.silent) {
        ElMessage.error(res.message || '请求失败')
      }
      return Promise.reject(new Error(res.message || '请求失败'))
    }
    return res
  },
  error => {
    if (error.response && error.response.status === 401) {
      localStorage.removeItem('token')
      window.location.href = '/login'
      return Promise.reject(error)
    }
    if (!error.config?.silent) {
      const msg = error.response?.data?.message || '网络错误，请稍后重试'
      ElMessage.error(msg)
    }
    return Promise.reject(error)
  }
)

// 认证
export const login = (username, password) =>
  api.post('/login', { username, password })

export const changePassword = (oldPassword, newPassword) =>
  api.post('/change-password', { oldPassword, newPassword })

// 分组
export const getGroups = () =>
  api.get('/groups')

export const createGroup = (name, icon) =>
  api.post('/groups', { name, icon })

export const updateGroup = (id, name, icon) =>
  api.put(`/groups/${id}`, { name, icon })

export const deleteGroup = id =>
  api.delete(`/groups/${id}`)

export const reorderGroups = ids =>
  api.patch('/groups/reorder', { ids })

// 书签
export const createBookmark = (groupId, data) =>
  api.post(`/groups/${groupId}/bookmarks`, data)

export const updateBookmark = (groupId, bookmarkId, data) =>
  api.put(`/groups/${groupId}/bookmarks/${bookmarkId}`, data)

export const deleteBookmark = (groupId, bookmarkId) =>
  api.delete(`/groups/${groupId}/bookmarks/${bookmarkId}`)

export const reorderBookmarks = (groupId, ids) =>
  api.patch(`/groups/${groupId}/bookmarks/reorder`, { ids })

// 常驻书签
export const getPinned = () =>
  api.get('/pinned')

export const togglePin = (id, pinned) =>
  api.patch(`/pinned/${id}/pin`, { pinned })

export const reorderPinned = ids =>
  api.patch('/pinned/reorder', { ids })

// 图片
export const getImages = (params = {}) =>
  api.get('/images', { params })

export const uploadImage = (formData, onProgress) =>
  api.post('/images/upload', formData, {
    headers: { 'Content-Type': 'multipart/form-data' },
    onUploadProgress: onProgress
  })

export const uploadZip = formData =>
  api.post('/images/upload-zip', formData, {
    headers: { 'Content-Type': 'multipart/form-data' }
  })

export const deleteImage = id =>
  api.delete(`/images/${id}`)

export const clearImages = category =>
  api.delete(`/images/clear/${category}`)

// 用户设置
export const getSettings = userId =>
  api.get(`/settings/${userId}`, { silent: true })

export const updateSettings = (userId, data) =>
  api.put(`/settings/${userId}`, data)

// 搜索
export const searchBookmarks = q =>
  api.get('/search', { params: { q } })

// 用户管理（管理员）
export const getUsers = () =>
  api.get('/users', { silent: true })

export const createUser = (username, password, avatar, role) =>
  api.post('/users', { username, password, avatar, role })

export const deleteUser = id =>
  api.delete(`/users/${id}`)

export const updateMe = username =>
  api.patch('/users/me', { username })

export default api