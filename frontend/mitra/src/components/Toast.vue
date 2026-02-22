<template>
  <Teleport to="body">
    <div class="toast-container" v-if="toast.toasts.length">
      <TransitionGroup name="toast">
        <div
          v-for="t in toast.toasts"
          :key="t.id"
          class="toast-item"
          :class="`toast-${t.type}`"
          @click="toast.remove(t.id)"
        >
          <i v-if="t.type === 'success'" class="toast-icon ri-checkbox-circle-line"></i>
          <i v-else-if="t.type === 'error'" class="toast-icon ri-close-circle-line"></i>
          <i v-else-if="t.type === 'warning'" class="toast-icon ri-alert-line"></i>
          <i v-else class="toast-icon ri-information-line"></i>
          <span class="toast-message">{{ t.message }}</span>
          <button class="toast-close" @click.stop="toast.remove(t.id)"><i class="ri-close-line"></i></button>
        </div>
      </TransitionGroup>
    </div>
  </Teleport>
</template>

<script setup>
import { useToastStore } from '../stores/toast'
const toast = useToastStore()
</script>

<style scoped>
.toast-container {
  position: fixed;
  top: 1.25rem;
  right: 1.25rem;
  z-index: 99999;
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
  max-width: 420px;
}

.toast-item {
  display: flex;
  align-items: center;
  gap: 0.625rem;
  padding: 0.875rem 1rem;
  border-radius: 0.75rem;
  font-size: 0.875rem;
  font-weight: 500;
  color: white;
  cursor: pointer;
  box-shadow: 0 10px 25px rgba(0,0,0,0.15);
  backdrop-filter: blur(8px);
  min-width: 280px;
}

.toast-icon { flex-shrink: 0; }
.toast-message { flex: 1; line-height: 1.4; }
.toast-close {
  background: none;
  border: none;
  color: rgba(255,255,255,0.7);
  cursor: pointer;
  font-size: 0.75rem;
  padding: 0.25rem;
  flex-shrink: 0;
  transition: color 0.2s;
}
.toast-close:hover { color: white; }

.toast-success { background: linear-gradient(135deg, #10b981, #059669); }
.toast-error { background: linear-gradient(135deg, #ef4444, #dc2626); }
.toast-warning { background: linear-gradient(135deg, #f59e0b, #d97706); }
.toast-info { background: linear-gradient(135deg, #4f6df5, #4338ca); }

/* Transitions */
.toast-enter-active { animation: toastIn 0.35s ease-out; }
.toast-leave-active { animation: toastOut 0.25s ease-in forwards; }

@keyframes toastIn {
  from { opacity: 0; transform: translateX(100px) scale(0.95); }
  to { opacity: 1; transform: translateX(0) scale(1); }
}
@keyframes toastOut {
  from { opacity: 1; transform: translateX(0) scale(1); }
  to { opacity: 0; transform: translateX(100px) scale(0.95); }
}
</style>
