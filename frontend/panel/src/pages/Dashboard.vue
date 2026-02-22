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

    <!-- Stat Cards -->
    <div class="db-stats">
      <div class="db-stat-card">
        <div class="db-stat-icon" style="background:rgba(99,102,241,.12);color:#6366f1"><i class="ri-store-2-line"></i></div>
        <div class="db-stat-info">
          <span class="db-stat-label">Total Outlet</span>
          <span class="db-stat-value">{{ s.total_outlets ?? 0 }}</span>
        </div>
      </div>
      <div class="db-stat-card">
        <div class="db-stat-icon" style="background:rgba(34,197,94,.12);color:#22c55e"><i class="ri-group-line"></i></div>
        <div class="db-stat-info">
          <span class="db-stat-label">Mitra Aktif</span>
          <span class="db-stat-value">{{ s.active_mitra ?? 0 }}</span>
        </div>
      </div>
      <div class="db-stat-card">
        <div class="db-stat-icon" style="background:rgba(139,92,246,.12);color:#8b5cf6"><i class="ri-handshake-line"></i></div>
        <div class="db-stat-info">
          <span class="db-stat-label">Partnership</span>
          <span class="db-stat-value">{{ s.total_partnerships ?? 0 }}</span>
        </div>
      </div>
      <div class="db-stat-card">
        <div class="db-stat-icon" style="background:rgba(249,115,22,.12);color:#f97316"><i class="ri-money-dollar-circle-line"></i></div>
        <div class="db-stat-info">
          <span class="db-stat-label">Revenue Bulan Ini</span>
          <span class="db-stat-value">{{ fc(s.monthly_revenue) }}</span>
        </div>
      </div>
    </div>

    <!-- Main Grid -->
    <div class="db-grid">
      <!-- Left Column -->
      <div class="db-main">
        <!-- Pengajuan Ringkasan -->
        <div class="db-card">
          <div class="db-card-head">
            <h3><i class="ri-file-list-3-line"></i> Pengajuan Kemitraan</h3>
            <router-link to="/applications" class="db-card-link">Lihat Semua →</router-link>
          </div>
          <div class="db-card-body">
            <div class="db-app-stats">
              <div class="db-app-stat pending">
                <span class="db-app-num">{{ s.pending_applications ?? 0 }}</span>
                <span class="db-app-label">Menunggu Review</span>
              </div>
              <div class="db-app-stat approved">
                <span class="db-app-num">{{ appStat('APPROVED') }}</span>
                <span class="db-app-label">Disetujui</span>
              </div>
              <div class="db-app-stat rejected">
                <span class="db-app-num">{{ appStat('REJECTED') }}</span>
                <span class="db-app-label">Ditolak</span>
              </div>
              <div class="db-app-stat total">
                <span class="db-app-num">{{ s.total_applications ?? 0 }}</span>
                <span class="db-app-label">Total</span>
              </div>
            </div>
          </div>
        </div>

        <!-- Partnership Status -->
        <div class="db-card">
          <div class="db-card-head">
            <h3><i class="ri-pie-chart-line"></i> Status Partnership</h3>
            <router-link to="/partnerships" class="db-card-link">Detail →</router-link>
          </div>
          <div class="db-card-body">
            <div v-if="Object.keys(s.partnerships_by_status || {}).length" class="db-bars">
              <div v-for="(count, status) in s.partnerships_by_status" :key="status" class="db-bar-row">
                <span class="db-bar-label">{{ psLabel(status) }}</span>
                <div class="db-bar-track">
                  <div class="db-bar-fill" :class="'ps-' + status" :style="{ width: barPct(count, s.total_partnerships) + '%' }"></div>
                </div>
                <span class="db-bar-count">{{ count }}</span>
              </div>
            </div>
            <div v-else class="db-empty"><i class="ri-bar-chart-box-line"></i><span>Belum ada data</span></div>
          </div>
        </div>

        <!-- Revenue Chart -->
        <div class="db-card">
          <div class="db-card-head">
            <h3><i class="ri-line-chart-line"></i> Revenue 6 Bulan Terakhir</h3>
          </div>
          <div class="db-card-body">
            <div v-if="s.revenue_chart?.length" class="db-bars">
              <div v-for="item in s.revenue_chart" :key="item.month" class="db-bar-row">
                <span class="db-bar-label" style="min-width:50px">{{ item.month }}</span>
                <div class="db-bar-track">
                  <div class="db-bar-fill rev" :style="{ width: revPct(item.revenue) + '%' }"></div>
                </div>
                <span class="db-bar-count" style="min-width:100px;text-align:right">{{ fc(item.revenue) }}</span>
              </div>
            </div>
            <div v-else class="db-empty"><i class="ri-line-chart-line"></i><span>Belum ada data revenue</span></div>
          </div>
        </div>
      </div>

      <!-- Right Sidebar -->
      <div class="db-side">
        <!-- Invoice Quick View -->
        <div class="db-card">
          <div class="db-card-head">
            <h3><i class="ri-bill-line"></i> Invoice</h3>
          </div>
          <div class="db-card-body">
            <div class="db-invoice-stats">
              <div class="db-inv-item">
                <div class="db-inv-icon pending"><i class="ri-time-line"></i></div>
                <div>
                  <div class="db-inv-num">{{ s.pending_invoices ?? 0 }}</div>
                  <div class="db-inv-label">Belum Dibayar</div>
                </div>
              </div>
              <div class="db-inv-item">
                <div class="db-inv-icon paid"><i class="ri-checkbox-circle-line"></i></div>
                <div>
                  <div class="db-inv-num">{{ s.paid_invoices ?? 0 }}</div>
                  <div class="db-inv-label">Sudah Bayar</div>
                </div>
              </div>
            </div>
          </div>
        </div>

        <!-- Lead Pipeline -->
        <div class="db-card">
          <div class="db-card-head">
            <h3><i class="ri-user-follow-line"></i> Lead Pipeline</h3>
          </div>
          <div class="db-card-body">
            <div v-if="s.leads_by_status && Object.keys(s.leads_by_status).length" class="db-lead-list">
              <div v-for="(count, status) in s.leads_by_status" :key="status" class="db-lead-row">
                <span class="db-lead-badge" :class="'ld-' + status">{{ formatStatus(status) }}</span>
                <span class="db-lead-count">{{ count }}</span>
              </div>
            </div>
            <div v-else class="db-empty"><i class="ri-user-add-line"></i><span>Belum ada leads</span></div>
          </div>
        </div>

        <!-- Quick Links -->
        <div class="db-card db-quick">
          <div class="db-card-head">
            <h3><i class="ri-flashlight-line"></i> Aksi Cepat</h3>
          </div>
          <div class="db-card-body">
            <router-link to="/outlets/create" class="db-quick-btn"><i class="ri-add-circle-line"></i> Tambah Outlet</router-link>
            <router-link to="/applications" class="db-quick-btn"><i class="ri-file-list-3-line"></i> Review Pengajuan</router-link>
            <router-link to="/partnerships" class="db-quick-btn"><i class="ri-handshake-line"></i> Lihat Partnership</router-link>
            <router-link to="/mitra" class="db-quick-btn"><i class="ri-group-line"></i> Kelola Mitra</router-link>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { dashboardApi } from '../services/api'

