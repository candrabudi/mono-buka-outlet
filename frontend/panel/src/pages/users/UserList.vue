<template>
  <div class="animate-in">
    <!-- Hero header -->
    <div class="user-hero">
      <div class="user-hero-content">
        <div class="user-hero-text">
          <h1 class="user-hero-title">
            <svg width="28" height="28" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M17 21v-2a4 4 0 0 0-4-4H5a4 4 0 0 0-4 4v2"/>
              <circle cx="9" cy="7" r="4"/>
              <path d="M23 21v-2a4 4 0 0 0-3-3.87"/>
              <path d="M16 3.13a4 4 0 0 1 0 7.75"/>
            </svg>
            User Management
          </h1>
          <p class="user-hero-sub">Kelola semua akun user admin, finance, dan mitra</p>
        </div>
        <button v-if="auth.hasRole('master')" @click="openCreateModal" class="user-hero-btn">
          <svg width="16" height="16" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4"/>
          </svg>
          Tambah User
        </button>
      </div>
      <div class="user-stats-bar">
        <div class="user-stat-item">
          <div class="user-stat-dot dot-all"></div>
          <span class="user-stat-label">Total</span>
          <span class="user-stat-value">{{ users.length }}</span>
        </div>
        <div class="user-stat-item">
          <div class="user-stat-dot dot-active"></div>
          <span class="user-stat-label">Aktif</span>
          <span class="user-stat-value">{{ activeCount }}</span>
        </div>
        <div class="user-stat-item">
          <div class="user-stat-dot dot-admin"></div>
          <span class="user-stat-label">Admin</span>
          <span class="user-stat-value">{{ adminCount }}</span>
        </div>
        <div class="user-stat-item">
          <div class="user-stat-dot dot-mitra"></div>
          <span class="user-stat-label">Mitra</span>
          <span class="user-stat-value">{{ mitraCount }}</span>
        </div>
      </div>
    </div>

    <!-- Filter & Search -->
    <div class="user-filter-strip">
      <div class="user-search-box">
        <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor">
          <circle cx="11" cy="11" r="8" stroke-width="2"/>
          <path d="m21 21-4.35-4.35" stroke-width="2" stroke-linecap="round"/>
        </svg>
        <input v-model="filters.search" @input="applyFilters" type="text" placeholder="Cari nama atau email..." />
      </div>
      <div class="user-filter-pills">
        <select v-model="filters.role" @change="applyFilters" class="user-filter-select">
          <option value="">Semua Role</option>
          <option value="master">Master</option>
          <option value="admin">Admin</option>
          <option value="finance">Finance</option>
          <option value="mitra">Mitra</option>
        </select>
        <select v-model="filters.status" @change="applyFilters" class="user-filter-select">
          <option value="">Semua Status</option>
          <option value="active">Aktif</option>
          <option value="inactive">Nonaktif</option>
        </select>
      </div>
    </div>

    <!-- Loading skeleton -->
    <div v-if="loading" class="user-table-card">
      <div class="user-table-skeleton">
        <div v-for="n in 5" :key="n" class="skeleton-row">
          <div class="skeleton" style="width:40px;height:40px;border-radius:12px"></div>
          <div style="flex:1">
            <div class="skeleton" style="height:16px;width:60%;border-radius:6px;margin-bottom:6px"></div>
            <div class="skeleton" style="height:12px;width:40%;border-radius:6px"></div>
          </div>
          <div class="skeleton" style="width:70px;height:28px;border-radius:8px"></div>
          <div class="skeleton" style="width:60px;height:28px;border-radius:8px"></div>
        </div>
      </div>
    </div>

    <!-- Empty state -->
    <div v-else-if="!filteredUsers.length" class="user-empty">
      <div class="user-empty-icon">
        <svg width="48" height="48" viewBox="0 0 24 24" fill="none" stroke="#94a3b8" stroke-width="1.5">
          <path d="M17 21v-2a4 4 0 0 0-4-4H5a4 4 0 0 0-4 4v2"/>
          <circle cx="9" cy="7" r="4"/>
          <path d="M23 21v-2a4 4 0 0 0-3-3.87"/>
          <path d="M16 3.13a4 4 0 0 1 0 7.75"/>
        </svg>
      </div>
      <h3>{{ filters.search || filters.role || filters.status ? 'Tidak ada hasil' : 'Belum ada user' }}</h3>
      <p>{{ filters.search || filters.role || filters.status ? 'Coba ubah filter pencarian' : 'Tambahkan user pertama Anda' }}</p>
      <button v-if="!filters.search && !filters.role && !filters.status && auth.hasRole('master')" @click="openCreateModal" class="btn btn-primary" style="margin-top:16px">
        <svg width="16" height="16" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4"/>
        </svg>
        Tambah User
      </button>
    </div>

    <!-- User table -->
    <div v-else class="user-table-card">
      <table class="user-table">
        <thead>
          <tr>
            <th>User</th>
            <th>Role</th>
            <th>Status</th>
            <th>Bergabung</th>
            <th v-if="auth.hasRole('master')" style="text-align:center;">Aksi</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="u in filteredUsers" :key="u.id" class="user-row">
            <td>
              <div class="user-cell">
                <div class="user-avatar" :style="avatarStyle(u)">
                  {{ u.name?.charAt(0)?.toUpperCase() }}
                </div>
                <div class="user-info">
                  <div class="user-name">{{ u.name }}</div>
                  <div class="user-email">{{ u.email }}</div>
                </div>
              </div>
            </td>
            <td>
              <span class="role-badge" :class="'role-' + u.role">
                {{ roleLabel(u.role) }}
              </span>
            </td>
            <td>
              <span class="status-badge" :class="u.is_active ? 'status-active' : 'status-inactive'">
                <span class="status-dot"></span>
                {{ u.is_active ? 'Aktif' : 'Nonaktif' }}
              </span>
            </td>
            <td>
              <span class="user-date">{{ formatDate(u.created_at) }}</span>
            </td>
            <td v-if="auth.hasRole('master')">
              <div class="user-actions">
                <button @click="openEditModal(u)" class="action-btn action-edit" title="Edit">
                  <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor">
                    <path d="M11 4H4a2 2 0 0 0-2 2v14a2 2 0 0 0 2 2h14a2 2 0 0 0 2-2v-7" stroke-width="2"/>
                    <path d="M18.5 2.5a2.121 2.121 0 0 1 3 3L12 15l-4 1 1-4 9.5-9.5z" stroke-width="2"/>
                  </svg>
                </button>
                <button @click="toggleActive(u)" class="action-btn" :class="u.is_active ? 'action-on' : 'action-off'" :title="u.is_active ? 'Nonaktifkan' : 'Aktifkan'">
                  <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor">
                    <path d="M18.36 6.64a9 9 0 1 1-12.73 0" stroke-width="2" stroke-linecap="round"/>
                    <line x1="12" y1="2" x2="12" y2="12" stroke-width="2" stroke-linecap="round"/>
                  </svg>
                </button>
                <button @click="confirmDelete(u)" class="action-btn action-delete" title="Hapus">
                  <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor">
                    <polyline points="3 6 5 6 21 6" stroke-width="2"/>
                    <path d="M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6m3 0V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2" stroke-width="2"/>
                  </svg>
                </button>
              </div>
            </td>
          </tr>
        </tbody>
      </table>
    </div>

    <!-- Create / Edit Modal -->
    <div v-if="showModal" class="modal-overlay" @click.self="closeModal">
      <div class="modal-content" style="max-width:520px;">
        <div class="modal-header">
          <h2>{{ editingUser ? 'Edit User' : 'Tambah User Baru' }}</h2>
          <button @click="closeModal" class="modal-close">
            <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <line x1="18" y1="6" x2="6" y2="18"/>
              <line x1="6" y1="6" x2="18" y2="18"/>
            </svg>
          </button>
        </div>
        <form @submit.prevent="editingUser ? updateUser() : createUser()" class="modal-body">
          <!-- Error alert -->
          <div v-if="formError" class="form-alert form-alert-danger">
            <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <circle cx="12" cy="12" r="10"/>
              <line x1="15" y1="9" x2="9" y2="15"/>
              <line x1="9" y1="9" x2="15" y2="15"/>
            </svg>
            {{ formError }}
          </div>

          <div class="form-group">
            <label class="form-label">Nama Lengkap <span class="required">*</span></label>
            <input v-model="form.name" class="form-input" placeholder="Masukkan nama lengkap" required>
          </div>

          <div class="form-group">
            <label class="form-label">Email <span class="required">*</span></label>
            <input v-model="form.email" type="email" class="form-input" placeholder="email@contoh.com" required>
          </div>

          <div v-if="!editingUser" class="form-row">
            <div class="form-group">
              <label class="form-label">Password <span class="required">*</span></label>
              <div class="input-password-wrap">
                <input v-model="form.password" :type="showPassword ? 'text' : 'password'" class="form-input" placeholder="Min. 8 karakter" required minlength="8">
                <button type="button" @click="showPassword = !showPassword" class="password-toggle" tabindex="-1">
                  <svg v-if="!showPassword" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="#94a3b8" stroke-width="2"><path d="M1 12s4-8 11-8 11 8 11 8-4 8-11 8-11-8-11-8z"/><circle cx="12" cy="12" r="3"/></svg>
                  <svg v-else width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="#94a3b8" stroke-width="2"><path d="M17.94 17.94A10.07 10.07 0 0 1 12 20c-7 0-11-8-11-8a18.45 18.45 0 0 1 5.06-5.94M9.9 4.24A9.12 9.12 0 0 1 12 4c7 0 11 8 11 8a18.5 18.5 0 0 1-2.16 3.19m-6.72-1.07a3 3 0 1 1-4.24-4.24"/><line x1="1" y1="1" x2="23" y2="23"/></svg>
                </button>
              </div>
            </div>
            <div class="form-group">
              <label class="form-label">Konfirmasi Password <span class="required">*</span></label>
              <input v-model="form.confirm_password" :type="showPassword ? 'text' : 'password'" class="form-input" placeholder="Ulangi password" required minlength="8">
            </div>
          </div>

          <div class="form-row">
            <div class="form-group">
              <label class="form-label">Role <span class="required">*</span></label>
              <select v-model="form.role" class="form-input" required>
                <option value="" disabled>Pilih role</option>
                <option value="master">Master</option>
                <option value="admin">Admin</option>
                <option value="finance">Finance</option>
                <option value="mitra">Mitra</option>
              </select>
            </div>
            <div class="form-group">
              <label class="form-label">No. HP</label>
              <input v-model="form.phone" type="tel" class="form-input" placeholder="08xxxxxxxxxx">
            </div>
          </div>

          <div class="modal-actions">
            <button type="button" @click="closeModal" class="btn btn-secondary">Batal</button>
            <button type="submit" class="btn btn-primary" :disabled="saving">
              <svg v-if="saving" class="spin" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <path d="M21 12a9 9 0 1 1-6.219-8.56"/>
              </svg>
              {{ saving ? 'Menyimpan...' : (editingUser ? 'Update' : 'Simpan') }}
            </button>
          </div>
        </form>
      </div>
    </div>

    <!-- Delete modal -->
    <div v-if="deleteTarget" class="modal-overlay" @click.self="deleteTarget = null">
      <div class="modal-content" style="max-width:420px;">
        <div class="modal-body" style="padding:2rem;text-align:center;">
          <div class="delete-icon-wrap">
            <svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="#ef4444">
              <polyline points="3 6 5 6 21 6" stroke-width="2"/>
              <path d="M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6m3 0V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2" stroke-width="2"/>
            </svg>
          </div>
          <h3 style="font-size:1.1rem;font-weight:700;margin-bottom:8px;">Hapus User</h3>
          <p style="color:#64748b;font-size:0.875rem;">Yakin ingin menghapus <strong>{{ deleteTarget.name }}</strong>? Aksi ini tidak dapat diurungkan.</p>
        </div>
        <div class="modal-footer" style="justify-content:center;gap:12px;padding-bottom:1.5rem;">
          <button @click="deleteTarget = null" class="btn btn-secondary">Batal</button>
          <button @click="doDelete" class="btn btn-danger" :disabled="deleting">
            {{ deleting ? 'Menghapus...' : 'Ya, Hapus' }}
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted } from 'vue'
import { useAuthStore } from '../../stores/auth'
import { useToastStore } from '../../stores/toast'
import { userApi } from '../../services/api'

