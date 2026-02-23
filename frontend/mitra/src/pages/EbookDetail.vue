<template>
  <div class="animate-in">
    <div v-if="loading" class="ed-loading">Memuat ebook...</div>
    <div v-else-if="!ebook" class="ed-loading">Ebook tidak ditemukan</div>

    <template v-else>
      <div class="ed-hero" :style="heroStyle">
        <div class="ed-hero-overlay"></div>
        <div class="ed-hero-content">
          <div class="ed-hero-price">{{ ebook.price === 0 ? 'GRATIS' : formatRp(ebook.price) }}</div>
          <h1 class="ed-hero-title">{{ ebook.title }}</h1>
          <p v-if="ebook.author" class="ed-hero-author">oleh {{ ebook.author }}</p>
          <div v-if="ebook.categories && ebook.categories.length" class="ed-hero-categories">
            <span v-for="cat in ebook.categories" :key="cat.id" class="ed-hero-category">{{ cat.name }}</span>
          </div>
          <div class="ed-hero-stats">
            <span><i class="ri-shopping-cart-line"></i> {{ ebook.total_sold }} terjual</span>
          </div>
        </div>
      </div>

      <div class="ed-content">
        <div class="ed-main">
          <div class="ed-section">
            <h2 class="ed-section-title">Deskripsi</h2>
            <div class="ed-desc" v-html="ebook.description || '<p>Belum ada deskripsi.</p>'"></div>
          </div>
        </div>

        <div class="ed-sidebar">
          <div class="ed-action-card">
            <!-- Not purchased -->
            <template v-if="!purchased">
              <h3 class="ed-action-title">{{ ebook.price === 0 ? 'Klaim Gratis' : 'Beli Ebook' }}</h3>
              <p class="ed-action-price">{{ ebook.price === 0 ? 'Gratis' : formatRp(ebook.price) }}</p>
              <button @click="doPurchase" class="ed-btn ed-btn-buy" :disabled="purchasing">
                {{ purchasing ? 'Memproses...' : (ebook.price === 0 ? 'Klaim Sekarang' : 'Beli Sekarang') }}
              </button>
            </template>

            <!-- Purchased -->
            <template v-else>
              <h3 class="ed-action-title"><i class="ri-checkbox-circle-line" style="color:#16a34a"></i> Ebook Sudah Dibeli</h3>
              <router-link :to="`/ebooks/${ebook.id}/read`" class="ed-btn ed-btn-read">
                <i class="ri-book-read-line"></i> Baca Online
              </router-link>
            </template>
          </div>
        </div>
      </div>
    </template>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ebookApi, midtransApi } from '../services/api'
import { useToastStore } from '../stores/toast'

const route = useRoute()
const router = useRouter()
const toast = useToastStore()

const ebook = ref(null)
const purchased = ref(false)
const loading = ref(true)
const purchasing = ref(false)

onMounted(async () => {
  try {
    const { data } = await ebookApi.get(route.params.id)
    ebook.value = data.data
    purchased.value = data.already_purchased
  } catch { toast.error('Gagal memuat ebook') }
  finally { loading.value = false }
})

const heroStyle = computed(() => {
  if (ebook.value?.cover_url) return { backgroundImage: `url(${ebook.value.cover_url})`, backgroundSize: 'cover', backgroundPosition: 'center' }
  return { background: 'linear-gradient(135deg, #1e1b4b, #4338ca)' }
})

async function doPurchase() {
  if (purchasing.value) return
  purchasing.value = true
  try {
    const { data } = await ebookApi.purchase(ebook.value.id)
    const order = data.data

    // Free ebook
    if (!order.midtrans_snap_token) {
      toast.success(data.message || 'Ebook berhasil diklaim!')
      purchased.value = true
      purchasing.value = false
      return
    }

    // Paid ebook — open Midtrans Snap
    const { data: keyData } = await midtransApi.getClientKey()
    const script = document.createElement('script')
    script.src = keyData.snap_url
    script.setAttribute('data-client-key', keyData.client_key)
    script.onload = () => {
      window.snap.pay(order.midtrans_snap_token, {
        onSuccess: () => { toast.success('Pembayaran berhasil!'); purchased.value = true; purchasing.value = false },
        onPending: () => { toast.info('Pembayaran pending, silakan selesaikan'); purchasing.value = false },
        onError: () => { toast.error('Pembayaran gagal'); purchasing.value = false },
        onClose: () => { toast.info('Popup pembayaran ditutup'); purchasing.value = false },
      })
    }
    document.head.appendChild(script)
  } catch (err) {
    toast.error(err.response?.data?.error || 'Gagal memproses pembelian')
    purchasing.value = false
  }
}

