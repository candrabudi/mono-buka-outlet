<template>
  <div class="st-page">
    <!-- Hero -->
    <div class="st-hero">
      <div class="st-hero-top">
        <div>
          <h1 class="st-hero-title">Pengaturan Sistem</h1>
          <p class="st-hero-sub">Kelola konfigurasi aplikasi dan integrasi pembayaran</p>
        </div>
      </div>
    </div>

    <!-- Tabs -->
    <div class="st-tab-bar">
      <button v-for="g in groups" :key="g.key" class="st-tab" :class="{ active: activeGroup === g.key }" @click="activeGroup = g.key">
        <span class="st-tab-icon" v-html="g.icon"></span>
        <span>{{ g.label }}</span>
        <span class="st-tab-count" v-if="groupSettings(g.key).length">{{ groupSettings(g.key).length }}</span>
      </button>
    </div>

    <!-- Settings Card -->
    <div class="st-card">
      <div class="st-card-head">
        <h3>{{ activeGroupLabel }}</h3>
        <button @click="saveSettings" class="st-save-btn" :disabled="saving">
          <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5"><polyline points="20 6 9 17 4 12"/></svg>
          {{ saving ? 'Menyimpan...' : 'Simpan Perubahan' }}
        </button>
      </div>
      <div class="st-card-body">
        <div v-for="s in groupSettings(activeGroup)" :key="s.key" class="st-field">
          <div class="st-field-left">
            <label :for="s.key">{{ s.label }}</label>
            <div class="st-field-desc">{{ s.description }}</div>
          </div>
          <div class="st-field-right">
            <select v-if="s.key === 'midtrans_environment'" v-model="formData[s.key]" :id="s.key" class="st-input st-select">
              <option value="sandbox">Sandbox</option>
              <option value="production">Production</option>
            </select>
            <input v-else v-model="formData[s.key]" :id="s.key" :type="isSecret(s.key) ? 'password' : 'text'" class="st-input" :placeholder="s.label" />
            <button v-if="isSecret(s.key)" @click="toggleShow(s.key)" class="st-toggle-eye" :title="shown[s.key] ? 'Sembunyikan' : 'Tampilkan'">
              <svg v-if="!shown[s.key]" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M1 12s4-8 11-8 11 8 11 8-4 8-11 8-11-8-11-8z"/><circle cx="12" cy="12" r="3"/></svg>
              <svg v-else width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M17.94 17.94A10.07 10.07 0 0112 20c-7 0-11-8-11-8a18.45 18.45 0 015.06-5.94M9.9 4.24A9.12 9.12 0 0112 4c7 0 11 8 11 8a18.5 18.5 0 01-2.16 3.19m-6.72-1.07a3 3 0 11-4.24-4.24"/><line x1="1" y1="1" x2="23" y2="23"/></svg>
            </button>
          </div>
        </div>
        <div v-if="!groupSettings(activeGroup).length" class="st-empty">
          <p>Tidak ada pengaturan untuk grup ini</p>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted } from 'vue'
import { settingApi } from '../../services/api'
import { useToastStore } from '../../stores/toast'

const toast = useToastStore()
const settings = ref([])
const activeGroup = ref('midtrans')
const saving = ref(false)
const formData = reactive({})
const shown = reactive({})

const groups = [
  { key: 'midtrans', label: 'Payment Gateway', icon: '<svg width="15" height="15" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><rect x="1" y="4" width="22" height="16" rx="2" ry="2"/><line x1="1" y1="10" x2="23" y2="10"/></svg>' },
  { key: 'general', label: 'Umum', icon: '<svg width="15" height="15" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><circle cx="12" cy="12" r="3"/><path d="M19.4 15a1.65 1.65 0 00.33 1.82l.06.06a2 2 0 010 2.83 2 2 0 01-2.83 0l-.06-.06a1.65 1.65 0 00-1.82-.33 1.65 1.65 0 00-1 1.51V21a2 2 0 01-2 2 2 2 0 01-2-2v-.09A1.65 1.65 0 009 19.4a1.65 1.65 0 00-1.82.33l-.06.06a2 2 0 01-2.83 0 2 2 0 010-2.83l.06-.06A1.65 1.65 0 004.68 15a1.65 1.65 0 00-1.51-1H3a2 2 0 01-2-2 2 2 0 012-2h.09A1.65 1.65 0 004.6 9a1.65 1.65 0 00-.33-1.82l-.06-.06a2 2 0 010-2.83 2 2 0 012.83 0l.06.06A1.65 1.65 0 009 4.68a1.65 1.65 0 001-1.51V3a2 2 0 012-2 2 2 0 012 2v.09a1.65 1.65 0 001 1.51 1.65 1.65 0 001.82-.33l.06-.06a2 2 0 012.83 0 2 2 0 010 2.83l-.06.06A1.65 1.65 0 0019.4 9a1.65 1.65 0 001.51 1H21a2 2 0 012 2 2 2 0 01-2 2h-.09a1.65 1.65 0 00-1.51 1z"/></svg>' },
]

const activeGroupLabel = computed(() => groups.find(g => g.key === activeGroup.value)?.label || '')

function groupSettings(group) {
  return settings.value.filter(s => s.group_name === group)
}

function isSecret(key) {
  return key.includes('key') || key.includes('secret')
}

function toggleShow(key) {
  shown[key] = !shown[key]
  // Toggle input type
  const el = document.getElementById(key)
  if (el) el.type = shown[key] ? 'text' : 'password'
}

