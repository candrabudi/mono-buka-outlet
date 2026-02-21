<template>
  <section class="login-section">
    <div class="login-wrapper">
      <!-- ═══ LEFT: Form Area ═══ -->
      <div class="login-left">
        <header class="login-logo">
          <div class="logo-icon">
            <svg width="32" height="32" viewBox="0 0 32 32" fill="none" xmlns="http://www.w3.org/2000/svg">
              <rect width="32" height="32" rx="8" fill="#fd9644"/>
              <path d="M9 10H23V12H17V22H15V12H9V10ZM19 14H23V16H19V14Z" fill="white"/>
            </svg>
          </div>
          <span class="logo-text">BukaOutlet</span>
        </header>

        <div class="login-form-wrapper">
          <!-- ═══ STEP 1: Login Form ═══ -->
          <div v-if="auth.otpStep === 'login'" class="login-form-content animate-fade-in">
            <header class="form-header">
              <h2>Masuk ke Admin Panel</h2>
              <p>Kelola kemitraan franchise dengan lebih mudah</p>
            </header>

            <form @submit.prevent="handleLogin">
              <div class="input-group">
                <label for="email">Email</label>
                <div class="input-wrapper">
                  <svg class="input-icon" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
                    <path d="M3 8l7.89 5.26a2 2 0 002.22 0L21 8M5 19h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v10a2 2 0 002 2z"/>
                  </svg>
                  <input
                    id="email"
                    v-model="form.email"
                    type="email"
                    placeholder="admin@franchise.com"
                    required
                    autocomplete="email"
                  />
                </div>
              </div>

              <div class="input-group">
                <label for="password">Password</label>
                <div class="input-wrapper">
                  <svg class="input-icon" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
                    <rect x="3" y="11" width="18" height="11" rx="2" ry="2"/>
                    <path d="M7 11V7a5 5 0 0110 0v4"/>
                  </svg>
                  <input
                    id="password"
                    v-model="form.password"
                    :type="showPassword ? 'text' : 'password'"
                    placeholder="••••••••"
                    required
                    autocomplete="current-password"
                  />
                  <button type="button" class="toggle-password" @click="showPassword = !showPassword" tabindex="-1">
                    <svg v-if="!showPassword" width="20" height="20" viewBox="0 0 22 20" fill="none">
                      <path d="M2 1L20 19" stroke="#718096" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round"/>
                      <path d="M9.58 8.59C9.21 8.96 9 9.47 9 10c0 .53.21 1.04.58 1.42.38.37.89.58 1.42.58.53 0 1.04-.21 1.41-.58" stroke="#718096" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round"/>
                      <path d="M8.36 3.37C9.22 3.12 10.11 3 11 3c4 0 7.33 2.33 10 7-.78 1.36-1.61 2.52-2.5 3.49M16.36 15.35C14.73 16.45 12.94 17 11 17c-4 0-7.33-2.33-10-7 1.37-2.4 2.91-4.17 4.63-5.34" stroke="#718096" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round"/>
                    </svg>
                    <svg v-else width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="#718096" stroke-width="1.5">
                      <path d="M1 12s4-8 11-8 11 8 11 8-4 8-11 8-11-8-11-8z"/>
                      <circle cx="12" cy="12" r="3"/>
                    </svg>
                  </button>
                </div>
              </div>



              <button type="submit" class="btn-login" :disabled="loading">
                <span v-if="loading" class="spinner"></span>
                {{ loading ? 'Memverifikasi...' : 'Masuk' }}
              </button>
            </form>



            <nav class="login-footer-links">
              <a href="#">Terms & Condition</a>
              <a href="#">Privacy Policy</a>
              <a href="#">Help</a>
            </nav>
            <p class="login-copyright">© 2026 BukaOutlet. All Rights Reserved.</p>
          </div>

          <!-- ═══ STEP 2: OTP Verification ═══ -->
          <div v-else class="login-form-content animate-fade-in">
            <header class="form-header">
              <div class="otp-icon-wrapper">
                <svg width="28" height="28" viewBox="0 0 24 24" fill="none" stroke="#fd9644" stroke-width="1.5">
                  <path d="M22 16.92v3a2 2 0 01-2.18 2 19.79 19.79 0 01-8.63-3.07 19.5 19.5 0 01-6-6 19.79 19.79 0 01-3.07-8.67A2 2 0 014.11 2h3a2 2 0 012 1.72c.13.81.37 1.6.65 2.36a2 2 0 01-.45 2.11L8.09 9.91a16 16 0 006 6l1.27-1.27a2 2 0 012.11-.45c.76.28 1.55.52 2.36.65a2 2 0 011.72 2z"/>
                </svg>
              </div>
              <h2>Verifikasi OTP</h2>
              <p>Kode 6 digit telah dikirim ke <strong>{{ auth.otpEmail }}</strong></p>
            </header>

            <form @submit.prevent="handleVerifyOtp">
              <div class="otp-inputs">
                <input
                  v-for="(_, i) in 6" :key="i"
                  ref="otpInputs"
                  type="text"
                  inputmode="numeric"
                  maxlength="1"
                  class="otp-box"
                  :value="otpDigits[i]"
                  @input="onOtpInput($event, i)"
                  @keydown="onOtpKeydown($event, i)"
                  @paste="onOtpPaste"
                  @focus="$event.target.select()"
                />
              </div>



              <button type="submit" class="btn-login" :disabled="loading || otpCode.length < 6">
                <span v-if="loading" class="spinner"></span>
                {{ loading ? 'Memverifikasi...' : 'Verifikasi' }}
              </button>
            </form>

            <div class="otp-actions">
              <button @click="handleResendOtp" class="resend-btn" :disabled="resendCooldown > 0">
                {{ resendCooldown > 0 ? `Kirim ulang (${resendCooldown}s)` : 'Kirim ulang kode' }}
              </button>
              <button @click="auth.cancelOtp(); error = ''; success = ''" class="back-btn">
                ← Kembali ke login
              </button>
            </div>
          </div>
        </div>
      </div>

      <!-- ═══ RIGHT: Illustration ═══ -->
      <div class="login-right">
        <ul class="floating-shapes">
          <li class="shape shape-1"><img src="/images/shapes/square.svg" alt="" /></li>
          <li class="shape shape-2"><img src="/images/shapes/vline.svg" alt="" /></li>
          <li class="shape shape-3"><img src="/images/shapes/dotted.svg" alt="" /></li>
        </ul>
        <div class="illustration-wrapper">
          <img src="/images/illustration/signin.svg" alt="Franchise Management" />
        </div>
        <div class="illustration-text">
          <h3>Cepat, Mudah & Terpercaya</h3>
          <p>
            BukaOutlet membantu Anda mengelola seluruh proses kemitraan —
            dari leads, perjanjian, hingga laporan keuangan. Semua dalam satu platform
            yang <span class="highlight">terintegrasi</span>.
          </p>
        </div>
      </div>
    </div>
  </section>
