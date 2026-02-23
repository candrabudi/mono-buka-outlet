<template>
  <div class="animate-in">
    <div class="form-hero">
      <h1 class="form-hero-title"><i :class="isEdit ? 'ri-pencil-line' : 'ri-book-open-line'"></i> {{ isEdit ? 'Edit Ebook' : 'Tambah Ebook' }}</h1>
      <p class="form-hero-sub">{{ isEdit ? 'Perbarui informasi ebook' : 'Tambahkan ebook baru untuk mitra' }}</p>
    </div>

    <form @submit.prevent="submit" class="ebook-form-card">
      <!-- Title -->
      <div class="form-group">
        <label class="form-label">Judul Ebook <span class="required">*</span></label>
        <input v-model="form.title" type="text" class="form-input" placeholder="Judul ebook..." />
      </div>

      <!-- Author -->
      <div class="form-group">
        <label class="form-label">Penulis</label>
        <input v-model="form.author" type="text" class="form-input" placeholder="Nama penulis..." />
      </div>

      <!-- Categories (multi-select) -->
      <div class="form-group">
        <label class="form-label">Kategori</label>
        <div class="cat-checkboxes" v-if="categories.length">
          <label v-for="cat in categories" :key="cat.id" class="cat-checkbox-item">
            <input type="checkbox" :value="cat.id" v-model="selectedCategoryIds" />
            <span class="cat-checkbox-label">{{ cat.name }}</span>
          </label>
        </div>
        <p v-else class="cat-hint">Belum ada kategori. <router-link to="/ebook-categories">Buat kategori</router-link></p>
      </div>

      <!-- Description — Rich Text Editor -->
      <div class="form-group">
        <label class="form-label">Deskripsi</label>
        <div class="editor-wrapper">
          <QuillEditor
            v-model:content="form.description"
            contentType="html"
            theme="snow"
            :toolbar="toolbarOptions"
            placeholder="Tulis deskripsi ebook..."
            style="min-height: 200px;"
          />
        </div>
      </div>

      <!-- Price -->
      <div class="form-group">
        <label class="form-label">Harga</label>
        <div class="price-input-wrap">
          <span class="price-prefix">Rp</span>
          <input
            :value="formatPriceInput(form.price)"
            @input="onPriceInput"
            type="text"
            class="form-input price-input"
            placeholder="0"
          />
        </div>
        <span class="form-hint">Isi 0 untuk ebook gratis</span>
      </div>

      <!-- Cover Upload -->
      <div class="form-group">
        <label class="form-label">Cover Image</label>
        <div class="upload-area" @click="$refs.coverInput.click()">
          <div v-if="uploadingCover" class="upload-placeholder">
            <div class="upload-spinner"></div>
            <span>Mengupload cover...</span>
          </div>
          <img v-else-if="form.cover_url" :src="form.cover_url" class="upload-preview" />
          <div v-else class="upload-placeholder">
            <svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="#94a3b8"><rect x="3" y="3" width="18" height="18" rx="2" stroke-width="2"/><circle cx="8.5" cy="8.5" r="1.5" stroke-width="2"/><polyline points="21 15 16 10 5 21" stroke-width="2"/></svg>
            <span>Klik untuk upload cover</span>
          </div>
        </div>
        <input ref="coverInput" type="file" accept="image/*" @change="uploadCover" style="display:none" />
      </div>

      <!-- PDF Upload -->
      <div class="form-group">
        <label class="form-label">File PDF <span class="required">*</span></label>
        <div class="upload-area upload-file" @click="$refs.fileInput.click()">
          <div v-if="uploadingFile" class="upload-placeholder">
            <div class="upload-spinner"></div>
            <span>Mengupload file PDF...</span>
          </div>
          <div v-else-if="form.file_url" class="file-uploaded">
            <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="#22c55e"><path d="M22 11.08V12a10 10 0 1 1-5.93-9.14" stroke-width="2"/><polyline points="22 4 12 14.01 9 11.01" stroke-width="2"/></svg>
            <span>File PDF sudah diupload ✅</span>
          </div>
          <div v-else class="upload-placeholder">
            <svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="#94a3b8"><path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z" stroke-width="2"/><polyline points="14 2 14 8 20 8" stroke-width="2"/></svg>
            <span>Klik untuk upload file PDF</span>
          </div>
        </div>
        <input ref="fileInput" type="file" accept=".pdf" @change="uploadFile" style="display:none" />
      </div>

      <!-- Active toggle -->
      <div class="form-group form-row">
        <label class="form-label">Status Aktif</label>
        <label class="toggle-switch">
          <input type="checkbox" v-model="form.is_active" />
          <span class="toggle-slider"></span>
        </label>
      </div>

      <!-- Actions -->
      <div class="form-actions">
        <router-link to="/ebooks" class="btn btn-secondary">Batal</router-link>
        <button type="submit" class="btn btn-primary" :disabled="submitting">
          {{ submitting ? 'Menyimpan...' : (isEdit ? 'Update' : 'Simpan') }}
        </button>
      </div>
    </form>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ebookApi, ebookCategoryApi, uploadApi } from '../../services/api'
