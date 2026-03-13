<template>
  <div class="animate-in">
    <!-- Hero -->
    <div class="ai-hero">
      <div class="ai-hero-top">
        <div>
          <h1 class="ai-hero-title"><i class="ri-robot-2-line"></i> AI Konsultan</h1>
          <p class="ai-hero-sub">Kelola knowledge base, system prompt, dan konfigurasi AI</p>
        </div>
        <button class="ai-hero-btn" @click="refreshCache" :disabled="refreshing">
          <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5"><path d="M23 4v6h-6"/><path d="M1 20v-6h6"/><path d="M3.51 9a9 9 0 0 1 14.85-3.36L23 10M1 14l4.64 4.36A9 9 0 0 0 20.49 15"/></svg>
          {{ refreshing ? 'Refreshing...' : 'Refresh Cache' }}
        </button>
      </div>
      <div class="ai-stats-bar">
        <div class="ai-stat"><span class="stat-dot dot-config"></span><span class="stat-label">Config</span><span class="stat-val">{{ configs.length }}</span></div>
        <div class="ai-stat"><span class="stat-dot dot-kb"></span><span class="stat-label">Knowledge</span><span class="stat-val">{{ knowledge.length }}</span></div>
        <div class="ai-stat"><span class="stat-dot dot-prompt"></span><span class="stat-label">Prompts</span><span class="stat-val">{{ prompts.length }}</span></div>
        <div class="ai-stat"><span class="stat-dot dot-cat"></span><span class="stat-label">Kategori</span><span class="stat-val">{{ categories.length }}</span></div>
      </div>
    </div>

    <!-- Tabs -->
    <div class="ai-tab-bar">
      <button v-for="tab in tabs" :key="tab.id" class="ai-tab" :class="{ active: activeTab === tab.id }" @click="activeTab = tab.id">
        <span class="ai-tab-icon"><i :class="tab.icon"></i></span>
        <span>{{ tab.label }}</span>
        <span class="ai-tab-count" v-if="tab.count">{{ tab.count }}</span>
      </button>
    </div>

    <!-- ═══ CONFIG ═══ -->
    <div v-if="activeTab === 'config'" class="ai-card">
      <div class="ai-card-head">
        <h3>Konfigurasi Model AI</h3>
        <button @click="saveConfig" class="ai-save-btn" :disabled="saving">
          <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5"><polyline points="20 6 9 17 4 12"/></svg>
          {{ saving ? 'Menyimpan...' : 'Simpan Perubahan' }}
        </button>
      </div>
      <div class="ai-card-body">
        <div v-for="cfg in configsSorted" :key="cfg.key" class="ai-field">
          <div class="ai-field-left">
            <label>{{ configLabels[cfg.key] || cfg.key }}</label>
            <div class="ai-field-desc">{{ cfg.description }}</div>
          </div>
          <div class="ai-field-right">
            <!-- API Key -->
            <template v-if="cfg.key === 'openai_api_key'">
              <input :type="showApiKey ? 'text' : 'password'" v-model="cfg.value" class="ai-input ai-input-mono" placeholder="sk-proj-..." />
              <button @click="showApiKey = !showApiKey" class="ai-toggle-eye" :title="showApiKey ? 'Sembunyikan' : 'Tampilkan'">
                <i :class="showApiKey ? 'ri-eye-off-line' : 'ri-eye-line'"></i>
              </button>
              <span class="ai-key-status" :class="cfg.value ? 'status-ok' : 'status-empty'">{{ cfg.value ? 'Terpasang' : 'Belum diisi' }}</span>
            </template>
            <!-- Model -->
            <SearchSelect v-else-if="cfg.key === 'openai_model'" v-model="cfg.value" :options="modelOptions" placeholder="Pilih model" :allow-empty="false" />
            <!-- Temperature -->
            <template v-else-if="cfg.key === 'openai_temperature'">
              <input type="range" v-model.number="cfg.value" min="0" max="1" step="0.1" class="ai-range" />
              <span class="ai-range-val">{{ cfg.value }}</span>
              <span class="ai-range-hint">{{ parseFloat(cfg.value) <= 0.3 ? 'Konsisten' : parseFloat(cfg.value) >= 0.8 ? 'Kreatif' : 'Seimbang' }}</span>
            </template>
            <!-- Default -->
            <input v-else type="text" v-model="cfg.value" class="ai-input" />
          </div>
        </div>
      </div>
    </div>

    <!-- ═══ KNOWLEDGE ═══ -->
    <div v-if="activeTab === 'knowledge'">
      <div class="ai-toolbar">
        <div class="ai-search">
          <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="#94a3b8"><circle cx="11" cy="11" r="8" stroke-width="2"/><path d="m21 21-4.35-4.35" stroke-width="2" stroke-linecap="round"/></svg>
          <input v-model="kbSearch" type="text" placeholder="Cari knowledge..." />
        </div>
        <SearchSelect
          v-model="kbCategoryFilter"
          :options="categoryFilterOpts"
          placeholder="Semua Kategori"
          empty-label="Semua Kategori"
          search-placeholder="Cari kategori..."
        />
        <button class="ai-save-btn" @click="showKBForm = true">
          <svg width="14" height="14" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4"/></svg>
          Tambah
        </button>
      </div>

      <div v-if="filteredKnowledge.length" class="ai-card">
        <div class="ai-card-body">
          <div v-for="kb in filteredKnowledge" :key="kb.id" class="ai-field" :class="{ 'field-inactive': !kb.is_active }">
            <div class="ai-field-left" style="flex:1;min-width:0;">
              <label>{{ kb.title }}</label>
              <div class="ai-field-desc" style="display:flex;align-items:center;gap:8px;flex-wrap:wrap;">
                <span v-if="kb.category" class="ai-badge badge-category">{{ kb.category.name }}</span>
                <span class="ai-badge" :class="'badge-p' + (kb.priority >= 8 ? 'high' : kb.priority >= 5 ? 'med' : 'low')">P{{ kb.priority }}</span>
                <span v-for="kw in (kb.keywords || []).slice(0, 3)" :key="kw" class="ai-badge badge-kw">{{ kw }}</span>
              </div>
              <div class="ai-field-desc" style="margin-top:4px;">{{ truncate(kb.content, 100) }}</div>
            </div>
            <div class="ai-field-actions">
              <button @click="editKB(kb)" class="action-btn action-edit" title="Edit">
                <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor"><path d="M11 4H4a2 2 0 0 0-2 2v14a2 2 0 0 0 2 2h14a2 2 0 0 0 2-2v-7" stroke-width="2"/><path d="M18.5 2.5a2.121 2.121 0 0 1 3 3L12 15l-4 1 1-4 9.5-9.5z" stroke-width="2"/></svg>
              </button>
              <button @click="deleteKB(kb)" class="action-btn action-delete" title="Hapus">
                <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor"><polyline points="3 6 5 6 21 6" stroke-width="2"/><path d="M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6m3 0V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2" stroke-width="2"/></svg>
              </button>
            </div>
          </div>
        </div>
      </div>
      <div v-else class="ai-empty">
        <i class="ri-file-unknow-line"></i>
        <h3>Belum ada knowledge base</h3>
        <p>Tambahkan knowledge untuk memperkaya respons AI</p>
      </div>
    </div>

    <!-- ═══ PROMPTS ═══ -->
    <div v-if="activeTab === 'prompts'">
      <div class="ai-toolbar" style="justify-content:flex-end;">
        <button class="ai-save-btn" @click="showPromptForm = true">
          <svg width="14" height="14" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4"/></svg>
          Tambah Prompt
        </button>
      </div>
      <div v-for="prompt in prompts" :key="prompt.id" class="ai-card" style="margin-bottom:16px;">
        <div class="ai-card-head">
          <div style="display:flex;align-items:center;gap:10px;">
            <h3>{{ prompt.name }}</h3>
            <span class="ai-badge" :class="prompt.is_active ? 'badge-active' : 'badge-inactive'">{{ prompt.is_active ? 'Active' : 'Inactive' }}</span>
          </div>
          <button @click="editPrompt(prompt)" class="ai-save-btn" style="font-size:0.78rem;padding:7px 16px;">
            <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor"><path d="M11 4H4a2 2 0 0 0-2 2v14a2 2 0 0 0 2 2h14a2 2 0 0 0 2-2v-7" stroke-width="2"/><path d="M18.5 2.5a2.121 2.121 0 0 1 3 3L12 15l-4 1 1-4 9.5-9.5z" stroke-width="2"/></svg>
            Edit
          </button>
        </div>
        <div class="ai-card-body" style="padding:0;">
          <pre class="ai-prompt-code">{{ truncate(prompt.prompt, 600) }}</pre>
        </div>
      </div>
      <div v-if="!prompts.length" class="ai-empty">
        <i class="ri-chat-off-line"></i>
        <h3>Belum ada system prompt</h3>
        <p>Buat system prompt untuk mengatur perilaku AI</p>
      </div>
    </div>

    <!-- ═══ CATEGORIES ═══ -->
    <div v-if="activeTab === 'categories'">
      <div class="ai-toolbar" style="justify-content:flex-end;">
        <button class="ai-save-btn" @click="showCatForm = true">
          <svg width="14" height="14" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4"/></svg>
          Tambah Kategori
        </button>
      </div>
      <div v-if="categories.length" class="ai-card">
        <div class="ai-card-body">
          <div v-for="cat in categories" :key="cat.id" class="ai-field">
            <div class="ai-field-left">
              <label>{{ cat.name }}</label>
              <div class="ai-field-desc">{{ cat.description || 'Tidak ada deskripsi' }} &middot; Slug: {{ cat.slug }} &middot; Urutan: {{ cat.sort_order }}</div>
            </div>
            <span class="ai-badge" :class="cat.is_active ? 'badge-active' : 'badge-inactive'">{{ cat.is_active ? 'Active' : 'Inactive' }}</span>
          </div>
        </div>
      </div>
      <div v-else class="ai-empty">
        <i class="ri-folder-unknow-line"></i>
        <h3>Belum ada kategori</h3>
        <p>Tambahkan kategori untuk mengelompokkan knowledge base</p>
      </div>
    </div>

    <!-- ═══ MODALS ═══ -->

    <!-- KB Form -->
    <div v-if="showKBForm" class="modal-overlay" @click.self="closeKBForm">
      <div class="modal-content" style="max-width:680px;">
        <div class="modal-header"><h3>{{ editingKB ? 'Edit Knowledge' : 'Tambah Knowledge' }}</h3><button class="modal-close" @click="closeKBForm">&times;</button></div>
        <div class="modal-body">
          <div class="fm-group"><label class="fm-label">Judul <span class="fm-req">*</span></label><input v-model="kbForm.title" class="ai-input" placeholder="Judul knowledge" /></div>
          <div class="fm-row">
            <div class="fm-group"><label class="fm-label">Slug <span class="fm-req">*</span></label><input v-model="kbForm.slug" class="ai-input" placeholder="slug-knowledge" /></div>
            <div class="fm-group"><label class="fm-label">Kategori</label><SearchSelect v-model="kbForm.category_id" :options="categoryFormOpts" placeholder="Pilih Kategori" search-placeholder="Cari kategori..." /></div>
          </div>
          <div class="fm-group"><label class="fm-label">Konten <span class="fm-req">*</span></label><textarea v-model="kbForm.content" class="ai-input ai-textarea" rows="10" placeholder="Konten knowledge (markdown supported)"></textarea></div>
          <div class="fm-row">
            <div class="fm-group"><label class="fm-label">Keywords (pisahkan koma)</label><input v-model="kbKeywordsInput" class="ai-input" placeholder="keyword1, keyword2, ..." /></div>
            <div class="fm-group"><label class="fm-label">Priority (1-10)</label><input v-model.number="kbForm.priority" type="number" min="1" max="10" class="ai-input" /></div>
          </div>
          <div v-if="editingKB" class="fm-group"><label class="fm-check"><input v-model="kbForm.is_active" type="checkbox" /><span>Aktif</span></label></div>
        </div>
        <div class="modal-footer"><button class="btn btn-secondary" @click="closeKBForm">Batal</button><button class="ai-save-btn" @click="saveKB" :disabled="saving">{{ saving ? 'Menyimpan...' : 'Simpan' }}</button></div>
      </div>
    </div>

    <!-- Prompt Form -->
    <div v-if="showPromptForm" class="modal-overlay" @click.self="closePromptForm">
      <div class="modal-content" style="max-width:720px;">
        <div class="modal-header"><h3>{{ editingPrompt ? 'Edit System Prompt' : 'Tambah System Prompt' }}</h3><button class="modal-close" @click="closePromptForm">&times;</button></div>
        <div class="modal-body">
          <div class="fm-group"><label class="fm-label">Nama Prompt <span class="fm-req">*</span></label><input v-model="promptForm.name" class="ai-input" placeholder="default" /></div>
          <div class="fm-group"><label class="fm-label">System Prompt <span class="fm-req">*</span></label><textarea v-model="promptForm.prompt" class="ai-input ai-textarea" rows="18" placeholder="Tulis system prompt untuk AI..." style="font-family:'JetBrains Mono',monospace;font-size:0.8rem;"></textarea></div>
          <div class="fm-group"><label class="fm-check"><input v-model="promptForm.is_active" type="checkbox" /><span>Aktifkan prompt ini</span></label></div>
        </div>
        <div class="modal-footer"><button class="btn btn-secondary" @click="closePromptForm">Batal</button><button class="ai-save-btn" @click="savePrompt" :disabled="saving">{{ saving ? 'Menyimpan...' : 'Simpan' }}</button></div>
      </div>
    </div>

    <!-- Category Form -->
    <div v-if="showCatForm" class="modal-overlay" @click.self="showCatForm = false">
      <div class="modal-content" style="max-width:480px;">
        <div class="modal-header"><h3>Tambah Kategori</h3><button class="modal-close" @click="showCatForm = false">&times;</button></div>
        <div class="modal-body">
          <div class="fm-group"><label class="fm-label">Nama <span class="fm-req">*</span></label><input v-model="catForm.name" class="ai-input" placeholder="Nama kategori" /></div>
          <div class="fm-group"><label class="fm-label">Slug <span class="fm-req">*</span></label><input v-model="catForm.slug" class="ai-input" placeholder="slug-kategori" /></div>
          <div class="fm-group"><label class="fm-label">Deskripsi</label><textarea v-model="catForm.description" class="ai-input ai-textarea" rows="3" placeholder="Deskripsi..."></textarea></div>
          <div class="fm-group"><label class="fm-label">Urutan</label><input v-model.number="catForm.sort_order" type="number" class="ai-input" /></div>
        </div>
        <div class="modal-footer"><button class="btn btn-secondary" @click="showCatForm = false">Batal</button><button class="ai-save-btn" @click="saveCategory" :disabled="saving">Simpan</button></div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { aiApi } from '../../services/api'
