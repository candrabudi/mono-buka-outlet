<template>
  <div class="animate-in">
    <div class="eo-hero">
      <h1 class="eo-hero-title"><i class="ri-file-list-3-line"></i> Pesanan Ebook</h1>
      <p class="eo-hero-sub">Pantau semua transaksi ebook dari mitra</p>
    </div>

    <div class="eo-filters">
      <SearchSelect
        v-model="filters.status"
        :options="statusOptions"
        placeholder="Semua Status"
        empty-label="Semua Status"
        @update:model-value="load"
      />
      <SearchSelect
        v-model="filters.download_status"
        :options="downloadOptions"
        placeholder="Semua Download"
        empty-label="Semua Download"
        @update:model-value="load"
      />
    </div>

    <div v-if="loading" class="eo-loading">Memuat data...</div>
    <div v-else-if="!orders.length" class="eo-empty">
      <div style="margin-bottom:16px;"><i class="ri-mail-open-line" style="font-size:3rem;color:#94a3b8"></i></div>
      <h3>Belum ada pesanan</h3>
    </div>

    <div v-else class="eo-table-wrap">
      <table class="eo-table">
        <thead>
          <tr>
            <th>Order</th>
            <th>Mitra</th>
            <th>Ebook</th>
            <th>Jumlah</th>
            <th>Status</th>
            <th>Bukti</th>
            <th>Download</th>
            <th>Tanggal</th>
            <th>Aksi</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="o in orders" :key="o.id">
            <td><span class="eo-order-num">{{ o.order_number }}</span></td>
            <td>
              <div class="eo-user-name">{{ o.user?.name }}</div>
              <div class="eo-user-email">{{ o.user?.email }}</div>
            </td>
            <td><span class="eo-ebook-title">{{ o.ebook?.title }}</span></td>
            <td><span class="eo-amount">{{ formatRp(o.amount) }}</span></td>
            <td><span class="eo-badge" :class="'badge-' + o.status.toLowerCase()">{{ o.status }}</span></td>
            <td>
              <a v-if="o.payment_proof_url" href="#" @click.prevent="proofTarget = o" class="eo-proof-link">
                <i class="ri-image-line"></i> Lihat
              </a>
              <span v-else class="eo-no-proof">—</span>
            </td>
            <td><span class="eo-badge" :class="'dl-' + o.download_status.toLowerCase()">{{ dlLabel(o.download_status) }}</span></td>
            <td><span class="eo-date">{{ formatDate(o.created_at) }}</span></td>
            <td>
              <div v-if="o.status === 'PENDING' && o.payment_proof_url" class="eo-actions">
                <button @click="doApprove(o)" class="act-btn act-approve" :disabled="processing" title="Approve">
                  <i class="ri-check-line"></i>
                </button>
                <button @click="rejectTarget = o" class="act-btn act-reject" :disabled="processing" title="Reject">
                  <i class="ri-close-line"></i>
                </button>
              </div>
              <span v-else-if="o.status === 'PENDING'" class="eo-hint">Menunggu bukti</span>
            </td>
          </tr>
        </tbody>
      </table>
    </div>

    <!-- Proof preview modal -->
    <div v-if="proofTarget" class="modal-overlay" @click.self="proofTarget = null">
      <div class="modal-proof">
        <div class="modal-proof-header">
          <h3>Bukti Pembayaran</h3>
          <button @click="proofTarget = null" class="modal-close"><i class="ri-close-line"></i></button>
        </div>
        <div class="modal-proof-info">
          <span><strong>{{ proofTarget.user?.name }}</strong> — {{ proofTarget.ebook?.title }}</span>
          <span class="eo-amount">{{ formatRp(proofTarget.amount) }}</span>
        </div>
        <img :src="proofTarget.payment_proof_url" alt="Bukti pembayaran" class="modal-proof-img" />
        <div v-if="proofTarget.status === 'PENDING'" class="modal-proof-actions">
          <button @click="doApprove(proofTarget); proofTarget = null" class="eo-btn-action eo-btn-approve" :disabled="processing">
            <i class="ri-check-double-line"></i> Approve Pembayaran
          </button>
          <button @click="rejectTarget = proofTarget; proofTarget = null" class="eo-btn-action eo-btn-reject-action">
            <i class="ri-close-circle-line"></i> Reject
          </button>
        </div>
      </div>
    </div>

    <!-- Reject modal -->
    <div v-if="rejectTarget" class="modal-overlay" @click.self="rejectTarget = null">
      <div class="modal-content">
        <div class="modal-body">
          <div style="margin-bottom:12px"><i class="ri-close-circle-line" style="font-size:2.5rem;color:#ef4444"></i></div>
          <h3>Reject Pembayaran?</h3>
          <p style="color:#64748b;font-size:0.85rem;margin-top:8px">Pesanan <strong>{{ rejectTarget.order_number }}</strong> dari {{ rejectTarget.user?.name }}</p>
          <textarea v-model="rejectNote" placeholder="Alasan penolakan (opsional)" class="eo-textarea"></textarea>
        </div>
        <div class="modal-actions">
          <button @click="rejectTarget = null" class="eo-btn-action eo-btn-secondary">Batal</button>
          <button @click="doReject" class="eo-btn-action eo-btn-reject-action" :disabled="processing">
            <i class="ri-close-circle-line"></i> {{ processing ? 'Memproses...' : 'Ya, Reject' }}
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { ebookOrderApi } from '../../services/api'
import { useToastStore } from '../../stores/toast'
import SearchSelect from '../../components/SearchSelect.vue'

