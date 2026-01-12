<template>
  <div v-if="loading" class="text-center py-20">
    <p class="text-xl text-gray-600">Loading...</p>
  </div>

  <div v-else-if="book" class="max-w-4xl mx-auto">
    <div class="grid md:grid-cols-3 gap-8 mb-8">
      <div>
        <img
          v-if="book.cover_url"
          :src="book.cover_url"
          :alt="book.title"
          class="w-full rounded-lg shadow-lg"
        />
        <div v-else class="w-full aspect-w-2 aspect-h-3 bg-gray-200 rounded-lg flex items-center justify-center">
          <span class="text-gray-400 text-6xl">📚</span>
        </div>

        <div v-if="isAuthenticated" class="mt-4 space-y-2">
          <button @click="toggleFavorite" class="btn w-full" :class="isFavorited ? 'btn-danger' : 'btn-primary'">
            {{ isFavorited ? '💔 Remove from Favorites' : '❤️ Add to Favorites' }}
          </button>

          <select v-model="readingStatus" @change="updateReadingStatus" class="input">
            <option value="">Add to Reading List</option>
            <option value="want_to_read">Want to Read</option>
            <option value="reading">Currently Reading</option>
            <option value="read">Already Read</option>
          </select>
        </div>
      </div>

      <div class="md:col-span-2">
        <h1 class="text-4xl font-bold mb-4">{{ book.title }}</h1>

        <div v-if="book.authors && book.authors.length" class="mb-4">
          <p class="text-xl text-gray-600">
            by {{ book.authors.map(a => a.name).join(', ') }}
          </p>
        </div>

        <div v-if="book.isbn" class="mb-4">
          <p class="text-sm text-gray-500">ISBN: {{ book.isbn }}</p>
        </div>

        <div class="mb-6">
          <h2 class="text-2xl font-semibold mb-2">Description</h2>
          <p class="text-gray-700 leading-relaxed">{{ book.description }}</p>
        </div>

        <div v-if="book.published_at" class="mb-6">
          <p class="text-sm text-gray-500">
            Published: {{ new Date(book.published_at).toLocaleDateString() }}
          </p>
        </div>
      </div>
    </div>

    <div class="card mt-8">
      <h2 class="text-2xl font-semibold mb-6">Comments</h2>

      <form v-if="isAuthenticated" @submit.prevent="submitComment" class="mb-8">
        <textarea
          v-model="newComment"
          rows="3"
          placeholder="Share your thoughts about this book..."
          class="input"
        ></textarea>
        <button type="submit" :disabled="!newComment.trim()" class="btn btn-primary mt-2">
          Post Comment
        </button>
      </form>

      <div v-if="comments.length === 0" class="text-gray-500 text-center py-8">
        No comments yet. Be the first to share your thoughts!
      </div>

      <div v-else class="space-y-4">
        <div v-for="comment in comments" :key="comment.id" class="border-b pb-4">
          <div class="flex justify-between items-start mb-2">
            <span class="font-semibold">{{ comment.username }}</span>
            <span class="text-sm text-gray-500">
              {{ new Date(comment.created_at).toLocaleDateString() }}
            </span>
          </div>
          <p class="text-gray-700">{{ comment.content }}</p>
        </div>
      </div>
    </div>
  </div>

  <div v-else class="text-center py-20">
    <p class="text-xl text-gray-600">Book not found</p>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { useBooksStore } from '@/stores/books'
import { useUserBooksStore } from '@/stores/userBooks'
import { useAuthStore } from '@/stores/auth'

const route = useRoute()
const booksStore = useBooksStore()
const userBooksStore = useUserBooksStore()
const authStore = useAuthStore()

const book = ref(null)
const loading = ref(true)
const newComment = ref('')
const readingStatus = ref('')
const isFavorited = ref(false)

const isAuthenticated = computed(() => authStore.isAuthenticated)
const comments = computed(() => userBooksStore.comments)

onMounted(async () => {
  const bookId = parseInt(route.params.id)
  await booksStore.fetchBook(bookId)
  book.value = booksStore.currentBook
  await userBooksStore.fetchBookComments(bookId)
  loading.value = false
})

async function submitComment() {
  if (!newComment.value.trim()) return

  try {
    await userBooksStore.createComment(book.value.id, newComment.value)
    newComment.value = ''
  } catch (error) {
    alert('Failed to post comment')
  }
}

async function toggleFavorite() {
  try {
    if (isFavorited.value) {
      await userBooksStore.removeFromFavorites(book.value.id)
      isFavorited.value = false
    } else {
      await userBooksStore.addToFavorites(book.value.id)
      isFavorited.value = true
    }
  } catch (error) {
    alert('Failed to update favorites')
  }
}

async function updateReadingStatus() {
  if (!readingStatus.value) return

  try {
    await userBooksStore.addToReadingList(book.value.id, readingStatus.value)
    alert('Added to reading list!')
  } catch (error) {
    alert('Failed to update reading list')
  }
}
</script>
