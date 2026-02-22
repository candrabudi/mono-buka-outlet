<template>
  <div class="af-page">
    <!-- Hero -->
    <div class="af-hero">
      <button class="af-back" @click="$router.back()">
        <i class="ri-arrow-left-line"></i>
        Kembali
      </button>
      <div>
        <h1 class="af-hero-title">Pengajuan Kemitraan</h1>
        <p class="af-hero-sub">Lengkapi data di bawah untuk mengajukan kemitraan</p>
      </div>
    </div>

    <!-- Loading -->
    <div v-if="loading" class="af-loading">
      <div class="af-skeleton shimmer" style="height:120px;border-radius:14px;margin-bottom:20px"></div>
      <div class="af-skeleton shimmer" style="height:300px;border-radius:14px"></div>
    </div>

    <template v-else>
      <!-- Summary Card -->
      <div class="af-summary">
        <div class="af-summary-left">
          <div class="af-summary-logo">
            <img v-if="outlet?.logo" :src="outlet.logo" :alt="outlet?.name" />
            <i v-else class="ri-home-5-line" style="font-size:24px;color:#94a3b8"></i>
          </div>
          <div>
            <div class="af-summary-outlet">{{ outlet?.name || '-' }}</div>
            <div class="af-summary-cat">{{ outlet?.category_name || outlet?.category || 'Franchise' }}</div>
          </div>
        </div>
        <div class="af-summary-divider"></div>
        <div class="af-summary-right">
          <div class="af-summary-pkg-label">Paket Dipilih</div>
          <div class="af-summary-pkg-name">{{ selectedPkg?.name || '-' }}</div>
          <div class="af-summary-pkg-price">{{ fc(selectedPkg?.price) }}</div>
        </div>
        <div class="af-summary-divider"></div>
        <div class="af-summary-specs">
          <div v-if="selectedPkg?.duration" class="af-spec">
            <span class="af-spec-label">Durasi</span>
            <span class="af-spec-val">{{ selectedPkg.duration }}</span>
          </div>
          <div v-if="selectedPkg?.estimated_bep" class="af-spec">
            <span class="af-spec-label">BEP</span>
            <span class="af-spec-val">{{ selectedPkg.estimated_bep }}</span>
          </div>
          <div v-if="selectedPkg?.net_profit" class="af-spec">
            <span class="af-spec-label">Laba</span>
            <span class="af-spec-val">{{ selectedPkg.net_profit }}</span>
          </div>
          <div v-if="selectedPkg?.minimum_dp" class="af-spec">
            <span class="af-spec-label">DP Min</span>
            <span class="af-spec-val">{{ fc(selectedPkg.minimum_dp) }}</span>
          </div>
        </div>
      </div>

      <!-- Form -->
      <div class="af-form-card">
        <h3 class="af-form-title">
          <i class="ri-edit-line"></i>
          Data Pengajuan
        </h3>

        <form @submit.prevent="submit" class="af-form">
          <!-- Step indicators -->
          <div class="af-steps">
            <div class="af-step" :class="{ active: step >= 1, done: step > 1 }"><span>1</span> Motivasi</div>
            <div class="af-step-line" :class="{ active: step > 1 }"></div>
            <div class="af-step" :class="{ active: step >= 2, done: step > 2 }"><span>2</span> Kontak & Detail</div>
            <div class="af-step-line" :class="{ active: step > 2 }"></div>
            <div class="af-step" :class="{ active: step >= 3 }"><span>3</span> Konfirmasi</div>
          </div>

          <!-- Step 1: Motivation -->
          <div v-show="step === 1" class="af-step-content">
            <div class="af-field">
              <label>Motivasi Bergabung <span class="af-req">*</span></label>
              <p class="af-hint">Ceritakan alasan Anda tertarik bergabung dengan outlet ini</p>
              <textarea v-model="form.motivation" placeholder="Saya tertarik bergabung karena..." rows="5" class="af-textarea"></textarea>
            </div>
            <div class="af-field">
              <label>Pengalaman Bisnis</label>
              <p class="af-hint">Jelaskan pengalaman bisnis Anda sebelumnya (jika ada)</p>
              <textarea v-model="form.experience" placeholder="Pengalaman saya di bidang bisnis..." rows="4" class="af-textarea"></textarea>
            </div>
            <div class="af-actions">
              <div></div>
              <button type="button" class="af-btn-next" @click="step = 2" :disabled="!form.motivation.trim()">
                Lanjutkan
                <i class="ri-arrow-right-s-line"></i>
              </button>
            </div>
          </div>

          <!-- Step 2: Contact + Details -->
          <div v-show="step === 2" class="af-step-content">
            <div class="af-field-row">
              <div class="af-field">
                <label>No. HP yang Bisa Dihubungi <span class="af-req">*</span></label>
                <p class="af-hint">Nomor aktif untuk komunikasi terkait pengajuan</p>
                <div class="af-input-prefix-wrap">
                  <span class="af-input-prefix"><i class="ri-phone-line"></i></span>
                  <input v-model="form.contact_phone" type="text" placeholder="08xxxxxxxxxx" class="af-input with-prefix" />
                </div>
                <span v-if="form.contact_phone && !isPhoneValid" class="af-field-error">Format: 08xx atau +62xx</span>
              </div>
              <div class="af-field">
                <label>Email yang Bisa Dihubungi <span class="af-req">*</span></label>
                <p class="af-hint">Email aktif untuk menerima informasi pengajuan</p>
                <div class="af-input-prefix-wrap">
                  <span class="af-input-prefix"><i class="ri-mail-line"></i></span>
                  <input v-model="form.contact_email" type="email" placeholder="email@contoh.com" class="af-input with-prefix" />
                </div>
                <span v-if="form.contact_email && !isEmailValid" class="af-field-error">Format email tidak valid</span>
              </div>
            </div>
            <div class="af-field">
              <label>Rencana Lokasi <span class="af-req">*</span></label>
              <p class="af-hint">Lokasi yang Anda rencanakan untuk mendirikan outlet</p>
              <input v-model="form.proposed_location" type="text" placeholder="Contoh: Jl. Sudirman No.10, Jakarta Pusat" class="af-input" />
            </div>
            <div class="af-field">
              <label>Budget Investasi <span class="af-req">*</span></label>
              <p class="af-hint">Berapa budget yang Anda siapkan untuk investasi ini</p>
              <div class="af-input-prefix-wrap">
                <span class="af-input-prefix">Rp</span>
                <input v-model.number="form.investment_budget" type="number" placeholder="150000000" min="0" class="af-input with-prefix" />
              </div>
              <p v-if="form.investment_budget" class="af-budget-formatted">{{ fc(form.investment_budget) }}</p>
            </div>
            <div class="af-actions">
              <button type="button" class="af-btn-back" @click="step = 1">
                <i class="ri-arrow-left-s-line"></i>
                Kembali
              </button>
              <button type="button" class="af-btn-next" @click="step = 3" :disabled="!canStep3">
                Lanjutkan
                <i class="ri-arrow-right-s-line"></i>
              </button>
            </div>
          </div>

          <!-- Step 3: Confirmation -->
          <div v-show="step === 3" class="af-step-content">
            <div class="af-confirm-grid">
              <div class="af-confirm-item">
                <div class="af-confirm-label">Outlet</div>
                <div class="af-confirm-val">{{ outlet?.name }}</div>
              </div>
              <div class="af-confirm-item">
                <div class="af-confirm-label">Paket</div>
                <div class="af-confirm-val">{{ selectedPkg?.name }} — {{ fc(selectedPkg?.price) }}</div>
              </div>
              <div class="af-confirm-item full">
                <div class="af-confirm-label">Motivasi</div>
                <div class="af-confirm-val">{{ form.motivation }}</div>
              </div>
              <div class="af-confirm-item full" v-if="form.experience">
                <div class="af-confirm-label">Pengalaman</div>
                <div class="af-confirm-val">{{ form.experience }}</div>
              </div>
              <div class="af-confirm-item">
                <div class="af-confirm-label">No. HP Kontak</div>
                <div class="af-confirm-val">{{ form.contact_phone }}</div>
              </div>
              <div class="af-confirm-item">
                <div class="af-confirm-label">Email Kontak</div>
                <div class="af-confirm-val">{{ form.contact_email }}</div>
              </div>
              <div class="af-confirm-item">
                <div class="af-confirm-label">Lokasi</div>
                <div class="af-confirm-val">{{ form.proposed_location }}</div>
              </div>
              <div class="af-confirm-item">
                <div class="af-confirm-label">Budget</div>
                <div class="af-confirm-val">{{ fc(form.investment_budget) }}</div>
              </div>
            </div>

            <div class="af-confirm-note">
              <i class="ri-error-warning-line"></i>
              Dengan mengirim pengajuan ini, Anda menyetujui bahwa data yang diisi adalah benar. Tim kami akan meninjau pengajuan Anda.
            </div>

            <div class="af-actions">
              <button type="button" class="af-btn-back" @click="step = 2">
                <i class="ri-arrow-left-s-line"></i>
                Kembali
              </button>
              <button type="submit" class="af-btn-submit" :disabled="submitting">
                <i class="ri-send-plane-line"></i>
                <span v-if="submitting" class="af-spin"></span>
                {{ submitting ? 'Mengirim...' : 'Kirim Pengajuan' }}
              </button>
            </div>
          </div>
        </form>
      </div>
    </template>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { outletApi, applicationApi } from '../services/api'
