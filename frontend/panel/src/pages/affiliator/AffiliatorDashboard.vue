<template>
  <div class="aff-page">
    <!-- Hero Header -->
    <div class="aff-hero">
      <div class="aff-hero-content">
        <div>
          <h1 class="aff-hero-title">Affiliator Dashboard</h1>
          <p class="aff-hero-sub">Pantau referral & partnership Anda</p>
        </div>
        <div class="aff-referral-box" v-if="referralCode">
          <div class="aff-referral-label">Kode Referral Anda</div>
          <div class="aff-referral-code-wrap">
            <span class="aff-referral-code">{{ referralCode }}</span>
            <button @click="copyCode" class="aff-copy-btn" :class="{ copied: justCopied }">
              <Check v-if="justCopied" :size="16" />
              <Copy v-else :size="16" />
              {{ justCopied ? 'Disalin!' : 'Salin' }}
            </button>
          </div>
        </div>
      </div>

      <!-- Stats Cards -->
      <div class="aff-stats" v-if="stats">
        <div class="aff-stat-card">
          <div class="aff-stat-icon icon-referral">
            <Users :size="22" />
          </div>
          <div class="aff-stat-info">
            <span class="aff-stat-value">{{ stats.total_referrals }}</span>
            <span class="aff-stat-label">Total Referral</span>
          </div>
        </div>
        <div class="aff-stat-card">
          <div class="aff-stat-icon icon-active">
            <CheckCircle :size="22" />
          </div>
          <div class="aff-stat-info">
            <span class="aff-stat-value">{{ stats.active_partnerships }}</span>
            <span class="aff-stat-label">Partnership Aktif</span>
          </div>
        </div>
        <div class="aff-stat-card">
          <div class="aff-stat-icon icon-pending">
            <Clock :size="22" />
          </div>
          <div class="aff-stat-info">
            <span class="aff-stat-value">{{ stats.pending_partnerships }}</span>
            <span class="aff-stat-label">Menunggu</span>
          </div>
        </div>
        <div class="aff-stat-card">
          <div class="aff-stat-icon icon-rate">
            <BarChart3 :size="22" />
          </div>
          <div class="aff-stat-info">
            <span class="aff-stat-value">{{ stats.conversion_rate?.toFixed(1) }}%</span>
            <span class="aff-stat-label">Konversi</span>
          </div>
        </div>
      </div>
    </div>

    <!-- Partnership Table -->
    <div class="aff-section">
      <div class="aff-section-header">
        <h2>Partnership Referral Saya</h2>
        <div class="aff-search-wrap">
          <Search class="aff-search-icon" :size="16" color="#94a3b8" />
          <input v-model="searchQuery" type="text" class="aff-search" placeholder="Cari mitra..." />
        </div>
      </div>

      <!-- Loading -->
      <div v-if="loading" class="aff-loading">
        <div v-for="n in 5" :key="n" class="aff-row-skel shimmer"></div>
      </div>

      <!-- Table -->
      <div v-else-if="filteredPartnerships.length" class="aff-table-wrap">
        <table class="aff-table">
          <thead>
            <tr>
              <th>Mitra</th>
              <th>Outlet</th>
              <th>Paket</th>
              <th>Status</th>
              <th>Progress</th>
              <th>Tanggal</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="p in filteredPartnerships" :key="p.id">
              <td>
                <div class="aff-mitra-cell">
                  <div class="aff-mitra-avatar" :style="avatarBg(p.mitra?.name || '')">{{ getInitial(p.mitra?.name) }}</div>
                  <div>
                    <div class="aff-mitra-name">{{ p.mitra?.name || '-' }}</div>
                    <div class="aff-mitra-email">{{ p.mitra?.email || '-' }}</div>
                  </div>
                </div>
              </td>
              <td>{{ p.outlet?.name || '-' }}</td>
              <td>{{ p.package?.name || '-' }}</td>
              <td>
                <span class="aff-status-badge" :class="statusClass(p.status)">{{ statusLabel(p.status) }}</span>
              </td>
              <td>
                <div class="aff-progress-bar">
                  <div class="aff-progress-fill" :style="{ width: p.progress_percentage + '%' }"></div>
                </div>
                <span class="aff-progress-text">{{ p.progress_percentage }}%</span>
              </td>
              <td class="aff-date-col">{{ formatDate(p.created_at) }}</td>
            </tr>
          </tbody>
        </table>
      </div>

      <!-- Empty -->
      <div v-else class="aff-empty">
        <div class="aff-empty-circle">
          <UserPlus :size="36" color="#94a3b8" :stroke-width="1.5" />
        </div>
        <h3>Belum ada partnership referral</h3>
        <p>Bagikan kode referral Anda untuk mendapatkan partnership.</p>
      </div>

      <!-- Pagination -->
      <div v-if="totalPages > 1" class="aff-pagination">
        <button @click="changePage(page - 1)" :disabled="page <= 1" class="aff-page-btn">&laquo; Prev</button>
        <span class="aff-page-info">Halaman {{ page }} dari {{ totalPages }}</span>
        <button @click="changePage(page + 1)" :disabled="page >= totalPages" class="aff-page-btn">Next &raquo;</button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { Copy, Check, Users, CheckCircle, Clock, BarChart3, Search, UserPlus } from 'lucide-vue-next'
