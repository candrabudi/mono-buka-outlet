<template>
  <div class="affm-page">
    <!-- Hero -->
    <div class="affm-hero">
      <div class="affm-hero-top">
        <div>
          <h1 class="affm-hero-title"><i class="ri-team-line"></i> Affiliator Management</h1>
          <p class="affm-hero-sub">Kelola komisi, saldo, dan penarikan affiliator</p>
        </div>
      </div>
      <div class="affm-stats-bar">
        <div class="affm-stat"><span class="stat-dot dot-aff"></span><span class="stat-label">Affiliator</span><span class="stat-val">{{ affiliatorList.length }}</span></div>
        <div class="affm-stat"><span class="stat-dot dot-pending"></span><span class="stat-label">Penarikan Menunggu</span><span class="stat-val">{{ pendingWithdrawals.length }}</span></div>
      </div>
    </div>

    <!-- Tabs -->
    <div class="affm-tab-bar">
      <button v-for="t in tabList" :key="t.key" class="affm-tab" :class="{ active: activeTab === t.key }" @click="switchTab(t.key)">
        <i :class="t.icon"></i>
        <span>{{ t.label }}</span>
        <span v-if="t.badge" class="affm-tab-badge">{{ t.badge }}</span>
      </button>
    </div>

    <!-- TAB: Give Commission -->
    <div v-if="activeTab === 'commission'" class="affm-card">
      <div class="affm-card-head">
        <h3><i class="ri-hand-coin-line"></i> Berikan Komisi / Bonus</h3>
        <p class="affm-card-desc">Pilih affiliator, tentukan partnership terkait, dan masukkan jumlah komisi</p>
      </div>
      <div class="affm-card-body">
        <form @submit.prevent="submitCommission" class="affm-form">
          <div class="affm-form-grid">
            <div class="affm-fg">
              <label>Affiliator <span class="req">*</span></label>
              <SearchSelect
                v-model="commForm.affiliator_id"
                :options="affiliatorOptions"
                placeholder="— Pilih Affiliator —"
                search-placeholder="Cari affiliator..."
                :allow-empty="false"
                @update:model-value="onAffiliatorChange"
              />
            </div>
            <div class="affm-fg">
              <label>Partnership (Mitra — Outlet)</label>
              <SearchSelect
                v-model="commForm.partnership_id"
                :options="partnershipOptions"
                placeholder="— Pilih Partnership (opsional) —"
                search-placeholder="Cari mitra / outlet..."
                empty-label="Tanpa Partnership"
              />
            </div>
            <div class="affm-fg">
              <label>Tipe <span class="req">*</span></label>
              <SearchSelect
                v-model="commForm.type"
                :options="typeOptions"
                placeholder="— Pilih Tipe —"
                :allow-empty="false"
              />
            </div>
            <div class="affm-fg">
              <label>Jumlah (Rp) <span class="req">*</span></label>
              <input v-model.number="commForm.amount" type="number" step="any" required placeholder="500000" class="affm-input" />
            </div>
            <div class="affm-fg affm-fg-full">
              <label>Deskripsi</label>
              <input v-model="commForm.description" type="text" placeholder="Contoh: Komisi referral mitra baru" class="affm-input" />
            </div>
          </div>

          <!-- Preview Card -->
          <div v-if="selectedAffiliator" class="affm-preview">
            <div class="affm-preview-left">
              <div class="affm-preview-avatar" :style="avatarBg(selectedAffiliator.name)">{{ getInitial(selectedAffiliator.name) }}</div>
              <div class="affm-preview-info">
                <div class="affm-preview-name">{{ selectedAffiliator.name }}</div>
                <div class="affm-preview-email">{{ selectedAffiliator.email }}</div>
                <div class="affm-preview-code" v-if="selectedAffiliator.referral_code">
                  <i class="ri-gift-line"></i> {{ selectedAffiliator.referral_code }}
                </div>
              </div>
            </div>
            <div class="affm-preview-mid" v-if="selectedPartnership">
              <div class="affm-preview-partnership">
                <span class="affm-preview-tag tag-mitra"><i class="ri-user-star-line"></i> {{ selectedPartnership.mitra?.name || '-' }}</span>
                <i class="ri-arrow-right-s-line affm-preview-arrow"></i>
                <span class="affm-preview-tag tag-outlet"><i class="ri-store-2-line"></i> {{ selectedPartnership.outlet?.name || '-' }}</span>
              </div>
              <div class="affm-preview-status">
                <span class="affm-mini-badge" :class="'st-'+selectedPartnership.status">{{ statusLabel(selectedPartnership.status) }}</span>
              </div>
            </div>
            <div class="affm-preview-right" v-if="commForm.amount > 0">
              <span class="affm-preview-amt-label">Akan diberikan</span>
              <span class="affm-preview-amt-value">Rp {{ formatNumber(commForm.amount) }}</span>
            </div>
          </div>

          <div class="affm-form-actions">
            <button type="submit" class="affm-btn-primary" :disabled="commLoading || !commForm.affiliator_id || !commForm.amount">
              <i :class="commLoading ? 'ri-loader-4-line ri-spin' : 'ri-send-plane-line'"></i>
              {{ commLoading ? 'Memproses...' : 'Berikan Komisi' }}
            </button>
          </div>
        </form>
      </div>
    </div>

    <!-- TAB: Withdrawals -->
    <div v-if="activeTab === 'withdrawals'" class="affm-card">
      <div class="affm-card-head">
        <h3><i class="ri-bank-line"></i> Permintaan Penarikan</h3>
        <SearchSelect
          v-model="wdFilter"
          :options="wdStatusOptions"
          placeholder="Semua Status"
          empty-label="Semua Status"
          @update:model-value="loadWithdrawals"
        />
      </div>
      <div class="affm-card-body affm-card-body-flush">
        <div v-if="wdLoading" class="affm-loading">
          <div v-for="n in 4" :key="n" class="affm-skel shimmer"></div>
        </div>
        <div v-else-if="withdrawalList.length" class="affm-table-wrap">
          <table class="affm-table">
            <thead>
              <tr>
                <th>Affiliator</th>
                <th>Jumlah</th>
                <th>Bank</th>
                <th>No. Rekening</th>
                <th>Status</th>
                <th>Tanggal</th>
                <th>Aksi</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="w in withdrawalList" :key="w.id">
                <td>
                  <div class="affm-cell-user">
                    <div class="affm-cell-avatar" :style="avatarBg(w.affiliator?.name || '')">{{ getInitial(w.affiliator?.name) }}</div>
                    <div>
                      <div class="affm-cell-name">{{ w.affiliator?.name || '-' }}</div>
                      <div class="affm-cell-email">{{ w.affiliator?.email || '' }}</div>
                    </div>
                  </div>
                </td>
                <td class="affm-cell-amount">Rp {{ formatNumber(w.amount) }}</td>
                <td>{{ w.bank_name }}</td>
                <td>
                  <div class="affm-cell-acc">{{ w.account_number }}</div>
                  <div class="affm-cell-acc-name">a.n. {{ w.account_holder }}</div>
                </td>
                <td><span class="affm-badge" :class="wdStatusClass(w.status)">{{ wdStatusLabel(w.status) }}</span></td>
                <td class="affm-cell-date">{{ formatDate(w.created_at) }}</td>
                <td>
                  <div class="affm-cell-actions" v-if="w.status === 'PENDING'">
                    <button @click="openWdModal(w, 'APPROVED')" class="affm-act-btn act-approve" title="Setujui"><i class="ri-check-line"></i></button>
                    <button @click="openWdModal(w, 'REJECTED')" class="affm-act-btn act-reject" title="Tolak"><i class="ri-close-line"></i></button>
                  </div>
                  <div class="affm-cell-actions" v-else-if="w.status === 'APPROVED'">
                    <button @click="openWdModal(w, 'TRANSFERRED')" class="affm-act-btn act-transfer" title="Tandai Ditransfer"><i class="ri-arrow-right-up-line"></i></button>
                  </div>
                  <span v-else class="affm-cell-done">—</span>
                </td>
              </tr>
            </tbody>
          </table>
        </div>
        <div v-else class="affm-empty">
          <i class="ri-inbox-line"></i>
          <p>Belum ada permintaan penarikan</p>
        </div>
      </div>
    </div>

    <!-- TAB: Commission History -->
    <div v-if="activeTab === 'history'" class="affm-card">
      <div class="affm-card-head">
        <h3><i class="ri-history-line"></i> Riwayat Komisi</h3>
        <SearchSelect
          v-model="historyAffId"
          :options="affiliatorOptions"
          placeholder="— Pilih Affiliator —"
          search-placeholder="Cari affiliator..."
          empty-label="— Pilih Affiliator —"
          @update:model-value="loadHistory"
        />
      </div>

      <!-- Balance Card -->
      <div v-if="historyAffId && histBalance !== null" class="affm-balance-strip">
        <div class="affm-balance-item">
          <i class="ri-wallet-3-line"></i>
          <span class="affm-balance-label">Saldo Tersedia</span>
          <span class="affm-balance-value">Rp {{ formatNumber(histBalance) }}</span>
        </div>
      </div>

      <div class="affm-card-body affm-card-body-flush">
        <div v-if="histLoading" class="affm-loading">
          <div v-for="n in 5" :key="n" class="affm-skel shimmer"></div>
        </div>
        <div v-else-if="historyList.length" class="affm-table-wrap">
          <table class="affm-table">
            <thead>
              <tr>
                <th>Tanggal</th>
                <th>Tipe</th>
                <th>Mitra — Outlet</th>
                <th>Jumlah</th>
                <th>Deskripsi</th>
                <th>Diberikan Oleh</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="c in historyList" :key="c.id">
                <td class="affm-cell-date">{{ formatDateTime(c.created_at) }}</td>
                <td><span class="affm-badge" :class="typeClass(c.type)">{{ typeLbl(c.type) }}</span></td>
                <td>
                  <div v-if="c.partnership" class="affm-cell-partnership">
                    <span class="affm-cell-mitra"><i class="ri-user-star-line"></i> {{ c.partnership.mitra?.name || '-' }}</span>
                    <span class="affm-cell-outlet"><i class="ri-store-2-line"></i> {{ c.partnership.outlet?.name || '-' }}</span>
                  </div>
                  <span v-else class="affm-cell-na">—</span>
                </td>
                <td class="affm-cell-amount affm-cell-green">+Rp {{ formatNumber(c.amount) }}</td>
                <td>{{ c.description || '-' }}</td>
                <td>{{ c.given_by_user?.name || '-' }}</td>
              </tr>
            </tbody>
          </table>
        </div>
        <div v-else-if="historyAffId" class="affm-empty">
          <i class="ri-inbox-line"></i>
          <p>Belum ada riwayat komisi</p>
        </div>
        <div v-else class="affm-empty">
          <i class="ri-user-search-line"></i>
          <p>Pilih affiliator untuk melihat riwayat komisi</p>
        </div>
      </div>
    </div>

    <!-- Withdrawal Process Modal -->
    <Teleport to="body">
      <div v-if="wdModal.show" class="affm-overlay" @click.self="wdModal.show = false">
        <div class="affm-modal">
          <div class="affm-modal-head">
            <h3>{{ wdModal.status === 'APPROVED' ? 'Setujui Penarikan' : wdModal.status === 'REJECTED' ? 'Tolak Penarikan' : 'Konfirmasi Transfer' }}</h3>
            <button @click="wdModal.show = false" class="affm-modal-x">&times;</button>
          </div>
          <div class="affm-modal-body">
            <div class="affm-modal-info"><span>Affiliator:</span><strong>{{ wdModal.withdrawal?.affiliator?.name }}</strong></div>
            <div class="affm-modal-info"><span>Jumlah:</span><strong class="affm-modal-amt">Rp {{ formatNumber(wdModal.withdrawal?.amount) }}</strong></div>
            <div class="affm-modal-info"><span>Bank:</span><strong>{{ wdModal.withdrawal?.bank_name }} - {{ wdModal.withdrawal?.account_number }}</strong></div>
            <div class="affm-fg" style="margin-top:16px">
              <label>Catatan Admin</label>
              <textarea v-model="wdModal.notes" rows="3" class="affm-input affm-textarea" :placeholder="wdModal.status === 'REJECTED' ? 'Alasan penolakan...' : 'Catatan opsional...'"></textarea>
            </div>
          </div>
          <div class="affm-modal-foot">
            <button @click="wdModal.show = false" class="affm-btn-sec">Batal</button>
            <button @click="processWithdrawal" class="affm-btn-primary" :class="{ 'btn-danger': wdModal.status === 'REJECTED' }" :disabled="wdModal.loading">
              <i :class="wdModal.loading ? 'ri-loader-4-line ri-spin' : 'ri-check-double-line'"></i>
              {{ wdModal.status === 'APPROVED' ? 'Setujui' : wdModal.status === 'REJECTED' ? 'Tolak' : 'Tandai Sudah Transfer' }}
            </button>
          </div>
        </div>
      </div>
    </Teleport>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { userApi, affiliatorMgmtApi, partnershipApi } from '../../services/api'
