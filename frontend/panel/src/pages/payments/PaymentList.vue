<template>
  <div class="py-page">
    <!-- Hero -->
    <div class="py-hero">
      <div class="py-hero-top">
        <div>
          <h1 class="py-hero-title">Pembayaran</h1>
          <p class="py-hero-sub">Kelola dan verifikasi pembayaran dari mitra</p>
        </div>
      </div>
      <div class="py-stats">
        <div class="py-stat">
          <span class="py-stat-dot" style="background:#818cf8"></span>
          <span class="py-stat-label">Total</span>
          <span class="py-stat-val">{{ allPayments.length }}</span>
        </div>
        <div class="py-stat">
          <span class="py-stat-dot" style="background:#22c55e"></span>
          <span class="py-stat-label">Verified</span>
          <span class="py-stat-val">{{ allPayments.filter(p => p.verified_status === 'VERIFIED').length }}</span>
        </div>
        <div class="py-stat">
          <span class="py-stat-dot" style="background:#f59e0b"></span>
          <span class="py-stat-label">Pending</span>
          <span class="py-stat-val">{{ allPayments.filter(p => p.verified_status === 'PENDING').length }}</span>
        </div>
      </div>
    </div>

    <!-- Toolbar -->
    <div class="py-toolbar">
      <div class="py-search-wrap">
        <svg class="py-search-ico" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="#94a3b8" stroke-width="2"><circle cx="11" cy="11" r="8"/><line x1="21" y1="21" x2="16.65" y2="16.65"/></svg>
        <input v-model="search" class="py-search" placeholder="Cari nama mitra atau tipe..." />
      </div>
      <select v-model="filterMitra" class="py-filter-select">
        <option value="">Semua Mitra</option>
        <option v-for="m in mitraOptions" :key="m.id" :value="m.id">{{ m.name }}</option>
      </select>
    </div>

    <!-- Table -->
    <div class="py-table-wrap">
      <table class="py-table" v-if="filtered.length">
        <thead>
          <tr>
            <th>Mitra</th>
            <th>Tipe</th>
            <th>Jumlah</th>
            <th>Status</th>
            <th>Tanggal</th>
            <th>Aksi</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="pay in filtered" :key="pay.id">
            <td>
              <div class="py-mitra-cell">
                <div class="py-mitra-avatar" :style="avatarBg(pay._mitraName)">{{ getInitial(pay._mitraName) }}</div>
                <div>
                  <div class="py-mitra-name">{{ pay._mitraName }}</div>
                  <div class="py-mitra-email">{{ pay._mitraEmail }}</div>
                </div>
              </div>
            </td>
            <td><span class="py-type-badge" :class="'pt-' + pay.type">{{ pay.type }}</span></td>
            <td><span class="py-amount">{{ fc(pay.amount) }}</span></td>
            <td>
              <span class="py-badge" :class="verifyBadge(pay.verified_status)">
                <span class="py-badge-dot"></span>
                {{ verifyLabel(pay.verified_status) }}
              </span>
            </td>
            <td><span class="py-date">{{ formatDate(pay.created_at) }}</span></td>
            <td>
              <div class="py-actions">
                <button v-if="pay.proof_url" @click="viewProof(pay)" class="py-detail-btn" title="Lihat Bukti">
                  <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M1 12s4-8 11-8 11 8 11 8-4 8-11 8-11-8-11-8z"/><circle cx="12" cy="12" r="3"/></svg>
                  Bukti
                </button>
                <button v-if="pay.verified_status !== 'VERIFIED'" @click="verify(pay)" class="py-verify-btn">
                  <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5"><polyline points="20 6 9 17 4 12"/></svg>
                  Verifikasi
                </button>
              </div>
            </td>
          </tr>
        </tbody>
      </table>
      <div v-else class="py-empty-row">
        <div class="py-empty-circle">
          <svg width="32" height="32" viewBox="0 0 24 24" fill="none" stroke="#94a3b8" stroke-width="1.5">
            <rect x="1" y="4" width="22" height="16" rx="2" ry="2"/>
            <line x1="1" y1="10" x2="23" y2="10"/>
          </svg>
        </div>
        <p>Belum ada data pembayaran</p>
      </div>
    </div>

    <!-- Proof Modal -->
    <Teleport to="body">
      <div v-if="showProofModal" class="py-overlay" @click.self="showProofModal = false">
        <div class="py-modal" @click.stop>
          <div class="py-modal-head">
            <div class="py-modal-title-group">
              <div class="py-modal-icon">
                <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M1 12s4-8 11-8 11 8 11 8-4 8-11 8-11-8-11-8z"/><circle cx="12" cy="12" r="3"/></svg>
              </div>
              <h3>Bukti Pembayaran</h3>
            </div>
            <button @click="showProofModal = false" class="py-modal-x">&times;</button>
          </div>
          <div class="py-modal-body" style="text-align:center">
            <img :src="proofUrl" alt="Bukti Pembayaran" class="py-proof-img" />
          </div>
          <div class="py-modal-foot">
            <a :href="proofUrl" target="_blank" class="py-detail-btn">
              <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M18 13v6a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2V8a2 2 0 0 1 2-2h6"/><polyline points="15 3 21 3 21 9"/><line x1="10" y1="14" x2="21" y2="3"/></svg>
              Buka di Tab Baru
            </a>
            <button @click="showProofModal = false" class="py-btn-sec">Tutup</button>
          </div>
        </div>
      </div>
    </Teleport>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { partnershipApi, paymentApi } from '../../services/api'
