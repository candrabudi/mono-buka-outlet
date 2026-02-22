<template>
  <div class="db-page">
    <!-- Hero -->
    <div class="db-hero">
      <div class="db-hero-inner">
        <div>
          <h1 class="db-hero-title">Selamat Datang Kembali 👋</h1>
          <p class="db-hero-sub">Ringkasan performa bisnis outlet franchise Anda</p>
        </div>
        <div class="db-hero-date">
          <i class="ri-calendar-line"></i> {{ todayDate }}
        </div>
      </div>
    </div>

    <!-- Loading -->
    <div v-if="loading" class="db-loading">
      <div class="db-spinner"></div>
      <span>Memuat dashboard...</span>
    </div>

    <template v-else>
      <!-- KPI Cards Row 1 -->
      <div class="db-kpi-row">
        <div class="db-kpi blue">
          <div class="db-kpi-top">
            <div class="db-kpi-icon"><i class="ri-store-2-fill"></i></div>
            <span class="db-kpi-label">Total Outlet</span>
          </div>
          <div class="db-kpi-value">{{ s.total_outlets ?? 0 }}</div>
          <div class="db-kpi-sub">Outlet terdaftar</div>
        </div>
        <div class="db-kpi green">
          <div class="db-kpi-top">
            <div class="db-kpi-icon"><i class="ri-team-fill"></i></div>
            <span class="db-kpi-label">Mitra Aktif</span>
          </div>
          <div class="db-kpi-value">{{ s.total_mitra ?? 0 }}</div>
          <div class="db-kpi-sub">dari {{ s.total_users ?? 0 }} total user</div>
        </div>
        <div class="db-kpi purple">
          <div class="db-kpi-top">
            <div class="db-kpi-icon"><i class="ri-handshake-fill"></i></div>
            <span class="db-kpi-label">Partnership</span>
          </div>
          <div class="db-kpi-value">{{ s.total_partnerships ?? 0 }}</div>
          <div class="db-kpi-sub">Kemitraan aktif</div>
        </div>
        <div class="db-kpi orange">
          <div class="db-kpi-top">
            <div class="db-kpi-icon"><i class="ri-money-dollar-circle-fill"></i></div>
            <span class="db-kpi-label">Total Income</span>
          </div>
          <div class="db-kpi-value">{{ fc(s.total_income_amount) }}</div>
          <div class="db-kpi-sub">Invoice terbayar</div>
        </div>
      </div>

      <!-- KPI Cards Row 2 — Financial -->
      <div class="db-kpi-row small">
        <div class="db-kpi-mini">
          <i class="ri-file-list-3-line text-amber"></i>
          <div>
            <div class="db-kpi-mini-val">{{ s.pending_applications ?? 0 }}</div>
            <div class="db-kpi-mini-label">Pengajuan Pending</div>
          </div>
        </div>
        <div class="db-kpi-mini">
          <i class="ri-bill-line text-red"></i>
          <div>
            <div class="db-kpi-mini-val">{{ s.pending_invoices ?? 0 }}</div>
            <div class="db-kpi-mini-label">Invoice Belum Bayar</div>
          </div>
        </div>
        <div class="db-kpi-mini">
          <i class="ri-wallet-3-line text-cyan"></i>
          <div>
            <div class="db-kpi-mini-val">{{ fc(s.pending_invoice_amount) }}</div>
            <div class="db-kpi-mini-label">Nominal Tertunda</div>
          </div>
        </div>
        <div class="db-kpi-mini">
          <i class="ri-bar-chart-box-line text-green"></i>
          <div>
            <div class="db-kpi-mini-val">{{ fc(s.monthly_revenue) }}</div>
            <div class="db-kpi-mini-label">Revenue Bulan Ini</div>
          </div>
        </div>
      </div>

      <!-- Main Grid -->
      <div class="db-grid">
        <!-- Revenue Chart -->
        <div class="db-card wide">
          <div class="db-card-head">
            <h3><i class="ri-line-chart-fill"></i> Tren Revenue (6 Bulan)</h3>
          </div>
          <div class="db-card-body chart-body">
            <Bar v-if="revenueChartData" :data="revenueChartData" :options="revenueChartOptions" />
            <div v-else class="db-empty"><i class="ri-line-chart-line"></i><span>Belum ada data revenue</span></div>
          </div>
        </div>

        <!-- Partnership Status Donut -->
        <div class="db-card">
          <div class="db-card-head">
            <h3><i class="ri-pie-chart-2-fill"></i> Status Partnership</h3>
            <router-link to="/partnerships" class="db-card-link">Detail →</router-link>
          </div>
          <div class="db-card-body">
            <div v-if="partnershipChartData" class="db-donut-wrap">
              <Doughnut :data="partnershipChartData" :options="donutOptions" />
            </div>
            <div class="db-donut-legend" v-if="s.partnerships_by_status">
              <div v-for="(count, status) in s.partnerships_by_status" :key="status" class="db-legend-item">
                <span class="db-legend-dot" :style="{background: psColor(status)}"></span>
                <span class="db-legend-label">{{ psLabel(status) }}</span>
                <span class="db-legend-val">{{ count }}</span>
              </div>
            </div>
            <div v-if="!Object.keys(s.partnerships_by_status || {}).length" class="db-empty"><i class="ri-pie-chart-line"></i><span>Belum ada data</span></div>
          </div>
        </div>
      </div>

      <!-- Second Grid -->
      <div class="db-grid mt">
        <!-- Pengajuan -->
        <div class="db-card">
          <div class="db-card-head">
            <h3><i class="ri-file-list-3-fill"></i> Pengajuan Kemitraan</h3>
            <router-link to="/applications" class="db-card-link">Lihat Semua →</router-link>
          </div>
          <div class="db-card-body">
            <div class="db-app-grid">
              <div class="db-app-tile pending">
                <i class="ri-time-fill"></i>
                <span class="db-app-num">{{ appStat('PENDING') }}</span>
                <span class="db-app-label">Pending</span>
              </div>
              <div class="db-app-tile review">
                <i class="ri-search-eye-fill"></i>
                <span class="db-app-num">{{ appStat('REVIEWED') }}</span>
                <span class="db-app-label">Ditinjau</span>
              </div>
              <div class="db-app-tile approved">
                <i class="ri-checkbox-circle-fill"></i>
                <span class="db-app-num">{{ appStat('APPROVED') }}</span>
                <span class="db-app-label">Disetujui</span>
              </div>
              <div class="db-app-tile rejected">
                <i class="ri-close-circle-fill"></i>
                <span class="db-app-num">{{ appStat('REJECTED') }}</span>
                <span class="db-app-label">Ditolak</span>
              </div>
            </div>
            <div class="db-app-total">
              <span>Total Pengajuan</span>
              <span class="db-app-total-val">{{ s.total_applications ?? 0 }}</span>
            </div>
          </div>
        </div>

        <!-- Invoice & Pembayaran -->
        <div class="db-card">
          <div class="db-card-head">
            <h3><i class="ri-secure-payment-fill"></i> Invoice & Pembayaran</h3>
          </div>
          <div class="db-card-body">
            <div class="db-fin-grid">
              <div class="db-fin-item">
                <div class="db-fin-icon green"><i class="ri-money-dollar-circle-fill"></i></div>
                <div>
                  <div class="db-fin-amount">{{ fc(s.total_income_amount) }}</div>
                  <div class="db-fin-label">Total Diterima</div>
                </div>
              </div>
              <div class="db-fin-item">
                <div class="db-fin-icon amber"><i class="ri-time-fill"></i></div>
                <div>
                  <div class="db-fin-amount">{{ fc(s.pending_invoice_amount) }}</div>
                  <div class="db-fin-label">Belum Dibayar</div>
                </div>
              </div>
              <div class="db-fin-item">
                <div class="db-fin-icon blue"><i class="ri-bank-card-fill"></i></div>
                <div>
                  <div class="db-fin-amount">{{ fc(s.total_payments_verified) }}</div>
                  <div class="db-fin-label">Pembayaran Terverifikasi</div>
                </div>
              </div>
            </div>
            <div class="db-fin-bar">
              <div class="db-fin-bar-fill" :style="{ width: invoicePaidPct + '%' }"></div>
            </div>
            <div class="db-fin-bar-info">
              <span>{{ s.paid_invoices ?? 0 }} terbayar</span>
              <span>dari {{ (s.paid_invoices ?? 0) + (s.pending_invoices ?? 0) }} invoice</span>
            </div>
          </div>
        </div>

        <!-- Lead Pipeline -->
        <div class="db-card wide">
          <div class="db-card-head">
            <h3><i class="ri-user-follow-fill"></i> Lead Pipeline</h3>
          </div>
          <div class="db-card-body">
            <div v-if="s.leads_by_status && Object.keys(s.leads_by_status).length" class="db-pipeline">
              <div v-for="(count, status) in s.leads_by_status" :key="status" class="db-pipe-item">
                <div class="db-pipe-bar-track">
                  <div class="db-pipe-bar-fill" :style="{ height: pipePct(count) + '%', background: leadColor(status) }"></div>
                </div>
                <div class="db-pipe-count">{{ count }}</div>
                <div class="db-pipe-label">{{ formatStatus(status) }}</div>
              </div>
            </div>
            <div v-else class="db-empty"><i class="ri-user-add-line"></i><span>Belum ada leads</span></div>
          </div>
        </div>
      </div>

      <!-- Quick Actions -->
      <div class="db-quick-row">
        <router-link to="/outlets/create" class="db-quick-btn"><i class="ri-add-circle-fill"></i> Tambah Outlet</router-link>
        <router-link to="/applications" class="db-quick-btn"><i class="ri-file-search-fill"></i> Review Pengajuan</router-link>
        <router-link to="/partnerships" class="db-quick-btn"><i class="ri-handshake-fill"></i> Partnership</router-link>
        <router-link to="/mitra" class="db-quick-btn"><i class="ri-group-fill"></i> Kelola Mitra</router-link>
        <router-link to="/meetings" class="db-quick-btn"><i class="ri-calendar-event-fill"></i> Meetings</router-link>
        <router-link to="/locations" class="db-quick-btn"><i class="ri-map-pin-fill"></i> Lokasi</router-link>
      </div>
    </template>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { dashboardApi } from '../services/api'
