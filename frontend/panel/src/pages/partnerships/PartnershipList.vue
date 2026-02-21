<template>
  <div class="ps-page">
    <!-- Hero -->
    <div class="ps-hero">
      <div class="ps-hero-top">
        <div>
          <h1 class="ps-hero-title">Partnership Management</h1>
          <p class="ps-hero-sub">Kelola kemitraan, pembayaran, dan perkembangan</p>
        </div>
        <button @click="openCreate" class="ps-btn-primary">
          <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><line x1="12" y1="5" x2="12" y2="19"/><line x1="5" y1="12" x2="19" y2="12"/></svg>
          Buat Partnership
        </button>
      </div>
      <div class="ps-stats">
        <div class="ps-stat"><span class="ps-stat-dot" style="background:#818cf8"></span><span class="ps-stat-label">Total</span><span class="ps-stat-val">{{ total }}</span></div>
        <div class="ps-stat"><span class="ps-stat-dot" style="background:#f59e0b"></span><span class="ps-stat-label">Pending</span><span class="ps-stat-val">{{ partnerships.filter(p=>p.status==='PENDING').length }}</span></div>
        <div class="ps-stat"><span class="ps-stat-dot" style="background:#22c55e"></span><span class="ps-stat-label">Running</span><span class="ps-stat-val">{{ partnerships.filter(p=>p.status==='RUNNING').length }}</span></div>
        <div class="ps-stat"><span class="ps-stat-dot" style="background:#06b6d4"></span><span class="ps-stat-label">Completed</span><span class="ps-stat-val">{{ partnerships.filter(p=>p.status==='COMPLETED').length }}</span></div>
      </div>
    </div>

    <!-- Toolbar -->
    <div class="ps-toolbar">
      <div class="ps-search-wrap">
        <svg class="ps-search-ico" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="#94a3b8" stroke-width="2"><circle cx="11" cy="11" r="8"/><line x1="21" y1="21" x2="16.65" y2="16.65"/></svg>
        <input v-model="search" class="ps-search" placeholder="Cari nama mitra..." />
      </div>
    </div>

    <!-- Table -->
    <div class="ps-table-wrap">
      <table class="ps-table">
        <thead>
          <tr>
            <th>Mitra</th>
            <th>Leader</th>
            <th>Outlet & Paket</th>
            <th>Pembayaran</th>
            <th>Progress</th>
            <th>Status</th>
            <th>Aksi</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="p in filtered" :key="p.id">
            <td>
              <div class="ps-mitra-cell">
                <div class="ps-mitra-avatar" :style="avatarBg(p.mitra?.name||'?')">{{ getInitial(p.mitra?.name||'?') }}</div>
                <div>
                  <div class="ps-mitra-name">{{ p.mitra?.name || '-' }}</div>
                  <div class="ps-mitra-email">{{ p.mitra?.email || '' }}</div>
                </div>
              </div>
            </td>
            <td>
              <div class="ps-lead-cell">{{ p.leader?.name || '-' }}</div>
              <div class="ps-lead-phone">{{ p.leader?.email || '' }}</div>
            </td>
            <td>
              <div class="ps-lead-cell">{{ p.outlet?.name || '-' }}</div>
              <div class="ps-lead-phone" v-if="p.package">{{ p.package?.name }} — {{ fc(p.package?.price||0) }}</div>
            </td>
            <td>
              <div class="ps-pay-cell">
                <div class="ps-pay-row">
                  <span class="ps-pay-amount">{{ fc(paidAmount(p)) }}</span>
                  <span class="ps-pay-sep">/</span>
                  <span class="ps-pay-total">{{ fc(p.package?.price||0) }}</span>
                </div>
                <div class="ps-progress-bar">
                  <div class="ps-progress-fill" :style="{width: payPercent(p)+'%', background: payPercent(p) >= 100 ? '#22c55e' : '#6366f1'}"></div>
                </div>
                <div class="ps-pay-footer">
                  <span class="ps-pay-pct">{{ payPercent(p) }}%</span>
                  <a v-if="latestPayLink(p)" :href="latestPayLink(p)" target="_blank" class="ps-pay-link" title="Buka link pembayaran">
                    <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M18 13v6a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2V8a2 2 0 0 1 2-2h6"/><polyline points="15 3 21 3 21 9"/><line x1="10" y1="14" x2="21" y2="3"/></svg>
                    Bayar
                  </a>
                </div>
              </div>
            </td>
            <td>
              <div class="ps-progress-cell">
                <div class="ps-progress-bar"><div class="ps-progress-fill" :style="{width: p.progress_percentage+'%'}"></div></div>
                <span class="ps-progress-text">{{ p.progress_percentage }}%</span>
              </div>
            </td>
            <td><span class="ps-badge" :class="'st-'+p.status">{{ p.status }}</span></td>
            <td>
              <router-link :to="`/partnerships/${p.id}`" class="ps-detail-btn">Detail →</router-link>
            </td>
          </tr>
          <tr v-if="!filtered.length">
            <td colspan="8" class="ps-empty-row">Belum ada partnership</td>
          </tr>
        </tbody>
      </table>
    </div>

    <!-- Create Modal -->
    <Teleport to="body">
      <div v-if="showModal" class="ps-overlay" @click.self="showModal=false">
        <div class="ps-modal" @click.stop>
          <div class="ps-modal-head">
            <div class="ps-modal-title-group">
              <div class="ps-modal-icon">
                <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M17 21v-2a4 4 0 0 0-4-4H5a4 4 0 0 0-4 4v2"/><circle cx="9" cy="7" r="4"/><path d="M23 21v-2a4 4 0 0 0-3-3.87"/><path d="M16 3.13a4 4 0 0 1 0 7.75"/></svg>
              </div>
              <h3>Buat Partnership Baru</h3>
            </div>
            <button @click="showModal=false" class="ps-modal-x">&times;</button>
          </div>
          <form @submit.prevent="create" class="ps-modal-body">
            <div class="ps-fg">
              <label>Pilih Leader <span class="req">*</span></label>
              <select v-model="form.leader_id" class="ps-input" required>
                <option value="">— Pilih Leader —</option>
                <option v-for="l in leaders" :key="l.id" :value="l.id">
                  {{ l.name }} — {{ l.email }}
                </option>
              </select>
              <div class="ps-field-hint">Leader yang akan menangani mitra ini</div>
            </div>

            <div class="ps-fg">
              <label>Pilih Mitra <span class="req">*</span></label>
              <select v-model="form.mitra_id" class="ps-input" required>
                <option value="">— Pilih Mitra —</option>
                <option v-for="m in mitras" :key="m.id" :value="m.id">
                  {{ m.name }} — {{ m.email }}
                </option>
              </select>
              <div class="ps-field-hint">User mitra yang akan menjadi partner</div>
            </div>

            <div class="ps-fg">
              <label>Pilih Outlet <span class="req">*</span></label>
              <select v-model="form.outlet_id" class="ps-input" required @change="onOutletChange">
                <option value="">— Pilih Outlet —</option>
                <option v-for="o in outlets" :key="o.id" :value="o.id">
                  {{ o.name }}
                </option>
              </select>
              <div class="ps-field-hint">Outlet franchise yang diambil mitra</div>
            </div>

            <div class="ps-fg">
              <label>Pilih Paket <span class="req">*</span></label>
              <select v-model="form.package_id" class="ps-input" required :disabled="!form.outlet_id || loadingPkgs">
                <option value="">{{ loadingPkgs ? 'Memuat paket...' : '— Pilih Paket —' }}</option>
                <option v-for="pkg in packages" :key="pkg.id" :value="pkg.id">
                  {{ pkg.name }} — Rp {{ pkg.price?.toLocaleString('id-ID') }}
                </option>
              </select>
              <div class="ps-field-hint">Paket investasi yang dipilih mitra</div>
            </div>

            <div class="ps-modal-foot">
              <button type="button" @click="showModal=false" class="ps-btn-sec">Batal</button>
              <button type="submit" class="ps-btn-primary" :disabled="saving">
                {{ saving ? 'Menyimpan...' : 'Buat Partnership' }}
              </button>
            </div>
          </form>
        </div>
      </div>
    </Teleport>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted } from 'vue'
