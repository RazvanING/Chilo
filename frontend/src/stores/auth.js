import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { authAPI } from '@/services/api'

export const useAuthStore = defineStore('auth', () => {
  const user = ref(null)
  const token = ref(localStorage.getItem('access_token') || null)
  const requires2FA = ref(false)
  const tempUserId = ref(null)

  const isAuthenticated = computed(() => !!token.value)
  const isAdmin = computed(() => user.value?.is_admin || false)

  function setAuth(authData) {
    token.value = authData.tokens.access_token
    user.value = authData.user
    localStorage.setItem('access_token', authData.tokens.access_token)
    localStorage.setItem('user', JSON.stringify(authData.user))
  }

  async function register(credentials) {
    const response = await authAPI.register(credentials)
    return response.data
  }

  async function login(credentials) {
    const response = await authAPI.login(credentials)

    if (response.data.data.requires_2fa) {
      requires2FA.value = true
      tempUserId.value = response.data.data.user_id
      return { requires2FA: true }
    }

    setAuth(response.data.data)
    return { requires2FA: false }
  }

  async function verify2FALogin(code) {
    const response = await authAPI.verify2FALogin(tempUserId.value, code)
    setAuth({ tokens: response.data.data.tokens, user: user.value })
    requires2FA.value = false
    tempUserId.value = null
  }

  async function logout() {
    token.value = null
    user.value = null
    requires2FA.value = false
    tempUserId.value = null
    localStorage.removeItem('access_token')
    localStorage.removeItem('user')
  }

  async function fetchUser() {
    try {
      const response = await authAPI.getMe()
      user.value = response.data.data
      localStorage.setItem('user', JSON.stringify(user.value))
    } catch (error) {
      await logout()
    }
  }

  async function setup2FA() {
    const response = await authAPI.setup2FA()
    return response.data.data
  }

  async function verify2FASetup(code) {
    await authAPI.verify2FASetup(code)
    if (user.value) {
      user.value.two_factor_enabled = true
    }
  }

  async function disable2FA(code) {
    await authAPI.disable2FA(code)
    if (user.value) {
      user.value.two_factor_enabled = false
    }
  }

  // Initialize from localStorage
  if (token.value) {
    const savedUser = localStorage.getItem('user')
    if (savedUser) {
      user.value = JSON.parse(savedUser)
    } else {
      fetchUser()
    }
  }

  return {
    user,
    token,
    requires2FA,
    isAuthenticated,
    isAdmin,
    register,
    login,
    verify2FALogin,
    logout,
    fetchUser,
    setup2FA,
    verify2FASetup,
    disable2FA
  }
})
