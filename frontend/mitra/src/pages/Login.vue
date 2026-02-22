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
          <!-- Step 1: Email + Password -->
          <div v-if="step === 'credentials'" class="login-form-content animate-fade-in" key="step1">
            <header class="form-header">
              <h2>Masuk ke Portal Mitra</h2>
              <p>Pantau kemitraan dan pembayaran Anda</p>
            </header>

            <form @submit.prevent="handleLogin">
              <div class="input-group">
                <label for="email">Email</label>
                <div class="input-wrapper">
                  <i class="input-icon ri-mail-line"></i>
                  <input id="email" v-model="form.email" type="email" placeholder="mitra@email.com" required autocomplete="email" />
                </div>
              </div>

              <div class="input-group">
                <label for="password">Password</label>
                <div class="input-wrapper">
                  <i class="input-icon ri-lock-line"></i>
                  <input id="password" v-model="form.password" :type="showPw ? 'text' : 'password'" placeholder="••••••••" required autocomplete="current-password" />
                  <button type="button" class="toggle-password" @click="showPw = !showPw" tabindex="-1">
                    <i :class="showPw ? 'ri-eye-line' : 'ri-eye-off-line'" style="color:#718096"></i>
                  </button>
                </div>
              </div>

              <button type="submit" class="btn-login" :disabled="loading">
                <span v-if="loading" class="spinner"></span>
                {{ loading ? 'Memverifikasi...' : 'Masuk' }}
              </button>
            </form>

            <div class="register-link">
              Belum punya akun? <router-link to="/register">Daftar di sini</router-link>
            </div>
          </div>

          <!-- Step 2: OTP Verification -->
          <div v-else-if="step === 'otp'" class="login-form-content animate-fade-in" key="step2">
            <header class="form-header">
              <div class="otp-icon-circle">
                <i class="ri-shield-keyhole-line" style="font-size:28px;color:#fd9644"></i>
              </div>
              <h2>Verifikasi OTP</h2>
              <p>Kode verifikasi telah dikirim ke <strong>{{ maskedEmail }}</strong></p>
            </header>

            <form @submit.prevent="handleVerifyOTP">
              <div class="otp-inputs">
                <input
                  v-for="(_, i) in 6" :key="i"
                  :ref="el => otpRefs[i] = el"
                  v-model="otpDigits[i]"
                  type="text" inputmode="numeric" maxlength="1"
                  class="otp-input"
                  @input="onOtpInput(i)"
                  @keydown="onOtpKeydown($event, i)"
                  @paste="onOtpPaste"
                />
              </div>

              <div class="otp-timer" v-if="resendCooldown > 0">
                Kirim ulang dalam <strong>{{ resendCooldown }}s</strong>
              </div>
              <button v-else type="button" class="btn-resend" @click="handleResendOTP" :disabled="resendLoading">
                {{ resendLoading ? 'Mengirim...' : 'Kirim ulang kode OTP' }}
              </button>

              <button type="submit" class="btn-login" :disabled="otpLoading || otpCode.length < 6">
                <span v-if="otpLoading" class="spinner"></span>
                {{ otpLoading ? 'Memverifikasi...' : 'Verifikasi' }}
              </button>
            </form>

            <button class="btn-back" @click="step = 'credentials'">
              <i class="ri-arrow-left-s-line"></i>
              Kembali ke login
            </button>
          </div>

          <nav class="login-footer-links">
            <a href="#">Terms & Condition</a>
            <a href="#">Privacy Policy</a>
            <a href="#">Help</a>
          </nav>
          <p class="login-copyright">© 2026 BukaOutlet. All Rights Reserved.</p>
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
          <h3>Portal Mitra BukaOutlet</h3>
          <p>Pantau perkembangan kemitraan, pembayaran, dan status outlet Anda. Semua informasi dalam satu platform yang <span class="highlight">terintegrasi</span>.</p>
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
import { authApi } from '../services/api'

const auth = useAuthStore()
const toast = useToastStore()
const router = useRouter()
const route = useRoute()

const step = ref('credentials')
const loading = ref(false)
const otpLoading = ref(false)
const resendLoading = ref(false)
const showPw = ref(false)
const form = reactive({ email: '', password: '' })

// OTP state
const otpDigits = reactive(['', '', '', '', '', ''])
const otpRefs = reactive([])
const otpEmail = ref('')
const resendCooldown = ref(0)
let cooldownTimer = null

