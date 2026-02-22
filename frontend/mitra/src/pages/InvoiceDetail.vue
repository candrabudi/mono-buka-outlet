<template>
  <div class="invd-page">
    <!-- Loading -->
    <div v-if="loading" class="invd-loading">
      <div class="invd-skel shimmer" style="height:200px;border-radius:16px"></div>
      <div class="invd-skel shimmer" style="height:400px;border-radius:16px;margin-top:20px"></div>
    </div>

    <template v-else-if="inv">
      <!-- Invoice Paper -->
      <div class="invd-paper">
        <!-- Header -->
        <div class="invd-header">
          <div class="invd-header-left">
            <div class="invd-brand">
              <div class="invd-brand-icon"><i class="ri-store-2-fill"></i></div>
              <div>
                <div class="invd-brand-name">BukaOutlet</div>
                <div class="invd-brand-sub">Platform Kemitraan Indonesia</div>
              </div>
            </div>
          </div>
          <div class="invd-header-right">
            <h1 class="invd-title">INVOICE</h1>
            <div class="invd-inv-number">{{ inv.invoice_number }}</div>
          </div>
        </div>

        <!-- Meta Bar -->
        <div class="invd-meta-bar">
          <div class="invd-meta-item">
            <span class="invd-meta-label">Tanggal Invoice</span>
            <span class="invd-meta-value">{{ formatDateFull(inv.created_at) }}</span>
          </div>
          <div class="invd-meta-item">
            <span class="invd-meta-label">Tipe</span>
            <span class="invd-meta-value"><span class="invd-type-badge" :class="'it-'+inv.type">{{ typeLabel(inv.type) }}</span></span>
          </div>
          <div class="invd-meta-item">
            <span class="invd-meta-label">Status</span>
            <span class="invd-meta-value">
              <span class="invd-status-badge" :class="'is-'+inv.status">
                <i :class="statusIcon(inv.status)"></i>
                {{ statusLabel(inv.status) }}
              </span>
            </span>
          </div>
          <div class="invd-meta-item" v-if="inv.expired_at && inv.status === 'PENDING'">
            <span class="invd-meta-label">Jatuh Tempo</span>
            <span class="invd-meta-value invd-due">{{ formatDateFull(inv.expired_at) }}</span>
          </div>
          <div class="invd-meta-item" v-if="inv.paid_at">
            <span class="invd-meta-label">Tanggal Bayar</span>
            <span class="invd-meta-value invd-paid-date">{{ formatDateFull(inv.paid_at) }}</span>
          </div>
        </div>

        <!-- Billing Info -->
        <div class="invd-billing">
          <div class="invd-billing-section">
            <h4 class="invd-billing-title">Penerbit Invoice</h4>
            <p class="invd-billing-line invd-billing-company">BukaOutlet Indonesia</p>
            <p class="invd-billing-line">Platform Kemitraan & Franchise</p>
            <p class="invd-billing-line">support@bukaoutlet.id</p>
          </div>
          <div class="invd-billing-section">
            <h4 class="invd-billing-title">Ditagihkan Kepada</h4>
            <p class="invd-billing-line invd-billing-company">{{ auth.userName }}</p>
            <p class="invd-billing-line">{{ auth.user?.email || '-' }}</p>
            <p class="invd-billing-line">{{ auth.user?.phone || '-' }}</p>
          </div>
        </div>

        <!-- Line Items Table -->
        <div class="invd-table-section">
          <table class="invd-table">
            <thead>
              <tr>
                <th style="width:50%">Deskripsi</th>
                <th style="width:15%;text-align:center">Qty</th>
                <th style="width:35%;text-align:right">Jumlah</th>
              </tr>
            </thead>
            <tbody>
              <tr>
                <td>
                  <div class="invd-item-name">{{ inv.description || typeLabel(inv.type) }}</div>
                  <div class="invd-item-desc">
                    Invoice {{ inv.invoice_number }} — {{ typeDesc(inv.type) }}
                  </div>
                </td>
                <td style="text-align:center">1</td>
                <td style="text-align:right" class="invd-item-amount">{{ fc(inv.amount) }}</td>
              </tr>
            </tbody>
          </table>
        </div>

        <!-- Summary -->
        <div class="invd-summary">
          <div class="invd-summary-row">
            <span>Subtotal</span>
            <span>{{ fc(inv.amount) }}</span>
          </div>
          <div class="invd-summary-row">
            <span>PPN (0%)</span>
            <span>Rp0</span>
          </div>
          <div class="invd-summary-row invd-total-row">
            <span>Total Tagihan</span>
            <span>{{ fc(inv.amount) }}</span>
          </div>
        </div>

        <!-- Payment Info -->
        <div class="invd-payment-info" v-if="inv.midtrans_payment_type || inv.midtrans_transaction_id">
          <h4 class="invd-payment-title">
            <i class="ri-bank-card-line"></i>
            Informasi Pembayaran
          </h4>
          <div class="invd-payment-grid">
            <div class="invd-payment-item" v-if="inv.midtrans_payment_type">
              <span class="invd-payment-label">Metode Pembayaran</span>
              <span class="invd-payment-value">{{ paymentMethodLabel(inv.midtrans_payment_type) }}</span>
            </div>
            <div class="invd-payment-item" v-if="inv.midtrans_transaction_id">
              <span class="invd-payment-label">ID Transaksi</span>
              <span class="invd-payment-value invd-mono">{{ inv.midtrans_transaction_id }}</span>
            </div>
            <div class="invd-payment-item" v-if="inv.midtrans_order_id">
              <span class="invd-payment-label">Order ID</span>
              <span class="invd-payment-value invd-mono">{{ inv.midtrans_order_id }}</span>
            </div>
            <div class="invd-payment-item" v-if="inv.midtrans_transaction_status">
              <span class="invd-payment-label">Status Transaksi</span>
              <span class="invd-payment-value">{{ inv.midtrans_transaction_status }}</span>
            </div>
          </div>
        </div>

        <!-- Proof of Payment -->
        <div class="invd-proof" v-if="inv.proof_url">
          <h4 class="invd-proof-title">
            <i class="ri-image-line"></i>
            Bukti Pembayaran
          </h4>
          <a :href="inv.proof_url" target="_blank" class="invd-proof-link">
            <img :src="inv.proof_url" alt="Bukti Pembayaran" class="invd-proof-img" />
          </a>
        </div>

        <!-- Footer -->
        <div class="invd-footer">
          <div class="invd-footer-note">
            <i class="ri-information-line"></i>
            <span>Invoice ini dibuat secara otomatis oleh sistem BukaOutlet dan berlaku sebagai bukti tagihan yang sah.</span>
          </div>
        </div>
      </div>

      <!-- Action Buttons -->
      <div class="invd-actions">
        <a v-if="inv.midtrans_redirect_url && inv.status === 'PENDING'" :href="inv.midtrans_redirect_url" target="_blank" class="invd-action-btn invd-btn-pay">
          <i class="ri-bank-card-line"></i>
          Bayar Sekarang
        </a>
        <router-link to="/invoices" class="invd-action-btn invd-btn-ghost">
          <i class="ri-arrow-left-s-line"></i>
          Kembali ke Daftar Invoice
        </router-link>
      </div>
    </template>

    <!-- Not Found -->
    <div v-else class="invd-not-found">
      <div class="invd-empty-circle"><i class="ri-file-text-line" style="font-size:32px;color:#94a3b8"></i></div>
      <p>Invoice tidak ditemukan</p>
      <router-link to="/invoices" class="invd-back-link">Kembali ke Daftar</router-link>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { invoiceApi } from '../services/api'
