<template>
  <div class="mt-page">
    <!-- Hero -->
    <div class="mt-hero">
      <div class="mt-hero-top">
        <div>
          <h1 class="mt-hero-title">Meeting Management</h1>
          <p class="mt-hero-sub">Kelola agenda, notulensi, dan action plan meeting</p>
        </div>
        <button @click="openCreate" class="mt-btn-primary">
          <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><line x1="12" y1="5" x2="12" y2="19"/><line x1="5" y1="12" x2="19" y2="12"/></svg>
          Buat Meeting
        </button>
      </div>
      <div class="mt-stats">
        <div class="mt-stat"><span class="mt-stat-dot" style="background:#818cf8"></span><span class="mt-stat-label">Total</span><span class="mt-stat-val">{{ total }}</span></div>
        <div class="mt-stat"><span class="mt-stat-dot" style="background:#22c55e"></span><span class="mt-stat-label">Scheduled</span><span class="mt-stat-val">{{ meetings.filter(m=>m.status==='scheduled').length }}</span></div>
        <div class="mt-stat"><span class="mt-stat-dot" style="background:#f59e0b"></span><span class="mt-stat-label">Ongoing</span><span class="mt-stat-val">{{ meetings.filter(m=>m.status==='ongoing').length }}</span></div>
        <div class="mt-stat"><span class="mt-stat-dot" style="background:#06b6d4"></span><span class="mt-stat-label">Completed</span><span class="mt-stat-val">{{ meetings.filter(m=>m.status==='completed').length }}</span></div>
      </div>
    </div>

    <!-- Toolbar -->
    <div class="mt-toolbar">
      <div class="mt-search-wrap">
        <svg class="mt-search-ico" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="#94a3b8" stroke-width="2"><circle cx="11" cy="11" r="8"/><line x1="21" y1="21" x2="16.65" y2="16.65"/></svg>
        <input v-model="search" class="mt-search" placeholder="Cari meeting..." />
      </div>
      <div class="mt-filters">
        <select v-model="filterType" class="mt-select">
          <option value="">Semua Tipe</option>
          <option value="edukasi">Edukasi</option>
          <option value="closing">Closing</option>
          <option value="review_lokasi">Review Lokasi</option>
          <option value="operasional">Operasional</option>
        </select>
        <select v-model="filterStatus" class="mt-select">
          <option value="">Semua Status</option>
          <option value="scheduled">Scheduled</option>
          <option value="ongoing">Ongoing</option>
          <option value="completed">Completed</option>
          <option value="waiting_decision">Waiting Decision</option>
          <option value="approved">Approved</option>
          <option value="cancelled">Cancelled</option>
        </select>
      </div>
    </div>

    <!-- Loading -->
    <div v-if="loading" class="mt-list">
      <div v-for="n in 4" :key="n" class="mt-card mt-card-skel">
        <div class="skel-bar shimmer" style="width:60%;height:18px"></div>
        <div class="skel-bar shimmer" style="width:40%;height:14px"></div>
        <div class="skel-bar shimmer" style="width:80%;height:14px"></div>
      </div>
    </div>

    <!-- Meeting List -->
    <div v-else-if="filtered.length" class="mt-list">
      <div v-for="m in filtered" :key="m.id" class="mt-card" @click="openDetail(m.id)">
        <div class="mt-card-head">
          <div class="mt-card-type" :class="'type-'+m.meeting_type">{{ typeLabel(m.meeting_type) }}</div>
          <div class="mt-card-status" :class="'st-'+m.status">{{ statusLabel(m.status) }}</div>
        </div>
        <h3 class="mt-card-title">{{ m.title }}</h3>
        <div class="mt-card-meta">
          <span class="mt-meta-item">
            <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="#94a3b8" stroke-width="2"><rect x="3" y="4" width="18" height="18" rx="2" ry="2"/><line x1="16" y1="2" x2="16" y2="6"/><line x1="8" y1="2" x2="8" y2="6"/><line x1="3" y1="10" x2="21" y2="10"/></svg>
            {{ formatDate(m.meeting_date) }}
          </span>
          <span class="mt-meta-item">
            <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="#94a3b8" stroke-width="2"><circle cx="12" cy="12" r="10"/><polyline points="12 6 12 12 16 14"/></svg>
            {{ m.duration }} menit
          </span>
          <span v-if="m.participants" class="mt-meta-item">
            <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="#94a3b8" stroke-width="2"><path d="M17 21v-2a4 4 0 0 0-4-4H5a4 4 0 0 0-4 4v2"/><circle cx="9" cy="7" r="4"/></svg>
            {{ m.participants.length }} peserta
          </span>
          <span v-if="m.location" class="mt-meta-item">
            <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="#94a3b8" stroke-width="2"><path d="M21 10c0 7-9 13-9 13s-9-6-9-13a9 9 0 0 1 18 0z"/><circle cx="12" cy="10" r="3"/></svg>
            {{ m.location }}
          </span>
        </div>
      </div>
    </div>

    <!-- Empty -->
    <div v-else class="mt-empty">
      <div class="mt-empty-ico">
        <svg width="36" height="36" viewBox="0 0 24 24" fill="none" stroke="#94a3b8" stroke-width="1.5"><rect x="3" y="4" width="18" height="18" rx="2" ry="2"/><line x1="16" y1="2" x2="16" y2="6"/><line x1="8" y1="2" x2="8" y2="6"/><line x1="3" y1="10" x2="21" y2="10"/></svg>
      </div>
      <h3>Belum ada meeting</h3>
      <p>Buat meeting pertama untuk memulai.</p>
    </div>

    <!-- Create/Edit Modal -->
    <Teleport to="body">
      <div v-if="showForm" class="mt-overlay" @click.self="showForm=false">
        <div class="mt-modal" @click.stop>
          <div class="mt-modal-head">
            <h3>{{ editMode ? 'Edit Meeting' : 'Buat Meeting Baru' }}</h3>
            <button @click="showForm=false" class="mt-modal-x">&times;</button>
          </div>
          <form @submit.prevent="saveMeeting" class="mt-modal-body">
            <div class="mt-fg">
              <label>Judul Meeting <span class="req">*</span></label>
              <input v-model="form.title" class="mt-input" placeholder="Judul meeting" />
            </div>
            <div class="mt-frow">
              <div class="mt-fg">
                <label>Tipe <span class="req">*</span></label>
                <select v-model="form.meeting_type" class="mt-input">
                  <option value="edukasi">Edukasi</option>
                  <option value="closing">Closing</option>
                  <option value="review_lokasi">Review Lokasi</option>
                  <option value="operasional">Operasional</option>
                </select>
              </div>
              <div class="mt-fg">
                <label>Durasi (menit)</label>
                <input v-model.number="form.duration" type="number" class="mt-input" placeholder="60" />
              </div>
            </div>
            <div class="mt-fg">
              <label>Tanggal & Waktu <span class="req">*</span></label>
              <input v-model="form.meeting_date" type="datetime-local" class="mt-input" />
            </div>
            <div class="mt-frow">
              <div class="mt-fg">
                <label>Lokasi</label>
                <input v-model="form.location" class="mt-input" placeholder="Lokasi fisik (opsional)" />
              </div>
              <div class="mt-fg">
                <label>Link Meeting</label>
                <input v-model="form.meeting_link" class="mt-input" placeholder="Zoom/Google Meet URL" />
              </div>
            </div>
            <div v-if="editMode" class="mt-fg">
              <label>Status</label>
              <select v-model="form.status" class="mt-input">
                <option value="scheduled">Scheduled</option>
                <option value="ongoing">Ongoing</option>
                <option value="completed">Completed</option>
                <option value="waiting_decision">Waiting Decision</option>
                <option value="approved">Approved</option>
                <option value="cancelled">Cancelled</option>
              </select>
            </div>
            <div class="mt-modal-foot">
              <button type="button" @click="showForm=false" class="mt-btn-sec">Batal</button>
              <button type="submit" class="mt-btn-primary" :disabled="saving">{{ saving ? 'Menyimpan...' : (editMode ? 'Update' : 'Simpan') }}</button>
            </div>
          </form>
        </div>
      </div>
    </Teleport>

    <!-- Detail Modal -->
    <Teleport to="body">
      <div v-if="showDetail" class="mt-overlay" @click.self="closeDetail">
        <div class="mt-modal mt-modal-lg" @click.stop>
          <div class="mt-modal-head">
            <div>
              <h3>{{ detail?.title }}</h3>
              <div class="mt-detail-badges">
                <span class="mt-card-type" :class="'type-'+detail?.meeting_type">{{ typeLabel(detail?.meeting_type) }}</span>
                <span class="mt-card-status" :class="'st-'+detail?.status">{{ statusLabel(detail?.status) }}</span>
              </div>
            </div>
            <div class="mt-detail-actions">
              <button @click="editMeeting" class="mt-btn-ghost" title="Edit">
                <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M11 4H4a2 2 0 0 0-2 2v14a2 2 0 0 0 2 2h14a2 2 0 0 0 2-2v-7"/><path d="M18.5 2.5a2.121 2.121 0 0 1 3 3L12 15l-4 1 1-4 9.5-9.5z"/></svg>
              </button>
              <button @click="deleteMeeting" class="mt-btn-ghost mt-btn-danger" title="Hapus">
                <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><polyline points="3 6 5 6 21 6"/><path d="M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6m3 0V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2"/></svg>
              </button>
              <button @click="closeDetail" class="mt-modal-x">&times;</button>
            </div>
          </div>
          <div class="mt-modal-body mt-detail-body">
            <div class="mt-detail-info">
              <div class="mt-info-item"><strong>Tanggal:</strong> {{ formatDate(detail?.meeting_date) }}</div>
              <div class="mt-info-item"><strong>Durasi:</strong> {{ detail?.duration }} menit</div>
              <div v-if="detail?.location" class="mt-info-item"><strong>Lokasi:</strong> {{ detail?.location }}</div>
              <div v-if="detail?.meeting_link" class="mt-info-item"><strong>Link:</strong> <a :href="detail?.meeting_link" target="_blank">{{ detail?.meeting_link }}</a></div>
            </div>

            <!-- Tabs -->
            <div class="mt-tabs">
              <button v-for="t in tabs" :key="t.key" class="mt-tab" :class="{active:activeTab===t.key}" @click="activeTab=t.key">{{ t.label }}</button>
            </div>

            <!-- Tab: Participants -->
            <div v-if="activeTab==='participants'" class="mt-tab-content">
              <div class="mt-section-head">
                <h4>Peserta Meeting</h4>
                <button @click="showAddParticipant=true" class="mt-btn-sm">+ Tambah</button>
              </div>
              <div v-if="detail?.participants?.length" class="mt-parti-list">
                <div v-for="p in detail.participants" :key="p.id" class="mt-parti-item">
                  <div class="mt-parti-info">
                    <div class="mt-parti-name">{{ p.user ? p.user.name : p.external_name }}</div>
                    <div class="mt-parti-meta">{{ p.role }} · {{ p.user ? 'Internal' : (p.external_email || p.external_phone || 'Eksternal') }}</div>
                  </div>
                  <button @click="removeParticipant(p.id)" class="mt-btn-x">×</button>
                </div>
              </div>
              <p v-else class="mt-empty-text">Belum ada peserta</p>
            </div>

            <!-- Tab: Notes -->
            <div v-if="activeTab==='notes'" class="mt-tab-content">
              <form @submit.prevent="saveNotes" class="mt-notes-form">
                <div class="mt-fg"><label>Tujuan</label><textarea v-model="notesForm.purpose" class="mt-textarea" rows="2"></textarea></div>
                <div class="mt-fg"><label>Ringkasan</label><textarea v-model="notesForm.summary" class="mt-textarea" rows="3"></textarea></div>
                <div class="mt-fg"><label>Poin Diskusi</label><textarea v-model="notesForm.discussion_points" class="mt-textarea" rows="3"></textarea></div>
                <div class="mt-frow">
                  <div class="mt-fg">
                    <label>Keputusan</label>
                    <select v-model="notesForm.decision" class="mt-input">
                      <option value="">Belum ada</option>
                      <option value="lanjut_dp">Lanjut DP</option>
                      <option value="revisi_proposal">Revisi Proposal</option>
                      <option value="survei_ulang">Survei Ulang</option>
                      <option value="tidak_lanjut">Tidak Lanjut</option>
                    </select>
                  </div>
                  <div class="mt-fg"><label>Next Step</label><textarea v-model="notesForm.next_step" class="mt-textarea" rows="2"></textarea></div>
                </div>
                <button type="submit" class="mt-btn-primary" :disabled="savingNotes">{{ savingNotes ? 'Menyimpan...' : 'Simpan Notulensi' }}</button>
              </form>
            </div>

            <!-- Tab: Action Plans -->
            <div v-if="activeTab==='actions'" class="mt-tab-content">
              <div class="mt-section-head">
                <h4>Action Plans</h4>
                <button @click="showAddAction=true" class="mt-btn-sm">+ Tambah</button>
              </div>
              <div v-if="detail?.action_plans?.length" class="mt-action-list">
                <div v-for="a in detail.action_plans" :key="a.id" class="mt-action-item">
                  <div class="mt-action-check">
                    <input type="checkbox" :checked="a.status==='done'" @change="toggleActionStatus(a)" />
                  </div>
                  <div class="mt-action-info" :class="{'done': a.status==='done'}">
                    <div class="mt-action-name">{{ a.task_name }}</div>
                    <div class="mt-action-meta">PIC: {{ a.pic }} · Deadline: {{ a.deadline ? formatShort(a.deadline) : '-' }}</div>
                  </div>
                  <button @click="removeAction(a.id)" class="mt-btn-x">×</button>
                </div>
              </div>
              <p v-else class="mt-empty-text">Belum ada action plan</p>
            </div>

            <!-- Tab: Files -->
            <div v-if="activeTab==='files'" class="mt-tab-content">
              <div class="mt-section-head">
                <h4>Dokumen</h4>
                <label class="mt-btn-sm mt-upload-label">
                  <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4"/><polyline points="17 8 12 3 7 8"/><line x1="12" y1="3" x2="12" y2="15"/></svg>
                  Upload
                  <input type="file" @change="uploadFile" hidden accept=".pdf,.docx,.xlsx,.jpg,.jpeg,.png,.doc,.xls" />
                </label>
              </div>
              <div v-if="detail?.files?.length" class="mt-file-list">
                <div v-for="f in detail.files" :key="f.id" class="mt-file-item">
                  <div class="mt-file-icon">{{ fileIcon(f.file_type) }}</div>
                  <div class="mt-file-info">
                    <a :href="f.file_path" target="_blank" class="mt-file-name">{{ f.file_name }}</a>
                    <div class="mt-file-meta">{{ f.file_type?.toUpperCase() }} · {{ formatShort(f.uploaded_at) }}</div>
                  </div>
                  <button @click="removeFile(f.id)" class="mt-btn-x">×</button>
                </div>
              </div>
              <p v-else class="mt-empty-text">Belum ada dokumen</p>
            </div>
          </div>
        </div>
      </div>
    </Teleport>

    <!-- Add Participant Modal -->
    <Teleport to="body">
      <div v-if="showAddParticipant" class="mt-overlay" @click.self="showAddParticipant=false">
        <div class="mt-modal mt-modal-sm" @click.stop>
          <div class="mt-modal-head"><h3>Tambah Peserta</h3><button @click="showAddParticipant=false" class="mt-modal-x">&times;</button></div>
          <form @submit.prevent="addParticipant" class="mt-modal-body">
            <div class="mt-fg">
              <label>Role <span class="req">*</span></label>
              <select v-model="partForm.role" class="mt-input">
                <option value="sales">Sales</option>
                <option value="brand_manager">Brand Manager</option>
                <option value="legal">Legal</option>
                <option value="finance">Finance</option>
                <option value="mitra">Mitra</option>
                <option value="owner">Owner</option>
              </select>
            </div>
            <div class="mt-fg">
              <label>Pilih User Internal (opsional)</label>
              <select v-model="partForm.user_id" class="mt-input">
                <option value="">-- Eksternal --</option>
                <option v-for="u in users" :key="u.id" :value="u.id">{{ u.name }} ({{ u.role }})</option>
              </select>
            </div>
            <div v-if="!partForm.user_id">
              <div class="mt-fg"><label>Nama</label><input v-model="partForm.external_name" class="mt-input" placeholder="Nama peserta" /></div>
              <div class="mt-frow">
                <div class="mt-fg"><label>Email</label><input v-model="partForm.external_email" class="mt-input" placeholder="email" /></div>
                <div class="mt-fg"><label>Telepon</label><input v-model="partForm.external_phone" class="mt-input" placeholder="08xx" /></div>
              </div>
            </div>
            <div class="mt-modal-foot"><button type="button" @click="showAddParticipant=false" class="mt-btn-sec">Batal</button><button type="submit" class="mt-btn-primary">Tambah</button></div>
          </form>
        </div>
      </div>
    </Teleport>

    <!-- Add Action Plan Modal -->
    <Teleport to="body">
      <div v-if="showAddAction" class="mt-overlay" @click.self="showAddAction=false">
        <div class="mt-modal mt-modal-sm" @click.stop>
          <div class="mt-modal-head"><h3>Tambah Action Plan</h3><button @click="showAddAction=false" class="mt-modal-x">&times;</button></div>
          <form @submit.prevent="addAction" class="mt-modal-body">
            <div class="mt-fg"><label>Task <span class="req">*</span></label><input v-model="actionForm.task_name" class="mt-input" placeholder="Nama tugas" /></div>
            <div class="mt-frow">
              <div class="mt-fg"><label>PIC <span class="req">*</span></label><input v-model="actionForm.pic" class="mt-input" placeholder="Penanggung jawab" /></div>
              <div class="mt-fg"><label>Deadline</label><input v-model="actionForm.deadline" type="date" class="mt-input" /></div>
            </div>
            <div class="mt-modal-foot"><button type="button" @click="showAddAction=false" class="mt-btn-sec">Batal</button><button type="submit" class="mt-btn-primary">Tambah</button></div>
          </form>
        </div>
      </div>
    </Teleport>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted } from 'vue'
