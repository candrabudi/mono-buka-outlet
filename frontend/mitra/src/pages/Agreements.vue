<template>
  <div class="agr-page">
    <div class="agr-hero">
      <div><h1 class="agr-hero-title">Agreement</h1><p class="agr-hero-sub">Dokumen perjanjian kemitraan Anda</p></div>
    </div>

    <div class="agr-table-wrap">
      <table class="agr-table" v-if="agreements.length">
        <thead><tr><th>Judul</th><th>Tipe</th><th>Versi</th><th>Status</th><th>Tanggal</th><th>Aksi</th></tr></thead>
        <tbody>
          <tr v-for="a in agreements" :key="a.id">
            <td><span class="agr-title">{{ a.title || 'Agreement' }}</span></td>
            <td><span class="agr-type-badge" :class="a.type==='CONTRACT'?'at-contract':'at-doc'">{{ a.type==='CONTRACT'?'Kontrak':'Dokumen' }}</span></td>
            <td>v{{ a.version }}</td>
            <td>
              <span class="agr-badge" :class="a.status==='SIGNED'?'as-signed':'as-pending'">
                <span class="agr-badge-dot"></span>
                {{ a.status==='SIGNED'?'Ditandatangani':'Belum' }}
              </span>
            </td>
            <td><span class="agr-date">{{ formatDate(a.created_at) }}</span></td>
            <td><a v-if="a.file_url" :href="a.file_url" target="_blank" class="agr-view-btn">
              <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M18 13v6a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2V8a2 2 0 0 1 2-2h6"/><polyline points="15 3 21 3 21 9"/><line x1="10" y1="14" x2="21" y2="3"/></svg>
              Lihat
            </a></td>
          </tr>
        </tbody>
      </table>
      <div v-else class="agr-empty">
        <div class="agr-empty-circle"><svg width="32" height="32" viewBox="0 0 24 24" fill="none" stroke="#94a3b8" stroke-width="1.5"><path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z"/><polyline points="14 2 14 8 20 8"/></svg></div>
        <p>Belum ada agreement</p>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { agreementApi } from '../services/api'
import { useToastStore } from '../stores/toast'

const toast = useToastStore()
const agreements = ref([])

onMounted(async () => {
  try { const { data } = await agreementApi.list(); agreements.value = data.data || [] }
  catch { toast.error('Gagal memuat agreement') }
})

function formatDate(d) { return d ? new Date(d).toLocaleDateString('id-ID',{day:'numeric',month:'short',year:'numeric'}) : '-' }
</script>

<style scoped>
.agr-hero { background: linear-gradient(135deg, #0f0c29 0%, #302b63 50%, #24243e 100%); border-radius: 16px; padding: 32px 36px 24px; margin-bottom: 24px; box-shadow: 0 4px 24px rgba(15,12,41,0.2); }
.agr-hero-title { font-size: 1.6rem; font-weight: 800; color: #fff; margin: 0 0 4px; }
.agr-hero-sub { font-size: .85rem; color: rgba(255,255,255,.5); margin: 0; }

.agr-table-wrap { background: #fff; border-radius: 16px; border: 1px solid #e8ecf1; overflow: hidden; }
.agr-table { width: 100%; border-collapse: collapse; }
.agr-table thead { background: #f8fafc; }
.agr-table th { padding: 14px 20px; font-size: .72rem; font-weight: 700; color: #64748b; text-transform: uppercase; letter-spacing: .05em; text-align: left; border-bottom: 1px solid #e8ecf1; }
.agr-table td { padding: 16px 20px; font-size: .85rem; color: #1e293b; border-bottom: 1px solid #f1f5f9; vertical-align: middle; }
.agr-table tr:last-child td { border-bottom: none; }
.agr-table tbody tr { transition: background .15s; }
.agr-table tbody tr:hover { background: #fafbfc; }

.agr-title { font-weight: 700; color: #0f172a; }

.agr-type-badge { font-size: .68rem; font-weight: 700; padding: 4px 10px; border-radius: 6px; text-transform: uppercase; letter-spacing: .03em; }
.at-contract { background: #ede9fe; color: #7c3aed; }
.at-doc { background: #f1f5f9; color: #475569; }

.agr-badge { font-size: .68rem; font-weight: 700; padding: 4px 10px; border-radius: 6px; text-transform: uppercase; display: inline-flex; align-items: center; gap: 5px; }
.agr-badge-dot { width: 6px; height: 6px; border-radius: 50%; }
.as-signed { background: #dcfce7; color: #16a34a; }
.as-signed .agr-badge-dot { background: #22c55e; }
.as-pending { background: #fef3c7; color: #d97706; }
.as-pending .agr-badge-dot { background: #f59e0b; }

.agr-date { font-size: .82rem; color: #64748b; }

.agr-view-btn { display: inline-flex; align-items: center; gap: 5px; font-size: .78rem; font-weight: 600; color: #6366f1; text-decoration: none; padding: 6px 14px; border-radius: 8px; border: 1px solid #e0e7ff; background: #eef2ff; transition: all .15s; }
.agr-view-btn:hover { background: #e0e7ff; }

.agr-empty { text-align: center; padding: 56px 20px; }
.agr-empty-circle { width: 72px; height: 72px; border-radius: 50%; background: linear-gradient(135deg, #f1f5f9, #e2e8f0); display: flex; align-items: center; justify-content: center; margin: 0 auto 16px; }
.agr-empty p { color: #94a3b8; font-size: .85rem; margin: 0; }

@media (max-width: 768px) { .agr-hero { padding: 24px 20px 18px; } .agr-table-wrap { overflow-x: auto; } }
</style>
