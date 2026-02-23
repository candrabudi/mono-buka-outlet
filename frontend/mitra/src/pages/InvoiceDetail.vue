<template>
  <div class="id-page">
    <!-- Loading -->
    <div v-if="loading" class="id-loading">
      <div class="id-skel-hero shimmer"></div>
      <div class="id-skel-row"><div class="id-skel-card shimmer"></div><div class="id-skel-card shimmer"></div></div>
    </div>

    <template v-else-if="inv">
      <!-- Hero -->
      <div class="id-hero">
        <button class="id-back" @click="$router.push('/invoices')">
          <i class="ri-arrow-left-line"></i>
          Kembali
        </button>
        <div class="id-hero-content">
          <div class="id-hero-icon">
            <i class="ri-receipt-line"></i>
          </div>
          <div>
            <div class="id-hero-badges">
              <span class="id-badge" :class="'st-'+inv.status">
                <span class="id-badge-dot"></span>
                {{ statusLabel(inv.status) }}
              </span>
              <span class="id-type-badge" :class="'tp-'+inv.type">{{ typeLabel(inv.type) }}</span>
            </div>
            <h1 class="id-hero-title">{{ inv.invoice_number }}</h1>
            <p class="id-hero-sub">
              <i class="ri-calendar-line"></i>
              Diterbitkan {{ formatDate(inv.created_at) }}
            </p>
          </div>
        </div>
        <div class="id-stats">
          <div class="id-stat">
            <span class="id-stat-label">Total Tagihan</span>
            <span class="id-stat-val">{{ fc(inv.amount) }}</span>
          </div>
          <div class="id-stat" v-if="inv.paid_at">
            <span class="id-stat-label">Tanggal Bayar</span>
            <span class="id-stat-val">{{ formatDate(inv.paid_at) }}</span>
          </div>
          <div class="id-stat" v-else-if="inv.expired_at && inv.status === 'PENDING'">
            <span class="id-stat-label">Jatuh Tempo</span>
            <span class="id-stat-val" style="color:#fbbf24">{{ formatDate(inv.expired_at) }}</span>
          </div>
          <div class="id-stat" v-if="inv.midtrans_payment_type">
            <span class="id-stat-label">Metode Pembayaran</span>
            <span class="id-stat-val">{{ paymentMethodLabel(inv.midtrans_payment_type) }}</span>
          </div>
        </div>
      </div>

      <!-- Content -->
      <div class="id-content">
        <!-- Left: Invoice Detail -->
        <div class="id-main">
          <!-- Billing Info -->
          <div class="id-section">
            <h3 class="id-section-title">
              <i class="ri-user-line"></i>
              Informasi Penagihan
            </h3>
            <div class="id-billing-grid">
              <div class="id-billing-card">
                <div class="id-billing-dot" style="background:#6366f1"></div>
                <div>
                  <div class="id-billing-heading">Penerbit</div>
                  <div class="id-billing-name">BukaOutlet Indonesia</div>
                  <div class="id-billing-sub">Platform Kemitraan & Franchise</div>
                  <div class="id-billing-sub">support@bukaoutlet.id</div>
                </div>
              </div>
              <div class="id-billing-card">
                <div class="id-billing-dot" style="background:#f59e0b"></div>
                <div>
                  <div class="id-billing-heading">Ditagihkan Kepada</div>
                  <div class="id-billing-name">{{ auth.userName }}</div>
                  <div class="id-billing-sub">{{ auth.user?.email || '-' }}</div>
                  <div class="id-billing-sub">{{ auth.user?.phone || '-' }}</div>
                </div>
              </div>
            </div>
          </div>

          <!-- Line Items -->
          <div class="id-section">
            <h3 class="id-section-title">
              <i class="ri-file-list-3-line"></i>
              Rincian Tagihan
            </h3>
            <div class="id-table-wrap">
              <table class="id-table">
                <thead>
                  <tr>
                    <th>Deskripsi</th>
                    <th style="text-align:center;width:80px">Qty</th>
                    <th style="text-align:right;width:160px">Jumlah</th>
                  </tr>
                </thead>
                <tbody>
                  <tr>
                    <td>
                      <div class="id-item-title">{{ inv.description || typeLabel(inv.type) }}</div>
                      <div class="id-item-sub">{{ typeDesc(inv.type) }}</div>
                    </td>
                    <td style="text-align:center;color:#64748b">1</td>
                    <td style="text-align:right" class="id-item-amount">{{ fc(inv.amount) }}</td>
                  </tr>
                </tbody>
                <tfoot>
                  <tr class="id-total-row">
                    <td colspan="2" style="text-align:right">Total Tagihan</td>
                    <td style="text-align:right">{{ fc(inv.amount) }}</td>
                  </tr>
                </tfoot>
              </table>
            </div>
          </div>

          <!-- Transaction Details -->
          <div class="id-section" v-if="inv.midtrans_transaction_id || inv.midtrans_order_id">
            <h3 class="id-section-title">
              <i class="ri-exchange-funds-line"></i>
              Detail Transaksi
            </h3>
            <div class="id-info-grid">
              <div class="id-info-row" v-if="inv.midtrans_order_id">
                <span class="id-info-label">Order ID</span>
                <span class="id-info-value id-mono">{{ inv.midtrans_order_id }}</span>
              </div>
              <div class="id-info-row" v-if="inv.midtrans_transaction_id">
                <span class="id-info-label">ID Transaksi</span>
                <span class="id-info-value id-mono">{{ inv.midtrans_transaction_id }}</span>
              </div>
              <div class="id-info-row" v-if="inv.midtrans_payment_type">
                <span class="id-info-label">Metode Pembayaran</span>
                <span class="id-info-value">{{ paymentMethodLabel(inv.midtrans_payment_type) }}</span>
              </div>
              <div class="id-info-row" v-if="inv.midtrans_transaction_status">
                <span class="id-info-label">Status Transaksi</span>
                <span class="id-info-value">{{ inv.midtrans_transaction_status }}</span>
              </div>
            </div>
          </div>

          <!-- Supported Payment Methods Info -->
          <div class="id-section" v-if="inv.status === 'PENDING'">
            <h3 class="id-section-title">
              <i class="ri-bank-card-2-line"></i>
              Metode Pembayaran yang Didukung
            </h3>
            <div class="id-payment-methods">
              <div class="id-pm-group">
                <div class="id-pm-label">💳 Kartu Kredit / Debit</div>
                <div class="id-pm-desc">Visa, Mastercard, JCB, Amex — dengan 3D Secure</div>
              </div>
              <div class="id-pm-group">
                <div class="id-pm-label">🏦 Transfer Bank (Virtual Account)</div>
                <div class="id-pm-desc">BCA, BNI, BRI, Mandiri, Permata</div>
              </div>
              <div class="id-pm-group">
                <div class="id-pm-label">📱 E-Wallet & QRIS</div>
                <div class="id-pm-desc">GoPay, ShopeePay, DANA, QRIS</div>
              </div>
              <div class="id-pm-group">
                <div class="id-pm-label">🛒 PayLater & Minimarket</div>
                <div class="id-pm-desc">Kredivo, Akulaku, Indomaret, Alfamart</div>
              </div>
            </div>
          </div>
        </div>

        <!-- Right: Sidebar -->
        <div class="id-sidebar">
          <!-- Pay CTA — Snap Popup -->
          <div class="id-pay-card" v-if="inv.midtrans_snap_token && inv.status === 'PENDING'">
            <div class="id-pay-icon-row">
              <div class="id-pay-icon"><i class="ri-secure-payment-line"></i></div>
              <div>
                <div class="id-pay-heading">Selesaikan Pembayaran</div>
                <p class="id-pay-desc">Pilih metode pembayaran favorit Anda — Credit Card, PayLater, Transfer Bank, E-Wallet, QRIS, dan lainnya.</p>
              </div>
            </div>
            <button @click="openSnapPopup" class="id-pay-btn" :disabled="snapLoading">
              <i class="ri-bank-card-line" v-if="!snapLoading"></i>
              <span v-if="snapLoading" class="id-pay-spinner"></span>
              {{ snapLoading ? 'Memuat...' : 'Bayar Sekarang' }}
            </button>
            <!-- Fallback link -->
            <a v-if="inv.midtrans_redirect_url" :href="inv.midtrans_redirect_url" target="_blank" class="id-fallback-link">
              <i class="ri-external-link-line"></i>
              Atau buka halaman pembayaran
            </a>
          </div>

          <!-- Success notice -->
          <div class="id-success-card" v-if="paymentSuccess">
            <div class="id-success-icon"><i class="ri-checkbox-circle-fill"></i></div>
            <div class="id-success-heading">Pembayaran Berhasil!</div>
            <p class="id-success-desc">Terima kasih, pembayaran Anda sedang diproses. Status akan diperbarui dalam beberapa saat.</p>
          </div>

          <!-- Proof -->
          <div class="id-sidebar-card" v-if="inv.proof_url">
            <h3 class="id-sidebar-title">
              <i class="ri-image-line"></i>
              Bukti Pembayaran
            </h3>
            <a :href="inv.proof_url" target="_blank" class="id-proof-btn">
              <i class="ri-external-link-line"></i>
              Lihat Bukti
            </a>
          </div>

          <!-- Invoice Note -->
          <div class="id-sidebar-card">
            <h3 class="id-sidebar-title">
              <i class="ri-information-line"></i>
              Catatan
            </h3>
            <p class="id-sidebar-note">Invoice ini dibuat secara otomatis oleh sistem BukaOutlet dan berlaku sebagai bukti tagihan yang sah.</p>
          </div>

          <!-- Quick Actions -->
          <div class="id-quick-actions">
            <router-link to="/invoices" class="id-action-btn">
              <i class="ri-arrow-left-s-line"></i>
              Kembali ke Daftar
            </router-link>
          </div>
        </div>
      </div>
    </template>

    <!-- Not Found -->
    <div v-else class="id-empty">
      <div class="id-empty-icon"><i class="ri-file-text-line"></i></div>
      <p>Invoice tidak ditemukan</p>
      <router-link to="/invoices" class="id-empty-link">Kembali ke Daftar</router-link>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { invoiceApi, midtransApi } from '../services/api'
