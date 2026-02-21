<template>
  <div class="oc-page">
    <!-- Hero Header -->
    <div class="oc-hero">
      <div class="oc-hero-content">
        <div>
          <h1 class="oc-hero-title">Kategori Outlet</h1>
          <p class="oc-hero-sub">Kelola kategori untuk mengelompokkan outlet</p>
        </div>
        <button @click="openCreate" class="oc-add-btn">
          <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor"><line x1="12" y1="5" x2="12" y2="19" stroke-width="2" stroke-linecap="round"/><line x1="5" y1="12" x2="19" y2="12" stroke-width="2" stroke-linecap="round"/></svg>
          Tambah Kategori
        </button>
      </div>
      <div class="oc-stats-bar">
        <div class="oc-stat">
          <span class="oc-stat-dot dot-total"></span>
          <span class="oc-stat-label">Total</span>
          <span class="oc-stat-value">{{ categories.length }}</span>
        </div>
        <div class="oc-stat">
          <span class="oc-stat-dot dot-active"></span>
          <span class="oc-stat-label">Aktif</span>
          <span class="oc-stat-value">{{ categories.filter(c => c.is_active).length }}</span>
        </div>
        <div class="oc-stat">
          <span class="oc-stat-dot dot-inactive"></span>
          <span class="oc-stat-label">Nonaktif</span>
          <span class="oc-stat-value">{{ categories.filter(c => !c.is_active).length }}</span>
        </div>
      </div>
    </div>

    <!-- Loading -->
    <div v-if="loading" class="oc-grid">
      <div v-for="n in 6" :key="n" class="oc-card oc-card-skeleton">
        <div class="skel-icon shimmer"></div>
        <div class="skel-title shimmer"></div>
        <div class="skel-desc shimmer"></div>
        <div class="skel-footer shimmer"></div>
      </div>
    </div>

    <!-- Category Grid -->
    <div v-else-if="categories.length" class="oc-grid">
      <div v-for="(cat, idx) in categories" :key="cat.id" class="oc-card" :class="{ 'oc-card-inactive': !cat.is_active }">
        <div class="oc-card-top">
          <div class="oc-card-icon" :style="iconBg(idx)">{{ cat.name.charAt(0) }}</div>
          <div class="oc-card-actions">
            <button @click="toggleActive(cat)" class="oc-toggle-btn" :class="cat.is_active ? 'toggle-on' : 'toggle-off'" :title="cat.is_active ? 'Nonaktifkan' : 'Aktifkan'">
              <div class="toggle-track"><div class="toggle-thumb"></div></div>
            </button>
          </div>
        </div>
        <div class="oc-card-body">
          <h3 class="oc-card-name">{{ cat.name }}</h3>
          <code class="oc-card-slug">{{ cat.slug }}</code>
          <p class="oc-card-desc">{{ cat.description || 'Belum ada deskripsi' }}</p>
        </div>
        <div class="oc-card-footer">
          <span class="oc-card-status" :class="cat.is_active ? 'st-active' : 'st-inactive'">
            <span class="st-dot"></span>
            {{ cat.is_active ? 'Aktif' : 'Nonaktif' }}
          </span>
          <div class="oc-card-btns">
            <button @click="openEdit(cat)" class="oc-icon-btn edit-btn" title="Edit">
              <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor"><path d="M11 4H4a2 2 0 0 0-2 2v14a2 2 0 0 0 2 2h14a2 2 0 0 0 2-2v-7" stroke-width="2"/><path d="M18.5 2.5a2.121 2.121 0 0 1 3 3L12 15l-4 1 1-4 9.5-9.5z" stroke-width="2"/></svg>
            </button>
            <button @click="confirmDelete(cat)" class="oc-icon-btn del-btn" title="Hapus">
              <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor"><polyline points="3 6 5 6 21 6" stroke-width="2"/><path d="M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6m3 0V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2" stroke-width="2"/></svg>
            </button>
          </div>
        </div>
      </div>
    </div>

    <!-- Empty state -->
    <div v-else class="oc-empty">
      <div class="oc-empty-circle">
        <svg width="36" height="36" viewBox="0 0 24 24" fill="none" stroke="#94a3b8" stroke-width="1.5"><path d="M22 19a2 2 0 0 1-2 2H4a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h5l2 3h9a2 2 0 0 1 2 2z"/></svg>
      </div>
      <h3>Belum ada kategori</h3>
      <p>Buat kategori pertama untuk mengelompokkan outlet.</p>
      <button @click="openCreate" class="oc-add-btn oc-add-btn-sm">+ Tambah Kategori</button>
    </div>

    <!-- Modal Create/Edit -->
    <Teleport to="body">
      <div v-if="showModal" class="oc-modal-overlay" @click.self="closeModal">
        <div class="oc-modal" @click.stop>
          <div class="oc-modal-header">
            <div class="oc-modal-title-group">
              <div class="oc-modal-icon">
                <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor"><path d="M22 19a2 2 0 0 1-2 2H4a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h5l2 3h9a2 2 0 0 1 2 2z" stroke-width="2"/></svg>
              </div>
              <h3>{{ editItem ? 'Edit Kategori' : 'Tambah Kategori Baru' }}</h3>
            </div>
            <button @click="closeModal" class="oc-modal-close">&times;</button>
          </div>
          <form @submit.prevent="handleSave" class="oc-modal-body">
            <div class="oc-form-group">
              <label>Nama Kategori <span class="req">*</span></label>
              <input v-model="formData.name" type="text" class="oc-input" :class="{ 'is-error': formErrors.name }" placeholder="Contoh: Franchise" maxlength="100" />
              <div v-if="formErrors.name" class="oc-field-error">{{ formErrors.name }}</div>
            </div>
            <div class="oc-form-group">
              <label>Deskripsi</label>
              <textarea v-model="formData.description" class="oc-input oc-textarea" placeholder="Deskripsi singkat kategori..." rows="3" maxlength="500"></textarea>
              <div class="oc-char-count">{{ (formData.description || '').length }}/500</div>
            </div>
            <div class="oc-modal-footer">
              <button type="button" @click="closeModal" class="btn btn-secondary">Batal</button>
              <button type="submit" class="oc-save-btn" :disabled="saving">
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
      <div v-if="deleteTarget" class="oc-modal-overlay" @click.self="deleteTarget = null">
        <div class="oc-modal oc-modal-sm" @click.stop>
          <div class="oc-delete-body">
            <div class="oc-delete-icon-wrap">
              <svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="#ef4444"><polyline points="3 6 5 6 21 6" stroke-width="2"/><path d="M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6m3 0V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2" stroke-width="2"/></svg>
            </div>
            <h3>Hapus Kategori</h3>
            <p>Yakin ingin menghapus <strong>{{ deleteTarget.name }}</strong>?<br/>Kategori yang masih digunakan outlet tidak bisa dihapus.</p>
          </div>
          <div class="oc-delete-footer">
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
import { ref, reactive, onMounted } from 'vue'
import { outletCategoryApi } from '../../services/api'
import { useToastStore } from '../../stores/toast'