import { useToastStore } from '../stores/toast'
import { useAuthStore } from '../stores/auth'

const route = useRoute()
const toast = useToastStore()
const auth = useAuthStore()
const inv = ref(null)
const loading = ref(true)

function fc(v) { return v ? new Intl.NumberFormat('id-ID', { style: 'currency', currency: 'IDR', minimumFractionDigits: 0 }).format(v) : 'Rp0' }
function formatDateFull(d) { return d ? new Date(d).toLocaleDateString('id-ID', { day: 'numeric', month: 'long', year: 'numeric' }) : '-' }
function typeLabel(t) { return { DP: 'Down Payment', CICILAN: 'Cicilan', PELUNASAN: 'Pelunasan', INVOICE: 'Invoice' }[t] || t }
function typeDesc(t) {
  return {
    DP: 'Pembayaran uang muka (Down Payment) untuk kemitraan',
    CICILAN: 'Pembayaran cicilan sesuai jadwal yang telah disepakati',
    PELUNASAN: 'Pembayaran pelunasan untuk menyelesaikan tagihan kemitraan',
    INVOICE: 'Tagihan pembayaran kemitraan'
  }[t] || 'Tagihan pembayaran'
}
function statusLabel(s) { return { PAID: 'Lunas', PENDING: 'Menunggu Pembayaran', EXPIRED: 'Kadaluarsa', FAILED: 'Gagal', CANCELED: 'Dibatalkan' }[s] || s }
function statusIcon(s) { return { PAID: 'ri-checkbox-circle-fill', PENDING: 'ri-time-fill', EXPIRED: 'ri-alarm-warning-fill', FAILED: 'ri-close-circle-fill', CANCELED: 'ri-forbid-fill' }[s] || 'ri-question-fill' }
function paymentMethodLabel(m) {
  const map = { bank_transfer: 'Transfer Bank', credit_card: 'Kartu Kredit', gopay: 'GoPay', shopeepay: 'ShopeePay', qris: 'QRIS', cstore: 'Minimarket', echannel: 'Mandiri Bill', bca_va: 'BCA Virtual Account', bni_va: 'BNI Virtual Account', bri_va: 'BRI Virtual Account', permata_va: 'Permata VA' }
  return map[m] || m
}