onMounted(async () => {
  try {
    const { data } = await settingApi.list()
    settings.value = data.data || []
    // Populate form
    for (const s of settings.value) {
      formData[s.key] = s.value
    }
  } catch {
    toast.error('Gagal memuat pengaturan')
  }
})

async function saveSettings() {
  saving.value = true
  try {
    // Collect only settings in active group
    const payload = {}
    for (const s of groupSettings(activeGroup.value)) {
      payload[s.key] = formData[s.key] || ''
    }
    await settingApi.bulkUpdate(payload)
    toast.success('Pengaturan berhasil disimpan')
  } catch (e) {
    toast.error(e.response?.data?.error || 'Gagal menyimpan')
  } finally {
    saving.value = false
  }
}
</script>

<style scoped>
/* ═══ HERO ═══ */
.st-hero { background: linear-gradient(135deg, #0f0c29 0%, #302b63 50%, #24243e 100%); border-radius: 16px; padding: 32px 36px 24px; margin-bottom: 24px; }
.st-hero-top { display: flex; justify-content: space-between; align-items: flex-start; }
.st-hero-title { font-size: 1.6rem; font-weight: 800; color: #fff; margin: 0 0 4px; }
.st-hero-sub { font-size: .85rem; color: rgba(255,255,255,.5); margin: 0; }

/* ═══ TAB BAR ═══ */
.st-tab-bar { display: flex; gap: 4px; background: #fff; border-radius: 14px; padding: 5px; margin-bottom: 20px; border: 1px solid #e8ecf1; }
.st-tab { flex: 1; display: flex; align-items: center; justify-content: center; gap: 8px; padding: 10px 14px; border-radius: 10px; border: none; background: transparent; color: #64748b; font-size: 0.82rem; font-weight: 500; cursor: pointer; transition: all .2s; font-family: inherit; }
.st-tab:hover { color: #334155; background: #f8fafc; }
.st-tab.active { background: linear-gradient(135deg, #6366f1, #8b5cf6); color: #fff; font-weight: 600; box-shadow: 0 4px 14px rgba(99,102,241,0.25); }
.st-tab-icon { display: flex; align-items: center; }
.st-tab-count { background: #f1f5f9; color: #64748b; font-size: 0.68rem; padding: 1px 7px; border-radius: 8px; font-weight: 700; }
.st-tab.active .st-tab-count { background: rgba(255,255,255,0.25); color: #fff; }

/* ═══ CARD ═══ */
.st-card { background: #fff; border-radius: 16px; border: 1px solid #e8ecf1; }
.st-card-head { display: flex; align-items: center; justify-content: space-between; padding: 20px 28px; border-bottom: 1px solid #f1f5f9; }
.st-card-head h3 { font-size: 1.05rem; font-weight: 700; color: #0f172a; margin: 0; }
.st-save-btn { display: inline-flex; align-items: center; gap: 6px; padding: 9px 20px; border-radius: 10px; border: none; background: linear-gradient(135deg, #6366f1, #8b5cf6); color: #fff; font-size: 0.82rem; font-weight: 600; cursor: pointer; transition: all .2s; box-shadow: 0 4px 14px rgba(99,102,241,0.2); }
.st-save-btn:hover { box-shadow: 0 6px 20px rgba(99,102,241,0.3); transform: translateY(-1px); }
.st-save-btn:disabled { opacity: .6; cursor: not-allowed; transform: none; }
.st-card-body { padding: 8px 0; }

/* ═══ FIELD ═══ */
.st-field { display: flex; align-items: center; justify-content: space-between; padding: 18px 28px; border-bottom: 1px solid #f8fafc; transition: background .15s; }
.st-field:last-child { border-bottom: none; }
.st-field:hover { background: #fafbfc; }
.st-field-left { flex: 1; min-width: 0; }
.st-field-left label { display: block; font-size: 0.88rem; font-weight: 600; color: #0f172a; margin-bottom: 2px; }
.st-field-desc { font-size: 0.75rem; color: #94a3b8; }
.st-field-right { display: flex; align-items: center; gap: 6px; width: 380px; flex-shrink: 0; }
.st-input { flex: 1; padding: 10px 14px; border-radius: 10px; border: 1.5px solid #e2e8f0; background: #fafbfc; color: #1e293b; font-size: 0.85rem; outline: none; transition: all .2s; font-family: inherit; box-sizing: border-box; }
.st-input:focus { border-color: #6366f1; background: #fff; box-shadow: 0 0 0 3px rgba(99,102,241,0.1); }
.st-select { appearance: none; background-image: url("data:image/svg+xml,%3csvg xmlns='http://www.w3.org/2000/svg' fill='none' viewBox='0 0 20 20'%3e%3cpath stroke='%236b7280' stroke-linecap='round' stroke-linejoin='round' stroke-width='1.5' d='M6 8l4 4 4-4'/%3e%3c/svg%3e"); background-position: right 12px center; background-repeat: no-repeat; background-size: 16px; padding-right: 36px; cursor: pointer; }
.st-toggle-eye { width: 36px; height: 36px; border-radius: 8px; border: 1px solid #e2e8f0; background: #fff; display: flex; align-items: center; justify-content: center; cursor: pointer; color: #94a3b8; transition: all .2s; flex-shrink: 0; }
.st-toggle-eye:hover { border-color: #6366f1; color: #6366f1; }

.st-empty { text-align: center; padding: 48px 20px; color: #94a3b8; font-size: 0.85rem; }

@media (max-width: 768px) {
  .st-hero { padding: 24px 20px 18px; }
  .st-field { flex-direction: column; align-items: flex-start; gap: 8px; }
  .st-field-right { width: 100%; }
}
</style>