import { meetingApi, userApi } from '../../services/api'
import { useToastStore } from '../../stores/toast'

const toast = useToastStore()
const meetings = ref([])
const total = ref(0)
const loading = ref(false)
const search = ref('')
const filterType = ref('')
const filterStatus = ref('')

const showForm = ref(false)
const editMode = ref(false)
const saving = ref(false)
const form = reactive({ title:'', meeting_type:'edukasi', meeting_date:'', duration:60, location:'', meeting_link:'', status:'scheduled' })
let editId = null

const showDetail = ref(false)
const detail = ref(null)
const activeTab = ref('participants')
const tabs = [
  { key:'participants', label:'Peserta' },
  { key:'notes', label:'Notulensi' },
  { key:'actions', label:'Action Plan' },
  { key:'files', label:'Dokumen' },
]

const showAddParticipant = ref(false)
const partForm = reactive({ user_id:'', external_name:'', external_email:'', external_phone:'', role:'mitra' })
const users = ref([])

const showAddAction = ref(false)
const actionForm = reactive({ task_name:'', pic:'', deadline:'' })

const notesForm = reactive({ purpose:'', summary:'', discussion_points:'', decision:'', next_step:'' })
const savingNotes = ref(false)

const typeLabels = { edukasi:'Edukasi', closing:'Closing', review_lokasi:'Review Lokasi', operasional:'Operasional' }
const statusLabels = { scheduled:'Scheduled', ongoing:'Ongoing', completed:'Completed', waiting_decision:'Waiting Decision', approved:'Approved', cancelled:'Cancelled' }
function typeLabel(t) { return typeLabels[t] || t }
function statusLabel(s) { return statusLabels[s] || s }