const toast = useToastStore()
const categories = ref([])
const loading = ref(false)
const showModal = ref(false)
const editItem = ref(null)
const saving = ref(false)
const deleteTarget = ref(null)
const deleting = ref(false)

const formData = reactive({ name: '', description: '' })
const formErrors = reactive({})

const iconGradients = [
  'linear-gradient(135deg, #667eea, #764ba2)',
  'linear-gradient(135deg, #f093fb, #f5576c)',
  'linear-gradient(135deg, #4facfe, #00f2fe)',
  'linear-gradient(135deg, #43e97b, #38f9d7)',
  'linear-gradient(135deg, #fa709a, #fee140)',
  'linear-gradient(135deg, #a18cd1, #fbc2eb)',
  'linear-gradient(135deg, #fd9644, #fc5c65)',
  'linear-gradient(135deg, #0abde3, #341f97)',
]

function iconBg(idx) {
  return { background: iconGradients[idx % iconGradients.length] }
}

onMounted(() => loadCategories())

async function loadCategories() {
  loading.value = true
  try {
    const { data } = await outletCategoryApi.list()
    categories.value = data.data || []
  } catch (e) {
    toast.error('Gagal memuat kategori')
  } finally {
    loading.value = false
  }
}

function openCreate() {
  editItem.value = null
  Object.assign(formData, { name: '', description: '' })
  Object.keys(formErrors).forEach(k => delete formErrors[k])
  showModal.value = true
}

function openEdit(cat) {
  editItem.value = cat
  Object.assign(formData, { name: cat.name, description: cat.description })
  Object.keys(formErrors).forEach(k => delete formErrors[k])
  showModal.value = true
}

function closeModal() {
  showModal.value = false
  editItem.value = null
}

