<template>
  <div class="animate-in">
    <!-- Hero header -->
    <div class="outlet-hero">
      <div class="outlet-hero-content">
        <div class="outlet-hero-text">
          <h1 class="outlet-hero-title">Outlet Management</h1>
          <p class="outlet-hero-sub">Kelola semua data outlet franchise & kemitraan Anda</p>
        </div>
        <router-link to="/outlets/create" class="outlet-hero-btn">
          <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4"/></svg>
          Tambah Outlet
        </router-link>
      </div>
      <!-- Mini stats -->
      <div class="outlet-stats-bar">
        <div class="outlet-stat-item">
          <div class="outlet-stat-dot dot-all"></div>
          <span class="outlet-stat-label">Total</span>
          <span class="outlet-stat-value">{{ meta.total || 0 }}</span>
        </div>
        <div class="outlet-stat-item">
          <div class="outlet-stat-dot dot-active"></div>
          <span class="outlet-stat-label">Aktif</span>
          <span class="outlet-stat-value">{{ activeCount }}</span>
        </div>
        <div class="outlet-stat-item">
          <div class="outlet-stat-dot dot-featured"></div>
          <span class="outlet-stat-label">Featured</span>
          <span class="outlet-stat-value">{{ featuredCount }}</span>
        </div>
      </div>
    </div>

    <!-- Filter & Search strip -->
    <div class="outlet-filter-strip">
      <div class="outlet-search-box">
        <svg class="outlet-search-icon" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor"><circle cx="11" cy="11" r="8" stroke-width="2"/><path d="m21 21-4.35-4.35" stroke-width="2" stroke-linecap="round"/></svg>
        <input v-model="filters.search" @input="debouncedLoad" type="text" placeholder="Cari outlet..." class="outlet-search-input" />
      </div>
      <div class="outlet-filter-pills">
        <select v-model="filters.category_id" @change="loadOutlets" class="outlet-filter-select">
          <option value="">Semua Kategori</option>
          <option v-for="cat in categories" :key="cat.id" :value="cat.id">
            {{ cat.name }}
          </option>
        </select>
        <select v-model="filters.active" @change="loadOutlets" class="outlet-filter-select">
          <option value="">Status</option>
          <option value="true">Aktif</option>
          <option value="false">Nonaktif</option>
        </select>
        <select v-model="filters.featured" @change="loadOutlets" class="outlet-filter-select">
          <option value="">Featured</option>
          <option value="true">Featured</option>
          <option value="false">Regular</option>
        </select>
      </div>
    </div>

    <!-- Loading skeleton -->
    <div v-if="loading" class="outlet-grid">
      <div v-for="n in 6" :key="n" class="outlet-card-skeleton">
        <div class="skeleton" style="height:140px;border-radius:12px 12px 0 0"></div>
        <div style="padding:16px">
          <div class="skeleton" style="height:20px;width:70%;margin-bottom:8px;border-radius:6px"></div>
          <div class="skeleton" style="height:14px;width:50%;margin-bottom:16px;border-radius:6px"></div>
          <div class="skeleton" style="height:32px;border-radius:8px"></div>
        </div>
      </div>
    </div>

    <!-- Empty state -->
    <div v-else-if="!outlets.length" class="outlet-empty">
      <div class="outlet-empty-icon">
        <svg width="48" height="48" viewBox="0 0 24 24" fill="none" stroke="#94a3b8" stroke-width="1.5"><path d="M3 9l9-7 9 7v11a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2z"/><polyline points="9 22 9 12 15 12 15 22"/></svg>
      </div>
      <h3>Belum ada outlet</h3>
      <p>Mulai tambahkan outlet pertama Anda</p>
      <router-link to="/outlets/create" class="btn btn-primary" style="margin-top:16px">
        <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4"/></svg>
        Tambah Outlet Baru
      </router-link>
    </div>

    <!-- Outlet cards grid -->
    <div v-else class="outlet-grid">
      <div v-for="o in outlets" :key="o.id" class="outlet-card" :class="{ 'outlet-card-featured': o.is_featured }">
        <!-- Card banner -->
        <div class="outlet-card-banner" :style="bannerStyle(o)">
          <div class="outlet-card-overlay"></div>
          <div class="outlet-card-badges">
            <span class="outlet-badge-cat">{{ o.category_name || o.category }}</span>
            <span v-if="o.is_featured" class="outlet-badge-star">★ Featured</span>
          </div>
          <div v-if="!o.is_active" class="outlet-inactive-ribbon">Nonaktif</div>
        </div>

        <!-- Card body -->
        <div class="outlet-card-body">
          <div class="outlet-card-top">
            <div class="outlet-card-avatar" :style="avatarGradient(o)">
              <img v-if="o.logo" :src="o.logo" :alt="o.name" />
              <span v-else>{{ o.name?.charAt(0)?.toUpperCase() }}</span>
            </div>
            <div class="outlet-card-info">
              <router-link :to="`/outlets/${o.id}`" class="outlet-card-name">{{ o.name }}</router-link>
              <div class="outlet-card-location" v-if="o.city">
                <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor"><path d="M21 10c0 7-9 13-9 13s-9-6-9-13a9 9 0 0 1 18 0z" stroke-width="2"/><circle cx="12" cy="10" r="3" stroke-width="2"/></svg>
                {{ o.city }}<span v-if="o.province">, {{ o.province }}</span>
              </div>
            </div>
          </div>

          <p class="outlet-card-desc">{{ o.short_description || 'Belum ada deskripsi singkat' }}</p>

          <!-- Investment highlight -->
          <div class="outlet-card-invest">
            <div class="outlet-invest-row">
              <div class="outlet-invest-col">
                <div class="outlet-invest-label">Investasi</div>
                <div class="outlet-invest-value">{{ formatCurrencyCompact(o.minimum_investment) }}</div>
              </div>
              <div class="outlet-invest-col" v-if="o.profit_sharing_percentage">
                <div class="outlet-invest-label">Profit Share</div>
                <div class="outlet-invest-value">{{ o.profit_sharing_percentage }}%</div>
              </div>
              <div class="outlet-invest-col" v-if="o.estimated_roi">
                <div class="outlet-invest-label">ROI</div>
                <div class="outlet-invest-value">{{ o.estimated_roi }}</div>
              </div>
            </div>
          </div>

          <!-- Footer meta -->
          <div class="outlet-card-meta" v-if="o.total_outlets">
            <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor"><path d="M3 9l9-7 9 7v11a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2z" stroke-width="2"/></svg>
            <span>{{ o.total_outlets }} outlet beroperasi</span>
          </div>
        </div>

        <!-- Card actions -->
        <div class="outlet-card-actions">
          <router-link :to="`/outlets/${o.id}`" class="outlet-action-btn outlet-action-edit" title="Edit">
            <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor"><path d="M11 4H4a2 2 0 0 0-2 2v14a2 2 0 0 0 2 2h14a2 2 0 0 0 2-2v-7" stroke-width="2"/><path d="M18.5 2.5a2.121 2.121 0 0 1 3 3L12 15l-4 1 1-4 9.5-9.5z" stroke-width="2"/></svg>
          </router-link>
          <button @click="toggleActive(o)" class="outlet-action-btn" :class="o.is_active ? 'outlet-action-on' : 'outlet-action-off'" :title="o.is_active ? 'Nonaktifkan' : 'Aktifkan'">
            <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor"><path d="M18.36 6.64a9 9 0 1 1-12.73 0" stroke-width="2" stroke-linecap="round"/><line x1="12" y1="2" x2="12" y2="12" stroke-width="2" stroke-linecap="round"/></svg>
          </button>
          <button @click="toggleFeatured(o)" class="outlet-action-btn" :class="o.is_featured ? 'outlet-action-star' : ''" title="Toggle Featured">
            <svg width="14" height="14" viewBox="0 0 24 24" :fill="o.is_featured ? 'currentColor' : 'none'" stroke="currentColor"><polygon points="12 2 15.09 8.26 22 9.27 17 14.14 18.18 21.02 12 17.77 5.82 21.02 7 14.14 2 9.27 8.91 8.26 12 2" stroke-width="2" stroke-linejoin="round"/></svg>
          </button>
          <button @click="confirmDelete(o)" class="outlet-action-btn outlet-action-delete" title="Hapus">
            <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor"><polyline points="3 6 5 6 21 6" stroke-width="2"/><path d="M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6m3 0V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2" stroke-width="2"/></svg>
          </button>
        </div>
      </div>
    </div>

    <!-- Pagination -->
    <div v-if="totalPages > 1" class="outlet-pagination">
      <button @click="goPage(filters.page - 1)" :disabled="filters.page <= 1" class="outlet-page-btn">
        <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor"><polyline points="15 18 9 12 15 6" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/></svg>
      </button>
      <span class="outlet-page-info">
        Halaman <strong>{{ filters.page }}</strong> dari <strong>{{ totalPages }}</strong>
        <span class="outlet-page-total">({{ meta.total }} outlet)</span>
      </span>
      <button @click="goPage(filters.page + 1)" :disabled="filters.page >= totalPages" class="outlet-page-btn">
        <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor"><polyline points="9 18 15 12 9 6" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/></svg>
      </button>
    </div>

    <!-- Delete modal -->
    <div v-if="deleteTarget" class="modal-overlay" @click.self="deleteTarget = null">
      <div class="modal-content" style="max-width:420px;">
        <div class="modal-body" style="padding:2rem;text-align:center;">
          <div style="width:56px;height:56px;border-radius:50%;background:#fee2e2;display:flex;align-items:center;justify-content:center;margin:0 auto 16px;">
            <svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="#ef4444"><polyline points="3 6 5 6 21 6" stroke-width="2"/><path d="M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6m3 0V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2" stroke-width="2"/></svg>
          </div>
          <h3 style="font-size:1.1rem;font-weight:700;margin-bottom:8px;">Hapus Outlet</h3>
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
import { ref, computed, onMounted } from 'vue'
import { outletApi, outletCategoryApi } from '../../services/api'
import { useToastStore } from '../../stores/toast'

