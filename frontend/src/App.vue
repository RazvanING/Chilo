<template>
  <div id="app" class="min-h-screen">
    <nav class="bg-white shadow-lg">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div class="flex justify-between h-16">
          <div class="flex items-center">
            <router-link to="/" class="text-2xl font-bold text-blue-600">
              Library App
            </router-link>
            <div class="hidden md:flex ml-10 space-x-4">
              <router-link to="/books" class="text-gray-700 hover:text-blue-600 px-3 py-2">
                Browse Books
              </router-link>
              <router-link v-if="isAuthenticated" to="/my-library" class="text-gray-700 hover:text-blue-600 px-3 py-2">
                My Library
              </router-link>
              <router-link v-if="isAuthenticated" to="/favorites" class="text-gray-700 hover:text-blue-600 px-3 py-2">
                Favorites
              </router-link>
              <router-link v-if="isAdmin" to="/admin" class="text-gray-700 hover:text-blue-600 px-3 py-2">
                Admin
              </router-link>
            </div>
          </div>

          <div class="flex items-center space-x-4">
            <template v-if="isAuthenticated">
              <router-link to="/profile" class="text-gray-700 hover:text-blue-600">
                Profile
              </router-link>
              <button @click="handleLogout" class="btn btn-secondary">
                Logout
              </button>
            </template>
            <template v-else>
              <router-link to="/login" class="text-gray-700 hover:text-blue-600">
                Login
              </router-link>
              <router-link to="/register" class="btn btn-primary">
                Sign Up
              </router-link>
            </template>
          </div>
        </div>
      </div>
    </nav>

    <main class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
      <router-view />
    </main>

    <footer class="bg-gray-800 text-white mt-20">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        <p class="text-center">&copy; 2024 Library App. All rights reserved.</p>
      </div>
    </footer>
  </div>
</template>

<script setup>
import { computed } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'

const router = useRouter()
const authStore = useAuthStore()

const isAuthenticated = computed(() => authStore.isAuthenticated)
const isAdmin = computed(() => authStore.isAdmin)

async function handleLogout() {
  await authStore.logout()
  router.push('/login')
}
</script>
