<template>
  <div class="mitra-page">
    <!-- Hero Header -->
    <div class="mitra-hero">
      <div class="mitra-hero-content">
        <div>
          <h1 class="mitra-hero-title">Management Mitra</h1>
          <p class="mitra-hero-sub">Kelola data mitra dan akses portal mitra</p>
        </div>
        <button @click="openCreate" class="mitra-add-btn">
          <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor"><line x1="12" y1="5" x2="12" y2="19" stroke-width="2" stroke-linecap="round"/><line x1="5" y1="12" x2="19" y2="12" stroke-width="2" stroke-linecap="round"/></svg>
          Tambah Mitra
        </button>
      </div>
      <div class="mitra-stats-bar">
        <div class="mitra-stat">
          <span class="mitra-stat-dot dot-total"></span>
          <span class="mitra-stat-label">Total Mitra</span>
          <span class="mitra-stat-value">{{ total }}</span>
        </div>
        <div class="mitra-stat">
          <span class="mitra-stat-dot dot-active"></span>
          <span class="mitra-stat-label">Aktif</span>
          <span class="mitra-stat-value">{{ mitras.filter(m => m.is_active).length }}</span>
        </div>
        <div class="mitra-stat">
          <span class="mitra-stat-dot dot-inactive"></span>
          <span class="mitra-stat-label">Nonaktif</span>
          <span class="mitra-stat-value">{{ mitras.filter(m => !m.is_active).length }}</span>
        </div>
      </div>
    </div>

    <!-- Search & Filter -->
    <div class="mitra-toolbar">
      <div class="mitra-search-wrap">
        <svg class="mitra-search-icon" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="#94a3b8" stroke-width="2"><circle cx="11" cy="11" r="8"/><line x1="21" y1="21" x2="16.65" y2="16.65"/></svg>
        <input v-model="searchQuery" type="text" class="mitra-search" placeholder="Cari nama, email, atau telepon..." />
      </div>
    </div>

    <!-- Loading -->
    <div v-if="loading" class="mitra-grid">
      <div v-for="n in 6" :key="n" class="mitra-card mitra-card-skel">
        <div class="skel-avatar shimmer"></div>
        <div class="skel-name shimmer"></div>
        <div class="skel-info shimmer"></div>
        <div class="skel-info-sm shimmer"></div>
      </div>
    </div>

    <!-- Mitra Grid -->
    <div v-else-if="filteredMitras.length" class="mitra-grid">
      <div v-for="m in filteredMitras" :key="m.id" class="mitra-card" :class="{ 'mitra-card-inactive': !m.is_active }">
        <div class="mitra-card-top">
          <div class="mitra-card-avatar" :style="avatarBg(m.name)">{{ getInitial(m.name) }}</div>
          <button @click="toggleActive(m)" class="mitra-toggle-btn" :class="m.is_active ? 'toggle-on' : 'toggle-off'" :title="m.is_active ? 'Nonaktifkan' : 'Aktifkan'">
            <div class="toggle-track"><div class="toggle-thumb"></div></div>
          </button>
        </div>

        <div class="mitra-card-body">
          <h3 class="mitra-card-name">{{ m.name }}</h3>
          <div class="mitra-card-info">
            <div class="mitra-info-row">
              <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="#94a3b8" stroke-width="2"><path d="M4 4h16c1.1 0 2 .9 2 2v12c0 1.1-.9 2-2 2H4c-1.1 0-2-.9-2-2V6c0-1.1.9-2 2-2z"/><polyline points="22,6 12,13 2,6"/></svg>
              <span>{{ m.email }}</span>
            </div>
            <div class="mitra-info-row">
              <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="#94a3b8" stroke-width="2"><path d="M22 16.92v3a2 2 0 0 1-2.18 2 19.79 19.79 0 0 1-8.63-3.07 19.5 19.5 0 0 1-6-6 19.79 19.79 0 0 1-3.07-8.67A2 2 0 0 1 4.11 2h3a2 2 0 0 1 2 1.72c.12.96.36 1.9.7 2.81a2 2 0 0 1-.45 2.11L8.09 9.91a16 16 0 0 0 6 6l1.27-1.27a2 2 0 0 1 2.11-.45c.91.36 1.85.64 2.81.7A2 2 0 0 1 22 16.92z"/></svg>
              <span>{{ m.phone || '-' }}</span>
            </div>
          </div>
        </div>

        <div class="mitra-card-footer">
          <span class="mitra-card-status" :class="m.is_active ? 'st-active' : 'st-inactive'">
            <span class="st-dot"></span>
            {{ m.is_active ? 'Aktif' : 'Nonaktif' }}
          </span>
          <div class="mitra-card-btns">
            <button @click="openEdit(m)" class="mitra-icon-btn edit-btn" title="Edit">
              <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor"><path d="M11 4H4a2 2 0 0 0-2 2v14a2 2 0 0 0 2 2h14a2 2 0 0 0 2-2v-7" stroke-width="2"/><path d="M18.5 2.5a2.121 2.121 0 0 1 3 3L12 15l-4 1 1-4 9.5-9.5z" stroke-width="2"/></svg>
            </button>
            <button @click="confirmDelete(m)" class="mitra-icon-btn del-btn" title="Hapus">
              <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor"><polyline points="3 6 5 6 21 6" stroke-width="2"/><path d="M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6m3 0V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2" stroke-width="2"/></svg>
            </button>
          </div>
        </div>

        <div class="mitra-card-date">
          Bergabung: {{ formatDate(m.created_at) }}
        </div>
      </div>
    </div>

    <!-- Empty state -->
    <div v-else class="mitra-empty">
      <div class="mitra-empty-circle">
        <svg width="36" height="36" viewBox="0 0 24 24" fill="none" stroke="#94a3b8" stroke-width="1.5"><path d="M17 21v-2a4 4 0 0 0-4-4H5a4 4 0 0 0-4 4v2"/><circle cx="9" cy="7" r="4"/><path d="M23 21v-2a4 4 0 0 0-3-3.87"/><path d="M16 3.13a4 4 0 0 1 0 7.75"/></svg>
      </div>
      <h3>Belum ada mitra</h3>
      <p>Tambahkan mitra pertama untuk memulai.</p>
      <button @click="openCreate" class="mitra-add-btn" style="margin-top:20px;font-size:0.82rem;">+ Tambah Mitra</button>
    </div>

    <!-- Modal Create/Edit -->
    <Teleport to="body">
      <div v-if="showModal" class="mitra-modal-overlay" @click.self="closeModal">
        <div class="mitra-modal" @click.stop>
          <div class="mitra-modal-header">
            <div class="mitra-modal-title-group">
              <div class="mitra-modal-icon">
                <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M20 21v-2a4 4 0 0 0-4-4H8a4 4 0 0 0-4 4v2"/><circle cx="12" cy="7" r="4"/></svg>
              </div>
              <h3>{{ editItem ? 'Edit Mitra' : 'Tambah Mitra Baru' }}</h3>
            </div>
            <button @click="closeModal" class="mitra-modal-close">&times;</button>
          </div>
          <form @submit.prevent="handleSave" class="mitra-modal-body">
            <div class="mitra-form-group">
              <label>Nama Lengkap <span class="req">*</span></label>
              <input v-model="form.name" type="text" class="mitra-input" :class="{ 'is-error': formErrors.name }" placeholder="Nama lengkap mitra" />
              <div v-if="formErrors.name" class="mitra-field-error">{{ formErrors.name }}</div>
            </div>
            <div class="mitra-form-row">
              <div class="mitra-form-group">
                <label>Email <span class="req">*</span></label>
                <input v-model="form.email" type="email" class="mitra-input" :class="{ 'is-error': formErrors.email }" placeholder="email@example.com" />
                <div v-if="formErrors.email" class="mitra-field-error">{{ formErrors.email }}</div>
              </div>
              <div class="mitra-form-group">
                <label>No. Handphone <span class="req">*</span></label>
                <input v-model="form.phone" type="text" class="mitra-input" :class="{ 'is-error': formErrors.phone }" placeholder="08xxxxxxxxxx" />
                <div v-if="formErrors.phone" class="mitra-field-error">{{ formErrors.phone }}</div>
              </div>
            </div>
            <div class="mitra-form-row">
              <div class="mitra-form-group">
                <label>Password {{ editItem ? '(kosongkan jika tidak ingin ubah)' : '' }} <span v-if="!editItem" class="req">*</span></label>
                <div class="mitra-pw-wrap">
                  <input v-model="form.password" :type="showPassword ? 'text' : 'password'" class="mitra-input mitra-pw-input" :class="{ 'is-error': formErrors.password }" placeholder="Minimal 8 karakter" />
                  <button type="button" class="mitra-pw-toggle" @click="showPassword = !showPassword" tabindex="-1">
                    <svg v-if="!showPassword" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="#94a3b8" stroke-width="2"><path d="M1 12s4-8 11-8 11 8 11 8-4 8-11 8-11-8-11-8z"/><circle cx="12" cy="12" r="3"/></svg>
                    <svg v-else width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="#94a3b8" stroke-width="2"><path d="M17.94 17.94A10.07 10.07 0 0 1 12 20c-7 0-11-8-11-8a18.45 18.45 0 0 1 5.06-5.94M9.9 4.24A9.12 9.12 0 0 1 12 4c7 0 11 8 11 8a18.5 18.5 0 0 1-2.16 3.19m-6.72-1.07a3 3 0 1 1-4.24-4.24"/><line x1="1" y1="1" x2="23" y2="23"/></svg>
                  </button>
                </div>
                <div v-if="formErrors.password" class="mitra-field-error">{{ formErrors.password }}</div>
              </div>
              <div class="mitra-form-group">
                <label>Konfirmasi Password <span v-if="!editItem" class="req">*</span></label>
                <div class="mitra-pw-wrap">
                  <input v-model="form.confirm_password" :type="showConfirmPassword ? 'text' : 'password'" class="mitra-input mitra-pw-input" :class="{ 'is-error': formErrors.confirm_password }" placeholder="Ulangi password" />
                  <button type="button" class="mitra-pw-toggle" @click="showConfirmPassword = !showConfirmPassword" tabindex="-1">
                    <svg v-if="!showConfirmPassword" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="#94a3b8" stroke-width="2"><path d="M1 12s4-8 11-8 11 8 11 8-4 8-11 8-11-8-11-8z"/><circle cx="12" cy="12" r="3"/></svg>
                    <svg v-else width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="#94a3b8" stroke-width="2"><path d="M17.94 17.94A10.07 10.07 0 0 1 12 20c-7 0-11-8-11-8a18.45 18.45 0 0 1 5.06-5.94M9.9 4.24A9.12 9.12 0 0 1 12 4c7 0 11 8 11 8a18.5 18.5 0 0 1-2.16 3.19m-6.72-1.07a3 3 0 1 1-4.24-4.24"/><line x1="1" y1="1" x2="23" y2="23"/></svg>
                  </button>
                </div>
                <div v-if="formErrors.confirm_password" class="mitra-field-error">{{ formErrors.confirm_password }}</div>
              </div>
            </div>
            <div v-if="form.password" class="mitra-pw-hints">
              <div class="pw-hint" :class="{ 'pw-ok': form.password.length >= 8 }">Minimal 8 karakter</div>
              <div class="pw-hint" :class="{ 'pw-ok': /[A-Z]/.test(form.password) }">Huruf besar</div>
              <div class="pw-hint" :class="{ 'pw-ok': /[a-z]/.test(form.password) }">Huruf kecil</div>
              <div class="pw-hint" :class="{ 'pw-ok': /[0-9]/.test(form.password) }">Angka</div>
            </div>
            <div class="mitra-modal-footer">
              <button type="button" @click="closeModal" class="btn btn-secondary">Batal</button>
              <button type="submit" class="mitra-save-btn" :disabled="saving">
                <svg v-if="saving" class="spin-icon" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor"><path d="M21 12a9 9 0 1 1-6.219-8.56" stroke-width="2.5" stroke-linecap="round"/></svg>
                {{ saving ? 'Menyimpan...' : (editItem ? 'Update' : 'Simpan') }}
              </button>
            </div>
          </form>
        </div>
      </div>
    </Teleport>

    <!-- Delete Confirm Modal -->
    <Teleport to="body">
      <div v-if="deleteTarget" class="mitra-modal-overlay" @click.self="deleteTarget = null">
        <div class="mitra-modal mitra-modal-sm" @click.stop>
          <div class="mitra-delete-body">
            <div class="mitra-delete-icon-wrap">
              <svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="#ef4444" stroke-width="2"><polyline points="3 6 5 6 21 6"/><path d="M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6m3 0V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2"/></svg>
            </div>
            <h3>Hapus Mitra</h3>
            <p>Yakin ingin menghapus <strong>{{ deleteTarget.name }}</strong>?<br/>Data mitra akan dihapus secara permanen.</p>
          </div>
          <div class="mitra-delete-footer">
            <button @click="deleteTarget = null" class="btn btn-secondary">Batal</button>
            <button @click="doDelete" class="btn btn-danger" :disabled="deleting">
              {{ deleting ? 'Menghapus...' : 'Ya, Hapus' }}
            </button>
          </div>
        </div>
      </div>
    </Teleport>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted } from 'vue'
