import { defineStore } from 'pinia'
import { getSettings, updateSettings as apiUpdateSettings } from '@/api'

export const useSettingsStore = defineStore('settings', {
  state: () => ({
    bgImage: '',
    displayMode: 'both',
    pageTitle: '个人导航页',
    pageFavicon: '',
    footer: ''
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
        document.title = this.pageTitle
        if (this.pageFavicon) {
          const link = document.querySelector('link[rel="icon"]')
          if (link) link.href = this.pageFavicon
        }
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
      document.title = this.pageTitle
      if (this.pageFavicon) {
        const link = document.querySelector('link[rel="icon"]')
        if (link) link.href = this.pageFavicon
      }
    }
  }
})
