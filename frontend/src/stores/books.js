import { defineStore } from 'pinia'
import { ref } from 'vue'
import { booksAPI, authorsAPI } from '@/services/api'

export const useBooksStore = defineStore('books', () => {
  const books = ref([])
  const currentBook = ref(null)
  const authors = ref([])
  const loading = ref(false)
  const error = ref(null)

  async function fetchBooks(page = 1, pageSize = 20) {
    loading.value = true
    error.value = null
    try {
      const response = await booksAPI.getAll({ page, page_size: pageSize })
      books.value = response.data.data || []
    } catch (err) {
      error.value = err.message
      books.value = []
    } finally {
      loading.value = false
    }
  }

  async function searchBooks(query, page = 1, pageSize = 20) {
    loading.value = true
    error.value = null
    try {
      const response = await booksAPI.search(query, { page, page_size: pageSize })
      books.value = response.data.data || []
    } catch (err) {
      error.value = err.message
      books.value = []
    } finally {
      loading.value = false
    }
  }

  async function fetchBook(id) {
    loading.value = true
    error.value = null
    try {
      const response = await booksAPI.getById(id)
      currentBook.value = response.data.data
    } catch (err) {
      error.value = err.message
      currentBook.value = null
    } finally {
      loading.value = false
    }
  }

  async function createBook(bookData) {
    loading.value = true
    error.value = null
    try {
      const response = await booksAPI.create(bookData)
      books.value.unshift(response.data)
      return response.data
    } catch (err) {
      error.value = err.message
      throw err
    } finally {
      loading.value = false
    }
  }

  async function updateBook(id, bookData) {
    loading.value = true
    error.value = null
    try {
      const response = await booksAPI.update(id, bookData)
      const index = books.value.findIndex(b => b.id === id)
      if (index !== -1) {
        books.value[index] = response.data.data
      }
      return response.data.data
    } catch (err) {
      error.value = err.message
      throw err
    } finally {
      loading.value = false
    }
  }

  async function deleteBook(id) {
    loading.value = true
    error.value = null
    try {
      await booksAPI.delete(id)
      books.value = books.value.filter(b => b.id !== id)
    } catch (err) {
      error.value = err.message
      throw err
    } finally {
      loading.value = false
    }
  }

  async function uploadBookCover(bookId, file) {
    try {
      const response = await booksAPI.uploadCover(bookId, file)
      return response.data.data.cover_url
    } catch (err) {
      error.value = err.message
      throw err
    }
  }

  async function fetchAuthors(page = 1, pageSize = 50) {
    try {
      const response = await authorsAPI.getAll({ page, page_size: pageSize })
      authors.value = response.data.data || []
    } catch (err) {
      error.value = err.message
      authors.value = []
    }
  }

  async function createAuthor(authorData) {
    try {
      const response = await authorsAPI.create(authorData)
      authors.value.push(response.data)
      return response.data
    } catch (err) {
      error.value = err.message
      throw err
    }
  }

  return {
    books,
    currentBook,
    authors,
    loading,
    error,
    fetchBooks,
    searchBooks,
    fetchBook,
    createBook,
    updateBook,
    deleteBook,
    uploadBookCover,
    fetchAuthors,
    createAuthor
  }
})