import { Bar, Doughnut } from 'vue-chartjs'
import { Chart as ChartJS, ArcElement, BarElement, CategoryScale, LinearScale, Tooltip, Legend } from 'chart.js'

ChartJS.register(ArcElement, BarElement, CategoryScale, LinearScale, Tooltip, Legend)

const s = ref({})
const loading = ref(true)

onMounted(async () => {
  try {
    const { data } = await dashboardApi.stats()
    s.value = data.data || {}
  } catch (e) { console.error(e) }
  finally { loading.value = false }
})

const todayDate = new Date().toLocaleDateString('id-ID', { weekday: 'long', day: 'numeric', month: 'long', year: 'numeric' })

function fc(v) { return 'Rp ' + (v || 0).toLocaleString('id-ID') }
function formatStatus(st) { return st.replace(/_/g, ' ') }
function appStat(status) { return (s.value.applications_by_status || {})[status] || 0 }
function psLabel(st) {
  const map = { PENDING: 'Pending', ACTIVE: 'Active', RUNNING: 'Running', COMPLETED: 'Selesai', CANCELLED: 'Batal' }
  return map[st] || st
}
function psColor(st) {
  const map = { PENDING: '#f59e0b', ACTIVE: '#3b82f6', RUNNING: '#6366f1', COMPLETED: '#22c55e', CANCELLED: '#ef4444' }
  return map[st] || '#94a3b8'
}
function leadColor(st) {
  const map = { NEW: '#3b82f6', RUNNING: '#22c55e', COMPLETED: '#10b981', DP_PAID: '#f59e0b', FULLY_PAID: '#8b5cf6' }
  return map[st] || '#94a3b8'
}