import { useToastStore } from '../../stores/toast'
import SearchSelect from '../../components/SearchSelect.vue'

const toast = useToastStore()
const activeTab = ref('config')
const saving = ref(false)
const refreshing = ref(false)
const loading = ref(false)
const knowledge = ref([])
const categories = ref([])
const prompts = ref([])
const configs = ref([])
const kbSearch = ref('')
const kbCategoryFilter = ref('')
const showKBForm = ref(false)
const showPromptForm = ref(false)
const showCatForm = ref(false)
const editingKB = ref(null)
const editingPrompt = ref(null)
const kbKeywordsInput = ref('')
const kbForm = ref({ title: '', slug: '', content: '', category_id: '', keywords: [], priority: 5, is_active: true })
const promptForm = ref({ name: '', prompt: '', is_active: true })
const catForm = ref({ name: '', slug: '', description: '', sort_order: 0 })
const showApiKey = ref(false)

const modelOptions = [
  { value: 'gpt-4o', label: 'GPT-4o (Recommended)' },
  { value: 'gpt-4o-mini', label: 'GPT-4o Mini (Faster)' },
  { value: 'gpt-4-turbo', label: 'GPT-4 Turbo' },
  { value: 'gpt-4', label: 'GPT-4' },
  { value: 'gpt-3.5-turbo', label: 'GPT-3.5 Turbo (Cheapest)' },
]
const categoryFilterOpts = computed(() => categories.value.map(c => ({ value: c.id, label: c.name })))
const categoryFormOpts = computed(() => categories.value.map(c => ({ value: c.id, label: c.name })))

