<template>
  <div>
    <section class="text-center py-20">
      <h1 class="text-5xl font-bold text-gray-900 mb-4">
        Welcome to Library App
      </h1>
      <p class="text-xl text-gray-600 mb-8">
        Discover, track, and share your reading journey
      </p>
      <div class="space-x-4">
        <router-link to="/books" class="btn btn-primary text-lg px-8 py-3">
          Browse Books
        </router-link>
        <router-link v-if="!isAuthenticated" to="/register" class="btn btn-secondary text-lg px-8 py-3">
          Get Started
        </router-link>
      </div>
    </section>

    <section class="grid md:grid-cols-3 gap-8 mt-20">
      <div class="card text-center">
        <div class="text-4xl mb-4">📚</div>
        <h3 class="text-xl font-semibold mb-2">Extensive Library</h3>
        <p class="text-gray-600">Access thousands of books across all genres</p>
      </div>

      <div class="card text-center">
        <div class="text-4xl mb-4">📖</div>
        <h3 class="text-xl font-semibold mb-2">Track Your Reading</h3>
        <p class="text-gray-600">Keep track of what you've read and want to read</p>
      </div>

      <div class="card text-center">
        <div class="text-4xl mb-4">💬</div>
        <h3 class="text-xl font-semibold mb-2">Share & Discuss</h3>
        <p class="text-gray-600">Comment on books and share your thoughts</p>
      </div>
    </section>
  </div>
</template>

<script setup>
import { computed } from 'vue'
import { useAuthStore } from '@/stores/auth'

const authStore = useAuthStore()
const isAuthenticated = computed(() => authStore.isAuthenticated)
</script>
