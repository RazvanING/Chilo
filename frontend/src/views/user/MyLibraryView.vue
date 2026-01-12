<template>
  <div>
    <h1 class="text-4xl font-bold mb-8">My Reading List</h1>

    <div class="mb-6 flex space-x-4">
      <button
        v-for="status in statuses"
        :key="status.value"
        @click="filterByStatus(status.value)"
        class="btn"
        :class="currentStatus === status.value ? 'btn-primary' : 'btn-secondary'"
      >
        {{ status.label }}
      </button>
    </div>

    <div v-if="loading" class="text-center py-20">
      <p class="text-xl text-gray-600">Loading...</p>
    </div>

    <div v-else-if="readingList.length === 0" class="text-center py-20">
      <p class="text-xl text-gray-600">No books in this list yet</p>
      <router-link to="/books" class="btn btn-primary mt-4">Browse Books</router-link>
    </div>

    <div v-else class="grid md:grid-cols-2 lg:grid-cols-4 gap-6">
      <div v-for="item in readingList" :key="item.id" class="card">
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
        <p class="text-sm text-gray-600 mb-4">Status: {{ formatStatus(item.status) }}</p>

        <div class="space-y-2">
          <router-link :to="`/books/${item.book_id}`" class="btn btn-primary w-full text-sm">
            View Details
          </router-link>
          <button @click="removeBook(item.book_id)" class="btn btn-danger w-full text-sm">
            Remove
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

const readingList = ref([])
const loading = ref(false)
const currentStatus = ref('')

const statuses = [
  { value: '', label: 'All' },
  { value: 'want_to_read', label: 'Want to Read' },
  { value: 'reading', label: 'Currently Reading' },
  { value: 'read', label: 'Already Read' }
]

onMounted(async () => {
  await fetchBooks()
})

async function fetchBooks(status = '') {
  loading.value = true
  await userBooksStore.fetchReadingList(status)
  readingList.value = userBooksStore.readingList
  loading.value = false
}

function filterByStatus(status) {
  currentStatus.value = status
  fetchBooks(status)
}

async function removeBook(bookId) {
  if (confirm('Remove this book from your reading list?')) {
    try {
      await userBooksStore.removeFromReadingList(bookId)
      readingList.value = readingList.value.filter(item => item.book_id !== bookId)
    } catch (error) {
      alert('Failed to remove book')
    }
  }
}

function formatStatus(status) {
  return status.split('_').map(word => word.charAt(0).toUpperCase() + word.slice(1)).join(' ')
}
</script>