import { useToastStore } from '../../stores/toast'
import SearchSelect from '../../components/SearchSelect.vue'

const toast = useToastStore()
const activeTab = ref('commission')

// ── Data ──
const affiliatorList = ref([])
const partnershipsList = ref([])

async function loadAffiliators() {
  try {
    const { data } = await userApi.list({ role: 'affiliator', limit: 200 })
    affiliatorList.value = data.data || []
  } catch { toast.error('Gagal memuat data affiliator') }
}

async function loadPartnerships() {
  try {
    const { data } = await partnershipApi.list({ limit: 500 })
    partnershipsList.value = data.data || []
  } catch { partnershipsList.value = [] }
}

// ── Computed Options ──
const affiliatorOptions = computed(() => affiliatorList.value.map(a => ({ value: a.id, label: a.name, sub: a.email })))
const typeOptions = [
  { value: 'COMMISSION', label: 'Komisi' },
  { value: 'BONUS', label: 'Bonus' },
  { value: 'ADJUSTMENT', label: 'Penyesuaian' },
]

const partnershipOptions = computed(() => {
  let list = partnershipsList.value
  // Filter by selected affiliator if set
  if (commForm.value.affiliator_id) {
    list = list.filter(p => p.affiliator_id === commForm.value.affiliator_id)
  }
  return list.map(p => ({
    value: p.id,
    label: `${p.mitra?.name || 'Mitra'} — ${p.outlet?.name || 'Outlet'}`,
    sub: statusLabel(p.status),
  }))
})