import { userApi } from '../../services/api'
import { useToastStore } from '../../stores/toast'

const toast = useToastStore()
const mitras = ref([])
const total = ref(0)
const loading = ref(false)
const showModal = ref(false)
const editItem = ref(null)
const saving = ref(false)
const deleteTarget = ref(null)
const deleting = ref(false)
const searchQuery = ref('')
const showPassword = ref(false)
const showConfirmPassword = ref(false)

const form = reactive({ name: '', email: '', phone: '', password: '', confirm_password: '' })
const formErrors = reactive({})

const avatarGradients = [
  'linear-gradient(135deg, #667eea, #764ba2)',
  'linear-gradient(135deg, #f093fb, #f5576c)',
  'linear-gradient(135deg, #4facfe, #00f2fe)',
  'linear-gradient(135deg, #43e97b, #38f9d7)',
  'linear-gradient(135deg, #fa709a, #fee140)',
  'linear-gradient(135deg, #a18cd1, #fbc2eb)',
  'linear-gradient(135deg, #fd9644, #fc5c65)',
  'linear-gradient(135deg, #0abde3, #341f97)',
]

function avatarBg(name) {
  const idx = name.charCodeAt(0) % avatarGradients.length
  return { background: avatarGradients[idx] }
}

function getInitial(name) {
  return name ? name.split(' ').map(n => n[0]).join('').substring(0, 2).toUpperCase() : '?'
}

