<template>
  <div class="animate-in">
    <router-link to="/leads" class="btn btn-ghost btn-sm mb-4">← Kembali ke Kanban</router-link>

    <div v-if="lead" class="space-y-5">
      <!-- Info Card -->
      <div class="card">
        <div class="card-body">
          <div class="flex items-start justify-between mb-4">
            <div>
              <h2 class="text-xl font-bold">{{ lead.full_name }}</h2>
              <p class="text-sm text-gray-500">{{ lead.email }} · {{ lead.phone || '-' }}</p>
            </div>
          </div>
          <div class="flex items-center gap-3 mb-3">
            <span class="badge" :class="statusBadge(lead.status)">{{ lead.status.replace(/_/g, ' ') }}</span>
            <span class="text-sm text-gray-500">Progress: <strong>{{ lead.progress_percentage }}%</strong></span>
          </div>
          <div class="progress-bar"><div class="progress-bar-fill" :style="{ width: lead.progress_percentage + '%' }"></div></div>
        </div>
      </div>

      <!-- Status Update -->
      <div class="card">
        <div class="card-header"><span>🔄 Update Status</span></div>
        <div class="card-body">
          <div class="flex flex-wrap gap-2">
            <button v-for="s in statuses" :key="s" @click="updateStatus(s)"
              class="btn btn-sm" :class="lead.status === s ? 'btn-primary' : 'btn-secondary'" :disabled="lead.status === s">
              {{ s.replace(/_/g, ' ') }}
            </button>
          </div>
        </div>
      </div>

      <!-- Notes -->
      <div class="card">
        <div class="card-header"><span>📝 Catatan</span></div>
        <div class="card-body">
          <p class="text-gray-600 text-sm">{{ lead.notes || 'Tidak ada catatan.' }}</p>
        </div>
      </div>

      <!-- Convert to Partnership -->
      <div v-if="['ACTIVE_PARTNERSHIP','RUNNING'].includes(lead.status)" class="card border-green-200 bg-green-50">
        <div class="card-body text-center">
          <p class="text-green-700 font-semibold mb-3">Lead ini siap dikonversi menjadi Partnership!</p>
          <button @click="convertToPartnership" class="btn btn-success" :disabled="converting">
            {{ converting ? 'Membuat...' : '🤝 Buat Partnership' }}
          </button>
        </div>
      </div>
    </div>
    <div v-else class="text-center py-16 text-gray-400">Memuat data...</div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { leadApi, partnershipApi } from '../../services/api'

const route = useRoute()
const router = useRouter()
const lead = ref(null)
const converting = ref(false)
const statuses = ['NEW','CONSULTATION','LOCATION_SUBMITTED','SURVEY_APPROVED','MEETING_DONE','READY_FOR_DP','DP_PAID','AGREEMENT_REVIEW','FULLY_PAID','ACTIVE_PARTNERSHIP','RUNNING','COMPLETED']

onMounted(loadLead)

async function loadLead() { const { data } = await leadApi.get(route.params.id); lead.value = data.data }

async function updateStatus(status) {
  try { await leadApi.updateStatus(lead.value.id, { status }); await loadLead() } catch (e) { console.error(e) }
}

async function convertToPartnership() {
  converting.value = true
  try {
    await partnershipApi.create({ lead_id: lead.value.id, brand_id: lead.value.brand_id || '00000000-0000-0000-0000-000000000000', mitra_id: lead.value.sales_id })
    router.push('/partnerships')
  } catch (e) { console.error(e) } finally { converting.value = false }
}

function statusBadge(s) { return ['RUNNING','COMPLETED'].includes(s) ? 'badge-success' : ['DP_PAID','FULLY_PAID'].includes(s) ? 'badge-warning' : 'badge-info' }
</script>
