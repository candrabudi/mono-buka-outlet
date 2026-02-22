<template>
  <div class="out-page">
    <!-- Hero -->
    <div class="out-hero">
      <div class="out-hero-top">
        <div>
          <h1 class="out-hero-title">Jelajahi Outlet</h1>
          <p class="out-hero-sub">Temukan peluang kemitraan terbaik untuk Anda</p>
        </div>
      </div>
      <div class="out-stats">
        <div class="out-stat"><span class="out-stat-dot" style="background:#818cf8"></span><span class="out-stat-label">Total Outlet</span><span class="out-stat-val">{{ total }}</span></div>
        <div class="out-stat"><span class="out-stat-dot" style="background:#22c55e"></span><span class="out-stat-label">Ditampilkan</span><span class="out-stat-val">{{ outlets.length }}</span></div>
      </div>
    </div>

    <!-- Search Bar -->
    <div class="out-toolbar">
      <div class="out-search-wrap">
        <svg class="out-search-icon" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><circle cx="11" cy="11" r="8"/><line x1="21" y1="21" x2="16.65" y2="16.65"/></svg>
        <input v-model="search" type="text" placeholder="Cari nama outlet atau kategori..." class="out-search-input" @input="debouncedFetch" />
        <button v-if="search" class="out-search-clear" @click="search=''; fetchOutlets()">
          <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><line x1="18" y1="6" x2="6" y2="18"/><line x1="6" y1="6" x2="18" y2="18"/></svg>
        </button>
      </div>
    </div>

    <!-- Loading -->
    <div v-if="loading" class="out-loading">
      <div class="out-skeleton-grid">
        <div v-for="n in 6" :key="n" class="out-skeleton-card">
          <div class="out-skeleton-banner shimmer"></div>
          <div class="out-skeleton-body">
            <div class="out-skeleton-line w60 shimmer"></div>
            <div class="out-skeleton-line w40 shimmer"></div>
            <div class="out-skeleton-line w80 shimmer"></div>
          </div>
        </div>
      </div>
    </div>

    <!-- Empty -->
    <div v-else-if="outlets.length === 0" class="out-empty-wrap">
      <div class="out-empty-card">
        <div class="out-empty-circle">
          <svg width="36" height="36" viewBox="0 0 24 24" fill="none" stroke="#94a3b8" stroke-width="1.5"><path d="M3 9l9-7 9 7v11a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2z"/><polyline points="9 22 9 12 15 12 15 22"/></svg>
        </div>
        <h3>{{ search ? 'Tidak ditemukan' : 'Belum ada outlet' }}</h3>
        <p>{{ search ? `Tidak ada outlet yang cocok dengan "${search}"` : 'Outlet yang tersedia akan ditampilkan di sini' }}</p>
        <button v-if="search" class="out-empty-btn" @click="search=''; fetchOutlets()">Reset Pencarian</button>
      </div>
    </div>

    <!-- Grid -->
    <div v-else class="out-grid">
      <div v-for="o in outlets" :key="o.id" class="out-card" @click="$router.push(`/outlets/${o.id}`)">
        <!-- Banner -->
        <div class="out-card-banner">
          <img v-if="o.banner" :src="o.banner" :alt="o.name" />
          <div v-else class="out-card-banner-ph">
            <svg width="36" height="36" viewBox="0 0 24 24" fill="none" stroke="#d6d3d1" stroke-width="1"><path d="M3 9l9-7 9 7v11a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2z"/><polyline points="9 22 9 12 15 12 15 22"/></svg>
          </div>
          <div class="out-card-banner-overlay"></div>
          <div class="out-card-banner-top">
            <span v-if="o.is_featured" class="out-card-feat">
              <svg width="10" height="10" viewBox="0 0 24 24" fill="currentColor"><path d="M12 2L15.09 8.26L22 9.27L17 14.14L18.18 21.02L12 17.77L5.82 21.02L7 14.14L2 9.27L8.91 8.26L12 2Z"/></svg>
              Featured
            </span>
          </div>
          <div class="out-card-banner-bot">
            <span class="out-card-invest">{{ fc(o.minimum_investment) }}</span>
          </div>
        </div>

        <!-- Body -->
        <div class="out-card-body">
          <div class="out-card-head">
            <div class="out-card-logo">
              <img v-if="o.logo" :src="o.logo" :alt="o.name" />
              <svg v-else width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="#94a3b8" stroke-width="1.5"><path d="M3 9l9-7 9 7v11a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2z"/></svg>
            </div>
            <div class="out-card-title-wrap">
              <h3 class="out-card-name">{{ o.name }}</h3>
              <span class="out-card-cat">{{ o.category_name || o.category || 'Franchise' }}</span>
            </div>
          </div>

          <p class="out-card-desc">{{ o.short_description || (o.description ? o.description.substring(0, 120) + '...' : 'Belum ada deskripsi') }}</p>

          <div class="out-card-chips">
            <span v-if="o.estimated_roi" class="out-chip"><svg width="10" height="10" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5"><polyline points="23 6 13.5 15.5 8.5 10.5 1 18"/></svg> ROI {{ o.estimated_roi }}</span>
            <span v-if="o.profit_sharing_percentage" class="out-chip"><svg width="10" height="10" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5"><circle cx="12" cy="12" r="10"/></svg> {{ o.profit_sharing_percentage }}%</span>
            <span v-if="o.total_outlets" class="out-chip"><svg width="10" height="10" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5"><path d="M3 9l9-7 9 7v11a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2z"/></svg> {{ o.total_outlets }} outlet</span>
          </div>

          <div class="out-card-footer">
            <span class="out-card-loc" v-if="o.city || o.province">
              <svg width="11" height="11" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M21 10c0 7-9 13-9 13s-9-6-9-13a9 9 0 0 1 18 0z"/><circle cx="12" cy="10" r="3"/></svg>
              {{ [o.city, o.province].filter(Boolean).join(', ') }}
            </span>
            <span class="out-card-cta">Lihat Detail →</span>
          </div>
        </div>
      </div>
    </div>

    <!-- Pagination -->
    <div v-if="total > limit && !loading" class="out-pagination">
      <button :disabled="page <= 1" @click="page--; fetchOutlets()" class="out-pag-btn">
        <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><polyline points="15 18 9 12 15 6"/></svg>
        Sebelumnya
      </button>
      <div class="out-pag-pages">
        <button v-for="p in pagesToShow" :key="p" class="out-pag-num" :class="{active: p === page}" @click="page = p; fetchOutlets()">{{ p }}</button>
      </div>
      <button :disabled="page >= Math.ceil(total / limit)" @click="page++; fetchOutlets()" class="out-pag-btn">
        Selanjutnya
        <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><polyline points="9 18 15 12 9 6"/></svg>
      </button>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { outletApi } from '../services/api'