import { useToastStore } from '../../stores/toast'
import { QuillEditor } from '@vueup/vue-quill'
import '@vueup/vue-quill/dist/vue-quill.snow.css'

const route = useRoute()
const router = useRouter()
const toast = useToastStore()

const isEdit = computed(() => !!route.params.id)
const submitting = ref(false)
const uploadingCover = ref(false)
const uploadingFile = ref(false)

const toolbarOptions = [
  ['bold', 'italic', 'underline', 'strike'],
  ['blockquote'],
  [{ header: [1, 2, 3, false] }],
  [{ list: 'ordered' }, { list: 'bullet' }],
  [{ color: [] }, { background: [] }],
  ['link'],
  ['clean'],
]

const form = ref({
  title: '',
  author: '',
  description: '',
  price: 0,
  cover_url: '',
  file_url: '',
  is_active: true,
})

const categories = ref([])
const selectedCategoryIds = ref([])

onMounted(async () => {
  // Load ebook categories master
  try {
    const { data } = await ebookCategoryApi.list({ active_only: 'true' })
    categories.value = data.data || []
  } catch {}

  if (isEdit.value) {
    try {
      const { data } = await ebookApi.get(route.params.id)
      const eb = data.data
      form.value = {
        title: eb.title,
        author: eb.author || '',
        description: eb.description || '',
        price: eb.price,
        cover_url: eb.cover_url || '',
        file_url: eb.file_url || '',
        is_active: eb.is_active,
      }
      selectedCategoryIds.value = (eb.categories || []).map(c => c.id)
    } catch { toast.error('Gagal memuat ebook') }
  }
})

// Price formatting
function formatPriceInput(v) {
  if (!v && v !== 0) return ''
  return Number(v).toLocaleString('id-ID')
}

function onPriceInput(e) {
  const raw = e.target.value.replace(/\D/g, '')
  form.value.price = parseInt(raw) || 0
}

async function uploadCover(e) {
  const file = e.target.files[0]
  if (!file) return
  uploadingCover.value = true
  try {
    const { data } = await uploadApi.upload(file)
    form.value.cover_url = data.data.url
    toast.success('Cover berhasil diupload')
  } catch { toast.error('Gagal upload cover') }
  finally { uploadingCover.value = false }
}

async function uploadFile(e) {
  const file = e.target.files[0]
  if (!file) return
  uploadingFile.value = true
  try {
    const { data } = await uploadApi.upload(file)
    form.value.file_url = data.data.url
    toast.success('File PDF berhasil diupload')
  } catch { toast.error('Gagal upload file') }
  finally { uploadingFile.value = false }
}

async function submit() {
  if (!form.value.title.trim()) {
    toast.error('Judul ebook wajib diisi')
    return
  }
  if (!form.value.file_url) {
    toast.error('File PDF wajib diupload terlebih dahulu')
    return
  }
  submitting.value = true
  const payload = { ...form.value, category_ids: selectedCategoryIds.value }
  try {
    if (isEdit.value) {
      await ebookApi.update(route.params.id, payload)
      toast.success('Ebook berhasil diupdate')
    } else {
      await ebookApi.create(payload)
      toast.success('Ebook berhasil ditambahkan')
    }
    router.push('/ebooks')
  } catch (err) {
    toast.error(err.response?.data?.error || 'Gagal menyimpan ebook')
  } finally { submitting.value = false }
}
</script>