async function handleSave() {
  Object.keys(formErrors).forEach(k => delete formErrors[k])
  if (!formData.name?.trim()) {
    formErrors.name = 'Nama kategori wajib diisi'
    return
  }
  saving.value = true
  try {
    const payload = { ...formData }
    if (editItem.value) {
      const { data } = await outletCategoryApi.update(editItem.value.id, payload)
      toast.success(data.message || 'Kategori berhasil diupdate')
    } else {
      const { data } = await outletCategoryApi.create(payload)
      toast.success(data.message || 'Kategori berhasil ditambah')
    }
    closeModal()
    loadCategories()
  } catch (e) {
    if (e.response?.status === 422) {
      const errs = e.response.data.errors || []
      errs.forEach(err => { formErrors[err.field] = err.message })
    } else {
      toast.error(e.response?.data?.error || 'Gagal menyimpan kategori')
    }
  } finally {
    saving.value = false
  }
}

async function toggleActive(cat) {
  try {
    const { data } = await outletCategoryApi.toggle(cat.id)
    toast.success(data.message)
    cat.is_active = data.data.is_active
  } catch { toast.error('Gagal mengubah status') }
}

function confirmDelete(cat) { deleteTarget.value = cat }

async function doDelete() {
  deleting.value = true
  try {
    const { data } = await outletCategoryApi.delete(deleteTarget.value.id)
    toast.success(data.message || 'Kategori berhasil dihapus')
    deleteTarget.value = null
    loadCategories()
  } catch (e) {
    toast.error(e.response?.data?.error || 'Gagal menghapus kategori')
  } finally { deleting.value = false }
}
</script>