function formatDate(d) {
  if (!d) return '-'
  return new Date(d).toLocaleDateString('id-ID', { day: 'numeric', month: 'short', year: 'numeric' })
}

const filteredMitras = computed(() => {
  const q = searchQuery.value.toLowerCase().trim()
  if (!q) return mitras.value
  return mitras.value.filter(m =>
    m.name.toLowerCase().includes(q) ||
    m.email.toLowerCase().includes(q) ||
    (m.phone || '').includes(q)
  )
})

onMounted(() => loadMitras())

async function loadMitras() {
  loading.value = true
  try {
    const { data } = await userApi.list({ role: 'mitra', limit: 100 })
    mitras.value = data.data || []
    total.value = data.total || 0
  } catch {
    toast.error('Gagal memuat data mitra')
  } finally {
    loading.value = false
  }
}

function openCreate() {
  editItem.value = null
  Object.assign(form, { name: '', email: '', phone: '', password: '', confirm_password: '' })
  Object.keys(formErrors).forEach(k => delete formErrors[k])
  showPassword.value = false
  showConfirmPassword.value = false
  showModal.value = true
}

function openEdit(m) {
  editItem.value = m
  Object.assign(form, { name: m.name, email: m.email, phone: m.phone, password: '', confirm_password: '' })
  Object.keys(formErrors).forEach(k => delete formErrors[k])
  showPassword.value = false
  showConfirmPassword.value = false
  showModal.value = true
}