onMounted(async () => {
  try {
    const { data } = await invoiceApi.getByID(route.params.id)
    inv.value = data.data
  } catch {
    toast.error('Gagal memuat detail invoice')
  } finally {
    loading.value = false
  }
})
</script>

<style scoped>
/* ═══ Page ═══ */
.invd-page { max-width: 800px; margin: 0 auto; }
.invd-loading { display: flex; flex-direction: column; gap: 20px; }
.invd-skel { background: #e8ecf1; border-radius: 16px; }
.shimmer { background: linear-gradient(90deg, #e8ecf1 25%, #f1f5f9 50%, #e8ecf1 75%); background-size: 200% 100%; animation: shimmer 1.5s infinite; }
@keyframes shimmer { 0% { background-position: 200% 0; } 100% { background-position: -200% 0; } }

/* ═══ Paper ═══ */
.invd-paper { background: #fff; border-radius: 20px; border: 1px solid #e8ecf1; overflow: hidden; box-shadow: 0 4px 24px rgba(0,0,0,.04); }

/* ═══ Header ═══ */
.invd-header { display: flex; justify-content: space-between; align-items: flex-start; padding: 36px 36px 28px; border-bottom: 2px solid #f1f5f9; }
.invd-brand { display: flex; align-items: center; gap: 12px; }
.invd-brand-icon { width: 44px; height: 44px; border-radius: 12px; background: linear-gradient(135deg, #6366f1, #8b5cf6); color: #fff; display: flex; align-items: center; justify-content: center; font-size: 22px; }
.invd-brand-name { font-size: 1.1rem; font-weight: 800; color: #0f172a; }
.invd-brand-sub { font-size: .72rem; color: #94a3b8; margin-top: 1px; }
.invd-header-right { text-align: right; }
.invd-title { font-size: 1.8rem; font-weight: 900; color: #6366f1; margin: 0; letter-spacing: .04em; }
.invd-inv-number { font-size: .85rem; font-weight: 700; color: #64748b; margin-top: 4px; font-family: 'JetBrains Mono', 'Fira Code', monospace; }

/* ═══ Meta Bar ═══ */
.invd-meta-bar { display: flex; flex-wrap: wrap; gap: 6px 20px; padding: 20px 36px; background: #f8fafc; border-bottom: 1px solid #f1f5f9; }
.invd-meta-item { display: flex; flex-direction: column; gap: 2px; min-width: 140px; }
.invd-meta-label { font-size: .68rem; font-weight: 600; color: #94a3b8; text-transform: uppercase; letter-spacing: .06em; }
.invd-meta-value { font-size: .82rem; font-weight: 600; color: #1e293b; }
.invd-due { color: #d97706; }
.invd-paid-date { color: #16a34a; }

.invd-type-badge { font-size: .68rem; font-weight: 700; padding: 3px 8px; border-radius: 5px; text-transform: uppercase; letter-spacing: .03em; }
.it-DP { background: #e0f2fe; color: #0284c7; }
.it-CICILAN { background: #ede9fe; color: #7c3aed; }
.it-PELUNASAN { background: #dcfce7; color: #16a34a; }
.it-INVOICE { background: #fef3c7; color: #d97706; }

.invd-status-badge { display: inline-flex; align-items: center; gap: 4px; font-size: .75rem; font-weight: 700; padding: 4px 10px; border-radius: 6px; }
.is-PAID { background: #dcfce7; color: #16a34a; }
.is-PENDING { background: #fef3c7; color: #d97706; }
.is-EXPIRED, .is-FAILED, .is-CANCELED { background: #f1f5f9; color: #64748b; }

/* ═══ Billing ═══ */
.invd-billing { display: grid; grid-template-columns: 1fr 1fr; gap: 32px; padding: 28px 36px; }
.invd-billing-title { font-size: .7rem; font-weight: 700; color: #94a3b8; text-transform: uppercase; letter-spacing: .06em; margin: 0 0 10px; }
.invd-billing-line { font-size: .82rem; color: #475569; margin: 0 0 3px; line-height: 1.5; }
.invd-billing-company { font-weight: 700; color: #0f172a; font-size: .88rem; }

/* ═══ Table ═══ */
.invd-table-section { padding: 0 36px; }
.invd-table { width: 100%; border-collapse: collapse; }
.invd-table thead { background: #f8fafc; }
.invd-table th { padding: 12px 16px; font-size: .7rem; font-weight: 700; color: #64748b; text-transform: uppercase; letter-spacing: .06em; text-align: left; border-top: 1px solid #e8ecf1; border-bottom: 1px solid #e8ecf1; }
.invd-table td { padding: 18px 16px; font-size: .85rem; color: #1e293b; border-bottom: 1px solid #f1f5f9; }
.invd-item-name { font-weight: 700; color: #0f172a; margin-bottom: 3px; }
.invd-item-desc { font-size: .78rem; color: #94a3b8; line-height: 1.4; }
.invd-item-amount { font-weight: 800; color: #0f172a; white-space: nowrap; }

/* ═══ Summary ═══ */
.invd-summary { margin: 0 36px 28px; margin-left: auto; max-width: 320px; padding-top: 4px; }
.invd-summary-row { display: flex; justify-content: space-between; align-items: center; padding: 8px 0; font-size: .85rem; color: #64748b; }
.invd-total-row { border-top: 2px solid #0f172a; margin-top: 4px; padding-top: 12px; font-size: 1.05rem; font-weight: 800; color: #0f172a; }

/* ═══ Payment Info ═══ */
.invd-payment-info { margin: 0 36px 28px; padding: 20px; background: #f8fafc; border-radius: 12px; border: 1px solid #e8ecf1; }
.invd-payment-title { font-size: .82rem; font-weight: 700; color: #0f172a; margin: 0 0 14px; display: flex; align-items: center; gap: 6px; }
.invd-payment-title i { color: #6366f1; }
.invd-payment-grid { display: grid; grid-template-columns: 1fr 1fr; gap: 14px; }
.invd-payment-item { display: flex; flex-direction: column; gap: 2px; }
.invd-payment-label { font-size: .7rem; color: #94a3b8; font-weight: 600; text-transform: uppercase; letter-spacing: .04em; }
.invd-payment-value { font-size: .82rem; color: #1e293b; font-weight: 600; }
.invd-mono { font-family: 'JetBrains Mono', 'Fira Code', monospace; font-size: .78rem; }

/* ═══ Proof ═══ */
.invd-proof { margin: 0 36px 28px; }
.invd-proof-title { font-size: .82rem; font-weight: 700; color: #0f172a; margin: 0 0 12px; display: flex; align-items: center; gap: 6px; }
.invd-proof-title i { color: #6366f1; }
.invd-proof-link { display: block; border-radius: 12px; overflow: hidden; border: 1px solid #e8ecf1; transition: box-shadow .2s; }
.invd-proof-link:hover { box-shadow: 0 4px 14px rgba(0,0,0,.08); }
.invd-proof-img { width: 100%; max-height: 300px; object-fit: contain; background: #f8fafc; display: block; }

/* ═══ Footer ═══ */
.invd-footer { padding: 20px 36px; background: #f8fafc; border-top: 1px solid #f1f5f9; }
.invd-footer-note { display: flex; align-items: flex-start; gap: 8px; font-size: .75rem; color: #94a3b8; line-height: 1.5; }
.invd-footer-note i { flex-shrink: 0; margin-top: 1px; color: #cbd5e1; }

/* ═══ Actions ═══ */
.invd-actions { display: flex; gap: 10px; margin-top: 20px; justify-content: center; }
.invd-action-btn { display: inline-flex; align-items: center; gap: 6px; padding: 12px 24px; border-radius: 12px; font-size: .85rem; font-weight: 700; text-decoration: none; transition: all .2s; }
.invd-btn-pay { background: linear-gradient(135deg, #6366f1, #8b5cf6); color: #fff; box-shadow: 0 2px 12px rgba(99,102,241,.3); }
.invd-btn-pay:hover { box-shadow: 0 4px 18px rgba(99,102,241,.4); transform: translateY(-1px); }
.invd-btn-ghost { background: #fff; color: #475569; border: 1px solid #e8ecf1; }
.invd-btn-ghost:hover { background: #f8fafc; }

/* ═══ Not Found ═══ */
.invd-not-found { text-align: center; padding: 80px 20px; background: #fff; border-radius: 16px; border: 1px solid #e8ecf1; }
.invd-empty-circle { width: 72px; height: 72px; border-radius: 50%; background: linear-gradient(135deg, #f1f5f9, #e2e8f0); display: flex; align-items: center; justify-content: center; margin: 0 auto 16px; }
.invd-not-found p { color: #94a3b8; font-size: .88rem; margin: 0 0 16px; }
.invd-back-link { font-size: .82rem; font-weight: 600; color: #6366f1; text-decoration: none; }
.invd-back-link:hover { text-decoration: underline; }

@media (max-width: 768px) {
  .invd-header { flex-direction: column; gap: 16px; padding: 24px 20px 20px; }
  .invd-header-right { text-align: left; }
  .invd-meta-bar { padding: 16px 20px; }
  .invd-billing { grid-template-columns: 1fr; gap: 20px; padding: 20px; }
  .invd-table-section { padding: 0 12px; overflow-x: auto; }
  .invd-summary { margin: 0 20px 20px; max-width: 100%; }
  .invd-payment-info { margin: 0 20px 20px; }
  .invd-payment-grid { grid-template-columns: 1fr; }
  .invd-proof { margin: 0 20px 20px; }
  .invd-footer { padding: 16px 20px; }
  .invd-actions { flex-direction: column; }
}
</style>
