import axios from 'axios'

// Create axios instances for different services
const authAxios = axios.create({
  baseURL: '/api/auth',
  headers: {
    'Content-Type': 'application/json'
  }
})

const salesAxios = axios.create({
  baseURL: '/api/sales',
  headers: {
    'Content-Type': 'application/json'
  }
})

// Auth API
export const authApi = {
  login: (email: string, password: string) => {
    const credentials = btoa(`${email}:${password}`)
    return authAxios.get('/token/54bb2165-71e1-41a6-af3e-7da4a0e1e2c1', {
      headers: {
        'Authorization': `Basic ${credentials}`
      }
    })
  },

  setAuthToken: (token: string | null) => {
    if (token) {
      salesAxios.defaults.headers.common['Authorization'] = `Bearer ${token}`
      authAxios.defaults.headers.common['Authorization'] = `Bearer ${token}`
    } else {
      delete salesAxios.defaults.headers.common['Authorization']
      delete authAxios.defaults.headers.common['Authorization']
    }
  }
}

// Users API
export const usersApi = {
  getUsers: () => salesAxios.get('/users'),
  
  createUser: (userData: {
    name: string
    email: string
    roles: string[]
    department?: string
    password: string
    passwordConfirm: string
  }) => salesAxios.post('/users', userData),
  
  updateUser: (id: string, userData: any) => salesAxios.put(`/users/${id}`, userData),
  
  deleteUser: (id: string) => salesAxios.delete(`/users/${id}`)
}

// Health check APIs
export const healthApi = {
  checkLiveness: () => salesAxios.get('/liveness'),
  checkReadiness: () => salesAxios.get('/readiness')
}

// Set up response interceptors for error handling
const setupInterceptors = (axiosInstance: typeof axios) => {
  axiosInstance.interceptors.response.use(
    (response) => response,
    (error) => {
      if (error.response?.status === 401) {
        // Handle unauthorized access
        window.location.href = '/login'
      }
      return Promise.reject(error)
    }
  )
}

setupInterceptors(authAxios)
setupInterceptors(salesAxios)