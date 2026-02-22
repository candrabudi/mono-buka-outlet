<template>
  <section class="login-section">
    <div class="login-wrapper">
      <!-- Left: Form -->
      <div class="login-left">
        <header class="login-logo">
          <div class="logo-icon">
            <i class="ri-store-2-fill" style="font-size:32px;color:#fff;background:#fd9644;border-radius:8px;width:32px;height:32px;display:flex;align-items:center;justify-content:center"></i>
          </div>
          <span class="logo-text">BukaOutlet</span>
        </header>

        <div class="login-form-wrapper">
          <div class="login-form-content animate-fade-in">
            <header class="form-header">
              <h2>Daftar Akun Mitra</h2>
              <p>Mulai perjalanan kemitraan Anda</p>
            </header>

            <form @submit.prevent="handleRegister">
              <div class="input-group">
                <label for="name">Nama Lengkap</label>
                <div class="input-wrapper">
                  <i class="input-icon ri-user-line"></i>
                  <input id="name" v-model="form.name" type="text" placeholder="Nama lengkap Anda" required @blur="touched.name = true" />
                </div>
                <span v-if="touched.name && errors.name" class="field-error">{{ errors.name }}</span>
              </div>

              <div class="input-group">
                <label for="email">Email</label>
                <div class="input-wrapper">
                  <i class="input-icon ri-mail-line"></i>
                  <input id="email" v-model="form.email" type="email" placeholder="mitra@email.com" required @blur="touched.email = true" />
                </div>
                <span v-if="touched.email && errors.email" class="field-error">{{ errors.email }}</span>
              </div>

              <div class="input-group">
                <label for="phone">No. Handphone</label>
                <div class="input-wrapper">
                  <i class="input-icon ri-phone-line"></i>
                  <input id="phone" v-model="form.phone" type="text" placeholder="08xxxxxxxxxx" @blur="touched.phone = true" />
                </div>
                <span v-if="touched.phone && errors.phone" class="field-error">{{ errors.phone }}</span>
              </div>

              <div class="input-group">
                <label for="password">Password</label>
                <div class="input-wrapper">
                  <i class="input-icon ri-lock-line"></i>
                  <input id="password" v-model="form.password" :type="showPw ? 'text' : 'password'" placeholder="Min. 8 karakter" required @input="touched.password = true" />
                  <button type="button" class="toggle-password" @click="showPw = !showPw" tabindex="-1">
                    <i :class="showPw ? 'ri-eye-line' : 'ri-eye-off-line'" style="color:#718096"></i>
                  </button>
                </div>
                <!-- Password Strength Indicator -->
                <div v-if="touched.password && form.password" class="pw-strength">
                  <div class="pw-strength-bar">
                    <div class="pw-strength-fill" :style="{ width: pwStrength.percent + '%', background: pwStrength.color }"></div>
                  </div>
                  <span class="pw-strength-label" :style="{ color: pwStrength.color }">{{ pwStrength.label }}</span>
                </div>
                <div v-if="touched.password && form.password" class="pw-rules">
                  <span :class="{ met: pwChecks.length }">✓ Min. 8 karakter</span>
                  <span :class="{ met: pwChecks.upper }">✓ Huruf besar</span>
                  <span :class="{ met: pwChecks.lower }">✓ Huruf kecil</span>
                  <span :class="{ met: pwChecks.digit }">✓ Angka</span>
                </div>
              </div>

              <div class="input-group">
                <label for="confirm_password">Konfirmasi Password</label>
                <div class="input-wrapper">
                  <i class="input-icon ri-lock-line"></i>
                  <input id="confirm_password" v-model="form.confirm_password" :type="showPwConfirm ? 'text' : 'password'" placeholder="Ulangi password" required @blur="touched.confirm = true" />
                  <button type="button" class="toggle-password" @click="showPwConfirm = !showPwConfirm" tabindex="-1">
                    <i :class="showPwConfirm ? 'ri-eye-line' : 'ri-eye-off-line'" style="color:#718096"></i>
                  </button>
                </div>
                <span v-if="touched.confirm && errors.confirm" class="field-error">{{ errors.confirm }}</span>
              </div>

              <button type="submit" class="btn-login" :disabled="loading || !canSubmit">
                <span v-if="loading" class="spinner"></span>
                {{ loading ? 'Mendaftar...' : 'Daftar Sekarang' }}
              </button>
            </form>

            <div class="register-login-link">
              Sudah punya akun? <router-link to="/login">Masuk di sini</router-link>
            </div>

            <nav class="login-footer-links">
              <a href="#">Terms & Condition</a>
              <a href="#">Privacy Policy</a>
              <a href="#">Help</a>
            </nav>
            <p class="login-copyright">© 2026 BukaOutlet. All Rights Reserved.</p>
          </div>
        </div>
      </div>

      <!-- Right: Illustration -->
      <div class="login-right">
        <ul class="floating-shapes">
          <li class="shape shape-1"><img src="/images/shapes/square.svg" alt="" /></li>
          <li class="shape shape-2"><img src="/images/shapes/vline.svg" alt="" /></li>
          <li class="shape shape-3"><img src="/images/shapes/dotted.svg" alt="" /></li>
        </ul>
        <div class="illustration-wrapper">
          <img src="/images/illustration/signin.svg" alt="Franchise" />
        </div>
        <div class="illustration-text">
          <h3>Bergabung Jadi Mitra</h3>
          <p>Daftarkan diri Anda dan mulai perjalanan kemitraan bersama BukaOutlet. Kelola outlet Anda dengan platform yang <span class="highlight">terintegrasi</span>.</p>
        </div>
      </div>
    </div>
  </section>