</template>

<script setup>
import { ref, reactive, computed, onUnmounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useAuthStore } from '../stores/auth'
import { useToastStore } from '../stores/toast'

const auth = useAuthStore()
const toast = useToastStore()
const router = useRouter()
const route = useRoute()
const loading = ref(false)
const showPassword = ref(false)
const form = reactive({ email: '', password: '' })
const otpDigits = ref(['', '', '', '', '', ''])
const otpInputs = ref([])
const resendCooldown = ref(0)
let cooldownTimer = null

const otpCode = computed(() => otpDigits.value.join(''))

async function handleLogin() {
  loading.value = true
  try {
    const data = await auth.login(form)
    toast.success(data.message || 'OTP telah dikirim ke email Anda')
    startResendCooldown()
  } catch (e) {
    const msg = e.response?.data?.error || (e.code === 'ERR_NETWORK' ? 'Server tidak dapat dihubungi' : e.message || 'Login gagal')
    toast.error(msg)
  } finally {
    loading.value = false
  }
}

async function handleVerifyOtp() {
  loading.value = true
  try {
    await auth.verifyOtp(otpCode.value)
    toast.success('Login berhasil!')
    router.push(route.query.redirect || '/')
  } catch (e) {
    const msg = e.response?.data?.error || (e.code === 'ERR_NETWORK' ? 'Server tidak dapat dihubungi' : e.message || 'Kode OTP tidak valid')
    toast.error(msg)
    otpDigits.value = ['', '', '', '', '', '']
    otpInputs.value[0]?.focus()
  } finally {
    loading.value = false
  }
}