function formatDate(d) { if(!d) return '-'; return new Date(d).toLocaleDateString('id-ID',{weekday:'short',day:'numeric',month:'short',year:'numeric',hour:'2-digit',minute:'2-digit'}) }
function formatShort(d) { if(!d) return '-'; return new Date(d).toLocaleDateString('id-ID',{day:'numeric',month:'short',year:'numeric'}) }
function fileIcon(t) { const icons={pdf:'📄',docx:'📝',doc:'📝',xlsx:'📊',xls:'📊',jpg:'🖼',jpeg:'🖼',png:'🖼'}; return icons[t]||'📎' }

const filtered = computed(() => {
  let list = meetings.value
  const q = search.value.toLowerCase().trim()
  if (q) list = list.filter(m => m.title.toLowerCase().includes(q))
  return list
})

onMounted(async () => {
  loadMeetings()
  try {
    const { data } = await userApi.list({ limit: 100 })
    users.value = data.data || []
  } catch {}
})

async function loadMeetings() {
  loading.value = true
  try {
    const params = { page: 1, limit: 50 }
    if (filterStatus.value) params.status = filterStatus.value
    if (filterType.value) params.type = filterType.value
    const { data } = await meetingApi.list(params)
    meetings.value = data.data || []
    total.value = data.total || 0
  } catch { toast.error('Gagal memuat meeting') }
  finally { loading.value = false }
}

