<template>
  <div class="ld-page">
    <!-- Hero -->
    <div class="ld-hero">
      <div class="ld-hero-top">
        <div>
          <h1 class="ld-hero-title">Lead Management</h1>
          <p class="ld-hero-sub">Kelola calon mitra & pipeline kemitraan</p>
        </div>
        <div class="ld-hero-actions">
          <button @click="viewMode='kanban'" class="ld-view-btn" :class="{active: viewMode==='kanban'}">
            <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><rect x="3" y="3" width="7" height="7"/><rect x="14" y="3" width="7" height="7"/><rect x="3" y="14" width="7" height="7"/><rect x="14" y="14" width="7" height="7"/></svg>
            Kanban
          </button>
          <button @click="viewMode='list'" class="ld-view-btn" :class="{active: viewMode==='list'}">
            <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><line x1="8" y1="6" x2="21" y2="6"/><line x1="8" y1="12" x2="21" y2="12"/><line x1="8" y1="18" x2="21" y2="18"/><line x1="3" y1="6" x2="3.01" y2="6"/><line x1="3" y1="12" x2="3.01" y2="12"/><line x1="3" y1="18" x2="3.01" y2="18"/></svg>
            List
          </button>
          <button @click="openCreate" class="ld-btn-primary">
            <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><line x1="12" y1="5" x2="12" y2="19"/><line x1="5" y1="12" x2="19" y2="12"/></svg>
            Tambah Lead
          </button>
        </div>
      </div>
      <div class="ld-stats">
        <div class="ld-stat"><span class="ld-stat-dot" style="background:#818cf8"></span><span class="ld-stat-label">Total</span><span class="ld-stat-val">{{ allLeads.length }}</span></div>
        <div class="ld-stat"><span class="ld-stat-dot" style="background:#22c55e"></span><span class="ld-stat-label">New</span><span class="ld-stat-val">{{ (kanban['NEW']||[]).length }}</span></div>
        <div class="ld-stat"><span class="ld-stat-dot" style="background:#f59e0b"></span><span class="ld-stat-label">In Progress</span><span class="ld-stat-val">{{ allLeads.filter(l=>l.status!=='NEW'&&l.status!=='COMPLETED').length }}</span></div>
        <div class="ld-stat"><span class="ld-stat-dot" style="background:#06b6d4"></span><span class="ld-stat-label">Completed</span><span class="ld-stat-val">{{ (kanban['COMPLETED']||[]).length }}</span></div>
      </div>
    </div>

    <!-- Search Toolbar -->
    <div class="ld-toolbar">
      <div class="ld-search-wrap">
        <svg class="ld-search-ico" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="#94a3b8" stroke-width="2"><circle cx="11" cy="11" r="8"/><line x1="21" y1="21" x2="16.65" y2="16.65"/></svg>
        <input v-model="search" class="ld-search" placeholder="Cari nama, email, telepon..." />
      </div>
    </div>

    <!-- Kanban View -->
    <div v-if="viewMode==='kanban'" class="kanban-board">
      <div v-for="status in visibleStatuses" :key="status" class="kanban-column" @dragover.prevent @drop="onDrop($event, status)">
        <div class="kanban-col-header">
          <span class="kanban-col-title">{{ formatStatus(status) }}</span>
          <span class="kanban-col-count">{{ filteredKanban(status).length }}</span>
        </div>
        <div class="kanban-col-body">
          <div v-for="lead in filteredKanban(status)" :key="lead.id" class="kanban-card" draggable="true" @dragstart="onDragStart($event, lead)">
            <div class="kc-name">{{ lead.full_name }}</div>
            <div class="kc-email">{{ lead.email }}</div>
            <div class="kc-progress-wrap">
              <div class="kc-progress"><div class="kc-progress-fill" :style="{width: lead.progress_percentage+'%'}"></div></div>
              <span class="kc-progress-text">{{ lead.progress_percentage }}%</span>
            </div>
            <div class="kc-footer">
              <span v-if="lead.phone" class="kc-phone">{{ lead.phone }}</span>
              <router-link :to="`/leads/${lead.id}`" class="kc-detail">Detail →</router-link>
            </div>
          </div>
          <div v-if="!filteredKanban(status).length" class="kanban-empty">Kosong</div>
        </div>
      </div>
    </div>

    <!-- List View -->
    <div v-else class="ld-table-wrap">
      <table class="ld-table">
        <thead>
          <tr><th>Nama</th><th>Email</th><th>Telepon</th><th>Status</th><th>Progress</th><th>Aksi</th></tr>
        </thead>
        <tbody>
          <tr v-for="lead in filteredAll" :key="lead.id">
            <td class="ld-name-cell">{{ lead.full_name }}</td>
            <td>{{ lead.email }}</td>
            <td>{{ lead.phone || '-' }}</td>
            <td><span class="ld-badge" :class="'st-'+lead.status">{{ formatStatus(lead.status) }}</span></td>
            <td>
              <div class="ld-progress-cell">
                <div class="ld-progress-bar"><div class="ld-progress-fill" :style="{width:lead.progress_percentage+'%'}"></div></div>
                <span>{{ lead.progress_percentage }}%</span>
              </div>
            </td>
            <td>
              <div class="ld-actions">
                <router-link :to="`/leads/${lead.id}`" class="ld-detail-btn">Detail</router-link>
                <button @click="confirmDelete(lead)" class="ld-del-btn" title="Hapus">
                  <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><polyline points="3 6 5 6 21 6"/><path d="M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6m3 0V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2"/></svg>
                </button>
              </div>
            </td>
          </tr>
          <tr v-if="!filteredAll.length"><td colspan="6" class="ld-empty-row">Belum ada lead</td></tr>
        </tbody>
      </table>
    </div>

    <!-- Create Lead Modal -->
    <Teleport to="body">
      <div v-if="showCreate" class="ld-overlay" @click.self="showCreate=false">
        <div class="ld-modal" @click.stop>
          <div class="ld-modal-head">
            <div class="ld-modal-title-group">
              <div class="ld-modal-icon">
                <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M20 21v-2a4 4 0 0 0-4-4H8a4 4 0 0 0-4 4v2"/><circle cx="12" cy="7" r="4"/></svg>
              </div>
              <h3>Tambah Lead Baru</h3>
            </div>
            <button @click="showCreate=false" class="ld-modal-x">&times;</button>
          </div>
          <form @submit.prevent="createLead" class="ld-modal-body">
            <div class="ld-fg">
              <label>Pilih Outlet / Brand <span class="req">*</span></label>
              <select v-model="newLead.brand_id" class="ld-input" required>
                <option value="">— Pilih Outlet —</option>
                <option v-for="o in outlets" :key="o.id" :value="o.id">{{ o.name }}</option>
              </select>
              <div class="ld-field-hint">Outlet franchise yang diminati calon mitra</div>
            </div>
            <div class="ld-fg">
              <label>Nama Lengkap <span class="req">*</span></label>
              <input v-model="newLead.full_name" class="ld-input" placeholder="Nama lengkap calon mitra" required />
            </div>
            <div class="ld-frow">
              <div class="ld-fg">
                <label>Email <span class="req">*</span></label>
                <input v-model="newLead.email" type="email" class="ld-input" placeholder="email@example.com" required />
              </div>
              <div class="ld-fg">
                <label>Telepon</label>
                <input v-model="newLead.phone" class="ld-input" placeholder="08xxxxxxxxxx" />
              </div>
            </div>
            <div class="ld-fg">
              <label>Catatan</label>
              <textarea v-model="newLead.notes" class="ld-textarea" rows="3" placeholder="Catatan tambahan (opsional)"></textarea>
            </div>
            <div class="ld-modal-foot">
              <button type="button" @click="showCreate=false" class="ld-btn-sec">Batal</button>
              <button type="submit" class="ld-btn-primary" :disabled="creating">{{ creating ? 'Menyimpan...' : 'Simpan' }}</button>
            </div>
          </form>
        </div>
      </div>
    </Teleport>

    <!-- Delete Confirm -->
    <Teleport to="body">
      <div v-if="deleteTarget" class="ld-overlay" @click.self="deleteTarget=null">
        <div class="ld-modal ld-modal-sm" @click.stop>
          <div class="ld-delete-body">
            <div class="ld-delete-icon"><svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="#ef4444" stroke-width="2"><polyline points="3 6 5 6 21 6"/><path d="M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6m3 0V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2"/></svg></div>
            <h3>Hapus Lead</h3>
            <p>Yakin ingin menghapus <strong>{{ deleteTarget.full_name }}</strong>?</p>
          </div>
          <div class="ld-delete-foot">
            <button @click="deleteTarget=null" class="ld-btn-sec">Batal</button>
            <button @click="doDelete" class="ld-btn-danger" :disabled="deleting">{{ deleting ? 'Menghapus...' : 'Ya, Hapus' }}</button>
          </div>
        </div>
      </div>
    </Teleport>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted } from 'vue'
