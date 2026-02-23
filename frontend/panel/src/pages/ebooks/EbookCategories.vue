<template>
  <div class="animate-in">
    <div class="cat-hero">
      <div class="cat-hero-row">
        <div>
          <h1 class="cat-hero-title"><i class="ri-price-tag-3-line"></i> Kategori Ebook</h1>
          <p class="cat-hero-sub">Kelola kategori master untuk ebook</p>
        </div>
        <button @click="openModal()" class="cat-hero-btn">
          <i class="ri-add-line"></i>
          Tambah Kategori
        </button>
      </div>
    </div>

    <!-- Loading -->
    <div v-if="loading" class="cat-loading">Memuat...</div>

    <!-- Empty -->
    <div v-else-if="!cats.length" class="cat-empty">
      <div style="font-size:2.5rem;margin-bottom:12px"><i class="ri-folder-open-line" style="font-size:2.5rem;color:#94a3b8"></i></div>
      <h3>Belum ada kategori</h3>
      <p>Buat kategori pertama untuk mengelompokkan ebook</p>
    </div>

    <!-- Table -->
    <div v-else class="cat-table-wrap">
      <table class="cat-table">
        <thead>
          <tr>
            <th>Nama</th>
            <th>Slug</th>
            <th>Status</th>
            <th style="width:140px">Aksi</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="cat in cats" :key="cat.id">
            <td><strong>{{ cat.name }}</strong></td>
            <td><code>{{ cat.slug }}</code></td>
            <td>
              <span class="cat-badge" :class="cat.is_active ? 'badge-active' : 'badge-inactive'">
                {{ cat.is_active ? 'Aktif' : 'Nonaktif' }}
              </span>
            </td>
            <td>
              <div class="cat-actions">
                <button @click="openModal(cat)" class="act-btn act-edit" title="Edit"><i class="ri-pencil-line"></i></button>
                <button @click="toggleActive(cat)" class="act-btn" :class="cat.is_active ? 'act-off' : 'act-on'" :title="cat.is_active ? 'Nonaktifkan' : 'Aktifkan'">
                  <i :class="cat.is_active ? 'ri-lock-line' : 'ri-lock-unlock-line'"></i>
                </button>
                <button @click="confirmDelete(cat)" class="act-btn act-del" title="Hapus"><i class="ri-delete-bin-line"></i></button>
              </div>
            </td>
          </tr>
        </tbody>
      </table>
    </div>

    <!-- Modal create/edit -->
    <div v-if="showModal" class="modal-overlay" @click.self="showModal = false">
      <div class="modal-content" style="max-width:480px">
        <div class="modal-header">
          <h3>{{ editTarget ? 'Edit Kategori' : 'Tambah Kategori' }}</h3>
          <button @click="showModal = false" class="modal-close">&times;</button>
        </div>
        <form @submit.prevent="saveCategory" class="modal-body">
          <div class="form-group">
            <label class="form-label">Nama Kategori <span class="required">*</span></label>
            <input v-model="catForm.name" type="text" class="form-input" placeholder="Contoh: Bisnis, Marketing..." required />
          </div>
          <div class="form-group form-row">
            <label class="form-label">Status Aktif</label>
            <label class="toggle-switch">
              <input type="checkbox" v-model="catForm.is_active" />
              <span class="toggle-slider"></span>
            </label>
          </div>
          <div class="modal-actions">
            <button type="button" @click="showModal = false" class="btn btn-secondary">Batal</button>
            <button type="submit" class="btn btn-primary" :disabled="saving">{{ saving ? 'Menyimpan...' : 'Simpan' }}</button>
          </div>
        </form>
      </div>
    </div>

    <!-- Delete confirm -->
    <div v-if="deleteTarget" class="modal-overlay" @click.self="deleteTarget = null">
      <div class="modal-content" style="max-width:400px">
        <div class="modal-body" style="text-align:center;padding:2rem">
          <div style="margin-bottom:12px"><i class="ri-delete-bin-line" style="font-size:2.5rem;color:#ef4444"></i></div>
          <h3>Hapus kategori <strong>{{ deleteTarget.name }}</strong>?</h3>
          <p style="color:#64748b;font-size:0.85rem;margin-top:8px">Kategori yang sedang digunakan tidak bisa dihapus.</p>
        </div>
        <div class="modal-actions" style="justify-content:center;padding-bottom:1.5rem">
          <button @click="deleteTarget = null" class="btn btn-secondary">Batal</button>
          <button @click="doDelete" class="btn btn-danger" :disabled="deleting">{{ deleting ? 'Menghapus...' : 'Hapus' }}</button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { ebookCategoryApi } from '../../services/api'