const tabList = computed(() => [
  { key: 'commission', label: 'Berikan Komisi', icon: 'ri-hand-coin-line', badge: 0 },
  { key: 'withdrawals', label: 'Penarikan', icon: 'ri-bank-line', badge: pendingWithdrawals.value.length || 0 },
  { key: 'history', label: 'Riwayat Komisi', icon: 'ri-history-line', badge: 0 },
])

function switchTab(key) {
  activeTab.value = key
  if (key === 'withdrawals') loadWithdrawals()
}

// ── Commission Form ──
const commForm = ref({ affiliator_id: '', partnership_id: '', type: 'COMMISSION', amount: null, description: '' })
const commLoading = ref(false)

const selectedAffiliator = computed(() => {
  if (!commForm.value.affiliator_id) return null
  return affiliatorList.value.find(a => a.id === commForm.value.affiliator_id)
})

const selectedPartnership = computed(() => {
  if (!commForm.value.partnership_id) return null
  return partnershipsList.value.find(p => p.id === commForm.value.partnership_id)
})

function onAffiliatorChange() {
  // Reset partnership when affiliator changes
  commForm.value.partnership_id = ''
}

async function submitCommission() {
  if (!commForm.value.affiliator_id || !commForm.value.amount) return
  commLoading.value = true
  try {
    const payload = {
      affiliator_id: commForm.value.affiliator_id,
      type: commForm.value.type,
      amount: commForm.value.amount,
      description: commForm.value.description,
    }
    if (commForm.value.partnership_id) {
      payload.partnership_id = commForm.value.partnership_id
    }
    await affiliatorMgmtApi.giveCommission(payload)
    toast.success('Komisi berhasil diberikan!')
    commForm.value = { affiliator_id: commForm.value.affiliator_id, partnership_id: '', type: 'COMMISSION', amount: null, description: '' }
  } catch (err) {
    toast.error(err.response?.data?.error || 'Gagal memberikan komisi')
  } finally { commLoading.value = false }
}

