<template>
  <div class="greeting">
    <h1 class="greeting-text">{{ greetingText }}</h1>
    <p class="time-text">{{ currentTime }}</p>
    <p class="date-text">{{ currentDate }}</p>
    <p v-if="settings.welcomeMessage" class="welcome-message">{{ settings.welcomeMessage }}</p>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { useAuthStore } from '@/stores/auth'
import { useSettingsStore } from '@/stores/settings'

const auth = useAuthStore()
const settings = useSettingsStore()

const now = ref(new Date())
let timer = null

const greetingText = computed(() => {
  const hour = now.value.getHours()
  const name = auth.user?.username || ''
  if (hour < 6) return `夜深了，${name}`
  if (hour < 9) return `早上好，${name}`
  if (hour < 12) return `上午好，${name}`
  if (hour < 14) return `中午好，${name}`
  if (hour < 18) return `下午好，${name}`
  if (hour < 22) return `晚上好，${name}`
  return `夜深了，${name}`
})

const currentTime = computed(() => {
  return now.value.toLocaleTimeString('zh-CN', { hour: '2-digit', minute: '2-digit' })
})

const currentDate = computed(() => {
  const weekDays = ['星期日', '星期一', '星期二', '星期三', '星期四', '星期五', '星期六']
  const d = now.value
  return `${d.getFullYear()}年${d.getMonth() + 1}月${d.getDate()}日 ${weekDays[d.getDay()]}`
})

onMounted(() => {
  timer = setInterval(() => { now.value = new Date() }, 1000)
})

onUnmounted(() => {
  clearInterval(timer)
})
</script>

<style scoped>
.greeting {
  text-align: center;
  padding: 40px 0 40px;
  animation: greetingFadeIn 0.6s ease;
}

@keyframes greetingFadeIn {
  from {
    opacity: 0;
    transform: translateY(12px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

.greeting-text {
  font-size: 32px;
  font-weight: 700;
  color: var(--color-text-primary);
  margin: 0 0 12px;
  letter-spacing: 1px;
  text-shadow: 0 1px 4px rgba(0, 0, 0, 0.3);
}

.time-text {
  font-size: 56px;
  font-weight: 300;
  color: var(--color-text-primary);
  margin: 0 0 4px;
  font-variant-numeric: tabular-nums;
  letter-spacing: 2px;
  text-shadow: 0 1px 6px rgba(0, 0, 0, 0.3);
}

.date-text {
  font-size: 15px;
  color: var(--color-text-secondary);
  margin: 0;
  text-shadow: 0 1px 2px rgba(0, 0, 0, 0.25);
}

.welcome-message {
  font-size: 14px;
  color: var(--color-text-secondary);
  margin: 12px 0 0;
  text-shadow: 0 1px 2px rgba(0, 0, 0, 0.25);
  font-style: italic;
  opacity: 0.8;
}

@media (max-width: 768px) {
  .greeting {
    padding: 40px 0 24px;
  }
  .greeting-text {
    font-size: 22px;
  }
  .time-text {
    font-size: 36px;
  }
  .date-text {
    font-size: 13px;
  }
}
</style>