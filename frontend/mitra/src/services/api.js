import axios from 'axios'
import router from '../router'
import { useToastStore } from '../stores/toast'

const api = axios.create({
  baseURL: import.meta.env.VITE_API_BASE_URL || '/api/v1/mitra',
  headers: { 'Content-Type': 'application/json' },
})

api.interceptors.request.use((config) => {
  const token = localStorage.getItem('mitra_token')
  if (token) config.headers.Authorization = `Bearer ${token}`
  return config
})

api.interceptors.response.use(
  (res) => res,
  (err) => {
    if (err.response?.status === 401) {
      localStorage.removeItem('mitra_token')
      localStorage.removeItem('mitra_user')
      const toast = useToastStore()
      toast.error('Sesi telah berakhir, silakan login kembali')
      router.push('/login')
    }
    return Promise.reject(err)
  }
)

export default api

// Auth
export const authApi = {
  login: (data) => api.post('/auth/login', data),
  verifyOtp: (data) => api.post('/auth/verify-otp', data),
  resendOtp: (data) => api.post('/auth/resend-otp', data),
  register: (data) => api.post('/auth/register', data),
  profile: () => api.get('/profile'),
  updateProfile: (data) => api.put('/profile', data),
  changePassword: (data) => api.post('/change-password', data),
}

// Outlets (browsing)
export const outletApi = {
  list: (params) => api.get('/outlets', { params }),
  getByID: (id) => api.get(`/outlets/${id}`),
  getPackages: (id) => api.get(`/outlets/${id}/packages`),
}

// Partnership Applications
export const applicationApi = {
  apply: (data) => api.post('/applications', data),
  myList: () => api.get('/applications'),
  getByID: (id) => api.get(`/applications/${id}`),
  cancel: (id) => api.post(`/applications/${id}/cancel`),
}

// Partnership (own)
export const partnershipApi = {
  getMine: () => api.get('/partnership'),
}

// Invoices (own)
export const invoiceApi = {
  list: () => api.get('/invoices'),
  getByID: (id) => api.get(`/invoices/${id}`),
}

// Midtrans
export const midtransApi = {
  getClientKey: () => api.get('/midtrans/client-key'),
}


// Agreements (own)
export const agreementApi = {
  list: () => api.get('/agreements'),
}

// Locations (own)
export const locationApi = {
  list: () => api.get('/locations'),
  getByID: (id) => api.get(`/locations/${id}`),
  create: (data) => api.post('/locations', data),
}
