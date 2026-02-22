<template>
  <header class="app-header" :class="{ 'sidebar-expanded': sidebarExpanded }">
    <div class="header-inner">
      <!-- Drawer toggle (visible when sidebar collapsed) -->
      <button v-if="!sidebarExpanded" type="button" class="header-drawer-btn" @click="$emit('toggleSidebar')">
        <svg width="16" height="40" viewBox="0 0 16 40" fill="none" xmlns="http://www.w3.org/2000/svg">
          <path d="M0 10C0 4.47715 4.47715 0 10 0H16V40H10C4.47715 40 0 35.5228 0 30V10Z" fill="var(--primary)"/>
          <path d="M10 15L6 20.0049L10 25.0098" stroke="#fff" stroke-width="1.2" stroke-linecap="round" stroke-linejoin="round" transform="rotate(180 8 20)"/>
        </svg>
      </button>

      <!-- Page title -->
      <div class="header-title-area">
        <h3 class="header-title">{{ title }}</h3>
        <p v-if="subtitle" class="header-subtitle">{{ subtitle }}</p>
      </div>

      <!-- Search bar (desktop only) -->
      <div class="header-search-wrapper">
        <div class="header-search">
          <span class="header-search-icon">
            <svg width="20" height="20" viewBox="0 0 20 20" fill="none" xmlns="http://www.w3.org/2000/svg">
              <circle cx="9.78639" cy="9.78602" r="8.23951" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round"/>
              <path d="M15.5176 15.9447L18.7479 19.1667" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round"/>
            </svg>
          </span>
          <input type="text" placeholder="Search..." class="header-search-input" />
        </div>
      </div>

      <!-- Right side actions -->
      <div class="header-actions">
        <!-- Date -->
        <span class="header-date">{{ today }}</span>

        <!-- Notification button -->
        <button type="button" class="header-action-btn">
          <svg width="24" height="24" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
            <path d="M12 22C13.38 22 14.56 21.17 15 20H9C9.44 21.17 10.62 22 12 22Z" fill="var(--primary)"/>
            <path fill-rule="evenodd" clip-rule="evenodd" d="M13.77 2.84C13.3 2.32 12.63 2 11.89 2C10.5 2 9.37 3.13 9.37 4.52V4.64C6.99 5.45 5.22 7.39 4.96 9.75L4.46 14.13C4.37 14.92 4.04 15.67 3.5 16.3C2.28 17.73 3.44 19.78 5.47 19.78H18.31C20.34 19.78 21.5 17.73 20.28 16.3C19.74 15.67 19.41 14.92 19.32 14.13L18.82 9.75C18.81 9.69 18.81 9.63 18.8 9.57C18.37 9.7 17.92 9.78 17.44 9.78C14.99 9.78 13 7.79 13 5.33C13 4.41 13.28 3.55 13.77 2.84Z" fill="currentColor"/>
            <circle cx="17.44" cy="5.33" r="3.33" fill="var(--primary)"/>
          </svg>
        </button>

        <!-- User -->
        <div class="header-user" @click="showProfileMenu = !showProfileMenu">
          <div class="header-user-avatar">{{ userInitial }}</div>
          <div class="header-user-info">
            <h3 class="header-user-name">{{ userName }}</h3>
            <p class="header-user-role">{{ userRole }}</p>
          </div>
          <svg width="16" height="16" viewBox="0 0 24 24" fill="none"><path d="M7 10L12 14L17 10" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/></svg>
        </div>

        <!-- Profile dropdown -->
        <div v-if="showProfileMenu" class="header-profile-dropdown">
          <button @click="showProfileMenu = false; $emit('logout')" class="header-profile-item text-orange-500">
            <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linecap="round"><path d="M15 10L13.7071 11.2929C13.3166 11.6834 13.3166 12.3166 13.7071 12.7071L15 14M14 12L22 12M6 20C3.79086 20 2 18.2091 2 16V8C2 5.79086 3.79086 4 6 4M6 20C8.20914 20 10 18.2091 10 16V8C10 5.79086 8.20914 4 6 4M6 20H14C16.2091 20 18 18.2091 18 16M6 4H14C16.2091 4 18 5.79086 18 8"/></svg>
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