async function handleResendOtp() {
  if (resendCooldown.value > 0) return
  try {
    const data = await auth.resendOtp()
    toast.success(data.message || 'Kode OTP baru telah dikirim')
    otpDigits.value = ['', '', '', '', '', '']
    startResendCooldown()
  } catch (e) {
    const msg = e.response?.data?.error || (e.code === 'ERR_NETWORK' ? 'Server tidak dapat dihubungi' : 'Gagal mengirim ulang OTP')
    toast.error(msg)
  }
}

function startResendCooldown() {
  resendCooldown.value = 60
  cooldownTimer = setInterval(() => {
    resendCooldown.value--
    if (resendCooldown.value <= 0) clearInterval(cooldownTimer)
  }, 1000)
}

function onOtpInput(e, i) {
  const val = e.target.value.replace(/\D/g, '')
  otpDigits.value[i] = val.charAt(0) || ''
  e.target.value = otpDigits.value[i]
  if (val && i < 5) otpInputs.value[i + 1]?.focus()
}

function onOtpKeydown(e, i) {
  if (e.key === 'Backspace' && !otpDigits.value[i] && i > 0) {
    otpDigits.value[i - 1] = ''
    otpInputs.value[i - 1]?.focus()
  }
}

function onOtpPaste(e) {
  const paste = e.clipboardData.getData('text').replace(/\D/g, '').slice(0, 6)
  paste.split('').forEach((ch, i) => { otpDigits.value[i] = ch })
  if (paste.length === 6) otpInputs.value[5]?.focus()
  e.preventDefault()
}

onUnmounted(() => { if (cooldownTimer) clearInterval(cooldownTimer) })
</script>

<style scoped>
@import url('https://fonts.googleapis.com/css2?family=Poppins:wght@300;400;500;600;700&family=Urbanist:wght@300;400;500;600;700&display=swap');

/* ── Base ── */
.login-section {
  background: #fff;
  min-height: 100vh;
  font-family: 'Urbanist', sans-serif;
}
.login-wrapper {
  display: flex;
  min-height: 100vh;
}

/* ── Left Panel ── */
.login-left {
  flex: 1;
  padding: 40px 20px 40px 48px;
  display: flex;
  flex-direction: column;
  max-width: 50%;
}
.login-logo {
  display: flex;
  align-items: center;
  gap: 10px;
}
.logo-icon {
  display: flex;
  align-items: center;
  justify-content: center;
}
.logo-text {
  font-family: 'Poppins', sans-serif;
  font-size: 20px;
  font-weight: 700;
  color: #1a202c;
}
.login-form-wrapper {
  max-width: 450px;
  margin: auto;
  padding: 60px 0 40px;
  width: 100%;
}

/* Form Header */
.form-header {
  text-align: center;
  margin-bottom: 32px;
}
.form-header h2 {
  font-family: 'Poppins', sans-serif;
  font-size: 28px;
  font-weight: 600;
  color: #1a202c;
  margin: 0 0 8px;
  line-height: 1.3;
}
.form-header p {
  font-size: 15px;
  font-weight: 500;
  color: #718096;
  margin: 0;
  line-height: 1.5;
}
.form-header p strong {
  color: #fd9644;
  word-break: break-all;
}