import { useToastStore } from '../stores/toast'
import { useAuthStore } from '../stores/auth'

const route = useRoute()
const toast = useToastStore()
const auth = useAuthStore()
const inv = ref(null)
const loading = ref(true)
const snapLoading = ref(false)
const paymentSuccess = ref(false)
let snapLoaded = false

function fc(v) { return v ? new Intl.NumberFormat('id-ID', { style: 'currency', currency: 'IDR', minimumFractionDigits: 0 }).format(v) : 'Rp0' }
function formatDate(d) { return d ? new Date(d).toLocaleDateString('id-ID', { day: 'numeric', month: 'long', year: 'numeric' }) : '-' }
function typeLabel(t) { return { DP: 'Down Payment', CICILAN: 'Cicilan', PELUNASAN: 'Pelunasan', INVOICE: 'Invoice' }[t] || t }
function typeDesc(t) {
  return {
    DP: 'Pembayaran uang muka (Down Payment) untuk proses aktivasi kemitraan outlet',
    CICILAN: 'Pembayaran cicilan sesuai jadwal yang telah disepakati bersama',
    PELUNASAN: 'Pembayaran pelunasan akhir untuk menyelesaikan seluruh tagihan kemitraan',
    INVOICE: 'Tagihan pembayaran terkait kemitraan outlet'
  }[t] || 'Tagihan pembayaran kemitraan'
}
function statusLabel(s) { return { PAID: 'Lunas', PENDING: 'Menunggu Pembayaran', EXPIRED: 'Kadaluarsa', FAILED: 'Gagal', CANCELED: 'Dibatalkan' }[s] || s }
function paymentMethodLabel(m) {
  const map = {
    bank_transfer: 'Transfer Bank', credit_card: 'Kartu Kredit', gopay: 'GoPay',
    shopeepay: 'ShopeePay', qris: 'QRIS', cstore: 'Minimarket', echannel: 'Mandiri Bill',
    bca_va: 'BCA VA', bni_va: 'BNI VA', bri_va: 'BRI VA', permata_va: 'Permata VA',
    dana: 'DANA', kredivo: 'Kredivo', akulaku: 'Akulaku'
  }
  return map[m] || m
}

