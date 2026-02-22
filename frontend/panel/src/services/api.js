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

export const leadApi = {
  list: (params) => api.get('/leads', { params }),
  get: (id) => api.get(`/leads/${id}`),
  kanban: (params) => api.get('/leads/kanban', { params }),
  create: (data) => api.post('/leads', data),
  update: (id, data) => api.put(`/leads/${id}`, data),
  updateStatus: (id, data) => api.patch(`/leads/${id}/status`, data),
  delete: (id) => api.delete(`/leads/${id}`),
}

export const partnershipApi = {
  list: (params) => api.get('/partnerships', { params }),
  get: (id) => api.get(`/partnerships/${id}`),
  create: (data) => api.post('/partnerships', data),
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
