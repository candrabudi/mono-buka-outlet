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
      <!-- KPI Cards -->
      <div class="db-kpi-row">
        <div class="db-kpi blue">
          <div class="db-kpi-icon"><i class="ri-store-2-fill"></i></div>
          <div class="db-kpi-info">
            <span class="db-kpi-label">Total Outlet</span>
            <span class="db-kpi-value">{{ s.total_outlets ?? 0 }}</span>
          </div>
        </div>
        <div class="db-kpi green">
          <div class="db-kpi-icon"><i class="ri-team-fill"></i></div>
          <div class="db-kpi-info">
            <span class="db-kpi-label">Mitra Aktif</span>
            <span class="db-kpi-value">{{ s.total_mitra ?? 0 }} <small>/ {{ s.total_users ?? 0 }}</small></span>
          </div>
        </div>
        <div class="db-kpi purple">
          <div class="db-kpi-icon"><i class="ri-handshake-fill"></i></div>
          <div class="db-kpi-info">
            <span class="db-kpi-label">Partnership</span>
            <span class="db-kpi-value">{{ s.total_partnerships ?? 0 }}</span>
          </div>
        </div>
        <div class="db-kpi orange">
          <div class="db-kpi-icon"><i class="ri-money-dollar-circle-fill"></i></div>
          <div class="db-kpi-info">
            <span class="db-kpi-label">Total Income</span>
            <span class="db-kpi-value">{{ fcShort(s.total_income_amount) }}</span>
          </div>
        </div>
        <div class="db-kpi cyan">
          <div class="db-kpi-icon"><i class="ri-file-list-3-fill"></i></div>
          <div class="db-kpi-info">
            <span class="db-kpi-label">Pengajuan Pending</span>
            <span class="db-kpi-value">{{ s.pending_applications ?? 0 }}</span>
          </div>
        </div>
        <div class="db-kpi red">
          <div class="db-kpi-icon"><i class="ri-bill-fill"></i></div>
          <div class="db-kpi-info">
            <span class="db-kpi-label">Invoice Tertunda</span>
            <span class="db-kpi-value">{{ s.pending_invoices ?? 0 }}</span>
          </div>
        </div>
      </div>

      <!-- Row: Revenue + Partnership Donut -->
      <div class="db-row-2-1">
        <!-- Revenue Chart -->
        <div class="db-card">
          <div class="db-card-head">
            <h3><i class="ri-bar-chart-grouped-fill"></i> Tren Revenue (6 Bulan)</h3>
          </div>
          <div class="db-card-body" style="height:280px">
            <Bar v-if="revenueChartData" :data="revenueChartData" :options="revenueChartOptions" />
            <div v-else class="db-empty"><i class="ri-line-chart-line"></i><span>Belum ada data revenue</span></div>
          </div>
        </div>

        <!-- Partnership Donut -->
        <div class="db-card">
          <div class="db-card-head">
            <h3><i class="ri-pie-chart-2-fill"></i> Status Partnership</h3>
            <router-link to="/partnerships" class="db-card-link">Detail →</router-link>
          </div>
          <div class="db-card-body">
            <div v-if="partnershipChartData" class="db-donut-wrap">
              <Doughnut :data="partnershipChartData" :options="donutOptions" />
            </div>
            <div class="db-donut-legend" v-if="s.partnerships_by_status && Object.keys(s.partnerships_by_status).length">
              <div v-for="(count, status) in s.partnerships_by_status" :key="status" class="db-legend-item">
                <span class="db-legend-dot" :style="{background: psColor(status)}"></span>
                <span class="db-legend-label">{{ psLabel(status) }}</span>
                <span class="db-legend-val">{{ count }}</span>
              </div>
            </div>
            <div v-if="!Object.keys(s.partnerships_by_status || {}).length" class="db-empty sm"><i class="ri-pie-chart-line"></i><span>Belum ada data</span></div>
          </div>
        </div>
      </div>

      <!-- Row: Pengajuan + Invoice + Lead -->
      <div class="db-row-3">
        <!-- Pengajuan -->
        <div class="db-card">
          <div class="db-card-head">
            <h3><i class="ri-file-search-fill"></i> Pengajuan Kemitraan</h3>
            <router-link to="/applications" class="db-card-link">Semua →</router-link>
          </div>
          <div class="db-card-body">
            <div class="db-app-grid">
              <div class="db-app-tile pending">
                <i class="ri-time-fill"></i>
                <span class="db-app-num">{{ appStat('PENDING') }}</span>
                <span class="db-app-lbl">Pending</span>
              </div>
              <div class="db-app-tile review">
                <i class="ri-search-eye-fill"></i>
                <span class="db-app-num">{{ appStat('REVIEWED') }}</span>
                <span class="db-app-lbl">Ditinjau</span>
              </div>
              <div class="db-app-tile approved">
                <i class="ri-checkbox-circle-fill"></i>
                <span class="db-app-num">{{ appStat('APPROVED') }}</span>
                <span class="db-app-lbl">Disetujui</span>
              </div>
              <div class="db-app-tile rejected">
                <i class="ri-close-circle-fill"></i>
                <span class="db-app-num">{{ appStat('REJECTED') }}</span>
                <span class="db-app-lbl">Ditolak</span>
              </div>
            </div>
            <div class="db-app-footer">
              <span>Total</span>
              <span class="fw800">{{ s.total_applications ?? 0 }}</span>
            </div>
          </div>
        </div>

        <!-- Invoice & Pembayaran -->
        <div class="db-card">
          <div class="db-card-head">
            <h3><i class="ri-secure-payment-fill"></i> Invoice & Pembayaran</h3>
          </div>
          <div class="db-card-body">
            <div class="db-fin-rows">
              <div class="db-fin-row">
                <div class="db-fin-ic green"><i class="ri-money-dollar-circle-fill"></i></div>
                <div class="db-fin-txt">
                  <div class="db-fin-amt">{{ fc(s.total_income_amount) }}</div>
                  <div class="db-fin-lbl">Total Diterima</div>
                </div>
              </div>
              <div class="db-fin-row">
                <div class="db-fin-ic amber"><i class="ri-time-fill"></i></div>
                <div class="db-fin-txt">
                  <div class="db-fin-amt">{{ fc(s.pending_invoice_amount) }}</div>
                  <div class="db-fin-lbl">Belum Dibayar</div>
                </div>
              </div>
              <div class="db-fin-row">
                <div class="db-fin-ic blue"><i class="ri-bank-card-fill"></i></div>
                <div class="db-fin-txt">
                  <div class="db-fin-amt">{{ fc(s.total_payments_verified) }}</div>
                  <div class="db-fin-lbl">Terverifikasi</div>
                </div>
              </div>
            </div>
            <div class="db-fin-progress">
              <div class="db-fin-bar"><div class="db-fin-fill" :style="{ width: invoicePaidPct + '%' }"></div></div>
              <div class="db-fin-info"><span>{{ s.paid_invoices ?? 0 }} terbayar</span><span>{{ (s.paid_invoices??0) + (s.pending_invoices??0) }} total</span></div>
            </div>
          </div>
        </div>

        <!-- Lead Pipeline -->
        <div class="db-card">
          <div class="db-card-head">
            <h3><i class="ri-user-follow-fill"></i> Lead Pipeline</h3>
          </div>
          <div class="db-card-body">
            <div v-if="s.leads_by_status && Object.keys(s.leads_by_status).length" class="db-lead-list">
              <div v-for="(count, status) in s.leads_by_status" :key="status" class="db-lead-row">
                <div class="db-lead-dot" :style="{background: leadColor(status)}"></div>
                <span class="db-lead-name">{{ formatStatus(status) }}</span>
                <div class="db-lead-bar-track">
                  <div class="db-lead-bar-fill" :style="{ width: pipePct(count) + '%', background: leadColor(status) }"></div>
                </div>
                <span class="db-lead-count">{{ count }}</span>
              </div>
            </div>
            <div v-else class="db-empty sm"><i class="ri-user-add-line"></i><span>Belum ada leads</span></div>
          </div>
        </div>
      </div>

      <!-- Quick Actions -->
      <div class="db-actions">
        <router-link to="/outlets/create" class="db-act"><i class="ri-add-circle-fill"></i> Tambah Outlet</router-link>
        <router-link to="/applications" class="db-act"><i class="ri-file-search-fill"></i> Review Pengajuan</router-link>
        <router-link to="/partnerships" class="db-act"><i class="ri-handshake-fill"></i> Partnership</router-link>
        <router-link to="/mitra" class="db-act"><i class="ri-group-fill"></i> Kelola Mitra</router-link>
        <router-link to="/meetings" class="db-act"><i class="ri-calendar-event-fill"></i> Meetings</router-link>
        <router-link to="/locations" class="db-act"><i class="ri-map-pin-fill"></i> Lokasi</router-link>
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
function fcShort(v) {
  if (!v) return 'Rp 0'
  if (v >= 1e9) return 'Rp ' + (v/1e9).toFixed(1) + ' M'
  if (v >= 1e6) return 'Rp ' + (v/1e6).toFixed(1) + ' jt'
  if (v >= 1e3) return 'Rp ' + (v/1e3).toFixed(0) + ' rb'
  return 'Rp ' + v
}
function formatStatus(st) { return st.replace(/_/g, ' ') }
function appStat(status) { return (s.value.applications_by_status || {})[status] || 0 }
function psLabel(st) {
  const map = { PENDING: 'Pending', ACTIVE: 'Active', RUNNING: 'Running', COMPLETED: 'Selesai', CANCELLED: 'Batal', AGREEMENT_SIGNED: 'Agreement' }
  return map[st] || st
}
function psColor(st) {
  const map = { PENDING: '#f59e0b', ACTIVE: '#3b82f6', RUNNING: '#6366f1', COMPLETED: '#22c55e', CANCELLED: '#ef4444', AGREEMENT_SIGNED: '#8b5cf6' }
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
      backgroundColor: 'rgba(99, 102, 241, 0.85)',
      borderRadius: 8,
      borderSkipped: false,
      barThickness: 28,
    }]
  }
})