// Load Midtrans Snap.js dynamically
async function loadSnapJS() {
  if (snapLoaded && window.snap) return true

  try {
    const { data } = await midtransApi.getClientKey()
    const clientKey = data.client_key
    const snapUrl = data.snap_url

    if (!clientKey || !snapUrl) {
      console.warn('Midtrans client key or snap URL missing')
      return false
    }

    return new Promise((resolve) => {
      // Check if script already exists
      const existing = document.querySelector('script[data-midtrans-snap]')
      if (existing) {
        snapLoaded = true
        resolve(true)
        return
      }

      const script = document.createElement('script')
      script.src = snapUrl
      script.setAttribute('data-client-key', clientKey)
      script.setAttribute('data-midtrans-snap', 'true')
      script.onload = () => {
        snapLoaded = true
        resolve(true)
      }
      script.onerror = () => {
        console.error('Failed to load Midtrans Snap.js')
        resolve(false)
      }
      document.head.appendChild(script)
    })
  } catch (err) {
    console.error('Failed to get Midtrans client key:', err)
    return false
  }
}

// Open Midtrans Snap popup
async function openSnapPopup() {
  if (!inv.value?.midtrans_snap_token) {
    toast.error('Token pembayaran tidak tersedia')
    return
  }

  snapLoading.value = true

  const loaded = await loadSnapJS()
  if (!loaded || !window.snap) {
    snapLoading.value = false
    // Fallback: open redirect URL
    if (inv.value.midtrans_redirect_url) {
      window.open(inv.value.midtrans_redirect_url, '_blank')
    } else {
      toast.error('Gagal memuat Midtrans. Coba lagi nanti.')
    }
    return
  }

  snapLoading.value = false

  window.snap.pay(inv.value.midtrans_snap_token, {
    onSuccess: (result) => {
      console.log('Payment success:', result)
      paymentSuccess.value = true
      toast.success('Pembayaran berhasil! Terima kasih.')
      // Refresh invoice data
      refreshInvoice()
    },
    onPending: (result) => {
      console.log('Payment pending:', result)
      toast.info('Pembayaran menunggu konfirmasi. Silakan selesaikan pembayaran.')
      refreshInvoice()
    },
    onError: (result) => {
      console.error('Payment error:', result)
      toast.error('Pembayaran gagal. Silakan coba lagi.')
    },
    onClose: () => {
      console.log('Snap popup closed')
      // Refresh just in case user completed payment before closing
      refreshInvoice()
    }
  })
}