const auth = useAuthStore()
const toast = useToastStore()

const users = ref([])
const loading = ref(false)
const showModal = ref(false)
const editingUser = ref(null)
const deleteTarget = ref(null)
const deleting = ref(false)
const saving = ref(false)
const showPassword = ref(false)
const formError = ref('')

const filters = reactive({ search: '', role: '', status: '' })

const defaultForm = { name: '', email: '', password: '', confirm_password: '', role: 'admin', phone: '' }
const form = reactive({ ...defaultForm })

// Computed
const activeCount = computed(() => users.value.filter(u => u.is_active).length)
const adminCount = computed(() => users.value.filter(u => ['master', 'admin'].includes(u.role)).length)
const mitraCount = computed(() => users.value.filter(u => u.role === 'mitra').length)

const filteredUsers = computed(() => {
  let result = [...users.value]
  if (filters.search) {
    const q = filters.search.toLowerCase()
    result = result.filter(u => u.name?.toLowerCase().includes(q) || u.email?.toLowerCase().includes(q))
  }
  if (filters.role) {
    result = result.filter(u => u.role === filters.role)
  }
  if (filters.status === 'active') {
    result = result.filter(u => u.is_active)
  } else if (filters.status === 'inactive') {
    result = result.filter(u => !u.is_active)
  }
  return result
})

