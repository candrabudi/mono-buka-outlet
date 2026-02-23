<template>
  <div class="animate-in">
    <div class="eb-hero">
      <div class="eb-hero-content">
        <div>
          <h1 class="eb-hero-title"><i class="ri-book-open-line"></i> Jelajahi Ebook</h1>
          <p class="eb-hero-sub">Temukan panduan & materi bermanfaat untuk bisnis Anda</p>
        </div>
      </div>
    </div>

    <div class="eb-search-bar">
      <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="#94a3b8"><circle cx="11" cy="11" r="8" stroke-width="2"/><path d="m21 21-4.35-4.35" stroke-width="2" stroke-linecap="round"/></svg>
      <input v-model="search" @input="debouncedLoad" type="text" placeholder="Cari ebook..." />
    </div>

    <div v-if="loading" class="eb-grid">
      <div v-for="n in 6" :key="n" class="eb-card skeleton-card">
        <div class="skeleton" style="height:200px"></div>
        <div style="padding:16px"><div class="skeleton" style="height:18px;width:70%;border-radius:6px;margin-bottom:8px"></div><div class="skeleton" style="height:14px;width:40%;border-radius:6px"></div></div>
      </div>
    </div>

    <div v-else-if="!ebooks.length" class="eb-empty">
      <div style="margin-bottom:16px;"><i class="ri-book-line" style="font-size:3rem;color:#94a3b8"></i></div>
      <h3>Belum ada ebook tersedia</h3>
      <p>Ebook akan segera hadir, nantikan ya!</p>
    </div>

    <div v-else class="eb-grid">
      <router-link v-for="eb in ebooks" :key="eb.id" :to="`/ebooks/${eb.id}`" class="eb-card">
        <div class="eb-cover" :style="coverStyle(eb)">
          <div class="eb-price-tag">{{ eb.price === 0 ? 'GRATIS' : formatRp(eb.price) }}</div>
          <div v-if="eb.already_purchased" class="eb-purchased-badge">✅ Sudah Dibeli</div>
        </div>
        <div class="eb-body">
          <h3 class="eb-title">{{ eb.title }}</h3>
          <p v-if="eb.author" class="eb-author">oleh {{ eb.author }}</p>
          <div v-if="eb.categories && eb.categories.length" class="eb-categories">
            <span v-for="cat in eb.categories" :key="cat.id" class="eb-category">{{ cat.name }}</span>
          </div>
          <p class="eb-desc">{{ stripHtml(eb.description, 80) }}</p>
          <div class="eb-footer">
            <span class="eb-sold"><i class="ri-shopping-cart-line"></i> {{ eb.total_sold }} terjual</span>
          </div>
        </div>
      </router-link>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { ebookApi } from '../services/api'
import { useToastStore } from '../stores/toast'

const toast = useToastStore()
const ebooks = ref([])
const loading = ref(false)
const search = ref('')

let timer = null
function debouncedLoad() { clearTimeout(timer); timer = setTimeout(load, 300) }

onMounted(load)

async function load() {
  loading.value = true
  try {
    const { data } = await ebookApi.list({ search: search.value })
    ebooks.value = data.data || []
  } catch { toast.error('Gagal memuat ebook') }
  finally { loading.value = false }
}

function formatRp(v) { return 'Rp ' + Number(v || 0).toLocaleString('id-ID') }
function stripHtml(html, n) {
  const text = (html || '').replace(/<[^>]*>/g, '').replace(/&nbsp;/g, ' ').trim()
  return text.length > n ? text.slice(0, n) + '...' : (text || 'Belum ada deskripsi')
}

function coverStyle(eb) {
  if (eb.cover_url) return { backgroundImage: `url(${eb.cover_url})`, backgroundSize: 'cover', backgroundPosition: 'center' }
  const g = ['linear-gradient(135deg,#667eea,#764ba2)','linear-gradient(135deg,#f093fb,#f5576c)','linear-gradient(135deg,#4facfe,#00f2fe)','linear-gradient(135deg,#43e97b,#38f9d7)','linear-gradient(135deg,#fa709a,#fee140)']
  let h = 0; for (let i = 0; i < (eb.title||'').length; i++) h = ((h << 5) - h) + eb.title.charCodeAt(i)
  return { background: g[Math.abs(h) % g.length] }
}
</script>

