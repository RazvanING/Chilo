import { defineStore } from 'pinia'
import { ref } from 'vue'
import { userBooksAPI, commentsAPI } from '@/services/api'

export const useUserBooksStore = defineStore('userBooks', () => {
  const readingList = ref([])
  const favorites = ref([])
  const comments = ref([])
  const loading = ref(false)
  const error = ref(null)

  async function fetchReadingList(status = '') {
    loading.value = true
    error.value = null
    try {
      const response = await userBooksAPI.getReadingList(status)
      readingList.value = response.data.data || []
    } catch (err) {
      error.value = err.message
      readingList.value = []
    } finally {
      loading.value = false
    }
  }

  async function addToReadingList(bookId, status) {
    try {
      await userBooksAPI.addToReadingList(bookId, status)
      await fetchReadingList()
    } catch (err) {
      error.value = err.message
      throw err
    }
  }

  async function removeFromReadingList(bookId) {
    try {
      await userBooksAPI.removeFromReadingList(bookId)
      readingList.value = readingList.value.filter(item => item.book_id !== bookId)
    } catch (err) {
      error.value = err.message
      throw err
    }
  }

  async function fetchFavorites() {
    loading.value = true
    error.value = null
    try {
      const response = await userBooksAPI.getFavorites()
      favorites.value = response.data.data || []
    } catch (err) {
      error.value = err.message
      favorites.value = []
    } finally {
      loading.value = false
    }
  }

  async function addToFavorites(bookId) {
    try {
      await userBooksAPI.addToFavorites(bookId)
      await fetchFavorites()
    } catch (err) {
      error.value = err.message
      throw err
    }
  }

  async function removeFromFavorites(bookId) {
    try {
      await userBooksAPI.removeFromFavorites(bookId)
      favorites.value = favorites.value.filter(item => item.book_id !== bookId)
    } catch (err) {
      error.value = err.message
      throw err
    }
  }

  async function fetchBookComments(bookId) {
    loading.value = true
    error.value = null
    try {
      const response = await commentsAPI.getBookComments(bookId)
      comments.value = response.data.data || []
    } catch (err) {
      error.value = err.message
      comments.value = []
    } finally {
      loading.value = false
    }
  }

  async function createComment(bookId, content) {
    try {
      const response = await commentsAPI.create(bookId, content)
      comments.value.unshift(response.data)
      return response.data
    } catch (err) {
      error.value = err.message
      throw err
    }
  }

  async function deleteComment(commentId) {
    try {
      await commentsAPI.delete(commentId)
      comments.value = comments.value.filter(c => c.id !== commentId)
    } catch (err) {
      error.value = err.message
      throw err
    }
  }

  return {
    readingList,
    favorites,
    comments,
    loading,
    error,
    fetchReadingList,
    addToReadingList,
    removeFromReadingList,
    fetchFavorites,
    addToFavorites,
    removeFromFavorites,
    fetchBookComments,
    createComment,
    deleteComment
  }
})
