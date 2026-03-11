import axios from 'axios'
import router from '../router'
import { useToastStore } from '../stores/toast'

const api = axios.create({
  baseURL: import.meta.env.VITE_API_BASE_URL || '/api/v1/admin',
  headers: { 'Content-Type': 'application/json' },
})

api.interceptors.request.use((config) => {
  const token = localStorage.getItem('admin_token')
  if (token) config.headers.Authorization = `Bearer ${token}`
  return config
})

api.interceptors.response.use(
  (res) => res,
  (err) => {
    if (err.response?.status === 401) {
      localStorage.removeItem('admin_token')
      localStorage.removeItem('admin_user')
      const toast = useToastStore()
      toast.error('Sesi telah berakhir, silakan login kembali')
      router.push('/login')
    }
    return Promise.reject(err)
  }
)

export default api

// Auth — OTP flow
export const authApi = {
  login: (data) => api.post('/auth/login', data),           // Step 1: email+pass → OTP sent
  verifyOtp: (data) => api.post('/auth/verify-otp', data),  // Step 2: OTP → JWT
  resendOtp: (data) => api.post('/auth/resend-otp', data),  // Resend OTP
  profile: () => api.get('/profile'),
}




export const outletApi = {
  list: (params) => api.get('/outlets', { params }),
  get: (id) => api.get(`/outlets/${id}`),
  create: (data) => api.post('/outlets', data),
  update: (id, data) => api.put(`/outlets/${id}`, data),
  delete: (id) => api.delete(`/outlets/${id}`),
  toggle: (id) => api.patch(`/outlets/${id}/toggle`),
  toggleFeatured: (id) => api.patch(`/outlets/${id}/featured`),
}

export const outletCategoryApi = {
  list: (params) => api.get('/outlet-categories', { params }),
  get: (id) => api.get(`/outlet-categories/${id}`),
  create: (data) => api.post('/outlet-categories', data),
  update: (id, data) => api.put(`/outlet-categories/${id}`, data),
  delete: (id) => api.delete(`/outlet-categories/${id}`),
  toggle: (id) => api.patch(`/outlet-categories/${id}/toggle`),
}

export const outletPackageApi = {
  listByOutlet: (outletId) => api.get(`/outlet-packages/outlet/${outletId}`),
  get: (id) => api.get(`/outlet-packages/${id}`),
  create: (data) => api.post('/outlet-packages', data),
  update: (id, data) => api.put(`/outlet-packages/${id}`, data),
  delete: (id) => api.delete(`/outlet-packages/${id}`),
}

export const partnershipApi = {
  list: (params) => api.get('/partnerships', { params }),
  get: (id) => api.get(`/partnerships/${id}`),
  create: (data) => api.post('/partnerships', data),
  updateStatus: (id, data) => api.patch(`/partnerships/${id}/status`, data),
}

export const paymentApi = {
  create: (data) => api.post('/payments', data),
  verify: (id, data) => api.patch(`/payments/${id}/verify`, data),
  byPartnership: (pid) => api.get(`/payments/partnership/${pid}`),
}

export const agreementApi = {
  create: (data) => api.post('/agreements', data),
  sign: (id) => api.patch(`/agreements/${id}/sign`),
  byPartnership: (pid) => api.get(`/agreements/partnership/${pid}`),
}

export const revenueApi = {
  create: (data) => api.post('/revenues', data),
  byPartnership: (pid) => api.get(`/revenues/partnership/${pid}`),
}

export const dashboardApi = {
  stats: (params) => api.get('/dashboard', { params }),
}

export const userApi = {
  list: (params) => api.get('/users', { params }),
  get: (id) => api.get(`/users/${id}`),
  create: (data) => api.post('/users', data),
  update: (id, data) => api.put(`/users/${id}`, data),
  delete: (id) => api.delete(`/users/${id}`),
  toggle: (id) => api.patch(`/users/${id}/toggle`),
}

export const uploadApi = {
  upload: (file) => {
    const formData = new FormData()
    formData.append('file', file)
    return api.post('/upload', formData, {
      headers: { 'Content-Type': 'multipart/form-data' },
    })
  },
}

export const meetingApi = {
  list: (params) => api.get('/meetings', { params }),
  get: (id) => api.get(`/meetings/${id}`),
  create: (data) => api.post('/meetings', data),
  update: (id, data) => api.put(`/meetings/${id}`, data),
  delete: (id) => api.delete(`/meetings/${id}`),

  // Participants
  addParticipant: (id, data) => api.post(`/meetings/${id}/participants`, data),
  deleteParticipant: (id, pid) => api.delete(`/meetings/${id}/participants/${pid}`),

  // Notes
  saveNotes: (id, data) => api.post(`/meetings/${id}/notes`, data),

  // Action Plans
  addActionPlan: (id, data) => api.post(`/meetings/${id}/action-plans`, data),
  updateActionPlan: (id, aid, data) => api.put(`/meetings/${id}/action-plans/${aid}`, data),
  deleteActionPlan: (id, aid) => api.delete(`/meetings/${id}/action-plans/${aid}`),

  // Files
  upload: (id, file) => {
    const formData = new FormData()
    formData.append('file', file)
    return api.post(`/meetings/${id}/upload`, formData, {
      headers: { 'Content-Type': 'multipart/form-data' },
    })
  },
  deleteFile: (id, fid) => api.delete(`/meetings/${id}/files/${fid}`),
}