// Watch filters
import { watch } from 'vue'
watch([filterType, filterStatus], () => loadMeetings())

function openCreate() {
  editMode.value = false; editId = null
  Object.assign(form, { title:'', meeting_type:'edukasi', meeting_date:'', duration:60, location:'', meeting_link:'', status:'scheduled' })
  showForm.value = true
}

function editMeeting() {
  if (!detail.value) return
  editMode.value = true; editId = detail.value.id
  const d = detail.value
  const dt = d.meeting_date ? new Date(d.meeting_date) : null
  const dateStr = dt ? dt.toISOString().slice(0,16) : ''
  Object.assign(form, { title:d.title, meeting_type:d.meeting_type, meeting_date:dateStr, duration:d.duration, location:d.location||'', meeting_link:d.meeting_link||'', status:d.status })
  showDetail.value = false
  showForm.value = true
}

async function saveMeeting() {
  if (!form.title || !form.meeting_date) { toast.error('Judul dan tanggal wajib diisi'); return }
  saving.value = true
  try {
    const payload = { ...form, meeting_date: new Date(form.meeting_date).toISOString() }
    if (editMode.value && editId) {
      await meetingApi.update(editId, payload)
      toast.success('Meeting berhasil diupdate')
    } else {
      await meetingApi.create(payload)
      toast.success('Meeting berhasil dibuat')
    }
    showForm.value = false
    loadMeetings()
  } catch(e) { toast.error(e.response?.data?.error || 'Gagal menyimpan') }
  finally { saving.value = false }
}

