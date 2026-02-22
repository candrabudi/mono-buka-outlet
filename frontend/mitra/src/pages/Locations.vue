<template>
  <div class="loc-page">
    <div class="loc-hero">
      <div><h1 class="loc-hero-title">Pengajuan Lokasi</h1><p class="loc-hero-sub">Daftar lokasi yang Anda ajukan</p></div>
    </div>

    <div class="loc-grid" v-if="locations.length">
      <div v-for="loc in locations" :key="loc.id" class="loc-card">
        <div class="loc-card-head">
          <span class="loc-badge" :class="'ls-'+loc.status">{{ statusLabel(loc.status) }}</span>
          <span class="loc-score" v-if="loc.total_score">{{ loc.total_score }} pts</span>
        </div>
        <h4 class="loc-name">{{ loc.nama_lokasi }}</h4>
        <div class="loc-meta">
          <span><i class="ri-map-pin-line"></i> {{ loc.kota || '-' }}</span>
          <span><i class="ri-layout-line" style="font-size:12px"></i> {{ loc.luas_tempat || 0 }} m²</span>
        </div>
      </div>
    </div>
    <div v-else class="loc-empty">
      <div class="loc-empty-circle"><i class="ri-map-pin-line"></i></div>
      <p>Belum ada pengajuan lokasi</p>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { locationApi } from '../services/api'
import { useToastStore } from '../stores/toast'

const toast = useToastStore()
const locations = ref([])

onMounted(async () => {
  try { const { data } = await locationApi.list(); locations.value = data.data || [] }
  catch { toast.error('Gagal memuat lokasi') }
})

function statusLabel(s) { return { DRAFT:'Draft', SUBMITTED:'Diajukan', IN_REVIEW:'Ditinjau', SURVEY_SCHEDULED:'Survei', SURVEYED:'Disurvei', APPROVED:'Disetujui', REJECTED:'Ditolak', REVISION_NEEDED:'Revisi' }[s] || s }
</script>

<style scoped>
.loc-hero { background: linear-gradient(135deg, #0f0c29 0%, #302b63 50%, #24243e 100%); border-radius: 16px; padding: 32px 36px 24px; margin-bottom: 24px; box-shadow: 0 4px 24px rgba(15,12,41,0.2); }
.loc-hero-title { font-size: 1.6rem; font-weight: 800; color: #fff; margin: 0 0 4px; }
.loc-hero-sub { font-size: .85rem; color: rgba(255,255,255,.5); margin: 0; }

.loc-grid { display: grid; grid-template-columns: repeat(auto-fill, minmax(300px, 1fr)); gap: 16px; }
.loc-card { background: #fff; border-radius: 16px; border: 1px solid #e8ecf1; padding: 20px; transition: all .2s; }
.loc-card:hover { box-shadow: 0 4px 16px rgba(0,0,0,.06); transform: translateY(-1px); }
.loc-card-head { display: flex; align-items: center; justify-content: space-between; margin-bottom: 12px; }
.loc-badge { font-size: .68rem; font-weight: 700; padding: 4px 10px; border-radius: 6px; text-transform: uppercase; }
.ls-DRAFT { background: #f1f5f9; color: #475569; }
.ls-SUBMITTED { background: #e0f2fe; color: #0284c7; }
.ls-IN_REVIEW { background: #fef3c7; color: #d97706; }
.ls-APPROVED { background: #dcfce7; color: #16a34a; }
.ls-REJECTED { background: #fee2e2; color: #dc2626; }
.ls-SURVEY_SCHEDULED,.ls-SURVEYED { background: #ede9fe; color: #7c3aed; }
.ls-REVISION_NEEDED { background: #fce7f3; color: #db2777; }
.loc-score { font-size: .82rem; font-weight: 800; color: #6366f1; }
.loc-name { font-size: .95rem; font-weight: 700; color: #0f172a; margin: 0 0 10px; }
.loc-meta { display: flex; gap: 16px; font-size: .78rem; color: #64748b; }
.loc-meta span { display: flex; align-items: center; gap: 4px; }

.loc-empty { text-align: center; padding: 80px 32px; background: white; border-radius: 16px; border: 2px dashed #e2e8f0; }
.loc-empty-circle { width: 72px; height: 72px; border-radius: 50%; background: linear-gradient(135deg, #f1f5f9, #e2e8f0); display: flex; align-items: center; justify-content: center; margin: 0 auto 16px; }
.loc-empty p { color: #94a3b8; font-size: .85rem; margin: 0; }

@media (max-width: 768px) { .loc-hero { padding: 24px 20px 18px; } .loc-grid { grid-template-columns: 1fr; } }
</style>
