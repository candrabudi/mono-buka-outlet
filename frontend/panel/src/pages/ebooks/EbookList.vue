<template>
  <div class="animate-in">
    <!-- Hero -->
    <div class="ebook-hero">
      <div class="ebook-hero-content">
        <div>
          <h1 class="ebook-hero-title"><i class="ri-book-open-line"></i> Ebook Management</h1>
          <p class="ebook-hero-sub">Kelola katalog ebook untuk mitra</p>
        </div>
        <router-link to="/ebooks/create" class="ebook-hero-btn">
          <svg width="16" height="16" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4"/></svg>
          Tambah Ebook
        </router-link>
      </div>
      <div class="ebook-stats-bar">
        <div class="ebook-stat"><span class="stat-dot dot-total"></span><span class="stat-label">Total</span><span class="stat-val">{{ total }}</span></div>
        <div class="ebook-stat"><span class="stat-dot dot-active"></span><span class="stat-label">Aktif</span><span class="stat-val">{{ activeCount }}</span></div>
        <div class="ebook-stat"><span class="stat-dot dot-sold"></span><span class="stat-label">Terjual</span><span class="stat-val">{{ totalSold }}</span></div>
      </div>
    </div>

    <!-- Search -->
    <div class="ebook-toolbar">
      <div class="ebook-search">
        <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="#94a3b8"><circle cx="11" cy="11" r="8" stroke-width="2"/><path d="m21 21-4.35-4.35" stroke-width="2" stroke-linecap="round"/></svg>
        <input v-model="search" @input="debouncedLoad" type="text" placeholder="Cari judul, penulis, atau kategori..." />
      </div>
      <select v-if="categories.length" v-model="selectedCategory" @change="load" class="ebook-filter-select">
        <option value="">Semua Kategori</option>
        <option v-for="cat in categories" :key="cat" :value="cat">{{ cat }}</option>
      </select>
    </div>

    <!-- Loading -->
    <div v-if="loading" class="ebook-grid">
      <div v-for="n in 6" :key="n" class="ebook-card skeleton-card">
        <div class="skeleton" style="height:180px;border-radius:12px 12px 0 0"></div>
        <div style="padding:16px"><div class="skeleton" style="height:20px;width:70%;border-radius:6px;margin-bottom:8px"></div><div class="skeleton" style="height:14px;width:40%;border-radius:6px"></div></div>
      </div>
    </div>

    <!-- Empty -->
    <div v-else-if="!ebooks.length" class="ebook-empty">
      <div class="ebook-empty-icon"><i class="ri-book-line" style="font-size:3rem;color:#94a3b8"></i></div>
      <h3>Belum ada ebook</h3>
      <p>Tambahkan ebook pertama untuk mitra Anda</p>
      <router-link to="/ebooks/create" class="btn btn-primary" style="margin-top:16px">Tambah Ebook</router-link>
    </div>

    <!-- Grid -->
    <div v-else class="ebook-grid">
      <div v-for="eb in ebooks" :key="eb.id" class="ebook-card" :class="{ 'card-inactive': !eb.is_active }">
        <div class="ebook-cover" :style="coverStyle(eb)">
          <div v-if="!eb.is_active" class="inactive-badge">Nonaktif</div>
          <div class="ebook-price-tag">{{ eb.price === 0 ? 'GRATIS' : formatRp(eb.price) }}</div>
        </div>
        <div class="ebook-body">
          <h3 class="ebook-title">{{ eb.title }}</h3>
          <p class="ebook-author" v-if="eb.author">{{ eb.author }}</p>
          <div v-if="eb.categories && eb.categories.length" class="ebook-category-badges">
            <span v-for="cat in eb.categories" :key="cat.id" class="ebook-category-badge">{{ cat.name }}</span>
          </div>
          <p class="ebook-desc">{{ stripHtml(eb.description, 80) }}</p>
          <div class="ebook-meta">
            <span><i class="ri-shopping-cart-line"></i> {{ eb.total_sold }} terjual</span>
          </div>
        </div>
        <div class="ebook-actions">
          <router-link :to="`/ebooks/${eb.id}`" class="action-btn action-edit" title="Edit">
            <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor"><path d="M11 4H4a2 2 0 0 0-2 2v14a2 2 0 0 0 2 2h14a2 2 0 0 0 2-2v-7" stroke-width="2"/><path d="M18.5 2.5a2.121 2.121 0 0 1 3 3L12 15l-4 1 1-4 9.5-9.5z" stroke-width="2"/></svg>
          </router-link>
          <button @click="toggleActive(eb)" class="action-btn" :class="eb.is_active ? 'action-on' : 'action-off'" :title="eb.is_active ? 'Nonaktifkan' : 'Aktifkan'">
            <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor"><path d="M18.36 6.64a9 9 0 1 1-12.73 0" stroke-width="2" stroke-linecap="round"/><line x1="12" y1="2" x2="12" y2="12" stroke-width="2" stroke-linecap="round"/></svg>
          </button>
          <button @click="confirmDelete(eb)" class="action-btn action-delete" title="Hapus">
            <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor"><polyline points="3 6 5 6 21 6" stroke-width="2"/><path d="M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6m3 0V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2" stroke-width="2"/></svg>
          </button>
        </div>
      </div>
    </div>

    <!-- Delete modal -->
    <div v-if="deleteTarget" class="modal-overlay" @click.self="deleteTarget = null">
      <div class="modal-content" style="max-width:420px;">
        <div class="modal-body" style="padding:2rem;text-align:center;">
          <div style="margin-bottom:12px;"><i class="ri-delete-bin-line" style="font-size:2.5rem;color:#ef4444"></i></div>
          <h3 style="font-size:1.1rem;font-weight:700;margin-bottom:8px;">Hapus Ebook</h3>
          <p style="color:#64748b;font-size:0.875rem;">Yakin ingin menghapus <strong>{{ deleteTarget.title }}</strong>?</p>
        </div>
        <div class="modal-footer" style="justify-content:center;gap:12px;padding-bottom:1.5rem;">
          <button @click="deleteTarget = null" class="btn btn-secondary">Batal</button>
          <button @click="doDelete" class="btn btn-danger" :disabled="deleting">{{ deleting ? 'Menghapus...' : 'Ya, Hapus' }}</button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { ebookApi } from '../../services/api'