/* Input Group */
.input-group {
  margin-bottom: 20px;
}
.input-group label {
  display: block;
  font-size: 14px;
  font-weight: 600;
  color: #1a202c;
  margin-bottom: 8px;
}
.input-wrapper {
  position: relative;
}
.input-icon {
  position: absolute;
  left: 16px;
  top: 50%;
  transform: translateY(-50%);
  color: #a0aec0;
  pointer-events: none;
}
.input-wrapper input {
  width: 100%;
  height: 56px;
  border: 1px solid #e2e8f0;
  border-radius: 12px;
  padding: 0 48px 0 48px;
  font-size: 15px;
  font-weight: 500;
  color: #1a202c;
  background: #fff;
  transition: all 0.2s ease;
  font-family: 'Urbanist', sans-serif;
  box-sizing: border-box;
}
.input-wrapper input::placeholder {
  color: #a0aec0;
  font-weight: 400;
}
.input-wrapper input:focus {
  outline: none;
  border-color: #fd9644;
  box-shadow: 0 0 0 3px rgba(253, 150, 68, 0.15);
}
.toggle-password {
  position: absolute;
  right: 16px;
  top: 50%;
  transform: translateY(-50%);
  background: none;
  border: none;
  cursor: pointer;
  padding: 0;
  display: flex;
  align-items: center;
}

/* Alerts */
.alert {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 12px 16px;
  border-radius: 10px;
  font-size: 14px;
  font-weight: 500;
  margin-bottom: 20px;
}
.alert-error {
  background: #FEF2F2;
  border: 1px solid #FECACA;
  color: #DC2626;
}
.alert-success {
  background: #F0FDF4;
  border: 1px solid #BBF7D0;
  color: #16A34A;
}

