<template>
  <div class="animate-in">
    <div class="eo-hero">
      <h1 class="eo-hero-title"><i class="ri-file-list-3-line"></i> Pesanan Ebook Saya</h1>
      <p class="eo-hero-sub">Riwayat pembelian & status download ebook Anda</p>
    </div>

    <div v-if="loading" class="eo-loading">Memuat data...</div>
    <div v-else-if="!orders.length" class="eo-empty">
      <div style="margin-bottom:16px;"><i class="ri-mail-open-line" style="font-size:3rem;color:#94a3b8"></i></div>
      <h3>Belum ada pesanan</h3>
      <p>Jelajahi ebook dan mulai pembelian pertama Anda</p>
      <router-link to="/ebooks" class="eo-btn eo-btn-browse" style="margin-top:16px"><i class="ri-book-open-line"></i> Jelajahi Ebook</router-link>
    </div>

    <div v-else class="eo-list">
      <div v-for="o in orders" :key="o.id" class="eo-card">
        <div class="eo-card-cover" :style="coverStyle(o.ebook)">
          <div class="eo-card-status" :class="'st-' + o.status.toLowerCase()">{{ o.status }}</div>
        </div>
        <div class="eo-card-body">
          <h3 class="eo-card-title">{{ o.ebook?.title }}</h3>
          <div class="eo-card-meta">
            <span>{{ formatRp(o.amount) }}</span>
            <span>•</span>
            <span>{{ formatDate(o.created_at) }}</span>
          </div>
          <div v-if="o.status === 'PAID'" class="eo-card-dl">
            <span class="eo-dl-status" :class="'dl-' + o.download_status.toLowerCase()">
              Download: {{ dlLabel(o.download_status) }}
            </span>
            <span v-if="o.download_status === 'REJECTED' && o.download_note" class="eo-dl-note">
              Alasan: {{ o.download_note }}
            </span>
          </div>
        </div>
        <div class="eo-card-actions" v-if="o.status === 'PAID'">
          <router-link :to="`/ebooks/${o.ebook_id}/read`" class="eo-btn eo-btn-read"><i class="ri-book-read-line"></i> Baca</router-link>
          <template v-if="o.download_status === 'NONE'">
            <button @click="requestDownload(o)" class="eo-btn eo-btn-dl" :disabled="processing === o.id">
              {{ processing === o.id ? '...' : '' }} <i class="ri-download-cloud-line"></i> Request Download
            </button>
          </template>
          <span v-else-if="o.download_status === 'REQUESTED'" class="eo-btn eo-btn-pending"><i class="ri-time-line"></i> Menunggu</span>
          <a v-else-if="o.download_status === 'APPROVED'" :href="downloadUrl(o.ebook_id)" class="eo-btn eo-btn-download">
            <i class="ri-download-line"></i> Download
          </a>
        </div>
        <div class="eo-card-actions" v-else-if="o.status === 'PENDING'">
          <div v-if="o.payment_proof_url" class="eo-proof-badge"><i class="ri-image-line"></i> Bukti dikirim</div>
          <button v-if="o.midtrans_snap_token" @click="resumePayment(o)" class="eo-btn eo-btn-pay" :disabled="paying"><i class="ri-bank-card-line"></i> {{ paying ? 'Memproses...' : 'Lanjut Bayar' }}</button>
          <label class="eo-btn eo-btn-proof" :class="{ disabled: uploading === o.id }">
            <i class="ri-upload-cloud-line"></i> {{ uploading === o.id ? 'Uploading...' : (o.payment_proof_url ? 'Ganti Bukti' : 'Upload Bukti') }}
            <input type="file" accept="image/*" @change="handleUploadProof($event, o)" style="display:none" :disabled="uploading === o.id" />
          </label>
          <button @click="cancelTarget = o" class="eo-btn eo-btn-cancel" :disabled="cancelling === o.id"><i class="ri-close-circle-line"></i> {{ cancelling === o.id ? 'Membatalkan...' : 'Batalkan' }}</button>
        </div>
        <div class="eo-card-actions" v-else-if="o.status === 'CANCELED'">
          <span class="eo-btn eo-btn-canceled"><i class="ri-close-circle-line"></i> Dibatalkan</span>
        </div>
      </div>
    </div>

    <!-- Cancel confirmation modal -->
    <div v-if="cancelTarget" class="modal-overlay" @click.self="cancelTarget = null">
      <div class="modal-content">
        <div class="modal-body">
          <div style="margin-bottom:12px"><i class="ri-close-circle-line" style="font-size:2.5rem;color:#ef4444"></i></div>
          <h3>Batalkan Pesanan?</h3>
          <p style="color:#64748b;font-size:0.85rem;margin-top:8px">Yakin ingin membatalkan pesanan <strong>{{ cancelTarget.ebook?.title }}</strong>?</p>
        </div>
        <div class="modal-actions">
          <button @click="cancelTarget = null" class="eo-btn eo-btn-secondary">Kembali</button>
          <button @click="confirmCancel" class="eo-btn eo-btn-danger" :disabled="cancelling"><i class="ri-close-circle-line"></i> {{ cancelling ? 'Membatalkan...' : 'Ya, Batalkan' }}</button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { ebookApi, uploadApi, midtransApi } from '../services/api'
