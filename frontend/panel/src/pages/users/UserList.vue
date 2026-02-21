<template>
  <div class="animate-in">
    <div class="flex items-center justify-between mb-5">
      <div></div>
      <button @click="showModal=true" class="btn btn-primary" v-if="auth.hasRole('master')">+ Tambah User</button>
    </div>
    <div class="card">
      <table class="data-table">
        <thead><tr><th>Nama</th><th>Email</th><th>Role</th><th>Status</th></tr></thead>
        <tbody>
          <tr v-for="u in users" :key="u.id">
            <td class="font-semibold">{{ u.name }}</td>
            <td class="text-gray-500">{{ u.email }}</td>
            <td><span class="badge badge-primary">{{ auth.roleLabel(u.role) }}</span></td>
            <td><span class="badge" :class="u.is_active?'badge-success':'badge-danger'">{{ u.is_active?'Aktif':'Nonaktif' }}</span></td>
          </tr>
          <tr v-if="!users.length"><td colspan="4" class="text-center py-8 text-gray-400">Belum ada data</td></tr>
        </tbody>
      </table>
    </div>
    <div v-if="showModal" class="modal-overlay" @click.self="showModal=false">
      <div class="modal-content">
        <div class="modal-header"><h2>Tambah User</h2><button @click="showModal=false" class="btn btn-ghost btn-icon">✕</button></div>
        <form @submit.prevent="createUser" class="modal-body space-y-4">
          <div class="form-group"><label class="form-label">Nama *</label><input v-model="form.name" class="form-input" required></div>
          <div class="form-group"><label class="form-label">Email *</label><input v-model="form.email" type="email" class="form-input" required></div>
          <div class="form-group"><label class="form-label">Password *</label><input v-model="form.password" type="password" class="form-input" required></div>
          <div class="form-group"><label class="form-label">Role *</label>
            <select v-model="form.role" class="form-input" required>
              <option value="master">Master</option><option value="admin">Admin</option>
              <option value="finance">Finance</option>
              <option value="mitra">Mitra</option>
            </select>
          </div>
          <div class="flex gap-3 justify-end"><button type="button" @click="showModal=false" class="btn btn-secondary">Batal</button><button type="submit" class="btn btn-primary">Simpan</button></div>
        </form>
      </div>
    </div>
  </div>
</template>
<script setup>
import { ref, reactive, onMounted } from 'vue'
import { useAuthStore } from '../../stores/auth'
import { userApi } from '../../services/api'
const auth = useAuthStore()
const users = ref([]), showModal = ref(false)
const form = reactive({ name:'', email:'', password:'', role:'admin' })
onMounted(loadUsers)
async function loadUsers(){ try{ const{data}=await userApi.list({page:1,limit:50}); users.value=data.data||[] }catch(e){console.error(e)} }
async function createUser(){ try{ await userApi.create(form); showModal.value=false; Object.assign(form,{name:'',email:'',password:'',role:'admin'}); await loadUsers() }catch(e){console.error(e)} }
</script>
