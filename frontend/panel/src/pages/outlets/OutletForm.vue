<template>
  <div class="animate-in">
    <!-- Hero header -->
    <div class="form-hero">
      <router-link to="/outlets" class="form-hero-back">
        <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor"><path d="M19 12H5m7-7-7 7 7 7" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/></svg>
      </router-link>
      <div>
        <h1 class="form-hero-title">{{ isEdit ? 'Edit Outlet' : 'Tambah Outlet Baru' }}</h1>
        <p class="form-hero-sub">{{ isEdit ? 'Perbarui informasi outlet' : 'Lengkapi data outlet di bawah ini' }}</p>
      </div>
    </div>

    <!-- Tabs -->
    <div class="form-tabs">
      <button v-for="(tab, i) in tabs" :key="i" type="button"
        class="form-tab" :class="{ active: activeTab === i }" @click="activeTab = i">
        <span class="form-tab-icon" v-html="tab.icon"></span>
        <span class="form-tab-label">{{ tab.label }}</span>
        <span v-if="tab.badge" class="form-tab-badge">{{ tab.badge }}</span>
      </button>
    </div>

    <form @submit.prevent="handleSubmit" novalidate>
      <div class="form-layout">
        <!-- Main left column -->
        <div class="form-main">

          <!-- Tab 0: Informasi Dasar -->
          <div v-show="activeTab === 0" class="form-tab-panel">
            <div class="form-row-2">
              <div class="form-field">
                <label class="form-field-label">Nama Outlet <span class="req">*</span></label>
                <input v-model="form.name" type="text" class="form-field-input" :class="{ 'is-error': errors.name }"
                  placeholder="Contoh: Kebab Turki Baba Rafi" @blur="validateField('name')" />
                <div v-if="errors.name" class="form-field-error">{{ errors.name }}</div>
              </div>
              <div class="form-field">
                <label class="form-field-label">Kategori <span class="req">*</span></label>
                <select v-model="form.category_id" class="form-field-input" :class="{ 'is-error': errors.category_id }" @change="validateField('category_id')">
                  <option value="">Pilih kategori...</option>
                  <option v-for="cat in categories" :key="cat.id" :value="cat.id">{{ cat.name }}</option>
                </select>
                <div v-if="errors.category_id" class="form-field-error">{{ errors.category_id }}</div>
              </div>
            </div>
            <div class="form-field">
              <label class="form-field-label">Deskripsi Singkat</label>
              <input v-model="form.short_description" type="text" class="form-field-input"
                placeholder="Ringkasan singkat outlet (maks 500 karakter)" maxlength="500" />
              <div class="form-field-counter">{{ (form.short_description || '').length }}/500</div>
            </div>
            <div class="form-field">
              <label class="form-field-label">Deskripsi Detail <span class="req">*</span></label>
              <div class="rte-wrapper" :class="{ 'is-error': errors.description }">
                <div class="rte-toolbar">
                  <button type="button" @click="execCmd('bold')" class="rte-btn" title="Bold"><strong>B</strong></button>
                  <button type="button" @click="execCmd('italic')" class="rte-btn" title="Italic"><em>I</em></button>
                  <button type="button" @click="execCmd('underline')" class="rte-btn" title="Underline"><u>U</u></button>
                  <div class="rte-sep"></div>
                  <button type="button" @click="execCmd('insertUnorderedList')" class="rte-btn" title="Bullet List">
                    <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor"><line x1="8" y1="6" x2="21" y2="6" stroke-width="2"/><line x1="8" y1="12" x2="21" y2="12" stroke-width="2"/><line x1="8" y1="18" x2="21" y2="18" stroke-width="2"/><line x1="3" y1="6" x2="3.01" y2="6" stroke-width="2" stroke-linecap="round"/><line x1="3" y1="12" x2="3.01" y2="12" stroke-width="2" stroke-linecap="round"/><line x1="3" y1="18" x2="3.01" y2="18" stroke-width="2" stroke-linecap="round"/></svg>
                  </button>
                  <button type="button" @click="execCmd('insertOrderedList')" class="rte-btn" title="Numbered List">
                    <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor"><line x1="10" y1="6" x2="21" y2="6" stroke-width="2"/><line x1="10" y1="12" x2="21" y2="12" stroke-width="2"/><line x1="10" y1="18" x2="21" y2="18" stroke-width="2"/><path d="M4 6h1v4M3 10h3M4 14v4h3" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round"/></svg>
                  </button>
                  <div class="rte-sep"></div>
                  <select @change="execHeading($event)" class="rte-select">
                    <option value="">Heading</option>
                    <option value="h2">H2</option>
                    <option value="h3">H3</option>
                    <option value="p">Normal</option>
                  </select>
                  <button type="button" @click="execCmd('removeFormat')" class="rte-btn" title="Clear">
                    <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor"><line x1="18" y1="6" x2="6" y2="18" stroke-width="2" stroke-linecap="round"/><line x1="6" y1="6" x2="18" y2="18" stroke-width="2" stroke-linecap="round"/></svg>
                  </button>
                </div>
                <div ref="editorRef" class="rte-content" contenteditable="true"
                  @input="onEditorInput" @blur="validateField('description')"></div>
              </div>
              <div v-if="errors.description" class="form-field-error">{{ errors.description }}</div>
            </div>
          </div>

          <!-- Tab 1: Investasi -->
          <div v-show="activeTab === 1" class="form-tab-panel">
            <div class="form-row-2">
              <div class="form-field">
                <label class="form-field-label">Harga Mulai Dari <span class="req">*</span></label>
                <div class="input-with-prefix">
                  <span class="input-prefix">Rp</span>
                  <input :value="formatIdrInput(form.minimum_investment)" type="text" class="form-field-input has-prefix"
                    :class="{ 'is-error': errors.minimum_investment }"
                    placeholder="10.000.000" @input="form.minimum_investment = parseIdrInput($event)" @blur="validateField('minimum_investment')" />
                </div>
                <div v-if="errors.minimum_investment" class="form-field-error">{{ errors.minimum_investment }}</div>
              </div>
            </div>
            <div class="form-row-2">
              <div class="form-field">
                <label class="form-field-label">Profit Sharing</label>
                <div class="input-with-suffix">
                  <input v-model.number="form.profit_sharing_percentage" type="number" step="0.1" class="form-field-input has-suffix"
                    :class="{ 'is-error': errors.profit_sharing_percentage }"
                    placeholder="0" @blur="validateField('profit_sharing_percentage')" />
                  <span class="input-suffix">%</span>
                </div>
                <div v-if="errors.profit_sharing_percentage" class="form-field-error">{{ errors.profit_sharing_percentage }}</div>
              </div>
              <div class="form-field">
                <label class="form-field-label">Estimasi ROI</label>
                <input v-model="form.estimated_roi" type="text" class="form-field-input" placeholder="Contoh: 6-12 bulan" />
              </div>
            </div>
          </div>

          <!-- Tab 2: Lokasi -->
          <div v-show="activeTab === 2" class="form-tab-panel">
            <div class="form-row-2">
              <div class="form-field">
                <label class="form-field-label">Kota <span class="req">*</span></label>
                <input v-model="form.city" type="text" class="form-field-input" :class="{ 'is-error': errors.city }"
                  placeholder="Contoh: Jakarta Selatan" @blur="validateField('city')" />
                <div v-if="errors.city" class="form-field-error">{{ errors.city }}</div>
              </div>
              <div class="form-field">
                <label class="form-field-label">Provinsi <span class="req">*</span></label>
                <input v-model="form.province" type="text" class="form-field-input" :class="{ 'is-error': errors.province }"
                  placeholder="Contoh: DKI Jakarta" @blur="validateField('province')" />
                <div v-if="errors.province" class="form-field-error">{{ errors.province }}</div>
              </div>
            </div>
            <div class="form-field">
              <label class="form-field-label">Alamat Lengkap</label>
              <textarea v-model="form.address" class="form-field-input" rows="2" placeholder="Alamat lengkap outlet"></textarea>
            </div>
            <div class="form-field">
              <label class="form-field-label">Persyaratan Lokasi</label>
              <textarea v-model="form.location_requirement" class="form-field-input" rows="2" placeholder="Contoh: Lahan min 3x3m, dekat jalan utama, area ramai"></textarea>
            </div>
          </div>

          <!-- Tab 3: Kontak -->
          <div v-show="activeTab === 3" class="form-tab-panel">
            <div class="form-row-2">
              <div class="form-field">
                <label class="form-field-label">Telepon</label>
                <input v-model="form.contact_phone" type="tel" class="form-field-input" placeholder="08xxxxxxxxxx" />
              </div>
              <div class="form-field">
                <label class="form-field-label">WhatsApp</label>
                <input v-model="form.contact_whatsapp" type="tel" class="form-field-input" placeholder="08xxxxxxxxxx" />
              </div>
            </div>
            <div class="form-row-2">
              <div class="form-field">
                <label class="form-field-label">Email</label>
                <input v-model="form.contact_email" type="email" class="form-field-input"
                  :class="{ 'is-error': errors.contact_email }"
                  placeholder="info@outlet.com" @blur="validateField('contact_email')" />
                <div v-if="errors.contact_email" class="form-field-error">{{ errors.contact_email }}</div>
              </div>
              <div class="form-field">
                <label class="form-field-label">Website</label>
                <input v-model="form.website" type="url" class="form-field-input"
                  :class="{ 'is-error': errors.website }"
                  placeholder="https://www.outlet.com" @blur="validateField('website')" />
                <div v-if="errors.website" class="form-field-error">{{ errors.website }}</div>
              </div>
            </div>
          </div>

          <!-- Tab 4: Paket -->
          <div v-show="activeTab === 4" class="form-tab-panel">
            <!-- Package list -->
            <div v-for="(pkg, idx) in packages" :key="pkg._key || pkg.id" class="pkg-card">
              <div class="pkg-card-header">
                <div>
                  <div class="pkg-card-name">{{ pkg.name }}</div>
                  <div class="pkg-card-price">{{ formatCurrency(pkg.price) }}</div>
                </div>
                <div class="pkg-card-actions">
                  <button type="button" @click="editPackage(idx)" class="pkg-act-btn pkg-act-edit" title="Edit">
                    <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor"><path d="M11 4H4a2 2 0 0 0-2 2v14a2 2 0 0 0 2 2h14a2 2 0 0 0 2-2v-7" stroke-width="2"/><path d="M18.5 2.5a2.121 2.121 0 0 1 3 3L12 15l-4 1 1-4 9.5-9.5z" stroke-width="2"/></svg>
                  </button>
                  <button type="button" @click="removePackage(idx)" class="pkg-act-btn pkg-act-del" title="Hapus">
                    <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor"><polyline points="3 6 5 6 21 6" stroke-width="2"/><path d="M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6m3 0V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2" stroke-width="2"/></svg>
                  </button>
                </div>
              </div>
              <div v-if="pkg.minimum_dp || pkg.duration || pkg.estimated_bep || pkg.net_profit" class="pkg-card-meta">
                <span v-if="pkg.minimum_dp">DP: {{ formatCurrency(pkg.minimum_dp) }}</span>
                <span v-if="pkg.duration">Durasi: {{ pkg.duration }}</span>
                <span v-if="pkg.estimated_bep">BEP: {{ pkg.estimated_bep }}</span>
                <span v-if="pkg.net_profit">Profit: {{ pkg.net_profit }}</span>
              </div>
              <div v-if="pkg.benefits && pkg.benefits.length" class="pkg-card-benefits">
                <div v-for="(b, bi) in pkg.benefits" :key="bi" class="pkg-benefit-item">
                  <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="#16a34a" stroke-width="2.5"><polyline points="20 6 9 17 4 12"/></svg>
                  {{ b }}
                </div>
              </div>
            </div>

            <!-- Add/Edit package form -->
            <div v-if="showPkgForm" class="pkg-form">
              <div class="pkg-form-title">{{ editPkgIdx !== null ? 'Edit Paket' : `Tambah Paket #${packages.length + 1}` }}</div>
              <div class="form-row-2">
                <div class="form-field">
                  <label class="form-field-label">Nama Paket <span class="req">*</span></label>
                  <input v-model="pkgForm.name" type="text" class="form-field-input" placeholder="Contoh: Paket Kemitraan 10 Juta" />
                </div>
                <div class="form-field">
                  <label class="form-field-label">Harga <span class="req">*</span></label>
                  <div class="input-with-prefix">
                    <span class="input-prefix">Rp</span>
                    <input :value="formatIdrInput(pkgForm.price)" type="text" class="form-field-input has-prefix"
                      placeholder="10.000.000" @input="pkgForm.price = parseIdrInput($event)" />
                  </div>
                </div>
              </div>
              <div class="form-row-2">
                <div class="form-field">
                  <label class="form-field-label">Minimal DP</label>
                  <div class="input-with-prefix">
                    <span class="input-prefix">Rp</span>
                    <input :value="formatIdrInput(pkgForm.minimum_dp)" type="text" class="form-field-input has-prefix"
                      placeholder="5.000.000" @input="pkgForm.minimum_dp = parseIdrInput($event)" />
                  </div>
                </div>
                <div class="form-field">
                  <label class="form-field-label">Durasi</label>
                  <input v-model="pkgForm.duration" type="text" class="form-field-input" placeholder="Contoh: 3 Tahun" />
                </div>
              </div>
              <div class="form-row-2">
                <div class="form-field">
                  <label class="form-field-label">Estimasi BEP</label>
                  <input v-model="pkgForm.estimated_bep" type="text" class="form-field-input" placeholder="Contoh: 6 - 12 Bulan" />
                </div>
              </div>
              <div class="form-row-2">
                <div class="form-field">
                  <label class="form-field-label">Net Profit</label>
                  <input v-model="pkgForm.net_profit" type="text" class="form-field-input" placeholder="Contoh: 10% - 15%" />
                </div>
                <div class="form-field">
                  <label class="form-field-label">Gambar Paket</label>
                  <div v-if="pkgForm.image" class="media-preview banner-preview">
                    <img :src="pkgForm.image" alt="Paket" @error="onImgError($event)" />
                    <button type="button" class="media-remove-btn" @click="pkgForm.image = ''">&times;</button>
                  </div>
                  <div v-else class="upload-dropzone upload-dropzone-sm" @click="$refs.pkgImageInput.click()">
                    <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="#94a3b8" stroke-width="1.5"><path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4"/><polyline points="17 8 12 3 7 8"/><line x1="12" y1="3" x2="12" y2="15"/></svg>
                    <span>Upload</span>
                  </div>
                  <input ref="pkgImageInput" type="file" accept="image/*" style="display:none" @change="handlePkgImageUpload" />
                </div>
              </div>
              <div class="form-field">
                <label class="form-field-label">Benefit</label>
                <div v-for="(item, bi) in pkgBenefits" :key="bi" class="benefit-row">
                  <input :value="item" type="text" class="form-field-input benefit-input"
                    placeholder="Contoh: Free Marketing Senilai 5 Juta"
                    @input="pkgBenefits[bi] = $event.target.value" />
                  <button type="button" class="benefit-remove-btn" @click="pkgBenefits.splice(bi, 1)" title="Hapus">
                    <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor"><line x1="18" y1="6" x2="6" y2="18" stroke-width="2" stroke-linecap="round"/><line x1="6" y1="6" x2="18" y2="18" stroke-width="2" stroke-linecap="round"/></svg>
                  </button>
                </div>
                <button type="button" class="benefit-add-btn" @click="pkgBenefits.push('')">
                  <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor"><line x1="12" y1="5" x2="12" y2="19" stroke-width="2" stroke-linecap="round"/><line x1="5" y1="12" x2="19" y2="12" stroke-width="2" stroke-linecap="round"/></svg>
                  Tambah Benefit
                </button>
              </div>
              <div class="pkg-form-actions">
                <button type="button" @click="showPkgForm = false" class="btn-secondary-sm">Selesai</button>
                <button type="button" @click="saveAndAddAnother" class="btn-outline-sm" v-if="editPkgIdx === null">
                  Simpan & Tambah Lagi
                </button>
                <button type="button" @click="savePackageLocal" class="btn-primary-sm">
                  {{ editPkgIdx !== null ? 'Update Paket' : 'Simpan Paket' }}
                </button>
              </div>
            </div>

            <!-- Add button -->
            <button v-if="!showPkgForm" type="button" @click="openAddPackage" class="pkg-add-btn">
              <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor"><line x1="12" y1="5" x2="12" y2="19" stroke-width="2" stroke-linecap="round"/><line x1="5" y1="12" x2="19" y2="12" stroke-width="2" stroke-linecap="round"/></svg>
              Tambah Paket
            </button>

            <div v-if="!packages.length && !showPkgForm" class="pkg-empty">
              Belum ada paket. Klik tombol di atas untuk menambahkan.
            </div>
          </div>

        </div>

        <!-- Sidebar right -->
        <div class="form-sidebar">
          <!-- Media card -->
          <div class="sidebar-card">
            <div class="sidebar-card-header">
              <div class="sidebar-icon icon-pink">
                <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor"><rect x="3" y="3" width="18" height="18" rx="2" ry="2" stroke-width="2"/><circle cx="8.5" cy="8.5" r="1.5" stroke-width="2"/><polyline points="21 15 16 10 5 21" stroke-width="2"/></svg>
              </div>
              <h4 class="sidebar-card-title">Media</h4>
            </div>
            <div class="sidebar-card-body">
              <div class="form-field">
                <label class="form-field-label-sm">Logo</label>
                <div v-if="form.logo" class="media-preview">
                  <img :src="form.logo" alt="Logo" @error="onImgError($event)" />
                  <button type="button" class="media-remove-btn" @click="form.logo = ''" title="Hapus">&times;</button>
                </div>
                <div v-else class="upload-dropzone" :class="{ 'is-uploading': uploadingLogo }" @click="!uploadingLogo && $refs.logoInput.click()">
                  <svg v-if="!uploadingLogo" width="22" height="22" viewBox="0 0 24 24" fill="none" stroke="#94a3b8" stroke-width="1.5"><path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4"/><polyline points="17 8 12 3 7 8"/><line x1="12" y1="3" x2="12" y2="15"/></svg>
                  <svg v-else class="spin-icon" width="22" height="22" viewBox="0 0 24 24" fill="none" stroke="#6366f1" stroke-width="2.5"><path d="M21 12a9 9 0 1 1-6.219-8.56" stroke-linecap="round"/></svg>
                  <span>{{ uploadingLogo ? 'Mengunggah...' : 'Upload Logo' }}</span>
                </div>
                <input ref="logoInput" type="file" accept="image/*" style="display:none" @change="handleUpload($event, 'logo')" />
              </div>
              <div class="form-field">
                <label class="form-field-label-sm">Banner</label>
                <div v-if="form.banner" class="media-preview banner-preview">
                  <img :src="form.banner" alt="Banner" @error="onImgError($event)" />
                  <button type="button" class="media-remove-btn" @click="form.banner = ''" title="Hapus">&times;</button>
                </div>
                <div v-else class="upload-dropzone" :class="{ 'is-uploading': uploadingBanner }" @click="!uploadingBanner && $refs.bannerInput.click()">
                  <svg v-if="!uploadingBanner" width="22" height="22" viewBox="0 0 24 24" fill="none" stroke="#94a3b8" stroke-width="1.5"><path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4"/><polyline points="17 8 12 3 7 8"/><line x1="12" y1="3" x2="12" y2="15"/></svg>
                  <svg v-else class="spin-icon" width="22" height="22" viewBox="0 0 24 24" fill="none" stroke="#6366f1" stroke-width="2.5"><path d="M21 12a9 9 0 1 1-6.219-8.56" stroke-linecap="round"/></svg>
                  <span>{{ uploadingBanner ? 'Mengunggah...' : 'Upload Banner' }}</span>
                </div>
                <input ref="bannerInput" type="file" accept="image/*" style="display:none" @change="handleUpload($event, 'banner')" />
              </div>
            </div>
          </div>

          <!-- Extra info -->
          <div class="sidebar-card">
            <div class="sidebar-card-header">
              <div class="sidebar-icon icon-cyan">
                <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor"><circle cx="12" cy="12" r="10" stroke-width="2"/><line x1="12" y1="16" x2="12" y2="12" stroke-width="2"/><line x1="12" y1="8" x2="12.01" y2="8" stroke-width="2"/></svg>
              </div>
              <h4 class="sidebar-card-title">Lainnya</h4>
            </div>
            <div class="sidebar-card-body">
              <div class="form-field">
                <label class="form-field-label-sm">Total Outlet</label>
                <input v-model.number="form.total_outlets" type="number" class="form-field-input" placeholder="0" min="0" />
              </div>
              <div class="form-field">
                <label class="form-field-label-sm">Tahun Berdiri</label>
                <input v-model.number="form.year_established" type="number" class="form-field-input" placeholder="2020" />
              </div>
            </div>
          </div>

          <!-- Submit sticky -->
          <div class="form-submit-card">
            <button type="submit" class="form-submit-btn" :disabled="submitting">
              <svg v-if="submitting" class="spin-icon" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor"><path d="M21 12a9 9 0 1 1-6.219-8.56" stroke-width="2.5" stroke-linecap="round"/></svg>
              <svg v-else width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor"><path d="M19 21H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h11l5 5v11a2 2 0 0 1-2 2z" stroke-width="2"/><polyline points="17 21 17 13 7 13 7 21" stroke-width="2"/><polyline points="7 3 7 8 15 8" stroke-width="2"/></svg>
              {{ submitting ? 'Menyimpan...' : (isEdit ? 'Update Outlet' : 'Simpan Outlet') }}
            </button>
            <router-link to="/outlets" class="form-cancel-btn">Batal</router-link>
          </div>
        </div>
      </div>
    </form>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted, nextTick } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { outletApi, outletCategoryApi, uploadApi, outletPackageApi } from '../../services/api'
