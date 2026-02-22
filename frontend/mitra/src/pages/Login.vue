<template>
  <section class="login-section">
    <div class="login-wrapper">
      <!-- Left: Form -->
      <div class="login-left">
        <header class="login-logo">
          <div class="logo-icon">
            <svg width="32" height="32" viewBox="0 0 32 32" fill="none">
              <rect width="32" height="32" rx="8" fill="#fd9644"/>
              <path d="M9 10H23V12H17V22H15V12H9V10ZM19 14H23V16H19V14Z" fill="white"/>
            </svg>
          </div>
          <span class="logo-text">BukaOutlet</span>
        </header>

        <div class="login-form-wrapper">
          <div class="login-form-content animate-fade-in">
            <header class="form-header">
              <h2>Masuk ke Portal Mitra</h2>
              <p>Pantau kemitraan dan pembayaran Anda</p>
            </header>

            <form @submit.prevent="handleLogin">
              <div class="input-group">
                <label for="email">Email</label>
                <div class="input-wrapper">
                  <svg class="input-icon" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
                    <path d="M3 8l7.89 5.26a2 2 0 002.22 0L21 8M5 19h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v10a2 2 0 002 2z"/>
                  </svg>
                  <input id="email" v-model="form.email" type="email" placeholder="mitra@email.com" required autocomplete="email" />
                </div>
              </div>

              <div class="input-group">
                <label for="password">Password</label>
                <div class="input-wrapper">
                  <svg class="input-icon" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
                    <rect x="3" y="11" width="18" height="11" rx="2" ry="2"/>
                    <path d="M7 11V7a5 5 0 0110 0v4"/>
                  </svg>
                  <input id="password" v-model="form.password" :type="showPw ? 'text' : 'password'" placeholder="••••••••" required autocomplete="current-password" />
                  <button type="button" class="toggle-password" @click="showPw = !showPw" tabindex="-1">
                    <svg v-if="!showPw" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="#718096" stroke-width="1.5"><path d="M17.94 17.94A10.07 10.07 0 0 1 12 20c-7 0-11-8-11-8a18.45 18.45 0 0 1 5.06-5.94M9.9 4.24A9.12 9.12 0 0 1 12 4c7 0 11 8 11 8a18.5 18.5 0 0 1-2.16 3.19m-6.72-1.07a3 3 0 1 1-4.24-4.24"/><line x1="1" y1="1" x2="23" y2="23"/></svg>
                    <svg v-else width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="#718096" stroke-width="1.5"><path d="M1 12s4-8 11-8 11 8 11 8-4 8-11 8-11-8-11-8z"/><circle cx="12" cy="12" r="3"/></svg>
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
          <h3>Portal Mitra BukaOutlet</h3>
          <p>Pantau perkembangan kemitraan, pembayaran, dan status outlet Anda. Semua informasi dalam satu platform yang <span class="highlight">terintegrasi</span>.</p>
        </div>
      </div>
    </div>
  </section>
</template>

<script setup>
import { ref, reactive } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useAuthStore } from '../stores/auth'
import { useToastStore } from '../stores/toast'
import { authApi } from '../services/api'

const auth = useAuthStore()
const toast = useToastStore()
const router = useRouter()
const route = useRoute()
const loading = ref(false)
const showPw = ref(false)
const form = reactive({ email: '', password: '' })

async function handleLogin() {
  loading.value = true
  try {
    const { data } = await authApi.login(form)
    const d = data.data
    // Direct login — no OTP for mitra
    auth.setAuth(d.token, d.user)
    toast.success('Login berhasil!')
    router.push(route.query.redirect || '/')
  } catch (e) {
    const msg = e.response?.data?.error || (e.code === 'ERR_NETWORK' ? 'Server tidak dapat dihubungi' : 'Login gagal')
    toast.error(msg)
  } finally {
    loading.value = false
  }
}
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
.form-header p { font-size: 15px; font-weight: 500; color: #718096; margin: 0; }
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
@media (max-width: 480px) { .form-header h2{font-size:22px} .login-form-wrapper{padding:40px 0 24px} }
</style>