const toast = useToastStore()
const outlets = ref([])
const loading = ref(false)
const meta = ref({ total: 0, page: 1, limit: 12 })
const deleteTarget = ref(null)
const deleting = ref(false)

const filters = ref({ search: '', category_id: '', active: '', featured: '', page: 1 })
const categories = ref([])

const activeCount = computed(() => outlets.value.filter(o => o.is_active).length)
const featuredCount = computed(() => outlets.value.filter(o => o.is_featured).length)
const totalPages = computed(() => Math.ceil((meta.value.total || 0) / (meta.value.limit || 12)))

let debounceTimer = null
function debouncedLoad() {
  clearTimeout(debounceTimer)
  debounceTimer = setTimeout(() => { filters.value.page = 1; loadOutlets() }, 300)
}

onMounted(() => {
  loadOutlets()
  loadCategories()
})

async function loadCategories() {
  try {
    const { data } = await outletCategoryApi.list({ active_only: 'true' })
    categories.value = data.data || []
  } catch (e) { console.error('Failed to load categories', e) }
}

async function loadOutlets() {
  loading.value = true
  try {
    const params = { page: filters.value.page, limit: 12 }
    if (filters.value.search) params.search = filters.value.search
    if (filters.value.category_id) params.category_id = filters.value.category_id
    if (filters.value.active) params.active = filters.value.active
    if (filters.value.featured) params.featured = filters.value.featured

    const { data } = await outletApi.list(params)
    outlets.value = data.data || []
    meta.value = data.meta || { total: 0, page: 1, limit: 12 }
  } catch (e) {
    toast.error('Gagal memuat data outlet')
  } finally {
    loading.value = false
  }
}