import { useToastStore } from '../../stores/toast'

const route = useRoute()
const router = useRouter()
const toast = useToastStore()

const isEdit = computed(() => !!route.params.id && route.params.id !== 'create')
const submitting = ref(false)
const editorRef = ref(null)
const activeTab = ref(0)
const categories = ref([])
const uploadingLogo = ref(false)
const uploadingBanner = ref(false)

const tabs = computed(() => [
  { label: 'Informasi Dasar', icon: '<svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor"><path d="M14.5 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V7.5L14.5 2z" stroke-width="2"/><polyline points="14 2 14 8 20 8" stroke-width="2"/></svg>' },
  { label: 'Investasi', icon: '<svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor"><path d="M12 1v22M17 5H9.5a3.5 3.5 0 0 0 0 7h5a3.5 3.5 0 0 1 0 7H6" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/></svg>' },
  { label: 'Lokasi', icon: '<svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor"><path d="M21 10c0 7-9 13-9 13s-9-6-9-13a9 9 0 0 1 18 0z" stroke-width="2"/><circle cx="12" cy="10" r="3" stroke-width="2"/></svg>' },
  { label: 'Kontak', icon: '<svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor"><path d="M22 16.92v3a2 2 0 0 1-2.18 2 19.79 19.79 0 0 1-8.63-3.07 19.5 19.5 0 0 1-6-6 19.79 19.79 0 0 1-3.07-8.67A2 2 0 0 1 4.11 2h3a2 2 0 0 1 2 1.72 12.5 12.5 0 0 0 .7 2.81 2 2 0 0 1-.45 2.11L8.09 9.91a16 16 0 0 0 6 6l1.27-1.27a2 2 0 0 1 2.11-.45c.91.34 1.85.57 2.81.7A2 2 0 0 1 22 16.92z" stroke-width="2"/></svg>' },
  { label: 'Paket', icon: '<svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor"><path d="M21 16V8a2 2 0 0 0-1-1.73l-7-4a2 2 0 0 0-2 0l-7 4A2 2 0 0 0 3 8v8a2 2 0 0 0 1 1.73l7 4a2 2 0 0 0 2 0l7-4A2 2 0 0 0 21 16z" stroke-width="2"/></svg>', badge: packages.value.length || null },
])