const s = ref({})

onMounted(async () => {
  try {
    const { data } = await dashboardApi.stats()
    s.value = data.data || {}
  } catch (e) { console.error(e) }
})

const todayDate = new Date().toLocaleDateString('id-ID', { weekday: 'long', day: 'numeric', month: 'long', year: 'numeric' })

function fc(v) { return 'Rp ' + (v || 0).toLocaleString('id-ID') }
function formatStatus(s) { return s.replace(/_/g, ' ') }
function appStat(status) { return (s.value.applications_by_status || {})[status] || 0 }
function psLabel(st) {
  const map = { PENDING: 'Pending', ACTIVE: 'Active', RUNNING: 'Running', COMPLETED: 'Selesai', CANCELLED: 'Batal' }
  return map[st] || st
}
function barPct(count, total) { return total ? Math.round((count / total) * 100) : 0 }
function revPct(r) {
  const m = Math.max(...(s.value.revenue_chart || [{ revenue: 1 }]).map(x => x.revenue))
  return m ? (r / m) * 100 : 0
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

/* Stat Cards */
.db-stats { display: grid; grid-template-columns: repeat(4, 1fr); gap: 16px; margin-bottom: 24px; }
.db-stat-card { background: #fff; border: 1px solid #e2e8f0; border-radius: 14px; padding: 20px; display: flex; align-items: center; gap: 16px; transition: all .2s; }
.db-stat-card:hover { box-shadow: 0 4px 20px rgba(0,0,0,.06); transform: translateY(-2px); }
.db-stat-icon { width: 48px; height: 48px; border-radius: 12px; display: flex; align-items: center; justify-content: center; font-size: 22px; flex-shrink: 0; }
.db-stat-label { font-size: .75rem; color: #94a3b8; font-weight: 500; text-transform: uppercase; letter-spacing: .3px; }
.db-stat-value { font-size: 1.35rem; font-weight: 800; color: #0f172a; display: block; margin-top: 2px; }

/* Grid */
.db-grid { display: grid; grid-template-columns: 1fr 380px; gap: 24px; }
.db-main { display: flex; flex-direction: column; gap: 20px; }
.db-side { display: flex; flex-direction: column; gap: 20px; }

/* Card */
.db-card { background: #fff; border: 1px solid #e2e8f0; border-radius: 14px; overflow: hidden; }
.db-card-head { display: flex; justify-content: space-between; align-items: center; padding: 16px 20px; border-bottom: 1px solid #f1f5f9; }
.db-card-head h3 { font-size: .9rem; font-weight: 700; color: #0f172a; margin: 0; display: flex; align-items: center; gap: 8px; }
.db-card-head h3 i { font-size: 1.1rem; color: #6366f1; }
.db-card-link { font-size: .78rem; font-weight: 600; color: #6366f1; text-decoration: none; }
.db-card-link:hover { text-decoration: underline; }
.db-card-body { padding: 20px; }

/* App Stats */
.db-app-stats { display: grid; grid-template-columns: repeat(4, 1fr); gap: 12px; }
.db-app-stat { text-align: center; padding: 16px 8px; border-radius: 12px; }
.db-app-stat.pending { background: #fef3c7; }
.db-app-stat.approved { background: #dcfce7; }
.db-app-stat.rejected { background: #fee2e2; }
.db-app-stat.total { background: #f1f5f9; }
.db-app-num { font-size: 1.5rem; font-weight: 800; color: #0f172a; display: block; }
.db-app-label { font-size: .72rem; font-weight: 600; color: #64748b; margin-top: 4px; display: block; }

/* Bar Charts */
.db-bars { display: flex; flex-direction: column; gap: 10px; }
.db-bar-row { display: flex; align-items: center; gap: 12px; }
.db-bar-label { font-size: .78rem; font-weight: 600; color: #475569; min-width: 70px; }
.db-bar-track { flex: 1; height: 8px; background: #f1f5f9; border-radius: 10px; overflow: hidden; }
.db-bar-fill { height: 100%; border-radius: 10px; transition: width .6s ease; min-width: 4px; }
.db-bar-fill.ps-PENDING { background: #f59e0b; }
.db-bar-fill.ps-ACTIVE, .db-bar-fill.ps-RUNNING { background: #3b82f6; }
.db-bar-fill.ps-COMPLETED { background: #22c55e; }
.db-bar-fill.ps-CANCELLED { background: #ef4444; }
.db-bar-fill.rev { background: linear-gradient(90deg, #6366f1, #8b5cf6); }
.db-bar-count { font-size: .82rem; font-weight: 700; color: #0f172a; min-width: 28px; }

/* Invoice Stats */
.db-invoice-stats { display: flex; flex-direction: column; gap: 14px; }
.db-inv-item { display: flex; align-items: center; gap: 14px; }
.db-inv-icon { width: 42px; height: 42px; border-radius: 10px; display: flex; align-items: center; justify-content: center; font-size: 20px; }
.db-inv-icon.pending { background: #fef3c7; color: #f59e0b; }
.db-inv-icon.paid { background: #dcfce7; color: #22c55e; }
.db-inv-num { font-size: 1.2rem; font-weight: 800; color: #0f172a; }
.db-inv-label { font-size: .75rem; color: #94a3b8; }

/* Lead List */
.db-lead-list { display: flex; flex-direction: column; gap: 8px; }
.db-lead-row { display: flex; justify-content: space-between; align-items: center; padding: 8px 14px; border-radius: 8px; background: #f8fafc; }
.db-lead-badge { font-size: .72rem; font-weight: 700; padding: 3px 10px; border-radius: 14px; text-transform: uppercase; letter-spacing: .3px; }
.ld-NEW { background: #dbeafe; color: #1e40af; }
.ld-RUNNING { background: #dcfce7; color: #166534; }
.ld-COMPLETED { background: #f0fdf4; color: #15803d; }
.ld-DP_PAID { background: #fef3c7; color: #92400e; }
.ld-FULLY_PAID { background: #ede9fe; color: #5b21b6; }
.db-lead-count { font-size: .9rem; font-weight: 800; color: #0f172a; }

/* Quick Links */
.db-quick .db-card-body { display: flex; flex-direction: column; gap: 8px; }
.db-quick-btn { display: flex; align-items: center; gap: 10px; padding: 10px 14px; border-radius: 10px; font-size: .82rem; font-weight: 600; color: #334155; text-decoration: none; background: #f8fafc; border: 1px solid #e2e8f0; transition: all .15s; }
.db-quick-btn:hover { background: #6366f1; color: #fff; border-color: #6366f1; }
.db-quick-btn i { font-size: 1rem; }

/* Empty */
.db-empty { text-align: center; padding: 24px; color: #94a3b8; display: flex; flex-direction: column; align-items: center; gap: 6px; }
.db-empty i { font-size: 2rem; opacity: .4; }
.db-empty span { font-size: .82rem; }

/* Responsive */
@media (max-width: 1280px) { .db-grid { grid-template-columns: 1fr; } }
@media (max-width: 768px) {
  .db-stats { grid-template-columns: repeat(2, 1fr); }
  .db-hero { padding: 20px; }
  .db-hero-inner { flex-direction: column; gap: 8px; align-items: flex-start; }
  .db-app-stats { grid-template-columns: repeat(2, 1fr); }
}
@media (max-width: 480px) { .db-stats { grid-template-columns: 1fr; } }
</style>