function closeModal() {
  showModal.value = false
  editItem.value = null
}

function validate() {
  Object.keys(formErrors).forEach(k => delete formErrors[k])
  if (!form.name?.trim()) formErrors.name = 'Nama wajib diisi'
  if (!form.email?.trim()) formErrors.email = 'Email wajib diisi'
  if (!form.phone?.trim()) formErrors.phone = 'No. handphone wajib diisi'

  const needPassword = !editItem.value
  const hasPassword = form.password && form.password.length > 0

  if (needPassword && !hasPassword) {
    formErrors.password = 'Password wajib diisi'
  }
  if (hasPassword) {
    const pwErrors = []
    if (form.password.length < 8) pwErrors.push('minimal 8 karakter')
    if (!/[A-Z]/.test(form.password)) pwErrors.push('huruf besar')
    if (!/[a-z]/.test(form.password)) pwErrors.push('huruf kecil')
    if (!/[0-9]/.test(form.password)) pwErrors.push('angka')
    if (pwErrors.length) {
      formErrors.password = 'Password harus mengandung: ' + pwErrors.join(', ')
    }
  }
  if (needPassword || hasPassword) {
    if (!form.confirm_password) {
      formErrors.confirm_password = 'Konfirmasi password wajib diisi'
    } else if (form.password !== form.confirm_password) {
      formErrors.confirm_password = 'Konfirmasi password tidak cocok'
    }
  }

  return Object.keys(formErrors).length === 0
}