const otpCode = computed(() => otpDigits.join(''))
const maskedEmail = computed(() => {
  const e = otpEmail.value
  const at = e.indexOf('@')
  if (at <= 2) return '***' + e.substring(at)
  return e.substring(0, 2) + '***' + e.substring(at)
})

function startCooldown(seconds = 60) {
  resendCooldown.value = seconds
  clearInterval(cooldownTimer)
  cooldownTimer = setInterval(() => {
    resendCooldown.value--
    if (resendCooldown.value <= 0) clearInterval(cooldownTimer)
  }, 1000)
}

function onOtpInput(i) {
  const val = otpDigits[i]
  if (val && i < 5) otpRefs[i + 1]?.focus()
}
function onOtpKeydown(e, i) {
  if (e.key === 'Backspace' && !otpDigits[i] && i > 0) {
    otpRefs[i - 1]?.focus()
  }
}
function onOtpPaste(e) {
  const text = (e.clipboardData || window.clipboardData).getData('text').replace(/\D/g, '').slice(0, 6)
  for (let i = 0; i < 6; i++) otpDigits[i] = text[i] || ''
  if (text.length >= 6) otpRefs[5]?.focus()
  e.preventDefault()
}

async function handleLogin() {
  loading.value = true
  try {
    const { data } = await authApi.login(form)
    otpEmail.value = data.data?.email || form.email
    step.value = 'otp'
    startCooldown(60)
    toast.success(data.message || 'Kode OTP telah dikirim ke email Anda')
    // Reset OTP digits
    for (let i = 0; i < 6; i++) otpDigits[i] = ''
    setTimeout(() => otpRefs[0]?.focus(), 200)
  } catch (e) {
    const resp = e.response
    if (resp?.status === 429) {
      const retryMs = resp.data?.retry_after_ms || 120000
      toast.error(`Terlalu banyak percobaan. Coba lagi dalam ${Math.ceil(retryMs / 1000)} detik`)
    } else {
      toast.error(resp?.data?.error || (e.code === 'ERR_NETWORK' ? 'Server tidak dapat dihubungi' : 'Login gagal'))
    }
  } finally {
    loading.value = false
  }
}

async function handleVerifyOTP() {
  if (otpCode.value.length < 6) return
  otpLoading.value = true
  try {
    const { data } = await authApi.verifyOtp({ email: otpEmail.value, code: otpCode.value })
    auth.setAuth(data.data.token, data.data.user)
    toast.success('Login berhasil!')
    router.push(route.query.redirect || '/')
  } catch (e) {
    const resp = e.response
    if (resp?.status === 429) {
      const retryMs = resp.data?.retry_after_ms || 120000
      toast.error(`Terlalu banyak percobaan. Coba lagi dalam ${Math.ceil(retryMs / 1000)} detik`)
    } else {
      toast.error(resp?.data?.error || 'Kode OTP tidak valid')
    }
    for (let i = 0; i < 6; i++) otpDigits[i] = ''
    otpRefs[0]?.focus()
  } finally {
    otpLoading.value = false
  }
}

async function handleResendOTP() {
  resendLoading.value = true
  try {
    const { data } = await authApi.resendOtp({ email: otpEmail.value })
    toast.success(data.message || 'Kode OTP baru telah dikirim')
    startCooldown(60)
    for (let i = 0; i < 6; i++) otpDigits[i] = ''
    otpRefs[0]?.focus()
  } catch (e) {
    toast.error(e.response?.data?.error || 'Gagal mengirim ulang OTP')
  } finally {
    resendLoading.value = false
  }
}

onUnmounted(() => clearInterval(cooldownTimer))
</script>