// ── Withdrawals ──
const withdrawalList = ref([])
const wdLoading = ref(false)
const wdFilter = ref('')
const wdStatusOptions = [
  { value: 'PENDING', label: 'Menunggu' },
  { value: 'APPROVED', label: 'Disetujui' },
  { value: 'TRANSFERRED', label: 'Ditransfer' },
  { value: 'REJECTED', label: 'Ditolak' },
]
const pendingWithdrawals = computed(() => withdrawalList.value.filter(w => w.status === 'PENDING'))

async function loadWithdrawals() {
  wdLoading.value = true
  try {
    const { data } = await affiliatorMgmtApi.getWithdrawals({ status: wdFilter.value, limit: 100 })
    withdrawalList.value = data.data || []
  } catch { toast.error('Gagal memuat data penarikan') }
  finally { wdLoading.value = false }
}

// ── Withdrawal Modal ──
const wdModal = ref({ show: false, withdrawal: null, status: '', notes: '', loading: false })

function openWdModal(w, status) {
  wdModal.value = { show: true, withdrawal: w, status, notes: '', loading: false }
}

async function processWithdrawal() {
  wdModal.value.loading = true
  try {
    await affiliatorMgmtApi.processWithdrawal(wdModal.value.withdrawal.id, {
      status: wdModal.value.status,
      admin_notes: wdModal.value.notes,
    })
    toast.success('Status penarikan berhasil diperbarui')
    wdModal.value.show = false
    await loadWithdrawals()
  } catch (err) {
    toast.error(err.response?.data?.error || 'Gagal memproses penarikan')
  } finally { wdModal.value.loading = false }
}