import { partnershipApi, userApi, outletApi, outletPackageApi, invoiceApi } from '../../services/api'
import { useToastStore } from '../../stores/toast'

const toast = useToastStore()
const partnerships = ref([])
const total = ref(0)
const search = ref('')
const showModal = ref(false)
const saving = ref(false)
const form = reactive({ leader_id: '', mitra_id: '', outlet_id: '', package_id: '' })

const leaders = ref([])
const mitras = ref([])
const outlets = ref([])
const packages = ref([])
const loadingPkgs = ref(false)

const avatarGradients = [
  'linear-gradient(135deg, #667eea, #764ba2)',
  'linear-gradient(135deg, #f093fb, #f5576c)',
  'linear-gradient(135deg, #4facfe, #00f2fe)',
  'linear-gradient(135deg, #43e97b, #38f9d7)',
  'linear-gradient(135deg, #fa709a, #fee140)',
  'linear-gradient(135deg, #a18cd1, #fbc2eb)',
]
function avatarBg(name) { const idx = (name||'?').charCodeAt(0) % avatarGradients.length; return { background: avatarGradients[idx] } }
function getInitial(name) { return name ? name.split(' ').map(n=>n[0]).join('').substring(0,2).toUpperCase() : '?' }
function formatDate(d) { return d ? new Date(d).toLocaleDateString('id-ID', { day:'numeric', month:'short', year:'numeric' }) : '-' }
function fc(v) { return new Intl.NumberFormat('id-ID',{style:'currency',currency:'IDR',minimumFractionDigits:0}).format(v||0) }