async function openDetail(id) {
  try {
    const { data } = await meetingApi.get(id)
    detail.value = data.data
    activeTab.value = 'participants'
    // Prefill notes
    const n = detail.value?.notes
    Object.assign(notesForm, { purpose:n?.purpose||'', summary:n?.summary||'', discussion_points:n?.discussion_points||'', decision:n?.decision||'', next_step:n?.next_step||'' })
    showDetail.value = true
  } catch { toast.error('Gagal memuat detail') }
}
function closeDetail() { showDetail.value = false; detail.value = null }

async function deleteMeeting() {
  if (!detail.value || !confirm('Yakin ingin menghapus meeting ini?')) return
  try {
    await meetingApi.delete(detail.value.id)
    toast.success('Meeting dihapus')
    closeDetail(); loadMeetings()
  } catch { toast.error('Gagal menghapus') }
}

// Participants
async function addParticipant() {
  if (!partForm.role) { toast.error('Role wajib'); return }
  try {
    await meetingApi.addParticipant(detail.value.id, { ...partForm })
    showAddParticipant.value = false
    Object.assign(partForm, { user_id:'', external_name:'', external_email:'', external_phone:'', role:'mitra' })
    await openDetail(detail.value.id)
  } catch(e) { toast.error(e.response?.data?.error || 'Gagal menambah peserta') }
}
async function removeParticipant(pid) {
  try { await meetingApi.deleteParticipant(detail.value.id, pid); await openDetail(detail.value.id) }
  catch { toast.error('Gagal menghapus peserta') }
}

