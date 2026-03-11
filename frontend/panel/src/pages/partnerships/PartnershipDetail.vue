<template>
  <div class="pd-page">
    <!-- Back button -->
    <router-link to="/partnerships" class="pd-back">
      <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><polyline points="15 18 9 12 15 6"/></svg>
      Kembali ke Partnership
    </router-link>

    <div v-if="p">
      <!-- Hero Header -->
      <div class="pd-hero">
        <div class="pd-hero-bg"></div>
        <div class="pd-hero-content">
          <div class="pd-hero-left">
            <div class="pd-avatar" :style="avatarBg(p.mitra?.name||'?')">{{ getInitial(p.mitra?.name||'?') }}</div>
            <div class="pd-hero-info">
              <h1 class="pd-hero-name">{{ p.mitra?.name || '-' }}</h1>
              <div class="pd-hero-email">{{ p.mitra?.email || '' }}</div>
              <div class="pd-hero-phone" v-if="p.mitra?.phone">{{ p.mitra.phone }}</div>
            </div>
          </div>
          <div class="pd-hero-right">
            <div class="pd-progress-ring-wrap">
              <svg class="pd-progress-ring" viewBox="0 0 80 80">
                <circle cx="40" cy="40" r="34" stroke="rgba(255,255,255,0.1)" stroke-width="6" fill="none"/>
                <circle cx="40" cy="40" r="34" stroke="url(#pd-grad)" stroke-width="6" fill="none"
                  stroke-linecap="round" :stroke-dasharray="213.6" :stroke-dashoffset="213.6 - (213.6 * p.progress_percentage / 100)"
                  transform="rotate(-90 40 40)"/>
                <defs><linearGradient id="pd-grad" x1="0" y1="0" x2="1" y2="1"><stop offset="0%" stop-color="#667eea"/><stop offset="100%" stop-color="#764ba2"/></linearGradient></defs>
              </svg>
              <div class="pd-progress-value">{{ p.progress_percentage }}%</div>
            </div>
            <span class="pd-status-badge" :class="'st-'+p.status">{{ statusLabel(p.status) }}</span>
          </div>
        </div>
      </div>

      <!-- Status Pipeline -->
      <div class="pd-pipeline-wrap">
        <div class="pd-pipeline-header">
          <h3 class="pd-pipeline-title"><i class="ri-flow-chart"></i> Status Partnership</h3>
        </div>
        <div class="pd-pipeline">
          <div v-for="(st, idx) in statusSteps" :key="st.val" class="pd-pipe-step"
            :class="{ active: p.status === st.val, completed: statusIdx(p.status) > idx, selected: selectedStatus === st.val && selectedStatus !== p.status }"
            @click="selectedStatus = st.val">
            <div class="pd-pipe-line" v-if="idx > 0" :class="{ done: statusIdx(p.status) >= idx }"></div>
            <div class="pd-pipe-dot" :style="pipeDotStyle(st, idx)">
              <i :class="statusIdx(p.status) > idx ? 'ri-check-line' : st.icon"></i>
            </div>
            <div class="pd-pipe-text">
              <div class="pd-pipe-label">{{ st.label }}</div>
              <div class="pd-pipe-prog">{{ st.progress }}%</div>
            </div>
          </div>
        </div>
        <transition name="pd-slide">
          <div v-if="selectedStatus && selectedStatus !== p.status" class="pd-status-confirm">
            <div class="pd-sc-info">
              <span class="pd-sc-from">{{ statusLabel(p.status) }}</span>
              <i class="ri-arrow-right-line pd-sc-arrow"></i>
              <span class="pd-sc-to">{{ statusLabel(selectedStatus) }}</span>
            </div>
            <button @click="doUpdateStatus" class="pd-sc-btn" :disabled="updatingStatus">
              <i :class="updatingStatus ? 'ri-loader-4-line ri-spin' : 'ri-check-double-line'"></i>
              {{ updatingStatus ? 'Memproses...' : 'Konfirmasi' }}
            </button>
          </div>
        </transition>
      </div>

      <!-- Info Cards Row -->
      <div class="pd-info-grid">
        <div class="pd-info-card">
          <div class="pd-info-icon ic-affiliator">
            <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M12 2L15.09 8.26L22 9.27L17 14.14L18.18 21.02L12 17.77L5.82 21.02L7 14.14L2 9.27L8.91 8.26L12 2Z"/></svg>
          </div>
          <div class="pd-info-body">
            <div class="pd-info-label">Affiliator</div>
            <div class="pd-info-value">{{ p.affiliator?.name || '-' }}</div>
            <div class="pd-info-sub" v-if="p.affiliator?.email">{{ p.affiliator.email }}</div>
          </div>
        </div>
        <div class="pd-info-card">
          <div class="pd-info-icon ic-outlet">
            <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M3 9l9-7 9 7v11a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2z"/><polyline points="9 22 9 12 15 12 15 22"/></svg>
          </div>
          <div class="pd-info-body">
            <div class="pd-info-label">Outlet</div>
            <div class="pd-info-value">{{ p.outlet?.name || '-' }}</div>
          </div>
        </div>
        <div class="pd-info-card">
          <div class="pd-info-icon ic-pkg">
            <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M21 16V8a2 2 0 0 0-1-1.73l-7-4a2 2 0 0 0-2 0l-7 4A2 2 0 0 0 3 8v8a2 2 0 0 0 1 1.73l7 4a2 2 0 0 0 2 0l7-4A2 2 0 0 0 21 16z"/><polyline points="3.27 6.96 12 12.01 20.73 6.96"/><line x1="12" y1="22.08" x2="12" y2="12"/></svg>
          </div>
          <div class="pd-info-body">
            <div class="pd-info-label">Paket</div>
            <div class="pd-info-value">{{ p.package?.name || '-' }}</div>
            <div class="pd-info-sub pd-price" v-if="p.package?.price">{{ fc(p.package.price) }}</div>
          </div>
        </div>
        <div class="pd-info-card">
          <div class="pd-info-icon ic-date">
            <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><rect x="3" y="4" width="18" height="18" rx="2" ry="2"/><line x1="16" y1="2" x2="16" y2="6"/><line x1="8" y1="2" x2="8" y2="6"/><line x1="3" y1="10" x2="21" y2="10"/></svg>
          </div>
          <div class="pd-info-body">
            <div class="pd-info-label">Tanggal Mulai</div>
            <div class="pd-info-value">{{ p.start_date ? formatDate(p.start_date) : 'Belum dimulai' }}</div>
          </div>
        </div>
      </div>

      <!-- Tab Bar -->
      <div class="pd-tab-bar">
        <button v-for="t in tabs" :key="t.key" class="pd-tab" :class="{ active: tab === t.key }" @click="tab = t.key">
          <span class="pd-tab-icon" v-html="t.icon"></span>
          <span>{{ t.label }}</span>
          <span class="pd-tab-count" v-if="t.count">{{ t.count }}</span>
        </button>
      </div>
      <!-- Tab Content -->
      <div class="pd-content-card">
        <!-- Invoices -->
        <template v-if="tab==='invoices'">
          <div class="pd-content-head">
            <h3>Daftar Invoice</h3>
            <div class="pd-content-head-actions">
              <button v-if="hasPendingInvoices" @click="syncAllInvoices" class="pd-sync-btn" :disabled="syncing" title="Sinkronkan status semua invoice pending dengan Midtrans">
                <svg :class="{ 'spin-icon': syncing }" width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5"><path d="M21 12a9 9 0 1 1-6.219-8.56" stroke-linecap="round"/></svg>
                {{ syncing ? 'Menyinkronkan...' : 'Sync Status' }}
              </button>
              <button @click="openInvModal" class="pd-add-btn">
                <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5"><line x1="12" y1="5" x2="12" y2="19"/><line x1="5" y1="12" x2="19" y2="12"/></svg>
                Buat Invoice
              </button>
            </div>
          </div>
          <div v-if="invoices.length" class="pd-items">
            <div v-for="inv in invoices" :key="inv.id" class="pd-item">
              <div class="pd-item-left">
                <div class="pd-item-icon" :class="inv.status==='PAID'?'pi-full':'pi-dp'">
                  <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z"/><polyline points="14 2 14 8 20 8"/><line x1="16" y1="13" x2="8" y2="13"/><line x1="16" y1="17" x2="8" y2="17"/></svg>
                </div>
                <div>
                  <div class="pd-item-title">
                    <span class="pd-inv-type-tag" :class="'invt-'+inv.type">{{ invTypeLabel(inv.type) }}</span>
                    {{ inv.invoice_number }}
                  </div>
                  <div class="pd-item-date">
                    {{ formatDate(inv.created_at) }} · {{ inv.description || '-' }}
                    <span v-if="inv.expired_at && inv.status==='PENDING'" class="pd-expiry-hint" :class="{ 'pd-expiry-soon': isExpiringSoon(inv) }">
                      · {{ expiryLabel(inv) }}
                    </span>
                  </div>
                </div>
              </div>
              <div class="pd-item-right">
                <div class="pd-item-amount">{{ fc(inv.amount) }}</div>
                <div class="pd-item-actions">
                  <span class="pd-badge" :class="invBadge(inv.status)">{{ invLabel(inv.status) }}</span>
                  <!-- Payment method label for PAID invoices -->
                  <span v-if="inv.status==='PAID'" class="pd-pay-method" :class="inv.midtrans_payment_type==='manual'?'pm-manual':'pm-auto'">
                    <svg v-if="inv.midtrans_payment_type==='manual'" width="10" height="10" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5"><path d="M16 21v-2a4 4 0 0 0-4-4H5a4 4 0 0 0-4 4v2"/><circle cx="8.5" cy="7" r="4"/><polyline points="17 11 19 13 23 9"/></svg>
                    <svg v-else width="10" height="10" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5"><polyline points="20 6 9 17 4 12"/></svg>
                    {{ inv.midtrans_payment_type==='manual' ? 'Manual' : 'Otomatis' }}
                  </span>
                  <!-- View proof for manual -->
                  <button v-if="inv.status==='PAID' && inv.midtrans_payment_type==='manual' && inv.proof_url" @click="openProofViewer(inv)" class="pd-proof-btn">
                    <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M1 12s4-8 11-8 11 8 11 8-4 8-11 8-11-8-11-8z"/><circle cx="12" cy="12" r="3"/></svg>
                    Lihat Bukti
                  </button>
                  <!-- Actions for PENDING -->
                  <button v-if="inv.midtrans_redirect_url && inv.status==='PENDING'" @click="copyLink(inv.midtrans_redirect_url)" class="pd-verify-btn" title="Salin link pembayaran">
                    <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><rect x="9" y="9" width="13" height="13" rx="2"/><path d="M5 15H4a2 2 0 0 1-2-2V4a2 2 0 0 1 2-2h9a2 2 0 0 1 2 2v1"/></svg>
                    Salin Link
                  </button>
                  <a v-if="inv.midtrans_redirect_url && inv.status==='PENDING'" :href="inv.midtrans_redirect_url" target="_blank" class="pd-verify-btn">Bayar →</a>
                  <button v-if="inv.status==='PENDING'" @click="openApproveModal(inv)" class="pd-approve-btn" title="Verifikasi pembayaran manual">
                    <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5"><polyline points="20 6 9 17 4 12"/></svg>
                    Verifikasi
                  </button>
                  <!-- Check Status from Midtrans -->
                  <button v-if="inv.status==='PENDING' && inv.midtrans_order_id" @click="checkInvoiceStatus(inv)" class="pd-check-btn" :disabled="checkingInv[inv.id]" title="Cek status pembayaran di Midtrans">
                    <svg :class="{ 'spin-icon': checkingInv[inv.id] }" width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5"><path d="M21 12a9 9 0 1 1-6.219-8.56" stroke-linecap="round"/></svg>
                    {{ checkingInv[inv.id] ? 'Mengecek...' : 'Cek Status' }}
                  </button>
                </div>
              </div>
            </div>
          </div>
          <div v-else class="pd-empty">
            <svg width="40" height="40" viewBox="0 0 24 24" fill="none" stroke="#94a3b8" stroke-width="1.5"><path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z"/><polyline points="14 2 14 8 20 8"/></svg>
            <p>Belum ada invoice</p>
          </div>
        </template>

        <!-- Agreements -->
        <template v-else-if="tab==='agreements'">
          <div class="pd-content-head">
            <h3>Dokumen Agreement</h3>
            <button @click="showAgrModal=true" class="pd-add-btn">
              <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5"><line x1="12" y1="5" x2="12" y2="19"/><line x1="5" y1="12" x2="19" y2="12"/></svg>
              Upload Agreement
            </button>
          </div>
          <div v-if="agreements.length" class="pd-items">
            <div v-for="a in agreements" :key="a.id" class="pd-item">
              <div class="pd-item-left">
                <div class="pd-item-icon" :class="a.type==='CONTRACT'?'pi-agr':'pi-doc'">
                  <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z"/><polyline points="14 2 14 8 20 8"/></svg>
                </div>
                <div>
                  <div class="pd-item-title">
                    <span class="pd-agr-type-tag" :class="a.type==='CONTRACT'?'agt-contract':'agt-document'">{{ a.type==='CONTRACT'?'Kontrak':'Dokumen' }}</span>
                    {{ a.title || ('Agreement Versi ' + a.version) }}
                  </div>
                  <div class="pd-item-date">{{ formatDate(a.created_at) }} · Versi {{ a.version }}</div>
                </div>
              </div>
              <div class="pd-item-right">
                <a v-if="a.file_url" :href="a.file_url" target="_blank" class="pd-verify-btn" style="margin-right:6px">Lihat File</a>
                <template v-if="a.type==='CONTRACT'">
                  <span class="pd-badge" :class="a.status==='SIGNED'?'bg-verified':'bg-pending'">{{ a.status==='SIGNED' ? 'Sudah Ditandatangani' : 'Belum Ditandatangani' }}</span>
                  <button v-if="a.status!=='SIGNED'" @click="signAgreement(a.id)" class="pd-verify-btn">Tandatangani</button>
                </template>
                <span v-else class="pd-badge bg-draft">Dokumen</span>
              </div>
            </div>
          </div>
          <div v-else class="pd-empty">
            <svg width="40" height="40" viewBox="0 0 24 24" fill="none" stroke="#475569" stroke-width="1.5"><path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z"/><polyline points="14 2 14 8 20 8"/></svg>
            <p>Belum ada agreement</p>
          </div>
        </template>

        <!-- Revenue -->
        <template v-else-if="tab==='revenue'">
          <div class="pd-content-head">
            <h3>Laporan Revenue</h3>
            <button @click="showRevModal=true" class="pd-add-btn">
              <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5"><line x1="12" y1="5" x2="12" y2="19"/><line x1="5" y1="12" x2="19" y2="12"/></svg>
              Input Revenue
            </button>
          </div>
          <div v-if="revenues.length" class="pd-items">
            <div v-for="r in revenues" :key="r.id" class="pd-rev-card">
              <div class="pd-rev-month">{{ r.month }}</div>
              <div class="pd-rev-grid">
                <div class="pd-rev-item"><span class="pd-rev-label">Revenue</span><span class="pd-rev-val rv-green">{{ fc(r.revenue) }}</span></div>
                <div class="pd-rev-item"><span class="pd-rev-label">Expense</span><span class="pd-rev-val rv-red">{{ fc(r.expense) }}</span></div>
                <div class="pd-rev-item"><span class="pd-rev-label">Profit</span><span class="pd-rev-val rv-profit">{{ fc(r.profit) }}</span></div>
                <div class="pd-rev-item"><span class="pd-rev-label">Mitra Share</span><span class="pd-rev-val rv-share">{{ fc(r.mitra_share) }}</span></div>
              </div>
            </div>
          </div>
          <div v-else class="pd-empty">
            <svg width="40" height="40" viewBox="0 0 24 24" fill="none" stroke="#475569" stroke-width="1.5"><polyline points="22 12 18 12 15 21 9 3 6 12 2 12"/></svg>
            <p>Belum ada data revenue</p>
          </div>
        </template>

        <!-- Lokasi -->
        <template v-else-if="tab==='locations'">
          <div class="pd-content-head">
            <h3><i class="ri-map-pin-line" style="margin-right:4px;color:#6366f1"></i> Peninjauan Lokasi</h3>
            <button v-if="canAddLocation" @click="showLocModal=true" class="pd-add-btn">
              <i class="ri-add-line"></i> Tambah Lokasi
            </button>
            <span v-else-if="locations.length" class="pd-loc-info-badge">
              <i class="ri-information-line"></i> Sudah ada lokasi aktif ({{ locStatusLabel(activeLocationStatus) }})
            </span>
          </div>
          <div v-if="locations.length" class="pd-items">
            <div v-for="loc in locations" :key="loc.id" class="pd-loc-card" @click="$router.push({name:'LocationDetail',params:{id:loc.id}})">
              <div class="pd-loc-head">
                <div class="pd-loc-left">
                  <span class="pd-loc-badge" :class="'lbs-'+loc.status">{{ locStatusLabel(loc.status) }}</span>
                  <span v-if="loc.tipe_bangunan" class="pd-loc-type"><i class="ri-building-2-line"></i> {{ loc.tipe_bangunan }}</span>
                </div>
                <div class="pd-loc-score" :style="{color:locScoreColor(loc.total_score)}"><i class="ri-bar-chart-fill"></i> {{ loc.total_score }} <span>{{ loc.score_category }}</span></div>
              </div>
              <div class="pd-loc-name">{{ loc.nama_lokasi }}</div>
              <div class="pd-loc-meta">
                <span><i class="ri-map-pin-line"></i> {{ loc.kota || '-' }}</span>
                <span><i class="ri-ruler-line"></i> {{ loc.luas_tempat }} m²</span>
                <span><i class="ri-money-dollar-circle-line"></i> {{ fc(loc.harga_sewa_per_tahun) }}/thn</span>
              </div>
            </div>
          </div>
          <div v-else class="pd-empty">
            <i class="ri-map-pin-add-line" style="font-size:2.5rem;color:#94a3b8"></i>
            <p>Belum ada pengajuan lokasi untuk partnership ini</p>
          </div>
        </template>
      </div>
    </div>

    <!-- Loading -->
    <div v-else class="pd-loading">
      <div class="pd-spinner"></div>
      <p>Memuat data partnership...</p>
    </div>

    <!-- Agreement Modal -->
    <Teleport to="body">
      <div v-if="showAgrModal" class="pd-overlay" @click.self="showAgrModal=false">
        <div class="pd-modal" @click.stop>
          <div class="pd-modal-head">
            <div class="pd-modal-icon mi-agr"><svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z"/><polyline points="14 2 14 8 20 8"/></svg></div>
            <h3>Upload Agreement</h3>
            <button @click="showAgrModal=false" class="pd-modal-x">&times;</button>
          </div>
          <form @submit.prevent="createAgreement" class="pd-modal-body">
            <div class="pd-fg">
              <label>Tipe Dokumen <span class="req">*</span></label>
              <select v-model="agrForm.type" class="pd-input" required>
                <option value="CONTRACT">Kontrak (perlu ditandatangani)</option>
                <option value="DOCUMENT">Dokumen Lainnya</option>
              </select>
            </div>
            <div class="pd-fg">
              <label>Judul <span class="req">*</span></label>
              <input v-model="agrForm.title" class="pd-input" :placeholder="agrForm.type==='CONTRACT'?'Contoh: Perjanjian Kerjasama Franchise':'Contoh: Lampiran Dokumen Tambahan'" required />
            </div>
            <div class="pd-fg">
              <label>Upload File <span class="req">*</span></label>
              <div class="pd-upload-zone" @click="$refs.agrFileInput.click()" @dragover.prevent @drop.prevent="onAgrFileDrop">
                <input ref="agrFileInput" type="file" accept=".pdf,.docx,.doc" hidden @change="onAgrFileChange">
                <div v-if="!agrFilePreview" class="pd-upload-placeholder">
                  <svg width="28" height="28" viewBox="0 0 24 24" fill="none" stroke="#94a3b8" stroke-width="1.5"><path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4"/><polyline points="17 8 12 3 7 8"/><line x1="12" y1="3" x2="12" y2="15"/></svg>
                  <span>Klik atau drag file di sini</span>
                  <span class="pd-upload-hint">PDF, DOCX (maks 10MB)</span>
                </div>
                <div v-else class="pd-upload-preview">
                  <div class="pd-file-name">{{ agrFileName }}</div>
                  <button type="button" @click.stop="clearAgrFile" class="pd-upload-clear">&times;</button>
                </div>
              </div>
            </div>
            <div class="pd-modal-foot">
              <button type="button" @click="showAgrModal=false" class="pd-btn-sec">Batal</button>
              <button type="submit" class="pd-btn-primary" :disabled="agrUploading || !agrFile || !agrForm.title">{{ agrUploading ? 'Mengupload...' : 'Upload' }}</button>
            </div>
          </form>
        </div>
      </div>
    </Teleport>

    <!-- Revenue Modal -->
    <Teleport to="body">
      <div v-if="showRevModal" class="pd-overlay" @click.self="showRevModal=false">
        <div class="pd-modal" @click.stop>
          <div class="pd-modal-head">
            <div class="pd-modal-icon mi-rev"><svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><polyline points="22 12 18 12 15 21 9 3 6 12 2 12"/></svg></div>
            <h3>Input Revenue Bulanan</h3>
            <button @click="showRevModal=false" class="pd-modal-x">&times;</button>
          </div>
          <form @submit.prevent="createRevenue" class="pd-modal-body">
            <div class="pd-fg"><label>Bulan <span class="req">*</span></label><input v-model="revForm.month" type="month" class="pd-input" required></div>
            <div class="pd-form-row">
              <div class="pd-fg"><label>Revenue (Rp) <span class="req">*</span></label><input v-model.number="revForm.revenue" type="number" class="pd-input" required></div>
              <div class="pd-fg"><label>Expense (Rp) <span class="req">*</span></label><input v-model.number="revForm.expense" type="number" class="pd-input" required></div>
            </div>
            <div class="pd-modal-foot"><button type="button" @click="showRevModal=false" class="pd-btn-sec">Batal</button><button type="submit" class="pd-btn-primary">Simpan</button></div>
          </form>
        </div>
      </div>
    </Teleport>

    <!-- Invoice Modal -->
    <Teleport to="body">
      <div v-if="showInvModal" class="pd-overlay" @click.self="closeInvModal">
        <div class="pd-modal pd-modal-lg" @click.stop>
          <div class="pd-modal-head">
            <div class="pd-modal-icon mi-pay"><svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z"/><polyline points="14 2 14 8 20 8"/><line x1="16" y1="13" x2="8" y2="13"/><line x1="16" y1="17" x2="8" y2="17"/></svg></div>
            <h3>{{ invSuccess ? 'Invoice Berhasil Dibuat!' : 'Buat Invoice Baru' }}</h3>
            <button @click="closeInvModal" class="pd-modal-x">&times;</button>
          </div>

          <!-- Success State -->
          <div v-if="invSuccess" class="pd-modal-body">
            <div class="pd-inv-success">
              <div class="pd-inv-success-icon">
                <svg width="48" height="48" viewBox="0 0 24 24" fill="none" stroke="#22c55e" stroke-width="2"><path d="M22 11.08V12a10 10 0 1 1-5.93-9.14"/><polyline points="22 4 12 14.01 9 11.01"/></svg>
              </div>
              <div class="pd-inv-success-info">
                <div class="pd-inv-number">{{ lastInvoice?.invoice_number }}</div>
                <div class="pd-inv-amt">{{ fc(lastInvoice?.amount || 0) }}</div>
                <div class="pd-inv-desc">{{ lastInvoice?.description }}</div>
              </div>
            </div>
            <div class="pd-inv-link-box">
              <label>Link Pembayaran Midtrans</label>
              <div class="pd-inv-link-row">
                <input :value="lastInvoice?.midtrans_redirect_url" class="pd-input" readonly @click="$event.target.select()" />
                <button type="button" @click="copyLink(lastInvoice?.midtrans_redirect_url)" class="pd-copy-btn">
                  <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><rect x="9" y="9" width="13" height="13" rx="2"/><path d="M5 15H4a2 2 0 0 1-2-2V4a2 2 0 0 1 2-2h9a2 2 0 0 1 2 2v1"/></svg>
                  {{ copied ? 'Tersalin!' : 'Salin' }}
                </button>
              </div>
              <div class="pd-inv-link-hint">Kirim link ini ke mitra untuk pembayaran</div>
            </div>
            <div class="pd-modal-foot">
              <a :href="lastInvoice?.midtrans_redirect_url" target="_blank" class="pd-btn-primary">Buka Halaman Pembayaran →</a>
              <button type="button" @click="closeInvModal" class="pd-btn-sec">Tutup</button>
            </div>
          </div>

          <!-- Form State -->
          <form v-else @submit.prevent="createInvoice" class="pd-modal-body">
            <!-- Package Info Card -->
            <div class="pd-inv-pkg-card" v-if="p">
              <div class="pd-inv-pkg-icon">
                <svg width="22" height="22" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M21 16V8a2 2 0 0 0-1-1.73l-7-4a2 2 0 0 0-2 0l-7 4A2 2 0 0 0 3 8v8a2 2 0 0 0 1 1.73l7 4a2 2 0 0 0 2 0l7-4A2 2 0 0 0 21 16z"/></svg>
              </div>
              <div class="pd-inv-pkg-info">
                <div class="pd-inv-pkg-label">Paket yang Diambil</div>
                <div class="pd-inv-pkg-name">{{ p.outlet?.name }} — {{ p.package?.name || 'Paket Standar' }}</div>
                <div class="pd-inv-pkg-price">Harga Paket: {{ fc(p.package?.price || 0) }}</div>
                <div class="pd-inv-pkg-mitra">Mitra: {{ p.mitra?.name || '-' }} ({{ p.mitra?.email || '-' }})</div>
              </div>
            </div>

            <div class="pd-fg">
              <label>Tipe Pembayaran <span class="req">*</span></label>
              <select v-model="invForm.type" class="pd-input" required @change="onInvTypeChange">
                <option value="DP">Down Payment (DP)</option>
                <option value="CICILAN">Cicilan</option>
                <option value="PELUNASAN">Pelunasan</option>
                <option value="INVOICE">Invoice Lainnya</option>
              </select>
            </div>
            <div class="pd-fg"><label>Jumlah Tagihan (Rp) <span class="req">*</span></label><input v-model.number="invForm.amount" type="number" class="pd-input" placeholder="Nominal invoice" required></div>
            <div class="pd-fg"><label>Deskripsi <span class="req">*</span></label><input v-model="invForm.description" class="pd-input" :placeholder="'Contoh: DP ' + (p?.package?.name || 'Paket Franchise')" required></div>
            <div class="pd-modal-foot"><button type="button" @click="closeInvModal" class="pd-btn-sec">Batal</button><button type="submit" class="pd-btn-primary" :disabled="creatingInv">{{ creatingInv ? 'Memproses...' : 'Buat Invoice & Dapatkan Link' }}</button></div>
          </form>
        </div>
      </div>
    </Teleport>

    <!-- Approve Confirmation Modal -->
    <Teleport to="body">
      <div v-if="showApproveModal" class="pd-overlay" @click.self="showApproveModal=false">
        <div class="pd-modal pd-modal-confirm" @click.stop>
          <div class="pd-confirm-body">
            <div class="pd-confirm-icon">
              <svg width="32" height="32" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <path d="M9 12l2 2 4-4"/>
                <circle cx="12" cy="12" r="10"/>
              </svg>
            </div>
            <h3 class="pd-confirm-title">Verifikasi Pembayaran</h3>
            <p class="pd-confirm-desc">Anda akan menandai invoice berikut sebagai <strong>sudah dibayar</strong> secara manual. Tindakan ini tidak dapat dibatalkan.</p>
            <div class="pd-confirm-invoice" v-if="approveTarget">
              <div class="pd-confirm-inv-row">
                <span class="pd-confirm-inv-label">Invoice</span>
                <span class="pd-confirm-inv-val">{{ approveTarget.invoice_number }}</span>
              </div>
              <div class="pd-confirm-inv-row">
                <span class="pd-confirm-inv-label">Tipe</span>
                <span class="pd-inv-type-tag" :class="'invt-'+approveTarget.type">{{ invTypeLabel(approveTarget.type) }}</span>
              </div>
              <div class="pd-confirm-inv-row">
                <span class="pd-confirm-inv-label">Nominal</span>
                <span class="pd-confirm-inv-amount">{{ fc(approveTarget.amount) }}</span>
              </div>
            </div>

            <!-- Upload Bukti -->
            <div class="pd-confirm-upload">
              <label class="pd-confirm-upload-label">
                <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4"/><polyline points="17 8 12 3 7 8"/><line x1="12" y1="3" x2="12" y2="15"/></svg>
                Upload Bukti Pembayaran <span class="req">*</span>
              </label>
              <div class="pd-upload-zone" @click="$refs.proofInput.click()" @dragover.prevent @drop.prevent="onProofDrop">
                <input ref="proofInput" type="file" accept="image/*,.pdf" hidden @change="onProofFileChange" />
                <div v-if="!proofPreview" class="pd-upload-placeholder">
                  <svg width="28" height="28" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5"><rect x="3" y="3" width="18" height="18" rx="2"/><circle cx="8.5" cy="8.5" r="1.5"/><polyline points="21 15 16 10 5 21"/></svg>
                  <span>Klik atau drag foto bukti transfer</span>
                  <span class="pd-upload-hint">JPG, PNG, PDF — maks. 10MB</span>
                </div>
                <div v-else class="pd-upload-preview">
                  <img v-if="proofPreview !== 'file'" :src="proofPreview" />
                  <span v-else class="pd-file-name">{{ proofFileName }}</span>
                  <button type="button" @click.stop="clearProofFile" class="pd-upload-clear">&times;</button>
                </div>
              </div>
            </div>
          </div>
          <div class="pd-confirm-footer">
            <button @click="closeApproveModal" class="pd-btn-sec">Batal</button>
            <button @click="confirmApprove" class="pd-btn-approve" :disabled="approvingInv || !proofFile">
              <svg v-if="!approvingInv" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5"><polyline points="20 6 9 17 4 12"/></svg>
              <div v-else class="pd-btn-spinner"></div>
              {{ approvingInv ? 'Mengupload & Memproses...' : 'Ya, Verifikasi Lunas' }}
            </button>
          </div>
        </div>
      </div>
    </Teleport>

    <!-- Proof Viewer Modal -->
    <Teleport to="body">
      <div v-if="showProofViewer" class="pd-overlay" @click.self="showProofViewer=false">
        <div class="pd-modal pd-modal-proof" @click.stop>
          <div class="pd-proof-header">
            <div class="pd-proof-header-info">
              <h3>Bukti Pembayaran</h3>
              <div class="pd-proof-meta">
                <span class="pd-inv-type-tag" :class="'invt-'+(proofViewTarget?.type||'INVOICE')">{{ invTypeLabel(proofViewTarget?.type) }}</span>
                <span>{{ proofViewTarget?.invoice_number }}</span>
                <span class="pd-proof-meta-sep">·</span>
                <span class="pd-proof-meta-amount">{{ fc(proofViewTarget?.amount||0) }}</span>
              </div>
            </div>
            <button @click="showProofViewer=false" class="pd-modal-x">&times;</button>
          </div>
          <div class="pd-proof-viewer-body">
            <img :src="proofViewTarget?.proof_url" class="pd-proof-image" @error="$event.target.style.display='none'" />
          </div>
          <div class="pd-proof-viewer-foot">
            <span class="pd-pay-method pm-manual" style="font-size:.72rem">
              <svg width="10" height="10" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5"><path d="M16 21v-2a4 4 0 0 0-4-4H5a4 4 0 0 0-4 4v2"/><circle cx="8.5" cy="7" r="4"/><polyline points="17 11 19 13 23 9"/></svg>
              Verifikasi Manual
            </span>
            <span v-if="proofViewTarget?.paid_at" class="pd-proof-paid-date">Dibayar {{ formatDate(proofViewTarget.paid_at) }}</span>
            <a :href="proofViewTarget?.proof_url" target="_blank" class="pd-verify-btn" style="margin-left:auto">Buka Full →</a>
          </div>
        </div>
      </div>
    </Teleport>

    <!-- Location Create Modal -->
    <Teleport to="body">
      <div v-if="showLocModal" class="pd-overlay" @click.self="showLocModal=false">
        <div class="pd-modal pd-modal-lg" style="max-width:820px" @click.stop>
          <div class="pd-modal-head">
            <div class="pd-modal-icon" style="background:#dbeafe;color:#2563eb"><i class="ri-map-pin-add-fill" style="font-size:1.1rem"></i></div>
            <h3>Pengajuan Lokasi Baru</h3>
            <button @click="showLocModal=false" class="pd-modal-x">&times;</button>
          </div>
          <form @submit.prevent="createLoc" class="pd-modal-body" style="max-height:68vh;overflow-y:auto">
            <div class="pd-loc-section">
              <div class="pd-loc-sh"><i class="ri-map-pin-2-fill" style="color:#2563eb"></i> Informasi Umum</div>
              <div class="pd-fg"><label>Nama Lokasi <span class="req">*</span></label><input v-model="locForm.nama_lokasi" class="pd-input" required placeholder="Ruko Jl. Sudirman No. 5" /></div>
              <div class="pd-fg"><label>Alamat Lengkap <span class="req">*</span></label><textarea v-model="locForm.alamat" class="pd-input" style="min-height:60px;resize:vertical" required placeholder="Alamat lengkap lokasi"></textarea></div>
              <div class="pd-form-row" style="grid-template-columns:1fr 1fr 1fr">
                <div class="pd-fg"><label>Provinsi</label><input v-model="locForm.provinsi" class="pd-input" placeholder="Jawa Timur" /></div>
                <div class="pd-fg"><label>Kota/Kab</label><input v-model="locForm.kota" class="pd-input" placeholder="Surabaya" /></div>
                <div class="pd-fg"><label>Kecamatan</label><input v-model="locForm.kecamatan" class="pd-input" placeholder="Gubeng" /></div>
              </div>
              <div class="pd-fg" style="max-width:200px"><label>Kode Pos</label><input v-model="locForm.kode_pos" class="pd-input" placeholder="60281" /></div>
            </div>
            <div class="pd-loc-section">
              <div class="pd-loc-sh"><i class="ri-building-2-fill" style="color:#7c3aed"></i> Detail Lokasi</div>
              <div class="pd-form-row" style="grid-template-columns:1fr 1fr 1fr">
                <div class="pd-fg"><label>Luas (m²)</label><input v-model.number="locForm.luas_tempat" type="number" class="pd-input" placeholder="50" /></div>
                <div class="pd-fg"><label>Sewa/Tahun (Rp)</label><input v-model.number="locForm.harga_sewa_per_tahun" type="number" class="pd-input" placeholder="50000000" /></div>
                <div class="pd-fg"><label>Durasi Sewa (thn)</label><input v-model.number="locForm.durasi_sewa" type="number" class="pd-input" placeholder="3" /></div>
              </div>
              <div class="pd-form-row" style="grid-template-columns:1fr 1fr 1fr">
                <div class="pd-fg"><label>Tipe Bangunan</label>
                  <select v-model="locForm.tipe_bangunan" class="pd-input">
                    <option value="">Pilih tipe</option>
                    <option value="ruko">Ruko</option><option value="stand">Stand</option><option value="mall">Mall</option><option value="kios">Kios</option><option value="lainnya">Lainnya</option>
                  </select>
                </div>
                <div class="pd-fg"><label>Lebar Jalan (m)</label><input v-model.number="locForm.lebar_jalan" type="number" step="any" class="pd-input" placeholder="8" /></div>
                <div class="pd-fg"><label>Jumlah Lantai</label><input v-model.number="locForm.jumlah_lantai" type="number" class="pd-input" placeholder="2" /></div>
              </div>
            </div>
            <div class="pd-loc-section">
              <div class="pd-loc-sh"><i class="ri-bar-chart-2-fill" style="color:#16a34a"></i> Traffic & Potensi</div>
              <div class="pd-form-row">
                <div class="pd-fg"><label>Est. Lalu Lintas/Hari</label><input v-model.number="locForm.estimasi_lalu_lintas" type="number" class="pd-input" placeholder="3000" /></div>
                <div class="pd-fg"><label>Kompetitor (500m)</label><input v-model.number="locForm.jumlah_kompetitor" type="number" class="pd-input" placeholder="2" /></div>
              </div>
              <div class="pd-form-row">
                <div class="pd-fg"><label>Dekat Dengan</label><input v-model="locForm.dekat_dengan" class="pd-input" placeholder="sekolah, kampus, pasar..." /></div>
                <div class="pd-fg"><label>Target Market</label><input v-model="locForm.target_market" class="pd-input" placeholder="Pelajar, pekerja..." /></div>
              </div>
            </div>
            <div class="pd-modal-foot">
              <button type="button" @click="showLocModal=false" class="pd-btn-sec">Batal</button>
              <button type="submit" class="pd-btn-primary" :disabled="creatingLoc">
                <i :class="creatingLoc?'ri-loader-4-line ri-spin':'ri-check-line'" style="margin-right:4px"></i>
                {{ creatingLoc ? 'Menyimpan...' : 'Simpan & Hitung Score' }}
              </button>
            </div>
          </form>
        </div>
      </div>
    </Teleport>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { partnershipApi, agreementApi, revenueApi, invoiceApi, uploadApi, locationApi } from '../../services/api'