import { affiliatorApi } from '../../services/api'
import { useToastStore } from '../../stores/toast'

const toast = useToastStore()
const stats = ref(null)
const partnerships = ref([])
const total = ref(0)
const page = ref(1)
const limit = ref(10)
const loading = ref(false)
const searchQuery = ref('')
const referralCode = ref('')
const justCopied = ref(false)

const totalPages = computed(() => Math.ceil(total.value / limit.value))

const filteredPartnerships = computed(() => {
  const q = searchQuery.value.toLowerCase().trim()
  if (!q) return partnerships.value
  return partnerships.value.filter(p =>
    (p.mitra?.name || '').toLowerCase().includes(q) ||
    (p.mitra?.email || '').toLowerCase().includes(q) ||
    (p.outlet?.name || '').toLowerCase().includes(q)
  )
})

onMounted(async () => {
  await Promise.all([loadDashboard(), loadPartnerships()])
})

async function loadDashboard() {
  try {
    const { data } = await affiliatorApi.dashboard()
    stats.value = data.data
    referralCode.value = data.data.referral_code || ''
  } catch {
    toast.error('Gagal memuat dashboard')
  }
}

async function loadPartnerships() {
  loading.value = true
  try {
    const { data } = await affiliatorApi.partnerships({ page: page.value, limit: limit.value })
    partnerships.value = data.data || []
    total.value = data.meta?.total || 0
  } catch {
    toast.error('Gagal memuat partnership')
  } finally {
    loading.value = false
  }
}

function changePage(p) {
  if (p < 1 || p > totalPages.value) return
  page.value = p
  loadPartnerships()
}

async function copyCode() {
  try {
    await navigator.clipboard.writeText(referralCode.value)
    justCopied.value = true
    setTimeout(() => { justCopied.value = false }, 2000)
  } catch {
    toast.error('Gagal menyalin kode referral')
  }
}

const avatarGradients = [
  'linear-gradient(135deg, #667eea, #764ba2)',
  'linear-gradient(135deg, #f093fb, #f5576c)',
  'linear-gradient(135deg, #4facfe, #00f2fe)',
  'linear-gradient(135deg, #43e97b, #38f9d7)',
  'linear-gradient(135deg, #fa709a, #fee140)',
  'linear-gradient(135deg, #a18cd1, #fbc2eb)',
]

function avatarBg(name) {
  const idx = (name || '?').charCodeAt(0) % avatarGradients.length
  return { background: avatarGradients[idx] }
}

function getInitial(name) {
  return name ? name.split(' ').map(n => n[0]).join('').substring(0, 2).toUpperCase() : '?'
}

function formatDate(d) {
  if (!d) return '-'
  return new Date(d).toLocaleDateString('id-ID', { day: 'numeric', month: 'short', year: 'numeric' })
}

function statusLabel(status) {
  const labels = {
    PENDING: 'Menunggu',
    DP_VERIFIED: 'DP Terverifikasi',
    AGREEMENT_SIGNED: 'Perjanjian',
    DEVELOPMENT: 'Pembangunan',
    RUNNING: 'Berjalan',
    COMPLETED: 'Selesai',
  }
  return labels[status] || status
}

function statusClass(status) {
  const classes = {
    PENDING: 'st-pending',
    DP_VERIFIED: 'st-verified',
    AGREEMENT_SIGNED: 'st-agreement',
    DEVELOPMENT: 'st-dev',
    RUNNING: 'st-running',
    COMPLETED: 'st-completed',
  }
  return classes[status] || ''
}
</script>

