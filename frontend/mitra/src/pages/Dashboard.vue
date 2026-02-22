<template>
  <div class="mt-dash">
    <!-- Hero -->
    <div class="mt-hero">
      <div class="mt-hero-top">
        <div>
          <h1 class="mt-hero-title">Selamat Datang, {{ auth.userName }}</h1>
          <p class="mt-hero-sub">Pantau perkembangan kemitraan Anda di sini</p>
        </div>
      </div>
      <div class="mt-stats">
        <div class="mt-stat">
          <span class="mt-stat-dot" style="background:#818cf8"></span>
          <span class="mt-stat-label">Status</span>
          <span class="mt-stat-val">{{ partnership?.status || '-' }}</span>
        </div>
        <div class="mt-stat">
          <span class="mt-stat-dot" style="background:#22c55e"></span>
          <span class="mt-stat-label">Progress</span>
          <span class="mt-stat-val">{{ partnership?.progress_percentage || 0 }}%</span>
        </div>
      </div>
    </div>

    <!-- Info Cards -->
    <div class="mt-info-grid">
      <div class="mt-info-card">
        <div class="mt-info-icon" style="background:#eff6ff;color:#3b82f6">
          <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M21 16V8a2 2 0 0 0-1-1.73l-7-4a2 2 0 0 0-2 0l-7 4A2 2 0 0 0 3 8v8a2 2 0 0 0 1 1.73l7 4a2 2 0 0 0 2 0l7-4A2 2 0 0 0 21 16z"/></svg>
        </div>
        <div>
          <div class="mt-info-label">Paket</div>
          <div class="mt-info-value">{{ partnership?.package?.name || '-' }}</div>
        </div>
      </div>
      <div class="mt-info-card">
        <div class="mt-info-icon" style="background:#f0fdf4;color:#22c55e">
          <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M3 9l9-7 9 7v11a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2z"/><polyline points="9 22 9 12 15 12 15 22"/></svg>
        </div>
        <div>
          <div class="mt-info-label">Outlet</div>
          <div class="mt-info-value">{{ partnership?.outlet?.name || '-' }}</div>
        </div>
      </div>
      <div class="mt-info-card">
        <div class="mt-info-icon" style="background:#faf5ff;color:#8b5cf6">
          <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M12 2L15.09 8.26L22 9.27L17 14.14L18.18 21.02L12 17.77L5.82 21.02L7 14.14L2 9.27L8.91 8.26L12 2Z"/></svg>
        </div>
        <div>
          <div class="mt-info-label">Leader</div>
          <div class="mt-info-value">{{ partnership?.leader?.name || '-' }}</div>
        </div>
      </div>
      <div class="mt-info-card">
        <div class="mt-info-icon" style="background:#fff8f1;color:#fd9644">
          <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><rect x="3" y="4" width="18" height="18" rx="2" ry="2"/><line x1="16" y1="2" x2="16" y2="6"/><line x1="8" y1="2" x2="8" y2="6"/><line x1="3" y1="10" x2="21" y2="10"/></svg>
        </div>
        <div>
          <div class="mt-info-label">Tanggal Mulai</div>
          <div class="mt-info-value">{{ partnership?.start_date ? formatDate(partnership.start_date) : 'Belum dimulai' }}</div>
        </div>
      </div>
    </div>

    <!-- Progress -->
    <div class="mt-progress-card" v-if="partnership">
      <h3>Progress Kemitraan</h3>
      <div class="mt-progress-bar-wrap">
        <div class="mt-progress-bar">
          <div class="mt-progress-fill" :style="{ width: (partnership.progress_percentage || 0) + '%' }"></div>
        </div>
        <span class="mt-progress-pct">{{ partnership.progress_percentage || 0 }}%</span>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useAuthStore } from '../stores/auth'
import { partnershipApi } from '../services/api'

const auth = useAuthStore()
const partnership = ref(null)

onMounted(async () => {
  try {
    const { data } = await partnershipApi.getMine()
    partnership.value = data.data
  } catch { /* no partnership yet */ }
})

function formatDate(d) {
  return d ? new Date(d).toLocaleDateString('id-ID', { day: 'numeric', month: 'long', year: 'numeric' }) : '-'
}
</script>

<style scoped>
.mt-hero { background: linear-gradient(135deg, #0f0c29 0%, #302b63 50%, #24243e 100%); border-radius: 16px; padding: 32px 36px 24px; margin-bottom: 24px; box-shadow: 0 4px 24px rgba(15,12,41,0.2); }
.mt-hero-top { margin-bottom: 20px; }
.mt-hero-title { font-size: 1.5rem; font-weight: 800; color: #fff; margin: 0 0 4px; }
.mt-hero-sub { font-size: .85rem; color: rgba(255,255,255,.5); margin: 0; }
.mt-stats { display: flex; gap: 28px; flex-wrap: wrap; padding-top: 16px; border-top: 1px solid rgba(255,255,255,.08); }
.mt-stat { display: flex; align-items: center; gap: 8px; }
.mt-stat-dot { width: 8px; height: 8px; border-radius: 50%; }
.mt-stat-label { font-size: .72rem; color: rgba(255,255,255,.4); text-transform: uppercase; letter-spacing: .05em; }
.mt-stat-val { font-size: .9rem; font-weight: 800; color: #fff; }

.mt-info-grid { display: grid; grid-template-columns: repeat(auto-fill, minmax(240px, 1fr)); gap: 16px; margin-bottom: 24px; }
.mt-info-card { background: #fff; border-radius: 14px; border: 1px solid #e8ecf1; padding: 20px; display: flex; align-items: center; gap: 14px; transition: all .2s; }
.mt-info-card:hover { box-shadow: 0 4px 16px rgba(0,0,0,.06); transform: translateY(-1px); }
.mt-info-icon { width: 42px; height: 42px; border-radius: 12px; display: flex; align-items: center; justify-content: center; flex-shrink: 0; }
.mt-info-label { font-size: .72rem; color: #94a3b8; font-weight: 600; text-transform: uppercase; letter-spacing: .04em; }
.mt-info-value { font-size: .92rem; font-weight: 700; color: #0f172a; margin-top: 2px; }

.mt-progress-card { background: #fff; border-radius: 14px; border: 1px solid #e8ecf1; padding: 24px; }
.mt-progress-card h3 { font-size: 1rem; font-weight: 700; margin: 0 0 16px; color: #0f172a; }
.mt-progress-bar-wrap { display: flex; align-items: center; gap: 14px; }
.mt-progress-bar { flex: 1; height: 10px; background: #e2e8f0; border-radius: 5px; overflow: hidden; }
.mt-progress-fill { height: 100%; background: linear-gradient(90deg, #6366f1, #8b5cf6); border-radius: 5px; transition: width .5s; }
.mt-progress-pct { font-size: .9rem; font-weight: 800; color: #6366f1; min-width: 40px; }

@media (max-width: 768px) {
  .mt-hero { padding: 24px 20px 18px; }
  .mt-info-grid { grid-template-columns: 1fr; }
}
</style>