import { useToastStore } from '../../stores/toast'

const toast = useToastStore()
const partnerships = ref([])
const allPayments = ref([])
const search = ref('')
const filterMitra = ref('')
const showProofModal = ref(false)
const proofUrl = ref('')

onMounted(async () => {
  try {
    const { data } = await partnershipApi.list({ page: 1, limit: 100 })
    partnerships.value = data.data || []
    await loadAllPayments()
  } catch {
    toast.error('Gagal memuat data')
  }
})

async function loadAllPayments() {
  const results = []
  await Promise.all(
    partnerships.value.map(async (p) => {
      try {
        const { data } = await paymentApi.byPartnership(p.id)
        const items = data.data || []
        items.forEach(pay => {
          pay._partnershipId = p.id
          pay._mitraName = p.mitra?.name || '-'
          pay._mitraEmail = p.mitra?.email || ''
        })
        results.push(...items)
      } catch { /* skip */ }
    })
  )
  // Sort by newest first
  results.sort((a, b) => new Date(b.created_at) - new Date(a.created_at))
  allPayments.value = results
}

const mitraOptions = computed(() => {
  const map = new Map()
  partnerships.value.forEach(p => {
    if (p.mitra) map.set(p.mitra_id || p.id, { id: p.id, name: p.mitra.name })
  })
  return Array.from(map.values()).sort((a, b) => a.name.localeCompare(b.name))
})

const filtered = computed(() => {
  let list = allPayments.value
  if (filterMitra.value) {
    list = list.filter(p => p._partnershipId === filterMitra.value)
  }
  const q = search.value.toLowerCase().trim()
  if (q) {
    list = list.filter(p =>
      p._mitraName.toLowerCase().includes(q) ||
      (p.type || '').toLowerCase().includes(q)
    )
  }
  return list
})

async function verify(pay) {
  try {
    await paymentApi.verify(pay.id, { status: 'VERIFIED' })
    toast.success('Pembayaran berhasil diverifikasi')
    await loadAllPayments()
  } catch (e) {
    toast.error(e.response?.data?.error || 'Gagal verifikasi')
  }
}

function viewProof(pay) {
  proofUrl.value = pay.proof_url
  showProofModal.value = true
}

// Helpers
const avatarGradients = [
  'linear-gradient(135deg, #667eea, #764ba2)',
  'linear-gradient(135deg, #f093fb, #f5576c)',
  'linear-gradient(135deg, #4facfe, #00f2fe)',
  'linear-gradient(135deg, #43e97b, #38f9d7)',
  'linear-gradient(135deg, #fa709a, #fee140)',
  'linear-gradient(135deg, #a18cd1, #fbc2eb)',
]
function avatarBg(name) { return { background: avatarGradients[(name || '?').charCodeAt(0) % avatarGradients.length] } }
function getInitial(name) { return name ? name.split(' ').map(n => n[0]).join('').substring(0, 2).toUpperCase() : '?' }
function formatDate(d) { return d ? new Date(d).toLocaleDateString('id-ID', { day: 'numeric', month: 'short', year: 'numeric' }) : '-' }
function fc(v) { return new Intl.NumberFormat('id-ID', { style: 'currency', currency: 'IDR', minimumFractionDigits: 0 }).format(v) }

function verifyBadge(s) {
  return { VERIFIED: 'st-verified', PENDING: 'st-pending', REJECTED: 'st-rejected' }[s] || 'st-pending'
}
function verifyLabel(s) {
  return { VERIFIED: 'Terverifikasi', PENDING: 'Menunggu', REJECTED: 'Ditolak' }[s] || s
}
</script>

