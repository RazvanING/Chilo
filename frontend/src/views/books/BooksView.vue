<template>
  <div>
    <div class="flex justify-between items-center mb-8">
      <h1 class="text-4xl font-bold">Browse Books</h1>
    </div>

    <div class="mb-6">
      <input
        v-model="searchQuery"
        @input="handleSearch"
        type="text"
        placeholder="Search books by title, author, or description..."
        class="input"
      />
    </div>

    <div v-if="loading" class="text-center py-20">
      <p class="text-xl text-gray-600">Loading books...</p>
    </div>

    <div v-else-if="books.length === 0" class="text-center py-20">
      <p class="text-xl text-gray-600">No books found</p>
    </div>

    <div v-else class="grid md:grid-cols-2 lg:grid-cols-4 gap-6">
      <div
        v-for="book in books"
        :key="book.id"
        class="card hover:shadow-xl transition-shadow cursor-pointer"
        @click="goToBook(book.id)"
      >
        <div class="aspect-w-2 aspect-h-3 mb-4">
          <img
            v-if="book.cover_url"
            :src="book.cover_url"
            :alt="book.title"
            class="w-full h-48 object-cover rounded"
          />
          <div v-else class="w-full h-48 bg-gray-200 rounded flex items-center justify-center">
            <span class="text-gray-400 text-4xl">📚</span>
          </div>
        </div>

        <h3 class="font-semibold text-lg mb-2 line-clamp-2">{{ book.title }}</h3>

        <p v-if="book.authors && book.authors.length" class="text-sm text-gray-600 mb-2">
          by {{ book.authors.map(a => a.name).join(', ') }}
        </p>

        <p class="text-sm text-gray-500 line-clamp-3">{{ book.description }}</p>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useBooksStore } from '@/stores/books'

const router = useRouter()
const booksStore = useBooksStore()

const searchQuery = ref('')
const loading = ref(false)
const books = ref([])
let searchTimeout = null

onMounted(async () => {
  loading.value = true
  await booksStore.fetchBooks()
  books.value = booksStore.books
  loading.value = false
})

function handleSearch() {
  clearTimeout(searchTimeout)
  searchTimeout = setTimeout(async () => {
    if (searchQuery.value.trim()) {
      loading.value = true
      await booksStore.searchBooks(searchQuery.value)
      books.value = booksStore.books
      loading.value = false
    } else {
      loading.value = true
      await booksStore.fetchBooks()
      books.value = booksStore.books
      loading.value = false
    }
  }, 300)
}

function goToBook(id) {
  router.push(`/books/${id}`)
}
</script>