import { useToastStore } from '../../stores/toast'

const toast = useToastStore()
const route = useRoute()
const id = route.params.id
const p = ref(null)
const tab = ref('invoices')
const agreements = ref([]), revenues = ref([]), invoices = ref([]), locations = ref([])
const showAgrModal = ref(false), showRevModal = ref(false), showInvModal = ref(false), showLocModal = ref(false)
const agrForm = reactive({ title:'', type:'CONTRACT', file_url:'' })
const agrFile = ref(null)
const agrFilePreview = ref('')
const agrFileName = ref('')
const agrUploading = ref(false)
const revForm = reactive({ month:'', revenue:0, expense:0 })
const invForm = reactive({ type:'DP', amount:0, description:'' })
const creatingInv = ref(false)
const invSuccess = ref(false)
const lastInvoice = ref(null)
const copied = ref(false)
const showApproveModal = ref(false)
const approveTarget = ref(null)
const approvingInv = ref(false)
const proofFile = ref(null)
const proofPreview = ref('')
const proofFileName = ref('')
const showProofViewer = ref(false)
const proofViewTarget = ref(null)
const selectedStatus = ref('')
const updatingStatus = ref(false)

const statusSteps = [
  { val:'PENDING', label:'Menunggu', progress:0, icon:'ri-time-line', color:'#f59e0b' },
  { val:'DP_VERIFIED', label:'DP Terverifikasi', progress:25, icon:'ri-money-dollar-circle-line', color:'#0ea5e9' },
  { val:'AGREEMENT_SIGNED', label:'Perjanjian', progress:50, icon:'ri-file-text-line', color:'#8b5cf6' },
  { val:'DEVELOPMENT', label:'Pembangunan', progress:75, icon:'ri-building-2-line', color:'#6366f1' },
  { val:'RUNNING', label:'Berjalan', progress:90, icon:'ri-run-line', color:'#22c55e' },
  { val:'COMPLETED', label:'Selesai', progress:100, icon:'ri-checkbox-circle-fill', color:'#10b981' },
]