const revenueChartOptions = {
  responsive: true,
  maintainAspectRatio: false,
  plugins: {
    legend: { display: false },
    tooltip: { callbacks: { label: (ctx) => 'Rp ' + (ctx.raw || 0).toLocaleString('id-ID') } }
  },
  scales: {
    x: { grid: { display: false }, ticks: { font: { size: 11, weight: '600' } } },
    y: {
      grid: { color: '#f1f5f9' },
      ticks: { font: { size: 10 }, callback: (v) => v >= 1e6 ? (v/1e6).toFixed(0) + 'jt' : v >= 1e3 ? (v/1e3).toFixed(0) + 'rb' : v }
    }
  }
}

// Partnership Donut
const partnershipChartData = computed(() => {
  const data = s.value.partnerships_by_status
  if (!data || !Object.keys(data).length) return null
  return {
    labels: Object.keys(data).map(psLabel),
    datasets: [{ data: Object.values(data), backgroundColor: Object.keys(data).map(psColor), borderWidth: 0, hoverOffset: 6 }]
  }
})

const donutOptions = {
  responsive: true,
  maintainAspectRatio: false,
  cutout: '65%',
  plugins: { legend: { display: false }, tooltip: { callbacks: { label: (ctx) => ctx.label + ': ' + ctx.raw } } }
}
</script>

