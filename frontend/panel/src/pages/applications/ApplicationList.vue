<template>
  <div class="ap-page">
    <!-- Hero -->
    <div class="ap-hero">
      <div class="ap-hero-top">
        <div>
          <h1 class="ap-hero-title">Pengajuan Partnership</h1>
          <p class="ap-hero-sub">Review dan kelola pengajuan kemitraan dari mitra</p>
        </div>
      </div>
      <div class="ap-stats">
        <div class="ap-stat"><span class="ap-stat-dot" style="background:#f59e0b"></span><span class="ap-stat-label">Pending</span><span class="ap-stat-val">{{ countByStatus('PENDING') }}</span></div>
        <div class="ap-stat"><span class="ap-stat-dot" style="background:#3b82f6"></span><span class="ap-stat-label">Reviewed</span><span class="ap-stat-val">{{ countByStatus('REVIEWED') }}</span></div>
        <div class="ap-stat"><span class="ap-stat-dot" style="background:#22c55e"></span><span class="ap-stat-label">Approved</span><span class="ap-stat-val">{{ countByStatus('APPROVED') }}</span></div>
        <div class="ap-stat"><span class="ap-stat-dot" style="background:#ef4444"></span><span class="ap-stat-label">Rejected</span><span class="ap-stat-val">{{ countByStatus('REJECTED') }}</span></div>
      </div>
    </div>

    <!-- Toolbar -->
    <div class="ap-toolbar">
      <div class="ap-search-wrap">
        <svg class="ap-search-ico" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="#94a3b8" stroke-width="2"><circle cx="11" cy="11" r="8"/><line x1="21" y1="21" x2="16.65" y2="16.65"/></svg>
        <input v-model="search" class="ap-search" placeholder="Cari nama mitra atau outlet..." />
      </div>
      <div class="ap-filter-wrap">
        <button v-for="s in statusFilters" :key="s.value" class="ap-filter-btn" :class="{ active: statusFilter === s.value }" @click="statusFilter = s.value">{{ s.label }}</button>
      </div>
    </div>

    <!-- Table -->
    <div class="ap-table-wrap">
      <table class="ap-table">
        <thead>
          <tr>
            <th>Mitra</th>
            <th>Outlet / Paket</th>
            <th>Lokasi Diusulkan</th>
            <th>Budget</th>
            <th>Tanggal</th>
            <th>Status</th>
            <th>Aksi</th>
          </tr>
        </thead>
        <tbody>
          <tr v-if="loading">
            <td colspan="7" class="ap-loading-row">
              <div class="ap-spinner"></div>
              Memuat data...
            </td>
          </tr>
          <tr v-else-if="!filtered.length">
            <td colspan="7" class="ap-empty-row">Belum ada pengajuan</td>
          </tr>
          <tr v-for="a in filtered" :key="a.id" class="ap-row" @click="openDetail(a)">
            <td>
              <div class="ap-mitra-cell">
                <div class="ap-mitra-avatar" :style="avatarBg(a.mitra?.name || '?')">{{ getInitial(a.mitra?.name || '?') }}</div>
                <div>
                  <div class="ap-mitra-name">{{ a.mitra?.name || '-' }}</div>
                  <div class="ap-mitra-email">{{ a.contact_email || a.mitra?.email }}</div>
                </div>
              </div>
            </td>
            <td>
              <div class="ap-outlet-name">{{ a.outlet?.name || '-' }}</div>
              <div class="ap-package-name" v-if="a.package">{{ a.package.name }} — {{ fc(a.package.price) }}</div>
            </td>
            <td>
              <div class="ap-location">{{ a.proposed_location || '-' }}</div>
            </td>
            <td>
              <div class="ap-budget">{{ fc(a.investment_budget) }}</div>
            </td>
            <td>
              <div class="ap-date">{{ fDate(a.created_at) }}</div>
            </td>
            <td>
              <span class="ap-badge" :class="'st-' + a.status">{{ statusLabel(a.status) }}</span>
            </td>
            <td @click.stop>
              <div class="ap-actions">
                <button v-if="a.status === 'PENDING'" class="ap-btn ap-btn-review" @click="openReview(a)">Review</button>
                <button v-if="a.status === 'PENDING' || a.status === 'REVIEWED'" class="ap-btn ap-btn-approve" @click="reviewApp(a.id, 'APPROVED')">
                  <i class="ri-check-line"></i> Approve
                </button>
                <button v-if="a.status === 'PENDING' || a.status === 'REVIEWED'" class="ap-btn ap-btn-reject" @click="reviewApp(a.id, 'REJECTED')">
                  <i class="ri-close-line"></i>
                </button>
              </div>
            </td>
          </tr>
        </tbody>
      </table>
    </div>

    <!-- Review Modal -->
    <Teleport to="body">
      <div v-if="showModal" class="ap-overlay" @click.self="showModal = false">
        <div class="ap-modal" @click.stop>
          <div class="ap-modal-head">
            <h3>Review Pengajuan</h3>
            <button @click="showModal = false" class="ap-modal-x">&times;</button>
          </div>
          <div class="ap-modal-body" v-if="selectedApp">
            <div class="ap-detail-grid">
              <div class="ap-detail-item">
                <label>Mitra</label>
                <p>{{ selectedApp.mitra?.name }} ({{ selectedApp.contact_email }})</p>
              </div>
              <div class="ap-detail-item">
                <label>No. HP</label>
                <p>{{ selectedApp.contact_phone || '-' }}</p>
              </div>
              <div class="ap-detail-item">
                <label>Outlet</label>
                <p>{{ selectedApp.outlet?.name || '-' }}</p>
              </div>
              <div class="ap-detail-item">
                <label>Paket</label>
                <p>{{ selectedApp.package?.name || '-' }} — {{ fc(selectedApp.package?.price || 0) }}</p>
              </div>
              <div class="ap-detail-item full">
                <label>Lokasi Diusulkan</label>
                <p>{{ selectedApp.proposed_location || '-' }}</p>
              </div>
              <div class="ap-detail-item">
                <label>Budget Investasi</label>
                <p>{{ fc(selectedApp.investment_budget) }}</p>
              </div>
              <div class="ap-detail-item full">
                <label>Motivasi</label>
                <p class="ap-detail-text">{{ selectedApp.motivation || '-' }}</p>
              </div>
              <div class="ap-detail-item full">
                <label>Pengalaman</label>
                <p class="ap-detail-text">{{ selectedApp.experience || '-' }}</p>
              </div>
            </div>

            <div class="ap-review-section">
              <label>Catatan Admin</label>
              <textarea v-model="adminNotes" class="ap-textarea" rows="3" placeholder="Tulis catatan review..."></textarea>
            </div>

            <div class="ap-modal-actions">
              <button class="ap-btn ap-btn-approve" @click="submitReview('APPROVED')" :disabled="reviewLoading">
                <i class="ri-check-double-line"></i> Approve & Buat Partnership
              </button>
              <button class="ap-btn ap-btn-review" @click="submitReview('REVIEWED')" :disabled="reviewLoading">
                <i class="ri-eye-line"></i> Mark as Reviewed
              </button>
              <button class="ap-btn ap-btn-reject" @click="submitReview('REJECTED')" :disabled="reviewLoading">
                <i class="ri-close-line"></i> Reject
              </button>
            </div>
          </div>
        </div>
      </div>
    </Teleport>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted } from 'vue'