</template>

<script setup>
import { ref, reactive, computed } from 'vue'
import { useRouter } from 'vue-router'
import { authApi } from '../services/api'
import { useToastStore } from '../stores/toast'

const toast = useToastStore()
const router = useRouter()
const loading = ref(false)
const showPw = ref(false)
const showPwConfirm = ref(false)
const form = reactive({ name: '', email: '', phone: '', password: '', confirm_password: '' })
const touched = reactive({ name: false, email: false, phone: false, password: false, confirm: false })

// Password strength checks
const pwChecks = computed(() => ({
  length: form.password.length >= 8,
  upper: /[A-Z]/.test(form.password),
  lower: /[a-z]/.test(form.password),
  digit: /[0-9]/.test(form.password),
}))

const pwStrength = computed(() => {
  const checks = Object.values(pwChecks.value).filter(Boolean).length
  if (checks <= 1) return { percent: 25, color: '#ef4444', label: 'Lemah' }
  if (checks === 2) return { percent: 50, color: '#f59e0b', label: 'Sedang' }
  if (checks === 3) return { percent: 75, color: '#3b82f6', label: 'Baik' }
  return { percent: 100, color: '#22c55e', label: 'Kuat' }
})

// Field errors
const errors = computed(() => {
  const e = {}
  if (form.name && form.name.trim().length < 3) e.name = 'Nama lengkap minimal 3 karakter'
  if (form.email && !/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(form.email)) e.email = 'Format email tidak valid'
  if (form.phone) {
    const phone = form.phone.replace(/[\s-]/g, '')
    if (!/^(\+62|62|08)[0-9]{8,13}$/.test(phone)) e.phone = 'Format nomor tidak valid (08xx atau +62xx)'
  }
  if (form.confirm_password && form.password !== form.confirm_password) e.confirm = 'Password tidak cocok'
  return e
})

const canSubmit = computed(() => {
  return form.name.trim().length >= 3
    && /^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(form.email)
    && form.password.length >= 8
    && form.password === form.confirm_password
    && Object.keys(errors.value).length === 0
})

