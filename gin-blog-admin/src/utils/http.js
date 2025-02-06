import axios from 'axios'
import { useAuthStore } from '@/store'

export const request = axios.create(
  {
    baseURL: import.meta.env.VITE_BASE_API,
    timeout: 12000,
  },
)