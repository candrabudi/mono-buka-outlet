<template>
  <div class="ad-page">
    <!-- Loading -->
    <div v-if="loading" class="ad-loading">
      <div class="ad-skel-hero shimmer"></div>
      <div class="ad-skel-row"><div class="ad-skel-card shimmer"></div><div class="ad-skel-card shimmer"></div></div>
    </div>

    <template v-else-if="app">
      <!-- Hero -->
      <div class="ad-hero">
        <button class="ad-back" @click="$router.push('/applications')">
          <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5"><path d="M19 12H5"/><polyline points="12 19 5 12 12 5"/></svg>
          Kembali
        </button>
        <div class="ad-hero-content">
          <div class="ad-hero-logo" v-if="app.outlet?.logo">
            <img :src="app.outlet.logo" :alt="app.outlet?.name" />
          </div>
          <div>
            <div class="ad-hero-badges">
              <span class="ad-badge" :class="'as-'+app.status">
                <span class="ad-badge-dot"></span>
                {{ statusLabel(app.status) }}
              </span>
            </div>
            <h1 class="ad-hero-title">{{ app.outlet?.name || 'Pengajuan Kemitraan' }}</h1>
            <p class="ad-hero-sub">
              <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><rect x="3" y="4" width="18" height="18" rx="2" ry="2"/><line x1="16" y1="2" x2="16" y2="6"/><line x1="8" y1="2" x2="8" y2="6"/><line x1="3" y1="10" x2="21" y2="10"/></svg>
              Diajukan {{ formatDate(app.created_at) }}
            </p>
          </div>
        </div>
        <div class="ad-stats">
          <div class="ad-stat" v-if="app.package?.name">
            <span class="ad-stat-label">Paket</span>
            <span class="ad-stat-val">{{ app.package.name }}</span>
          </div>
          <div class="ad-stat" v-if="app.package?.price">
            <span class="ad-stat-label">Harga Paket</span>
            <span class="ad-stat-val">{{ fc(app.package.price) }}</span>
          </div>
          <div class="ad-stat" v-if="app.investment_budget">
            <span class="ad-stat-label">Budget Investasi</span>
            <span class="ad-stat-val">{{ fc(app.investment_budget) }}</span>
          </div>
        </div>
      </div>

      <!-- Content -->
      <div class="ad-content">
        <!-- Left: Detail Info -->
        <div class="ad-main">
          <!-- Outlet Info -->
          <div class="ad-section">
            <h3 class="ad-section-title">
              <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M3 9l9-7 9 7v11a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2z"/><polyline points="9 22 9 12 15 12 15 22"/></svg>
              Informasi Outlet
            </h3>
            <div class="ad-info-grid">
              <div class="ad-info-item">
                <span class="ad-info-label">Nama Outlet</span>
                <span class="ad-info-value">{{ app.outlet?.name || '-' }}</span>
              </div>
              <div class="ad-info-item" v-if="app.outlet?.category_name">
                <span class="ad-info-label">Kategori</span>
                <span class="ad-info-value">{{ app.outlet.category_name }}</span>
              </div>
              <div class="ad-info-item" v-if="app.outlet?.city || app.outlet?.province">
                <span class="ad-info-label">Lokasi Outlet</span>
                <span class="ad-info-value">{{ [app.outlet?.city, app.outlet?.province].filter(Boolean).join(', ') }}</span>
              </div>
            </div>
          </div>

          <!-- Paket -->
          <div class="ad-section" v-if="app.package">
            <h3 class="ad-section-title">
              <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M20.59 13.41l-7.17 7.17a2 2 0 0 1-2.83 0L2 12V2h10l8.59 8.59a2 2 0 0 1 0 2.82z"/><line x1="7" y1="7" x2="7.01" y2="7"/></svg>
              Paket Kemitraan
            </h3>
            <div class="ad-info-grid">
              <div class="ad-info-item">
                <span class="ad-info-label">Nama Paket</span>
                <span class="ad-info-value">{{ app.package.name }}</span>
              </div>
              <div class="ad-info-item">
                <span class="ad-info-label">Harga</span>
                <span class="ad-info-value ad-info-highlight">{{ fc(app.package.price) }}</span>
              </div>
              <div class="ad-info-item" v-if="app.package.duration">
                <span class="ad-info-label">Durasi</span>
                <span class="ad-info-value">{{ app.package.duration }}</span>
              </div>
              <div class="ad-info-item" v-if="app.package.estimated_bep">
                <span class="ad-info-label">Estimasi BEP</span>
                <span class="ad-info-value">{{ app.package.estimated_bep }}</span>
              </div>
              <div class="ad-info-item" v-if="app.package.minimum_dp">
                <span class="ad-info-label">Minimum DP</span>
                <span class="ad-info-value">{{ fc(app.package.minimum_dp) }}</span>
              </div>
              <div class="ad-info-item" v-if="app.package.net_profit">
                <span class="ad-info-label">Estimasi Laba</span>
                <span class="ad-info-value">{{ app.package.net_profit }}</span>
              </div>
            </div>
            <p class="ad-pkg-desc" v-if="app.package.description">{{ app.package.description }}</p>
          </div>

          <!-- Lokasi & Budget -->
          <div class="ad-section">
            <h3 class="ad-section-title">
              <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M21 10c0 7-9 13-9 13s-9-6-9-13a9 9 0 0 1 18 0z"/><circle cx="12" cy="10" r="3"/></svg>
              Lokasi & Budget
            </h3>
            <div class="ad-info-grid">
              <div class="ad-info-item" v-if="app.proposed_location">
                <span class="ad-info-label">Lokasi yang Diusulkan</span>
                <span class="ad-info-value">{{ app.proposed_location }}</span>
              </div>
              <div class="ad-info-item" v-if="app.investment_budget">
                <span class="ad-info-label">Budget Investasi</span>
                <span class="ad-info-value ad-info-highlight">{{ fc(app.investment_budget) }}</span>
              </div>
            </div>
          </div>

          <!-- Kontak -->
          <div class="ad-section" v-if="app.contact_phone || app.contact_email">
            <h3 class="ad-section-title">
              <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M22 16.92v3a2 2 0 0 1-2.18 2 19.79 19.79 0 0 1-8.63-3.07 19.5 19.5 0 0 1-6-6 19.79 19.79 0 0 1-3.07-8.67A2 2 0 0 1 4.11 2h3a2 2 0 0 1 2 1.72"/></svg>
              Kontak yang Bisa Dihubungi
            </h3>
            <div class="ad-info-grid">
              <div class="ad-info-item" v-if="app.contact_phone">
                <span class="ad-info-label">No. Handphone</span>
                <span class="ad-info-value">{{ app.contact_phone }}</span>
              </div>
              <div class="ad-info-item" v-if="app.contact_email">
                <span class="ad-info-label">Email</span>
                <span class="ad-info-value">{{ app.contact_email }}</span>
              </div>
            </div>
          </div>
          <div class="ad-section" v-if="app.motivation || app.experience">
            <h3 class="ad-section-title">
              <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><circle cx="12" cy="12" r="10"/><line x1="12" y1="16" x2="12" y2="12"/><line x1="12" y1="8" x2="12.01" y2="8"/></svg>
              Motivasi & Pengalaman
            </h3>
            <div class="ad-text-section" v-if="app.motivation">
              <h4 class="ad-text-label">Motivasi</h4>
              <p class="ad-text-value">{{ app.motivation }}</p>
            </div>
            <div class="ad-text-section" v-if="app.experience">
              <h4 class="ad-text-label">Pengalaman</h4>
              <p class="ad-text-value">{{ app.experience }}</p>
            </div>
          </div>
        </div>

        <!-- Right: Status Timeline -->
        <div class="ad-sidebar">
          <!-- Status Card -->
          <div class="ad-status-card">
            <h3 class="ad-status-title">Status Pengajuan</h3>
            <div class="ad-timeline">
              <div class="ad-tl-item" :class="{ active: true, done: true }">
                <div class="ad-tl-dot"></div>
                <div class="ad-tl-content">
                  <div class="ad-tl-label">Pengajuan Dikirim</div>
                  <div class="ad-tl-date">{{ formatDate(app.created_at) }}</div>
                </div>
              </div>
              <div class="ad-tl-item" :class="{ active: isReviewed, done: isReviewed }">
                <div class="ad-tl-dot"></div>
                <div class="ad-tl-content">
                  <div class="ad-tl-label">Sedang Direview</div>
                  <div class="ad-tl-date" v-if="app.reviewed_at">{{ formatDate(app.reviewed_at) }}</div>
                  <div class="ad-tl-date" v-else>Menunggu</div>
                </div>
              </div>
              <div class="ad-tl-item" :class="{ active: isFinal, done: isFinal, approved: app.status === 'APPROVED', rejected: app.status === 'REJECTED' }">
                <div class="ad-tl-dot"></div>
                <div class="ad-tl-content">
                  <div class="ad-tl-label">{{ app.status === 'REJECTED' ? 'Ditolak' : app.status === 'APPROVED' ? 'Disetujui' : 'Keputusan' }}</div>
                  <div class="ad-tl-date" v-if="app.reviewed_at && isFinal">{{ formatDate(app.reviewed_at) }}</div>
                  <div class="ad-tl-date" v-else>Menunggu</div>
                </div>
              </div>
            </div>
          </div>

          <!-- Admin Notes -->
          <div class="ad-notes-card" v-if="app.admin_notes">
            <h3 class="ad-notes-title">
              <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M21 15a2 2 0 0 1-2 2H7l-4 4V5a2 2 0 0 1 2-2h14a2 2 0 0 1 2 2z"/></svg>
              Catatan Admin
            </h3>
            <p class="ad-notes-text">{{ app.admin_notes }}</p>
          </div>

          <!-- Quick Actions -->
          <div class="ad-actions-card">
            <router-link :to="`/outlets/${app.outlet_id}`" class="ad-action-btn">
              <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M3 9l9-7 9 7v11a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2z"/><polyline points="9 22 9 12 15 12 15 22"/></svg>
              Lihat Outlet
            </router-link>
            <router-link to="/applications" class="ad-action-btn ad-action-ghost">
              <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><polyline points="15 18 9 12 15 6"/></svg>
              Semua Pengajuan
            </router-link>
          </div>
        </div>
      </div>
    </template>

    <!-- Not Found -->
    <div v-else class="ad-not-found">
      <div class="ad-empty-circle">
        <svg width="32" height="32" viewBox="0 0 24 24" fill="none" stroke="#94a3b8" stroke-width="1.5"><circle cx="12" cy="12" r="10"/><line x1="15" y1="9" x2="9" y2="15"/><line x1="9" y1="9" x2="15" y2="15"/></svg>
      </div>
      <p>Pengajuan tidak ditemukan</p>
      <router-link to="/applications" class="ad-back-link">Kembali ke Daftar</router-link>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { applicationApi } from '../services/api'