function goPage(p) {
  if (p < 1 || p > totalPages.value) return
  filters.value.page = p
  loadOutlets()
}

async function toggleActive(o) {
  try {
    const { data } = await outletApi.toggle(o.id)
    toast.success(data.message)
    o.is_active = data.data.is_active
  } catch { toast.error('Gagal mengubah status') }
}

async function toggleFeatured(o) {
  try {
    const { data } = await outletApi.toggleFeatured(o.id)
    toast.success(data.message)
    o.is_featured = data.data.is_featured
  } catch { toast.error('Gagal mengubah featured') }
}

function confirmDelete(o) { deleteTarget.value = o }
async function doDelete() {
  deleting.value = true
  try {
    const { data } = await outletApi.delete(deleteTarget.value.id)
    toast.success(data.message || 'Outlet berhasil dihapus')
    deleteTarget.value = null
    loadOutlets()
  } catch { toast.error('Gagal menghapus outlet') }
  finally { deleting.value = false }
}

function formatCurrency(v) {
  return new Intl.NumberFormat('id-ID', { style: 'currency', currency: 'IDR', minimumFractionDigits: 0 }).format(v || 0)
}

function formatCurrencyCompact(v) {
  const n = v || 0
  if (n >= 1_000_000_000) {
    const val = n / 1_000_000_000
    return 'Rp ' + (val % 1 === 0 ? val.toFixed(0) : val.toFixed(1).replace('.', ',')) + 'M'
  }
  if (n >= 1_000_000) {
    const val = n / 1_000_000
    return 'Rp ' + (val % 1 === 0 ? val.toFixed(0) : val.toFixed(1).replace('.', ',')) + 'jt'
  }
  if (n >= 1_000) {
    const val = n / 1_000
    return 'Rp ' + (val % 1 === 0 ? val.toFixed(0) : val.toFixed(1).replace('.', ',')) + 'rb'
  }
  return 'Rp ' + n.toLocaleString('id-ID')
}