const configLabels = { openai_api_key: 'API Key OpenAI', openai_model: 'Model AI', openai_temperature: 'Temperature (Kreativitas)', openai_max_tokens: 'Max Tokens' }
const configOrder = ['openai_api_key', 'openai_model', 'openai_temperature', 'openai_max_tokens']
const configsSorted = computed(() => [...configs.value].sort((a, b) => {
  const ai = configOrder.indexOf(a.key), bi = configOrder.indexOf(b.key)
  return (ai === -1 ? 99 : ai) - (bi === -1 ? 99 : bi)
}))

const tabs = computed(() => [
  { id: 'config', label: 'Konfigurasi', icon: 'ri-settings-4-line' },
  { id: 'knowledge', label: 'Knowledge Base', icon: 'ri-book-open-line', count: knowledge.value.length },
  { id: 'prompts', label: 'System Prompts', icon: 'ri-chat-settings-line', count: prompts.value.length },
  { id: 'categories', label: 'Kategori', icon: 'ri-folder-line', count: categories.value.length },
])

const filteredKnowledge = computed(() => {
  let items = knowledge.value
  if (kbSearch.value) {
    const q = kbSearch.value.toLowerCase()
    items = items.filter(k => k.title.toLowerCase().includes(q) || k.content.toLowerCase().includes(q) || (k.keywords && k.keywords.some(kw => kw.toLowerCase().includes(q))))
  }
  if (kbCategoryFilter.value) items = items.filter(k => k.category_id === kbCategoryFilter.value)
  return items
})

