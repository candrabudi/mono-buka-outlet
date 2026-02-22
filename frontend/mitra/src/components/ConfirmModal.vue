<template>
  <Teleport to="body">
    <Transition name="cm-fade">
      <div v-if="modelValue" class="cm-overlay" @click.self="cancel">
        <Transition name="cm-scale">
          <div v-if="modelValue" class="cm-dialog" :class="variant">
            <!-- Icon -->
            <div class="cm-icon-wrap" :class="variant">
              <slot name="icon">
                <i v-if="variant === 'danger'" class="ri-close-circle-line" style="font-size:28px"></i>
                <i v-else-if="variant === 'warning'" class="ri-alert-line" style="font-size:28px"></i>
                <i v-else class="ri-information-line" style="font-size:28px"></i>
              </slot>
            </div>

            <!-- Content -->
            <h3 class="cm-title">{{ title }}</h3>
            <p class="cm-message">{{ message }}</p>
            <slot name="body"></slot>

            <!-- Actions -->
            <div class="cm-actions">
              <button class="cm-btn cm-btn-cancel" @click="cancel" :disabled="loading">
                {{ cancelText }}
              </button>
              <button class="cm-btn cm-btn-confirm" :class="variant" @click="confirm" :disabled="loading">
                <span v-if="loading" class="cm-spinner"></span>
                {{ loading ? loadingText : confirmText }}
              </button>
            </div>
          </div>
        </Transition>
      </div>
    </Transition>
  </Teleport>
</template>

<script setup>
defineProps({
  modelValue: { type: Boolean, default: false },
  title: { type: String, default: 'Konfirmasi' },
  message: { type: String, default: 'Apakah Anda yakin?' },
  confirmText: { type: String, default: 'Ya, Lanjutkan' },
  cancelText: { type: String, default: 'Batal' },
  loadingText: { type: String, default: 'Memproses...' },
  variant: { type: String, default: 'danger' }, // danger | warning | info
  loading: { type: Boolean, default: false },
})

const emit = defineEmits(['update:modelValue', 'confirm', 'cancel'])

function cancel() {
  emit('update:modelValue', false)
  emit('cancel')
}
function confirm() {
  emit('confirm')
}
</script>

<style scoped>
.cm-overlay{position:fixed;inset:0;z-index:9999;background:rgba(15,23,42,.5);backdrop-filter:blur(4px);display:flex;align-items:center;justify-content:center;padding:20px}
.cm-dialog{background:#fff;border-radius:20px;padding:32px 28px 24px;max-width:420px;width:100%;text-align:center;box-shadow:0 25px 50px rgba(0,0,0,.15)}

.cm-icon-wrap{width:64px;height:64px;border-radius:50%;display:flex;align-items:center;justify-content:center;margin:0 auto 18px}
.cm-icon-wrap.danger{background:#fef2f2;color:#ef4444}
.cm-icon-wrap.warning{background:#fffbeb;color:#f59e0b}
.cm-icon-wrap.info{background:#eff6ff;color:#3b82f6}

.cm-title{font-size:1.1rem;font-weight:800;color:#0f172a;margin:0 0 8px}
.cm-message{font-size:.88rem;color:#64748b;line-height:1.6;margin:0 0 24px}

.cm-actions{display:flex;gap:10px}
.cm-btn{flex:1;padding:12px 0;border:none;border-radius:12px;font-size:.85rem;font-weight:700;cursor:pointer;transition:all .2s;font-family:inherit;display:inline-flex;align-items:center;justify-content:center;gap:6px}
.cm-btn:disabled{opacity:.5;cursor:not-allowed}
.cm-btn-cancel{background:#f1f5f9;color:#475569}
.cm-btn-cancel:hover:not(:disabled){background:#e2e8f0}
.cm-btn-confirm.danger{background:linear-gradient(135deg,#ef4444,#dc2626);color:#fff;box-shadow:0 2px 8px rgba(239,68,68,.25)}
.cm-btn-confirm.danger:hover:not(:disabled){box-shadow:0 4px 14px rgba(239,68,68,.35);transform:translateY(-1px)}
.cm-btn-confirm.warning{background:linear-gradient(135deg,#f59e0b,#d97706);color:#fff;box-shadow:0 2px 8px rgba(245,158,11,.25)}
.cm-btn-confirm.warning:hover:not(:disabled){box-shadow:0 4px 14px rgba(245,158,11,.35);transform:translateY(-1px)}
.cm-btn-confirm.info{background:linear-gradient(135deg,#3b82f6,#2563eb);color:#fff;box-shadow:0 2px 8px rgba(59,130,246,.25)}
.cm-btn-confirm.info:hover:not(:disabled){box-shadow:0 4px 14px rgba(59,130,246,.35);transform:translateY(-1px)}

.cm-spinner{width:14px;height:14px;border:2px solid rgba(255,255,255,.3);border-top-color:#fff;border-radius:50%;animation:cm-spin .7s linear infinite}
@keyframes cm-spin{to{transform:rotate(360deg)}}

/* Transitions */
.cm-fade-enter-active,.cm-fade-leave-active{transition:opacity .2s ease}
.cm-fade-enter-from,.cm-fade-leave-to{opacity:0}
.cm-scale-enter-active{transition:all .25s cubic-bezier(.34,1.56,.64,1)}
.cm-scale-leave-active{transition:all .15s ease}
.cm-scale-enter-from{opacity:0;transform:scale(.9) translateY(10px)}
.cm-scale-leave-to{opacity:0;transform:scale(.95)}
</style>
