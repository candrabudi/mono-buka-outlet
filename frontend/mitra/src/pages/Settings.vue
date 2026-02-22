<template>
  <div class="set-page">
    <!-- Hero -->
    <div class="set-hero">
      <div class="set-hero-top">
        <div>
          <h1 class="set-hero-title">Pengaturan Profil</h1>
          <p class="set-hero-sub">Kelola informasi akun dan keamanan Anda</p>
        </div>
      </div>
    </div>

    <div class="set-grid">
      <!-- Profile Info Card -->
      <div class="set-card">
        <div class="set-card-header">
          <div class="set-card-icon" style="background:#eff6ff;color:#3b82f6">
            <i class="ri-user-line"></i>
          </div>
          <div>
            <h3 class="set-card-title">Informasi Profil</h3>
            <p class="set-card-sub">Ubah nama, email, dan nomor handphone</p>
          </div>
        </div>

        <form @submit.prevent="handleUpdateProfile" class="set-form">
          <div class="set-field">
            <label for="prof-name">Nama Lengkap</label>
            <div class="set-input-wrap">
              <i class="ri-user-line"></i>
              <input id="prof-name" v-model="profile.name" type="text" placeholder="Nama lengkap" />
            </div>
            <span v-if="profileErrors.name" class="set-error">{{ profileErrors.name }}</span>
          </div>

          <div class="set-field">
            <label for="prof-email">Email</label>
            <div class="set-input-wrap">
              <i class="ri-mail-line"></i>
              <input id="prof-email" v-model="profile.email" type="email" placeholder="email@example.com" />
            </div>
            <span v-if="profileErrors.email" class="set-error">{{ profileErrors.email }}</span>
          </div>

          <div class="set-field">
            <label for="prof-phone">No. Handphone</label>
            <div class="set-input-wrap">
              <i class="ri-phone-line"></i>
              <input id="prof-phone" v-model="profile.phone" type="text" placeholder="08xxxxxxxxxx" />
            </div>
            <span v-if="profileErrors.phone" class="set-error">{{ profileErrors.phone }}</span>
          </div>

          <div class="set-actions">
            <button type="submit" class="set-btn set-btn-primary" :disabled="profileLoading || !profileChanged">
              <span v-if="profileLoading" class="set-spinner"></span>
              {{ profileLoading ? 'Menyimpan...' : 'Simpan Perubahan' }}
            </button>
            <button type="button" class="set-btn set-btn-ghost" @click="resetProfile" :disabled="!profileChanged">Batal</button>
          </div>
        </form>
      </div>

      <!-- Change Password Card -->
      <div class="set-card">
        <div class="set-card-header">
          <div class="set-card-icon" style="background:#fef3c7;color:#d97706">
            <i class="ri-lock-line"></i>
          </div>
          <div>
            <h3 class="set-card-title">Ubah Password</h3>
            <p class="set-card-sub">Pastikan password baru kuat dan unik</p>
          </div>
        </div>

        <form @submit.prevent="handleChangePassword" class="set-form">
          <div class="set-field">
            <label for="pw-old">Password Saat Ini</label>
            <div class="set-input-wrap">
              <i class="ri-lock-line"></i>
              <input id="pw-old" v-model="pw.old_password" :type="showOld ? 'text' : 'password'" placeholder="Password saat ini" required />
              <button type="button" class="set-toggle-pw" @click="showOld = !showOld" tabindex="-1">
                <i :class="showOld ? 'ri-eye-line' : 'ri-eye-off-line'" style="color:#94a3b8"></i>
              </button>
            </div>
          </div>

          <div class="set-field">
            <label for="pw-new">Password Baru</label>
            <div class="set-input-wrap">
              <i class="ri-lock-line"></i>
              <input id="pw-new" v-model="pw.new_password" :type="showNew ? 'text' : 'password'" placeholder="Min. 8 karakter" required @input="pwTouched = true" />
              <button type="button" class="set-toggle-pw" @click="showNew = !showNew" tabindex="-1">
                <i :class="showNew ? 'ri-eye-line' : 'ri-eye-off-line'" style="color:#94a3b8"></i>
              </button>
            </div>
            <!-- Strength indicator -->
            <div v-if="pwTouched && pw.new_password" class="set-pw-strength">
              <div class="set-pw-bar"><div class="set-pw-fill" :style="{ width: pwStrength.percent + '%', background: pwStrength.color }"></div></div>
              <span class="set-pw-label" :style="{ color: pwStrength.color }">{{ pwStrength.label }}</span>
            </div>
            <div v-if="pwTouched && pw.new_password" class="set-pw-rules">
              <span :class="{ met: pwChecks.length }">✓ 8+ karakter</span>
              <span :class="{ met: pwChecks.upper }">✓ Huruf besar</span>
              <span :class="{ met: pwChecks.lower }">✓ Huruf kecil</span>
              <span :class="{ met: pwChecks.digit }">✓ Angka</span>
            </div>
          </div>

          <div class="set-field">
            <label for="pw-confirm">Konfirmasi Password Baru</label>
            <div class="set-input-wrap">
              <i class="ri-lock-line"></i>
              <input id="pw-confirm" v-model="pw.confirm_new_password" :type="showConfirm ? 'text' : 'password'" placeholder="Ulangi password baru" required />
              <button type="button" class="set-toggle-pw" @click="showConfirm = !showConfirm" tabindex="-1">
                <i :class="showConfirm ? 'ri-eye-line' : 'ri-eye-off-line'" style="color:#94a3b8"></i>
              </button>
            </div>
            <span v-if="pw.confirm_new_password && pw.new_password !== pw.confirm_new_password" class="set-error">Password tidak cocok</span>
          </div>

          <div class="set-actions">
            <button type="submit" class="set-btn set-btn-warning" :disabled="pwLoading || !canSubmitPw">
              <span v-if="pwLoading" class="set-spinner"></span>
              {{ pwLoading ? 'Mengubah...' : 'Ubah Password' }}
            </button>
          </div>
        </form>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted } from 'vue'