import { useToastStore } from '../../stores/toast'

const toast = useToastStore()
const cats = ref([])
const loading = ref(false)
const showModal = ref(false)
const editTarget = ref(null)
const saving = ref(false)
const deleteTarget = ref(null)
const deleting = ref(false)

const catForm = ref({ name: '', is_active: true })

onMounted(load)

async function load() {
  loading.value = true
  try {
    const { data } = await ebookCategoryApi.list()
    cats.value = data.data || []
  } catch { toast.error('Gagal memuat kategori') }
  finally { loading.value = false }
}

function openModal(cat = null) {
  editTarget.value = cat
  catForm.value = cat ? { name: cat.name, is_active: cat.is_active } : { name: '', is_active: true }
  showModal.value = true
}

async function saveCategory() {
  if (!catForm.value.name.trim()) {
    toast.error('Nama kategori wajib diisi')
    return
  }
  saving.value = true
  try {
    if (editTarget.value) {
      await ebookCategoryApi.update(editTarget.value.id, catForm.value)
      toast.success('Kategori berhasil diupdate')
    } else {
      await ebookCategoryApi.create(catForm.value)
      toast.success('Kategori berhasil ditambahkan')
    }
    showModal.value = false
    load()
  } catch (err) {
    toast.error(err.response?.data?.error || 'Gagal menyimpan kategori')
  } finally { saving.value = false }
}

async function toggleActive(cat) {
  try {
    await ebookCategoryApi.toggle(cat.id)
    toast.success(cat.is_active ? 'Kategori dinonaktifkan' : 'Kategori diaktifkan')
    load()
  } catch { toast.error('Gagal mengubah status') }
}

function confirmDelete(cat) { deleteTarget.value = cat }

async function doDelete() {
  deleting.value = true
  try {
    await ebookCategoryApi.delete(deleteTarget.value.id)
    toast.success('Kategori berhasil dihapus')
    deleteTarget.value = null
    load()
  } catch (err) {
    toast.error(err.response?.data?.error || 'Gagal menghapus kategori')
  } finally { deleting.value = false }
}
</script>