import { leadApi, outletApi } from '../../services/api'
import { useToastStore } from '../../stores/toast'

const toast = useToastStore()
const kanban = ref({})
const allLeads = ref([])
const outlets = ref([])
const search = ref('')
const viewMode = ref('kanban')
const showCreate = ref(false)
const creating = ref(false)
const deleteTarget = ref(null)
const deleting = ref(false)
const newLead = reactive({ brand_id: '', full_name: '', email: '', phone: '', notes: '' })
let dragged = null

const visibleStatuses = ['NEW','CONSULTATION','LOCATION_SUBMITTED','SURVEY_APPROVED','MEETING_DONE','READY_FOR_DP','DP_PAID','AGREEMENT_REVIEW','FULLY_PAID','ACTIVE_PARTNERSHIP','RUNNING','COMPLETED']

function formatStatus(s) { return (s||'').replace(/_/g, ' ') }

function filteredKanban(status) {
  const list = kanban.value[status] || []
  const q = search.value.toLowerCase().trim()
  if (!q) return list
  return list.filter(l => l.full_name.toLowerCase().includes(q) || l.email.toLowerCase().includes(q) || (l.phone||'').includes(q))
}

const filteredAll = computed(() => {
  const q = search.value.toLowerCase().trim()
  if (!q) return allLeads.value
  return allLeads.value.filter(l => l.full_name.toLowerCase().includes(q) || l.email.toLowerCase().includes(q) || (l.phone||'').includes(q))
})