import { useToastStore } from '../stores/toast'

const route = useRoute()
const router = useRouter()
const toast = useToastStore()

const outlet = ref(null)
const packages = ref([])
const selectedPkg = ref(null)
const loading = ref(true)
const submitting = ref(false)
const step = ref(1)
const form = reactive({ motivation: '', experience: '', proposed_location: '', investment_budget: 0, contact_phone: '', contact_email: '' })

const isPhoneValid = computed(() => /^(\+62|62|08)[0-9]{8,13}$/.test(form.contact_phone.replace(/[\s-]/g, '')))
const isEmailValid = computed(() => /^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(form.contact_email))
const canStep3 = computed(() => form.contact_phone.trim() && isPhoneValid.value && form.contact_email.trim() && isEmailValid.value && form.proposed_location.trim() && form.investment_budget > 0)

function fc(v) { return v ? new Intl.NumberFormat('id-ID', { style: 'currency', currency: 'IDR', minimumFractionDigits: 0 }).format(v) : '-' }

async function loadData() {
  loading.value = true
  try {
    const [outletRes, pkgRes] = await Promise.all([
      outletApi.getByID(route.params.id),
      outletApi.getPackages(route.params.id),
    ])
    outlet.value = outletRes.data.data
    packages.value = pkgRes.data.data || []
    const pkgId = route.query.package
    if (pkgId) selectedPkg.value = packages.value.find(p => p.id === pkgId) || packages.value[0] || null
    else selectedPkg.value = packages.value[0] || null
  } catch { toast.error('Gagal memuat data outlet') }
  finally { loading.value = false }
}