<style scoped>
/* ═══ HERO ═══ */
.aff-hero {
  background: linear-gradient(135deg, #0f0c29 0%, #302b63 50%, #24243e 100%);
  border-radius: 16px; padding: 32px 36px 28px; margin-bottom: 24px;
  box-shadow: 0 4px 24px rgba(15,12,41,0.2);
}
.aff-hero-content {
  display: flex; align-items: flex-start; justify-content: space-between;
  gap: 16px; flex-wrap: wrap; margin-bottom: 28px;
}
.aff-hero-title { font-size: 1.6rem; font-weight: 800; color: #fff; margin: 0 0 4px; letter-spacing: -0.02em; }
.aff-hero-sub { font-size: 0.85rem; color: rgba(255,255,255,0.5); margin: 0; }

/* Referral Code Box */
.aff-referral-box {
  background: rgba(255,255,255,0.08);
  border: 1px solid rgba(255,255,255,0.12);
  border-radius: 14px; padding: 16px 22px;
  backdrop-filter: blur(10px);
}
.aff-referral-label { font-size: 0.72rem; text-transform: uppercase; letter-spacing: 0.06em; color: rgba(255,255,255,0.4); margin-bottom: 8px; }
.aff-referral-code-wrap { display: flex; align-items: center; gap: 12px; }
.aff-referral-code {
  font-size: 1.3rem; font-weight: 800; color: #a78bfa;
  letter-spacing: 0.08em; font-family: 'SF Mono', 'Fira Code', monospace;
}
.aff-copy-btn {
  display: inline-flex; align-items: center; gap: 6px;
  padding: 8px 16px; font-size: 0.78rem; font-weight: 600;
  border-radius: 10px; border: 1px solid rgba(255,255,255,0.15);
  background: rgba(255,255,255,0.06); color: #e2e8f0; cursor: pointer;
  transition: all 0.2s;
}
.aff-copy-btn:hover { background: rgba(255,255,255,0.12); border-color: rgba(255,255,255,0.25); }
.aff-copy-btn.copied { background: rgba(34,197,94,0.15); border-color: rgba(34,197,94,0.3); color: #4ade80; }

/* Stats Grid */
.aff-stats { display: grid; grid-template-columns: repeat(4, 1fr); gap: 16px; }
.aff-stat-card {
  background: rgba(255,255,255,0.06); border: 1px solid rgba(255,255,255,0.08);
  border-radius: 14px; padding: 20px; display: flex; align-items: center; gap: 16px;
  backdrop-filter: blur(8px); transition: all 0.2s;
}
.aff-stat-card:hover { background: rgba(255,255,255,0.1); border-color: rgba(255,255,255,0.15); transform: translateY(-2px); }
.aff-stat-icon {
  width: 48px; height: 48px; border-radius: 14px; display: flex; align-items: center; justify-content: center;
}
.icon-referral { background: rgba(99,102,241,0.15); color: #818cf8; }
.icon-active { background: rgba(34,197,94,0.15); color: #4ade80; }
.icon-pending { background: rgba(251,191,36,0.15); color: #fbbf24; }
.icon-rate { background: rgba(168,85,247,0.15); color: #c084fc; }
.aff-stat-info { display: flex; flex-direction: column; }
.aff-stat-value { font-size: 1.5rem; font-weight: 800; color: #fff; line-height: 1.2; }
.aff-stat-label { font-size: 0.72rem; color: rgba(255,255,255,0.45); text-transform: uppercase; letter-spacing: 0.04em; margin-top: 2px; }

/* ═══ SECTION ═══ */
.aff-section {
  background: #fff; border-radius: 16px; padding: 28px;
  border: 1px solid #e8ecf1; box-shadow: 0 1px 4px rgba(0,0,0,0.04);
}
.aff-section-header { display: flex; align-items: center; justify-content: space-between; margin-bottom: 24px; gap: 16px; flex-wrap: wrap; }
.aff-section-header h2 { font-size: 1.1rem; font-weight: 700; color: #0f172a; margin: 0; }
.aff-search-wrap { position: relative; }
.aff-search-icon { position: absolute; left: 14px; top: 50%; transform: translateY(-50%); }
.aff-search {
  padding: 10px 14px 10px 40px; border: 1.5px solid #e2e8f0; border-radius: 12px;
  font-size: 0.85rem; background: #fafbfc; color: #1e293b; outline: none; min-width: 280px;
}
.aff-search:focus { border-color: #6366f1; box-shadow: 0 0 0 3px rgba(99,102,241,0.1); }

/* ═══ TABLE ═══ */
.aff-table-wrap { overflow-x: auto; }
.aff-table { width: 100%; border-collapse: collapse; }
.aff-table thead th {
  text-align: left; padding: 12px 16px; font-size: 0.72rem; font-weight: 700;
  color: #94a3b8; text-transform: uppercase; letter-spacing: 0.05em;
  border-bottom: 2px solid #f1f5f9;
}
.aff-table tbody td { padding: 16px; border-bottom: 1px solid #f8fafc; font-size: 0.85rem; color: #334155; }
.aff-table tbody tr:hover { background: #f8fafc; }

.aff-mitra-cell { display: flex; align-items: center; gap: 12px; }
.aff-mitra-avatar {
  width: 36px; height: 36px; border-radius: 10px; display: flex; align-items: center; justify-content: center;
  font-size: 0.75rem; font-weight: 700; color: white; flex-shrink: 0;
}
.aff-mitra-name { font-weight: 600; color: #0f172a; font-size: 0.85rem; }
.aff-mitra-email { font-size: 0.75rem; color: #94a3b8; }

.aff-status-badge {
  display: inline-flex; padding: 5px 12px; border-radius: 8px;
  font-size: 0.72rem; font-weight: 700; text-transform: uppercase; letter-spacing: 0.04em;
}
.st-pending { background: #fef3c7; color: #92400e; }
.st-verified { background: #dbeafe; color: #1e40af; }
.st-agreement { background: #ede9fe; color: #5b21b6; }
.st-dev { background: #fce7f3; color: #9d174d; }
.st-running { background: #dcfce7; color: #166534; }
.st-completed { background: #f0fdf4; color: #14532d; }

.aff-progress-bar {
  width: 80px; height: 6px; border-radius: 3px; background: #e2e8f0;
  overflow: hidden; margin-bottom: 4px;
}
.aff-progress-fill { height: 100%; background: linear-gradient(90deg, #6366f1, #8b5cf6); border-radius: 3px; transition: width 0.3s; }
.aff-progress-text { font-size: 0.72rem; color: #94a3b8; font-weight: 600; }
.aff-date-col { white-space: nowrap; font-size: 0.8rem; color: #94a3b8; }

/* ═══ PAGINATION ═══ */
.aff-pagination { display: flex; align-items: center; justify-content: center; gap: 16px; margin-top: 24px; padding-top: 20px; border-top: 1px solid #f1f5f9; }
.aff-page-btn {
  padding: 8px 18px; border-radius: 10px; font-size: 0.82rem; font-weight: 600;
  border: 1.5px solid #e2e8f0; background: #fff; color: #475569; cursor: pointer;
}
.aff-page-btn:disabled { opacity: 0.4; cursor: not-allowed; }
.aff-page-btn:not(:disabled):hover { border-color: #6366f1; color: #6366f1; background: #eef2ff; }
.aff-page-info { font-size: 0.82rem; color: #94a3b8; }

/* ═══ SKELETON ═══ */
.aff-loading { display: flex; flex-direction: column; gap: 12px; }
.aff-row-skel { height: 56px; border-radius: 12px; }
.shimmer {
  background: linear-gradient(90deg, #e8ecf1 25%, #f1f5f9 50%, #e8ecf1 75%);
  background-size: 200% 100%; animation: shimmer 1.5s infinite;
}
@keyframes shimmer { 0% { background-position: 200% 0 } 100% { background-position: -200% 0 } }

/* ═══ EMPTY ═══ */
.aff-empty {
  display: flex; flex-direction: column; align-items: center; justify-content: center;
  text-align: center; padding: 60px 32px;
}
.aff-empty-circle {
  width: 80px; height: 80px; border-radius: 50%;
  background: linear-gradient(135deg, #f1f5f9, #e2e8f0);
  display: flex; align-items: center; justify-content: center; margin-bottom: 20px;
}
.aff-empty h3 { font-size: 1.05rem; font-weight: 700; color: #0f172a; margin: 0; }
.aff-empty p { font-size: 0.85rem; color: #94a3b8; margin: 6px 0 0; }

@media (max-width: 1024px) {
  .aff-stats { grid-template-columns: repeat(2, 1fr); }
}
@media (max-width: 768px) {
  .aff-hero { padding: 24px 20px 20px; }
  .aff-hero-content { flex-direction: column; }
  .aff-stats { grid-template-columns: 1fr; }
  .aff-section { padding: 20px; }
  .aff-section-header { flex-direction: column; align-items: stretch; }
  .aff-search { min-width: unset; width: 100%; }
}
</style>