onMounted(async () => {
  load()
  try { const { data } = await outletApi.list({ limit: 100 }); outlets.value = data.data || [] } catch {}
})

async function load() {
  try {
    const { data } = await leadApi.kanban()
    kanban.value = data.data || {}
    // Flatten for list view
    const flat = []
    for (const st of visibleStatuses) { if (kanban.value[st]) flat.push(...kanban.value[st]) }
    allLeads.value = flat
  } catch { toast.error('Gagal memuat data lead') }
}

function onDragStart(e, lead) { dragged = lead; e.dataTransfer.effectAllowed = 'move' }
async function onDrop(e, status) {
  if (!dragged || dragged.status === status) return
  try { await leadApi.updateStatus(dragged.id, { status }); toast.success('Status diupdate'); await load() }
  catch { toast.error('Gagal update status') }
  dragged = null
}

function openCreate() {
  Object.assign(newLead, { brand_id: '', full_name: '', email: '', phone: '', notes: '' })
  showCreate.value = true
}

async function createLead() {
  if (!newLead.brand_id) { toast.error('Pilih outlet/brand terlebih dahulu'); return }
  creating.value = true
  try {
    await leadApi.create(newLead)
    toast.success('Lead berhasil ditambah')
    showCreate.value = false
    await load()
  } catch (e) { toast.error(e.response?.data?.error || 'Gagal menambah lead') }
  finally { creating.value = false }
}

function confirmDelete(lead) { deleteTarget.value = lead }
async function doDelete() {
  deleting.value = true
  try {
    await leadApi.delete(deleteTarget.value.id)
    toast.success('Lead berhasil dihapus')
    deleteTarget.value = null
    await load()
  } catch { toast.error('Gagal menghapus') }
  finally { deleting.value = false }
}
</script>