async function loadAll() {
  loading.value = true
  try {
    const [kbRes, catRes, promptRes, configRes] = await Promise.all([aiApi.listKnowledge(), aiApi.listCategories(), aiApi.listPrompts(), aiApi.getConfig()])
    knowledge.value = kbRes.data.items || []
    categories.value = catRes.data.categories || []
    prompts.value = promptRes.data.prompts || []
    configs.value = configRes.data.configs || []
  } catch { toast.error('Gagal memuat data AI') }
  finally { loading.value = false }
}

async function saveConfig() {
  saving.value = true
  try { await aiApi.updateConfig({ configs: configs.value }); toast.success('Konfigurasi AI berhasil disimpan') }
  catch { toast.error('Gagal menyimpan konfigurasi') }
  finally { saving.value = false }
}

function editKB(kb) { editingKB.value = kb; kbForm.value = { ...kb, category_id: kb.category_id || '' }; kbKeywordsInput.value = (kb.keywords || []).join(', '); showKBForm.value = true }
function closeKBForm() { showKBForm.value = false; editingKB.value = null; kbForm.value = { title: '', slug: '', content: '', category_id: '', keywords: [], priority: 5, is_active: true }; kbKeywordsInput.value = '' }

async function saveKB() {
  if (!kbForm.value.title || !kbForm.value.slug || !kbForm.value.content) { toast.error('Judul, slug, dan konten wajib diisi'); return }
  saving.value = true
  const data = { ...kbForm.value, keywords: kbKeywordsInput.value.split(',').map(k => k.trim()).filter(Boolean) }
  try {
    if (editingKB.value) { await aiApi.updateKnowledge(editingKB.value.id, data); toast.success('Knowledge berhasil diupdate') }
    else { await aiApi.createKnowledge(data); toast.success('Knowledge berhasil ditambahkan') }
    closeKBForm(); loadAll()
  } catch { toast.error('Gagal menyimpan knowledge') }
  finally { saving.value = false }
}

