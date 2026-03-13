<template>
  <div class="ss-wrap" ref="wrapRef">
    <button
      type="button"
      class="ss-trigger"
      :class="{ open: isOpen, filled: !!modelValue, disabled }"
      @click="toggle"
      :disabled="disabled"
    >
      <span class="ss-trigger-text" :class="{ placeholder: !selectedLabel }">{{ selectedLabel || placeholder }}</span>
      <svg class="ss-chevron" :class="{ open: isOpen }" width="16" height="16" viewBox="0 0 20 20" fill="none"><path d="M6 8l4 4 4-4" stroke="#94a3b8" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round"/></svg>
    </button>

    <Transition name="ss-drop">
      <div v-if="isOpen" class="ss-dropdown">
        <div class="ss-search-box">
          <Search :size="14" class="ss-search-icon" />
          <input
            ref="searchRef"
            v-model="query"
            type="text"
            class="ss-search-input"
            :placeholder="searchPlaceholder"
            @keydown.escape="close"
            @keydown.enter.prevent="selectFirst"
          />
        </div>
        <div class="ss-options" ref="optionsRef">
          <button
            v-if="allowEmpty"
            type="button"
            class="ss-option"
            :class="{ selected: !modelValue }"
            @click="select('')"
          >
            <span class="ss-option-text ss-option-empty">{{ emptyLabel }}</span>
            <Check v-if="!modelValue" :size="14" class="ss-check" />
          </button>
          <button
            v-for="opt in filteredOptions"
            :key="opt.value"
            type="button"
            class="ss-option"
            :class="{ selected: opt.value === modelValue }"
            @click="select(opt.value)"
          >
            <div class="ss-option-content">
              <span class="ss-option-text">{{ opt.label }}</span>
              <span v-if="opt.sub" class="ss-option-sub">{{ opt.sub }}</span>
            </div>
            <Check v-if="opt.value === modelValue" :size="14" class="ss-check" />
          </button>
          <div v-if="!filteredOptions.length" class="ss-no-result">
            <SearchX :size="18" />
            <span>Tidak ditemukan</span>
          </div>
        </div>
      </div>
    </Transition>
  </div>
</template>

<script setup>
import { ref, computed, watch, nextTick, onMounted, onUnmounted } from 'vue'
import { Search, Check, SearchX } from 'lucide-vue-next'

const props = defineProps({
  modelValue: { type: [String, Number], default: '' },
  options: { type: Array, default: () => [] },
  placeholder: { type: String, default: '— Pilih —' },
  searchPlaceholder: { type: String, default: 'Cari...' },
  emptyLabel: { type: String, default: '— Tidak dipilih —' },
  allowEmpty: { type: Boolean, default: true },
  disabled: { type: Boolean, default: false },
})

const emit = defineEmits(['update:modelValue'])

const isOpen = ref(false)
const query = ref('')
const wrapRef = ref(null)
const searchRef = ref(null)
const optionsRef = ref(null)

const selectedLabel = computed(() => {
  const opt = props.options.find(o => o.value === props.modelValue)
  return opt ? opt.label : ''
})

const filteredOptions = computed(() => {
  const q = query.value.toLowerCase().trim()
  if (!q) return props.options
  return props.options.filter(o =>
    o.label.toLowerCase().includes(q) ||
    (o.sub && o.sub.toLowerCase().includes(q))
  )
})

function toggle() {
  if (props.disabled) return
  isOpen.value ? close() : open()
}

function open() {
  isOpen.value = true
  query.value = ''
  nextTick(() => searchRef.value?.focus())
}

function close() {
  isOpen.value = false
  query.value = ''
}

function select(val) {
  emit('update:modelValue', val)
  close()
}

function selectFirst() {
  if (filteredOptions.value.length) {
    select(filteredOptions.value[0].value)
  }
}

function handleClickOutside(e) {
  if (wrapRef.value && !wrapRef.value.contains(e.target)) {
    close()
  }
}