/* Login Button */
.btn-login {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  width: 100%;
  height: 52px;
  border: none;
  border-radius: 12px;
  background: linear-gradient(135deg, #fd9644 0%, #fc5c65 100%);
  color: #fff;
  font-size: 16px;
  font-weight: 600;
  font-family: 'Urbanist', sans-serif;
  cursor: pointer;
  transition: all 0.25s ease;
  margin-top: 4px;
}
.btn-login:hover:not(:disabled) {
  background: linear-gradient(135deg, #e8832e 0%, #e04e55 100%);
  transform: translateY(-1px);
  box-shadow: 0 4px 12px rgba(253, 150, 68, 0.35);
}
.btn-login:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}
.spinner {
  width: 18px;
  height: 18px;
  border: 2.5px solid rgba(255,255,255,0.3);
  border-top-color: #fff;
  border-radius: 50%;
  animation: spin 0.7s linear infinite;
}
@keyframes spin { to { transform: rotate(360deg); } }

/* OTP Inputs */
.otp-icon-wrapper {
  width: 56px;
  height: 56px;
  border-radius: 16px;
  background: #FFF8F1;
  display: flex;
  align-items: center;
  justify-content: center;
  margin: 0 auto 16px;
}
.otp-inputs {
  display: flex;
  gap: 12px;
  justify-content: center;
  margin-bottom: 24px;
}
.otp-box {
  width: 56px;
  height: 60px;
  text-align: center;
  font-size: 24px;
  font-weight: 700;
  font-family: 'Poppins', sans-serif;
  color: #1a202c;
  background: #F5F5F5;
  border: 2px solid transparent;
  border-radius: 14px;
  outline: none;
  transition: all 0.2s ease;
  caret-color: #fd9644;
  box-sizing: border-box;
}
.otp-box:focus {
  border-color: #fd9644;
  background: #fff;
  box-shadow: 0 0 0 3px rgba(253, 150, 68, 0.15);
}
.otp-actions {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-top: 24px;
}
.resend-btn {
  background: none;
  border: none;
  color: #fd9644;
  font-size: 14px;
  font-weight: 700;
  cursor: pointer;
  padding: 0;
  font-family: 'Urbanist', sans-serif;
}
.resend-btn:disabled {
  color: #a0aec0;
  cursor: not-allowed;
}
.back-btn {
  background: none;
  border: none;
  color: #718096;
  font-size: 14px;
  font-weight: 500;
  cursor: pointer;
  padding: 0;
  font-family: 'Urbanist', sans-serif;
  transition: color 0.2s;
}
.back-btn:hover { color: #1a202c; }

/* Demo Box */
.demo-box {
  margin-top: 32px;
  padding: 16px 20px;
  background: #FFF8F1;
  border-radius: 12px;
  border: 1px solid #FEECDB;
}
.demo-title {
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 13px;
  font-weight: 700;
  color: #e8832e;
  margin-bottom: 10px;
}
.demo-grid {
  display: grid;
  grid-template-columns: auto 1fr;
  gap: 4px 12px;
  font-size: 13px;
}
.demo-label {
  font-weight: 600;
  color: #c4680e;
}
.demo-value {
  color: #9a5316;
  font-family: monospace;
  font-size: 12px;
}
.demo-pw {
  color: #a0aec0 !important;
  font-weight: 500 !important;
  font-size: 12px;
}

/* Footer */
.login-footer-links {
  display: flex;
  justify-content: center;
  flex-wrap: wrap;
  gap: 24px;
  padding-top: 48px;
}
.login-footer-links a {
  font-size: 13px;
  color: #718096;
  text-decoration: none;
  transition: color 0.2s;
}
.login-footer-links a:hover { color: #1a202c; }
.login-copyright {
  text-align: center;
  font-size: 13px;
  color: #a0aec0;
  margin-top: 16px;
}

/* ── Right Panel ── */
.login-right {
  flex: 1;
  background: #FFF9F3;
  padding: 80px;
  position: relative;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  overflow: hidden;
}
.floating-shapes {
  list-style: none;
  padding: 0;
  margin: 0;
}
.shape {
  position: absolute;
  opacity: 0.5;
  animation: float 6s ease-in-out infinite;
}
.shape-1 { top: 40px; left: 32px; animation-delay: 0s; }
.shape-2 { top: 56px; right: 48px; animation-delay: 2s; }
.shape-3 { bottom: 28px; left: 32px; animation-delay: 4s; }
@keyframes float {
  0%, 100% { transform: translateY(0); }
  50% { transform: translateY(-12px); }
}
.illustration-wrapper {
  width: 100%;
  max-width: 580px;
}
.illustration-wrapper img {
  width: 100%;
  height: auto;
}
.illustration-text {
  text-align: center;
  max-width: 480px;
  padding: 0 8px;
  margin-top: 24px;
}
.illustration-text h3 {
  font-family: 'Poppins', sans-serif;
  font-size: 30px;
  font-weight: 600;
  color: #1a202c;
  margin: 0 0 12px;
}
.illustration-text p {
  font-size: 14px;
  font-weight: 500;
  color: #718096;
  line-height: 1.7;
  margin: 0;
}
.illustration-text .highlight {
  color: #fd9644;
  font-weight: 700;
}

/* ── Animation ── */
.animate-fade-in {
  animation: fadeSlideUp 0.4s ease-out;
}
@keyframes fadeSlideUp {
  from { opacity: 0; transform: translateY(16px); }
  to   { opacity: 1; transform: translateY(0); }
}

/* ── Responsive ── */
@media (max-width: 1024px) {
  .login-right { display: none; }
  .login-left {
    max-width: 100%;
    padding: 32px 20px;
  }
}
@media (max-width: 480px) {
  .form-header h2 { font-size: 22px; }
  .otp-inputs { gap: 8px; }
  .otp-box { width: 46px; height: 52px; font-size: 20px; }
  .login-form-wrapper { padding: 40px 0 24px; }
}
</style>