// Notes
async function saveNotes() {
  savingNotes.value = true
  try {
    await meetingApi.saveNotes(detail.value.id, { ...notesForm })
    toast.success('Notulensi disimpan')
    await openDetail(detail.value.id)
  } catch { toast.error('Gagal menyimpan notulensi') }
  finally { savingNotes.value = false }
}

// Action Plans
async function addAction() {
  if (!actionForm.task_name || !actionForm.pic) { toast.error('Task dan PIC wajib'); return }
  try {
    await meetingApi.addActionPlan(detail.value.id, { ...actionForm })
    showAddAction.value = false
    Object.assign(actionForm, { task_name:'', pic:'', deadline:'' })
    await openDetail(detail.value.id)
  } catch { toast.error('Gagal menambah action plan') }
}
async function toggleActionStatus(a) {
  const newStatus = a.status === 'done' ? 'pending' : 'done'
  try {
    await meetingApi.updateActionPlan(detail.value.id, a.id, { task_name: a.task_name, pic: a.pic, status: newStatus })
    await openDetail(detail.value.id)
  } catch { toast.error('Gagal update status') }
}
async function removeAction(aid) {
  try { await meetingApi.deleteActionPlan(detail.value.id, aid); await openDetail(detail.value.id) }
  catch { toast.error('Gagal menghapus') }
}

// Files
async function uploadFile(e) {
  const file = e.target.files[0]
  if (!file) return
  try {
    await meetingApi.upload(detail.value.id, file)
    toast.success('File berhasil diupload')
    await openDetail(detail.value.id)
  } catch { toast.error('Gagal upload file') }
  e.target.value = ''
}
async function removeFile(fid) {
  try { await meetingApi.deleteFile(detail.value.id, fid); await openDetail(detail.value.id) }
  catch { toast.error('Gagal menghapus file') }
}
</script>

