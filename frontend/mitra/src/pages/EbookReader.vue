<template>
  <div class="reader-wrapper" ref="wrapperRef" @keydown="handleKey" tabindex="0">
    <!-- Top toolbar -->
    <div class="reader-toolbar">
      <router-link :to="`/ebooks/${$route.params.id}`" class="reader-back">
        <i class="ri-arrow-left-s-line"></i> Kembali
      </router-link>
      <span class="reader-title">{{ title }}</span>
      <div class="reader-controls">
        <button @click="zoomOut" class="ctrl-btn" title="Zoom Out"><i class="ri-subtract-line"></i></button>
        <span class="zoom-label">{{ Math.round(scale * 100) }}%</span>
        <button @click="zoomIn" class="ctrl-btn" title="Zoom In"><i class="ri-add-line"></i></button>
        <div class="divider"></div>
        <button @click="prevPage" class="ctrl-btn" :disabled="currentPage <= 1" title="Previous"><i class="ri-arrow-up-s-line"></i></button>
        <span class="page-label">{{ currentPage }} / {{ totalPages }}</span>
        <button @click="nextPage" class="ctrl-btn" :disabled="currentPage >= totalPages" title="Next"><i class="ri-arrow-down-s-line"></i></button>
        <div class="divider"></div>
        <button @click="toggleFullscreen" class="ctrl-btn" title="Fullscreen"><i class="ri-fullscreen-line"></i></button>
      </div>
    </div>

    <!-- PDF viewport -->
    <div class="reader-viewport" ref="viewportRef" @scroll="onScroll">
      <div v-if="loading" class="reader-loading">
        <div class="spinner"></div>
        <p>Memuat ebook...</p>
      </div>
      <div v-else-if="error" class="reader-error">
        <i class="ri-error-warning-line" style="font-size:3rem;color:#ef4444"></i>
        <p>{{ error }}</p>
        <router-link :to="`/ebooks/${$route.params.id}`" class="btn-back">Kembali ke Detail</router-link>
      </div>
      <div v-else class="pages-container" ref="pagesRef">
        <canvas
          v-for="p in totalPages"
          :key="p"
          :ref="el => { if (el) canvasRefs[p] = el }"
          class="pdf-page"
          :data-page="p"
        ></canvas>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, onUnmounted, watch, nextTick } from 'vue'
import { useRoute } from 'vue-router'
import { ebookApi } from '../services/api'
import { useToastStore } from '../stores/toast'

const route = useRoute()
const toast = useToastStore()

const wrapperRef = ref(null)
const viewportRef = ref(null)
const pagesRef = ref(null)
const canvasRefs = ref({})

const title = ref('Baca Ebook')
const loading = ref(true)
const error = ref(null)
const scale = ref(1.2)
const currentPage = ref(1)
const totalPages = ref(0)

let pdfDoc = null
let rendering = false
let renderQueue = []

const token = localStorage.getItem('mitra_token')
const baseUrl = import.meta.env.VITE_API_BASE_URL || '/api/v1/mitra'
const pdfUrl = `${baseUrl}/ebooks/${route.params.id}/read?token=${token}`

onMounted(async () => {
  try {
    const { data } = await ebookApi.get(route.params.id)
    title.value = data.data?.title || 'Baca Ebook'
    if (!data.already_purchased) {
      error.value = 'Anda belum membeli ebook ini'
      loading.value = false
      return
    }
    await loadPdf()
  } catch {
    error.value = 'Gagal memuat ebook'
    loading.value = false
  }
  wrapperRef.value?.focus()
})

onUnmounted(() => {
  if (pdfDoc) pdfDoc.destroy()
})

async function loadPdf() {
  try {
    const pdfjsLib = await import('https://cdnjs.cloudflare.com/ajax/libs/pdf.js/4.0.379/pdf.min.mjs')
    pdfjsLib.GlobalWorkerOptions.workerSrc = 'https://cdnjs.cloudflare.com/ajax/libs/pdf.js/4.0.379/pdf.worker.min.mjs'

    const loadingTask = pdfjsLib.getDocument(pdfUrl)
    pdfDoc = await loadingTask.promise
    totalPages.value = pdfDoc.numPages
    loading.value = false

    await nextTick()
    renderAllPages()
  } catch (e) {
    console.error('PDF load error:', e)
    error.value = 'Gagal memuat file PDF'
    loading.value = false
  }
}

async function renderPage(pageNum) {
  const canvas = canvasRefs.value[pageNum]
  if (!canvas || !pdfDoc) return

  const page = await pdfDoc.getPage(pageNum)
  const viewport = page.getViewport({ scale: scale.value })

  canvas.width = viewport.width
  canvas.height = viewport.height

  const ctx = canvas.getContext('2d')
  await page.render({ canvasContext: ctx, viewport }).promise
}

async function renderAllPages() {
  if (rendering) return
  rendering = true
  for (let i = 1; i <= totalPages.value; i++) {
    await renderPage(i)
  }
  rendering = false
}