import { authApi } from '../services/api'
import { useAuthStore } from '../stores/auth'
import { useToastStore } from '../stores/toast'

const auth = useAuthStore()
const toast = useToastStore()

// ╔═══════════════════════════════════╗
// ║        Profile Section            ║
// ╚═══════════════════════════════════╝
const profile = reactive({ name: '', email: '', phone: '' })
const originalProfile = reactive({ name: '', email: '', phone: '' })
const profileLoading = ref(false)

const profileChanged = computed(() =>
  profile.name !== originalProfile.name ||
  profile.email !== originalProfile.email ||
  profile.phone !== originalProfile.phone
)

const profileErrors = computed(() => {
  const e = {}
  if (profile.name && profile.name.trim().length < 3) e.name = 'Nama minimal 3 karakter'
  if (profile.email && !/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(profile.email)) e.email = 'Format email tidak valid'
  if (profile.phone) {
    const phone = profile.phone.replace(/[\s-]/g, '')
    if (!/^(\+62|62|08)[0-9]{8,13}$/.test(phone)) e.phone = 'Format: 08xx atau +62xx'
  }
  return e
})

async function loadProfile() {
  try {
    const { data } = await authApi.profile()
    const u = data.data
    profile.name = u.name || ''
    profile.email = u.email || ''
    profile.phone = u.phone || ''
    Object.assign(originalProfile, { name: profile.name, email: profile.email, phone: profile.phone })
  } catch { /* handled by interceptor */ }
}

function resetProfile() {
  Object.assign(profile, { ...originalProfile })
}

async function handleUpdateProfile() {
  if (Object.keys(profileErrors.value).length > 0) {
    toast.error('Periksa kembali data yang diisi')
    return
  }
  profileLoading.value = true
  try {
    const { data } = await authApi.updateProfile(profile)
    const u = data.data
    auth.setAuth(auth.token, u)
    Object.assign(originalProfile, { name: u.name, email: u.email, phone: u.phone })
    profile.name = u.name
    profile.email = u.email
    profile.phone = u.phone
    toast.success(data.message || 'Profil berhasil diupdate')
  } catch (e) {
    toast.error(e.response?.data?.error || 'Gagal update profil')
  } finally {
    profileLoading.value = false
  }
}

// ╔═══════════════════════════════════╗
// ║       Password Section            ║
// ╚═══════════════════════════════════╝
const pw = reactive({ old_password: '', new_password: '', confirm_new_password: '' })
const pwLoading = ref(false)
const pwTouched = ref(false)
const showOld = ref(false)
const showNew = ref(false)
const showConfirm = ref(false)

const pwChecks = computed(() => ({
  length: pw.new_password.length >= 8,
  upper: /[A-Z]/.test(pw.new_password),
  lower: /[a-z]/.test(pw.new_password),
  digit: /[0-9]/.test(pw.new_password),
}))

const pwStrength = computed(() => {
  const checks = Object.values(pwChecks.value).filter(Boolean).length
  if (checks <= 1) return { percent: 25, color: '#ef4444', label: 'Lemah' }
  if (checks === 2) return { percent: 50, color: '#f59e0b', label: 'Sedang' }
  if (checks === 3) return { percent: 75, color: '#3b82f6', label: 'Baik' }
  return { percent: 100, color: '#22c55e', label: 'Kuat' }
})

const canSubmitPw = computed(() =>
  pw.old_password.length > 0 &&
  pw.new_password.length >= 8 &&
  pwChecks.value.upper && pwChecks.value.lower && pwChecks.value.digit &&
  pw.new_password === pw.confirm_new_password
)

async function handleChangePassword() {
  if (!canSubmitPw.value) return
  pwLoading.value = true
  try {
    const { data } = await authApi.changePassword(pw)
    toast.success(data.message || 'Password berhasil diubah')
    pw.old_password = ''
    pw.new_password = ''
    pw.confirm_new_password = ''
    pwTouched.value = false
  } catch (e) {
    toast.error(e.response?.data?.error || 'Gagal mengubah password')
  } finally {
    pwLoading.value = false
  }
}