// Invoice data per partnership
const invoiceMap = ref({})

function paidAmount(p) {
  const invs = invoiceMap.value[p.id] || []
  return invs.filter(i => i.status === 'PAID').reduce((sum, i) => sum + (i.amount || 0), 0)
}
function payPercent(p) {
  const price = p.package?.price || 0
  if (!price) return 0
  return Math.min(100, Math.round(paidAmount(p) / price * 100))
}
function latestPayLink(p) {
  const invs = invoiceMap.value[p.id] || []
  const pending = invs.find(i => i.status === 'PENDING' && i.midtrans_redirect_url)
  return pending?.midtrans_redirect_url || ''
}

const filtered = computed(() => {
  const q = search.value.toLowerCase().trim()
  if (!q) return partnerships.value
  return partnerships.value.filter(p =>
    (p.mitra?.name||'').toLowerCase().includes(q) ||
    (p.mitra?.email||'').toLowerCase().includes(q) ||
    (p.leader?.name||'').toLowerCase().includes(q) ||
    (p.outlet?.name||'').toLowerCase().includes(q)
  )
})

onMounted(async () => {
  loadPartnerships()
  try {
    const [leaderRes, mitraRes, outletRes] = await Promise.all([
      userApi.list({ role: 'leader', limit: 100 }),
      userApi.list({ role: 'mitra', limit: 100 }),
      outletApi.list({ limit: 100 }),
    ])
    leaders.value = leaderRes.data.data || []
    mitras.value = mitraRes.data.data || []
    outlets.value = outletRes.data.data || []
  } catch {}
})

async function loadPartnerships() {
  try {
    const { data } = await partnershipApi.list({ page: 1, limit: 50 })
    partnerships.value = data.data || []
    total.value = data.meta?.total || partnerships.value.length
    // Load invoices for each partnership
    loadInvoices()
  } catch {
    toast.error('Gagal memuat partnership')
  }
}

async function loadInvoices() {
  const map = {}
  await Promise.all(
    partnerships.value.map(async (p) => {
      try {
        const { data } = await invoiceApi.getByPartnership(p.id)
        map[p.id] = data.data || []
      } catch {
        map[p.id] = []
      }
    })
  )
  invoiceMap.value = map
}

function openCreate() {
  Object.assign(form, { leader_id: '', mitra_id: '', outlet_id: '', package_id: '' })
  packages.value = []
  showModal.value = true
}

async function onOutletChange() {
  form.package_id = ''
  packages.value = []
  if (!form.outlet_id) return
  loadingPkgs.value = true
  try {
    const { data } = await outletPackageApi.listByOutlet(form.outlet_id)
    packages.value = data.data || []
  } catch {}
  finally { loadingPkgs.value = false }
}

async function create() {
  if (!form.leader_id || !form.mitra_id) {
    toast.error('Pilih Leader dan Mitra')
    return
  }
  if (!form.outlet_id || !form.package_id) {
    toast.error('Pilih Outlet dan Paket')
    return
  }
  saving.value = true
  try {
    const payload = {
      leader_id: form.leader_id,
      mitra_id: form.mitra_id,
      outlet_id: form.outlet_id,
      package_id: form.package_id,
    }
    await partnershipApi.create(payload)
    toast.success('Partnership berhasil dibuat')
    showModal.value = false
    loadPartnerships()
  } catch (e) {
    toast.error(e.response?.data?.error || 'Gagal membuat partnership')
  } finally {
    saving.value = false
  }
}
</script>

