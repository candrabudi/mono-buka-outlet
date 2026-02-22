import { defineStore } from 'pinia'
import { ref, computed } from 'vue'

export const useAuthStore = defineStore('auth', () => {
  const token = ref(localStorage.getItem('mitra_token') || '')
  const user = ref(JSON.parse(localStorage.getItem('mitra_user') || 'null'))

  const isAuthenticated = computed(() => !!token.value)
  const userName = computed(() => user.value?.name || '')
  const userInitial = computed(() => userName.value?.charAt(0)?.toUpperCase() || '?')
  const userEmail = computed(() => user.value?.email || '')

  function setAuth(tokenVal, userVal) {
    token.value = tokenVal
    user.value = userVal
    localStorage.setItem('mitra_token', tokenVal)
    localStorage.setItem('mitra_user', JSON.stringify(userVal))
  }

  function logout() {
    token.value = ''
    user.value = null
    localStorage.removeItem('mitra_token')
    localStorage.removeItem('mitra_user')
  }

  return {
    token, user, isAuthenticated, userName, userInitial, userEmail,
    setAuth, logout,
  }
})
