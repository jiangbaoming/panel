import { defineStore } from 'pinia'
import { login as apiLogin, updateMe as apiUpdateMe } from '@/api'
import { getSettings as apiGetSettings } from '@/api'

export const useAuthStore = defineStore('auth', {
  state: () => ({
    token: localStorage.getItem('token') || '',
    user: JSON.parse(localStorage.getItem('user') || 'null'),
    settings: null
  }),
  getters: {
    isLoggedIn: (state) => !!state.token,
    isAdmin: (state) => state.user?.role === 'admin',
    userId: (state) => state.user?.id
  },
  actions: {
    async login(username, password) {
      const res = await apiLogin(username, password)
      this.token = res.token
      this.user = { id: res.id, username: res.username, avatar: res.avatar, role: res.role }
      localStorage.setItem('token', res.token)
      localStorage.setItem('user', JSON.stringify(this.user))
    },
    logout() {
      this.token = ''
      this.user = null
      this.settings = null
      localStorage.removeItem('token')
      localStorage.removeItem('user')
    },
    async updateUsername(username) {
      const res = await apiUpdateMe(username)
      this.user.username = username
      localStorage.setItem('user', JSON.stringify(this.user))
      return res
    }
  }
})