const invoicePaidPct = computed(() => {
  const total = (s.value.paid_invoices || 0) + (s.value.pending_invoices || 0)
  return total ? Math.round((s.value.paid_invoices / total) * 100) : 0
})

function pipePct(count) {
  const max = Math.max(...Object.values(s.value.leads_by_status || { x: 1 }))
  return max ? (count / max) * 100 : 0
}

// Revenue Bar Chart
const revenueChartData = computed(() => {
  const chart = s.value.revenue_chart
  if (!chart || !chart.length) return null
  return {
    labels: chart.map(c => c.month),
    datasets: [{
      label: 'Revenue',
      data: chart.map(c => c.revenue || 0),
      backgroundColor: 'rgba(99, 102, 241, 0.8)',
      borderRadius: 8,
      borderSkipped: false,
      barThickness: 32,
    }]
  }
})

const revenueChartOptions = {
  responsive: true,
  maintainAspectRatio: false,
  plugins: {
    legend: { display: false },
    tooltip: {
      callbacks: { label: (ctx) => 'Rp ' + (ctx.raw || 0).toLocaleString('id-ID') }
    }
  },
  scales: {
    x: { grid: { display: false }, ticks: { font: { size: 12, weight: 600 } } },
    y: {
      grid: { color: '#f1f5f9' },
      ticks: {
        font: { size: 11 },
        callback: (v) => v >= 1e6 ? (v / 1e6).toFixed(0) + 'jt' : v >= 1e3 ? (v / 1e3).toFixed(0) + 'rb' : v
      }
    }
  }
}

