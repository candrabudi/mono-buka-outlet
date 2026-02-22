<template>
  <div class="nf-page">
    <div class="nf-card">
      <h1>404</h1>
      <p>Halaman tidak ditemukan</p>
      <router-link to="/" class="nf-btn">Kembali ke Dashboard</router-link>
    </div>
  </div>
</template>

<style scoped>
.nf-page { display: flex; align-items: center; justify-content: center; min-height: 100vh; background: #f8fafc; }
.nf-card { text-align: center; padding: 60px 40px; }
.nf-card h1 { font-size: 6rem; font-weight: 900; color: #e2e8f0; margin: 0 0 10px; line-height: 1; }
.nf-card p { font-size: 1.1rem; color: #94a3b8; margin: 0 0 32px; }
.nf-btn { display: inline-block; padding: 12px 28px; background: linear-gradient(135deg, #6366f1, #8b5cf6); color: #fff; border-radius: 12px; text-decoration: none; font-weight: 600; font-size: .9rem; }
</style>