<style scoped>
@import url('https://fonts.googleapis.com/css2?family=Poppins:wght@300;400;500;600;700&family=Urbanist:wght@300;400;500;600;700&display=swap');
.login-section { background: #fff; min-height: 100vh; font-family: 'Urbanist', sans-serif; }
.login-wrapper { display: flex; min-height: 100vh; }
.login-left { flex: 1; padding: 40px 20px 40px 48px; display: flex; flex-direction: column; max-width: 50%; }
.login-logo { display: flex; align-items: center; gap: 10px; }
.logo-icon { display: flex; align-items: center; justify-content: center; }
.logo-text { font-family: 'Poppins', sans-serif; font-size: 20px; font-weight: 700; color: #1a202c; }
.login-form-wrapper { max-width: 450px; margin: auto; padding: 60px 0 40px; width: 100%; }
.form-header { text-align: center; margin-bottom: 32px; }
.form-header h2 { font-family: 'Poppins', sans-serif; font-size: 28px; font-weight: 600; color: #1a202c; margin: 0 0 8px; }
.form-header p { font-size: 15px; font-weight: 500; color: #718096; margin: 0; line-height: 1.5; }
.form-header strong { color: #fd9644; }
.input-group { margin-bottom: 20px; }
.input-group label { display: block; font-size: 14px; font-weight: 600; color: #1a202c; margin-bottom: 8px; }
.input-wrapper { position: relative; }
.input-icon { position: absolute; left: 16px; top: 50%; transform: translateY(-50%); color: #a0aec0; pointer-events: none; }
.input-wrapper input { width: 100%; height: 56px; border: 1px solid #e2e8f0; border-radius: 12px; padding: 0 48px; font-size: 15px; font-weight: 500; color: #1a202c; background: #fff; transition: all .2s; font-family: 'Urbanist', sans-serif; box-sizing: border-box; }
.input-wrapper input::placeholder { color: #a0aec0; font-weight: 400; }
.input-wrapper input:focus { outline: none; border-color: #fd9644; box-shadow: 0 0 0 3px rgba(253,150,68,.15); }
.toggle-password { position: absolute; right: 16px; top: 50%; transform: translateY(-50%); background: none; border: none; cursor: pointer; padding: 0; display: flex; }
.btn-login { display: flex; align-items: center; justify-content: center; gap: 8px; width: 100%; height: 52px; border: none; border-radius: 12px; background: linear-gradient(135deg, #fd9644, #fc5c65); color: #fff; font-size: 16px; font-weight: 600; font-family: 'Urbanist', sans-serif; cursor: pointer; transition: all .25s; margin-top: 4px; }
.btn-login:hover:not(:disabled) { background: linear-gradient(135deg, #e8832e, #e04e55); transform: translateY(-1px); box-shadow: 0 4px 12px rgba(253,150,68,.35); }
.btn-login:disabled { opacity: .6; cursor: not-allowed; }
.spinner { width: 18px; height: 18px; border: 2.5px solid rgba(255,255,255,.3); border-top-color: #fff; border-radius: 50%; animation: spin .7s linear infinite; }
@keyframes spin { to { transform: rotate(360deg); } }
.register-link { text-align: center; margin-top: 24px; font-size: 14px; color: #718096; }
.register-link a { color: #fd9644; font-weight: 700; text-decoration: none; }
.register-link a:hover { text-decoration: underline; }

/* OTP Styles */
.otp-icon-circle { width: 64px; height: 64px; border-radius: 50%; background: #FFF6EE; display: flex; align-items: center; justify-content: center; margin: 0 auto 16px; }
.otp-inputs { display: flex; gap: 10px; justify-content: center; margin-bottom: 20px; }
.otp-input { width: 52px; height: 60px; border: 2px solid #e2e8f0; border-radius: 14px; text-align: center; font-size: 24px; font-weight: 700; color: #1a202c; font-family: 'Poppins', sans-serif; transition: all .2s; background: #fafafa; box-sizing: border-box; }
.otp-input:focus { outline: none; border-color: #fd9644; box-shadow: 0 0 0 3px rgba(253,150,68,.15); background: #fff; }
.otp-timer { text-align: center; font-size: 14px; color: #718096; margin-bottom: 20px; }
.otp-timer strong { color: #fd9644; }
.btn-resend { display: block; width: 100%; background: none; border: none; color: #fd9644; font-weight: 600; font-size: 14px; cursor: pointer; margin-bottom: 20px; padding: 8px; font-family: 'Urbanist', sans-serif; transition: color .15s; }
.btn-resend:hover:not(:disabled) { color: #e8832e; text-decoration: underline; }
.btn-resend:disabled { opacity: .5; cursor: not-allowed; }
.btn-back { display: flex; align-items: center; justify-content: center; gap: 6px; width: 100%; background: none; border: none; color: #718096; font-size: 14px; font-weight: 600; cursor: pointer; margin-top: 16px; padding: 8px; font-family: 'Urbanist', sans-serif; transition: color .15s; }
.btn-back:hover { color: #1a202c; }

.login-footer-links { display: flex; justify-content: center; flex-wrap: wrap; gap: 24px; padding-top: 48px; }
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
@media (max-width: 480px) { .form-header h2{font-size:22px} .login-form-wrapper{padding:40px 0 24px} .otp-input{width:44px;height:52px;font-size:20px} }
</style>