// Package state
const packages = ref([])
const deletedPackageIds = ref([])
const showPkgForm = ref(false)
const editPkgIdx = ref(null)
const pkgForm = reactive({ name: '', price: 0, minimum_dp: 0, duration: '', estimated_bep: '', net_profit: '', image: '' })
const pkgBenefits = ref([''])
let pkgKeyCounter = 0

const form = reactive({
  name: '', category_id: '', short_description: '', description: '',
  minimum_investment: null, profit_sharing_percentage: 0,
  estimated_roi: '', location_requirement: '',
  address: '', city: '', province: '',
  contact_phone: '', contact_email: '', contact_whatsapp: '', website: '',
  logo: '', banner: '', total_outlets: 0, year_established: null,
})

const errors = reactive({})

onMounted(async () => {
  // Load categories
  try {
    const { data } = await outletCategoryApi.list({ active_only: 'true' })
    categories.value = data.data || []
  } catch (e) { console.error('Failed to load categories', e) }

  if (isEdit.value) {
    try {
      const { data } = await outletApi.get(route.params.id)
      const o = data.data
      Object.keys(form).forEach(k => { if (o[k] != null) form[k] = o[k] })
      await nextTick()
      if (editorRef.value) editorRef.value.innerHTML = form.description || ''

      // Load packages
      try {
        const pkgRes = await outletPackageApi.listByOutlet(route.params.id)
        packages.value = (pkgRes.data.data || []).map(p => ({ ...p, _existing: true }))
      } catch (e) { console.error('Failed to load packages', e) }
    } catch { toast.error('Gagal memuat data outlet'); router.push('/outlets') }
  }
})