async function handleRegister() {
  if (!canSubmit.value) {
    toast.error('Mohon lengkapi semua field dengan benar')
    return
  }
  loading.value = true
  try {
    const { data } = await authApi.register(form)
    toast.success(data.message || 'Registrasi berhasil! Silakan login.')
    router.push('/login')
  } catch (e) {
    const resp = e.response
    if (resp?.status === 429) {
      const retryMs = resp.data?.retry_after_ms || 120000
      toast.error(`Terlalu banyak percobaan. Coba lagi dalam ${Math.ceil(retryMs / 1000)} detik`)
    } else {
      toast.error(resp?.data?.error || 'Registrasi gagal')
    }
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
@import url('https://fonts.googleapis.com/css2?family=Poppins:wght@300;400;500;600;700&family=Urbanist:wght@300;400;500;600;700&display=swap');
.login-section { background: #fff; min-height: 100vh; font-family: 'Urbanist', sans-serif; }
.login-wrapper { display: flex; min-height: 100vh; }
.login-left { flex: 1; padding: 40px 20px 40px 48px; display: flex; flex-direction: column; max-width: 50%; overflow-y: auto; }
.login-logo { display: flex; align-items: center; gap: 10px; }
.logo-icon { display: flex; align-items: center; justify-content: center; }
.logo-text { font-family: 'Poppins', sans-serif; font-size: 20px; font-weight: 700; color: #1a202c; }
.login-form-wrapper { max-width: 450px; margin: auto; padding: 32px 0; width: 100%; }
.form-header { text-align: center; margin-bottom: 28px; }
.form-header h2 { font-family: 'Poppins', sans-serif; font-size: 28px; font-weight: 600; color: #1a202c; margin: 0 0 8px; }
.form-header p { font-size: 15px; font-weight: 500; color: #718096; margin: 0; }
.input-group { margin-bottom: 16px; }
.input-group label { display: block; font-size: 14px; font-weight: 600; color: #1a202c; margin-bottom: 6px; }
.input-wrapper { position: relative; }
.input-icon { position: absolute; left: 16px; top: 50%; transform: translateY(-50%); color: #a0aec0; pointer-events: none; }
.input-wrapper input { width: 100%; height: 52px; border: 1px solid #e2e8f0; border-radius: 12px; padding: 0 48px; font-size: 15px; font-weight: 500; color: #1a202c; background: #fff; transition: all .2s; font-family: 'Urbanist', sans-serif; box-sizing: border-box; }
.input-wrapper input::placeholder { color: #a0aec0; font-weight: 400; }
.input-wrapper input:focus { outline: none; border-color: #fd9644; box-shadow: 0 0 0 3px rgba(253,150,68,.15); }
.toggle-password { position: absolute; right: 16px; top: 50%; transform: translateY(-50%); background: none; border: none; cursor: pointer; padding: 0; display: flex; }
.field-error { display: block; font-size: 12px; color: #ef4444; margin-top: 6px; font-weight: 500; }

/* Password Strength */
.pw-strength { display: flex; align-items: center; gap: 10px; margin-top: 8px; }
.pw-strength-bar { flex: 1; height: 4px; background: #e2e8f0; border-radius: 2px; overflow: hidden; }
.pw-strength-fill { height: 100%; border-radius: 2px; transition: all .3s ease; }
.pw-strength-label { font-size: 12px; font-weight: 700; min-width: 50px; }
.pw-rules { display: flex; flex-wrap: wrap; gap: 6px 12px; margin-top: 8px; }
.pw-rules span { font-size: 11px; color: #cbd5e1; font-weight: 600; transition: color .2s; }
.pw-rules span.met { color: #22c55e; }

.btn-login { display: flex; align-items: center; justify-content: center; gap: 8px; width: 100%; height: 52px; border: none; border-radius: 12px; background: linear-gradient(135deg, #fd9644, #fc5c65); color: #fff; font-size: 16px; font-weight: 600; font-family: 'Urbanist', sans-serif; cursor: pointer; transition: all .25s; margin-top: 8px; }
.btn-login:hover:not(:disabled) { background: linear-gradient(135deg, #e8832e, #e04e55); transform: translateY(-1px); box-shadow: 0 4px 12px rgba(253,150,68,.35); }
.btn-login:disabled { opacity: .5; cursor: not-allowed; }
.spinner { width: 18px; height: 18px; border: 2.5px solid rgba(255,255,255,.3); border-top-color: #fff; border-radius: 50%; animation: spin .7s linear infinite; }
@keyframes spin { to { transform: rotate(360deg); } }
.register-login-link { text-align: center; margin-top: 24px; font-size: 14px; color: #718096; }
.register-login-link a { color: #fd9644; font-weight: 700; text-decoration: none; }
.register-login-link a:hover { text-decoration: underline; }
.login-footer-links { display: flex; justify-content: center; flex-wrap: wrap; gap: 24px; padding-top: 36px; }
.login-footer-links a { font-size: 13px; color: #718096; text-decoration: none; }
.login-footer-links a:hover { color: #1a202c; }
.login-copyright { text-align: center; font-size: 13px; color: #a0aec0; margin-top: 16px; }
.login-right { flex: 1; background: #FFF9F3; padding: 80px; position: relative; display: flex; flex-direction: column; align-items: center; justify-content: center; overflow: hidden; }
.floating-shapes { list-style: none; padding: 0; margin: 0; }
.shape { position: absolute; opacity: 0.5; animation: float 6s ease-in-out infinite; }
.shape-1 { top: 40px; left: 32px; }
.shape-2 { top: 56px; right: 48px; animation-delay: 2s; }
.shape-3 { bottom: 28px; left: 32px; animation-delay: 4s; }
@keyframes float { 0%,100%{transform:translateY(0)} 50%{transform:translateY(-12px)} }
.illustration-wrapper { width: 100%; max-width: 580px; }
.illustration-wrapper img { width: 100%; height: auto; }
.illustration-text { text-align: center; max-width: 480px; margin-top: 24px; }
.illustration-text h3 { font-family: 'Poppins', sans-serif; font-size: 30px; font-weight: 600; color: #1a202c; margin: 0 0 12px; }
.illustration-text p { font-size: 14px; font-weight: 500; color: #718096; line-height: 1.7; margin: 0; }
.illustration-text .highlight { color: #fd9644; font-weight: 700; }
.animate-fade-in { animation: fadeSlideUp .4s ease-out; }
@keyframes fadeSlideUp { from{opacity:0;transform:translateY(16px)} to{opacity:1;transform:translateY(0)} }
@media (max-width: 1024px) { .login-right{display:none} .login-left{max-width:100%;padding:32px 20px} }
@media (max-width: 480px) { .form-header h2{font-size:22px} .login-form-wrapper{padding:24px 0} }
</style>
