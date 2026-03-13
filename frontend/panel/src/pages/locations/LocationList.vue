<template>
  <div class="lc-page">
    <!-- Hero -->
    <div class="lc-hero">
      <div class="lc-hero-top">
        <div>
          <h1 class="lc-hero-title">Pengajuan Lokasi</h1>
          <p class="lc-hero-sub">Daftar semua pengajuan lokasi dari seluruh partnership</p>
        </div>
      </div>
      <div class="lc-stats">
        <div class="lc-stat"><span class="lc-stat-dot" style="background:#818cf8"></span><span class="lc-stat-label">Total</span><span class="lc-stat-val">{{ stats[0].count }}</span></div>
        <div class="lc-stat"><span class="lc-stat-dot" style="background:#22c55e"></span><span class="lc-stat-label">Disetujui</span><span class="lc-stat-val">{{ stats[1].count }}</span></div>
        <div class="lc-stat"><span class="lc-stat-dot" style="background:#f59e0b"></span><span class="lc-stat-label">Menunggu</span><span class="lc-stat-val">{{ stats[2].count }}</span></div>
        <div class="lc-stat"><span class="lc-stat-dot" style="background:#ef4444"></span><span class="lc-stat-label">Ditolak</span><span class="lc-stat-val">{{ stats[3].count }}</span></div>
      </div>
    </div>

    <!-- Toolbar -->
    <div class="lc-toolbar">
      <div class="lc-search-wrap">
        <i class="ri-search-line lc-search-ico"></i>
        <input v-model="search" class="lc-search" placeholder="Cari lokasi atau alamat..." @input="onFilter" />
      </div>
      <div class="lc-filters">
        <SearchSelect
          v-model="filterStatus"
          :options="statusSelectOptions"
          placeholder="Semua Status"
          empty-label="Semua Status"
          @update:model-value="onFilter"
        />
        <input v-model="filterKota" class="lc-select lc-kota-input" placeholder="Filter kota..." @input="onFilter" />
      </div>
    </div>

    <!-- Loading -->
    <div v-if="loading" class="lc-list">
      <div v-for="n in 4" :key="n" class="lc-card lc-card-skel">
        <div class="skel-bar shimmer" style="width:60%;height:18px"></div>
        <div class="skel-bar shimmer" style="width:40%;height:14px"></div>
      </div>
    </div>

    <!-- Location List -->
    <div v-else-if="list.length" class="lc-list">
      <div v-for="loc in list" :key="loc.id" class="lc-card" @click="$router.push({name:'LocationDetail',params:{id:loc.id}})">
        <div class="lc-card-head">
          <div class="lc-card-left">
            <span class="lc-badge" :class="'ls-'+loc.status">{{ statusLabel(loc.status) }}</span>
            <span v-if="loc.tipe_bangunan" class="lc-type-tag"><i class="ri-building-2-line"></i> {{ loc.tipe_bangunan }}</span>
          </div>
          <div class="lc-score-pill" :style="{background:scoreColor(loc.total_score)+'18',color:scoreColor(loc.total_score)}">
            <i class="ri-bar-chart-fill"></i> {{ loc.total_score }}<span class="lc-score-cat">{{ loc.score_category }}</span>
          </div>
        </div>
        <h3 class="lc-card-title">{{ loc.nama_lokasi }}</h3>
        <div class="lc-card-meta">
          <span class="lc-meta-item"><i class="ri-map-pin-line"></i> {{ loc.kota || loc.alamat?.substring(0,40) || '-' }}</span>
          <span class="lc-meta-item"><i class="ri-ruler-line"></i> {{ loc.luas_tempat }} m²</span>
          <span class="lc-meta-item"><i class="ri-money-dollar-circle-line"></i> {{ fc(loc.harga_sewa_per_tahun) }}/thn</span>
          <span class="lc-meta-item"><i class="ri-user-line"></i> {{ loc.mitra?.name || '-' }}</span>
        </div>
        <div class="lc-card-foot">
          <span class="lc-card-date"><i class="ri-calendar-line"></i> {{ formatDate(loc.created_at) }}</span>
          <div class="lc-score-bar-wrap">
            <div class="lc-score-bar"><div class="lc-score-fill" :style="{width:loc.total_score+'%',background:scoreColor(loc.total_score)}"></div></div>
          </div>
        </div>
      </div>
    </div>

    <!-- Empty -->
    <div v-else class="lc-empty">
      <div class="lc-empty-ico"><i class="ri-map-pin-add-line"></i></div>
      <h3>Belum ada pengajuan lokasi</h3>
      <p>Tambahkan lokasi dari halaman detail partnership.</p>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { locationApi } from '../../services/api'