onMounted(loadUsers)

async function loadUsers() {
  loading.value = true
  try {
    const { data } = await userApi.list({ page: 1, limit: 100 })
    users.value = data.data || []
  } catch (e) {
    toast.error('Gagal memuat data user')
  } finally {
    loading.value = false
  }
}

function applyFilters() {
  // Client-side filtering via computed
}

// Create
function openCreateModal() {
  editingUser.value = null
  Object.assign(form, { ...defaultForm })
  formError.value = ''
  showPassword.value = false
  showModal.value = true
}

async function createUser() {
  formError.value = ''

  if (form.password !== form.confirm_password) {
    formError.value = 'Konfirmasi password tidak cocok'
    return
  }
  if (form.password.length < 8) {
    formError.value = 'Password minimal 8 karakter'
    return
  }

  saving.value = true
  try {
    await userApi.create({
      name: form.name,
      email: form.email,
      password: form.password,
      confirm_password: form.confirm_password,
      role: form.role,
      phone: form.phone,
    })
    toast.success('User berhasil ditambahkan')
    closeModal()
    await loadUsers()
  } catch (e) {
    const msg = e.response?.data?.error || 'Gagal menambah user'
    formError.value = msg
  } finally {
    saving.value = false
  }
}

// Edit
function openEditModal(u) {
  editingUser.value = u
  Object.assign(form, {
    name: u.name,
    email: u.email,
    password: '',
    confirm_password: '',
    role: u.role,
    phone: u.phone || '',
  })
  formError.value = ''
  showModal.value = true
}

