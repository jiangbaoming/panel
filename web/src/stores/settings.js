import { defineStore } from 'pinia'
import { getSettings, updateSettings as apiUpdateSettings } from '@/api'

export const useSettingsStore = defineStore('settings', {
  state: () => ({
    bgImage: '',
    displayMode: 'both',
    pageTitle: '个人导航页',
    pageFavicon: '',
    footer: '',
    welcomeMessage: ''
  }),
  actions: {
    async fetch(userId) {
      try {
        const s = await getSettings(userId)
        this.bgImage = s.bg_image || ''
        this.displayMode = s.display_mode || 'both'
        this.pageTitle = s.page_title || '个人导航页'
        this.pageFavicon = s.page_favicon || ''
        this.footer = s.footer || ''
        this.welcomeMessage = s.welcome_message || ''
        document.title = this.pageTitle
        updateFavicon(this.pageFavicon)
      } catch (e) {
        // 默认值已设置
      }
    },
    async save(userId, data) {
      const s = await apiUpdateSettings(userId, data)
      this.bgImage = s.bg_image || ''
      this.displayMode = s.display_mode || 'both'
      this.pageTitle = s.page_title || '个人导航页'
      this.pageFavicon = s.page_favicon || ''
      this.footer = s.footer || ''
      this.welcomeMessage = s.welcome_message || ''
      document.title = this.pageTitle
      updateFavicon(this.pageFavicon)
    }
  }
})

function updateFavicon(url) {
  const valid = url && (url.startsWith('http') || url.startsWith('/') || url.startsWith('data:'))
  const href = valid ? url : 'data:,'
  let link = document.querySelector('link[rel="icon"]')
  if (!link) {
    link = document.createElement('link')
    link.rel = 'icon'
    document.head.appendChild(link)
  }
  link.href = href
}