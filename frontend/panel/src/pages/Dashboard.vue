<template>
  <div class="animate-in space-y-6">
    <!-- Stat Widgets — BanCo style -->
    <div class="grid grid-cols-1 md:grid-cols-2 xl:grid-cols-4 gap-6">
      <!-- Total Leads -->
      <div class="widget-card">
        <div class="widget-header">
          <div class="widget-icon widget-icon-blue">
            <svg width="24" height="24" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
              <ellipse cx="11.78" cy="17.56" rx="7.78" ry="4.44" fill="currentColor" opacity="0.3"/>
              <circle cx="11.78" cy="6.44" r="4.44" fill="currentColor"/>
            </svg>
          </div>
          <span class="widget-title">Total Leads</span>
        </div>
        <div class="widget-body">
          <p class="widget-value blue">{{ stats.total_leads ?? '—' }}</p>
          <div class="widget-trend up">
            <svg width="16" height="14" viewBox="0 0 16 14" fill="none"><path d="M13.43 0.52H7.57C6.28 0.52 5.48 1.92 6.13 3.03L9.06 8.05C9.7 9.16 11.3 9.16 11.94 8.05L14.87 3.03C15.52 1.92 14.72 0.52 13.43 0.52Z" fill="currentColor"/><path opacity="0.4" d="M2.17 13.05H8.03C9.32 13.05 10.12 11.65 9.47 10.54L6.54 5.52C5.9 4.41 4.3 4.41 3.66 5.52L0.73 10.54C0.08 11.65 0.88 13.05 2.17 13.05Z" fill="currentColor"/></svg>
            <span>Active pipeline</span>
          </div>
        </div>
      </div>

      <!-- Mitra Aktif -->
      <div class="widget-card">
        <div class="widget-header">
          <div class="widget-icon widget-icon-green">
            <svg width="20" height="20" viewBox="0 0 20 20" fill="none" xmlns="http://www.w3.org/2000/svg">
              <path d="M18 11C18 15.97 13.97 20 9 20C4.03 20 0 15.97 0 11C0 6.03 4.03 2 9 2C13.97 2 18 6.03 18 11Z" fill="currentColor" opacity="0.3"/>
              <path d="M19.8 8.01C19.01 4.08 15.92 0.99 11.99 0.2C10.9 -0.02 10 0.9 10 2V8C10 9.1 10.9 10 12 10H18C19.1 10 20.02 9.1 19.8 8.01Z" fill="currentColor"/>
            </svg>
          </div>
          <span class="widget-title">Mitra Aktif</span>
        </div>
        <div class="widget-body">
          <p class="widget-value green">{{ stats.active_mitra ?? '—' }}</p>
          <div class="widget-trend up">
            <svg width="16" height="14" viewBox="0 0 16 14" fill="none"><path d="M13.43 0.52H7.57C6.28 0.52 5.48 1.92 6.13 3.03L9.06 8.05C9.7 9.16 11.3 9.16 11.94 8.05L14.87 3.03C15.52 1.92 14.72 0.52 13.43 0.52Z" fill="currentColor"/><path opacity="0.4" d="M2.17 13.05H8.03C9.32 13.05 10.12 11.65 9.47 10.54L6.54 5.52C5.9 4.41 4.3 4.41 3.66 5.52L0.73 10.54C0.08 11.65 0.88 13.05 2.17 13.05Z" fill="currentColor"/></svg>
            <span>Running partnerships</span>
          </div>
        </div>
      </div>

      <!-- Total Investasi -->
      <div class="widget-card">
        <div class="widget-header">
          <div class="widget-icon widget-icon-purple">
            <svg width="20" height="20" viewBox="0 0 20 20" fill="none"><path d="M20 4C20 1.79 18.21 0 16 0H4C1.79 0 0 1.79 0 4V14C0 16.21 1.79 18 4 18H16C18.21 18 20 16.21 20 14V4Z" fill="currentColor" opacity="0.3"/><path d="M6 9C6 7.34 4.66 6 3 6H0V12H3C4.66 12 6 10.66 6 9Z" fill="currentColor"/></svg>
          </div>
          <span class="widget-title">Total Investasi</span>
        </div>
        <div class="widget-body">
          <p class="widget-value purple">{{ formatCurrencyShort(stats.total_investment) }}</p>
          <div class="widget-trend up">
            <svg width="16" height="14" viewBox="0 0 16 14" fill="none"><path d="M13.43 0.52H7.57C6.28 0.52 5.48 1.92 6.13 3.03L9.06 8.05C9.7 9.16 11.3 9.16 11.94 8.05L14.87 3.03C15.52 1.92 14.72 0.52 13.43 0.52Z" fill="currentColor"/><path opacity="0.4" d="M2.17 13.05H8.03C9.32 13.05 10.12 11.65 9.47 10.54L6.54 5.52C5.9 4.41 4.3 4.41 3.66 5.52L0.73 10.54C0.08 11.65 0.88 13.05 2.17 13.05Z" fill="currentColor"/></svg>
            <span>Cumulative value</span>
          </div>
        </div>
      </div>

      <!-- Revenue Bulan Ini -->
      <div class="widget-card">
        <div class="widget-header">
          <div class="widget-icon widget-icon-orange">
            <svg width="20" height="20" viewBox="0 0 20 20" fill="none"><path d="M0 4C0 1.79 1.79 0 4 0H16C18.21 0 20 1.79 20 4V16C20 18.21 18.21 20 16 20H4C1.79 20 0 18.21 0 16V4Z" fill="currentColor" opacity="0.3"/><path d="M14 9C12.9 9 12 9.9 12 11V13C12 14.1 12.9 15 14 15C15.1 15 16 14.1 16 13V11C16 9.9 15.1 9 14 9Z" fill="currentColor"/><path d="M6 5C4.9 5 4 5.9 4 7V13C4 14.1 4.9 15 6 15C7.1 15 8 14.1 8 13V7C8 5.9 7.1 5 6 5Z" fill="currentColor"/></svg>
          </div>
          <span class="widget-title">Revenue Bulan Ini</span>
        </div>
        <div class="widget-body">
          <p class="widget-value orange">{{ formatCurrencyShort(stats.monthly_revenue) }}</p>
          <div class="widget-trend up">
            <svg width="16" height="14" viewBox="0 0 16 14" fill="none"><path d="M13.43 0.52H7.57C6.28 0.52 5.48 1.92 6.13 3.03L9.06 8.05C9.7 9.16 11.3 9.16 11.94 8.05L14.87 3.03C15.52 1.92 14.72 0.52 13.43 0.52Z" fill="currentColor"/><path opacity="0.4" d="M2.17 13.05H8.03C9.32 13.05 10.12 11.65 9.47 10.54L6.54 5.52C5.9 4.41 4.3 4.41 3.66 5.52L0.73 10.54C0.08 11.65 0.88 13.05 2.17 13.05Z" fill="currentColor"/></svg>
            <span>This month</span>
          </div>
        </div>
      </div>
    </div>

    <!-- Charts — BanCo layout -->
    <div class="grid grid-cols-1 xl:grid-cols-3 gap-6">
      <!-- Lead Pipeline (2/3 width) -->
      <div class="xl:col-span-2 card-bankco">
        <div class="card-bankco-header">
          <h3 class="card-bankco-title">Lead Pipeline</h3>
          <router-link to="/leads" class="card-bankco-action">Lihat Semua →</router-link>
        </div>
        <div class="card-bankco-body">
          <div v-if="stats.leads_by_status && Object.keys(stats.leads_by_status).length" class="space-y-4">
            <div v-for="(count, status) in stats.leads_by_status" :key="status" class="pipeline-row">
              <span class="pipeline-badge" :class="getStatusClass(status)">{{ formatStatus(status) }}</span>
              <div class="pipeline-bar-wrapper">
                <div class="pipeline-bar">
                  <div class="pipeline-bar-fill" :class="getStatusClass(status)" :style="{ width: pct(count) + '%' }"></div>
                </div>
              </div>
              <span class="pipeline-count">{{ count }}</span>
            </div>
          </div>
          <div v-else class="empty-placeholder">
            <div class="empty-placeholder-icon">📊</div>
            <p>Belum ada data lead</p>
          </div>
        </div>
      </div>

      <!-- Revenue (1/3 width) -->
      <div class="card-bankco">
        <div class="card-bankco-header">
          <h3 class="card-bankco-title">Revenue</h3>
          <span class="card-bankco-sub">6 Bulan Terakhir</span>
        </div>
        <div class="card-bankco-body">
          <div v-if="stats.revenue_chart?.length" class="space-y-4">
            <div v-for="item in stats.revenue_chart" :key="item.month" class="revenue-row">
              <span class="revenue-month">{{ item.month }}</span>
              <div class="revenue-bar-wrapper">
                <div class="revenue-bar">
                  <div class="revenue-bar-fill" :style="{ width: revPct(item.revenue) + '%' }"></div>
                </div>
              </div>
              <span class="revenue-amount">{{ formatCurrency(item.revenue) }}</span>
            </div>
          </div>
          <div v-else class="empty-placeholder">
            <div class="empty-placeholder-icon">📈</div>
            <p>Belum ada data revenue</p>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { dashboardApi } from '../services/api'

const stats = ref({})

onMounted(() => { loadStats() })

async function loadStats() {
  try {
    const { data } = await dashboardApi.stats()
    stats.value = data.data || {}
  } catch (e) { console.error(e) }
}

function formatCurrency(v) { return new Intl.NumberFormat('id-ID', { style: 'currency', currency: 'IDR', minimumFractionDigits: 0 }).format(v || 0) }
function formatCurrencyShort(v) { return new Intl.NumberFormat('id-ID', { style: 'currency', currency: 'IDR', minimumFractionDigits: 0, notation: 'compact' }).format(v || 0) }
function formatStatus(s) { return s.replace(/_/g, ' ') }
function getStatusClass(s) {
  return {
    NEW: 'status-info', RUNNING: 'status-success', COMPLETED: 'status-success',
    DP_PAID: 'status-warning', FULLY_PAID: 'status-primary'
  }[s] || 'status-neutral'
}
function pct(c) { const m = Math.max(...Object.values(stats.value.leads_by_status || {0:1})); return (c/m)*100 }
function revPct(r) { const m = Math.max(...(stats.value.revenue_chart||[{revenue:1}]).map(x=>x.revenue)); return (r/m)*100 }
</script>