onMounted(() => document.addEventListener('mousedown', handleClickOutside))
onUnmounted(() => document.removeEventListener('mousedown', handleClickOutside))
</script>

<style scoped>
.ss-wrap { position: relative; width: 100%; }

.ss-trigger {
  width: 100%; display: flex; align-items: center; justify-content: space-between;
  padding: 10px 14px; border: 1.5px solid #e2e8f0; border-radius: 11px;
  font-size: 0.85rem; color: #1e293b; background: #fafbfc;
  cursor: pointer; outline: none; transition: all 0.2s;
  font-family: inherit; text-align: left;
}
.ss-trigger:hover { border-color: #cbd5e1; }
.ss-trigger.open { border-color: #6366f1; box-shadow: 0 0 0 3px rgba(99,102,241,0.1); background: #fff; }
.ss-trigger.disabled { opacity: 0.5; cursor: not-allowed; background: #f1f5f9; }
.ss-trigger-text { flex: 1; overflow: hidden; text-overflow: ellipsis; white-space: nowrap; }
.ss-trigger-text.placeholder { color: #94a3b8; }

.ss-chevron { flex-shrink: 0; transition: transform 0.2s; }
.ss-chevron.open { transform: rotate(180deg); }

/* Dropdown */
.ss-dropdown {
  position: absolute; top: calc(100% + 6px); left: 0; right: 0; z-index: 50;
  background: #fff; border: 1.5px solid #e2e8f0; border-radius: 14px;
  box-shadow: 0 12px 40px rgba(0,0,0,0.1), 0 2px 8px rgba(0,0,0,0.05);
  overflow: hidden;
}

/* Search */
.ss-search-box {
  display: flex; align-items: center; gap: 8px;
  padding: 10px 14px; border-bottom: 1px solid #f1f5f9;
}
.ss-search-icon { color: #94a3b8; flex-shrink: 0; }
.ss-search-input {
  flex: 1; border: none; outline: none; font-size: 0.82rem;
  color: #1e293b; background: transparent; font-family: inherit;
}
.ss-search-input::placeholder { color: #cbd5e1; }

/* Options */
.ss-options { max-height: 220px; overflow-y: auto; padding: 4px; }
.ss-options::-webkit-scrollbar { width: 5px; }
.ss-options::-webkit-scrollbar-track { background: transparent; }
.ss-options::-webkit-scrollbar-thumb { background: #e2e8f0; border-radius: 3px; }

.ss-option {
  width: 100%; display: flex; align-items: center; justify-content: space-between; gap: 8px;
  padding: 9px 12px; border: none; background: transparent;
  cursor: pointer; border-radius: 9px; transition: all 0.12s;
  font-family: inherit; text-align: left;
}
.ss-option:hover { background: #f8fafc; }
.ss-option.selected { background: #eef2ff; }

.ss-option-content { flex: 1; min-width: 0; }
.ss-option-text { display: block; font-size: 0.84rem; color: #1e293b; font-weight: 500; overflow: hidden; text-overflow: ellipsis; white-space: nowrap; }
.ss-option-sub { display: block; font-size: 0.72rem; color: #94a3b8; margin-top: 1px; overflow: hidden; text-overflow: ellipsis; white-space: nowrap; }
.ss-option-empty { color: #94a3b8; font-style: italic; font-size: 0.82rem; }
.ss-check { color: #6366f1; flex-shrink: 0; }

.ss-no-result {
  display: flex; align-items: center; justify-content: center; gap: 8px;
  padding: 24px 16px; color: #cbd5e1; font-size: 0.82rem;
}

/* Transition */
.ss-drop-enter-active { transition: all 0.15s ease; }
.ss-drop-leave-active { transition: all 0.1s ease; }
.ss-drop-enter-from { opacity: 0; transform: translateY(-6px); }
.ss-drop-leave-to { opacity: 0; transform: translateY(-4px); }
</style>