import { useToastStore } from '../stores/toast'

const toast = useToastStore()
const orders = ref([])
const loading = ref(false)
const processing = ref(null)
const paying = ref(false)
const cancelling = ref(null)
const cancelTarget = ref(null)
const uploading = ref(null)

onMounted(load)

async function load() {
  loading.value = true
  try {
    const { data } = await ebookApi.myOrders()
    orders.value = data.data || []
  } catch { toast.error('Gagal memuat pesanan') }
  finally { loading.value = false }
}

async function requestDownload(o) {
  processing.value = o.id
  try {
    await ebookApi.requestDownload(o.id)
    o.download_status = 'REQUESTED'
    toast.success('Request download berhasil dikirim')
  } catch (err) {
    toast.error(err.response?.data?.error || 'Gagal request download')
  } finally { processing.value = null }
}

async function resumePayment(o) {
  if (paying.value) return
  paying.value = true
  try {
    const { data: keyData } = await midtransApi.getClientKey()
    const script = document.createElement('script')
    script.src = keyData.snap_url
    script.setAttribute('data-client-key', keyData.client_key)
    script.onload = () => {
      window.snap.pay(o.midtrans_snap_token, {
        onSuccess: () => { toast.success('Pembayaran berhasil!'); paying.value = false; load() },
        onPending: () => { toast.info('Pembayaran pending'); paying.value = false },
        onError: () => { toast.error('Pembayaran gagal'); paying.value = false },
        onClose: () => { paying.value = false },
      })
    }
    document.head.appendChild(script)
  } catch { toast.error('Gagal memuat pembayaran'); paying.value = false }
}

async function confirmCancel() {
  const o = cancelTarget.value
  if (!o || cancelling.value) return
  cancelling.value = o.id
  try {
    await ebookApi.cancelOrder(o.id)
    toast.success('Pesanan berhasil dibatalkan')
    cancelTarget.value = null
    load()
  } catch (err) {
    toast.error(err.response?.data?.error || 'Gagal membatalkan pesanan')
  } finally { cancelling.value = null }
}

async function handleUploadProof(event, o) {
  const file = event.target.files?.[0]
  if (!file) return
  uploading.value = o.id
  try {
    const { data: uploadData } = await uploadApi.upload(file)
    const url = uploadData.url || uploadData.data?.url
    await ebookApi.uploadProof(o.id, { payment_proof_url: url })
    toast.success('Bukti pembayaran berhasil diupload')
    load()
  } catch (err) {
    toast.error(err.response?.data?.error || 'Gagal upload bukti')
  } finally {
    uploading.value = null
    event.target.value = ''
  }
}

function downloadUrl(ebookId) {
  const token = localStorage.getItem('mitra_token')
  const base = import.meta.env.VITE_API_BASE_URL || '/api/v1/mitra'
  return `${base}/ebooks/${ebookId}/download?token=${token}`
}

function formatRp(v) { return 'Rp ' + Number(v || 0).toLocaleString('id-ID') }
function formatDate(d) { return d ? new Date(d).toLocaleDateString('id-ID', { day: '2-digit', month: 'short', year: 'numeric' }) : '-' }
function dlLabel(s) { return { NONE: 'Belum diminta', REQUESTED: 'Menunggu approval', APPROVED: 'Disetujui', REJECTED: 'Ditolak' }[s] || s }

function coverStyle(eb) {
  if (eb?.cover_url) return { backgroundImage: `url(${eb.cover_url})`, backgroundSize: 'cover', backgroundPosition: 'center' }
  const g = ['linear-gradient(135deg,#667eea,#764ba2)','linear-gradient(135deg,#f093fb,#f5576c)','linear-gradient(135deg,#4facfe,#00f2fe)','linear-gradient(135deg,#43e97b,#38f9d7)']
  let h = 0; for (let i = 0; i < (eb?.title||'').length; i++) h = ((h << 5) - h) + eb.title.charCodeAt(i)
  return { background: g[Math.abs(h) % g.length] }
}
</script>