async function submit() {
  if (!form.motivation.trim()) return toast.error('Motivasi wajib diisi')
  if (!form.proposed_location.trim()) return toast.error('Lokasi wajib diisi')
  if (!form.investment_budget) return toast.error('Budget wajib diisi')
  submitting.value = true
  try {
    await applicationApi.apply({
      outlet_id: route.params.id,
      package_id: selectedPkg.value?.id,
      ...form,
    })
    toast.success('Pengajuan berhasil dikirim! Tim kami akan meninjau.')
    router.push('/applications')
  } catch (e) { toast.error(e.response?.data?.error || 'Gagal mengirim pengajuan') }
  finally { submitting.value = false }
}

onMounted(loadData)
</script>

<style scoped>
/* ═══ Hero ═══ */
.af-hero{background:linear-gradient(135deg,#0f0c29 0%,#302b63 50%,#24243e 100%);border-radius:16px;padding:28px 36px 24px;margin-bottom:24px;box-shadow:0 4px 24px rgba(15,12,41,.2)}
.af-back{display:inline-flex;align-items:center;gap:5px;font-size:.78rem;font-weight:600;color:rgba(255,255,255,.45);background:rgba(255,255,255,.06);border:1px solid rgba(255,255,255,.1);border-radius:8px;padding:6px 14px;cursor:pointer;margin-bottom:18px;transition:all .15s;backdrop-filter:blur(4px)}
.af-back:hover{color:#fff;border-color:rgba(255,255,255,.25);background:rgba(255,255,255,.1)}
.af-hero-title{font-size:1.6rem;font-weight:800;color:#fff;margin:0 0 4px}
.af-hero-sub{font-size:.85rem;color:rgba(255,255,255,.5);margin:0}

/* ═══ Loading ═══ */
.af-loading{max-width:740px}
.shimmer{background:linear-gradient(90deg,#f1f5f9 25%,#e8ecf1 50%,#f1f5f9 75%)!important;background-size:200% 100%;animation:shimmer 1.5s infinite}
@keyframes shimmer{0%{background-position:200% 0}100%{background-position:-200% 0}}

/* ═══ Summary ═══ */
.af-summary{background:#fff;border-radius:14px;border:1px solid #e8ecf1;padding:22px 24px;margin-bottom:24px;display:flex;align-items:center;gap:24px;flex-wrap:wrap}
.af-summary-left{display:flex;align-items:center;gap:14px}
.af-summary-logo{width:48px;height:48px;border-radius:12px;background:#f1f5f9;border:1px solid #e8ecf1;display:flex;align-items:center;justify-content:center;overflow:hidden;flex-shrink:0}
.af-summary-logo img{width:100%;height:100%;object-fit:cover}
.af-summary-outlet{font-size:.92rem;font-weight:700;color:#0f172a}
.af-summary-cat{font-size:.72rem;color:#94a3b8;text-transform:capitalize}
.af-summary-divider{width:1px;height:48px;background:#e8ecf1;flex-shrink:0}
.af-summary-right{min-width:0}
.af-summary-pkg-label{font-size:.65rem;color:#94a3b8;text-transform:uppercase;letter-spacing:.05em;font-weight:600}
.af-summary-pkg-name{font-size:.88rem;font-weight:700;color:#0f172a}
.af-summary-pkg-price{font-size:.95rem;font-weight:800;color:#6366f1}
.af-summary-specs{display:flex;gap:20px;flex-wrap:wrap}
.af-spec{display:flex;flex-direction:column}
.af-spec-label{font-size:.62rem;color:#94a3b8;text-transform:uppercase;letter-spacing:.04em}
.af-spec-val{font-size:.82rem;font-weight:700;color:#1e293b}

/* ═══ Form Card ═══ */
.af-form-card{background:#fff;border-radius:14px;border:1px solid #e8ecf1;padding:28px;max-width:740px}
.af-form-title{font-size:.92rem;font-weight:700;color:#0f172a;margin:0 0 24px;display:flex;align-items:center;gap:8px}
.af-form-title svg{color:#6366f1}

/* ═══ Steps ═══ */
.af-steps{display:flex;align-items:center;gap:0;margin-bottom:32px}
.af-step{display:flex;align-items:center;gap:8px;font-size:.78rem;font-weight:600;color:#94a3b8;transition:all .2s}
.af-step span{width:28px;height:28px;border-radius:50%;border:2px solid #e8ecf1;display:flex;align-items:center;justify-content:center;font-size:.72rem;font-weight:700;transition:all .2s}
.af-step.active{color:#0f172a}
.af-step.active span{border-color:#6366f1;background:#6366f1;color:#fff}
.af-step.done span{border-color:#22c55e;background:#22c55e;color:#fff}
.af-step-line{flex:1;height:2px;background:#e8ecf1;margin:0 12px;transition:background .3s}
.af-step-line.active{background:#6366f1}

/* ═══ Fields ═══ */
.af-step-content{animation:fadeIn .3s ease}
@keyframes fadeIn{from{opacity:0;transform:translateY(8px)}to{opacity:1;transform:translateY(0)}}
.af-field{margin-bottom:22px}
.af-field label{display:block;font-size:.85rem;font-weight:700;color:#0f172a;margin-bottom:4px}
.af-req{color:#ef4444}
.af-hint{font-size:.75rem;color:#94a3b8;margin:0 0 8px;line-height:1.4}
.af-input,.af-textarea{width:100%;border:1px solid #e8ecf1;border-radius:10px;padding:11px 14px;font-size:.85rem;color:#1e293b;box-sizing:border-box;font-family:inherit;transition:all .15s}
.af-input:focus,.af-textarea:focus{outline:none;border-color:#818cf8;box-shadow:0 0 0 3px rgba(129,140,248,.12)}
.af-textarea{resize:vertical;min-height:80px}
.af-input-prefix-wrap{position:relative}
.af-input-prefix{position:absolute;left:14px;top:50%;transform:translateY(-50%);font-size:.85rem;font-weight:700;color:#94a3b8}
.af-input.with-prefix{padding-left:40px}
.af-field-error{display:block;font-size:.72rem;color:#ef4444;font-weight:500;margin-top:4px}
.af-field-row{display:grid;grid-template-columns:1fr 1fr;gap:16px}
.af-budget-formatted{font-size:.78rem;color:#6366f1;font-weight:600;margin:6px 0 0}

/* ═══ Actions ═══ */
.af-actions{display:flex;justify-content:space-between;align-items:center;margin-top:28px;padding-top:20px;border-top:1px solid #f1f5f9}
.af-btn-back{display:inline-flex;align-items:center;gap:5px;padding:10px 18px;border:1px solid #e8ecf1;border-radius:10px;background:#fff;font-size:.82rem;font-weight:600;color:#475569;cursor:pointer;transition:all .15s}
.af-btn-back:hover{border-color:#94a3b8;color:#1e293b}
.af-btn-next{display:inline-flex;align-items:center;gap:5px;padding:10px 22px;border:none;border-radius:10px;background:linear-gradient(135deg,#6366f1,#8b5cf6);color:#fff;font-size:.85rem;font-weight:600;cursor:pointer;transition:all .2s;box-shadow:0 2px 8px rgba(99,102,241,.25)}
.af-btn-next:hover:not(:disabled){box-shadow:0 4px 14px rgba(99,102,241,.35);transform:translateY(-1px)}
.af-btn-next:disabled{opacity:.4;cursor:not-allowed}
.af-btn-submit{display:inline-flex;align-items:center;gap:6px;padding:12px 28px;border:none;border-radius:10px;background:linear-gradient(135deg,#22c55e,#16a34a);color:#fff;font-size:.88rem;font-weight:700;cursor:pointer;transition:all .2s;box-shadow:0 2px 8px rgba(34,197,94,.3)}
.af-btn-submit:hover:not(:disabled){box-shadow:0 6px 20px rgba(34,197,94,.35);transform:translateY(-1px)}
.af-btn-submit:disabled{opacity:.5;cursor:not-allowed}
.af-spin{width:14px;height:14px;border:2px solid rgba(255,255,255,.3);border-top-color:#fff;border-radius:50%;animation:spin .7s linear infinite}
@keyframes spin{to{transform:rotate(360deg)}}

/* ═══ Confirm ═══ */
.af-confirm-grid{display:grid;grid-template-columns:1fr 1fr;gap:16px;margin-bottom:20px}
.af-confirm-item{background:#f8fafc;border-radius:10px;padding:14px 16px}
.af-confirm-item.full{grid-column:1/-1}
.af-confirm-label{font-size:.68rem;color:#94a3b8;text-transform:uppercase;letter-spacing:.04em;font-weight:600;margin-bottom:4px}
.af-confirm-val{font-size:.85rem;color:#0f172a;font-weight:600;line-height:1.5;white-space:pre-wrap}
.af-confirm-note{display:flex;align-items:flex-start;gap:10px;padding:14px 16px;background:#fffbeb;border:1px solid #fef3c7;border-radius:10px;font-size:.8rem;color:#92400e;line-height:1.5}
.af-confirm-note svg{flex-shrink:0;margin-top:1px;color:#f59e0b}

@media(max-width:768px){.af-hero{padding:24px 20px 18px}.af-summary{flex-direction:column;align-items:flex-start}.af-summary-divider{width:100%;height:1px}.af-confirm-grid{grid-template-columns:1fr}.af-steps{overflow-x:auto}}
</style>