function onScroll() {
  if (!viewportRef.value || !totalPages.value) return
  const container = viewportRef.value
  const scrollTop = container.scrollTop
  const containerHeight = container.clientHeight

  for (let i = 1; i <= totalPages.value; i++) {
    const canvas = canvasRefs.value[i]
    if (!canvas) continue
    const rect = canvas.getBoundingClientRect()
    const containerRect = container.getBoundingClientRect()
    const relativeTop = rect.top - containerRect.top
    if (relativeTop > -rect.height / 2 && relativeTop < containerHeight / 2) {
      currentPage.value = i
      break
    }
  }
}

function scrollToPage(num) {
  const canvas = canvasRefs.value[num]
  if (canvas) {
    canvas.scrollIntoView({ behavior: 'smooth', block: 'start' })
    currentPage.value = num
  }
}

function prevPage() { if (currentPage.value > 1) scrollToPage(currentPage.value - 1) }
function nextPage() { if (currentPage.value < totalPages.value) scrollToPage(currentPage.value + 1) }

function zoomIn() {
  scale.value = Math.min(scale.value + 0.2, 3)
}
function zoomOut() {
  scale.value = Math.max(scale.value - 0.2, 0.5)
}

watch(scale, async () => {
  await nextTick()
  renderAllPages()
})

function toggleFullscreen() {
  if (document.fullscreenElement) {
    document.exitFullscreen()
  } else {
    wrapperRef.value?.requestFullscreen()
  }
}

function handleKey(e) {
  if (e.key === 'ArrowLeft' || e.key === 'ArrowUp') { prevPage(); e.preventDefault() }
  else if (e.key === 'ArrowRight' || e.key === 'ArrowDown') { nextPage(); e.preventDefault() }
  else if (e.key === '+' || e.key === '=') { zoomIn(); e.preventDefault() }
  else if (e.key === '-') { zoomOut(); e.preventDefault() }
}
</script>

<style scoped>
.reader-wrapper {
  display: flex; flex-direction: column; height: calc(100vh - 90px);
  background: #111827; border-radius: 16px; overflow: hidden; outline: none;
}

/* Toolbar */
.reader-toolbar {
  display: flex; align-items: center; gap: 12px;
  padding: 10px 20px; background: #1f2937;
  border-bottom: 1px solid rgba(255,255,255,0.06); flex-shrink: 0;
}
.reader-back {
  display: flex; align-items: center; gap: 4px;
  color: rgba(255,255,255,0.6); text-decoration: none;
  font-size: 0.8rem; font-weight: 600; white-space: nowrap;
  transition: color 0.2s;
}
.reader-back:hover { color: white; }
.reader-back i { font-size: 1.1rem; }
.reader-title {
  flex: 1; font-size: 0.85rem; font-weight: 700; color: white;
  white-space: nowrap; overflow: hidden; text-overflow: ellipsis;
}
.reader-controls { display: flex; align-items: center; gap: 4px; flex-shrink: 0; }
.ctrl-btn {
  width: 32px; height: 32px; border-radius: 8px; border: none;
  background: rgba(255,255,255,0.06); color: rgba(255,255,255,0.7);
  font-size: 0.9rem; cursor: pointer;
  display: flex; align-items: center; justify-content: center;
  transition: all 0.15s;
}
.ctrl-btn:hover { background: rgba(255,255,255,0.15); color: white; }
.ctrl-btn:disabled { opacity: 0.3; cursor: not-allowed; }
.zoom-label, .page-label {
  font-size: 0.72rem; color: rgba(255,255,255,0.5);
  min-width: 50px; text-align: center; font-weight: 600;
  font-variant-numeric: tabular-nums;
}
.divider { width: 1px; height: 20px; background: rgba(255,255,255,0.1); margin: 0 4px; }

/* Viewport */
.reader-viewport {
  flex: 1; overflow: auto; background: #0f172a;
  display: flex; justify-content: center;
}
.reader-viewport::-webkit-scrollbar { width: 8px; }
.reader-viewport::-webkit-scrollbar-track { background: #1e293b; }
.reader-viewport::-webkit-scrollbar-thumb { background: #475569; border-radius: 4px; }
.reader-viewport::-webkit-scrollbar-thumb:hover { background: #64748b; }

.pages-container {
  display: flex; flex-direction: column; align-items: center;
  padding: 24px 0; gap: 16px;
}

.pdf-page {
  box-shadow: 0 4px 24px rgba(0,0,0,0.5);
  border-radius: 4px;
  max-width: 100%;
}

/* Loading & Error */
.reader-loading, .reader-error {
  display: flex; flex-direction: column; align-items: center;
  justify-content: center; height: 100%; width: 100%;
  color: rgba(255,255,255,0.5); gap: 16px; padding: 32px;
}
.spinner {
  width: 40px; height: 40px;
  border: 3px solid rgba(255,255,255,0.1);
  border-top-color: #6366f1; border-radius: 50%;
  animation: spin 1s linear infinite;
}
@keyframes spin { to { transform: rotate(360deg); } }
.btn-back {
  background: linear-gradient(135deg, #6366f1, #8b5cf6);
  color: white; padding: 10px 24px; border-radius: 10px;
  text-decoration: none; font-size: 0.85rem; font-weight: 600;
}

@media (max-width: 768px) {
  .reader-wrapper { height: calc(100vh - 70px); border-radius: 0; }
  .reader-toolbar { padding: 8px 12px; gap: 8px; }
  .reader-title { display: none; }
  .zoom-label { display: none; }
  .divider { display: none; }
}
</style>
