<template>
  <div>
    <div class="flex justify-between items-center mb-8">
      <h1 class="text-4xl font-bold">Manage Authors</h1>
      <button @click="showAddForm = true" class="btn btn-primary">
        + Add New Author
      </button>
    </div>

    <div v-if="showAddForm" class="card mb-6">
      <h2 class="text-2xl font-semibold mb-4">{{ editingAuthor ? 'Edit Author' : 'Add New Author' }}</h2>

      <form @submit.prevent="handleSubmitAuthor" class="space-y-4">
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-2">Name *</label>
          <input v-model="authorForm.name" type="text" required class="input" />
        </div>

        <div>
          <label class="block text-sm font-medium text-gray-700 mb-2">Bio</label>
          <textarea v-model="authorForm.bio" rows="3" class="input"></textarea>
        </div>

        <div class="flex space-x-4">
          <button type="submit" class="btn btn-primary">
            {{ editingAuthor ? 'Update' : 'Create' }}
          </button>
          <button type="button" @click="cancelAuthorForm" class="btn btn-secondary">
            Cancel
          </button>
        </div>
      </form>
    </div>

    <div v-if="loading" class="text-center py-20">
      <p class="text-xl text-gray-600">Loading...</p>
    </div>

    <div v-else class="grid md:grid-cols-2 lg:grid-cols-3 gap-6">
      <div v-for="author in authors" :key="author.id" class="card">
        <h3 class="text-xl font-semibold mb-2">{{ author.name }}</h3>
        <p class="text-gray-600 mb-4 line-clamp-3">{{ author.bio || 'No bio available' }}</p>
        <div class="flex space-x-2">
          <button @click="editAuthor(author)" class="btn btn-secondary text-sm">
            Edit
          </button>
          <button @click="deleteAuthor(author.id)" class="btn btn-danger text-sm">
            Delete
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useBooksStore } from '@/stores/books'

const booksStore = useBooksStore()

const authors = ref([])
const loading = ref(false)
const showAddForm = ref(false)
const editingAuthor = ref(null)

const authorForm = ref({
  name: '',
  bio: ''
})

onMounted(async () => {
  await fetchAuthors()
})

async function fetchAuthors() {
  loading.value = true
  await booksStore.fetchAuthors(1, 100)
  authors.value = booksStore.authors
  loading.value = false
}

async function handleSubmitAuthor() {
  try {
    if (editingAuthor.value) {
      await booksStore.updateAuthor(editingAuthor.value.id, authorForm.value)
    } else {
      await booksStore.createAuthor(authorForm.value)
    }
    await fetchAuthors()
    cancelAuthorForm()
  } catch (error) {
    alert('Failed to save author')
  }
}

function editAuthor(author) {
  editingAuthor.value = author
  authorForm.value = {
    name: author.name,
    bio: author.bio || ''
  }
  showAddForm.value = true
}

async function deleteAuthor(authorId) {
  if (confirm('Are you sure you want to delete this author?')) {
    try {
      await booksStore.deleteAuthor(authorId)
      authors.value = authors.value.filter(a => a.id !== authorId)
    } catch (error) {
      alert('Failed to delete author')
    }
  }
}

function cancelAuthorForm() {
  showAddForm.value = false
  editingAuthor.value = null
  authorForm.value = { name: '', bio: '' }
}
</script>
