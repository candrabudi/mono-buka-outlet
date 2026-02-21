import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { authApi } from '../services/api'

export const useAuthStore = defineStore('auth', () => {
  const token = ref(localStorage.getItem('admin_token') || '')
  const user = ref(JSON.parse(localStorage.getItem('admin_user') || 'null'))

  // OTP state
  const otpEmail = ref('')
  const otpSessionId = ref('')
  const otpStep = ref('login') // 'login' | 'otp'

  const isAuthenticated = computed(() => !!token.value)
  const userRole = computed(() => user.value?.role || '')
  const userName = computed(() => user.value?.name || '')
  const userInitial = computed(() => userName.value?.charAt(0)?.toUpperCase() || '?')

  const ADMIN_ROLES = ['master', 'admin', 'finance']
  const isAdmin = computed(() => ADMIN_ROLES.includes(userRole.value))

  function hasRole(...roles) {
    return roles.includes(userRole.value)
  }

  function roleLabel(role) {
    const labels = {
      master: 'Master',
      admin: 'Admin',
      finance: 'Finance',
      mitra: 'Mitra',
    }
    return labels[role] || role
  }

  // Step 1: login → sends OTP to email
  async function login(credentials) {
    const { data } = await authApi.login(credentials)
    const d = data.data
    otpEmail.value = d.email
    otpSessionId.value = d.session_id
    otpStep.value = 'otp'
    return data
  }

  // Step 2: verify OTP → get JWT
  async function verifyOtp(code) {
    const { data } = await authApi.verifyOtp({ email: otpEmail.value, code })
    const d = data.data
    token.value = d.token
    user.value = d.user
    localStorage.setItem('admin_token', d.token)
    localStorage.setItem('admin_user', JSON.stringify(d.user))
    otpStep.value = 'login'
    otpEmail.value = ''
    otpSessionId.value = ''
    return data
  }

  // Resend OTP
  async function resendOtp() {
    const { data } = await authApi.resendOtp({ email: otpEmail.value })
    otpSessionId.value = data.data.session_id
    return data
  }

  // Cancel OTP step → go back
  function cancelOtp() {
    otpStep.value = 'login'
    otpEmail.value = ''
    otpSessionId.value = ''
  }

  async function fetchProfile() {
    try {
      const { data } = await authApi.profile()
      user.value = data.data
      localStorage.setItem('admin_user', JSON.stringify(data.data))
    } catch {
      logout()
    }
  }

  function logout() {
    token.value = ''
    user.value = null
    localStorage.removeItem('admin_token')
    localStorage.removeItem('admin_user')
  }

  return {
    token, user, isAuthenticated, userRole, userName, userInitial, isAdmin,
    otpEmail, otpSessionId, otpStep,
    hasRole, roleLabel, login, verifyOtp, resendOtp, cancelOtp, fetchProfile, logout,
  }
})
