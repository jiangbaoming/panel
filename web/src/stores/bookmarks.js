import { defineStore } from 'pinia'
import {
  getGroups, createGroup, updateGroup, deleteGroup, reorderGroups,
  createBookmark, updateBookmark, deleteBookmark, reorderBookmarks,
  getPinned, togglePin, reorderPinned
} from '@/api'

export const useBookmarksStore = defineStore('bookmarks', {
  state: () => ({
    groups: [],
    pinned: [],
    loading: false
  }),
  actions: {
    async fetchAll() {
      this.loading = true
      try {
        const [groups, pinned] = await Promise.all([getGroups(), getPinned()])
        this.groups = groups
        this.pinned = pinned
      } finally {
        this.loading = false
      }
    },
    async addGroup(name, icon) {
      const group = await createGroup(name, icon)
      this.groups.push({ ...group, bookmarks: group.bookmarks || [] })
    },
    async editGroup(id, name, icon) {
      await updateGroup(id, name, icon)
      const g = this.groups.find(g => g.id === id)
      if (g) { g.name = name; g.icon = icon }
    },
    async removeGroup(id) {
      await deleteGroup(id)
      this.groups = this.groups.filter(g => g.id !== id)
    },
    async sortGroups(ids) {
      await reorderGroups(ids)
    },
    async addBookmark(groupId, data) {
      const bm = await createBookmark(groupId, data)
      const g = this.groups.find(g => g.id === groupId)
      if (g) {
        if (!g.bookmarks) g.bookmarks = []
        g.bookmarks.push(bm)
      }
    },
    async editBookmark(groupId, bookmarkId, data) {
      await updateBookmark(groupId, bookmarkId, data)
      const g = this.groups.find(g => g.id === groupId)
      if (g) {
        const b = g.bookmarks.find(b => b.id === bookmarkId)
        if (b) Object.assign(b, data)
      }
    },
    async removeBookmark(groupId, bookmarkId) {
      await deleteBookmark(groupId, bookmarkId)
      const g = this.groups.find(g => g.id === groupId)
      if (g) g.bookmarks = g.bookmarks.filter(b => b.id !== bookmarkId)
    },
    async sortBookmarks(groupId, ids) {
      await reorderBookmarks(groupId, ids)
    },
    // 常驻
    async togglePinned(id, pinned) {
      await togglePin(id, pinned)
      if (!pinned) {
        this.pinned = this.pinned.filter(b => b.id !== id)
      } else {
        await this.fetchAll()
      }
    },
    async sortPinned(ids) {
      await reorderPinned(ids)
    }
  }
})
