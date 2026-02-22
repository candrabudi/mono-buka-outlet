<template>
  <div class="od-page">
    <!-- Loading -->
    <div v-if="loading" class="od-loading-wrap">
      <div class="od-skeleton-hero shimmer"></div>
      <div class="od-skeleton-row"><div class="od-skeleton-card shimmer"></div><div class="od-skeleton-card shimmer"></div><div class="od-skeleton-card shimmer"></div></div>
    </div>

    <template v-else-if="outlet">
      <!-- Hero -->
      <div class="od-hero" :style="outlet.banner ? { backgroundImage: `linear-gradient(to right, rgba(15,12,41,.92) 0%, rgba(48,43,99,.85) 50%, rgba(36,36,62,.8) 100%), url(${outlet.banner})`, backgroundSize: 'cover', backgroundPosition: 'center' } : {}">
        <button class="od-back" @click="$router.push('/outlets')">
          <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5"><path d="M19 12H5"/><polyline points="12 19 5 12 12 5"/></svg>
          Kembali
        </button>
        <div class="od-hero-content">
          <div class="od-hero-logo" v-if="outlet.logo">
            <img :src="outlet.logo" :alt="outlet.name" />
          </div>
          <div>
            <div class="od-hero-badges">
              <span class="od-hero-cat">{{ outlet.category_name || outlet.category || 'Franchise' }}</span>
              <span v-if="outlet.is_featured" class="od-hero-feat">
                <svg width="10" height="10" viewBox="0 0 24 24" fill="currentColor"><path d="M12 2L15.09 8.26L22 9.27L17 14.14L18.18 21.02L12 17.77L5.82 21.02L7 14.14L2 9.27L8.91 8.26L12 2Z"/></svg>
                Featured
              </span>
            </div>
            <h1 class="od-hero-title">{{ outlet.name }}</h1>
            <p class="od-hero-sub" v-if="outlet.city || outlet.province">
              <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M21 10c0 7-9 13-9 13s-9-6-9-13a9 9 0 0 1 18 0z"/><circle cx="12" cy="10" r="3"/></svg>
              {{ [outlet.city, outlet.province].filter(Boolean).join(', ') }}
            </p>
          </div>
        </div>
        <div class="od-stats">
          <div class="od-stat">
            <span class="od-stat-num">{{ fc(outlet.minimum_investment) }}</span>
            <span class="od-stat-label">Investasi Mulai</span>
          </div>
          <div class="od-stat" v-if="outlet.estimated_roi">
            <span class="od-stat-num">{{ outlet.estimated_roi }}</span>
            <span class="od-stat-label">Est. ROI</span>
          </div>
          <div class="od-stat" v-if="outlet.profit_sharing_percentage">
            <span class="od-stat-num">{{ outlet.profit_sharing_percentage }}%</span>
            <span class="od-stat-label">Profit Sharing</span>
          </div>
          <div class="od-stat" v-if="outlet.total_outlets">
            <span class="od-stat-num">{{ outlet.total_outlets }}</span>
            <span class="od-stat-label">Total Outlet</span>
          </div>
        </div>
      </div>

      <!-- Content -->
      <div class="od-content">
        <!-- Left: Info -->
        <div class="od-main">
          <!-- About -->
          <div class="od-section-card">
            <h3 class="od-section-title">
              <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><circle cx="12" cy="12" r="10"/><line x1="12" y1="16" x2="12" y2="12"/><line x1="12" y1="8" x2="12.01" y2="8"/></svg>
              Tentang Outlet
            </h3>
            <div v-if="outlet.description" class="od-text od-html" v-html="outlet.description"></div>
            <p v-else class="od-text">Belum ada deskripsi</p>
          </div>

          <!-- Requirements -->
          <div class="od-section-card" v-if="outlet.location_requirement">
            <h3 class="od-section-title">
              <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M21 10c0 7-9 13-9 13s-9-6-9-13a9 9 0 0 1 18 0z"/><circle cx="12" cy="10" r="3"/></svg>
              Persyaratan Lokasi
            </h3>
            <p class="od-text">{{ outlet.location_requirement }}</p>
          </div>

          <!-- Contact -->
          <div class="od-section-card" v-if="outlet.contact_phone || outlet.contact_email || outlet.contact_whatsapp || outlet.website">
            <h3 class="od-section-title">
              <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M22 16.92v3a2 2 0 0 1-2.18 2 19.79 19.79 0 0 1-8.63-3.07 19.5 19.5 0 0 1-6-6 19.79 19.79 0 0 1-3.07-8.67A2 2 0 0 1 4.11 2h3a2 2 0 0 1 2 1.72"/></svg>
              Informasi Kontak
            </h3>
            <div class="od-contact-list">
              <div v-if="outlet.contact_phone" class="od-contact-item">
                <div class="od-contact-icon"><svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M22 16.92v3a2 2 0 0 1-2.18 2 19.79 19.79 0 0 1-8.63-3.07 19.5 19.5 0 0 1-6-6 19.79 19.79 0 0 1-3.07-8.67A2 2 0 0 1 4.11 2h3a2 2 0 0 1 2 1.72c.127.96.361 1.903.7 2.81a2 2 0 0 1-.45 2.11L8.09 9.91"/></svg></div>
                <div><div class="od-contact-label">Telepon</div><div class="od-contact-val">{{ outlet.contact_phone }}</div></div>
              </div>
              <div v-if="outlet.contact_email" class="od-contact-item">
                <div class="od-contact-icon"><svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M4 4h16c1.1 0 2 .9 2 2v12c0 1.1-.9 2-2 2H4c-1.1 0-2-.9-2-2V6c0-1.1.9-2 2-2z"/><polyline points="22,6 12,13 2,6"/></svg></div>
                <div><div class="od-contact-label">Email</div><div class="od-contact-val">{{ outlet.contact_email }}</div></div>
              </div>
              <div v-if="outlet.contact_whatsapp" class="od-contact-item">
                <div class="od-contact-icon wa"><svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M21 11.5a8.38 8.38 0 0 1-.9 3.8 8.5 8.5 0 0 1-7.6 4.7 8.38 8.38 0 0 1-3.8-.9L3 21l1.9-5.7a8.38 8.38 0 0 1-.9-3.8 8.5 8.5 0 0 1 4.7-7.6 8.38 8.38 0 0 1 3.8-.9h.5a8.48 8.48 0 0 1 8 8v.5z"/></svg></div>
                <div><div class="od-contact-label">WhatsApp</div><div class="od-contact-val">{{ outlet.contact_whatsapp }}</div></div>
              </div>
              <div v-if="outlet.website" class="od-contact-item">
                <div class="od-contact-icon"><svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><circle cx="12" cy="12" r="10"/><line x1="2" y1="12" x2="22" y2="12"/></svg></div>
                <div><div class="od-contact-label">Website</div><div class="od-contact-val">{{ outlet.website }}</div></div>
              </div>
            </div>
          </div>
        </div>

        <!-- Right: Packages -->
        <div class="od-sidebar">
          <div class="od-pkg-header">
            <h3>
              <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M20.59 13.41l-7.17 7.17a2 2 0 0 1-2.83 0L2 12V2h10l8.59 8.59a2 2 0 0 1 0 2.82z"/><line x1="7" y1="7" x2="7.01" y2="7"/></svg>
              Paket Kemitraan
            </h3>
            <span class="od-pkg-count">{{ packages.length }} paket</span>
          </div>

          <div v-if="loadingPkgs" class="od-pkg-loading">
            <div v-for="n in 2" :key="n" class="od-pkg-skeleton shimmer"></div>
          </div>
          <div v-else-if="packages.length === 0" class="od-pkg-empty">
            <p>Belum ada paket tersedia</p>
          </div>
          <div v-else class="od-pkg-list">
            <div v-for="(pkg, idx) in packages" :key="pkg.id" class="od-pkg-card">
              <div class="od-pkg-ribbon" v-if="idx === 0">Populer</div>
              <div class="od-pkg-top">
                <h4 class="od-pkg-name">{{ pkg.name }}</h4>
                <div class="od-pkg-price">{{ fc(pkg.price) }}</div>
              </div>
              <p class="od-pkg-desc">{{ pkg.description || '-' }}</p>

              <div class="od-pkg-specs">
                <div v-if="pkg.duration" class="od-pkg-spec">
                  <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><circle cx="12" cy="12" r="10"/><polyline points="12 6 12 12 16 14"/></svg>
                  <span>{{ pkg.duration }}</span>
                </div>
                <div v-if="pkg.estimated_bep" class="od-pkg-spec">
                  <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><polyline points="23 6 13.5 15.5 8.5 10.5 1 18"/><polyline points="17 6 23 6 23 12"/></svg>
                  <span>BEP {{ pkg.estimated_bep }}</span>
                </div>
                <div v-if="pkg.net_profit" class="od-pkg-spec">
                  <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><line x1="12" y1="1" x2="12" y2="23"/><path d="M17 5H9.5a3.5 3.5 0 0 0 0 7h5a3.5 3.5 0 0 1 0 7H6"/></svg>
                  <span>Laba {{ pkg.net_profit }}</span>
                </div>
                <div v-if="pkg.minimum_dp" class="od-pkg-spec">
                  <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><rect x="1" y="4" width="22" height="16" rx="2" ry="2"/><line x1="1" y1="10" x2="23" y2="10"/></svg>
                  <span>DP min {{ fc(pkg.minimum_dp) }}</span>
                </div>
              </div>

              <ul v-if="pkg.benefits && pkg.benefits.length" class="od-pkg-benefits">
                <li v-for="(b, i) in pkg.benefits" :key="i">
                  <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="#22c55e" stroke-width="3"><polyline points="20 6 9 17 4 12"/></svg>
                  {{ b }}
                </li>
              </ul>

              <button class="od-pkg-cta" @click="$router.push({ name: 'ApplyForm', params: { id: outlet.id }, query: { package: pkg.id } })">
                <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><line x1="22" y1="2" x2="11" y2="13"/><polygon points="22 2 15 22 11 13 2 9 22 2"/></svg>
                Ajukan Kemitraan
              </button>
            </div>
          </div>
        </div>
      </div>
    </template>

    <!-- Not Found -->
    <div v-else class="od-not-found">
      <div class="out-empty-circle">
        <svg width="32" height="32" viewBox="0 0 24 24" fill="none" stroke="#94a3b8" stroke-width="1.5"><circle cx="12" cy="12" r="10"/><line x1="15" y1="9" x2="9" y2="15"/><line x1="9" y1="9" x2="15" y2="15"/></svg>
      </div>
      <p>Outlet tidak ditemukan</p>
      <button class="out-empty-btn" @click="$router.push('/outlets')">Kembali ke Daftar</button>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { outletApi } from '../services/api'