import { applicationApi } from '../../services/api'
import { useToastStore } from '../../stores/toast'

const toast = useToastStore()
const loading = ref(false)
const apps = ref([])
const search = ref('')
const statusFilter = ref('')
const showModal = ref(false)
const selectedApp = ref(null)
const adminNotes = ref('')
const reviewLoading = ref(false)

const statusFilters = [
  { label: 'Semua', value: '' },
  { label: 'Pending', value: 'PENDING' },
  { label: 'Reviewed', value: 'REVIEWED' },
  { label: 'Approved', value: 'APPROVED' },
  { label: 'Rejected', value: 'REJECTED' },
]

const filtered = computed(() => {
  let list = apps.value
  if (statusFilter.value) list = list.filter(a => a.status === statusFilter.value)
  if (search.value) {
    const q = search.value.toLowerCase()
    list = list.filter(a =>
      (a.mitra?.name || '').toLowerCase().includes(q) ||
      (a.outlet?.name || '').toLowerCase().includes(q) ||
      (a.contact_email || '').toLowerCase().includes(q)
    )
  }
  return list
})

function countByStatus(s) { return apps.value.filter(a => a.status === s).length }

function statusLabel(s) {
  const map = { PENDING: 'Menunggu', REVIEWED: 'Ditinjau', APPROVED: 'Disetujui', REJECTED: 'Ditolak', CANCELLED: 'Dibatalkan' }
  return map[s] || s
}