<style scoped>
/* ═══ HERO ═══ */
.mt-hero { background: linear-gradient(135deg, #0f0c29 0%, #302b63 50%, #24243e 100%); border-radius: 16px; padding: 32px 36px 24px; margin-bottom: 24px; }
.mt-hero-top { display:flex; justify-content:space-between; align-items:flex-start; gap:16px; flex-wrap:wrap; margin-bottom:20px; }
.mt-hero-title { font-size:1.6rem; font-weight:800; color:#fff; margin:0 0 4px; }
.mt-hero-sub { font-size:.85rem; color:rgba(255,255,255,.5); margin:0; }
.mt-stats { display:flex; gap:28px; flex-wrap:wrap; padding-top:16px; border-top:1px solid rgba(255,255,255,.08); }
.mt-stat { display:flex; align-items:center; gap:8px; }
.mt-stat-dot { width:8px; height:8px; border-radius:50%; }
.mt-stat-label { font-size:.72rem; color:rgba(255,255,255,.4); text-transform:uppercase; letter-spacing:.05em; }
.mt-stat-val { font-size:.9rem; font-weight:800; color:#fff; }

/* ═══ BUTTONS ═══ */
.mt-btn-primary { display:inline-flex; align-items:center; gap:8px; padding:11px 24px; font-size:.85rem; font-weight:700; border-radius:12px; border:none; cursor:pointer; background:linear-gradient(135deg,#6366f1,#8b5cf6); color:#fff; white-space:nowrap; }
.mt-btn-primary:disabled { opacity:.6; cursor:not-allowed; }
.mt-btn-sec { padding:11px 22px; border-radius:12px; font-size:.85rem; font-weight:600; background:#f1f5f9; color:#475569; border:none; cursor:pointer; }
.mt-btn-sm { display:inline-flex; align-items:center; gap:6px; padding:7px 16px; font-size:.78rem; font-weight:600; border-radius:8px; border:none; cursor:pointer; background:#eef2ff; color:#6366f1; }
.mt-btn-ghost { width:34px; height:34px; border-radius:10px; display:flex; align-items:center; justify-content:center; border:1px solid #e8ecf1; background:#fff; cursor:pointer; color:#94a3b8; }
.mt-btn-ghost:hover { border-color:#6366f1; color:#6366f1; background:#eef2ff; }
.mt-btn-danger:hover { border-color:#ef4444!important; color:#ef4444!important; background:#fef2f2!important; }
.mt-btn-x { width:28px; height:28px; border-radius:8px; display:flex; align-items:center; justify-content:center; border:none; background:transparent; cursor:pointer; color:#94a3b8; font-size:1.2rem; }
.mt-btn-x:hover { background:#fef2f2; color:#ef4444; }

/* ═══ TOOLBAR ═══ */
.mt-toolbar { display:flex; justify-content:space-between; align-items:center; gap:12px; flex-wrap:wrap; margin-bottom:20px; }
.mt-search-wrap { position:relative; min-width:260px; flex:1; max-width:380px; }
.mt-search-ico { position:absolute; left:14px; top:50%; transform:translateY(-50%); }
.mt-search { width:100%; padding:10px 14px 10px 40px; border:1.5px solid #e2e8f0; border-radius:12px; font-size:.85rem; background:#fff; color:#1e293b; outline:none; box-sizing:border-box; }
.mt-search:focus { border-color:#6366f1; }
.mt-filters { display:flex; gap:10px; }
.mt-select { padding:10px 14px; border:1.5px solid #e2e8f0; border-radius:12px; font-size:.82rem; background:#fff; color:#1e293b; outline:none; cursor:pointer; }

/* ═══ LIST ═══ */
.mt-list { display:flex; flex-direction:column; gap:12px; }
.mt-card { background:#fff; border-radius:14px; border:1px solid #e8ecf1; padding:20px 24px; cursor:pointer; transition:border-color .2s; }
.mt-card:hover { border-color:#c7d2fe; }
.mt-card-head { display:flex; justify-content:space-between; align-items:center; margin-bottom:10px; }
.mt-card-type { font-size:.68rem; font-weight:700; text-transform:uppercase; letter-spacing:.04em; padding:4px 10px; border-radius:6px; }
.type-edukasi { background:#e0f2fe; color:#0284c7; }
.type-closing { background:#dcfce7; color:#16a34a; }
.type-review_lokasi { background:#fef3c7; color:#d97706; }
.type-operasional { background:#ede9fe; color:#7c3aed; }
.mt-card-status { font-size:.68rem; font-weight:600; padding:4px 10px; border-radius:6px; }
.st-scheduled { background:#e0f2fe; color:#0284c7; }
.st-ongoing { background:#fef3c7; color:#d97706; }
.st-completed { background:#dcfce7; color:#16a34a; }
.st-waiting_decision { background:#fef9c3; color:#a16207; }
.st-approved { background:#d1fae5; color:#059669; }
.st-cancelled { background:#fee2e2; color:#dc2626; }
.mt-card-title { font-size:1.05rem; font-weight:700; color:#0f172a; margin:0 0 10px; }
.mt-card-meta { display:flex; flex-wrap:wrap; gap:16px; }
.mt-meta-item { display:flex; align-items:center; gap:6px; font-size:.78rem; color:#64748b; }

/* ═══ SKELETON ═══ */
.mt-card-skel { display:flex; flex-direction:column; gap:12px; min-height:100px; }
.skel-bar { border-radius:6px; }
.shimmer { background:linear-gradient(90deg,#e8ecf1 25%,#f1f5f9 50%,#e8ecf1 75%); background-size:200% 100%; animation:shimmer 1.5s infinite; }
@keyframes shimmer { 0%{background-position:200% 0} 100%{background-position:-200% 0} }

/* ═══ EMPTY ═══ */
.mt-empty { display:flex; flex-direction:column; align-items:center; padding:80px 32px; background:#fff; border-radius:16px; border:2px dashed #e2e8f0; text-align:center; }
.mt-empty-ico { width:80px; height:80px; border-radius:50%; background:linear-gradient(135deg,#f1f5f9,#e2e8f0); display:flex; align-items:center; justify-content:center; margin-bottom:20px; }
.mt-empty h3 { font-weight:700; color:#0f172a; margin:0; }
.mt-empty p { color:#94a3b8; font-size:.85rem; margin:6px 0 0; }
.mt-empty-text { color:#94a3b8; font-size:.85rem; text-align:center; padding:24px; }

/* ═══ MODAL ═══ */
.mt-overlay { position:fixed; inset:0; background:rgba(0,0,0,.5); display:flex; align-items:center; justify-content:center; z-index:1000; backdrop-filter:blur(4px); }
.mt-modal { background:#fff; border-radius:18px; width:100%; max-width:580px; box-shadow:0 24px 80px rgba(0,0,0,.2); max-height:90vh; display:flex; flex-direction:column; }
.mt-modal-lg { max-width:780px; }
.mt-modal-sm { max-width:480px; }
.mt-modal-head { display:flex; justify-content:space-between; align-items:center; padding:22px 28px; border-bottom:1px solid #f1f5f9; }
.mt-modal-head h3 { font-size:1.1rem; font-weight:700; margin:0; color:#0f172a; }
.mt-modal-x { width:34px; height:34px; border-radius:10px; display:flex; align-items:center; justify-content:center; border:none; background:transparent; font-size:1.4rem; color:#94a3b8; cursor:pointer; }
.mt-modal-x:hover { background:#f1f5f9; color:#0f172a; }
.mt-modal-body { padding:24px 28px; overflow-y:auto; flex:1; }
.mt-modal-foot { display:flex; justify-content:flex-end; gap:10px; padding-top:16px; }

.mt-fg { margin-bottom:16px; }
.mt-fg label { display:block; font-size:.82rem; font-weight:600; margin-bottom:6px; color:#334155; }
.mt-frow { display:grid; grid-template-columns:1fr 1fr; gap:14px; }
.req { color:#ef4444; }
.mt-input { width:100%; padding:10px 14px; border:1.5px solid #e2e8f0; border-radius:10px; font-size:.85rem; background:#fafbfc; color:#1e293b; outline:none; box-sizing:border-box; font-family:inherit; }
.mt-input:focus { border-color:#6366f1; background:#fff; }
.mt-textarea { width:100%; padding:10px 14px; border:1.5px solid #e2e8f0; border-radius:10px; font-size:.85rem; background:#fafbfc; color:#1e293b; outline:none; box-sizing:border-box; font-family:inherit; resize:vertical; }
.mt-textarea:focus { border-color:#6366f1; }

/* ═══ DETAIL ═══ */
.mt-detail-badges { display:flex; gap:8px; margin-top:6px; }
.mt-detail-actions { display:flex; gap:8px; align-items:center; }
.mt-detail-body { padding-top:0!important; }
.mt-detail-info { display:flex; flex-wrap:wrap; gap:12px 28px; padding:16px 0; border-bottom:1px solid #f1f5f9; margin-bottom:16px; }
.mt-info-item { font-size:.82rem; color:#64748b; }
.mt-info-item strong { color:#334155; }
.mt-info-item a { color:#6366f1; text-decoration:none; }

/* ═══ TABS ═══ */
.mt-tabs { display:flex; gap:4px; border-bottom:2px solid #f1f5f9; margin-bottom:20px; }
.mt-tab { padding:10px 18px; font-size:.82rem; font-weight:600; color:#94a3b8; background:none; border:none; border-bottom:2px solid transparent; margin-bottom:-2px; cursor:pointer; }
.mt-tab.active { color:#6366f1; border-bottom-color:#6366f1; }

.mt-section-head { display:flex; justify-content:space-between; align-items:center; margin-bottom:14px; }
.mt-section-head h4 { font-size:.92rem; font-weight:700; margin:0; color:#0f172a; }

/* ═══ PARTICIPANTS ═══ */
.mt-parti-list { display:flex; flex-direction:column; gap:8px; }
.mt-parti-item { display:flex; align-items:center; justify-content:space-between; padding:12px 16px; background:#f8fafc; border-radius:10px; }
.mt-parti-name { font-size:.88rem; font-weight:600; color:#0f172a; }
.mt-parti-meta { font-size:.72rem; color:#94a3b8; margin-top:2px; }

/* ═══ ACTION ═══ */
.mt-action-list { display:flex; flex-direction:column; gap:8px; }
.mt-action-item { display:flex; align-items:center; gap:12px; padding:12px 16px; background:#f8fafc; border-radius:10px; }
.mt-action-check input { width:18px; height:18px; accent-color:#6366f1; cursor:pointer; }
.mt-action-info { flex:1; }
.mt-action-info.done .mt-action-name { text-decoration:line-through; color:#94a3b8; }
.mt-action-name { font-size:.88rem; font-weight:600; color:#0f172a; }
.mt-action-meta { font-size:.72rem; color:#94a3b8; margin-top:2px; }

/* ═══ FILES ═══ */
.mt-file-list { display:flex; flex-direction:column; gap:8px; }
.mt-file-item { display:flex; align-items:center; gap:12px; padding:12px 16px; background:#f8fafc; border-radius:10px; }
.mt-file-icon { font-size:1.5rem; }
.mt-file-info { flex:1; }
.mt-file-name { font-size:.85rem; font-weight:600; color:#6366f1; text-decoration:none; }
.mt-file-name:hover { text-decoration:underline; }
.mt-file-meta { font-size:.72rem; color:#94a3b8; margin-top:2px; }
.mt-upload-label { cursor:pointer; }

.mt-notes-form { max-width:100%; }

@media (max-width:768px) {
  .mt-hero { padding:24px 20px 18px; }
  .mt-toolbar { flex-direction:column; }
  .mt-search-wrap { max-width:100%; }
  .mt-filters { width:100%; }
  .mt-frow { grid-template-columns:1fr; }
  .mt-modal { max-width:95vw; }
}
</style>
