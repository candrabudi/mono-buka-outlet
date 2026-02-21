<template>
  <div class="ld-page" v-if="loc">
    <!-- Back -->
    <router-link to="/locations" class="ld-back"><i class="ri-arrow-left-line"></i> Kembali ke Daftar</router-link>

    <!-- Hero -->
    <div class="ld-hero">
      <div class="ld-hero-top">
        <div>
          <h1 class="ld-hero-title">{{ loc.nama_lokasi }}</h1>
          <p class="ld-hero-addr"><i class="ri-map-pin-line"></i> {{ loc.alamat }}</p>
          <div class="ld-hero-meta">
            <span><i class="ri-building-line"></i> {{ loc.kota }}, {{ loc.provinsi }}</span>
            <span><i class="ri-ruler-line"></i> {{ loc.luas_tempat }} m²</span>
            <span><i class="ri-home-4-line"></i> {{ loc.tipe_bangunan || '-' }}</span>
          </div>
        </div>
        <div class="ld-hero-score-ring">
          <svg viewBox="0 0 80 80" class="ld-ring-svg">
            <circle cx="40" cy="40" r="34" fill="none" stroke="rgba(255,255,255,.15)" stroke-width="6"/>
            <circle cx="40" cy="40" r="34" fill="none" :stroke="scoreColor(loc.total_score)" stroke-width="6" stroke-linecap="round" :stroke-dasharray="213.6" :stroke-dashoffset="213.6-(213.6*loc.total_score/100)" transform="rotate(-90 40 40)"/>
          </svg>
          <div class="ld-ring-center">
            <div class="ld-ring-val">{{ loc.total_score }}</div>
            <div class="ld-ring-label">Score</div>
          </div>
        </div>
      </div>
      <div class="ld-hero-badges">
        <span class="lc-badge" :class="'ls-'+loc.status">{{ statusLabel(loc.status) }}</span>
        <span class="ld-score-cat" :style="{color:scoreColor(loc.total_score)}">{{ loc.score_category || '-' }}</span>
        <span v-if="loc.partnership" class="ld-partnership-tag"><i class="ri-handshake-line"></i> Partnership: {{ partnershipLabel(loc.partnership.status) }}</span>
      </div>
    </div>

    <!-- Info Cards -->
    <div class="ld-info-grid">
      <div class="ld-info-card"><div class="ld-info-icon lic-blue"><i class="ri-money-dollar-circle-line"></i></div><div><div class="ld-info-label">Sewa/Tahun</div><div class="ld-info-val">{{ fc(loc.harga_sewa_per_tahun) }}</div></div></div>
      <div class="ld-info-card"><div class="ld-info-icon lic-green"><i class="ri-user-line"></i></div><div><div class="ld-info-label">Mitra</div><div class="ld-info-val">{{ loc.mitra?.name || '-' }}</div></div></div>
      <div class="ld-info-card"><div class="ld-info-icon lic-violet"><i class="ri-road-map-line"></i></div><div><div class="ld-info-label">Lebar Jalan</div><div class="ld-info-val">{{ loc.lebar_jalan }} m</div></div></div>
      <div class="ld-info-card"><div class="ld-info-icon lic-amber"><i class="ri-walk-line"></i></div><div><div class="ld-info-label">Traffic/Hari</div><div class="ld-info-val">{{ (loc.estimasi_lalu_lintas||0).toLocaleString() }}</div></div></div>
    </div>

    <!-- Tabs -->
    <div class="ld-tabs">
      <button v-for="t in tabs" :key="t.key" class="ld-tab" :class="{active:tab===t.key}" @click="tab=t.key">
        <i :class="t.icon"></i> {{ t.label }}
      </button>
    </div>

    <!-- Tab Content -->
    <div class="ld-content-card">

      <!-- OVERVIEW -->
      <template v-if="tab==='overview'">
        <h3 class="ld-section-title"><i class="ri-file-list-3-line"></i> Detail Lokasi</h3>
        <div class="ld-detail-grid">
          <div class="ld-detail-item"><span class="ld-dl"><i class="ri-time-line"></i> Durasi Sewa</span><span class="ld-dv">{{ loc.durasi_sewa }} tahun</span></div>
          <div class="ld-detail-item"><span class="ld-dl"><i class="ri-stack-line"></i> Jumlah Lantai</span><span class="ld-dv">{{ loc.jumlah_lantai }}</span></div>
          <div class="ld-detail-item"><span class="ld-dl"><i class="ri-store-2-line"></i> Kompetitor (500m)</span><span class="ld-dv">{{ loc.jumlah_kompetitor }}</span></div>
          <div class="ld-detail-item"><span class="ld-dl"><i class="ri-map-pin-range-line"></i> Dekat Dengan</span><span class="ld-dv">{{ loc.dekat_dengan || '-' }}</span></div>
          <div class="ld-detail-item"><span class="ld-dl"><i class="ri-group-line"></i> Target Market</span><span class="ld-dv">{{ loc.target_market || '-' }}</span></div>
          <div v-if="loc.partnership" class="ld-detail-item"><span class="ld-dl"><i class="ri-handshake-line"></i> Partnership</span><span class="ld-dv">{{ partnershipLabel(loc.partnership.status) }}</span></div>
        </div>
        <div class="ld-status-section">
          <h3 class="ld-section-title"><i class="ri-flow-chart"></i> Status Pipeline</h3>
          <div class="ld-pipeline">
            <div v-for="(st, idx) in statusSteps" :key="st.val" class="ld-pipe-step"
              :class="{ active: loc.status === st.val, completed: statusIndex(loc.status) > idx, selected: newStatus === st.val && newStatus !== loc.status }"
              @click="newStatus = st.val">
              <div class="ld-pipe-connector" v-if="idx > 0"></div>
              <div class="ld-pipe-dot" :style="dotStyle(st, idx)">
                <i :class="statusIndex(loc.status) > idx ? 'ri-check-line' : st.icon"></i>
              </div>
              <div class="ld-pipe-info">
                <div class="ld-pipe-label">{{ st.label }}</div>
                <div class="ld-pipe-desc">{{ st.desc }}</div>
              </div>
            </div>
          </div>
          <transition name="slide-fade">
            <div v-if="newStatus !== loc.status" class="ld-status-confirm">
              <div class="ld-status-change-info">
                <span class="ld-sc-from"><i class="ri-arrow-right-s-line"></i> {{ statusLabel(loc.status) }}</span>
                <i class="ri-arrow-right-line ld-sc-arrow"></i>
                <span class="ld-sc-to">{{ statusLabel(newStatus) }}</span>
              </div>
              <button @click="updateStatus" class="lc-btn-primary ld-confirm-btn" :disabled="updatingStatus">
                <i :class="updatingStatus ? 'ri-loader-4-line ri-spin' : 'ri-check-double-line'"></i>
                {{ updatingStatus ? 'Memproses...' : 'Konfirmasi Perubahan' }}
              </button>
            </div>
          </transition>
        </div>
      </template>

      <!-- SCORING -->
      <template v-else-if="tab==='scoring'">
        <div class="ld-score-header">
          <h3 class="ld-section-title"><i class="ri-bar-chart-2-fill"></i> Scoring Lokasi</h3>
          <button @click="recalculate" class="lc-btn-primary" style="font-size:.78rem;padding:8px 18px"><i class="ri-refresh-line"></i> Hitung Ulang</button>
        </div>
        <div class="ld-score-grid">
          <div v-for="s in scoreItems" :key="s.key" class="ld-score-card">
            <div class="ld-sc-head">
              <span class="ld-sc-label"><i :class="s.icon"></i> {{ s.label }}</span>
              <span class="ld-sc-weight">{{ s.weight }}%</span>
            </div>
            <div class="ld-sc-bar"><div class="ld-sc-fill" :style="{width:s.value+'%',background:scoreColor(s.value)}"></div></div>
            <div class="ld-sc-val">{{ s.value }}<span class="ld-sc-max">/100</span></div>
          </div>
        </div>
        <div class="ld-score-total">
          <div class="ld-st-label">Total Score</div>
          <div class="ld-st-val" :style="{color:scoreColor(loc.total_score)}">{{ loc.total_score }}</div>
          <div class="ld-st-cat" :style="{color:scoreColor(loc.total_score)}">{{ loc.score_category }}</div>
        </div>
      </template>

      <!-- SURVEY -->
      <template v-else-if="tab==='survey'">
        <div class="ld-section-head-row">
          <h3 class="ld-section-title"><i class="ri-survey-line"></i> Data Survei</h3>
          <button @click="showSurveyModal=true" class="lc-btn-primary" style="font-size:.78rem;padding:8px 18px"><i class="ri-add-line"></i> Tambah Survei</button>
        </div>
        <div v-if="loc.surveys?.length" class="ld-survey-list">
          <div v-for="s in loc.surveys" :key="s.id" class="ld-survey-card">
            <div class="ld-sv-head">
              <span class="ld-sv-date"><i class="ri-calendar-line"></i> {{ formatDate(s.survey_date) }}</span>
              <span class="ld-sv-by"><i class="ri-user-3-line"></i> {{ s.surveyor?.name || '-' }}</span>
            </div>
            <div class="ld-sv-body">
              <p>{{ s.hasil_survey || 'Belum ada hasil.' }}</p>
              <p v-if="s.catatan_survey" class="ld-sv-note"><i class="ri-sticky-note-line"></i> {{ s.catatan_survey }}</p>
            </div>
            <div class="ld-sv-foot">
              <span><i class="ri-money-dollar-circle-line"></i> Omzet: <strong>{{ fc(s.estimasi_omzet) }}</strong>/bln</span>
              <span><i class="ri-timer-line"></i> BEP: <strong>{{ s.estimasi_bep }}</strong> bulan</span>
            </div>
          </div>
        </div>
        <div v-else class="ld-empty-text"><i class="ri-survey-line" style="font-size:1.5rem;display:block;margin-bottom:8px"></i>Belum ada data survei</div>
      </template>

      <!-- APPROVAL -->
      <template v-else-if="tab==='approval'">
        <div class="ld-section-head-row">
          <h3 class="ld-section-title"><i class="ri-shield-check-line"></i> Riwayat Approval</h3>
          <button @click="showApprovalModal=true" class="lc-btn-primary" style="font-size:.78rem;padding:8px 18px"><i class="ri-add-line"></i> Beri Keputusan</button>
        </div>
        <div v-if="loc.approvals?.length" class="ld-approval-list">
          <div v-for="a in loc.approvals" :key="a.id" class="ld-appr-card" :class="'apc-'+a.decision">
            <div class="ld-appr-decision">
              <i :class="{'ri-checkbox-circle-fill':a.decision==='approved','ri-close-circle-fill':a.decision==='rejected','ri-refresh-line':a.decision==='revision'}"></i>
              {{ {approved:'Disetujui',rejected:'Ditolak',revision:'Revisi'}[a.decision] || a.decision }}
            </div>
            <div class="ld-appr-by"><i class="ri-user-line"></i> {{ a.approver?.name || '-' }} · {{ formatDate(a.approved_at) }}</div>
            <div v-if="a.note" class="ld-appr-note">{{ a.note }}</div>
          </div>
        </div>
        <div v-else class="ld-empty-text"><i class="ri-shield-check-line" style="font-size:1.5rem;display:block;margin-bottom:8px"></i>Belum ada riwayat approval</div>
      </template>

      <!-- DOCUMENTS -->
      <template v-else-if="tab==='documents'">
        <div class="ld-section-head-row">
          <h3 class="ld-section-title"><i class="ri-folder-image-line"></i> Dokumen & Foto</h3>
          <label class="lc-btn-primary" style="font-size:.78rem;padding:8px 18px;cursor:pointer">
            <i class="ri-upload-cloud-line"></i> Upload
            <input ref="fileUploadInput" type="file" hidden accept="image/*,.pdf" @change="onFileUpload" />
          </label>
        </div>
        <div v-if="loc.files?.length" class="ld-files-grid">
          <div v-for="f in loc.files" :key="f.id" class="ld-file-card">
            <img v-if="isImage(f.file_url)" :src="f.file_url" class="ld-file-thumb" @click="window.open(f.file_url,'_blank')" />
            <div v-else class="ld-file-icon"><i class="ri-file-3-line"></i></div>
            <div class="ld-file-info">
              <span class="ld-file-type">{{ f.file_type }}</span>
              <span class="ld-file-label">{{ f.label || 'File' }}</span>
            </div>
            <button @click="deleteFile(f.id)" class="ld-file-del"><i class="ri-delete-bin-6-line"></i></button>
          </div>
        </div>
        <div v-else class="ld-empty-text"><i class="ri-folder-image-line" style="font-size:1.5rem;display:block;margin-bottom:8px"></i>Belum ada dokumen</div>
      </template>
    </div>

    <!-- Survey Modal -->
    <Teleport to="body">
      <div v-if="showSurveyModal" class="lc-overlay" @click.self="showSurveyModal=false">
        <div class="lc-modal" style="max-width:540px" @click.stop>
          <div class="lc-modal-head">
            <h3><i class="ri-survey-fill" style="color:#7c3aed;margin-right:6px"></i> Input Hasil Survei</h3>
            <button @click="showSurveyModal=false" class="lc-modal-x"><i class="ri-close-line"></i></button>
          </div>
          <form @submit.prevent="createSurvey" class="lc-modal-body">
            <div class="lc-frow">
              <div class="lc-fg"><label><i class="ri-calendar-line"></i> Tanggal Survei <span class="req">*</span></label><input v-model="surveyForm.survey_date" type="datetime-local" class="lc-input" required /></div>
              <div class="lc-fg"><label><i class="ri-user-3-line"></i> Surveyor ID</label><input v-model="surveyForm.survey_by" class="lc-input" placeholder="UUID user" /></div>
            </div>
            <div class="lc-fg"><label><i class="ri-file-text-line"></i> Hasil Survei <span class="req">*</span></label><textarea v-model="surveyForm.hasil_survey" class="lc-input lc-textarea" required placeholder="Tulis hasil survei lokasi..."></textarea></div>
            <div class="lc-fg"><label><i class="ri-sticky-note-line"></i> Catatan Tambahan</label><textarea v-model="surveyForm.catatan_survey" class="lc-input lc-textarea" placeholder="Catatan lain..."></textarea></div>
            <div class="lc-frow">
              <div class="lc-fg"><label><i class="ri-money-dollar-circle-line"></i> Est. Omzet/Bulan</label><input v-model.number="surveyForm.estimasi_omzet" type="number" class="lc-input" /></div>
              <div class="lc-fg"><label><i class="ri-timer-line"></i> Est. BEP (bulan)</label><input v-model.number="surveyForm.estimasi_bep" type="number" class="lc-input" /></div>
            </div>
            <div class="lc-modal-foot"><button type="button" @click="showSurveyModal=false" class="lc-btn-sec">Batal</button><button type="submit" class="lc-btn-primary"><i class="ri-check-line"></i> Simpan Survei</button></div>
          </form>
        </div>
      </div>
    </Teleport>

    <!-- Approval Modal -->
    <Teleport to="body">
      <div v-if="showApprovalModal" class="lc-overlay" @click.self="showApprovalModal=false">
        <div class="lc-modal" style="max-width:480px" @click.stop>
          <div class="lc-modal-head">
            <h3><i class="ri-shield-check-fill" style="color:#16a34a;margin-right:6px"></i> Keputusan Approval</h3>
            <button @click="showApprovalModal=false" class="lc-modal-x"><i class="ri-close-line"></i></button>
          </div>
          <form @submit.prevent="submitApproval" class="lc-modal-body">
            <div class="lc-fg"><label><i class="ri-checkbox-circle-line"></i> Keputusan <span class="req">*</span></label>
              <select v-model="approvalForm.decision" class="lc-input" required>
                <option value="">Pilih keputusan</option>
                <option value="approved">✅ Setujui</option>
                <option value="rejected">❌ Tolak</option>
                <option value="revision">🔄 Minta Revisi</option>
              </select>
            </div>
            <div class="lc-fg"><label><i class="ri-chat-3-line"></i> Catatan {{ approvalForm.decision==='rejected'?'(wajib)':'' }}</label><textarea v-model="approvalForm.note" class="lc-input lc-textarea" :required="approvalForm.decision==='rejected'" placeholder="Alasan keputusan..."></textarea></div>
            <div class="lc-modal-foot"><button type="button" @click="showApprovalModal=false" class="lc-btn-sec">Batal</button><button type="submit" class="lc-btn-primary"><i class="ri-check-line"></i> Submit</button></div>
          </form>
        </div>
      </div>
    </Teleport>
  </div>
  <div v-else class="ld-loading"><i class="ri-loader-4-line ri-spin" style="font-size:1.5rem"></i> Memuat data...</div>