const route = useRoute()
const outlet = ref(null)
const packages = ref([])
const loading = ref(true)
const loadingPkgs = ref(true)

function fc(v) { return v ? new Intl.NumberFormat('id-ID', { style: 'currency', currency: 'IDR', minimumFractionDigits: 0 }).format(v) : '-' }

async function fetchOutlet() {
  loading.value = true
  try { const { data } = await outletApi.getByID(route.params.id); outlet.value = data.data }
  catch { outlet.value = null }
  finally { loading.value = false }
}

async function fetchPackages() {
  loadingPkgs.value = true
  try { const { data } = await outletApi.getPackages(route.params.id); packages.value = data.data || [] }
  catch { packages.value = [] }
  finally { loadingPkgs.value = false }
}

onMounted(() => { fetchOutlet(); fetchPackages() })
</script>

<style scoped>
/* ═══ Hero ═══ */
.od-hero{background:linear-gradient(135deg,#0f0c29 0%,#302b63 50%,#24243e 100%);border-radius:16px;padding:28px 36px 24px;margin-bottom:24px;box-shadow:0 4px 24px rgba(15,12,41,.2)}
.od-back{display:inline-flex;align-items:center;gap:5px;font-size:.78rem;font-weight:600;color:rgba(255,255,255,.45);background:rgba(255,255,255,.06);border:1px solid rgba(255,255,255,.1);border-radius:8px;padding:6px 14px;cursor:pointer;margin-bottom:18px;transition:all .15s;backdrop-filter:blur(4px)}
.od-back:hover{color:#fff;border-color:rgba(255,255,255,.25);background:rgba(255,255,255,.1)}
.od-hero-content{display:flex;align-items:center;gap:16px;margin-bottom:20px}
.od-hero-logo{width:56px;height:56px;border-radius:14px;overflow:hidden;flex-shrink:0;border:2px solid rgba(255,255,255,.15);background:rgba(255,255,255,.08)}
.od-hero-logo img{width:100%;height:100%;object-fit:cover}
.od-hero-badges{display:flex;gap:8px;margin-bottom:6px}
.od-hero-cat{font-size:.65rem;font-weight:700;color:rgba(255,255,255,.5);text-transform:uppercase;letter-spacing:.06em;background:rgba(255,255,255,.08);padding:3px 10px;border-radius:20px}
.od-hero-feat{font-size:.62rem;font-weight:700;color:#fbbf24;background:rgba(251,191,36,.12);padding:3px 10px;border-radius:20px;display:inline-flex;align-items:center;gap:3px}
.od-hero-title{font-size:1.6rem;font-weight:800;color:#fff;margin:0}
.od-hero-sub{font-size:.82rem;color:rgba(255,255,255,.4);margin:6px 0 0;display:flex;align-items:center;gap:5px}
.od-stats{display:flex;gap:32px;flex-wrap:wrap;padding-top:18px;border-top:1px solid rgba(255,255,255,.08)}
.od-stat{display:flex;flex-direction:column;gap:2px}
.od-stat-num{font-size:1rem;font-weight:800;color:#fff}
.od-stat-label{font-size:.68rem;color:rgba(255,255,255,.35);text-transform:uppercase;letter-spacing:.05em}

/* ═══ Content ═══ */
.od-content{display:grid;grid-template-columns:1fr 380px;gap:24px;align-items:start}
.od-main{display:flex;flex-direction:column;gap:20px}
.od-sidebar{display:flex;flex-direction:column;gap:16px}

/* ═══ Section Card ═══ */
.od-section-card{background:#fff;border-radius:14px;border:1px solid #e8ecf1;padding:22px}
.od-section-title{font-size:.88rem;font-weight:700;color:#0f172a;margin:0 0 14px;display:flex;align-items:center;gap:8px}
.od-section-title svg{color:#6366f1}
.od-text{font-size:.85rem;color:#475569;line-height:1.7;margin:0;white-space:pre-wrap}

/* Rich HTML content from editor */
.od-html{font-size:.85rem;color:#475569;line-height:1.7}
.od-html :deep(p){margin:0 0 12px}
.od-html :deep(p:last-child){margin-bottom:0}
.od-html :deep(h1),.od-html :deep(h2),.od-html :deep(h3),.od-html :deep(h4){color:#0f172a;margin:16px 0 8px;font-weight:700}
.od-html :deep(h1){font-size:1.3rem}
.od-html :deep(h2){font-size:1.15rem}
.od-html :deep(h3){font-size:1rem}
.od-html :deep(ul),.od-html :deep(ol){padding-left:20px;margin:8px 0}
.od-html :deep(li){margin-bottom:4px}
.od-html :deep(a){color:#6366f1;text-decoration:underline}
.od-html :deep(a:hover){color:#4f46e5}
.od-html :deep(strong),.od-html :deep(b){font-weight:700;color:#334155}
.od-html :deep(img){max-width:100%;height:auto;border-radius:8px;margin:8px 0}
.od-html :deep(table){width:100%;border-collapse:collapse;margin:8px 0}
.od-html :deep(th),.od-html :deep(td){padding:8px 12px;border:1px solid #e2e8f0;font-size:.82rem;text-align:left}
.od-html :deep(th){background:#f8fafc;font-weight:600;color:#334155}
.od-html :deep(blockquote){border-left:3px solid #818cf8;padding:8px 16px;margin:8px 0;background:#f8fafc;color:#475569;font-style:italic}

/* ═══ Contact ═══ */
.od-contact-list{display:flex;flex-direction:column;gap:12px}
.od-contact-item{display:flex;align-items:center;gap:12px}
.od-contact-icon{width:36px;height:36px;border-radius:10px;background:#f5f3ff;color:#6366f1;display:flex;align-items:center;justify-content:center;flex-shrink:0}
.od-contact-icon.wa{background:#ecfdf5;color:#22c55e}
.od-contact-label{font-size:.68rem;color:#94a3b8;text-transform:uppercase;letter-spacing:.04em}
.od-contact-val{font-size:.85rem;font-weight:600;color:#1e293b}

/* ═══ Packages ═══ */
.od-pkg-header{display:flex;align-items:center;justify-content:space-between;padding:0 2px}
.od-pkg-header h3{font-size:.92rem;font-weight:700;color:#0f172a;margin:0;display:flex;align-items:center;gap:8px}
.od-pkg-header h3 svg{color:#6366f1}
.od-pkg-count{font-size:.72rem;color:#94a3b8;font-weight:600}

.od-pkg-loading{display:flex;flex-direction:column;gap:12px}
.od-pkg-skeleton{height:200px;border-radius:14px;background:#f1f5f9}
.shimmer{background:linear-gradient(90deg,#f1f5f9 25%,#e8ecf1 50%,#f1f5f9 75%);background-size:200% 100%;animation:shimmer 1.5s infinite}
@keyframes shimmer{0%{background-position:200% 0}100%{background-position:-200% 0}}

.od-pkg-empty{text-align:center;padding:40px 20px;background:#fff;border-radius:14px;border:1px solid #e8ecf1}
.od-pkg-empty p{color:#94a3b8;font-size:.85rem;margin:0}

.od-pkg-list{display:flex;flex-direction:column;gap:16px}
.od-pkg-card{background:#fff;border-radius:14px;border:1px solid #e8ecf1;padding:22px;position:relative;overflow:hidden;transition:all .2s}
.od-pkg-card:hover{box-shadow:0 4px 16px rgba(0,0,0,.06);transform:translateY(-1px)}
.od-pkg-ribbon{position:absolute;top:12px;right:-28px;background:linear-gradient(135deg,#6366f1,#8b5cf6);color:#fff;font-size:.58rem;font-weight:700;text-transform:uppercase;letter-spacing:.06em;padding:3px 32px;transform:rotate(45deg);box-shadow:0 2px 6px rgba(99,102,241,.3)}
.od-pkg-top{margin-bottom:10px}
.od-pkg-name{font-size:.95rem;font-weight:700;color:#0f172a;margin:0 0 4px}
.od-pkg-price{font-size:1.15rem;font-weight:800;color:#6366f1}
.od-pkg-desc{font-size:.8rem;color:#64748b;line-height:1.5;margin:0 0 14px}

.od-pkg-specs{display:grid;grid-template-columns:1fr 1fr;gap:8px;margin-bottom:14px}
.od-pkg-spec{display:flex;align-items:center;gap:5px;font-size:.75rem;color:#475569;font-weight:500}
.od-pkg-spec svg{color:#94a3b8;flex-shrink:0}

.od-pkg-benefits{list-style:none;padding:0;margin:0 0 16px}
.od-pkg-benefits li{display:flex;align-items:flex-start;gap:6px;font-size:.8rem;color:#475569;margin-bottom:5px;line-height:1.4}
.od-pkg-benefits svg{margin-top:2px;flex-shrink:0}

.od-pkg-cta{width:100%;padding:11px 0;border:none;border-radius:10px;background:linear-gradient(135deg,#6366f1,#8b5cf6);color:#fff;font-size:.85rem;font-weight:600;cursor:pointer;transition:all .2s;display:inline-flex;align-items:center;justify-content:center;gap:6px;box-shadow:0 2px 8px rgba(99,102,241,.25)}
.od-pkg-cta:hover{box-shadow:0 6px 20px rgba(99,102,241,.35);transform:translateY(-1px)}

/* ═══ Not Found ═══ */
.od-not-found{text-align:center;padding:56px 20px}
.out-empty-circle{width:80px;height:80px;border-radius:50%;background:linear-gradient(135deg,#f1f5f9,#e2e8f0);display:flex;align-items:center;justify-content:center;margin:0 auto 20px}
.od-not-found p{color:#94a3b8;font-size:.85rem;margin:0 0 20px}
.out-empty-btn{padding:8px 20px;border:1px solid #e8ecf1;border-radius:10px;background:#fff;font-size:.82rem;font-weight:600;color:#6366f1;cursor:pointer}

/* ═══ Skeleton ═══ */
.od-skeleton-hero{height:220px;border-radius:16px;background:#f1f5f9;margin-bottom:24px}
.od-skeleton-row{display:grid;grid-template-columns:1fr 1fr 1fr;gap:16px}
.od-skeleton-card{height:140px;border-radius:14px;background:#f1f5f9}

@media(max-width:900px){.od-content{grid-template-columns:1fr}.od-hero{padding:24px 20px 18px}}
</style>
