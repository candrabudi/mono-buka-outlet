<template>
  <div class="md-page">
    <!-- Hero -->
    <div class="md-hero">
      <div class="md-hero-content">
        <div class="md-hero-left">
          <div class="md-avatar" :style="avatarBg">{{ auth.userInitial }}</div>
          <div>
            <h1 class="md-hero-title">Selamat Datang, {{ auth.userName }} 👋</h1>
            <p class="md-hero-sub">Pantau perkembangan kemitraan & outlet Anda</p>
          </div>
        </div>
        <div class="md-hero-right" v-if="partnership">
          <div class="md-ring-wrap">
            <svg viewBox="0 0 80 80" class="md-ring-svg">
              <circle cx="40" cy="40" r="34" fill="none" stroke="rgba(255,255,255,.12)" stroke-width="6"/>
              <circle cx="40" cy="40" r="34" fill="none" stroke="url(#mdGrad)" stroke-width="6" stroke-linecap="round"
                :stroke-dasharray="213.6" :stroke-dashoffset="213.6 - (213.6 * (partnership.progress_percentage||0) / 100)"
                transform="rotate(-90 40 40)"/>
              <defs><linearGradient id="mdGrad" x1="0" y1="0" x2="1" y2="1">
                <stop offset="0%" stop-color="#667eea"/><stop offset="100%" stop-color="#764ba2"/>
              </linearGradient></defs>
            </svg>
            <div class="md-ring-center">
              <div class="md-ring-val">{{ partnership.progress_percentage||0 }}%</div>
              <div class="md-ring-label">Progress</div>
            </div>
          </div>
          <span class="md-status-badge" :class="'st-'+partnership.status">{{ statusLabel(partnership.status) }}</span>
        </div>
      </div>
    </div>

    <!-- Loading -->
    <div v-if="loading" class="md-loading"><i class="ri-loader-4-line ri-spin"></i> Memuat data...</div>

    <template v-else>
      <!-- KPI Cards -->
      <div class="md-kpi-row">
        <div class="md-kpi" @click="$router.push('/invoices')">
          <div class="md-kpi-icon ki-blue"><i class="ri-bill-line"></i></div>
          <div class="md-kpi-body">
            <div class="md-kpi-label">Invoice</div>
            <div class="md-kpi-val">{{ invoices.length }}</div>
          </div>
        </div>
        <div class="md-kpi" @click="$router.push('/applications')">
          <div class="md-kpi-icon ki-violet"><i class="ri-file-list-3-line"></i></div>
          <div class="md-kpi-body">
            <div class="md-kpi-label">Pengajuan</div>
            <div class="md-kpi-val">{{ applications.length }}</div>
          </div>
        </div>
        <div class="md-kpi" @click="$router.push('/agreements')">
          <div class="md-kpi-icon ki-green"><i class="ri-file-shield-2-line"></i></div>
          <div class="md-kpi-body">
            <div class="md-kpi-label">Perjanjian</div>
            <div class="md-kpi-val">{{ agreements.length }}</div>
          </div>
        </div>
        <div class="md-kpi" @click="$router.push('/locations')">
          <div class="md-kpi-icon ki-amber"><i class="ri-map-pin-2-line"></i></div>
          <div class="md-kpi-body">
            <div class="md-kpi-label">Lokasi</div>
            <div class="md-kpi-val">{{ locations.length }}</div>
          </div>
        </div>
      </div>

      <!-- Row: Outlet Info + Partnership Info -->
      <div class="md-grid-2">
        <!-- Outlet Saya -->
        <div class="md-card">
          <div class="md-card-head">
            <h3><i class="ri-store-2-fill"></i> Outlet Saya</h3>
            <router-link v-if="partnership?.outlet" :to="'/outlets/' + partnership.outlet.id" class="md-card-link">Detail →</router-link>
          </div>
          <div class="md-card-body" v-if="partnership?.outlet">
            <div class="md-outlet-hero">
              <div class="md-outlet-avatar"><i class="ri-store-3-fill"></i></div>
              <div>
                <div class="md-outlet-name">{{ partnership.outlet.name }}</div>
                <div class="md-outlet-desc">{{ partnership.outlet.description || 'Outlet Franchise' }}</div>
              </div>
            </div>
            <div class="md-detail-grid">
              <div class="md-detail-item">
                <span class="md-dl"><i class="ri-inbox-line"></i> Paket</span>
                <span class="md-dv">{{ partnership.package?.name || '-' }}</span>
              </div>
              <div class="md-detail-item">
                <span class="md-dl"><i class="ri-money-dollar-circle-line"></i> Harga Paket</span>
                <span class="md-dv">{{ fc(partnership.package?.price) }}</span>
              </div>
              <div class="md-detail-item">
                <span class="md-dl"><i class="ri-calendar-line"></i> Tanggal Mulai</span>
                <span class="md-dv">{{ partnership.start_date ? fmtDate(partnership.start_date) : 'Belum dimulai' }}</span>
              </div>
              <div class="md-detail-item">
                <span class="md-dl"><i class="ri-star-line"></i> Affiliator</span>
                <span class="md-dv">{{ partnership.affiliator?.name || '-' }}</span>
              </div>
            </div>
          </div>
          <div class="md-card-body md-empty" v-else>
            <i class="ri-store-2-line"></i>
            <p>Belum ada outlet terdaftar</p>
            <router-link to="/outlets" class="md-btn-primary-sm">Jelajahi Outlet</router-link>
          </div>
        </div>

        <!-- Ringkasan Keuangan -->
        <div class="md-card">
          <div class="md-card-head">
            <h3><i class="ri-wallet-3-fill"></i> Ringkasan Keuangan</h3>
            <router-link to="/invoices" class="md-card-link">Lihat Invoice →</router-link>
          </div>
          <div class="md-card-body">
            <div class="md-fin-rows">
              <div class="md-fin-row">
                <div class="md-fin-ic fi-green"><i class="ri-checkbox-circle-fill"></i></div>
                <div class="md-fin-info">
                  <div class="md-fin-label">Terbayar</div>
                  <div class="md-fin-val">{{ fc(paidAmount) }}</div>
                </div>
                <span class="md-fin-count">{{ paidCount }} invoice</span>
              </div>
              <div class="md-fin-row">
                <div class="md-fin-ic fi-amber"><i class="ri-time-fill"></i></div>
                <div class="md-fin-info">
                  <div class="md-fin-label">Menunggu Bayar</div>
                  <div class="md-fin-val">{{ fc(pendingAmount) }}</div>
                </div>
                <span class="md-fin-count">{{ pendingCount }} invoice</span>
              </div>
            </div>
            <div class="md-fin-progress" v-if="invoices.length">
              <div class="md-fin-bar"><div class="md-fin-fill" :style="{width: paymentPct + '%'}"></div></div>
              <div class="md-fin-pct">{{ paymentPct }}% lunas</div>
            </div>
          </div>
        </div>
      </div>

      <!-- Status Pipeline -->
      <div class="md-card" v-if="partnership" style="margin-bottom:20px">
        <div class="md-card-head">
          <h3><i class="ri-flow-chart"></i> Tahapan Kemitraan</h3>
        </div>
        <div class="md-card-body">
          <div class="md-pipeline">
            <div v-for="(st, idx) in pipelineSteps" :key="st.val" class="md-pipe-step"
              :class="{ active: partnership.status === st.val, done: pipeIdx(partnership.status) > idx }">
              <div class="md-pipe-connector" v-if="idx > 0" :class="{ filled: pipeIdx(partnership.status) >= idx }"></div>
              <div class="md-pipe-dot" :style="dotStyle(st, idx)">
                <i :class="pipeIdx(partnership.status) > idx ? 'ri-check-line' : st.icon"></i>
              </div>
              <div class="md-pipe-label">{{ st.label }}</div>
            </div>
          </div>
        </div>
      </div>

      <!-- Quick Actions -->
      <div class="md-card" style="margin-bottom:20px">
        <div class="md-card-head"><h3><i class="ri-apps-2-fill"></i> Menu Cepat</h3></div>
        <div class="md-card-body">
          <div class="md-actions">
            <router-link to="/outlets" class="md-action"><i class="ri-store-2-line"></i><span>Jelajahi Outlet</span></router-link>
            <router-link to="/applications" class="md-action"><i class="ri-file-list-3-line"></i><span>Pengajuan</span></router-link>
            <router-link to="/invoices" class="md-action"><i class="ri-bill-line"></i><span>Invoice</span></router-link>
            <router-link to="/agreements" class="md-action"><i class="ri-file-shield-2-line"></i><span>Perjanjian</span></router-link>
            <router-link to="/locations" class="md-action"><i class="ri-map-pin-2-line"></i><span>Lokasi</span></router-link>
            <router-link to="/ebooks" class="md-action"><i class="ri-book-open-line"></i><span>Ebook</span></router-link>
            <router-link to="/settings" class="md-action"><i class="ri-settings-3-line"></i><span>Pengaturan</span></router-link>
          </div>
        </div>
      </div>
    </template>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useAuthStore } from '../stores/auth'