</template>

<script setup>
import { ref, reactive, computed, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { locationApi, uploadApi } from '../../services/api'
import { useToastStore } from '../../stores/toast'

const toast = useToastStore()
const route = useRoute()
const id = route.params.id
const loc = ref(null)
const tab = ref('overview')
const newStatus = ref('')
const showSurveyModal = ref(false)
const showApprovalModal = ref(false)
const surveyForm = reactive({ survey_date:'', survey_by:'', hasil_survey:'', catatan_survey:'', estimasi_omzet:0, estimasi_bep:0 })
const approvalForm = reactive({ decision:'', note:'' })

const tabs = [
  { key:'overview', label:'Overview', icon:'ri-eye-line' },
  { key:'scoring', label:'Scoring', icon:'ri-bar-chart-2-line' },
  { key:'survey', label:'Survei', icon:'ri-survey-line' },
  { key:'approval', label:'Approval', icon:'ri-shield-check-line' },
  { key:'documents', label:'Dokumen', icon:'ri-folder-image-line' },
]

const allStatuses = ['DRAFT','SUBMITTED','IN_REVIEW','SURVEY_SCHEDULED','SURVEYED','APPROVED','REJECTED','REVISION_NEEDED']
const updatingStatus = ref(false)

const statusSteps = [
  { val:'DRAFT', label:'Draft', desc:'Pengajuan awal', icon:'ri-draft-line', color:'#94a3b8' },
  { val:'SUBMITTED', label:'Diajukan', desc:'Menunggu review', icon:'ri-send-plane-line', color:'#3b82f6' },
  { val:'IN_REVIEW', label:'Ditinjau', desc:'Sedang ditinjau', icon:'ri-search-eye-line', color:'#f59e0b' },
  { val:'SURVEY_SCHEDULED', label:'Survei', desc:'Dijadwalkan survei', icon:'ri-calendar-check-line', color:'#8b5cf6' },
  { val:'SURVEYED', label:'Disurvei', desc:'Sudah disurvei', icon:'ri-checkbox-circle-line', color:'#6366f1' },
  { val:'APPROVED', label:'Disetujui', desc:'Lokasi disetujui', icon:'ri-checkbox-circle-fill', color:'#22c55e' },
  { val:'REJECTED', label:'Ditolak', desc:'Lokasi ditolak', icon:'ri-close-circle-fill', color:'#ef4444' },
  { val:'REVISION_NEEDED', label:'Revisi', desc:'Perlu perbaikan', icon:'ri-error-warning-line', color:'#f97316' },
]

function statusIndex(s) { return statusSteps.findIndex(st => st.val === s) }
function dotStyle(st, idx) {
  const current = statusIndex(loc.value?.status)
  if (idx === current) return { background: st.color, color: '#fff', boxShadow: `0 0 0 4px ${st.color}22` }
  if (idx < current) return { background: statusSteps[current]?.color || '#22c55e', color: '#fff' }
  if (newStatus.value === st.val && newStatus.value !== loc.value?.status) return { background: st.color + '20', color: st.color, border: `2px solid ${st.color}` }
  return {}
}

const scoreItems = computed(() => [
  { key:'traffic', label:'Traffic', weight:30, value: loc.value?.score_traffic||0, icon:'ri-walk-line' },
  { key:'sewa', label:'Harga Sewa', weight:20, value: loc.value?.score_sewa||0, icon:'ri-money-dollar-circle-line' },
  { key:'kompetitor', label:'Kompetitor', weight:20, value: loc.value?.score_kompetitor||0, icon:'ri-store-2-line' },
  { key:'akses', label:'Akses Jalan', weight:15, value: loc.value?.score_akses||0, icon:'ri-road-map-line' },
  { key:'market', label:'Target Market', weight:15, value: loc.value?.score_market||0, icon:'ri-group-line' },
])

function statusLabel(s) { return { DRAFT:'Draft', SUBMITTED:'Diajukan', IN_REVIEW:'Ditinjau', SURVEY_SCHEDULED:'Survei', SURVEYED:'Disurvei', APPROVED:'Disetujui', REJECTED:'Ditolak', REVISION_NEEDED:'Revisi' }[s] || s }
function partnershipLabel(s) { return { PENDING:'Menunggu', DP_VERIFIED:'DP Terverifikasi', AGREEMENT_SIGNED:'Agreement Signed', DEVELOPMENT:'Pembangunan', RUNNING:'Berjalan', COMPLETED:'Selesai' }[s] || s }
function scoreColor(s) { if (s >= 80) return '#22c55e'; if (s >= 65) return '#6366f1'; if (s >= 50) return '#f59e0b'; return '#ef4444' }
function fc(n) { return 'Rp' + (n||0).toLocaleString('id-ID') }
function formatDate(d) { return d ? new Date(d).toLocaleDateString('id-ID', { day:'2-digit', month:'short', year:'numeric' }) : '-' }
function isImage(url) { return /\.(jpg|jpeg|png|gif|webp|svg)$/i.test(url||'') }

async function loadData() {
  try { const { data } = await locationApi.getByID(id); loc.value = data.data; newStatus.value = loc.value.status } catch(e) { toast.error('Gagal memuat data') }
}
async function updateStatus() {
  updatingStatus.value = true
  try { await locationApi.updateStatus(id, newStatus.value); toast.success('Status berhasil diupdate'); await loadData() } catch(e) { toast.error(e.response?.data?.error || 'Gagal') }
  finally { updatingStatus.value = false }
}
async function recalculate() {
  try { await locationApi.recalculate(id); toast.success('Score dihitung ulang'); await loadData() } catch(e) { toast.error('Gagal') }
}
async function createSurvey() {
  try {
    const payload = { ...surveyForm }; if (payload.survey_date) payload.survey_date = new Date(payload.survey_date).toISOString(); if (!payload.survey_by) delete payload.survey_by
    await locationApi.createSurvey(id, payload); toast.success('Survei disimpan'); showSurveyModal.value = false; Object.assign(surveyForm, { survey_date:'', survey_by:'', hasil_survey:'', catatan_survey:'', estimasi_omzet:0, estimasi_bep:0 }); await loadData()
  } catch(e) { toast.error(e.response?.data?.error || 'Gagal') }
}
async function submitApproval() {
  try { await locationApi.approve(id, approvalForm); toast.success('Keputusan disimpan'); showApprovalModal.value = false; Object.assign(approvalForm, { decision:'', note:'' }); await loadData() } catch(e) { toast.error(e.response?.data?.error || 'Gagal') }
}
async function onFileUpload(e) {
  const file = e.target.files[0]; if (!file) return
  try { const { data: upRes } = await uploadApi.upload(file); const url = upRes.data?.url || upRes.url || ''; await locationApi.addFile(id, { file_url: url, file_type: file.type.startsWith('image/')?'photo':'document', label: file.name }); toast.success('File diupload'); await loadData() } catch(e) { toast.error('Gagal upload') }
}
async function deleteFile(fileId) {
  try { await locationApi.deleteFile(id, fileId); toast.success('File dihapus'); await loadData() } catch(e) { toast.error('Gagal hapus') }
}

onMounted(loadData)
</script>

<style scoped>
.ld-page { margin: 0; }
.ld-back { display:inline-flex; align-items:center; gap:4px; font-size:.82rem; color:#6366f1; text-decoration:none; font-weight:600; margin-bottom:16px; padding:6px 12px; border-radius:8px; transition:all .2s; }
.ld-back:hover { background:#eef2ff; }
.ld-loading { display:flex; align-items:center; justify-content:center; gap:8px; padding:80px; color:#94a3b8; font-size:.9rem; }

/* ═══ HERO ═══ */
.ld-hero { background:linear-gradient(135deg,#0f0c29 0%,#302b63 50%,#24243e 100%); border-radius:16px; padding:28px 32px; margin-bottom:20px; color:#fff; }
.ld-hero-top { display:flex; align-items:center; justify-content:space-between; gap:20px; margin-bottom:14px; }
.ld-hero-title { font-size:1.4rem; font-weight:800; margin:0; }
.ld-hero-addr { font-size:.82rem; opacity:.7; margin:4px 0 10px; display:flex; align-items:center; gap:4px; }
.ld-hero-addr i { font-size:.85rem; }
.ld-hero-meta { display:flex; gap:14px; font-size:.75rem; opacity:.5; }
.ld-hero-meta i { margin-right:2px; }
.ld-hero-badges { display:flex; gap:10px; align-items:center; padding-top:14px; border-top:1px solid rgba(255,255,255,.08); }
.ld-score-cat { font-size:.78rem; font-weight:700; }
.ld-partnership-tag { display:inline-flex; align-items:center; gap:4px; font-size:.72rem; font-weight:600; color:rgba(255,255,255,.7); background:rgba(255,255,255,.1); padding:4px 10px; border-radius:6px; }
.ld-partnership-tag i { font-size:.8rem; }

/* Score Ring */
.ld-hero-score-ring { position:relative; width:80px; height:80px; flex-shrink:0; }
.ld-ring-svg { width:80px; height:80px; filter:drop-shadow(0 2px 8px rgba(0,0,0,.2)); }
.ld-ring-center { position:absolute; inset:0; display:flex; flex-direction:column; align-items:center; justify-content:center; }
.ld-ring-val { font-size:1.3rem; font-weight:800; }
.ld-ring-label { font-size:.55rem; opacity:.7; text-transform:uppercase; letter-spacing:.08em; }

/* Info Cards */
.ld-info-grid { display:grid; grid-template-columns:repeat(4,1fr); gap:12px; margin-bottom:20px; }
.ld-info-card { display:flex; gap:12px; align-items:center; background:#fff; border:1px solid #eef1f6; border-radius:12px; padding:14px 16px; }
.ld-info-icon { width:36px; height:36px; border-radius:10px; display:flex; align-items:center; justify-content:center; flex-shrink:0; font-size:1.05rem; }
.lic-blue { background:#dbeafe; color:#2563eb; }
.lic-green { background:#dcfce7; color:#16a34a; }
.lic-violet { background:#ede9fe; color:#7c3aed; }
.lic-amber { background:#fef3c7; color:#d97706; }
.ld-info-label { font-size:.65rem; color:#94a3b8; font-weight:600; text-transform:uppercase; letter-spacing:.04em; }
.ld-info-val { font-size:.88rem; font-weight:700; color:#0f172a; margin-top:2px; }

/* Tabs */
.ld-tabs { display:flex; gap:4px; border-bottom:2px solid #f1f5f9; margin-bottom:20px; }
.ld-tab { padding:10px 18px; font-size:.82rem; font-weight:600; color:#94a3b8; background:none; border:none; border-bottom:2px solid transparent; margin-bottom:-2px; cursor:pointer; display:flex; align-items:center; gap:5px; font-family:inherit; transition:all .2s; }
.ld-tab i { font-size:.9rem; }
.ld-tab:hover { color:#334155; }
.ld-tab.active { color:#6366f1; border-bottom-color:#6366f1; }

/* Content */
.ld-content-card { background:#fff; border-radius:14px; border:1px solid #eef1f6; padding:24px; min-height:200px; }
.ld-section-head-row { display:flex; align-items:center; justify-content:space-between; margin-bottom:16px; }
.ld-section-head-row .ld-section-title { margin-bottom:0; }
.ld-section-title { display:flex; align-items:center; gap:6px; font-size:.92rem; font-weight:700; color:#0f172a; margin:0 0 14px; }
.ld-section-title i { color:#6366f1; font-size:1rem; }
.ld-empty-text { color:#94a3b8; font-size:.85rem; text-align:center; padding:40px; }

/* Detail Grid */
.ld-detail-grid { display:grid; grid-template-columns:repeat(2,1fr); gap:10px; margin-bottom:20px; }
.ld-detail-item { display:flex; justify-content:space-between; padding:10px 14px; background:#f8fafc; border-radius:8px; }
.ld-dl { font-size:.78rem; color:#64748b; font-weight:500; display:flex; align-items:center; gap:4px; }
.ld-dl i { font-size:.8rem; color:#94a3b8; }
.ld-dv { font-size:.82rem; font-weight:700; color:#0f172a; }

/* Map */
.ld-map-container { border-radius:12px; overflow:hidden; border:1px solid #eef1f6; margin-bottom:20px; }
.ld-map-frame { width:100%; height:300px; border:none; }

/* Status Pipeline */
.ld-status-section { margin-top:20px; padding-top:20px; border-top:1px solid #f1f5f9; }
.ld-pipeline { display:flex; flex-wrap:wrap; gap:0; margin-bottom:16px; }
.ld-pipe-step { position:relative; display:flex; align-items:center; gap:10px; padding:10px 16px 10px 10px; border-radius:10px; cursor:pointer; transition:all .2s; flex: 1; min-width:130px; }
.ld-pipe-step:hover { background:#f8fafc; }
.ld-pipe-step.active { background:#f0f9ff; }
.ld-pipe-step.selected { background:#faf5ff; }
.ld-pipe-connector { position:absolute; left:-8px; top:50%; width:16px; height:2px; background:#e2e8f0; }
.ld-pipe-step.completed .ld-pipe-connector { background:#22c55e; }
.ld-pipe-dot { width:32px; height:32px; border-radius:50%; display:flex; align-items:center; justify-content:center; font-size:.85rem; background:#f1f5f9; color:#94a3b8; flex-shrink:0; transition:all .25s; border:2px solid transparent; }
.ld-pipe-step.active .ld-pipe-dot { transform:scale(1.1); }
.ld-pipe-info { min-width:0; }
.ld-pipe-label { font-size:.75rem; font-weight:700; color:#334155; white-space:nowrap; }
.ld-pipe-step.active .ld-pipe-label { color:#0f172a; }
.ld-pipe-step.completed .ld-pipe-label { color:#64748b; }
.ld-pipe-desc { font-size:.6rem; color:#94a3b8; white-space:nowrap; overflow:hidden; text-overflow:ellipsis; }

.ld-status-confirm { display:flex; align-items:center; justify-content:space-between; gap:16px; padding:14px 20px; background:linear-gradient(135deg,#f8fafc,#eef2ff); border:1.5px solid #c7d2fe; border-radius:12px; margin-top:4px; }
.ld-status-change-info { display:flex; align-items:center; gap:8px; font-size:.82rem; font-weight:600; }
.ld-sc-from { color:#64748b; display:flex; align-items:center; gap:2px; }
.ld-sc-arrow { color:#6366f1; font-size:1rem; }
.ld-sc-to { color:#6366f1; font-weight:700; }
.ld-confirm-btn { font-size:.8rem !important; padding:9px 20px !important; }

.slide-fade-enter-active { transition:all .3s ease; }
.slide-fade-leave-active { transition:all .2s ease; }
.slide-fade-enter-from { transform:translateY(-8px); opacity:0; }
.slide-fade-leave-to { transform:translateY(-8px); opacity:0; }

/* Score */
.ld-score-header { display:flex; align-items:center; justify-content:space-between; margin-bottom:16px; }
.ld-score-header .ld-section-title { margin-bottom:0; }
.ld-score-grid { display:grid; grid-template-columns:repeat(5,1fr); gap:12px; margin-bottom:20px; }
.ld-score-card { background:#f8fafc; border:1px solid #eef1f6; border-radius:12px; padding:16px; text-align:center; }
.ld-sc-head { display:flex; justify-content:space-between; margin-bottom:10px; align-items:center; }
.ld-sc-label { font-size:.72rem; font-weight:700; color:#334155; display:flex; align-items:center; gap:4px; }
.ld-sc-label i { color:#94a3b8; }
.ld-sc-weight { font-size:.6rem; font-weight:700; color:#94a3b8; background:#f1f5f9; padding:2px 6px; border-radius:4px; }
.ld-sc-bar { height:8px; background:#e2e8f0; border-radius:4px; overflow:hidden; margin-bottom:8px; }
.ld-sc-fill { height:100%; border-radius:4px; transition:width .4s; }
.ld-sc-val { font-size:1.3rem; font-weight:800; color:#0f172a; }
.ld-sc-max { font-size:.7rem; color:#94a3b8; font-weight:500; }
.ld-score-total { text-align:center; padding:24px; background:#f8fafc; border-radius:14px; border:1px solid #eef1f6; }
.ld-st-label { font-size:.7rem; color:#94a3b8; font-weight:600; text-transform:uppercase; letter-spacing:.06em; }
.ld-st-val { font-size:2.4rem; font-weight:800; margin:4px 0; }
.ld-st-cat { font-size:.88rem; font-weight:700; }

/* Survey */
.ld-survey-card { background:#f8fafc; border:1px solid #eef1f6; border-radius:12px; padding:16px; margin-bottom:10px; }
.ld-sv-head { display:flex; justify-content:space-between; margin-bottom:10px; }
.ld-sv-date { font-size:.82rem; font-weight:700; color:#0f172a; display:flex; align-items:center; gap:4px; }
.ld-sv-date i { color:#6366f1; }
.ld-sv-by { font-size:.72rem; color:#94a3b8; display:flex; align-items:center; gap:3px; }
.ld-sv-body { font-size:.82rem; color:#334155; line-height:1.6; }
.ld-sv-body p { margin:0 0 6px; }
.ld-sv-note { color:#64748b; font-style:italic; display:flex; align-items:flex-start; gap:4px; }
.ld-sv-note i { margin-top:3px; color:#94a3b8; }
.ld-sv-foot { display:flex; gap:20px; margin-top:10px; padding-top:10px; border-top:1px solid #eef1f6; font-size:.78rem; color:#64748b; }
.ld-sv-foot i { color:#94a3b8; }
.ld-sv-foot strong { color:#0f172a; }

/* Approval */
.ld-appr-card { border-radius:12px; padding:16px; margin-bottom:10px; border-left:4px solid; }
.apc-approved { background:#f0fdf4; border-color:#22c55e; }
.apc-rejected { background:#fef2f2; border-color:#ef4444; }
.apc-revision { background:#fffbeb; border-color:#f59e0b; }
.ld-appr-decision { font-size:.88rem; font-weight:700; color:#0f172a; margin-bottom:4px; display:flex; align-items:center; gap:5px; }
.apc-approved .ld-appr-decision i { color:#22c55e; }
.apc-rejected .ld-appr-decision i { color:#ef4444; }
.apc-revision .ld-appr-decision i { color:#f59e0b; }
.ld-appr-by { font-size:.72rem; color:#94a3b8; margin-bottom:6px; display:flex; align-items:center; gap:3px; }
.ld-appr-note { font-size:.82rem; color:#334155; background:rgba(0,0,0,.03); padding:10px 14px; border-radius:8px; margin-top:8px; }

/* Files */
.ld-files-grid { display:grid; grid-template-columns:repeat(auto-fill,minmax(160px,1fr)); gap:12px; }
.ld-file-card { position:relative; background:#f8fafc; border:1px solid #eef1f6; border-radius:12px; overflow:hidden; }
.ld-file-thumb { width:100%; height:120px; object-fit:cover; cursor:pointer; }
.ld-file-icon { height:120px; display:flex; align-items:center; justify-content:center; font-size:2.5rem; color:#94a3b8; }
.ld-file-info { padding:10px 12px; }
.ld-file-type { font-size:.6rem; font-weight:700; text-transform:uppercase; color:#94a3b8; }
.ld-file-label { display:block; font-size:.72rem; font-weight:600; color:#334155; margin-top:2px; overflow:hidden; text-overflow:ellipsis; white-space:nowrap; }
.ld-file-del { position:absolute; top:6px; right:6px; width:26px; height:26px; border-radius:50%; border:none; background:rgba(239,68,68,.9); color:#fff; cursor:pointer; font-size:.8rem; display:flex; align-items:center; justify-content:center; opacity:0; transition:all .2s; }
.ld-file-card:hover .ld-file-del { opacity:1; }

/* Shared Badges */
.lc-badge { font-size:.68rem; font-weight:700; text-transform:uppercase; letter-spacing:.04em; padding:4px 10px; border-radius:6px; }
.ls-DRAFT { background:rgba(255,255,255,.2); color:#fff; }
.ls-SUBMITTED { background:rgba(219,234,254,.3); color:#bfdbfe; }
.ls-IN_REVIEW { background:rgba(254,243,199,.3); color:#fde68a; }
.ls-SURVEY_SCHEDULED { background:rgba(224,231,255,.3); color:#c7d2fe; }
.ls-SURVEYED { background:rgba(237,233,254,.3); color:#ddd6fe; }
.ls-APPROVED { background:rgba(220,252,231,.3); color:#bbf7d0; }
.ls-REJECTED { background:rgba(254,226,226,.3); color:#fecaca; }
.ls-REVISION_NEEDED { background:rgba(255,237,213,.3); color:#fed7aa; }

/* Shared modal, form, btn from LocationList are globally scoped via Remix Icon CDN */
.lc-overlay { position:fixed; inset:0; background:rgba(0,0,0,.5); display:flex; align-items:flex-start; justify-content:center; z-index:1000; backdrop-filter:blur(4px); padding-top:32px; overflow-y:auto; }
.lc-modal { background:#fff; border-radius:18px; width:100%; max-width:580px; box-shadow:0 24px 80px rgba(0,0,0,.2); margin-bottom:40px; }
.lc-modal-head { display:flex; justify-content:space-between; align-items:center; padding:22px 28px; border-bottom:1px solid #f1f5f9; }
.lc-modal-head h3 { font-size:1.1rem; font-weight:700; margin:0; color:#0f172a; display:flex; align-items:center; }
.lc-modal-x { width:36px; height:36px; border-radius:10px; display:flex; align-items:center; justify-content:center; border:none; background:transparent; font-size:1.3rem; color:#94a3b8; cursor:pointer; }
.lc-modal-x:hover { background:#f1f5f9; color:#0f172a; }
.lc-modal-body { padding:24px 28px; }
.lc-modal-foot { display:flex; justify-content:flex-end; gap:10px; padding-top:18px; border-top:1px solid #f1f5f9; margin-top:6px; }
.lc-fg { margin-bottom:12px; }
.lc-fg label { display:block; font-size:.78rem; font-weight:600; margin-bottom:5px; color:#334155; }
.lc-fg label i { color:#94a3b8; margin-right:3px; }
.req { color:#ef4444; }
.lc-input { width:100%; padding:10px 14px; border:1.5px solid #e2e8f0; border-radius:10px; font-size:.85rem; background:#fff; color:#1e293b; outline:none; box-sizing:border-box; font-family:inherit; }
.lc-input:focus { border-color:#6366f1; box-shadow:0 0 0 3px rgba(99,102,241,.08); }
.lc-textarea { min-height:70px; resize:vertical; }
.lc-frow { display:grid; grid-template-columns:1fr 1fr; gap:12px; }
.lc-btn-primary { display:inline-flex; align-items:center; gap:6px; padding:11px 24px; font-size:.85rem; font-weight:700; border-radius:12px; border:none; cursor:pointer; background:linear-gradient(135deg,#6366f1,#8b5cf6); color:#fff; white-space:nowrap; font-family:inherit; }
.lc-btn-sec { padding:11px 22px; border-radius:12px; font-size:.85rem; font-weight:600; background:#f1f5f9; color:#475569; border:none; cursor:pointer; font-family:inherit; }

.ri-spin { animation:riSpin .8s linear infinite; }
@keyframes riSpin { from{transform:rotate(0deg)} to{transform:rotate(360deg)} }

@media (max-width:768px) {
  .ld-hero { padding:20px; }
  .ld-hero-top { flex-direction:column; text-align:center; }
  .ld-info-grid { grid-template-columns:repeat(2,1fr); }
  .ld-score-grid { grid-template-columns:repeat(2,1fr); }
  .ld-detail-grid { grid-template-columns:1fr; }
  .ld-tabs { overflow-x:auto; }
  .lc-frow { grid-template-columns:1fr; }
}
</style>