import { useToastStore } from '../stores/toast'

const route = useRoute()
const toast = useToastStore()
const app = ref(null)
const loading = ref(true)

function fc(v) { return v ? new Intl.NumberFormat('id-ID', { style: 'currency', currency: 'IDR', minimumFractionDigits: 0 }).format(v) : '-' }
function formatDate(d) { return d ? new Date(d).toLocaleDateString('id-ID', { day: 'numeric', month: 'long', year: 'numeric' }) : '-' }
function statusLabel(s) { return { PENDING: 'Menunggu', REVIEWED: 'Direview', APPROVED: 'Disetujui', REJECTED: 'Ditolak' }[s] || s }

const isReviewed = computed(() => ['REVIEWED', 'APPROVED', 'REJECTED'].includes(app.value?.status))
const isFinal = computed(() => ['APPROVED', 'REJECTED'].includes(app.value?.status))

onMounted(async () => {
  try {
    const { data } = await applicationApi.getByID(route.params.id)
    app.value = data.data
  } catch {
    app.value = null
    toast.error('Gagal memuat detail pengajuan')
  } finally {
    loading.value = false
  }
})
</script>

<style scoped>
/* ═══ Hero ═══ */
.ad-hero{background:linear-gradient(135deg,#0f0c29 0%,#302b63 50%,#24243e 100%);border-radius:16px;padding:28px 36px 24px;margin-bottom:24px;box-shadow:0 4px 24px rgba(15,12,41,.2)}
.ad-back{display:inline-flex;align-items:center;gap:5px;font-size:.78rem;font-weight:600;color:rgba(255,255,255,.45);background:rgba(255,255,255,.06);border:1px solid rgba(255,255,255,.1);border-radius:8px;padding:6px 14px;cursor:pointer;margin-bottom:18px;transition:all .15s;backdrop-filter:blur(4px)}
.ad-back:hover{color:#fff;border-color:rgba(255,255,255,.25);background:rgba(255,255,255,.1)}
.ad-hero-content{display:flex;align-items:center;gap:16px;margin-bottom:20px}
.ad-hero-logo{width:56px;height:56px;border-radius:14px;overflow:hidden;flex-shrink:0;border:2px solid rgba(255,255,255,.15);background:rgba(255,255,255,.08)}
.ad-hero-logo img{width:100%;height:100%;object-fit:cover}
.ad-hero-badges{display:flex;gap:8px;margin-bottom:6px}
.ad-hero-title{font-size:1.5rem;font-weight:800;color:#fff;margin:0}
.ad-hero-sub{font-size:.82rem;color:rgba(255,255,255,.4);margin:6px 0 0;display:flex;align-items:center;gap:5px}
.ad-stats{display:flex;gap:32px;flex-wrap:wrap;padding-top:18px;border-top:1px solid rgba(255,255,255,.08)}
.ad-stat{display:flex;flex-direction:column;gap:2px}
.ad-stat-label{font-size:.68rem;color:rgba(255,255,255,.35);text-transform:uppercase;letter-spacing:.05em}
.ad-stat-val{font-size:.95rem;font-weight:800;color:#fff}

/* ═══ Badge ═══ */
.ad-badge{font-size:.68rem;font-weight:700;padding:4px 12px;border-radius:6px;text-transform:uppercase;letter-spacing:.03em;display:inline-flex;align-items:center;gap:5px}
.ad-badge-dot{width:6px;height:6px;border-radius:50%}
.as-PENDING{background:rgba(245,158,11,.15);color:#fbbf24}
.as-PENDING .ad-badge-dot{background:#f59e0b}
.as-REVIEWED{background:rgba(59,130,246,.15);color:#60a5fa}
.as-REVIEWED .ad-badge-dot{background:#3b82f6}
.as-APPROVED{background:rgba(34,197,94,.15);color:#4ade80}
.as-APPROVED .ad-badge-dot{background:#22c55e}
.as-REJECTED{background:rgba(239,68,68,.15);color:#f87171}
.as-REJECTED .ad-badge-dot{background:#ef4444}

/* ═══ Content ═══ */
.ad-content{display:grid;grid-template-columns:1fr 340px;gap:24px;align-items:start}
.ad-main{display:flex;flex-direction:column;gap:20px}
.ad-sidebar{display:flex;flex-direction:column;gap:16px}

/* ═══ Section ═══ */
.ad-section{background:#fff;border-radius:14px;border:1px solid #e8ecf1;padding:24px}
.ad-section-title{font-size:.9rem;font-weight:700;color:#0f172a;margin:0 0 16px;display:flex;align-items:center;gap:8px}
.ad-section-title svg{color:#6366f1}
.ad-info-grid{display:grid;grid-template-columns:1fr 1fr;gap:14px}
.ad-info-item{display:flex;flex-direction:column;gap:2px}
.ad-info-label{font-size:.7rem;color:#94a3b8;font-weight:600;text-transform:uppercase;letter-spacing:.04em}
.ad-info-value{font-size:.88rem;font-weight:600;color:#1e293b}
.ad-info-highlight{color:#6366f1;font-weight:800}
.ad-pkg-desc{font-size:.82rem;color:#64748b;line-height:1.6;margin:14px 0 0;padding-top:14px;border-top:1px solid #f1f5f9}

.ad-text-section{margin-bottom:16px}
.ad-text-section:last-child{margin-bottom:0}
.ad-text-label{font-size:.78rem;font-weight:700;color:#475569;margin:0 0 6px;text-transform:uppercase;letter-spacing:.03em}
.ad-text-value{font-size:.85rem;color:#475569;line-height:1.7;margin:0;white-space:pre-wrap}

/* ═══ Status Timeline ═══ */
.ad-status-card{background:#fff;border-radius:14px;border:1px solid #e8ecf1;padding:24px}
.ad-status-title{font-size:.9rem;font-weight:700;color:#0f172a;margin:0 0 20px}
.ad-timeline{display:flex;flex-direction:column;gap:0}
.ad-tl-item{display:flex;gap:14px;padding-bottom:24px;position:relative}
.ad-tl-item:last-child{padding-bottom:0}
.ad-tl-item:not(:last-child)::after{content:'';position:absolute;left:9px;top:22px;bottom:0;width:2px;background:#e2e8f0}
.ad-tl-item.done:not(:last-child)::after{background:#818cf8}
.ad-tl-dot{width:20px;height:20px;border-radius:50%;border:2px solid #e2e8f0;background:#fff;flex-shrink:0;position:relative;z-index:1;transition:all .2s}
.ad-tl-item.active .ad-tl-dot{border-color:#818cf8;background:#818cf8;box-shadow:0 0 0 4px rgba(129,140,248,.15)}
.ad-tl-item.active .ad-tl-dot::after{content:'';position:absolute;inset:4px;border-radius:50%;background:#fff}
.ad-tl-item.approved .ad-tl-dot{border-color:#22c55e;background:#22c55e;box-shadow:0 0 0 4px rgba(34,197,94,.15)}
.ad-tl-item.rejected .ad-tl-dot{border-color:#ef4444;background:#ef4444;box-shadow:0 0 0 4px rgba(239,68,68,.15)}
.ad-tl-content{flex:1;padding-top:1px}
.ad-tl-label{font-size:.85rem;font-weight:600;color:#0f172a}
.ad-tl-item:not(.active) .ad-tl-label{color:#94a3b8}
.ad-tl-date{font-size:.75rem;color:#94a3b8;margin-top:2px}

/* ═══ Notes ═══ */
.ad-notes-card{background:#fff;border-radius:14px;border:1px solid #e8ecf1;padding:24px}
.ad-notes-title{font-size:.9rem;font-weight:700;color:#0f172a;margin:0 0 12px;display:flex;align-items:center;gap:8px}
.ad-notes-title svg{color:#6366f1}
.ad-notes-text{font-size:.85rem;color:#475569;line-height:1.7;margin:0;padding:14px;background:#f8fafc;border-radius:10px;border-left:3px solid #818cf8}

/* ═══ Actions ═══ */
.ad-actions-card{display:flex;flex-direction:column;gap:8px}
.ad-action-btn{display:flex;align-items:center;justify-content:center;gap:6px;padding:11px 0;border-radius:10px;font-size:.85rem;font-weight:600;text-decoration:none;transition:all .2s;background:linear-gradient(135deg,#6366f1,#8b5cf6);color:#fff;box-shadow:0 2px 8px rgba(99,102,241,.25)}
.ad-action-btn:hover{box-shadow:0 4px 14px rgba(99,102,241,.35);transform:translateY(-1px)}
.ad-action-ghost{background:#fff;color:#475569;border:1px solid #e8ecf1;box-shadow:none}
.ad-action-ghost:hover{background:#f8fafc;box-shadow:0 2px 8px rgba(0,0,0,.04)}

/* ═══ Not Found ═══ */
.ad-not-found{text-align:center;padding:80px 20px;background:#fff;border-radius:16px;border:1px solid #e8ecf1}
.ad-empty-circle{width:80px;height:80px;border-radius:50%;background:linear-gradient(135deg,#f1f5f9,#e2e8f0);display:flex;align-items:center;justify-content:center;margin:0 auto 20px}
.ad-not-found p{color:#94a3b8;font-size:.85rem;margin:0 0 16px}
.ad-back-link{display:inline-flex;align-items:center;gap:6px;padding:8px 20px;border:1px solid #e8ecf1;border-radius:10px;font-size:.82rem;font-weight:600;color:#6366f1;text-decoration:none;transition:all .15s}
.ad-back-link:hover{border-color:#818cf8;background:#f5f3ff}

/* ═══ Skeleton ═══ */
.ad-skel-hero{height:200px;border-radius:16px;background:#f1f5f9;margin-bottom:24px}
.ad-skel-row{display:grid;grid-template-columns:1fr 340px;gap:24px}
.ad-skel-card{height:250px;border-radius:14px;background:#f1f5f9}
.shimmer{background:linear-gradient(90deg,#f1f5f9 25%,#e8ecf1 50%,#f1f5f9 75%);background-size:200% 100%;animation:shimmer 1.5s infinite}
@keyframes shimmer{0%{background-position:200% 0}100%{background-position:-200% 0}}

@media(max-width:900px){
  .ad-content{grid-template-columns:1fr}
  .ad-hero{padding:24px 20px 18px}
  .ad-info-grid{grid-template-columns:1fr}
  .ad-skel-row{grid-template-columns:1fr}
}
</style>