async function handleUpload(event, field) {
  const file = event.target.files?.[0]
  if (!file) return
  const loading = field === 'logo' ? uploadingLogo : uploadingBanner
  loading.value = true
  // Show local preview immediately
  const localPreview = URL.createObjectURL(file)
  form[field] = localPreview
  try {
    const { data } = await uploadApi.upload(file)
    const serverUrl = data.data?.url || data.url || ''
    if (serverUrl) {
      form[field] = serverUrl
    }
  } catch (e) {
    toast.error(e.response?.data?.error || 'Gagal mengunggah file')
    form[field] = '' // clear on actual upload failure
  } finally {
    loading.value = false
    event.target.value = ''
    URL.revokeObjectURL(localPreview)
  }
}

function onImgError(event) {
  // Show a gray placeholder instead of clearing the URL
  event.target.style.display = 'none'
  const parent = event.target.parentElement
  if (parent && !parent.querySelector('.img-fallback')) {
    const fallback = document.createElement('div')
    fallback.className = 'img-fallback'
    fallback.innerHTML = '<svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="#94a3b8" stroke-width="1.5"><rect x="3" y="3" width="18" height="18" rx="2"/><circle cx="8.5" cy="8.5" r="1.5"/><path d="m21 15-5-5L5 21"/></svg>'
    parent.insertBefore(fallback, event.target)
  }
}