async function deleteKB(kb) {
  if (!confirm(`Hapus knowledge "${kb.title}"?`)) return
  try { await aiApi.deleteKnowledge(kb.id); toast.success('Knowledge dihapus'); loadAll() }
  catch { toast.error('Gagal menghapus') }
}

function editPrompt(p) { editingPrompt.value = p; promptForm.value = { ...p }; showPromptForm.value = true }
function closePromptForm() { showPromptForm.value = false; editingPrompt.value = null; promptForm.value = { name: '', prompt: '', is_active: true } }

async function savePrompt() {
  if (!promptForm.value.name || !promptForm.value.prompt) { toast.error('Nama dan prompt wajib diisi'); return }
  saving.value = true
  try {
    if (editingPrompt.value) { await aiApi.updatePrompt(editingPrompt.value.id, promptForm.value); toast.success('System prompt diupdate') }
    else { await aiApi.createPrompt(promptForm.value); toast.success('System prompt ditambahkan') }
    closePromptForm(); loadAll()
  } catch { toast.error('Gagal menyimpan prompt') }
  finally { saving.value = false }
}

async function saveCategory() {
  if (!catForm.value.name || !catForm.value.slug) { toast.error('Nama dan slug wajib diisi'); return }
  saving.value = true
  try { await aiApi.createCategory(catForm.value); toast.success('Kategori ditambahkan'); showCatForm.value = false; catForm.value = { name: '', slug: '', description: '', sort_order: 0 }; loadAll() }
  catch { toast.error('Gagal menyimpan kategori') }
  finally { saving.value = false }
}

async function refreshCache() {
  refreshing.value = true
  try { await aiApi.invalidateCache(); toast.success('Cache AI berhasil di-refresh') }
  catch { toast.error('Gagal refresh cache') }
  finally { refreshing.value = false }
}

function truncate(str, len) { if (!str) return ''; return str.length > len ? str.slice(0, len) + '...' : str }

onMounted(loadAll)
</script>