<style scoped>
.form-hero { background: linear-gradient(135deg, #0f0c29 0%, #302b63 50%, #24243e 100%); border-radius: 16px; padding: 32px 36px; margin-bottom: 24px; }
.form-hero-title { font-size: 1.6rem; font-weight: 800; color: #fff; }
.form-hero-sub { font-size: 0.85rem; color: rgba(255,255,255,0.5); margin-top: 6px; }

.ebook-form-card { background: white; border-radius: 16px; padding: 32px; border: 1px solid #e8ecf1; }
.form-group { margin-bottom: 24px; }
.form-label { display: block; font-size: 0.85rem; font-weight: 600; color: #334155; margin-bottom: 8px; }
.required { color: #ef4444; }
.form-input { width: 100%; padding: 12px 16px; border: 1.5px solid #e2e8f0; border-radius: 12px; font-size: 0.85rem; color: #1e293b; font-family: inherit; box-sizing: border-box; }
.form-input:focus { border-color: #818cf8; outline: none; box-shadow: 0 0 0 3px rgba(129,140,248,0.12); }
.form-hint { font-size: 0.75rem; color: #94a3b8; margin-top: 4px; display: block; }
.form-row { display: flex; align-items: center; justify-content: space-between; }

/* Price input with Rp prefix */
.price-input-wrap { display: flex; align-items: center; border: 1.5px solid #e2e8f0; border-radius: 12px; overflow: hidden; background: white; }
.price-input-wrap:focus-within { border-color: #818cf8; box-shadow: 0 0 0 3px rgba(129,140,248,0.12); }
.price-prefix { padding: 12px 14px; background: #f8fafc; color: #64748b; font-weight: 700; font-size: 0.85rem; border-right: 1.5px solid #e2e8f0; user-select: none; }
.price-input { border: none !important; border-radius: 0 !important; box-shadow: none !important; font-weight: 600; font-size: 1rem; letter-spacing: 0.02em; }

/* Rich text editor wrapper */
.editor-wrapper { border: 1.5px solid #e2e8f0; border-radius: 12px; overflow: hidden; }
.editor-wrapper:focus-within { border-color: #818cf8; box-shadow: 0 0 0 3px rgba(129,140,248,0.12); }

/* Upload */
.upload-area { border: 2px dashed #e2e8f0; border-radius: 12px; cursor: pointer; overflow: hidden; transition: border-color 0.2s; }
.upload-area:hover { border-color: #818cf8; }
.upload-preview { width: 100%; max-height: 200px; object-fit: cover; display: block; }
.upload-placeholder { display: flex; flex-direction: column; align-items: center; gap: 8px; padding: 32px; color: #94a3b8; font-size: 0.8rem; }
.upload-file .file-uploaded { display: flex; align-items: center; gap: 8px; padding: 20px; color: #22c55e; font-weight: 600; font-size: 0.85rem; }

/* Toggle */
.toggle-switch { position: relative; display: inline-block; width: 48px; height: 26px; }
.toggle-switch input { opacity: 0; width: 0; height: 0; }
.toggle-slider { position: absolute; cursor: pointer; inset: 0; background: #cbd5e1; border-radius: 26px; transition: 0.3s; }
.toggle-slider::before { content: ''; position: absolute; height: 20px; width: 20px; left: 3px; bottom: 3px; background: white; border-radius: 50%; transition: 0.3s; }
.toggle-switch input:checked + .toggle-slider { background: #6366f1; }
.toggle-switch input:checked + .toggle-slider::before { transform: translateX(22px); }

/* Actions */
.form-actions { display: flex; justify-content: flex-end; gap: 12px; padding-top: 16px; border-top: 1px solid #f1f5f9; }
.btn { padding: 11px 24px; border-radius: 12px; font-size: 0.85rem; font-weight: 700; border: none; cursor: pointer; text-decoration: none; display: inline-flex; align-items: center; gap: 6px; }
.btn-primary { background: linear-gradient(135deg, #6366f1, #8b5cf6); color: white; box-shadow: 0 4px 12px rgba(99,102,241,0.3); }
.btn-secondary { background: #f1f5f9; color: #475569; }
.btn:disabled { opacity: 0.5; cursor: not-allowed; }

.cat-checkboxes { display: flex; flex-wrap: wrap; gap: 10px; }
.cat-checkbox-item { display: flex; align-items: center; gap: 6px; padding: 8px 14px; background: #f8fafc; border: 1.5px solid #e2e8f0; border-radius: 10px; cursor: pointer; transition: 0.2s; }
.cat-checkbox-item:has(input:checked) { background: #eef2ff; border-color: #818cf8; }
.cat-checkbox-item input[type="checkbox"] { accent-color: #6366f1; width: 16px; height: 16px; cursor: pointer; }
.cat-checkbox-label { font-size: 0.82rem; font-weight: 600; color: #334155; user-select: none; }
.cat-hint { font-size: 0.82rem; color: #94a3b8; }
.cat-hint a { color: #6366f1; text-decoration: underline; }

.upload-spinner { width: 24px; height: 24px; border: 3px solid #e2e8f0; border-top-color: #6366f1; border-radius: 50%; animation: spin 0.8s linear infinite; }
@keyframes spin { to { transform: rotate(360deg); } }
</style>

<style>
/* Quill overrides (unscoped to reach inner elements) */
.ql-toolbar.ql-snow { border: none !important; border-bottom: 1px solid #e2e8f0 !important; background: #f8fafc; border-radius: 12px 12px 0 0; }
.ql-container.ql-snow { border: none !important; font-family: inherit; font-size: 0.88rem; }
.ql-editor { min-height: 180px; color: #1e293b; line-height: 1.7; padding: 16px; }
.ql-editor.ql-blank::before { color: #94a3b8; font-style: normal; }
</style>
