<template>
  <div class="max-w-2xl mx-auto">
    <h1 class="text-4xl font-bold mb-8">My Profile</h1>

    <div class="card mb-6">
      <h2 class="text-2xl font-semibold mb-4">Account Information</h2>
      <div class="space-y-2">
        <p><span class="font-medium">Email:</span> {{ user?.email }}</p>
        <p><span class="font-medium">Username:</span> {{ user?.username }}</p>
        <p><span class="font-medium">Role:</span> {{ user?.is_admin ? 'Admin' : 'User' }}</p>
      </div>
    </div>

    <div class="card">
      <h2 class="text-2xl font-semibold mb-4">Two-Factor Authentication</h2>

      <div v-if="user?.two_factor_enabled" class="space-y-4">
        <div class="bg-green-100 border border-green-400 text-green-700 px-4 py-3 rounded">
          2FA is enabled for your account
        </div>

        <form @submit.prevent="handleDisable2FA" class="space-y-4">
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-2">
              Enter 2FA Code to Disable
            </label>
            <input
              v-model="disableCode"
              type="text"
              maxlength="6"
              class="input"
              placeholder="000000"
            />
          </div>
          <button type="submit" class="btn btn-danger">
            Disable 2FA
          </button>
        </form>
      </div>

      <div v-else class="space-y-4">
        <p class="text-gray-600">
          Enable two-factor authentication for additional security
        </p>

        <div v-if="!qrCode">
          <button @click="handleSetup2FA" class="btn btn-primary">
            Enable 2FA
          </button>
        </div>

        <div v-else class="space-y-4">
          <div class="bg-blue-100 border border-blue-400 text-blue-700 px-4 py-3 rounded">
            Scan this QR code with your authenticator app
          </div>

          <div class="flex justify-center">
            <img :src="qrCode" alt="QR Code" class="w-64 h-64" />
          </div>

          <form @submit.prevent="handleVerify2FA" class="space-y-4">
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-2">
                Enter the 6-digit code from your authenticator app
              </label>
              <input
                v-model="verifyCode"
                type="text"
                maxlength="6"
                class="input"
                placeholder="000000"
              />
            </div>
            <button type="submit" class="btn btn-primary">
              Verify and Enable
            </button>
          </form>
        </div>
      </div>

      <div v-if="error" class="mt-4 bg-red-100 border border-red-400 text-red-700 px-4 py-3 rounded">
        {{ error }}
      </div>

      <div v-if="success" class="mt-4 bg-green-100 border border-green-400 text-green-700 px-4 py-3 rounded">
        {{ success }}
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue'
import { useAuthStore } from '@/stores/auth'

const authStore = useAuthStore()

const user = computed(() => authStore.user)
const qrCode = ref('')
const verifyCode = ref('')
const disableCode = ref('')
const error = ref('')
const success = ref('')

async function handleSetup2FA() {
  error.value = ''
  try {
    const setup = await authStore.setup2FA()
    qrCode.value = setup.qr_code
  } catch (err) {
    error.value = err.response?.data?.message || 'Failed to setup 2FA'
  }
}

async function handleVerify2FA() {
  error.value = ''
  success.value = ''

  try {
    await authStore.verify2FASetup(verifyCode.value)
    success.value = '2FA enabled successfully!'
    qrCode.value = ''
    verifyCode.value = ''
  } catch (err) {
    error.value = err.response?.data?.message || 'Invalid verification code'
  }
}

async function handleDisable2FA() {
  error.value = ''
  success.value = ''

  try {
    await authStore.disable2FA(disableCode.value)
    success.value = '2FA disabled successfully!'
    disableCode.value = ''
  } catch (err) {
    error.value = err.response?.data?.message || 'Failed to disable 2FA'
  }
}
</script>
