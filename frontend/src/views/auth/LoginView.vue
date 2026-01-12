<template>
  <div class="max-w-md mx-auto mt-10">
    <div class="card">
      <h2 class="text-3xl font-bold mb-6 text-center">Login</h2>

      <div v-if="error" class="bg-red-100 border border-red-400 text-red-700 px-4 py-3 rounded mb-4">
        {{ error }}
      </div>

      <form v-if="!requires2FA" @submit.prevent="handleLogin" class="space-y-4">
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-2">Email</label>
          <input
            v-model="email"
            type="email"
            required
            class="input"
            placeholder="your@email.com"
          />
        </div>

        <div>
          <label class="block text-sm font-medium text-gray-700 mb-2">Password</label>
          <input
            v-model="password"
            type="password"
            required
            class="input"
            placeholder="••••••••"
          />
        </div>

        <button type="submit" :disabled="loading" class="btn btn-primary w-full">
          {{ loading ? 'Logging in...' : 'Login' }}
        </button>
      </form>

      <form v-else @submit.prevent="handle2FAVerification" class="space-y-4">
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-2">
            Enter 2FA Code
          </label>
          <input
            v-model="twoFactorCode"
            type="text"
            required
            maxlength="6"
            class="input"
            placeholder="000000"
          />
        </div>

        <button type="submit" :disabled="loading" class="btn btn-primary w-full">
          {{ loading ? 'Verifying...' : 'Verify' }}
        </button>
      </form>

      <div class="mt-6 text-center">
        <p class="text-gray-600">
          Don't have an account?
          <router-link to="/register" class="text-blue-600 hover:underline">
            Sign up
          </router-link>
        </p>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useAuthStore } from '@/stores/auth'

const router = useRouter()
const route = useRoute()
const authStore = useAuthStore()

const email = ref('')
const password = ref('')
const twoFactorCode = ref('')
const requires2FA = ref(false)
const loading = ref(false)
const error = ref('')

async function handleLogin() {
  loading.value = true
  error.value = ''

  try {
    const result = await authStore.login({
      email: email.value,
      password: password.value
    })

    if (result.requires2FA) {
      requires2FA.value = true
    } else {
      const redirect = route.query.redirect || '/'
      router.push(redirect)
    }
  } catch (err) {
    error.value = err.response?.data?.message || 'Login failed'
  } finally {
    loading.value = false
  }
}

async function handle2FAVerification() {
  loading.value = true
  error.value = ''

  try {
    await authStore.verify2FALogin(twoFactorCode.value)
    const redirect = route.query.redirect || '/'
    router.push(redirect)
  } catch (err) {
    error.value = err.response?.data?.message || '2FA verification failed'
  } finally {
    loading.value = false
  }
}
</script>