<style scoped>
.eb-hero { background: linear-gradient(135deg, #1e1b4b 0%, #312e81 50%, #4338ca 100%); border-radius: 16px; padding: 32px 36px; margin-bottom: 24px; }
.eb-hero-title { font-size: 1.6rem; font-weight: 800; color: #fff; }
.eb-hero-sub { font-size: 0.85rem; color: rgba(255,255,255,0.55); margin-top: 6px; }

.eb-search-bar { display: flex; align-items: center; gap: 10px; background: white; border: 1.5px solid #e2e8f0; border-radius: 12px; padding: 0 16px; height: 46px; margin-bottom: 24px; box-shadow: 0 1px 3px rgba(0,0,0,0.04); }
.eb-search-bar:focus-within { border-color: #818cf8; box-shadow: 0 0 0 3px rgba(129,140,248,0.12); }
.eb-search-bar input { flex: 1; border: none; outline: none; font-size: 0.85rem; color: #1e293b; background: none; font-family: inherit; }

.eb-grid { display: grid; grid-template-columns: repeat(auto-fill, minmax(260px, 1fr)); gap: 24px; }
.eb-card { background: #fff; border-radius: 16px; overflow: hidden; border: 1px solid #e8ecf1; text-decoration: none; display: flex; flex-direction: column; transition: transform 0.2s, box-shadow 0.2s; }
.eb-card:hover { transform: translateY(-4px); box-shadow: 0 12px 32px rgba(99,102,241,0.12); }

.eb-cover { height: 200px; position: relative; display: flex; align-items: center; justify-content: center; }
.eb-price-tag { position: absolute; bottom: 12px; left: 12px; background: rgba(0,0,0,0.75); color: #fbbf24; font-size: 0.78rem; font-weight: 800; padding: 5px 14px; border-radius: 8px; backdrop-filter: blur(8px); }
.eb-purchased-badge { position: absolute; top: 12px; right: 12px; background: rgba(22,163,74,0.9); color: white; font-size: 0.7rem; font-weight: 700; padding: 5px 12px; border-radius: 8px; }

.eb-body { padding: 18px 20px 16px; flex: 1; display: flex; flex-direction: column; }
.eb-title { font-size: 1rem; font-weight: 700; color: #0f172a; margin: 0 0 4px; }
.eb-author { font-size: 0.78rem; color: #6366f1; margin: 0 0 4px; font-weight: 600; }
.eb-categories { display: flex; flex-wrap: wrap; gap: 4px; margin-bottom: 8px; }
.eb-category { display: inline-block; padding: 3px 10px; font-size: 0.65rem; font-weight: 700; color: #6366f1; background: #eef2ff; border-radius: 6px; }
.eb-desc { font-size: 0.8rem; color: #64748b; line-height: 1.5; margin: 0 0 12px; flex: 1; }
.eb-footer { font-size: 0.75rem; color: #94a3b8; }

.eb-empty { text-align: center; padding: 80px 32px; background: white; border-radius: 16px; border: 2px dashed #e2e8f0; }
.eb-empty h3 { font-size: 1.1rem; font-weight: 700; color: #0f172a; margin: 0; }
.eb-empty p { font-size: 0.85rem; color: #94a3b8; margin-top: 6px; }

.skeleton-card { overflow: hidden; }
.skeleton { background: linear-gradient(90deg, #f1f5f9 25%, #e2e8f0 50%, #f1f5f9 75%); background-size: 200% 100%; animation: shimmer 1.5s infinite; }
@keyframes shimmer { 0% { background-position: 200% 0; } 100% { background-position: -200% 0; } }

@media (max-width: 768px) { .eb-grid { grid-template-columns: 1fr; } .eb-hero { padding: 24px 20px; } }
</style>