function formatRp(v) { return 'Rp ' + Number(v || 0).toLocaleString('id-ID') }
</script>

<style scoped>
.ed-loading { text-align: center; padding: 80px; color: #64748b; }

.ed-hero { position: relative; border-radius: 16px; overflow: hidden; min-height: 240px; display: flex; align-items: flex-end; margin-bottom: 24px; }
.ed-hero-overlay { position: absolute; inset: 0; background: linear-gradient(180deg, rgba(0,0,0,0.1) 0%, rgba(0,0,0,0.8) 100%); }
.ed-hero-content { position: relative; z-index: 2; padding: 32px; width: 100%; }
.ed-hero-price { display: inline-block; background: rgba(251,191,36,0.2); color: #fbbf24; font-size: 0.85rem; font-weight: 800; padding: 6px 16px; border-radius: 8px; margin-bottom: 12px; backdrop-filter: blur(8px); }
.ed-hero-title { font-size: 1.8rem; font-weight: 800; color: white; margin: 0 0 4px; }
.ed-hero-author { font-size: 0.9rem; color: rgba(255,255,255,0.7); margin: 0 0 8px; }
.ed-hero-categories { display: flex; flex-wrap: wrap; gap: 6px; margin-bottom: 8px; }
.ed-hero-category { display: inline-block; padding: 4px 14px; font-size: 0.72rem; font-weight: 700; color: white; background: rgba(99,102,241,0.4); border-radius: 8px; backdrop-filter: blur(8px); }
.ed-hero-stats { font-size: 0.8rem; color: rgba(255,255,255,0.5); display: flex; gap: 16px; }

.ed-content { display: grid; grid-template-columns: 1fr 340px; gap: 24px; align-items: start; }
.ed-main { background: white; border-radius: 16px; padding: 28px; border: 1px solid #e8ecf1; }
.ed-section-title { font-size: 1rem; font-weight: 700; color: #0f172a; margin: 0 0 12px; }
.ed-desc { font-size: 0.88rem; color: #475569; line-height: 1.7; margin: 0; }
.ed-desc :deep(h1), .ed-desc :deep(h2), .ed-desc :deep(h3) { color: #0f172a; margin: 1em 0 0.5em; font-weight: 700; }
.ed-desc :deep(h1) { font-size: 1.4rem; }
.ed-desc :deep(h2) { font-size: 1.2rem; }
.ed-desc :deep(h3) { font-size: 1.05rem; }
.ed-desc :deep(p) { margin: 0 0 0.8em; }
.ed-desc :deep(ul), .ed-desc :deep(ol) { padding-left: 1.5em; margin: 0 0 0.8em; }
.ed-desc :deep(li) { margin-bottom: 0.3em; }
.ed-desc :deep(blockquote) { border-left: 3px solid #6366f1; padding-left: 1em; margin: 1em 0; color: #64748b; font-style: italic; }
.ed-desc :deep(a) { color: #6366f1; text-decoration: underline; }
.ed-desc :deep(strong) { font-weight: 700; color: #0f172a; }

.ed-sidebar { position: sticky; top: 24px; }
.ed-action-card { background: white; border-radius: 16px; padding: 28px; border: 1px solid #e8ecf1; text-align: center; }
.ed-action-title { font-size: 1rem; font-weight: 700; color: #0f172a; margin: 0 0 8px; }
.ed-action-price { font-size: 1.5rem; font-weight: 800; color: #f59e0b; margin: 0 0 20px; }

.ed-btn { display: block; width: 100%; padding: 14px; border-radius: 12px; font-size: 0.9rem; font-weight: 700; border: none; cursor: pointer; text-decoration: none; text-align: center; margin-bottom: 10px; }
.ed-btn:last-child { margin-bottom: 0; }
.ed-btn-buy { background: linear-gradient(135deg, #f59e0b, #ef4444); color: white; box-shadow: 0 4px 20px rgba(245,158,11,0.3); }
.ed-btn-read { background: linear-gradient(135deg, #6366f1, #8b5cf6); color: white; box-shadow: 0 4px 20px rgba(99,102,241,0.3); }
.ed-btn:disabled { opacity: 0.5; cursor: not-allowed; }

@media (max-width: 768px) {
  .ed-content { grid-template-columns: 1fr; }
  .ed-hero-title { font-size: 1.4rem; }
}
</style>