function validateField(field) {
  delete errors[field]
  switch (field) {
    case 'name':
      if (!form.name?.trim()) errors.name = 'Nama outlet wajib diisi'
      else if (form.name.trim().length < 3) errors.name = 'Nama outlet minimal 3 karakter'
      else if (form.name.length > 255) errors.name = 'Nama outlet maksimal 255 karakter'
      break
    case 'category_id':
      if (!form.category_id) errors.category_id = 'Kategori wajib dipilih'
      break
    case 'description':
      form.description = editorRef.value?.innerHTML || ''
      if (!form.description?.replace(/<[^>]*>/g, '').trim()) errors.description = 'Deskripsi outlet wajib diisi'
      break
    case 'minimum_investment':
      if (!form.minimum_investment || form.minimum_investment <= 0) errors.minimum_investment = 'Harga mulai dari harus lebih dari 0'
      break
    case 'profit_sharing_percentage':
      if (form.profit_sharing_percentage < 0 || form.profit_sharing_percentage > 100)
        errors.profit_sharing_percentage = 'Harus antara 0 - 100%'
      break
    case 'city':
      if (!form.city?.trim()) errors.city = 'Kota wajib diisi'
      break
    case 'province':
      if (!form.province?.trim()) errors.province = 'Provinsi wajib diisi'
      break
    case 'contact_email':
      if (form.contact_email && !/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(form.contact_email))
        errors.contact_email = 'Format email tidak valid'
      break
    case 'website':
      if (form.website && !form.website.startsWith('http'))
        errors.website = 'Harus diawali http:// atau https://'
      break
  }
}

function validateAll() {
  ['name','category_id','description','minimum_investment','city','province','profit_sharing_percentage','contact_email','website']
    .forEach(f => validateField(f))
  return Object.keys(errors).length === 0
}

function execCmd(cmd) { document.execCommand(cmd, false, null); editorRef.value?.focus() }
function execHeading(e) { if (e.target.value) document.execCommand('formatBlock', false, e.target.value); e.target.value = ''; editorRef.value?.focus() }
function onEditorInput() { form.description = editorRef.value?.innerHTML || '' }

async function handleSubmit() {
  if (!validateAll()) {
    toast.error('Mohon perbaiki error pada form')
    // Switch to tab with errors
    if (errors.name || errors.category_id || errors.description) activeTab.value = 0
    else if (errors.minimum_investment || errors.profit_sharing_percentage) activeTab.value = 1
    else if (errors.city || errors.province) activeTab.value = 2
    else if (errors.contact_email || errors.website) activeTab.value = 3
    return
  }
  submitting.value = true
  try {
    const payload = { ...form }
    if (!payload.year_established) payload.year_established = null

    const res = isEdit.value
      ? await outletApi.update(route.params.id, payload)
      : await outletApi.create(payload)

    const outletId = res.data.data?.id || route.params.id
    await savePackagesToServer(outletId)

    toast.success(res.data.message || (isEdit.value ? 'Outlet diupdate' : 'Outlet dibuat'))
    router.push('/outlets')
  } catch (e) {
    if (e.response?.status === 422 && e.response.data?.errors) {
      e.response.data.errors.forEach(err => { if (err.field !== '_general') errors[err.field] = err.message })
      toast.error('Validasi gagal, periksa form')
    } else { toast.error(e.response?.data?.error || 'Gagal menyimpan') }
  } finally { submitting.value = false }
}

async function savePackagesToServer(outletId) {
  for (const id of deletedPackageIds.value) {
    try { await outletPackageApi.delete(id) } catch (e) { console.error('Failed to delete package', e) }
  }
  for (let i = 0; i < packages.value.length; i++) {
    const pkg = packages.value[i]
    const pkgPayload = {
      outlet_id: outletId, name: pkg.name, price: pkg.price, minimum_dp: pkg.minimum_dp || 0,
      duration: pkg.duration || '', image: pkg.image || '',
      estimated_bep: pkg.estimated_bep || '', net_profit: pkg.net_profit || '',
      description: pkg.description || '', benefits: pkg.benefits || [], sort_order: i,
    }
    try {
      if (pkg._existing && pkg.id) await outletPackageApi.update(pkg.id, pkgPayload)
      else await outletPackageApi.create(pkgPayload)
    } catch (e) { console.error('Failed to save package', e) }
  }
}

function openAddPackage() {
  editPkgIdx.value = null
  Object.assign(pkgForm, { name: '', price: 0, minimum_dp: 0, duration: '', estimated_bep: '', net_profit: '', image: '' })
  pkgBenefits.value = ['']
  showPkgForm.value = true
}

function editPackage(idx) {
  editPkgIdx.value = idx
  const pkg = packages.value[idx]
  Object.assign(pkgForm, {
    name: pkg.name, price: pkg.price, minimum_dp: pkg.minimum_dp || 0, duration: pkg.duration || '',
    estimated_bep: pkg.estimated_bep || '', net_profit: pkg.net_profit || '',
    image: pkg.image || '',
  })
  pkgBenefits.value = (pkg.benefits && pkg.benefits.length) ? [...pkg.benefits] : ['']
  showPkgForm.value = true
}