onMounted(loadProfile)
</script>

<style scoped>
.set-hero { background: linear-gradient(135deg, #0f0c29 0%, #302b63 50%, #24243e 100%); border-radius: 16px; padding: 32px 36px 24px; margin-bottom: 24px; box-shadow: 0 4px 24px rgba(15,12,41,0.2); }
.set-hero-title { font-size: 1.6rem; font-weight: 800; color: #fff; margin: 0 0 4px; }
.set-hero-sub { font-size: .85rem; color: rgba(255,255,255,.5); margin: 0; }

.set-grid { display: grid; grid-template-columns: 1fr 1fr; gap: 24px; align-items: start; }

.set-card { background: #fff; border-radius: 16px; border: 1px solid #e8ecf1; padding: 28px; transition: box-shadow .2s; }
.set-card:hover { box-shadow: 0 4px 20px rgba(0,0,0,.04); }
.set-card-header { display: flex; align-items: flex-start; gap: 14px; margin-bottom: 24px; }
.set-card-icon { width: 42px; height: 42px; border-radius: 12px; display: flex; align-items: center; justify-content: center; flex-shrink: 0; }
.set-card-title { font-size: 1rem; font-weight: 700; color: #0f172a; margin: 0 0 2px; }
.set-card-sub { font-size: .78rem; color: #94a3b8; margin: 0; }

.set-form { display: flex; flex-direction: column; gap: 18px; }
.set-field { display: flex; flex-direction: column; gap: 6px; }
.set-field label { font-size: .82rem; font-weight: 600; color: #334155; }
.set-input-wrap { position: relative; }
.set-input-icon { position: absolute; left: 14px; top: 50%; transform: translateY(-50%); color: #94a3b8; pointer-events: none; }
.set-input-wrap input { width: 100%; height: 48px; border: 1px solid #e2e8f0; border-radius: 12px; padding: 0 42px 0 40px; font-size: .88rem; font-weight: 500; color: #1e293b; background: #fafbfc; transition: all .2s; box-sizing: border-box; font-family: inherit; }
.set-input-wrap input::placeholder { color: #94a3b8; font-weight: 400; }
.set-input-wrap input:focus { outline: none; border-color: #818cf8; box-shadow: 0 0 0 3px rgba(129,140,248,.12); background: #fff; }
.set-toggle-pw { position: absolute; right: 12px; top: 50%; transform: translateY(-50%); background: none; border: none; cursor: pointer; padding: 4px; display: flex; }
.set-error { font-size: .75rem; color: #ef4444; font-weight: 500; }

.set-actions { display: flex; gap: 10px; padding-top: 6px; }
.set-btn { display: inline-flex; align-items: center; justify-content: center; gap: 6px; padding: 10px 22px; border: none; border-radius: 10px; font-size: .85rem; font-weight: 600; cursor: pointer; transition: all .2s; font-family: inherit; }
.set-btn:disabled { opacity: .45; cursor: not-allowed; }
.set-btn-primary { background: linear-gradient(135deg, #6366f1, #8b5cf6); color: #fff; box-shadow: 0 2px 8px rgba(99,102,241,.25); }
.set-btn-primary:hover:not(:disabled) { box-shadow: 0 4px 14px rgba(99,102,241,.35); transform: translateY(-1px); }
.set-btn-warning { background: linear-gradient(135deg, #f59e0b, #d97706); color: #fff; box-shadow: 0 2px 8px rgba(245,158,11,.25); }
.set-btn-warning:hover:not(:disabled) { box-shadow: 0 4px 14px rgba(245,158,11,.35); transform: translateY(-1px); }
.set-btn-ghost { background: #f1f5f9; color: #475569; }
.set-btn-ghost:hover:not(:disabled) { background: #e2e8f0; }
.set-spinner { width: 16px; height: 16px; border: 2px solid rgba(255,255,255,.3); border-top-color: #fff; border-radius: 50%; animation: spin .7s linear infinite; }
@keyframes spin { to { transform: rotate(360deg); } }

/* Password strength */
.set-pw-strength { display: flex; align-items: center; gap: 8px; }
.set-pw-bar { flex: 1; height: 4px; background: #e2e8f0; border-radius: 2px; overflow: hidden; }
.set-pw-fill { height: 100%; border-radius: 2px; transition: all .3s; }
.set-pw-label { font-size: .72rem; font-weight: 700; min-width: 44px; }
.set-pw-rules { display: flex; flex-wrap: wrap; gap: 4px 10px; }
.set-pw-rules span { font-size: .68rem; color: #cbd5e1; font-weight: 600; transition: color .2s; }
.set-pw-rules span.met { color: #22c55e; }

@media (max-width: 900px) {
  .set-grid { grid-template-columns: 1fr; }
  .set-hero { padding: 24px 20px 18px; }
}
</style>