<style scoped>
/* ═══ HERO ═══ */
.ld-hero { background: linear-gradient(135deg, #0f0c29 0%, #302b63 50%, #24243e 100%); border-radius: 16px; padding: 32px 36px 24px; margin-bottom: 24px; }
.ld-hero-top { display: flex; justify-content: space-between; align-items: flex-start; gap: 16px; flex-wrap: wrap; margin-bottom: 20px; }
.ld-hero-title { font-size: 1.6rem; font-weight: 800; color: #fff; margin: 0 0 4px; }
.ld-hero-sub { font-size: .85rem; color: rgba(255,255,255,.5); margin: 0; }
.ld-hero-actions { display: flex; gap: 8px; align-items: center; flex-wrap: wrap; }
.ld-view-btn { display: inline-flex; align-items: center; gap: 6px; padding: 9px 16px; font-size: .78rem; font-weight: 600; border-radius: 10px; border: 1px solid rgba(255,255,255,.15); background: transparent; color: rgba(255,255,255,.5); cursor: pointer; }
.ld-view-btn.active { background: rgba(255,255,255,.1); color: #fff; border-color: rgba(255,255,255,.3); }
.ld-stats { display: flex; gap: 28px; flex-wrap: wrap; padding-top: 16px; border-top: 1px solid rgba(255,255,255,.08); }
.ld-stat { display: flex; align-items: center; gap: 8px; }
.ld-stat-dot { width: 8px; height: 8px; border-radius: 50%; }
.ld-stat-label { font-size: .72rem; color: rgba(255,255,255,.4); text-transform: uppercase; letter-spacing: .05em; }
.ld-stat-val { font-size: .9rem; font-weight: 800; color: #fff; }

/* ═══ BUTTONS ═══ */
.ld-btn-primary { display: inline-flex; align-items: center; gap: 8px; padding: 11px 24px; font-size: .85rem; font-weight: 700; border-radius: 12px; border: none; cursor: pointer; background: linear-gradient(135deg, #6366f1, #8b5cf6); color: #fff; white-space: nowrap; }
.ld-btn-primary:disabled { opacity: .6; cursor: not-allowed; }
.ld-btn-sec { padding: 11px 22px; border-radius: 12px; font-size: .85rem; font-weight: 600; background: #f1f5f9; color: #475569; border: none; cursor: pointer; }
.ld-btn-danger { padding: 11px 22px; border-radius: 12px; font-size: .85rem; font-weight: 600; background: linear-gradient(135deg, #ef4444, #dc2626); color: #fff; border: none; cursor: pointer; }
.ld-btn-danger:disabled { opacity: .6; }

/* ═══ TOOLBAR ═══ */
.ld-toolbar { margin-bottom: 20px; }
.ld-search-wrap { position: relative; max-width: 400px; }
.ld-search-ico { position: absolute; left: 14px; top: 50%; transform: translateY(-50%); }
.ld-search { width: 100%; padding: 10px 14px 10px 40px; border: 1.5px solid #e2e8f0; border-radius: 12px; font-size: .85rem; background: #fff; color: #1e293b; outline: none; box-sizing: border-box; }
.ld-search:focus { border-color: #6366f1; }

/* ═══ KANBAN ═══ */
.kanban-board { display: flex; gap: 14px; overflow-x: auto; padding-bottom: 16px; }
.kanban-column { min-width: 240px; max-width: 280px; flex-shrink: 0; background: #f8fafc; border-radius: 14px; border: 1px solid #e8ecf1; display: flex; flex-direction: column; }
.kanban-col-header { display: flex; justify-content: space-between; align-items: center; padding: 14px 16px; border-bottom: 1px solid #e8ecf1; }
.kanban-col-title { font-size: .72rem; font-weight: 700; color: #475569; text-transform: uppercase; letter-spacing: .04em; }
.kanban-col-count { font-size: .68rem; font-weight: 700; background: #e0e7ff; color: #6366f1; padding: 2px 8px; border-radius: 10px; }
.kanban-col-body { padding: 12px; flex: 1; display: flex; flex-direction: column; gap: 10px; min-height: 100px; }
.kanban-card { background: #fff; border-radius: 12px; border: 1px solid #e8ecf1; padding: 14px 16px; cursor: grab; transition: border-color .2s; }
.kanban-card:hover { border-color: #c7d2fe; }
.kanban-card:active { cursor: grabbing; }
.kc-name { font-size: .88rem; font-weight: 700; color: #0f172a; margin-bottom: 2px; }
.kc-email { font-size: .72rem; color: #94a3b8; margin-bottom: 8px; overflow: hidden; text-overflow: ellipsis; white-space: nowrap; }
.kc-progress-wrap { display: flex; align-items: center; gap: 8px; margin-bottom: 8px; }
.kc-progress { flex: 1; height: 5px; background: #e2e8f0; border-radius: 3px; overflow: hidden; }
.kc-progress-fill { height: 100%; background: linear-gradient(90deg, #6366f1, #8b5cf6); border-radius: 3px; }
.kc-progress-text { font-size: .68rem; font-weight: 700; color: #475569; }
.kc-footer { display: flex; justify-content: space-between; align-items: center; }
.kc-phone { font-size: .68rem; color: #94a3b8; }
.kc-detail { font-size: .72rem; font-weight: 600; color: #6366f1; text-decoration: none; }
.kanban-empty { text-align: center; padding: 24px; font-size: .78rem; color: #94a3b8; }

/* ═══ TABLE ═══ */
.ld-table-wrap { background: #fff; border-radius: 16px; border: 1px solid #e8ecf1; overflow: hidden; }
.ld-table { width: 100%; border-collapse: collapse; }
.ld-table thead { background: #f8fafc; }
.ld-table th { padding: 14px 20px; font-size: .72rem; font-weight: 700; color: #64748b; text-transform: uppercase; letter-spacing: .05em; text-align: left; border-bottom: 1px solid #e8ecf1; }
.ld-table td { padding: 14px 20px; font-size: .85rem; color: #1e293b; border-bottom: 1px solid #f1f5f9; }
.ld-table tr:last-child td { border-bottom: none; }
.ld-table tbody tr { transition: background .15s; }
.ld-table tbody tr:hover { background: #fafbfc; }
.ld-name-cell { font-weight: 700; }
.ld-badge { font-size: .66rem; font-weight: 700; padding: 3px 8px; border-radius: 6px; text-transform: uppercase; }
.st-NEW { background: #e0f2fe; color: #0284c7; }
.st-CONSULTATION { background: #fef3c7; color: #d97706; }
.st-LOCATION_SUBMITTED, .st-SURVEY_APPROVED { background: #ede9fe; color: #7c3aed; }
.st-MEETING_DONE, .st-READY_FOR_DP { background: #fce7f3; color: #db2777; }
.st-DP_PAID, .st-AGREEMENT_REVIEW, .st-FULLY_PAID { background: #dcfce7; color: #16a34a; }
.st-ACTIVE_PARTNERSHIP, .st-RUNNING { background: #d1fae5; color: #059669; }
.st-COMPLETED { background: #cffafe; color: #0891b2; }
.ld-progress-cell { display: flex; align-items: center; gap: 8px; }
.ld-progress-bar { width: 60px; height: 5px; background: #e2e8f0; border-radius: 3px; overflow: hidden; }
.ld-progress-fill { height: 100%; background: linear-gradient(90deg, #6366f1, #8b5cf6); border-radius: 3px; }
.ld-progress-cell span { font-size: .78rem; font-weight: 700; color: #475569; }
.ld-actions { display: flex; gap: 6px; }
.ld-detail-btn { font-size: .78rem; font-weight: 600; color: #6366f1; text-decoration: none; padding: 5px 12px; border-radius: 8px; border: 1px solid #e0e7ff; background: #eef2ff; }
.ld-detail-btn:hover { background: #e0e7ff; }
.ld-del-btn { width: 30px; height: 30px; border-radius: 8px; display: flex; align-items: center; justify-content: center; border: 1px solid #e8ecf1; background: #fff; cursor: pointer; color: #94a3b8; }
.ld-del-btn:hover { background: #fef2f2; border-color: #ef4444; color: #ef4444; }
.ld-empty-row { text-align: center; padding: 48px 20px !important; color: #94a3b8; }

/* ═══ MODAL ═══ */
.ld-overlay { position: fixed; inset: 0; background: rgba(0,0,0,.5); display: flex; align-items: center; justify-content: center; z-index: 1000; backdrop-filter: blur(4px); }
.ld-modal { background: #fff; border-radius: 18px; width: 100%; max-width: 560px; box-shadow: 0 24px 80px rgba(0,0,0,.2); }
.ld-modal-sm { max-width: 400px; }
.ld-modal-head { display: flex; align-items: center; justify-content: space-between; padding: 22px 28px; border-bottom: 1px solid #f1f5f9; }
.ld-modal-title-group { display: flex; align-items: center; gap: 12px; }
.ld-modal-icon { width: 40px; height: 40px; border-radius: 12px; background: linear-gradient(135deg, #eef2ff, #e0e7ff); display: flex; align-items: center; justify-content: center; color: #6366f1; }
.ld-modal-head h3 { font-size: 1.1rem; font-weight: 700; margin: 0; color: #0f172a; }
.ld-modal-x { width: 34px; height: 34px; border-radius: 10px; display: flex; align-items: center; justify-content: center; border: none; background: transparent; font-size: 1.4rem; color: #94a3b8; cursor: pointer; }
.ld-modal-x:hover { background: #f1f5f9; color: #0f172a; }
.ld-modal-body { padding: 24px 28px; }
.ld-modal-foot { display: flex; justify-content: flex-end; gap: 10px; padding-top: 16px; }

.ld-fg { margin-bottom: 16px; }
.ld-fg label { display: block; font-size: .82rem; font-weight: 600; margin-bottom: 6px; color: #334155; }
.ld-frow { display: grid; grid-template-columns: 1fr 1fr; gap: 14px; }
.req { color: #ef4444; }
.ld-input { width: 100%; padding: 11px 14px; border: 1.5px solid #e2e8f0; border-radius: 12px; font-size: .85rem; background: #fafbfc; color: #1e293b; outline: none; box-sizing: border-box; font-family: inherit; }
.ld-input:focus { border-color: #6366f1; background: #fff; }
select.ld-input { appearance: none; -webkit-appearance: none; background-image: url("data:image/svg+xml,%3csvg xmlns='http://www.w3.org/2000/svg' fill='none' viewBox='0 0 20 20'%3e%3cpath stroke='%236b7280' stroke-linecap='round' stroke-linejoin='round' stroke-width='1.5' d='M6 8l4 4 4-4'/%3e%3c/svg%3e"); background-position: right 12px center; background-repeat: no-repeat; background-size: 16px; padding-right: 36px; cursor: pointer; }
.ld-textarea { width: 100%; padding: 11px 14px; border: 1.5px solid #e2e8f0; border-radius: 12px; font-size: .85rem; background: #fafbfc; color: #1e293b; outline: none; box-sizing: border-box; font-family: inherit; resize: vertical; }
.ld-textarea:focus { border-color: #6366f1; }
.ld-field-hint { font-size: .72rem; color: #94a3b8; margin-top: 4px; }

/* DELETE */
.ld-delete-body { text-align: center; padding: 32px 28px 20px; }
.ld-delete-icon { width: 60px; height: 60px; border-radius: 50%; background: rgba(239,68,68,.08); display: flex; align-items: center; justify-content: center; margin: 0 auto 16px; }
.ld-delete-body h3 { font-size: 1.1rem; font-weight: 700; margin: 0 0 8px; color: #0f172a; }
.ld-delete-body p { color: #64748b; font-size: .85rem; margin: 0; line-height: 1.6; }
.ld-delete-foot { display: flex; justify-content: center; gap: 12px; padding: 0 28px 28px; }

@media (max-width: 768px) {
  .ld-hero { padding: 24px 20px 18px; }
  .ld-hero-top { flex-direction: column; align-items: flex-start; }
  .ld-frow { grid-template-columns: 1fr; }
  .kanban-board { padding-bottom: 8px; }
}
</style>