<style scoped>
/* ═══ HERO ═══ */
.ai-hero { background: linear-gradient(180deg, #1a1c2e 0%, #22243a 100%); border-radius: var(--radius-xl); padding: 32px 36px 24px; margin-bottom: 24px; box-shadow: 0 4px 24px rgba(26,28,46,0.3); }
.ai-hero-top { display: flex; justify-content: space-between; align-items: flex-start; gap: 16px; flex-wrap: wrap; }
.ai-hero-title { font-size: 1.6rem; font-weight: 800; color: #fff !important; margin: 0 0 4px; display: flex; align-items: center; gap: 10px; }
.ai-hero-title i { color: #fff; }
.ai-hero-sub { font-size: .85rem; color: rgba(255,255,255,.5); margin: 0; }
.ai-hero-btn { display: inline-flex; align-items: center; gap: 8px; padding: 11px 24px; font-size: 0.85rem; font-weight: 700; border-radius: 12px; background: var(--gradient-primary); color: white; border: none; cursor: pointer; box-shadow: 0 4px 20px rgba(253,150,68,0.3); transition: all .2s; font-family: inherit; }
.ai-hero-btn:hover { box-shadow: 0 6px 24px rgba(253,150,68,0.4); transform: translateY(-1px); }
.ai-hero-btn:disabled { opacity: .6; cursor: not-allowed; transform: none; }
.ai-stats-bar { display: flex; gap: 28px; padding-top: 18px; border-top: 1px solid rgba(255,255,255,0.08); margin-top: 20px; }
.ai-stat { display: flex; align-items: center; gap: 8px; }
.stat-dot { width: 8px; height: 8px; border-radius: 50%; }
.dot-config { background: var(--primary); box-shadow: 0 0 8px rgba(253,150,68,0.5); }
.dot-kb { background: var(--success); box-shadow: 0 0 8px rgba(34,197,94,0.5); }
.dot-prompt { background: var(--warning); box-shadow: 0 0 8px rgba(245,158,11,0.5); }
.dot-cat { background: #f472b6; box-shadow: 0 0 8px rgba(244,114,182,0.5); }
.stat-label { font-size: 0.72rem; color: rgba(255,255,255,0.4); text-transform: uppercase; letter-spacing: 0.05em; }
.stat-val { font-size: 0.9rem; font-weight: 800; color: white; }

/* ═══ TAB BAR ═══ */
.ai-tab-bar { display: flex; gap: 4px; background: #fff; border-radius: 14px; padding: 5px; margin-bottom: 20px; border: 1px solid var(--bgray-200); }
.ai-tab { flex: 1; display: flex; align-items: center; justify-content: center; gap: 8px; padding: 10px 14px; border-radius: 10px; border: none; background: transparent; color: var(--bgray-600); font-size: 0.82rem; font-weight: 500; cursor: pointer; transition: all .2s; font-family: inherit; }
.ai-tab:hover { color: var(--bgray-900); background: var(--bgray-50); }
.ai-tab.active { background: var(--gradient-primary); color: #fff; font-weight: 600; box-shadow: 0 4px 14px rgba(253,150,68,0.25); }
.ai-tab-icon { display: flex; align-items: center; }
.ai-tab-count { background: var(--bgray-100); color: var(--bgray-600); font-size: 0.68rem; padding: 1px 7px; border-radius: 8px; font-weight: 700; }
.ai-tab.active .ai-tab-count { background: rgba(255,255,255,0.25); color: #fff; }

/* ═══ CARD ═══ */
.ai-card { background: #fff; border-radius: var(--radius-xl); border: 1px solid var(--bgray-200); }
.ai-card-head { display: flex; align-items: center; justify-content: space-between; padding: 20px 28px; border-bottom: 1px solid var(--bgray-200); }
.ai-card-head h3 { font-size: 1.05rem; font-weight: 700; color: var(--bgray-900) !important; margin: 0; }
.ai-card-body { padding: 8px 0; }

/* ═══ SAVE BUTTON ═══ */
.ai-save-btn { display: inline-flex; align-items: center; gap: 6px; padding: 9px 20px; border-radius: 10px; border: none; background: var(--gradient-primary); color: #fff; font-size: 0.82rem; font-weight: 600; cursor: pointer; transition: all .2s; box-shadow: 0 4px 14px rgba(253,150,68,0.25); font-family: inherit; }
.ai-save-btn:hover { box-shadow: 0 6px 20px rgba(253,150,68,0.35); transform: translateY(-1px); }
.ai-save-btn:disabled { opacity: .6; cursor: not-allowed; transform: none; }

/* ═══ FIELD ═══ */
.ai-field { display: flex; align-items: center; justify-content: space-between; padding: 18px 28px; border-bottom: 1px solid var(--bgray-100); transition: background .15s; gap: 16px; }
.ai-field:last-child { border-bottom: none; }
.ai-field:hover { background: var(--bgray-50); }
.ai-field.field-inactive { opacity: 0.5; }
.ai-field-left { flex-shrink: 1; min-width: 0; }
.ai-field-left label { display: block; font-size: 0.88rem; font-weight: 600; color: var(--bgray-900) !important; margin-bottom: 2px; }
.ai-field-desc { font-size: 0.75rem; color: var(--bgray-500); line-height: 1.4; }
.ai-field-right { display: flex; align-items: center; gap: 8px; width: 380px; flex-shrink: 0; }
.ai-field-actions { display: flex; gap: 4px; flex-shrink: 0; }

/* ═══ INPUT ═══ */
.ai-input { flex: 1; padding: 10px 14px; border-radius: 10px; border: 1.5px solid var(--bgray-200); background: var(--bgray-50); color: var(--bgray-900); font-size: 0.85rem; outline: none; transition: all .2s; font-family: inherit; box-sizing: border-box; width: 100%; }
.ai-input:focus { border-color: var(--primary); background: #fff; box-shadow: 0 0 0 3px rgba(253,150,68,0.15); }
.ai-input-mono { font-family: 'JetBrains Mono', 'Fira Code', monospace; font-size: 0.8rem; }
.ai-select { appearance: none; background-image: url("data:image/svg+xml,%3csvg xmlns='http://www.w3.org/2000/svg' fill='none' viewBox='0 0 20 20'%3e%3cpath stroke='%236b7280' stroke-linecap='round' stroke-linejoin='round' stroke-width='1.5' d='M6 8l4 4 4-4'/%3e%3c/svg%3e"); background-position: right 12px center; background-repeat: no-repeat; background-size: 16px; padding-right: 36px; cursor: pointer; }
.ai-textarea { resize: vertical; line-height: 1.6; min-height: 80px; }

/* ═══ TOGGLE EYE ═══ */
.ai-toggle-eye { width: 36px; height: 36px; border-radius: 8px; border: 1px solid var(--bgray-200); background: #fff; display: flex; align-items: center; justify-content: center; cursor: pointer; color: var(--bgray-500); transition: all .2s; flex-shrink: 0; font-size: 16px; }
.ai-toggle-eye:hover { border-color: var(--primary); color: var(--primary); }
.ai-key-status { font-size: 0.72rem; font-weight: 600; white-space: nowrap; padding: 4px 10px; border-radius: 8px; }
.status-ok { color: #16a34a; background: #f0fdf4; }
.status-empty { color: #d97706; background: #fffbeb; }

/* ═══ RANGE ═══ */
.ai-range { flex: 1; height: 6px; appearance: none; background: var(--bgray-200); border-radius: 99px; outline: none; cursor: pointer; }
.ai-range::-webkit-slider-thumb { appearance: none; width: 20px; height: 20px; border-radius: 50%; background: var(--gradient-primary); cursor: pointer; box-shadow: 0 2px 6px rgba(253,150,68,0.3); }
.ai-range-val { font-weight: 800; font-size: 0.95rem; color: var(--primary); min-width: 28px; text-align: center; }
.ai-range-hint { font-size: 0.72rem; color: var(--bgray-500); white-space: nowrap; background: var(--bgray-50); padding: 3px 10px; border-radius: 6px; }

/* ═══ BADGES ═══ */
.ai-badge { font-size: 0.68rem; padding: 2px 9px; border-radius: 6px; font-weight: 700; white-space: nowrap; }
.badge-category { background: var(--primary-100); color: var(--primary-dark); }
.badge-phigh { background: #fee2e2; color: #991b1b; }
.badge-pmed { background: #fef3c7; color: #92400e; }
.badge-plow { background: #d1fae5; color: #065f46; }
.badge-kw { background: var(--bgray-100); color: var(--bgray-600); }
.badge-active { background: #d1fae5; color: #065f46; }
.badge-inactive { background: var(--bgray-100); color: var(--bgray-500); }

/* ═══ TOOLBAR ═══ */
.ai-toolbar { display: flex; gap: 12px; margin-bottom: 20px; align-items: center; }
.ai-search { display: flex; align-items: center; gap: 10px; background: white; border: 1.5px solid var(--bgray-200); border-radius: var(--radius-lg); padding: 0 16px; height: 42px; flex: 1; }
.ai-search:focus-within { border-color: var(--primary); box-shadow: 0 0 0 3px rgba(253,150,68,0.12); }
.ai-search input { flex: 1; border: none; outline: none; font-size: 0.85rem; color: var(--bgray-900); background: none; font-family: inherit; }
.ai-filter-select { height: 42px; padding: 0 34px 0 14px; border-radius: var(--radius-lg); border: 1.5px solid var(--bgray-200); font-size: 0.82rem; font-weight: 500; color: var(--bgray-700); background: white; cursor: pointer; appearance: none; background-image: url("data:image/svg+xml,%3csvg xmlns='http://www.w3.org/2000/svg' fill='none' viewBox='0 0 20 20'%3e%3cpath stroke='%236b7280' stroke-linecap='round' stroke-linejoin='round' stroke-width='1.5' d='M6 8l4 4 4-4'/%3e%3c/svg%3e"); background-position: right 10px center; background-repeat: no-repeat; background-size: 1.2em; font-family: inherit; }

/* ═══ PROMPT CODE ═══ */
.ai-prompt-code { background: var(--secondary); color: #cbd5e1; padding: 20px 28px; font-family: 'JetBrains Mono', 'Fira Code', monospace; font-size: 0.78rem; line-height: 1.7; overflow-x: auto; white-space: pre-wrap; max-height: 320px; overflow-y: auto; margin: 0; border-radius: 0 0 var(--radius-xl) var(--radius-xl); }

/* ═══ ACTION BUTTONS ═══ */
.action-btn { width: 34px; height: 34px; border-radius: 8px; border: 1px solid var(--bgray-200); background: #fff; display: flex; align-items: center; justify-content: center; cursor: pointer; color: var(--bgray-500); transition: all .15s; }
.action-btn:hover { border-color: var(--primary); color: var(--primary); background: var(--primary-50); }
.action-delete:hover { border-color: var(--danger); color: var(--danger); background: #fee2e2; }

/* ═══ EMPTY ═══ */
.ai-empty { display: flex; flex-direction: column; align-items: center; padding: 80px 32px; background: white; border-radius: var(--radius-xl); border: 2px dashed var(--bgray-200); text-align: center; }
.ai-empty i { font-size: 3rem; color: var(--bgray-500); margin-bottom: 16px; }
.ai-empty h3 { font-size: 1.15rem; font-weight: 700; color: var(--bgray-900) !important; margin: 0; }
.ai-empty p { font-size: 0.85rem; color: var(--bgray-500); margin-top: 6px; }

/* ═══ MODAL ═══ */
.modal-overlay { position: fixed; inset: 0; background: rgba(26,28,46,0.6); backdrop-filter: blur(4px); display: flex; align-items: center; justify-content: center; z-index: 9999; padding: 20px; }
.modal-content { background: #fff; border-radius: var(--radius-xl); width: 100%; max-height: 90vh; display: flex; flex-direction: column; box-shadow: var(--shadow-xl); }
.modal-header { display: flex; align-items: center; justify-content: space-between; padding: 20px 28px; border-bottom: 1px solid var(--bgray-200); }
.modal-header h3 { font-size: 1.05rem; font-weight: 700; margin: 0; color: var(--bgray-900) !important; }
.modal-close { width: 32px; height: 32px; border: none; background: var(--bgray-100); border-radius: 8px; font-size: 1.2rem; cursor: pointer; display: flex; align-items: center; justify-content: center; color: var(--bgray-600); transition: all .15s; }
.modal-close:hover { background: var(--bgray-200); color: var(--bgray-900); }
.modal-body { padding: 24px 28px; overflow-y: auto; flex: 1; }
.modal-footer { padding: 16px 28px; border-top: 1px solid var(--bgray-200); display: flex; justify-content: flex-end; gap: 10px; }

/* ═══ FORM ═══ */
.fm-group { margin-bottom: 16px; }
.fm-row { display: grid; grid-template-columns: 1fr 1fr; gap: 16px; }
.fm-label { display: block; font-size: 0.8rem; font-weight: 600; color: var(--bgray-700); margin-bottom: 6px; }
.fm-req { color: var(--danger); }
.fm-check { display: flex; align-items: center; gap: 8px; font-size: 0.85rem; cursor: pointer; color: var(--bgray-700); }
.fm-check input[type="checkbox"] { width: 18px; height: 18px; accent-color: var(--primary); }

.btn { display: inline-flex; align-items: center; gap: 6px; padding: 9px 20px; border-radius: 10px; font-size: 0.82rem; font-weight: 600; cursor: pointer; transition: all .2s; font-family: inherit; }
.btn-secondary { background: var(--bgray-100); border: 1px solid var(--bgray-200); color: var(--bgray-700); }
.btn-secondary:hover { background: var(--bgray-200); }

@media (max-width: 768px) {
  .ai-hero { padding: 24px 20px 18px; }
  .ai-hero-title { font-size: 1.3rem; }
  .ai-stats-bar { gap: 16px; flex-wrap: wrap; }
  .ai-field { flex-direction: column; align-items: flex-start; }
  .ai-field-right { width: 100%; }
  .ai-toolbar { flex-direction: column; }
  .fm-row { grid-template-columns: 1fr; }
}
</style>