function categoryLabel(o) {
  if (o.category_name) return `${o.category_icon || ''} ${o.category_name}`
  return o.category || ''
}

const gradients = [
  'linear-gradient(135deg, #667eea 0%, #764ba2 100%)',
  'linear-gradient(135deg, #f093fb 0%, #f5576c 100%)',
  'linear-gradient(135deg, #4facfe 0%, #00f2fe 100%)',
  'linear-gradient(135deg, #43e97b 0%, #38f9d7 100%)',
  'linear-gradient(135deg, #fa709a 0%, #fee140 100%)',
  'linear-gradient(135deg, #a18cd1 0%, #fbc2eb 100%)',
  'linear-gradient(135deg, #fd9644 0%, #fc5c65 100%)',
]
function hashStr(s) { let h = 0; for (let i = 0; i < (s||'').length; i++) h = ((h << 5) - h) + s.charCodeAt(i); return Math.abs(h) }
function bannerStyle(o) {
  if (o.banner) return { backgroundImage: `url(${o.banner})`, backgroundSize: 'cover', backgroundPosition: 'center' }
  return { background: gradients[hashStr(o.name) % gradients.length] }
}
function avatarGradient(o) {
  if (o.logo) return {}
  return { background: gradients[(hashStr(o.name) + 3) % gradients.length] }
}
</script>