<style scoped>
/* ═══ HERO ═══ */
.py-hero { background: linear-gradient(135deg, #0f0c29 0%, #302b63 50%, #24243e 100%); border-radius: 16px; padding: 32px 36px 24px; margin-bottom: 24px; box-shadow: 0 4px 24px rgba(15,12,41,0.2); }
.py-hero-top { display: flex; justify-content: space-between; align-items: flex-start; gap: 16px; flex-wrap: wrap; margin-bottom: 20px; }
.py-hero-title { font-size: 1.6rem; font-weight: 800; color: #fff; margin: 0 0 4px; letter-spacing: -0.02em; }
.py-hero-sub { font-size: .85rem; color: rgba(255,255,255,.5); margin: 0; }
.py-stats { display: flex; gap: 28px; flex-wrap: wrap; padding-top: 16px; border-top: 1px solid rgba(255,255,255,.08); }
.py-stat { display: flex; align-items: center; gap: 8px; }
.py-stat-dot { width: 8px; height: 8px; border-radius: 50%; }
.py-stat-label { font-size: .72rem; color: rgba(255,255,255,.4); text-transform: uppercase; letter-spacing: .05em; }
.py-stat-val { font-size: .9rem; font-weight: 800; color: #fff; }

/* ═══ TOOLBAR ═══ */
.py-toolbar { display: flex; gap: 12px; margin-bottom: 20px; flex-wrap: wrap; }
.py-search-wrap { position: relative; flex: 1; max-width: 400px; }
.py-search-ico { position: absolute; left: 14px; top: 50%; transform: translateY(-50%); }
.py-search { width: 100%; padding: 10px 14px 10px 40px; border: 1.5px solid #e2e8f0; border-radius: 12px; font-size: .85rem; background: #fff; color: #1e293b; outline: none; box-sizing: border-box; font-family: inherit; }
.py-search:focus { border-color: #6366f1; box-shadow: 0 0 0 3px rgba(99,102,241,.1); }

.py-filter-select {
  padding: 10px 36px 10px 14px; border: 1.5px solid #e2e8f0; border-radius: 12px;
  font-size: .85rem; background: #fff; color: #1e293b; outline: none; font-family: inherit;
  appearance: none; -webkit-appearance: none; cursor: pointer;
  background-image: url("data:image/svg+xml,%3csvg xmlns='http://www.w3.org/2000/svg' fill='none' viewBox='0 0 20 20'%3e%3cpath stroke='%236b7280' stroke-linecap='round' stroke-linejoin='round' stroke-width='1.5' d='M6 8l4 4 4-4'/%3e%3c/svg%3e");
  background-position: right 12px center; background-repeat: no-repeat; background-size: 16px;
}
.py-filter-select:focus { border-color: #6366f1; box-shadow: 0 0 0 3px rgba(99,102,241,.1); }

/* ═══ TABLE ═══ */
.py-table-wrap { background: #fff; border-radius: 16px; border: 1px solid #e8ecf1; overflow: hidden; }
.py-table { width: 100%; border-collapse: collapse; }
.py-table thead { background: #f8fafc; }
.py-table th { padding: 14px 20px; font-size: .72rem; font-weight: 700; color: #64748b; text-transform: uppercase; letter-spacing: .05em; text-align: left; border-bottom: 1px solid #e8ecf1; }
.py-table td { padding: 16px 20px; font-size: .85rem; color: #1e293b; border-bottom: 1px solid #f1f5f9; vertical-align: middle; }
.py-table tr:last-child td { border-bottom: none; }
.py-table tbody tr { transition: background .15s; }
.py-table tbody tr:hover { background: #fafbfc; }

/* Mitra cell */
.py-mitra-cell { display: flex; align-items: center; gap: 12px; }
.py-mitra-avatar { width: 36px; height: 36px; border-radius: 10px; display: flex; align-items: center; justify-content: center; font-size: .72rem; font-weight: 800; color: #fff; flex-shrink: 0; }
.py-mitra-name { font-weight: 700; font-size: .85rem; color: #0f172a; }
.py-mitra-email { font-size: .7rem; color: #94a3b8; }

.py-amount { font-weight: 700; font-size: .88rem; color: #0f172a; }
.py-date { font-size: .82rem; color: #64748b; }

/* Type badge */
.py-type-badge { font-size: .68rem; font-weight: 700; padding: 4px 10px; border-radius: 6px; text-transform: uppercase; letter-spacing: .03em; }
.pt-DP { background: #e0f2fe; color: #0284c7; }
.pt-PELUNASAN { background: #dcfce7; color: #16a34a; }
.pt-CICILAN { background: #ede9fe; color: #7c3aed; }
.pt-INVOICE { background: #fef3c7; color: #d97706; }

/* Status badge */
.py-badge { font-size: .68rem; font-weight: 700; padding: 4px 10px; border-radius: 6px; text-transform: uppercase; letter-spacing: .03em; display: inline-flex; align-items: center; gap: 5px; }
.py-badge-dot { width: 6px; height: 6px; border-radius: 50%; }
.st-verified { background: #dcfce7; color: #16a34a; }
.st-verified .py-badge-dot { background: #22c55e; }
.st-pending { background: #fef3c7; color: #d97706; }
.st-pending .py-badge-dot { background: #f59e0b; }
.st-rejected { background: #fee2e2; color: #dc2626; }
.st-rejected .py-badge-dot { background: #ef4444; }

/* Actions */
.py-actions { display: flex; gap: 6px; }
.py-detail-btn {
  display: inline-flex; align-items: center; gap: 5px;
  font-size: .78rem; font-weight: 600; color: #6366f1; text-decoration: none;
  padding: 6px 14px; border-radius: 8px; border: 1px solid #e0e7ff; background: #eef2ff;
  cursor: pointer; transition: all .15s; font-family: inherit;
}
.py-detail-btn:hover { background: #e0e7ff; }

.py-verify-btn {
  display: inline-flex; align-items: center; gap: 5px;
  font-size: .78rem; font-weight: 600; color: #fff;
  padding: 6px 14px; border-radius: 8px; border: none;
  background: linear-gradient(135deg, #22c55e, #16a34a); cursor: pointer;
  box-shadow: 0 2px 8px rgba(34,197,94,.25); transition: all .15s; font-family: inherit;
}
.py-verify-btn:hover { box-shadow: 0 4px 14px rgba(34,197,94,.35); transform: translateY(-1px); }

/* Empty */
.py-empty-row { text-align: center; padding: 56px 20px; }
.py-empty-circle { width: 72px; height: 72px; border-radius: 50%; background: linear-gradient(135deg, #f1f5f9, #e2e8f0); display: flex; align-items: center; justify-content: center; margin: 0 auto 16px; }
.py-empty-row p { color: #94a3b8; font-size: .85rem; margin: 0; }

/* ═══ MODAL ═══ */
.py-overlay { position: fixed; inset: 0; background: rgba(0,0,0,.5); display: flex; align-items: center; justify-content: center; z-index: 1000; backdrop-filter: blur(4px); animation: fadeIn .2s ease; }
@keyframes fadeIn { from { opacity: 0 } to { opacity: 1 } }
.py-modal { background: #fff; border-radius: 18px; width: 100%; max-width: 600px; box-shadow: 0 24px 80px rgba(0,0,0,.2); animation: slideUp .3s ease; overflow: hidden; }
@keyframes slideUp { from { transform: translateY(20px); opacity: 0 } to { transform: translateY(0); opacity: 1 } }
.py-modal-head { display: flex; align-items: center; justify-content: space-between; padding: 22px 28px; border-bottom: 1px solid #f1f5f9; }
.py-modal-title-group { display: flex; align-items: center; gap: 12px; }
.py-modal-icon { width: 40px; height: 40px; border-radius: 12px; background: linear-gradient(135deg, #eef2ff, #e0e7ff); display: flex; align-items: center; justify-content: center; color: #6366f1; }
.py-modal-head h3 { font-size: 1.1rem; font-weight: 700; margin: 0; color: #0f172a; }
.py-modal-x { width: 34px; height: 34px; border-radius: 10px; display: flex; align-items: center; justify-content: center; border: none; background: transparent; font-size: 1.4rem; color: #94a3b8; cursor: pointer; }
.py-modal-x:hover { background: #f1f5f9; color: #0f172a; }
.py-modal-body { padding: 24px 28px; }
.py-modal-foot { display: flex; justify-content: flex-end; gap: 10px; padding: 16px 28px; border-top: 1px solid #f1f5f9; }
.py-btn-sec { padding: 10px 22px; border-radius: 12px; font-size: .85rem; font-weight: 600; background: #f1f5f9; color: #475569; border: none; cursor: pointer; font-family: inherit; }
.py-btn-sec:hover { background: #e2e8f0; }
.py-proof-img { max-width: 100%; max-height: 55vh; border-radius: 10px; object-fit: contain; }

/* ═══ RESPONSIVE ═══ */
@media (max-width: 768px) {
  .py-hero { padding: 24px 20px 18px; }
  .py-toolbar { flex-direction: column; }
  .py-search-wrap { max-width: 100%; }
  .py-filter-select { width: 100%; }
  .py-table-wrap { overflow-x: auto; }
}
</style>