function removePackage(idx) {
  const pkg = packages.value[idx]
  if (pkg._existing && pkg.id) deletedPackageIds.value.push(pkg.id)
  packages.value.splice(idx, 1)
}

function savePackageLocal() {
  if (!pkgForm.name?.trim()) { toast.error('Nama paket wajib diisi'); return }
  if (!pkgForm.price || pkgForm.price <= 0) { toast.error('Harga harus lebih dari 0'); return }

  const benefits = pkgBenefits.value.map(b => b.trim()).filter(Boolean)
  const pkgData = {
    name: pkgForm.name, price: pkgForm.price, minimum_dp: pkgForm.minimum_dp,
    duration: pkgForm.duration, estimated_bep: pkgForm.estimated_bep,
    net_profit: pkgForm.net_profit, image: pkgForm.image, benefits,
  }

  if (editPkgIdx.value !== null) {
    packages.value[editPkgIdx.value] = { ...packages.value[editPkgIdx.value], ...pkgData }
  } else {
    pkgData._key = 'new-' + (++pkgKeyCounter)
    packages.value.push(pkgData)
  }
  showPkgForm.value = false
}

function saveAndAddAnother() {
  if (!pkgForm.name?.trim()) { toast.error('Nama paket wajib diisi'); return }
  if (!pkgForm.price || pkgForm.price <= 0) { toast.error('Harga harus lebih dari 0'); return }

  const benefits = pkgBenefits.value.map(b => b.trim()).filter(Boolean)
  const pkgData = {
    name: pkgForm.name, price: pkgForm.price, minimum_dp: pkgForm.minimum_dp,
    duration: pkgForm.duration, estimated_bep: pkgForm.estimated_bep,
    net_profit: pkgForm.net_profit, image: pkgForm.image, benefits,
    _key: 'new-' + (++pkgKeyCounter),
  }
  packages.value.push(pkgData)
  toast.success(`Paket "${pkgData.name}" ditambahkan`)
  // Reset form for next package
  Object.assign(pkgForm, { name: '', price: 0, minimum_dp: 0, duration: '', estimated_bep: '', net_profit: '', image: '' })
  pkgBenefits.value = ['']
  editPkgIdx.value = null
}

async function handlePkgImageUpload(event) {
  const file = event.target.files?.[0]
  if (!file) return
  const localPreview = URL.createObjectURL(file)
  pkgForm.image = localPreview
  try {
    const { data } = await uploadApi.upload(file)
    const serverUrl = data.data?.url || data.url || ''
    if (serverUrl) pkgForm.image = serverUrl
  } catch (e) {
    toast.error(e.response?.data?.error || 'Gagal mengunggah')
    pkgForm.image = ''
  } finally {
    event.target.value = ''
    URL.revokeObjectURL(localPreview)
  }
}

function formatCurrency(v) {
  if (!v) return ''
  return new Intl.NumberFormat('id-ID', { style: 'currency', currency: 'IDR', minimumFractionDigits: 0 }).format(v)
}

function formatIdrInput(v) {
  if (!v && v !== 0) return ''
  const num = typeof v === 'string' ? parseInt(v.replace(/\D/g, '')) : v
  if (!num && num !== 0) return ''
  return num.toLocaleString('id-ID')
}

function parseIdrInput(event) {
  const raw = event.target.value.replace(/\D/g, '')
  const num = parseInt(raw) || 0
  event.target.value = num ? num.toLocaleString('id-ID') : ''
  return num
}
</script>

