<template>
  <div class="app-page">
    <!-- Hero -->
    <div class="app-hero">
      <div class="app-hero-top">
        <div>
          <h1 class="app-hero-title">Pengajuan Kemitraan</h1>
          <p class="app-hero-sub">Pantau status pengajuan kemitraan Anda</p>
        </div>
        <router-link to="/outlets" class="app-new-btn">
          <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5"><line x1="12" y1="5" x2="12" y2="19"/><line x1="5" y1="12" x2="19" y2="12"/></svg>
          Ajukan Baru
        </router-link>
      </div>
      <div class="app-stats">
        <div class="app-stat"><span class="app-stat-dot" style="background:#818cf8"></span><span class="app-stat-label">Total</span><span class="app-stat-val">{{ apps.length }}</span></div>
        <div class="app-stat"><span class="app-stat-dot" style="background:#f59e0b"></span><span class="app-stat-label">Menunggu</span><span class="app-stat-val">{{ apps.filter(a=>a.status==='PENDING').length }}</span></div>
        <div class="app-stat"><span class="app-stat-dot" style="background:#22c55e"></span><span class="app-stat-label">Disetujui</span><span class="app-stat-val">{{ apps.filter(a=>a.status==='APPROVED').length }}</span></div>
        <div class="app-stat"><span class="app-stat-dot" style="background:#ef4444"></span><span class="app-stat-label">Ditolak</span><span class="app-stat-val">{{ apps.filter(a=>a.status==='REJECTED').length }}</span></div>
      </div>
    </div>

    <!-- Loading -->
    <div v-if="loading" class="app-loading">
      <div v-for="n in 3" :key="n" class="app-skel shimmer"></div>
    </div>

    <!-- Empty -->
    <div v-else-if="apps.length === 0" class="app-empty-wrap">
      <div class="app-empty-card">
        <div class="app-empty-circle">
          <svg width="36" height="36" viewBox="0 0 24 24" fill="none" stroke="#94a3b8" stroke-width="1.5"><path d="M9 5H7a2 2 0 0 0-2 2v12a2 2 0 0 0 2 2h10a2 2 0 0 0 2-2V7a2 2 0 0 0-2-2h-2"/><rect x="9" y="3" width="6" height="4" rx="1"/></svg>
        </div>
        <h3>Belum ada pengajuan</h3>
        <p>Mulai jelajahi outlet untuk mengajukan kemitraan</p>
        <router-link to="/outlets" class="app-empty-btn">
          <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><circle cx="11" cy="11" r="8"/><line x1="21" y1="21" x2="16.65" y2="16.65"/></svg>
          Jelajahi Outlet
        </router-link>
      </div>
    </div>

    <!-- Cards List -->
    <div v-else class="app-list">
      <div v-for="a in apps" :key="a.id" class="app-card" @click="$router.push(`/applications/${a.id}`)" style="cursor:pointer">
        <div class="app-card-left">
          <div class="app-card-logo">
            <img v-if="a.outlet?.logo" :src="a.outlet.logo" :alt="a.outlet?.name" />
            <svg v-else width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="#94a3b8" stroke-width="1.5"><path d="M3 9l9-7 9 7v11a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2z"/></svg>
          </div>
          <div class="app-card-info">
            <div class="app-card-outlet">{{ a.outlet?.name || '-' }}</div>
            <div class="app-card-pkg">{{ a.package?.name || '-' }} · {{ fc(a.package?.price) }}</div>
          </div>
        </div>
        <div class="app-card-center">
          <div class="app-card-detail" v-if="a.proposed_location">
            <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M21 10c0 7-9 13-9 13s-9-6-9-13a9 9 0 0 1 18 0z"/><circle cx="12" cy="10" r="3"/></svg>
            {{ a.proposed_location }}
          </div>
          <div class="app-card-detail" v-if="a.investment_budget">
            <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><line x1="12" y1="1" x2="12" y2="23"/><path d="M17 5H9.5a3.5 3.5 0 0 0 0 7h5a3.5 3.5 0 0 1 0 7H6"/></svg>
            {{ fc(a.investment_budget) }}
          </div>
        </div>
        <div class="app-card-right">
          <span class="app-badge" :class="'as-'+a.status">
            <span class="app-badge-dot"></span>
            {{ statusLabel(a.status) }}
          </span>
          <span class="app-card-date">{{ formatDate(a.created_at) }}</span>
        </div>

        <!-- Admin Notes -->
        <div v-if="a.admin_notes" class="app-card-notes">
          <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M21 15a2 2 0 0 1-2 2H7l-4 4V5a2 2 0 0 1 2-2h14a2 2 0 0 1 2 2z"/></svg>
          <span>{{ a.admin_notes }}</span>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { applicationApi } from '../services/api'
