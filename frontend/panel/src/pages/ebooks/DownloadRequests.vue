<template>
  <div class="animate-in">
    <div class="dl-hero">
      <h1 class="dl-hero-title">📥 Request Download Ebook</h1>
      <p class="dl-hero-sub">Approve atau reject permintaan download dari mitra</p>
    </div>

    <div v-if="loading" class="dl-loading">Memuat data...</div>

    <div v-else-if="!requests.length" class="dl-empty">
      <div style="font-size:3rem;margin-bottom:16px;"></div>
      <h3>Tidak ada request pending</h3>
      <p>Semua permintaan download sudah diproses</p>
    </div>

    <div v-else class="dl-table-wrap">
      <table class="dl-table">
        <thead>
          <tr>
            <th>Mitra</th>
            <th>Ebook</th>
            <th>Tanggal Request</th>
            <th>Aksi</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="req in requests" :key="req.id">
            <td>
              <div class="dl-user">
                <div class="dl-user-avatar">{{ (req.user?.name || '?').charAt(0).toUpperCase() }}</div>
                <div>
                  <div class="dl-user-name">{{ req.user?.name }}</div>
                  <div class="dl-user-email">{{ req.user?.email }}</div>
                </div>
              </div>
            </td>
            <td>
              <div class="dl-ebook">
                <img v-if="req.ebook?.cover_url" :src="req.ebook.cover_url" class="dl-ebook-cover" />
                <div>
                  <div class="dl-ebook-title">{{ req.ebook?.title }}</div>
                  <div class="dl-ebook-price">{{ formatRp(req.amount) }}</div>
                </div>
              </div>
            </td>
            <td>
              <span class="dl-date">{{ formatDate(req.download_requested_at) }}</span>
            </td>
            <td>
              <div class="dl-actions">
                <button @click="approve(req)" class="dl-btn dl-btn-approve" :disabled="processing === req.id">
                  Approve
                </button>
                <button @click="openReject(req)" class="dl-btn dl-btn-reject" :disabled="processing === req.id">
                  Reject
                </button>
              </div>
            </td>
          </tr>
        </tbody>
      </table>
    </div>

    <!-- Reject modal -->
    <div v-if="rejectTarget" class="modal-overlay" @click.self="rejectTarget = null">
      <div class="modal-content" style="max-width:420px">
        <div class="modal-body" style="padding:2rem">
          <h3 style="font-size:1.1rem;font-weight:700;margin-bottom:16px;">Alasan Penolakan</h3>
          <textarea v-model="rejectNote" class="form-textarea" rows="3" placeholder="Tuliskan alasan penolakan..." style="width:100%;box-sizing:border-box;"></textarea>
        </div>
        <div class="modal-footer" style="justify-content:flex-end;gap:12px;padding:0 2rem 1.5rem">
          <button @click="rejectTarget = null" class="btn btn-secondary">Batal</button>
          <button @click="doReject" class="btn btn-danger" :disabled="!rejectNote.trim() || processing">Reject</button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { ebookOrderApi } from '../../services/api'
import { useToastStore } from '../../stores/toast'

const toast = useToastStore()
const requests = ref([])
const loading = ref(false)
const processing = ref(null)
const rejectTarget = ref(null)
const rejectNote = ref('')

onMounted(load)

async function load() {
  loading.value = true
  try {
    const { data } = await ebookOrderApi.downloadRequests()
    requests.value = data.data || []
  } catch { toast.error('Gagal memuat data') }
  finally { loading.value = false }
}

async function approve(req) {
  processing.value = req.id
  try {
    await ebookOrderApi.approveDownload(req.id, { note: '' })
    toast.success('Download disetujui')
    requests.value = requests.value.filter(r => r.id !== req.id)
  } catch { toast.error('Gagal approve') }
  finally { processing.value = null }
}

function openReject(req) { rejectTarget.value = req; rejectNote.value = '' }

async function doReject() {
  processing.value = rejectTarget.value.id
  try {
    await ebookOrderApi.rejectDownload(rejectTarget.value.id, { note: rejectNote.value })
    toast.success('Download ditolak')
    requests.value = requests.value.filter(r => r.id !== rejectTarget.value.id)
    rejectTarget.value = null
  } catch { toast.error('Gagal reject') }
  finally { processing.value = null }
}