<style scoped>
/* ═══ HERO ═══ */
.form-hero {
  background: linear-gradient(135deg, #1a1c2e 0%, #2d2f48 60%, #3d2060 100%);
  border-radius: 1rem; padding: 24px 28px 20px;
  margin-bottom: 0; display: flex; flex-wrap: wrap;
  align-items: center; gap: 16px;
  border-bottom-left-radius: 0; border-bottom-right-radius: 0;
}
.form-hero-back {
  width: 40px; height: 40px; border-radius: 10px;
  background: rgba(255,255,255,0.1); color: white;
  display: flex; align-items: center; justify-content: center;
  transition: background 0.2s; text-decoration: none;
}
.form-hero-back:hover { background: rgba(255,255,255,0.2); }
.form-hero-title { font-size: 1.25rem; font-weight: 800; color: white; }
.form-hero-sub { font-size: 0.8rem; color: rgba(255,255,255,0.5); margin-top: 2px; }

/* ═══ TABS ═══ */
.form-tabs {
  display: flex; gap: 0; background: #fff;
  border: 1px solid #e8eaee; border-top: none;
  border-radius: 0 0 14px 14px; margin-bottom: 24px;
  overflow-x: auto; padding: 0 8px;
}
.form-tab {
  display: flex; align-items: center; gap: 8px;
  padding: 14px 20px; border: none; background: none;
  font-size: 0.82rem; font-weight: 600; color: #94a3b8;
  cursor: pointer; transition: all .2s ease;
  border-bottom: 3px solid transparent; white-space: nowrap;
  position: relative;
}
.form-tab:hover { color: #475569; }
.form-tab.active {
  color: #6366f1; border-bottom-color: #6366f1;
}
.form-tab-icon { display: flex; align-items: center; }
.form-tab.active .form-tab-icon svg { stroke: #6366f1; }
.form-tab-badge {
  background: #6366f1; color: white; font-size: 0.65rem;
  padding: 1px 7px; border-radius: 10px; font-weight: 700;
  min-width: 18px; text-align: center;
}

/* ═══ LAYOUT ═══ */
.form-layout {
  display: grid; grid-template-columns: 1fr 300px; gap: 24px;
  align-items: start;
}

/* ═══ TAB PANELS ═══ */
.form-tab-panel {
  background: white; border-radius: 14px; border: 1px solid #e8eaee;
  padding: 32px 36px; min-height: 200px;
}

/* ═══ SIDEBAR ═══ */
.sidebar-card {
  background: white; border-radius: 14px; border: 1px solid #e8eaee;
  margin-bottom: 16px; overflow: hidden;
}
.sidebar-card-header {
  display: flex; align-items: center; gap: 10px;
  padding: 14px 18px; border-bottom: 1px solid #f1f5f9;
}
.sidebar-icon {
  width: 32px; height: 32px; border-radius: 8px; flex-shrink: 0;
  display: flex; align-items: center; justify-content: center;
}
.sidebar-card-title { font-size: 0.82rem; font-weight: 700; color: #1e293b; }
.sidebar-card-body { padding: 16px 18px; }

.icon-pink { background: #fdf2f8; color: #ec4899; }
.icon-cyan { background: #ecfeff; color: #06b6d4; }

/* ═══ FORM FIELDS ═══ */
.form-field { margin-bottom: 16px; }
.form-field:last-child { margin-bottom: 0; }
.form-field-label {
  display: block; font-size: 0.82rem; font-weight: 600;
  color: #475569; margin-bottom: 8px;
}
.form-field-label-sm { display: block; font-size: 0.75rem; font-weight: 600; color: #64748b; margin-bottom: 5px; }
.req { color: #ef4444; }

.form-field-input {
  width: 100%; padding: 12px 16px; font-size: 0.85rem;
  border: 1.5px solid #e2e8f0; border-radius: 10px;
  background: #fafbfc; color: #1a1c2e; transition: all 0.2s;
  outline: none; font-family: inherit;
}
.form-field-input:focus { border-color: #6366f1; background: #fff; box-shadow: 0 0 0 3px rgba(99,102,241,0.08); }
.form-field-input::placeholder { color: #b0b8c1; }
.form-field-input.is-error { border-color: #ef4444; background: #fef2f2; }
select.form-field-input {
  appearance: none;
  background-image: url("data:image/svg+xml,%3csvg xmlns='http://www.w3.org/2000/svg' fill='none' viewBox='0 0 20 20'%3e%3cpath stroke='%236b7280' stroke-linecap='round' stroke-linejoin='round' stroke-width='1.5' d='M6 8l4 4 4-4'/%3e%3c/svg%3e");
  background-position: right 10px center; background-repeat: no-repeat; background-size: 1.4em; padding-right: 36px;
}
textarea.form-field-input { resize: vertical; min-height: 60px; }

.form-field-error { font-size: 0.72rem; color: #ef4444; margin-top: 5px; font-weight: 500; }
.form-field-hint { font-size: 0.72rem; color: #94a3b8; margin-top: 5px; }
.form-field-counter { font-size: 0.68rem; color: #b0b8c1; margin-top: 4px; text-align: right; }

.form-row-2 { display: grid; grid-template-columns: 1fr 1fr; gap: 20px; margin-bottom: 16px; }
.form-row-2:last-child { margin-bottom: 0; }

/* Prefix/Suffix */
.input-with-prefix, .input-with-suffix { position: relative; }
.input-prefix, .input-suffix {
  position: absolute; top: 50%; transform: translateY(-50%);
  font-size: 0.8rem; font-weight: 700; color: #94a3b8; pointer-events: none;
}
.input-prefix { left: 14px; }
.input-suffix { right: 14px; }
.form-field-input.has-prefix { padding-left: 40px; }
.form-field-input.has-suffix { padding-right: 36px; }

/* Media preview */
.media-preview {
  margin-top: 8px; border-radius: 10px; overflow: hidden;
  border: 1px solid #e2e8f0; height: 80px; position: relative;
}
.media-preview img { width: 100%; height: 100%; object-fit: contain; background: #f8fafc; }
.banner-preview { height: 100px; }
.banner-preview img { object-fit: cover; }
.media-remove-btn {
  position: absolute; top: 6px; right: 6px;
  width: 24px; height: 24px; border-radius: 50%;
  background: rgba(0,0,0,0.5); color: #fff;
  border: none; cursor: pointer; font-size: 1rem;
  display: flex; align-items: center; justify-content: center;
  transition: background .15s ease;
}
.media-remove-btn:hover { background: rgba(239,68,68,0.8); }

/* Upload dropzone */
.upload-dropzone {
  display: flex; flex-direction: column; align-items: center; justify-content: center;
  gap: 6px; padding: 16px 12px;
  border: 2px dashed #e2e8f0; border-radius: 10px;
  cursor: pointer; margin-top: 8px; transition: all .2s ease;
  color: #94a3b8; font-size: 0.78rem; font-weight: 500;
}
.upload-dropzone:hover { border-color: #6366f1; background: rgba(99,102,241,.03); color: #6366f1; }
.upload-dropzone:hover svg { stroke: #6366f1; }
.upload-dropzone.is-uploading { border-color: #c7d2fe; background: #eef2ff; cursor: wait; }
.upload-dropzone-sm { padding: 12px 10px; }
.img-fallback { width: 100%; height: 100%; display: flex; align-items: center; justify-content: center; background: #f1f5f9; }
.spin-icon { animation: spin 0.8s linear infinite; }
@keyframes spin { from { transform: rotate(0deg); } to { transform: rotate(360deg); } }

/* ═══ RICH TEXT EDITOR ═══ */
.rte-wrapper { border: 1.5px solid #e2e8f0; border-radius: 10px; overflow: hidden; transition: border-color 0.2s; }
.rte-wrapper:focus-within { border-color: #6366f1; box-shadow: 0 0 0 3px rgba(99,102,241,0.08); }
.rte-wrapper.is-error { border-color: #ef4444; }

.rte-toolbar {
  display: flex; align-items: center; gap: 2px; padding: 6px 10px;
  background: #f8fafc; border-bottom: 1px solid #f1f5f9; flex-wrap: wrap;
}
.rte-btn {
  width: 30px; height: 30px; display: flex; align-items: center; justify-content: center;
  border: none; background: none; border-radius: 6px; cursor: pointer;
  color: #64748b; font-size: 0.8rem; transition: all 0.15s;
}
.rte-btn:hover { background: #e2e8f0; color: #1a1c2e; }
.rte-sep { width: 1px; height: 20px; background: #e2e8f0; margin: 0 3px; }
.rte-select {
  height: 30px; padding: 0 6px; border: 1px solid #e2e8f0; border-radius: 6px;
  font-size: 0.7rem; color: #64748b; background: white; cursor: pointer; outline: none;
}

.rte-content {
  min-height: 180px; padding: 14px 16px; outline: none;
  font-size: 0.85rem; line-height: 1.7; color: #1a1c2e; background: #fff;
}
.rte-content:empty:before { content: 'Tulis deskripsi detail outlet di sini...'; color: #b0b8c1; }
.rte-content h2 { font-size: 1.15rem; font-weight: 700; margin: 0.5em 0; }
.rte-content h3 { font-size: 1rem; font-weight: 600; margin: 0.5em 0; }
.rte-content ul, .rte-content ol { padding-left: 1.5rem; margin: 0.5em 0; }

/* ═══ SUBMIT ═══ */
.form-submit-card {
  position: sticky; top: 100px;
  background: white; border-radius: 14px; border: 1px solid #e8eaee;
  padding: 18px; display: flex; flex-direction: column; gap: 10px;
}
.form-submit-btn {
  width: 100%; padding: 12px; border: none; border-radius: 10px;
  background: linear-gradient(135deg, #6366f1 0%, #4f46e5 100%);
  color: white; font-size: 0.85rem; font-weight: 700; cursor: pointer;
  display: flex; align-items: center; justify-content: center; gap: 8px;
  box-shadow: 0 4px 16px rgba(99,102,241,0.25); transition: all 0.25s;
}
.form-submit-btn:hover:not(:disabled) { transform: translateY(-2px); box-shadow: 0 8px 24px rgba(99,102,241,0.35); }
.form-submit-btn:disabled { opacity: 0.6; cursor: not-allowed; }
.form-cancel-btn {
  display: block; text-align: center; padding: 10px;
  font-size: 0.8rem; font-weight: 600; color: #94a3b8;
  text-decoration: none; border-radius: 10px; transition: all 0.2s;
}
.form-cancel-btn:hover { background: #f8fafc; color: #64748b; }

@keyframes spin { to { transform: rotate(360deg); } }
.spin-icon { animation: spin 1s linear infinite; }

/* ═══ PACKAGE CARDS ═══ */
.pkg-card {
  border: 1px solid #e2e8f0; border-radius: 12px; padding: 16px;
  margin-bottom: 12px; background: #fafbfc; transition: box-shadow .15s ease;
}
.pkg-card:hover { box-shadow: 0 2px 8px rgba(0,0,0,0.06); }
.pkg-card-header { display: flex; justify-content: space-between; align-items: flex-start; }
.pkg-card-name { font-weight: 600; font-size: 0.92rem; color: #1e293b; }
.pkg-card-price { font-size: 0.85rem; color: #6366f1; font-weight: 600; margin-top: 2px; }
.pkg-card-actions { display: flex; gap: 6px; }
.pkg-act-btn {
  width: 30px; height: 30px; border-radius: 8px; border: 1px solid #e2e8f0;
  background: #fff; cursor: pointer; display: flex; align-items: center; justify-content: center;
  transition: all .15s ease; color: #64748b;
}
.pkg-act-edit:hover { border-color: #6366f1; color: #6366f1; background: #eef2ff; }
.pkg-act-del:hover { border-color: #ef4444; color: #ef4444; background: #fef2f2; }
.pkg-card-meta {
  display: flex; gap: 16px; flex-wrap: wrap; margin-top: 10px; padding-top: 10px;
  border-top: 1px solid #f1f5f9; font-size: 0.78rem; color: #64748b;
}
.pkg-card-benefits { margin-top: 10px; display: flex; flex-direction: column; gap: 6px; }
.pkg-benefit-item { display: flex; align-items: center; gap: 6px; font-size: 0.8rem; color: #475569; }

.pkg-form {
  border: 2px solid #6366f1; border-radius: 12px; padding: 20px;
  background: #fafbff; margin-bottom: 12px;
}
.pkg-form-title { font-weight: 600; font-size: 0.95rem; color: #1e293b; margin-bottom: 16px; }
.pkg-form-actions { display: flex; justify-content: flex-end; gap: 10px; margin-top: 16px; }

.btn-primary-sm {
  padding: 8px 20px; border-radius: 8px; font-size: 0.82rem; font-weight: 600;
  border: none; cursor: pointer; background: linear-gradient(135deg, #6366f1, #4f46e5);
  color: #fff; transition: all .15s ease;
}
.btn-primary-sm:hover { transform: translateY(-1px); box-shadow: 0 3px 10px rgba(99,102,241,.3); }
.btn-secondary-sm {
  padding: 8px 20px; border-radius: 8px; font-size: 0.82rem; font-weight: 600;
  border: 1px solid #e2e8f0; cursor: pointer; background: #fff; color: #64748b;
  transition: all .15s ease;
}
.btn-secondary-sm:hover { border-color: #cbd5e1; background: #f8fafc; }
.btn-outline-sm {
  padding: 8px 20px; border-radius: 8px; font-size: 0.82rem; font-weight: 600;
  border: 1.5px solid #6366f1; cursor: pointer; background: #fff; color: #6366f1;
}
.btn-outline-sm:hover { background: #eef2ff; }

/* Benefit inputs */
.benefit-row {
  display: flex; align-items: center; gap: 8px; margin-bottom: 10px;
}
.benefit-input { flex: 1; margin-bottom: 0 !important; }
.benefit-remove-btn {
  width: 34px; height: 34px; border-radius: 8px; border: 1px solid #e2e8f0;
  background: #fff; cursor: pointer; display: flex; align-items: center; justify-content: center;
  color: #94a3b8; transition: all .15s ease; flex-shrink: 0;
}
.benefit-remove-btn:hover { border-color: #ef4444; color: #ef4444; background: #fef2f2; }
.benefit-add-btn {
  display: flex; align-items: center; gap: 6px;
  padding: 8px 14px; border: 1.5px dashed #d1d5db; border-radius: 8px;
  background: transparent; color: #6366f1; font-weight: 600; font-size: 0.8rem;
  cursor: pointer; transition: all .15s ease; margin-top: 4px;
}
.benefit-add-btn:hover { border-color: #6366f1; background: rgba(99,102,241,.03); }

.pkg-add-btn {
  display: flex; align-items: center; justify-content: center; gap: 8px;
  width: 100%; padding: 14px; border: 2px dashed #d1d5db; border-radius: 12px;
  background: transparent; color: #6366f1; font-weight: 600; font-size: 0.85rem;
  cursor: pointer; transition: all .2s ease;
}
.pkg-add-btn:hover { border-color: #6366f1; background: rgba(99,102,241,.03); }

.pkg-empty {
  text-align: center; padding: 32px 16px; color: #94a3b8;
  font-size: 0.82rem;
}

@media (max-width: 960px) {
  .form-layout { grid-template-columns: 1fr; }
  .form-sidebar { order: -1; }
  .form-submit-card { position: static; }
  .form-tabs { gap: 0; }
  .form-tab { padding: 12px 14px; font-size: 0.75rem; }
  .form-tab-label { display: none; }
  .form-row-2 { grid-template-columns: 1fr; }
}
</style>