async function updateUser() {
  formError.value = ''
  saving.value = true
  try {
    await userApi.update(editingUser.value.id, {
      name: form.name,
      email: form.email,
      phone: form.phone,
    })
    toast.success('User berhasil diupdate')
    closeModal()
    await loadUsers()
  } catch (e) {
    formError.value = e.response?.data?.error || 'Gagal update user'
  } finally {
    saving.value = false
  }
}

// Toggle Active
async function toggleActive(u) {
  try {
    const { data } = await userApi.toggle(u.id)
    u.is_active = data.data.is_active
    toast.success(data.message || 'Status berhasil diubah')
  } catch {
    toast.error('Gagal mengubah status')
  }
}

// Delete
function confirmDelete(u) { deleteTarget.value = u }
async function doDelete() {
  deleting.value = true
  try {
    await userApi.delete(deleteTarget.value.id)
    toast.success('User berhasil dihapus')
    deleteTarget.value = null
    await loadUsers()
  } catch {
    toast.error('Gagal menghapus user')
  } finally {
    deleting.value = false
  }
}

function closeModal() {
  showModal.value = false
  editingUser.value = null
  formError.value = ''
}

function roleLabel(role) {
  const labels = { master: 'Master', admin: 'Admin', finance: 'Finance', mitra: 'Mitra' }
  return labels[role] || role
}

function formatDate(d) {
  if (!d) return '-'
  return new Date(d).toLocaleDateString('id-ID', { day: 'numeric', month: 'short', year: 'numeric' })
}

const avatarGradients = [
  'linear-gradient(135deg, #667eea, #764ba2)',
  'linear-gradient(135deg, #f093fb, #f5576c)',
  'linear-gradient(135deg, #4facfe, #00f2fe)',
  'linear-gradient(135deg, #43e97b, #38f9d7)',
  'linear-gradient(135deg, #fa709a, #fee140)',
  'linear-gradient(135deg, #a18cd1, #fbc2eb)',
  'linear-gradient(135deg, #fd9644, #fc5c65)',
]
function avatarStyle(u) {
  let h = 0
  for (let i = 0; i < (u.name||'').length; i++) h = ((h << 5) - h) + u.name.charCodeAt(i)
  return { background: avatarGradients[Math.abs(h) % avatarGradients.length] }
}
</script>