function statusIdx(s) { return statusSteps.findIndex(st => st.val === s) }
function pipeDotStyle(st, idx) {
  const cur = statusIdx(p.value?.status)
  if (idx === cur) return { background: st.color, color: '#fff', boxShadow: `0 0 0 4px ${st.color}22` }
  if (idx < cur) return { background: '#22c55e', color: '#fff' }
  if (selectedStatus.value === st.val && selectedStatus.value !== p.value?.status) return { background: st.color + '18', color: st.color, border: `2px solid ${st.color}` }
  return {}
}

async function doUpdateStatus() {
  updatingStatus.value = true
  try {
    const step = statusSteps.find(s => s.val === selectedStatus.value)
    await partnershipApi.updateStatus(id, { status: selectedStatus.value, progress_percentage: step?.progress || 0 })
    toast.success('Status berhasil diubah ke ' + statusLabel(selectedStatus.value))
    await loadAll()
    selectedStatus.value = p.value.status
  } catch(e) { toast.error(e.response?.data?.error || 'Gagal mengubah status') }
  finally { updatingStatus.value = false }
}

// Midtrans status check
const checkingInv = reactive({})
const syncing = ref(false)
const hasPendingInvoices = computed(() => invoices.value.some(i => i.status === 'PENDING'))

const creatingLoc = ref(false)
const locForm = reactive({
  nama_lokasi:'', alamat:'', provinsi:'', kota:'', kecamatan:'', kode_pos:'',
  luas_tempat:0, harga_sewa_per_tahun:0, durasi_sewa:1,
  tipe_bangunan:'', lebar_jalan:0, jumlah_lantai:1, estimasi_lalu_lintas:0,
  dekat_dengan:'', jumlah_kompetitor:0, target_market:''
})