// ── Commission History ──
const historyAffId = ref('')
const historyList = ref([])
const histBalance = ref(null)
const histLoading = ref(false)

async function loadHistory() {
  if (!historyAffId.value) { historyList.value = []; histBalance.value = null; return }
  histLoading.value = true
  try {
    const { data } = await affiliatorMgmtApi.getCommissions(historyAffId.value, { limit: 100 })
    historyList.value = data.data || []
    histBalance.value = data.balance ?? null
  } catch { toast.error('Gagal memuat riwayat komisi') }
  finally { histLoading.value = false }
}

// ── Init ──
onMounted(async () => {
  await Promise.all([loadAffiliators(), loadWithdrawals(), loadPartnerships()])
})

// ── Helpers ──
const avatarGradients = [
  'linear-gradient(135deg, #667eea, #764ba2)',
  'linear-gradient(135deg, #f093fb, #f5576c)',
  'linear-gradient(135deg, #4facfe, #00f2fe)',
  'linear-gradient(135deg, #43e97b, #38f9d7)',
  'linear-gradient(135deg, #fa709a, #fee140)',
]
function avatarBg(name) { return { background: avatarGradients[(name||'?').charCodeAt(0) % avatarGradients.length] } }
function getInitial(name) { return name ? name.split(' ').map(n=>n[0]).join('').substring(0,2).toUpperCase() : '?' }
function formatNumber(n) { return n ? Number(n).toLocaleString('id-ID') : '0' }
function formatDate(d) { return d ? new Date(d).toLocaleDateString('id-ID', { day: 'numeric', month: 'short', year: 'numeric' }) : '-' }
function formatDateTime(d) { return d ? new Date(d).toLocaleDateString('id-ID', { day: 'numeric', month: 'short', year: 'numeric', hour: '2-digit', minute: '2-digit' }) : '-' }

function statusLabel(s) {
  return { PENDING:'Menunggu', DP_VERIFIED:'DP Terverifikasi', AGREEMENT_SIGNED:'Perjanjian', DEVELOPMENT:'Pembangunan', RUNNING:'Berjalan', COMPLETED:'Selesai' }[s] || s
}
function wdStatusLabel(s) { return { PENDING: 'Menunggu', APPROVED: 'Disetujui', TRANSFERRED: 'Ditransfer', REJECTED: 'Ditolak' }[s] || s }
function wdStatusClass(s) { return { PENDING: 'badge-pending', APPROVED: 'badge-approved', TRANSFERRED: 'badge-transferred', REJECTED: 'badge-rejected' }[s] || '' }
function typeLbl(t) { return { COMMISSION: 'Komisi', BONUS: 'Bonus', ADJUSTMENT: 'Penyesuaian' }[t] || t }
function typeClass(t) { return { COMMISSION: 'badge-commission', BONUS: 'badge-bonus', ADJUSTMENT: 'badge-adjustment' }[t] || '' }
</script>