import { useToastStore } from '../../stores/toast'
import SearchSelect from '../../components/SearchSelect.vue'

const toast = useToastStore()
const list = ref([])
const loading = ref(true)
const search = ref('')
const filterStatus = ref('')
const filterKota = ref('')

const statuses = [
  { val:'DRAFT', label:'Draft' },
  { val:'SUBMITTED', label:'Diajukan' },
  { val:'IN_REVIEW', label:'Ditinjau' },
  { val:'SURVEY_SCHEDULED', label:'Survei Dijadwalkan' },
  { val:'SURVEYED', label:'Sudah Disurvei' },
  { val:'APPROVED', label:'Disetujui' },
  { val:'REJECTED', label:'Ditolak' },
  { val:'REVISION_NEEDED', label:'Revisi' },
]
const statusSelectOptions = statuses.map(s => ({ value: s.val, label: s.label }))

const stats = computed(() => {
  const all = list.value
  return [
    { label:'Total', count: all.length },
    { label:'Disetujui', count: all.filter(l => l.status==='APPROVED').length },
    { label:'Menunggu', count: all.filter(l => ['SUBMITTED','IN_REVIEW','SURVEY_SCHEDULED'].includes(l.status)).length },
    { label:'Ditolak', count: all.filter(l => l.status==='REJECTED').length },
  ]
})

function statusLabel(s) {
  return { DRAFT:'Draft', SUBMITTED:'Diajukan', IN_REVIEW:'Ditinjau', SURVEY_SCHEDULED:'Survei', SURVEYED:'Disurvei', APPROVED:'Disetujui', REJECTED:'Ditolak', REVISION_NEEDED:'Revisi' }[s] || s
}

function scoreColor(s) {
  if (s >= 80) return '#22c55e'
  if (s >= 65) return '#6366f1'
  if (s >= 50) return '#f59e0b'
  return '#ef4444'
}

function fc(n) { return 'Rp' + (n||0).toLocaleString('id-ID') }
function formatDate(d) { return d ? new Date(d).toLocaleDateString('id-ID', { day:'2-digit', month:'short', year:'numeric' }) : '-' }

let debounce = null
function onFilter() {
  clearTimeout(debounce)
  debounce = setTimeout(loadData, 300)
}

async function loadData() {
  loading.value = true
  try {
    const { data } = await locationApi.getAll({ status: filterStatus.value, kota: filterKota.value, search: search.value })
    list.value = data.data || []
  } catch(e) { toast.error('Gagal memuat data lokasi') }
  finally { loading.value = false }
}

onMounted(loadData)
</script>