function locStatusLabel(s) { return { DRAFT:'Draft', SUBMITTED:'Diajukan', IN_REVIEW:'Ditinjau', SURVEY_SCHEDULED:'Survei Dijadwalkan', SURVEYED:'Sudah Disurvei', APPROVED:'Disetujui', REJECTED:'Ditolak', REVISION_NEEDED:'Perlu Revisi' }[s] || s }
function locScoreColor(s) { if (s >= 80) return '#22c55e'; if (s >= 65) return '#6366f1'; if (s >= 50) return '#f59e0b'; return '#ef4444' }

const activeLocationStatus = computed(() => {
  const active = locations.value.find(l => l.status !== 'REJECTED')
  return active?.status || ''
})
const canAddLocation = computed(() => {
  return !locations.value.some(l => l.status !== 'REJECTED')
})

async function createLoc() {
  creatingLoc.value = true
  try {
    const payload = { ...locForm, mitra_id: p.value.mitra_id, partnership_id: id }
    await locationApi.create(payload)
    toast.success('Lokasi berhasil ditambahkan')
    showLocModal.value = false
    Object.assign(locForm, { nama_lokasi:'', alamat:'', provinsi:'', kota:'', kecamatan:'', kode_pos:'', luas_tempat:0, harga_sewa_per_tahun:0, durasi_sewa:1, tipe_bangunan:'', lebar_jalan:0, jumlah_lantai:1, estimasi_lalu_lintas:0, dekat_dengan:'', jumlah_kompetitor:0, target_market:'' })
    const locR = await locationApi.getByPartnership(id); locations.value = locR.data.data||[]
  } catch(e) { toast.error(e.response?.data?.error || 'Gagal menyimpan lokasi') }
  finally { creatingLoc.value = false }
}

const tabs = computed(() => [
  { key: 'invoices', label: 'Invoice', count: invoices.value.length, icon: '<svg width="15" height="15" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z"/><polyline points="14 2 14 8 20 8"/><line x1="16" y1="13" x2="8" y2="13"/><line x1="16" y1="17" x2="8" y2="17"/></svg>' },
  { key: 'agreements', label: 'Agreement', count: agreements.value.length, icon: '<svg width="15" height="15" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z"/><polyline points="14 2 14 8 20 8"/></svg>' },
  { key: 'revenue', label: 'Revenue', count: revenues.value.length, icon: '<svg width="15" height="15" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><polyline points="22 12 18 12 15 21 9 3 6 12 2 12"/></svg>' },
  { key: 'locations', label: 'Lokasi', count: locations.value.length, icon: '<i class="ri-map-pin-line" style="font-size:15px"></i>' },
])

const gradients = [
  'linear-gradient(135deg, #667eea, #764ba2)',
  'linear-gradient(135deg, #f093fb, #f5576c)',
  'linear-gradient(135deg, #4facfe, #00f2fe)',
  'linear-gradient(135deg, #43e97b, #38f9d7)',
  'linear-gradient(135deg, #fa709a, #fee140)',
]
function avatarBg(name) { return { background: gradients[(name||'?').charCodeAt(0) % gradients.length] } }
function getInitial(name) { return name ? name.split(' ').map(n=>n[0]).join('').substring(0,2).toUpperCase() : '?' }
function formatDate(d) { return d ? new Date(d).toLocaleDateString('id-ID', { day:'numeric', month:'short', year:'numeric' }) : '-' }
function fc(v) { return new Intl.NumberFormat('id-ID',{style:'currency',currency:'IDR',minimumFractionDigits:0}).format(v) }
function statusLabel(s) {
  const map = { PENDING:'Menunggu', DP_VERIFIED:'DP Terverifikasi', AGREEMENT_SIGNED:'Perjanjian Ditandatangani', DEVELOPMENT:'Pembangunan', RUNNING:'Berjalan', COMPLETED:'Selesai' }
  return map[s] || s
}