export const settingApi = {
  list: (params) => api.get('/settings', { params }),
  getByKey: (key) => api.get(`/settings/${key}`),
  bulkUpdate: (settings) => api.put('/settings', { settings }),
}

export const invoiceApi = {
  create: (data) => api.post('/invoices', data),
  getByPartnership: (id) => api.get(`/invoices/partnership/${id}`),
  getByID: (id) => api.get(`/invoices/${id}`),
  approve: (id, data) => api.put(`/invoices/${id}/approve`, data),
  checkStatus: (id) => api.get(`/invoices/${id}/check-status`),
  syncPending: () => api.post('/invoices/sync-pending'),
}

export const locationApi = {
  create: (data) => api.post('/location-submissions', data),
  getAll: (params) => api.get('/location-submissions', { params }),
  getByPartnership: (pid) => api.get(`/location-submissions/partnership/${pid}`),
  getByID: (id) => api.get(`/location-submissions/${id}`),
  update: (id, data) => api.put(`/location-submissions/${id}`, data),
  updateStatus: (id, status) => api.patch(`/location-submissions/${id}/status`, { status }),
  remove: (id) => api.delete(`/location-submissions/${id}`),
  recalculate: (id) => api.post(`/location-submissions/${id}/recalculate`),
  createSurvey: (id, data) => api.post(`/location-submissions/${id}/surveys`, data),
  addFile: (id, data) => api.post(`/location-submissions/${id}/files`, data),
  deleteFile: (id, fileId) => api.delete(`/location-submissions/${id}/files/${fileId}`),
  approve: (id, data) => api.post(`/location-submissions/${id}/approve`, data),
}

export const applicationApi = {
  list: (params) => api.get('/partnership-applications', { params }),
  get: (id) => api.get(`/partnership-applications/${id}`),
  review: (id, data) => api.patch(`/partnership-applications/${id}/review`, data),
}

export const ebookApi = {
  list: (params) => api.get('/ebooks', { params }),
  get: (id) => api.get(`/ebooks/${id}`),
  create: (data) => api.post('/ebooks', data),
  update: (id, data) => api.put(`/ebooks/${id}`, data),
  delete: (id) => api.delete(`/ebooks/${id}`),
  toggle: (id) => api.patch(`/ebooks/${id}/toggle`),
}

export const ebookOrderApi = {
  list: (params) => api.get('/ebook-orders', { params }),
  downloadRequests: () => api.get('/ebook-orders/download-requests'),
  approveDownload: (id, data) => api.patch(`/ebook-orders/${id}/approve-download`, data),
  rejectDownload: (id, data) => api.patch(`/ebook-orders/${id}/reject-download`, data),
  approvePayment: (id) => api.patch(`/ebook-orders/${id}/approve-payment`),
  rejectPayment: (id, data) => api.patch(`/ebook-orders/${id}/reject-payment`, data),
}

export const ebookCategoryApi = {
  list: (params) => api.get('/ebook-categories', { params }),
  get: (id) => api.get(`/ebook-categories/${id}`),
  create: (data) => api.post('/ebook-categories', data),
  update: (id, data) => api.put(`/ebook-categories/${id}`, data),
  delete: (id) => api.delete(`/ebook-categories/${id}`),
  toggle: (id) => api.patch(`/ebook-categories/${id}/toggle`),
}

// AI Konsultan Admin
export const aiApi = {
  // Knowledge Base
  listKnowledge: () => api.get('/ai/knowledge'),
  createKnowledge: (data) => api.post('/ai/knowledge', data),
  updateKnowledge: (id, data) => api.put(`/ai/knowledge/${id}`, data),
  deleteKnowledge: (id) => api.delete(`/ai/knowledge/${id}`),

  // Categories
  listCategories: () => api.get('/ai/categories'),
  createCategory: (data) => api.post('/ai/categories', data),

  // System Prompts
  listPrompts: () => api.get('/ai/prompts'),
  createPrompt: (data) => api.post('/ai/prompts', data),
  updatePrompt: (id, data) => api.put(`/ai/prompts/${id}`, data),

  // Config
  getConfig: () => api.get('/ai/config'),
  updateConfig: (data) => api.put('/ai/config', data),

  // Cache
  invalidateCache: () => api.post('/ai/cache/invalidate'),
}

// Affiliator
export const affiliatorApi = {
  dashboard: () => api.get('/affiliator/dashboard'),
  partnerships: (params) => api.get('/affiliator/partnerships', { params }),
  referralCode: () => api.get('/affiliator/referral-code'),
}