function fc(n) { return 'Rp ' + (n || 0).toLocaleString('id-ID') }

function fDate(d) {
  if (!d) return '-'
  return new Date(d).toLocaleDateString('id-ID', { day: 'numeric', month: 'short', year: 'numeric' })
}

function getInitial(name) { return (name || '?')[0].toUpperCase() }

function avatarBg(name) {
  const colors = ['#6366f1','#8b5cf6','#ec4899','#f43f5e','#f97316','#14b8a6','#06b6d4','#3b82f6']
  const idx = (name || '').charCodeAt(0) % colors.length
  return { background: colors[idx] }
}

function openDetail(a) {
  openReview(a)
}

function openReview(a) {
  selectedApp.value = a
  adminNotes.value = a.admin_notes || ''
  showModal.value = true
}

async function submitReview(status) {
  if (!selectedApp.value) return
  reviewLoading.value = true
  try {
    await applicationApi.review(selectedApp.value.id, { status, admin_notes: adminNotes.value })
    toast.success(status === 'APPROVED' ? 'Pengajuan disetujui! Partnership berhasil dibuat.' : `Status diubah ke ${statusLabel(status)}`)
    showModal.value = false
    await fetchApps()
  } catch (e) {
    toast.error(e.response?.data?.error || 'Gagal mengubah status')
  } finally {
    reviewLoading.value = false
  }
}

async function reviewApp(id, status) {
  if (status === 'APPROVED' && !confirm('Setujui pengajuan ini? Partnership akan otomatis dibuat.')) return
  if (status === 'REJECTED' && !confirm('Tolak pengajuan ini?')) return
  try {
    await applicationApi.review(id, { status, admin_notes: '' })
    toast.success(status === 'APPROVED' ? 'Disetujui! Partnership dibuat.' : 'Pengajuan ditolak.')
    await fetchApps()
  } catch (e) {
    toast.error(e.response?.data?.error || 'Gagal')
  }
}

async function fetchApps() {
  loading.value = true
  try {
    const { data } = await applicationApi.list({ page: 1, limit: 200 })
    apps.value = data.data || []
  } catch (e) {
    toast.error('Gagal memuat pengajuan')
  } finally {
    loading.value = false
  }
}

onMounted(fetchApps)
</script>

<style scoped>
.ap-page { padding: 0; }