import { partnershipApi, invoiceApi, applicationApi, agreementApi, locationApi } from '../services/api'

const auth = useAuthStore()
const loading = ref(true)
const partnership = ref(null)
const invoices = ref([])
const applications = ref([])
const agreements = ref([])
const locations = ref([])

onMounted(async () => {
  const results = await Promise.allSettled([
    partnershipApi.getMine(),
    invoiceApi.list(),
    applicationApi.myList(),
    agreementApi.list(),
    locationApi.list(),
  ])
  const pData = results[0].status === 'fulfilled' ? results[0].value.data.data : null
  partnership.value = Array.isArray(pData) ? (pData[0] || null) : pData
  invoices.value = results[1].status === 'fulfilled' ? (results[1].value.data.data || []) : []
  applications.value = results[2].status === 'fulfilled' ? (results[2].value.data.data || []) : []
  agreements.value = results[3].status === 'fulfilled' ? (results[3].value.data.data || []) : []
  locations.value = results[4].status === 'fulfilled' ? (results[4].value.data.data || []) : []
  loading.value = false
})

const avatarBg = computed(() => {
  const colors = ['linear-gradient(135deg,#667eea,#764ba2)','linear-gradient(135deg,#f093fb,#f5576c)','linear-gradient(135deg,#4facfe,#00f2fe)','linear-gradient(135deg,#43e97b,#38f9d7)']
  return { background: colors[(auth.userName||'?').charCodeAt(0) % colors.length] }
})