<style scoped>
.db-page { padding: 0; }

/* Hero */
.db-hero { background: linear-gradient(135deg, #0f172a 0%, #1e293b 100%); border-radius: 16px; padding: 24px 28px; margin-bottom: 20px; }
.db-hero-inner { display: flex; justify-content: space-between; align-items: center; }
.db-hero-title { font-size: 1.4rem; font-weight: 700; color: #fff; margin: 0 0 4px; }
.db-hero-sub { font-size: .82rem; color: #94a3b8; margin: 0; }
.db-hero-date { font-size: .8rem; color: #94a3b8; display: flex; align-items: center; gap: 6px; white-space: nowrap; }

/* Loading */
.db-loading { display: flex; align-items: center; justify-content: center; gap: 12px; padding: 60px; color: #94a3b8; }
.db-spinner { width: 24px; height: 24px; border: 3px solid #e2e8f0; border-top-color: #6366f1; border-radius: 50%; animation: spin .6s linear infinite; }
@keyframes spin { to { transform: rotate(360deg) } }

/* KPI Row — 6 cards */
.db-kpi-row { display: grid; grid-template-columns: repeat(6, 1fr); gap: 12px; margin-bottom: 20px; }
.db-kpi { background: #fff; border: 1px solid #e2e8f0; border-radius: 12px; padding: 16px; display: flex; align-items: center; gap: 12px; transition: all .2s; position: relative; overflow: hidden; }
.db-kpi:hover { box-shadow: 0 6px 20px rgba(0,0,0,.05); transform: translateY(-1px); }
.db-kpi::before { content:''; position: absolute; top: 0; left: 0; width: 3px; height: 100%; }
.db-kpi.blue::before { background: #3b82f6; } .db-kpi.green::before { background: #22c55e; }
.db-kpi.purple::before { background: #8b5cf6; } .db-kpi.orange::before { background: #f97316; }
.db-kpi.cyan::before { background: #06b6d4; } .db-kpi.red::before { background: #ef4444; }
.db-kpi-icon { width: 38px; height: 38px; border-radius: 10px; display: flex; align-items: center; justify-content: center; font-size: 18px; flex-shrink: 0; }
.db-kpi.blue .db-kpi-icon { background: rgba(59,130,246,.1); color: #3b82f6; }
.db-kpi.green .db-kpi-icon { background: rgba(34,197,94,.1); color: #22c55e; }
.db-kpi.purple .db-kpi-icon { background: rgba(139,92,246,.1); color: #8b5cf6; }
.db-kpi.orange .db-kpi-icon { background: rgba(249,115,22,.1); color: #f97316; }
.db-kpi.cyan .db-kpi-icon { background: rgba(6,182,212,.1); color: #06b6d4; }
.db-kpi.red .db-kpi-icon { background: rgba(239,68,68,.1); color: #ef4444; }
.db-kpi-label { font-size: .68rem; color: #94a3b8; font-weight: 600; text-transform: uppercase; letter-spacing: .3px; }
.db-kpi-value { font-size: 1.15rem; font-weight: 800; color: #0f172a; display: block; margin-top: 1px; }
.db-kpi-value small { font-size: .7rem; color: #94a3b8; font-weight: 500; }

/* 2:1 Row */
.db-row-2-1 { display: grid; grid-template-columns: 2fr 1fr; gap: 16px; margin-bottom: 16px; }

/* 3 Column Row */
.db-row-3 { display: grid; grid-template-columns: 1fr 1fr 1fr; gap: 16px; margin-bottom: 16px; }

/* Card */
.db-card { background: #fff; border: 1px solid #e2e8f0; border-radius: 12px; overflow: hidden; }
.db-card-head { display: flex; justify-content: space-between; align-items: center; padding: 14px 18px; border-bottom: 1px solid #f1f5f9; }
.db-card-head h3 { font-size: .82rem; font-weight: 700; color: #0f172a; margin: 0; display: flex; align-items: center; gap: 7px; }
.db-card-head h3 i { font-size: 1rem; color: #6366f1; }
.db-card-link { font-size: .75rem; font-weight: 600; color: #6366f1; text-decoration: none; }
.db-card-link:hover { text-decoration: underline; }
.db-card-body { padding: 18px; }

/* Donut */
.db-donut-wrap { height: 160px; margin-bottom: 14px; }
.db-donut-legend { display: flex; flex-direction: column; gap: 5px; }
.db-legend-item { display: flex; align-items: center; gap: 8px; font-size: .78rem; }
.db-legend-dot { width: 8px; height: 8px; border-radius: 50%; flex-shrink: 0; }
.db-legend-label { flex: 1; color: #475569; font-weight: 500; }
.db-legend-val { font-weight: 800; color: #0f172a; }

/* App Tiles */
.db-app-grid { display: grid; grid-template-columns: repeat(2, 1fr); gap: 8px; margin-bottom: 12px; }
.db-app-tile { display: flex; flex-direction: column; align-items: center; padding: 12px 6px; border-radius: 10px; text-align: center; }
.db-app-tile i { font-size: 1.1rem; margin-bottom: 4px; }
.db-app-tile.pending { background: #fef3c7; color: #92400e; }
.db-app-tile.review { background: #dbeafe; color: #1e40af; }
.db-app-tile.approved { background: #dcfce7; color: #166534; }
.db-app-tile.rejected { background: #fee2e2; color: #991b1b; }
.db-app-num { font-size: 1.2rem; font-weight: 800; }
.db-app-lbl { font-size: .65rem; font-weight: 600; margin-top: 1px; }
.db-app-footer { display: flex; justify-content: space-between; padding: 8px 12px; background: #f8fafc; border-radius: 8px; font-size: .78rem; font-weight: 600; color: #64748b; }
.fw800 { font-weight: 800; color: #0f172a; }

/* Finance */
.db-fin-rows { display: flex; flex-direction: column; gap: 12px; margin-bottom: 14px; }
.db-fin-row { display: flex; align-items: center; gap: 12px; }
.db-fin-ic { width: 36px; height: 36px; border-radius: 8px; display: flex; align-items: center; justify-content: center; font-size: 17px; flex-shrink: 0; }
.db-fin-ic.green { background: rgba(34,197,94,.1); color: #22c55e; }
.db-fin-ic.amber { background: rgba(245,158,11,.1); color: #f59e0b; }
.db-fin-ic.blue { background: rgba(59,130,246,.1); color: #3b82f6; }
.db-fin-amt { font-size: .88rem; font-weight: 800; color: #0f172a; }
.db-fin-lbl { font-size: .68rem; color: #94a3b8; }
.db-fin-progress { margin-top: 2px; }
.db-fin-bar { height: 6px; background: #f1f5f9; border-radius: 10px; overflow: hidden; margin-bottom: 4px; }
.db-fin-fill { height: 100%; background: linear-gradient(90deg, #22c55e, #10b981); border-radius: 10px; transition: width .6s; }
.db-fin-info { display: flex; justify-content: space-between; font-size: .68rem; color: #94a3b8; }

/* Lead Pipeline */
.db-lead-list { display: flex; flex-direction: column; gap: 8px; }
.db-lead-row { display: flex; align-items: center; gap: 8px; }
.db-lead-dot { width: 8px; height: 8px; border-radius: 50%; flex-shrink: 0; }
.db-lead-name { font-size: .75rem; font-weight: 600; color: #475569; min-width: 70px; text-transform: uppercase; font-size: .68rem; }
.db-lead-bar-track { flex: 1; height: 6px; background: #f1f5f9; border-radius: 10px; overflow: hidden; }
.db-lead-bar-fill { height: 100%; border-radius: 10px; transition: width .5s; min-width: 4px; }
.db-lead-count { font-size: .82rem; font-weight: 800; color: #0f172a; min-width: 22px; text-align: right; }

/* Quick Actions */
.db-actions { display: flex; flex-wrap: wrap; gap: 8px; margin-top: 4px; }
.db-act { display: inline-flex; align-items: center; gap: 6px; padding: 8px 16px; border-radius: 8px; font-size: .78rem; font-weight: 600; color: #475569; text-decoration: none; background: #fff; border: 1px solid #e2e8f0; transition: all .15s; }
.db-act:hover { background: #6366f1; color: #fff; border-color: #6366f1; }
.db-act i { font-size: .95rem; }

/* Empty */
.db-empty { text-align: center; padding: 24px; color: #cbd5e1; display: flex; flex-direction: column; align-items: center; gap: 6px; }
.db-empty.sm { padding: 16px; }
.db-empty i { font-size: 1.8rem; opacity: .3; }
.db-empty span { font-size: .78rem; }

/* Responsive */
@media (max-width: 1280px) { .db-row-3 { grid-template-columns: 1fr 1fr; } .db-row-2-1 { grid-template-columns: 1fr 1fr; } }
@media (max-width: 1024px) { .db-kpi-row { grid-template-columns: repeat(3, 1fr); } }
@media (max-width: 768px) {
  .db-kpi-row { grid-template-columns: repeat(2, 1fr); }
  .db-row-2-1, .db-row-3 { grid-template-columns: 1fr; }
  .db-hero { padding: 18px; }
  .db-hero-inner { flex-direction: column; gap: 6px; align-items: flex-start; }
}
@media (max-width: 480px) { .db-kpi-row { grid-template-columns: 1fr; } }
</style>
