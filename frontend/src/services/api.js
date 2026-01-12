import axios from 'axios'

const api = axios.create({
  baseURL: import.meta.env.VITE_API_URL || '/api',
  headers: {
    'Content-Type': 'application/json'
  }
})

// Request interceptor to add auth token
api.interceptors.request.use(
  (config) => {
    const token = localStorage.getItem('access_token')
    if (token) {
      config.headers.Authorization = `Bearer ${token}`
    }
    return config
  },
  (error) => {
    return Promise.reject(error)
  }
)

// Response interceptor to handle auth errors
api.interceptors.response.use(
  (response) => response,
  (error) => {
    if (error.response?.status === 401) {
      localStorage.removeItem('access_token')
      localStorage.removeItem('user')
      window.location.href = '/login'
    }
    return Promise.reject(error)
  }
)

// Auth API
export const authAPI = {
  register: (data) => api.post('/auth/register', data),
  login: (data) => api.post('/auth/login', data),
  getMe: () => api.get('/auth/me'),
  setup2FA: () => api.post('/auth/2fa/setup'),
  verify2FASetup: (code) => api.post('/auth/2fa/verify', { code }),
  verify2FALogin: (userId, code) => api.post('/auth/2fa/verify-login', { user_id: userId, code }),
  disable2FA: (code) => api.post('/auth/2fa/disable', { code }),
  makeAdmin: (userId) => api.post(`/auth/users/${userId}/make-admin`)
}

// Books API
export const booksAPI = {
  getAll: (params) => api.get('/books', { params }),
  search: (query, params) => api.get('/books/search', { params: { q: query, ...params } }),
  getById: (id) => api.get(`/books/${id}`),
  create: (data) => api.post('/books', data),
  update: (id, data) => api.put(`/books/${id}`, data),
  delete: (id) => api.delete(`/books/${id}`),
  uploadCover: (id, file) => {
    const formData = new FormData()
    formData.append('cover', file)
    return api.post(`/books/${id}/cover`, formData, {
      headers: { 'Content-Type': 'multipart/form-data' }
    })
  }
}

// Authors API
export const authorsAPI = {
  getAll: (params) => api.get('/authors', { params }),
  getById: (id) => api.get(`/authors/${id}`),
  create: (data) => api.post('/authors', data),
  update: (id, data) => api.put(`/authors/${id}`, data),
  delete: (id) => api.delete(`/authors/${id}`)
}

// User Books API
export const userBooksAPI = {
  getReadingList: (status) => api.get('/user/reading-list', { params: { status } }),
  addToReadingList: (bookId, status) => api.post(`/user/books/${bookId}/reading-list`, { status }),
  removeFromReadingList: (bookId) => api.delete(`/user/books/${bookId}/reading-list`),
  getFavorites: () => api.get('/user/favorites'),
  addToFavorites: (bookId) => api.post(`/user/books/${bookId}/favorites`),
  removeFromFavorites: (bookId) => api.delete(`/user/books/${bookId}/favorites`)
}

// Comments API
export const commentsAPI = {
  getBookComments: (bookId) => api.get(`/books/${bookId}/comments`),
  create: (bookId, content) => api.post(`/books/${bookId}/comments`, { content }),
  update: (commentId, content) => api.put(`/comments/${commentId}`, { content }),
  delete: (commentId) => api.delete(`/comments/${commentId}`)
}

export default api