async function refreshInvoice() {
  try {
    const { data } = await invoiceApi.getByID(route.params.id)
    inv.value = data.data
  } catch { /* silent */ }
}

onMounted(async () => {
  try {
    const { data } = await invoiceApi.getByID(route.params.id)
    inv.value = data.data
    // Preload Snap.js if invoice is pending
    if (inv.value?.status === 'PENDING' && inv.value?.midtrans_snap_token) {
      loadSnapJS()
    }
  } catch {
    toast.error('Gagal memuat detail invoice')
  } finally {
    loading.value = false
  }
})
</script>

<style scoped>
/* ═══ Hero — same pattern as ApplicationDetail ═══ */
.id-hero { background: linear-gradient(135deg, #0f0c29 0%, #302b63 50%, #24243e 100%); border-radius: 16px; padding: 28px 36px 24px; margin-bottom: 24px; box-shadow: 0 4px 24px rgba(15,12,41,.2); }
.id-back { display: inline-flex; align-items: center; gap: 5px; font-size: .78rem; font-weight: 600; color: rgba(255,255,255,.45); background: rgba(255,255,255,.06); border: 1px solid rgba(255,255,255,.1); border-radius: 8px; padding: 6px 14px; cursor: pointer; margin-bottom: 18px; transition: all .15s; backdrop-filter: blur(4px); }
.id-back:hover { color: #fff; border-color: rgba(255,255,255,.25); background: rgba(255,255,255,.1); }

.id-hero-content { display: flex; align-items: center; gap: 16px; margin-bottom: 20px; }
.id-hero-icon { width: 56px; height: 56px; border-radius: 14px; flex-shrink: 0; border: 2px solid rgba(255,255,255,.12); background: rgba(255,255,255,.06); display: flex; align-items: center; justify-content: center; font-size: 24px; color: rgba(255,255,255,.5); }
.id-hero-badges { display: flex; gap: 8px; margin-bottom: 6px; flex-wrap: wrap; }
.id-hero-title { font-size: 1.5rem; font-weight: 800; color: #fff; margin: 0; }
.id-hero-sub { font-size: .82rem; color: rgba(255,255,255,.4); margin: 6px 0 0; display: flex; align-items: center; gap: 5px; }

.id-badge { font-size: .68rem; font-weight: 700; padding: 4px 12px; border-radius: 6px; text-transform: uppercase; letter-spacing: .03em; display: inline-flex; align-items: center; gap: 5px; }
.id-badge-dot { width: 6px; height: 6px; border-radius: 50%; }
.st-PAID { background: rgba(34,197,94,.15); color: #4ade80; }
.st-PAID .id-badge-dot { background: #22c55e; }
.st-PENDING { background: rgba(245,158,11,.15); color: #fbbf24; }
.st-PENDING .id-badge-dot { background: #f59e0b; }
.st-EXPIRED, .st-FAILED, .st-CANCELED { background: rgba(148,163,184,.15); color: #94a3b8; }
.st-EXPIRED .id-badge-dot, .st-FAILED .id-badge-dot, .st-CANCELED .id-badge-dot { background: #94a3b8; }

.id-type-badge { font-size: .68rem; font-weight: 700; padding: 4px 10px; border-radius: 6px; text-transform: uppercase; letter-spacing: .03em; }
.tp-DP { background: rgba(56,189,248,.12); color: #38bdf8; }
.tp-CICILAN { background: rgba(167,139,250,.12); color: #a78bfa; }
.tp-PELUNASAN { background: rgba(52,211,153,.12); color: #34d399; }
.tp-INVOICE { background: rgba(251,191,36,.12); color: #fbbf24; }

.id-stats { display: flex; gap: 32px; flex-wrap: wrap; padding-top: 18px; border-top: 1px solid rgba(255,255,255,.08); }
.id-stat { display: flex; flex-direction: column; gap: 2px; }
.id-stat-label { font-size: .68rem; color: rgba(255,255,255,.35); text-transform: uppercase; letter-spacing: .05em; }
.id-stat-val { font-size: .95rem; font-weight: 800; color: #fff; }

/* ═══ Content Grid ═══ */
.id-content { display: grid; grid-template-columns: 1fr 340px; gap: 24px; align-items: start; }
.id-main { display: flex; flex-direction: column; gap: 20px; min-width: 0; }
.id-sidebar { display: flex; flex-direction: column; gap: 16px; }

/* ═══ Section Cards ═══ */
.id-section { background: #fff; border-radius: 14px; border: 1px solid #e8ecf1; padding: 24px 28px; }
.id-section-title { font-size: .88rem; font-weight: 700; color: #0f172a; margin: 0 0 18px; display: flex; align-items: center; gap: 8px; }
.id-section-title i { color: #6366f1; font-size: 16px; }

/* Billing */
.id-billing-grid { display: grid; grid-template-columns: 1fr 1fr; gap: 20px; }
.id-billing-card { display: flex; gap: 12px; padding: 16px; background: #f8fafc; border-radius: 12px; border: 1px solid #f1f5f9; }
.id-billing-dot { width: 4px; border-radius: 2px; flex-shrink: 0; }
.id-billing-heading { font-size: .66rem; font-weight: 700; color: #94a3b8; text-transform: uppercase; letter-spacing: .06em; margin-bottom: 6px; }
.id-billing-name { font-size: .88rem; font-weight: 700; color: #0f172a; margin-bottom: 2px; }
.id-billing-sub { font-size: .78rem; color: #64748b; line-height: 1.5; }

/* Table */
.id-table-wrap { overflow-x: auto; }
.id-table { width: 100%; border-collapse: collapse; }
.id-table thead { background: #f8fafc; }
.id-table th { padding: 11px 16px; font-size: .68rem; font-weight: 700; color: #64748b; text-transform: uppercase; letter-spacing: .06em; text-align: left; border: 1px solid #e8ecf1; border-left: none; border-right: none; }
.id-table td { padding: 18px 16px; font-size: .85rem; color: #1e293b; border-bottom: 1px solid #f1f5f9; }
.id-item-title { font-weight: 700; color: #0f172a; margin-bottom: 4px; }
.id-item-sub { font-size: .75rem; color: #94a3b8; line-height: 1.5; max-width: 400px; }
.id-item-amount { font-weight: 800; color: #0f172a; white-space: nowrap; }

/* Total Row */
.id-total-row td { padding: 14px 16px; font-size: .95rem; font-weight: 800; color: #0f172a; border-top: 2px solid #0f172a; border-bottom: none; }

/* Transaction Info */
.id-info-grid { display: flex; flex-direction: column; gap: 0; }
.id-info-row { display: flex; justify-content: space-between; align-items: center; padding: 12px 0; border-bottom: 1px solid #f1f5f9; }
.id-info-row:last-child { border-bottom: none; }
.id-info-label { font-size: .78rem; color: #94a3b8; font-weight: 500; }
.id-info-value { font-size: .82rem; color: #1e293b; font-weight: 600; }
.id-mono { font-family: 'JetBrains Mono', 'Fira Code', monospace; font-size: .75rem; color: #475569; }

/* ═══ Payment Methods Grid ═══ */
.id-payment-methods { display: grid; grid-template-columns: 1fr 1fr; gap: 12px; }
.id-pm-group { padding: 14px 16px; background: #f8fafc; border-radius: 10px; border: 1px solid #f1f5f9; transition: all .15s; }
.id-pm-group:hover { border-color: #e0e7ff; background: #fafbff; }
.id-pm-label { font-size: .82rem; font-weight: 700; color: #0f172a; margin-bottom: 4px; }
.id-pm-desc { font-size: .72rem; color: #94a3b8; line-height: 1.4; }

/* ═══ Sidebar Cards ═══ */
.id-pay-card { background: linear-gradient(135deg, #312e81, #4338ca); border-radius: 14px; padding: 22px; color: #fff; }
.id-pay-icon-row { display: flex; gap: 12px; margin-bottom: 16px; }
.id-pay-icon { width: 40px; height: 40px; border-radius: 10px; background: rgba(255,255,255,.1); display: flex; align-items: center; justify-content: center; font-size: 20px; color: rgba(255,255,255,.7); flex-shrink: 0; }
.id-pay-heading { font-size: .88rem; font-weight: 700; }
.id-pay-desc { font-size: .75rem; color: rgba(255,255,255,.5); margin: 4px 0 0; line-height: 1.5; }
.id-pay-btn { display: flex; align-items: center; justify-content: center; gap: 8px; width: 100%; padding: 13px; border-radius: 10px; font-size: .88rem; font-weight: 700; color: #312e81; background: #fff; border: none; cursor: pointer; box-shadow: 0 4px 14px rgba(0,0,0,.12); transition: all .2s; font-family: inherit; }
.id-pay-btn:hover:not(:disabled) { transform: translateY(-1px); box-shadow: 0 6px 20px rgba(0,0,0,.18); }
.id-pay-btn:disabled { opacity: .7; cursor: not-allowed; transform: none; }

.id-pay-spinner { width: 16px; height: 16px; border: 2px solid #c7d2fe; border-top-color: #4338ca; border-radius: 50%; animation: spin .6s linear infinite; }
@keyframes spin { to { transform: rotate(360deg); } }

.id-fallback-link { display: flex; align-items: center; justify-content: center; gap: 5px; margin-top: 12px; font-size: .72rem; color: rgba(255,255,255,.4); text-decoration: none; transition: color .15s; }
.id-fallback-link:hover { color: rgba(255,255,255,.7); }

/* Success Card */
.id-success-card { background: linear-gradient(135deg, #065f46, #047857); border-radius: 14px; padding: 22px; color: #fff; text-align: center; animation: fadeInUp .4s ease; }
.id-success-icon { font-size: 36px; color: #34d399; margin-bottom: 8px; }
.id-success-heading { font-size: .95rem; font-weight: 700; margin-bottom: 6px; }
.id-success-desc { font-size: .78rem; color: rgba(255,255,255,.6); margin: 0; line-height: 1.5; }
@keyframes fadeInUp { from { opacity: 0; transform: translateY(12px); } to { opacity: 1; transform: translateY(0); } }

.id-sidebar-card { background: #fff; border-radius: 14px; border: 1px solid #e8ecf1; padding: 22px; }
.id-sidebar-title { font-size: .85rem; font-weight: 700; color: #0f172a; margin: 0 0 12px; display: flex; align-items: center; gap: 6px; }
.id-sidebar-title i { color: #6366f1; }
.id-sidebar-note { font-size: .78rem; color: #64748b; line-height: 1.6; margin: 0; }

.id-proof-btn { display: inline-flex; align-items: center; gap: 6px; font-size: .82rem; font-weight: 600; color: #6366f1; text-decoration: none; padding: 10px 18px; border-radius: 10px; background: #eef2ff; transition: all .15s; }
.id-proof-btn:hover { background: #e0e7ff; }

/* Quick Actions */
.id-quick-actions { display: flex; flex-direction: column; gap: 8px; }
.id-action-btn { display: flex; align-items: center; gap: 6px; padding: 12px 18px; border-radius: 10px; font-size: .82rem; font-weight: 600; color: #475569; text-decoration: none; background: #fff; border: 1px solid #e8ecf1; transition: all .15s; }
.id-action-btn:hover { background: #f8fafc; border-color: #cbd5e1; }

/* ═══ Loading & Empty ═══ */
.id-loading { display: flex; flex-direction: column; gap: 0; }
.id-skel-hero { height: 190px; border-radius: 16px; margin-bottom: 24px; }
.id-skel-row { display: grid; grid-template-columns: 1fr 340px; gap: 24px; }
.id-skel-card { height: 300px; border-radius: 14px; }
.shimmer { background: linear-gradient(90deg, #e8ecf1 25%, #f1f5f9 50%, #e8ecf1 75%); background-size: 200% 100%; animation: shimmer 1.5s infinite; }
@keyframes shimmer { 0% { background-position: 200% 0; } 100% { background-position: -200% 0; } }

.id-empty { text-align: center; padding: 80px 20px; background: #fff; border-radius: 14px; border: 1px solid #e8ecf1; }
.id-empty-icon { width: 72px; height: 72px; border-radius: 50%; background: linear-gradient(135deg, #f1f5f9, #e2e8f0); display: flex; align-items: center; justify-content: center; margin: 0 auto 16px; font-size: 28px; color: #94a3b8; }
.id-empty p { color: #94a3b8; font-size: .85rem; margin: 0 0 16px; }
.id-empty-link { font-size: .82rem; font-weight: 600; color: #6366f1; text-decoration: none; }
.id-empty-link:hover { text-decoration: underline; }

/* ═══ Responsive ═══ */
@media (max-width: 900px) {
  .id-hero { padding: 24px 20px 18px; }
  .id-content { grid-template-columns: 1fr; }
  .id-skel-row { grid-template-columns: 1fr; }
  .id-billing-grid { grid-template-columns: 1fr; }
  .id-payment-methods { grid-template-columns: 1fr; }
  .id-section { padding: 20px; }
  .id-summary { width: 100%; }
}
</style>