onMounted(loadAll)

async function loadAll() {
  try {
    const { data } = await partnershipApi.get(id); p.value = data.data; selectedStatus.value = p.value.status
  } catch { toast.error('Gagal memuat data partnership') }
  const [agrR,revR,invR] = await Promise.all([
    agreementApi.byPartnership(id).catch(()=>({data:{data:[]}})),
    revenueApi.byPartnership(id).catch(()=>({data:{data:[]}})),
    invoiceApi.getByPartnership(id).catch(()=>({data:{data:[]}})),
  ])
  agreements.value = agrR.data.data||[]; revenues.value = revR.data.data||[]; invoices.value = invR.data.data||[]
  // Load locations
  try { const locR = await locationApi.getByPartnership(id); locations.value = locR.data.data||[] } catch { locations.value = [] }
  // Auto-sync pending invoices on load
  syncOnLoad()
}

async function syncOnLoad() {
  const pending = invoices.value.filter(i => i.status === 'PENDING' && i.midtrans_order_id)
  if (!pending.length) return
  for (const inv of pending) {
    try {
      const { data } = await invoiceApi.checkStatus(inv.id)
      if (data.synced) {
        const idx = invoices.value.findIndex(i => i.id === inv.id)
        if (idx !== -1) invoices.value[idx] = data.data
      }
    } catch {}
  }
}

async function checkInvoiceStatus(inv) {
  checkingInv[inv.id] = true
  try {
    const { data } = await invoiceApi.checkStatus(inv.id)
    if (data.synced) {
      toast.success(data.message)
      // Update invoice in list
      const idx = invoices.value.findIndex(i => i.id === inv.id)
      if (idx !== -1) invoices.value[idx] = data.data
    } else {
      toast.info(data.message || 'Status belum berubah')
    }
  } catch (e) {
    toast.error('Gagal cek status Midtrans')
  } finally {
    checkingInv[inv.id] = false
  }
}

async function syncAllInvoices() {
  syncing.value = true
  try {
    const { data } = await invoiceApi.syncPending()
    toast.success(data.message || 'Sinkronisasi selesai')
    // Reload all invoices
    const invR = await invoiceApi.getByPartnership(id)
    invoices.value = invR.data.data || []
  } catch (e) {
    toast.error('Gagal menyinkronkan invoice')
  } finally {
    syncing.value = false
  }
}

function isExpiringSoon(inv) {
  if (!inv.expired_at) return false
  const diff = new Date(inv.expired_at) - new Date()
  return diff > 0 && diff < 3 * 60 * 60 * 1000 // < 3 hours
}

function expiryLabel(inv) {
  if (!inv.expired_at) return ''
  const diff = new Date(inv.expired_at) - new Date()
  if (diff <= 0) return 'Kedaluwarsa'
  const hours = Math.floor(diff / (1000 * 60 * 60))
  const mins = Math.floor((diff % (1000 * 60 * 60)) / (1000 * 60))
  if (hours > 0) return `${hours}j ${mins}m tersisa`
  return `${mins} menit tersisa`
}

// File upload helpers
function onAgrFileChange(e) { handleFile(e.target.files[0], 'agr') }
function onAgrFileDrop(e) { handleFile(e.dataTransfer.files[0], 'agr') }

function handleFile(file, type) {
  if (!file) return
  if (file.size > 10 * 1024 * 1024) { toast.error('Ukuran file maksimal 10MB'); return }
  agrFile.value = file
  agrFileName.value = file.name
  agrFilePreview.value = 'file'
}
function clearAgrFile() { agrFile.value = null; agrFilePreview.value = ''; agrFileName.value = '' }
async function createAgreement() {
  if (!agrFile.value) { toast.error('Pilih file terlebih dahulu'); return }
  agrUploading.value = true
  try {
    const { data: upRes } = await uploadApi.upload(agrFile.value)
    const fileUrl = upRes.data?.url || upRes.url || ''
    await agreementApi.create({ partnership_id:id, title: agrForm.title, type: agrForm.type, file_url: fileUrl })
    toast.success('Agreement berhasil diupload')
    showAgrModal.value = false
    clearAgrFile()
    agrForm.title = ''
    agrForm.type = 'CONTRACT'
    await loadAll()
  } catch(e) {
    toast.error(e.response?.data?.error || 'Gagal upload agreement')
  } finally {
    agrUploading.value = false
  }
}
async function signAgreement(aid) {
  try { await agreementApi.sign(aid); toast.success('Agreement berhasil ditandatangani'); await loadAll() }
  catch(e){ toast.error('Gagal menandatangani') }
}
async function createRevenue() {
  try { await revenueApi.create({ partnership_id:id, ...revForm }); toast.success('Revenue berhasil diinput'); showRevModal.value=false; Object.assign(revForm,{month:'',revenue:0,expense:0}); await loadAll() }
  catch(e){ toast.error(e.response?.data?.error || 'Gagal input revenue') }
}

function invBadge(s) {
  return { PAID:'bg-verified', PENDING:'bg-pending', EXPIRED:'bg-draft', FAILED:'bg-draft', CANCELED:'bg-draft' }[s] || 'bg-draft'
}
function invLabel(s) {
  return { PAID:'Lunas', PENDING:'Menunggu', EXPIRED:'Kadaluarsa', FAILED:'Gagal', CANCELED:'Dibatalkan' }[s] || s
}
function invTypeLabel(t) {
  return { DP:'DP', CICILAN:'Cicilan', PELUNASAN:'Pelunasan', INVOICE:'Invoice' }[t] || t || 'Invoice'
}
function onInvTypeChange() {
  const pkgName = p.value?.package?.name || 'Paket Franchise'
  const outletName = p.value?.outlet?.name || ''
  const labels = { DP: `Down Payment ${pkgName}`, CICILAN: `Cicilan ${pkgName}`, PELUNASAN: `Pelunasan ${pkgName}`, INVOICE: `Pembayaran ${pkgName} — ${outletName}` }
  invForm.description = labels[invForm.type] || ''
}
async function createInvoice() {
  creatingInv.value = true
  try {
    const { data } = await invoiceApi.create({ partnership_id: id, type: invForm.type, amount: invForm.amount, description: invForm.description })
    lastInvoice.value = data.data
    invSuccess.value = true
    Object.assign(invForm, { type: 'DP', amount: 0, description: '' })
    await loadAll()
  } catch(e) {
    toast.error(e.response?.data?.error || 'Gagal membuat invoice')
  } finally {
    creatingInv.value = false
  }
}

function openInvModal() {
  invSuccess.value = false
  lastInvoice.value = null
  copied.value = false
  // Pre-fill from package
  invForm.type = 'DP'
  if (p.value?.package) {
    invForm.amount = p.value.package.price || 0
    invForm.description = `Down Payment ${p.value.package.name || 'Paket'}`
  }
  showInvModal.value = true
}

function closeInvModal() {
  showInvModal.value = false
  invSuccess.value = false
  lastInvoice.value = null
  copied.value = false
}

function openApproveModal(inv) {
  approveTarget.value = inv
  approvingInv.value = false
  clearProofFile()
  showApproveModal.value = true
}

function closeApproveModal() {
  showApproveModal.value = false
  approveTarget.value = null
  clearProofFile()
}

async function confirmApprove() {
  if (!approveTarget.value || !proofFile.value) return
  approvingInv.value = true
  try {
    // Upload proof file first
    const { data: upRes } = await uploadApi.upload(proofFile.value)
    const proofUrl = upRes.data?.url || upRes.url || ''
    // Then approve with proof URL
    await invoiceApi.approve(approveTarget.value.id, { proof_url: proofUrl })
    toast.success('Invoice berhasil diverifikasi sebagai LUNAS')
    closeApproveModal()
    await loadAll()
  } catch(e) {
    toast.error(e.response?.data?.error || 'Gagal memverifikasi invoice')
  } finally {
    approvingInv.value = false
  }
}

function onProofFileChange(e) {
  const file = e.target.files[0]
  if (!file) return
  handleProofFile(file)
}
function onProofDrop(e) {
  const file = e.dataTransfer.files[0]
  if (!file) return
  handleProofFile(file)
}
function handleProofFile(file) {
  if (file.size > 10 * 1024 * 1024) { toast.error('Ukuran file maksimal 10MB'); return }
  proofFile.value = file
  proofFileName.value = file.name
  if (file.type.startsWith('image/')) {
    const reader = new FileReader()
    reader.onload = (e) => { proofPreview.value = e.target.result }
    reader.readAsDataURL(file)
  } else {
    proofPreview.value = 'file'
  }
}
function clearProofFile() {
  proofFile.value = null
  proofPreview.value = ''
  proofFileName.value = ''
}

function openProofViewer(inv) {
  proofViewTarget.value = inv
  showProofViewer.value = true
}

function copyLink(url) {
  if (!url) return
  navigator.clipboard.writeText(url).then(() => {
    copied.value = true
    toast.success('Link pembayaran berhasil disalin!')
    setTimeout(() => { copied.value = false }, 2000)
  })
}
</script>

<style scoped>
/* ═══════════════════════════════════════════
   PARTNERSHIP DETAIL — Premium Design System
   ═══════════════════════════════════════════ */