// Partnership Donut
const partnershipChartData = computed(() => {
  const data = s.value.partnerships_by_status
  if (!data || !Object.keys(data).length) return null
  const labels = Object.keys(data).map(psLabel)
  const values = Object.values(data)
  const colors = Object.keys(data).map(psColor)
  return {
    labels,
    datasets: [{ data: values, backgroundColor: colors, borderWidth: 0, hoverOffset: 6 }]
  }
})

const donutOptions = {
  responsive: true,
  maintainAspectRatio: false,
  cutout: '65%',
  plugins: {
    legend: { display: false },
    tooltip: { callbacks: { label: (ctx) => ctx.label + ': ' + ctx.raw } }
  }
}
</script>

<style scoped>
.db-page { padding: 0; }

/* Hero */
.db-hero { background: linear-gradient(135deg, #0f172a 0%, #1e293b 100%); border-radius: 16px; padding: 28px 32px; margin-bottom: 24px; }
.db-hero-inner { display: flex; justify-content: space-between; align-items: center; }
.db-hero-title { font-size: 1.5rem; font-weight: 700; color: #fff; margin: 0 0 4px; }
.db-hero-sub { font-size: .85rem; color: #94a3b8; margin: 0; }
.db-hero-date { font-size: .82rem; color: #94a3b8; display: flex; align-items: center; gap: 6px; }

/* Loading */
.db-loading { display: flex; align-items: center; justify-content: center; gap: 12px; padding: 60px; color: #94a3b8; font-size: .9rem; }
.db-spinner { width: 24px; height: 24px; border: 3px solid #e2e8f0; border-top-color: #6366f1; border-radius: 50%; animation: spin .6s linear infinite; }
@keyframes spin { to { transform: rotate(360deg) } }

/* KPI Cards */
.db-kpi-row { display: grid; grid-template-columns: repeat(4, 1fr); gap: 16px; margin-bottom: 16px; }
.db-kpi-row.small { margin-bottom: 24px; }
.db-kpi { background: #fff; border: 1px solid #e2e8f0; border-radius: 14px; padding: 20px; position: relative; overflow: hidden; transition: all .2s; }
.db-kpi:hover { box-shadow: 0 8px 24px rgba(0,0,0,.06); transform: translateY(-2px); }
.db-kpi::after { content: ''; position: absolute; top: 0; left: 0; width: 4px; height: 100%; border-radius: 4px 0 0 4px; }
.db-kpi.blue::after { background: #3b82f6; }
.db-kpi.green::after { background: #22c55e; }
.db-kpi.purple::after { background: #8b5cf6; }
.db-kpi.orange::after { background: #f97316; }
.db-kpi-top { display: flex; align-items: center; gap: 10px; margin-bottom: 12px; }
.db-kpi-icon { width: 40px; height: 40px; border-radius: 10px; display: flex; align-items: center; justify-content: center; font-size: 20px; }
.db-kpi.blue .db-kpi-icon { background: rgba(59,130,246,.1); color: #3b82f6; }
.db-kpi.green .db-kpi-icon { background: rgba(34,197,94,.1); color: #22c55e; }
.db-kpi.purple .db-kpi-icon { background: rgba(139,92,246,.1); color: #8b5cf6; }
.db-kpi.orange .db-kpi-icon { background: rgba(249,115,22,.1); color: #f97316; }
.db-kpi-label { font-size: .75rem; font-weight: 600; color: #94a3b8; text-transform: uppercase; letter-spacing: .5px; }
.db-kpi-value { font-size: 1.6rem; font-weight: 800; color: #0f172a; line-height: 1; }
.db-kpi-sub { font-size: .72rem; color: #94a3b8; margin-top: 6px; }

/* Mini KPI */
.db-kpi-mini { background: #fff; border: 1px solid #e2e8f0; border-radius: 12px; padding: 16px 18px; display: flex; align-items: center; gap: 14px; transition: all .15s; }
.db-kpi-mini:hover { box-shadow: 0 4px 12px rgba(0,0,0,.04); }
.db-kpi-mini > i { font-size: 24px; }
.text-amber { color: #f59e0b; }
.text-red { color: #ef4444; }
.text-cyan { color: #06b6d4; }
.text-green { color: #22c55e; }
.db-kpi-mini-val { font-size: 1.1rem; font-weight: 800; color: #0f172a; }
.db-kpi-mini-label { font-size: .72rem; color: #94a3b8; font-weight: 500; }

/* Grid */
.db-grid { display: grid; grid-template-columns: 1fr 1fr; gap: 20px; }
.db-grid.mt { margin-top: 20px; }

/* Card */
.db-card { background: #fff; border: 1px solid #e2e8f0; border-radius: 14px; overflow: hidden; }
.db-card.wide { grid-column: 1 / -1; }
.db-card-head { display: flex; justify-content: space-between; align-items: center; padding: 16px 20px; border-bottom: 1px solid #f1f5f9; }
.db-card-head h3 { font-size: .88rem; font-weight: 700; color: #0f172a; margin: 0; display: flex; align-items: center; gap: 8px; }
.db-card-head h3 i { font-size: 1.1rem; color: #6366f1; }
.db-card-link { font-size: .78rem; font-weight: 600; color: #6366f1; text-decoration: none; }
.db-card-link:hover { text-decoration: underline; }
.db-card-body { padding: 20px; }
.chart-body { height: 260px; }

/* Donut */
.db-donut-wrap { height: 180px; margin-bottom: 16px; }
.db-donut-legend { display: flex; flex-direction: column; gap: 6px; }
.db-legend-item { display: flex; align-items: center; gap: 8px; font-size: .82rem; }
.db-legend-dot { width: 10px; height: 10px; border-radius: 50%; flex-shrink: 0; }
.db-legend-label { flex: 1; color: #475569; font-weight: 500; }
.db-legend-val { font-weight: 800; color: #0f172a; }

/* App Tiles */
.db-app-grid { display: grid; grid-template-columns: repeat(4, 1fr); gap: 10px; margin-bottom: 14px; }
.db-app-tile { display: flex; flex-direction: column; align-items: center; padding: 14px 8px; border-radius: 12px; text-align: center; }
.db-app-tile i { font-size: 1.3rem; margin-bottom: 6px; }
.db-app-tile.pending { background: #fef3c7; color: #92400e; }
.db-app-tile.review { background: #dbeafe; color: #1e40af; }
.db-app-tile.approved { background: #dcfce7; color: #166534; }
.db-app-tile.rejected { background: #fee2e2; color: #991b1b; }
.db-app-num { font-size: 1.3rem; font-weight: 800; display: block; }
.db-app-label { font-size: .68rem; font-weight: 600; margin-top: 2px; }
.db-app-total { display: flex; justify-content: space-between; padding: 10px 14px; background: #f8fafc; border-radius: 8px; font-size: .82rem; font-weight: 600; color: #475569; }
.db-app-total-val { font-weight: 800; color: #0f172a; }

/* Financial */
.db-fin-grid { display: flex; flex-direction: column; gap: 14px; margin-bottom: 16px; }
.db-fin-item { display: flex; align-items: center; gap: 14px; }
.db-fin-icon { width: 42px; height: 42px; border-radius: 10px; display: flex; align-items: center; justify-content: center; font-size: 20px; }
.db-fin-icon.green { background: rgba(34,197,94,.1); color: #22c55e; }
.db-fin-icon.amber { background: rgba(245,158,11,.1); color: #f59e0b; }
.db-fin-icon.blue { background: rgba(59,130,246,.1); color: #3b82f6; }
.db-fin-amount { font-size: 1rem; font-weight: 800; color: #0f172a; }
.db-fin-label { font-size: .72rem; color: #94a3b8; }
.db-fin-bar { height: 8px; background: #f1f5f9; border-radius: 10px; overflow: hidden; margin-bottom: 6px; }
.db-fin-bar-fill { height: 100%; background: linear-gradient(90deg, #22c55e, #10b981); border-radius: 10px; transition: width .6s; }
.db-fin-bar-info { display: flex; justify-content: space-between; font-size: .72rem; color: #94a3b8; font-weight: 500; }

/* Pipeline Vertical */
.db-pipeline { display: flex; align-items: flex-end; gap: 0; justify-content: center; height: 160px; padding: 0 20px; }
.db-pipe-item { flex: 1; display: flex; flex-direction: column; align-items: center; gap: 6px; max-width: 80px; }
.db-pipe-bar-track { width: 32px; height: 120px; background: #f1f5f9; border-radius: 8px; overflow: hidden; display: flex; align-items: flex-end; }
.db-pipe-bar-fill { width: 100%; border-radius: 8px; transition: height .6s ease; min-height: 6px; }
.db-pipe-count { font-size: .88rem; font-weight: 800; color: #0f172a; }
.db-pipe-label { font-size: .65rem; font-weight: 600; color: #94a3b8; text-transform: uppercase; text-align: center; }

/* Quick Actions */
.db-quick-row { display: flex; flex-wrap: wrap; gap: 10px; margin-top: 24px; }
.db-quick-btn { display: inline-flex; align-items: center; gap: 8px; padding: 10px 18px; border-radius: 10px; font-size: .82rem; font-weight: 600; color: #334155; text-decoration: none; background: #fff; border: 1px solid #e2e8f0; transition: all .15s; }
.db-quick-btn:hover { background: #6366f1; color: #fff; border-color: #6366f1; transform: translateY(-1px); box-shadow: 0 4px 12px rgba(99,102,241,.2); }
.db-quick-btn i { font-size: 1rem; }

/* Empty */
.db-empty { text-align: center; padding: 28px; color: #94a3b8; display: flex; flex-direction: column; align-items: center; gap: 6px; }
.db-empty i { font-size: 2rem; opacity: .3; }
.db-empty span { font-size: .82rem; }

/* Responsive */
@media (max-width: 1280px) { .db-grid { grid-template-columns: 1fr; } .db-card.wide { grid-column: auto; } }
@media (max-width: 768px) {
  .db-kpi-row { grid-template-columns: repeat(2, 1fr); }
  .db-hero { padding: 20px; }
  .db-hero-inner { flex-direction: column; gap: 8px; align-items: flex-start; }
  .db-app-grid { grid-template-columns: repeat(2, 1fr); }
}
@media (max-width: 480px) { .db-kpi-row { grid-template-columns: 1fr; } }
</style>
