<template>
  <div class="inv-page">
    <div class="inv-hero">
      <div class="inv-hero-top">
        <div>
          <h1 class="inv-hero-title">Invoice</h1>
          <p class="inv-hero-sub">Riwayat tagihan dan status pembayaran Anda</p>
        </div>
      </div>
      <div class="inv-stats">
        <div class="inv-stat"><span class="inv-stat-dot" style="background:#818cf8"></span><span class="inv-stat-label">Total</span><span class="inv-stat-val">{{ invoices.length }}</span></div>
        <div class="inv-stat"><span class="inv-stat-dot" style="background:#22c55e"></span><span class="inv-stat-label">Lunas</span><span class="inv-stat-val">{{ invoices.filter(i=>i.status==='PAID').length }}</span></div>
        <div class="inv-stat"><span class="inv-stat-dot" style="background:#f59e0b"></span><span class="inv-stat-label">Pending</span><span class="inv-stat-val">{{ invoices.filter(i=>i.status==='PENDING').length }}</span></div>
      </div>
    </div>

    <div class="inv-table-wrap">
      <table class="inv-table" v-if="invoices.length">
        <thead>
          <tr><th>No. Invoice</th><th>Tipe</th><th>Jumlah</th><th>Status</th><th>Tanggal</th><th>Aksi</th></tr>
        </thead>
        <tbody>
          <tr v-for="inv in invoices" :key="inv.id">
            <td><span class="inv-number">{{ inv.invoice_number }}</span></td>
            <td><span class="inv-type-badge" :class="'it-'+inv.type">{{ typeLabel(inv.type) }}</span></td>
            <td><span class="inv-amount">{{ fc(inv.amount) }}</span></td>
            <td>
              <span class="inv-badge" :class="'is-'+inv.status">
                <span class="inv-badge-dot"></span>
                {{ statusLabel(inv.status) }}
              </span>
            </td>
            <td><span class="inv-date">{{ formatDate(inv.created_at) }}</span></td>
            <td>
              <a v-if="inv.midtrans_redirect_url && inv.status==='PENDING'" :href="inv.midtrans_redirect_url" target="_blank" class="inv-pay-btn">
                <i class="ri-bank-card-line" style="font-size:14px"></i>
                Bayar
              </a>
            </td>
          </tr>
        </tbody>
      </table>
      <div v-else class="inv-empty">
        <div class="inv-empty-circle"><i class="ri-file-text-line"></i></div>
        <p>Belum ada invoice</p>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { invoiceApi } from '../services/api'
import { useToastStore } from '../stores/toast'

const toast = useToastStore()
const invoices = ref([])

onMounted(async () => {
  try { const { data } = await invoiceApi.list(); invoices.value = data.data || [] }
  catch { toast.error('Gagal memuat invoice') }
})

function fc(v) { return new Intl.NumberFormat('id-ID',{style:'currency',currency:'IDR',minimumFractionDigits:0}).format(v) }
function formatDate(d) { return d ? new Date(d).toLocaleDateString('id-ID',{day:'numeric',month:'short',year:'numeric'}) : '-' }
function typeLabel(t) { return {DP:'DP',CICILAN:'Cicilan',PELUNASAN:'Pelunasan',INVOICE:'Invoice'}[t]||t }
function statusLabel(s) { return {PAID:'Lunas',PENDING:'Menunggu',EXPIRED:'Kadaluarsa',FAILED:'Gagal',CANCELED:'Dibatalkan'}[s]||s }
</script>