/* ─── PAGE & BACK LINK ─── */
.pd-page { max-width: 100%; }
.pd-back { display: inline-flex; align-items: center; gap: 6px; color: #64748b; font-size: 0.8rem; font-weight: 600; text-decoration: none; margin-bottom: 20px; padding: 6px 12px 6px 8px; border-radius: 8px; transition: all .2s; }
.pd-back:hover { color: #0f172a; background: #f1f5f9; }

/* ─── STATUS PIPELINE ─── */
.pd-pipeline-wrap { background:#fff; border:1px solid #eef1f6; border-radius:16px; padding:20px 24px; margin-bottom:20px; }
.pd-pipeline-header { margin-bottom:16px; }
.pd-pipeline-title { font-size:.88rem; font-weight:700; color:#0f172a; margin:0; display:flex; align-items:center; gap:6px; }
.pd-pipeline-title i { color:#6366f1; font-size:1.05rem; }
.pd-pipeline { display:flex; align-items:flex-start; gap:0; }
.pd-pipe-step { position:relative; display:flex; flex-direction:column; align-items:center; gap:6px; flex:1; padding:8px 4px; cursor:pointer; transition:all .2s; border-radius:10px; }
.pd-pipe-step:hover { background:#f8fafc; }
.pd-pipe-step.active { background:#f0f9ff; }
.pd-pipe-step.selected { background:#faf5ff; }
.pd-pipe-line { position:absolute; top:24px; left:-50%; width:100%; height:2px; background:#e2e8f0; z-index:0; }
.pd-pipe-line.done { background:#22c55e; }
.pd-pipe-dot { position:relative; z-index:1; width:36px; height:36px; border-radius:50%; display:flex; align-items:center; justify-content:center; font-size:.9rem; background:#f1f5f9; color:#94a3b8; flex-shrink:0; transition:all .25s; border:2px solid transparent; }
.pd-pipe-step.active .pd-pipe-dot { transform:scale(1.1); }
.pd-pipe-text { text-align:center; }
.pd-pipe-label { font-size:.68rem; font-weight:700; color:#334155; line-height:1.2; }
.pd-pipe-step.active .pd-pipe-label { color:#0f172a; }
.pd-pipe-step.completed .pd-pipe-label { color:#64748b; }
.pd-pipe-prog { font-size:.58rem; color:#94a3b8; font-weight:600; }
.pd-status-confirm { display:flex; align-items:center; justify-content:space-between; gap:16px; padding:14px 20px; background:linear-gradient(135deg,#f8fafc,#eef2ff); border:1.5px solid #c7d2fe; border-radius:12px; margin-top:12px; }
.pd-sc-info { display:flex; align-items:center; gap:8px; font-size:.82rem; font-weight:600; }
.pd-sc-from { color:#64748b; }
.pd-sc-arrow { color:#6366f1; font-size:1rem; }
.pd-sc-to { color:#6366f1; font-weight:700; }
.pd-sc-btn { display:inline-flex; align-items:center; gap:6px; padding:9px 20px; font-size:.8rem; font-weight:700; border-radius:10px; border:none; cursor:pointer; background:linear-gradient(135deg,#6366f1,#8b5cf6); color:#fff; white-space:nowrap; font-family:inherit; transition:all .15s; }
.pd-sc-btn:hover { box-shadow:0 4px 14px rgba(99,102,241,.3); }
.pd-sc-btn:disabled { opacity:.6; cursor:not-allowed; }
.pd-slide-enter-active { transition:all .3s ease; }
.pd-slide-leave-active { transition:all .2s ease; }
.pd-slide-enter-from { transform:translateY(-8px); opacity:0; }
.pd-slide-leave-to { transform:translateY(-8px); opacity:0; }
.ri-spin { animation:riSpin .8s linear infinite; }
@keyframes riSpin { from{transform:rotate(0deg)} to{transform:rotate(360deg)} }

/* ─── HERO HEADER ─── */
.pd-hero { background: linear-gradient(135deg, #0f0c29 0%, #302b63 50%, #24243e 100%); border-radius: 20px; overflow: hidden; margin-bottom: 20px; position: relative; }
.pd-hero::before { content: ''; position: absolute; inset: 0; background: radial-gradient(ellipse at 80% 20%, rgba(99,102,241,0.15) 0%, transparent 60%); pointer-events: none; }
.pd-hero-bg { display: none; }
.pd-hero-content { display: flex; align-items: center; justify-content: space-between; padding: 28px 32px; position: relative; z-index: 1; }
.pd-hero-left { display: flex; align-items: center; gap: 18px; }
.pd-avatar { width: 56px; height: 56px; border-radius: 14px; display: flex; align-items: center; justify-content: center; color: #fff; font-weight: 800; font-size: 1.05rem; letter-spacing: .5px; box-shadow: 0 4px 16px rgba(0,0,0,0.2); }
.pd-hero-info { display: flex; flex-direction: column; gap: 2px; }
.pd-hero-name { font-size: 1.25rem; font-weight: 800; color: #fff; margin: 0; letter-spacing: -.01em; }
.pd-hero-email { font-size: 0.78rem; color: rgba(255,255,255,.5); }
.pd-hero-phone { font-size: 0.72rem; color: rgba(255,255,255,.3); }
.pd-hero-right { display: flex; align-items: center; gap: 16px; }
.pd-progress-ring-wrap { position: relative; width: 72px; height: 72px; }
.pd-progress-ring { width: 72px; height: 72px; filter: drop-shadow(0 2px 8px rgba(99,102,241,0.3)); }
.pd-progress-value { position: absolute; inset: 0; display: flex; align-items: center; justify-content: center; font-size: 0.92rem; font-weight: 800; color: #fff; letter-spacing: -.02em; }
.pd-status-badge { display: inline-flex; align-items: center; padding: 6px 14px; border-radius: 8px; font-size: 0.68rem; font-weight: 700; text-transform: uppercase; letter-spacing: .05em; white-space: nowrap; }
.st-PENDING { background: rgba(254,243,199,0.9); color: #b45309; }
.st-DP_VERIFIED { background: rgba(224,242,254,0.9); color: #0369a1; }
.st-AGREEMENT_SIGNED { background: rgba(237,233,254,0.9); color: #6d28d9; }
.st-DEVELOPMENT { background: rgba(252,231,243,0.9); color: #be185d; }
.st-RUNNING { background: rgba(220,252,231,0.9); color: #15803d; }
.st-COMPLETED { background: rgba(209,250,229,0.9); color: #047857; }

/* ─── INFO GRID ─── */
.pd-info-grid { display: grid; grid-template-columns: repeat(4, 1fr); gap: 14px; margin-bottom: 20px; }
.pd-info-card { background: #fff; border-radius: 14px; border: 1px solid #eef1f6; padding: 16px; display: flex; gap: 12px; transition: all .25s ease; }
.pd-info-card:hover { box-shadow: 0 6px 24px rgba(0,0,0,0.06); border-color: #e0e7ff; transform: translateY(-2px); }
.pd-info-icon { width: 38px; height: 38px; border-radius: 10px; display: flex; align-items: center; justify-content: center; flex-shrink: 0; }
.ic-affiliator { background: linear-gradient(135deg, #fef9c3, #fef3c7); color: #ca8a04; }
.ic-outlet { background: linear-gradient(135deg, #dbeafe, #eff6ff); color: #2563eb; }
.ic-pkg { background: linear-gradient(135deg, #ede9fe, #f5f3ff); color: #7c3aed; }
.ic-date { background: linear-gradient(135deg, #dcfce7, #f0fdf4); color: #16a34a; }
.pd-info-body { min-width: 0; }
.pd-info-label { font-size: 0.65rem; color: #94a3b8; text-transform: uppercase; letter-spacing: .06em; font-weight: 700; margin-bottom: 4px; }
.pd-info-value { font-size: 0.85rem; font-weight: 700; color: #0f172a; line-height: 1.3; overflow: hidden; text-overflow: ellipsis; white-space: nowrap; }
.pd-info-sub { font-size: 0.72rem; color: #64748b; margin-top: 3px; overflow: hidden; text-overflow: ellipsis; white-space: nowrap; }
.pd-price { color: #059669; font-weight: 800; }

/* ─── TAB BAR ─── */
.pd-tab-bar { display: flex; gap: 4px; background: #fff; border-radius: 14px; padding: 5px; margin-bottom: 16px; border: 1px solid #eef1f6; box-shadow: 0 1px 3px rgba(0,0,0,0.02); }
.pd-tab { flex: 1; display: flex; align-items: center; justify-content: center; gap: 7px; padding: 10px 12px; border-radius: 10px; border: none; background: transparent; color: #64748b; font-size: 0.8rem; font-weight: 600; cursor: pointer; transition: all .2s ease; font-family: inherit; position: relative; }
.pd-tab:hover { color: #334155; background: #f8fafc; }
.pd-tab.active { background: linear-gradient(135deg, #6366f1, #8b5cf6); color: #fff; box-shadow: 0 4px 14px rgba(99,102,241,0.3); }
.pd-tab-icon { display: flex; align-items: center; opacity: .75; }
.pd-tab.active .pd-tab-icon { opacity: 1; }
.pd-tab-count { background: #f1f5f9; color: #64748b; font-size: 0.62rem; padding: 2px 7px; border-radius: 6px; font-weight: 800; min-width: 18px; text-align: center; }
.pd-tab.active .pd-tab-count { background: rgba(255,255,255,0.25); color: #fff; }

/* ─── CONTENT CARD ─── */
.pd-content-card { background: #fff; border-radius: 16px; border: 1px solid #eef1f6; padding: 24px; min-height: 200px; box-shadow: 0 1px 3px rgba(0,0,0,0.02); }
.pd-content-head { display: flex; align-items: center; justify-content: space-between; margin-bottom: 16px; padding-bottom: 14px; border-bottom: 1px solid #f1f5f9; }
.pd-content-head h3 { font-size: 0.95rem; font-weight: 700; color: #0f172a; margin: 0; letter-spacing: -.01em; }
.pd-add-btn { display: inline-flex; align-items: center; gap: 6px; padding: 8px 16px; border-radius: 10px; border: none; background: linear-gradient(135deg, #6366f1, #8b5cf6); color: #fff; font-size: 0.75rem; font-weight: 700; cursor: pointer; transition: all .2s; box-shadow: 0 3px 12px rgba(99,102,241,0.25); font-family: inherit; }
.pd-add-btn:hover { box-shadow: 0 6px 20px rgba(99,102,241,0.35); transform: translateY(-1px); }

/* ─── ITEM LIST ─── */
.pd-items { display: flex; flex-direction: column; gap: 8px; }
.pd-item { display: flex; align-items: center; justify-content: space-between; padding: 14px 16px; background: #fafbfd; border: 1px solid #eef1f6; border-radius: 12px; transition: all .2s ease; }
.pd-item:hover { background: #f1f5f9; border-color: #dde3ed; box-shadow: 0 2px 8px rgba(0,0,0,0.03); }
.pd-item-left { display: flex; align-items: center; gap: 12px; min-width: 0; flex: 1; }
.pd-item-icon { width: 36px; height: 36px; border-radius: 10px; display: flex; align-items: center; justify-content: center; flex-shrink: 0; }
.pi-dp { background: linear-gradient(135deg, #dbeafe, #eff6ff); color: #2563eb; }
.pi-full { background: linear-gradient(135deg, #d1fae5, #f0fdf4); color: #16a34a; }
.pi-agr { background: linear-gradient(135deg, #ede9fe, #f5f3ff); color: #7c3aed; }
.pi-doc { background: linear-gradient(135deg, #e0f2fe, #f0f9ff); color: #0284c7; }
.pd-agr-type-tag { display: inline-flex; align-items: center; padding: 2px 8px; border-radius: 4px; font-size: 0.6rem; font-weight: 700; text-transform: uppercase; letter-spacing: .03em; margin-right: 6px; }
.agt-contract { background: #fef3c7; color: #b45309; }
.agt-document { background: #e0f2fe; color: #0369a1; }
.pd-item-title { font-size: 0.82rem; font-weight: 700; color: #0f172a; line-height: 1.3; }
.pd-item-date { font-size: 0.7rem; color: #94a3b8; margin-top: 2px; line-height: 1.3; }
.pd-item-right { display: flex; align-items: center; gap: 10px; flex-shrink: 0; }
.pd-item-amount { font-size: 0.88rem; font-weight: 800; color: #059669; letter-spacing: -.01em; white-space: nowrap; }
.pd-item-actions { display: flex; align-items: center; gap: 6px; flex-wrap: wrap; }
.pd-badge { display: inline-flex; align-items: center; padding: 4px 10px; border-radius: 6px; font-size: 0.65rem; font-weight: 700; text-transform: uppercase; letter-spacing: .04em; white-space: nowrap; }
.bg-verified { background: #d1fae5; color: #065f46; }
.bg-pending { background: #fef3c7; color: #92400e; }
.bg-draft { background: #f1f5f9; color: #475569; }
.pd-verify-btn { display: inline-flex; align-items: center; gap: 4px; padding: 5px 12px; border-radius: 8px; border: 1px solid #c7d2fe; background: #eef2ff; color: #4f46e5; font-size: 0.7rem; font-weight: 700; cursor: pointer; transition: all .2s; text-decoration: none; white-space: nowrap; font-family: inherit; }
.pd-verify-btn:hover { background: #e0e7ff; border-color: #a5b4fc; box-shadow: 0 2px 8px rgba(99,102,241,0.15); }

/* ─── INVOICE TYPE TAG ─── */
.pd-inv-type-tag { display: inline-flex; padding: 2px 7px; border-radius: 4px; font-size: 0.58rem; font-weight: 800; text-transform: uppercase; letter-spacing: .05em; margin-right: 6px; vertical-align: middle; }
.invt-DP { background: #dbeafe; color: #1d4ed8; }
.invt-CICILAN { background: #ede9fe; color: #6d28d9; }
.invt-PELUNASAN { background: #dcfce7; color: #15803d; }
.invt-INVOICE { background: #f1f5f9; color: #475569; }

/* ─── APPROVE BUTTON & MODAL ─── */
.pd-approve-btn { display: inline-flex; align-items: center; gap: 4px; padding: 5px 12px; border-radius: 8px; border: 1px solid #86efac; background: #f0fdf4; color: #15803d; font-size: 0.7rem; font-weight: 700; cursor: pointer; transition: all .2s; white-space: nowrap; font-family: inherit; }
.pd-approve-btn:hover { background: #dcfce7; border-color: #4ade80; box-shadow: 0 2px 8px rgba(34,197,94,0.15); }

/* ─── MIDTRANS STATUS CHECK ─── */
.pd-content-head-actions { display: flex; align-items: center; gap: 8px; }
.pd-sync-btn { display: inline-flex; align-items: center; gap: 6px; padding: 7px 14px; border-radius: 8px; border: 1px solid #99f6e4; background: #f0fdfa; color: #0d9488; font-size: 0.72rem; font-weight: 700; cursor: pointer; transition: all .2s; white-space: nowrap; font-family: inherit; }
.pd-sync-btn:hover { background: #ccfbf1; border-color: #5eead4; box-shadow: 0 2px 8px rgba(20,184,166,0.15); }
.pd-sync-btn:disabled { opacity: 0.6; cursor: not-allowed; }
.pd-check-btn { display: inline-flex; align-items: center; gap: 4px; padding: 5px 10px; border-radius: 8px; border: 1px solid #fde68a; background: #fffbeb; color: #b45309; font-size: 0.68rem; font-weight: 700; cursor: pointer; transition: all .2s; white-space: nowrap; font-family: inherit; }
.pd-check-btn:hover { background: #fef3c7; border-color: #fbbf24; box-shadow: 0 2px 8px rgba(245,158,11,0.15); }
.pd-check-btn:disabled { opacity: 0.6; cursor: not-allowed; }
.pd-expiry-hint { font-size: 0.65rem; color: #94a3b8; font-weight: 600; }
.pd-expiry-soon { color: #ef4444; }
@keyframes spin { to { transform: rotate(360deg); } }
.spin-icon { animation: spin 1s linear infinite; }

.pd-modal-confirm { max-width: 420px; text-align: center; overflow: hidden; }
.pd-confirm-body { padding: 32px 28px 20px; }
.pd-confirm-icon { width: 64px; height: 64px; border-radius: 50%; background: linear-gradient(135deg, #dcfce7, #bbf7d0); color: #16a34a; display: flex; align-items: center; justify-content: center; margin: 0 auto 16px; }
.pd-confirm-title { font-size: 1.05rem; font-weight: 800; color: #0f172a; margin: 0 0 8px; letter-spacing: -.01em; }
.pd-confirm-desc { font-size: 0.8rem; color: #64748b; margin: 0 0 18px; line-height: 1.55; }
.pd-confirm-desc strong { color: #0f172a; }
.pd-confirm-invoice { background: #f8fafc; border: 1px solid #eef1f6; border-radius: 12px; padding: 14px 16px; text-align: left; }
.pd-confirm-inv-row { display: flex; align-items: center; justify-content: space-between; padding: 5px 0; }
.pd-confirm-inv-row + .pd-confirm-inv-row { border-top: 1px solid #f1f5f9; }
.pd-confirm-inv-label { font-size: 0.72rem; color: #94a3b8; font-weight: 600; }
.pd-confirm-inv-val { font-size: 0.82rem; font-weight: 700; color: #0f172a; }
.pd-confirm-inv-amount { font-size: 0.92rem; font-weight: 800; color: #059669; }
.pd-confirm-upload { margin-top: 16px; text-align: left; }
.pd-confirm-upload-label { display: flex; align-items: center; gap: 6px; font-size: 0.78rem; font-weight: 700; color: #334155; margin-bottom: 8px; }
.pd-confirm-upload-label .req { color: #ef4444; }
.pd-confirm-footer { display: flex; gap: 8px; padding: 16px 28px 24px; justify-content: center; }
.pd-btn-approve { display: inline-flex; align-items: center; gap: 6px; padding: 11px 24px; border-radius: 12px; border: none; background: linear-gradient(135deg, #22c55e, #16a34a); color: #fff; font-size: 0.82rem; font-weight: 700; cursor: pointer; transition: all .2s; box-shadow: 0 4px 14px rgba(34,197,94,0.3); font-family: inherit; }
.pd-btn-approve:hover { box-shadow: 0 6px 20px rgba(34,197,94,0.4); transform: translateY(-1px); }
.pd-btn-approve:disabled { opacity: .5; cursor: not-allowed; transform: none; box-shadow: none; }
.pd-btn-spinner { width: 16px; height: 16px; border: 2px solid rgba(255,255,255,0.3); border-top-color: #fff; border-radius: 50%; animation: pd-spin 0.7s linear infinite; }

/* ─── PAYMENT METHOD LABEL ─── */
.pd-pay-method { display: inline-flex; align-items: center; gap: 3px; padding: 3px 8px; border-radius: 5px; font-size: 0.6rem; font-weight: 700; text-transform: uppercase; letter-spacing: .04em; white-space: nowrap; }
.pm-auto { background: #dcfce7; color: #15803d; }
.pm-manual { background: #fef3c7; color: #92400e; }

/* ─── PROOF BUTTON ─── */
.pd-proof-btn { display: inline-flex; align-items: center; gap: 4px; padding: 4px 10px; border-radius: 6px; border: 1px solid #e0e7ff; background: #f5f3ff; color: #6d28d9; font-size: 0.65rem; font-weight: 700; cursor: pointer; transition: all .2s; white-space: nowrap; font-family: inherit; }
.pd-proof-btn:hover { background: #ede9fe; border-color: #c4b5fd; box-shadow: 0 2px 6px rgba(109,40,217,0.12); }

/* ─── PROOF VIEWER MODAL ─── */
.pd-modal-proof { max-width: 560px; overflow: hidden; }
.pd-proof-header { display: flex; align-items: center; gap: 12px; padding: 18px 24px; border-bottom: 1px solid #f1f5f9; }
.pd-proof-header-info { flex: 1; min-width: 0; }
.pd-proof-header-info h3 { font-size: 0.95rem; font-weight: 700; color: #0f172a; margin: 0 0 4px; }
.pd-proof-meta { display: flex; align-items: center; gap: 6px; font-size: 0.75rem; font-weight: 600; color: #64748b; flex-wrap: wrap; }
.pd-proof-meta-sep { color: #cbd5e1; }
.pd-proof-meta-amount { color: #059669; font-weight: 800; }
.pd-proof-viewer-body { padding: 16px; background: #f8fafc; display: flex; align-items: center; justify-content: center; min-height: 200px; }
.pd-proof-image { max-width: 100%; max-height: 400px; border-radius: 10px; object-fit: contain; box-shadow: 0 4px 20px rgba(0,0,0,0.08); }
.pd-proof-viewer-foot { display: flex; align-items: center; gap: 10px; padding: 14px 24px; border-top: 1px solid #f1f5f9; flex-wrap: wrap; }
.pd-proof-paid-date { font-size: 0.72rem; color: #94a3b8; font-weight: 500; }

/* ─── REVENUE CARDS ─── */
.pd-rev-card { background: #fafbfd; border: 1px solid #eef1f6; border-radius: 12px; padding: 16px; margin-bottom: 8px; transition: all .2s ease; }
.pd-rev-card:hover { background: #f1f5f9; border-color: #dde3ed; }
.pd-rev-month { font-size: 0.82rem; font-weight: 700; color: #334155; margin-bottom: 12px; padding-bottom: 10px; border-bottom: 1px solid #eef1f6; }
.pd-rev-grid { display: grid; grid-template-columns: repeat(4, 1fr); gap: 10px; }
.pd-rev-item { display: flex; flex-direction: column; gap: 3px; }
.pd-rev-label { font-size: 0.62rem; color: #94a3b8; text-transform: uppercase; letter-spacing: .05em; font-weight: 700; }
.pd-rev-val { font-size: 0.85rem; font-weight: 800; }
.rv-green { color: #059669; }
.rv-red { color: #dc2626; }
.rv-profit { color: #2563eb; }
.rv-share { color: #7c3aed; }

/* ─── EMPTY STATE ─── */
.pd-empty { display: flex; flex-direction: column; align-items: center; gap: 12px; padding: 48px 20px; color: #cbd5e1; }
.pd-empty svg { opacity: .5; }
.pd-empty p { font-size: 0.82rem; margin: 0; color: #94a3b8; font-weight: 500; }

/* ─── LOADING ─── */
.pd-loading { display: flex; flex-direction: column; align-items: center; justify-content: center; height: 300px; gap: 16px; color: #94a3b8; }
.pd-spinner { width: 32px; height: 32px; border: 3px solid #eef1f6; border-top-color: #6366f1; border-radius: 50%; animation: pd-spin 0.7s linear infinite; }
@keyframes pd-spin { to { transform: rotate(360deg); } }

/* ─── MODAL ─── */
.pd-overlay { position: fixed; inset: 0; background: rgba(15,23,42,.55); backdrop-filter: blur(6px); display: flex; align-items: center; justify-content: center; z-index: 1000; animation: pd-fadeIn .15s ease-out; }
@keyframes pd-fadeIn { from { opacity: 0; } to { opacity: 1; } }
.pd-modal { background: #fff; border-radius: 20px; width: 100%; max-width: 480px; box-shadow: 0 24px 80px rgba(0,0,0,.18), 0 0 0 1px rgba(0,0,0,.04); animation: pd-slideUp .2s ease-out; }
@keyframes pd-slideUp { from { transform: translateY(16px); opacity: 0; } to { transform: translateY(0); opacity: 1; } }
.pd-modal-head { display: flex; align-items: center; gap: 12px; padding: 20px 24px; border-bottom: 1px solid #f1f5f9; }
.pd-modal-icon { width: 38px; height: 38px; border-radius: 10px; display: flex; align-items: center; justify-content: center; }
.mi-pay { background: linear-gradient(135deg, #dbeafe, #eff6ff); color: #2563eb; }
.mi-agr { background: linear-gradient(135deg, #ede9fe, #f5f3ff); color: #7c3aed; }
.mi-rev { background: linear-gradient(135deg, #dcfce7, #f0fdf4); color: #16a34a; }
.pd-modal-head h3 { flex: 1; font-size: 1rem; font-weight: 700; color: #0f172a; margin: 0; }
.pd-modal-x { width: 32px; height: 32px; border-radius: 8px; display: flex; align-items: center; justify-content: center; background: none; border: none; color: #94a3b8; font-size: 1.3rem; cursor: pointer; transition: all .15s; }
.pd-modal-x:hover { background: #f1f5f9; color: #0f172a; }
.pd-modal-body { padding: 20px 24px; }
.pd-fg { margin-bottom: 16px; }
.pd-fg label { display: block; font-size: 0.78rem; font-weight: 600; margin-bottom: 6px; color: #334155; }
.pd-fg .req { color: #ef4444; }
.pd-input { width: 100%; padding: 10px 14px; border-radius: 10px; border: 1.5px solid #e2e8f0; background: #fafbfc; color: #1e293b; font-size: 0.82rem; outline: none; transition: all .2s; box-sizing: border-box; font-family: inherit; appearance: none; -webkit-appearance: none; }
.pd-input:focus { border-color: #6366f1; background: #fff; box-shadow: 0 0 0 3px rgba(99,102,241,0.08); }
select.pd-input { background-image: url("data:image/svg+xml,%3csvg xmlns='http://www.w3.org/2000/svg' fill='none' viewBox='0 0 20 20'%3e%3cpath stroke='%236b7280' stroke-linecap='round' stroke-linejoin='round' stroke-width='1.5' d='M6 8l4 4 4-4'/%3e%3c/svg%3e"); background-position: right 12px center; background-repeat: no-repeat; background-size: 14px; padding-right: 34px; cursor: pointer; }
.pd-form-row { display: grid; grid-template-columns: 1fr 1fr; gap: 12px; }
.pd-modal-foot { display: flex; gap: 8px; justify-content: flex-end; padding-top: 16px; border-top: 1px solid #f8fafc; }
.pd-btn-sec { padding: 10px 20px; border-radius: 10px; font-size: .8rem; font-weight: 600; background: #f1f5f9; color: #475569; border: none; cursor: pointer; transition: all .15s; font-family: inherit; }
.pd-btn-sec:hover { background: #e2e8f0; }
.pd-btn-primary { display: inline-flex; align-items: center; gap: 6px; padding: 10px 22px; border-radius: 10px; border: none; background: linear-gradient(135deg, #6366f1, #8b5cf6); color: #fff; font-size: 0.8rem; font-weight: 700; cursor: pointer; transition: all .2s; box-shadow: 0 3px 12px rgba(99,102,241,0.25); text-decoration: none; font-family: inherit; }
.pd-btn-primary:hover { box-shadow: 0 6px 20px rgba(99,102,241,0.35); transform: translateY(-1px); }
.pd-btn-primary:disabled { opacity: .5; cursor: not-allowed; transform: none; box-shadow: none; }

/* ─── UPLOAD ZONE ─── */
.pd-upload-zone { border: 2px dashed #e2e8f0; border-radius: 12px; padding: 24px 16px; cursor: pointer; transition: all .2s; text-align: center; background: #fafbfc; }
.pd-upload-zone:hover { border-color: #818cf8; background: #f5f3ff; }
.pd-upload-placeholder { display: flex; flex-direction: column; align-items: center; gap: 8px; color: #94a3b8; font-size: 0.78rem; font-weight: 600; }
.pd-upload-hint { font-size: 0.68rem; color: #cbd5e1; font-weight: 500; }
.pd-upload-preview { position: relative; display: flex; align-items: center; justify-content: center; min-height: 60px; }
.pd-upload-preview img { max-height: 110px; max-width: 100%; border-radius: 8px; object-fit: contain; box-shadow: 0 2px 8px rgba(0,0,0,0.08); }
.pd-file-name { font-size: 0.82rem; font-weight: 700; color: #334155; }
.pd-upload-clear { position: absolute; top: -8px; right: -8px; width: 22px; height: 22px; border-radius: 50%; border: 2px solid #fff; background: #ef4444; color: #fff; font-size: 0.85rem; cursor: pointer; display: flex; align-items: center; justify-content: center; line-height: 1; transition: all .15s; box-shadow: 0 2px 6px rgba(239,68,68,0.3); }
.pd-upload-clear:hover { background: #dc2626; transform: scale(1.1); }

/* ─── INVOICE MODAL ─── */
.pd-modal-lg { max-width: 520px; }
.pd-inv-pkg-card { display: flex; gap: 14px; padding: 14px 16px; background: linear-gradient(135deg, #f0f4ff 0%, #faf5ff 100%); border-radius: 12px; border: 1px solid #e0e7ff; margin-bottom: 16px; }
.pd-inv-pkg-icon { width: 42px; height: 42px; border-radius: 12px; background: linear-gradient(135deg, #6366f1, #8b5cf6); color: #fff; display: flex; align-items: center; justify-content: center; flex-shrink: 0; box-shadow: 0 3px 10px rgba(99,102,241,0.25); }
.pd-inv-pkg-info { flex: 1; min-width: 0; }
.pd-inv-pkg-label { font-size: 0.62rem; color: #94a3b8; text-transform: uppercase; letter-spacing: .06em; font-weight: 700; margin-bottom: 3px; }
.pd-inv-pkg-name { font-size: 0.88rem; font-weight: 700; color: #0f172a; margin-bottom: 2px; overflow: hidden; text-overflow: ellipsis; white-space: nowrap; }
.pd-inv-pkg-price { font-size: 0.8rem; font-weight: 800; color: #059669; }
.pd-inv-pkg-mitra { font-size: 0.72rem; color: #64748b; margin-top: 2px; }

.pd-inv-success { display: flex; align-items: center; gap: 16px; padding: 18px 20px; background: linear-gradient(135deg, #f0fdf4, #ecfdf5); border-radius: 12px; border: 1px solid #bbf7d0; margin-bottom: 16px; }
.pd-inv-success-icon { flex-shrink: 0; }
.pd-inv-number { font-size: 0.92rem; font-weight: 800; color: #059669; letter-spacing: -.01em; }
.pd-inv-amt { font-size: 1.2rem; font-weight: 800; color: #0f172a; margin: 2px 0; letter-spacing: -.02em; }
.pd-inv-desc { font-size: 0.78rem; color: #64748b; }

.pd-inv-link-box { margin-bottom: 14px; }
.pd-inv-link-box label { display: block; font-size: 0.75rem; font-weight: 700; color: #475569; margin-bottom: 6px; }
.pd-inv-link-row { display: flex; gap: 6px; }
.pd-inv-link-row .pd-input { font-size: 0.75rem; color: #6366f1; background: #f8f7ff; cursor: pointer; font-weight: 600; }
.pd-copy-btn { display: inline-flex; align-items: center; gap: 4px; padding: 10px 14px; border-radius: 10px; border: none; background: linear-gradient(135deg, #6366f1, #8b5cf6); color: #fff; font-size: 0.75rem; font-weight: 700; cursor: pointer; white-space: nowrap; transition: all .2s; box-shadow: 0 3px 10px rgba(99,102,241,0.2); font-family: inherit; }
.pd-copy-btn:hover { box-shadow: 0 4px 16px rgba(99,102,241,0.35); }
.pd-inv-link-hint { font-size: 0.68rem; color: #94a3b8; margin-top: 6px; font-weight: 500; font-style: italic; }

/* ─── RESPONSIVE ─── */
@media (max-width: 1024px) {
  .pd-info-grid { grid-template-columns: repeat(2, 1fr); }
}
@media (max-width: 768px) {
  .pd-hero-content { flex-direction: column; gap: 20px; align-items: flex-start; padding: 24px 20px; }
  .pd-hero-right { width: 100%; justify-content: flex-start; }
  .pd-info-grid { grid-template-columns: 1fr 1fr; gap: 10px; }
  .pd-rev-grid { grid-template-columns: repeat(2, 1fr); }
  .pd-item { flex-direction: column; align-items: flex-start; gap: 10px; }
  .pd-item-right { width: 100%; justify-content: space-between; }
  .pd-content-card { padding: 16px; }
  .pd-modal { margin: 16px; max-width: calc(100% - 32px); }
  .pd-tab { font-size: 0.72rem; padding: 8px 6px; gap: 4px; }
  .pd-tab-count { font-size: 0.58rem; padding: 1px 5px; }
}
@media (max-width: 480px) {
  .pd-info-grid { grid-template-columns: 1fr; }
  .pd-rev-grid { grid-template-columns: 1fr 1fr; }
  .pd-content-head { flex-direction: column; gap: 10px; align-items: flex-start; }
}

/* ═══ LOCATION CARDS ═══ */
.pd-loc-info-badge { display:inline-flex; align-items:center; gap:5px; font-size:.75rem; font-weight:600; color:#6366f1; background:#eef2ff; padding:6px 14px; border-radius:8px; border:1px solid #c7d2fe; }
.pd-loc-card { background:#f8fafc; border:1px solid #eef1f6; border-radius:12px; padding:16px 18px; margin-bottom:10px; cursor:pointer; transition:all .2s; }
.pd-loc-card:hover { border-color:#c7d2fe; box-shadow:0 2px 10px rgba(99,102,241,.06); }
.pd-loc-head { display:flex; justify-content:space-between; align-items:center; margin-bottom:8px; }
.pd-loc-left { display:flex; gap:8px; align-items:center; }
.pd-loc-badge { font-size:.65rem; font-weight:700; text-transform:uppercase; letter-spacing:.04em; padding:3px 8px; border-radius:5px; }
.lbs-DRAFT { background:#f1f5f9; color:#475569; }
.lbs-SUBMITTED { background:#e0f2fe; color:#0284c7; }
.lbs-IN_REVIEW { background:#fef3c7; color:#92400e; }
.lbs-SURVEY_SCHEDULED { background:#e0e7ff; color:#4338ca; }
.lbs-SURVEYED { background:#ede9fe; color:#6d28d9; }
.lbs-APPROVED { background:#dcfce7; color:#15803d; }
.lbs-REJECTED { background:#fee2e2; color:#b91c1c; }
.lbs-REVISION_NEEDED { background:#ffedd5; color:#c2410c; }
.pd-loc-type { font-size:.65rem; font-weight:600; color:#475569; background:#e2e8f0; padding:3px 8px; border-radius:5px; display:flex; align-items:center; gap:3px; text-transform:capitalize; }
.pd-loc-type i { font-size:.7rem; }
.pd-loc-score { display:flex; align-items:center; gap:4px; font-size:.78rem; font-weight:700; }
.pd-loc-score i { font-size:.85rem; }
.pd-loc-score span { font-weight:500; opacity:.7; font-size:.68rem; }
.pd-loc-name { font-size:.95rem; font-weight:700; color:#0f172a; margin-bottom:8px; }
.pd-loc-meta { display:flex; gap:14px; flex-wrap:wrap; font-size:.72rem; color:#64748b; }
.pd-loc-meta span { display:flex; align-items:center; gap:4px; }
.pd-loc-meta i { color:#94a3b8; font-size:.8rem; }

/* Location form sections */
.pd-loc-section { background:#f8fafc; border:1px solid #eef1f6; border-radius:12px; padding:16px 18px; margin-bottom:14px; }
.pd-loc-sh { display:flex; align-items:center; gap:8px; font-size:.82rem; font-weight:700; color:#0f172a; margin-bottom:12px; }
.pd-loc-sh i { font-size:1rem; }

.ri-spin { animation:riSpin .8s linear infinite; }
@keyframes riSpin { from{transform:rotate(0deg)} to{transform:rotate(360deg)} }
</style>