import { useToastStore } from '../../stores/toast'

const toast = useToastStore()
const ebooks = ref([])
const loading = ref(false)
const total = ref(0)
const search = ref('')
const selectedCategory = ref('')
const categories = ref([])
const deleteTarget = ref(null)
const deleting = ref(false)

const activeCount = computed(() => ebooks.value.filter(e => e.is_active).length)
const totalSold = computed(() => ebooks.value.reduce((sum, e) => sum + (e.total_sold || 0), 0))

let timer = null
function debouncedLoad() { clearTimeout(timer); timer = setTimeout(load, 300) }

onMounted(load)

async function load() {
  loading.value = true
  try {
    const { data } = await ebookApi.list({ search: search.value, limit: 50 })
    let allEbooks = data.data || []
    total.value = data.total || 0
    // Build category filter from all ebooks
    const cats = new Set()
    allEbooks.forEach(e => { (e.categories || []).forEach(c => cats.add(c.name)) })
    categories.value = [...cats].sort()
    // Client-side filter by selected category
    if (selectedCategory.value) {
      allEbooks = allEbooks.filter(e =>
        (e.categories || []).some(c => c.name === selectedCategory.value)
      )
    }
    ebooks.value = allEbooks
  } catch { toast.error('Gagal memuat ebook') }
  finally { loading.value = false }
}

async function toggleActive(eb) {
  try {
    const { data } = await ebookApi.toggle(eb.id)
    eb.is_active = data.data.is_active
    toast.success(eb.is_active ? 'Ebook diaktifkan' : 'Ebook dinonaktifkan')
  } catch { toast.error('Gagal mengubah status') }
}

function confirmDelete(eb) { deleteTarget.value = eb }
async function doDelete() {
  deleting.value = true
  try {
    await ebookApi.delete(deleteTarget.value.id)
    toast.success('Ebook berhasil dihapus')
    deleteTarget.value = null
    load()
  } catch { toast.error('Gagal menghapus ebook') }
  finally { deleting.value = false }
}