<style scoped>
/* ═══ HERO ═══ */
.affm-hero { background: linear-gradient(135deg, #0f0c29 0%, #302b63 50%, #24243e 100%); border-radius: 16px; padding: 32px 36px 24px; margin-bottom: 24px; }
.affm-hero-top { display: flex; justify-content: space-between; align-items: flex-start; gap: 16px; flex-wrap: wrap; margin-bottom: 20px; }
.affm-hero-title { font-size: 1.6rem; font-weight: 800; color: #fff; margin: 0 0 4px; display: flex; align-items: center; gap: 10px; }
.affm-hero-sub { font-size: .85rem; color: rgba(255,255,255,.5); margin: 0; }
.affm-stats-bar { display: flex; gap: 28px; flex-wrap: wrap; padding-top: 16px; border-top: 1px solid rgba(255,255,255,.08); }
.affm-stat { display: flex; align-items: center; gap: 8px; }
.stat-dot { width: 8px; height: 8px; border-radius: 50%; }
.dot-aff { background: #818cf8; box-shadow: 0 0 8px rgba(129,140,248,.5); }
.dot-pending { background: #f59e0b; box-shadow: 0 0 8px rgba(245,158,11,.5); }
.stat-label { font-size: .72rem; color: rgba(255,255,255,.4); text-transform: uppercase; letter-spacing: .05em; }
.stat-val { font-size: .9rem; font-weight: 800; color: #fff; }

/* ═══ TABS ═══ */
.affm-tab-bar { display: flex; gap: 4px; background: #fff; border-radius: 14px; padding: 5px; margin-bottom: 20px; border: 1px solid #e8ecf1; }
.affm-tab { flex: 1; display: flex; align-items: center; justify-content: center; gap: 8px; padding: 11px 14px; border-radius: 10px; border: none; background: transparent; color: #64748b; font-size: .82rem; font-weight: 500; cursor: pointer; transition: all .2s; font-family: inherit; position: relative; }
.affm-tab:hover { color: #334155; background: #f8fafc; }
.affm-tab.active { background: linear-gradient(135deg, #6366f1, #8b5cf6); color: #fff; font-weight: 600; box-shadow: 0 4px 14px rgba(99,102,241,.25); }
.affm-tab i { font-size: 1rem; }
.affm-tab-badge { position: absolute; top: 4px; right: 8px; min-width: 18px; height: 18px; display: flex; align-items: center; justify-content: center; font-size: .65rem; font-weight: 700; background: #ef4444; color: #fff; border-radius: 10px; padding: 0 5px; }
.affm-tab.active .affm-tab-badge { background: rgba(255,255,255,.3); }

/* ═══ CARD ═══ */
.affm-card { background: #fff; border-radius: 16px; border: 1px solid #e8ecf1; box-shadow: 0 1px 4px rgba(0,0,0,.04); }
.affm-card-head { display: flex; align-items: center; justify-content: space-between; padding: 22px 28px; border-bottom: 1px solid #f1f5f9; flex-wrap: wrap; gap: 12px; }
.affm-card-head h3 { font-size: 1.05rem; font-weight: 700; color: #0f172a; margin: 0; display: flex; align-items: center; gap: 8px; }
.affm-card-head h3 i { color: #6366f1; font-size: 1.1rem; }
.affm-card-desc { color: #94a3b8; font-size: .82rem; margin: 4px 0 0; width: 100%; }
.affm-card-body { padding: 24px 28px; }
.affm-card-body-flush { padding: 0; }

/* ═══ FORM ═══ */
.affm-form-grid { display: grid; grid-template-columns: 1fr 1fr; gap: 16px; }
.affm-fg { display: flex; flex-direction: column; gap: 6px; }
.affm-fg label { font-size: .78rem; font-weight: 600; color: #475569; }
.affm-fg-full { grid-column: 1 / -1; }
.req { color: #ef4444; }
.affm-input { padding: 10px 14px; border: 1.5px solid #e2e8f0; border-radius: 11px; font-size: .85rem; color: #1e293b; background: #fafbfc; outline: none; transition: all .2s; font-family: inherit; box-sizing: border-box; width: 100%; }
.affm-input:focus { border-color: #6366f1; box-shadow: 0 0 0 3px rgba(99,102,241,.1); background: #fff; }
.affm-textarea { resize: vertical; min-height: 80px; }

/* ═══ PREVIEW ═══ */
.affm-preview { display: flex; align-items: center; gap: 20px; margin-top: 20px; padding: 20px 24px; background: linear-gradient(135deg, #eef2ff 0%, #faf5ff 100%); border: 1px solid #e0e7ff; border-radius: 14px; flex-wrap: wrap; }
.affm-preview-left { display: flex; align-items: center; gap: 14px; }
.affm-preview-avatar { width: 44px; height: 44px; border-radius: 12px; display: flex; align-items: center; justify-content: center; font-size: .8rem; font-weight: 700; color: #fff; flex-shrink: 0; }
.affm-preview-info { min-width: 0; }
.affm-preview-name { font-weight: 700; color: #1e293b; font-size: .9rem; }
.affm-preview-email { font-size: .78rem; color: #94a3b8; }
.affm-preview-code { display: inline-flex; align-items: center; gap: 4px; margin-top: 4px; font-size: .72rem; color: #7c3aed; font-weight: 600; background: rgba(124,58,237,.08); padding: 3px 10px; border-radius: 6px; }
.affm-preview-mid { flex: 1; min-width: 200px; display: flex; flex-direction: column; gap: 6px; padding: 0 16px; border-left: 2px solid #e0e7ff; }
.affm-preview-partnership { display: flex; align-items: center; gap: 8px; flex-wrap: wrap; }
.affm-preview-tag { display: inline-flex; align-items: center; gap: 4px; font-size: .78rem; font-weight: 600; padding: 4px 10px; border-radius: 8px; }
.tag-mitra { background: rgba(99,102,241,.1); color: #4338ca; }
.tag-outlet { background: rgba(16,185,129,.1); color: #047857; }
.affm-preview-arrow { color: #94a3b8; }
.affm-preview-status { }
.affm-mini-badge { font-size: .65rem; font-weight: 700; padding: 2px 8px; border-radius: 6px; text-transform: uppercase; letter-spacing: .03em; }
.st-PENDING { background: #fef3c7; color: #92400e; }
.st-DP_VERIFIED { background: #dbeafe; color: #1e40af; }
.st-AGREEMENT_SIGNED { background: #f3e8ff; color: #6b21a8; }
.st-DEVELOPMENT { background: #eef2ff; color: #4338ca; }
.st-RUNNING { background: #dcfce7; color: #166534; }
.st-COMPLETED { background: #d1fae5; color: #065f46; }
.affm-preview-right { text-align: right; margin-left: auto; }
.affm-preview-amt-label { display: block; font-size: .72rem; color: #94a3b8; }
.affm-preview-amt-value { font-size: 1.2rem; font-weight: 800; color: #059669; }

/* ═══ BUTTONS ═══ */
.affm-form-actions { display: flex; justify-content: flex-end; margin-top: 24px; }
.affm-btn-primary { display: inline-flex; align-items: center; gap: 8px; padding: 11px 24px; border-radius: 12px; font-size: .85rem; font-weight: 600; border: none; cursor: pointer; background: linear-gradient(135deg, #6366f1, #8b5cf6); color: #fff; transition: all .2s; box-shadow: 0 4px 14px rgba(99,102,241,.2); font-family: inherit; }
.affm-btn-primary:not(:disabled):hover { box-shadow: 0 6px 20px rgba(99,102,241,.3); transform: translateY(-1px); }
.affm-btn-primary:disabled { opacity: .5; cursor: not-allowed; transform: none; }
.affm-btn-primary.btn-danger { background: linear-gradient(135deg, #ef4444, #dc2626); box-shadow: 0 4px 14px rgba(239,68,68,.2); }
.affm-btn-sec { display: inline-flex; align-items: center; gap: 8px; padding: 11px 22px; border-radius: 12px; font-size: .85rem; font-weight: 600; background: #f1f5f9; color: #475569; border: none; cursor: pointer; font-family: inherit; }
.affm-btn-sec:hover { background: #e2e8f0; }

/* ═══ TABLE ═══ */
.affm-table-wrap { overflow-x: auto; }
.affm-table { width: 100%; border-collapse: collapse; }
.affm-table thead th { text-align: left; padding: 12px 16px; font-size: .72rem; font-weight: 700; color: #94a3b8; text-transform: uppercase; letter-spacing: .05em; border-bottom: 2px solid #f1f5f9; }
.affm-table tbody td { padding: 14px 16px; border-bottom: 1px solid #f8fafc; font-size: .85rem; color: #334155; }
.affm-table tbody tr:hover { background: #fafbfc; }
.affm-cell-user { display: flex; align-items: center; gap: 10px; }
.affm-cell-avatar { width: 34px; height: 34px; border-radius: 10px; display: flex; align-items: center; justify-content: center; font-size: .7rem; font-weight: 700; color: #fff; flex-shrink: 0; }
.affm-cell-name { font-weight: 600; color: #0f172a; font-size: .85rem; }
.affm-cell-email { font-size: .72rem; color: #94a3b8; }
.affm-cell-amount { font-weight: 700; white-space: nowrap; }
.affm-cell-green { color: #059669; }
.affm-cell-date { white-space: nowrap; color: #94a3b8; font-size: .8rem; }
.affm-cell-acc { font-family: 'SF Mono', monospace; font-size: .82rem; }
.affm-cell-acc-name { font-size: .72rem; color: #94a3b8; }
.affm-cell-done { color: #cbd5e1; }
.affm-cell-partnership { display: flex; flex-direction: column; gap: 2px; }
.affm-cell-mitra { font-size: .8rem; font-weight: 600; color: #4338ca; display: flex; align-items: center; gap: 4px; }
.affm-cell-mitra i { font-size: .75rem; }
.affm-cell-outlet { font-size: .75rem; color: #059669; display: flex; align-items: center; gap: 4px; }
.affm-cell-outlet i { font-size: .7rem; }
.affm-cell-na { color: #cbd5e1; font-size: .82rem; }
.affm-cell-actions { display: flex; gap: 6px; }
.affm-act-btn { width: 32px; height: 32px; border-radius: 8px; border: none; display: flex; align-items: center; justify-content: center; cursor: pointer; transition: all .15s; }
.act-approve { background: rgba(34,197,94,.1); color: #16a34a; }
.act-approve:hover { background: rgba(34,197,94,.2); }
.act-reject { background: rgba(239,68,68,.1); color: #dc2626; }
.act-reject:hover { background: rgba(239,68,68,.2); }
.act-transfer { background: rgba(99,102,241,.1); color: #6366f1; }
.act-transfer:hover { background: rgba(99,102,241,.2); }

/* ═══ BADGES ═══ */
.affm-badge { display: inline-flex; padding: 4px 10px; border-radius: 8px; font-size: .7rem; font-weight: 700; text-transform: uppercase; letter-spacing: .04em; }
.badge-pending { background: #fef3c7; color: #92400e; }
.badge-approved { background: #dbeafe; color: #1e40af; }
.badge-transferred { background: #dcfce7; color: #166534; }
.badge-rejected { background: #fef2f2; color: #991b1b; }
.badge-commission { background: #eef2ff; color: #4338ca; }
.badge-bonus { background: #fdf4ff; color: #86198f; }
.badge-adjustment { background: #fff7ed; color: #c2410c; }

/* ═══ BALANCE ═══ */
.affm-balance-strip { display: flex; gap: 24px; padding: 18px 28px; background: linear-gradient(135deg, #f0fdf4 0%, #ecfdf5 100%); border-bottom: 1px solid #bbf7d0; }
.affm-balance-item { display: flex; align-items: center; gap: 10px; }
.affm-balance-item i { font-size: 1.2rem; color: #059669; }
.affm-balance-label { font-size: .78rem; color: #64748b; }
.affm-balance-value { font-size: 1.15rem; font-weight: 800; color: #059669; }

/* ═══ MODAL ═══ */
.affm-overlay { position: fixed; inset: 0; z-index: 9999; background: rgba(0,0,0,.45); backdrop-filter: blur(4px); display: flex; align-items: center; justify-content: center; padding: 20px; }
.affm-modal { background: #fff; border-radius: 18px; width: 100%; max-width: 480px; box-shadow: 0 20px 60px rgba(0,0,0,.15); animation: modalIn .2s ease; }
@keyframes modalIn { from { opacity: 0; transform: scale(.95) translateY(10px); } }
.affm-modal-head { display: flex; align-items: center; justify-content: space-between; padding: 22px 28px; border-bottom: 1px solid #f1f5f9; }
.affm-modal-head h3 { font-size: 1.05rem; font-weight: 700; color: #0f172a; margin: 0; }
.affm-modal-x { width: 34px; height: 34px; border-radius: 10px; display: flex; align-items: center; justify-content: center; border: none; background: transparent; font-size: 1.4rem; color: #94a3b8; cursor: pointer; }
.affm-modal-x:hover { background: #f1f5f9; color: #0f172a; }
.affm-modal-body { padding: 20px 28px; }
.affm-modal-info { display: flex; justify-content: space-between; padding: 8px 0; border-bottom: 1px solid #f1f5f9; font-size: .85rem; }
.affm-modal-info span { color: #94a3b8; }
.affm-modal-info strong { color: #0f172a; }
.affm-modal-amt { color: #059669 !important; font-size: 1rem; }
.affm-modal-foot { display: flex; justify-content: flex-end; gap: 10px; padding: 16px 28px 22px; }

/* ═══ EMPTY / LOADING ═══ */
.affm-empty { display: flex; flex-direction: column; align-items: center; padding: 60px 20px; text-align: center; }
.affm-empty i { font-size: 2.5rem; color: #94a3b8; margin-bottom: 10px; }
.affm-empty p { color: #94a3b8; font-size: .85rem; margin: 0; }
.affm-loading { display: flex; flex-direction: column; gap: 10px; padding: 20px 28px; }
.affm-skel { height: 52px; border-radius: 12px; }
.shimmer { background: linear-gradient(90deg, #e8ecf1 25%, #f1f5f9 50%, #e8ecf1 75%); background-size: 200% 100%; animation: shimmer 1.5s infinite; }
@keyframes shimmer { 0% { background-position: 200% 0 } 100% { background-position: -200% 0 } }
.ri-spin { animation: spin 1s linear infinite; }
@keyframes spin { to { transform: rotate(360deg); } }

@media (max-width: 768px) {
  .affm-hero { padding: 24px 20px 18px; }
  .affm-form-grid { grid-template-columns: 1fr; }
  .affm-tab-bar { flex-direction: column; }
  .affm-card-head { flex-direction: column; align-items: flex-start; }
  .affm-card-body { padding: 16px; }
  .affm-preview { flex-direction: column; align-items: flex-start; }
  .affm-preview-mid { border-left: none; padding: 12px 0 0; border-top: 2px solid #e0e7ff; }
  .affm-preview-right { margin-left: 0; }
}
</style>