const toast = useToastStore()
const orders = ref([])
const loading = ref(false)
const processing = ref(false)
const filters = ref({ status: '', download_status: '' })
const proofTarget = ref(null)
const rejectTarget = ref(null)
const rejectNote = ref('')

const statusOptions = [
  { value: 'PENDING', label: 'Pending' },
  { value: 'PAID', label: 'Paid' },
  { value: 'EXPIRED', label: 'Expired' },
  { value: 'FAILED', label: 'Failed' },
  { value: 'CANCELED', label: 'Canceled' },
]
const downloadOptions = [
  { value: 'NONE', label: 'Belum Request' },
  { value: 'REQUESTED', label: 'Menunggu Approval' },
  { value: 'APPROVED', label: 'Approved' },
  { value: 'REJECTED', label: 'Rejected' },
]

onMounted(load)

async function load() {
  loading.value = true
  try {
    const { data } = await ebookOrderApi.list(filters.value)
    orders.value = data.data || []
  } catch { toast.error('Gagal memuat pesanan') }
  finally { loading.value = false }
}

async function doApprove(o) {
  if (processing.value) return
  processing.value = true
  try {
    await ebookOrderApi.approvePayment(o.id)
    toast.success('Pembayaran disetujui')
    load()
  } catch (err) {
    toast.error(err.response?.data?.error || 'Gagal approve')
  } finally { processing.value = false }
}

async function doReject() {
  if (processing.value || !rejectTarget.value) return
  processing.value = true
  try {
    await ebookOrderApi.rejectPayment(rejectTarget.value.id, { note: rejectNote.value })
    toast.success('Pembayaran ditolak')
    rejectTarget.value = null
    rejectNote.value = ''
    load()
  } catch (err) {
    toast.error(err.response?.data?.error || 'Gagal reject')
  } finally { processing.value = false }
}

function formatRp(v) { return 'Rp ' + Number(v || 0).toLocaleString('id-ID') }
function formatDate(d) { return d ? new Date(d).toLocaleDateString('id-ID', { day: '2-digit', month: 'short', year: 'numeric' }) : '-' }
function dlLabel(s) {
  const map = { NONE: 'Belum', REQUESTED: 'Menunggu', APPROVED: 'Disetujui', REJECTED: 'Ditolak' }
  return map[s] || s
}
</script>