async function handleSave() {
  if (!validate()) return
  saving.value = true
  try {
    if (editItem.value) {
      const payload = { name: form.name, email: form.email, phone: form.phone }
      if (form.password) payload.password = form.password
      await userApi.update(editItem.value.id, payload)
      toast.success('Mitra berhasil diupdate')
    } else {
      await userApi.create({
        name: form.name, email: form.email, phone: form.phone,
        password: form.password, confirm_password: form.confirm_password, role: 'mitra',
      })
      toast.success('Mitra berhasil ditambahkan')
    }
    closeModal()
    loadMitras()
  } catch (e) {
    toast.error(e.response?.data?.error || 'Gagal menyimpan mitra')
  } finally {
    saving.value = false
  }
}

async function toggleActive(m) {
  try {
    const { data } = await userApi.toggle(m.id)
    toast.success(data.message)
    m.is_active = data.data.is_active
  } catch {
    toast.error('Gagal mengubah status')
  }
}

function confirmDelete(m) { deleteTarget.value = m }

async function doDelete() {
  deleting.value = true
  try {
    const { data } = await userApi.delete(deleteTarget.value.id)
    toast.success(data.message || 'Mitra berhasil dihapus')
    deleteTarget.value = null
    loadMitras()
  } catch (e) {
    toast.error(e.response?.data?.error || 'Gagal menghapus mitra')
  } finally {
    deleting.value = false
  }
}
</script>