<style scoped>
/* ═══ HERO ═══ */
.lc-hero { background: linear-gradient(135deg, #0f0c29 0%, #302b63 50%, #24243e 100%); border-radius: 16px; padding: 32px 36px 24px; margin-bottom: 24px; }
.lc-hero-top { display:flex; justify-content:space-between; align-items:flex-start; gap:16px; flex-wrap:wrap; margin-bottom:20px; }
.lc-hero-title { font-size:1.6rem; font-weight:800; color:#fff; margin:0 0 4px; }
.lc-hero-sub { font-size:.85rem; color:rgba(255,255,255,.5); margin:0; }
.lc-stats { display:flex; gap:28px; flex-wrap:wrap; padding-top:16px; border-top:1px solid rgba(255,255,255,.08); }
.lc-stat { display:flex; align-items:center; gap:8px; }
.lc-stat-dot { width:8px; height:8px; border-radius:50%; }
.lc-stat-label { font-size:.72rem; color:rgba(255,255,255,.4); text-transform:uppercase; letter-spacing:.05em; }
.lc-stat-val { font-size:.9rem; font-weight:800; color:#fff; }

/* ═══ TOOLBAR ═══ */
.lc-toolbar { display:flex; justify-content:space-between; align-items:center; gap:12px; flex-wrap:wrap; margin-bottom:20px; }
.lc-search-wrap { position:relative; min-width:260px; flex:1; max-width:380px; }
.lc-search-ico { position:absolute; left:14px; top:50%; transform:translateY(-50%); color:#94a3b8; font-size:1rem; }
.lc-search { width:100%; padding:10px 14px 10px 40px; border:1.5px solid #e2e8f0; border-radius:12px; font-size:.85rem; background:#fff; color:#1e293b; outline:none; box-sizing:border-box; font-family:inherit; }
.lc-search:focus { border-color:#6366f1; }
.lc-filters { display:flex; gap:10px; }
.lc-select { padding:10px 14px; border:1.5px solid #e2e8f0; border-radius:12px; font-size:.82rem; background:#fff; color:#1e293b; outline:none; cursor:pointer; font-family:inherit; }
.lc-kota-input { min-width:140px; }

/* ═══ CARD LIST ═══ */
.lc-list { display:flex; flex-direction:column; gap:12px; }
.lc-card { background:#fff; border-radius:14px; border:1px solid #e8ecf1; padding:20px 24px; cursor:pointer; transition:all .2s; }
.lc-card:hover { border-color:#c7d2fe; box-shadow:0 2px 12px rgba(99,102,241,.06); }
.lc-card-head { display:flex; justify-content:space-between; align-items:center; margin-bottom:10px; }
.lc-card-left { display:flex; gap:8px; align-items:center; }
.lc-card-title { font-size:1.05rem; font-weight:700; color:#0f172a; margin:0 0 10px; }
.lc-card-meta { display:flex; flex-wrap:wrap; gap:16px; }
.lc-meta-item { display:flex; align-items:center; gap:5px; font-size:.78rem; color:#64748b; }
.lc-meta-item i { color:#94a3b8; font-size:.9rem; }
.lc-card-foot { display:flex; justify-content:space-between; align-items:center; margin-top:14px; padding-top:12px; border-top:1px solid #f5f7fa; }
.lc-card-date { font-size:.72rem; color:#94a3b8; display:flex; align-items:center; gap:4px; }
.lc-card-date i { font-size:.8rem; }
.lc-score-bar-wrap { width:120px; }
.lc-score-bar { height:5px; background:#f1f5f9; border-radius:3px; overflow:hidden; }
.lc-score-fill { height:100%; border-radius:3px; transition:width .3s; }

.lc-score-pill { display:inline-flex; align-items:center; gap:5px; padding:5px 12px; border-radius:8px; font-size:.72rem; font-weight:700; }
.lc-score-pill i { font-size:.8rem; }
.lc-score-cat { font-weight:500; margin-left:4px; opacity:.75; }
.lc-type-tag { display:inline-flex; align-items:center; gap:4px; padding:4px 10px; border-radius:6px; font-size:.68rem; font-weight:600; background:#f1f5f9; color:#475569; text-transform:capitalize; }
.lc-type-tag i { font-size:.75rem; }

.lc-badge { font-size:.68rem; font-weight:700; text-transform:uppercase; letter-spacing:.04em; padding:4px 10px; border-radius:6px; }
.ls-DRAFT { background:#f1f5f9; color:#475569; }
.ls-SUBMITTED { background:#e0f2fe; color:#0284c7; }
.ls-IN_REVIEW { background:#fef3c7; color:#92400e; }
.ls-SURVEY_SCHEDULED { background:#e0e7ff; color:#4338ca; }
.ls-SURVEYED { background:#ede9fe; color:#6d28d9; }
.ls-APPROVED { background:#dcfce7; color:#15803d; }
.ls-REJECTED { background:#fee2e2; color:#b91c1c; }
.ls-REVISION_NEEDED { background:#ffedd5; color:#c2410c; }

/* ═══ SKELETON ═══ */
.lc-card-skel { display:flex; flex-direction:column; gap:12px; min-height:80px; }
.skel-bar { border-radius:6px; }
.shimmer { background:linear-gradient(90deg,#e8ecf1 25%,#f1f5f9 50%,#e8ecf1 75%); background-size:200% 100%; animation:shimmer 1.5s infinite; }
@keyframes shimmer { 0%{background-position:200% 0} 100%{background-position:-200% 0} }

/* ═══ EMPTY ═══ */
.lc-empty { display:flex; flex-direction:column; align-items:center; padding:80px 32px; background:#fff; border-radius:16px; border:2px dashed #e2e8f0; text-align:center; }
.lc-empty-ico { width:80px; height:80px; border-radius:50%; background:linear-gradient(135deg,#f1f5f9,#e2e8f0); display:flex; align-items:center; justify-content:center; margin-bottom:20px; font-size:2rem; color:#94a3b8; }
.lc-empty h3 { font-weight:700; color:#0f172a; margin:0; }
.lc-empty p { color:#94a3b8; font-size:.85rem; margin:6px 0 0; }

@media (max-width:768px) {
  .lc-hero { padding:24px 20px 18px; }
  .lc-hero-title { font-size:1.3rem; }
  .lc-toolbar { flex-direction:column; }
  .lc-search-wrap { max-width:100%; }
  .lc-filters { width:100%; }
  .lc-card-meta { gap:10px; }
}
</style>