import { useToastStore } from '../stores/toast'

const toast = useToastStore()
const apps = ref([])
const loading = ref(true)

function fc(v) { return v ? new Intl.NumberFormat('id-ID', { style: 'currency', currency: 'IDR', minimumFractionDigits: 0 }).format(v) : '-' }
function formatDate(d) { return d ? new Date(d).toLocaleDateString('id-ID', { day: 'numeric', month: 'short', year: 'numeric' }) : '-' }
function statusLabel(s) { return { PENDING: 'Menunggu', REVIEWED: 'Direview', APPROVED: 'Disetujui', REJECTED: 'Ditolak' }[s] || s }

onMounted(async () => {
  try { const { data } = await applicationApi.myList(); apps.value = data.data || [] }
  catch { toast.error('Gagal memuat pengajuan') }
  finally { loading.value = false }
})
</script>

<style scoped>
/* ═══ Hero ═══ */
.app-hero{background:linear-gradient(135deg,#0f0c29 0%,#302b63 50%,#24243e 100%);border-radius:16px;padding:32px 36px 24px;margin-bottom:24px;box-shadow:0 4px 24px rgba(15,12,41,.2)}
.app-hero-top{display:flex;justify-content:space-between;align-items:flex-start;margin-bottom:20px;flex-wrap:wrap;gap:12px}
.app-hero-title{font-size:1.6rem;font-weight:800;color:#fff;margin:0 0 4px}
.app-hero-sub{font-size:.85rem;color:rgba(255,255,255,.5);margin:0}
.app-new-btn{display:inline-flex;align-items:center;gap:6px;padding:8px 18px;background:rgba(255,255,255,.08);border:1px solid rgba(255,255,255,.12);border-radius:10px;color:#fff;font-size:.78rem;font-weight:600;text-decoration:none;transition:all .15s;backdrop-filter:blur(4px)}
.app-new-btn:hover{background:rgba(255,255,255,.14);border-color:rgba(255,255,255,.25)}
.app-stats{display:flex;gap:28px;flex-wrap:wrap;padding-top:16px;border-top:1px solid rgba(255,255,255,.08)}
.app-stat{display:flex;align-items:center;gap:8px}
.app-stat-dot{width:8px;height:8px;border-radius:50%}
.app-stat-label{font-size:.72rem;color:rgba(255,255,255,.4);text-transform:uppercase;letter-spacing:.05em}
.app-stat-val{font-size:.9rem;font-weight:800;color:#fff}

/* ═══ Loading ═══ */
.app-loading{display:flex;flex-direction:column;gap:14px}
.app-skel{height:90px;border-radius:14px;background:#f1f5f9}
.shimmer{background:linear-gradient(90deg,#f1f5f9 25%,#e8ecf1 50%,#f1f5f9 75%);background-size:200% 100%;animation:shimmer 1.5s infinite}
@keyframes shimmer{0%{background-position:200% 0}100%{background-position:-200% 0}}

/* ═══ Empty ═══ */
.app-empty-wrap{display:flex;justify-content:center;padding:20px 0}
.app-empty-card{text-align:center;background:#fff;border-radius:16px;border:1px solid #e8ecf1;padding:56px 40px;max-width:420px;width:100%}
.app-empty-circle{width:80px;height:80px;border-radius:50%;background:linear-gradient(135deg,#f1f5f9,#e2e8f0);display:flex;align-items:center;justify-content:center;margin:0 auto 20px}
.app-empty-card h3{font-size:1rem;font-weight:700;color:#0f172a;margin:0 0 6px}
.app-empty-card p{font-size:.85rem;color:#94a3b8;margin:0 0 20px;line-height:1.5}
.app-empty-btn{display:inline-flex;align-items:center;gap:6px;padding:10px 22px;border:none;border-radius:10px;background:linear-gradient(135deg,#6366f1,#8b5cf6);color:#fff;font-size:.82rem;font-weight:600;text-decoration:none;transition:all .2s;box-shadow:0 2px 8px rgba(99,102,241,.25)}
.app-empty-btn:hover{box-shadow:0 4px 14px rgba(99,102,241,.35);transform:translateY(-1px)}

/* ═══ Cards ═══ */
.app-list{display:flex;flex-direction:column;gap:14px}
.app-card{background:#fff;border-radius:14px;border:1px solid #e8ecf1;padding:18px 22px;display:grid;grid-template-columns:1fr 1fr auto;align-items:center;gap:16px;transition:all .2s}
.app-card:hover{box-shadow:0 4px 16px rgba(0,0,0,.06);transform:translateY(-1px)}

.app-card-left{display:flex;align-items:center;gap:12px}
.app-card-logo{width:42px;height:42px;border-radius:11px;background:#f1f5f9;border:1px solid #e8ecf1;display:flex;align-items:center;justify-content:center;overflow:hidden;flex-shrink:0}
.app-card-logo img{width:100%;height:100%;object-fit:cover}
.app-card-outlet{font-size:.88rem;font-weight:700;color:#0f172a}
.app-card-pkg{font-size:.75rem;color:#64748b;margin-top:2px}

.app-card-center{display:flex;flex-direction:column;gap:4px}
.app-card-detail{display:flex;align-items:center;gap:5px;font-size:.78rem;color:#64748b}
.app-card-detail svg{color:#94a3b8;flex-shrink:0}

.app-card-right{display:flex;flex-direction:column;align-items:flex-end;gap:6px}
.app-card-date{font-size:.72rem;color:#94a3b8}

/* ═══ Badge ═══ */
.app-badge{font-size:.68rem;font-weight:700;padding:4px 12px;border-radius:6px;text-transform:uppercase;letter-spacing:.03em;display:inline-flex;align-items:center;gap:5px;white-space:nowrap}
.app-badge-dot{width:6px;height:6px;border-radius:50%}
.as-PENDING{background:#fef3c7;color:#d97706}
.as-PENDING .app-badge-dot{background:#f59e0b}
.as-REVIEWED{background:#dbeafe;color:#1e40af}
.as-REVIEWED .app-badge-dot{background:#3b82f6}
.as-APPROVED{background:#dcfce7;color:#16a34a}
.as-APPROVED .app-badge-dot{background:#22c55e}
.as-REJECTED{background:#fee2e2;color:#991b1b}
.as-REJECTED .app-badge-dot{background:#ef4444}

/* ═══ Notes ═══ */
.app-card-notes{grid-column:1/-1;display:flex;align-items:flex-start;gap:8px;padding:12px 14px;background:#f8fafc;border-radius:8px;font-size:.8rem;color:#475569;line-height:1.5;margin-top:4px}
.app-card-notes svg{color:#6366f1;flex-shrink:0;margin-top:2px}

@media(max-width:768px){.app-hero{padding:24px 20px 18px}.app-card{grid-template-columns:1fr;gap:10px}.app-card-right{align-items:flex-start;flex-direction:row;gap:12px}}
</style>
