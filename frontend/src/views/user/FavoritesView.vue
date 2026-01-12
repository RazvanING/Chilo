<template>
  <div>
    <h1 class="text-4xl font-bold mb-8">My Favorites</h1>

    <div v-if="loading" class="text-center py-20">
      <p class="text-xl text-gray-600">Loading...</p>
    </div>

    <div v-else-if="favorites.length === 0" class="text-center py-20">
      <p class="text-xl text-gray-600">No favorite books yet</p>
      <router-link to="/books" class="btn btn-primary mt-4">Browse Books</router-link>
    </div>

    <div v-else class="grid md:grid-cols-2 lg:grid-cols-4 gap-6">
      <div v-for="item in favorites" :key="item.id" class="card">
        <div class="aspect-w-2 aspect-h-3 mb-4">
          <img
            v-if="item.book?.cover_url"
            :src="item.book.cover_url"
            :alt="item.book.title"
            class="w-full h-48 object-cover rounded"
          />
          <div v-else class="w-full h-48 bg-gray-200 rounded flex items-center justify-center">
            <span class="text-gray-400 text-4xl">📚</span>
          </div>
        </div>

        <h3 class="font-semibold text-lg mb-2 line-clamp-2">{{ item.book?.title }}</h3>

        <div class="space-y-2">
          <router-link :to="`/books/${item.book_id}`" class="btn btn-primary w-full text-sm">
            View Details
          </router-link>
          <button @click="removeBook(item.book_id)" class="btn btn-danger w-full text-sm">
            Remove from Favorites
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useUserBooksStore } from '@/stores/userBooks'

const userBooksStore = useUserBooksStore()

const favorites = ref([])
const loading = ref(false)

onMounted(async () => {
  loading.value = true
  await userBooksStore.fetchFavorites()
  favorites.value = userBooksStore.favorites
  loading.value = false
})

async function removeBook(bookId) {
  if (confirm('Remove this book from favorites?')) {
    try {
      await userBooksStore.removeFromFavorites(bookId)
      favorites.value = favorites.value.filter(item => item.book_id !== bookId)
    } catch (error) {
      alert('Failed to remove from favorites')
    }
  }
}
</script>