<style scoped>
.eo-hero { background: linear-gradient(135deg, #0f0c29, #302b63, #24243e); border-radius: 16px; padding: 32px 36px; margin-bottom: 24px; }
.eo-hero-title { font-size: 1.6rem; font-weight: 800; color: #fff; }
.eo-hero-sub { font-size: 0.85rem; color: rgba(255,255,255,0.5); margin-top: 6px; }

.eo-filters { display: flex; gap: 12px; margin-bottom: 24px; flex-wrap: wrap; }
.eo-select { height: 42px; padding: 0 34px 0 14px; border-radius: 12px; border: 1.5px solid #e2e8f0; font-size: 0.8rem; font-weight: 500; color: #475569; background: white; cursor: pointer; appearance: none; background-image: url("data:image/svg+xml,%3csvg xmlns='http://www.w3.org/2000/svg' fill='none' viewBox='0 0 20 20'%3e%3cpath stroke='%236b7280' stroke-linecap='round' stroke-linejoin='round' stroke-width='1.5' d='M6 8l4 4 4-4'/%3e%3c/svg%3e"); background-position: right 10px center; background-repeat: no-repeat; background-size: 1.4em; }

.eo-loading, .eo-empty { text-align: center; padding: 80px 32px; background: white; border-radius: 16px; border: 2px dashed #e2e8f0; }
.eo-empty h3 { font-size: 1.1rem; font-weight: 700; color: #0f172a; margin: 0; }

.eo-table-wrap { background: white; border-radius: 16px; border: 1px solid #e8ecf1; overflow: hidden; }
.eo-table { width: 100%; border-collapse: collapse; }
.eo-table th { text-align: left; padding: 14px 16px; font-size: 0.72rem; font-weight: 700; color: #64748b; text-transform: uppercase; background: #f8fafc; border-bottom: 1px solid #e8ecf1; }
.eo-table td { padding: 14px 16px; border-bottom: 1px solid #f1f5f9; font-size: 0.85rem; }

.eo-order-num { font-family: monospace; font-size: 0.78rem; color: #6366f1; font-weight: 600; }
.eo-user-name { font-weight: 600; color: #0f172a; font-size: 0.85rem; }
.eo-user-email { font-size: 0.72rem; color: #94a3b8; }
.eo-ebook-title { font-weight: 600; color: #0f172a; }
.eo-amount { font-weight: 700; color: #f59e0b; }
.eo-date { font-size: 0.8rem; color: #64748b; }
.eo-hint { font-size: 0.72rem; color: #94a3b8; font-style: italic; }

.eo-badge { padding: 4px 12px; border-radius: 8px; font-size: 0.7rem; font-weight: 700; text-transform: uppercase; }
.badge-pending { background: #fef3c7; color: #92400e; }
.badge-paid { background: #dcfce7; color: #166534; }
.badge-expired { background: #f1f5f9; color: #64748b; }
.badge-canceled { background: #f1f5f9; color: #64748b; }
.badge-failed { background: #fee2e2; color: #991b1b; }
.dl-none { background: #f1f5f9; color: #64748b; }
.dl-requested { background: #fef3c7; color: #92400e; }
.dl-approved { background: #dcfce7; color: #166534; }
.dl-rejected { background: #fee2e2; color: #991b1b; }

.eo-proof-link { color: #6366f1; font-size: 0.8rem; font-weight: 600; text-decoration: none; display: flex; align-items: center; gap: 4px; }
.eo-proof-link:hover { text-decoration: underline; }
.eo-no-proof { color: #cbd5e1; font-size: 0.8rem; }

.eo-actions { display: flex; gap: 6px; }
.act-btn { width: 32px; height: 32px; border-radius: 8px; border: none; cursor: pointer; display: flex; align-items: center; justify-content: center; font-size: 1rem; }
.act-approve { background: #dcfce7; color: #16a34a; }
.act-approve:hover { background: #16a34a; color: white; }
.act-reject { background: #fee2e2; color: #dc2626; }
.act-reject:hover { background: #dc2626; color: white; }
.act-btn:disabled { opacity: 0.5; cursor: not-allowed; }

/* Modals */
.modal-overlay { position: fixed; inset: 0; background: rgba(0,0,0,0.5); display: flex; align-items: center; justify-content: center; z-index: 999; backdrop-filter: blur(4px); }

.modal-proof { background: white; border-radius: 16px; max-width: 520px; width: 92%; overflow: hidden; box-shadow: 0 20px 60px rgba(0,0,0,0.2); }
.modal-proof-header { display: flex; justify-content: space-between; align-items: center; padding: 20px 24px 12px; }
.modal-proof-header h3 { font-size: 1rem; font-weight: 700; color: #0f172a; margin: 0; }
.modal-close { background: none; border: none; font-size: 1.2rem; color: #94a3b8; cursor: pointer; }
.modal-proof-info { padding: 0 24px 16px; display: flex; justify-content: space-between; align-items: center; font-size: 0.85rem; color: #475569; }
.modal-proof-img { width: 100%; max-height: 400px; object-fit: contain; background: #f8fafc; }
.modal-proof-actions { display: flex; gap: 10px; padding: 16px 24px; }

.modal-content { background: white; border-radius: 16px; max-width: 400px; width: 90%; overflow: hidden; box-shadow: 0 20px 60px rgba(0,0,0,0.2); }
.modal-body { padding: 2rem; text-align: center; }
.modal-body h3 { font-size: 1.1rem; font-weight: 700; color: #0f172a; margin: 0; }
.modal-actions { display: flex; gap: 10px; padding: 0 2rem 1.5rem; }

.eo-textarea { width: 100%; min-height: 60px; border: 1.5px solid #e2e8f0; border-radius: 10px; padding: 10px 14px; font-size: 0.85rem; color: #0f172a; margin-top: 16px; resize: vertical; font-family: inherit; }

.eo-btn-action { flex: 1; padding: 12px; border-radius: 10px; font-size: 0.85rem; font-weight: 700; border: none; cursor: pointer; display: flex; align-items: center; justify-content: center; gap: 6px; }
.eo-btn-approve { background: linear-gradient(135deg, #22c55e, #16a34a); color: white; }
.eo-btn-reject-action { background: linear-gradient(135deg, #ef4444, #dc2626); color: white; }
.eo-btn-secondary { background: #f1f5f9; color: #475569; }
.eo-btn-action:disabled { opacity: 0.5; cursor: not-allowed; }

@media (max-width: 768px) {
  .eo-table-wrap { overflow-x: auto; }
  .eo-hero { padding: 24px 20px; }
}
</style>