/* Hero */
.ap-hero { background: linear-gradient(135deg, #0f172a 0%, #1e293b 100%); border-radius: 16px; padding: 28px 32px 20px; margin-bottom: 20px; }
.ap-hero-title { font-size: 1.5rem; font-weight: 700; color: #fff; margin: 0 0 4px; }
.ap-hero-sub { font-size: .85rem; color: #94a3b8; margin: 0; }
.ap-stats { display: flex; gap: 24px; margin-top: 16px; padding-top: 16px; border-top: 1px solid rgba(255,255,255,.08); }
.ap-stat { display: flex; align-items: center; gap: 8px; }
.ap-stat-dot { width: 8px; height: 8px; border-radius: 50%; flex-shrink: 0; }
.ap-stat-label { font-size: .78rem; color: #94a3b8; }
.ap-stat-val { font-size: .85rem; font-weight: 700; color: #fff; }

/* Toolbar */
.ap-toolbar { display: flex; align-items: center; gap: 16px; margin-bottom: 16px; flex-wrap: wrap; }
.ap-search-wrap { position: relative; flex: 1; min-width: 200px; }
.ap-search-ico { position: absolute; left: 14px; top: 50%; transform: translateY(-50%); }
.ap-search { width: 100%; padding: 10px 14px 10px 40px; border: 1px solid #e2e8f0; border-radius: 10px; font-size: .85rem; background: #fff; transition: border .2s; }
.ap-search:focus { outline: none; border-color: #6366f1; box-shadow: 0 0 0 3px rgba(99,102,241,.1); }
.ap-filter-wrap { display: flex; gap: 6px; }
.ap-filter-btn { padding: 7px 14px; border: 1px solid #e2e8f0; border-radius: 8px; font-size: .78rem; font-weight: 600; color: #64748b; background: #fff; cursor: pointer; transition: all .2s; }
.ap-filter-btn:hover { border-color: #6366f1; color: #6366f1; }
.ap-filter-btn.active { background: #6366f1; color: #fff; border-color: #6366f1; }

/* Table */
.ap-table-wrap { background: #fff; border-radius: 12px; border: 1px solid #e2e8f0; overflow: hidden; }
.ap-table { width: 100%; border-collapse: collapse; }
.ap-table th { text-align: left; padding: 12px 16px; font-size: .75rem; font-weight: 600; color: #64748b; text-transform: uppercase; letter-spacing: .5px; background: #f8fafc; border-bottom: 1px solid #e2e8f0; }
.ap-table td { padding: 14px 16px; border-bottom: 1px solid #f1f5f9; font-size: .85rem; color: #334155; vertical-align: middle; }
.ap-row { cursor: pointer; transition: background .15s; }
.ap-row:hover { background: #f8fafc; }
.ap-loading-row, .ap-empty-row { text-align: center; padding: 40px 16px !important; color: #94a3b8; font-size: .85rem; }

/* Mitra Cell */
.ap-mitra-cell { display: flex; align-items: center; gap: 10px; }
.ap-mitra-avatar { width: 36px; height: 36px; border-radius: 10px; display: flex; align-items: center; justify-content: center; color: #fff; font-weight: 700; font-size: .85rem; flex-shrink: 0; }
.ap-mitra-name { font-weight: 600; color: #0f172a; font-size: .85rem; }
.ap-mitra-email { font-size: .75rem; color: #94a3b8; }
.ap-outlet-name { font-weight: 600; color: #0f172a; }
.ap-package-name { font-size: .75rem; color: #64748b; margin-top: 2px; }
.ap-location { font-size: .82rem; color: #475569; max-width: 180px; overflow: hidden; text-overflow: ellipsis; white-space: nowrap; }
.ap-budget { font-weight: 700; color: #0f172a; }
.ap-date { font-size: .8rem; color: #64748b; white-space: nowrap; }

/* Badge */
.ap-badge { padding: 4px 12px; border-radius: 20px; font-size: .72rem; font-weight: 700; text-transform: uppercase; letter-spacing: .3px; }
.st-PENDING { background: #fef3c7; color: #92400e; }
.st-REVIEWED { background: #dbeafe; color: #1e40af; }
.st-APPROVED { background: #dcfce7; color: #166534; }
.st-REJECTED { background: #fee2e2; color: #991b1b; }
.st-CANCELLED { background: #f1f5f9; color: #64748b; }

/* Actions */
.ap-actions { display: flex; gap: 6px; align-items: center; }
.ap-btn { display: inline-flex; align-items: center; gap: 4px; padding: 6px 12px; border: none; border-radius: 8px; font-size: .78rem; font-weight: 600; cursor: pointer; transition: all .15s; }
.ap-btn-review { background: #dbeafe; color: #1e40af; }
.ap-btn-review:hover { background: #bfdbfe; }
.ap-btn-approve { background: #dcfce7; color: #166534; }
.ap-btn-approve:hover { background: #bbf7d0; }
.ap-btn-reject { background: #fee2e2; color: #991b1b; padding: 6px 8px; }
.ap-btn-reject:hover { background: #fecaca; }

/* Spinner */
.ap-spinner { width: 20px; height: 20px; border: 2.5px solid #e2e8f0; border-top-color: #6366f1; border-radius: 50%; animation: spin .6s linear infinite; display: inline-block; margin-right: 8px; vertical-align: middle; }
@keyframes spin { to { transform: rotate(360deg) } }

/* Modal */
.ap-overlay { position: fixed; inset: 0; background: rgba(0,0,0,.5); backdrop-filter: blur(4px); display: flex; align-items: center; justify-content: center; z-index: 1000; padding: 20px; }
.ap-modal { background: #fff; border-radius: 16px; width: 100%; max-width: 640px; max-height: 90vh; overflow-y: auto; box-shadow: 0 25px 50px rgba(0,0,0,.15); }
.ap-modal-head { display: flex; align-items: center; justify-content: space-between; padding: 20px 24px; border-bottom: 1px solid #e2e8f0; }
.ap-modal-head h3 { font-size: 1.1rem; font-weight: 700; color: #0f172a; margin: 0; }
.ap-modal-x { width: 32px; height: 32px; border: none; background: #f1f5f9; border-radius: 8px; font-size: 1.2rem; cursor: pointer; display: flex; align-items: center; justify-content: center; color: #64748b; transition: all .15s; }
.ap-modal-x:hover { background: #e2e8f0; }
.ap-modal-body { padding: 24px; }

/* Detail Grid */
.ap-detail-grid { display: grid; grid-template-columns: 1fr 1fr; gap: 16px; margin-bottom: 20px; }
.ap-detail-item label { display: block; font-size: .72rem; font-weight: 600; color: #94a3b8; text-transform: uppercase; letter-spacing: .5px; margin-bottom: 4px; }
.ap-detail-item p { font-size: .88rem; color: #0f172a; margin: 0; font-weight: 500; }
.ap-detail-item.full { grid-column: 1 / -1; }
.ap-detail-text { background: #f8fafc; padding: 10px 14px; border-radius: 8px; font-size: .82rem !important; line-height: 1.6; color: #475569 !important; white-space: pre-wrap; }

/* Review */
.ap-review-section { margin-bottom: 20px; }
.ap-review-section label { display: block; font-size: .78rem; font-weight: 600; color: #334155; margin-bottom: 6px; }
.ap-textarea { width: 100%; padding: 10px 14px; border: 1px solid #e2e8f0; border-radius: 10px; font-size: .85rem; resize: vertical; font-family: inherit; box-sizing: border-box; transition: border .2s; }
.ap-textarea:focus { outline: none; border-color: #6366f1; box-shadow: 0 0 0 3px rgba(99,102,241,.1); }

.ap-modal-actions { display: flex; gap: 8px; flex-wrap: wrap; }
.ap-modal-actions .ap-btn { padding: 10px 18px; font-size: .82rem; }
.ap-modal-actions .ap-btn:disabled { opacity: .5; cursor: not-allowed; }

@media (max-width: 768px) {
  .ap-hero { padding: 20px; }
  .ap-stats { flex-wrap: wrap; gap: 12px; }
  .ap-detail-grid { grid-template-columns: 1fr; }
  .ap-filter-wrap { width: 100%; overflow-x: auto; }
  .ap-table { font-size: .8rem; }
}
</style>