<style scoped>
/* ═══ HERO ═══ */
.oc-hero {
  background: linear-gradient(135deg, #0f0c29 0%, #302b63 50%, #24243e 100%);
  border-radius: 16px;
  padding: 32px 36px 24px;
  margin-bottom: 28px;
  box-shadow: 0 4px 24px rgba(15,12,41,0.2);
}
.oc-hero-content {
  display: flex; align-items: flex-start; justify-content: space-between;
  gap: 16px; flex-wrap: wrap; margin-bottom: 24px;
}
.oc-hero-title {
  font-size: 1.6rem; font-weight: 800; color: #fff;
  margin: 0 0 4px; letter-spacing: -0.02em;
}
.oc-hero-sub { font-size: 0.85rem; color: rgba(255,255,255,0.5); margin: 0; }
.oc-add-btn {
  display: inline-flex; align-items: center; gap: 8px;
  padding: 11px 24px; font-size: 0.85rem; font-weight: 700;
  border-radius: 12px; border: none; cursor: pointer;
  background: linear-gradient(135deg, #6366f1, #8b5cf6);
  color: white; box-shadow: 0 4px 20px rgba(99,102,241,0.35);
  white-space: nowrap;
}
.oc-add-btn-sm { margin-top: 20px; font-size: 0.82rem; padding: 10px 20px; }

/* Stats bar */
.oc-stats-bar {
  display: flex; gap: 28px; flex-wrap: wrap;
  padding-top: 18px; border-top: 1px solid rgba(255,255,255,0.08);
}
.oc-stat { display: flex; align-items: center; gap: 8px; }
.oc-stat-dot { width: 8px; height: 8px; border-radius: 50%; }
.dot-total { background: #818cf8; box-shadow: 0 0 8px rgba(129,140,248,0.5); }
.dot-active { background: #4ade80; box-shadow: 0 0 8px rgba(74,222,128,0.5); }
.dot-inactive { background: #f87171; box-shadow: 0 0 8px rgba(248,113,113,0.4); }
.oc-stat-label { font-size: 0.72rem; color: rgba(255,255,255,0.4); text-transform: uppercase; letter-spacing: 0.05em; }
.oc-stat-value { font-size: 0.9rem; font-weight: 800; color: white; }

/* ═══ GRID ═══ */
.oc-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(280px, 1fr));
  gap: 20px;
}

/* ═══ CARD ═══ */
.oc-card {
  background: #fff;
  border-radius: 16px;
  border: 1px solid #e8ecf1;
  overflow: hidden;
  display: flex;
  flex-direction: column;
  box-shadow: 0 1px 4px rgba(0,0,0,0.04);
}
.oc-card-inactive { opacity: 0.6; }

.oc-card-top {
  display: flex; align-items: center; justify-content: space-between;
  padding: 20px 20px 0;
}
.oc-card-icon {
  width: 52px; height: 52px; border-radius: 14px;
  display: flex; align-items: center; justify-content: center;
  font-size: 1.4rem; font-weight: 800; color: white;
  box-shadow: 0 4px 14px rgba(0,0,0,0.12);
}

/* Toggle switch */
.oc-toggle-btn {
  background: none; border: none; cursor: pointer; padding: 0;
}
.toggle-track {
  width: 40px; height: 22px; border-radius: 12px;
  background: #cbd5e1; position: relative;
  transition: background 0.2s;
}
.toggle-thumb {
  width: 18px; height: 18px; border-radius: 50%;
  background: white; position: absolute;
  top: 2px; left: 2px;
  transition: transform 0.2s;
  box-shadow: 0 1px 4px rgba(0,0,0,0.15);
}
.toggle-on .toggle-track { background: #22c55e; }
.toggle-on .toggle-thumb { transform: translateX(18px); }

.oc-card-body {
  padding: 16px 20px 0;
  flex: 1;
}
.oc-card-name {
  font-size: 1.05rem; font-weight: 700; color: #0f172a;
  margin: 0 0 6px; letter-spacing: -0.01em;
}
.oc-card-slug {
  display: inline-block;
  background: linear-gradient(135deg, #eef2ff, #e0e7ff);
  color: #4f46e5; padding: 3px 10px; border-radius: 6px;
  font-size: 0.72rem; font-weight: 600;
  margin-bottom: 10px;
}
.oc-card-desc {
  font-size: 0.8rem; color: #64748b; line-height: 1.55;
  margin: 0;
  display: -webkit-box; -webkit-line-clamp: 2; -webkit-box-orient: vertical;
  overflow: hidden; min-height: 36px;
}

.oc-card-footer {
  display: flex; align-items: center; justify-content: space-between;
  padding: 14px 20px;
  margin-top: 14px;
  border-top: 1px solid #f1f5f9;
}
.oc-card-status {
  display: flex; align-items: center; gap: 6px;
  font-size: 0.72rem; font-weight: 600; text-transform: uppercase;
  letter-spacing: 0.04em;
}
.st-dot { width: 7px; height: 7px; border-radius: 50%; }
.st-active { color: #16a34a; }
.st-active .st-dot { background: #22c55e; }
.st-inactive { color: #94a3b8; }
.st-inactive .st-dot { background: #cbd5e1; }

.oc-card-btns { display: flex; gap: 6px; }
.oc-icon-btn {
  width: 32px; height: 32px; border-radius: 8px;
  display: flex; align-items: center; justify-content: center;
  border: 1px solid #e8ecf1; background: #fff;
  cursor: pointer; color: #94a3b8;
}
.edit-btn:hover { background: #eef2ff; border-color: #6366f1; color: #6366f1; }
.del-btn:hover { background: #fef2f2; border-color: #ef4444; color: #ef4444; }

/* ═══ SKELETON ═══ */
.oc-card-skeleton {
  padding: 24px 20px; display: flex; flex-direction: column; gap: 12px;
  min-height: 180px;
}
.skel-icon { width: 52px; height: 52px; border-radius: 14px; }
.skel-title { width: 60%; height: 18px; border-radius: 6px; }
.skel-desc { width: 90%; height: 14px; border-radius: 6px; }
.skel-footer { width: 40%; height: 14px; border-radius: 6px; margin-top: auto; }
.shimmer {
  background: linear-gradient(90deg, #e8ecf1 25%, #f1f5f9 50%, #e8ecf1 75%);
  background-size: 200% 100%;
  animation: shimmer 1.5s infinite;
}
@keyframes shimmer { 0% { background-position: 200% 0 } 100% { background-position: -200% 0 } }

/* ═══ EMPTY ═══ */
.oc-empty {
  display: flex; flex-direction: column; align-items: center;
  justify-content: center; text-align: center;
  padding: 80px 32px; background: white;
  border-radius: 16px; border: 2px dashed #e2e8f0;
  min-height: 320px;
}
.oc-empty-circle {
  width: 80px; height: 80px; border-radius: 50%;
  background: linear-gradient(135deg, #f1f5f9, #e2e8f0);
  display: flex; align-items: center; justify-content: center;
  margin-bottom: 24px;
}
.oc-empty h3 { font-size: 1.15rem; font-weight: 700; color: #0f172a; margin: 0; }
.oc-empty p { font-size: 0.85rem; color: #94a3b8; margin: 6px 0 0; }

/* ═══ MODAL ═══ */
.oc-modal-overlay {
  position: fixed; inset: 0;
  background: rgba(0,0,0,.5);
  display: flex; align-items: center; justify-content: center;
  z-index: 1000;
  backdrop-filter: blur(6px); -webkit-backdrop-filter: blur(6px);
  animation: fadeIn .2s ease;
}
@keyframes fadeIn { from { opacity: 0 } to { opacity: 1 } }
.oc-modal {
  background: #fff; border-radius: 18px;
  width: 100%; max-width: 520px;
  box-shadow: 0 24px 80px rgba(0,0,0,.2);
  animation: slideUp .3s ease;
}
.oc-modal-sm { max-width: 400px; }
@keyframes slideUp { from { transform: translateY(20px); opacity: 0 } to { transform: translateY(0); opacity: 1 } }

.oc-modal-header {
  display: flex; align-items: center; justify-content: space-between;
  padding: 22px 28px; border-bottom: 1px solid #f1f5f9;
}
.oc-modal-title-group { display: flex; align-items: center; gap: 12px; }
.oc-modal-icon {
  width: 40px; height: 40px; border-radius: 12px;
  background: linear-gradient(135deg, #eef2ff, #e0e7ff);
  display: flex; align-items: center; justify-content: center;
  color: #6366f1;
}
.oc-modal-header h3 { font-size: 1.1rem; font-weight: 700; margin: 0; color: #0f172a; }
.oc-modal-close {
  width: 34px; height: 34px; border-radius: 10px;
  display: flex; align-items: center; justify-content: center;
  border: none; background: transparent;
  font-size: 1.4rem; color: #94a3b8; cursor: pointer;
}
.oc-modal-close:hover { background: #f1f5f9; color: #0f172a; }

.oc-modal-body { padding: 24px 28px; }
.oc-form-group { margin-bottom: 18px; }
.oc-form-group:last-of-type { margin-bottom: 0; }
.oc-form-group label {
  display: block; font-size: 0.82rem; font-weight: 600;
  margin-bottom: 8px; color: #334155;
}
.req { color: #ef4444; }
.oc-input {
  width: 100%; padding: 11px 16px;
  border: 1.5px solid #e2e8f0; border-radius: 12px;
  font-size: 0.85rem; background: #fafbfc;
  color: #1e293b; outline: none; box-sizing: border-box;
  font-family: inherit;
}
.oc-input:focus { border-color: #6366f1; background: #fff; box-shadow: 0 0 0 3px rgba(99,102,241,0.1); }
.oc-input.is-error { border-color: #ef4444; background: #fef2f2; }
.oc-textarea { resize: vertical; min-height: 80px; }
.oc-char-count { text-align: right; font-size: 0.7rem; color: #94a3b8; margin-top: 4px; }
.oc-field-error { color: #ef4444; font-size: 0.75rem; margin-top: 5px; font-weight: 500; }

.oc-modal-footer {
  display: flex; justify-content: flex-end; gap: 10px;
  padding-top: 16px;
}
.oc-save-btn {
  display: inline-flex; align-items: center; gap: 6px;
  padding: 11px 28px;
  background: linear-gradient(135deg, #6366f1, #8b5cf6);
  color: #fff; border: none; border-radius: 12px;
  font-size: 0.85rem; font-weight: 600; cursor: pointer;
}
.oc-save-btn:disabled { opacity: 0.6; cursor: not-allowed; }

/* Delete modal */
.oc-delete-body {
  text-align: center; padding: 32px 28px 20px;
}
.oc-delete-icon-wrap {
  width: 60px; height: 60px; border-radius: 50%;
  background: rgba(239,68,68,0.08);
  display: flex; align-items: center; justify-content: center;
  margin: 0 auto 16px;
}
.oc-delete-body h3 { font-size: 1.1rem; font-weight: 700; margin: 0 0 8px; color: #0f172a; }
.oc-delete-body p { color: #64748b; font-size: 0.85rem; margin: 0; line-height: 1.6; }
.oc-delete-footer {
  display: flex; justify-content: center; gap: 12px;
  padding: 0 28px 28px;
}

/* Buttons */
.btn {
  padding: 11px 22px; border-radius: 12px; font-size: 0.85rem;
  font-weight: 600; cursor: pointer; border: none;
}
.btn-secondary { background: #f1f5f9; color: #475569; }
.btn-secondary:hover { background: #e2e8f0; }
.btn-danger { background: linear-gradient(135deg, #ef4444, #dc2626); color: #fff; }
.btn-danger:disabled { opacity: 0.6; cursor: not-allowed; }

.spin-icon { animation: spin 0.8s linear infinite; }
@keyframes spin { to { transform: rotate(360deg) } }

/* ═══ RESPONSIVE ═══ */
@media (max-width: 768px) {
  .oc-hero { padding: 24px 20px 18px; }
  .oc-hero-content { flex-direction: column; align-items: flex-start; }
  .oc-grid { grid-template-columns: 1fr; gap: 16px; }
}
</style>