<style scoped>
/* ═══ HERO ═══ */
.ps-hero { background: linear-gradient(135deg, #0f0c29 0%, #302b63 50%, #24243e 100%); border-radius: 16px; padding: 32px 36px 24px; margin-bottom: 24px; }
.ps-hero-top { display: flex; justify-content: space-between; align-items: flex-start; gap: 16px; flex-wrap: wrap; margin-bottom: 20px; }
.ps-hero-title { font-size: 1.6rem; font-weight: 800; color: #fff; margin: 0 0 4px; }
.ps-hero-sub { font-size: .85rem; color: rgba(255,255,255,.5); margin: 0; }
.ps-stats { display: flex; gap: 28px; flex-wrap: wrap; padding-top: 16px; border-top: 1px solid rgba(255,255,255,.08); }
.ps-stat { display: flex; align-items: center; gap: 8px; }
.ps-stat-dot { width: 8px; height: 8px; border-radius: 50%; }
.ps-stat-label { font-size: .72rem; color: rgba(255,255,255,.4); text-transform: uppercase; letter-spacing: .05em; }
.ps-stat-val { font-size: .9rem; font-weight: 800; color: #fff; }

/* ═══ BUTTONS ═══ */
.ps-btn-primary { display: inline-flex; align-items: center; gap: 8px; padding: 11px 24px; font-size: .85rem; font-weight: 700; border-radius: 12px; border: none; cursor: pointer; background: linear-gradient(135deg, #6366f1, #8b5cf6); color: #fff; white-space: nowrap; }
.ps-btn-primary:disabled { opacity: .6; cursor: not-allowed; }
.ps-btn-sec { padding: 11px 22px; border-radius: 12px; font-size: .85rem; font-weight: 600; background: #f1f5f9; color: #475569; border: none; cursor: pointer; }

/* ═══ TOOLBAR ═══ */
.ps-toolbar { margin-bottom: 20px; }
.ps-search-wrap { position: relative; max-width: 400px; }
.ps-search-ico { position: absolute; left: 14px; top: 50%; transform: translateY(-50%); }
.ps-search { width: 100%; padding: 10px 14px 10px 40px; border: 1.5px solid #e2e8f0; border-radius: 12px; font-size: .85rem; background: #fff; color: #1e293b; outline: none; box-sizing: border-box; }
.ps-search:focus { border-color: #6366f1; }

/* ═══ TABLE ═══ */
.ps-table-wrap { background: #fff; border-radius: 16px; border: 1px solid #e8ecf1; overflow: hidden; }
.ps-table { width: 100%; border-collapse: collapse; }
.ps-table thead { background: #f8fafc; }
.ps-table th { padding: 14px 20px; font-size: .72rem; font-weight: 700; color: #64748b; text-transform: uppercase; letter-spacing: .05em; text-align: left; border-bottom: 1px solid #e8ecf1; }
.ps-table td { padding: 16px 20px; font-size: .85rem; color: #1e293b; border-bottom: 1px solid #f1f5f9; vertical-align: middle; }
.ps-table tr:last-child td { border-bottom: none; }
.ps-table tbody tr { transition: background .15s; }
.ps-table tbody tr:hover { background: #fafbfc; }

.ps-mitra-cell { display: flex; align-items: center; gap: 12px; }
.ps-mitra-avatar { width: 38px; height: 38px; border-radius: 10px; display: flex; align-items: center; justify-content: center; font-size: .78rem; font-weight: 800; color: #fff; flex-shrink: 0; }
.ps-mitra-name { font-weight: 700; font-size: .88rem; color: #0f172a; }
.ps-mitra-email { font-size: .72rem; color: #94a3b8; }
.ps-lead-cell { font-weight: 600; font-size: .85rem; }
.ps-lead-phone { font-size: .72rem; color: #94a3b8; }

.ps-progress-cell { display: flex; align-items: center; gap: 10px; }
.ps-progress-bar { width: 80px; height: 6px; background: #e2e8f0; border-radius: 3px; overflow: hidden; }
.ps-progress-fill { height: 100%; background: linear-gradient(90deg, #6366f1, #8b5cf6); border-radius: 3px; transition: width .3s; }
.ps-progress-text { font-size: .78rem; font-weight: 700; color: #334155; }

.ps-badge { font-size: .68rem; font-weight: 700; padding: 4px 10px; border-radius: 6px; text-transform: uppercase; letter-spacing: .03em; }
.st-PENDING { background: #fef3c7; color: #d97706; }
.st-DP_VERIFIED { background: #e0f2fe; color: #0284c7; }
.st-AGREEMENT_SIGNED { background: #ede9fe; color: #7c3aed; }
.st-DEVELOPMENT { background: #fce7f3; color: #db2777; }
.st-RUNNING { background: #dcfce7; color: #16a34a; }
.st-COMPLETED { background: #d1fae5; color: #059669; }

.ps-date-cell { font-size: .82rem; color: #64748b; }
.ps-detail-btn { font-size: .82rem; font-weight: 600; color: #6366f1; text-decoration: none; padding: 6px 14px; border-radius: 8px; border: 1px solid #e0e7ff; background: #eef2ff; }
.ps-detail-btn:hover { background: #e0e7ff; }

/* ═══ PAYMENT CELL ═══ */
.ps-pay-cell { min-width: 140px; }
.ps-pay-row { display: flex; align-items: baseline; gap: 2px; margin-bottom: 4px; }
.ps-pay-amount { font-size: .75rem; font-weight: 700; color: #0f172a; }
.ps-pay-sep { font-size: .7rem; color: #cbd5e1; }
.ps-pay-total { font-size: .68rem; color: #94a3b8; }
.ps-pay-footer { display: flex; align-items: center; justify-content: space-between; margin-top: 4px; }
.ps-pay-pct { font-size: .7rem; font-weight: 700; color: #475569; }
.ps-pay-link { display: inline-flex; align-items: center; gap: 3px; font-size: .68rem; font-weight: 600; color: #6366f1; text-decoration: none; padding: 2px 8px; background: #eef2ff; border-radius: 6px; transition: all .15s; }
.ps-pay-link:hover { background: #e0e7ff; }

.ps-empty-row { text-align: center; padding: 48px 20px !important; color: #94a3b8; font-size: .85rem; }

/* ═══ MODAL ═══ */
.ps-overlay { position: fixed; inset: 0; background: rgba(0,0,0,.5); display: flex; align-items: center; justify-content: center; z-index: 1000; backdrop-filter: blur(4px); }
.ps-modal { background: #fff; border-radius: 18px; width: 100%; max-width: 520px; box-shadow: 0 24px 80px rgba(0,0,0,.2); }
.ps-modal-head { display: flex; align-items: center; justify-content: space-between; padding: 22px 28px; border-bottom: 1px solid #f1f5f9; }
.ps-modal-title-group { display: flex; align-items: center; gap: 12px; }
.ps-modal-icon { width: 40px; height: 40px; border-radius: 12px; background: linear-gradient(135deg, #eef2ff, #e0e7ff); display: flex; align-items: center; justify-content: center; color: #6366f1; }
.ps-modal-head h3 { font-size: 1.1rem; font-weight: 700; margin: 0; color: #0f172a; }
.ps-modal-x { width: 34px; height: 34px; border-radius: 10px; display: flex; align-items: center; justify-content: center; border: none; background: transparent; font-size: 1.4rem; color: #94a3b8; cursor: pointer; }
.ps-modal-x:hover { background: #f1f5f9; color: #0f172a; }
.ps-modal-body { padding: 24px 28px; }
.ps-modal-foot { display: flex; justify-content: flex-end; gap: 10px; padding-top: 16px; }

.ps-fg { margin-bottom: 18px; }
.ps-fg label { display: block; font-size: .82rem; font-weight: 600; margin-bottom: 6px; color: #334155; }
.req { color: #ef4444; }
.ps-input { width: 100%; padding: 11px 14px; border: 1.5px solid #e2e8f0; border-radius: 12px; font-size: .85rem; background: #fafbfc; color: #1e293b; outline: none; box-sizing: border-box; font-family: inherit; appearance: none; -webkit-appearance: none; }
.ps-input:focus { border-color: #6366f1; background: #fff; }
select.ps-input { background-image: url("data:image/svg+xml,%3csvg xmlns='http://www.w3.org/2000/svg' fill='none' viewBox='0 0 20 20'%3e%3cpath stroke='%236b7280' stroke-linecap='round' stroke-linejoin='round' stroke-width='1.5' d='M6 8l4 4 4-4'/%3e%3c/svg%3e"); background-position: right 12px center; background-repeat: no-repeat; background-size: 16px; padding-right: 36px; cursor: pointer; }
.ps-field-hint { font-size: .72rem; color: #94a3b8; margin-top: 4px; }

@media (max-width: 768px) {
  .ps-hero { padding: 24px 20px 18px; }
  .ps-hero-top { flex-direction: column; align-items: flex-start; }
  .ps-table-wrap { overflow-x: auto; }
}
</style>