const outlets = ref([])
const loading = ref(true)
const search = ref('')
const page = ref(1)
const limit = ref(9)
const total = ref(0)
let debounceTimer = null

function fc(v) { return v ? new Intl.NumberFormat('id-ID', { style: 'currency', currency: 'IDR', minimumFractionDigits: 0 }).format(v) : '-' }

const pagesToShow = computed(() => {
  const totalPages = Math.ceil(total.value / limit.value)
  const pages = []
  const start = Math.max(1, page.value - 2)
  const end = Math.min(totalPages, start + 4)
  for (let i = start; i <= end; i++) pages.push(i)
  return pages
})

async function fetchOutlets() {
  loading.value = true
  try {
    const { data } = await outletApi.list({ search: search.value, page: page.value, limit: limit.value })
    outlets.value = data.data || []
    total.value = data.total || data.data?.length || 0
  } catch { outlets.value = [] }
  finally { loading.value = false }
}

function debouncedFetch() {
  clearTimeout(debounceTimer)
  debounceTimer = setTimeout(() => { page.value = 1; fetchOutlets() }, 400)
}

onMounted(fetchOutlets)
</script>

<style scoped>
/* ═══ Hero ═══ */
.out-hero{background:linear-gradient(135deg,#0f0c29 0%,#302b63 50%,#24243e 100%);border-radius:16px;padding:32px 36px 24px;margin-bottom:24px;box-shadow:0 4px 24px rgba(15,12,41,.2)}
.out-hero-top{margin-bottom:20px}
.out-hero-title{font-size:1.6rem;font-weight:800;color:#fff;margin:0 0 4px}
.out-hero-sub{font-size:.85rem;color:rgba(255,255,255,.5);margin:0}
.out-stats{display:flex;gap:28px;flex-wrap:wrap;padding-top:16px;border-top:1px solid rgba(255,255,255,.08)}
.out-stat{display:flex;align-items:center;gap:8px}
.out-stat-dot{width:8px;height:8px;border-radius:50%}
.out-stat-label{font-size:.72rem;color:rgba(255,255,255,.4);text-transform:uppercase;letter-spacing:.05em}
.out-stat-val{font-size:.9rem;font-weight:800;color:#fff}

/* ═══ Toolbar ═══ */
.out-toolbar{margin-bottom:24px}
.out-search-wrap{position:relative;max-width:420px}
.out-search-icon{position:absolute;left:14px;top:50%;transform:translateY(-50%);color:#94a3b8;pointer-events:none}
.out-search-input{width:100%;height:46px;border:1px solid #e8ecf1;border-radius:12px;padding:0 40px 0 40px;font-size:.85rem;color:#1e293b;background:#fff;box-sizing:border-box;transition:all .2s}
.out-search-input:focus{outline:none;border-color:#818cf8;box-shadow:0 0 0 3px rgba(129,140,248,.12)}
.out-search-clear{position:absolute;right:12px;top:50%;transform:translateY(-50%);width:26px;height:26px;border:none;background:#f1f5f9;border-radius:6px;cursor:pointer;display:flex;align-items:center;justify-content:center;color:#64748b;transition:all .15s}
.out-search-clear:hover{background:#e2e8f0;color:#1e293b}

/* ═══ Loading Skeleton ═══ */
.out-skeleton-grid{display:grid;grid-template-columns:repeat(auto-fill,minmax(300px,1fr));gap:20px}
.out-skeleton-card{background:#fff;border-radius:16px;border:1px solid #e8ecf1;overflow:hidden}
.out-skeleton-banner{height:160px;background:#f1f5f9}
.out-skeleton-body{padding:18px}
.out-skeleton-line{height:12px;border-radius:6px;background:#f1f5f9;margin-bottom:10px}
.out-skeleton-line.w60{width:60%}.out-skeleton-line.w40{width:40%}.out-skeleton-line.w80{width:80%}
.shimmer{background:linear-gradient(90deg,#f1f5f9 25%,#e8ecf1 50%,#f1f5f9 75%);background-size:200% 100%;animation:shimmer 1.5s infinite}
@keyframes shimmer{0%{background-position:200% 0}100%{background-position:-200% 0}}

/* ═══ Empty ═══ */
.out-empty-wrap{display:flex;justify-content:center;padding:20px 0}
.out-empty-card{text-align:center;background:#fff;border-radius:16px;border:1px solid #e8ecf1;padding:56px 40px;max-width:420px;width:100%}
.out-empty-circle{width:80px;height:80px;border-radius:50%;background:linear-gradient(135deg,#f1f5f9,#e2e8f0);display:flex;align-items:center;justify-content:center;margin:0 auto 20px}
.out-empty-card h3{font-size:1rem;font-weight:700;color:#0f172a;margin:0 0 6px}
.out-empty-card p{font-size:.85rem;color:#94a3b8;margin:0 0 20px;line-height:1.5}
.out-empty-btn{padding:8px 20px;border:1px solid #e8ecf1;border-radius:10px;background:#fff;font-size:.82rem;font-weight:600;color:#6366f1;cursor:pointer;transition:all .15s}
.out-empty-btn:hover{border-color:#818cf8;background:#f5f3ff}

/* ═══ Grid ═══ */
.out-grid{display:grid;grid-template-columns:repeat(auto-fill,minmax(310px,1fr));gap:20px}

/* ═══ Card ═══ */
.out-card{background:#fff;border-radius:16px;border:1px solid #e8ecf1;overflow:hidden;cursor:pointer;transition:all .25s ease}
.out-card:hover{box-shadow:0 12px 32px rgba(99,102,241,.12);transform:translateY(-4px);border-color:#c7d2fe}

.out-card-banner{height:160px;position:relative;overflow:hidden;background:#f8fafc}
.out-card-banner img{width:100%;height:100%;object-fit:cover;transition:transform .4s ease}
.out-card:hover .out-card-banner img{transform:scale(1.05)}
.out-card-banner-ph{height:100%;display:flex;align-items:center;justify-content:center;background:linear-gradient(135deg,#fafaf9,#f5f5f4)}
.out-card-banner-overlay{position:absolute;inset:0;background:linear-gradient(to top,rgba(0,0,0,.5) 0%,transparent 60%);pointer-events:none}
.out-card-banner-top{position:absolute;top:10px;right:10px;display:flex;gap:6px}
.out-card-feat{background:linear-gradient(135deg,#f59e0b,#f97316);color:#fff;padding:4px 10px;border-radius:20px;font-size:.62rem;font-weight:700;display:inline-flex;align-items:center;gap:3px;text-transform:uppercase;letter-spacing:.04em;box-shadow:0 2px 8px rgba(245,158,11,.35)}
.out-card-banner-bot{position:absolute;bottom:10px;left:14px}
.out-card-invest{font-size:.9rem;font-weight:800;color:#fff;text-shadow:0 1px 4px rgba(0,0,0,.3)}

.out-card-body{padding:16px 18px 18px}
.out-card-head{display:flex;align-items:center;gap:10px;margin-bottom:10px}
.out-card-logo{width:38px;height:38px;border-radius:10px;background:#f1f5f9;display:flex;align-items:center;justify-content:center;overflow:hidden;flex-shrink:0;border:1px solid #e8ecf1}
.out-card-logo img{width:100%;height:100%;object-fit:cover}
.out-card-title-wrap{min-width:0}
.out-card-name{font-size:.92rem;font-weight:700;color:#0f172a;line-height:1.25;white-space:nowrap;overflow:hidden;text-overflow:ellipsis}
.out-card-cat{font-size:.7rem;color:#94a3b8;text-transform:capitalize}

.out-card-desc{font-size:.8rem;color:#64748b;line-height:1.55;margin:0 0 12px;display:-webkit-box;-webkit-line-clamp:2;line-clamp:2;-webkit-box-orient:vertical;overflow:hidden;min-height:2.4em}

.out-card-chips{display:flex;flex-wrap:wrap;gap:6px;margin-bottom:14px}
.out-chip{display:inline-flex;align-items:center;gap:4px;padding:3px 9px;border-radius:6px;font-size:.65rem;font-weight:700;background:#f5f3ff;color:#6366f1;text-transform:uppercase;letter-spacing:.03em}

.out-card-footer{display:flex;align-items:center;justify-content:space-between;padding-top:12px;border-top:1px solid #f1f5f9}
.out-card-loc{display:flex;align-items:center;gap:3px;font-size:.72rem;color:#94a3b8}
.out-card-cta{font-size:.75rem;font-weight:700;color:#6366f1;transition:color .15s}
.out-card:hover .out-card-cta{color:#4f46e5}

/* ═══ Pagination ═══ */
.out-pagination{display:flex;align-items:center;justify-content:center;gap:8px;margin-top:32px;padding:16px 0}
.out-pag-btn{display:inline-flex;align-items:center;gap:6px;padding:8px 16px;border:1px solid #e8ecf1;border-radius:10px;background:#fff;cursor:pointer;font-size:.82rem;font-weight:600;color:#475569;transition:all .15s}
.out-pag-btn:hover:not(:disabled){border-color:#818cf8;color:#6366f1}
.out-pag-btn:disabled{opacity:.35;cursor:not-allowed}
.out-pag-pages{display:flex;gap:4px}
.out-pag-num{width:36px;height:36px;border:1px solid #e8ecf1;border-radius:10px;background:#fff;cursor:pointer;font-size:.82rem;font-weight:600;color:#475569;transition:all .15s;display:flex;align-items:center;justify-content:center}
.out-pag-num:hover{border-color:#818cf8;color:#6366f1}
.out-pag-num.active{background:linear-gradient(135deg,#6366f1,#8b5cf6);color:#fff;border-color:transparent;box-shadow:0 2px 8px rgba(99,102,241,.3)}

@media(max-width:768px){.out-hero{padding:24px 20px 18px}.out-grid{grid-template-columns:1fr}.out-skeleton-grid{grid-template-columns:1fr}}
</style>