function formatRp(v) {
  return 'Rp ' + Number(v || 0).toLocaleString('id-ID')
}

function stripHtml(html, maxLen) {
  const text = (html || '').replace(/<[^>]*>/g, '').replace(/&nbsp;/g, ' ').trim()
  return text.length > maxLen ? text.slice(0, maxLen) + '...' : (text || 'Belum ada deskripsi')
}

function coverStyle(eb) {
  if (eb.cover_url) return { backgroundImage: `url(${eb.cover_url})`, backgroundSize: 'cover', backgroundPosition: 'center' }
  const gradients = ['linear-gradient(135deg,#667eea,#764ba2)','linear-gradient(135deg,#f093fb,#f5576c)','linear-gradient(135deg,#4facfe,#00f2fe)','linear-gradient(135deg,#43e97b,#38f9d7)','linear-gradient(135deg,#fa709a,#fee140)']
  let h = 0; for (let i = 0; i < (eb.title||'').length; i++) h = ((h << 5) - h) + eb.title.charCodeAt(i)
  return { background: gradients[Math.abs(h) % gradients.length] }
}
</script>

<style scoped>
.ebook-hero { background: linear-gradient(135deg, #0f0c29 0%, #302b63 50%, #24243e 100%); border-radius: 16px; padding: 32px 36px 24px; margin-bottom: 24px; box-shadow: 0 4px 24px rgba(15,12,41,0.2); }
.ebook-hero-content { display: flex; align-items: flex-start; justify-content: space-between; gap: 16px; flex-wrap: wrap; margin-bottom: 24px; }
.ebook-hero-title { font-size: 1.6rem; font-weight: 800; color: #fff; }
.ebook-hero-sub { font-size: 0.85rem; color: rgba(255,255,255,0.5); margin-top: 6px; }
.ebook-hero-btn { display: inline-flex; align-items: center; gap: 8px; padding: 11px 24px; font-size: 0.85rem; font-weight: 700; border-radius: 12px; background: linear-gradient(135deg, #f59e0b, #ef4444); color: white; text-decoration: none; box-shadow: 0 4px 20px rgba(245,158,11,0.3); }
.ebook-stats-bar { display: flex; gap: 28px; padding-top: 18px; border-top: 1px solid rgba(255,255,255,0.08); }
.ebook-stat { display: flex; align-items: center; gap: 8px; }
.stat-dot { width: 8px; height: 8px; border-radius: 50%; }
.dot-total { background: #818cf8; box-shadow: 0 0 8px rgba(129,140,248,0.5); }
.dot-active { background: #4ade80; box-shadow: 0 0 8px rgba(74,222,128,0.5); }
.dot-sold { background: #fbbf24; box-shadow: 0 0 8px rgba(251,191,36,0.5); }
.stat-label { font-size: 0.72rem; color: rgba(255,255,255,0.4); text-transform: uppercase; letter-spacing: 0.05em; }
.stat-val { font-size: 0.9rem; font-weight: 800; color: white; }

.ebook-toolbar { display: flex; gap: 12px; margin-bottom: 24px; }
.ebook-search { display: flex; align-items: center; gap: 10px; background: white; border: 1.5px solid #e2e8f0; border-radius: 12px; padding: 0 16px; height: 46px; flex: 1; box-shadow: 0 1px 3px rgba(0,0,0,0.04); }
.ebook-search:focus-within { border-color: #818cf8; box-shadow: 0 0 0 3px rgba(129,140,248,0.12); }
.ebook-search input { flex: 1; border: none; outline: none; font-size: 0.85rem; color: #1e293b; background: none; font-family: inherit; }

.ebook-grid { display: grid; grid-template-columns: repeat(auto-fill, minmax(280px, 1fr)); gap: 24px; margin-bottom: 28px; }
.ebook-card { background: #fff; border-radius: 16px; overflow: hidden; border: 1px solid #e8ecf1; display: flex; flex-direction: column; box-shadow: 0 1px 4px rgba(0,0,0,0.04); transition: transform 0.2s, box-shadow 0.2s; }
.ebook-card:hover { transform: translateY(-2px); box-shadow: 0 8px 24px rgba(0,0,0,0.08); }
.card-inactive { opacity: 0.6; }

.ebook-cover { height: 180px; position: relative; display: flex; align-items: center; justify-content: center; color: white; font-size: 3rem; }
.inactive-badge { position: absolute; top: 12px; right: 12px; background: rgba(239,68,68,0.9); color: white; font-size: 0.65rem; font-weight: 700; padding: 4px 12px; border-radius: 8px; text-transform: uppercase; }
.ebook-price-tag { position: absolute; bottom: 12px; left: 12px; background: rgba(0,0,0,0.7); color: #fbbf24; font-size: 0.75rem; font-weight: 800; padding: 5px 14px; border-radius: 8px; backdrop-filter: blur(8px); }

.ebook-body { padding: 18px 20px 16px; flex: 1; }
.ebook-title { font-size: 1rem; font-weight: 700; color: #0f172a; margin: 0 0 4px; white-space: nowrap; overflow: hidden; text-overflow: ellipsis; }
.ebook-author { font-size: 0.8rem; color: #64748b; margin: 0 0 6px; }
.ebook-category-badges { display: flex; flex-wrap: wrap; gap: 4px; margin-bottom: 8px; }
.ebook-category-badge { display: inline-block; padding: 3px 10px; font-size: 0.68rem; font-weight: 700; color: #6366f1; background: #eef2ff; border-radius: 6px; }
.ebook-desc { font-size: 0.78rem; color: #94a3b8; line-height: 1.5; margin: 0 0 10px; display: -webkit-box; -webkit-line-clamp: 2; -webkit-box-orient: vertical; overflow: hidden; }
.ebook-meta { font-size: 0.75rem; color: #94a3b8; display: flex; gap: 12px; }
.ebook-filter-select { height: 46px; padding: 0 34px 0 14px; border-radius: 12px; border: 1.5px solid #e2e8f0; font-size: 0.82rem; font-weight: 500; color: #475569; background: white; cursor: pointer; appearance: none; background-image: url("data:image/svg+xml,%3csvg xmlns='http://www.w3.org/2000/svg' fill='none' viewBox='0 0 20 20'%3e%3cpath stroke='%236b7280' stroke-linecap='round' stroke-linejoin='round' stroke-width='1.5' d='M6 8l4 4 4-4'/%3e%3c/svg%3e"); background-position: right 10px center; background-repeat: no-repeat; background-size: 1.2em; }

.ebook-actions { display: flex; border-top: 1px solid #f1f5f9; background: #fafbfc; }
.action-btn { flex: 1; display: flex; align-items: center; justify-content: center; padding: 11px 0; border: none; background: none; cursor: pointer; color: #b0b8c4; text-decoration: none; }
.action-btn:not(:last-child) { border-right: 1px solid #f1f5f9; }
.action-edit { color: #94a3b8; }
.action-on { color: #22c55e; }
.action-off { color: #d1d5db; }
.action-delete { color: #cbd5e1; }

.ebook-empty { display: flex; flex-direction: column; align-items: center; padding: 80px 32px; background: white; border-radius: 16px; border: 2px dashed #e2e8f0; text-align: center; }
.ebook-empty-icon { font-size: 3rem; margin-bottom: 16px; }
.ebook-empty h3 { font-size: 1.15rem; font-weight: 700; color: #0f172a; margin: 0; }
.ebook-empty p { font-size: 0.85rem; color: #94a3b8; margin-top: 6px; }

.skeleton-card { overflow: hidden; }
.skeleton { background: linear-gradient(90deg, #f1f5f9 25%, #e2e8f0 50%, #f1f5f9 75%); background-size: 200% 100%; animation: shimmer 1.5s infinite; }
@keyframes shimmer { 0% { background-position: 200% 0; } 100% { background-position: -200% 0; } }

@media (max-width: 768px) {
  .ebook-hero { padding: 24px 20px 18px; }
  .ebook-hero-title { font-size: 1.3rem; }
  .ebook-grid { grid-template-columns: 1fr; }
}
</style>
