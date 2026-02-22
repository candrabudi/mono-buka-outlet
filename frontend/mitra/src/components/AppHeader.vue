<template>
  <header class="app-header" :class="{ 'sidebar-expanded': sidebarExpanded }">
    <div class="header-inner">
      <!-- Drawer toggle (visible when sidebar collapsed) -->
      <button v-if="!sidebarExpanded" type="button" class="header-drawer-btn" @click="$emit('toggleSidebar')">
        <i class="ri-menu-line"></i>
      </button>

      <!-- Page title -->
      <div class="header-title-area">
        <h3 class="header-title">{{ title }}</h3>
        <p v-if="subtitle" class="header-subtitle">{{ subtitle }}</p>
      </div>

      <!-- Search bar (desktop only) -->
      <div class="header-search-wrapper">
        <div class="header-search">
          <span class="header-search-icon"><i class="ri-search-line"></i></span>
          <input type="text" placeholder="Search..." class="header-search-input" />
        </div>
      </div>

      <!-- Right side actions -->
      <div class="header-actions">
        <!-- Date -->
        <span class="header-date">{{ today }}</span>

        <!-- Notification button -->
        <button type="button" class="header-action-btn">
          <i class="ri-notification-3-fill" style="font-size:20px"></i>
        </button>

        <!-- User -->
        <div class="header-user" @click="showProfileMenu = !showProfileMenu">
          <div class="header-user-avatar">{{ userInitial }}</div>
          <div class="header-user-info">
            <h3 class="header-user-name">{{ userName }}</h3>
            <p class="header-user-role">{{ userRole }}</p>
          </div>
          <i class="ri-arrow-down-s-line"></i>
        </div>

        <!-- Profile dropdown -->
        <div v-if="showProfileMenu" class="header-profile-dropdown">
          <button @click="showProfileMenu = false; $emit('logout')" class="header-profile-item text-orange-500">
            <i class="ri-logout-box-r-line"></i>
            <span>Logout</span>
          </button>
        </div>
        <div v-if="showProfileMenu" class="header-profile-overlay" @click="showProfileMenu = false"></div>
      </div>
    </div>
  </header>
</template>

<script setup>
import { ref } from 'vue'

defineProps({
  title: { type: String, default: '' },
  subtitle: { type: String, default: '' },
  sidebarExpanded: { type: Boolean, default: true },
  userName: { type: String, default: '' },
  userInitial: { type: String, default: '' },
  userRole: { type: String, default: '' },
})

defineEmits(['toggleSidebar', 'logout'])

const today = new Date().toLocaleDateString('id-ID', {
  weekday: 'short', year: 'numeric', month: 'short', day: 'numeric'
})

const showProfileMenu = ref(false)
</script>