<style scoped>
/* ═══ HERO HEADER ═══ */
.outlet-hero {
  background: linear-gradient(135deg, #0f0c29 0%, #302b63 50%, #24243e 100%);
  border-radius: 16px;
  padding: 32px 36px 24px;
  margin-bottom: 24px;
  box-shadow: 0 4px 24px rgba(15,12,41,0.2);
}
.outlet-hero-content {
  display: flex; align-items: flex-start; justify-content: space-between; gap: 16px;
  flex-wrap: wrap; margin-bottom: 24px;
}
.outlet-hero-title { font-size: 1.6rem; font-weight: 800; color: #fff; line-height: 1.2; letter-spacing: -0.02em; }
.outlet-hero-sub { font-size: 0.85rem; color: rgba(255,255,255,0.5); margin-top: 6px; }
.outlet-hero-btn {
  display: inline-flex; align-items: center; gap: 8px;
  padding: 11px 24px; font-size: 0.85rem; font-weight: 700;
  border-radius: 12px; border: none; cursor: pointer;
  background: linear-gradient(135deg, #f59e0b 0%, #ef4444 100%);
  color: white; text-decoration: none;
  box-shadow: 0 4px 20px rgba(245,158,11,0.3);
  white-space: nowrap;
}

.outlet-stats-bar {
  display: flex; gap: 28px; flex-wrap: wrap;
  padding-top: 18px;
  border-top: 1px solid rgba(255,255,255,0.08);
}
.outlet-stat-item { display: flex; align-items: center; gap: 8px; }
.outlet-stat-dot { width: 8px; height: 8px; border-radius: 50%; }
.dot-all { background: #818cf8; box-shadow: 0 0 8px rgba(129,140,248,0.5); }
.dot-active { background: #4ade80; box-shadow: 0 0 8px rgba(74,222,128,0.5); }
.dot-featured { background: #fbbf24; box-shadow: 0 0 8px rgba(251,191,36,0.5); }
.outlet-stat-label { font-size: 0.72rem; color: rgba(255,255,255,0.4); text-transform: uppercase; letter-spacing: 0.05em; }
.outlet-stat-value { font-size: 0.9rem; font-weight: 800; color: white; }

/* ═══ FILTER STRIP ═══ */
.outlet-filter-strip {
  display: flex; gap: 12px; align-items: center; flex-wrap: wrap;
  margin-bottom: 24px;
}
.outlet-search-box {
  display: flex; align-items: center; gap: 10px;
  background: white; border: 1.5px solid #e2e8f0; border-radius: 12px;
  padding: 0 16px; height: 46px; flex: 1; min-width: 200px;
  box-shadow: 0 1px 3px rgba(0,0,0,0.04);
}
.outlet-search-box:focus-within { border-color: #818cf8; box-shadow: 0 0 0 3px rgba(129,140,248,0.12); }
.outlet-search-icon { color: #94a3b8; flex-shrink: 0; }
.outlet-search-input {
  flex: 1; border: none; outline: none; font-size: 0.85rem;
  color: #1e293b; background: none; font-family: inherit;
}
.outlet-search-input::placeholder { color: #94a3b8; }

.outlet-filter-pills { display: flex; gap: 8px; flex-wrap: wrap; }
.outlet-filter-select {
  height: 46px; padding: 0 34px 0 14px; border-radius: 12px;
  border: 1.5px solid #e2e8f0; font-size: 0.8rem; font-weight: 500;
  color: #475569; background: white; cursor: pointer;
  appearance: none;
  background-image: url("data:image/svg+xml,%3csvg xmlns='http://www.w3.org/2000/svg' fill='none' viewBox='0 0 20 20'%3e%3cpath stroke='%236b7280' stroke-linecap='round' stroke-linejoin='round' stroke-width='1.5' d='M6 8l4 4 4-4'/%3e%3c/svg%3e");
  background-position: right 10px center; background-repeat: no-repeat; background-size: 1.4em;
  box-shadow: 0 1px 3px rgba(0,0,0,0.04);
}
.outlet-filter-select:focus { border-color: #818cf8; outline: none; }

/* ═══ CARD GRID ═══ */
.outlet-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(340px, 1fr));
  gap: 24px;
  margin-bottom: 28px;
}

.outlet-card {
  background: #ffffff;
  border-radius: 16px;
  overflow: hidden;
  border: 1px solid #e8ecf1;
  position: relative;
  display: flex;
  flex-direction: column;
  box-shadow: 0 1px 4px rgba(0,0,0,0.04);
}
.outlet-card-featured {
  border-color: #f59e0b;
  box-shadow: 0 0 0 1px rgba(245,158,11,0.3), 0 4px 16px rgba(245,158,11,0.1);
}

/* ═══ BANNER ═══ */
.outlet-card-banner {
  height: 140px; position: relative; overflow: hidden; flex-shrink: 0;
}
.outlet-card-overlay {
  position: absolute; inset: 0;
  background: linear-gradient(180deg, rgba(0,0,0,0.05) 0%, rgba(0,0,0,0.4) 100%);
}
.outlet-card-badges {
  position: absolute; top: 12px; left: 12px; display: flex; gap: 6px; z-index: 2;
}
.outlet-badge-cat {
  padding: 5px 12px; font-size: 0.65rem; font-weight: 700; letter-spacing: 0.04em;
  border-radius: 8px;
  background: rgba(255,255,255,0.85);
  color: #1e293b;
  backdrop-filter: blur(12px);
  -webkit-backdrop-filter: blur(12px);
  text-transform: uppercase;
}
.outlet-badge-star {
  padding: 5px 12px; font-size: 0.65rem; font-weight: 700; letter-spacing: 0.04em;
  border-radius: 8px;
  background: linear-gradient(135deg, #fbbf24, #f59e0b);
  color: #78350f;
  white-space: nowrap;
  text-transform: uppercase;
}
.outlet-inactive-ribbon {
  position: absolute; top: 12px; right: -22px;
  background: linear-gradient(135deg, #ef4444, #dc2626);
  color: white; font-size: 0.6rem; font-weight: 700;
  padding: 4px 30px; transform: rotate(45deg); letter-spacing: 0.06em;
  text-transform: uppercase;
  box-shadow: 0 2px 8px rgba(239,68,68,0.3);
}

/* ═══ CARD BODY ═══ */
.outlet-card-body {
  padding: 18px 20px 16px;
  flex: 1;
  display: flex;
  flex-direction: column;
}

.outlet-card-top { display: flex; align-items: center; gap: 14px; margin-bottom: 14px; }
.outlet-card-avatar {
  width: 48px; height: 48px; border-radius: 14px; overflow: hidden;
  display: flex; align-items: center; justify-content: center;
  font-weight: 800; font-size: 1.2rem; color: white; flex-shrink: 0;
  border: 3px solid white;
  box-shadow: 0 3px 12px rgba(0,0,0,0.12);
  margin-top: -36px; position: relative; z-index: 3;
}
.outlet-card-avatar img { width: 100%; height: 100%; object-fit: cover; }

.outlet-card-info { flex: 1; min-width: 0; }
.outlet-card-name {
  font-size: 1rem; font-weight: 700; color: #0f172a;
  text-decoration: none; display: block;
  white-space: nowrap; overflow: hidden; text-overflow: ellipsis;
  letter-spacing: -0.01em;
}
.outlet-card-name:hover { color: #6366f1; }
.outlet-card-location {
  display: flex; align-items: center; gap: 4px;
  font-size: 0.72rem; color: #94a3b8; margin-top: 3px;
}

.outlet-card-desc {
  font-size: 0.8rem; color: #64748b; line-height: 1.6;
  display: -webkit-box; -webkit-line-clamp: 2; -webkit-box-orient: vertical; overflow: hidden;
  margin-bottom: 16px; min-height: 38px;
}

/* ═══ INVESTMENT STATS (PREMIUM) ═══ */
.outlet-card-invest {
  background: linear-gradient(135deg, #1e1b4b 0%, #312e81 100%);
  border-radius: 12px;
  padding: 14px 16px;
  margin-bottom: 14px;
}
.outlet-invest-row {
  display: flex;
  gap: 8px;
}
.outlet-invest-col {
  flex: 1;
  min-width: 0;
}
.outlet-invest-col:not(:last-child) {
  border-right: 1px solid rgba(255,255,255,0.1);
  padding-right: 8px;
}
.outlet-invest-label {
  font-size: 0.6rem; font-weight: 600; color: rgba(255,255,255,0.45);
  text-transform: uppercase; letter-spacing: 0.06em;
  margin-bottom: 3px;
}
.outlet-invest-value {
  font-size: 0.85rem; font-weight: 800; color: #fbbf24;
  white-space: nowrap; overflow: hidden; text-overflow: ellipsis;
}

/* ═══ FOOTER META ═══ */
.outlet-card-meta {
  display: flex; align-items: center; gap: 6px;
  padding-top: 12px; border-top: 1px solid #f1f5f9;
  margin-top: auto;
  font-size: 0.72rem; font-weight: 600; color: #94a3b8;
}

/* ═══ CARD ACTIONS ═══ */
.outlet-card-actions {
  display: flex;
  border-top: 1px solid #f1f5f9;
  flex-shrink: 0;
  background: #fafbfc;
}
.outlet-action-btn {
  flex: 1; display: flex; align-items: center; justify-content: center;
  padding: 11px 0; border: none; background: none; cursor: pointer;
  color: #b0b8c4; font-size: 0.8rem;
  text-decoration: none;
}
.outlet-action-btn:not(:last-child) { border-right: 1px solid #f1f5f9; }
.outlet-action-edit { color: #94a3b8; }
.outlet-action-on { color: #22c55e; }
.outlet-action-off { color: #d1d5db; }
.outlet-action-star { color: #f59e0b; }
.outlet-action-delete { color: #cbd5e1; }

/* ═══ SKELETON ═══ */
.outlet-card-skeleton {
  background: white; border-radius: 16px; border: 1px solid #e8ecf1; overflow: hidden;
}

/* ═══ EMPTY STATE ═══ */
.outlet-empty {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  text-align: center;
  padding: 80px 32px;
  background: white;
  border-radius: 16px;
  border: 2px dashed #e2e8f0;
  min-height: 360px;
}
.outlet-empty-icon {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 80px;
  height: 80px;
  border-radius: 50%;
  background: linear-gradient(135deg, #f1f5f9 0%, #e2e8f0 100%);
  margin-bottom: 24px;
}
.outlet-empty h3 { font-size: 1.15rem; font-weight: 700; color: #0f172a; margin: 0; }
.outlet-empty p { font-size: 0.85rem; color: #94a3b8; margin-top: 6px; margin-bottom: 0; }

/* ═══ PAGINATION ═══ */
.outlet-pagination {
  display: flex; align-items: center; justify-content: center; gap: 16px;
  padding: 12px 0;
}
.outlet-page-btn {
  width: 42px; height: 42px; border-radius: 12px;
  border: 1.5px solid #e2e8f0; background: white; cursor: pointer;
  display: flex; align-items: center; justify-content: center;
  color: #64748b;
  box-shadow: 0 1px 3px rgba(0,0,0,0.04);
}
.outlet-page-btn:disabled { opacity: 0.3; cursor: not-allowed; }
.outlet-page-info { font-size: 0.85rem; color: #64748b; }
.outlet-page-info strong { color: #0f172a; font-weight: 700; }
.outlet-page-total { color: #94a3b8; font-size: 0.75rem; }

/* ═══ RESPONSIVE ═══ */
@media (max-width: 768px) {
  .outlet-hero { padding: 24px 20px 18px; margin-bottom: 18px; }
  .outlet-hero-title { font-size: 1.3rem; }
  .outlet-grid { grid-template-columns: 1fr; gap: 18px; }
  .outlet-filter-strip { flex-direction: column; }
  .outlet-search-box { width: 100%; }
  .outlet-filter-pills { width: 100%; }
  .outlet-filter-select { flex: 1; min-width: 0; }
  .outlet-invest-value { font-size: 0.78rem; }
}
</style>