<style scoped>
/* ═══ HERO HEADER ═══ */
.user-hero {
  background: linear-gradient(135deg, #0f0c29 0%, #302b63 50%, #24243e 100%);
  border-radius: 16px;
  padding: 32px 36px 24px;
  margin-bottom: 24px;
  box-shadow: 0 4px 24px rgba(15,12,41,0.2);
}
.user-hero-content {
  display: flex; align-items: flex-start; justify-content: space-between; gap: 16px;
  flex-wrap: wrap; margin-bottom: 24px;
}
.user-hero-title {
  font-size: 1.6rem; font-weight: 800; color: #fff; line-height: 1.2;
  letter-spacing: -0.02em; display: flex; align-items: center; gap: 12px;
}
.user-hero-sub { font-size: 0.85rem; color: rgba(255,255,255,0.5); margin-top: 6px; }
.user-hero-btn {
  display: inline-flex; align-items: center; gap: 8px;
  padding: 11px 24px; font-size: 0.85rem; font-weight: 700;
  border-radius: 12px; border: none; cursor: pointer;
  background: linear-gradient(135deg, #f59e0b 0%, #ef4444 100%);
  color: white; box-shadow: 0 4px 20px rgba(245,158,11,0.3);
  white-space: nowrap; font-family: inherit;
}
.user-hero-btn:hover { transform: translateY(-1px); box-shadow: 0 6px 24px rgba(245,158,11,0.4); }

.user-stats-bar {
  display: flex; gap: 28px; flex-wrap: wrap;
  padding-top: 18px;
  border-top: 1px solid rgba(255,255,255,0.08);
}
.user-stat-item { display: flex; align-items: center; gap: 8px; }
.user-stat-dot { width: 8px; height: 8px; border-radius: 50%; }
.dot-all { background: #818cf8; box-shadow: 0 0 8px rgba(129,140,248,0.5); }
.dot-active { background: #4ade80; box-shadow: 0 0 8px rgba(74,222,128,0.5); }
.dot-admin { background: #a78bfa; box-shadow: 0 0 8px rgba(167,139,250,0.5); }
.dot-mitra { background: #fbbf24; box-shadow: 0 0 8px rgba(251,191,36,0.5); }
.user-stat-label { font-size: 0.72rem; color: rgba(255,255,255,0.4); text-transform: uppercase; letter-spacing: 0.05em; }
.user-stat-value { font-size: 0.9rem; font-weight: 800; color: white; }

/* ═══ FILTER STRIP ═══ */
.user-filter-strip {
  display: flex; gap: 12px; align-items: center; flex-wrap: wrap;
  margin-bottom: 24px;
}
.user-search-box {
  display: flex; align-items: center; gap: 10px;
  background: white; border: 1.5px solid #e2e8f0; border-radius: 12px;
  padding: 0 16px; height: 46px; flex: 1; min-width: 200px;
  box-shadow: 0 1px 3px rgba(0,0,0,0.04);
}
.user-search-box:focus-within { border-color: #818cf8; box-shadow: 0 0 0 3px rgba(129,140,248,0.12); }
.user-search-box svg { color: #94a3b8; flex-shrink: 0; }
.user-search-box input {
  flex: 1; border: none; outline: none; font-size: 0.85rem;
  color: #1e293b; background: none; font-family: inherit;
}
.user-search-box input::placeholder { color: #94a3b8; }

.user-filter-pills { display: flex; gap: 8px; flex-wrap: wrap; }
.user-filter-select {
  height: 46px; padding: 0 34px 0 14px; border-radius: 12px;
  border: 1.5px solid #e2e8f0; font-size: 0.8rem; font-weight: 500;
  color: #475569; background: white; cursor: pointer; font-family: inherit;
  appearance: none;
  background-image: url("data:image/svg+xml,%3csvg xmlns='http://www.w3.org/2000/svg' fill='none' viewBox='0 0 20 20'%3e%3cpath stroke='%236b7280' stroke-linecap='round' stroke-linejoin='round' stroke-width='1.5' d='M6 8l4 4 4-4'/%3e%3c/svg%3e");
  background-position: right 10px center; background-repeat: no-repeat; background-size: 1.4em;
  box-shadow: 0 1px 3px rgba(0,0,0,0.04);
}
.user-filter-select:focus { border-color: #818cf8; outline: none; }

/* ═══ TABLE ═══ */
.user-table-card {
  background: #ffffff;
  border-radius: 16px;
  border: 1px solid #e8ecf1;
  overflow: hidden;
  box-shadow: 0 1px 4px rgba(0,0,0,0.04);
}
.user-table {
  width: 100%; border-collapse: collapse;
}
.user-table thead th {
  padding: 14px 20px; font-size: 0.72rem; font-weight: 700;
  color: #94a3b8; text-transform: uppercase; letter-spacing: 0.06em;
  background: #fafbfc; border-bottom: 1px solid #f1f5f9;
  text-align: left; white-space: nowrap;
}
.user-table tbody td {
  padding: 16px 20px; border-bottom: 1px solid #f8fafc;
  vertical-align: middle;
}
.user-row { transition: background 0.15s; }
.user-row:hover { background: #fefefe; }
.user-row:last-child td { border-bottom: none; }

/* User cell */
.user-cell { display: flex; align-items: center; gap: 14px; }
.user-avatar {
  width: 42px; height: 42px; border-radius: 12px;
  display: flex; align-items: center; justify-content: center;
  font-weight: 800; font-size: 1rem; color: white; flex-shrink: 0;
  box-shadow: 0 2px 8px rgba(0,0,0,0.12);
}
.user-info { min-width: 0; }
.user-name { font-size: 0.88rem; font-weight: 700; color: #0f172a; white-space: nowrap; overflow: hidden; text-overflow: ellipsis; }
.user-email { font-size: 0.78rem; color: #94a3b8; margin-top: 2px; }

/* Role badge */
.role-badge {
  display: inline-flex; align-items: center;
  padding: 5px 14px; font-size: 0.7rem; font-weight: 700;
  border-radius: 8px; letter-spacing: 0.04em; text-transform: uppercase;
}
.role-master { background: linear-gradient(135deg, #fef3c7, #fde68a); color: #92400e; }
.role-admin { background: linear-gradient(135deg, #ede9fe, #ddd6fe); color: #5b21b6; }
.role-finance { background: linear-gradient(135deg, #d1fae5, #a7f3d0); color: #065f46; }
.role-mitra { background: linear-gradient(135deg, #dbeafe, #bfdbfe); color: #1e40af; }

/* Status badge */
.status-badge {
  display: inline-flex; align-items: center; gap: 6px;
  padding: 5px 14px; font-size: 0.72rem; font-weight: 600;
  border-radius: 20px;
}
.status-dot { width: 6px; height: 6px; border-radius: 50%; }
.status-active { background: #f0fdf4; color: #16a34a; }
.status-active .status-dot { background: #22c55e; box-shadow: 0 0 6px rgba(34,197,94,0.5); }
.status-inactive { background: #fef2f2; color: #dc2626; }
.status-inactive .status-dot { background: #ef4444; }

.user-date { font-size: 0.8rem; color: #94a3b8; }

/* Actions */
.user-actions { display: flex; gap: 6px; justify-content: center; }
.action-btn {
  width: 34px; height: 34px; border-radius: 10px; border: 1px solid #f1f5f9;
  display: flex; align-items: center; justify-content: center;
  background: #fafbfc; cursor: pointer; transition: all 0.15s;
}
.action-btn:hover { background: #f1f5f9; }
.action-edit { color: #94a3b8; }
.action-edit:hover { color: #6366f1; border-color: #c7d2fe; background: #eef2ff; }
.action-on { color: #22c55e; }
.action-on:hover { background: #f0fdf4; border-color: #bbf7d0; }
.action-off { color: #d1d5db; }
.action-off:hover { color: #f59e0b; background: #fffbeb; border-color: #fde68a; }
.action-delete { color: #cbd5e1; }
.action-delete:hover { color: #ef4444; border-color: #fecaca; background: #fef2f2; }

/* ═══ SKELETON ═══ */
.user-table-skeleton { padding: 20px; }
.skeleton-row {
  display: flex; align-items: center; gap: 16px; padding: 12px 0;
  border-bottom: 1px solid #f8fafc;
}
.skeleton-row:last-child { border-bottom: none; }
.skeleton {
  background: linear-gradient(90deg, #f1f5f9 25%, #e2e8f0 50%, #f1f5f9 75%);
  background-size: 200% 100%;
  animation: shimmer 1.5s infinite;
}
@keyframes shimmer { 0% { background-position: 200% 0; } 100% { background-position: -200% 0; } }

/* ═══ EMPTY STATE ═══ */
.user-empty {
  display: flex; flex-direction: column; align-items: center; justify-content: center;
  text-align: center; padding: 80px 32px;
  background: white; border-radius: 16px; border: 2px dashed #e2e8f0; min-height: 360px;
}
.user-empty-icon {
  display: flex; align-items: center; justify-content: center;
  width: 80px; height: 80px; border-radius: 50%;
  background: linear-gradient(135deg, #f1f5f9, #e2e8f0); margin-bottom: 24px;
}
.user-empty h3 { font-size: 1.15rem; font-weight: 700; color: #0f172a; margin: 0; }
.user-empty p { font-size: 0.85rem; color: #94a3b8; margin-top: 6px; margin-bottom: 0; }

/* ═══ MODAL ═══ */
.modal-close {
  width: 36px; height: 36px; border-radius: 10px; border: none;
  background: rgba(255,255,255,0.06); cursor: pointer; color: #94a3b8;
  display: flex; align-items: center; justify-content: center;
  transition: all 0.15s;
}
.modal-close:hover { background: rgba(255,255,255,0.12); color: #fff; }

.form-row { display: grid; grid-template-columns: 1fr 1fr; gap: 16px; }
.form-group { margin-bottom: 16px; }
.form-label { display: block; font-size: 0.8rem; font-weight: 600; color: #374151; margin-bottom: 6px; }
.required { color: #ef4444; }
.form-input {
  width: 100%; height: 44px; padding: 0 14px; border-radius: 10px;
  border: 1.5px solid #e2e8f0; font-size: 0.85rem; color: #1e293b;
  background: #fff; font-family: inherit; box-sizing: border-box;
  transition: border-color 0.15s, box-shadow 0.15s;
}
.form-input:focus { border-color: #818cf8; outline: none; box-shadow: 0 0 0 3px rgba(129,140,248,0.12); }
.form-input::placeholder { color: #cbd5e1; }
select.form-input { cursor: pointer; appearance: none;
  background-image: url("data:image/svg+xml,%3csvg xmlns='http://www.w3.org/2000/svg' fill='none' viewBox='0 0 20 20'%3e%3cpath stroke='%236b7280' stroke-linecap='round' stroke-linejoin='round' stroke-width='1.5' d='M6 8l4 4 4-4'/%3e%3c/svg%3e");
  background-position: right 10px center; background-repeat: no-repeat; background-size: 1.2em;
}

.input-password-wrap { position: relative; }
.input-password-wrap .form-input { padding-right: 44px; }
.password-toggle {
  position: absolute; right: 8px; top: 50%; transform: translateY(-50%);
  width: 32px; height: 32px; border: none; background: none; cursor: pointer;
  display: flex; align-items: center; justify-content: center; border-radius: 8px;
}
.password-toggle:hover { background: #f1f5f9; }

.form-alert {
  display: flex; align-items: center; gap: 10px;
  padding: 12px 16px; border-radius: 10px; font-size: 0.82rem; font-weight: 500;
  margin-bottom: 16px;
}
.form-alert-danger { background: #fef2f2; color: #dc2626; border: 1px solid #fecaca; }

.modal-actions {
  display: flex; gap: 12px; justify-content: flex-end;
  padding-top: 8px; margin-top: 8px; border-top: 1px solid #f1f5f9;
}

.delete-icon-wrap {
  width: 56px; height: 56px; border-radius: 50%; background: #fee2e2;
  display: flex; align-items: center; justify-content: center;
  margin: 0 auto 16px;
}

.spin { animation: spin 1s linear infinite; }
@keyframes spin { from { transform: rotate(0deg); } to { transform: rotate(360deg); } }

/* ═══ RESPONSIVE ═══ */
@media (max-width: 768px) {
  .user-hero { padding: 24px 20px 18px; margin-bottom: 18px; }
  .user-hero-title { font-size: 1.3rem; }
  .user-hero-title svg { display: none; }
  .user-filter-strip { flex-direction: column; }
  .user-search-box { width: 100%; }
  .user-filter-pills { width: 100%; }
  .user-filter-select { flex: 1; min-width: 0; }
  .form-row { grid-template-columns: 1fr; }

  .user-table-card { overflow-x: auto; }
  .user-table { min-width: 640px; }
}
</style>
