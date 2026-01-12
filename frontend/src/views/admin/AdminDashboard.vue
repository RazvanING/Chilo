<template>
  <div>
    <h1 class="text-4xl font-bold mb-8">Admin Dashboard</h1>

    <div class="grid md:grid-cols-3 gap-6 mb-8">
      <div class="card text-center">
        <div class="text-4xl mb-2">📚</div>
        <h3 class="text-3xl font-bold">{{ books.length }}</h3>
        <p class="text-gray-600">Total Books</p>
      </div>

      <div class="card text-center">
        <div class="text-4xl mb-2">✍️</div>
        <h3 class="text-3xl font-bold">{{ authors.length }}</h3>
        <p class="text-gray-600">Total Authors</p>
      </div>

      <div class="card text-center">
        <div class="text-4xl mb-2">⚙️</div>
        <h3 class="text-3xl font-bold">Admin</h3>
        <p class="text-gray-600">Management</p>
      </div>
    </div>

    <div class="card mb-6">
      <div class="flex justify-between items-center mb-6">
        <h2 class="text-2xl font-semibold">Manage Books</h2>
        <router-link to="/admin/books/new" class="btn btn-primary">
          + Add New Book
        </router-link>
      </div>

      <div v-if="loading" class="text-center py-10">
        <p class="text-gray-600">Loading...</p>
      </div>

      <div v-else-if="books.length === 0" class="text-center py-10">
        <p class="text-gray-600">No books yet</p>
      </div>

      <div v-else class="overflow-x-auto">
        <table class="min-w-full divide-y divide-gray-200">
          <thead class="bg-gray-50">
            <tr>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Title</th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Authors</th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">ISBN</th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Actions</th>
            </tr>
          </thead>
          <tbody class="bg-white divide-y divide-gray-200">
            <tr v-for="book in books" :key="book.id">
              <td class="px-6 py-4 whitespace-nowrap">{{ book.title }}</td>
              <td class="px-6 py-4 whitespace-nowrap">
                {{ book.authors?.map(a => a.name).join(', ') }}
              </td>
              <td class="px-6 py-4 whitespace-nowrap">{{ book.isbn || '-' }}</td>
              <td class="px-6 py-4 whitespace-nowrap space-x-2">
                <router-link :to="`/admin/books/${book.id}/edit`" class="text-blue-600 hover:underline">
                  Edit
                </router-link>
                <button @click="deleteBook(book.id)" class="text-red-600 hover:underline">
                  Delete
                </button>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>

    <div class="card">
      <div class="flex justify-between items-center mb-6">
        <h2 class="text-2xl font-semibold">Manage Authors</h2>
        <router-link to="/admin/authors" class="btn btn-primary">
          Manage Authors
        </router-link>
      </div>

      <div v-if="authors.length === 0" class="text-center py-10">
        <p class="text-gray-600">No authors yet</p>
      </div>

      <div v-else class="grid md:grid-cols-2 lg:grid-cols-3 gap-4">
        <div v-for="author in authors.slice(0, 6)" :key="author.id" class="p-4 bg-gray-50 rounded">
          <h4 class="font-semibold">{{ author.name }}</h4>
          <p class="text-sm text-gray-600 line-clamp-2">{{ author.bio || 'No bio available' }}</p>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useBooksStore } from '@/stores/books'

const booksStore = useBooksStore()

const books = ref([])
const authors = ref([])
const loading = ref(false)

onMounted(async () => {
  loading.value = true
  await booksStore.fetchBooks(1, 100)
  books.value = booksStore.books
  await booksStore.fetchAuthors(1, 100)
  authors.value = booksStore.authors
  loading.value = false
})

async function deleteBook(bookId) {
  if (confirm('Are you sure you want to delete this book?')) {
    try {
      await booksStore.deleteBook(bookId)
      books.value = books.value.filter(b => b.id !== bookId)
    } catch (error) {
      alert('Failed to delete book')
    }
  }
}
</script>