const paidInvoices = computed(() => invoices.value.filter(i => i.status === 'PAID'))
const pendingInvoices = computed(() => invoices.value.filter(i => i.status === 'PENDING'))
const paidAmount = computed(() => paidInvoices.value.reduce((s, i) => s + (i.amount || 0), 0))
const pendingAmount = computed(() => pendingInvoices.value.reduce((s, i) => s + (i.amount || 0), 0))
const paidCount = computed(() => paidInvoices.value.length)
const pendingCount = computed(() => pendingInvoices.value.length)
const paymentPct = computed(() => {
  const total = invoices.value.length
  return total ? Math.round((paidCount.value / total) * 100) : 0
})

const pipelineSteps = [
  { val:'PENDING', label:'Menunggu', icon:'ri-time-line', color:'#f59e0b' },
  { val:'DP_VERIFIED', label:'DP Terverifikasi', icon:'ri-money-dollar-circle-line', color:'#0ea5e9' },
  { val:'AGREEMENT_SIGNED', label:'Perjanjian', icon:'ri-file-text-line', color:'#8b5cf6' },
  { val:'DEVELOPMENT', label:'Pembangunan', icon:'ri-building-2-line', color:'#6366f1' },
  { val:'RUNNING', label:'Berjalan', icon:'ri-run-line', color:'#22c55e' },
  { val:'COMPLETED', label:'Selesai', icon:'ri-checkbox-circle-fill', color:'#10b981' },
]
function pipeIdx(s) { return pipelineSteps.findIndex(st => st.val === s) }
function dotStyle(st, idx) {
  const cur = pipeIdx(partnership.value?.status)
  if (idx === cur) return { background: st.color, color: '#fff', boxShadow: `0 0 0 4px ${st.color}22` }
  if (idx < cur) return { background: '#22c55e', color: '#fff' }
  return {}
}
function statusLabel(s) {
  return { PENDING:'Menunggu', DP_VERIFIED:'DP Terverifikasi', AGREEMENT_SIGNED:'Perjanjian Ditandatangani', DEVELOPMENT:'Pembangunan', RUNNING:'Berjalan', COMPLETED:'Selesai' }[s] || s
}
function fc(n) { return 'Rp ' + (n||0).toLocaleString('id-ID') }
function fmtDate(d) { return d ? new Date(d).toLocaleDateString('id-ID', { day:'numeric', month:'long', year:'numeric' }) : '-' }
</script>

<style scoped>
.md-page { margin:0; }

