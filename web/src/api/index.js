import axios from 'axios'

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

// 响应拦截器：401 自动跳转登录
api.interceptors.response.use(
  response => response,
  error => {
    if (error.response && error.response.status === 401) {
      localStorage.removeItem('token')
      window.location.href = '/login'
    }
    return Promise.reject(error)
  }
)

// 认证
export const login = (username, password) =>
  api.post('/login', { username, password }).then(r => r.data)

export const changePassword = (oldPassword, newPassword) =>
  api.post('/change-password', { oldPassword, newPassword }).then(r => r.data)

// 分组
export const getGroups = () =>
  api.get('/groups').then(r => r.data)

export const createGroup = (name, icon) =>
  api.post('/groups', { name, icon }).then(r => r.data)

export const updateGroup = (id, name, icon) =>
  api.put(`/groups/${id}`, { name, icon }).then(r => r.data)

export const deleteGroup = id =>
  api.delete(`/groups/${id}`).then(r => r.data)

export const reorderGroups = ids =>
  api.patch('/groups/reorder', { ids }).then(r => r.data)

// 书签
export const createBookmark = (groupId, data) =>
  api.post(`/groups/${groupId}/bookmarks`, data).then(r => r.data)

export const updateBookmark = (groupId, bookmarkId, data) =>
  api.put(`/groups/${groupId}/bookmarks/${bookmarkId}`, data).then(r => r.data)

export const deleteBookmark = (groupId, bookmarkId) =>
  api.delete(`/groups/${groupId}/bookmarks/${bookmarkId}`).then(r => r.data)

export const reorderBookmarks = (groupId, ids) =>
  api.patch(`/groups/${groupId}/bookmarks/reorder`, { ids }).then(r => r.data)

// 常驻书签
export const getPinned = () =>
  api.get('/pinned').then(r => r.data)

export const togglePin = (id, pinned) =>
  api.patch(`/pinned/${id}/pin`, { pinned }).then(r => r.data)

export const reorderPinned = ids =>
  api.patch('/pinned/reorder', { ids }).then(r => r.data)

// 图片
export const getImages = (params = {}) =>
  api.get('/images', { params }).then(r => r.data)

export const uploadImage = (formData, onProgress) =>
  api.post('/images/upload', formData, {
    headers: { 'Content-Type': 'multipart/form-data' },
    onUploadProgress: onProgress
  }).then(r => r.data)

export const uploadZip = formData =>
  api.post('/images/upload-zip', formData, {
    headers: { 'Content-Type': 'multipart/form-data' }
  }).then(r => r.data)

export const deleteImage = id =>
  api.delete(`/images/${id}`).then(r => r.data)

export const clearImages = category =>
  api.delete(`/images/clear/${category}`).then(r => r.data)

// 用户设置
export const getSettings = userId =>
  api.get(`/settings/${userId}`).then(r => r.data)

export const updateSettings = (userId, data) =>
  api.put(`/settings/${userId}`, data).then(r => r.data)

// 搜索
export const searchBookmarks = q =>
  api.get('/search', { params: { q } }).then(r => r.data)

// 用户管理（管理员）
export const getUsers = () =>
  api.get('/users').then(r => r.data)

export const createUser = (username, password, avatar, role) =>
  api.post('/users', { username, password, avatar, role }).then(r => r.data)

export const deleteUser = id =>
  api.delete(`/users/${id}`).then(r => r.data)

export const updateMe = username =>
  api.patch('/users/me', { username }).then(r => r.data)

export default api
