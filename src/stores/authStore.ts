import { create } from 'zustand'
import { persist } from 'zustand/middleware'
import { authApi } from '../services/api'

interface User {
  id: string
  name: string
  email: string
  roles: string[]
}

interface AuthState {
  user: User | null
  token: string | null
  isAuthenticated: boolean
  login: (email: string, password: string) => Promise<void>
  logout: () => void
  setUser: (user: User) => void
}

export const useAuthStore = create<AuthState>()(
  persist(
    (set, get) => ({
      user: null,
      token: null,
      isAuthenticated: false,

      login: async (email: string, password: string) => {
        try {
          const response = await authApi.login(email, password)
          const { token } = response.data
          
          // Set token in API headers
          authApi.setAuthToken(token)
          
          // Get user info (you might need to implement this endpoint)
          // For now, we'll extract basic info from the token or use mock data
          const user: User = {
            id: '1',
            name: email.split('@')[0],
            email,
            roles: ['ADMIN'] // This should come from your auth response
          }

          set({
            user,
            token,
            isAuthenticated: true
          })
        } catch (error) {
          console.error('Login failed:', error)
          throw error
        }
      },

      logout: () => {
        authApi.setAuthToken(null)
        set({
          user: null,
          token: null,
          isAuthenticated: false
        })
      },

      setUser: (user: User) => {
        set({ user })
      }
    }),
    {
      name: 'auth-storage',
      partialize: (state) => ({
        user: state.user,
        token: state.token,
        isAuthenticated: state.isAuthenticated
      })
    }
  )
)