/* Hero */
.md-hero { background:linear-gradient(135deg,#0f0c29 0%,#302b63 50%,#24243e 100%); border-radius:18px; padding:28px 32px; margin-bottom:22px; color:#fff; }
.md-hero-content { display:flex; align-items:center; justify-content:space-between; gap:20px; flex-wrap:wrap; }
.md-hero-left { display:flex; align-items:center; gap:16px; }
.md-avatar { width:52px; height:52px; border-radius:14px; display:flex; align-items:center; justify-content:center; font-size:1.2rem; font-weight:800; color:#fff; flex-shrink:0; }
.md-hero-title { font-size:1.35rem; font-weight:800; margin:0 0 2px; }
.md-hero-sub { font-size:.82rem; opacity:.5; margin:0; }
.md-hero-right { display:flex; align-items:center; gap:14px; }
.md-ring-wrap { position:relative; width:72px; height:72px; }
.md-ring-svg { width:72px; height:72px; filter:drop-shadow(0 2px 6px rgba(0,0,0,.2)); }
.md-ring-center { position:absolute; inset:0; display:flex; flex-direction:column; align-items:center; justify-content:center; }
.md-ring-val { font-size:1.1rem; font-weight:800; }
.md-ring-label { font-size:.5rem; opacity:.6; text-transform:uppercase; letter-spacing:.06em; }
.md-status-badge { font-size:.7rem; font-weight:700; padding:5px 12px; border-radius:8px; text-transform:uppercase; letter-spacing:.03em; }
.st-PENDING { background:rgba(245,158,11,.15); color:#f59e0b; }
.st-DP_VERIFIED { background:rgba(14,165,233,.15); color:#0ea5e9; }
.st-AGREEMENT_SIGNED { background:rgba(139,92,246,.15); color:#a78bfa; }
.st-DEVELOPMENT { background:rgba(99,102,241,.15); color:#818cf8; }
.st-RUNNING { background:rgba(34,197,94,.15); color:#22c55e; }
.st-COMPLETED { background:rgba(16,185,129,.15); color:#10b981; }

/* Loading */
.md-loading { display:flex; align-items:center; justify-content:center; gap:8px; padding:60px; color:#94a3b8; font-size:.88rem; }
.ri-spin { animation:spin .8s linear infinite; }
@keyframes spin { to { transform:rotate(360deg) } }

/* KPI */
.md-kpi-row { display:grid; grid-template-columns:repeat(4,1fr); gap:14px; margin-bottom:22px; }
.md-kpi { background:#fff; border:1px solid #eef1f6; border-radius:14px; padding:18px 20px; display:flex; align-items:center; gap:14px; cursor:pointer; transition:all .2s; }
.md-kpi:hover { box-shadow:0 4px 16px rgba(0,0,0,.06); transform:translateY(-2px); }
.md-kpi-icon { width:42px; height:42px; border-radius:12px; display:flex; align-items:center; justify-content:center; font-size:1.15rem; flex-shrink:0; }
.ki-blue { background:#dbeafe; color:#2563eb; }
.ki-violet { background:#ede9fe; color:#7c3aed; }
.ki-green { background:#dcfce7; color:#16a34a; }
.ki-amber { background:#fef3c7; color:#d97706; }
.md-kpi-label { font-size:.68rem; color:#94a3b8; font-weight:600; text-transform:uppercase; letter-spacing:.04em; }
.md-kpi-val { font-size:1.2rem; font-weight:800; color:#0f172a; }

/* Grid */
.md-grid-2 { display:grid; grid-template-columns:1fr 1fr; gap:18px; margin-bottom:22px; }

/* Card */
.md-card { background:#fff; border:1px solid #eef1f6; border-radius:16px; overflow:hidden; }
.md-card-head { display:flex; justify-content:space-between; align-items:center; padding:16px 22px; border-bottom:1px solid #f1f5f9; }
.md-card-head h3 { font-size:.88rem; font-weight:700; color:#0f172a; margin:0; display:flex; align-items:center; gap:6px; }
.md-card-head h3 i { color:#6366f1; font-size:1rem; }
.md-card-link { font-size:.75rem; font-weight:600; color:#6366f1; text-decoration:none; }
.md-card-link:hover { text-decoration:underline; }
.md-card-body { padding:20px 22px; }
.md-empty { text-align:center; padding:32px 20px; color:#94a3b8; }
.md-empty i { font-size:2.2rem; display:block; margin-bottom:8px; opacity:.4; }
.md-empty p { font-size:.85rem; margin:0 0 14px; }
.md-btn-primary-sm { display:inline-flex; align-items:center; gap:5px; padding:9px 20px; font-size:.8rem; font-weight:700; border-radius:10px; background:linear-gradient(135deg,#6366f1,#8b5cf6); color:#fff; text-decoration:none; transition:all .15s; }
.md-btn-primary-sm:hover { box-shadow:0 4px 14px rgba(99,102,241,.3); }

/* Outlet Info */
.md-outlet-hero { display:flex; align-items:center; gap:14px; margin-bottom:16px; padding-bottom:14px; border-bottom:1px solid #f1f5f9; }
.md-outlet-avatar { width:48px; height:48px; border-radius:14px; background:linear-gradient(135deg,#6366f1,#8b5cf6); color:#fff; display:flex; align-items:center; justify-content:center; font-size:1.3rem; flex-shrink:0; }
.md-outlet-name { font-size:1rem; font-weight:800; color:#0f172a; }
.md-outlet-desc { font-size:.75rem; color:#94a3b8; margin-top:2px; }
.md-detail-grid { display:grid; grid-template-columns:1fr 1fr; gap:8px; }
.md-detail-item { display:flex; justify-content:space-between; padding:10px 12px; background:#f8fafc; border-radius:8px; }
.md-dl { font-size:.75rem; color:#64748b; font-weight:500; display:flex; align-items:center; gap:4px; }
.md-dl i { font-size:.8rem; color:#94a3b8; }
.md-dv { font-size:.78rem; font-weight:700; color:#0f172a; }

/* Finance */
.md-fin-rows { display:flex; flex-direction:column; gap:12px; margin-bottom:16px; }
.md-fin-row { display:flex; align-items:center; gap:12px; }
.md-fin-ic { width:38px; height:38px; border-radius:10px; display:flex; align-items:center; justify-content:center; font-size:1.1rem; flex-shrink:0; }
.fi-green { background:rgba(34,197,94,.1); color:#22c55e; }
.fi-amber { background:rgba(245,158,11,.1); color:#f59e0b; }
.md-fin-info { flex:1; }
.md-fin-label { font-size:.68rem; color:#94a3b8; font-weight:600; }
.md-fin-val { font-size:.95rem; font-weight:800; color:#0f172a; }
.md-fin-count { font-size:.7rem; color:#94a3b8; font-weight:600; white-space:nowrap; }
.md-fin-progress { }
.md-fin-bar { height:8px; background:#e2e8f0; border-radius:4px; overflow:hidden; margin-bottom:4px; }
.md-fin-fill { height:100%; background:linear-gradient(90deg,#22c55e,#10b981); border-radius:4px; transition:width .5s; }
.md-fin-pct { font-size:.72rem; font-weight:700; color:#22c55e; }

/* Pipeline */
.md-pipeline { display:flex; align-items:flex-start; gap:0; }
.md-pipe-step { position:relative; display:flex; flex-direction:column; align-items:center; gap:6px; flex:1; padding:8px 4px; }
.md-pipe-connector { position:absolute; top:20px; left:-50%; width:100%; height:2px; background:#e2e8f0; z-index:0; }
.md-pipe-connector.filled { background:#22c55e; }
.md-pipe-dot { position:relative; z-index:1; width:34px; height:34px; border-radius:50%; display:flex; align-items:center; justify-content:center; font-size:.85rem; background:#f1f5f9; color:#94a3b8; transition:all .25s; }
.md-pipe-step.active .md-pipe-dot { transform:scale(1.1); }
.md-pipe-label { font-size:.62rem; font-weight:700; color:#64748b; text-align:center; line-height:1.2; }
.md-pipe-step.active .md-pipe-label { color:#0f172a; }
.md-pipe-step.done .md-pipe-label { color:#94a3b8; }

/* Quick Actions */
.md-actions { display:grid; grid-template-columns:repeat(auto-fill,minmax(120px,1fr)); gap:10px; }
.md-action { display:flex; flex-direction:column; align-items:center; gap:8px; padding:18px 12px; border-radius:14px; background:#f8fafc; border:1px solid #eef1f6; text-decoration:none; color:#334155; font-size:.75rem; font-weight:600; transition:all .2s; }
.md-action:hover { background:#eef2ff; border-color:#c7d2fe; color:#6366f1; transform:translateY(-2px); }
.md-action i { font-size:1.4rem; color:#6366f1; }

/* Responsive */
@media (max-width:1024px) { .md-kpi-row { grid-template-columns:repeat(2,1fr); } }
@media (max-width:768px) {
  .md-hero { padding:20px; }
  .md-hero-content { flex-direction:column; text-align:center; }
  .md-hero-left { flex-direction:column; }
  .md-grid-2 { grid-template-columns:1fr; }
  .md-kpi-row { grid-template-columns:1fr 1fr; }
  .md-detail-grid { grid-template-columns:1fr; }
}
</style>