function formatRp(v) { return 'Rp ' + Number(v || 0).toLocaleString('id-ID') }
function formatDate(d) { return d ? new Date(d).toLocaleDateString('id-ID', { day: '2-digit', month: 'short', year: 'numeric', hour: '2-digit', minute: '2-digit' }) : '-' }
</script>

<style scoped>
.dl-hero { background: linear-gradient(135deg, #0f0c29, #302b63, #24243e); border-radius: 16px; padding: 32px 36px; margin-bottom: 24px; }
.dl-hero-title { font-size: 1.6rem; font-weight: 800; color: #fff; }
.dl-hero-sub { font-size: 0.85rem; color: rgba(255,255,255,0.5); margin-top: 6px; }

.dl-loading, .dl-empty { text-align: center; padding: 80px 32px; background: white; border-radius: 16px; border: 2px dashed #e2e8f0; }
.dl-empty h3 { font-size: 1.1rem; font-weight: 700; color: #0f172a; margin: 0; }
.dl-empty p { font-size: 0.85rem; color: #94a3b8; margin-top: 6px; }

.dl-table-wrap { background: white; border-radius: 16px; border: 1px solid #e8ecf1; overflow: hidden; }
.dl-table { width: 100%; border-collapse: collapse; }
.dl-table th { text-align: left; padding: 14px 20px; font-size: 0.75rem; font-weight: 700; color: #64748b; text-transform: uppercase; letter-spacing: 0.05em; background: #f8fafc; border-bottom: 1px solid #e8ecf1; }
.dl-table td { padding: 16px 20px; border-bottom: 1px solid #f1f5f9; vertical-align: middle; }
.dl-table tr:last-child td { border-bottom: none; }

.dl-user { display: flex; align-items: center; gap: 12px; }
.dl-user-avatar { width: 36px; height: 36px; border-radius: 10px; background: linear-gradient(135deg, #6366f1, #8b5cf6); color: white; display: flex; align-items: center; justify-content: center; font-weight: 800; font-size: 0.85rem; flex-shrink: 0; }
.dl-user-name { font-size: 0.85rem; font-weight: 600; color: #0f172a; }
.dl-user-email { font-size: 0.75rem; color: #94a3b8; }

.dl-ebook { display: flex; align-items: center; gap: 12px; }
.dl-ebook-cover { width: 40px; height: 52px; border-radius: 6px; object-fit: cover; flex-shrink: 0; }
.dl-ebook-title { font-size: 0.85rem; font-weight: 600; color: #0f172a; }
.dl-ebook-price { font-size: 0.75rem; color: #f59e0b; font-weight: 700; }

.dl-date { font-size: 0.8rem; color: #64748b; }

.dl-actions { display: flex; gap: 8px; }
.dl-btn { padding: 8px 16px; border-radius: 10px; font-size: 0.78rem; font-weight: 700; border: none; cursor: pointer; }
.dl-btn-approve { background: #f0fdf4; color: #16a34a; }
.dl-btn-approve:hover { background: #dcfce7; }
.dl-btn-reject { background: #fef2f2; color: #dc2626; }
.dl-btn-reject:hover { background: #fee2e2; }
.dl-btn:disabled { opacity: 0.5; cursor: not-allowed; }

.form-textarea { padding: 12px 16px; border: 1.5px solid #e2e8f0; border-radius: 12px; font-size: 0.85rem; font-family: inherit; resize: vertical; }
.form-textarea:focus { border-color: #818cf8; outline: none; }
.btn { padding: 11px 24px; border-radius: 12px; font-size: 0.85rem; font-weight: 700; border: none; cursor: pointer; }
.btn-secondary { background: #f1f5f9; color: #475569; }
.btn-danger { background: #ef4444; color: white; }
.btn:disabled { opacity: 0.5; }

@media (max-width: 768px) {
  .dl-table-wrap { overflow-x: auto; }
  .dl-hero { padding: 24px 20px; }
}
</style>
