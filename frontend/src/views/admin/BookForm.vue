<template>
  <div class="max-w-2xl mx-auto">
    <h1 class="text-4xl font-bold mb-8">{{ isEdit ? 'Edit Book' : 'Add New Book' }}</h1>

    <div v-if="error" class="bg-red-100 border border-red-400 text-red-700 px-4 py-3 rounded mb-4">
      {{ error }}
    </div>

    <form @submit.prevent="handleSubmit" class="card space-y-6">
      <div>
        <label class="block text-sm font-medium text-gray-700 mb-2">Title *</label>
        <input v-model="form.title" type="text" required class="input" />
      </div>

      <div>
        <label class="block text-sm font-medium text-gray-700 mb-2">Description *</label>
        <textarea v-model="form.description" rows="5" required class="input"></textarea>
      </div>

      <div>
        <label class="block text-sm font-medium text-gray-700 mb-2">ISBN</label>
        <input v-model="form.isbn" type="text" maxlength="13" class="input" placeholder="1234567890123" />
      </div>

      <div>
        <label class="block text-sm font-medium text-gray-700 mb-2">Published Date *</label>
        <input v-model="form.published_at" type="date" required class="input" />
      </div>

      <div>
        <label class="block text-sm font-medium text-gray-700 mb-2">Authors *</label>
        <select v-model="form.author_ids" multiple class="input" required>
          <option v-for="author in authors" :key="author.id" :value="author.id">
            {{ author.name }}
          </option>
        </select>
        <p class="text-sm text-gray-500 mt-1">Hold Ctrl/Cmd to select multiple authors</p>
      </div>

      <div>
        <label class="block text-sm font-medium text-gray-700 mb-2">Book Cover</label>
        <input type="file" accept="image/*" @change="handleFileChange" class="input" />
        <img v-if="coverPreview" :src="coverPreview" alt="Cover preview" class="mt-2 w-48 rounded" />
      </div>

      <div class="flex space-x-4">
        <button type="submit" :disabled="loading" class="btn btn-primary">
          {{ loading ? 'Saving...' : (isEdit ? 'Update Book' : 'Create Book') }}
        </button>
        <router-link to="/admin" class="btn btn-secondary">
          Cancel
        </router-link>
      </div>
    </form>
  </div>
</template>

<script setup>
import { ref, onMounted, computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useBooksStore } from '@/stores/books'

const route = useRoute()
const router = useRouter()
const booksStore = useBooksStore()

const isEdit = computed(() => !!route.params.id)
const bookId = computed(() => parseInt(route.params.id))

const form = ref({
  title: '',
  description: '',
  isbn: '',
  published_at: '',
  author_ids: []
})

const authors = ref([])
const coverFile = ref(null)
const coverPreview = ref('')
const loading = ref(false)
const error = ref('')

onMounted(async () => {
  await booksStore.fetchAuthors(1, 100)
  authors.value = booksStore.authors

  if (isEdit.value) {
    await booksStore.fetchBook(bookId.value)
    const book = booksStore.currentBook
    if (book) {
      form.value = {
        title: book.title,
        description: book.description,
        isbn: book.isbn || '',
        published_at: book.published_at ? book.published_at.split('T')[0] : '',
        author_ids: book.authors?.map(a => a.id) || []
      }
      if (book.cover_url) {
        coverPreview.value = book.cover_url
      }
    }
  }
})

function handleFileChange(event) {
  const file = event.target.files[0]
  if (file) {
    coverFile.value = file
    const reader = new FileReader()
    reader.onload = (e) => {
      coverPreview.value = e.target.result
    }
    reader.readAsDataURL(file)
  }
}

async function handleSubmit() {
  loading.value = true
  error.value = ''

  try {
    let book
    if (isEdit.value) {
      book = await booksStore.updateBook(bookId.value, form.value)
    } else {
      book = await booksStore.createBook(form.value)
    }

    // Upload cover if selected
    if (coverFile.value && book) {
      await booksStore.uploadBookCover(book.id, coverFile.value)
    }

    router.push('/admin')
  } catch (err) {
    error.value = err.response?.data?.message || 'Failed to save book'
  } finally {
    loading.value = false
  }
}
</script>