<style scoped>
.cat-hero { background: linear-gradient(135deg, #0f0c29 0%, #302b63 50%, #24243e 100%); border-radius: 16px; padding: 32px 36px; margin-bottom: 24px; }
.cat-hero-row { display: flex; align-items: center; justify-content: space-between; gap: 16px; flex-wrap: wrap; }
.cat-hero-title { font-size: 1.6rem; font-weight: 800; color: #fff; }
.cat-hero-sub { font-size: 0.85rem; color: rgba(255,255,255,0.5); margin-top: 6px; }
.cat-hero-btn { display: inline-flex; align-items: center; gap: 8px; padding: 11px 24px; font-size: 0.85rem; font-weight: 700; border-radius: 12px; border: none; cursor: pointer; background: linear-gradient(135deg, #f59e0b, #ef4444); color: white; box-shadow: 0 4px 20px rgba(245,158,11,0.3); }

.cat-loading { text-align: center; padding: 60px; color: #94a3b8; font-size: 0.9rem; }
.cat-empty { display: flex; flex-direction: column; align-items: center; padding: 60px; background: white; border-radius: 16px; border: 2px dashed #e2e8f0; text-align: center; }
.cat-empty h3 { font-size: 1.1rem; font-weight: 700; color: #0f172a; margin: 0; }
.cat-empty p { font-size: 0.85rem; color: #94a3b8; margin-top: 6px; }

.cat-table-wrap { background: white; border-radius: 16px; overflow: hidden; border: 1px solid #e8ecf1; }
.cat-table { width: 100%; border-collapse: collapse; }
.cat-table th { padding: 14px 20px; text-align: left; font-size: 0.75rem; font-weight: 700; text-transform: uppercase; color: #64748b; background: #f8fafc; border-bottom: 1px solid #e8ecf1; }
.cat-table td { padding: 14px 20px; font-size: 0.88rem; color: #1e293b; border-bottom: 1px solid #f1f5f9; }
.cat-table code { background: #f1f5f9; padding: 3px 8px; border-radius: 6px; font-size: 0.78rem; color: #6366f1; }

.cat-badge { padding: 4px 12px; border-radius: 8px; font-size: 0.72rem; font-weight: 700; }
.badge-active { background: #dcfce7; color: #16a34a; }
.badge-inactive { background: #fee2e2; color: #dc2626; }

.cat-actions { display: flex; gap: 6px; }
.act-btn { width: 32px; height: 32px; border: none; border-radius: 8px; cursor: pointer; background: #f1f5f9; font-size: 0.85rem; display: flex; align-items: center; justify-content: center; transition: 0.2s; }
.act-btn:hover { background: #e2e8f0; }
.act-del:hover { background: #fee2e2; }

.modal-overlay { position: fixed; inset: 0; background: rgba(0,0,0,0.5); z-index: 100; display: flex; align-items: center; justify-content: center; }
.modal-content { background: white; border-radius: 16px; overflow: hidden; width: 90%; animation: fadeIn 0.2s; }
.modal-header { display: flex; justify-content: space-between; align-items: center; padding: 20px 24px; border-bottom: 1px solid #f1f5f9; }
.modal-header h3 { font-size: 1.1rem; font-weight: 700; margin: 0; }
.modal-close { background: none; border: none; font-size: 1.5rem; cursor: pointer; color: #94a3b8; }
.modal-body { padding: 24px; }
.modal-actions { display: flex; justify-content: flex-end; gap: 12px; padding: 0 24px 24px; }
@keyframes fadeIn { from { opacity: 0; transform: translateY(-10px); } }

.form-group { margin-bottom: 20px; }
.form-label { display: block; font-size: 0.82rem; font-weight: 600; color: #334155; margin-bottom: 8px; }
.required { color: #ef4444; }
.form-input { width: 100%; padding: 12px 16px; border: 1.5px solid #e2e8f0; border-radius: 12px; font-size: 0.85rem; color: #1e293b; font-family: inherit; box-sizing: border-box; }
.form-input:focus { border-color: #818cf8; outline: none; box-shadow: 0 0 0 3px rgba(129,140,248,0.12); }
.form-row { display: flex; align-items: center; justify-content: space-between; }

.toggle-switch { position: relative; display: inline-block; width: 48px; height: 26px; }
.toggle-switch input { opacity: 0; width: 0; height: 0; }
.toggle-slider { position: absolute; cursor: pointer; inset: 0; background: #cbd5e1; border-radius: 26px; transition: 0.3s; }
.toggle-slider::before { content: ''; position: absolute; height: 20px; width: 20px; left: 3px; bottom: 3px; background: white; border-radius: 50%; transition: 0.3s; }
.toggle-switch input:checked + .toggle-slider { background: #6366f1; }
.toggle-switch input:checked + .toggle-slider::before { transform: translateX(22px); }

.btn { padding: 11px 24px; border-radius: 12px; font-size: 0.85rem; font-weight: 700; border: none; cursor: pointer; }
.btn-primary { background: linear-gradient(135deg, #6366f1, #8b5cf6); color: white; box-shadow: 0 4px 12px rgba(99,102,241,0.3); }
.btn-secondary { background: #f1f5f9; color: #475569; }
.btn-danger { background: #ef4444; color: white; }
.btn:disabled { opacity: 0.5; cursor: not-allowed; }
</style>