<style scoped>
.eo-hero { background: linear-gradient(135deg, #1e1b4b, #312e81, #4338ca); border-radius: 16px; padding: 32px 36px; margin-bottom: 24px; }
.eo-hero-title { font-size: 1.6rem; font-weight: 800; color: #fff; }
.eo-hero-sub { font-size: 0.85rem; color: rgba(255,255,255,0.55); margin-top: 6px; }

.eo-loading, .eo-empty { text-align: center; padding: 80px 32px; background: white; border-radius: 16px; border: 2px dashed #e2e8f0; }
.eo-empty h3 { font-size: 1.1rem; font-weight: 700; color: #0f172a; margin: 0; }
.eo-empty p { font-size: 0.85rem; color: #94a3b8; margin-top: 6px; }

.eo-list { display: flex; flex-direction: column; gap: 16px; }
.eo-card { background: white; border-radius: 16px; border: 1px solid #e8ecf1; overflow: hidden; display: flex; align-items: stretch; }

.eo-card-cover { width: 120px; flex-shrink: 0; position: relative; display: flex; align-items: center; justify-content: center; }
.eo-card-status { position: absolute; top: 8px; left: 8px; padding: 3px 10px; border-radius: 6px; font-size: 0.65rem; font-weight: 700; text-transform: uppercase; }
.st-pending { background: #fef3c7; color: #92400e; }
.st-paid { background: #dcfce7; color: #166534; }
.st-expired { background: #f1f5f9; color: #64748b; }
.st-failed { background: #fee2e2; color: #991b1b; }

.eo-card-body { flex: 1; padding: 16px 20px; display: flex; flex-direction: column; justify-content: center; min-width: 0; }
.eo-card-title { font-size: 0.95rem; font-weight: 700; color: #0f172a; margin: 0 0 4px; white-space: nowrap; overflow: hidden; text-overflow: ellipsis; }
.eo-card-meta { font-size: 0.78rem; color: #94a3b8; display: flex; gap: 8px; margin-bottom: 8px; }
.eo-card-dl { display: flex; flex-direction: column; gap: 4px; }
.eo-dl-status { font-size: 0.75rem; font-weight: 600; }
.dl-none { color: #94a3b8; }
.dl-requested { color: #d97706; }
.dl-approved { color: #16a34a; }
.dl-rejected { color: #dc2626; }
.eo-dl-note { font-size: 0.72rem; color: #94a3b8; font-style: italic; }

.eo-card-actions { display: flex; flex-direction: column; gap: 6px; padding: 12px 16px; justify-content: center; border-left: 1px solid #f1f5f9; }
.eo-btn { padding: 8px 16px; border-radius: 10px; font-size: 0.78rem; font-weight: 700; border: none; cursor: pointer; text-decoration: none; text-align: center; white-space: nowrap; }
.eo-btn-read { background: #eef2ff; color: #4f46e5; }
.eo-btn-dl { background: #f0fdf4; color: #16a34a; }
.eo-btn-download { background: linear-gradient(135deg, #6366f1, #8b5cf6); color: white; }
.eo-btn-pending { background: #fef3c7; color: #92400e; cursor: default; }
.eo-btn-pay { background: linear-gradient(135deg, #f59e0b, #ef4444); color: white; }
.eo-btn-browse { display: inline-block; background: linear-gradient(135deg, #6366f1, #8b5cf6); color: white; padding: 12px 24px; border-radius: 12px; text-decoration: none; }
.eo-btn-cancel { background: #fee2e2; color: #dc2626; }
.eo-btn-canceled { background: #f1f5f9; color: #94a3b8; cursor: default; }
.eo-btn-proof { background: #eef2ff; color: #4f46e5; cursor: pointer; display: flex; align-items: center; gap: 4px; }
.eo-btn-proof.disabled { opacity: 0.5; cursor: not-allowed; }
.eo-proof-badge { font-size: 0.7rem; font-weight: 600; color: #16a34a; display: flex; align-items: center; gap: 4px; white-space: nowrap; }
.eo-btn:disabled { opacity: 0.5; cursor: not-allowed; }
.st-canceled { background: #f1f5f9; color: #64748b; }

.modal-overlay { position: fixed; inset: 0; background: rgba(0,0,0,0.5); display: flex; align-items: center; justify-content: center; z-index: 999; backdrop-filter: blur(4px); }
.modal-content { background: white; border-radius: 16px; max-width: 400px; width: 90%; overflow: hidden; box-shadow: 0 20px 60px rgba(0,0,0,0.2); }
.modal-body { padding: 2rem; text-align: center; }
.modal-body h3 { font-size: 1.1rem; font-weight: 700; color: #0f172a; margin: 0; }
.modal-actions { display: flex; gap: 10px; padding: 0 2rem 1.5rem; }
.modal-actions .eo-btn { flex: 1; padding: 12px; border-radius: 10px; font-size: 0.85rem; font-weight: 700; }
.eo-btn-secondary { background: #f1f5f9; color: #475569; border: none; cursor: pointer; }
.eo-btn-danger { background: linear-gradient(135deg, #ef4444, #dc2626); color: white; border: none; cursor: pointer; }

@media (max-width: 768px) {
  .eo-card { flex-direction: column; }
  .eo-card-cover { width: 100%; height: 120px; }
  .eo-card-actions { flex-direction: row; border-left: none; border-top: 1px solid #f1f5f9; }
}
</style>