<style scoped>
.inv-hero { background: linear-gradient(135deg, #0f0c29 0%, #302b63 50%, #24243e 100%); border-radius: 16px; padding: 32px 36px 24px; margin-bottom: 24px; box-shadow: 0 4px 24px rgba(15,12,41,0.2); }
.inv-hero-top { margin-bottom: 20px; }
.inv-hero-title { font-size: 1.6rem; font-weight: 800; color: #fff; margin: 0 0 4px; }
.inv-hero-sub { font-size: .85rem; color: rgba(255,255,255,.5); margin: 0; }
.inv-stats { display: flex; gap: 28px; flex-wrap: wrap; padding-top: 16px; border-top: 1px solid rgba(255,255,255,.08); }
.inv-stat { display: flex; align-items: center; gap: 8px; }
.inv-stat-dot { width: 8px; height: 8px; border-radius: 50%; }
.inv-stat-label { font-size: .72rem; color: rgba(255,255,255,.4); text-transform: uppercase; letter-spacing: .05em; }
.inv-stat-val { font-size: .9rem; font-weight: 800; color: #fff; }

.inv-table-wrap { background: #fff; border-radius: 16px; border: 1px solid #e8ecf1; overflow: hidden; }
.inv-table { width: 100%; border-collapse: collapse; }
.inv-table thead { background: #f8fafc; }
.inv-table th { padding: 14px 20px; font-size: .72rem; font-weight: 700; color: #64748b; text-transform: uppercase; letter-spacing: .05em; text-align: left; border-bottom: 1px solid #e8ecf1; }
.inv-table td { padding: 16px 20px; font-size: .85rem; color: #1e293b; border-bottom: 1px solid #f1f5f9; vertical-align: middle; }
.inv-table tr:last-child td { border-bottom: none; }
.inv-table tbody tr { transition: background .15s; }
.inv-table tbody tr:hover { background: #fafbfc; }

.inv-number { font-weight: 700; color: #0f172a; font-size: .82rem; }
.inv-amount { font-weight: 700; font-size: .88rem; color: #0f172a; }
.inv-date { font-size: .82rem; color: #64748b; }

.inv-type-badge { font-size: .68rem; font-weight: 700; padding: 4px 10px; border-radius: 6px; text-transform: uppercase; letter-spacing: .03em; }
.it-DP { background: #e0f2fe; color: #0284c7; }
.it-CICILAN { background: #ede9fe; color: #7c3aed; }
.it-PELUNASAN { background: #dcfce7; color: #16a34a; }
.it-INVOICE { background: #fef3c7; color: #d97706; }

.inv-badge { font-size: .68rem; font-weight: 700; padding: 4px 10px; border-radius: 6px; text-transform: uppercase; letter-spacing: .03em; display: inline-flex; align-items: center; gap: 5px; }
.inv-badge-dot { width: 6px; height: 6px; border-radius: 50%; }
.is-PAID { background: #dcfce7; color: #16a34a; }
.is-PAID .inv-badge-dot { background: #22c55e; }
.is-PENDING { background: #fef3c7; color: #d97706; }
.is-PENDING .inv-badge-dot { background: #f59e0b; }
.is-EXPIRED,.is-FAILED,.is-CANCELED { background: #f1f5f9; color: #64748b; }
.is-EXPIRED .inv-badge-dot,.is-FAILED .inv-badge-dot,.is-CANCELED .inv-badge-dot { background: #94a3b8; }

.inv-pay-btn { display: inline-flex; align-items: center; gap: 5px; font-size: .78rem; font-weight: 600; color: #fff; text-decoration: none; padding: 6px 14px; border-radius: 8px; background: linear-gradient(135deg, #6366f1, #8b5cf6); box-shadow: 0 2px 8px rgba(99,102,241,.25); transition: all .15s; }
.inv-pay-btn:hover { box-shadow: 0 4px 14px rgba(99,102,241,.35); transform: translateY(-1px); }

.inv-empty { text-align: center; padding: 56px 20px; }
.inv-empty-circle { width: 72px; height: 72px; border-radius: 50%; background: linear-gradient(135deg, #f1f5f9, #e2e8f0); display: flex; align-items: center; justify-content: center; margin: 0 auto 16px; }
.inv-empty p { color: #94a3b8; font-size: .85rem; margin: 0; }

@media (max-width: 768px) { .inv-hero { padding: 24px 20px 18px; } .inv-table-wrap { overflow-x: auto; } }
</style>