<style scoped>
/* ═══ HERO ═══ */
.mitra-hero {
  background: linear-gradient(135deg, #0f0c29 0%, #302b63 50%, #24243e 100%);
  border-radius: 16px; padding: 32px 36px 24px; margin-bottom: 24px;
  box-shadow: 0 4px 24px rgba(15,12,41,0.2);
}
.mitra-hero-content {
  display: flex; align-items: flex-start; justify-content: space-between;
  gap: 16px; flex-wrap: wrap; margin-bottom: 24px;
}
.mitra-hero-title { font-size: 1.6rem; font-weight: 800; color: #fff; margin: 0 0 4px; letter-spacing: -0.02em; }
.mitra-hero-sub { font-size: 0.85rem; color: rgba(255,255,255,0.5); margin: 0; }
.mitra-add-btn {
  display: inline-flex; align-items: center; gap: 8px;
  padding: 11px 24px; font-size: 0.85rem; font-weight: 700;
  border-radius: 12px; border: none; cursor: pointer;
  background: linear-gradient(135deg, #6366f1, #8b5cf6);
  color: white; box-shadow: 0 4px 20px rgba(99,102,241,0.35); white-space: nowrap;
}
.mitra-stats-bar { display: flex; gap: 28px; flex-wrap: wrap; padding-top: 18px; border-top: 1px solid rgba(255,255,255,0.08); }
.mitra-stat { display: flex; align-items: center; gap: 8px; }
.mitra-stat-dot { width: 8px; height: 8px; border-radius: 50%; }
.dot-total { background: #818cf8; box-shadow: 0 0 8px rgba(129,140,248,0.5); }
.dot-active { background: #4ade80; box-shadow: 0 0 8px rgba(74,222,128,0.5); }
.dot-inactive { background: #f87171; box-shadow: 0 0 8px rgba(248,113,113,0.4); }
.mitra-stat-label { font-size: 0.72rem; color: rgba(255,255,255,0.4); text-transform: uppercase; letter-spacing: 0.05em; }
.mitra-stat-value { font-size: 0.9rem; font-weight: 800; color: white; }

/* ═══ TOOLBAR ═══ */
.mitra-toolbar { margin-bottom: 20px; }
.mitra-search-wrap {
  position: relative; max-width: 400px;
}
.mitra-search-icon { position: absolute; left: 14px; top: 50%; transform: translateY(-50%); }
.mitra-search {
  width: 100%; padding: 11px 14px 11px 40px;
  border: 1.5px solid #e2e8f0; border-radius: 12px;
  font-size: 0.85rem; background: #fff; color: #1e293b;
  outline: none; box-sizing: border-box;
}
.mitra-search:focus { border-color: #6366f1; box-shadow: 0 0 0 3px rgba(99,102,241,0.1); }

/* ═══ GRID ═══ */
.mitra-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
  gap: 20px;
}

/* ═══ CARD ═══ */
.mitra-card {
  background: #fff; border-radius: 16px; border: 1px solid #e8ecf1;
  overflow: hidden; display: flex; flex-direction: column;
  box-shadow: 0 1px 4px rgba(0,0,0,0.04);
  transition: border-color 0.2s;
}
.mitra-card:hover { border-color: #c7d2fe; }
.mitra-card-inactive { opacity: 0.6; }

.mitra-card-top { display: flex; align-items: center; justify-content: space-between; padding: 20px 20px 0; }
.mitra-card-avatar {
  width: 56px; height: 56px; border-radius: 16px;
  display: flex; align-items: center; justify-content: center;
  font-size: 1.2rem; font-weight: 800; color: white;
  box-shadow: 0 4px 14px rgba(0,0,0,0.12);
}

.mitra-toggle-btn { background: none; border: none; cursor: pointer; padding: 0; }
.toggle-track { width: 40px; height: 22px; border-radius: 12px; background: #cbd5e1; position: relative; transition: background 0.2s; }
.toggle-thumb { width: 18px; height: 18px; border-radius: 50%; background: white; position: absolute; top: 2px; left: 2px; transition: transform 0.2s; box-shadow: 0 1px 4px rgba(0,0,0,0.15); }
.toggle-on .toggle-track { background: #22c55e; }
.toggle-on .toggle-thumb { transform: translateX(18px); }

.mitra-card-body { padding: 16px 20px 0; flex: 1; }
.mitra-card-name { font-size: 1.05rem; font-weight: 700; color: #0f172a; margin: 0 0 12px; }
.mitra-card-info { display: flex; flex-direction: column; gap: 8px; }
.mitra-info-row {
  display: flex; align-items: center; gap: 8px;
  font-size: 0.8rem; color: #64748b;
}
.mitra-info-row span {
  overflow: hidden; text-overflow: ellipsis; white-space: nowrap;
}

.mitra-card-footer {
  display: flex; align-items: center; justify-content: space-between;
  padding: 14px 20px; margin-top: 16px; border-top: 1px solid #f1f5f9;
}
.mitra-card-status { display: flex; align-items: center; gap: 6px; font-size: 0.72rem; font-weight: 600; text-transform: uppercase; letter-spacing: 0.04em; }
.st-dot { width: 7px; height: 7px; border-radius: 50%; }
.st-active { color: #16a34a; }
.st-active .st-dot { background: #22c55e; }
.st-inactive { color: #94a3b8; }
.st-inactive .st-dot { background: #cbd5e1; }
.mitra-card-btns { display: flex; gap: 6px; }
.mitra-icon-btn {
  width: 32px; height: 32px; border-radius: 8px;
  display: flex; align-items: center; justify-content: center;
  border: 1px solid #e8ecf1; background: #fff; cursor: pointer; color: #94a3b8;
}
.edit-btn:hover { background: #eef2ff; border-color: #6366f1; color: #6366f1; }
.del-btn:hover { background: #fef2f2; border-color: #ef4444; color: #ef4444; }
.mitra-card-date {
  padding: 8px 20px; background: #f8fafc;
  font-size: 0.7rem; color: #94a3b8; text-align: right;
}

/* ═══ SKELETON ═══ */
.mitra-card-skel { padding: 24px 20px; display: flex; flex-direction: column; gap: 12px; min-height: 200px; }
.skel-avatar { width: 56px; height: 56px; border-radius: 16px; }
.skel-name { width: 60%; height: 18px; border-radius: 6px; }
.skel-info { width: 80%; height: 14px; border-radius: 6px; }
.skel-info-sm { width: 50%; height: 14px; border-radius: 6px; }
.shimmer {
  background: linear-gradient(90deg, #e8ecf1 25%, #f1f5f9 50%, #e8ecf1 75%);
  background-size: 200% 100%; animation: shimmer 1.5s infinite;
}
@keyframes shimmer { 0% { background-position: 200% 0 } 100% { background-position: -200% 0 } }

/* ═══ EMPTY ═══ */
.mitra-empty {
  display: flex; flex-direction: column; align-items: center; justify-content: center;
  text-align: center; padding: 80px 32px; background: white;
  border-radius: 16px; border: 2px dashed #e2e8f0; min-height: 320px;
}
.mitra-empty-circle {
  width: 80px; height: 80px; border-radius: 50%;
  background: linear-gradient(135deg, #f1f5f9, #e2e8f0);
  display: flex; align-items: center; justify-content: center; margin-bottom: 24px;
}
.mitra-empty h3 { font-size: 1.15rem; font-weight: 700; color: #0f172a; margin: 0; }
.mitra-empty p { font-size: 0.85rem; color: #94a3b8; margin: 6px 0 0; }

/* ═══ MODAL ═══ */
.mitra-modal-overlay {
  position: fixed; inset: 0; background: rgba(0,0,0,.5);
  display: flex; align-items: center; justify-content: center;
  z-index: 1000; backdrop-filter: blur(6px); animation: fadeIn .2s ease;
}
@keyframes fadeIn { from { opacity: 0 } to { opacity: 1 } }
.mitra-modal {
  background: #fff; border-radius: 18px; width: 100%; max-width: 560px;
  box-shadow: 0 24px 80px rgba(0,0,0,.2); animation: slideUp .3s ease;
}
.mitra-modal-sm { max-width: 400px; }
@keyframes slideUp { from { transform: translateY(20px); opacity: 0 } to { transform: translateY(0); opacity: 1 } }
.mitra-modal-header {
  display: flex; align-items: center; justify-content: space-between;
  padding: 22px 28px; border-bottom: 1px solid #f1f5f9;
}
.mitra-modal-title-group { display: flex; align-items: center; gap: 12px; }
.mitra-modal-icon {
  width: 40px; height: 40px; border-radius: 12px;
  background: linear-gradient(135deg, #eef2ff, #e0e7ff);
  display: flex; align-items: center; justify-content: center; color: #6366f1;
}
.mitra-modal-header h3 { font-size: 1.1rem; font-weight: 700; margin: 0; color: #0f172a; }
.mitra-modal-close {
  width: 34px; height: 34px; border-radius: 10px;
  display: flex; align-items: center; justify-content: center;
  border: none; background: transparent; font-size: 1.4rem; color: #94a3b8; cursor: pointer;
}
.mitra-modal-close:hover { background: #f1f5f9; color: #0f172a; }
.mitra-modal-body { padding: 24px 28px; }
.mitra-form-group { margin-bottom: 18px; }
.mitra-form-group:last-of-type { margin-bottom: 0; }
.mitra-form-group label { display: block; font-size: 0.82rem; font-weight: 600; margin-bottom: 8px; color: #334155; }
.mitra-form-row { display: grid; grid-template-columns: 1fr 1fr; gap: 14px; }
.req { color: #ef4444; }
.mitra-input {
  width: 100%; padding: 11px 16px; border: 1.5px solid #e2e8f0; border-radius: 12px;
  font-size: 0.85rem; background: #fafbfc; color: #1e293b; outline: none;
  box-sizing: border-box; font-family: inherit;
}
.mitra-input:focus { border-color: #6366f1; background: #fff; box-shadow: 0 0 0 3px rgba(99,102,241,0.1); }
.mitra-input.is-error { border-color: #ef4444; background: #fef2f2; }
.mitra-field-error { color: #ef4444; font-size: 0.75rem; margin-top: 5px; font-weight: 500; }
.mitra-pw-wrap { position: relative; }
.mitra-pw-input { padding-right: 44px; }
.mitra-pw-toggle {
  position: absolute; right: 8px; top: 50%; transform: translateY(-50%);
  background: none; border: none; cursor: pointer; padding: 4px;
  display: flex; align-items: center; justify-content: center;
  border-radius: 6px; color: #94a3b8;
}
.mitra-pw-toggle:hover { background: #f1f5f9; }
.mitra-modal-footer { display: flex; justify-content: flex-end; gap: 10px; padding-top: 16px; }
.mitra-save-btn {
  display: inline-flex; align-items: center; gap: 6px;
  padding: 11px 28px; background: linear-gradient(135deg, #6366f1, #8b5cf6);
  color: #fff; border: none; border-radius: 12px; font-size: 0.85rem; font-weight: 600; cursor: pointer;
}
.mitra-save-btn:disabled { opacity: 0.6; cursor: not-allowed; }

/* Delete modal */
.mitra-delete-body { text-align: center; padding: 32px 28px 20px; }
.mitra-delete-icon-wrap {
  width: 60px; height: 60px; border-radius: 50%; background: rgba(239,68,68,0.08);
  display: flex; align-items: center; justify-content: center; margin: 0 auto 16px;
}
.mitra-delete-body h3 { font-size: 1.1rem; font-weight: 700; margin: 0 0 8px; color: #0f172a; }
.mitra-delete-body p { color: #64748b; font-size: 0.85rem; margin: 0; line-height: 1.6; }
.mitra-delete-footer { display: flex; justify-content: center; gap: 12px; padding: 0 28px 28px; }

.btn { padding: 11px 22px; border-radius: 12px; font-size: 0.85rem; font-weight: 600; cursor: pointer; border: none; }
.btn-secondary { background: #f1f5f9; color: #475569; }
.btn-secondary:hover { background: #e2e8f0; }
.btn-danger { background: linear-gradient(135deg, #ef4444, #dc2626); color: #fff; }
.btn-danger:disabled { opacity: 0.6; cursor: not-allowed; }
.spin-icon { animation: spin 0.8s linear infinite; }
@keyframes spin { to { transform: rotate(360deg) } }
.mitra-pw-hints { display: flex; gap: 10px; flex-wrap: wrap; margin-bottom: 16px; }
.pw-hint { font-size: 0.72rem; color: #94a3b8; padding: 3px 10px; border-radius: 6px; background: #f1f5f9; font-weight: 500; }
.pw-hint.pw-ok { color: #16a34a; background: #f0fdf4; }

@media (max-width: 768px) {
  .mitra-hero { padding: 24px 20px 18px; }
  .mitra-hero-content { flex-direction: column; align-items: flex-start; }
  .mitra-grid { grid-template-columns: 1fr; gap: 16px; }
  .mitra-form-row { grid-template-columns: 1fr; }
}
</style>
